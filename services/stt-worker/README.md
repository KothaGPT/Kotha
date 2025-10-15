# Kotha STT Worker (Offline)

Python FastAPI service that provides a local WebSocket streaming interface for speech-to-text using Vosk (Bangla-focused). This enables the "Self-hosted" login option and offline transcription for the MVP.

## Endpoints

- GET `/health` → `{ status: "ok" }` when the model is loaded
- WS `/ws` → binary PCM 16kHz/16-bit chunks in, JSON partial/final results out

Example WS messages:

- Partial: `{ "final": false, "result": { "partial": "..." } }`
- Final: `{ "final": true,  "result": { "text": "...", ... } }`

## Quick start

1. Install Python 3.11+ and create a venv (recommended)

```bash
cd services/stt-worker
python -m venv .venv
source .venv/bin/activate  # Windows: .venv\\Scripts\\activate
pip install --upgrade pip
pip install -r requirements.txt
```

2. Download a Vosk Bangla model

```bash
# Provide a proper model URL (zip/tar.gz). Check alphacephei.com or trusted mirrors.
./download-models.sh https://example.com/vosk-model-bn-small.zip
```

3. Run the worker

```bash
export VOSK_MODEL_DIR="$(pwd)/models/vosk-bn"   # optional if using default
export STT_WORKER_PORT=8765                      # default 8765
python main.py
```

Now `/health` should return OK and `/ws` should accept streaming audio.

## Integration notes

- The Electron app checks server health via the main process (`window.api.checkServerHealth()`), which should target `http://127.0.0.1:8765/health`.
- Send 16kHz, 16-bit PCM chunks over the websocket for best results (`SAMPLE_RATE=16000`).
- You can point to a different model directory with `VOSK_MODEL_DIR`.

## Next steps (beyond MVP)

- Add a punctuation-restorer microservice and connect it post-ASR.
- Support alternative engines (Whisper-ONNX, Coqui) behind a unified interface.
- Add word boosting / custom vocabulary handling.
