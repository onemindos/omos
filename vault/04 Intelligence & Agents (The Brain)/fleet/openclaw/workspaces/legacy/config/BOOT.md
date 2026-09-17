# BOOT.md — Legacy Startup Ritual

Run once on Gateway restart. Keep short.

1. Check NATS bridge health: attempt `nexus.system.health` ping
2. Verify NATS connectivity — confirm nexus-bridge can reach NATS at localhost:4222
3. Check that onemind-watch and onemind-pull services are active
4. Load today's memory: read `memory/YYYY-MM-DD.md`
5. Read AMBIENT_STATE.md if it exists for current environmental context
6. Send startup notice to Zeus: "⚡ Legacy is online. [system status summary]"
7. Schedule morning brief cron if not already set (0 7 * * *)


---

## Related

- [[00-24 SO (Sovereign Operations)/02 Infrastructure (VPS, Servers & Networking) (The Body)/02 Infrastructure|02 Infrastructure]]
