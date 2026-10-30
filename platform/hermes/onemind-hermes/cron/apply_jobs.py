"""Idempotent applier for cron/jobs.yaml. Runs inside the onemind-hermes pod, as the `hermes` user.

    python3 apply_jobs.py [--dry-run] [--jobs /opt/data/cron-defs/jobs.yaml]

Uses Hermes' own cron API (cron.jobs: create_job, update_job, pause_job, resume_job), the same
code `hermes cron create/edit/pause/resume` call, so every value goes through Hermes' validation.

- A job missing by name is created (paused from the start when `paused: true`).
- An existing job is updated only for fields that differ; nothing is rewritten when it matches,
  so run history and the next-run time survive a re-apply.
- The enabled/paused state is made to match `paused` in the file.
- Jobs not in the file are listed as unmanaged and never touched.
A second run with no changes prints only "unchanged" lines.
"""

import os
import sys

sys.path.insert(0, "/opt/hermes")
from ruamel.yaml import YAML  # noqa: E402  (PyYAML is not in Hermes' venv)
from cron import jobs as cj  # noqa: E402

DRY = "--dry-run" in sys.argv
PATH = sys.argv[sys.argv.index("--jobs") + 1] if "--jobs" in sys.argv else "/opt/data/cron-defs/jobs.yaml"
SCRIPTS = os.path.join(os.environ.get("HERMES_HOME", "/opt/data"), "scripts")

# desired key -> job record key
PLAIN = {"prompt": "prompt", "deliver": "deliver", "script": "script", "monitor_script": "monitor_script",
         "provider": "provider", "model": "model", "reasoning_effort": "reasoning_effort"}


def norm(value):
    return (value.strip() or None) if isinstance(value, str) else (value or None)


def drift(spec, job):
    ch = {}
    for key, rec in PLAIN.items():
        if key in spec and norm(spec[key]) != norm(job.get(rec)):
            ch[rec] = norm(spec[key])
    if sorted(spec.get("skills") or []) != sorted(job.get("skills") or []):
        ch["skills"] = list(spec.get("skills") or [])
    if bool(spec.get("no_agent", False)) != bool(job.get("no_agent")):
        ch["no_agent"] = bool(spec.get("no_agent", False))
    if cj.parse_schedule(spec["schedule"]) != job.get("schedule"):
        ch["schedule"] = spec["schedule"]
    return ch


def main():
    with open(PATH, encoding="utf-8") as fh:
        specs = (YAML(typ="safe").load(fh) or {}).get("jobs", [])
    existing = {j["name"]: j for j in cj.list_jobs(include_disabled=True)}
    rc = 0
    for spec in specs:
        name = spec["name"]
        for key in ("script", "monitor_script"):
            if spec.get(key) and not os.path.isfile(os.path.join(SCRIPTS, spec[key])):
                print("ERROR %s: %s=%s not found in %s" % (name, key, spec[key], SCRIPTS))
                rc = 1
        job = existing.get(name)
        paused = bool(spec.get("paused", False))
        reason = (spec.get("paused_reason") or "").strip() or None
        if job is None:
            print("%s %s: create%s" % ("WOULD" if DRY else "DO", name, " (paused)" if paused else ""))
            if not DRY and rc == 0:
                cj.create_job(
                    spec.get("prompt"), spec["schedule"], name=name, deliver=spec.get("deliver"),
                    skills=spec.get("skills") or None, model=spec.get("model"), provider=spec.get("provider"),
                    script=spec.get("script"), no_agent=bool(spec.get("no_agent", False)),
                    monitor_script=spec.get("monitor_script"), reasoning_effort=spec.get("reasoning_effort"),
                    paused=paused, paused_reason=reason if paused else None)
            continue
        changes = drift(spec, job)
        if changes:
            print("%s %s: update %s" % ("WOULD" if DRY else "DO", name, ", ".join(sorted(changes))))
            if not DRY:
                cj.update_job(job["id"], changes)
        enabled = bool(job.get("enabled", True)) and job.get("state") != "paused"
        want_pause = paused and (enabled or (reason and job.get("paused_reason") != reason))
        want_resume = not paused and job.get("state") == "paused"
        if want_pause:
            print("%s %s: pause" % ("WOULD" if DRY else "DO", name))
            if not DRY:
                cj.pause_job(job["id"], reason)
        elif want_resume:
            print("%s %s: resume" % ("WOULD" if DRY else "DO", name))
            if not DRY:
                cj.resume_job(job["id"])
        if not (changes or want_pause or want_resume):
            print("unchanged %s" % name)
    managed = {s["name"] for s in specs}
    for name in sorted(set(existing) - managed):
        print("unmanaged (left alone) %s" % name)
    return rc


if __name__ == "__main__":
    sys.exit(main())
