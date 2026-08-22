"""
OneMind OS — Presence Mode Voice Agent
=======================================
WebSocket audio bridge: Flutter app ↔ STT ↔ LLM (Hermes) ↔ TTS

Architecture:
  Flutter sends raw PCM 16-bit mono 16kHz over WebSocket (port 8765)
  → VAD buffers into utterances
  → STT: whisper CLI → Hermes /api/audio/transcribe → fallback stub
  → LLM: Hermes API server port 8642 (OpenAI-compat, SSE streaming)
  → TTS: macOS `say -v Samantha` → PCM chunks back over WebSocket
  All events → NATS subjects throughout

Usage:
    python3 om_presence.py [--host 0.0.0.0] [--port 8765]
"""

from __future__ import annotations

import asyncio
import base64
import json
import logging
import math
import os
import struct
import subprocess
import sys
import tempfile
import time
import uuid
from pathlib import Path
from typing import Optional

import httpx
import nats
import websockets
from nats.js.errors import NotFoundError, BucketNotFoundError
from nats.js.api import KeyValueConfig
from websockets.legacy.server import WebSocketServerProtocol

# ---------------------------------------------------------------------------
# Configuration
# ---------------------------------------------------------------------------

NATS_URL = os.environ.get("NATS_URL", "nats://127.0.0.1:4222")
HERMES_DASHBOARD_URL = os.environ.get("HERMES_DASHBOARD_URL", "http://127.0.0.1:9119")
HERMES_API_URL = os.environ.get("HERMES_API_URL", "http://127.0.0.1:8642")
HERMES_DASHBOARD_TOKEN = os.environ.get(
    "HERMES_DASHBOARD_TOKEN", "h01xQ91f6qTW7V0SGlLZRYsBgZLR1SF1oV9uCpxAElo"
)
HERMES_API_KEY = os.environ.get("HERMES_API_KEY", "JdFfBxXEPLYy26flc6rdckkrvLp3DaEkKIQhebDiWTc")
HERMES_MODEL = os.environ.get("HERMES_MODEL", "claude-sonnet-4.6")

WS_HOST = os.environ.get("OM_PRESENCE_HOST", "127.0.0.1")
WS_PORT = int(os.environ.get("OM_PRESENCE_PORT", "8765"))

# Audio / VAD
SAMPLE_RATE = 16000          # Hz
SAMPLE_WIDTH = 2             # bytes (int16)
CHUNK_SAMPLES = 1024         # samples per WebSocket frame
VAD_SILENCE_MS = 800         # ms of silence → end of utterance
VAD_MIN_UTTERANCE_MS = 200   # discard shorter utterances (noise)
VAD_DEFAULT_THRESHOLD = 300  # RMS energy threshold

# TTS chunk size when streaming PCM back to Flutter
TTS_CHUNK_BYTES = 4096

# NATS buckets / streams
KV_PRESENCE_BUCKET = "presence"
KV_CONFIG_BUCKET = "om-config"
EVT_STREAM_SUBJECT = "evt.so.voice.utterance"

LOG_DIR = Path(__file__).parent / "logs"
LOG_FILE = LOG_DIR / "om-presence.log"

# ---------------------------------------------------------------------------
# Logging
# ---------------------------------------------------------------------------

LOG_DIR.mkdir(exist_ok=True)

logger = logging.getLogger("om-presence")
logger.setLevel(logging.DEBUG)
_fmt = logging.Formatter("%(asctime)s %(levelname)-8s %(message)s")
_fh = logging.FileHandler(LOG_FILE)
_fh.setFormatter(_fmt)
_fh.setLevel(logging.DEBUG)
_ch = logging.StreamHandler(sys.stdout)
_ch.setFormatter(_fmt)
_ch.setLevel(logging.INFO)
logger.addHandler(_fh)
logger.addHandler(_ch)


# ---------------------------------------------------------------------------
# NATS helpers
# ---------------------------------------------------------------------------

