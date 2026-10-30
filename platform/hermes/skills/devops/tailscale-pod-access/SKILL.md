---
name: tailscale-pod-access
description: "Use to reach a machine on the OneMind tailnet (Mac, S24, Windows, Home Assistant host) from the agent pod, and to diagnose why that fails."
version: 2.0.0
author: Legacy
metadata:
  hermes:
    tags: [tailscale, ssh, networking, tailnet, egress]
    related_skills: [doks-cluster-ops, s24-termux-ssh, android-adb-management, mac-repo-hygiene, find-zeus]
triggers:
  - "SSH to the Mac"
  - "reach a tailnet device"
  - "tailnet egress"
  - "tailscale is down"
---

# Tailscale access from the agent pod

Rewritten from the Legacy-pod version. That one SSHed from the pod to personal
machines with a pinned key, a userspace `tailscaled` and a hard-coded address
list. Keep the working parts; drop the rest.

## Read this first: what `onemind-hermes` can reach today

- The `onemind-hermes` pod runs **no `tailscaled`**. It reaches the tailnet only
  through the cluster's Tailscale operator **egress** Services in namespace
  `hermes` (`omos-mbp` :22, `omos-s24` :8022, `omos-ha-edge` :8123), defined in
  `platform/hermes/agent/egress.yaml`.
- Its NetworkPolicy (`platform/hermes/onemind-hermes/network.yaml`) allows egress
  only to DNS, LiteLLM, Matrix Synapse and public HTTPS (443). **Connections to
  the egress Services are blocked** until the owner adds an egress rule for them
  in Git. Test with a short `timeout 5 nc -vz omos-mbp.hermes.svc.cluster.local 22`;
  a timeout means the policy, not the device.
- The pod also has no SSH key, no cluster ServiceAccount token and no kubeconfig.
  A key for the target machines is a credential: it comes from a SealedSecret the
  owner provisions, never from this skill, and never goes into a reply or memory.
- The tailnet policy (`platform/core/policy.hujson`) separately limits each
  target and port.

If any of these is missing, say so and propose the Git change (egress rule, key
SealedSecret). Do not build a workaround such as a userspace `tailscaled`, a
tunnel through the Mac or an exit node.

## Reaching a device, once access is provisioned

```bash
# Hostnames are the egress Service names; they resolve only inside the cluster.
ssh -o BatchMode=yes -o ConnectTimeout=10 zeus@omos-mbp.hermes.svc.cluster.local '<command>'
ssh -p 8022 -o BatchMode=yes zeus@omos-s24.hermes.svc.cluster.local '<command>'     # Termux must be running
curl -s --max-time 10 http://omos-ha-edge.hermes.svc.cluster.local:8123/            # HA; prefer the bus (ha-nats-bridge)
```

- Pin the host key once (`ssh-keyscan` over the same path) and keep it in the
  pod's `known_hosts` rather than disabling host-key checking.
- macOS non-interactive SSH has a minimal `PATH`: prefix commands with
  `export PATH=/opt/homebrew/bin:/usr/local/bin:$PATH`.
- Windows machines need OpenSSH Server enabled by the owner on the machine.
- iPad, TV and phones without Termux have no SSH; there is nothing to reach.
- Anything that changes a device (install, delete, reboot, settings) is the
  owner's call: see the "execute"/"override" rule in the soul, and the device
  skills (`android-adb-management`, `s24-app-control`, `mac-disk-hygiene`).

## Adding a device (Git, not `kubectl apply`)

Add an `ExternalName` Service to `platform/hermes/agent/egress.yaml` with the
label `tailscale.com/proxy-group: egress` and the annotation
`tailscale.com/tailnet-ip`, allow the target and port in `platform/core/policy.hujson`
and open the pod's NetworkPolicy to it, in one PR. Argo applies it. The operator
fills in `externalName`.

If egress breaks: `kubectl -n hermes get proxygroup egress` and
`kubectl -n tailscale get pods | grep egress` (read-only), then check the
operator's OAuth secret via the runbook, not by hand.

## Diagnosing a peer (read-only)

On a host that runs `tailscale` (the owner's Mac, not this pod):
`tailscale status --json` shows per peer `Online`, `CurAddr` (empty with a
`Relay` set means DERP only), `Active` and `sshHostKeys`.

- `sshHostKeys` empty: the peer does not publish Tailscale SSH host keys. The
  "host key verification failed" message is misleading; the peer has to enable
  Remote Login (macOS: System Settings, General, Sharing) or Tailscale SSH
  (`tailscale set --ssh`, which changes only that setting; `tailscale up --ssh`
  demands every non-default flag again). The cache can be stale; just try.
- `Online` and `Active` true but SSH fails: the problem is SSH configuration,
  not connectivity.
- Taildrop (`tailscale file cp <file> <peer>:`) delivers a file even when SSH is
  blocked, as long as the peer is online. The owner has to accept it.

## Pitfalls

- Host aliases from `~/.ssh/config` do not resolve under a non-interactive
  `ProxyCommand`; use the full egress hostname.
- When a command is given a `ProxyCommand` that runs a binary, use the binary's
  absolute path: the SSH subprocess does not inherit `PATH`.
- Background processes started over SSH die with the session unless started with
  `setsid ... </dev/null &`.
- A throwaway preview server on the Mac (venv inside the repo, `nohup`, a free
  port found with `lsof -iTCP -sTCP:LISTEN`, killed and the venv removed
  afterwards) leaves no system footprint. Starting and stopping it is a change on
  the owner's machine: say what you will run first. Removing the venv or repo
  files is a deletion and needs "execute" or "override".
- iOS devices on the Mac's USB: `libimobiledevice` on the Mac (`idevice_id -l`,
  `ideviceinfo`, `ideviceinstaller list --user`, `idevicebackup2 backup`) works
  over SSH, but needs the owner present for trust prompts and the Developer Disk
  Image for screenshots. Backups and wipes are the owner's call.
