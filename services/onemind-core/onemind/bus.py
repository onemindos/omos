"""Authenticated NATS connection with the platform CA and a bounded startup."""
import asyncio
import logging
import os
import ssl
import nats


async def connect():
    async def error(error):
        logging.warning("NATS transport unavailable (%s)", type(error).__name__)

    context = ssl.create_default_context(cafile=os.environ.get("NATS_CA_FILE"))
    return await asyncio.wait_for(nats.connect(
        servers=os.environ["NATS_URL"], user=os.environ["NATS_USER"], password=os.environ["NATS_PASSWORD"],
        tls=context, tls_hostname=os.environ.get("NATS_TLS_SERVER_NAME", "nats.nats.svc.cluster.local"),
        connect_timeout=5, max_reconnect_attempts=3, error_cb=error,
    ), timeout=20)
