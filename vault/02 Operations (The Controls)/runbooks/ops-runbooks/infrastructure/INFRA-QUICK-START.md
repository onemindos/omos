# Monitoring Quick Start

## 1. Configure Uptime Kuma (10 minutes)

Access: http://YOUR_IP:3001

Add monitors for:
- Agno API: http://localhost:8080/health (60s interval)
- API Gateway: http://localhost:8000/health (60s interval)
- NATS: http://localhost:8222/healthz (60s interval)
- PostgreSQL: Port check localhost:5432 (60s interval)

## 2. Setup Discord Alerts (5 minutes)

1. Create Discord webhook in your server
2. Settings → Notifications → Add Discord
3. Paste webhook URL
4. Test notification

## 3. Schedule Resource Monitoring (5 minutes)

```bash
crontab -e

# Add:
*/15 * * * * /opt/onemind/production/3-monitoring/monitor-resources.sh >> /var/log/onemind-resources.log 2>&1

# Set Discord webhook:
echo 'export DISCORD_WEBHOOK_URL="YOUR_WEBHOOK"' >> ~/.bashrc
source ~/.bashrc
```

Done! See [../REFERENCE.md](../REFERENCE.md) for full monitoring setup including Prometheus/Grafana.


---

## Related

- [[00-24 SO (Sovereign Operations)/02 Infrastructure (VPS, Servers & Networking) (The Body)/02 Infrastructure|02 Infrastructure]]
