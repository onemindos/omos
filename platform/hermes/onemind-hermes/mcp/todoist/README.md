# Todoist hosted MCP: startup and recovery

Todoist is a major operational integration, not a new service or a declaration
that it replaces ClickUp or owns all task authority. This change migrates no task
data and makes no Todoist API writes.

## Public configuration and scope

`server.json` owns only `mcp_servers.todoist` in the **default** Hermes home
(`/opt/data`). The endpoint is `https://ai.todoist.net/mcp`, `auth: oauth`, enabled.
The working public configuration was inspected on `onemind-hermes-6bf856b5-wlk99`
on 2026-11-04 without printing credentials. Its four existing exclusions are
preserved: `search`, `fetch`, `export-project-template`, `import-project-template`.
There is **no new tool allowlist**; remote sampling is explicitly disabled.
Normal task tools remain available subject to the user's authorization.

The existing `initial-config` init container seeds missing `config.yaml` from
`config.defaults.yaml` at mode 0600, repairs only the home/config ownership before
the Hermes shim drops privileges, then runs Zoho, Whop and Todoist reconcilers
outside the seed guard on **every pod start**, for existing and empty PVCs.
Todoist calls the supported `hermes config set mcp_servers.todoist '<public JSON>'`
CLI, never replaces the whole MCP map, and suppresses CLI output on failure.
The hashed `onemind-hermes-todoist` ConfigMap mounts read-only at `/todoist`;
Kustomize updates its Deployment reference when the manifest or reconciler changes.
No additional pod or scheduler requests are introduced; the existing init limits
remain CPU 500m / memory 192Mi (requests 10m / 32Mi).

Zoho, Whop, ClickUp, Dropbox, gateway MCPs, model/settings and other profiles are
not managed by this reconciler. Existing private OAuth state is not read, copied,
reset, or reauthorized by it. The native Todoist CLI's authentication is separate.

## Recovery boundaries

- **Pod replacement, same PVC:** Git restores the public definition and the private
  authorization remains in `/opt/data/mcp-tokens/todoist*.json` on
  `onemind-hermes-data`. Keep these files private (0600; directory 0700). Never put
  their contents, auth headers or OAuth callbacks into Git, a ConfigMap, a PR,
  transcript or public config export.
- **Empty PVC:** Git can recreate public settings, not the private grant. Restore
  the encrypted PVC backup through the existing recovery procedure, or authorize
  Todoist again on the official default-profile cloud instance with
  `/opt/hermes/bin/hermes mcp login todoist`. Do not start login during unattended
  init; do not copy a desktop/other-profile authorization or overwrite a live PVC.
- **Revoked/expired authorization:** startup reconciliation cannot repair a revoked
  grant. Use the same login command with an owner available for browser approval,
  then test the connection and read `user-info` to verify the intended account.
  Avoid task mutation as an authentication test. Start a fresh session to reload
  tool configuration after an approved rollout.
- **Verification after approved merge/deploy:** check Argo revision and generated
  ConfigMap reference, successful `initial-config` logs, default-home Todoist
  sampling disabled, then `/opt/hermes/bin/hermes mcp test todoist` and a read-only
  `user-info`. Do not dump the complete config or token files to diagnose this.

`backup.yaml` declares encrypted full-PVC VolSync/Restic coverage using
`onemind-hermes-restic`; this includes the credential directory by scope, but
is not proof that a particular grant has been captured or restored. The parent's
read-only observation on 2026-11-04 found last successful sync
**2026-11-04T01:25:36Z**, before Todoist authorization files created/updated
**2026-11-04T01:42:28Z–01:42:52Z**. That snapshot **does not establish backup of this
new authorization**. The next scheduled sync was **2026-11-05T01:25:00Z**.
Acceptance remains pending: verify a successful sync after authorization, and a
separately authorized restore drill plus connection/account verification. No
backup was triggered and no restore, rollout or cluster mutation is claimed here.

## Local contracts

```sh
python -m unittest discover -s platform/hermes/onemind-hermes/mcp/todoist -v
python -m unittest discover -s scripts/tests -p 'test_gitops_todoist_mcp.py' -v
kustomize build platform/hermes/onemind-hermes
make validate
```

Tests require PyYAML; the reconciler itself uses only stdlib. If the actual Hermes
CLI is installed, tests call it twice on an isolated drifted fixture and assert
unrelated settings, synthetic private file bytes/modes/owner/mtime and another
profile are unchanged. Fresh and existing startup code is also executed twice
with rewritten fixture paths, the real Hermes configuration CLI, simulated root/chown ordering,
and unrelated Zoho/Whop subprocesses mocked. This is not a root container or
live-volume recovery test. Without a CLI, its direct test explicitly skips;
public invocation, startup ordering, mount/default/backup contracts still run.
Fixtures stay under `$TMPDIR` / `$RUNNER_TEMP` (fallback home cache); no live
credentials are read and scratch fixtures are retained. The standard
`test_gitops*.py` quality job discovers these contracts without a new workflow.
