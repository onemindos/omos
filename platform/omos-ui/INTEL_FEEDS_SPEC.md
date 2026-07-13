# Intel Feeds Panel — Implementation Spec

**Component:** `IntelFeedsPanel`
**Route / tab:** `/intel` (or Intel tab in the main nav)
**NATS dependencies:** KV bucket `feed-config`, subjects `tel.feed.*.status`, `intel.>`

---

## 1. Layout

```
┌─────────────────────────────────────────────────────────────────────┐
│  INTEL FEEDS                                      [Refresh] status  │
├─────────────────────────────────────────────────────────────────────┤
│                                                                     │
│  ┌──────────────────────────────┐  ┌──────────────────────────────┐ │
│  │ ⛈  NOAA Weather Alerts      │  │ 🌍  USGS Seismic             │ │
│  │ ● Active  last: 2m ago       │  │ ● Active  last: 45s ago      │ │
│  │ 12 msgs / last poll          │  │ 7 msgs / last poll           │ │
│  │ 847 total                    │  │ 2,341 total                  │ │
│  │                      [■ ON]  │  │                      [■ ON]  │ │
│  │ ▸ Settings                   │  │ ▸ Settings                   │ │
│  └──────────────────────────────┘  └──────────────────────────────┘ │
│                                                                     │
│  ┌──────────────────────────────┐  ┌──────────────────────────────┐ │
│  │ ⚔  ACLED Conflict           │  │ 🔥  NASA FIRMS Fire          │ │
│  │                              │  │                              │ │
│  │ ┌──────────────────────────┐ │  │ ┌──────────────────────────┐ │ │
│  │ │ ⚠ Configure API key to  │ │  │ │ ⚠ Configure API key to  │ │ │
│  │ │   activate this feed     │ │  │ │   activate this feed     │ │ │
│  │ └──────────────────────────┘ │  │ └──────────────────────────┘ │ │
│  │                      [□ OFF] │  │                      [□ OFF] │ │
│  │ ▸ Settings                   │  │ ▸ Settings                   │ │
│  └──────────────────────────────┘  └──────────────────────────────┘ │
│                                                                     │
├─────────────────────────────────────────────────────────────────────┤
│  STREAM STATS — intel.>                     47.2 msg/s  total: 89k  │
│  noaa     ████████████  31.4/s   12,847 msgs                        │
│  usgs     ███           8.1/s    24,103 msgs                        │
│  acled    —             0/s      0 msgs   (no API key)              │
│  firms    —             0/s      0 msgs   (no API key)              │
└─────────────────────────────────────────────────────────────────────┘
```

**Expanded settings — NOAA card:**
```
  ▾ Settings
  ┌────────────────────────────────────────┐
  │ Poll interval (s)  [  300  ]           │
  │ Areas              [  US   ]           │
  │                             [Save]     │
  └────────────────────────────────────────┘
```

**Expanded settings — ACLED card (with API key prompt):**
```
  ▾ Settings
  ┌────────────────────────────────────────┐
  │ Poll interval (s)  [ 3600  ]           │
  │ API key            [••••••••••] [Show] │
  │ Email              [user@example.com ] │
  │                             [Save]     │
  └────────────────────────────────────────┘
```

---

## 2. Feed Definitions

Declare this as a static config object — it drives card rendering and tells the component which KV keys to read/write.

