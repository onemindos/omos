"""Connection settings shared by the subscription thread and the tools.

OneMind change (see ONEMIND-CHANGES.md): the plugin used to dial a local leaf node on
nats://localhost:4222 with no credentials and no TLS. On onemind-hermes it connects directly to the
cluster bus, so the URL, credentials and CA come from the environment. With none of the
variables set it behaves as before.

  NATS_URL                  default nats://localhost:4222
  NATS_USER / NATS_PASSWORD from a Secret
  NATS_CA_FILE              PEM bundle that verifies the server certificate
  NATS_TLS_SERVER_NAME      name to verify (the Service DNS name)
  NATS_CONTEXT_CLIENT_NAME  shows in /connz (default "onemind-hermes-nats-context")
"""
from __future__ import annotations

import os
from typing import Any, Dict


def nats_url() -> str:
    return os.environ.get("NATS_URL", "").strip() or "nats://localhost:4222"


def connect_options() -> Dict[str, Any]:
    opts: Dict[str, Any] = {}
    user = os.environ.get("NATS_USER", "").strip()
    if user:
        opts["user"] = user
        opts["password"] = os.environ.get("NATS_PASSWORD", "")
    ca_file = os.environ.get("NATS_CA_FILE", "").strip()
    if ca_file:
        import ssl

        opts["tls"] = ssl.create_default_context(cafile=ca_file)
        server_name = os.environ.get("NATS_TLS_SERVER_NAME", "").strip()
        if server_name:
            opts["tls_hostname"] = server_name
    opts["name"] = os.environ.get("NATS_CONTEXT_CLIENT_NAME", "").strip() or "onemind-hermes-nats-context"
    return opts
