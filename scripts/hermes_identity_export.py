#!/usr/bin/env python3
"""Export Hermes' identity files from the live data volume as a redacted tar stream on stdout.

Runs INSIDE the Hermes pod (read-only; it never writes to the volume):

    kubectl -n hermes exec -i <pod> -c hermes -- /opt/hermes/.venv/bin/python - < scripts/hermes_identity_export.py > export.tar

The rules are in docs/operations/hermes-export-rules.md. In short:
  * skills: export the ones we created (rules R1 to R8), skip upstream and hub skills;
  * config.yaml: parse, redact by path, key name and value shape, fail closed on anything unknown;
  * every exported file is scanned for credentials; any finding fails the export.

On failure the tar holds only EXPORT-REPORT.md (paths and rule names, never values) and the exit code is 2.
Standard library plus PyYAML (already in the Hermes venv).
"""
import collections
import hashlib
import io
import json
import os
import re
import sys
import tarfile

import yaml

DATA_DIR = os.environ.get("HERMES_DATA_DIR", "/opt/data")
IMAGE_SKILL_ROOTS = [p for p in os.environ.get(
    "HERMES_IMAGE_SKILL_ROOTS", "/opt/hermes/skills:/opt/hermes/optional-skills").split(":") if p]
MAX_FILE_BYTES = 512 * 1024
SKIP_NAMES = {".DS_Store"}
SKIP_DIRS = {"__pycache__", "node_modules"}

OURS_AUTHORS = {"legacy", "legacy (mother)", "onemind os / zeus"}

# Config rules ---------------------------------------------------------------------------------
EXPLICIT_PATHS = (
    "mcp.servers.notion.args",
    "mcp_servers.notion.args",
    "mcp_servers.github.headers.Authorization",
    "mcp_servers.wordpress.headers.Authorization",
    "mcp_servers.composio.headers.x-consumer-api-key",
    "mcp_servers.ghl.headers.Authorization",
    "mcp_servers.ghl-pro.env.GHL_API_KEY",
    "platforms.matrix.token",
    "dashboard.basic_auth.password_hash",
    "dashboard.basic_auth.secret",
)
IDENTIFIER_PATHS = ("TELEGRAM_ALLOWED_USERS", "platforms.matrix.user_id")
SENSITIVE_KEY = re.compile(r"token|secret|password|passwd|api[_-]?key|authorization|credential|cookie|bearer", re.I)
# A value that is clearly a reference and not a secret: op:// ref, ${PLACEHOLDER}, an env var name, oauth, bool, empty.
SAFE_VALUE = re.compile(r"^(op://\S+|\$\{[A-Za-z0-9_]+\}|[A-Z][A-Z0-9]*_[A-Z0-9_]+|(?i:oauth|true|false|none|null))?$")
VALUE_SHAPES = re.compile(
    r"ntn_[A-Za-z0-9]{6,}|ghp_[A-Za-z0-9]{10,}|gho_[A-Za-z0-9]{10,}|github_pat_[A-Za-z0-9_]{10,}"
    r"|ck_[A-Za-z0-9_-]{8,}|sk-[A-Za-z0-9_-]{16,}|AKIA[A-Z0-9]{12,}|tskey-[A-Za-z0-9-]{8,}"
    r"|eyJ[A-Za-z0-9_-]{10,}\.[A-Za-z0-9_-]{5,}(?:\.[A-Za-z0-9_-]*)?"
    r"|(?:Bearer|Basic)\s+[A-Za-z0-9._~+/=-]{8,}"
    r"|\b\d{8,12}:[A-Za-z0-9_-]{30,}\b"
    r"|-----BEGIN [A-Z ]*PRIVATE KEY-----")
# Credential-shaped assignment, used when scanning skills and other text.
ASSIGNMENT = re.compile(
    r"(password|passwd|secret|api[_-]?key|token|app_pass|bearer|authorization)[\"' ]*[:=][\"' ]*([A-Za-z0-9_\-./+]{8,})", re.I)
PLACEHOLDER_VALUE = re.compile(
    r"^(\$\{?[A-Za-z_][A-Za-z0-9_]*\}?|[A-Z][A-Z0-9]*_[A-Z0-9_]+|<[^>]+>"
    r"|(?i:your[_-].*|x{3,}|\*+|changeme|example|placeholder|token|secret|password))$")


def placeholder(path):
    return "${" + re.sub(r"[^A-Za-z0-9]+", "_", ".".join(str(p) for p in path)).strip("_").upper() + "}"


# Skills ----------------------------------------------------------------------------------------
def skill_dirs(root):
    found = {}
    for r, ds, fs in os.walk(root):
        ds[:] = sorted(d for d in ds if not d.startswith(".") and d not in SKIP_DIRS)
        if "SKILL.md" in fs:
            found.setdefault(os.path.basename(r), r)
    return found


