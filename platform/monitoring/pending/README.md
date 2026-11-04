# Pending Alertmanager receiver: OneMind portal

`alertmanagerconfig-onemind-backstage.yaml.disabled` sends critical and warning alerts to
the portal (`POST /api/events/http/alertmanager`). It is **not deployed**:

- The `.disabled` suffix means no Argo CD source or kustomization includes it. Argo's
  `platform-monitoring` Application reads `platform/monitoring` recursively and only treats
  `.yaml`, `.yml`, `.json` and `.jsonnet` files as manifests.
- It depends on a secret that does not exist yet. Do not enable it before both SealedSecrets below exist.
- `scripts/tests/test_gitops_alert_routing.py` fails if the file is renamed to a manifest
  extension, listed in a kustomization, or if an equivalent AlertmanagerConfig appears under
  `platform/monitoring` without both valid SealedSecrets and the backstage kustomization entry. It also
  checks that `scripts/ops/seal-alert-webhook.sh` cannot be reached by any Argo CD source.

Alertmanager is managed by the kube-prometheus-stack release (Argo Application
`monitoring-kube-prometheus-stack`). Its `alertmanagerConfigSelector: {}` makes the operator
merge every `AlertmanagerConfig` in the `monitoring` namespace into the live config, with no
namespace matcher injected (`alertmanagerConfigMatcherStrategy: None`). No other
`AlertmanagerConfig` exists today (the `openchoreo` one was removed with the OpenChoreo
observability plane). A new config with `continue: true` adds a route and leaves the
default routing untouched.

## Secret (Zeus creates it; the same random value goes in two namespaces)

| Namespace | Secret name | Key | File | Used by |
|---|---|---|---|---|
| `monitoring` | `onemind-alert-webhook` | `token` | `platform/monitoring/onemind-alert-webhook.sealedsecret.yaml` | this AlertmanagerConfig (`httpConfig.authorization`) |
| `backstage` | `onemind-alert-webhook` | `ONEMIND_ALERT_WEBHOOK_TOKEN` | `platform/backstage/alert-webhook.sealedsecret.yaml` | portal env var, optional |

`scripts/ops/seal-alert-webhook.sh` creates both files. It generates the token in memory (never
printed, never written unsealed, never in a command line), seals it with the controller certificate
(`kubeseal --fetch-cert`, `--scope namespace-wide`, as in `CLAUDE.md`), refuses to overwrite, and is a
dry run unless you pass `--apply`. It is not part of anything Argo syncs and changes nothing in the cluster.

## Enable (one PR; the steps in order)

1. Check the portal image pinned in `platform/backstage/deployment.yaml` is the release that contains
   F3. Before that the route does not exist and there is nothing to enable.
2. `git switch -c ops/enable-alert-webhook origin/main`
3. Dry run, and read the output (every line should say `ok`):
   `scripts/ops/seal-alert-webhook.sh`
4. Write the two files: `scripts/ops/seal-alert-webhook.sh --apply`
5. Review the diff. `git add -N platform/monitoring/onemind-alert-webhook.sealedsecret.yaml platform/backstage/alert-webhook.sealedsecret.yaml && git diff`
   shows two `kind: SealedSecret` files named `onemind-alert-webhook` in `monitoring` and `backstage`,
   each with only `encryptedData` (no `data:` or `stringData:`). Nothing else should have changed.
6. Add `  - alert-webhook.sealedsecret.yaml` to `resources:` in `platform/backstage/kustomization.yaml`.
7. In `platform/backstage/deployment.yaml` add under the container's `envFrom` (keep it optional so the
   pod still starts if the secret is ever removed):

   ```yaml
   - secretRef:
       name: onemind-alert-webhook
       optional: true
   ```

8. `git mv platform/monitoring/pending/alertmanagerconfig-onemind-backstage.yaml.disabled platform/monitoring/alertmanagerconfig-onemind-backstage.yaml`
9. `python3 -m unittest scripts/tests/test_gitops_alert_routing.py` and `make validate`, then open the PR.
   The guard test fails if step 8 is done without both valid SealedSecrets and step 6.
10. Merge. Argo syncs `platform-monitoring` (secret and AlertmanagerConfig) and the backstage app (secret,
    env var, pod restart).

If the secret in `monitoring` is missing when the AlertmanagerConfig is applied, the operator
skips that one config and logs the reason; it does not stop the other receivers.

## Verify (after Argo shows both apps Synced and Healthy)

1. Route is loaded. `kubectl --context do-nyc3-omos-cloud -n monitoring get alertmanagerconfig onemind-backstage`
   exists. In Alertmanager (`kubectl --context do-nyc3-omos-cloud -n monitoring port-forward svc/monitoring-kube-prometheus-alertmanager 9093:9093`,
   then `http://localhost:9093/#/status`) the config lists a receiver ending in `onemind-backstage`.
2. Portal has the token. With `kubectl --context do-nyc3-omos-cloud -n backstage port-forward svc/backstage 7007:7007`,
   `curl -s -o /dev/null -w '%{http_code}\n' -X POST localhost:7007/api/events/http/alertmanager -H 'content-type: application/json' -d '{}'`
   must print `401` (token loaded, request unauthenticated). `503` means the portal has no token
   (the `envFrom` or the secret is missing; restart the pod after fixing). Neither command prints the token.
3. Synthetic test alert, through the same port-forward to Alertmanager (harmless; it expires in 2 minutes):

   ```bash
   END="$(date -u -v+2M +%Y-%m-%dT%H:%M:%SZ)"   # Linux: date -u -d '+2 min' +%Y-%m-%dT%H:%M:%SZ
   curl -s -X POST localhost:9093/api/v2/alerts -H 'Content-Type: application/json' \
     -d "[{\"labels\":{\"alertname\":\"OneMindWebhookTest\",\"severity\":\"warning\"},\"annotations\":{\"summary\":\"Synthetic test alert for the portal webhook. Safe to ignore.\"},\"endsAt\":\"$END\"}]"
   ```

4. Expected in the portal within about a minute: one new inbox item "Alert: OneMindWebhookTest", severity
   high, topic `alerts`, and a Signal toast for admins. After the alert expires (about 2 minutes) the same
   item updates to "Resolved: OneMindWebhookTest" at low severity. A repeat does not create a second item.
   Other routes may also notify for a `warning` alert (for example Telegram); that is existing behaviour.
5. If nothing arrives: `kubectl --context do-nyc3-omos-cloud -n monitoring logs alertmanager-monitoring-kube-prometheus-alertmanager-0 -c alertmanager | grep -i onemind-backstage`.
   A 401 means the two secrets differ (re-seal both, see Rotate); a 503 means the portal has no token;
   no lines means the route did not load (step 1).

## Roll back

- Fast: revert the enable PR, or `git mv` the receiver back into this folder and merge. Alertmanager drops
  the route, the portal simply receives nothing, and the other receivers were never touched.
- Also remove the SealedSecrets: delete both files and the kustomization and `envFrom` lines. The `envFrom`
  is optional, so leaving the secret out never stops the pod. Deleting a file does not delete the live Secret
  unless the Argo Application prunes; check with `kubectl --context do-nyc3-omos-cloud get secret -A | grep onemind-alert-webhook`.
- Rotate the token: delete both sealed files, run the helper again (it refuses to overwrite), and commit both
  files in one PR so the two namespaces never hold different values.
