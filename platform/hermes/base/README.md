# hermes base: what the `hermes` namespace shares

Argo CD Application `platform-hermes-base` (sync wave 2, `prune: false`, `selfHeal: true`, **no
finalizer**). Added 2026-11-01 while retiring the old Hermes (`hermes-legacy`).

| Object | Why it is shared |
|---|---|
| `Namespace/hermes` | holds `onemind-hermes`, `onemind-mcp`, the old `hermes-legacy` and the volumes. `Prune=false,Delete=false` |
| `SealedSecret/hermes-dashboard-auth` | dashboard login (`envFrom`) of `onemind-hermes`; also mounted by `hermes-legacy` |
| `SealedSecret/omos-nats-creds` | `NATS_USER` / `NATS_PASSWORD` of `onemind-hermes` (direct client of the bus) |
| `SealedSecret/hermes-nats-ca` | NATS CA volume of `onemind-hermes`; also mounted by `hermes-legacy` |
| `Service/omos-ha-edge`, `omos-mbp`, `omos-s24` | Tailscale egress to the edge devices (the tool gateway's Home Assistant target uses `hermes/omos-ha-edge`) |

Every object carries `argocd.argoproj.io/sync-options: Prune=false,Delete=false`. Deleting an
Argo Application, or removing a file from this directory, therefore never deletes them. To
delete one on purpose, remove the annotation in a first commit, then the file in a second.
(A deleted SealedSecret takes its Secret with it: the Secret carries an `ownerReference` to the
SealedSecret, so Kubernetes garbage-collects the Secret.)

Not here, on purpose: `ghcr` (pull secret of `hermes-legacy` only), `cloudtak-agent-token` (mounted by
no pod; owner decision W05 in `platform/hermes/onemind-hermes/RETIREMENT-AUDIT.md`), `hermes-env`,
`hermes-leaf-conf`, `hermes-observer` and its RBAC, the `hermes-dashboard` Service and Ingress (selector:
the old pod), the Paperclip proxy pieces. They stay with `platform-legacy` and go with the old instance.
`onemind-mcp` (`platform/hermes/mcp`) is its own Application, `platform-onemind-mcp`.

## Ownership transfer (how nothing is deleted)

Facts checked live on 2026-11-01 (read-only): Argo CD v3.4.4, `application.resourceTrackingMethod`
unset, which is the v3 default, **annotation tracking**. A resource belongs to the Application
named in its `argocd.argoproj.io/tracking-id` annotation, format
`<app>:<group>/<kind>:<namespace>/<name>`. Today every object above carries
`platform-legacy:...`. An Application prunes only live objects whose tracking-id names it and
that are no longer in its desired state. Deleting the Application with the cascading finalizer
deletes everything tracked by it (except `Delete=false` objects, once the annotation is live).

The hazard in a plain move of a manifest between Applications is a race: if `platform-legacy`
syncs first, it sees the object missing from its desired state and prunes it before the new
Application has re-tagged it. The live objects do not carry `Prune=false` yet, because that
annotation only reaches the cluster when an Application applies it. So the move is done in two steps and
the objects never leave `platform-legacy`'s desired state until the new owner holds them.

**Phase 1 (this change).**

1. The manifests move to `platform/hermes/base/` (and `platform/hermes/mcp/` stays where it is).
   `platform-legacy` lists `platform/hermes/base` and `platform/hermes/mcp` as extra sources, so
   its desired state is the same set of objects as before (a source moved, not an object removed).
   Nothing is pruned, nothing is recreated. ConfigMap hash suffixes are unchanged.
2. New Applications `platform-hermes-base` and `platform-onemind-mcp` render the same files. Both
   Applications and the old one now show `SharedResourceWarning` (an Argo condition, not a failure).
3. The objects gain `Prune=false,Delete=false` (all except generated ConfigMaps, which must stay
   prunable because their names change with their content). Both owners apply the annotation, the
   tracking-id is rewritten to whichever Application applied last. A manual sync of each new Application
   (`argocd app sync platform-hermes-base platform-onemind-mcp`, or a hard refresh) makes it apply
   everything, including objects it considers in sync, which pins the tracking-id to the new owner.
   SealedSecrets: only `metadata.annotations` of the SealedSecret changes; `spec.encryptedData` and the
   template are untouched, so the controller re-unseals the identical Secret in place (no delete).

**Phase 2 (separate small change, after the verification below passes).** Remove the two
transfer-only sources from `platform/argocd/platform-legacy.yaml`. Because every object is then
tracked by its new owner, and also carries `Prune=false`, `platform-legacy` has nothing it will prune.
(Belt and braces: either condition alone prevents the delete.)

**Phase 3 (the actual retirement, owner approval).** Remove the old instance from
`platform/hermes/agent` and `bridge`, then the `hermes-data` PVC and the Application, as in
`RETIREMENT-AUDIT.md`. The `hermes` Namespace, `onemind-hermes`, `onemind-mcp`, the three Secrets and
the egress Services are no longer reachable from `platform-legacy`.

### Verification (read-only)

```sh
# 1. Every moved object is tracked by its new owner (annotation only; do not print Secret data).
kubectl --context do-nyc3-omos-cloud get ns hermes -o jsonpath='{.metadata.annotations.argocd\.argoproj\.io/tracking-id}{"\n"}'
for n in hermes-dashboard-auth omos-nats-creds hermes-nats-ca onemind-mcp-auth onemind-mcp-registry; do
  kubectl --context do-nyc3-omos-cloud -n hermes get sealedsecret $n \
    -o jsonpath='{.metadata.name}{"  "}{.metadata.annotations.argocd\.argoproj\.io/tracking-id}{"  "}{.metadata.annotations.argocd\.argoproj\.io/sync-options}{"\n"}'
done
kubectl --context do-nyc3-omos-cloud -n hermes get svc omos-ha-edge omos-mbp omos-s24 onemind-mcp \
  -o custom-columns=N:.metadata.name,TRACK:'.metadata.annotations.argocd\.argoproj\.io/tracking-id'
kubectl --context do-nyc3-omos-cloud -n hermes get deploy,networkpolicy,pdb onemind-mcp \
  -o custom-columns=K:.kind,TRACK:'.metadata.annotations.argocd\.argoproj\.io/tracking-id'
# 2. The Secrets were not recreated (creationTimestamp must predate the merge) and still have their owner.
kubectl --context do-nyc3-omos-cloud -n hermes get secret hermes-dashboard-auth omos-nats-creds hermes-nats-ca onemind-mcp-auth onemind-mcp-registry \
  -o custom-columns=N:.metadata.name,CREATED:.metadata.creationTimestamp,OWNER:.metadata.ownerReferences[0].kind
# 3. Pods untouched: same names/ages, no restarts.
kubectl --context do-nyc3-omos-cloud -n hermes get pods
# 4. Applications.
kubectl --context do-nyc3-omos-cloud -n argocd get app platform-hermes-base platform-onemind-mcp platform-legacy platform-onemind-hermes
```

Proceed to phase 2 only when every tracking-id above names `platform-hermes-base` or
`platform-onemind-mcp` and `sync-options` shows `Prune=false,Delete=false`.
