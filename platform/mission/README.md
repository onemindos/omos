# platform/mission — task service (FEATURES D2, H4)

Applied by the Argo Application `platform-mission` (`platform/argocd/platform-mission.yaml`),
**manual sync only**, into the existing namespace `onemind-core` (beside the policy service,
which BUS-DESIGN pairs it with). Source and rules: `services/mission/README.md`.

| File | What |
|---|---|
| `deployment.yaml` | One replica, `Recreate`, ecosystem pool (`onemind.io/compute-role: ecosystem`), requests 10m/32Mi, limits 300m/128Mi, non-root, read-only root filesystem |
| `configmap.yaml` | `principals.json` (NATS identity to principal, for `svc.mission.report`). No secrets |
| `network-policy.yaml` | No ingress: the service talks only over NATS (outbound) |

State: KV `mission-state` (`platform/nats/streams/kv-mission-state.yaml`, R3, history 5,
`preventDelete`) holds every record; the full history is the `MISSION` stream (`evt.mission.>`,
90 days). No PVC, so no VolSync source; NATS backups cover both. Exposure: none (no Service,
no ingress).

## Before the first sync (all owner steps)

1. **Image.** `.github/workflows/build-mission.yml` builds `ghcr.io/onemindos/mission:sha-<commit>`
   on main. Replace the placeholder tag in `deployment.yaml` with that tag and pin the digest
   from the run's `image-digest` artifact (`...:sha-<commit>@sha256:<digest>`).
2. **NATS user `mission`.** Add it to the A19 identities (omos #562,
   `platform/nats/identities/identities.yaml`) with the permissions in
   `docs/proposals/2026-10-29-mission-service-owner-proposals.md` (D2-P5), and export
   `svc.mission.get`, `svc.mission.list`, `svc.mission.report` to CLIENT and AGENTS
   (`share: true`). The policy user needs publish on `svc.mission.apply`.
3. **SealedSecret** (commit the sealed file here; never the plaintext):

   ```bash
   kubeseal --fetch-cert --controller-name=sealed-secrets-controller \
     --controller-namespace=kube-system > /tmp/pub-cert.pem
   # Password of the `mission` NATS user from your password manager, exported in your shell.
   kubectl create secret generic mission-service-runtime -n onemind-core \
     --from-literal=NATS_USER=mission \
     --from-literal=NATS_PASSWORD="$MISSION_NATS_PASSWORD" \
     --dry-run=client -o json \
     | kubeseal --cert /tmp/pub-cert.pem --format yaml --scope namespace-wide \
     > platform/mission/mission-service-runtime.sealedsecret.yaml
   ```

4. **Streams and bucket.** `CMD` (`cmd.>`), `EVT` (`evt.so.>`, `evt.drone.>`, `evt.robot.>`),
   `MISSION` (`evt.mission.>`) and KV `mission-state` exist already. The service creates two
   durable consumers (`CMD/mission-commands`, `EVT/mission-events`, deliver-new on first start).
5. **Policy service.** Sync `platform-policy` (#554) first; without it no command reaches
   this service, and the MAVLink bridge from #550 at or after `2b724093`, which accepts the
   `task.assign` `params_hash` on task-bearing commands (POL-16 / D2-P1).

6. **Arrival (D2-P2).** `MISSION_ARRIVAL` is `"false"`: a goto task completes only on the
   bridge's own `done_ok`. Set it to `"true"` only after ratifying D2-P2, because arrival
   trusts any publisher of `tel.fleet.*.*.position`.

Until these are done, leave the Application unsynced. It deletes nothing if removed (no
resources-finalizer).