```ts
interface FeedDef {
  id: string;           // "noaa" | "usgs" | "acled" | "firms"
  label: string;
  icon: string;         // emoji or SVG import
  statusSubject: string;// "tel.feed.noaa.status"
  kvPrefix: string;     // "feed.noaa"
  gated: boolean;       // requires API key to be useful
  settings: SettingDef[];
}

interface SettingDef {
  key: string;          // KV key suffix — e.g. "poll_interval"
  label: string;
  type: "text" | "number" | "password";
  placeholder?: string;
}

const FEEDS: FeedDef[] = [
  {
    id: "noaa", label: "NOAA Weather Alerts", icon: "⛈",
    statusSubject: "tel.feed.noaa.status",
    kvPrefix: "feed.noaa", gated: false,
    settings: [
      { key: "poll_interval", label: "Poll interval (s)", type: "number", placeholder: "300" },
      { key: "areas",         label: "Areas",             type: "text",   placeholder: "US" },
    ],
  },
  {
    id: "usgs", label: "USGS Seismic", icon: "🌍",
    statusSubject: "tel.feed.usgs.status",
    kvPrefix: "feed.usgs", gated: false,
    settings: [
      { key: "poll_interval",  label: "Poll interval (s)",  type: "number", placeholder: "60"  },
      { key: "min_magnitude",  label: "Min magnitude",      type: "number", placeholder: "2.5" },
    ],
  },
  {
    id: "acled", label: "ACLED Conflict", icon: "⚔",
    statusSubject: "tel.feed.acled.status",
    kvPrefix: "feed.acled", gated: true,
    settings: [
      { key: "poll_interval", label: "Poll interval (s)", type: "number",   placeholder: "3600" },
      { key: "api_key",       label: "API key",           type: "password", placeholder: ""     },
      { key: "email",         label: "Email",             type: "text",     placeholder: ""     },
    ],
  },
  {
    id: "firms", label: "NASA FIRMS Fire", icon: "🔥",
    statusSubject: "tel.feed.firms.status",
    kvPrefix: "feed.firms", gated: true,
    settings: [
      { key: "poll_interval", label: "Poll interval (s)", type: "number",   placeholder: "600"   },
      { key: "api_key",       label: "API key",           type: "password", placeholder: ""      },
      { key: "area",          label: "Area",              type: "text",     placeholder: "world" },
      { key: "days",          label: "Days",              type: "number",   placeholder: "1"     },
    ],
  },
];
```

---

## 3. NATS.ws Integration Patterns

### 3.1 Initialization

The panel receives an already-connected `NatsConnection` (`nc`) and `JetStreamClient` (`js`) as props or from a shared context. It does not manage the NATS connection itself.

```ts
// props shape
interface IntelFeedsPanelProps {
  nc: NatsConnection;  // from nats.ws connect()
  js: JetStreamClient; // nc.jetstream()
}
```

### 3.2 Reading all feed config on mount

Use `kv.history("feed.>")` to get all current values in one pass, then apply them to local state. This is preferable to calling `kv.get()` per key — it's one round trip and returns all keys matching the prefix.

```ts
const kv = await js.views.kv("feed-config");

// Snapshot — get all current values at mount
const iter = await kv.history({ keys: "feed.>" });
const config: Record<string, string> = {};
for await (const entry of iter) {
  if (entry.operation !== "DEL" && entry.operation !== "PURGE") {
    config[entry.key] = entry.string();
  }
}
// config["feed.noaa.enabled"] === "true", etc.
```

### 3.3 Watching KV for live config changes

After the history pass completes, replace it with a watcher to receive live updates. The watcher picks up from the current sequence, so there is no gap.

```ts
const watcher = await kv.watch({ keys: "feed.>" });
(async () => {
  for await (const entry of watcher) {
    if (entry.operation === "PUT") {
      dispatch({ type: "CONFIG_UPDATE", key: entry.key, value: entry.string() });
    } else if (entry.operation === "DEL" || entry.operation === "PURGE") {
      dispatch({ type: "CONFIG_DELETE", key: entry.key });
    }
  }
})();

// Cleanup on unmount:
watcher.stop();
```

One watcher on `feed.>` covers all four feeds. Handle the lifecycle carefully: call `watcher.stop()` in the component's cleanup effect.

### 3.4 Writing config (toggle + save)

All writes use `kv.put(key, value)`. Values are always strings — coerce numbers and booleans before writing.

```ts
// Toggle enabled
await kv.put(`feed.${id}.enabled`, enabled ? "true" : "false");

// Save a settings field
await kv.put(`${feed.kvPrefix}.${setting.key}`, inputValue.trim());
```

`kv.put` is fire-and-forget for the UI — the watcher will echo the change back and reconcile state. Do not optimistically update state before the echo; let the watch path be the single source of truth.

**Error handling:** wrap in try/catch and surface a toast. A failed put means the KV bucket is unreachable, which the operator needs to know.

### 3.5 Subscribing to feed status

Each worker publishes `tel.feed.{id}.status` every 30 seconds. Subscribe with a wildcard to cover all feeds in one subscription.

