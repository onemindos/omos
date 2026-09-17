#!/usr/bin/env python3
"""
verify-duplicates.py — SHA256-based duplicate verification for _duplicates/ folder

Purpose: Verify that every file in _duplicates/ has a byte-identical copy
somewhere else in the codex vault. Produces a report showing SAFE matches
and ORPHAN files (no match found outside _duplicates/).

Usage:
    python3 _system/development/verify-duplicates.py

Output:
    /tmp/dup_verify/originals_sha256.tsv   — hash index of vault files (excl _duplicates)
    /tmp/dup_verify/duplicates_sha256.tsv   — hash index of _duplicates files
    /tmp/dup_verify/verification_report.txt — final report with orphan list
    /tmp/dup_verify/orphans.tsv             — just the orphan files (if any)
    /tmp/dup_verify/errors.log              — files that couldn't be read

No files in the vault are modified. All output goes to /tmp/dup_verify/.
"""

import hashlib
import os
import sys
import time
from datetime import datetime, timedelta
from pathlib import Path

# ─── Configuration ───────────────────────────────────────────────────────────

VAULT_ROOT = Path("/Users/zeuslegacy/Documents/onemind-codex")
DUPLICATES_DIR = VAULT_ROOT / "_duplicates"
OUTPUT_DIR = Path("/tmp/dup_verify")

# Directories to SKIP when building the originals index
SKIP_DIRS = {
    "_duplicates",
    ".git",
    ".obsidian",
    ".trash",
    "node_modules",
    "__pycache__",
    ".venv",
    "venv",
}

# Files to skip everywhere
SKIP_FILES = {".DS_Store"}
SKIP_PREFIXES = ("._",)

CHUNK_SIZE = 65536  # 64 KB read chunks for hashing
PROGRESS_INTERVAL = 1000  # Print progress every N files


# ─── Helpers ─────────────────────────────────────────────────────────────────

def sha256_file(filepath: Path) -> str | None:
    """Compute SHA256 hash of a file. Returns None on read error."""
    h = hashlib.sha256()
    try:
        with open(filepath, "rb") as f:
            while True:
                chunk = f.read(CHUNK_SIZE)
                if not chunk:
                    break
                h.update(chunk)
        return h.hexdigest()
    except (PermissionError, OSError):
        return None


def should_skip_file(name: str) -> bool:
    """Check if a filename should be skipped."""
    if name in SKIP_FILES:
        return True
    for prefix in SKIP_PREFIXES:
        if name.startswith(prefix):
            return True
    return False


def format_size(size_bytes: int) -> str:
    """Format bytes into human-readable string."""
    if size_bytes < 1024:
        return f"{size_bytes} B"
    elif size_bytes < 1024 ** 2:
        return f"{size_bytes / 1024:.1f} KB"
    elif size_bytes < 1024 ** 3:
        return f"{size_bytes / 1024**2:.1f} MB"
    else:
        return f"{size_bytes / 1024**3:.2f} GB"


def format_duration(seconds: float) -> str:
    """Format seconds into human-readable duration."""
    td = timedelta(seconds=int(seconds))
    parts = []
    hours, remainder = divmod(td.seconds, 3600)
    minutes, secs = divmod(remainder, 60)
    if td.days:
        parts.append(f"{td.days}d")
    if hours:
        parts.append(f"{hours}h")
    if minutes:
        parts.append(f"{minutes}m")
    parts.append(f"{secs}s")
    return " ".join(parts)


# ─── Index Building ──────────────────────────────────────────────────────────

