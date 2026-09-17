#!/usr/bin/env python3
"""
verify-duplicates.py
====================
Byte-for-byte (SHA-256) audit of:
  _duplicates/ vs rest of onemind-codex vault

Goal: confirm every file in _duplicates/ has an identical copy
elsewhere in the vault before deleting _duplicates/.

Usage:
  python3 /Users/zeuslegacy/Documents/onemind-codex/_system/verify-duplicates.py

Output: /tmp/duplicates-audit/
  UNMATCHED.txt          — files in _duplicates/ with NO match elsewhere
  MATCHED-by-subfolder.txt — match rate per top-level subfolder
  report.json            — full machine-readable results
"""

import hashlib
import json
import os
import sys
import time
from collections import defaultdict
from datetime import datetime

# ─── CONFIG ───────────────────────────────────────────────────────────────────

VAULT_ROOT      = "/Users/zeuslegacy/Documents/onemind-codex"
DUPLICATES_DIR  = os.path.join(VAULT_ROOT, "_duplicates")
OUTPUT_DIR      = "/tmp/duplicates-audit"

# Dirs to skip when building the canonical index
SKIP_DIRS = {
    "_duplicates",
    ".git",
    ".obsidian",
    ".Spotlight-V100",
    ".fseventsd",
    ".TemporaryItems",
}

# Files to always skip (by name)
SKIP_FILENAMES = {
    ".DS_Store",
    ".localized",
}

# File prefixes to skip (AppleDouble resource forks)
SKIP_PREFIXES = ("._",)

CHUNK_SIZE      = 65536   # 64 KB reads
PROGRESS_EVERY  = 500     # print progress every N files
PROGRESS_EVERY_A = 1000   # canonical index prints less often (larger set)

# ─── HELPERS ──────────────────────────────────────────────────────────────────

def should_skip_file(name: str) -> bool:
    if name in SKIP_FILENAMES:
        return True
    for prefix in SKIP_PREFIXES:
        if name.startswith(prefix):
            return True
    return False


def sha256_file(path: str) -> str | None:
    """Return hex SHA-256 of file contents, or None on error."""
    h = hashlib.sha256()
    try:
        with open(path, "rb") as f:
            while chunk := f.read(CHUNK_SIZE):
                h.update(chunk)
        return h.hexdigest()
    except (OSError, PermissionError) as e:
        return None


def fmt_size(n_bytes: int) -> str:
    for unit in ("B", "KB", "MB", "GB"):
        if n_bytes < 1024:
            return f"{n_bytes:.1f} {unit}"
        n_bytes /= 1024
    return f"{n_bytes:.1f} TB"


def fmt_time(seconds: float) -> str:
    if seconds < 60:
        return f"{seconds:.0f}s"
    m, s = divmod(int(seconds), 60)
    if m < 60:
        return f"{m}m {s}s"
    h, m = divmod(m, 60)
    return f"{h}h {m}m {s}s"


def eta(done: int, total: int, elapsed: float) -> str:
    if done == 0:
        return "?"
    rate = done / elapsed
    remaining = (total - done) / rate
    return fmt_time(remaining)


# ─── PHASE A: build canonical hash index ──────────────────────────────────────

def build_canonical_index() -> tuple[set, int, list]:
    """
    Walk vault root, skip _duplicates/ and hidden dirs,
    SHA-256 every file → return (hash_set, file_count, errors).
    """
    print("\n" + "=" * 60)
    print("PHASE A — Building canonical index (rest of vault)")
    print("  Scanning:", VAULT_ROOT)
    print("  Skipping: _duplicates/, .git/, .obsidian/, and system dirs")
    print("=" * 60)

    canonical = set()
    errors = []
    count = 0
    t0 = time.time()

    for dirpath, dirnames, filenames in os.walk(VAULT_ROOT, topdown=True):
        # Prune skip dirs in-place (prevents os.walk from descending)
        rel = os.path.relpath(dirpath, VAULT_ROOT)
        top_component = rel.split(os.sep)[0]

        # Skip the duplicates dir entirely
        if top_component == "_duplicates":
            dirnames.clear()
            continue

        # Skip other ignored dirs
        dirnames[:] = [
            d for d in dirnames
            if d not in SKIP_DIRS and not d.startswith("._")
        ]

        for name in filenames:
            if should_skip_file(name):
                continue

            fpath = os.path.join(dirpath, name)

            # Skip symlinks
            if os.path.islink(fpath):
                continue

            digest = sha256_file(fpath)
            if digest is None:
                errors.append(fpath)
                continue

            canonical.add(digest)
            count += 1

            if count % PROGRESS_EARLY == 0:
                elapsed = time.time() - t0
                print(f"  [{count:>8,}] hashed  |  {len(canonical):,} unique  |  "
                      f"{fmt_time(elapsed)} elapsed", end="\r", flush=True)

    elapsed = time.time() - t0
    print(f"\n  DONE  |  {count:,} files hashed  |  {len(canonical):,} unique hashes  |  {fmt_time(elapsed)}")
    if errors:
        print(f"  ERRORS: {len(errors)} files unreadable (will be listed in report)")

    return canonical, count, errors


