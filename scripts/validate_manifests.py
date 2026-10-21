"""Schema-validate plain Kubernetes manifests with kubeconform.

Default mode checks only files changed between YAML_BASE_SHA and HEAD (same contract as
scripts/validate_changed_yaml.py). `--all` checks every manifest under the manifest roots.

What counts as a "plain manifest":
  * lives under MANIFEST_ROOTS,
  * is .yaml/.yml/.json,
  * has a top-level `apiVersion:` and `kind:`,
  * contains no Helm/Go template markers (`{{`), and
  * is not a kustomization, Helm values file, Chart.yaml or Argo CD ApplicationSet generator input.

Schemas: upstream Kubernetes (yannh) plus the datreeio CRDs-catalog for CRDs. Kinds that have
no published schema (for example several custom CRDs) are skipped via -ignore-missing-schemas,
so this gate catches typos and wrong field types in core resources and well-known CRDs
(Argo CD, cert-manager, SealedSecret, Tailscale, ...) rather than blocking on unknown CRDs.
"""
import os
from pathlib import Path
import re
import shutil
import subprocess
import sys

MANIFEST_ROOTS = (
    "platform/",
    "platform/takserver/deploy/",
    "platform/takserver/fedhub/",
)
# Paths under the roots that are not standalone manifests.
EXCLUDE = (
    re.compile(r"(^|/)values[^/]*\.(ya?ml)$"),
    re.compile(r"(^|/)(kustomization|Chart|chart-lock|foundation-lock)\.(ya?ml|json)$"),
    re.compile(r"^platform/values/"),
    re.compile(r"^platform/ecosystem/(values|charts|openchoreo-overrides)/"),
    re.compile(r"^platform/matrix/"),  # Synapse/ESS config and scripts, not cluster manifests
)
CRD_SCHEMA = (
    "https://raw.githubusercontent.com/datreeio/CRDs-catalog/main/"
    "{{.Group}}/{{.ResourceKind}}_{{.ResourceAPIVersion}}.json"
)
K8S_VERSION = os.environ.get("KUBECONFORM_K8S_VERSION", "1.36.0")


def is_manifest(path: str) -> bool:
    p = Path(path)
    if not path.startswith(MANIFEST_ROOTS) or p.suffix not in (".yaml", ".yml", ".json"):
        return False
    if any(rx.search(path) for rx in EXCLUDE) or not p.is_file():
        return False
    text = p.read_text(errors="replace")
    if "{{" in text:
        return False
    if p.suffix == ".json":
        return '"apiVersion"' in text and '"kind"' in text
    return bool(re.search(r"^apiVersion:", text, re.M) and re.search(r"^kind:", text, re.M))


def changed_files() -> list[str]:
    base = os.environ.get("YAML_BASE_SHA", "")
    if not base or set(base) == {"0"}:
        base = subprocess.check_output(["git", "rev-parse", "HEAD^"], text=True).strip()
    if not re.fullmatch(r"[0-9a-f]{40}", base):
        raise SystemExit("Expected a commit SHA for manifest validation")
    out = subprocess.check_output(["git", "diff", "--name-only", "--diff-filter=ACMR", "-z", base, "HEAD"])
    return [p for p in out.decode().split("\0") if p]


def all_files() -> list[str]:
    out = subprocess.check_output(["git", "ls-files", "-z", *MANIFEST_ROOTS])
    return [p for p in out.decode().split("\0") if p]


def main() -> int:
    if shutil.which("kubeconform") is None:
        print("kubeconform is not installed", file=sys.stderr)
        return 2
    candidates = all_files() if "--all" in sys.argv else changed_files()
    files = sorted(p for p in candidates if is_manifest(p))
    if not files:
        print("No changed plain manifests to validate")
        return 0
    cmd = [
        "kubeconform",
        "-summary",
        "-ignore-missing-schemas",
        "-strict",
        "-kubernetes-version", K8S_VERSION,
        "-schema-location", "default",
        "-schema-location", CRD_SCHEMA,
    ]
    if os.environ.get("KUBECONFORM_SKIP"):
        cmd += ["-skip", os.environ["KUBECONFORM_SKIP"]]
    cmd += files
    print(f"kubeconform: validating {len(files)} manifest file(s)")
    return subprocess.run(cmd).returncode


if __name__ == "__main__":
    sys.exit(main())