def build_index(root_dir: Path, index_path: Path, skip_dirs: set[str],
                label: str, errors_log) -> dict[str, list[str]]:
    """
    Walk root_dir, SHA256 every file, save to TSV, return hash→[paths] dict.

    Supports resuming: if index_path already exists, loads it and skips
    already-hashed files.
    """
    # Load existing index for resume
    existing: dict[str, list[str]] = {}
    existing_paths: set[str] = set()
    if index_path.exists():
        with open(index_path, "r") as f:
            for line in f:
                line = line.rstrip("\n")
                if not line or line.startswith("#"):
                    continue
                parts = line.split("\t", 2)
                if len(parts) == 3:
                    h, sz, p = parts
                    existing.setdefault(h, []).append(p)
                    existing_paths.add(p)
        print(f"  Loaded {len(existing_paths)} previously indexed files from {index_path.name}")

    # Count total files first for progress estimate
    print(f"  Counting files in {label}...")
    total_files = 0
    for dirpath, dirnames, filenames in os.walk(root_dir):
        # Prune skipped directories
        rel = os.path.relpath(dirpath, root_dir)
        top_component = rel.split(os.sep)[0] if rel != "." else ""
        if top_component in skip_dirs:
            dirnames.clear()
            continue
        dirnames[:] = [d for d in dirnames if d not in skip_dirs]
        for fname in filenames:
            if not should_skip_file(fname):
                total_files += 1
    print(f"  Found {total_files:,} files to process")

    # Open TSV in append mode for new entries
    tsv_file = open(index_path, "a")
    if not existing_paths:
        tsv_file.write(f"# {label} SHA256 index — generated {datetime.now().isoformat()}\n")
        tsv_file.write("# sha256\tsize_bytes\trelative_path\n")

    hash_index = dict(existing)  # Start with what we loaded
    processed = len(existing_paths)
    new_hashed = 0
    total_bytes = 0
    errors = 0
    start_time = time.time()
    last_folder = ""

    for dirpath, dirnames, filenames in os.walk(root_dir):
        rel_dir = os.path.relpath(dirpath, root_dir)
        top_component = rel_dir.split(os.sep)[0] if rel_dir != "." else ""

        if top_component in skip_dirs:
            dirnames.clear()
            continue
        dirnames[:] = [d for d in dirnames if d not in skip_dirs]

        # Show folder transitions for top-level dirs
        current_top = rel_dir.split(os.sep)[0] if rel_dir != "." else "(root)"
        if current_top != last_folder:
            elapsed = time.time() - start_time
            rate = processed / elapsed if elapsed > 0 else 0
            remaining = (total_files - processed) / rate if rate > 0 else 0
            print(f"  [{processed:,}/{total_files:,}] Entering: {current_top}/ "
                  f"({format_duration(elapsed)} elapsed, ~{format_duration(remaining)} remaining)")
            last_folder = current_top

        for fname in filenames:
            if should_skip_file(fname):
                continue

            filepath = Path(dirpath) / fname
            rel_path = str(filepath.relative_to(root_dir))

            # Skip if already indexed
            if rel_path in existing_paths:
                processed += 1
                continue

            # Get file size
            try:
                file_size = filepath.stat().st_size
            except (PermissionError, OSError) as e:
                errors_log.write(f"STAT_ERROR\t{rel_path}\t{e}\n")
                errors += 1
                processed += 1
                continue

            # Hash the file
            file_hash = sha256_file(filepath)
            if file_hash is None:
                errors_log.write(f"READ_ERROR\t{rel_path}\n")
                errors += 1
                processed += 1
                continue

            # Record
            hash_index.setdefault(file_hash, []).append(rel_path)
            tsv_file.write(f"{file_hash}\t{file_size}\t{rel_path}\n")
            total_bytes += file_size
            new_hashed += 1
            processed += 1

            # Progress
            if new_hashed % PROGRESS_INTERVAL == 0:
                elapsed = time.time() - start_time
                rate = processed / elapsed if elapsed > 0 else 0
                remaining = (total_files - processed) / rate if rate > 0 else 0
                pct = (processed / total_files * 100) if total_files > 0 else 0
                print(f"  [{processed:,}/{total_files:,}] {pct:.1f}% | "
                      f"{format_size(total_bytes)} hashed | "
                      f"~{format_duration(remaining)} remaining")

    tsv_file.close()
    elapsed = time.time() - start_time

    print(f"  ✅ {label} complete: {processed:,} files indexed "
          f"({new_hashed:,} new, {len(existing_paths):,} resumed) in {format_duration(elapsed)}")
    if errors:
        print(f"  ⚠️  {errors} files had read errors (see errors.log)")

    return hash_index


