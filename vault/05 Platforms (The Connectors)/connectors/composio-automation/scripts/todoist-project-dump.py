#!/usr/bin/env python3
"""
todoist-project-dump.py
Dumps all tasks from a Todoist project (grouped by section) to stdout.
Useful for quick audits, sending to an agent, or copying into a markdown file.

Usage:
  python3 todoist-project-dump.py --project-id 6g6P3vvjx9643JrW
  python3 todoist-project-dump.py --project-id PROJECT_ID --format json
  python3 todoist-project-dump.py --project-id PROJECT_ID --format markdown

Requires: composio CLI authenticated
"""

import argparse
import json
import subprocess
import sys
from collections import defaultdict


def composio(slug, params):
    result = subprocess.run(
        ["composio", "execute", slug, "-d", json.dumps(params)],
        capture_output=True, text=True, timeout=30
    )
    if result.returncode != 0:
        raise RuntimeError(f"composio {slug} failed:\n{result.stderr}")
    return json.loads(result.stdout)


def get_tasks_for_project(project_id):
    raw = composio("TODOIST_GET_PROJECT_FULL", {"project_id": project_id})
    data = raw.get("data", {})
    project = data.get("project", {})
    sections = data.get("sections", [])
    tasks = data.get("tasks", [])
    return project, sections, tasks


def main():
    parser = argparse.ArgumentParser(description="Dump all tasks from a Todoist project")
    parser.add_argument("--project-id", required=True, help="Todoist project ID")
    parser.add_argument("--format", choices=["text", "json", "markdown"], default="text")
    args = parser.parse_args()

    print(f"📥 Fetching project {args.project_id}...", file=sys.stderr)

    try:
        project, sections, tasks = get_tasks_for_project(args.project_id)
    except Exception as e:
        # Fallback: use filter if project_full fails
        print(f"⚠️  get_project_full failed ({e}), falling back to task list...", file=sys.stderr)
        raw = composio("TODOIST_GET_ALL_TASKS", {})
        data = raw.get("data", {})
        tasks = data.get("tasks", data.get("items", []))
        tasks = [t for t in tasks if t.get("project_id") == args.project_id]
        sections = []
        project = {"name": args.project_id}

    if args.format == "json":
        print(json.dumps({"project": project, "sections": sections, "tasks": tasks}, indent=2))
        return

    # Group tasks by section
    section_map = {s["id"]: s["name"] for s in sections}
    by_section = defaultdict(list)
    for t in tasks:
        sid = t.get("section_id") or "none"
        by_section[sid].append(t)

    project_name = project.get("name", args.project_id)

    PRIORITY = {4: "[P1]", 3: "[P2]", 2: "[P3]", 1: "[P4]"}

    if args.format == "markdown":
        print(f"# {project_name}\n")
        print(f"> Project ID: `{args.project_id}`  |  Total tasks: {len(tasks)}\n")
        # Sections first
        for sec in sections:
            sec_tasks = by_section.get(sec["id"], [])
            print(f"## {sec['name']}  ({len(sec_tasks)})\n")
            for t in sorted(sec_tasks, key=lambda x: -x.get("priority", 1)):
                p = PRIORITY.get(t.get("priority", 1), "")
                due = t.get("due", {})
                due_str = f" · due {due['date']}" if due else ""
                print(f"- {p} {t['content']}{due_str}")
            print()
        # No-section tasks
        no_sec = by_section.get("none", [])
        if no_sec:
            print(f"## (No Section)  ({len(no_sec)})\n")
            for t in sorted(no_sec, key=lambda x: -x.get("priority", 1)):
                p = PRIORITY.get(t.get("priority", 1), "")
                print(f"- {p} {t['content']}")
    else:
        # Plain text
        print(f"\n{'='*60}")
        print(f"  PROJECT: {project_name}")
        print(f"  ID: {args.project_id}  |  Tasks: {len(tasks)}")
        print(f"{'='*60}")
        for sec in sections:
            sec_tasks = by_section.get(sec["id"], [])
            print(f"\n📂 {sec['name']}  ({len(sec_tasks)} tasks)")
            for t in sorted(sec_tasks, key=lambda x: -x.get("priority", 1)):
                p = PRIORITY.get(t.get("priority", 1), "")
                print(f"   {p} {t['content'][:70]}")
        no_sec = by_section.get("none", [])
        if no_sec:
            print(f"\n⚠️  (No Section)  ({len(no_sec)} tasks)")
            for t in no_sec:
                print(f"   • {t['content'][:70]}")
        print(f"\n{'='*60}")


if __name__ == "__main__":
    main()
