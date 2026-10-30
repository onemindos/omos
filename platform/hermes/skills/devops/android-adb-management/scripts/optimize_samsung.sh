#!/bin/bash
# Samsung System Optimization Script — Tactical Performance Tuning
# Applies performance, privacy, and tactical settings via ADB.
# Usage: bash optimize_samsung.sh <ADB_SERIAL>
# Run on the Mac (ADB bridge) via Tailscale SSH from the onemind-hermes pod, after Zeus approves (ZEUS_APPROVED=1).

export PATH="/opt/homebrew/bin:$PATH"

# Guard: this script changes the device. Zeus must have written "execute" or
# "override" in chat for this action; only then set ZEUS_APPROVED=1 to run it.
if [ "${ZEUS_APPROVED:-}" != "1" ]; then
  echo "Refusing to run: needs Zeus's explicit execute/override for this action (set ZEUS_APPROVED=1)." >&2
  exit 3
fi

if [ -z "$1" ]; then
  echo "Usage: $0 <ADB_SERIAL>"
  echo "Example: $0 <SERIAL>"
  exit 1
fi

D="$1"
ADB="adb -s $D shell"

echo "========================================"
echo "  SAMSUNG OPTIMIZATION — Serial: $D"
echo "========================================"
echo ""

echo "=== ANIMATIONS — INSTANT ==="
$ADB settings put global window_animation_scale 0.0
$ADB settings put global transition_animation_scale 0.0
$ADB settings put global animator_duration_scale 0.0
echo "  [OK] Animations disabled"

echo "=== BATTERY / SCREEN ==="
$ADB settings put system screen_brightness_mode 0
$ADB settings put system screen_brightness 42
$ADB settings put system screen_off_timeout 600000
echo "  [OK] Screen: manual 42%, 10min timeout"

echo "=== NETWORKING ==="
$ADB settings put global private_dns_mode opportunistic
echo "  [OK] Private DNS: Automatic (opportunistic DoT — works with Tailscale)"
echo "  NOTE: Do NOT use 'hostname' mode — conflicts with Tailscale MagicDNS"

echo "=== DEVELOPER OPTIONS ==="
$ADB settings put global development_settings_enabled 1
$ADB settings put global stay_awake_on_plugged_in 1
$ADB settings put global package_verifier_enable 0
echo "  [OK] Dev settings visible, stay-awake on, verifier off"

echo "=== LOCATION — HIGH ACCURACY ==="
$ADB settings put secure location_mode 3
$ADB settings put secure location_providers_allowed +gps
echo "  [OK] GPS high accuracy + GPS on"

echo "=== GNSS RAW MEASUREMENTS (ATAK) ==="
$ADB settings put secure gnss_raw_meas_full_tracking 1 2>/dev/null
echo "  [OK] GNSS raw measurements enabled"

echo "=== DARK MODE ==="
$ADB settings put secure ui_night_mode 2
echo "  [OK] Dark mode"

echo "=== RAM MANAGEMENT ==="
$ADB settings put global cached_apps_freezer enabled
echo "  [OK] Cached app freezer enabled"

echo "=== BLUETOOTH OFF ==="
$ADB settings put global bluetooth_on 0
$ADB shell svc bluetooth disable 2>/dev/null
echo "  [OK] Bluetooth disabled"

echo "=== AUTO ROTATE ==="
$ADB settings put system accelerometer_rotation 1
echo "  [OK] Auto-rotate on"

echo "=== GRANT FINE LOCATION TO PLAY STORE ==="
$ADB shell pm grant com.android.vending android.permission.ACCESS_FINE_LOCATION 2>/dev/null
$ADB shell pm grant com.google.android.gms android.permission.ACCESS_FINE_LOCATION 2>/dev/null
echo "  [OK] Fine location granted to Play Store + GMS"

echo ""
echo "=== VERIFICATION ==="
echo "Animations: $($ADB settings get global window_animation_scale) / $($ADB settings get global transition_animation_scale) / $($ADB settings get global animator_duration_scale)"
echo "Location mode: $($ADB settings get secure location_mode)"
echo "GPS: $($ADB settings get secure location_providers_allowed)"
echo "DNS mode: $($ADB settings get global private_dns_mode)"
echo "Screen timeout: $($ADB settings get system screen_off_timeout)"
echo "Dark mode: $($ADB settings get secure ui_night_mode)"
echo "Stay awake: $($ADB settings get global stay_awake_on_plugged_in)"
echo "Bluetooth: $($ADB settings get global bluetooth_on)"
echo "Auto-rotate: $($ADB settings get system accelerometer_rotation)"
echo "GNSS raw: $($ADB settings get secure gnss_raw_meas_full_tracking 2>/dev/null)"

echo ""
echo "========================================"
echo "  OPTIMIZATION COMPLETE"
echo "========================================"
