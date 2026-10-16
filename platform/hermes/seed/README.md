# Hermes identity seeding

`seed.sh` fills an **empty** Hermes data volume with the identity files kept in Git, so a rebuild does not start from Hermes' defaults. It is Step 2b of `docs/operations/hermes-recoverability-plan.md`.

## What it does

| From the repo (`platform/hermes/agent/`) | To the volume |
|------------------------------------------|---------------|
| `agents/legacy/SOUL.md` | `/opt/data/SOUL.md` |
| `cron/jobs.json` | `/opt/data/cron/jobs.json` |
| `skills/**` | `/opt/data/skills/**` |
| `config/config.template.yaml` | `/opt/data/config.yaml`, **only if no `${PLACEHOLDERS}` are left** |

- It never overwrites an existing file, so a normal restart changes nothing.
- `HERMES_RESEED=1` is the deliberate reset: it overwrites the seeded files.
- A config template that still has `${PLACEHOLDERS}` (credentials were redacted on export) is not installed; the script says so and leaves it for a person to fill in.
  This is deliberately conservative. Hermes expands `${VAR}` in `config.yaml` from the environment, so once the credentials live in the `hermes-env` Secret (see `docs/operations/hermes-credential-cleanup.md`) the template can be installed as is; relax the check then.
- It never touches memory, the databases, sessions, or anything else on the volume.

## Test

```bash
sh platform/hermes/seed/test-seed.sh
```

Covers: empty volume, existing files kept, placeholder config refused, clean config installed once, reseed, and an empty seed.

## Not wired in yet, on purpose

Today Git holds only the soul. The skills, the config template and the cron jobs arrive with the first successful export (`.github/workflows/hermes-identity-export.yml`), which waits on the credential rotation. Until then, wiring this into the Hermes pod would restart the live agent to seed one file.

To wire it in later:

1. Make the identity files available at `/seed`, laid out as above. The soul and cron fit in a ConfigMap; skills and config do not (1 MiB limit), so use a small image built by CI from `platform/hermes/agent/` and pin it by digest.
2. Add an init container to `platform/hermes/agent/deployment.yaml`, next to `install-gateway-deps`, that mounts `/seed` read-only and the `data` volume at `/opt/data` and runs `sh /seed.sh`. Because it never overwrites, it is safe on every start.
3. Run the restore drill's empty-volume case once and compare with the live agent.
