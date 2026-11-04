#!/usr/bin/env python3
"""Tailnet health helper for the OneMind agent. Standard library only.

  mesh-health.py probe [--target host:port ...] [--timeout 4]
      TCP-connect to the cluster egress Services. Silent on success.
  mesh-health.py peers FILE|- [--offline-threshold N] [--expect-offline NAME ...]
      Summarize `tailscale status --json` (a file, or "-" for stdin).
      Prints names and OS only, never addresses.

Exit 0 = nothing to report, 1 = at least one finding, 2 = bad usage/input.
Empty output on green so it can run as a silent cron watchdog.
"""
import argparse
import json
import socket
import sys
from collections import defaultdict

DEFAULT_TARGETS = [
    "omos-mbp.hermes.svc.cluster.local:22",
    "omos-s24.hermes.svc.cluster.local:8022",
    "omos-ha-edge.hermes.svc.cluster.local:8123",
]


def probe(targets, timeout):
    findings = []
    for target in targets:
        host, _, port = target.rpartition(":")
        if not host or not port.isdigit():
            print(f"bad target {target!r}, expected host:port", file=sys.stderr)
            return 2, []
        try:
            with socket.create_connection((host, int(port)), timeout=timeout):
                pass
        except socket.gaierror:
            findings.append(f"{target}: DNS failure (Service missing?)")
        except socket.timeout:
            findings.append(f"{target}: timeout (NetworkPolicy, tailnet grant or proxy down)")
        except ConnectionRefusedError:
            findings.append(f"{target}: refused (path works, service on the device is down)")
        except OSError as exc:
            findings.append(f"{target}: {exc.__class__.__name__}: {exc}")
    return (1 if findings else 0), findings


def peers(source, threshold, expected_offline):
    try:
        if source == "-":
            text = sys.stdin.read()
        else:
            with open(source, encoding="utf-8") as handle:
                text = handle.read()
        data = json.loads(text)
    except (OSError, ValueError) as exc:
        print(f"cannot read status JSON: {exc}", file=sys.stderr)
        return 2, []
    findings = []
    state = data.get("BackendState", "unknown")
    if state != "Running":
        findings.append(f"backend state is {state}")
    by_os = defaultdict(lambda: [0, 0])  # online, offline
    unexpected_offline = []
    relay_only = []
    for peer in (data.get("Peer") or {}).values():
        name = peer.get("HostName") or peer.get("DNSName", "?").split(".")[0]
        os_name = (peer.get("OS") or "unknown").lower()
        if peer.get("Online"):
            by_os[os_name][0] += 1
            if peer.get("Relay") and not peer.get("CurAddr"):
                relay_only.append(name)
        else:
            by_os[os_name][1] += 1
            if name not in expected_offline:
                unexpected_offline.append(name)
    if unexpected_offline:
        findings.append("offline: " + ", ".join(sorted(unexpected_offline)))
    if len(unexpected_offline) > threshold:
        findings.append(f"{len(unexpected_offline)} peers offline (threshold {threshold})")
    if relay_only:
        findings.append("relay only (no direct path): " + ", ".join(sorted(relay_only)))
    if findings:
        findings.append(
            "by OS (online/offline): "
            + ", ".join(f"{k} {v[0]}/{v[1]}" for k, v in sorted(by_os.items()))
        )
    return (1 if findings else 0), findings


def main(argv=None):
    parser = argparse.ArgumentParser(
        description=__doc__, formatter_class=argparse.RawDescriptionHelpFormatter
    )
    sub = parser.add_subparsers(dest="cmd", required=True)
    p = sub.add_parser("probe")
    p.add_argument("--target", action="append", help="host:port (repeatable); default: the three egress Services")
    p.add_argument("--timeout", type=float, default=4.0)
    s = sub.add_parser("peers")
    s.add_argument("source", help="status JSON file, or - for stdin")
    s.add_argument("--offline-threshold", type=int, default=5)
    s.add_argument("--expect-offline", action="append", default=[], metavar="NAME")
    args = parser.parse_args(argv)
    if args.cmd == "probe":
        code, findings = probe(args.target or DEFAULT_TARGETS, args.timeout)
    else:
        code, findings = peers(args.source, args.offline_threshold, set(args.expect_offline))
    for line in findings:
        print(line)
    return code


if __name__ == "__main__":
    sys.exit(main())
