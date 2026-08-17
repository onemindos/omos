#!/usr/bin/env python3
"""
Clear demo contacts from every connected iTAK/ATAK client.

Uses the official TAK delete CoT (type t-x-d-d-r) with a <link> to each target
UID. This is the ONLY reliable way to remove a contact from clients — sending a
stale/expired copy of the original event does NOT remove it (iTAK keeps its local
cache). Disconnect/reconnect does not clear local cache either.

Certs resolve to ../certs unless overridden via TAK_CERT / TAK_KEY env vars.
"""
import ssl, socket, time, os

_BASE = os.path.dirname(os.path.abspath(__file__))
CERT = os.environ.get("TAK_CERT", os.path.join(_BASE, "..", "certs", "iphone17-cert.pem"))
KEY  = os.environ.get("TAK_KEY",  os.path.join(_BASE, "..", "certs", "iphone17-key.pem"))
HOST = os.environ.get("TAK_HOST", "127.0.0.1")
PORT = int(os.environ.get("TAK_PORT", "8089"))

# Every UID the demos can spawn, with its CoT type (type must match for delete)
TARGETS = [
    ("REAPER-1",   "a-f-A-M-F-Q"),
    ("CHOPPER-1",  "a-f-A-M-H-Q"),
    ("HOSTILE-01", "a-h-G-U-C-I"),
    ("HOSTILE-02", "a-h-G-U-C-A"),
    ("HOSTILE-03", "a-h-G-U-C-I"),
    ("HOSTILE-04", "a-h-G-U-C-A"),
    ("STRIKE-01",  "b-a-o-tbl"),
    ("STRIKE-02",  "b-a-o-tbl"),
    ("STRIKE-03",  "b-a-o-tbl"),
    ("STRIKE-04",  "b-a-o-tbl"),
    ("KA-SIM",     "t-x-c-t"),
    ("KA-CHOP",    "t-x-c-t"),
]

def ts(off=0):
    return time.strftime("%Y-%m-%dT%H:%M:%SZ", time.gmtime(time.time() + off))

ctx = ssl.create_default_context()
ctx.check_hostname = False
ctx.verify_mode = ssl.CERT_NONE
ctx.load_cert_chain(CERT, keyfile=KEY)

print("Connecting...")
with socket.create_connection((HOST, PORT), timeout=10) as raw:
    with ctx.wrap_socket(raw) as s:
        time.sleep(3)
        for uid, ctype in TARGETS:
            xml = (
                f'<event version="2.0" uid="del-{uid}" type="t-x-d-d-r" how="m-g"'
                f' time="{ts()}" start="{ts()}" stale="{ts(30)}">'
                f'<point lat="0" lon="0" hae="0" ce="9999999" le="9999999"/>'
                f'<detail><link uid="{uid}" type="{ctype}" relation="p-p"/></detail>'
                f'</event>'
            )
            s.sendall(xml.encode('utf-8'))
            print(f"  deleted {uid}")
            time.sleep(0.2)
        time.sleep(3)
        print("Map clean")
