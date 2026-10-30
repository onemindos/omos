# Pod→Phone Session-Based APK Install

When the pod has a direct ADB tunnel to the phone (via `s24-app-control` skill) and uses the Python `adb_shell` library, the standard `pm install` command often hangs due to Play Protect / Samsung Auto Blocker confirmation dialogs. The session-based install method avoids this.

## The Problem

```python
# This hangs — Android shows a confirmation dialog that blocks the shell
d.shell('pm install -r /data/local/tmp/app.apk', timeout_s=60)
# → TcpTimeoutException: Reading from localhost:15555 timed out
```

## The Fix: Session-Based Install

```python
import re
from adb_shell.adb_device import AdbDeviceTcp
from adb_shell.auth.sign_pythonrsa import PythonRSASigner

# CRITICAL: transport timeout must be 60s+ for install ops
d = AdbDeviceTcp('localhost', 15555, default_transport_timeout_s=60)
d.connect(rsa_keys=[signer], auth_timeout_s=10)

# Step 1: Create install session
out = d.shell('pm install-create -r -g 2>&1')
sid = re.search(r'\[(\d+)\]', out).group(1)

# Step 2: Write APK — -S SIZE is REQUIRED on Android 16
size = 21621788  # bytes, get via stat -c %s on device
d.shell(f'pm install-write -S {size} {sid} base /data/local/tmp/app.apk')

# Step 3: Commit
result = d.shell(f'pm install-commit {sid} 2>&1', timeout_s=30)
# Success → "Success"
# INSTALL_FAILED_INVALID_APK → bad APK (e.g., downloaded HTML instead of APK)
# INSTALL_PARSE_FAILED_NOT_APK → corrupted file or wrong format
```

## Downloading APKs via ADB Shell

**GitHub releases API** (never use `/latest/download/` — it redirects to a 9-byte HTML page):

```python
# FROM THE POD — query GitHub API for real asset URL
import json, urllib.request

url = "https://api.github.com/repos/Trixarian/NetherSX2-patch/releases/latest"
with urllib.request.urlopen(url) as r:
    data = json.loads(r.read())
for asset in data['assets']:
    if asset['name'].endswith('.apk'):
        real_url = asset['browser_download_url']  # ← use THIS, not /latest/download/
        break

# Then download to phone via ADB shell
d.shell(f'curl -L -o /data/local/tmp/app.apk "{real_url}"')
```

**CDN-protected websites** (Bunny Shield, Cloudflare) — curl gets JS challenge pages, not APKs. Fall back to GitHub releases or Play Store.

## Play Store Install + Auto-Tap

When Play Store is the only source (Dolphin, etc.), open the store page and auto-tap the Install button:

```python
# Open Play Store to app page
d.shell('am start -a android.intent.action.VIEW -d "market://details?id=PACKAGE_ID" com.android.vending')
time.sleep(5)  # Wait for page load

# Auto-tap Install button — S24 Ultra positions (1440×3120)
for x, y in [(900, 1500), (900, 1400), (850, 1450), (950, 1450)]:
    d.shell(f'input tap {x} {y}')
    time.sleep(0.3)

# Verify
time.sleep(5)
result = d.shell('pm list packages PACKAGE_ID')
print('Installed:' if result.strip() else 'NOT installed')
```

## adb_shell Timeout Double-Trap

The `adb_shell` library has two separate timeouts:

| Parameter | Controls | Default | For install |
|-----------|----------|---------|-------------|
| `default_transport_timeout_s` | All TCP reads (constructor) | 9s | Set to **60+** |
| `timeout_s` | Shell command execution | None | Set to **30+** |

If either is too short, you get `TcpTimeoutException`. The transport timeout is set in the constructor; it can't be changed per-command.
