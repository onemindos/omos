# Device Registry

The authoritative roster of every device on the tactical stack. **No secrets here** — track intent,
identity, and state only. One cert per physical device. Retire = revoke cert + mark retired.

Choosing what hardware to buy per mission tier → [eud-hardware.md](eud-hardware.md).
Open thread — personal daily-driver pick + PTT plan → [pickup-daily-driver-and-ptt.md](pickup-daily-driver-and-ptt.md).

## TAK devices

| device_id | callsign | tak_user | platform | tailscale_ip | group | role | cert issued | provisioned | connection | status |
|---|---|---|---|---|---|---|---|---|---|---|
| omos-iphone17 | zeus | omos-iphone17 | iTAK (iOS) | 100.122.203.68 | Red | Team Lead | yes | soft-cert data package imported | Tailscale | **live** (connected 2026-07-02) |
| omos-ipad | zeus-pad | omos-ipad | iTAK (iOS) | 100.93.208.36 | Red | Team Member | yes | import `OneMind-iPad.zip` | Tailscale | cert-issued |
| omos-tablet | zeus-tab | omos-tablet | ATAK (Android, 32-bit) | 100.91.64.68 | Red | Team Member | yes | import `OneMind-Tablet.zip` | Tailscale | cert-issued |

**Connecting a device = import its soft-cert data package.** See [SKILL.md](SKILL.md) for the full
procedure (and the dead-ends to avoid). Build/regenerate with
[scripts/build-itak-datapackage.sh](scripts/build-itak-datapackage.sh).

Ready-to-AirDrop packages (gitignored) in `deploy/exports/`: `OneMind-iPhone17.zip`, `OneMind-iPad.zip`,
`OneMind-Tablet.zip` — each bundles the device client cert + CA truststore (pass: the keystore password, default until step-2 rotation). Raw
`*.p12` + `ca.pem` also staged there. All registered as cert users in group **Red**.

Statuses: `planned` → `cert-pending` → `cert-issued` → `provisioned` → `live` → `retired`.

TAK server reachable from all devices over Tailscale at **100.74.182.117:8089** (SSL/TLS).
Server cert SANs include the Tailscale IP + hostname + LAN IP → **no hostname-mismatch prompt**; connect directly.
Stale Tailscale android peers to prune in admin console: `omos-pad`, `omos-tab`, `omos-tv` (offline 21–48d).

## Voice (Murmur) users

| user | device(s) | registered | notes |
|---|---|---|---|
| SuperUser | — | yes | admin only; do not operate day-to-day as SuperUser |

## Group / role convention

| Purpose | TAK group | Role |
|---|---|---|
| Chief / admin | Red | Team Lead |
| Trusted personal device | Red | Team Member |
| Drone / UAV controller | Air | Team Member |
| Sensor / tracker | Sensor | Team Member |
| Guest / civilian | (separate group) | never admin |

## Naming convention

`<owner-or-function>-<platform-or-index>` — e.g. `zeus-phone`, `zeus-tablet`, `field-android-01`,
`uav-controller-01`, `taktracker-01`. Cert common name = device_id = tak_user (keeps it 1:1).
