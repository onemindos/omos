#!/usr/bin/env python3
"""
OneMind TAK Demo v3 — clean, direct, full control
One persistent SSL connection. Official CoT types. Proper delete.
"""
import ssl, socket, time, math, threading, os

# Certs resolve relative to this script (../certs) unless overridden by env vars.
_BASE = os.path.dirname(os.path.abspath(__file__))
CERT     = os.environ.get("TAK_CERT", os.path.join(_BASE, "..", "certs", "iphone17-cert.pem"))
KEY      = os.environ.get("TAK_KEY",  os.path.join(_BASE, "..", "certs", "iphone17-key.pem"))
HOST     = os.environ.get("TAK_HOST", "127.0.0.1")
PORT     = int(os.environ.get("TAK_PORT", "8089"))
BASE_LAT = 41.01865
BASE_LON = -73.64153

# ── Geometry ──────────────────────────────────────────────────────────────────
def move(lat, lon, dist_m, brg):
    R = 6371000
    d, b = dist_m / R, math.radians(brg)
    lat2 = math.asin(math.sin(math.radians(lat)) * math.cos(d) +
                     math.cos(math.radians(lat)) * math.sin(d) * math.cos(b))
    lon2 = math.radians(lon) + math.atan2(
        math.sin(b) * math.sin(d) * math.cos(math.radians(lat)),
        math.cos(d) - math.sin(math.radians(lat)) * math.sin(lat2))
    return math.degrees(lat2), math.degrees(lon2)

def brg(la1, lo1, la2, lo2):
    dl = math.radians(lo2 - lo1)
    y = math.sin(dl) * math.cos(math.radians(la2))
    x = (math.cos(math.radians(la1)) * math.sin(math.radians(la2)) -
         math.sin(math.radians(la1)) * math.cos(math.radians(la2)) * math.cos(dl))
    return (math.degrees(math.atan2(y, x)) + 360) % 360

def dist(la1, lo1, la2, lo2):
    R = 6371000
    return R * 2 * math.asin(math.sqrt(
        math.sin(math.radians((la2-la1)/2))**2 +
        math.cos(math.radians(la1)) * math.cos(math.radians(la2)) *
        math.sin(math.radians((lo2-lo1)/2))**2))

# ── CoT builders ─────────────────────────────────────────────────────────────
def ts(offset=0):
    return time.strftime("%Y-%m-%dT%H:%M:%SZ", time.gmtime(time.time() + offset))

def cot_reaper(lat, lon, alt, course, speed):
    # a-f-A-M-F-Q-r = Friendly / Air / Military / Fixed Wing / Attack / Rotary (UAV symbol)
    return (
        f'<event version="2.0" uid="REAPER-1" type="a-f-A-M-F-Q" how="m-g"'
        f' time="{ts()}" start="{ts()}" stale="{ts(20)}">'
        f'<point lat="{lat:.7f}" lon="{lon:.7f}" hae="{alt:.0f}" ce="10" le="9999999"/>'
        f'<detail>'
        f'<contact callsign="REAPER-1"/>'
        f'<uid Droid="REAPER-1"/>'
        f'<__group name="Red" role="Team Member"/>'
        f'<track speed="{speed:.1f}" course="{course:.1f}"/>'
        f'<remarks>OneMind ISR/Strike UAV</remarks>'
        f'</detail>'
        f'</event>'
    )

def cot_hostile(uid, callsign, lat, lon, ctype="a-h-G-U-C-I"):
    # a-h-G-U-C-I = Hostile / Ground / Unit / Combat / Infantry
    # a-h-G-U-C-A = Hostile / Ground / Unit / Combat / Armor
    return (
        f'<event version="2.0" uid="{uid}" type="{ctype}" how="m-g"'
        f' time="{ts()}" start="{ts()}" stale="{ts(300)}">'
        f'<point lat="{lat:.7f}" lon="{lon:.7f}" hae="0" ce="10" le="9999999"/>'
        f'<detail>'
        f'<contact callsign="{callsign}"/>'
        f'<uid Droid="{callsign}"/>'
        f'<__group name="Red" role="Team Member"/>'
        f'<track speed="1.2" course="180.0"/>'
        f'<remarks>HOSTILE — REAPER-1 ISR detection</remarks>'
        f'</detail>'
        f'</event>'
    )