class NATSBus:
    """Thin wrapper around NATS connection + JetStream."""

    def __init__(self) -> None:
        self.nc: Optional[nats.NATS] = None
        self.js = None
        self._presence_kv = None
        self._config_kv = None

    async def connect(self) -> None:
        self.nc = await nats.connect(
            NATS_URL,
            reconnect_time_wait=2,
            max_reconnect_attempts=-1,
            error_cb=self._on_error,
            disconnected_cb=self._on_disconnect,
            reconnected_cb=self._on_reconnect,
        )
        self.js = self.nc.jetstream()
        await self._ensure_buckets()
        logger.info("NATS connected: %s", NATS_URL)

    async def _on_error(self, err):
        logger.warning("NATS error: %s", err)

    async def _on_disconnect(self):
        logger.warning("NATS disconnected")

    async def _on_reconnect(self):
        logger.info("NATS reconnected")
        await self._ensure_buckets()

    async def _ensure_buckets(self) -> None:
        """Create KV buckets if they don't exist."""
        try:
            self._presence_kv = await self.js.key_value(KV_PRESENCE_BUCKET)
        except (NotFoundError, BucketNotFoundError):
            try:
                self._presence_kv = await self.js.create_key_value(
                    KeyValueConfig(bucket=KV_PRESENCE_BUCKET, history=5)
                )
                logger.info("Created KV bucket: %s", KV_PRESENCE_BUCKET)
            except Exception as exc:
                logger.warning("Could not create KV bucket %s: %s", KV_PRESENCE_BUCKET, exc)

        try:
            self._config_kv = await self.js.key_value(KV_CONFIG_BUCKET)
        except (NotFoundError, BucketNotFoundError):
            self._config_kv = None
        except Exception as exc:
            logger.debug("Could not open KV bucket %s: %s", KV_CONFIG_BUCKET, exc)
            self._config_kv = None

    async def publish(self, subject: str, payload: dict) -> None:
        if self.nc and not self.nc.is_closed:
            try:
                await self.nc.publish(subject, json.dumps(payload).encode())
            except Exception as exc:
                logger.debug("NATS publish error on %s: %s", subject, exc)

    async def kv_put(self, bucket_name: str, key: str, value: str) -> None:
        try:
            if bucket_name == KV_PRESENCE_BUCKET and self._presence_kv:
                await self._presence_kv.put(key, value.encode())
        except Exception as exc:
            logger.debug("KV put error %s/%s: %s", bucket_name, key, exc)

    async def kv_get(self, bucket_name: str, key: str) -> Optional[str]:
        try:
            if bucket_name == KV_CONFIG_BUCKET and self._config_kv:
                entry = await self._config_kv.get(key)
                return entry.value.decode() if entry and entry.value else None
        except Exception:
            pass
        return None

    async def get_vad_threshold(self) -> int:
        raw = await self.kv_get(KV_CONFIG_BUCKET, "presence/vad_threshold")
        try:
            return int(raw) if raw else VAD_DEFAULT_THRESHOLD
        except (ValueError, TypeError):
            return VAD_DEFAULT_THRESHOLD

    async def close(self) -> None:
        if self.nc and not self.nc.is_closed:
            await self.nc.drain()


# ---------------------------------------------------------------------------
# VAD (energy-based Voice Activity Detection)
# ---------------------------------------------------------------------------

