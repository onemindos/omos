# Linear Issues: Flutter Mobile Testing

Two Linear issues ready to add for future mobile testing tasks.

## Issue 1: Samsung Galaxy A11 Testing 🤖

**File:** `docs/planning/LINEAR_ISSUE_SAMSUNG_A11_TESTING.md`

**Title:** Test OneMind OS Flutter app on Samsung Galaxy A11

**Summary:**
- Set up Android SDK on Mac
- Enable Developer Mode on Samsung A11
- Build and install APK
- Test all 24 screens and core functionality
- Document any issues found

**Priority:** Medium
**Estimate:** 2 hours
**Labels:** testing, mobile, android, device-testing

**Prerequisites:**
1. `brew install --cask android-commandlinetools`
2. `flutter doctor --android-licenses`
3. Enable USB Debugging on A11

**Key Testing Areas:**
- Chat functionality (SSE streaming)
- All 24 screens accessible
- Performance on 2-3GB RAM
- Network handling (WiFi + cellular)
- UI/UX on 6.4" display

---

## Issue 2: Ubuntu 24.04 Linux Build 🐧

**File:** `docs/planning/LINEAR_ISSUE_LINUX_BUILD.md`

**Title:** Build OneMind OS Flutter desktop app on Ubuntu 24.04

**Summary:**
- Install Flutter SDK on Ubuntu 24.04
- Install Linux development dependencies
- Build Flutter desktop app
- Test desktop functionality
- Create AppImage for distribution (optional)

**Priority:** Low (can be done later)
**Estimate:** 3 hours
**Labels:** testing, desktop, linux, build

**Prerequisites:**
1. Ubuntu 24.04 machine available
2. Install: clang, cmake, ninja-build, libgtk-3-dev, etc.
3. Install Flutter SDK (via official installer or snap)

**Key Testing Areas:**
- Desktop window management
- Keyboard/mouse interactions
- Performance on Linux desktop
- AppImage portability

---

## How to Add to Linear

### Via Linear Web UI:

1. Go to: https://linear.app/your-workspace/team/OMOS/new
2. Click "New Issue"
3. Copy content from the markdown files above
4. Set priority, labels, assignee
5. Create issue

### Via Linear CLI (if installed):

```bash
# Issue 1: Samsung A11 Testing
linear issue create \
  --title "Test OneMind OS Flutter app on Samsung Galaxy A11" \
  --description "$(cat docs/planning/LINEAR_ISSUE_SAMSUNG_A11_TESTING.md)" \
  --priority 2 \
  --label testing,mobile,android,device-testing \
  --estimate 2

# Issue 2: Linux Build
linear issue create \
  --title "Build OneMind OS Flutter desktop app on Ubuntu 24.04" \
  --description "$(cat docs/planning/LINEAR_ISSUE_LINUX_BUILD.md)" \
  --priority 3 \
  --label testing,desktop,linux,build \
  --estimate 3
```

---

## Current Status

| Task | Status | When |
|------|--------|------|
| **GitHub Secrets Setup** | ⏸️ Pending | Do now for automated deployment |
| **Android SDK Install** | ⏸️ Pending | Do now to build APK |
| **Samsung A11 Testing** | 📋 Issue created | Do later (after APK ready) |
| **Linux Build** | 📋 Issue created | Do later (when Ubuntu available) |

---

## Quick Actions

**Do Now:**
1. Set up GitHub secrets for automated deployment:
   - `VM_HOST = 100.102.21.44`
   - `VM_USER = ubuntu`
   - `VM_SSH_KEY = <paste ~/.ssh/onemind-key.pem>`

2. Install Android SDK:
   ```bash
   brew install --cask android-commandlinetools
   flutter doctor --android-licenses
   ```

**Do Later:**
- Add Issue 1 to Linear when ready to test on Samsung A11
- Add Issue 2 to Linear when Ubuntu 24.04 machine is ready

---

## Files Created

1. `docs/planning/LINEAR_ISSUE_SAMSUNG_A11_TESTING.md` - Complete Samsung A11 testing guide
2. `docs/planning/LINEAR_ISSUE_LINUX_BUILD.md` - Complete Linux build guide
3. `docs/planning/LINEAR_ISSUES_SUMMARY.md` - This file
4. `.github/workflows/flutter-deploy-all.yml` - Automated CI/CD for all platforms
5. `scripts/build-linux.sh` - Linux build script
6. `docs/AUTOMATED-DEPLOYMENT.md` - Full automation guide

All documentation is ready for when you're ready to tackle these tasks!
