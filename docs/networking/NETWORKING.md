# OneMind networking model and cut-over runbook

Status: target design plus the first (additive) PR. Written 2026-10-09 against cluster
`omos-cloud` (DOKS nyc3), tailnet `tail717ba0.ts.net`, Tailscale operator v1.98.4.

## 1. Target model

| Rule | Mechanism |
|------|-----------|
| Private by default | Everything is reachable only over the tailnet. Web UIs use `ingressClassName: tailscale`; TCP services use `Service` with `loadBalancerClass: tailscale`. |
| Public only when it must be | A short list of sites goes out through a Cloudflare Tunnel (`cloudflared` -> ingress-nginx ClusterIP). Everything else is not on the internet. |
| No public load balancers | No DigitalOcean `LoadBalancer` Services. Nothing listens on a public IP, so the cluster can move to a home server behind NAT with no inbound ports. |
| One way to add a public host | Annotate the Ingress (external-dns writes the Cloudflare CNAME to the tunnel). One way to add a private host: Tailscale Ingress or Service. |
| Everything through Argo CD | `platform-cloudflared`, `platform-external-dns`, and the `tak-tailnet` Service are in Git. Only the three secrets are created by hand. |

```
 Public users ──► Cloudflare edge ──► tunnel (outbound from cluster) ──► cloudflared x2
                                                                           │
                                                      ingress-nginx ClusterIP :443 ──► existing Ingress objects
 Tailnet users / EUDs ──► Tailscale ──► operator proxy pods ──► tailscale Ingress / tak-tailnet Service ──► pods
 You (kubectl) ──► tailnet API-server proxy   (public control plane locked by DOKS firewall, section 7)
```

### What this PR adds (nothing user-visible changes on merge)

| Piece | Path | State after merge |
|-------|------|-------------------|
| cloudflared connector + PDB | `platform/networking/cloudflared/cloudflared.yaml`, `platform/argocd/platform-cloudflared.yaml` | Deployment at **0 replicas** (inert) |
| external-dns | `platform/argocd/platform-external-dns.yaml` | Deployment at **0 replicas**; annotation filter matches nothing |
| TAK on tailnet | `platform/takserver/deploy/tak-tailscale.yaml` (synced by existing `platform-takserver`) | New tailnet node `tak`; `tak-external` and its DO LB untouched |

Why inert via replicas 0 rather than a guard: the Deployments need a Secret that does
not exist yet. An optional `secretKeyRef` would start the pod and crash-loop (Argo
Degraded, noisy alerts). Replicas 0 keeps Argo Healthy and makes activation an explicit,
reviewable one-line PR after the Secret exists. external-dns additionally ships with an
annotation filter so that even when running it only touches objects that opt in.

Why plain manifests for cloudflared, Helm for external-dns: the only Cloudflare chart for
a token tunnel (`cloudflare-tunnel-remote` 0.1.2) was last released 2024-09, defaults to
`cloudflared:latest` and has no PDB or hardening. The Deployment is small, so plain
manifests with a digest-pinned image are better. external-dns is maintained upstream with
a current chart (1.23.0, 2026-09-30), so it uses Helm, pinned.

## 2. Inventory of current public endpoints

All are Cloudflare-proxied A records to the ingress-nginx LB `143.244.220.155`, except TAK.
Verified against live Ingress objects on 2026-10-09.

