# OMOS-VOICE — Archived Service

**Status:** ARCHIVED (2026-08-12)
**Reason:** Replaced by ATAK-native voice architecture. omos-voice repo moved to Mac Trash.

## What It Was
Low-latency voice to action pipeline for OneMind.
Mumble PTT → RealtimeSTT (Whisper) → NATS → Legacy → Piper TTS → Mumble.

## Why Archived
- Mumble is a separate app — not TAK-native
- LiveKit/Pipecat upgrade was started but never deployed to cluster
- Brain milestone (Legacy ↔ voice.transcript) was never completed
- ATAK has native PTT plugin framework — cleaner path

## What Was Preserved
- NATS subject contract → saved to SO Wiki: Voice — NATS Subject Contract & ATAK Architecture
- Piper TTS model preserved in Mac Trash (recoverable if needed)

## Target Architecture (replaces this)
ATAK PTT native plugin
  → audio capture → faster-whisper/Parakeet (local, sovereign)
  → NATS voice.transcript
  → Legacy reasons → cmd.* / voice.say
  → NATS cot.in → TAK map + speaker

No Mumble. No LiveKit. No bridge. TAK-native.
