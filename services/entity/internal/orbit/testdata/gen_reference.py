"""Independent SGP4 reference for orbit_test.go (FEATURES A12 acceptance: within 5 km).

Uses python-sgp4 (MIT, Brandon Rhodes' port of Vallado's reference code) with WGS72, and a
closed-form (Bowring) TEME -> WGS84 geodetic conversion written separately from the Go one.
Run: uv run --with sgp4==2.27 python gen_reference.py stations.tle > reference.json (python-sgp4 2.27)
"""
import json, math, sys
from datetime import datetime, timezone
from sgp4.api import Satrec, WGS72

QUERIES = ["2026-10-28T14:18:15.123Z", "2026-10-28T18:00:00.400Z", "2026-10-29T12:00:00.000Z",
           "2026-10-30T03:33:33.700Z", "2026-11-01T00:00:00.000Z", "2026-11-03T21:07:00.000Z"]

def gmst(jd):
    t = (jd - 2451545.0) / 36525.0
    s = -6.2e-6 * t**3 + 0.093104 * t**2 + (876600.0 * 3600 + 8640184.812866) * t + 67310.54841
    return math.fmod(math.radians(s / 240.0), 2 * math.pi) % (2 * math.pi)

def geodetic(x, y, z, th):
    a, f = 6378.137, 1 / 298.257223563
    b = a * (1 - f); e2 = f * (2 - f); ep2 = (a * a - b * b) / (b * b)
    xe = math.cos(th) * x + math.sin(th) * y
    ye = -math.sin(th) * x + math.cos(th) * y
    p = math.hypot(xe, ye)
    tq = math.atan2(z * a, p * b)
    lat = math.atan2(z + ep2 * b * math.sin(tq) ** 3, p - e2 * a * math.cos(tq) ** 3)
    n = a / math.sqrt(1 - e2 * math.sin(lat) ** 2)
    return math.degrees(lat), math.degrees(math.atan2(ye, xe)), p / math.cos(lat) - n

def epoch_ms(l1):
    yy = int(l1[18:20]); year = 2000 + yy if yy < 57 else 1900 + yy
    jan1 = datetime(year, 1, 1, tzinfo=timezone.utc).timestamp() * 1000
    return int(jan1 + round((float(l1[20:32]) - 1) * 86400e3))

lines = [l.rstrip() for l in open(sys.argv[1] if len(sys.argv) > 1 else "stations.tle") if l.strip()]
out = []
for i in range(len(lines) - 1):
    l1, l2 = lines[i], lines[i + 1]
    if not (l1.startswith("1 ") and l2.startswith("2 ")):
        continue
    sat = Satrec.twoline2rv(l1, l2, WGS72)
    ep = epoch_ms(l1)
    for q in QUERIES:
        qms = int(datetime.fromisoformat(q.replace("Z", "+00:00")).timestamp() * 1000)
        t = ep + round((qms - ep) / 1000) * 1000
        jd = t / 86400e3 + 2440587.5
        e, r, _ = sat.sgp4(math.floor(jd - 0.5) + 0.5, jd - (math.floor(jd - 0.5) + 0.5))
        if e != 0:
            continue
        lat, lon, alt = geodetic(*r, gmst(jd))
        out.append({"norad": l1[2:7].strip(), "query": q, "t_ms": t,
                    "lat": round(lat, 6), "lon": round(lon, 6), "alt_km": round(alt, 3)})
json.dump(out, sys.stdout, indent=1)
