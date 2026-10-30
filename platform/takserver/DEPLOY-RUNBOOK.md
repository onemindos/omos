# TAK Server 5.7 Cluster — Deploy Runbook (OneMind / Option B: build-from-distro)

**Status:** ✅ Proven end-to-end on local `kind` cluster (2026-07-05). All 8 pods 1/1 Ready;
Ignite grid `servers=2, clients=4, state=ACTIVE`; API serving WebTAK on 8443/8446; messaging
CoT TLS input on 8089; 64-table TAK schema. This runbook reproduces that, then promotes to DOKS.

This is **Option B**: we build every TAK image ourselves from the official distro (no tak.gov
registry pull). This gives us a self-owned, air-gappable image supply chain for OneMind.

---

## 0. Source material (paths on Zeus's Mac)

| Thing | Path |
|---|---|
| TAK distro (5.7-RELEASE-43) | `/Users/zeus/ATAK/server/standard/takserver-docker-5.7-RELEASE-43/tak/` |
| Cluster blueprint (upstream) | `/Users/zeus/omos-infra/tak-server/` |
| Helm chart | `.../chart/` |
| Dockerfiles | `services/takserver/docker/` (moved 2026-10-17) |
| Build context (staged) | `/Users/zeus/omos-infra/tak-server/<distro-root>/` |
| Canonical build script (reference) | `.../upstream/scripts/build-generic.py` (lines 715–847) |
| Generated deploy files | `/Users/zeus/omos-infra/tak-server/deploy/` (CoreConfig.xml, TAKIgniteConfig.xml) |
| Generated certs | `/Users/zeus/omos-infra/tak-server/certs-out/` |
| kind values | `/Users/zeus/omos-infra/tak-server/values-kind.yaml` |

Always: `export PATH="/opt/homebrew/bin:$PATH"` (docker, kind, kubectl, helm, doctl live there).

---

## 1. The 7 images (all built from distro)

| Image | Dockerfile | Build context | Notes |
|---|---|---|---|
| `takserver-base:5.8-onemind` | `Dockerfile.takserver-base` | `<distro-root>` | base layer; runs `db-connection-configuration.sh`; bakes distro certs into `/opt/tak/certs/files/` |
| `takserver-config:5.8-onemind` | `Dockerfile.takserver-config` | `<distro-root>` | cluster leader (spring profile `config`) |
| `takserver-api:5.8-onemind` | `Dockerfile.takserver-api` | `<distro-root>` | WebTAK/API (`api,consolelog`) |
| `takserver-messaging:5.8-onemind` | `Dockerfile.takserver-messaging` | `<distro-root>` | CoT TLS input 8089 |
| `takserver-plugins:5.8-onemind` | `Dockerfile.takserver-plugins` | `<distro-root>` | launcher `takserver-plugins-cluster.sh`; **needs `PLUGIN_MANAGER_MAX_HEAP` env** |
| `takserver-database-setup:5.8-onemind` | `Dockerfile.database-setup` | `<distro-root>` | `SchemaManager.jar` + `generic-cluster-database-configuration.sh` |
| `docker-devtest-local.artifacts.tak.gov/takserver-cluster/takserver-ignite:2.17.0` | `Dockerfile.takserver-ignite` | `upstream` | downloads Apache Ignite 2.17.0 + copies k8s discovery jars; needs `docker-files/scripts/takserver-ignite/ignite.sh` in context |

**Staging inputs into <distro-root>** (needed for db-setup + CA images):
```bash
TAK=/Users/zeus/ATAK/server/standard/takserver-docker-5.7-RELEASE-43/tak
cp "$TAK/db-utils/generic-cluster-database-configuration.sh" <distro-root>/takserver-schemamanager/
mkdir -p <distro-root>/takserver-core/certs && cp -R "$TAK/certs/." <distro-root>/takserver-core/certs/
```

