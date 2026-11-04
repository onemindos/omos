"""Per-credential Codex usage for onemind-hermes (codex-1, codex-2, ...).

`hermes usage --json` reports only the credential that happens to be active. Hermes' own
fetcher, agent.account_usage.fetch_account_usage(), accepts an explicit api_key, and the
credential pool is readable with hermes_cli.auth.read_credential_pool(). So this script walks
the openai-codex pool and fetches each entry's usage with read_only=True, the variant the model
picker uses: a 401 is reported, never repaired, and nothing is refreshed, rotated or written.

Tokens are never printed. Usage: legacy_codex_usage.py [--json]   (default: one text line each)
"""

import json
import sys

sys.path.insert(0, "/opt/hermes")


def collect():
    from agent.account_usage import fetch_account_usage
    from hermes_cli.auth import read_credential_pool

    rows = []
    for entry in read_credential_pool("openai-codex"):
        row = {
            "label": entry.get("label") or entry.get("id") or "?",
            "priority": entry.get("priority"),
            "requests": entry.get("request_count"),
            "pool_status": entry.get("last_status"),
            "cooldown_until": entry.get("last_error_reset_at"),
        }
        token = entry.get("access_token") or ""
        try:
            snap = fetch_account_usage(
                "openai-codex", base_url=entry.get("base_url"), api_key=token, read_only=True)
        except Exception as exc:  # noqa: BLE001
            snap, row["error"] = None, type(exc).__name__
        if snap is None:
            row.setdefault("error", "usage unavailable (token expired or fetch failed)")
        else:
            row["plan"] = snap.plan
            row["windows"] = [
                {"label": w.label, "used_percent": w.used_percent,
                 "reset_at": w.reset_at.isoformat() if w.reset_at else None}
                for w in snap.windows
            ]
            row["details"] = list(snap.details)
        rows.append(row)
    return rows


def render(rows):
    lines = []
    for row in rows:
        if row.get("error"):
            lines.append("%s: %s" % (row["label"], row["error"]))
            continue
        parts = []
        for w in row["windows"]:
            used = "?" if w["used_percent"] is None else "%.0f" % w["used_percent"]
            reset = " (resets %s UTC)" % w["reset_at"][:16].replace("T", " ") if w["reset_at"] else ""
            parts.append("%s %s%%%s" % (w["label"], used, reset))
        wins = ", ".join(parts) or "no windows reported"
        extra = ""
        if row.get("cooldown_until"):
            extra = "; pool cooldown until %s" % str(row["cooldown_until"])[:16]
        lines.append("%s (%s): %s%s" % (row["label"], row.get("plan") or "plan ?", wins, extra))
    return lines or ["no openai-codex credentials in the pool"]


def main():
    try:
        rows = collect()
    except Exception as exc:  # noqa: BLE001
        print("codex usage unavailable: %s" % type(exc).__name__)
        return 0
    if "--json" in sys.argv[1:]:
        print(json.dumps(rows, indent=2))
    else:
        print("\n".join(render(rows)))
    return 0


if __name__ == "__main__":
    sys.exit(main())