class VAD:
    """
    Buffers incoming PCM int16 chunks.
    Emits utterances when a silence gap >= VAD_SILENCE_MS is detected.
    """

    def __init__(self, threshold: int = VAD_DEFAULT_THRESHOLD) -> None:
        self.threshold = threshold
        self._buffer: list[bytes] = []
        self._last_voiced_ts: float = 0.0
        self._has_voice = False

    def _rms(self, pcm: bytes) -> float:
        if not pcm:
            return 0.0
        n = len(pcm) // SAMPLE_WIDTH
        if n == 0:
            return 0.0
        samples = struct.unpack(f"<{n}h", pcm[:n * SAMPLE_WIDTH])
        return math.sqrt(sum(s * s for s in samples) / n)

    def push(self, chunk: bytes) -> Optional[bytes]:
        """
        Feed a PCM chunk. Returns a complete utterance (bytes) when silence
        ends the speech segment; otherwise returns None.
        """
        energy = self._rms(chunk)
        now = time.monotonic()

        if energy >= self.threshold:
            self._last_voiced_ts = now
            self._has_voice = True
            self._buffer.append(chunk)
            return None

        # Below threshold — silence
        if not self._has_voice:
            return None  # not started yet

        self._buffer.append(chunk)
        silence_ms = (now - self._last_voiced_ts) * 1000

        if silence_ms >= VAD_SILENCE_MS:
            utterance = b"".join(self._buffer)
            self._buffer.clear()
            self._has_voice = False

            # Discard very short utterances (noise/click)
            duration_ms = len(utterance) / (SAMPLE_RATE * SAMPLE_WIDTH) * 1000
            if duration_ms < VAD_MIN_UTTERANCE_MS:
                logger.debug("VAD: discarding %dms utterance (noise)", duration_ms)
                return None
            logger.debug("VAD: utterance ready, %.0fms", duration_ms)
            return utterance

        return None

    def reset(self) -> None:
        self._buffer.clear()
        self._has_voice = False
        self._last_voiced_ts = 0.0


# ---------------------------------------------------------------------------
# STT (Speech-to-Text, graceful degradation)
# ---------------------------------------------------------------------------

async def _pcm_to_wav(pcm: bytes) -> bytes:
    """Wrap raw PCM int16 mono 16kHz in a minimal WAV container."""
    num_samples = len(pcm) // SAMPLE_WIDTH
    data_size = num_samples * SAMPLE_WIDTH
    header = struct.pack(
        "<4sI4s4sIHHIIHH4sI",
        b"RIFF",
        36 + data_size,
        b"WAVE",
        b"fmt ",
        16,          # chunk size
        1,           # PCM
        1,           # mono
        SAMPLE_RATE,
        SAMPLE_RATE * SAMPLE_WIDTH,  # byte rate
        SAMPLE_WIDTH,
        SAMPLE_WIDTH * 8,            # bits per sample
        b"data",
        data_size,
    )
    return header + pcm


async def stt_whisper(pcm: bytes) -> Optional[str]:
    """Try whisper CLI transcription."""
    whisper_path = None
    for candidate in ["/opt/homebrew/bin/whisper", "/usr/local/bin/whisper"]:
        if Path(candidate).exists():
            whisper_path = candidate
            break
    if not whisper_path:
        # Also check PATH
        proc = await asyncio.create_subprocess_exec(
            "which", "whisper",
            stdout=asyncio.subprocess.PIPE,
            stderr=asyncio.subprocess.DEVNULL,
        )
        out, _ = await proc.communicate()
        if proc.returncode == 0:
            whisper_path = out.decode().strip()

    if not whisper_path:
        return None

    wav_data = await _pcm_to_wav(pcm)
    with tempfile.NamedTemporaryFile(suffix=".wav", delete=False) as f:
        f.write(wav_data)
        wav_path = f.name

    txt_path = wav_path.replace(".wav", ".txt")
    try:
        proc = await asyncio.create_subprocess_exec(
            whisper_path,
            "--model", "tiny.en",
            "--language", "en",
            "--output_format", "txt",
            "--output_dir", str(Path(wav_path).parent),
            wav_path,
            stdout=asyncio.subprocess.PIPE,
            stderr=asyncio.subprocess.PIPE,
        )
        _, stderr = await asyncio.wait_for(proc.communicate(), timeout=30)
        if proc.returncode != 0:
            logger.debug("whisper error: %s", stderr.decode()[:200])
            return None
        if Path(txt_path).exists():
            text = Path(txt_path).read_text().strip()
            return text if text else None
    except asyncio.TimeoutError:
        logger.warning("whisper timed out")
        return None
    except Exception as exc:
        logger.debug("whisper exception: %s", exc)
        return None
    finally:
        for p in [wav_path, txt_path]:
            try:
                Path(p).unlink(missing_ok=True)
            except Exception:
                pass
    return None


