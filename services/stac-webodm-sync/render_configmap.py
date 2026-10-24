#!/usr/bin/env python3
"""Write platform/stac/webodm-sync-code.configmap.yaml from sync.py.

The CronJob runs sync.py from a ConfigMap on an upstream image (no image of our own to build and pin). Run this
after editing sync.py:  python services/stac-webodm-sync/render_configmap.py
test_sync.py fails when the committed ConfigMap differs from what this renders.
"""
import os

HERE = os.path.dirname(os.path.abspath(__file__))
SOURCE = os.path.join(HERE, "sync.py")
TARGET = os.path.join(HERE, "..", "..", "platform", "stac", "webodm-sync-code.configmap.yaml")

HEADER = """\
# GENERATED from services/stac-webodm-sync/sync.py by services/stac-webodm-sync/render_configmap.py. Do not edit by hand:
# edit sync.py, run the script, commit both. Mounted at /app in the webodm-stac-sync CronJob (webodm-sync.yaml).
apiVersion: v1
kind: ConfigMap
metadata:
  name: webodm-sync-code
  namespace: stac
  labels:
    app.kubernetes.io/name: stac
    app.kubernetes.io/component: webodm-sync
data:
  sync.py: |
"""


def render() -> str:
    with open(SOURCE, encoding="utf-8") as f:
        body = f.read()
    lines = [("    " + line if line.strip() else "") for line in body.rstrip("\n").split("\n")]
    return HEADER + "\n".join(lines) + "\n"


if __name__ == "__main__":
    with open(TARGET, "w", encoding="utf-8") as f:
        f.write(render())
    print("wrote", os.path.normpath(TARGET))
