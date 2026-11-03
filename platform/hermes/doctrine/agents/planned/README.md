# Planned agents

Forge and the nine servants below are **designed but not running**. Their profile folders on the Hermes volume are empty (Forge included) and their cron jobs were paused. They are kept here because the four overseer souls name them.

## Hierarchy

```
Zeus -> Legacy -> overseers -> servants

SO  (technical)  -> forge, grid, guardian, oracle, trinity, spartan
HP  (personal)   -> spartan
LE  (home)       -> eden, haven, heritage
GE  (business)   -> empire
```

Legacy and the four overseers (so, hp, le, ge) live one level up in `../`. The registry with a status for every agent is `../../AGENTS.yaml`.

## Before activating one

- **Model fields are stale.** Several `CONFIG.md` files and registry entries name Bedrock models; the platform now uses DigitalOcean inference. They are marked `model_status: stale` in the registry.
- Create the profile on the volume, confirm its NATS session registers on the cluster bus, then move its entry from `planned:` to `agents:` in the registry and its folder out of `planned/`.
- Update `CRON-STATE`-style assumptions: all servant cron jobs were paused so they would not spend inference credits.
