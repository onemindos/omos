# Native Grafana folder links

Status: implementation prepared; live acceptance pending normal merge/Argo.
User requested the folder fix before Prometheus alert integration.

## Cause and bounded fix

Installed `@backstage-community/plugin-grafana` 1.1.0 constructs `folderUrl` as
`domain + dashboard.folderUrl`. Grafana omits folder metadata for root dashboards;
this produced a blank `...ts.netundefined/` folder link despite working title links.
Provision the three OneMind service dashboards under `/tmp/dashboards/OneMind`
using the dedicated `onemind.io/grafana-folder` ConfigMap annotation and the pinned
chart's `sidecar.dashboards.provider.foldersFromFilesStructure` option. The chart
omits conflicting `folder`/`folderUid` fields when enabled. Unannotated dashboard
ConfigMaps retain their root location; no default folder change or manual UI move.

This fixes the three mapped service dashboards without patching node_modules or
rebuilding the portal image. It is not a general fix for every root-level result
in the upstream plugin; newly mapped boards should use a real provisioned folder.

[Grafana provisioning documentation](https://grafana.com/docs/grafana/latest/administration/provisioning/)
documents filesystem folder provisioning, preserved dashboard UIDs and the need to
omit fixed folder/folderUid. The checked-in chart version remains 87.16.1 and the
bundled Grafana subchart 12.7.2 was inspected for the exact rendered behavior.

## Safety and acceptance

- Three dashboard JSON payloads, UIDs, selectors and all panel queries unchanged.
- No credentials/RBAC/proxy methods/authentication/image/PVC/backup changes.
- Helm change limited to Grafana provider ConfigMap and Grafana Deployment
  (sidecar annotation environment and config checksum); brief Recreate restart.
- Before change, Grafana search lists 31 dashboards, three `onemind-*` service UIDs
  at root, no dashboard ConfigMap already uses a folder annotation.
- After change require all 31 original dashboard UIDs retained, only the three
  service UIDs moved, proper OneMind folderTitle and relative folderUrl returned
  to the actual Backstage Viewer, 2/2 portal and healthy PVC-backed Grafana.
- Signed-in Backstage must show readable OneMind folder links on each of the
  three entities; follow a folder link and verify its real Grafana listing.
- CPU/memory data and existing anonymous/machine denial boundaries remain intact.

No Prometheus alert changes are included. Broader alert integration is next;
Grafana-managed alerts are still distinct from Prometheus alerts. Historical
cutover evidence remains in `grafana-durable-visibility-2026-10-20.md`.