def tree_hash(path):
    h = hashlib.sha256()
    for r, ds, fs in os.walk(path):
        ds[:] = sorted(d for d in ds if d not in SKIP_DIRS)
        for f in sorted(fs):
            fp = os.path.join(r, f)
            h.update(os.path.relpath(fp, path).encode())
            with open(fp, "rb") as fh:
                h.update(fh.read())
    return h.hexdigest()


def read_json(path, default):
    try:
        with open(path, encoding="utf-8") as fh:
            return json.load(fh)
    except (OSError, ValueError):
        return default


def frontmatter_author(skill_path):
    try:
        with open(os.path.join(skill_path, "SKILL.md"), encoding="utf-8", errors="ignore") as fh:
            head = fh.read(4000)
    except OSError:
        return ""
    m = re.search(r"^author:\s*(.+)$", head, re.M)
    return m.group(1).strip().strip("\"'") if m else ""


def classify_skills(data_dir, image_roots):
    """Return (export, review, skipped) where export/review are lists of (name, path, why) and skipped is a Counter."""
    skills_dir = os.path.join(data_dir, "skills")
    usage = read_json(os.path.join(skills_dir, ".usage.json"), {})
    hub = set((read_json(os.path.join(skills_dir, ".hub", "lock.json"), {}) or {}).get("installed", {}))
    image = {}
    for root in image_roots:
        for name, path in skill_dirs(root).items():
            image.setdefault(name, path)
    export, review, skipped = [], [], collections.Counter()
    for name, path in sorted(skill_dirs(skills_dir).items()):
        entry = usage.get(name) or {}
        created_by = entry.get("created_by")
        patches = entry.get("patch_count") or 0
        author = frontmatter_author(path)
        in_image = name in image
        if in_image and tree_hash(path) == tree_hash(image[name]):
            skipped["R1 identical to the image"] += 1
        elif name in hub or created_by == "installed":
            skipped["R2 hub-installed"] += 1
        elif created_by in ("agent", "learn"):
            export.append((name, path, "R3 created_by=%s" % created_by))
        elif author.lower() in OURS_AUTHORS:
            export.append((name, path, "R4 author=%s" % author))
        elif in_image and patches > 0:
            export.append((name, path, "R5 upstream, patched x%d" % patches))
        elif in_image:
            skipped["R6 upstream, older version only"] += 1
        elif not author:
            export.append((name, path, "R7 no author, not in image"))
        else:
            review.append((name, path, "R8 author=%s" % author))
    return export, review, skipped


# Scanning --------------------------------------------------------------------------------------
def scan_text(text):
    """Return a list of (line_number, rule) for credential-shaped content. Never returns the value."""
    findings = []
    for i, line in enumerate(text.splitlines(), 1):
        if VALUE_SHAPES.search(line):
            findings.append((i, "secret-shaped value"))
            continue
        for m in ASSIGNMENT.finditer(line):
            if not PLACEHOLDER_VALUE.match(m.group(2)):
                findings.append((i, "credential-shaped assignment"))
                break
    return findings


# Config ----------------------------------------------------------------------------------------
def _under(path_str, prefixes):
    return any(path_str == p or path_str.startswith(p + ".") for p in prefixes)


def redact_config(obj):
    """Return (redacted_copy, redactions, problems). redactions are (path, placeholder); problems are strings."""
    redactions, problems = [], []

    def walk(o, path):
        p = ".".join(path)
        if isinstance(o, dict):
            return {k: walk(v, path + [str(k)]) for k, v in o.items()}
        if isinstance(o, list):
            return [walk(v, path + ["[%d]" % i]) for i, v in enumerate(o)]
        if not isinstance(o, str):
            if path and _under(p, IDENTIFIER_PATHS):
                redactions.append((p, placeholder(path)))
                return placeholder(path)
            return o
        flat = re.sub(r"\.?\[\d+\]", "", p)
        in_list = bool(path) and path[-1].startswith("[")
        # An explicit path holding a scalar is replaced whole. Items of a list are only cleaned by value
        # shape below, so flags such as "-y" survive and just the secret inside an argument goes.
        if _under(flat, IDENTIFIER_PATHS) or (flat in EXPLICIT_PATHS and not in_list):
            redactions.append((p, placeholder(path)))
            return placeholder(path)
        if VALUE_SHAPES.search(o):
            redactions.append((p, placeholder(path)))
            return VALUE_SHAPES.sub(placeholder(path), o)
        key = path[-1] if path else ""
        if SENSITIVE_KEY.search(key) and not SAFE_VALUE.match(o):
            problems.append("no rule for sensitive-looking key: %s" % p)
        return o

    return walk(obj, []), redactions, problems


def dump_yaml(obj):
    return yaml.safe_dump(obj, sort_keys=False, default_flow_style=False, allow_unicode=True)


