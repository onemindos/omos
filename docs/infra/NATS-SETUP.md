# NATS-SETUP.md — OneMind NATS Server Foundation (P1-A)

> Phase 1 deliverable. Implemented: 2026-06-10. Updated: 2026-06-10 (P1-A card t_e473f24d).
> **Passwords are NOT in this file.** See `~/.nats/config/accounts.creds` (chmod 600, never commit).

---

## 1. Server

- Binary: `/opt/homebrew/bin/nats-server` (v2.14.1)
- Config: `~/.nats/config/onemind.conf`
- Data dir: `~/.nats/data`
- Log: `~/.nats/nats-server.log`
- Managed by: LaunchAgent `io.nats.server` (`~/Library/LaunchAgents/io.nats.server.plist`)
- Server name: `onemind-mac`
- JetStream domain: `onemind`

### Restart / reload

```sh
# Reload config without restart (for most changes):
launchctl kill SIGHUP "gui/$(id -u)/io.nats.server"

# Full restart (new config with account changes):
launchctl kickstart -k "gui/$(id -u)/io.nats.server"

# Verify PID after restart:
ps aux | grep nats-server | grep -v grep
```

---

## 2. Account Layout (Doctrine §5 config-based)

Four accounts. No JWTs. Config-based staging — designed to migrate to operator/JWT when
the fleet exceeds ~20 nodes (see NATS-DOCTRINE.md §5).

| Account   | Purpose                                      | JetStream | User       |
|-----------|----------------------------------------------|-----------|------------|
| `SYS`     | System account — server $SYS.> events        | no        | n/a        |
| `CORE`    | Agents + meta services (Legacy, workers)     | yes       | `core`     |
| `FABRIC`  | Telemetry + IoT (sensors, bridges, ESP32)    | yes       | `fabric`   |
| `INTERFACE` | Apps (Flutter, web cockpit)                | yes       | `interface`|
| `GLOBAL`  | **DEPRECATED** — backward-compat unauthenticated clients | yes | `onemind` |

`no_auth_user: onemind` is kept so all existing clients (OpenClaw subscribers, LaunchAgents)
continue to work against GLOBAL until migrated. Remove this once all clients use named accounts.

### Cross-account exports / imports (minimal, exactly as Doctrine specifies)

```
CORE    exports: evt.>  (stream),  svc.>  (service)
FABRIC  exports: tel.>  (stream)

CORE    imports: tel.>  from FABRIC
INTERFACE imports: evt.>  from CORE (stream)
              svc.>  from CORE (service)
              tel.>  from FABRIC (stream)
```

### Credentials

Stored in `~/.nats/config/accounts.creds` (chmod 600). **Never commit this file.**

Format:
```
core_password: <strong random hex>
fabric_password: <strong random hex>
interface_password: <strong random hex>
```

---

## 3. Server Config (passwords redacted)

```hcl
# ~/.nats/config/onemind.conf (passwords redacted for docs — real passwords in accounts.creds)
server_name: onemind-mac
port: 4222
http_port: 8222
host: 0.0.0.0

system_account: SYS

accounts {
  SYS {}

  CORE {
    jetstream: enabled
    users: [ { user: core, password: "<see accounts.creds>" } ]
    exports: [
      { stream: "evt.>" }
      { service: "svc.>" }
    ]
    imports: [
      { stream: { account: FABRIC, subject: "tel.>" } }
    ]
  }

  FABRIC {
    jetstream: enabled
    users: [ { user: fabric, password: "<see accounts.creds>" } ]
    exports: [ { stream: "tel.>" } ]
  }

  INTERFACE {
    jetstream: enabled
    users: [
      {
        user: interface, password: "<see accounts.creds>",
        permissions: {
          publish:   { deny: ["tel.>", "cmd.>"] }
          subscribe: { allow: ["evt.>", "tel.>", "svc.>", "_INBOX.>"] }
        }
      }
    ]
    imports: [
      { stream: { account: CORE,   subject: "evt.>" } }
      { service: { account: CORE,  subject: "svc.>" } }
      { stream: { account: FABRIC, subject: "tel.>" } }
    ]
  }

  # DEPRECATED: migrate clients to accounts, then remove
  GLOBAL {
    jetstream: enabled
    users: [ { user: onemind } ]
  }
}

no_auth_user: onemind

jetstream {
  store_dir: "/Users/onemind/.nats/data"
  max_mem: 1G
  max_file: 10G
  domain: onemind
}

leafnodes  { port: 7422 }
websocket  { port: 8080; no_tls: true; same_origin: false }
mqtt       { port: 1883 }

debug: false
trace: false
logtime: true
log_file: "/Users/onemind/.nats/nats-server.log"
```

---

## 4. Stream Setup

EVT and CMD are created in the **CORE** account; TEL in the **FABRIC** account.

### EVT — Events (facts, past tense)

```sh
nats --user core --password '<see accounts.creds>' stream add EVT \
  --subjects "evt.>" \
  --storage file \
  --retention limits \
  --max-age 30d \
  --replicas 1 \
  --discard old \
  --dupe-window 2m \
  --defaults
```

Result: `Stream EVT was created` — subjects: `evt.>`, retention: Limits, max_age: 30d, storage: File.

### TEL — Telemetry

```sh
nats --user fabric --password '<see accounts.creds>' stream add TEL \
  --subjects "tel.>" \
  --storage file \
  --retention limits \
  --max-age 7d \
  --replicas 1 \
  --discard old \
  --dupe-window 2m \
  --defaults
```

Result: `Stream TEL was created` — subjects: `tel.>`, retention: Limits, max_age: 7d, storage: File.

