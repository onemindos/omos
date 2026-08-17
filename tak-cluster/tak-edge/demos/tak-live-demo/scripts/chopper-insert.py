#!/usr/bin/env python3
"""
OneMind — Chopper insert demo
CHOPPER-1 (friendly rotary) flies from Velipoje to the marker.
Persistent SSL connection, holds subscription open, cinematic ~60s flight.
"""
import ssl, socket, time, math, threading, os

_BASE = os.path.dirname(os.path.abspath(__file__))
CERT = os.environ.get("TAK_CERT", os.path.join(_BASE, "..", "certs", "iphone17-cert.pem"))
KEY  = os.environ.get("TAK_KEY",  os.path.join(_BASE, "..", "certs", "iphone17-key.pem"))
HOST = os.environ.get("TAK_HOST", "127.0.0.1")
PORT = int(os.environ.get("TAK_PORT", "8089"))

# Start / end
VELIPOJE_LAT, VELIPOJE_LON = 41.8647, 19.4300
MARKER_LAT,   MARKER_LON   = 41.5788, 19.4659

SPEED_MPH      = 160
SPEED_MS       = SPEED_MPH * 0.44704   # 160 mph = 71.53 m/s
TICK           = 2.0     # update interval

def bearing(la1, lo1, la2, lo2):
    dl = math.radians(lo2 - lo1)
    y = math.sin(dl) * math.cos(math.radians(la2))
    x = (math.cos(math.radians(la1)) * math.sin(math.radians(la2)) -
         math.sin(math.radians(la1)) * math.cos(math.radians(la2)) * math.cos(dl))
    return (math.degrees(math.atan2(y, x)) + 360) % 360

def dist_m(la1, lo1, la2, lo2):
    R = 6371000
    return R * 2 * math.asin(math.sqrt(
        math.sin(math.radians((la2-la1)/2))**2 +
        math.cos(math.radians(la1)) * math.cos(math.radians(la2)) *
        math.sin(math.radians((lo2-lo1)/2))**2))

def ts(off=0):
    return time.strftime("%Y-%m-%dT%H:%M:%SZ", time.gmtime(time.time() + off))

def cot_chopper(lat, lon, alt, course, speed):
    # a-f-A-M-H-Q = Friendly / Air / Military / Rotary Wing (helicopter) — blue
    return (
        f'<event version="2.0" uid="CHOPPER-1" type="a-f-A-M-H-Q" how="m-g"'
        f' time="{ts()}" start="{ts()}" stale="{ts(20)}">'
        f'<point lat="{lat:.7f}" lon="{lon:.7f}" hae="{alt:.0f}" ce="10" le="9999999"/>'
        f'<detail>'
        f'<contact callsign="CHOPPER-1"/>'
        f'<uid Droid="CHOPPER-1"/>'
        f'<__group name="Red" role="Team Member"/>'
        f'<track speed="{speed:.1f}" course="{course:.1f}"/>'
        f'<remarks>OneMind rotary insert — Velipoje inbound</remarks>'
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

ctx = ssl.create_default_context()
ctx.check_hostname = False
ctx.verify_mode = ssl.CERT_NONE
ctx.load_cert_chain(CERT, keyfile=KEY)
raw = socket.create_connection((HOST, PORT), timeout=10)
s   = ctx.wrap_socket(raw)

stop_ka = threading.Event()
def _ka():
    while not stop_ka.is_set():
        try:
            s.sendall((
                f'<event version="2.0" uid="KA-CHOP" type="t-x-c-t" how="m-g"'
                f' time="{ts()}" start="{ts()}" stale="{ts(30)}">'
                f'<point lat="0" lon="0" hae="0" ce="9999999" le="9999999"/>'
                f'<detail/></event>'
            ).encode())
        except: pass
        time.sleep(8)
threading.Thread(target=_ka, daemon=True).start()

def send(xml):
    s.sendall(xml.encode('utf-8'))

print("Connecting...")
time.sleep(4)
print("Live.\n")

total_dist     = dist_m(VELIPOJE_LAT, VELIPOJE_LON, MARKER_LAT, MARKER_LON)
course         = bearing(VELIPOJE_LAT, VELIPOJE_LON, MARKER_LAT, MARKER_LON)
FLIGHT_SECONDS = total_dist / SPEED_MS      # real time at SPEED_MPH
steps          = max(1, int(FLIGHT_SECONDS / TICK))

print(f"━━━ CHOPPER-1 INBOUND ━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
print(f"    Velipoje → Marker  |  {total_dist/1000:.1f} km  |  hdg {course:.0f}°  |  {SPEED_MPH} mph  |  {FLIGHT_SECONDS/60:.1f} min")

for i in range(steps + 1):
    t   = i / steps
    lat = VELIPOJE_LAT + (MARKER_LAT - VELIPOJE_LAT) * t
    lon = VELIPOJE_LON + (MARKER_LON - VELIPOJE_LON) * t
    send(cot_chopper(lat, lon, 150, course, SPEED_MS))
    pct = int(t * 100)
    print(f"    inbound {pct:3d}%  ({lat:.4f}, {lon:.4f})", end='\r')
    time.sleep(TICK)
print()

# Hold on target — hover at marker for 10s
print("    CHOPPER-1 ON TARGET — holding over marker (10s)")
hold_start = time.time()
while time.time() - hold_start < 10:
    send(cot_chopper(MARKER_LAT, MARKER_LON, 80, course, 0))
    time.sleep(2)

print("\n━━━ CHOPPER-1 INSERT COMPLETE ━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
# Leave chopper on target, or clean up. Keep it for now.
stop_ka.set()
send(cot_delete("KA-CHOP", "t-x-c-t"))
s.close()
