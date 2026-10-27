"""Merge only OneMind's private MCP entry; preserve the volume-managed config."""
import importlib
import importlib.util
import os
from pathlib import Path
import sys
import stat
import tempfile

import yaml

URL = "http://onemind-mcp.hermes.svc.cluster.local:8000/mcp"
HEADER = "Bearer ${ONEMIND_MCP_TOKEN}"


def supports_header_expansion():
    for name in ("tools.mcp_tool_config", "tools.mcp_tool"):
        try:
            interpolate = getattr(importlib.import_module(name), "_interpolate_env_vars", None)
            if interpolate:
                rendered = interpolate({"headers": {"Authorization": "Bearer ${PATH}"}})
                if rendered["headers"]["Authorization"] == "Bearer " + os.environ["PATH"]:
                    return True
        except Exception:
            continue
    return False


def merge(config, tool_search_supported):
    servers = config.setdefault("mcp_servers", {})
    if not isinstance(servers, dict):
        raise ValueError("mcp_servers must be a mapping")
    prior = servers.get("onemind")
    if prior is not None and (not isinstance(prior, dict) or prior.get("url") != URL):
        raise ValueError("Existing unrelated onemind MCP entry: choose a different key before activation")
    entry = dict(prior or {})
    headers = entry.get("headers", {})
    if not isinstance(headers, dict):
        raise ValueError("MCP headers must be a mapping")
    headers = {key: value for key, value in headers.items() if key.lower() != "authorization"}
    entry.update({"url": URL, "headers": {**headers, "Authorization": HEADER}})
    servers["onemind"] = entry
    if tool_search_supported:
        tools = config.setdefault("tools", {})
        if not isinstance(tools, dict):
            raise ValueError("tools must be a mapping")
        search = tools.get("tool_search")
        if isinstance(search, dict):
            search = {**search, "enabled": "on"}
        elif search is None or isinstance(search, bool):
            search = {"enabled": "on"}
        else:
            raise ValueError("Unsupported Tool Search configuration")
        tools["tool_search"] = search
    return config


def main(config_path, token_present, *, header_supported=None, search_supported=None):
    if not token_present:
        print("OneMind MCP credential absent; existing Hermes configuration unchanged")
        return 0
    if header_supported is None:
        header_supported = supports_header_expansion()
    if not header_supported:
        raise RuntimeError("This Hermes revision cannot expand MCP header references; activation refused")
    if search_supported is None:
        search_supported = importlib.util.find_spec("tools.tool_search") is not None
    path = Path(config_path)
    raw = path.read_bytes()
    config = yaml.safe_load(raw) or {}
    if not isinstance(config, dict):
        raise ValueError("Expected an existing Hermes mapping")
    merged = merge(config, search_supported)
    updated = yaml.safe_dump(merged, sort_keys=False, allow_unicode=True).encode()
    if updated == raw:
        print("OneMind MCP already configured")
        return 0
    original = path.stat()
    backup = path.with_name(path.name + ".before-onemind-mcp")
    try:
        descriptor = os.open(backup, os.O_WRONLY | os.O_CREAT | os.O_EXCL | os.O_NOFOLLOW, 0o600)
    except FileExistsError:
        if not stat.S_ISREG(backup.lstat().st_mode) or backup.stat().st_mode & 0o077:
            raise RuntimeError("Existing MCP backup must be a private regular file")
    else:
        with os.fdopen(descriptor, "wb") as handle:
            handle.write(raw)
            handle.flush()
            os.fsync(handle.fileno())
        os.chown(backup, original.st_uid, original.st_gid)
    descriptor, temp_name = tempfile.mkstemp(prefix=path.name + ".mcp.", dir=path.parent)
    try:
        with os.fdopen(descriptor, "wb") as handle:
            handle.write(updated)
            handle.flush()
            os.fsync(handle.fileno())
        # Keep credential-bearing legacy config private even if its previous mode was broad.
        os.chmod(temp_name, 0o600)
        os.chown(temp_name, original.st_uid, original.st_gid)
        os.replace(temp_name, path)
    finally:
        if os.path.exists(temp_name):
            os.unlink(temp_name)
    print("OneMind MCP configured with runtime credential reference; other settings preserved")
    return 0


if __name__ == "__main__":
    if len(sys.argv) != 2:
        raise SystemExit("Usage: onemind_mcp_config.py CONFIG_PATH")
    raise SystemExit(main(sys.argv[1], bool(os.getenv("ONEMIND_MCP_TOKEN", ""))))