```ts
const statusSub = nc.subscribe("tel.feed.*.status");
(async () => {
  for await (const msg of statusSub) {
    const status: FeedStatus = JSON.parse(msg.string());
    // msg.subject = "tel.feed.noaa.status" → extract id from tokens
    const id = msg.subject.split(".")[2]; // "noaa"
    dispatch({ type: "STATUS_UPDATE", id, status });
  }
})();

// Cleanup:
statusSub.unsubscribe();
```

The status payload shape:
```ts
interface FeedStatus {
  feed_id:        string;
  enabled:        boolean;
  poll_interval:  number;
  last_poll_ts:   string;  // ISO-8601
  last_poll_msgs: number;
  last_error:     string;
  total_msgs:     number;
}
```

Workers that have not yet published (first 30s after boot, or workers not running) will have no status entry. Model this as `status: undefined` in component state — show "No status yet" rather than an error.

Status messages are core NATS (not JetStream), so they are ephemeral — last 30s heartbeat only. If no message arrives within 90 seconds of the last one, mark the feed as "stale".

---

## 4. Per-Feed Card Spec

### 4.1 Status indicator

```
enabled=true  + last_error=""      + last_poll within 90s  → green dot  "Active"
enabled=true  + last_error!=""                              → yellow dot "Error"
enabled=true  + no status received yet                      → grey dot   "No data"
enabled=true  + last_poll older than 90s                    → yellow dot "Stale"
enabled=false                                               → grey dot   "Disabled"
```

Implement as a pure function `feedDotColor(config, status): "green" | "yellow" | "grey"`.

### 4.2 Last poll display

Parse `last_poll_ts` with `Date.parse()` and show relative time ("2m ago", "45s ago"). Refresh the display every 10 seconds with a `setInterval` — do not re-subscribe to NATS for this.

### 4.3 Toggle switch

```
ON  → kv.put("feed.noaa.enabled", "true")
OFF → kv.put("feed.noaa.enabled", "false")
```

The switch renders its current value from KV config state (the `enabled` key), not from the status heartbeat. The two are different: KV is the operator's intent; the status heartbeat is the worker's current reality. Both can temporarily disagree (e.g. worker not yet restarted after config change).

Disable the toggle while a KV write is in flight (track with a local `saving` flag per feed).

### 4.4 Expandable settings

Use a `<details>`/`<summary>` or a controlled `open` boolean per card. On expand, populate inputs from the current KV config state — they should reflect what is actually in KV, not ephemeral local edits from a prior session.

Each field in `settings` maps to `${kvPrefix}.${key}` in KV. If the KV key does not exist, show the placeholder as a hint.

**Save button:** writes all dirty fields in the settings form. Only write keys that differ from the current KV values. After all puts complete, clear the dirty state.

### 4.5 API key field (ACLED, FIRMS)

- Input type `password` — browser masks by default.
- Add a toggle button to temporarily reveal (`type="text"`).
- If the KV key exists, pre-fill with the actual stored value fetched from KV (it is plaintext in KV — see §5).
- After save, revert display to masked.
- The "Configure API key to activate" banner (§6) disappears as soon as `feed.acled.api_key` appears in KV.

---

## 5. Security Note

API keys stored in KV bucket `feed-config` are plaintext. This is acceptable for this deployment model: the NATS cluster is inside DOKS, access is gated by NATS user authentication, and the Cloudflare tunnel endpoint (`nats-ws.onemindos.dev`) is additionally protected by Tailscale. omos-ui itself requires NATS credentials to connect.

UI obligations:
- Always render `api_key` fields as `type="password"` (masked) by default.
- Only reveal the key while the user has explicitly clicked "Show".
- Re-mask on blur or after save.
- Never log the key value to the console.
- Do not transmit the key anywhere other than the NATS KV put.

No additional encryption is required. If the threat model changes (multi-tenant, public cluster), the correct fix is NATS KV permissions scoping (read/write ACL per subject prefix), not UI-layer encryption.

---

## 6. No-API-Key State

For gated feeds (ACLED, FIRMS), check whether the `api_key` KV entry exists and is non-empty:

```ts
const hasApiKey = (config: Record<string, string>, id: string): boolean => {
  const val = config[`feed.${id}.api_key`];
  return typeof val === "string" && val.trim().length > 0;
};
```