# ─── PHASE B: audit _duplicates/ ──────────────────────────────────────────────

def audit_duplicates(canonical: set) -> dict:
    """
    Walk _duplicates/, SHA-256 each file, check against canonical.
    Returns results dict keyed by top-level subfolder.
    """
    print("\n" + "=" * 60)
    print("PHASE B — Auditing _duplicates/")
    print("  Target:", DUPLICATES_DIR)
    print("=" * 60)

    # Pre-count for ETA
    print("  Pre-counting files...", end="", flush=True)
    total_files = sum(
        1
        for dirpath, _, filenames in os.walk(DUPLICATES_DIR)
        for name in filenames
        if not should_skip_file(name) and not os.path.islink(os.path.join(dirpath, name))
    )
    print(f" {total_files:,} files found")

    # Results structure: {subfolder: {matched: [], unmatched: [], errors: []}}
    results = defaultdict(lambda: {"matched": 0, "unmatched": [], "errors": [], "bytes_matched": 0, "bytes_unmatched": 0})

    count = 0
    unmatched_total = 0
    t0 = time.time()

    for dirpath, dirnames, filenames in os.walk(DUPLICATES_DIR, topdown=True):
        dirnames[:] = [d for d in dirnames if not d.startswith("._")]

        # Determine top-level subfolder
        rel = os.path.relpath(dirpath, DUPLICATES_DIR)
        top_sub = rel.split(os.sep)[0] if rel != "." else "(root)"

        for name in filenames:
            if should_skip_file(name):
                continue

            fpath = os.path.join(dirpath, name)

            if os.path.islink(fpath):
                continue

            try:
                fsize = os.path.getsize(fpath)
            except OSError:
                fsize = 0

            digest = sha256_file(fpath)

            if digest is None:
                results[top_sub]["errors"].append({
                    "path": fpath,
                    "size": fsize,
                    "rel_path": os.path.relpath(fpath, DUPLICATES_DIR),
                })
                count += 1
                continue

            if digest in canonical:
                results[top_sub]["matched"] += 1
                results[top_sub]["bytes_matched"] += fsize
            else:
                results[top_sub]["unmatched"].append({
                    "path": fpath,
                    "size": fsize,
                    "sha256": digest,
                    "rel_path": os.path.relpath(fpath, DUPLICATES_DIR),
                })
                results[top_sub]["bytes_unmatched"] += fsize
                unmatched_total += 1

            count += 1

            if count % PROGRESS_EVERY == 0:
                elapsed = time.time() - t0
                pct = (count / total_files * 100) if total_files else 0
                print(
                    f"  [{count:>8,}/{total_files:,}]  {pct:.1f}%  |  "
                    f"unmatched so far: {unmatched_total:,}  |  "
                    f"ETA: {eta(count, total_files, elapsed)}        ",
                    end="\r", flush=True
                )

    elapsed = time.time() - t0
    print(f"\n  DONE  |  {count:,} files audited  |  "
          f"unmatched: {unmatched_total:,}  |  {fmt_time(elapsed)}")

    return dict(results)


# ─── PHASE C: write reports ────────────────────────────────────────────────────