| Hostname | Ingress / Service | Backend | Recommendation | Note |
|----------|-------------------|---------|----------------|------|
| `zeusdelacruz.com`, `www.zeusdelacruz.com` | (none) | (none) | **Retired 2026-10-10** | zeus-site service removed; DNS cleanup is a manual owner step (`docs/operations/RETIRE-ZEUS-SITE.md`). |
| `auth.onemindos.dev` | auth/authentik-ingress | Authentik | **Keep public via tunnel** | IdP / OIDC callbacks for users off the tailnet. Private mirror already exists: `onemind-auth.tail717ba0.ts.net`. |
| `matrix.onemindos.com`, `element.onemindos.com` | matrix/matrix-ingress | Synapse, Element | **Keep public via tunnel** | Needed for federation and non-tailnet users. Confirm federation uses `.well-known` delegation to 443 (a tunnel cannot carry raw 8448). |
| `onemindos.dev`, `tiles.onemindos.dev` | onetak/onetak-ingress | onetak-api (CloudTAK), onetak-tiles | **Tailnet only** (tunnel + Cloudflare Access as interim) | Operational C2 UI. Move once every user has Tailscale. Needs a `tailscale` Ingress (future PR). |
| `ws.onemindos.dev` | nats/nats-ws-ingress | nats:8080 websocket | **Tailnet only** | `omos-nats` already exposes 8080 on the tailnet. Only keep public if a browser client off the tailnet needs it. |
| `enroll.onemindos.dev` | onetak/tak-portal-ingress | tak-portal:8080 | **Tailnet only** | EUDs are on Tailscale now. |
| `fedhub.onemindos.dev` | takserver/fedhub-ingress | fedhub-ui | **Tailnet only** | Admin UI. |
| `valhalla`, `nominatim`, `overpass`, `tileserver` `.onemindos.dev` | onetak/geo-services-ingress | geo services (valhalla/nominatim currently at 0 replicas) | **Tailnet only** | Internal tooling; has no reason to be public. |
| `tak.onemindos.dev` (unproxied A -> `152.42.155.95`) | takserver/tak-external (DO LB `onemind-tak`) | tak-lb-proxy | **Tailnet for EUDs** (new `tak` node) | 8089/8443/8446/64738 move to the tailnet. |
| TAK federation 9101/9102 on the same LB | tak-external | fedhub-fed | **Open decision** | Server-to-server TCP with peers outside our tailnet; cannot ride an HTTP tunnel. Either peers join the tailnet (shared node), or this LB stays (about $12/mo). |
| `media.onemindos.dev` | none | none | Remove from docs | `CLAUDE.md` lists it but no Ingress exists. |

Savings: two DO LBs (`lb-small`, about $12/mo each). About $24/mo if federation also
moves to the tailnet; about $12/mo if the TAK LB stays only for 9101/9102.

## 3. Secrets the operator creates (names only; never commit values)

| Secret | Namespace | Key | Used by | Scope / source |
|--------|-----------|-----|---------|----------------|
| `cloudflared-tunnel-token` | `cloudflare` | `token` | cloudflared | Token of a new remotely-managed tunnel (section 4, step 2) |
| `external-dns-cloudflare` | `external-dns` | `api-token` | external-dns | **New dedicated** Cloudflare API token: Zone:DNS:Edit + Zone:Zone:Read, restricted to zones `onemindos.dev`, `onemindos.com` |

Do not reuse `cert-manager/cloudflare-api-token`: it is in another namespace (copying it
spreads one credential across blast radii) and it has different rotation needs. The
`external-dns` namespace is created by Argo on first sync, so create it first if you want
the Secret before merge: `kubectl create namespace external-dns`.

Create without leaving values in shell history (copy the token to the clipboard first):

```bash
pbpaste | tr -d '\n' | kubectl -n cloudflare create secret generic cloudflared-tunnel-token --from-file=token=/dev/stdin
pbpaste | tr -d '\n' | kubectl -n external-dns create secret generic external-dns-cloudflare --from-file=api-token=/dev/stdin
```

To keep them in Git instead, seal them (see `CLAUDE.md` Secrets Convention) into
`platform/networking/cloudflared/cloudflared-tunnel-token.sealedsecret.yaml`. For
external-dns there is no manifest directory (Helm-only Application), so use the
out-of-band method or add a directory source in a later PR.

## 4. Cut-over runbook

Each step lists a rollback. Do not start step N+1 until step N is verified. The existing
LB path stays fully working throughout; nothing public is removed before step 9.

**0. Safety net.** In Cloudflare DNS for each zone, Advanced -> Export DNS records (BIND
file) and keep it. Record the current values: proxied A records -> `143.244.220.155`;
`tak` -> `152.42.155.95` (unproxied).

**1. Merge this PR.** Expected: Argo creates the `external-dns` Deployment (0), the
`cloudflared-tunnel` Deployment (0) and PDB, the `tak-tailnet` Service, and the Tailscale
operator adds node `tak`. No DNS, ingress or LB change. Verify:
`kubectl -n cloudflare get deploy cloudflared-tunnel`, `kubectl -n takserver get svc tak-tailnet`
(EXTERNAL-IP becomes `100.x` plus `tak.tail717ba0.ts.net`).
Rollback: revert the PR (Argo prunes; the operator deletes node `tak`).

**2. Create the tunnel and secrets.** Cloudflare Zero Trust -> Networks -> Tunnels ->
Create tunnel (type Cloudflared, name `omos-cloud`). Copy the token (it appears once) and
create `cloudflared-tunnel-token` (section 3). Note the tunnel ID. Create the
`external-dns-cloudflare` Secret. Do not add public hostnames yet.
Rollback: delete the tunnel in the dashboard; delete the Secrets.

