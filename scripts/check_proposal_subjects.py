#!/usr/bin/env python3
"""Check every NATS subject named in the grammar owner proposals against the locked grammar.

The proposal files (docs/proposals/2026-10-29-grammar-register/) list their subjects in fenced
``subjects`` blocks, one per line: ``<tag> <subject>``. Tags:

  proposed  the subject the proposal asks the owner to ratify
  interim   the subject used until the owner rules
  current   a live or already-built subject the proposal keeps unchanged
  filter    a stream subject filter (``evt.entity.>``)
  legacy    an out-of-grammar or to-be-retired subject, named so it can be ruled on
  alternative  a subject considered and not proposed (reported, not enforced)
  action    an action type name that looks like a subject (``task.assign``); not a subject

Checks (exit 1 on any error):
  * proposed, interim, current and filter subjects match a pattern in the omos ``CLAUDE.md``
    grammar block or ``docs/NATS-DOCTRINE.md`` sections 2 and 9.1;
  * proposed and interim subjects contain no stop-line word (policy service catalogue list);
  * proposed and interim ``cmd.fleet`` verbs are in the closed verb list (stop line SL5);
  * proposed and interim ``tel.*`` subjects have a class as token two (the BUS-P1 reading);
  * every subject-like inline code span in a file is listed in one of that file's blocks.

It also reports which stream on this checkout stores each subject today
(``platform/nats/streams/*.yaml``). The script only reads files; it changes nothing.

Usage: python scripts/check_proposal_subjects.py [--self-test] [FILE ...]
"""

from __future__ import annotations

import glob
import os
import re
import sys

REPO = os.path.dirname(os.path.dirname(os.path.abspath(__file__)))
DEFAULT_GLOB = os.path.join(REPO, "docs/proposals/2026-10-29-grammar-register/*.md")
STREAM_GLOB = os.path.join(REPO, "platform/nats/streams/*.yaml")

TAGS = ("proposed", "interim", "current", "filter", "legacy", "alternative", "action")
STRICT_TAGS = ("proposed", "interim", "current", "filter")
NEW_TAGS = ("proposed", "interim")

# {...} as a whole token means "one or more tokens" (svc endpoint, agents session tail, tak.cot.>).
PATTERNS = [
    ("G:tel", "tel.{node}.{class}.{metric}", "CLAUDE.md grammar"),
    ("G:det", "det.{domain}.{type}.{id}", "CLAUDE.md grammar"),
    ("G:evt", "evt.{domain}.{entity}.{event}", "CLAUDE.md grammar"),
    ("G:cmd", "cmd.{domain}.{entity}.{action}", "CLAUDE.md grammar"),
    ("G:intel", "intel.{domain}.{source}.{type}", "CLAUDE.md grammar"),
    ("G:ent", "ent.{entity_type}.{entity_id}.state", "CLAUDE.md grammar"),
    ("G:svc", "svc.{service}.{...}", "CLAUDE.md grammar; endpoint may contain dots (doctrine 2)"),
    ("G:agents", "agents.{prompt|hb|status|attachments}.{agent}.{owner}.{...}", "CLAUDE.md grammar"),
    ("G:tak-inject", "tak.cot.inject", "CLAUDE.md grammar"),
    ("G:tak-cot", "tak.cot.{...}", "CLAUDE.md grammar (tak.cot.>)"),
    ("D2:srv", "$SRV.{PING|INFO}.agents", "NATS-DOCTRINE 2"),
    ("D9.1:tel-fleet", "tel.fleet.{class}.{id}.{metric}", "NATS-DOCTRINE 9.1"),
    ("D9.1:cmd-fleet", "cmd.fleet.{class}.{id}.{verb}", "NATS-DOCTRINE 9.1"),
    ("D9.1:world-mil", "evt.world.air.mil.{id}", "NATS-DOCTRINE 9.1"),
    ("D9.1:world-news", "evt.world.news.{country}.{id}", "NATS-DOCTRINE 9.1"),
    ("D9.1:map-geofence", "cmd.map.geofence.{id}.{verb}", "NATS-DOCTRINE 9.1"),
    ("D9.1:map-annotation", "cmd.map.annotation.{id}.{verb}", "NATS-DOCTRINE 9.1"),
    ("D9.1:evt-geofence", "evt.map.geofence.{id}.{enter|exit}", "NATS-DOCTRINE 9.1"),
]

# Policy service catalogue (feat/policy-service services/policy/internal/policy/catalogue.go).
STOP_WORDS = {"strike", "engage", "engagement", "intercept", "designate", "designation",
              "attack", "kill", "weapon", "weapons", "hostile"}
# Stop line SL5 (onemind-xr docs/design/JADC2-MAP.md 5.2; asset-command.md 3).
SL5_VERBS = {"goto", "hold", "rtl", "land", "orbit", "observe", "estop", "arm", "disarm",
             "takeoff", "estop_release", "geofence_override"}
# BUS-P1 reading: token two of tel.* is a class from today's TEL/ENV/SYS filters, plus feed.
TEL_CLASSES = {"fleet", "drone", "robot", "edge", "iot", "industrial", "sat", "nature", "bio",
               "weather", "soil", "water", "host", "feed"}