**Build all** (from `/Users/zeus/omos/services/takserver`; the CI build is `.github/workflows/build-takserver.yml`):
```bash
docker build -f docker/Dockerfile.takserver-base      -t takserver-base:5.8-onemind          <distro-root>
docker build -f docker/Dockerfile.takserver-config    -t takserver-config:5.8-onemind        <distro-root>
docker build -f docker/Dockerfile.takserver-api       -t takserver-api:5.8-onemind           <distro-root>
docker build -f docker/Dockerfile.takserver-messaging -t takserver-messaging:5.8-onemind     <distro-root>
docker build -f docker/Dockerfile.takserver-plugins   -t takserver-plugins:5.8-onemind       <distro-root>
docker build -f docker/Dockerfile.database-setup      -t takserver-database-setup:5.8-onemind <distro-root>
docker build -f docker/Dockerfile.takserver-ignite \
  -t docker-devtest-local.artifacts.tak.gov/takserver-cluster/takserver-ignite:2.17.0 upstream
```

---

## 2. Generate cluster CA + certs

```bash
docker build -f docker/Dockerfile.ca -t takserver-ca-setup:5.8-onemind <distro-root>
docker rm -f tak-ca-setup 2>/dev/null
docker run --name tak-ca-setup \
  --env CA_NAME=onemind-tak --env STATE=FL --env CITY=Tampa --env ORGANIZATIONAL_UNIT=OneMind \
  takserver-ca-setup:5.8-onemind
docker cp tak-ca-setup:/files ./certs-out    # ~31 files, storepass = the keystore password (default until step-2 rotation; see Secret tak-keystores in ns takserver)
```
Required outputs (verify present): `takserver.jks`, `truststore-root.jks`, `fed-truststore.jks`, `admin.pem`.
`takserver.jks` = single PrivateKeyEntry, alias **`takserver`** (JwtUtils uses the keystore's first alias).

---

## 3. Config files

### `deploy/CoreConfig.xml`  (from distro, DB pointed at the in-cluster postgres)
```bash
TAK=/Users/zeus/ATAK/server/standard/takserver-docker-5.7-RELEASE-43/tak
sed 's#jdbc:postgresql://tak-database:5432/cot#jdbc:postgresql://takserver-postgresql:5432/cot#;
     s/username="martiuser"/username="postgres"/;
     s/password="martipassword"/password="postgres"/' \
  "$TAK/CoreConfig.xml" > deploy/CoreConfig.xml
```
> Cert paths inside CoreConfig can stay at the distro default `/opt/tak/certs/files/...` — the
> **base image bakes the certs there**, so keystores resolve at runtime. (Repointing to
> `/certs-configmap/` is optional and, because of subPath mounts, does NOT take effect without pod
> recreation — see gotcha #8. Leaving the baked path is simplest.)

### `deploy/TAKIgniteConfig.xml`  ← **THE critical cluster toggle**
```xml
<?xml version="1.0" encoding="UTF-8" standalone="yes"?>
<TAKIgniteConfiguration xmlns="http://bbn.com/marti/xml/config"
    clusterEnabled="true"
    clusterKubernetes="true"
    igniteClusterNamespace="takserver">
</TAKIgniteConfiguration>
```
Why: `tak.server.ignite.IgniteConfigurationHolder` does
`if (isClusterKubernetes()) { finder = new TcpDiscoveryKubernetesIpFinder(); finder.setServiceName("takserver-ignite"); finder.setNamespace(getIgniteClusterNamespace()); }`
**else** it falls back to a localhost VmIpFinder (`127.0.0.1:47500+`) and never joins the grid.
The service name `takserver-ignite` is **hardcoded** → your release name MUST make the ignite
service resolve to `takserver-ignite` (helm release name `takserver` does this).

---

## 4. Helm chart patches (already applied in repo)

Both in `chart/templates/`:

1. **`takserver-core-deployment.yaml` + `takserver-plugins-deployment.yaml`** — make pull policy overridable:
   `imagePullPolicy: Always` → `imagePullPolicy: {{ .Values.imagePullPolicy | default "Always" }}`

2. **`takserver-plugins-deployment.yaml`** — add memory + heap env (plugins launcher does
   `-Xmx${PLUGIN_MANAGER_MAX_HEAP}m`; unset → `-Xmxm` → JVM won't start):
   ```yaml
   env:
     - name: "PLUGIN_MANAGER_MAX_HEAP"
       value: "{{ .Values.takserver.plugins.maxHeapMb | default 1024 }}"
   resources:
     requests: { cpu: ..., memory: {{ .Values.takserver.plugins.resources.requests.memory }} }
     limits:   { cpu: ..., memory: {{ .Values.takserver.plugins.resources.limits.memory }} }
   ```

Ignite pull policy is already templated in the subchart: set `ignite.image.pullPolicy` in values.

---

## 5. Namespace objects (configmaps + pull secret)

```bash
kc() { kubectl --context kind-tak -n takserver "$@"; }   # DOKS: swap context
kubectl --context kind-tak create ns takserver 2>/dev/null

kc create configmap cert-migration     --from-file=certs-out/
kc create configmap core-config        --from-file=CoreConfig.xml=deploy/CoreConfig.xml
kc create configmap tak-ignite-config  --from-file=TAKIgniteConfig.xml=deploy/TAKIgniteConfig.xml
kc create configmap ignite-extensions  --from-file=ignite-extensions.jar="$TAK/ignite-extensions-5.7-RELEASE-43.jar"

# kind (local images, no real registry): dummy secret to satisfy imagePullSecrets ref
kc create secret docker-registry reg-cred --docker-server=localhost --docker-username=x --docker-password=x
```
Configmap sizes are all < 1 MB (certs ~120 KB, ignite-extensions.jar 1.8 KB) — safe.

---

## 6. Deploy (helm)

```bash
helm dependency build chart     # deps (postgresql 12.1.6, nats 1.0.1, ignite 2.17.0) vendored as charts/*.tgz
helm --kube-context kind-tak upgrade --install takserver chart \
  --namespace takserver --create-namespace -f values-kind.yaml --timeout 15m
```
`values-kind.yaml` overrides: local image repos+tags, `imagePullPolicy: IfNotPresent`,
`ignite.image.{repository,tag,pullPolicy}`, small resources (12 GB node), nats single-node,
postgres = `postgis/postgis:15-3.4`. See the file for exact values.

---

## 7. Post-deploy fixes (REQUIRED — ordering matters)

The bitnami postgres subchart generates a **random** superuser password, which won't match
CoreConfig's `postgres/postgres`, so the db-setup hook silently fails (no schema). Fix:

```bash
# 7a. reset postgres superuser pw to match CoreConfig
PW=$(kc get secret takserver-postgresql -o jsonpath='{.data.postgres-password}' | base64 -d)
kc exec takserver-postgresql-0 -- env PGPASSWORD="$PW" psql -U postgres -d postgres \
  -c "ALTER USER postgres PASSWORD 'postgres';"

# 7b. create the TAK schema (94 migrations, 64 tables) with working creds
kubectl --context kind-tak -n takserver run tak-schema --restart=Never --attach \
  --image=takserver-database-setup:5.8-onemind --image-pull-policy=IfNotPresent \
  --env POSTGRES_HOST=takserver-postgresql --env POSTGRES_PORT=5432 \
  --env POSTGRES_DB=cot --env POSTGRES_USER=postgres --env POSTGRES_PASSWORD=postgres \
  --env region=local --env identifier=local
kc delete pod tak-schema

# 7c. recreate core pods so they boot against the ready DB
kc delete pod -l app=takserver-config; kc delete pod -l app=takserver-messaging; kc delete pod -l app=takserver-api
```
> **DOKS improvement:** set `global.postgresql.auth.postgresPassword: postgres` (and
> `postgresql.auth.postgresPassword: postgres`) in values **before** first install so the
> superuser password matches from the start — then 7a is unnecessary and the db-setup hook
> succeeds on its own. (Still verify schema exists; re-run SchemaManager if not.)

---

## 8. Verify green

```bash
kc get pods                          # expect 8x 1/1 (+ db-setup Completed)
kc logs deploy/takserver-config | grep "Topology snapshot" | tail -1   # servers=2, clients=4, ACTIVE
kc logs deploy/takserver-api | grep "Tomcat started"                    # 8443/8444/8446 https
kc exec deploy/takserver-api -- curl -sk https://localhost:8446/Marti/api/version   # WebTAK responds
```

---

## GOTCHAS (all hit + solved on 2026-07-05)

1. **Disk full → Docker daemon hangs on boot.** Freed 660 GB; daemon came up clean.
2. **Docker Desktop k8s conflict.** `KubernetesEnabled:true` wedged the daemon. Set false; use `kind`.
3. **kind node v1.36.1 segfaults (exit 139).** Use `kindest/node:v1.31.2`.
4. **Ignite `kind load` corruption** (`unexpected media type text/html` on unpack). Fix: `docker save img -o /tmp/x.tar && kind load image-archive /tmp/x.tar --name tak`.
5. **plugins `-Xmxm`.** Launcher uses `-Xmx${PLUGIN_MANAGER_MAX_HEAP}m`; set the env (patch #4.2).
6. **postgres password mismatch.** bitnami random superuser pw ≠ CoreConfig. Reset (7a) or preset in values.
7. **Missing schema.** db-setup hook ran under the bad pw → no tables. Re-run SchemaManager (7b).
8. **subPath configmap freeze.** configmaps mounted via `subPath` (`/CoreConfig.xml`, `/TAKIgniteConfig.xml`) do NOT hot-update; recreate pods after any configmap change.
9. **Ignite discovery = localhost.** Missing `clusterKubernetes="true"` in TAKIgniteConfig (§3).
10. **Ignite pull policy Always.** Set `ignite.image.pullPolicy: IfNotPresent` for local images.

---

## DOKS PROMOTION — PROVEN (2026-07-05, cluster `onemind-hub` nyc3)

**Result:** all 8 pods 1/1 Ready on DOKS; API serving WebTAK on 8446; messaging CoT TLS 8089;
Ignite grid `servers=2, clients=4, ACTIVE`; 64-table schema auto-created (no pw reset needed).

Facts: cluster `onemind-hub`, k8s 1.36.0-do.2, context `do-nyc3-onemind-hub`, node pool
3× s-4vcpu-8gb (~6.5Gi alloc/node), DOCR `registry.digitalocean.com/onemind-hub` (same region),
default SC `do-block-storage` (+ `-retain` variants). `takserver` ns was clean.

### ⚠️ Arch mismatch — THE big DOKS gotcha
Mac builds **arm64**; DOKS nodes are **amd64** → `ImagePullBackOff` with
`no match for platform in manifest`. Rebuild every image `linux/amd64` via buildx (fast — TAK
content is arch-independent JARs; only OS/JRE layers differ). The core Dockerfiles default their
`FROM` to the tak.gov registry via **build-args** — you MUST override them to the local base:

```bash
export PATH="/opt/homebrew/bin:$PATH"; cd /Users/zeus/omos-infra/tak-cluster
REG=registry.digitalocean.com/onemind-hub; D=upstream/docker-files
doctl registry login

# 1. amd64 base (load locally so derivatives FROM it)
docker buildx build --platform linux/amd64 -f $D/Dockerfile.takserver-base -t takserver-base:5.8-onemind --load <distro-root>

# 2. amd64 ignite (independent) + db-setup -> push
docker buildx build --platform linux/amd64 -f $D/Dockerfile.takserver-ignite -t $REG/takserver-ignite:2.17.0 --push upstream
docker buildx build --platform linux/amd64 -f $D/Dockerfile.database-setup   -t $REG/takserver-database-setup:5.8-onemind --push <distro-root>

# 3. amd64 core images — MUST pass base build-args (else FROM hits tak.gov = text/html error)
for n in config api messaging plugins; do
  case $n in config)DF=Dockerfile.takserver-config;;api)DF=Dockerfile.takserver-api;;messaging)DF=Dockerfile.takserver-messaging;;plugins)DF=Dockerfile.takserver-plugins;;esac
  docker buildx build --platform linux/amd64 \
    --build-arg TAKSERVER_IMAGE_REPO=takserver-base --build-arg TAKSERVER_IMAGE_TAG=5.8-onemind \
    -f $D/$DF -t $REG/takserver-$n:5.8-onemind --push <distro-root>
done
```
> zsh note: unquoted `$VAR` does NOT word-split — inline `--build-arg` flags, don't stuff them in a var.
> Verify arch: `docker buildx imagetools inspect $REG/takserver-config:5.8-onemind | grep Platform`.

### DOKS-specific fixes (all captured in `values-doks.yaml` + template patches)
1. **postgres pw preset works** → set `postgresPassword: postgres` (global + subchart auth) → db-setup
   hook succeeds on its own, schema auto-created, **no ALTER USER / manual SchemaManager needed.** ✅
2. **Ignite subchart has NO imagePullSecrets** → its SA (`takserver-ignite`) can't pull from DOCR (401).
   Fix (durable, applied): patched `third-party/ignite/templates/service-account.yaml` to render
   `imagePullSecrets` from `.Values.ignite.imagePullSecrets`; set that value to `registry-onemind-hub`.
   (Quick live patch if needed: `kubectl -n takserver patch sa takserver-ignite -p '{"imagePullSecrets":[{"name":"registry-onemind-hub"}]}'` then delete the ignite pod.)
3. **Real pull secret:** `doctl registry kubernetes-manifest --namespace takserver | kubectl apply -f -` → `registry-onemind-hub`.

### DOKS deploy sequence (clean, GitOps-ready)
```bash
CTX=do-nyc3-onemind-hub
kubectl --context $CTX create ns takserver
doctl registry kubernetes-manifest --namespace takserver | kubectl --context $CTX apply -f -
kubectl --context $CTX -n takserver create configmap cert-migration    --from-file=certs-out/
kubectl --context $CTX -n takserver create configmap core-config       --from-file=CoreConfig.xml=deploy/CoreConfig.xml
kubectl --context $CTX -n takserver create configmap tak-ignite-config --from-file=TAKIgniteConfig.xml=deploy/TAKIgniteConfig.xml
kubectl --context $CTX -n takserver create configmap ignite-extensions --from-file=ignite-extensions.jar="$TAK/ignite-extensions-5.7-RELEASE-43.jar"
helm --kube-context $CTX upgrade --install takserver chart -n takserver -f values-doks.yaml --timeout 15m
# NOTE: helm blocks on the db-setup post-install hook; that's expected. Watch pods separately.
```
> The 2026-07-05 first run left the release in `pending-install` because the CLI wait was
> interrupted; with amd64 images + the ignite SA fix in place, a fresh `helm upgrade --install`
> (or ArgoCD sync) reconciles cleanly.

---

## DOKS PROMOTION DELTAS (reference)

- **Registry:** push all 7 images to DOCR `registry.digitalocean.com/<name>/...` (or GHCR).
  Retag + push; set `values-doks.yaml` image repos accordingly.
- **imagePullSecret:** real one — `doctl registry kubernetes-manifest | kubectl apply -f -`
  creates `registry-<name>`; set `imagePullSecret` to it; `imagePullPolicy: Always`.
- **Storage:** default SC `do-block-storage` (verify). postgres + ignite PVCs use it.
- **Resources:** scale UP — production-values requests up to 8 Gi + 1–4 CPU/pod. Size the
  DOKS node pool to fit config(1.5Gi)+messaging(3Gi)+api(3Gi)+plugins(1.5Gi)+2×ignite(2Gi)+pg+nats.
- **Postgres:** preset `postgresPassword` in values (avoids the §7a reset). Consider a DO Managed
  PG later; for now the subchart is fine.
- **Ingress:** enable real ingress/LoadBalancer (`takserver.ingress.enabled: true`) → DO LB.
- **clusterKubernetes/service name:** unchanged — release name `takserver` keeps `takserver-ignite` valid.
- Everything else (configmaps, certs, SchemaManager, patches) is identical.