async def stt_hermes(pcm: bytes, client: httpx.AsyncClient) -> Optional[str]:
    """Try Hermes /api/audio/transcribe endpoint."""
    try:
        wav_data = await _pcm_to_wav(pcm)
        b64 = base64.b64encode(wav_data).decode()
        data_url = f"data:audio/wav;base64,{b64}"
        resp = await client.post(
            f"{HERMES_DASHBOARD_URL}/api/audio/transcribe",
            json={"data_url": data_url, "mime_type": "audio/wav"},
            headers={"Authorization": f"Bearer {HERMES_DASHBOARD_TOKEN}"},
            timeout=20.0,
        )
        if resp.status_code == 200:
            data = resp.json()
            text = data.get("transcript", "").strip()
            return text if text else None
    except Exception as exc:
        logger.debug("Hermes STT error: %s", exc)
    return None


async def transcribe(pcm: bytes, client: httpx.AsyncClient) -> str:
    """
    Graceful STT degradation:
    1. whisper CLI
    2. Hermes /api/audio/transcribe
    3. Stub fallback
    """
    text = await stt_whisper(pcm)
    if text:
        logger.debug("STT via whisper: %r", text)
        return text

    text = await stt_hermes(pcm, client)
    if text:
        logger.debug("STT via Hermes: %r", text)
        return text

    logger.warning("STT unavailable — returning stub")
    return "[audio received - STT unavailable]"


# ---------------------------------------------------------------------------
# LLM streaming via Hermes API server (OpenAI-compat SSE)
# ---------------------------------------------------------------------------

async def stream_llm(
    text: str,
    client: httpx.AsyncClient,
    on_chunk,          # async callable(str) per token
) -> str:
    """Stream a chat completion from Hermes. Returns full response text."""
    full_response: list[str] = []
    try:
        async with client.stream(
            "POST",
            f"{HERMES_API_URL}/v1/chat/completions",
            json={
                "model": HERMES_MODEL,
                "messages": [{"role": "user", "content": text}],
                "stream": True,
                "max_tokens": 1024,
            },
            headers={
                "Authorization": f"Bearer {HERMES_API_KEY}",
                "Accept": "text/event-stream",
            },
            timeout=httpx.Timeout(60.0, connect=10.0),
        ) as resp:
            resp.raise_for_status()
            async for line in resp.aiter_lines():
                if not line.startswith("data:"):
                    continue
                payload = line[5:].strip()
                if payload == "[DONE]":
                    break
                try:
                    chunk = json.loads(payload)
                    delta = chunk["choices"][0]["delta"]
                    token = delta.get("content", "")
                    if token:
                        full_response.append(token)
                        await on_chunk(token)
                except (json.JSONDecodeError, KeyError, IndexError):
                    pass
    except Exception as exc:
        logger.error("LLM streaming error: %s", exc)
        error_msg = "I encountered an error reaching the AI. Please try again."
        await on_chunk(error_msg)
        return error_msg

    return "".join(full_response)


# ---------------------------------------------------------------------------
# TTS via macOS `say -v Samantha`
# ---------------------------------------------------------------------------

