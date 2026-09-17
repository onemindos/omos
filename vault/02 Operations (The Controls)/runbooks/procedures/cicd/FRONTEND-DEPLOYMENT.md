# OneMind OS - Multi-Platform Deployment Guide

## Platforms

| Platform | Status | Build Command | Deploy Target |
|----------|--------|---------------|---------------|
| **Web** | ✅ Ready | `flutter build web --release` | Docker + nginx (port 3005) |
| **iOS** | ✅ Tested | `flutter build ipa --release` | App Store / TestFlight |
| **Android** | 🟡 Ready | `flutter build apk --release` | Play Store / APK |
| **macOS** | 🟡 Ready | `flutter build macos --release` | Mac App Store / DMG |
| **Windows** | 🟡 Ready | `flutter build windows --release` | Microsoft Store / MSI |
| **Linux** | 🟡 Ready | `flutter build linux --release` | Snap Store / AppImage |

---

## 1. Web Deployment (Docker + nginx)

### Local Testing
```bash
cd mobile
flutter run -d chrome
# Or serve build output:
flutter build web --release
python3 -m http.server 8000 --directory build/web
```

### Production Deployment (VM)
```bash
# On your Mac (build)
cd /Users/zeuslegacy/Documents/onemind-os
flutter build web --release

# Deploy to VM
git add -A && git commit -m "build: update Flutter web build"
git push

# On VM (100.102.21.44)
cd ~/onemind-os
git pull
docker compose -f infra/compose/docker-compose.core.yml up -d --build flutter-web
```

**Access:** `http://100.102.21.44:3005` or via Cloudflare Tunnel

---

## 2. iOS Deployment

### Requirements
- macOS with Xcode 15+
- Apple Developer Account ($99/year)
- Signing certificates configured

### Build for TestFlight/App Store
```bash
cd mobile

# Update version in pubspec.yaml
# version: 1.0.0+1  # version+build_number

# Clean build
flutter clean
flutter pub get
cd ios && pod install && cd ..

# Build archive
flutter build ipa --release \
  --export-options-plist=ios/ExportOptions.plist

# Output: build/ios/ipa/onemind_mobile.ipa
```

### Upload to TestFlight
```bash
# Option 1: Xcode (GUI)
open ios/Runner.xcworkspace
# Product → Archive → Distribute App → TestFlight

# Option 2: Command line
xcrun altool --upload-app \
  --type ios \
  --file build/ios/ipa/onemind_mobile.ipa \
  --apiKey YOUR_API_KEY \
  --apiIssuer YOUR_ISSUER_ID
```

### Direct Install (Development)
```bash
# Connect iPhone via USB
flutter install -d Zeus

# Or wireless (if paired):
flutter run -d Zeus --release
```

---

## 3. Android Deployment

### Requirements
- Android SDK installed (via Flutter or Android Studio)
- Keystore for signing (generate once, keep secure)

### Generate Signing Key (First Time Only)
```bash
cd mobile/android
keytool -genkey -v \
  -keystore ~/onemind-release-key.jks \
  -keyalg RSA \
  -keysize 2048 \
  -validity 10000 \
  -alias onemind

# Save password in 1Password!
```

### Configure Signing (android/key.properties)
```properties
storePassword=YOUR_KEYSTORE_PASSWORD
keyPassword=YOUR_KEY_PASSWORD
keyAlias=onemind
storeFile=/Users/zeuslegacy/onemind-release-key.jks
```

### Build APK/AAB
```bash
cd mobile

# APK (for direct install)
flutter build apk --release
# Output: build/app/outputs/flutter-apk/app-release.apk

# AAB (for Play Store)
flutter build appbundle --release
# Output: build/app/outputs/bundle/release/app-release.aab
```

### Install APK
```bash
# Connect Android device via USB
flutter install -d <device-id>

# Or transfer APK and install manually
adb install build/app/outputs/flutter-apk/app-release.apk
```