If `hasApiKey` is false:
- Show the warning banner: "Configure API key to activate this feed."
- The toggle switch is rendered but disabled (clicking it writes the `enabled` flag, but the worker will sleep until the key is set — the banner prevents confusion).
- The Settings section is auto-expanded to guide the operator to the API key field.

Once the key is saved and the watcher echoes it back, the banner disappears automatically (no refresh needed).

---

## 7. Stream Stats Panel

Subscribe to `intel.>` on core NATS. This gives live delivery of new messages as they are published by the workers. **Do not use JetStream for this panel** — the consumer would need to replay history, which is not wanted for a rate display. Core subscription only sees new messages from the moment of subscribe, which is exactly right for a rate counter.

```ts
const intelSub = nc.subscribe("intel.>");
(async () => {
  for await (const msg of intelSub) {
    // Subject pattern: intel.{domain}.{source}.{type}...
    // e.g. intel.weather.noaa.alert.tx.tornado-warning
    //      intel.seismic.usgs.m3plus.us2024abcd
    //      intel.conflict.acled.sdn.battles
    //      intel.fire.nasa.firms.world
    const tokens = msg.subject.split(".");
    // Map domain token to feed id
    const feedId = domainToFeedId(tokens[1], tokens[2]);
    dispatch({ type: "INTEL_MSG", feedId, bytes: msg.data.length });
  }
})();
```

Domain-to-feed mapping:
```ts
function domainToFeedId(domain: string, source: string): string {
  if (domain === "weather" && source === "noaa") return "noaa";
  if (domain === "seismic" && source === "usgs")  return "usgs";
  if (domain === "conflict" && source === "acled") return "acled";
  if (domain === "fire")                           return "firms";
  return "unknown";
}
```

**Rate calculation:** maintain a sliding 10-second window. Keep a circular buffer of `{ ts: number, feedId: string }` entries. On each render tick (every second via `setInterval`), evict entries older than 10 seconds and compute `count / 10` for the per-feed rate.

**Counters to display per feed:**
- `msg/s` — rate over the last 10 seconds
- `total` — cumulative count since panel mount (reset on remount)

**Panel-level totals:** sum of all feeds.

---

## 8. State Shape

```ts
interface IntelFeedsState {
  // KV config — all feed.* keys, string values
  config: Record<string, string>;
  configLoaded: boolean;

  // Status heartbeats — keyed by feed id
  status: Record<string, FeedStatus | undefined>;

  // Last time a status message arrived per feed (for stale detection)
  statusReceivedAt: Record<string, number>;

  // Stream stats
  intelCounts: Record<string, number>;  // total since mount
  intelRate:   Record<string, number>;  // msg/s, last 10s window

  // Per-feed UI state
  expanded:  Record<string, boolean>;
  saving:    Record<string, boolean>;
  error:     Record<string, string | null>;
}
```

---

## 9. Cleanup

All NATS subscriptions and watchers must be stopped in the component's cleanup effect:

```ts
useEffect(() => {
  // ... setup watcher, status sub, intel sub ...
  return () => {
    watcher.stop();
    statusSub.unsubscribe();
    intelSub.unsubscribe();
  };
}, [nc, js]);
```

The KV and subscription setup is async. Use a `mounted` ref to guard against state updates after unmount if the async init races with unmount.

---

## 10. Edge Cases

| Scenario | Behavior |
|---|---|
| NATS disconnected | KV and subs will error/stop. Show a panel-level "NATS disconnected" banner using the connection status from context. Do not retry internally — the global connection handler owns reconnect. |
| `feed-config` KV bucket does not exist | `js.views.kv("feed-config")` throws. Catch and show "KV bucket not provisioned — run nats-migration/provision.sh". |
| Worker not running / never published status | `status[id]` is `undefined`. Show grey dot, "No status yet". |
| KV value for a number field is empty string | Treat as unset; fall through to placeholder. |
| Concurrent writes (two operators) | Last write wins — NATS KV is optimistic. Acceptable for single-operator context. |
| Very high intel message rate (FIRMS world, 1000s/s) | The `intelSub` dispatch path must be cheap. Batch into a local array and flush to state every 1s via `setInterval` rather than calling `dispatch` on every message. |
