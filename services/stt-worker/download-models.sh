#!/usr/bin/env bash
set -euo pipefail

# This script downloads a Vosk Bangla model into services/stt-worker/models/vosk-bn
# Usage:
#   ./download-models.sh <MODEL_URL>
# If MODEL_URL is not provided, the script will prompt for one.

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
MODELS_DIR="$SCRIPT_DIR/models"
TARGET_DIR="$MODELS_DIR/vosk-bn"
TMP_DIR="$MODELS_DIR/.tmp"

MODEL_URL="${1:-${VOSK_MODEL_URL:-}}"

if [[ -z "${MODEL_URL}" ]]; then
  echo "ERROR: No model URL provided."
  echo "Provide a Vosk Bangla model URL, for example a zip/tar.gz from alphacephei.com or a trusted mirror."
  echo "Usage: ./download-models.sh https://example.com/vosk-model-bn-small.zip"
  exit 1
fi

mkdir -p "$TARGET_DIR" "$TMP_DIR"
cd "$TMP_DIR"

FILENAME="$(basename "$MODEL_URL")"

echo "Downloading model from: $MODEL_URL"
curl -L "$MODEL_URL" -o "$FILENAME"

# Try unzip first, then tar if needed
if unzip -t "$FILENAME" >/dev/null 2>&1; then
  echo "Detected zip archive, extracting..."
  unzip -q "$FILENAME"
else
  echo "Attempting tar extraction..."
  tar -xzf "$FILENAME" || tar -xJf "$FILENAME" || tar -xvf "$FILENAME"
fi

# Move the extracted directory (first directory) into TARGET_DIR
FIRST_DIR="$(find . -maxdepth 1 -type d ! -name '.' | head -n1)"
if [[ -z "$FIRST_DIR" ]]; then
  echo "ERROR: Could not find extracted folder."
  exit 1
fi

# Clean target and move
rm -rf "$TARGET_DIR"
mkdir -p "$TARGET_DIR"
rsync -a "$FIRST_DIR"/ "$TARGET_DIR"/

cd "$SCRIPT_DIR"
rm -rf "$TMP_DIR"

echo "✅ Model downloaded to: $TARGET_DIR"
echo "Set VOSK_MODEL_DIR=$TARGET_DIR when running the worker if needed."