# ─── Cross-Reference ─────────────────────────────────────────────────────────

def cross_reference(originals: dict[str, list[str]],
                    duplicates: dict[str, list[str]],
                    dup_index_path: Path) -> tuple[list, list, dict]:
    """
    Compare every duplicate file against the originals index.

    Returns: (safe_files, orphan_files, subfolder_stats)
    """
    safe_files = []    # (dup_path, size, hash, original_path)
    orphan_files = []  # (dup_path, size, hash)
    subfolder_stats = {}  # subfolder → {safe: N, orphan: N, safe_bytes: N, orphan_bytes: N}

    # Build a quick lookup: hash → exists in originals
    original_hashes = set(originals.keys())

    # We need size info — reload from TSV
    dup_file_info = {}  # path → (hash, size)
    with open(dup_index_path, "r") as f:
        for line in f:
            line = line.rstrip("\n")
            if not line or line.startswith("#"):
                continue
            parts = line.split("\t", 2)
            if len(parts) == 3:
                h, sz, p = parts
                dup_file_info[p] = (h, int(sz))

    total_dup_files = len(dup_file_info)
    checked = 0

    for dup_path, (file_hash, file_size) in dup_file_info.items():
        # Determine subfolder (first path component)
        parts = dup_path.split(os.sep)
        subfolder = parts[0] if len(parts) > 1 else "(root files)"

        if subfolder not in subfolder_stats:
            subfolder_stats[subfolder] = {
                "safe": 0, "orphan": 0,
                "safe_bytes": 0, "orphan_bytes": 0
            }

        if file_hash in original_hashes:
            # SAFE — found a match
            original_path = originals[file_hash][0]  # First matching original
            safe_files.append((dup_path, file_size, file_hash, original_path))
            subfolder_stats[subfolder]["safe"] += 1
            subfolder_stats[subfolder]["safe_bytes"] += file_size
        else:
            # ORPHAN — no match outside _duplicates
            orphan_files.append((dup_path, file_size, file_hash))
            subfolder_stats[subfolder]["orphan"] += 1
            subfolder_stats[subfolder]["orphan_bytes"] += file_size

        checked += 1
        if checked % 10000 == 0:
            pct = checked / total_dup_files * 100
            print(f"  Cross-referenced {checked:,}/{total_dup_files:,} ({pct:.1f}%)")

    return safe_files, orphan_files, subfolder_stats


# ─── Report Generation ───────────────────────────────────────────────────────

