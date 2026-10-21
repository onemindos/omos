# Portability: leaving DigitalOcean without losing anything

Goal: be able to run the whole OneMind platform on any Kubernetes cluster
(another cloud, a hybrid cloud + home setup, or a dedicated home server) and be
able to prove it by rebuilding it from Git plus backups.

| Read in this order | What it answers |
|---|---|
| [DO-LOCK-IN.md](DO-LOCK-IN.md) | What exactly ties us to DigitalOcean today, how many of each, and what replaces it |
| [REBUILD-FROM-ZERO.md](REBUILD-FROM-ZERO.md) | Step-by-step disaster recovery onto an empty cluster, what is NOT covered yet, and a quarterly drill |
| [HOME-CLUSTER-PLAN.md](HOME-CLUSTER-PLAN.md) | Which home software and hardware to pick, and a phased, reversible migration |

Related existing documents (inputs, not replaced by these):
[../STACK.md](../STACK.md), [../networking/NETWORKING.md](../networking/NETWORKING.md),
[../../platform/backups/OFFSITE-COPY.md](../../platform/backups/OFFSITE-COPY.md),
[../../platform/backups/tools/README.md](../../platform/backups/tools/README.md),
[../operations/home-migration.md](../operations/home-migration.md),
[../operations/cluster-review-2026-10-09.md](../operations/cluster-review-2026-10-09.md).

Ground rules used throughout: nothing is deleted until the replacement has been
verified; every phase has a way back; no secret values ever appear in these
documents or in Git in plaintext.

Facts in these documents were checked read-only against the live cluster
(`do-nyc3-omos-cloud`) and the repository on 2026-10-10. Counts drift; the
commands to recheck them are included.
