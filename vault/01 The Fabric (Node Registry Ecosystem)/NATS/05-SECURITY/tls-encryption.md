---
title: "TLS Encryption"
section: "05-SECURITY"
type: note
tags: [nats, security, tls, mtls, certificates, encryption]
source: https://docs.nats.io/running-a-nats-service/configuration/securing_nats/tls
created: 2026-06-14
---

# TLS Encryption

Encrypt NATS connections and optionally authenticate clients by certificate (mTLS).

See also: [[NATS/05-SECURITY/authentication|Authentication]] · [[NATS/05-SECURITY/README|Security Index]]

---

## Server-Side TLS

```text
# nats-server.conf
tls {
  cert_file: "./server-cert.pem"
  key_file:  "./server-key.pem"
  ca_file:   "./ca.pem"          # optional, for verifying client certs
  verify:    true                # require + verify client certs (mTLS)
  timeout:   2                   # TLS handshake timeout (seconds)
}
```

- Without `verify`, TLS encrypts the channel but clients aren't authenticated by cert.
- With **`verify: true`**, clients **must** present a cert signed by `ca_file` — this is **mTLS** (mutual TLS).
- With **`verify_and_map: true`**, the client cert's identity (CN/SAN) is **mapped to a NATS user**, so the certificate *is* the authentication (no password needed).

When TLS is configured, the server's `INFO` sets `tls_required: true`; clients must upgrade to TLS before `CONNECT`.

---

## Client-Side

```go
nc, _ := nats.Connect("tls://nats.example.com:4222",
    nats.RootCAs("ca.pem"),
    nats.ClientCert("client-cert.pem", "client-key.pem"))  // for mTLS
```

```python
import ssl
ctx = ssl.create_default_context(cafile="ca.pem")
ctx.load_cert_chain("client-cert.pem", "client-key.pem")   # mTLS
nc = await nats.connect("tls://nats.example.com:4222", tls=ctx)
```

---

## TLS Everywhere — per-link

TLS is configured **independently per connection type**, each with its own `tls { }` block:

| Link | Block |
|------|-------|
| Client ↔ server | top-level `tls { }` |
| Cluster routes | `cluster { tls { } }` |
| Gateways (supercluster) | `gateway { tls { } }` |
| Leaf nodes | `leafnodes { tls { } }` |
| WebSocket | `websocket { tls { } }` |
| MQTT | `mqtt { tls { } }` |
| Monitoring (HTTPS) | `https_port` + `tls { }` |

Secure **every** hop in a real deployment — an unencrypted cluster route or leaf link is a hole even if client TLS is on.

---

## mTLS as Authentication (`verify_and_map`)

```text
tls {
  cert_file: "server-cert.pem"
  key_file:  "server-key.pem"
  ca_file:   "ca.pem"
  verify_and_map: true
}
authorization {
  users: [
    { user: "CN=client.example.com,O=OneMind" }   # matches the cert subject
  ]
}
```

The cert's subject/SAN maps to a NATS user — certificates become identities, no shared secrets. Strong for fleets where each device has its own cert.

---

## Operational Notes

- **Handshake timeout** (`timeout`) — raise on slow/edge links.
- **Cipher suites / curve preferences** — configurable for compliance.
- **Cert rotation** — `nats-server --signal reload` (SIGHUP) **does** reload TLS cert files in place (unlike accounts/system_account, which need a full restart). Verify with a fresh connection after rotation.
- **OCSP stapling** — supported for cert revocation checking.
- **Insecure skip-verify** exists for dev only — never in production.

---

## OneMind Relevance

The OneMind bus binds to the **Tailscale** interface and relies on Tailscale's WireGuard encryption for transport today (the doctrine: "Tailscale-only, never expose publicly"). TLS at the NATS layer becomes important when (a) leaf nodes connect over non-Tailscale links, (b) browser/WebSocket clients connect over the internet via Funnel, or (c) per-device cert identity is wanted for the fleet. mTLS `verify_and_map` is the natural path to per-node certificate identity as the fleet scales.

## Sources
- Securing NATS / TLS — https://docs.nats.io/running-a-nats-service/configuration/securing_nats/tls
- TLS client cert mapping — https://docs.nats.io/running-a-nats-service/configuration/securing_nats/auth_intro/tls_mutual_auth
- Config reload — https://docs.nats.io/running-a-nats-service/nats_admin/signals