async def tts_say(text: str, ws: WebSocketServerProtocol) -> None:
    """
    Synthesize text using macOS `say -v Samantha`, write AIFF to temp file,
    convert to 16-bit mono 16kHz PCM, then stream chunks back over WebSocket.
    Falls back gracefully if afconvert is unavailable.
    """
    aiff_path = None
    pcm_path = None
    try:
        with tempfile.NamedTemporaryFile(suffix=".aiff", delete=False) as f:
            aiff_path = f.name
        with tempfile.NamedTemporaryFile(suffix=".wav", delete=False) as f:
            pcm_path = f.name

        # Synthesize to AIFF (standard format — --data-format flag not supported on all macOS)
        proc = await asyncio.create_subprocess_exec(
            "say",
            "-v", "Samantha",
            "-o", aiff_path,
            text,
            stdout=asyncio.subprocess.DEVNULL,
            stderr=asyncio.subprocess.PIPE,
        )
        _, stderr = await asyncio.wait_for(proc.communicate(), timeout=60)
        if proc.returncode != 0:
            logger.error("say error: %s", stderr.decode()[:200])
            return

        # Convert AIFF → 16-bit mono 16kHz WAV via afconvert
        conv_proc = await asyncio.create_subprocess_exec(
            "afconvert",
            "-f", "WAVE",
            "-d", "LEI16@16000",
            "-c", "1",
            aiff_path,
            pcm_path,
            stdout=asyncio.subprocess.DEVNULL,
            stderr=asyncio.subprocess.PIPE,
        )
        _, conv_err = await asyncio.wait_for(conv_proc.communicate(), timeout=30)

        if conv_proc.returncode != 0:
            logger.warning("afconvert error: %s — sending raw AIFF instead", conv_err.decode()[:100])
            # Fall back: stream the raw AIFF file as binary
            audio_bytes = Path(aiff_path).read_bytes()
        else:
            audio_bytes = Path(pcm_path).read_bytes()
            # Strip 44-byte WAV header to send raw PCM to Flutter
            if audio_bytes[:4] == b"RIFF":
                audio_bytes = audio_bytes[44:]

        # Notify client that audio is coming
        await ws.send(json.dumps({"type": "audio_ready"}))

        # Stream in chunks
        for offset in range(0, len(audio_bytes), TTS_CHUNK_BYTES):
            chunk = audio_bytes[offset: offset + TTS_CHUNK_BYTES]
            await ws.send(chunk)
            await asyncio.sleep(0)  # yield to event loop

        logger.debug("TTS streamed %d bytes", len(audio_bytes))

    except asyncio.TimeoutError:
        logger.error("TTS timed out")
    except Exception as exc:
        logger.error("TTS error: %s", exc)
    finally:
        for p in [aiff_path, pcm_path]:
            if p:
                try:
                    Path(p).unlink(missing_ok=True)
                except Exception:
                    pass


# ---------------------------------------------------------------------------
# Session
# ---------------------------------------------------------------------------

