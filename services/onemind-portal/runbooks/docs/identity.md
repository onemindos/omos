# Identity and access

## Access and dependencies

Authentik is reached at `https://onemind-auth.tail717ba0.ts.net/` and runs in namespace `auth`. Tailscale provides private connectivity; Authentik provides identity. The existing `akadmin` account is the initial admitted user. The catalog maps its signed OIDC subject through `onemind.io/oidc-sub`.

## Investigate

Distinguish a hostname/connectivity problem, a login failure, and an authorization denial. Check the workstation's tailnet connection and the identity provider page before changing Backstage.

```sh
kubectl -n auth get deployments,pods,services,ingresses
kubectl -n backstage get service authentik-mesh
kubectl -n argocd get application platform-backstage
```

The P3 identity check reads the internal readiness endpoint. A passing check does not validate the external hostname or OIDC callback. Backstage's OIDC requests use the Tailscale egress ProxyGroup; CoreDNS rewrites only the identity hostname to the egress Service. TLS must still validate the original hostname. The operator owns the Service's ExternalName.

For a callback failure, compare the canonical Backstage base URL and the Authentik application's allowed redirect URI. For an admitted user with denied access, inspect the signed-subject mapping and group membership deliberately. Do not bypass email verification, enable guest access, or map arbitrary users to administrators.

## Recovery

Obtain the existing password through the approved local credential store. Do not put it in chat, Git or logs. If the account cannot be recovered, use the identity provider's documented administrator recovery procedure with an authorized administrator. Preserve Authentik database storage and keys. Changing a hostname requires coordinated routing, base URL and redirect configuration; do not perform a partial rename.

## Verify completion

Complete the real login flow and confirm the intended user identity and permissions. Verify private HTTPS certificate validation and anonymous API denial. Readiness alone is insufficient.
