# App annotation contract

Every curated app is a catalog `Component` in `services/onemind-portal/catalog/apps.yaml` (the mounted copy,
`platform/backstage/catalog/apps.yaml`, must be identical). This page defines the `onemind.io/*` annotations that
tell the portal where an app runs, how to check it and who to call. Flagship F7.

The values below are enforced by `services/onemind-portal/scripts/test-app-contract.cjs` (the validator is
`scripts/app-contract.cjs`; the list of hosting values is `plugins/plugin-onemind-common/src/hostingValues.json`,
read by both the validator and the portal). **Facts are text only. Never put a password, key, token or recovery
code in any annotation.** The validator rejects values that look like credentials; gitleaks is the backstop.

A fact nobody has recorded stays absent. The app page then shows "Not recorded". Do not fill a field with a guess.

## `onemind.io/hosting` (required on every Component)

Where the app runs. Canonical values:

| Value | Meaning |
| --- | --- |
| `kubernetes` | Runs on the OneMind cluster (DOKS). |
| `digitalocean` | Other DigitalOcean infrastructure that is not Kubernetes (droplet, managed database, Spaces). |
| `saas` | A provider-hosted service we use as a customer (Zoho, Gmail, Dropbox, 1Password, Monarch, Apple). |
| `home` | Home, edge or on-premises service reached over the tailnet or a leaf node. |
| `device` | A physical device or vehicle firmware. |
| `desktop` | A desktop application. |
| `mobile` | A mobile application. |
| `web` | A reference website; not an app OneMind runs or depends on. |
| `planned` | Recorded for future use; not deployed anywhere. |
| `unknown` | Not established. Use this instead of guessing; the portal shows it visibly. |

Legacy values already in the catalog stay valid and are not renamed. The portal maps them to a kind:

| Legacy value | Kind it means |
| --- | --- |
| `digitalocean-nyc3` | `kubernetes` when the Component has an `argocd/app-name`, `onemind.io/kubernetes-workloads` or a Kubernetes selector or id annotation; otherwise `digitalocean`. |
| `external-provider-unverified` | `saas` |
| `mesh-endpoint-unverified` | `home` |
| `native-application` | `desktop` |
| `vehicle-firmware` | `device` |
| `external-website` | `web` |
| `planned-not-deployed` | `planned` |

A Component with `onemind.io/documented-status: Planned` is shown as `planned` whatever its hosting value says.
New entries should use the canonical values. The F8 add-app template already writes `digitalocean-nyc3`, `digitalocean`,
`saas` or `home`.

## Health

| Annotation | Meaning |
| --- | --- |
| `onemind.io/health-check` | The Gatus endpoint name: lower case letters, digits and hyphens, unique across the catalog. It is also the Prometheus `name` label. |
| `onemind.io/health-url` | What Gatus probes: an `http(s)` URL with no credentials. Required when `health-check` is set, and the other way round. |

Gatus runs inside the cluster, so use an in-cluster Service URL (`http://svc.ns.svc.cluster.local/health`) for a cluster
app and a public HTTPS URL for a SaaS app. Tailnet `*.ts.net` URLs are not reachable from the pod. The launcher URL is
deliberately not reused: it is where a person goes, not necessarily what proves the service is up.

`platform/gatus/endpoints.yaml` is generated from these two annotations
(`node services/onemind-portal/scripts/generate-gatus-endpoints.cjs`); see [platform/gatus/README.md](../../platform/gatus/README.md).
Before F7 the F8 template wrote the URL into `health-check`; it now writes the name there and the URL into `health-url`.

## Facts (all optional; text only)

| Annotation | What to write |
| --- | --- |
| `onemind.io/data-export` | How to get your data out, in one sentence ("Settings, Export, full archive"). |
| `onemind.io/backup-location` | Where a copy of the data lives ("Dropbox export folder in the vault"). Not a path to secrets. |
| `onemind.io/recovery-runbook` | A runbook path under the runbook site (`apps/<name>/`) or an `https://` URL. Generated for every cluster app. |
| `onemind.io/renewal` | `YYYY-MM-DD`, or `monthly`, `yearly`, `none`. |
| `onemind.io/plan` | The plan or tier and billing period ("Standard, billed yearly"). No card numbers. |
| `onemind.io/access` | Who has access, as names or roles ("Zeus (owner); Pat (viewer)"). Never how they sign in. |

The owner is `spec.owner`, not an annotation. These facts appear on the **App facts** card of the app page and nowhere else
is checked live: a plan or renewal in the catalog is a record, not an observation.

## Tests

- `test-app-contract.cjs`: every Component has a valid `onemind.io/hosting`; the other annotations are well formed;
  the Gatus names are unique; the two catalog copies are identical.
- `test-gatus-endpoints.cjs`: `platform/gatus/endpoints.yaml` is exactly what the generator produces.
- `plugin-onemind-common/src/hosting.test.ts`: the portal's reading of every value.
