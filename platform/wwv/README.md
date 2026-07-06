# WWV (OneMind Globe) — DOKS Deployment

## Architecture
- **Image**: `ghcr.io/onemindos/wwv:latest` (built from github.com/silvertakana/worldwideview)
- **Hostname**: `wwv.onemindos.dev`
- **Namespace**: `wwv`
- **Tunnel**: `66bbea70` (shared OneMind hub tunnel)
- **Access Policy**: "Zeus Only" via Cloudflare Access (configure in CF dashboard)

## Deploy Steps

1. **Create namespace + secret**:
   ```bash
   kubectl apply -f namespace.yaml
   kubectl create secret generic wwv-env --from-env-file=secret.env -n wwv
   kubectl create secret docker-registry ghcr-pull-secret \
     --docker-server=ghcr.io \
     --docker-username=<gh-user> \
     --docker-password=<gh-pat> \
     -n wwv
   ```

2. **Apply manifests**:
   ```bash
   kubectl apply -f deployment.yaml
   kubectl apply -f cloudflare-route.yaml
   ```

3. **Add Cloudflare Tunnel route** (via CF dashboard or API):
   - Hostname: `wwv.onemindos.dev`
   - Service: `http://wwv.wwv.svc.cluster.local:3000`
   - Access Policy: "Zeus Only" (same as hermes.onemindos.dev)

4. **Build & push image** (from worldwideview repo):
   ```bash
   docker build -t ghcr.io/onemindos/wwv:latest .
   docker push ghcr.io/onemindos/wwv:latest
   ```

## PostgreSQL
WWV requires PostgreSQL. Options:
- Use managed DO database (recommended)
- Deploy a PG pod in the wwv namespace
- Use the existing platform PG if one exists

Set `DATABASE_URL` in the wwv-env secret accordingly.

## Notes
- The app runs DB migrations on container startup via docker-entrypoint.sh
- Health endpoint: `/api/health` (returns unhealthy until DB is connected, but app still serves)
- Cesium Ion token is optional but recommended for full globe imagery
