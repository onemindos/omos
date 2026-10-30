# Pod-to-S24 Device Access Patterns (Zeus S24 Ultra)

> Reference only. The cloud agent reaches the phone only over Tailscale/SSH while the Mac and phone are online. Installing, removing or reconfiguring anything on the device needs Zeus to write "execute" or "override" for that specific action (SOUL rules 4 and 8). Host and serial values are placeholders: read them live from `tailscale status` and `adb devices`.

**Session Context:** July 2026 device query revealed uid 10000 user mapping blocker in Hermes pod.

## The Problem: uid 10000 No System Mapping

The Hermes pod runs containerized as uid 10000, which has no system user mapping in the container. Attempts to SSH/SCP directly from pod to external hosts fail:

```bash
$ ssh -p 8022 <phone-tailnet-host> 'whoami'
No user exists for uid 10000
# Exit code 255 — connection attempt aborted before even reaching remote
```

This is NOT a credential issue, key issue, or network issue. It's a uid mapping issue at the container runtime level. The SSH client itself fails to initialize.

## Workaround: Route Through Mac Jump Host

**All S24 queries from pod must go through the Mac (omos-mbp, <mac-tailnet-host>) as the SSH bridge.**

### Pattern 1: Verify Device Connectivity (Pod-Native)

This works directly from pod because Tailscale is a pod service:

```bash
/opt/data/home/.local/bin/tailscale --socket=/opt/data/home/.local/state/tailscale/tailscaled.sock status | grep <phone-tailnet-name>

# Expected output (if online):
# <phone-tailnet-host>  <phone-tailnet-name>       Zeus-Delacruz@ android   -
```

Use this first to confirm device is reachable before attempting SSH.

### Pattern 2: ADB from Pod via Mac

```bash
ssh -i /opt/data/home/.ssh/id_ed25519 \
    -o StrictHostKeyChecking=no \
    -o BatchMode=yes \
    zeus@<mac-tailnet-host> \
    'export PATH="/opt/homebrew/bin:$PATH"; adb -s <SERIAL> shell "your_command"'
```

**Critical details:**
- `export PATH="/opt/homebrew/bin:$PATH"` — ADB lives at `/opt/homebrew/bin/adb` on Mac, not in default PATH
- `-o BatchMode=yes` — fails immediately instead of prompting if uid 10000 issue strikes
- Device serial: `<SERIAL>` (Zeus's S24 Ultra)

### Pattern 3: File Transfer via Mac

To push files from pod to device:

```bash
# Write to temp file on pod
echo "some content" > /tmp/payload.txt

# SCP to Mac first
scp -i /opt/data/home/.ssh/id_ed25519 \
    -o StrictHostKeyChecking=no \
    /tmp/payload.txt \
    zeus@<mac-tailnet-host>:/tmp/

# Then push to device from Mac
ssh zeus@<mac-tailnet-host> 'export PATH=/opt/homebrew/bin:$PATH; adb -s <SERIAL> push /tmp/payload.txt /sdcard/'
```

## Alternative: Direct SSH When Termux sshd Is Running

If Termux SSH is running (port 8022 open), you can:

1. First verify Tailscale on device is up:
   ```bash
   /opt/data/home/.local/bin/tailscale --socket=/opt/data/home/.local/state/tailscale/tailscaled.sock ssh <phone-tailnet-host> 'echo test' 2>&1
   ```

2. If that works but returns no output, the issue is likely sshd not running. Check with:
   ```bash
   timeout 3 bash -c "cat < /dev/null > /dev/tcp/<phone-tailnet-host>/8022" && echo "Port 8022 OPEN" || echo "Port 8022 CLOSED"
   ```

3. If port is closed, SSH daemon is down. Use ADB workaround (Pattern 2) instead.

## Why ADB Is More Reliable Than SSH from Pod

- **SSH from pod:** uid 10000 blocker, must tunnel through Mac
- **ADB via Mac:** Works reliably when device is USB-connected to Mac or Tailscale-reachable from Mac (which is always true here)
- **Fallback chain:** Tailscale connectivity (verified) → Mac SSH (works) → ADB on Mac (device USB or network) → device shell

## Confirmed Working Commands (via Mac Bridge)

Query hardware specs:
```bash
ssh zeus@<mac-tailnet-host> 'export PATH=/opt/homebrew/bin:$PATH; D=<SERIAL>; \
  adb -s $D shell getprop ro.model; \
  adb -s $D shell cat /proc/cpuinfo | head -5; \
  adb -s $D shell cat /proc/meminfo | grep MemTotal; \
  adb -s $D shell dumpsys battery | grep -E "level|capacity"; \
  adb -s $D shell df /sdcard | tail -1'
```

Query ATAK app status:
```bash
ssh zeus@<mac-tailnet-host> 'export PATH=/opt/homebrew/bin:$PATH; D=<SERIAL>; \
  adb -s $D shell pm list packages | grep -i atak; \
  adb -s $D shell ps -A | grep com.atakmap'
```

## Future Sessions: When to Use This Pattern

- Device query from pod that requires SSH or ADB
- Direct SSH times out or fails with uid 10000 error
- Termux SSH port is not responding (sshd not running)
- Need to push/pull files between pod and device
- Need to run shell commands on device without interactive Termux

## See Also

- `s24-termux-ssh` SKILL.md — Termux SSH setup and manual startup
- `android-adb-management` main SKILL.md — Full device audit and management
- the pod Tailscale socket, as in `s24-termux-ssh`