### Upload to Play Store
1. Go to [Google Play Console](https://play.google.com/console)
2. Create app listing
3. Upload `app-release.aab`
4. Submit for review

---

## 4. macOS Deployment

### Requirements
- macOS with Xcode
- Apple Developer Account
- Mac App Store signing certificates

### Build for Mac App Store
```bash
cd mobile

# Clean and build
flutter clean
flutter pub get
cd macos && pod install && cd ..

flutter build macos --release

# Output: build/macos/Build/Products/Release/onemind_mobile.app
```

### Create DMG (for direct download)
```bash
# Install create-dmg
brew install create-dmg

# Create DMG
create-dmg \
  --volname "OneMind OS" \
  --window-pos 200 120 \
  --window-size 800 400 \
  --icon-size 100 \
  --icon "onemind_mobile.app" 200 190 \
  --hide-extension "onemind_mobile.app" \
  --app-drop-link 600 185 \
  "OneMindOS-1.0.0.dmg" \
  "build/macos/Build/Products/Release/onemind_mobile.app"
```

---

## 5. Windows Deployment

### Requirements
- Windows 10+ or Windows Server
- Visual Studio 2022 with C++ workload
- Flutter Windows desktop support enabled

### Build for Windows
```bash
cd mobile

flutter build windows --release

# Output: build/windows/x64/runner/Release/
```

### Create MSI Installer
```bash
# Install WiX Toolset
# https://wixtoolset.org/

# Create installer with Inno Setup (easier alternative)
# Download from: https://jrsoftware.org/isdl.php

# Create installer script (setup.iss)
```

### MSI Installer Script (setup.iss)
```ini
[Setup]
AppName=OneMind OS
AppVersion=1.0.0
DefaultDirName={pf}\OneMindOS
OutputDir=build\windows\installer
OutputBaseFilename=OneMindOS-Setup

[Files]
Source: "build\windows\x64\runner\Release\*"; DestDir: "{app}"; Flags: recursesubdirs

[Icons]
Name: "{autoprograms}\OneMind OS"; Filename: "{app}\onemind_mobile.exe"
```

---

## 6. Linux Deployment

### Build for Linux
```bash
cd mobile

flutter build linux --release

# Output: build/linux/x64/release/bundle/
```

### Create AppImage
```bash
# Install appimagetool
wget https://github.com/AppImage/AppImageKit/releases/download/continuous/appimagetool-x86_64.AppImage
chmod +x appimagetool-x86_64.AppImage

# Create AppDir structure
mkdir -p AppDir/usr/bin
cp -r build/linux/x64/release/bundle/* AppDir/usr/bin/

# Create desktop file
cat > AppDir/onemind.desktop << 'EOF'
[Desktop Entry]
Name=OneMind OS
Exec=onemind_mobile
Icon=onemind
Type=Application
Categories=Utility;
EOF

# Build AppImage
./appimagetool-x86_64.AppImage AppDir OneMindOS-1.0.0-x86_64.AppImage
```

### Create Snap Package
```bash
# Install snapcraft
sudo snap install snapcraft --classic

# Create snapcraft.yaml in mobile/snap/
# Run:
snapcraft

# Publish to Snap Store
snapcraft login
snapcraft upload OneMindOS_1.0.0_amd64.snap
```

---

## Deployment Scripts

### scripts/deploy-web.sh
```bash
#!/bin/bash
set -e

echo "🌐 Building Flutter web..."
cd mobile
flutter build web --release

echo "📦 Building Docker image..."
cd ..
docker build -f Dockerfile.flutter -t onemind-flutter-web:latest .

echo "🚀 Deploying to VM..."
git add -A
git commit -m "build: update Flutter web build" || true
git push

ssh -i ~/.ssh/onemind-key.pem ubuntu@100.102.21.44 << 'EOF'
cd ~/onemind-os
git pull
docker compose -f infra/compose/docker-compose.core.yml up -d --build flutter-web
EOF

echo "✅ Deployed! Access at: http://100.102.21.44:3005"
```

### scripts/deploy-ios.sh
```bash
#!/bin/bash
set -e

echo "📱 Building iOS app..."
cd mobile
flutter clean
flutter pub get
cd ios && pod install && cd ..

flutter build ipa --release

echo "📤 Uploading to TestFlight..."
xcrun altool --upload-app \
  --type ios \
  --file build/ios/ipa/onemind_mobile.ipa \
  --apiKey $APPLE_API_KEY \
  --apiIssuer $APPLE_ISSUER_ID

echo "✅ Uploaded to TestFlight!"
```

### scripts/deploy-android.sh
```bash
#!/bin/bash
set -e

echo "🤖 Building Android app..."
cd mobile
flutter clean
flutter pub get

flutter build appbundle --release

echo "📦 Bundle ready: build/app/outputs/bundle/release/app-release.aab"
echo "Upload to Play Console: https://play.google.com/console"
```

---

## GitHub Actions CI/CD

### .github/workflows/build-all-platforms.yml
```yaml
name: Build All Platforms

on:
  push:
    branches: [main]
    paths:
      - 'mobile/**'
  workflow_dispatch:

jobs:
  build-web:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4
      - uses: subosito/flutter-action@v2
        with:
          flutter-version: '3.x'
      - run: |
          cd mobile
          flutter pub get
          flutter build web --release
      - uses: actions/upload-artifact@v4
        with:
          name: web-build
          path: mobile/build/web

  build-ios:
    runs-on: macos-latest
    steps:
      - uses: actions/checkout@v4
      - uses: subosito/flutter-action@v2
      - run: |
          cd mobile
          flutter pub get
          cd ios && pod install && cd ..
          flutter build ios --release --no-codesign

  build-android:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4
      - uses: actions/setup-java@v4
        with:
          distribution: 'zulu'
          java-version: '17'
      - uses: subosito/flutter-action@v2
      - run: |
          cd mobile
          flutter pub get
          flutter build apk --release
      - uses: actions/upload-artifact@v4
        with:
          name: android-apk
          path: mobile/build/app/outputs/flutter-apk/app-release.apk

  build-desktop:
    strategy:
      matrix:
        os: [macos-latest, windows-latest, ubuntu-latest]
        include:
          - os: macos-latest
            platform: macos
          - os: windows-latest
            platform: windows
          - os: ubuntu-latest
            platform: linux
    runs-on: ${{ matrix.os }}
    steps:
      - uses: actions/checkout@v4
      - uses: subosito/flutter-action@v2
      - run: |
          cd mobile
          flutter pub get
          flutter build ${{ matrix.platform }} --release
```

---

## Version Management

### Bump version in pubspec.yaml
```yaml
version: 1.0.1+2  # version+build_number
#        ↑   ↑
#        |   └── Build number (increments for each build)
#        └────── Semantic version (major.minor.patch)
```

### Release Checklist
- [ ] Update version in `mobile/pubspec.yaml`
- [ ] Update `CHANGELOG.md`
- [ ] Test on all target platforms
- [ ] Create git tag: `git tag v1.0.1 && git push --tags`
- [ ] Build production artifacts
- [ ] Deploy/upload to stores
- [ ] Update Cloudflare Tunnel (if web URL changed)

---

## Access URLs

| Platform | URL/Location |
|----------|-------------|
| **Web (Dev)** | `http://localhost:3005` |
| **Web (Prod)** | `http://100.102.21.44:3005` or `https://app.onemindos.dev` |
| **iOS** | TestFlight → App Store |
| **Android** | APK / Play Store |
| **macOS** | DMG / Mac App Store |
| **Windows** | MSI installer |
| **Linux** | AppImage / Snap Store |

---

## Monitoring & Updates

### Check running version
```bash
# Web
curl http://100.102.21.44:3005/version.json

# Mobile (in app)
# Settings → About → Version 1.0.0 (Build 1)
```

### Hot reload (development)
```bash
# All platforms support hot reload during development
flutter run -d <device>
# Press 'r' to hot reload
# Press 'R' to hot restart
```

---

## Cost Breakdown

| Platform | Store Fee | Annual Cost | Notes |
|----------|-----------|-------------|-------|
| Web | Free | $0 | Self-hosted |
| iOS | $99/year | $99 | Apple Developer Program |
| Android | $25 one-time | $0 | Google Play Developer |
| macOS | $99/year | $0 | Included with iOS |
| Windows | $19 one-time | $0 | Microsoft Store (optional) |
| Linux | Free | $0 | Open distribution |

**Total:** ~$118/year (iOS + Android one-time)

---

## Next Steps

1. ✅ Web deployed to port 3005
2. 🟡 Add responsive layouts (see RESPONSIVE-DESIGN.md)
3. 🟡 Test on all platforms
4. 🟡 Set up CI/CD for automated builds
5. 🟡 Submit to app stores

Run this to make scripts executable:
```bash
chmod +x scripts/deploy-*.sh
```


---

## Related

- [[00-24 SO (Sovereign Operations)/02 Infrastructure (VPS, Servers & Networking) (The Body)/02 Infrastructure|02 Infrastructure]]
