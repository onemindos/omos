# Hermes follow-ups — 2026-11-04

Work parked by the owner on 2026-11-04, to continue later. State below was checked live on that date; verify it again
before acting.

## Where things stand

`hermes/onemind-hermes` (persona Legacy) runs Hermes v0.22.7 and passed the live check: model reply, NATS prompt and
heartbeat, `ops` 38 / `data` 14 / `tak` 25 tools, role delegation with tool calls, Telegram, Matrix, API server, dashboard,
Desktop reconnect, text-to-speech. The old `hermes-legacy` and Paperclip are removed.

## Needs the owner

| # | Item | What it unlocks | How |
|---|---|---|---|
| 1 | Keys session (about 15 minutes) | `execute` approvals through `@gate`; Grafana, GitHub, DigitalOcean, Cloudflare, n8n, Home Assistant, CloudTAK and Argo CD tools | `scripts/gateway/keys-session.sh`; order and key paths in `docs/operations/tool-gateway.md`. The execute gate stays at `replicas: 0` until it runs |
| 2 | Nightly config export | Password-free export to `onemindos/onemind-hermes-config` | Enable deploy keys in the org settings (recommended), or use a fine-grained token (`legacy-export` step; the script's `push()` then needs HTTPS). The job stays paused until then |
| 3 | `approvals.mode` | Decides who judges risky commands | Live is `smart` with a 3-entry allowlist; the seed says `manual`. `manual` can't be answered in Telegram, NATS or cron sessions |
| 4 | Speech-to-text | Dictation and wake word | Groq or OpenAI key (recommended), or local `faster-whisper` (about 200 MB plus a model, on a tight node) |
| 5 | Realtime voice | Desktop realtime voice | Needs `OPENAI_API_KEY`; optional |
| 6 | Mac plugins `hermes-project-stewardship`, `home-dashboard` | Those Desktop panels | No server half on the pod. Third-party code: install the server half or remove the plugins |
| 7 | Device logins | Legacy logging in to devices | Mac: SSH key (`legacy-ssh` step; the pod mount is not wired yet). S24: Termux `sshd`, SSH key, Wireless Debugging pairing. Windows laptop: online, OpenSSH Server, SSH key. Home Assistant: online, long-lived token |
| 8 | Old gateway on the Mac (`ai.hermes.gateway`, identity `legacy-edge`) | Less to run | Off unless Legacy needs local Mac hands that SSH can't give |
| 9 | Tailscale console | Tidy device list | Remove `hermes-agent`, `hermes-legacy`, `paperclip` |
| 10 | `gh` plus a GitHub token in the pod | "Create PR" from Desktop | Optional |

## For an agent, no owner input needed

- Read the `platform-argocd` drift (`argocd-cm`, `argocd-server`, `argocd-repo-server`, `argocd-application-controller`)
  before any sync. It is a manual-sync app, last synced 2026-10-25, and a sync may restart Argo. The `legacy` Argo CD key
  depends on it.
- `onemind-mcp`: pin the image built after #641. Harmless until the next change there.
- Legacy's Argo app returns to OutOfSync after each config change, because `prune: false` leaves old hashed ConfigMaps.
  Decide on a Git fix.
- `platform/core/policy.hujson` still names the `hermes-agent` host.
- `nats-platform/adapter.py` imports `_get_approval_config` from the wrong module, so the approval timeout falls back to
  300 s; a background `delegate_task` finishing after a NATS prompt can't be delivered back to the NATS chat.
- `docs/architecture/onemind-tool-gateway.md` still has Paperclip rows.
- Unverified until keys exist: Home Assistant long-lived token at `/api/mcp`; GitHub, DigitalOcean and Cloudflare tool
  lists against the live servers. Karakeep, TAK Portal and TAK Server have no adapter yet.
- `USER.md` is close to its limit again, which will block new memory writes.

## Not started

- TAK bridge: #560, then a certificate to switch it on.
- Revise #585 with the Mumble voice bot (data and chat through the bridge, voice over Mumble in ATAK).
- `hermes peer` for other Hermes instances. Legacy's API server isn't exposed on the tailnet yet.
- #587 waits for a maintenance window from the owner. #562 and #564 stay open; `nats-context` stays off until #562.
