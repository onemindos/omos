# OneMind TAK Live Demo

Scripted, real-time tactical demonstrations pushed to live iTAK/ATAK clients over
the TAK Server TLS CoT feed. Deploy a UAV, spawn hostiles, run a fire mission, and
fly an ISR orbit over any connected team member — driven entirely by Python
injecting Cursor-on-Target (CoT) events.

Proven live on: **Zeus-iPhone** (CT), **Carku** (Albania), **Paulino** (NYC).

---

## Layout

```
tak-live-demo/
├── README.md
├── scripts/
│   ├── onemind-fire-mission.py   # main demo: deploy → hostiles → fire mission → ISR orbit
│   ├── chopper-insert.py         # CHOPPER-1 flies point-to-point (Velipoje → marker)
│   ├── clear-map.py              # wipe all demo contacts (proper t-x-d-d-r delete)
│   ├── add-member.sh             # provision a new team member end-to-end (+ group fix)
│   └── who.sh                    # show who's connected and their auth group
├── certs/                        # injection client cert (gitignored — secrets)
│   ├── iphone17-cert.pem
│   └── iphone17-key.pem
└── packages/                     # generated enrollment ZIPs (gitignored)
    └── OneMind-*.zip
```

---

## Quick start

```bash
cd /Users/zeus/ATAK/demos/tak-live-demo/scripts

# 1. Clear anything left on the map
python3 clear-map.py

# 2. Point the demo at a target and run it
#    (edit AO_LAT/AO_LON near the bottom of the script to the member's live coords)
python3 onemind-fire-mission.py
```

Get a member's live coordinates:

```bash
docker exec tak-database bash -c "PGPASSWORD=martipassword psql -U martiuser -d cot -h localhost -t -c \
\"SELECT ROUND(ST_Y(event_pt::geometry)::numeric,5), ROUND(ST_X(event_pt::geometry)::numeric,5) \
FROM cot_router WHERE detail ILIKE '%CALLSIGN%' AND time > NOW() - INTERVAL '2 minutes' ORDER BY time DESC LIMIT 1;\""
```

---

## The fire-mission scenario (`onemind-fire-mission.py`)

Runs entirely around the AO center (set to the target member's position):

1. **Phase 1 — Deploy.** REAPER-1 (friendly, blue) flies in from the west, arrives
   over the AO, and orbits at 300 m.
2. **Phase 2 — Detect.** 4 hostiles (red) appear in a clean column 500 m east,
   150 m apart. REAPER-1 flies one extra recon orbit (~22 s) before engaging.
3. **Phase 3 — Fire mission.** REAPER-1 works down the column south→north, striking
   HOSTILE-01→04. Each strike drops a KIA marker, then the hostile is removed.
4. **Phase 4 — ISR.** REAPER-1 tightens to a 150 m orbit over the member for 20 s.
5. **Cleanup.** REAPER-1 and strike markers are removed via proper delete.

Tunables live at the top of the phase section: hostile count/line, orbit radii,
speeds, and durations.

---

## Hard-won lessons (READ before editing)

### 1. Removing contacts — use the official delete CoT
iTAK/ATAK caches contacts **locally**. To remove one you must send a delete event:

```xml
<event uid="del-TARGET" type="t-x-d-d-r" ...>
  <detail><link uid="TARGET" type="ORIGINAL_TYPE" relation="p-p"/></detail>
</event>
```

- Sending a stale/expired copy of the original event does **not** remove it.
- Moving it to 0,0 does **not** remove it.
- Disconnect/reconnect does **not** clear the client's local cache.
- `type` in the `<link>` must match the target's original CoT type.

### 2. Group assignment — the __ANON__ trap
TAKServer reads `UserAuthenticationFile.xml` **only at startup**. If you add a user
while the server is running and they connect before a restart, the server accepts
their valid cert but puts them in **`__ANON__`**, not their assigned group. Members
in different groups can't see each other.

- **Symptom:** new member connects fine, sets team color Red, but nobody sees them.
- **Fix:** register the user, then **restart takserver**, *then* have them connect.
- `add-member.sh` does this automatically. `who.sh` shows each client's real group.
- Team color (Red/Blue/Cyan in iTAK settings) is **cosmetic only** — it is NOT the
  authorization group. Both must be right: auth group from the file, color in-app.

### 3. Symbols & color
Color is driven by CoT **affiliation**, enforced by iTAK — you can't override it:
`a-f-*` = blue (friendly), `a-h-*` = red (hostile), `a-n-*` = green (neutral),
`a-u-*` = yellow (unknown). iTAK on iOS flattens most 2525 *shape* variants to a
generic icon, so lean on affiliation color + callsign labels, not fine symbology.
Used here: `a-f-A-M-F-Q` (attack air), `a-f-A-M-H-Q` (helo), `a-h-G-U-C-I`
(infantry), `a-h-G-U-C-A` (armor), `b-a-o-tbl` (strike/emergency marker).

### 4. Connection must persist
Open **one** SSL connection and hold it for the whole scenario, with a keepalive
ping every ~8 s. Open/send/close per event races the server's broadcast and events
never reach clients. All scripts here use a persistent connection + keepalive.

### 5. Geo-fences don't push
iTAK strips `u-d-c-c` (geocircle) shapes from server-pushed CoT. Draw shapes in-app
or ship them in a data package instead.

---

## Adding a new team member

```bash
cd /Users/zeus/ATAK/demos/tak-live-demo/scripts
./add-member.sh <username> [public_host]     # e.g. ./add-member.sh ruiz 74.108.23.135
```

This generates a TAK-CA-signed cert, registers it in the auth file with the Red
group, builds `packages/OneMind-<User>.zip`, and restarts the server so the group
loads. Send the ZIP; in iTAK they import it, set Team=Red and their callsign, and
connect — landing directly in Red.

---

## Environment

- **Server:** TAKServer 5.7-RELEASE-43 (Docker: `takserver` + `tak-database`).
- **Feed:** TLS CoT on `:8089` (local `127.0.0.1`, public `74.108.23.135`).
- **Injection cert:** `certs/iphone17-cert.pem` (derived from `omos-iphone17.p12`,
  Red group). Override paths with `TAK_CERT` / `TAK_KEY` / `TAK_HOST` / `TAK_PORT`.
- **CA / user files:** `../../server/standard/takserver-docker-5.7-RELEASE-43/tak/`.
```
