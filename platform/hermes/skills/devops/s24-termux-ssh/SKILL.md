---
name: s24-termux-ssh
description: "SSH access to Zeus's Samsung S24 Ultra through Termux, direct over Tailscale or via the Mac as jump host."
version: 2.0.0
author: Legacy
tags: [phone, ssh, termux, android, remote-access]
---

# S24 Ultra Termux SSH Access

Zeus's phone is a remote endpoint (not a Hermes host). Termux runs `sshd` on port
8022 as a normal app user, so SSH gives file work, Python/Node scripts and
`termux-api` calls. It cannot read system settings or launch other apps; use ADB
for that (`s24-app-control`, `android-adb-management`).

## Reach and authority (read first)

- The cloud agent reaches the phone only over Tailscale/SSH, and only when the
  phone is online with Tailscale up and `sshd` running. When it is offline, report
  that and stop; do not loop on retries or try to wake it.
- Read-only commands (`whoami`, `uptime`, `pgrep sshd`, `df`) need no approval.
- Any delete (files, packages, app data) or physical-world action (camera,
  microphone, notifications to people, calls, SMS, toggling radios, lock or power
  state, anything that changes what the phone does outside software) needs Zeus to
  write "execute" or "override" in chat for that specific action (SOUL rules 4 and
  8). Having a working SSH session is not authority.

## Connection

Look values up live; do not trust remembered ones (SOUL rule 6).

- Host: the phone's Tailscale MagicDNS name or IP from `tailscale status`
  (`<phone-tailnet-host>`). The LAN address is DHCP and changes; use it only from
  the same Wi-Fi.
- Port 8022. No username is needed (Termux accepts the app user).
- Auth: SSH key. The pod identity and the Mac identity are authorized on the
  phone. Password login exists as a fallback; ask Zeus for it, never store it.
- Phone: aarch64, Android, Termux with bash, Python, Node, Git.

### Direct from the pod (preferred)

```bash
TS_SOCK=/opt/data/home/.local/state/tailscale/tailscaled.sock
ssh -o "ProxyCommand=/opt/data/home/.local/bin/tailscale --socket=$TS_SOCK nc %h %p" \
    -i /opt/data/home/.ssh/<key> -p 8022 <phone-tailnet-host> 'COMMAND'
```

### Via the Mac as jump host (phone's Tailscale down but on the same LAN)

Same ProxyCommand to `zeus@<mac-tailnet-host>`, then
`ssh -p 8022 <phone-lan-ip> "COMMAND"` from there. Works only if the Mac is online.

### Via ADB (SSH down, phone plugged into the Mac)

`ssh zeus@<mac-tailnet-host> 'export PATH=/opt/homebrew/bin:$PATH; adb -s <SERIAL> shell "COMMAND"'`

## Auto-start

Termux:Boot (`com.termux.boot`) runs `~/.termux/boot/start-sshd.sh`, which does
`termux-wake-lock` then `sshd`. Battery optimization is disabled for `com.termux`
and `com.termux.boot`. Always-on VPN for Tailscale cannot be set over SSH or ADB;
Zeus toggles it in Settings, VPN, Tailscale, Always-on VPN.

## Health checks

Keepalive and health crons (tailscale mesh keepalive, S24 health monitor) belong in
the Hermes cron config in Git, not hand-made scripts. They only read state:
tailscaled alive, phone online, `ssh ... 'echo OK'` works. Tailscale auth keys
expire; ask Zeus before any key is regenerated.

## Pitfalls

1. Tailscale plus Private DNS in "hostname" mode breaks package mirrors (Android
   marks the DNS broken). Private DNS must be "Automatic". Zeus changes this on
   the phone.
2. ADB `input text` cannot type shell metacharacters; use SSH for complex commands.
3. `run-as com.termux` fails (not debuggable). Use SSH, or push to `/sdcard` and
   move the file from inside Termux.
4. After a partial `pkg install`, "cannot locate symbol" OpenSSL errors are
   fixed by `DEBIAN_FRONTEND=noninteractive apt -o Dpkg::Options::="--force-confnew" full-upgrade -y`.
5. `uv` does not support Android; use pip. Hermes cannot install natively
   (psutil fails to build).
6. `settings get/put` over SSH fails (needs INTERACT_ACROSS_USERS). Use ADB.
7. Do not write device passwords, keys or authorized_keys contents into Git; the
   pod and Mac public keys are enough to name, never to paste.

## Verification

```bash
ssh ... -p 8022 <phone-tailnet-host> 'echo OK; whoami; uptime'
ssh ... -p 8022 <phone-tailnet-host> 'pgrep sshd'
```