def write_reports(results: dict, canonical_count: int, canonical_errors: list):
    os.makedirs(OUTPUT_DIR, exist_ok=True)
    now = datetime.now().strftime("%Y-%m-%d %H:%M:%S")

    # Aggregate totals
    total_matched      = sum(r["matched"]              for r in results.values())
    total_unmatched    = sum(len(r["unmatched"])        for r in results.values())
    total_errors       = sum(len(r["errors"])           for r in results.values())
    total_files        = total_matched + total_unmatched + total_errors
    bytes_unmatched    = sum(r["bytes_unmatched"]       for r in results.values())
    bytes_matched      = sum(r["bytes_matched"]         for r in results.values())
    match_pct          = (total_matched / total_files * 100) if total_files else 0

    # ── UNMATCHED.txt ──────────────────────────────────────────────────────────
    unmatched_path = os.path.join(OUTPUT_DIR, "UNMATCHED.txt")
    with open(unmatched_path, "w", encoding="utf-8") as fh:
        fh.write(f"UNMATCHED FILES IN _duplicates/\n")
        fh.write(f"Generated: {now}\n")
        fh.write(f"Vault root: {VAULT_ROOT}\n")
        fh.write(f"{'=' * 70}\n\n")
        fh.write(f"SUMMARY\n")
        fh.write(f"  Total audited : {total_files:,}\n")
        fh.write(f"  Matched       : {total_matched:,}  ({match_pct:.1f}%)\n")
        fh.write(f"  UNMATCHED     : {total_unmatched:,}  ({fmt_size(bytes_unmatched)})\n")
        fh.write(f"  Read errors   : {total_errors:,}\n\n")

        if total_unmatched == 0:
            fh.write("✅ ALL FILES HAVE A MATCH ELSEWHERE IN THE VAULT.\n")
            fh.write("   Safe to delete _duplicates/ (within codex vault scope).\n")
        else:
            fh.write(f"⚠️  {total_unmatched:,} files ({fmt_size(bytes_unmatched)}) have NO match in vault (excl. _duplicates/).\n\n")
            fh.write("NOTE: Large groups like 'onemind-library/' and 'legacy-media-code/'\n")
            fh.write("      will likely appear here because their originals live OUTSIDE\n")
            fh.write("      this vault (in Documents/onemind-library or OM drives).\n\n")

            fh.write(f"{'=' * 70}\n")
            fh.write("UNMATCHED FILES BY SUBFOLDER\n")
            fh.write(f"{'=' * 70}\n\n")

            # Sort subfolders by unmatched count descending
            for sub, r in sorted(results.items(), key=lambda x: len(x[1]["unmatched"]), reverse=True):
                unmatched_list = r["unmatched"]
                if not unmatched_list:
                    continue
                sub_total = r["matched"] + len(unmatched_list)
                sub_pct = (len(unmatched_list) / sub_total * 100) if sub_total else 0
                fh.write(f"── {sub}  ({len(unmatched_list):,} unmatched / {sub_total:,} total | {fmt_size(r['bytes_unmatched'])})\n")
                # Sort by size descending so biggest unknowns are first
                for item in sorted(unmatched_list, key=lambda x: x["size"], reverse=True):
                    fh.write(f"   [{fmt_size(item['size']):>10}]  {item['rel_path']}\n")
                fh.write("\n")

        if total_errors > 0:
            fh.write(f"\n{'=' * 70}\n")
            fh.write("READ ERRORS (could not hash — neither matched nor unmatched)\n")
            fh.write(f"{'=' * 70}\n\n")
            for sub, r in results.items():
                for item in r["errors"]:
                    fh.write(f"  ERROR: {item['rel_path']}\n")

    # ── MATCHED-by-subfolder.txt ───────────────────────────────────────────────
    breakdown_path = os.path.join(OUTPUT_DIR, "MATCHED-by-subfolder.txt")
    with open(breakdown_path, "w", encoding="utf-8") as fh:
        fh.write(f"MATCH BREAKDOWN BY SUBFOLDER\n")
        fh.write(f"Generated: {now}\n")
        fh.write(f"{'=' * 70}\n\n")
        fh.write(f"{'SUBFOLDER':<45} {'MATCH%':>7}  {'MATCHED':>8}  {'UNMATCHED':>9}  {'TOTAL':>8}\n")
        fh.write(f"{'-' * 70}\n")

        for sub, r in sorted(results.items(), key=lambda x: x[0]):
            sub_matched   = r["matched"]
            sub_unmatched  = len(r["unmatched"])
            sub_errors    = len(r["errors"])
            sub_total     = sub_matched + sub_unmatched + sub_errors
            sub_pct       = (sub_matched / sub_total * 100) if sub_total else 0
            flag = "⚠️ " if sub_unmatched > 0 else "✅"
            fh.write(f"{flag} {sub:<43} {sub_pct:>7.1f}%  {sub_matched:>8,}  {sub_unmatched:>9,}  {sub_total:>8,}\n")

        fh.write(f"{'-' * 70}\n")
        fh.write(f"{'TOTAL':<45} {match_pct:>7.1f}%  {total_matched:>8,}  {total_unmatched:>9,}  {total_files:>8,}\n")

    # ── report.json ────────────────────────────────────────────────────────────
    report_path = os.path.join(OUTPUT_DIR, "report.json")
    report = {
        "generated": now,
        "vault_root": VAULT_ROOT,
        "duplicates_dir": DUPLICATES_DIR,
        "canonical_index": {
            "unique_hashes": canonical_count,
            "read_errors": len(canonical_errors),
        },
        "audit": {
            "total_files": total_files,
            "matched": total_matched,
            "unmatched": total_unmatched,
            "read_errors": total_errors,
            "match_pct": round(match_pct, 2),
            "bytes_matched": bytes_matched,
            "bytes_unmatched": bytes_unmatched,
        },
        "by_subfolder": {
            sub: {
                "matched": r["matched"],
                "unmatched": len(r["unmatched"]),
                "errors": len(r["errors"]),
                "bytes_matched": r["bytes_matched"],
                "bytes_unmatched": r["bytes_unmatched"],
                "unmatched_files": [
                    {"rel_path": x["rel_path"], "size": x["size"], "sha256": x["sha256"]}
                    for x in sorted(r["unmatched"], key=lambda x: x["size"], reverse=True)
                ],
            }
            for sub, r in results.items()
        },
    }
    with open(report_path, "w", encoding="utf-8") as fh:
        json.dump(report, fh, indent=2)

    return total_unmatched, total_matched, total_files, bytes_unmatched, match_pct