class PresenceSession:
    """Represents one active Presence Mode session."""

    def __init__(
        self,
        ws: WebSocketServerProtocol,
        bus: NATSBus,
        http: httpx.AsyncClient,
    ) -> None:
        self.ws = ws
        self.bus = bus
        self.http = http
        self.session_id = f"presence_{int(time.time())}"
        self.vad = VAD()
        self._interrupted = asyncio.Event()
        self._tts_lock = asyncio.Lock()

    # ------------------------------------------------------------------
    # WebSocket helpers
    # ------------------------------------------------------------------

    async def send_json(self, payload: dict) -> None:
        try:
            await self.ws.send(json.dumps(payload))
        except Exception:
            pass

    async def set_state(self, state: str) -> None:
        await self.send_json({"type": "state", "value": state})
        await self.bus.publish(
            f"presence.session.{self.session_id}.thinking",
            {"status": state, "session_id": self.session_id},
        )

    # ------------------------------------------------------------------
    # Session lifecycle
    # ------------------------------------------------------------------

    async def start(self) -> None:
        ts = time.time()
        await self.bus.kv_put(KV_PRESENCE_BUCKET, "active_session", self.session_id)
        await self.bus.publish(
            f"presence.session.{self.session_id}.started",
            {"session_id": self.session_id, "ts": ts},
        )
        # Load VAD threshold from config KV
        threshold = await self.bus.get_vad_threshold()
        self.vad.threshold = threshold
        logger.info("Session started: %s (VAD threshold=%d)", self.session_id, threshold)
        await self.set_state("listening")

    async def end(self) -> None:
        await self.bus.kv_put(KV_PRESENCE_BUCKET, "active_session", "")
        await self.bus.publish(
            f"presence.session.{self.session_id}.ended",
            {"session_id": self.session_id, "ts": time.time()},
        )
        logger.info("Session ended: %s", self.session_id)

    # ------------------------------------------------------------------
    # Pipeline: utterance → STT → LLM → TTS
    # ------------------------------------------------------------------

    async def handle_utterance(self, pcm: bytes) -> None:
        """Full STT → LLM → TTS pipeline for one utterance."""
        self._interrupted.clear()

        # STT
        await self.set_state("thinking")
        transcript = await transcribe(pcm, self.http)
        logger.info("Zeus said: %r", transcript)

        ts = time.time()
        zeus_transcript = {
            "type": "transcript",
            "speaker": "zeus",
            "text": transcript,
            "ts": ts,
        }
        await self.send_json(zeus_transcript)
        await self.bus.publish(
            f"presence.session.{self.session_id}.transcript",
            {"text": transcript, "speaker": "zeus", "ts": ts, "session_id": self.session_id},
        )
        # Publish to EVT stream for ClickHouse
        await self.bus.publish(
            EVT_STREAM_SUBJECT,
            {
                "session_id": self.session_id,
                "speaker": "zeus",
                "text": transcript,
                "ts": ts,
                "duration_samples": len(pcm) // SAMPLE_WIDTH,
            },
        )

        if transcript == "[audio received - STT unavailable]":
            await self.set_state("listening")
            return

        # LLM streaming
        response_chunks: list[str] = []

        async def on_chunk(token: str) -> None:
            if self._interrupted.is_set():
                return
            response_chunks.append(token)
            await self.send_json({"type": "chunk", "text": token})
            await self.bus.publish(
                f"presence.session.{self.session_id}.chunk",
                {"text": token, "session_id": self.session_id},
            )

        full_response = await stream_llm(transcript, self.http, on_chunk)

        if self._interrupted.is_set():
            logger.info("Interrupted after LLM response")
            await self.set_state("listening")
            return

        # Publish full Legacy transcript
        ts_resp = time.time()
        await self.send_json({
            "type": "transcript",
            "speaker": "legacy",
            "text": full_response,
            "ts": ts_resp,
        })
        await self.bus.publish(
            f"presence.session.{self.session_id}.transcript",
            {
                "text": full_response,
                "speaker": "legacy",
                "ts": ts_resp,
                "session_id": self.session_id,
            },
        )
        await self.bus.publish(
            EVT_STREAM_SUBJECT,
            {
                "session_id": self.session_id,
                "speaker": "legacy",
                "text": full_response,
                "ts": ts_resp,
            },
        )
        logger.info("Legacy replied (%d chars)", len(full_response))

        # TTS
        if not self._interrupted.is_set():
            await self.set_state("speaking")
            async with self._tts_lock:
                await tts_say(full_response, self.ws)

        await self.set_state("listening")

    # ------------------------------------------------------------------
    # Control frame handler
    # ------------------------------------------------------------------

    async def handle_control(self, payload: dict) -> None:
        action = payload.get("action")
        if action == "stop":
            logger.info("Control: stop — ending session")
            raise SessionStop()
        elif action == "interrupt":
            logger.info("Control: interrupt")
            self._interrupted.set()
            self.vad.reset()
            await self.set_state("listening")
        elif action == "start":
            logger.info("Control: start (no-op, already running)")

    # ------------------------------------------------------------------
    # Main receive loop
    # ------------------------------------------------------------------

    async def run(self) -> None:
        utterance_task: Optional[asyncio.Task] = None
        try:
            async for message in self.ws:
                if isinstance(message, bytes):
                    # PCM audio chunk from Flutter
                    utterance = self.vad.push(message)
                    if utterance:
                        # Cancel previous pipeline if still running
                        if utterance_task and not utterance_task.done():
                            self._interrupted.set()
                            utterance_task.cancel()
                            try:
                                await utterance_task
                            except asyncio.CancelledError:
                                pass
                        self._interrupted.clear()
                        utterance_task = asyncio.create_task(
                            self.handle_utterance(utterance)
                        )
                elif isinstance(message, str):
                    try:
                        payload = json.loads(message)
                    except json.JSONDecodeError:
                        logger.warning("Unreadable text frame: %r", message[:80])
                        continue
                    if payload.get("type") == "control":
                        await self.handle_control(payload)

        except SessionStop:
            pass
        except websockets.exceptions.ConnectionClosed as exc:
            logger.info("WebSocket closed: %s", exc)
        finally:
            if utterance_task and not utterance_task.done():
                utterance_task.cancel()
                try:
                    await utterance_task
                except asyncio.CancelledError:
                    pass


class SessionStop(Exception):
    pass


