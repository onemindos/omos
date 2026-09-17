---
id: gateway-auth-runbook
title: Gateway Auth & Device Pairing — Runbook
tags: [openclaw, gateway, auth, security, runbook, #p1]
owner: Zeus
status: active
created: 2026-05-12
updated: 2026-05-12
---

## Gateway Auth & Device Pairing — Runbook

> **This document exists because of a 72-hour token rotation loop resolved on 2026-05-12.** Read it before touching gateway auth, device pairing, or `openclaw.json`.

---

## The Golden Rules

1. `gateway.auth.token` must **never** exist in `openclaw.json`
2. `OPENCLAW_GATEWAY_TOKEN` in `gateway.systemd.env` is the **only** source of truth for the gateway token
3. The gateway identity device (`0b8946a3...`) must always be in `paired.json` with **full operator scopes**
4. `defaultSpawnContext` must stay `isolated` — never `fork`

Break any of these and you will get the mismatch/pairing loop.

---

## Active Credentials

| Item | Value |
|------|-------|
| Gateway token | `onemind-zeus-gateway-edb055e0cea082e8f40d5bed75d20a1f` |
| Control UI URL | `https://omos-core.tail717ba0.ts.net` (Tailscale only) |
| Exec approvals | Telegram DM via `@onemindos_bot` |

---

## Device Registry

| Device | ID (first 20 chars) | What it is | Required scopes |
|--------|---------------------|------------|-----------------|
| Gateway server identity | `0b8946a3b3133061b5ad...` | omos-core's own internal identity | Full operator |
| Zeus Mac browser | `850cae8d09569e1e7194...` | Control UI from `100.114.142.55` | Full operator |

Both must be in `~/.openclaw/devices/paired.json` at all times. If either disappears, re-approve from `openclaw devices list`.

---

## Root Cause — What Was Breaking

OpenClaw's `ensureGatewayStartupAuth` runs on **every config write** (not just restarts). If `gateway.auth.token` exists in `openclaw.json`, it compares the file value to the in-memory value. Any mismatch → rotates the token → writes new value to file → triggers a gateway restart → browser gets mismatch error.

**What was triggering config writes:**
- Stale ghost devices from the old skill-agent architecture connecting and triggering security audit scope-upgrade flows
- `onemind-session-logger` hook writing on every message turn
- Cron jobs running agent tasks that called `openclaw config set`

**The chain:**
```
ghost device connects
  → scope-upgrade security audit
    → config write
      → ensureGatewayStartupAuth runs
        → auth.token in file rotates
          → gateway restarts
            → browser token mismatch
              → pairing loop
```

**The fix:** Remove `auth.token` from the file. No file value = nothing to compare = nothing to rotate.

---

## Config Invariants (verify after any major change)

```bash
ssh onemindos-core 'python3 -c "
import json
cfg = json.load(open(\"/home/ubuntu/.openclaw/openclaw.json\"))
token_in_file = cfg[\"gateway\"][\"auth\"].get(\"token\")
spawn_ctx = cfg[\"channels\"][\"discord\"][\"threadBindings\"][\"defaultSpawnContext\"]
approvals = cfg[\"channels\"][\"telegram\"][\"execApprovals\"][\"target\"]
print(\"auth.token in file:\", token_in_file or \"NOT SET (good)\")
print(\"defaultSpawnContext:\", spawn_ctx, \"(want: isolated)\")
print(\"execApprovals.target:\", approvals, \"(want: dm)\")
"'
```

Expected output:
```
auth.token in file: NOT SET (good)
defaultSpawnContext: isolated (want: isolated)
execApprovals.target: dm (want: dm)
```

---

## If You See "gateway token mismatch"

1. Check if `auth.token` appeared in the file:
```bash
ssh onemindos-core 'python3 -c "import json; cfg=json.load(open(\"/home/ubuntu/.openclaw/openclaw.json\")); print(cfg[\"gateway\"][\"auth\"].get(\"token\", \"NOT SET\"))"'
```

2. If it exists, remove it:
```bash
ssh onemindos-core 'python3 -c "
import json
cfg = json.load(open(\"/home/ubuntu/.openclaw/openclaw.json\"))
cfg[\"gateway\"][\"auth\"].pop(\"token\", None)
with open(\"/home/ubuntu/.openclaw/openclaw.json\", \"w\") as f:
    json.dump(cfg, f, indent=2)
print(\"Removed. Token now env-only.\")
"'
```

3. Verify env var still set:
```bash
ssh onemindos-core 'grep OPENCLAW_GATEWAY_TOKEN ~/.openclaw/gateway.systemd.env'
```

4. Restart:
```bash
ssh onemindos-core 'systemctl --user restart openclaw-gateway'
```

---

## If You See "device pairing required"

1. List pending requests:
```bash
ssh onemindos-core 'openclaw devices list 2>&1'
```

2. Approve the Mac browser request (find the requestId from the list):
```bash
ssh onemindos-core 'openclaw devices approve <requestId>'
```

3. If the gateway identity device (`0b8946a3...`) is missing from paired, re-add it:
```bash
ssh onemindos-core 'python3 -c "
import json, time, secrets
dev = json.load(open(\"/home/ubuntu/.openclaw/identity/device.json\"))
paired = json.load(open(\"/home/ubuntu/.openclaw/devices/paired.json\"))
FULL = [\"operator.admin\",\"operator.read\",\"operator.write\",\"operator.approvals\",\"operator.pairing\"]
now = int(time.time()*1000)
tok = secrets.token_urlsafe(32)
paired[dev[\"deviceId\"]] = {
    \"deviceId\": dev[\"deviceId\"], \"publicKey\": dev[\"publicKeyPem\"],
    \"platform\": \"linux\", \"clientId\": \"cli\", \"clientMode\": \"cli\",
    \"role\": \"operator\", \"roles\": [\"operator\"], \"scopes\": FULL, \"approvedScopes\": FULL,
    \"tokens\": {\"operator\": {\"token\": tok, \"role\": \"operator\", \"scopes\": FULL, \"createdAtMs\": now}},
    \"createdAtMs\": now, \"approvedAtMs\": now
}
json.dump(paired, open(\"/home/ubuntu/.openclaw/devices/paired.json\",\"w\"), indent=2)
auth = json.load(open(\"/home/ubuntu/.openclaw/identity/device-auth.json\"))
auth[\"tokens\"][\"operator\"][\"token\"] = tok
json.dump(auth, open(\"/home/ubuntu/.openclaw/identity/device-auth.json\",\"w\"), indent=2)
print(\"Gateway identity re-paired:\", dev[\"deviceId\"][:20])
"'
```

---

## If You See "device identity required"

The gateway's own `identity/device.json` is pointing to a device that's been removed from `paired.json`.

Fix: delete the identity files so the gateway regenerates fresh ones on next restart (it will use shared-secret loopback auth instead, which bypasses pairing):

```bash
ssh onemindos-core 'rm ~/.openclaw/identity/device.json ~/.openclaw/identity/device-auth.json && systemctl --user restart openclaw-gateway'
```

Then re-add the new identity device to paired.json using the snippet above.

---

## Do Not Run `openclaw doctor --fix`

`doctor --fix` rotates the gateway token, disables skills, and rewrites config. It will break the auth setup above. If config validation fails, fix the specific issue manually rather than running doctor.

---

## Related

- [[project_openclaw_gateway_fixes]] — Memory entry with root cause summary
- `CLAUDE.md` — Gateway Auth Critical Rules section
- `~/.openclaw/gateway.systemd.env` — Live env vars including `OPENCLAW_GATEWAY_TOKEN`
- `~/.openclaw/devices/paired.json` — Live device registry
