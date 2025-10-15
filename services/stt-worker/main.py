import os
import json
import asyncio
from fastapi import FastAPI, WebSocket, WebSocketDisconnect
from fastapi.responses import JSONResponse
from fastapi.middleware.cors import CORSMiddleware
from vosk import Model, KaldiRecognizer

APP_PORT = int(os.getenv("STT_WORKER_PORT", "8765"))
MODEL_DIR = os.getenv("VOSK_MODEL_DIR", os.path.join(os.path.dirname(__file__), "models", "vosk-bn"))
SAMPLE_RATE = int(os.getenv("STT_SAMPLE_RATE", "16000"))

app = FastAPI(title="Kotha STT Worker", version="0.1.0")

# Optional CORS (useful when testing from browser)
app.add_middleware(
    CORSMiddleware,
    allow_origins=["*"],
    allow_credentials=True,
    allow_methods=["*"],
    allow_headers=["*"],
)

_model = None


def get_model() -> Model:
    global _model
    if _model is None:
        if not os.path.isdir(MODEL_DIR):
            raise RuntimeError(
                f"Vosk model not found at '{MODEL_DIR}'. Run services/stt-worker/download-models.sh or set VOSK_MODEL_DIR."
            )
        _model = Model(MODEL_DIR)
    return _model


@app.get("/health")
async def health():
    try:
        get_model()
        return {"status": "ok"}
    except Exception as e:
        return JSONResponse(status_code=500, content={"status": "error", "error": str(e)})


@app.websocket("/ws")
async def websocket_endpoint(ws: WebSocket):
    await ws.accept()
    # Initialize recognizer per-connection
    try:
        rec = KaldiRecognizer(get_model(), SAMPLE_RATE)
        rec.SetWords(True)
    except Exception as e:
        await ws.send_json({"error": f"failed_to_init_recognizer: {str(e)}"})
        await ws.close()
        return

    try:
        while True:
            try:
                data = await ws.receive_bytes()
            except WebSocketDisconnect:
                break

            if rec.AcceptWaveform(data):
                result = rec.Result()
                # Vosk returns JSON string; forward with a final flag
                try:
                    parsed = json.loads(result)
                except Exception:
                    parsed = {"text": result}
                await ws.send_json({"final": True, "result": parsed})
            else:
                partial = rec.PartialResult()
                try:
                    parsed = json.loads(partial)
                except Exception:
                    parsed = {"partial": partial}
                await ws.send_json({"final": False, "result": parsed})
    finally:
        await ws.close()


if __name__ == "__main__":
    import uvicorn

    uvicorn.run(
        "main:app",
        host="0.0.0.0",
        port=APP_PORT,
        reload=False,
        log_level="info",
    )
