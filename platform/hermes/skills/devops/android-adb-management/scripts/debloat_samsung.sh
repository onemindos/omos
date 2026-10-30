#!/bin/bash
# S24 Ultra / Samsung Galaxy Debloat Script
# Removes ~93 bloatware packages from a stock Samsung device.
# Usage: bash debloat_samsung.sh <ADB_SERIAL>
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
UNINST="adb -s $D shell pm uninstall --user 0"
DISABLE="adb -s $D shell pm disable-user --user 0"
SUCCESS=0
FAIL=0

echo "========================================"
echo "  SAMSUNG DEBLOAT — Serial: $D"
echo "========================================"
echo ""

echo "=== PHASE 1: PREVIOUS OWNER / CARRIER JUNK ==="
for pkg in \
  com.facebook.katana \
  com.facebook.appmanager \
  com.facebook.services \
  com.facebook.system \
  com.linkedin.android \
  com.spotify.music \
  com.microsoft.office.outlook \
  com.microsoft.office.officehubrow \
  com.microsoft.skydrive \
  com.microsoft.appmanager \
  com.att.myWireless \
  com.rsupport.rs.activity.rsupport.aas2 \
  com.hiya.star \
  com.verizon.mips.services \
  com.vzw.apnlib
do
  result=$($UNINST $pkg 2>&1)
  if echo "$result" | grep -q "Success"; then
    echo "  [OK] $pkg"
    SUCCESS=$((SUCCESS+1))
  else
    echo "  [SKIP] $pkg — $result"
    FAIL=$((FAIL+1))
  fi
done

echo ""
echo "=== PHASE 2: SAMSUNG BLOAT ==="
for pkg in \
  com.samsung.android.bixby.agent \
  com.samsung.android.bixby.wakeup \
  com.samsung.android.bixbyvision.framework \
  com.samsung.android.bixby.ondevice.enus \
  com.samsung.android.bixby.ondevice.esus \
  com.samsung.android.app.spage \
  com.samsung.android.tvplus \
  com.samsung.android.themecenter \
  com.samsung.android.themestore \
  com.sec.android.app.kidshome \
  com.samsung.android.kidsinstaller \
  com.samsung.android.game.gamehome \
  com.samsung.android.game.gametools \
  com.samsung.android.game.gos \
  com.samsung.android.spay \
  com.samsung.android.coldwalletservice \
  com.sec.android.app.samsungapps \
  com.samsung.android.aremoji \
  com.samsung.android.aremojieditor \
  com.samsung.android.stickercenter \
  com.sec.android.mimage.avatarstickers \
  com.sec.penup \
  com.samsung.android.app.dressroom \
  com.samsung.android.app.sketchbook \
  com.samsung.android.app.tips \
  com.samsung.ecomm \
  com.samsung.android.voc \
  com.samsung.android.app.find \
  com.samsung.android.oneconnect \
  com.samsung.android.sree \
  com.sec.android.app.voicenote \
  com.samsung.android.app.watchmanager \
  com.samsung.android.app.watchmanagerstub \
  com.samsung.android.app.parentalcare \
  com.sec.android.app.clockpackage \
  com.samsung.android.app.reminder \
  com.samsung.android.calendar \
  com.samsung.android.rubin.app \
  com.samsung.android.mhs.ai \
  com.samsung.android.app.contacts \
  com.samsung.android.messaging \
  com.samsung.android.allshare.service.mediashare \
  com.samsung.android.smartmirroring \
  com.sec.android.easyMover \
  com.sec.android.easyMover.Agent \
  com.samsung.android.smartswitchassistant \
  com.samsung.android.motionphoto.app \
  com.samsung.android.mocca \
  com.samsung.android.app.sharelive \
  com.samsung.android.aware.service \
  com.samsung.android.service.peoplestripe \
  com.samsung.android.privacydashboard \
  com.samsung.android.mapsagent \
  com.samsung.android.mcfds \
  com.samsung.android.mcfserver \
  com.samsung.android.mcf.autohotspot \
  com.samsung.android.appseparation \
  com.samsung.android.app.clipboardedge \
  com.samsung.android.app.taskedge \
  com.samsung.android.app.aodservice \
  com.samsung.android.dynamiclock \
  com.samsung.android.widget.pictureframe \
  com.samsung.android.wallpaper.live \
  com.samsung.android.app.soundpicker \
  com.samsung.android.app.smartcapture \
  com.samsung.android.photoremasterservice \
  com.samsung.android.singletake.service \
  com.samsung.android.globalpostprocmgr \
  com.samsung.android.app.camera.sticker.facearavatar.preload
do
  result=$($UNINST $pkg 2>&1)
  if echo "$result" | grep -q "Success"; then
    echo "  [OK] $pkg"
    SUCCESS=$((SUCCESS+1))
  else
    echo "  [SKIP] $pkg"
    FAIL=$((FAIL+1))
  fi
done

echo ""
echo "=== PHASE 3: GOOGLE BLOAT ==="
for pkg in \
  com.google.android.youtube \
  com.google.android.apps.youtube.music \
  com.google.android.videos \
  com.google.android.apps.docs \
  com.google.android.apps.photos \
  com.google.android.apps.bard \
  com.google.android.apps.aiwallpapers \
  com.google.android.apps.tachyon \
  com.google.android.apps.turbo \
  com.google.android.googlequicksearchbox \
  com.google.android.apps.messaging \
  com.google.android.apps.restore
do
  result=$($UNINST $pkg 2>&1)
  if echo "$result" | grep -q "Success"; then
    echo "  [OK] $pkg"
    SUCCESS=$((SUCCESS+1))
  else
    echo "  [SKIP] $pkg"
    FAIL=$((FAIL+1))
  fi
done

echo ""
echo "=== PHASE 4: ADDITIONAL SAMSUNG REMOVALS ==="
for pkg in \
  com.sec.android.app.sbrowser \
  com.sec.android.app.shealth \
  com.samsung.android.app.notes \
  com.samsung.android.nmt.apps.t2t.languagepack.enesus \
  com.google.android.safetycore \
  com.google.android.contactkeys
do
  result=$($UNINST $pkg 2>&1)
  if echo "$result" | grep -q "Success"; then
    echo "  [OK] $pkg"
    SUCCESS=$((SUCCESS+1))
  else
    echo "  [SKIP] $pkg"
    FAIL=$((FAIL+1))
  fi
done

echo ""
echo "=== PHASE 5: POST-DEBLOAT CLEANUP ==="
# Turn off Bluetooth (tactical — don't need it broadcasting)
adb -s $D shell settings put global bluetooth_on 0 2>/dev/null
adb -s $D shell svc bluetooth disable 2>/dev/null
echo "  [OK] Bluetooth off"

# Enable auto-rotate
adb -s $D shell settings put system accelerometer_rotation 1
echo "  [OK] Auto-rotate on"

echo ""
echo "========================================"
echo "  DEBLOAT COMPLETE"
echo "  Removed: $SUCCESS"
echo "  Skipped: $FAIL"
echo "========================================"
echo ""
echo "Run optimize_samsung.sh next for system tuning."
