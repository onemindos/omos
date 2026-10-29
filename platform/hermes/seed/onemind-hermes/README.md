# onemind-hermes identity seed

One agent, one soul: **Legacy**. These files seed the new Hermes instance `onemind-hermes`. They are a review draft; nothing here is wired into a deployment.

| File | Goes to | What it is |
|------|---------|-----------|
| `SOUL.md` | `/opt/data/SOUL.md` | The single persona. GE, HP, LE and SO are folded in as areas of responsibility, not agents. |
| `MEMORY.md` | `/opt/data/memories/MEMORY.md` | Environment and project facts, cleaned. |
| `USER.md` | `/opt/data/memories/USER.md` | Facts about Zeus and his working preferences, cleaned. |

Entries in the memory files are separated by a line containing `§`, the Hermes memory format.

## Sources merged

- Git: `platform/hermes/agent/agents/{legacy,legacy-edge,ge,hp,le,so,_template,planned}`, `doctrine/{CREED,PILLARS,PROFILES,EVOLUTION}.md`, `AGENTS.yaml`, `README.md`, `docs/MISSION.md`, `CLAUDE.md`.
- Live pod `hermes-legacy` (read-only `cat`, 2026-10-29): `/opt/data/SOUL.md`, the four profile souls, `memories/MEMORY.md` and `USER.md`.

## Drift: live versus Git

| Item | Result |
|------|--------|
| `/opt/data/SOUL.md` vs `agents/legacy/SOUL.md` | Identical. |
| Profile souls `generational-empire`, `holistic-performance`, `legacy-evolution`, `sovereign-operations` vs `agents/{ge,hp,le,so}/SOUL.md` | Identical (all four). |
| `memories/MEMORY.md` | Not in Git. Last edited 2026-10-16. A `MEMORY.md.pre-edit-20261016` backup sits beside it. Contains plaintext credentials and stale platform facts (see below). |
| `memories/USER.md` | Not in Git. Last edited 2026-10-02. Mixes user preferences with platform notes. |
| Live memory says "all 10 agents on DOKS pod profiles"; Git registry says the overseers are `configured` and servants `planned` | Stale for a single-agent design; dropped. |
| Source souls vs the platform | The Git souls themselves are stale in places: `legacy-edge` and `_template` say AWS Bedrock; the edge soul lists LAN and Tailscale addresses; the overseer oaths and Creed assume a mother/servant hierarchy; PILLARS numeric domain ranges are flagged by that file as unconfirmed. |

## What was dropped, and why

- **Credentials:** a TAK enrolment password and a WordPress application password were in the live `MEMORY.md`. Neither is reproduced. Both should be treated as exposed and rotated (see below).
- **Stale TAK facts:** public TAK address and port, the "CloudTAK on atoc" host, and "ALL voice/video in takserver ns". TAK is tailnet-only since 2026-10-23, and MediaMTX is in `onetak` per `docs/STACK.md`. Replaced with the current facts, marked to verify.
- **WordPress droplet IP and the `zeusdelacruz.com` Astro/Pages note:** the droplet entry carried a credential; `CLAUDE.md` records `zeusdelacruz.com` ingress as retired 2026-10-10. Dropped rather than guessed; Zeus to confirm if the brand site still matters.
- **Mac Tailscale IP and SSH key name:** internal addressing, not durable. Kept only the in-cluster SSH service name.
- **Multi-agent hierarchy:** oaths of allegiance, "Mother", "servants", escalation channels between overseers, per-overseer bus subjects, and the numeric domain routing rules. A short note that more agents may be added later is kept.
- **Legacy-Edge persona:** reduced to one paragraph. The edge is described as an extension that defers to Legacy, not a second soul.
- **Stack specifics from the old soul and edge soul** (NATS ports and accounts, ClickHouse port, LAN addresses, `ai.hermes.gateway` details): these change, and `docs/STACK.md` holds them.
- **"Omnipresent / sees everything / agent-of-agents / mother" framing:** replaced with plain statements of role.
- **Inference model names and Bedrock claims:** the provider is a runtime-config matter.
- **EVOLUTION roadmap checklists and PR templates:** only the loop and the self-improvement rule remain.
- **Housekeeping memory:** `MEMORY.md.lock` / `USER.md.lock` files and the `.pre-edit` backup are not carried over.

## Items for the owner to double-check

1. Rotate the TAK enrolment password and the WordPress application password found in the old live `MEMORY.md`. The same pod memory would be copied by any whole-volume clone, so do not clone `/opt/data/memories/`.
2. Telegram is mentioned as retained (per `CLAUDE.md`); Matrix is stated as the day-to-day channel (per old `USER.md`). Confirm which is primary for the new instance.
3. The Mercury "close 5 LLCs, open OneMind LLC" line is kept as unverified plan text. Drop it if it is done or private.
4. Which geo-stack, CloudTAK and OneTAK hostnames are still live was not verified here.
5. The soul never deletes and never acts in the physical world (drones, home devices, locks, vehicles, payments) unless Zeus writes "execute" or "override" in chat for that specific action; approval does not carry to the next action, and without the word she proposes and waits. There is no standing exception (the earlier namespace-deletion exception was removed).
6. Soul says the revenue target is $50K per month and education is the cornerstone, from old memory. Confirm still current.

## How these are copied into the pod (after review)

The existing `platform/hermes/seed/seed.sh` copies `agents/legacy/SOUL.md` and never overwrites existing files. For `onemind-hermes`, once the pod and its empty `/opt/data` volume exist:

```bash
POD=$(kubectl --context do-nyc3-omos-cloud -n hermes get pod -l app=onemind-hermes -o name | head -1)  # adjust label
kubectl --context do-nyc3-omos-cloud -n hermes cp platform/hermes/seed/onemind-hermes/SOUL.md  ${POD#pod/}:/opt/data/SOUL.md -c hermes
kubectl --context do-nyc3-omos-cloud -n hermes exec ${POD} -c hermes -- mkdir -p /opt/data/memories
kubectl --context do-nyc3-omos-cloud -n hermes cp platform/hermes/seed/onemind-hermes/MEMORY.md ${POD#pod/}:/opt/data/memories/MEMORY.md -c hermes
kubectl --context do-nyc3-omos-cloud -n hermes cp platform/hermes/seed/onemind-hermes/USER.md   ${POD#pod/}:/opt/data/memories/USER.md -c hermes
```

Preferred long-term path: ship them through the seed init container described in `platform/hermes/seed/README.md` (soul and memories in the seed layout), so a rebuild reproduces them. Copy only after Zeus reviews; the README itself is not copied. Do not run the commands above against `hermes-legacy`.
