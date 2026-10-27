"""Exercise the built non-root image over HTTP with isolated mock upstream apps."""
from pathlib import Path
import subprocess
import secrets
import re
import shutil
import sys
import tempfile
import time
import uuid

ROOT = Path(__file__).resolve().parents[1]
MOCK = '''import sys, json, threading
sys.path.insert(0, "/app")
from http.server import BaseHTTPRequestHandler, ThreadingHTTPServer
class Handler(BaseHTTPRequestHandler):
    def do_GET(self):
        data = ({"status":"success","data":{"resultType":"vector","result":[]}}
                if self.path.startswith("/api/v1/query")
                else {"connections":4,"uptime":"2d","streams":2})
        self.send_response(200)
        self.send_header("Content-Type", "application/json")
        self.end_headers()
        self.wfile.write(json.dumps(data).encode())
    def log_message(self, *args): pass
threading.Thread(target=ThreadingHTTPServer(("127.0.0.1",9000),Handler).serve_forever,daemon=True).start()
from onemind_mcp.server import main
main()
'''


def main():
    image = sys.argv[1] if len(sys.argv) > 1 else "onemind-mcp:check"
    if not re.fullmatch(r"[A-Za-z0-9_./:@-]+", image):
        raise ValueError("Expected an image reference")
    name = "onemind-mcp-test-" + uuid.uuid4().hex[:10]
    with tempfile.TemporaryDirectory() as directory:
        context = Path(directory)
        (context / "launcher.py").write_text(MOCK)
        shutil.copyfile(ROOT / "config/apps.yaml", context / "apps.yaml")
        shutil.copyfile(ROOT / "scripts/smoke.py", context / "smoke.py")
        (context / "Dockerfile").write_text(
            f"FROM {image}\nCOPY --chmod=0444 launcher.py smoke.py /\n"
            "COPY --chmod=0444 apps.yaml /etc/onemind/apps.yaml\n")
        check_image = name + ":fixtures"
        subprocess.run(["docker", "build", "--network", "none", "-t", check_image, directory],
                       check=True, capture_output=True)
        args = ["docker", "run", "-d", "--name", name, "--network", "none", "--read-only",
                "--cap-drop", "ALL", "--security-opt", "no-new-privileges", "--memory", "256m"]
        for value in [f"ONEMIND_MCP_TOKEN={secrets.token_urlsafe(32)}",
                      "ONEMIND_ALLOWED_HOSTS=127.0.0.1:8000", "NATS_MONITOR_URL=http://127.0.0.1:9000",
                      "PROMETHEUS_URL=http://127.0.0.1:9000", "ONEMIND_MCP_URL=http://127.0.0.1:8000/mcp"]:
            args += ["-e", value]
        try:
            subprocess.run(args + [check_image, "python", "/launcher.py"], check=True, capture_output=True)
            for _ in range(30):
                result = subprocess.run(["docker", "exec", name, "python", "-c",
                    'import httpx; assert httpx.get("http://127.0.0.1:8000/readyz",trust_env=False).status_code==200'], capture_output=True)
                if result.returncode == 0:
                    break
                time.sleep(0.3)
            else:
                raise RuntimeError("Container readiness did not pass")
            subprocess.run(["docker", "exec", name, "python", "/smoke.py"], check=True)
            subprocess.run(["docker", "exec", name, "python", "-c",
                'import httpx; assert httpx.post("http://127.0.0.1:8000/mcp",json={},trust_env=False).status_code==401; print("PASS unauthenticated rejection")'], check=True)
        finally:
            subprocess.run(["docker", "rm", "-f", name], capture_output=True)
            subprocess.run(["docker", "image", "rm", check_image], capture_output=True)


if __name__ == "__main__":
    main()
