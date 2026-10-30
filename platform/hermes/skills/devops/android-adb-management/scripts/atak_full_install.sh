#!/bin/bash
# ATAK-CIV Full Loadout Install Script
# Installs ATAK-CIV + all plugins + TakTracker via ADB.
# Usage: bash atak_full_install.sh <ADB_SERIAL> [ATAK_DIR]
# Run on the Mac (ADB bridge) via Tailscale SSH from the onemind-hermes pod, after Zeus approves (ZEUS_APPROVED=1).
#
# Default ATAK_DIR: /Users/zeus/Desktop/atak  (note trailing space in folder name)

export PATH="/opt/homebrew/bin:$PATH"

# Guard: this script changes the device. Zeus must have written "execute" or
# "override" in chat for this action; only then set ZEUS_APPROVED=1 to run it.
if [ "${ZEUS_APPROVED:-}" != "1" ]; then
  echo "Refusing to run: needs Zeus's explicit execute/override for this action (set ZEUS_APPROVED=1)." >&2
  exit 3
fi

if [ -z "$1" ]; then
  echo "Usage: $0 <ADB_SERIAL> [ATAK_DIR]"
  echo "Example: $0 <SERIAL>"
  exit 1
fi

D="$1"
ATAK_DIR="${2:-/Users/zeus/Desktop/atak }"
LOG="/tmp/atak_install.log"
echo "" > $LOG

install_apk() {
  local name="$1"
  local path="$2"
  echo -n "Installing $name... "
  result=$(adb -s $D install -r "$path" 2>&1)
  if echo "$result" | grep -q "Success"; then
    echo "OK"
    echo "OK $name" >> $LOG
  else
    echo "FAIL: $result"
    echo "FAIL $name: $result" >> $LOG
  fi
}

echo "========================================"
echo "  ATAK-CIV FULL LOADOUT INSTALL"
echo "  Serial: $D"
echo "  Source: $ATAK_DIR"
echo "========================================"
echo ""

# 1. ATAK-CIV Core (Full, with bundled maps)
echo "=== ATAK-CIV CORE ==="
if [ -f "$ATAK_DIR/ATAK-5.7.0.10-e0d26f0e-civ-release.apk" ]; then
  install_apk "ATAK-CIV-Full" "$ATAK_DIR/ATAK-5.7.0.10-e0d26f0e-civ-release.apk"
elif [ -f "$ATAK_DIR/ATAK-5.7.0.10-e0d26f0e-civSmall-release.apk" ]; then
  echo "Full CIV not found, falling back to civSmall..."
  install_apk "ATAK-CIV-Small" "$ATAK_DIR/ATAK-5.7.0.10-e0d26f0e-civSmall-release.apk"
else
  echo "ERROR: No ATAK-CIV APK found in $ATAK_DIR"
  exit 1
fi

# 2. TakTracker
echo ""
echo "=== TAKTRACKER ==="
for apk in "$ATAK_DIR"/TakTracker-*.apk; do
  install_apk "TakTracker" "$apk"
  break
done

# 3. All Plugins
echo ""
echo "=== PLUGINS ==="
for apk in "$ATAK_DIR"/ATAK-CIV_5.7.0_loadout/Plugins/*.apk; do
  name=$(basename "$apk" .apk)
  install_apk "$name" "$apk"
done

# 4. UAS Plugins (subfolder)
echo ""
echo "=== UAS PLUGINS ==="
for apk in "$ATAK_DIR"/ATAK-CIV_5.7.0_loadout/Plugins/UAS_Tool_5.7.0/*.apk; do
  name=$(basename "$apk" .apk)
  install_apk "$name" "$apk"
done

# 5. Push non-APK files to device storage
echo ""
echo "=== NON-APK FILES ==="
ESP32_BIN="$ATAK_DIR/ATAK-CIV_5.7.0_loadout/Plugins/UAS_Tool_5.7.0/ESP32C3_crsfToCot.bin"
if [ -f "$ESP32_BIN" ]; then
  adb -s $D push "$ESP32_BIN" /sdcard/Download/ESP32C3_crsfToCot.bin
  echo "OK ESP32 firmware -> /sdcard/Download/"
fi

DJI_LICENSE=$(find "$ATAK_DIR" -name "license_com.atakmap*" -print -quit 2>/dev/null)
if [ -n "$DJI_LICENSE" ]; then
  adb -s $D push "$DJI_LICENSE" /sdcard/Download/DJI_LDM_License.txt
  echo "OK DJI LDM License -> /sdcard/Download/"
fi

# 6. Summary
echo ""
echo "========================================"
echo "  INSTALL COMPLETE"
echo "========================================"
echo ""
echo "Success: $(grep -c "^OK" $LOG)"
echo "Failed:  $(grep -c "^FAIL" $LOG)"
echo ""
echo "=== FAILURES (if any) ==="
grep "^FAIL" $LOG || echo "  None"
echo ""
echo "Verify: adb -s $D shell pm list packages | grep atakmap"
