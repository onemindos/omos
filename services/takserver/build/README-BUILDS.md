# TAK Image Builds

Build scripts for custom Docker images derived from official TAK distributions, pushed to GHCR.

## Prerequisites

- Docker installed and running
- Authenticated to GHCR:
  ```bash
  docker login ghcr.io -u <github-username> -p <GHCR_PAT>
  ```
- Source zips in `~/Downloads/`:
  - `takserver-fedhub-docker-5.8-RELEASE-79.zip`
  - `takserver-ldap-5.8-RELEASE-79.zip`

---

## Build FedHub

```bash
bash services/takserver/build/build-fedhub.sh
```

Produces: `ghcr.io/onemindos/fedhub:5.8`

The image contains all four FedHub micro-services (policy, broker, UI, plugin-manager)
started in sequence by `configureInDocker.sh`. Ports: 9100 (UI), 9101 (v1 fed), 9102 (v2 fed).

### After Build: Generate FedHub Certs

```bash
cd ~/Downloads
unzip takserver-fedhub-docker-5.8-RELEASE-79.zip
cd takserver-fedhub-docker-5.8-RELEASE-79/tak/federation-hub/certs
export CA_NAME=onemindos-fedhub
./generateClusterCerts.sh
cd ~/Downloads
```

Seal the generated certs as a K8s Secret:

```bash
kubectl -n takserver create secret generic fedhub-certs \
  --from-file=takserver-fedhub-docker-5.8-RELEASE-79/tak/federation-hub/certs/files/ \
  --dry-run=client -o yaml | kubeseal -o yaml > \
  platform/takserver/fedhub/fedhub-certs.sealedsecret.yaml
```

Then deploy:

```bash
kubectl -n takserver apply -f platform/takserver/fedhub/
```

---

## Build TAK LDAP

```bash
bash services/takserver/build/build-tak-ldap.sh
```

Produces: `ghcr.io/onemindos/tak-ldap:5.8`

Debian 12 OpenLDAP with TAK-specific schema (`takCallSign`, `takRole`, `takColor`) loaded
by bootstrap.sh on first start.

### After Build: Create Admin Password SealedSecret

```bash
kubectl -n takserver create secret generic tak-ldap-creds \
  --from-literal=LDAP_ADMIN_PASSWORD='<strong-password>' \
  --dry-run=client -o yaml | kubeseal -o yaml > \
  platform/takserver/deploy/tak-ldap-creds.sealedsecret.yaml
```

Then deploy:

```bash
kubectl -n takserver apply -f platform/takserver/deploy/tak-ldap.yaml
kubectl -n takserver apply -f platform/takserver/deploy/tak-ldap-creds.sealedsecret.yaml
```

### After LDAP is Running: Patch TAK CoreConfig

Edit the `core-config` ConfigMap in the `takserver` namespace to add the `<auth>` block:

```xml
<auth default="ldap">
  <ldap url="ldap://tak-ldap:389"
        userstring="uid={username},ou=People,dc=onemindos,dc=dev"
        style="DS"
        serviceAccountDN="cn=ldapadm,dc=onemindos,dc=dev"
        serviceAccountCredential="<LDAP_ADMIN_PASSWORD>"
        callsignAttribute="takCallSign"
        colorAttribute="takColor"
        roleAttribute="takRole"
        updateinterval="60000" />
  <File location="UserAuthenticationFile.xml"/>
</auth>
```

Then restart in order: `takserver-config` → `takserver-messaging` → `takserver-plugins`.