SUBJECT_PREFIXES = {"tel", "det", "evt", "cmd", "intel", "ent", "svc", "agents", "tak", "saas",
                    "voice", "zeus", "op", "sys", "task", "omos", "fabric", "int", "$SRV"}
FILE_SUFFIXES = {"yaml", "yml", "md", "json", "go", "py", "fga", "sh", "txt"}
CODE_SPAN = re.compile(r"`([^`\s]+)`")
BLOCK = re.compile(r"^```subjects\s*$(.*?)^```\s*$", re.M | re.S)


def tokens(subject: str) -> list[str]:
    """Split on dots that are not inside {...}."""
    out, depth, cur = [], 0, ""
    for ch in subject:
        if ch == "{":
            depth += 1
        elif ch == "}":
            depth -= 1
        if ch == "." and depth == 0:
            out.append(cur)
            cur = ""
        else:
            cur += ch
    out.append(cur)
    return out


def is_var(tok: str) -> bool:
    return (tok.startswith("{") and tok.endswith("}")) or tok == "*"


def alternatives(tok: str) -> set[str] | None:
    """Literal alternatives of a token, or None for a free variable."""
    if tok == "*":
        return None
    if tok.startswith("{") and tok.endswith("}"):
        inner = tok[1:-1]
        return set(inner.split("|")) if "|" in inner else None
    return set(tok.split("|"))


def match_pattern(subject: str, pattern: str) -> bool:
    s, p = tokens(subject), tokens(pattern)
    if s and s[-1] == ">":
        s = s[:-1]
        prefix_only = True
    else:
        prefix_only = False
    for i, ptok in enumerate(p):
        if ptok == "{...}":
            return len(s) > i or prefix_only
        if i >= len(s):
            return prefix_only
        stok = s[i]
        palt, salt = alternatives(ptok), alternatives(stok)
        if palt is None:
            if stok in (">",):
                return False
            continue
        if salt is None or not salt <= palt:
            return False
    return len(s) == len(p)


def overlaps(a: str, b: str) -> bool:
    """True when some concrete subject matches both a and b (variables act as wildcards)."""
    ta, tb = tokens(a), tokens(b)
    for i in range(max(len(ta), len(tb))):
        x = ta[i] if i < len(ta) else None
        y = tb[i] if i < len(tb) else None
        if x == ">":
            return y is not None
        if y == ">":
            return x is not None
        if x is None or y is None:
            return False
        ax, ay = alternatives(x), alternatives(y)
        if ax is not None and ay is not None and not ax & ay:
            return False
    return True


def classify(subject: str) -> list[str]:
    return [pid for pid, pat, _ in PATTERNS if match_pattern(subject, pat)]


def literal_words(subject: str) -> set[str]:
    words = set()
    for tok in tokens(subject):
        alts = alternatives(tok)
        if alts:
            for a in alts:
                words.update(re.split(r"[_\-]", a.lower()))
                words.add(a.lower())
    return words


def load_streams() -> dict[str, list[str]]:
    streams = {}
    for path in sorted(glob.glob(STREAM_GLOB)):
        text = open(path, encoding="utf-8").read()
        if not re.search(r"^kind:\s*Stream\s*$", text, re.M):
            continue  # KeyValue buckets store $KV.* subjects only
        name = re.search(r"^spec:\s*\n\s+name:\s*\"?([A-Za-z0-9_-]+)", text, re.M)
        subj = re.search(r"^\s+subjects:\s*\n((?:\s+-\s*.+\n)+)", text, re.M)
        if not name or not subj:
            continue
        streams[name.group(1)] = re.findall(r"-\s*\"?([^\"\n]+?)\"?\s*$", subj.group(1), re.M)
    return streams


def stored_by(subject: str, streams: dict[str, list[str]]) -> list[str]:
    return sorted(n for n, fs in streams.items() if any(overlaps(subject, f) for f in fs))


def subject_like(span: str) -> bool:
    span = span.replace("\\|", "|")
    toks = tokens(span)
    if len(toks) < 2 or toks[0] not in SUBJECT_PREFIXES:
        return False
    if toks[-1] in FILE_SUFFIXES or "/" in span or "=" in span or "(" in span:
        return False
    return True