# ─── MAIN ─────────────────────────────────────────────────────────────────────

# Fix the typo in phase A progress printing
PROGRESS_EARLY = PROGRESS_EVERY_A

def main():
    print("=" * 60)
    print("  DUPLICATES AUDIT — SHA-256 Byte-for-Byte Verification")
    print("=" * 60)
    print(f"  Audit target : {DUPLICATES_DIR}")
    print(f"  Vault scope  : {VAULT_ROOT} (excl. _duplicates/)")
    print(f"  Output dir   : {OUTPUT_DIR}")
    print(f"  Started      : {datetime.now().strftime('%Y-%m-%d %H:%M:%S')}")
    print()
    print("  ⏱  This will take 15–40 minutes (95 GB of SHA-256 hashing).")
    print("     Safe to leave running. Progress prints every 500 files.")
    print()

    if not os.path.isdir(DUPLICATES_DIR):
        print(f"ERROR: _duplicates/ not found at {DUPLICATES_DIR}")
        sys.exit(1)

    t_total = time.time()

    # Phase A
    canonical, canonical_file_count, canonical_errors = build_canonical_index()

    # Phase B
    results = audit_duplicates(canonical)

    # Phase C
    total_unmatched, total_matched, total_files, bytes_unmatched, match_pct = write_reports(
        results, canonical_file_count, canonical_errors
    )

    elapsed_total = time.time() - t_total

    # ── Final console summary ──────────────────────────────────────────────────
    print("\n" + "=" * 60)
    print("  FINAL SUMMARY")
    print("=" * 60)
    print(f"  Files in _duplicates/ audited : {total_files:,}")
    print(f"  ✅ Matched (safe to delete)    : {total_matched:,}")
    print(f"  ⚠️  UNMATCHED (review required) : {total_unmatched:,}  ({fmt_size(bytes_unmatched)})")
    print(f"  Match rate                    : {match_pct:.1f}%")
    print(f"  Total time                    : {fmt_time(elapsed_total)}")
    print()
    print(f"  Reports written to: {OUTPUT_DIR}/")
    print(f"    UNMATCHED.txt           — files with no match (full list)")
    print(f"    MATCHED-by-subfolder.txt — breakdown by subfolder")
    print(f"    report.json             — machine-readable full results")
    print()

    if total_unmatched == 0:
        print("  ✅  ALL CLEAR — every file in _duplicates/ has a match")
        print("      elsewhere in the vault. Safe to delete.")
    else:
        print(f"  ⚠️   {total_unmatched:,} files have NO match within the codex vault.")
        print("      Review UNMATCHED.txt to decide if these are:")
        print("      a) Originals outside codex scope (onemind-library, OM drives)")
        print("      b) Truly unique files that must be moved out before deleting")

    print("=" * 60)


if __name__ == "__main__":
    main()
