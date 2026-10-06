# Osiris

Archived from the Fabric registry on 2026-10-06. Source edited: 2026-09-03 07:01:04Z.

These are historical documentation claims, not current deployment instructions or live health.

- Category: 📡 Intel
- Documented host: 🐳 Docker Local
- Recorded status: 📚 Reference
- Recorded address: https://github.com/simplifaisoul/osiris

Lives in omos-infra/platform/intel/osiris — pending standalone repo extraction

## Local Deployment (Docker Desktop)
Run locally to review features and concepts. Clone from GitHub link above, or use the pre-built onemindos image.
### docker-compose.yml
```yaml
version: '3.8'
services:
  osiris:
    image: ghcr.io/onemindos/osiris:latest
    # OR build from source:
    # build: .  # clone https://github.com/simplifaisoul/osiris
    ports:
      - "3000:3000"
    environment:
      NODE_OPTIONS: "--dns-result-order=ipv4first"
      NODE_ENV: production
      PORT: "3000"
      HOSTNAME: "0.0.0.0"
    restart: unless-stopped

  osiris-intel:
    image: ghcr.io/onemindos/osiris-intel:latest
    ports:
      - "4000:4000"
    restart: unless-stopped

  osiris-cache:
    image: nginx:alpine
    ports:
      - "8080:8080"
    restart: unless-stopped
```
### Quick Start
```shell
# Save the above as docker-compose.yml then:
docker compose up -d

# View logs
docker compose logs -f

# Stop
docker compose down
```
<callout icon="📌">
	Removed from DOKS cluster — reference project only. Run locally when reviewing features to absorb into OneMind.
</callout>