# ---------------------------------------------------------------------------
# Hermes probe
# ---------------------------------------------------------------------------

async def probe_hermes(client: httpx.AsyncClient) -> bool:
    """Verify Hermes API server is reachable and auth works."""
    try:
        resp = await client.post(
            f"{HERMES_API_URL}/v1/chat/completions",
            json={
                "model": HERMES_MODEL,
                "messages": [{"role": "user", "content": "ping"}],
                "stream": False,
                "max_tokens": 5,
            },
            headers={"Authorization": f"Bearer {HERMES_API_KEY}"},
            timeout=15.0,
        )
        if resp.status_code == 200:
            data = resp.json()
            content = data["choices"][0]["message"]["content"]
            logger.info("Hermes API probe OK — response: %r", content[:60])
            return True
        logger.warning("Hermes API probe: HTTP %d", resp.status_code)
    except Exception as exc:
        logger.warning("Hermes API probe failed: %s", exc)
    return False


async def probe_hermes_dashboard(client: httpx.AsyncClient) -> bool:
    """Verify Hermes dashboard transcribe endpoint."""
    try:
        resp = await client.get(
            f"{HERMES_DASHBOARD_URL}/api/status",
            timeout=5.0,
        )
        if resp.status_code == 200:
            data = resp.json()
            logger.info(
                "Hermes dashboard OK — version=%s state=%s",
                data.get("version"),
                data.get("gateway_state"),
            )
            return True
    except Exception as exc:
        logger.warning("Hermes dashboard probe failed: %s", exc)
    return False


# ---------------------------------------------------------------------------
# WebSocket server
# ---------------------------------------------------------------------------

class PresenceServer:
    """Single-connection WebSocket server for Presence Mode."""

    def __init__(self) -> None:
        self.bus = NATSBus()
        self._current_session: Optional[PresenceSession] = None

    async def handler(self, ws: WebSocketServerProtocol) -> None:
        remote = ws.remote_address
        logger.info("New connection from %s", remote)

        # Reject concurrent sessions
        if self._current_session is not None:
            await ws.send(json.dumps({
                "type": "error",
                "message": "A Presence session is already active. Only one session at a time.",
            }))
            await ws.close(1008, "Session busy")
            return

        async with httpx.AsyncClient() as http:
            session = PresenceSession(ws, self.bus, http)
            self._current_session = session
            try:
                await session.start()
                await session.run()
            finally:
                await session.end()
                self._current_session = None
                logger.info("Connection from %s closed", remote)

    async def run(self) -> None:
        await self.bus.connect()

        async with httpx.AsyncClient() as probe_client:
            api_ok = await probe_hermes(probe_client)
            dash_ok = await probe_hermes_dashboard(probe_client)

        if not api_ok:
            logger.warning(
                "Hermes API server at %s not responding — LLM calls will fail",
                HERMES_API_URL,
            )
        if not dash_ok:
            logger.warning(
                "Hermes dashboard at %s not responding — Hermes STT fallback unavailable",
                HERMES_DASHBOARD_URL,
            )

        logger.info("Presence Mode WebSocket listening on ws://%s:%d", WS_HOST, WS_PORT)

        async with websockets.serve(
            self.handler,
            WS_HOST,
            WS_PORT,
            ping_interval=20,
            ping_timeout=30,
            max_size=2 * 1024 * 1024,  # 2 MB max frame (generous for audio chunks)
        ):
            await asyncio.Future()  # run forever


# ---------------------------------------------------------------------------
# Entry point
# ---------------------------------------------------------------------------

def main() -> None:
    import argparse

    global WS_HOST, WS_PORT

    parser = argparse.ArgumentParser(description="OneMind Presence Mode voice agent")
    parser.add_argument("--host", default=WS_HOST, help="WebSocket bind host")
    parser.add_argument("--port", type=int, default=WS_PORT, help="WebSocket port")
    args = parser.parse_args()

    WS_HOST = args.host
    WS_PORT = args.port

    server = PresenceServer()
    try:
        asyncio.run(server.run())
    except KeyboardInterrupt:
        logger.info("Presence Mode shutting down")


if __name__ == "__main__":
    main()
