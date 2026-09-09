# Hub Legacy (mother) image = Hermes base + NATS platform plugin + Legacy persona.
#
# Bakes in everything the hub needs to BE Legacy on the fabric:
#   - nats-py + synadia-ai libs (not in Hermes's own deps)
#   - the nats-platform plugin (gateway adapter that registers AgentService
#     at agents.prompt.hermes.zeus.legacy and routes prompts through Hermes's
#     full pipeline — persona, tools, skills, session memory)
#   - SOUL.md (Legacy's persona/identity)
#   - config.yaml (model + NATS platform enabled, session=legacy)
#   - a cont-init script that seeds these into the PVC before the gateway starts
#
# `gateway run --replace` (via s6 supervisor) loads the plugin, reads SOUL.md,
# and registers properly. The hub is the MOTHER; the Mac is a daughter.
#
# Build (linux/amd64 for DOKS):
#   docker buildx build --platform linux/amd64 \
#     -t ghcr.io/onemindos/hermes-legacy:v2 -t ghcr.io/onemindos/hermes-legacy:latest --push .
FROM ghcr.io/onemindos/hermes:legacy-v1

# Proven versions (match the Mac's working set). venv is uv-managed (no pip).
RUN uv pip install --python /opt/hermes/.venv/bin/python --no-cache \
      nats-py==2.15.0 \
      synadia-ai-agent-service==0.4.1 \
      synadia-ai-agents==0.7.1

# nats-platform plugin (gateway adapter — NOT part of Hermes core)
COPY plugins/nats-platform /opt/hermes/plugins/nats-platform

# Legacy personas (role-specific, single-voice) + hub config.
# Seeded into HERMES_HOME=/opt/data by cont-init based on session name:
#   session=legacy      → SOUL.mother.md (the main authoritative primary)
#   session=legacy-edge → SOUL.edge.md   (any non-mother edge node)
COPY SOUL.mother.md /opt/hermes/seed/SOUL.mother.md
COPY SOUL.edge.md /opt/hermes/seed/SOUL.edge.md
COPY config.hub.yaml /opt/hermes/seed/config.yaml

# cont-init script: seeds persona/config/plugin into the PVC BEFORE the s6
# supervisor starts the gateway. Runs after 01-hermes-setup (volume chown)
# and before 02-reconcile-profiles (gateway service setup).
COPY 01-seed-legacy /etc/cont-init.d/01-seed-legacy

# No args override — s6 supervisor runs the gateway as normal.
# The nats-platform plugin (enabled in config.yaml) registers on the bus.
