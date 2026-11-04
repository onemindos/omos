---
name: android-adb-management
description: "Audit, debloat, optimize and prepare Android devices (Samsung S-series) for ATAK using ADB through the Mac bridge; includes Knox Guard diagnosis and APK install patterns."
version: 2.0.0
author: Legacy
tags: [android, adb, samsung, debloat, optimize, tactical, atak, knox]
---

# Android ADB Device Management

Prepare and maintain Android devices (Zeus's S24 Ultra, tablets, field units) for
ATAK use: audit, remove bloatware, apply tactical settings, install ATAK and
plugins, and diagnose Samsung Knox Guard.

```
onemind-hermes (cloud) --Tailscale SSH--> Mac (ADB bridge) --USB--> Android device
```

Alternative with no Mac and no cable: the wireless ADB tunnel in `s24-app-control`.

## Reach and authority (read first)

- The cloud agent reaches the Mac and the phone only over Tailscale/SSH, and only
  while they are online (the phone also has to be plugged into the Mac for the USB
  bridge). If anything is offline, report it and stop.
- Read-only inspection (`adb devices`, `getprop`, `pm list packages`, `dumpsys`,
  the device audit script) needs no approval.
- Every uninstall or disable (debloat), settings change (optimize), APK install,
  `pm clear`, file delete, factory reset, reflash, account sign-in, or anything
  that acts on the physical device in a way that is not just reading it needs Zeus
  to write "execute" or "override" in chat for that specific action
  (SOUL rules 4 and 8). "Debloat the phone" approves proposing the package list,
  not running it; show the exact list first, then wait. Approval for one device or
  one script does not carry to the next.
- Removing a package is a deletion: prefer `pm disable-user --user 0` (reversible)
  or `pm uninstall -k --user 0` (keeps data) over anything destructive.
- Never write passwords, certificate passphrases, device IMEI/ICCID/serials, SIM
  data or account names into Git or memory. Use `<SERIAL>` style placeholders in
  notes and read live values from `adb devices`.

## Device prerequisites (Zeus does these on the phone)

1. Developer options: Settings, About phone, Software information, tap Build number 7 times.
2. USB debugging on.
3. Samsung One UI 6.1+: turn off Auto Blocker (or just "Block USB software
   commands"). "Blocked by Auto Blocker" is the number-one cause of failed ADB.
4. On connect, accept the USB debugging prompt. USB mode must be "Transferring
   Files" (MTP); other modes do not enable ADB.

## Running ADB through the Mac

The Mac's zsh breaks inline scripts with variable expansion (`ADB="adb -s $D
shell"; $ADB cmd` fails). Write the script to a temp file on the Mac, then run
it; keep `<SERIAL>` as an argument.

```bash
TSOCK=/opt/data/home/.local/state/tailscale/tailscaled.sock
ssh -o ProxyCommand="/opt/data/home/.local/bin/tailscale --socket=$TSOCK nc %h %p" \
    -i /opt/data/home/.ssh/<key> zeus@<mac-tailnet-host> \
    'export PATH="/opt/homebrew/bin:$PATH"; adb devices'
```

Use a quoted heredoc (`<< 'EOF'`) for scripts so the pod does not expand `$`.
Resolve the Mac from `tailscale status`; do not hard-code addresses.

## Scripts in this skill

| Script | Purpose | Class |
|---|---|---|
| `scripts/device_audit.sh <SERIAL>` | Read-only readiness check | read-only |
| `scripts/debloat_samsung.sh <SERIAL>` | Removes about 93 carrier/Samsung/Google packages | needs execute/override |
| `scripts/optimize_samsung.sh <SERIAL>` | Tactical settings (animations 0, GPS high accuracy, dark mode, stay awake, Bluetooth off, etc.) | needs execute/override |
| `scripts/atak_full_install.sh <SERIAL> <APK dir>` | ATAK-CIV, plugins, TakTracker | needs execute/override |

The three that change the device refuse to run unless `ZEUS_APPROVED=1` is set,
and you set it only after Zeus wrote "execute" or "override" for that action.
Copy with `scp`, then run on the Mac. Review a script before running; never run
one from an unverified source (SOUL rule 4 and the untrusted-download rule).

Keep (do not remove): Play Store (`com.android.vending`), Play Services, Google
Maps (ATAK layer dependency), WebView, Chrome, Knox Secure Folder, camera, core
Knox framework. Samsung Pay framework (`spayfw`) cannot be removed
(`DELETE_FAILED_OWNER_BLOCKED`).

## Knox Guard (carrier MDM) diagnosis, read-only

Knox Guard (`com.samsung.android.kgclient`) is a system priv-app that survives
factory reset and cannot be removed by `pm`, `disable-user` or `dpm
remove-active-admin`.

```bash
adb -s <SERIAL> shell ps -A | grep kgclient
adb -s <SERIAL> shell dumpsys device_policy | grep -E "maximumTimeToUnlock|mSuspendPersonalApps|mAlwaysOnVpnPackage|isNetworkLoggingEnabled|maximumFailedPasswordsForWipe"
adb -s <SERIAL> shell logcat -d -t 50 | grep -i "kgclient\|KnoxGuard"
```

Dormant (no process, no logs, zeros/null/false) means proceed. Active means the
fix is a reflash with unlocked (XAA) firmware or Knox Mobile Enrollment
deregistration by the seller; both are Zeus's decision and hands-on work.

## Accounts and backups

- A dedicated Google account for the tactical device is fine for Play Store and
  Play Services; do not sync personal data to it. Zeus creates and signs in to it.
- Skip the Samsung account for tactical use unless Zeus wants Find My Mobile.
- Backups are local (Smart Switch to the Mac). No cloud backup by default.

## Post-debloat readiness checklist (read-only checks)

Screen lock set; encryption on (`ro.crypto.state=encrypted`); SIM `READY`; NFC and
Bluetooth off for tactical use; GPS in `location_providers_allowed`; Private DNS
mode `opportunistic` (hostname mode breaks Tailscale and Termux mirrors);
animations 0; dark mode; stay-awake on charge. The audit script prints these.

## ATAK-CIV loadout

Details in `references/atak-loadout-inventory.md` and
`references/s24-atak-plugin-inventory.md`. Summary: install ATAK-CIV (full with
offline maps, or small), then plugins by tier (core, situational, drone, hardware),
then TakTracker (complementary, not redundant). If the DIUBLUE UAS Tool installs
first, the standard UAS Tool fails with `INSTALL_FAILED_VERSION_DOWNGRADE`; that
is fine. Install order after debloat: ATAK-CIV, TakTracker, Tailscale, GPSTest,
F-Droid, Termux, OsmAnd+.

Map sources: `references/atak-map-sources.md`. APK sources and verification:
`references/apk-download-sources.md`. Downloads need Zeus's yes (name the file,
source and size); verify checksums or signatures before installing.

## ATAK client configuration

- ATAK's own `shared_prefs` are not writable by ADB (`run-as` fails). Put config
  in `/sdcard/atak/config/prefs/` for ATAK to import.
- Pref XML must have no comments, and `locationSourcePreference` is an int
  (`<int name="locationSourcePreference" value="0"/>`, 0 = internal GPS). A string
  value or a comment makes ATAK ignore the block and show "No GNSS". "No GNSS" is
  an ATAK source setting, not a dead GPS chip; confirm with
  `dumpsys location | grep "gps provider"`.
- TAK is reachable only over the tailnet (`tak.tail717ba0.ts.net:8089`). The old
  public address is retired. `connectString0` uses commas (`host,8089,ssl`).
- Certificates and their passphrases are credentials issued by the owner. The
  agent does not generate, store or print them; ATAK still asks Zeus for the
  passphrase in its UI. Do not drive enrollment with screenshot-and-tap loops;
  give Zeus the steps.
- A PIN-locked screen blocks `am start` for third-party apps. Zeus unlocks it.

## Termux and SSH on the device

Setup, `Termux:Boot` and the pitfalls are in `s24-termux-ssh`. Key points:
Termux needs 30 to 60 seconds to bootstrap; `input text` mangles shell
metacharacters (use it only for simple commands like `sshd`); grant storage with
`appops set com.termux android:manage_external_storage allow`; Termux SSH cannot
launch other apps or read system settings.

## APK install patterns

- GitHub `/releases/latest/download/...` often returns a tiny placeholder. Use the
  API (`api.github.com/repos/OWNER/REPO/releases/latest`) and read
  `assets[].browser_download_url`.
- Session install avoids confirmation dialogs: `pm install-create`, then
  `pm install-write -S <bytes> <session> base <apk>` (the `-S` size is required on
  Android 16), then `pm install-commit`. `adb_shell` needs a 60 second or more
  transport timeout. See `references/pod-apk-session-install.md`.
- F-Droid bulk installs: `references/fdroid-bulk-install-session.md`; picks:
  `references/fdroid-onemind-recommendations.md`.
- `pm install-existing <pkg>` restores a package still on the system partition
  (Galaxy Store, Samsung Messages). Default SMS app and the Samsung Account cannot
  be set or sideloaded over ADB; Zeus does that in Settings.

## Other pitfalls

- Android 14+ blocks IMEI/ICCID reads over ADB; dial `*#06#` on the phone. Do not
  record the numbers anywhere.
- Play Store "confirm location" is an account-region issue, not GPS.
- ADB `screencap` is 1440x3120 on the S24 Ultra; vision tools may report half
  scale, so double coordinates before `input tap`.
- Android 14 forbids cross-package `am start` from Termux (not the shell user).
- A cleanup that removed Galaxy Store also blocks Good Lock (no sideload path).
- The Hermes pod's uid has no passwd entry; if SSH says "No user exists for uid",
  see `references/pod-to-s24-access-patterns.md`.

## References

`references/pod-to-s24-access-patterns.md`, `atak-loadout-inventory.md`,
`s24-atak-plugin-inventory.md`, `atak-map-sources.md`, `apk-download-sources.md`,
`fdroid-onemind-recommendations.md`, `fdroid-bulk-install-session.md`,
`pod-apk-session-install.md`. Related skills: `s24-app-control`, `s24-termux-ssh`,
`tak-server`, `find-zeus`.