### CMD — Commands (work queue)

> **Note:** The pre-existing `COMMANDS` stream had subject `cmd.mode.>` which overlaps with `cmd.>`.
> It had 0 messages and was moved to `_legacy.cmd.mode.>` before creating CMD:
>
> ```sh
> nats stream edit COMMANDS --subjects "_legacy.cmd.mode.>" --force
> ```

```sh
nats --user core --password '<see accounts.creds>' stream add CMD \
  --subjects "cmd.>" \
  --storage file \
  --retention workq \
  --max-age 24h \
  --replicas 1 \
  --discard old \
  --dupe-window 2m \
  --defaults
```

Result: `Stream CMD was created` — subjects: `cmd.>`, retention: WorkQueue, max_age: 24h, storage: File.

---

## 5. KV Buckets

```sh
nats --user core --password '<see accounts.creds>' kv add om-config
nats --user core --password '<see accounts.creds>' kv add om-schema
```

Both created with default settings (file storage, history=1, unlimited size/age).
`om-config`: runtime configuration for OneMind services.
`om-schema`: JSON Schema registry (see NATS-DOCTRINE.md §7 "Schema Registry").

---

## 6. Verification Transcript (2026-06-10)

### All three accounts authenticate

```
$ nats account info --user core --password '<redacted>'
  User: core  Account: CORE (CORE)  Connected Server: onemind-mac v2.14.1

$ nats account info --user fabric --password '<redacted>'
  User: fabric  Account: FABRIC (FABRIC)  Connected Server: onemind-mac v2.14.1

$ nats account info --user interface --password '<redacted>'
  User: interface  Account: INTERFACE (INTERFACE)  Connected Server: onemind-mac v2.14.1
```

### Streams EVT, TEL, CMD present

```
$ nats stream ls  (filtered to new streams)
  EVT  — 2026-06-10 13:04:29  subjects: evt.>   max_age: 30d
  TEL  — 2026-06-10 13:04:35  subjects: tel.>   max_age: 7d
  CMD  — 2026-06-10 13:05:32  subjects: cmd.>   max_age: 24h  retention: WorkQueue
```

### KV buckets om-config, om-schema present

```
$ nats kv ls  (filtered to new buckets)
  om-config   — 2026-06-10 13:05:36
  om-schema   — 2026-06-10 13:05:36
```

### EVT stream captures messages

```
$ nats pub evt.sys.test.fired '{"t":1}'
  13:05:53 Published 7 bytes to "evt.sys.test.fired"

$ nats stream get EVT 1
  Item: EVT#1 received 2026-06-10 17:05:53 on Subject evt.sys.test.fired
  {"t":1}
```

### Unauthenticated GLOBAL compat

```
$ nats pub devices.test.ping '{}'
  13:06:02 Published 2 bytes to "devices.test.ping"   ✓
```

### WebSocket still up

```
$ curl -s -o /dev/null -w '%{http_code}' http://127.0.0.1:8080
  400   (HTTP 400 = WS upgrade required; port is alive)  ✓
```

### MQTT port still listening

```
$ lsof -iTCP:1883 -sTCP:LISTEN
  nats-serv  72088  onemind  10u  IPv6  TCP  *:ibm-mqisdp (LISTEN)  ✓
```

### Account isolation verified

```
# interface CANNOT publish to tel.> — blocked by explicit permissions deny
$ nats --user interface --password '<redacted>' pub tel.barn.temp.c '{"val":22.5}'
nats: error: nats: permissions violation: Permissions Violation for Publish to "tel.barn.temp.c"  ✓

# interface CANNOT subscribe to arbitrary FABRIC subjects (not imported)
$ nats --user interface --password '<redacted>' subscribe "FABRIC.>" --count=0 --timeout=1s
nats: error: nats: permissions violation: Permissions Violation for Subscription to "FABRIC.>"  ✓

# interface CAN subscribe to evt.> (imported from CORE) — verified via accountz API:
# curl http://127.0.0.1:8222/accountz?acc=INTERFACE
# shows imports: evt.> from CORE (stream), tel.> from FABRIC (stream), svc.> from CORE (service)  ✓
```

---

## 7. Rollback Note

If the new config causes problems, restore and restart:

```sh
# Restore backup (created automatically before editing)
cp ~/.nats/config/onemind.conf.bak-20260610 ~/.nats/config/onemind.conf

# Validate the restored config
/opt/homebrew/bin/nats-server -c ~/.nats/config/onemind.conf -t

# Restart the server
launchctl kickstart -k "gui/$(id -u)/io.nats.server"

# Verify server is up
ps aux | grep nats-server | grep -v grep
nats account info
```

The backup at `~/.nats/config/onemind.conf.bak-20260610` contains the original flat/GLOBAL-only
config with no named accounts. All existing streams and KV data are in `~/.nats/data` and are NOT
touched by any config change.

**Important:** The three new streams (EVT, TEL, CMD) and two new KV buckets (om-config, om-schema)
were created against the GLOBAL account and persist in `~/.nats/data` regardless of which config
version is active. Rollback does not delete them.

---

## 8. Migration Path (GLOBAL → accounts)

When migrating existing services off GLOBAL:

1. Update each service's NATS connection URL to include `user` + `password` from `accounts.creds`.
2. Test the service connects to the correct account.
3. Once all services using a given subject namespace are migrated, they'll be isolated.
4. Remove `GLOBAL` account and `no_auth_user: onemind` from `onemind.conf` when no clients remain.
5. Future: operator/JWT mode when fleet exceeds ~20 nodes (see NATS-DOCTRINE.md §5).
