#!/usr/bin/env python3
"""Hermes -> NATS agent host (Synadia Agent Protocol).

Puts a Hermes instance ON THE BUS as a first-class agent so it can be prompted
over NATS from anywhere on the fabric. This is the piece `hermes gateway run`
does NOT provide (that's chat-platform only). Runs INSIDE the Hermes container,
using its bundled venv (synadia_ai + nats-py + the hermes CLI are all present).

Registers:  agents.prompt.{agent}.{owner}.{session}
  default:  agents.prompt.hermes.zeus.legacy   (the hub's main Hermes = Legacy)

On each prompt it drives the local Hermes one-shot (`hermes -z <prompt>`) and
streams the reply back. Hermes's own memory/session state lives in HERMES_HOME
(/opt/data, a persistent volume), so continuity survives across prompts and
pod restarts.

Env:
  NATS_URL                 default nats://localhost:4222 (the in-pod leaf)
  HERMES_NATS_OWNER        default zeus
  HERMES_NATS_SESSION_NAME default legacy
  HERMES_HOME              default /opt/data
  HERMES_BIN               default /opt/hermes/bin/hermes
  HERMES_MODEL             default us.anthropic.claude-haiku-4-5-20251001-v1:0 (Bedrock)
  HERMES_PROMPT_TIMEOUT    default 180 (seconds)

Ref: omos-voice/scripts/legacy_stub.py (API template); the Synadia Agent
Protocol; docs/legacy-agent-fabric.md.
"""
from __future__ import annotations

import asyncio
import os
import signal
from pathlib import Path

import nats
from synadia_ai.agent_service import AgentService, PromptStream
from synadia_ai.agents import Envelope

NATS_URL = os.environ.get("NATS_URL", "nats://localhost:4222")
OWNER = os.environ.get("HERMES_NATS_OWNER", "zeus")
SESSION = os.environ.get("HERMES_NATS_SESSION_NAME", "legacy")
HERMES_HOME = os.environ.get("HERMES_HOME", "/opt/data")
HERMES_BIN = os.environ.get("HERMES_BIN", "/opt/hermes/bin/hermes")
HERMES_MODEL = os.environ.get("HERMES_MODEL", "us.anthropic.claude-haiku-4-5-20251001-v1:0")
PROMPT_TIMEOUT = float(os.environ.get("HERMES_PROMPT_TIMEOUT", "180"))


def ensure_model_config() -> None:
    """Self-heal: write a minimal Bedrock model config if HERMES_HOME has none.

    Keeps the deploy reproducible (GitOps) — a fresh PVC boots with a valid
    model instead of an invalid default that Bedrock rejects.
    """
    cfg = Path(HERMES_HOME) / "config.yaml"
    if cfg.exists():
        return
    cfg.write_text(
        "model:\n"
        "  provider: bedrock\n"
        f"  default: {HERMES_MODEL}\n"
        "  context_length: 131072\n"
        "  base_url: https://bedrock-runtime.us-east-1.amazonaws.com\n",
        encoding="utf-8",
    )
    print(f"[host] wrote default model config -> {cfg} ({HERMES_MODEL})", flush=True)


async def run_hermes(prompt: str) -> str:
    """Drive the local Hermes one-shot and return its reply text."""
    env = {**os.environ, "HOME": HERMES_HOME, "TERM": "dumb", "NO_COLOR": "1"}
    proc = await asyncio.create_subprocess_exec(
        HERMES_BIN, "-z", prompt,
        stdout=asyncio.subprocess.PIPE,
        stderr=asyncio.subprocess.PIPE,
        cwd=HERMES_HOME,
        env=env,
    )
    try:
        out, err = await asyncio.wait_for(proc.communicate(), timeout=PROMPT_TIMEOUT)
    except asyncio.TimeoutError:
        proc.kill()
        return "Legacy timed out generating a reply."
    text = (out or b"").decode(errors="replace").strip()
    if not text:
        tail = (err or b"").decode(errors="replace").strip()
        text = tail[-800:] if tail else "(no output)"
    return text


async def main() -> None:
    ensure_model_config()
    nc = await nats.connect(NATS_URL, name=f"hermes-host-{SESSION}")
    print(f"[host] connected to {NATS_URL} as agent=hermes owner={OWNER} session={SESSION}", flush=True)

    async def handler(envelope: Envelope, stream: PromptStream) -> None:
        text = (envelope.prompt or "").strip()
        print(f"[host] prompt: {text!r}", flush=True)
        reply = await run_hermes(text)
        print(f"[host] reply ({len(reply)} chars)", flush=True)
        await stream.send(reply)

    service = AgentService(
        agent="hermes", owner=OWNER, session_name=SESSION, nc=nc,
        description=f"Hermes on the hub ({SESSION}) — driven one-shot via NATS",
        heartbeat_interval_s=30,
    )
    service.on_prompt(handler)
    await service.start()
    print(f"[host] listening on {service.subject.prompt}", flush=True)

    stop = asyncio.Event()
    loop = asyncio.get_running_loop()
    for s in (signal.SIGINT, signal.SIGTERM):
        loop.add_signal_handler(s, stop.set)
    try:
        await stop.wait()
    finally:
        await service.stop()
        await nc.close()


if __name__ == "__main__":
    asyncio.run(main())
