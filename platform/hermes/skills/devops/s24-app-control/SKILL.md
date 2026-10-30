---
name: s24-app-control
description: "Launch apps, take screenshots and send input on Zeus's S24 Ultra through an ADB tunnel over Tailscale."
version: 2.0.0
author: Legacy
tags: [phone, android, apps, s24, adb, screenshot, pipeline]
---

# S24 App Control

Control Zeus's phone from the cloud pod without the Mac: Tailscale to the phone,
an SSH tunnel to its wireless ADB port, then a Python `adb-shell` client running
as the Android shell user.

```
pod --Tailscale--> phone sshd (Termux, 8022) --SSH tunnel--> localhost:<tunnel port> -> phone adb (5555)
```

## Reach and authority (read first)

- This works only while the phone is online, on Tailscale (Always-on VPN), with
  Wireless Debugging on and Termux `sshd` running. Wireless Debugging must stay
  on; a phone reboot turns it off and only Zeus can turn it back on. When the
  tunnel is down, report it. Do not try to work around it.
- Read-only actions are fine without approval: status, battery, package list,
  `dumpsys`, and screenshots of non-sensitive screens. Take a screenshot only
  when it is needed for the task, keep it in the pod scratch directory, and do
  not share it outside the owner's channel.
- These need Zeus to write "execute" or "override" for that specific action
  (SOUL rules 4 and 8):
  - anything that touches the physical world (camera capture, microphone, flashlight,
    calls, SMS, payments, unlocking, power or radio toggles, driving a drone or
    other device through an app);
  - sending messages from the phone's apps, installing or removing apps, changing
    system settings, clearing app data, deleting files.
  - Launching an app or tapping through a screen is not exempt if the result is one
    of the above. A working tunnel is not authority.
- The ADB key is the Mac's key (the phone trusts it). Do not regenerate it, copy
  it into Git, or print it.

## Components (names only; scripts live on the pod volume or in Git)

- A guardian cron checks three layers and self-heals each: Tailscale, the SSH
  tunnel, the ADB connection. State and log files sit under the pod `logs` dir.
  It alerts on repeated failure. Define this in the Hermes cron config in Git;
  the alert goes to the configured delivery channel only.
- A control script supports `status`, `screenshot`, `launch <alias or package/activity>`,
  `shell "<cmd>"` and `camera`. Aliases include geocam, atak, camera, maps,
  chrome, settings, termux, shizuku.
- Use SSH for file work in Termux and `termux-api`; use ADB for app launch,
  screenshots, input, settings and package management.

## Direct Python ADB

```python
from adb_shell.adb_device import AdbDeviceTcp
from adb_shell.auth.sign_pythonrsa import PythonRSASigner
# load the private/public key pair from the pod ADB key directory
signer = PythonRSASigner(pub, priv)
d = AdbDeviceTcp('127.0.0.1', <tunnel-port>, default_transport_timeout_s=60)
d.connect(rsa_keys=[signer], auth_timeout_s=10)
d.shell('dumpsys battery')
d.shell('screencap -p /sdcard/capture.png')
d.pull('/sdcard/capture.png', '<pod-scratch>/screen.png')   # two args
d.shell('rm /sdcard/capture.png')                            # your own temp file only
d.close()
```

Transport timeout of 60 seconds or more is needed for `pm install`.

## Launch and input cheat sheet

```
am start -n com.partech.geocamera/.MainActivity          # GeoCam
am start -n com.atakmap.app.civ/.ATAKActivity            # ATAK-CIV
am start -n com.android.chrome/com.google.android.apps.chrome.Main
am start -a android.intent.action.VIEW -d "https://example.org"
input keyevent 3 | 4 | 26     # home | back | power
input tap X Y | input swipe X1 Y1 X2 Y2 300
```

S24 Ultra screen is 1440x3120. Vision tools often report half-scale coordinates;
multiply by 2 before `input tap`. Do not drive ATAK certificate enrollment by tap
automation; give Zeus clear manual steps instead.

## Capabilities

| Operation | Path | Works |
|---|---|---|
| App launch, screenshot, input, settings, files | ADB | yes |
| Camera view | ADB screenshot of GeoCam | yes |
| `termux-camera-photo` | SSH | no (0-byte output on Android 16) |
| `am start` from Termux SSH | SSH | no (blocked cross-package launch) |
| rish via Shizuku | SSH | not running; ADB gives the same shell user |

## Pitfalls

1. Wireless Debugging off means no remote fix; ask Zeus.
2. Tunnel drops on network change or pod restart; the guardian rebuilds it within minutes.
3. A PIN-locked screen blocks `am start` of third-party activities. Zeus unlocks it.
4. `pull()` needs both arguments.