def cot_strike(uid, callsign, lat, lon):
    return (
        f'<event version="2.0" uid="{uid}" type="b-a-o-tbl" how="h-e"'
        f' time="{ts()}" start="{ts()}" stale="{ts(300)}">'
        f'<point lat="{lat:.7f}" lon="{lon:.7f}" hae="0" ce="10" le="9999999"/>'
        f'<detail>'
        f'<contact callsign="{callsign}"/>'
        f'<uid Droid="{callsign}"/>'
        f'<__group name="Red" role="Team Member"/>'
        f'<remarks>STRIKE COMPLETE — neutralized by REAPER-1</remarks>'
        f'</detail>'
        f'</event>'
    )

def cot_delete(uid, ctype):
    return (
        f'<event version="2.0" uid="del-{uid}" type="t-x-d-d-r" how="m-g"'
        f' time="{ts()}" start="{ts()}" stale="{ts(30)}">'
        f'<point lat="0" lon="0" hae="0" ce="9999999" le="9999999"/>'
        f'<detail><link uid="{uid}" type="{ctype}" relation="p-p"/></detail>'
        f'</event>'
    )

# ── Connection ────────────────────────────────────────────────────────────────
ctx = ssl.create_default_context()
ctx.check_hostname = False
ctx.verify_mode = ssl.CERT_NONE
ctx.load_cert_chain(CERT, keyfile=KEY)
raw = socket.create_connection((HOST, PORT), timeout=10)
s   = ctx.wrap_socket(raw)

# Keepalive — holds subscription open
stop_ka = threading.Event()
def _ka():
    while not stop_ka.is_set():
        try:
            s.sendall((
                f'<event version="2.0" uid="KA-SIM" type="t-x-c-t" how="m-g"'
                f' time="{ts()}" start="{ts()}" stale="{ts(30)}">'
                f'<point lat="0" lon="0" hae="0" ce="9999999" le="9999999"/>'
                f'<detail/></event>'
            ).encode())
        except: pass
        time.sleep(8)
threading.Thread(target=_ka, daemon=True).start()

def send(xml):
    s.sendall(xml.encode('utf-8'))
    time.sleep(0.1)

def fly(from_lat, from_lon, to_lat, to_lon, alt, speed):
    d = dist(from_lat, from_lon, to_lat, to_lon)
    b = brg(from_lat, from_lon, to_lat, to_lon)
    steps = max(4, int(d / (speed * 2)))
    print(f"    flying {d:.0f}m hdg {b:.0f}°")
    for i in range(steps + 1):
        t = i / steps
        lat = from_lat + (to_lat - from_lat) * t
        lon = from_lon + (to_lon - from_lon) * t
        send(cot_reaper(lat, lon, alt, b, speed))
        time.sleep(2)
    return to_lat, to_lon

def orbit(clat, clon, radius, alt, speed, duration, label):
    print(f"    orbiting {label} — {duration}s")
    start  = time.time()
    angle  = 0.0
    dps    = (speed / (2 * math.pi * radius)) * 360 * 2
    while time.time() - start < duration:
        lat, lon  = move(clat, clon, radius, angle)
        course    = (angle + 90) % 360
        send(cot_reaper(lat, lon, alt, course, speed))
        angle = (angle + dps) % 360
        print(f"      {int(time.time()-start)}s/{duration}s", end='\r')
        time.sleep(2)
    print()
    return move(clat, clon, radius, angle)

# ── Wait for server to register ───────────────────────────────────────────────
print("Connecting...")
time.sleep(4)
print("Live.\n")

# ── AO CENTER = PAULINO's live position (entire demo runs around him) ─────────
CARKU_LAT = 40.70314
CARKU_LON = -73.98690
AO_LAT = CARKU_LAT
AO_LON = CARKU_LON

