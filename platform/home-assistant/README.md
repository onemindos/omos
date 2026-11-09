# Home Assistant (cloud)

Home Assistant Container, run by Argo CD (`platform-home-assistant`). It is the
single place to view everything. The real devices stay on the home Home
Assistant box, which remains the device hub; this instance only mirrors it.

- URL: https://ha.tail717ba0.ts.net (tailnet only, Tailscale ingress)
- Image: `ghcr.io/home-assistant/home-assistant:2026.11.0`, pinned by digest
- Storage: PVC `home-assistant-config` (5 Gi, `do-block-storage-retain`) at `/config`
- Node: pinned to the `onemind-hub-core-v2` pool (the ecosystem node has no free volume slots)
- Frigate and LLM Vision are not part of this; they come later.

## First boot

1. Wait for the pod to be Ready (the first start can take several minutes).
2. Open the URL and create the owner account in the UI. No secret is needed in Git.
3. HACS 2.0.5 is already in `/config/custom_components/hacs` (installed once by the
   init container, checksum-verified). Add the integration under Settings, Devices
   and services, Add integration, HACS, and complete its GitHub device login.
4. In HACS install **Remote Home Assistant**, restart, then link the home box by
   its Tailscale address.
5. The pod's network policy allows only public HTTPS egress today. Reaching the
   home box over the tailnet needs a reviewed egress rule and a Tailscale egress
   Service (see `platform/hermes` for the existing egress pattern). Add it in the
   same change that links the home box.

## Notes

- The init container writes `configuration.yaml` only if it does not exist, with
  `http.use_x_forwarded_for` and `trusted_proxies: 10.109.0.0/16` (the Cilium
  pod network, which contains the Tailscale ingress proxies). After first boot
  the file belongs to the UI and is never overwritten.
- The container runs as root, as the official image requires. The namespace is
  therefore `baseline` for Pod Security instead of `restricted`.
- Backup: `backup.yaml` holds a VolSync source that is **paused** until the
  restic repository Secret `home-assistant-config-restic` is sealed and added to
  `kustomization.yaml` (see the TODO there). Snapshots are crash-consistent.
- Do not delete the PVC or namespace without a fresh snapshot and approval.

## HTTP settings live in `.storage/http` (2026.11+)

Home Assistant no longer reads `http:` from `configuration.yaml`. A YAML
block is migrated once as an unconfirmed "pending" trial that reverts after
five minutes unless promoted in the UI, and through the Tailscale proxy the
UI is unreachable until `trusted_proxies` is active. The init container
therefore seeds `.storage/http` directly with `use_x_forwarded_for` and
`trusted_proxies: 10.109.0.0/16` as the stable config (fresh volumes only).
Change proxy settings later in the UI and confirm them within five minutes.
The first live volume was fixed by hand on 2026-11-08 (backup at
`.storage/http.bak-2026-11-08`).
