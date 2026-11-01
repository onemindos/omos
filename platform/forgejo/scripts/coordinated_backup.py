"""Quiesce exactly one Forgejo writer for an independently encrypted VolSync backup.

This program has no data mount, Secret reads, arbitrary workload selector or CLI.
Kubernetes RBAC constrains the two named resources further. It is not agent access.
"""
import json
from pathlib import Path
import signal
import ssl
import time
import urllib.request
import urllib.error
import uuid

NS = "/namespaces/forgejo/"
SCALE = "/apis/apps/v1" + NS + "deployments/forgejo/scale"
SOURCE = "/apis/volsync.backube/v1alpha1" + NS + "replicationsources/forgejo-data"
PODS = "/api/v1" + NS + "pods?labelSelector=app.kubernetes.io%2Fname%3Dforgejo"


class API:
    def __init__(self):
        account = Path("/var/run/secrets/kubernetes.io/serviceaccount")
        self.token = (account / "token").read_text().strip()
        self.ssl = ssl.create_default_context(cafile=str(account / "ca.crt"))
        self.host = "https://kubernetes.default.svc"

    def __call__(self, path, patch=None):
        headers = {"Authorization": "Bearer " + self.token}
        if patch is not None:
            headers["Content-Type"] = "application/merge-patch+json"
        request = urllib.request.Request(
            self.host + path, headers=headers,
            data=json.dumps(patch).encode() if patch is not None else None,
            method="PATCH" if patch is not None else "GET",
        )
        with urllib.request.urlopen(request, context=self.ssl, timeout=20) as response:
            return json.load(response)


def wait_for(check, timeout, clock=time.monotonic, sleep=time.sleep):
    deadline = clock() + timeout
    while not check():
        if clock() >= deadline:
            raise TimeoutError("Bounded backup wait expired")
        sleep(5)


def coordinate(api, wait=wait_for):
    before = api(SCALE)
    replicas = before["spec"]["replicas"]
    if replicas not in (0, 1):
        raise ValueError("Expected one writer or a deliberately stopped application")
    if api(SOURCE).get("spec", {}).get("paused", False):
        raise ValueError("Backup source is paused")
    marker = "forgejo-" + uuid.uuid4().hex
    changed = False
    try:
        if replicas:
            # resourceVersion prevents scaling a concurrently changed resource.
            # Treat a lost response as possibly committed, and attempt recovery.
            changed = True
            try:
                api(SCALE, {"metadata": {"resourceVersion": before["metadata"]["resourceVersion"]},
                            "spec": {"replicas": 0}})
            except urllib.error.HTTPError as error:
                if 400 <= error.code < 500:
                    changed = False  # Explicit rejection; preserve external state.
                raise
        wait(lambda: not api(PODS).get("items"), 180)
        api(SOURCE, {"spec": {"trigger": {"manual": marker}}})
        wait(lambda: api(SOURCE).get("status", {}).get("lastManualSync") == marker, 1800)
        print("Coordinated encrypted backup completed: " + marker, flush=True)
    finally:
        if changed:
            current = api(SCALE)
            if current["spec"]["replicas"] == 0:
                api(SCALE, {"metadata": {"resourceVersion": current["metadata"]["resourceVersion"]},
                            "spec": {"replicas": replicas}})
                print("Original Forgejo replica count restored", flush=True)
            else:
                print("Replica count changed externally; preserved operator state", flush=True)


def interrupted(_signum, _frame):
    raise InterruptedError("Backup interrupted; attempting replica recovery")


def main():
    signal.signal(signal.SIGTERM, interrupted)
    try:
        coordinate(API())
    except Exception as error:
        # Do not echo bearer tokens, URLs/query parameters or API response bodies.
        print("Forgejo backup failed (" + type(error).__name__ + "). Check Job status and restore desired replicas through Argo.", flush=True)
        return 1
    return 0


if __name__ == "__main__":
    raise SystemExit(main())
