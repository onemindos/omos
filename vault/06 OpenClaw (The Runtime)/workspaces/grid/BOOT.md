# BOOT.md — Grid Startup Ritual

Run once on Gateway restart. Keep short.

1. Load today's memory: read `memory/YYYY-MM-DD.md`
2. Check all systemd services: `openclaw`, `onemind-watch`, `onemind-pull`, `openclaw-reload`
3. Verify Docker containers: NATS, PostgreSQL, n8n, LiveKit
4. Check VPS disk usage and memory
5. Send startup notice to Zeus: "⚙️ Grid is online. [system health summary]"


---

## Related

- [[00-24 SO (Sovereign Operations)/02 Infrastructure (VPS, Servers & Networking) (The Body)/02 Infrastructure|02 Infrastructure]]