def check_file(path: str, streams: dict[str, list[str]]) -> tuple[list[str], list[str]]:
    text = open(path, encoding="utf-8").read()
    rows, errors, listed = [], [], set()
    rel = os.path.relpath(path, REPO)
    for block in BLOCK.finditer(text):
        base_line = text[: block.start()].count("\n") + 2
        for off, line in enumerate(block.group(1).strip("\n").splitlines()):
            line = line.strip()
            if not line or line.startswith("#"):
                continue
            parts = line.split()
            where = f"{rel}:{base_line + off}"
            if len(parts) < 2 or parts[0] not in TAGS:
                errors.append(f"{where}: bad subjects line {line!r}")
                continue
            tag, subj = parts[0], parts[1]
            listed.add(subj)
            pats = classify(subj)
            notes = []
            if tag in STRICT_TAGS and not pats:
                errors.append(f"{where}: {tag} subject {subj} matches no grammar pattern")
            hits = literal_words(subj) & STOP_WORDS
            if hits:
                msg = f"stop-line word {sorted(hits)}"
                notes.append(msg)
                if tag in NEW_TAGS:
                    errors.append(f"{where}: {tag} subject {subj} has {msg}")
            toks = tokens(subj)
            if tag in NEW_TAGS and toks[:2] == ["cmd", "fleet"] and len(toks) == 5:
                verbs = alternatives(toks[4])
                if verbs is not None and not verbs <= SL5_VERBS:
                    errors.append(f"{where}: verb(s) {sorted(verbs - SL5_VERBS)} not in SL5")
            if tag in NEW_TAGS and toks[0] == "tel" and len(toks) > 1:
                cls = alternatives(toks[1])
                if cls is not None and not cls <= TEL_CLASSES:
                    errors.append(f"{where}: tel token two {toks[1]} is not a class (BUS-P1)")
            store = ",".join(stored_by(subj, streams)) or "none"
            verdict = ",".join(pats) if pats else "OFF-GRAMMAR"
            rows.append(f"{where}\t{tag}\t{subj}\t{verdict}\tstored-today={store}"
                        + (f"\t{'; '.join(notes)}" if notes else ""))
    for n, line in enumerate(text.splitlines(), 1):
        if line.startswith("```"):
            continue
        for span in CODE_SPAN.findall(line):
            norm = span.replace("\\|", "|")
            if subject_like(norm) and norm not in listed:
                errors.append(f"{rel}:{n}: subject {norm} is not listed in a subjects block")
    return rows, errors


def self_test() -> int:
    cases = [
        ("tel.fleet.drone.{id}.position", True), ("tel.{node}.{class}.{metric}", True),
        ("cmd.fleet.{class}.{id}.{verb}", True), ("cmd.projector.disable", False),
        ("cmd.projector.{instance}.{disable|enable}", True), ("svc.policy.command.propose", True),
        ("evt.entity.>", True), ("intel.weather.noaa.alert.{state}.{event}", False),
        ("tak.cot.mesh.tx", True), ("ent.{entity_type}.{entity_id}.state", True),
        ("ent.manual.upsert", False), ("agents.discover", False), ("omos.air.{id}", False),
        ("evt.map.geofence.{id}.{enter|exit}", True), ("evt.map.geofence.{id}.created", False),
        ("$SRV.PING.agents", True), ("tel.feed.>", True), ("cmd.>", True),
    ]
    bad = 0
    for subj, want in cases:
        got = bool(classify(subj))
        if got != want:
            print(f"SELF-TEST FAIL {subj}: matched={got} want={want} ({classify(subj)})")
            bad += 1
    ov = [("evt.entity.{id}.created", "evt.so.>", False), ("evt.sys.{source}.rejected", "evt.sys.>", True),
          ("tel.feed.{source}.status", "tel.fleet.>", False), ("tak.cot.mesh.tx", "tak.cot.>", True),
          ("cmd.fleet.{class}.{id}.{verb}", "cmd.>", True), ("ent.{t}.{id}.state", "ent.*.*.state", True)]
    for a, b, want in ov:
        if overlaps(a, b) != want:
            print(f"SELF-TEST FAIL overlap {a} ~ {b}: want {want}")
            bad += 1
    if "intercept" not in literal_words("cmd.drone.x47.intercept"):
        print("SELF-TEST FAIL stop-line word not found")
        bad += 1
    print(f"self-test: {len(cases) + len(ov) + 1 - bad} passed, {bad} failed")
    return 1 if bad else 0


def main(argv: list[str]) -> int:
    if "--self-test" in argv:
        rc = self_test()
        argv = [a for a in argv if a != "--self-test"]
        if rc:
            return rc
    files = argv or sorted(glob.glob(DEFAULT_GLOB))
    if not files:
        print("no proposal files found")
        return 1
    streams = load_streams()
    print(f"streams on this checkout: {len(streams)} ({', '.join(sorted(streams))})")
    all_rows, all_errors = [], []
    for f in files:
        rows, errors = check_file(f, streams)
        all_rows += rows
        all_errors += errors
    print("location\ttag\tsubject\tgrammar\tstored\tnotes")
    for r in all_rows:
        print(r)
    counts = {t: sum(1 for r in all_rows if r.split("\t")[1] == t) for t in TAGS}
    off = sum(1 for r in all_rows if "\tOFF-GRAMMAR\t" in r)
    print(f"\nsubjects checked: {len(all_rows)} ({', '.join(f'{k} {v}' for k, v in counts.items())}); "
          f"off-grammar (legacy, alternative or action only): {off}")
    for e in all_errors:
        print(f"ERROR {e}")
    print(f"result: {'FAIL' if all_errors else 'PASS'} ({len(all_errors)} errors)")
    return 1 if all_errors else 0


if __name__ == "__main__":
    sys.exit(main(sys.argv[1:]))
