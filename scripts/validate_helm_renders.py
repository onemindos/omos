"""`helm template` render check for Helm charts that Argo CD or the repo deploys.

Covers:
  1. Argo CD Applications whose source is a remote Helm chart (HTTP repo or OCI). The chart is
     pulled at the pinned targetRevision and rendered with the Application's own values
     (valueFiles resolved from this checkout, inline `values`, `valuesObject`).
  2. Local charts (any Chart.yaml outside third-party/), rendered with their production values.

Default mode renders only what a change could affect, using the diff between YAML_BASE_SHA and
HEAD. `--all` renders everything. Exit code is non-zero if any render fails.
"""
from __future__ import annotations

import os
from pathlib import Path
import re
import subprocess
import sys
import tempfile

import yaml  # installed with yamllint (PyYAML)

ROOT = Path(__file__).resolve().parents[1]
APP_DIRS = (
    "platform/argocd",
    "platform/ecosystem/applications",
    "platform/ecosystem/foundation-applications",
)
# Local chart dir -> values file used to render it (relative to repo root).
LOCAL_VALUES = {
    "platform/takserver/chart": ["platform/takserver/values-onemind.yaml"],
    "platform/ecosystem/charts/onemind-core": [],
}
KUBE_VERSION = os.environ.get("HELM_KUBE_VERSION", "1.36.0")


def changed_files() -> set[str]:
    base = os.environ.get("YAML_BASE_SHA", "")
    if not base or set(base) == {"0"}:
        base = subprocess.check_output(["git", "rev-parse", "HEAD^"], text=True, cwd=ROOT).strip()
    if not re.fullmatch(r"[0-9a-f]{40}", base):
        raise SystemExit("Expected a commit SHA for Helm render validation")
    out = subprocess.check_output(
        ["git", "diff", "--name-only", "--diff-filter=ACMR", "-z", base, "HEAD"], cwd=ROOT)
    return {p for p in out.decode().split("\0") if p}


def sources_of(app: dict) -> list[dict]:
    spec = app.get("spec", {})
    return spec.get("sources") or ([spec["source"]] if spec.get("source") else [])


def load_apps() -> list[tuple[str, dict]]:
    apps = []
    for d in APP_DIRS:
        for path in sorted((ROOT / d).glob("*.y*ml")):
            for doc in yaml.safe_load_all(path.read_text()):
                if isinstance(doc, dict) and doc.get("kind") == "Application":
                    apps.append((str(path.relative_to(ROOT)), doc))
    return apps


def run_helm(args: list[str], label: str) -> bool:
    proc = subprocess.run(["helm", "template", *args, "--include-crds", "--kube-version", KUBE_VERSION],
                          capture_output=True, text=True)
    if proc.returncode != 0:
        print(f"FAIL  {label}\n{proc.stderr.strip()[-1500:]}")
        return False
    print(f"ok    {label}")
    return True


def render_app(rel_path: str, app: dict, only_if_changed: set[str] | None) -> bool | None:
    name = app["metadata"]["name"]
    srcs = sources_of(app)
    ok: bool | None = None
    for src in srcs:
        chart = src.get("chart")
        if not chart:
            continue
        helm = src.get("helm", {}) or {}
        value_files = [v.replace("$values/", "") for v in helm.get("valueFiles", [])]
        if only_if_changed is not None and not ({rel_path, *value_files} & only_if_changed):
            continue
        namespace = app.get("spec", {}).get("destination", {}).get("namespace") or name
        release = helm.get("releaseName") or name
        version = str(src.get("targetRevision"))
        repo = src["repoURL"]
        with tempfile.TemporaryDirectory(prefix="helm-render-") as tmp:
            args = [release]
            if repo.startswith(("http://", "https://")):
                args += [chart, "--repo", repo]
            else:
                args += [f"oci://{repo.rstrip('/')}/{chart}"]
            args += ["--version", version, "--namespace", namespace]
            for vf in value_files:
                path = ROOT / vf
                if not path.is_file():
                    print(f"FAIL  {name}: values file {vf} does not exist")
                    return False
                args += ["-f", str(path)]
            if helm.get("values"):
                p = Path(tmp) / "inline-values.yaml"
                p.write_text(helm["values"])
                args += ["-f", str(p)]
            if helm.get("valuesObject"):
                p = Path(tmp) / "values-object.yaml"
                p.write_text(yaml.safe_dump(helm["valuesObject"]))
                args += ["-f", str(p)]
            result = run_helm(args, f"{name} ({chart} {version})")
        ok = result if ok is None else (ok and result)
    return ok


def render_local(chart_dir: str, values: list[str], only_if_changed: set[str] | None) -> bool | None:
    if only_if_changed is not None:
        prefix = chart_dir.rstrip("/") + "/"
        if not any(f.startswith(prefix) or f in values for f in only_if_changed):
            return None
    if not (ROOT / chart_dir / "Chart.yaml").is_file():
        return None
    args = [Path(chart_dir).name, str(ROOT / chart_dir)]
    for v in values:
        args += ["-f", str(ROOT / v)]
    return run_helm(args, f"local chart {chart_dir}")


def main() -> int:
    only = None if "--all" in sys.argv else changed_files()
    results: list[bool] = []
    for rel_path, app in load_apps():
        r = render_app(rel_path, app, only)
        if r is not None:
            results.append(r)
    for chart_dir, values in LOCAL_VALUES.items():
        r = render_local(chart_dir, values, only)
        if r is not None:
            results.append(r)
    # Guard: any Chart.yaml not listed above should be added to LOCAL_VALUES.
    for chart in ROOT.glob("**/Chart.yaml"):
        rel = str(chart.parent.relative_to(ROOT))
        if rel not in LOCAL_VALUES and "third-party" not in rel:
            print(f"warn  local chart {rel} is not registered in LOCAL_VALUES (not rendered)")
    if not results:
        print("No Helm-affecting changes to render")
        return 0
    print(f"Rendered {len(results)} target(s); {results.count(False)} failed")
    return 0 if all(results) else 1


if __name__ == "__main__":
    sys.exit(main())