# Collection ------------------------------------------------------------------------------------
def collect_files(path):
    """Yield (relative_path, bytes, note) for the regular files under a skill directory."""
    for r, ds, fs in os.walk(path):
        ds[:] = sorted(d for d in ds if d not in SKIP_DIRS and not d.startswith("."))
        for f in sorted(fs):
            fp = os.path.join(r, f)
            rel = os.path.relpath(fp, path)
            if f in SKIP_NAMES or os.path.islink(fp):
                continue
            if os.path.getsize(fp) > MAX_FILE_BYTES:
                yield rel, None, "skipped: larger than %d KB" % (MAX_FILE_BYTES // 1024)
                continue
            with open(fp, "rb") as fh:
                yield rel, fh.read(), ""


def build_export(data_dir=DATA_DIR, image_roots=IMAGE_SKILL_ROOTS):
    """Return (files, ok, report). files maps repo-relative path to bytes (report only when not ok)."""
    files, problems, redactions, skipped_files = {}, [], [], []

    export, review, skipped = classify_skills(data_dir, image_roots)
    for name, path, why in export:
        rel_skill = os.path.relpath(path, os.path.join(data_dir, "skills"))
        for rel, data, note in collect_files(path):
            target = "skills/%s/%s" % (rel_skill, rel)
            if data is None:
                skipped_files.append("%s (%s)" % (target, note))
                continue
            for line, rule in scan_text(data.decode("utf-8", errors="ignore")):
                problems.append("%s:%d: %s" % (target, line, rule))
            files[target] = data

    def add_text(src, target):
        fp = os.path.join(data_dir, src)
        if os.path.isfile(fp) and not os.path.islink(fp):
            with open(fp, "rb") as fh:
                data = fh.read()
            for line, rule in scan_text(data.decode("utf-8", errors="ignore")):
                problems.append("%s:%d: %s" % (target, line, rule))
            files[target] = data

    add_text("SOUL.md", "live/SOUL.md")
    add_text(os.path.join("cron", "jobs.json"), "cron/jobs.json")

    def add_config(src, target):
        fp = os.path.join(data_dir, src)
        if not (os.path.isfile(fp) and not os.path.islink(fp)):
            return
        with open(fp, encoding="utf-8") as fh:
            loaded = yaml.safe_load(fh) or {}
        redacted, reds, probs = redact_config(loaded)
        text = dump_yaml(redacted)
        for line, rule in scan_text(text):
            probs.append("%s: redacted output still has a %s (line %d)" % (target, rule, line))
        redactions.extend((target, p, ph) for p, ph in reds)
        problems.extend("%s: %s" % (target, pr) for pr in probs)
        files[target] = text.encode("utf-8")

    add_config("config.yaml", "config/config.template.yaml")
    profiles = os.path.join(data_dir, "profiles")
    if os.path.isdir(profiles):
        for name in sorted(os.listdir(profiles)):
            add_config(os.path.join("profiles", name, "config.yaml"),
                       "profiles/%s/config.template.yaml" % name)

    ok = not problems
    lines = ["# Hermes identity export report", "", "Status: **%s**" % ("OK" if ok else "FAILED"), ""]
    lines += ["Skills: %d exported, %d for review, %d skipped." % (len(export), len(review), sum(skipped.values())), ""]
    if skipped:
        lines += ["Skipped:"] + ["- %s: %d" % (k, v) for k, v in sorted(skipped.items())] + [""]
    lines += ["## Exported skills", ""] + ["- `%s` (%s)" % (n, why) for n, _, why in export] + [""]
    lines += ["## Needs review (not exported)", ""] + (["- `%s` (%s)" % (n, why) for n, _, why in review] or ["- none"]) + [""]
    lines += ["## Config redactions (paths only)", ""]
    lines += ["- `%s` `%s` -> `%s`" % r for r in redactions] or ["- none"]
    lines += ["", "## Files skipped for size", ""] + (["- %s" % s for s in skipped_files] or ["- none"])
    if problems:
        lines += ["", "## Problems (the export was refused)", ""] + ["- %s" % p for p in problems]
        lines += ["", "Rotate or remove the credentials at the source, or add a rule, then run again."]
    report = "\n".join(lines) + "\n"
    if not ok:
        return {"EXPORT-REPORT.md": report.encode("utf-8")}, False, report
    files["EXPORT-REPORT.md"] = report.encode("utf-8")
    return files, True, report


def write_tar(files, stream):
    with tarfile.open(fileobj=stream, mode="w|") as tar:
        for name in sorted(files):
            info = tarfile.TarInfo(name)
            info.size = len(files[name])
            info.mtime = 0
            info.mode = 0o644
            tar.addfile(info, io.BytesIO(files[name]))


def main():
    files, ok, _ = build_export()
    write_tar(files, sys.stdout.buffer)
    sys.stdout.buffer.flush()
    return 0 if ok else 2


if __name__ == "__main__":
    sys.exit(main())
