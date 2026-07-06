# Voice Platform

This folder documents the voice deployment path for OneMind OS.

- Source application: [onemindos/omos-voice](https://github.com/onemindos/omos-voice)
- Current role: keep Mumble PTT support, then evolve the voice stack toward Pipecat + LiveKit
- Target deployment: DigitalOcean Kubernetes (DOKS)

## Pipeline direction

- Mumble remains the proven field PTT path
- Pipecat becomes the pipeline orchestrator for STT, brain routing, TTS, and interruption handling
- LiveKit provides the WebRTC transport for browser, iOS, and ATAK clients
- NATS remains the shared bus between voice and Legacy services

See `livekit-values.yaml` for the LiveKit Helm values template used on DOKS.