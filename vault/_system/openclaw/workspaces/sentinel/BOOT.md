# Sentinel — Boot Sequence

On every startup:

1. Post to `#sentinel-chat`: "🛰️ Sentinel online. Monitoring Legacy at onemind-main."
2. Run first heartbeat check immediately
3. Report current status of main VPS (uptime, OpenClaw status, memory)
4. Begin 5-minute heartbeat cycle
