# Native Grafana folder links

Status: scoped live acceptance passed on 2026-10-20 at 23:01–23:04 UTC.
[PR #344](https://github.com/onemindos/omos/pull/344) merged as
`48d841b30f9e1304a3f2b863992b51f7af7503b1`; Argo reconciled normally.
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


## Observed production acceptance

Both `platform-monitoring` and `monitoring-kube-prometheus-stack` reached
Synced/Healthy at the exact merge revision above (chart source remains 87.16.1).
Grafana generation 8, pod `monitoring-grafana-6bc597f4b4-rc4r7`, 3/3 ready,
zero restarts, same imageID and Deployment UID, retained seeded `grafana-data`.
Portal remains generation 74, 2/2, same actual pinned imageIDs; DB specification
and controller revision unchanged, NATS 3/3. No imperative workload apply/force
sync or PVC deletion occurred.

All **31 original dashboard UIDs** remain; other 28 dashboards' folder identity
is unchanged. The three service boards now return `folderTitle: OneMind`,
`folderUrl: /dashboards/f/cg1uarp0h61vkc/onemind` to both Backstage replicas'
existing Viewer identity. This generated folder UID is observed, not a new Git
contract; discovery supplies the link rather than hardcoding it in the catalog.

Signed-in native cards for PostgreSQL, NATS and Backstage show the named OneMind
folder link, not a blank `undefined/` URL. The destination was opened in the
existing authenticated Grafana browser and lists exactly the three intended
service dashboards. Actual Prometheus CPU/memory queries still return finite
values for all three workloads after the restart. Dashboard JSON payloads are
byte-for-byte unchanged by this fix; no panel query or catalog selector changed.

Anonymous catalog/proxy GETs remain 401; existing machine catalog GET 200 and
proxy GET 403. The Backstage Kubernetes identity still cannot read Secrets,
exec, patch Deployments or delete pods; read Deployments remains allowed.
Both portal credentials remain non-admin Viewer with dashboard read but no
write/service-account-write permission. Existing backups remain enabled.

Required portable gates passed on exact implementation head
`5b9f75d1eb9bec352223d9ee405ffc37bf055a99`; 9 Grafana contracts pass. Complete
Helm renders differ only in provider ConfigMap and Grafana Deployment. Receipt
`/tmp/onemind-grafana-folder-check/run.json` reports `check_passed`.
`make validate`/optional whole-tree diagnostic still fail on pre-existing YAML
lint debt starting `platform/onetak/media-mtx.yaml:86`; hosted CI remains
billing-blocked, not green. No billing/check-policy/bypass changes were made.

Local redacted acceptance evidence: `/tmp/onemind-grafana-ops/folder-checks.json`,
`folder-apps-after.json`, `folder-deploy-after.json`, `final-checks.json` and
`denial-checks.json`. Screenshots: `/Users/zeus/.codex/visualizations/2026/10/20/grafana-visibility-2026-10-20/folder-list.png`
and `folder-native-postgres.png`. Local receipts are not permanent remote artifacts.

Prometheus alerts remain next and **not deployed by this repair**. Assess a native
Prometheus alert plugin against Backstage 1.55.3/new frontend composition and scoped
rule labels, authenticated read-only proxy, explicit unknown/error states and
real firing/pending/no-match tests. Do not turn an empty Grafana-managed rule list
into a false all-clear for Prometheus alerts.
