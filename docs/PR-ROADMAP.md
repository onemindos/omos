# Control Plane PR Roadmap

This is the scoped work plan for delegated coding agents. Keep PRs small. `Mother` owns architecture and review.

## PR 1 — Foundation docs + schemas

Repo: `omos-core`

- control-plane doctrine docs
- registry schema stubs
- UI projection schema stubs
- command/workflow schema stubs
- no runtime behavior yet

## PR 2 — `om-control-plane` MVP

Repo: `omos-core`

- NEX workload skeleton
- NATS monitor poller
- NEX workload snapshot
- om-unity entity count poller
- TAK/Legacy health probe
- publish `ui.now.summary`, `ui.systems.index`, `ui.assets.index`, `ui.feed.recent`

## PR 3 — `om-tak-projector` MVP

Repo: `omos-core`

- generalize current `om-overwatch-tak` pattern
- subscribe `ent.*.*.state`
- honor `visibility.tak`
- support `cmd.entity.{id}.tak.show/hide`
- publish CoT to `tak.cot.inject`

## PR 4 — Flutter data-spine cleanup

Repo: `omos-atoc`

- one NATS transport service
- generated/typed contract imports
- retire page-owned raw WebSocket NATS paths
- keep Wire as inspector-only

## PR 5 — Control-plane UI MVP

Repo: `omos-atoc`

- nav: Now, Assets, Feed, Work, Connect, Systems, Insights, Protect
- consume `ui.now.summary`, `ui.assets.index`, `ui.systems.index`, `ui.feed.recent`
- entity detail drawer
- Show in TAK / Hide in TAK buttons publishing commands

## PR 6 — Connector manager + Home Assistant

Repos: `omos-core`, `omos-atoc`

- integration registry
- Home Assistant connector page
- credential/test/discover/start/stop flow
- entity discovery into registry

## Review rules

- No PR may subscribe normal UI to raw `evt.>`, `tel.>`, `tak.cot.>`, `$SYS.>` firehose subjects.
- No PR may bypass command policy for effectors.
- No PR may make Flutter source-of-truth for registry state.
- No PR may merge if it breaks current om-overwatch -> TAK loop.
- Reference/ignore folders stay out of git.
