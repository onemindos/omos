# GeoServer on DOKS

Geospatial layer engine (WMS/WMTS/WFS/WCS) for the MapLibre COP — self-hosted
map layers, tactical overlays, PostGIS-backed feature layers.

Migrated off the Mac (2026-07-05). The local `kartoza/geoserver:2.25.2`
container (`/Users/zeus/opt/geoserver/`) was **empty** (0 workspaces/layers), so
this is a fresh deploy — nothing to migrate.

## Components (`deployment.yaml`)
- Namespace `geoserver`
- PVC `geoserver-data` (5Gi `do-block-storage`) → `/opt/geoserver/data_dir`
- Deployment `geoserver` (`kartoza/geoserver:2.25.2`, pinned)
- Service `geoserver:8080` (ClusterIP — internal only; MapLibre COP consumes in-cluster)

Probes: TCP startup+readiness (GeoServer's Wicket UI page is too heavy for a
10s httpGet probe → pods never went Ready with httpGet); httpGet liveness with a
generous 25s timeout.

## Admin secret (out-of-band, NOT in git)
```sh
kubectl -n geoserver create secret generic geoserver-admin \
  --from-literal=password='<strong-password>'   # user: admin  -> store in 1Password
```

## Deploy
```sh
kubectl apply -f deployment.yaml
kubectl -n geoserver create secret generic geoserver-admin --from-literal=password='<pw>' \
  --dry-run=client -o yaml | kubectl apply -f -
kubectl -n geoserver rollout status deploy/geoserver
```

## Verify
```sh
GSPW=$(kubectl -n geoserver get secret geoserver-admin -o jsonpath='{.data.password}' | base64 -d)
kubectl -n geoserver run t --rm -i --restart=Never --image=curlimages/curl -- \
  -s -u "admin:$GSPW" http://geoserver:8080/geoserver/rest/about/version.json
# web UI: http://geoserver:8080/geoserver/web/ (302 -> login)
```

## Exposure
Internal ClusterIP for now. To reach the web UI externally, add a Cloudflare
tunnel route (behind Access) like argocd/grafana — deferred until there are
layers worth managing.

## TODO
- Wire into Argo CD app-of-apps (GitOps auto-sync) during observability/GitOps hardening.
- Add PostGIS datastore + publish first layers when COP data is ready.
