# syntax=docker/dockerfile:1
# takml-vision — OneMind TAK-ML Vision Worker
#
# Bakes worker.py + torch/opencv deps into a versioned image so the pod
# boots in seconds (no pip install at boot like the old ConfigMap variant).
# Model stays on the PVC (/models) for hot-swap without image rebuilds.
#
# Build:  docker build -t ghcr.io/onemindos/takml-vision:1 .
# Push:   docker push ghcr.io/onemindos/takml-vision:1
FROM python:3.11-slim

# ffmpeg/libgl needed by opencv video + headless rendering
RUN apt-get update && apt-get install -y --no-install-recommends \
      ffmpeg libgl1 libglib2.0-0 \
    && rm -rf /var/lib/apt/lists/*

# pinned deps — torch CPU build is smaller + enough for YOLOv5s inference
RUN pip install --no-cache-dir \
      "torch==2.2.2" --index-url https://download.pytorch.org/whl/cpu \
      opencv-python-headless==4.10.0.84 \
      nats-py==2.7.0 \
      numpy==1.26.4

WORKDIR /app
COPY worker.py /app/worker.py

ENV PYTHONUNBUFFERED=1
ENTRYPOINT ["python3", "/app/worker.py"]
