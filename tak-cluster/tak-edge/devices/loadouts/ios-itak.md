# Loadout: ios-itak

iPhone/iPad on **iTAK**. iOS is sandboxed — no `adb`, no sideloading. Provisioning is
App Store + cert/profile import. USB helps only for file transfer / Apple Configurator.

## Steps
1. Install **iTAK** from the App Store (free).
2. Issue the device cert from `deploy/`:
   ```bash
   deploy/scripts/tak-add-device.sh zeus-phone zeus-phone Red "Team Member"
   ```
   Cert lands in `deploy/exports/zeus-phone.p12`.
3. Get the `.p12` + `ca.pem` onto the device:
   - **AirDrop** (fastest), or
   - **USB via Apple Configurator** (`cfgutil`) for a repeatable/scripted path, or
   - HTTPS enrollment: Safari → `https://192.168.1.183:8446/Marti/enrollment/`.
4. Install the cert profile: Settings → General → VPN & Device Management → install.
5. If prompted for CA trust: Settings → General → About → Certificate Trust Settings → enable full trust for **TAK-CA**.
6. iTAK → add server:
   ```
   Address:  192.168.1.183   (LAN)   |   100.74.182.117 (Tailscale/cellular)
   Port:     8089
   Protocol: SSL/TLS
   Cert:     zeus-phone.p12  (pass: atakatak unless changed)
   ```
7. Confirm on WebTAK, mark `live` in `devices/registry.md`.

## Scripted USB (optional, Apple Configurator)
```bash
brew install --cask apple-configurator
cfgutil list                              # find ECID/UDID of the tethered device
cfgutil -e <ECID> install-profile ca.mobileconfig   # if you wrap CA/cert as a .mobileconfig
```
See `devices/scripts/ios-enroll.md` for the manual walkthrough.
