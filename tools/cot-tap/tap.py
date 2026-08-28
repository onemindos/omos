#!/usr/bin/env python3
"""
OneMind CoT Tap — TAK Server -> NATS bus (outbound / b3 outbound leg).

Connects to a TAK Server CoT streaming port (8089/ssl) as a mutual-TLS client
(exactly like an EUD), receives the full CoT firehose for its groups, and
republishes each <event>...</event> to a NATS subject (default tak.cot.out).

This exists because TAK's cluster plugin manager (PluginClusterMessenger) is
send-only: sender-receiver plugins can INJECT (bus->TAK) but never RECEIVE
(TAK->bus) across pods. The tap is the authoritative client-protocol path for
the outbound leg, decoupled from plugin-cluster internals.

Env:
  TAK_HOST         (default takserver-messaging)
  TAK_PORT         (default 8089)
  TAK_CERT         (default /certs/tap.pem)
  TAK_KEY          (default /certs/tap.key)
  TAK_CA           (default /certs/ca.pem)   # optional; verify off by default
  NATS_URL         (default nats://tak-nats-leaf:4222)
  NATS_SUBJECT     (default tak.cot.out)
  VERIFY_SERVER    (default "false")         # set "true" to verify server cert vs CA
"""
import asyncio
import os
import ssl
import sys
import logging

import nats

LOG = logging.getLogger("cot-tap")
logging.basicConfig(
    level=logging.INFO,
    format="%(asctime)s %(levelname)s %(message)s",
    stream=sys.stdout,
)

TAK_HOST = os.getenv("TAK_HOST", "takserver-messaging")
TAK_PORT = int(os.getenv("TAK_PORT", "8089"))
TAK_CERT = os.getenv("TAK_CERT", "/certs/tap.pem")
TAK_KEY = os.getenv("TAK_KEY", "/certs/tap.key")
TAK_CA = os.getenv("TAK_CA", "/certs/ca.pem")
NATS_URL = os.getenv("NATS_URL", "nats://tak-nats-leaf:4222")
NATS_SUBJECT = os.getenv("NATS_SUBJECT", "tak.cot.out")
VERIFY_SERVER = os.getenv("VERIFY_SERVER", "false").lower() == "true"

EVENT_END = b"</event>"


def build_ssl_context() -> ssl.SSLContext:
    ctx = ssl.SSLContext(ssl.PROTOCOL_TLS_CLIENT)
    if VERIFY_SERVER and os.path.exists(TAK_CA):
        ctx.load_verify_locations(cafile=TAK_CA)
        ctx.check_hostname = False
        ctx.verify_mode = ssl.CERT_REQUIRED
    else:
        ctx.check_hostname = False
        ctx.verify_mode = ssl.CERT_NONE
    ctx.load_cert_chain(certfile=TAK_CERT, keyfile=TAK_KEY)
    return ctx


async def connect_nats():
    while True:
        try:
            nc = await nats.connect(
                NATS_URL,
                name="omos-cot-tap",
                max_reconnect_attempts=-1,
                reconnect_time_wait=2,
            )
            LOG.info("connected to NATS %s -> subject %s", NATS_URL, NATS_SUBJECT)
            return nc
        except Exception as e:
            LOG.error("NATS connect failed (%s); retry in 2s", e)
            await asyncio.sleep(2)


async def stream_tak(nc):
    ssl_ctx = build_ssl_context()
    reader, writer = await asyncio.open_connection(TAK_HOST, TAK_PORT, ssl=ssl_ctx)
    LOG.info("TLS connected to TAK %s:%s", TAK_HOST, TAK_PORT)
    buf = b""
    count = 0
    try:
        while True:
            chunk = await reader.read(8192)
            if not chunk:
                LOG.warning("TAK stream closed by server")
                break
            buf += chunk
            # frame on </event>
            while True:
                idx = buf.find(EVENT_END)
                if idx == -1:
                    break
                end = idx + len(EVENT_END)
                raw = buf[:end]
                buf = buf[end:]
                start = raw.find(b"<event")
                if start == -1:
                    continue
                event = raw[start:]
                try:
                    await nc.publish(NATS_SUBJECT, event)
                    count += 1
                    if count <= 5 or count % 50 == 0:
                        LOG.info("published CoT #%d (%d bytes) to %s", count, len(event), NATS_SUBJECT)
                except Exception as e:
                    LOG.error("NATS publish failed: %s", e)
    finally:
        writer.close()
        try:
            await writer.wait_closed()
        except Exception:
            pass


async def main():
    LOG.info(
        "OneMind CoT Tap starting — TAK %s:%s -> NATS %s [%s]",
        TAK_HOST, TAK_PORT, NATS_URL, NATS_SUBJECT,
    )
    nc = await connect_nats()
    while True:
        try:
            await stream_tak(nc)
        except Exception as e:
            LOG.error("TAK stream error: %s", e)
        LOG.info("reconnecting to TAK in 3s...")
        await asyncio.sleep(3)


if __name__ == "__main__":
    try:
        asyncio.run(main())
    except KeyboardInterrupt:
        pass