def generate_report(safe_files, orphan_files, subfolder_stats,
                    originals_count, duplicates_count, elapsed_total):
    """Generate the verification report and orphan list."""

    report_path = OUTPUT_DIR / "verification_report.txt"
    orphan_path = OUTPUT_DIR / "orphans.tsv"

    total_safe = len(safe_files)
    total_orphan = len(orphan_files)
    safe_bytes = sum(f[1] for f in safe_files)
    orphan_bytes = sum(f[1] for f in orphan_files)

    # Sort subfolders by total files descending
    sorted_subs = sorted(subfolder_stats.items(),
                         key=lambda x: x[1]["safe"] + x[1]["orphan"],
                         reverse=True)

    with open(report_path, "w") as r:
        r.write("=" * 72 + "\n")
        r.write("   DUPLICATES VERIFICATION REPORT\n")
        r.write("   SHA256 Byte-Identical Match Verification\n")
        r.write("=" * 72 + "\n\n")
        r.write(f"Date:              {datetime.now().strftime('%Y-%m-%d %H:%M:%S')}\n")
        r.write(f"Vault:             {VAULT_ROOT}\n")
        r.write(f"Duplicates folder: {DUPLICATES_DIR}\n")
        r.write(f"Total runtime:     {format_duration(elapsed_total)}\n\n")

        r.write(f"Originals indexed: {originals_count:,} unique hashes\n")
        r.write(f"Duplicates checked: {total_safe + total_orphan:,} files\n\n")

        r.write("-" * 72 + "\n")
        r.write("RESULTS\n")
        r.write("-" * 72 + "\n\n")

        safe_pct = (total_safe / (total_safe + total_orphan) * 100) if (total_safe + total_orphan) > 0 else 0
        r.write(f"  ✅ SAFE (byte-identical match found):  {total_safe:>8,} files  ({format_size(safe_bytes):>10})  [{safe_pct:.2f}%]\n")
        r.write(f"  ❌ ORPHAN (NO match found outside):    {total_orphan:>8,} files  ({format_size(orphan_bytes):>10})  [{100-safe_pct:.2f}%]\n\n")

        if total_orphan == 0:
            r.write("  🎉 VERDICT: ALL files in _duplicates/ have a byte-identical copy\n")
            r.write("     elsewhere in the vault. _duplicates/ is 100% SAFE TO DELETE.\n\n")
        else:
            r.write(f"  ⚠️  VERDICT: {total_orphan:,} files have NO match outside _duplicates/.\n")
            r.write("     DO NOT delete _duplicates/ until these orphans are resolved.\n")
            r.write("     See orphans.tsv for the full list.\n\n")

        r.write("-" * 72 + "\n")
        r.write("BREAKDOWN BY SUBFOLDER\n")
        r.write("-" * 72 + "\n\n")
        r.write(f"  {'Subfolder':<45} {'Safe':>7} {'Orphan':>7} {'Total':>7} {'Status'}\n")
        r.write(f"  {'-'*45} {'-'*7} {'-'*7} {'-'*7} {'-'*8}\n")

        for sub, stats in sorted_subs:
            total = stats["safe"] + stats["orphan"]
            status = "✅ CLEAN" if stats["orphan"] == 0 else f"⚠️  {stats['orphan']} orphans"
            display_name = sub[:44]
            r.write(f"  {display_name:<45} {stats['safe']:>7,} {stats['orphan']:>7,} {total:>7,} {status}\n")

        r.write(f"\n  {'TOTAL':<45} {total_safe:>7,} {total_orphan:>7,} {total_safe+total_orphan:>7,}\n")

        if total_orphan > 0:
            r.write("\n" + "-" * 72 + "\n")
            r.write("ORPHAN FILES — No byte-identical match found outside _duplicates/\n")
            r.write("-" * 72 + "\n\n")
            # Sort orphans by subfolder then path
            orphan_files_sorted = sorted(orphan_files, key=lambda x: x[0])
            for dup_path, file_size, file_hash in orphan_files_sorted:
                r.write(f"  {format_size(file_size):>10}  {file_hash[:16]}...  {dup_path}\n")

        r.write("\n" + "=" * 72 + "\n")
        r.write("END OF REPORT\n")
        r.write("=" * 72 + "\n")

    # Write orphans TSV
    with open(orphan_path, "w") as o:
        o.write("# Orphan files — no byte-identical match outside _duplicates/\n")
        o.write("# sha256\tsize_bytes\trelative_path\n")
        for dup_path, file_size, file_hash in sorted(orphan_files, key=lambda x: x[0]):
            o.write(f"{file_hash}\t{file_size}\t{dup_path}\n")

    return report_path


# ─── Main ────────────────────────────────────────────────────────────────────

