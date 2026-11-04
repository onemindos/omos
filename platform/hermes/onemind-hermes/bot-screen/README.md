# Bot Screen image

This image adds the upstream Bot Screen Debian package list to the gateway's
existing digest-pinned Hermes v0.22.5 base. It is image-time installation, not
an apt init container, a live pod modification, or a sudo grant. Packages survive
pod recreation and new PVCs. The upstream s6 entrypoint remains responsible for
dropping the gateway to UID 10000; no desktop is started at boot.

- Install only TigerVNC, Xfce core, X11 utilities, D-Bus and fonts (the exact list
  is in `Dockerfile`). Avoid the `xfce4` metapackage/session manager, screen lockers
  and power manager.
- Stage `cua-driver` using the base image's hash-locked PM authority. The base is a
  sealed bundle: a normal `ensure()` installs into a mutable home overlay, hidden
  by the PVC on restart. Explicitly stage and promote the verified driver entry
  and its fact into `/opt/hermes/tools`, retaining all existing package facts.
- Reuse the base's pinned **full Chromium**, not another browser download or a
  `chrome-headless-shell`. The offline smoke test must open it on the real screen.
- Hermes already bridges RFB via its authenticated `/api/display/ws` route, using
  single-use tickets and the profile's private `0600` Unix socket. No websockify,
  TCP VNC listener, new Service/Ingress port, or NetworkPolicy exception is needed.
  Each profile has its own screen/cookies, but profiles sharing the OS user are
  **not security boundaries**. The human-control lease is not OS isolation.

Authoritative references:
[Bot Screen](https://hermes-agent.nousresearch.com/docs/user-guide/features/bot-screen),
and the pinned image's `tools/bot_desktop/{runtime.py,browser.py,launcher.sh}` and
`hermes_cli/web_routers/display.py`.

## Build and gate

The existing `.github/workflows/build-hermes-runtime.yml` is **manual only**.
Dispatch on the reviewed branch/commit with input `variant: bot-screen`. Its
original `runtime-fix` default and correction tests are preserved; selecting
Bot Screen does not activate the unrelated native-1Password runtime correction.

The build loads a local amd64 image, then runs `smoke_test.py` as UID 10000 with
all capabilities dropped, no-new-privileges, no external network, a 3 GiB cgroup
limit and a fresh temporary home outside any live PVC. The gate checks:

1. All upstream required desktop binaries resolve.
2. PM can resolve and execute the pinned cua-driver without lazy downloads.
3. The actual authenticated display WebSocket router imports.
4. Hermes starts Xfce/Xvnc; `xdpyinfo` reaches the published display.
5. The private `0600` RFB Unix socket returns a real protocol greeting.
6. The dock's headed Chromium argv opens `about:blank` and becomes CDP-ready.

Only after this gate passes does the workflow authenticate to GHCR and push the
**same tested image**. Its artifact `hermes-runtime-image.txt` contains the
registry-verified digest. Never use the local Docker config ID as a registry
digest. Pin that artifact's exact `ghcr.io/onemindos/hermes-agent@sha256:...`
reference in the main gateway container of `../deployment.yaml`; init containers
can retain their existing base since they do not host screens.

Cheap local checks:

```sh
python platform/hermes/onemind-hermes/bot-screen/test_contract.py -v
python -m py_compile platform/hermes/onemind-hermes/bot-screen/smoke_test.py
git diff --check
make validate
```

## Capacity and owner acceptance after merge

The existing gateway requests (250m CPU / 1Gi memory) and limits (2 CPU / 3Gi)
remain unchanged. Upstream measured roughly 220 MB for the desktop plus
0.5–1 GB for a headed browser; concurrent profile screens can exceed this pod's
budget. The default 1536 MB free-memory start gate and idle stop remain enabled;
do not disable them to hide a capacity failure. Monitor headroom/OOMs and propose
capacity changes separately if needed. This PR adds no workload, PVC or secrets.

Source/image smoke success is not live acceptance. After an owner-authorized
merge and Argo rollout, verify the exact running image digest and healthy gateway,
then open Desktop → Bots → Screen. Verify start, visible desktop/browser, human
takeover, hand-back and persistence across a pod recreation. Computer Use must be
enabled in the intended bot's toolsets; installing a driver does not change the
bot's tool permissions. This PR does not edit live or seeded config.

Rollback is a reviewed Git revert of the gateway image reference to its prior
upstream digest. Do not delete profile state/cookies, force a pod rollout or
apply manifests manually. The previously inspected Argo app may already report
unrelated ConfigMap drift; diagnose separately, not by altering it here.