**3. Enable the tunnel.** PR: `replicas: 2` in `platform/networking/cloudflared/cloudflared.yaml`.
Verify: two Ready pods on different nodes, tunnel shows HEALTHY with 2 connectors.
Rollback: replicas back to 0 (nothing was routed yet, so no traffic impact).

**4. Test without touching production DNS.** Add one throwaway public hostname on the
tunnel, e.g. `tunnel-check.onemindos.dev` -> service
`https://ingress-nginx-controller.ingress-nginx.svc.cluster.local:443`, with No TLS Verify,
Origin Server Name `auth.onemindos.dev` and HTTP Host Header `auth.onemindos.dev`.
Open it: you should get the real site. (A hostname with no existing DNS record is
required: the dashboard creates a proxied CNAME and refuses to overwrite.)
Rollback: delete the hostname and its CNAME.

**5. Source-IP allowlist (blocks 4 ingresses until fixed).** Traffic through the tunnel
reaches ingress-nginx from the cloudflared **pod IP** (cluster pod CIDR `10.109.0.0/16`),
not from Cloudflare ranges. `auth/authentik-ingress`, `nats/nats-ws-ingress`,
`onetak/onetak-ingress` and `onetak/geo-services-ingress` carry
`nginx.ingress.kubernetes.io/whitelist-source-range` with Cloudflare ranges only, so they
would return 403 through the tunnel. In a follow-up PR append `10.109.0.0/16` to that
annotation on those four Ingresses (additive, the LB path keeps working). Also set
`use-forwarded-headers: "true"` and `forwarded-for-header: CF-Connecting-IP` with
`proxy-real-ip-cidr: 10.109.0.0/16` in `platform/values/ingress-nginx.yaml` (a Helm-CLI
release: needs `helm upgrade`) so logs and rate limits see real client IPs. Re-check the
CIDR first: `doctl kubernetes cluster get omos-cloud -o json | grep cluster_subnet`.
Rollback: revert the annotation; `helm upgrade` the previous values.

**6. Switch hostnames one at a time, lowest risk first:** `auth`, `matrix`/`element`, then whatever else stays public. For each hostname:

1. Add the route in the tunnel: origin `https://ingress-nginx-controller.ingress-nginx.svc.cluster.local:443`,
   No TLS Verify, Origin Server Name = hostname.
2. Change its DNS record from `A 143.244.220.155` to `CNAME <TUNNEL_ID>.cfargotunnel.com`,
   proxied. Edit the existing record in place (atomic) in the dashboard. external-dns will
   not adopt records it did not create (TXT ownership), so for the first cut-over of an
   existing name the dashboard edit is the safe route.
3. Test: page loads, login works (auth), websocket upgrades (ws), cert/HSTS OK.
4. For new hostnames after this, add to the Ingress instead (external-dns creates the CNAME):
   ```yaml
   annotations:
     dns.onemindos.dev/managed: "true"
     external-dns.alpha.kubernetes.io/target: "<TUNNEL_ID>.cfargotunnel.com"
     external-dns.alpha.kubernetes.io/cloudflare-proxied: "true"
   ```
   Prerequisite: PR with `replicaCount: 1` in `platform/argocd/platform-external-dns.yaml`
   once Secret `external-dns-cloudflare` exists. Optional first run: add `--dry-run` via
   `extraArgs` and read the logs.

Rollback per hostname: set the record back to `A 143.244.220.155` proxied (values from step
0). The LB is still up, so rollback is a DNS edit; propagation is seconds on the proxied
path. TLS certs for matrix come from cert-manager DNS-01 and are unaffected.

**7. Move private hosts to the tailnet.** For each "tailnet only" row in section 2: add a
`tailscale` Ingress (pattern in `platform/monitoring`, `platform/n8n`), verify it from a
tailnet device, then delete the public DNS record and the nginx Ingress rule. Rollback:
restore the record and rule from Git/export.

**8. Soak.** Run at least a week with LB and tunnel both live. Check `kubectl -n cloudflare logs`
for errors and Cloudflare tunnel metrics; confirm nothing still reaches `143.244.220.155`
(nginx access log).

**9. Cut-over PR (later, separate).** Remove the public paths: ingress-nginx Service
`type: ClusterIP` (helm upgrade of the CLI release, which makes DO delete the LB), delete
`tak-external` once EUDs are migrated and the federation decision is made, remove the legacy
`platform/core/cloudflared.yaml` + `config.yaml`, and delete the stale A records and the `tak`
record. Rollback: re-create the Service as `LoadBalancer` (a new IP; DO LBs are not
recoverable, so re-point DNS).

