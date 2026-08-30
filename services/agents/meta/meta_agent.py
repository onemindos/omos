#!/usr/bin/env python3
"""omos-meta — the Meta Agent (conductor) for the Legacy Agent Fabric.

Speaks the Synadia Agent Protocol as a *caller*: discovers agents on the NATS
fabric, watches liveness, prompts any of them, and can fan a prompt out to all.

Protocol (NATS-native, no registry):
  - discovery : request $SRV.INFO.agents, every live agent replies with its info
  - prompt    : request/stream to an agent's discovered `prompt` endpoint subject
  - liveness  : heartbeats on agents.hb.*.*.*  (offline after 3 missed beats)

Ref: https://nats.io/blog/nats-native-protocol-for-ai-agents/
"""
import argparse
import asyncio
import json
import os
import time
from typing import Any

import nats

NATS_URL = os.environ.get("NATS_URL", "nats://localhost:4222")
SRV_INFO_AGENTS = "$SRV.INFO.agents"
HB_WILDCARD = "agents.hb.*.*.*"


async def _connect():
    # Extra NATS_* env (user/pass/creds/tls) can be wired here for the cluster bus.
    opts: dict[str, Any] = {"servers": [NATS_URL], "name": "omos-meta"}
    if os.environ.get("NATS_USER"):
        opts["user"] = os.environ["NATS_USER"]
        opts["password"] = os.environ.get("NATS_PASSWORD", "")
    if os.environ.get("NATS_CREDS"):
        opts["user_credentials"] = os.environ["NATS_CREDS"]
    return await nats.connect(**opts)


async def discover(nc, window: float = 1.0) -> list[dict]:
    """Collect $SRV.INFO.agents replies from every live agent for `window` seconds."""
    inbox = nc.new_inbox()
    results: list[dict] = []
    sub = await nc.subscribe(inbox)

    async def _drain():
        async for msg in sub.messages:
            try:
                results.append(json.loads(msg.data))
            except json.JSONDecodeError:
                pass

    task = asyncio.create_task(_drain())
    await nc.publish(SRV_INFO_AGENTS, b"", reply=inbox)
    await asyncio.sleep(window)
    await sub.unsubscribe()
    task.cancel()
    return results


def _prompt_subject(info: dict) -> str | None:
    for ep in info.get("endpoints", []):
        if ep.get("name") == "prompt" and ep.get("subject"):
            return ep["subject"]
    return None


def _label(info: dict) -> str:
    md = info.get("metadata", {}) or {}
    ident = md.get("agent") or info.get("name", "?")
    owner = md.get("owner", "?")
    return f"{ident}/{owner} (id={info.get('id', '?')[:8]} v{info.get('version', '?')})"


async def _collect_stream(nc, subject: str, text: str, timeout: float = 30.0) -> str:
    """Send a prompt and accumulate streamed `response` chunks until the terminator."""
    inbox = nc.new_inbox()
    sub = await nc.subscribe(inbox)
    await nc.publish(subject, text.encode(), reply=inbox)
    out: list[str] = []
    deadline = time.monotonic() + timeout
    while time.monotonic() < deadline:
        try:
            msg = await asyncio.wait_for(sub.next_msg(timeout=timeout), timeout=timeout)
        except asyncio.TimeoutError:
            break
        if not msg.data:  # empty-body terminator (§6.5)
            break
        try:
            chunk = json.loads(msg.data)
            if chunk.get("type") == "response":
                out.append(str(chunk.get("data", "")))
            elif chunk.get("type") == "error":
                out.append(f"[error] {chunk.get('data')}")
                break
        except json.JSONDecodeError:
            out.append(msg.data.decode(errors="replace"))
    await sub.unsubscribe()
    return "".join(out)


# ---- CLI commands ----
async def cmd_list(_args):
    nc = await _connect()
    agents = await discover(nc)
    if not agents:
        print("no agents found on the fabric")
    for a in agents:
        print(f"- {_label(a)}  prompt={_prompt_subject(a)}")
    await nc.close()


async def cmd_watch(_args):
    nc = await _connect()
    seen: dict[str, float] = {}
    print(f"watching {HB_WILDCARD} (Ctrl-C to stop)...")
    sub = await nc.subscribe(HB_WILDCARD)
    async for msg in sub.messages:
        try:
            hb = json.loads(msg.data)
        except json.JSONDecodeError:
            continue
        key = f"{hb.get('agent')}/{hb.get('owner')}/{hb.get('name')}:{hb.get('instance_id','')[:8]}"
        now = time.monotonic()
        gap = f"+{now - seen[key]:.0f}s" if key in seen else "new"
        seen[key] = now
        print(f"  [{time.strftime('%H:%M:%S')}] {key}  ({gap}, cadence={hb.get('cadence','?')}s)")


async def cmd_prompt(args):
    nc = await _connect()
    agents = await discover(nc)
    target = next(
        (a for a in agents
         if (a.get("metadata", {}) or {}).get("agent") == args.agent
         and (a.get("metadata", {}) or {}).get("owner") == args.owner),
        None,
    )
    subject = _prompt_subject(target) if target else (
        f"agents.prompt.{args.agent}.{args.owner}.{args.name}")
    print(f"→ {subject}")
    reply = await _collect_stream(nc, subject, args.text)
    print(reply or "(no response)")
    await nc.close()


async def cmd_fanout(args):
    nc = await _connect()
    agents = await discover(nc)
    subjects = {s for s in (_prompt_subject(a) for a in agents) if s}
    print(f"fanning out to {len(subjects)} agent(s)...")
    async def one(subj):
        return subj, await _collect_stream(nc, subj, args.text)
    for subj, reply in await asyncio.gather(*(one(s) for s in subjects)):
        print(f"\n=== {subj} ===\n{reply or '(no response)'}")
    await nc.close()


def main():
    p = argparse.ArgumentParser(prog="meta_agent", description="omos-meta — Meta Agent (caller)")
    sub = p.add_subparsers(dest="cmd", required=True)
    sub.add_parser("list", help="discover agents on the fabric")
    sub.add_parser("watch", help="live heartbeat board")
    pp = sub.add_parser("prompt", help="prompt one agent")
    pp.add_argument("--agent", required=True); pp.add_argument("--owner", required=True)
    pp.add_argument("--name", default="default"); pp.add_argument("text")
    pf = sub.add_parser("fanout", help="prompt every agent and collect responses")
    pf.add_argument("text")
    args = p.parse_args()
    fn = {"list": cmd_list, "watch": cmd_watch, "prompt": cmd_prompt, "fanout": cmd_fanout}[args.cmd]
    try:
        asyncio.run(fn(args))
    except KeyboardInterrupt:
        pass


if __name__ == "__main__":
    main()
