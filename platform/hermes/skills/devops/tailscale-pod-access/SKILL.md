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

## Read this first: what `onemind-hermes` can reach

- The `onemind-hermes` pod runs **no `tailscaled`**. It reaches the tailnet only
  through the Tailscale operator **egress** Services in namespace `hermes`
  (`platform/hermes/base/egress.yaml`), each a proxy pod for one device:

  | Device | Service (`<name>.hermes.svc.cluster.local`) | Port | What listens there |
  |---|---|---|---|
  | MacBook Pro | `omos-mbp` | 22 | SSH (Remote Login). ADB for the phone runs through this Mac |
  | S24 Ultra | `omos-s24` | 8022 | Termux `sshd` (Termux must be running). Wireless ADB is tunnelled over this SSH, no separate port |
  | Home Assistant box | `omos-ha-edge` | 8123 | REST API; box offline since about 2026-09-19 |
  | Windows laptop | `omos-dell-zeus` | 22 | OpenSSH Server, only if the owner enabled it |

- Its NetworkPolicy (`platform/hermes/onemind-hermes/network.yaml`) opens exactly
  those device/port pairs, plus DNS, the tool gateway, NATS, Matrix and public
  HTTPS. Nothing else on the tailnet is reachable: no other port on these devices
  (not VNC 5900, not RDP) and no other device. The tailnet policy
  (`platform/core/policy.hujson`) allows the same ports.
- Devices with no Service (Dell edge, Intel box, TV, iPad) have been offline for
  months or offer nothing to reach. Ask the owner before expecting them.
- A probe result tells you which layer: connected = path and device are fine;
  refused = path fine, service on the device is off; timeout = proxy, tailnet
  policy or the device is offline (the proxy accepts the connection only when the
  device answers). Test with `timeout 5 nc -vz omos-mbp.hermes.svc.cluster.local 22`
  or `python3 scripts/mesh-health.py probe` (tailscale-continuous-mesh-access).
- The pod has **no SSH key, no ADB key, no Home Assistant token**, no cluster
  ServiceAccount token and no kubeconfig. Those are credentials: they come from a
  SealedSecret the owner provisions, never from this skill, and never go into a
  reply or memory. Per device the owner must supply: an SSH key whose public half
  is in the device's `authorized_keys` (Mac, S24 Termux, Windows OpenSSH), the
  one-time ADB pairing and Wireless Debugging on (S24), and a long-lived access
  token (Home Assistant).

If a credential or a path is missing, say so and propose the Git change (egress
rule, key SealedSecret). Do not build a workaround such as a userspace `tailscaled`, a
tunnel through the Mac or an exit node.

## Reaching a device, once access is provisioned

```bash
# Hostnames are the egress Service names; they resolve only inside the cluster.
ssh -o BatchMode=yes -o ConnectTimeout=10 zeus@omos-mbp.hermes.svc.cluster.local '<command>'
ssh -p 8022 -o BatchMode=yes omos-s24.hermes.svc.cluster.local '<command>'     # Termux must be running
curl -s --max-time 10 http://omos-ha-edge.hermes.svc.cluster.local:8123/            # HA; prefer the bus (ha-nats-bridge)
ssh -o BatchMode=yes <windows-user>@omos-dell-zeus.hermes.svc.cluster.local 'hostname'  # needs OpenSSH Server on the laptop
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

Add an `ExternalName` Service to `platform/hermes/base/egress.yaml` with the
label `tailscale.com/proxy-group: egress` and the annotation
`tailscale.com/tailnet-ip`, allow the target and port in `platform/core/policy.hujson`
and add a rule to the pod's NetworkPolicy for the proxy pod
(`tailscale.com/parent-resource: <Service name>`, that one port), in one PR. Argo applies it. The operator
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
