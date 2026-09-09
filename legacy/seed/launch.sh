#!/command/with-contenv sh
# Seed Legacy's persona + config + nats-platform plugin into the PVC on first
# boot, then launch the gateway headless. The gateway's nats-platform adapter
# registers AgentService at agents.prompt.hermes.zeus.legacy and routes prompts
# through Hermes's full pipeline (persona, tools, skills, session memory).
set -e

# --- Seed (self-heal: only copies if missing) ---
mkdir -p /opt/data/plugins

if [ -d /opt/hermes/plugins/nats-platform ] && [ ! -d /opt/data/plugins/nats-platform ]; then
    cp -r /opt/hermes/plugins/nats-platform /opt/data/plugins/nats-platform
    echo "[seed] installed nats-platform plugin"
fi

if [ -f /opt/hermes/seed/SOUL.md ] && [ ! -f /opt/data/SOUL.md ]; then
    cp /opt/hermes/seed/SOUL.md /opt/data/SOUL.md
    echo "[seed] installed SOUL.md (Legacy persona)"
fi

if [ -f /opt/hermes/seed/config.yaml ] && [ ! -f /opt/data/config.yaml ]; then
    cp /opt/hermes/seed/config.yaml /opt/data/config.yaml
    echo "[seed] installed config.yaml (model + NATS platform, session=legacy)"
fi

# --- Launch the gateway headless ---
# The nats-platform plugin (enabled in config.yaml) registers on the bus.
# --no-supervise: don't spawn the dashboard; we only want the gateway.
exec hermes gateway run --replace --force --no-supervise
