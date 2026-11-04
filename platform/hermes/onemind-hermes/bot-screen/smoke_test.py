"""Isolated, offline image acceptance; never run against a real Hermes home."""
import os
from pathlib import Path
import socket
import stat
import subprocess
import time

from tools.bot_desktop import browser, runtime
from tools.computer_use.cua_backend_driver import resolve_cua_driver_cmd


def main():
    assert os.getuid() == 10000, "exercise the gateway UID, not root"
    home = Path(os.environ["HERMES_HOME"])
    assert str(home).startswith("/opt/bot-screen/scratch/"), "isolated home required"
    home.mkdir(parents=True, exist_ok=True)
    assert not runtime.missing_binaries(), runtime.missing_binaries()
    driver = resolve_cua_driver_cmd()
    assert driver, "PM-pinned cua-driver must resolve without runtime downloads"
    subprocess.run([driver, "--version"], check=True, timeout=20)
    # Import the actual authenticated WebSocket route; no websockify or VNC port needed.
    from hermes_cli.web_routers.display import router
    assert any(route.path == "/api/display/ws" for route in router.routes)
    exe = browser.executable()
    assert exe and "headless" not in exe, "upstream headed Chromium required"
    started = False
    proc = None
    try:
        st = runtime.start(wait_seconds=45)
        started = True
        assert st.running and st.installed, st.as_dict()
        env = {**os.environ, **runtime.published_env()}
        subprocess.run(["xdpyinfo"], env=env, check=True, stdout=subprocess.DEVNULL, timeout=10)
        path = Path(st.socket)
        assert stat.S_IMODE(path.stat().st_mode) == 0o600
        with socket.socket(socket.AF_UNIX) as client:
            client.settimeout(10)
            client.connect(str(path))
            assert client.recv(12).startswith(b"RFB 003."), "RFB greeting missing"
        # Launch the same headed browser argv as the dock, with no network dependency.
        proc = subprocess.Popen(browser.dock_argv(exe, str(home / "browser-test")) + ["about:blank"], env=env)
        port_file = home / "browser-test/DevToolsActivePort"
        for _ in range(100):
            if port_file.exists():
                break
            assert proc.poll() is None, "headed Chromium exited"
            time.sleep(0.1)
        assert port_file.exists(), "headed Chromium did not become ready"
        print("PASS: UID10000, packages, cua-driver, WebSocket route, Xfce/X11, private RFB, headed Chromium", flush=True)
    finally:
        if proc is not None:
            proc.terminate()
            proc.wait(timeout=15)
        if started:
            runtime.stop()


if __name__ == "__main__":
    main()
