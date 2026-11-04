#!/bin/bash
# Android Device Full Audit Script — 40-Point Readiness Check
# Usage: bash device_audit.sh <ADB_SERIAL>
# Run on the Mac (ADB bridge) via Tailscale SSH from the onemind-hermes pod, after Zeus approves (ZEUS_APPROVED=1).
# Output is plain text, designed for parsing or direct reading.

export PATH="/opt/homebrew/bin:$PATH"

if [ -z "$1" ]; then
  echo "Usage: $0 <ADB_SERIAL>"
  echo "Example: $0 <SERIAL>"
  exit 1
fi

D="$1"

echo "=== 1. DEVICE IDENTITY ==="
adb -s $D shell getprop ro.product.model
adb -s $D shell getprop ro.product.device
adb -s $D shell getprop ro.build.version.release
adb -s $D shell getprop ro.build.version.security_patch
adb -s $D shell getprop ro.product.carrier
adb -s $D shell getprop ro.boot.carrier

echo "=== 2. STORAGE ==="
adb -s $D shell df -h /data | tail -1

echo "=== 3. RAM ==="
adb -s $D shell cat /proc/meminfo | head -3

echo "=== 4. BATTERY ==="
adb -s $D shell dumpsys battery | grep -E "level|status|temperature|health"

echo "=== 5. SIM STATUS ==="
adb -s $D shell getprop gsm.sim.state
adb -s $D shell getprop gsm.operator.alpha
adb -s $D shell getprop gsm.operator.numeric

echo "=== 6. WIFI ==="
adb -s $D shell dumpsys wifi 2>&1 | grep -E "SSID|mNetworkInfo" | head -3
adb -s $D shell ip addr show wlan0 2>&1 | grep "inet "

echo "=== 7. LOCATION ==="
adb -s $D shell settings get secure location_mode
adb -s $D shell settings get secure location_providers_allowed
adb -s $D shell dumpsys location 2>&1 | grep "last location" | head -3

echo "=== 8. SCREEN/DISPLAY ==="
adb -s $D shell settings get system screen_off_timeout
adb -s $D shell settings get system screen_brightness
adb -s $D shell settings get system screen_brightness_mode
adb -s $D shell settings get secure ui_night_mode
adb -s $D shell wm size
adb -s $D shell wm density

echo "=== 9. ANIMATIONS ==="
adb -s $D shell settings get global window_animation_scale
adb -s $D shell settings get global transition_animation_scale
adb -s $D shell settings get global animator_duration_scale

echo "=== 10. DNS ==="
adb -s $D shell settings get global private_dns_mode
adb -s $D shell settings get global private_dns_specifier

echo "=== 11. DEVELOPER OPTIONS ==="
adb -s $D shell settings get global development_settings_enabled
adb -s $D shell settings get global stay_awake_on_plugged_in
adb -s $D shell settings get global package_verifier_enable

echo "=== 12. ACCOUNTS ==="
adb -s $D shell dumpsys account 2>&1 | grep -E "Account.*type" | head -10

echo "=== 13. DEVICE ADMINS ==="
adb -s $D shell dumpsys device_policy 2>&1 | grep -A1 "Enabled Device Admins"

echo "=== 14. USB DEBUGGING ==="
adb -s $D shell settings get global adb_enabled

echo "=== 15. LOCK SCREEN ==="
adb -s $D shell dumpsys trust 2>&1 | grep -E "isTrustUsuallyManaged|trusted" | head -3

echo "=== 16. RUNNING SERVICES ==="
adb -s $D shell dumpsys activity services 2>&1 | grep -c "ServiceRecord"

echo "=== 17. PACKAGE COUNT ==="
adb -s $D shell pm list packages | wc -l
adb -s $D shell pm list packages -3 | wc -l

echo "=== 18. 3RD PARTY APPS ==="
adb -s $D shell pm list packages -3

echo "=== 19. CRITICAL APPS ==="
for pkg in com.android.vending com.google.android.gms com.google.android.apps.maps \
           com.google.android.webview com.android.chrome \
           com.samsung.knox.securefolder com.onepassword.android; do
  result=$(adb -s $D shell pm list packages $pkg 2>&1)
  if [ -n "$result" ]; then
    echo "  [PRESENT] $pkg"
  else
    echo "  [MISSING] $pkg"
  fi
done

echo "=== 20. BLUETOOTH ==="
adb -s $D shell settings get global bluetooth_on

echo "=== 21. NFC ==="
adb -s $D shell dumpsys nfc 2>&1 | grep -E "mState|isEnabled" | head -2

echo "=== 22. AIRPLANE MODE ==="
adb -s $D shell settings get global airplane_mode_on

echo "=== 23. DO NOT DISTURB ==="
adb -s $D shell settings get global zen_mode

echo "=== 24. AUTO ROTATE ==="
adb -s $D shell settings get system accelerometer_rotation

echo "=== 25. GNSS RAW ==="
adb -s $D shell settings get secure gnss_raw_meas_full_tracking 2>/dev/null

echo "=== 26. KNOX GUARD RUNNING ==="
adb -s $D shell ps -A 2>&1 | grep -c kgclient

echo "=== 27. ENCRYPTION ==="
adb -s $D shell getprop ro.crypto.state

echo "=== 28. UPTIME ==="
adb -s $D shell uptime

echo "=== 29. BLOAT CHECKS ==="
echo "Facebook: $(adb -s $D shell pm list packages com.facebook 2>&1)"
echo "Bixby: $(adb -s $D shell pm list packages com.samsung.android.bixby 2>&1)"
echo "Samsung Pay: $(adb -s $D shell pm list packages com.samsung.android.spay 2>&1)"
echo "Spotify: $(adb -s $D shell pm list packages com.spotify 2>&1)"
echo "Microsoft: $(adb -s $D shell pm list packages com.microsoft 2>&1)"
echo "S Browser: $(adb -s $D shell pm list packages com.sec.android.app.sbrowser 2>&1)"
echo "Kids Home: $(adb -s $D shell pm list packages com.sec.android.app.kidshome 2>&1)"
echo "Google Photos: $(adb -s $D shell pm list packages com.google.android.apps.photos 2>&1)"
echo "YouTube: $(adb -s $D shell pm list packages com.google.android.youtube 2>&1)"

echo "=== 30. KNOX GUARD POLICIES ==="
adb -s $D shell dumpsys device_policy 2>&1 | grep -E "maximumTimeToUnlock|mSuspendPersonalApps|mAlwaysOnVpnPackage|isNetworkLoggingEnabled"

echo "DONE"
