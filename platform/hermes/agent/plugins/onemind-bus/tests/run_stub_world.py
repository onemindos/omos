"""Run the stub entity, layer, policy and map services against a NATS server until Ctrl-C.

For a real-model check of the onemind-bus tools on a local stack (no cluster): point a local
Hermes with the plugin enabled at the same server, ask it the D11 requests, and watch the
printed view changes. The stub policy never dispatches anything.

    python run_stub_world.py nats://127.0.0.1:4222
"""
from __future__ import annotations

import json
import sys
import time
from pathlib import Path

sys.path.insert(0, str(Path(__file__).resolve().parent))

from stubs import World  # noqa: E402


def main() -> int:
    url = sys.argv[1] if len(sys.argv) > 1 else "nats://127.0.0.1:4222"
    w = World(url, Path(__file__).resolve().parent / "fixtures" / "entities.json")
    w.start()
    print(f"stub world on {url}; Ctrl-C to stop", flush=True)
    seen = 0
    try:
        while True:
            time.sleep(0.5)
            for subject, raw in w.map_msgs[seen:]:
                print(json.dumps({"map_action": subject, "payload": json.loads(raw)}), flush=True)
                print(json.dumps({"view": {"camera": w.view.camera, "visible": sorted(w.view.visible),
                                           "selected": sorted(w.view.selected),
                                           "annotations": sorted(w.view.annotations)}}), flush=True)
            seen = len(w.map_msgs)
            for p in w.proposals:
                if not p.get("_printed"):
                    p["_printed"] = True
                    print(json.dumps({"proposal": {k: v for k, v in p.items() if k != "_printed"}}), flush=True)
    except KeyboardInterrupt:
        pass
    finally:
        w.stop()
    return 0


if __name__ == "__main__":
    raise SystemExit(main())