## 5. TAK clients over Tailscale

The TAK CA and client certificates do not change. Only the server address the EUD dials
changes, from `tak.onemindos.dev` to `tak.tail717ba0.ts.net`. Existing devices usually need
a changed connect string, not re-enrollment.

1. Install Tailscale on the EUD and sign in to the tailnet; confirm `tak.tail717ba0.ts.net`
   resolves (MagicDNS) and that Tailscale stays connected (always-on VPN, battery
   optimisation exempt).
2. Add the tailnet server in ATAK/iTAK: host `tak.tail717ba0.ts.net`, port 8089, protocol SSL.
   Use a data package whose preferences contain `connectString0 =
   tak.tail717ba0.ts.net:8089:ssl` (keep the old server entry until verified).
3. New devices: enroll against `tak.tail717ba0.ts.net:8446` (cert enrollment). The package
   generator `tak-portal` builds data packages from `TAK_HOST` (`platform/onetak/tak-portal.yaml`,
   currently `tak.onemindos.dev`). Change it to `tak.tail717ba0.ts.net` in the cut-over PR, not before.
4. Verify CoT flows (own marker visible on a second device; CoT appears on `tak.cot.>`).
5. WebTAK/API: `https://tak.tail717ba0.ts.net:8443` (client cert required).
6. If a client rejects the server on hostname mismatch, re-issue the TAK server certificate
   with SAN `tak.tail717ba0.ts.net` (keep the old name too). Check this on one device first.
7. After all devices are confirmed, delete the old server entry and remove `tak-external`
   in the step 9 PR.

Rollback: devices keep the old server entry, so switching back is selecting it. `tak-external`
is untouched until step 9.

Notes: the tailnet ACL (`platform/core/policy.hujson`) is currently allow-all; consider
restricting `tak` to a group for EUD users before relying on it. Mumble UDP is not forwarded
(the proxy is TCP-only); Mumble falls back to TCP. The `tak-tailnet` Service has four ports
(8089, 8443, 8446, 64738/TCP). Multi-port is supported and already used live by
`nats/omos-nats` (4222/7422/8080).

## 6. Home server (behind NAT)

The model needs no inbound ports: cloudflared and the Tailscale proxies dial out. On the
home cluster: drop `service.beta.kubernetes.io/do-loadbalancer-*` annotations (the tunnel
reaches ingress-nginx over its ClusterIP), reuse the same Applications and Secrets, create
a new tunnel token or reuse the existing one (connectors can run in both places only
during cut-over), and re-point nothing in DNS because hostnames already target the tunnel.
Images used here are multi-arch (amd64 + arm64).

## 7. Kubernetes API control-plane firewall

The control plane is currently reachable from the whole internet
(`doctl kubernetes cluster get omos-cloud` shows no firewall). Restrict it after the
tailnet API-server proxy (`apiServerProxyConfig.mode: "true"`, context
`tailscale-operator.tail717ba0.ts.net`) is confirmed working from your machines, since that
becomes the primary path. Verified flags (`doctl kubernetes cluster update --help`):

```bash
doctl kubernetes cluster update omos-cloud \
  --enable-control-plane-firewall \
  --control-plane-firewall-allowed-addresses "<YOUR_IP_1>/32,<YOUR_IP_2>/32"
```

You must supply the addresses; this document does not know them. Include the home/office
egress IPs you use for `doctl kubernetes cluster kubeconfig save`, and any CI that talks to
the public API (the in-cluster runner and Argo use the in-cluster service address, not the
public endpoint). The flag replaces the whole allowed list on each update, so always pass
the complete set. Confirm with `doctl kubernetes cluster get omos-cloud -o json` afterwards.
Test the tailnet context and `kubectl get nodes` before closing your existing session.
Rollback: `--control-plane-firewall-allowed-addresses "0.0.0.0/0"` (or disable it in the
DigitalOcean control panel / API if your IP changed and you are locked out).

## 8. Open questions

1. Is federation (9101/9102) needed with external peers, and can they join the tailnet?
2. Does any user reach CloudTAK (`onemindos.dev`) or `ws.onemindos.dev` without Tailscale?
3. Does Matrix federation use `.well-known` delegation to 443 only?
4. Should `tak.onemindos.dev` be deleted at cut-over, or kept for a temporary public path?
5. Should the tailnet ACL be tightened before TAK relies on it?