# ── 4 hostiles in a straight line east of the AO, evenly spaced ───────────────
# Line runs south-to-north, ~500m east of Paulino, 150m apart — a clean column
LINE_BEARING_FROM_AO = 90     # line sits due east of AO
LINE_RANGE           = 500    # 500m out
SPACING              = 150    # 150m between each hostile
# anchor point 500m east of AO
LA, LO = move(AO_LAT, AO_LON, LINE_RANGE, LINE_BEARING_FROM_AO)
# build column running north (bearing 0) through the anchor, centered
H1_LAT, H1_LON = move(LA, LO, SPACING * 1.5, 180)  # southernmost
H2_LAT, H2_LON = move(LA, LO, SPACING * 0.5, 180)
H3_LAT, H3_LON = move(LA, LO, SPACING * 0.5, 0)
H4_LAT, H4_LON = move(LA, LO, SPACING * 1.5, 0)    # northernmost

HOSTILES = [
    ("HOSTILE-01", "a-h-G-U-C-I", H1_LAT, H1_LON),  # infantry
    ("HOSTILE-02", "a-h-G-U-C-A", H2_LAT, H2_LON),  # armor
    ("HOSTILE-03", "a-h-G-U-C-I", H3_LAT, H3_LON),  # infantry
    ("HOSTILE-04", "a-h-G-U-C-A", H4_LAT, H4_LON),  # armor
]

r_lat = AO_LAT
r_lon = AO_LON

# ═══════════════════════════════════════════════════════════════════════════════
print("━━━ PHASE 1 — REAPER-1 DEPLOYING ━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
s_lat, s_lon = move(AO_LAT, AO_LON, 900, 270)  # inbound from west
r_lat, r_lon = fly(s_lat, s_lon, AO_LAT, AO_LON, 300, 50)
r_lat, r_lon = orbit(AO_LAT, AO_LON, 300, 250, 35, 10, "AO")

# ═══════════════════════════════════════════════════════════════════════════════
print("\n━━━ PHASE 2 — HOSTILES DETECTED ━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
for _ in range(4):
    for uid, ctype, hlat, hlon in HOSTILES:
        send(cot_hostile(uid, uid, hlat, hlon, ctype))
    time.sleep(1)
print("    4 hostiles on map — column formation, 500m east")
r_lat, r_lon = orbit(AO_LAT, AO_LON, 300, 250, 35, 22, "AO")

# ═══════════════════════════════════════════════════════════════════════════════
print("\n━━━ PHASE 3 — FIRE MISSION (4 TARGETS) ━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
# Engage the column south-to-north, one pass down the line
for idx, (uid, ctype, hlat, hlon) in enumerate(HOSTILES, start=1):
    print(f"  ENGAGING {uid}")
    r_lat, r_lon = fly(r_lat, r_lon, hlat, hlon, 80, 65)
    print("  *** STRIKE ***")
    sid = f"STRIKE-0{idx}"
    for _ in range(3):
        send(cot_strike(sid, f"{sid} [KIA]", hlat, hlon))
        time.sleep(0.4)
    send(cot_delete(uid, ctype))
    time.sleep(0.5)
    send(cot_delete(uid, ctype))
    print(f"  {uid} NEUTRALIZED ✓")
    time.sleep(2.5)

# ═══════════════════════════════════════════════════════════════════════════════
print("\n━━━ PHASE 4 — ISR ORBIT OVER PAULINO ━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
r_lat, r_lon = fly(r_lat, r_lon, CARKU_LAT, CARKU_LON, 200, 40)
orbit(CARKU_LAT, CARKU_LON, 150, 150, 20, 20, "PAULINO")

# Cleanup
send(cot_delete("REAPER-1",  "a-f-A-M-F-Q"))
for idx in range(1, 5):
    send(cot_delete(f"STRIKE-0{idx}", "b-a-o-tbl"))
send(cot_delete("KA-SIM",    "t-x-c-t"))

print("\n━━━ DEMO COMPLETE ━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
stop_ka.set()
s.close()
