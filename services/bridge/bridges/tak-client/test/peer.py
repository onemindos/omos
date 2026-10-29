#!/usr/bin/env python3
# SPDX-License-Identifier: Apache-2.0
"""Scripted TAK Server streaming peer for testing tak-client end to end.

An implementation independent of the Go code: TLS with a required client
certificate, XML streaming, the t-x-takp-v/q/r negotiation and TAK Protocol v1
streaming frames encoded with the `takproto` package (MIT,
https://github.com/snstac/takproto), following atak-civ takproto/README.txt.

It sends --send EUD position reports after the client connects (and, in v1
mode, after negotiation), then prints one JSON line per event the client sends
(uid, type, proto) until --expect events of a non-control type arrived or
--timeout passed. Exit code 0 when the expectation was met.

    python -I peer.py --port 18089 --cert srv.pem --key srv.key --ca ca.pem --mode v1 --send 2 --expect 1
"""
import argparse
import json
import socket
import ssl
import sys
import time
import xml.etree.ElementTree as ET
from datetime import datetime, timedelta, timezone

import takproto
from takproto.proto import TakMessage

DECL = b'<?xml version="1.0" encoding="UTF-8" standalone="yes"?>\n'
CONTROL = ("t-x-takp-v", "t-x-takp-q", "t-x-takp-r", "t-x-c-t", "t-x-c-t-r")


def iso(dt):
    return dt.strftime("%Y-%m-%dT%H:%M:%S.%f")[:-3] + "Z"


def event_xml(uid, cot_type, detail="", stale_s=120, callsign=None):
    now = datetime.now(timezone.utc)
    if callsign:
        detail = (f'<contact endpoint="*:-1:stcp" callsign="{callsign}"/>'
                  '<__group name="Cyan" role="Team Member"/>'
                  '<takv device="Pixel" platform="ATAK-CIV" os="34" version="5.6.0"/>'
                  '<track speed="1.5" course="90"/>' + detail)
    return (f'<event version="2.0" uid="{uid}" type="{cot_type}" how="m-g" time="{iso(now)}" '
            f'start="{iso(now)}" stale="{iso(now + timedelta(seconds=stale_s))}">'
            '<point lat="30.2672" lon="-97.7431" hae="152.4" ce="9.9" le="9999999"/>'
            f'<detail>{detail}</detail></event>').encode()


class Peer:
    def __init__(self, conn, out):
        self.conn = conn
        self.out = out
        self.proto = 0
        self.buf = b""

    def send(self, xml):
        if self.proto == 1:
            self.conn.sendall(bytes(takproto.xml2proto(xml, takproto.TAKProtoVer.STREAM)))
        else:
            self.conn.sendall(DECL + xml)

    def _fill(self):
        chunk = self.conn.recv(65536)
        if not chunk:
            raise EOFError("client closed")
        self.buf += chunk

    def next_event(self):
        """Returns (uid, type, raw) for the next client event in the current framing."""
        if self.proto == 0:
            while b"</event>" not in self.buf:
                self._fill()
            end = self.buf.index(b"</event>") + len(b"</event>")
            raw, self.buf = self.buf[:end], self.buf[end:]
            raw = raw[raw.index(b"<event"):]
            ev = ET.fromstring(raw)
            return ev.get("uid"), ev.get("type"), raw
        while len(self.buf) < 2:
            self._fill()
        if self.buf[0] != 0xBF:
            raise ValueError(f"v1 frame starts with {self.buf[0]:#x}")
        # varint length
        n, shift, i = 0, 0, 1
        while True:
            while len(self.buf) <= i:
                self._fill()
            b = self.buf[i]
            n |= (b & 0x7F) << shift
            i += 1
            if b < 0x80:
                break
            shift += 7
        while len(self.buf) < i + n:
            self._fill()
        payload, self.buf = self.buf[i:i + n], self.buf[i + n:]
        m = TakMessage()
        m.ParseFromString(bytes(payload))
        return m.cotEvent.uid, m.cotEvent.type, payload

    def emit(self, uid, cot_type):
        print(json.dumps({"uid": uid, "type": cot_type, "proto": self.proto}), file=self.out, flush=True)


def main():
    ap = argparse.ArgumentParser()
    ap.add_argument("--port", type=int, required=True)
    ap.add_argument("--cert", required=True)
    ap.add_argument("--key", required=True)
    ap.add_argument("--ca", required=True)
    ap.add_argument("--mode", choices=["xml", "v1"], default="xml")
    ap.add_argument("--send", type=int, default=1)
    ap.add_argument("--expect", type=int, default=1)
    ap.add_argument("--timeout", type=float, default=30)
    a = ap.parse_args()

    ctx = ssl.SSLContext(ssl.PROTOCOL_TLS_SERVER)
    ctx.load_cert_chain(a.cert, a.key)
    ctx.load_verify_locations(a.ca)
    ctx.verify_mode = ssl.CERT_REQUIRED
    srv = socket.create_server(("127.0.0.1", a.port))
    srv.settimeout(a.timeout)
    print(json.dumps({"listening": a.port, "mode": a.mode}), flush=True)
    raw, _ = srv.accept()
    conn = ctx.wrap_socket(raw, server_side=True)
    conn.settimeout(a.timeout)
    cn = dict(x[0] for x in conn.getpeercert()["subject"]).get("commonName")
    print(json.dumps({"client_cn": cn, "tls": conn.version()}), flush=True)
    p = Peer(conn, sys.stdout)

    if a.mode == "v1":
        p.send(event_xml("protouid", "t-x-takp-v",
                         '<TakControl><TakProtocolSupport version="1"/>'
                         '<TakServerVersionInfo serverVersion="peer.py-takproto" apiVersion="3"/></TakControl>', 60))
        while True:
            uid, cot_type, _ = p.next_event()
            p.emit(uid, cot_type)
            if cot_type == "t-x-takp-q":
                p.send(event_xml("protouid", "t-x-takp-r", '<TakControl><TakResponse status="true"/></TakControl>', 60))
                p.proto = 1
                print(json.dumps({"negotiated": 1}), flush=True)
                break

    for i in range(a.send):
        p.send(event_xml(f"ANDROID-peer-{i}", "a-f-G-U-C", callsign=f"PEER{i}"))
    print(json.dumps({"sent": a.send, "proto": p.proto}), flush=True)

    got = 0
    deadline = time.time() + a.timeout
    try:
        while got < a.expect and time.time() < deadline:
            uid, cot_type, _ = p.next_event()
            p.emit(uid, cot_type)
            if cot_type == "t-x-c-t":
                p.send(event_xml("takPong", "t-x-c-t-r", "", 20))
            if cot_type not in CONTROL:
                got += 1
    except (socket.timeout, EOFError) as e:
        print(json.dumps({"error": str(e)}), flush=True)
    print(json.dumps({"received_non_control": got}), flush=True)
    time.sleep(0.5)
    conn.close()
    sys.exit(0 if got >= a.expect else 1)


if __name__ == "__main__":
    main()
