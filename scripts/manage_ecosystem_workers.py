#!/usr/bin/env python3
"""Plan/apply the isolated CPU pool with private remote state and local credentials."""
import argparse
import base64
import json
import os
from pathlib import Path
import subprocess
import shutil

ROOT = Path(__file__).resolve().parents[1] / "infra" / "ecosystem-workers"
EXPECTED_CLUSTER = "051f1134-ff05-451c-a52a-aab459feec6e"


def capture(argv):
    return subprocess.check_output(argv, text=True).strip()


def access_token():
    # Older doctl builds print auth help with exit code zero for `auth token`.
    # macOS can have both an older local binary and a current Homebrew binary.
    candidates = [shutil.which("doctl"), "/opt/homebrew/bin/doctl", "/usr/local/bin/doctl"]
    for binary in dict.fromkeys(candidates):
        if not binary or not Path(binary).is_file():
            continue
        result = subprocess.run([binary, "auth", "token"], text=True,
                                stdout=subprocess.PIPE, stderr=subprocess.PIPE)
        token = result.stdout.strip()
        if (result.returncode == 0 and len(token) >= 40 and
                token.isascii() and all(33 <= ord(c) <= 126 for c in token)):
            return token
    raise RuntimeError("A doctl version supporting 'auth token' and a valid login are required; update doctl")


def environment():
    context = capture(["kubectl", "config", "current-context"])
    if context != "do-nyc3-omos-cloud":
        raise RuntimeError("Select the verified do-nyc3-omos-cloud context first")
    secret = json.loads(capture(["kubectl", "get", "secret", "platform-backup-credentials",
                                 "-n", "ops", "-o", "json"]))["data"]
    decode = lambda key: base64.b64decode(secret[key]).decode()
    env = os.environ.copy()
    env.update(DIGITALOCEAN_TOKEN=access_token(),
               AWS_ACCESS_KEY_ID=decode("S3_ACCESS_KEY"),
               AWS_SECRET_ACCESS_KEY=decode("S3_SECRET_KEY"),
               AWS_EC2_METADATA_DISABLED="true")
    endpoint = decode("S3_ENDPOINT")
    if not endpoint.startswith("https://"):
        endpoint = "https://" + endpoint
    return env, decode("S3_BUCKET"), endpoint


def validate_scope(plan, env):
    data = json.loads(subprocess.check_output(["terraform", "show", "-json", str(plan)], cwd=ROOT, env=env))
    changed = [r for r in data.get("resource_changes", []) if r["change"]["actions"] != ["no-op"]]
    if len(changed) != 1:
        raise RuntimeError("Expected exactly one new node-pool resource")
    r = changed[0]
    if r["type"] != "digitalocean_kubernetes_node_pool" or r["change"]["after"]["cluster_id"] != EXPECTED_CLUSTER:
        raise RuntimeError("Plan changes resources outside the authorized CPU pool")
    actions = r["change"]["actions"]
    if actions == ["create"]:
        return "one new CPU pool"
    if actions == ["update"]:
        # An in-place update of the existing pool is allowed only for its autoscaling range.
        before, after = r["change"]["before"] or {}, r["change"]["after"] or {}
        changed_keys = sorted(k for k in set(before) | set(after) if before.get(k) != after.get(k)
                              and not (after.get(k) is None and k in (r["change"].get("after_unknown") or {})))
        if not changed_keys or not set(changed_keys) <= {"min_nodes", "max_nodes"}:
            raise RuntimeError("Pool update touches more than the autoscaling range: %s" % changed_keys)
        return "in-place change to the pool's %s only" % ", ".join(changed_keys)
    raise RuntimeError("Unexpected plan action for the CPU pool: %s" % actions)


def main():
    parser = argparse.ArgumentParser()
    parser.add_argument("action", choices=["plan", "apply"])
    parser.add_argument("--plan", type=Path, default=Path("/tmp/onemind-ecosystem-workers.tfplan"))
    args = parser.parse_args()
    env, bucket, endpoint = environment()
    if args.action == "plan":
        subprocess.run(["terraform", "init", "-input=false", "-backend-config=bucket=" + bucket,
                        "-backend-config=key=terraform/ecosystem-workers.tfstate",
                        "-backend-config=region=us-east-1", '-backend-config=endpoints={s3="' + endpoint + '"}',
                        "-backend-config=skip_credentials_validation=true",
                        "-backend-config=skip_region_validation=true",
                        "-backend-config=skip_requesting_account_id=true",
                        "-backend-config=skip_metadata_api_check=true", "-backend-config=use_path_style=true"],
                       cwd=ROOT, env=env, check=True)
        subprocess.run(["terraform", "validate"], cwd=ROOT, env=env, check=True)
        subprocess.run(["terraform", "plan", "-input=false", "-out=" + str(args.plan)], cwd=ROOT, env=env, check=True)
        scope = validate_scope(args.plan, env)
        print("Verified: %s; nothing else changes." % scope)
    else:
        validate_scope(args.plan, env)
        subprocess.run(["terraform", "apply", "-input=false", str(args.plan)], cwd=ROOT, env=env, check=True)


if __name__ == "__main__":
    main()
