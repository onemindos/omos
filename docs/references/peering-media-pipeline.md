# Peering — Distributed Media Pipeline on NATS JetStream

**Repo:** (third-party)  
**Language:** Go + React  
**Stack:** NATS JetStream + Redis + MinIO + React  
**Role in OneMind:** Reference architecture for video/media processing pipeline over NATS

Production-grade distributed media processing system. Supports subtitle overlay, video transcoding, and audio extraction with real-time SSE status streaming.

## Architecture

- NATS JetStream — job queue and event bus
- MinIO — object storage for media files
- Redis — job state cache
- React UI with SSE for real-time status

## Why it matters for OneMind

Reference for the om-streaming-bridge / video surveillance pipeline. When OneMind processes drone video or field camera feeds:

```
Video ingested → evt.media.{node}.frame.received (NATS)
    → Peering-style transcoder worker (JetStream consumer)
    → Processed clip to MinIO
    → evt.media.{node}.clip.ready
    → Vision agent (YOLO/VLM analysis)
```

## See Also

- `simple_edge_videostream` (iot-edge) — NAT-traversal webcam stream for ingestion
- `om-streaming-bridge` (omos-services) — OneMind's own streaming bridge service
