# Deploy All Platforms - Status

## Quick Status

| Platform | Status | Deploy Target | Build Time |
|----------|--------|---------------|------------|
| **Web** | ✅ Built | VM Port 3005 | ~30s |
| **iOS** | 🔄 Ready | TestFlight | ~2min |
| **Android** | 🔄 Ready | APK/Play Store | ~1min |
| **macOS** | 🔄 Ready | DMG/Mac Store | ~1min |
| **Windows** | ⏸️ Needs Windows | MSI installer | ~2min |
| **Linux** | ⏸️ Needs Linux | AppImage/Snap | ~1min |

## Automated Deployment (Recommended)

**Push to main → Builds all platforms automatically!**

```bash
git add .
git commit -m "your changes"
git push origin main
```

See `docs/AUTOMATED-DEPLOYMENT.md` for setup.

## Manual Build (Mac)

```bash
# Build all Mac-compatible platforms at once
./scripts/build-all-mac.sh
```

This will build: Web, iOS, Android (APK + AAB), and macOS.

## Manual Build (Ubuntu 24.04)

```bash
# On your Ubuntu machine
./scripts/build-linux.sh
```

This will build the Linux desktop app.

## Deploy Web to VM

```bash
# Manual deployment (until SSH is fixed):
# 1. SSH into VM
# 2. cd ~/onemind-os && git pull
# 3. docker compose -f infra/compose/docker-compose.core.yml up -d --build flutter-web
# 4. Access at http://100.102.21.44:3005
```

## Full documentation: See `mobile/DEPLOYMENT.md`


---

## Related

- [[00-24 SO (Sovereign Operations)/02 Infrastructure (VPS, Servers & Networking) (The Body)/02 Infrastructure|02 Infrastructure]]