def main():
    print("=" * 72)
    print("  DUPLICATES VERIFICATION — SHA256 Byte-Identical Check")
    print("=" * 72)
    print()
    print(f"  Vault root:       {VAULT_ROOT}")
    print(f"  Duplicates dir:   {DUPLICATES_DIR}")
    print(f"  Output dir:       {OUTPUT_DIR}")
    print(f"  Hash algorithm:   SHA256")
    print(f"  Chunk size:       {CHUNK_SIZE // 1024} KB")
    print()

    # Validate paths
    if not VAULT_ROOT.is_dir():
        print(f"ERROR: Vault root not found: {VAULT_ROOT}")
        sys.exit(1)
    if not DUPLICATES_DIR.is_dir():
        print(f"ERROR: Duplicates dir not found: {DUPLICATES_DIR}")
        sys.exit(1)

    # Create output directory
    OUTPUT_DIR.mkdir(parents=True, exist_ok=True)

    originals_tsv = OUTPUT_DIR / "originals_sha256.tsv"
    duplicates_tsv = OUTPUT_DIR / "duplicates_sha256.tsv"
    errors_path = OUTPUT_DIR / "errors.log"

    total_start = time.time()

    # Open shared error log
    errors_log = open(errors_path, "a")
    errors_log.write(f"\n# Run started {datetime.now().isoformat()}\n")

    # ── Step 1: Index originals (vault minus _duplicates) ──
    print("━" * 72)
    print("STEP 1/3: Indexing ORIGINALS (vault excluding _duplicates)")
    print("━" * 72)
    originals = build_index(
        root_dir=VAULT_ROOT,
        index_path=originals_tsv,
        skip_dirs=SKIP_DIRS,
        label="Originals",
        errors_log=errors_log
    )
    print()

    # ── Step 2: Index duplicates ──
    print("━" * 72)
    print("STEP 2/3: Indexing DUPLICATES (_duplicates/ folder)")
    print("━" * 72)
    duplicates = build_index(
        root_dir=DUPLICATES_DIR,
        index_path=duplicates_tsv,
        skip_dirs=set(),  # Don't skip anything inside _duplicates
        label="Duplicates",
        errors_log=errors_log
    )
    print()

    errors_log.close()

    # ── Step 3: Cross-reference ──
    print("━" * 72)
    print("STEP 3/3: Cross-referencing duplicates against originals")
    print("━" * 72)
    safe_files, orphan_files, subfolder_stats = cross_reference(
        originals, duplicates, duplicates_tsv
    )
    print()

    # ── Generate report ──
    total_elapsed = time.time() - total_start
    report_path = generate_report(
        safe_files, orphan_files, subfolder_stats,
        originals_count=len(originals),
        duplicates_count=len(duplicates),
        elapsed_total=total_elapsed
    )

    # ── Summary ──
    total_safe = len(safe_files)
    total_orphan = len(orphan_files)
    safe_pct = (total_safe / (total_safe + total_orphan) * 100) if (total_safe + total_orphan) > 0 else 0

    print("=" * 72)
    print("  VERIFICATION COMPLETE")
    print("=" * 72)
    print()
    print(f"  ✅ SAFE:   {total_safe:>8,} files ({safe_pct:.2f}%)")
    print(f"  ❌ ORPHAN: {total_orphan:>8,} files ({100-safe_pct:.2f}%)")
    print(f"  ⏱  Time:   {format_duration(total_elapsed)}")
    print()

    if total_orphan == 0:
        print("  🎉 ALL files in _duplicates/ are verified duplicates.")
        print("     The folder is SAFE TO DELETE.")
    else:
        print(f"  ⚠️  {total_orphan:,} orphan files found — DO NOT DELETE yet.")
        print(f"     Review: {OUTPUT_DIR / 'orphans.tsv'}")

    print()
    print(f"  Full report: {report_path}")
    print(f"  Orphan list: {OUTPUT_DIR / 'orphans.tsv'}")
    print(f"  Error log:   {OUTPUT_DIR / 'errors.log'}")
    print()


if __name__ == "__main__":
    main()
