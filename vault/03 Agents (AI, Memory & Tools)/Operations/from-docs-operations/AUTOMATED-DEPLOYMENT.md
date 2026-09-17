# Automated Multi-Platform Deployment

## Overview

Every push to `main` automatically builds and deploys to **all 6 platforms**:
- **Web** → Deployed to VM (http://100.102.21.44:3005)
- **iOS** → Build artifacts uploaded (ready for TestFlight)
- **Android** → APK + AAB uploaded (ready for Play Store)
- **macOS** → Build artifacts uploaded
- **Windows** → Build artifacts uploaded
- **Linux** → Build artifacts uploaded

## Setup (One-Time)

### 1. Add GitHub Secrets

Go to: https://github.com/YOUR_USERNAME/onemind-os/settings/secrets/actions

Add these secrets:

| Secret Name | Value | Purpose |
|-------------|-------|---------|
| `VM_HOST` | `100.102.21.44` | VM IP for web deployment |
| `VM_USER` | `ubuntu` | SSH username |
| `VM_SSH_KEY` | `<your private key>` | SSH private key from `~/.ssh/onemind-key.pem` |

**To get your SSH key:**
```bash
cat ~/.ssh/onemind-key.pem
# Copy the entire output (including BEGIN/END lines)
```

### 2. Optional: Auto-Upload to App Stores

#### TestFlight (iOS)

1. Generate API key at: https://appstoreconnect.apple.com/access/api
2. Add secrets:
   - `APPLE_API_KEY` - Your API key
   - `APPLE_ISSUER_ID` - Your issuer ID
3. Uncomment the "Upload to TestFlight" step in `.github/workflows/flutter-deploy-all.yml`

#### Play Store (Android)

1. Create service account at: https://console.cloud.google.com/iam-admin/serviceaccounts
2. Download JSON key
3. Add secret:
   - `PLAY_STORE_SERVICE_ACCOUNT` - Paste JSON contents
4. Uncomment the "Upload to Play Store" step in `.github/workflows/flutter-deploy-all.yml`

## Usage

### Automatic Deployment

Just push to main:

```bash
git add .
git commit -m "feat: new feature"
git push origin main
```

GitHub Actions will:
1. Build all 6 platforms in parallel
2. Deploy web to VM automatically
3. Upload artifacts for other platforms
4. Show deployment summary

### Manual Trigger

You can also trigger manually from GitHub:
1. Go to: https://github.com/YOUR_USERNAME/onemind-os/actions
2. Select "Flutter Deploy All Platforms"
3. Click "Run workflow"

## Build Artifacts

After each run, download artifacts from:
https://github.com/YOUR_USERNAME/onemind-os/actions

| Artifact | File | Use |
|----------|------|-----|
| `web-build` | Built web files | Already deployed to VM |
| `ios-ipa` | onemind_mobile.ipa | Upload to TestFlight/App Store |
| `android-apk` | app-release.apk | Direct install on devices |
| `android-aab` | app-release.aab | Upload to Play Store |
| `macos-app` | onemind_mobile.app | Create DMG or upload to Mac Store |
| `windows-build` | Release folder | Create MSI installer |
| `linux-build` | Bundle folder | Create AppImage/Snap |

## Build Times

| Platform | Time | Runner |
|----------|------|--------|
| Web | ~2min | ubuntu-latest |
| iOS | ~5min | macos-latest |
| Android | ~3min | ubuntu-latest |
| macOS | ~4min | macos-latest |
| Windows | ~4min | windows-latest |
| Linux | ~3min | ubuntu-latest |

**Total parallel time:** ~5min (slowest runner)

## Cost

GitHub Actions free tier:
- **Public repos:** Unlimited minutes
- **Private repos:** 2,000 minutes/month

Estimated usage:
- 1 deployment = ~20 build minutes (all platforms in parallel)
- 100 deployments/month = 2,000 minutes

## Monitoring

View build status:
- Badge: `![Build](https://github.com/YOUR_USERNAME/onemind-os/actions/workflows/flutter-deploy-all.yml/badge.svg)`
- Dashboard: https://github.com/YOUR_USERNAME/onemind-os/actions

## Troubleshooting

### Web deployment fails

**Error:** SSH connection timeout

**Fix:** Check VM SSH access:
```bash
ssh -i ~/.ssh/onemind-key.pem ubuntu@100.102.21.44
```

If that works, verify the `VM_SSH_KEY` secret is correct.

### iOS build fails

**Error:** No signing certificate

**Fix:** This is expected for the `--no-codesign` build. For App Store builds:
1. Add signing certificates to Xcode
2. Remove `--no-codesign` flag
3. Add provisioning profile to workflow

### Android build fails

**Error:** Keystore not found

**Fix:** For signed builds, add keystore:
1. Base64 encode keystore: `base64 ~/onemind-release-key.jks`
2. Add secret: `ANDROID_KEYSTORE_BASE64`
3. Add secrets: `KEYSTORE_PASSWORD`, `KEY_ALIAS`, `KEY_PASSWORD`
4. Update workflow to decode and use keystore

## Local Builds (Faster)

For faster iteration, build locally:

**Mac (all platforms):**
```bash
./scripts/build-all-mac.sh
```

**Ubuntu 24.04 (Linux only):**
```bash
./scripts/build-linux.sh
```

**Windows (Windows only):**
```powershell
cd mobile
flutter build windows --release
```
