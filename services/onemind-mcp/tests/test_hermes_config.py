import importlib.util
from pathlib import Path
import secrets

import pytest
import yaml

path = Path(__file__).resolve().parents[3] / "platform/hermes/agent/onemind_mcp_config.py"
spec = importlib.util.spec_from_file_location("onemind_mcp_config", path)
config = importlib.util.module_from_spec(spec)
spec.loader.exec_module(config)


def test_preserves_configuration_and_merges_only_owned_keys(tmp_path):
    path = tmp_path / "config.yaml"
    original = {"model": {"provider": "existing"}, "channels": {"matrix": True},
                "mcp_servers": {"other": {"url": "https://other.example/mcp"}},
                "tools": {"tool_search": {"listing_max_tokens": 1200}, "other_setting": True}}
    path.write_text(yaml.safe_dump(original))
    before = path.read_bytes()
    config.main(path, True, header_supported=True, search_supported=True)
    actual = yaml.safe_load(path.read_text())
    assert actual["model"] == original["model"]
    assert actual["channels"] == original["channels"]
    assert actual["mcp_servers"]["other"] == original["mcp_servers"]["other"]
    assert actual["mcp_servers"]["onemind"]["headers"]["Authorization"] == config.HEADER
    assert actual["tools"]["tool_search"] == {"enabled": "on", "listing_max_tokens": 1200}
    assert path.with_name("config.yaml.before-onemind-mcp").read_bytes() == before
    assert path.stat().st_mode & 0o777 == 0o600
    assert path.with_name("config.yaml.before-onemind-mcp").stat().st_mode & 0o777 == 0o600
    config.main(path, True, header_supported=True, search_supported=True)
    assert path.with_name("config.yaml.before-onemind-mcp").read_bytes() == before


def test_missing_credentials_is_noop_and_unsupported_headers_refuse(tmp_path):
    path = tmp_path / "config.yaml"
    path.write_text("mcp_servers: {}\n")
    before = path.read_bytes()
    config.main(path, False)
    assert path.read_bytes() == before
    with pytest.raises(RuntimeError, match="activation refused"):
        config.main(path, True, header_supported=False)
    assert path.read_bytes() == before


def test_never_inlines_credentials_and_no_tool_search_fallback(tmp_path, monkeypatch):
    token = secrets.token_urlsafe(32)
    monkeypatch.setenv("ONEMIND_MCP_TOKEN", token)
    path = tmp_path / "config.yaml"
    path.write_text("mcp_servers: {}\n")
    config.main(path, True, header_supported=True, search_supported=False)
    assert token not in path.read_text()
    assert "${ONEMIND_MCP_TOKEN}" in path.read_text()
    assert "tools" not in yaml.safe_load(path.read_text())


def test_existing_unrelated_entry_is_never_overwritten():
    with pytest.raises(ValueError, match="unrelated"):
        config.merge({"mcp_servers": {"onemind": {"url": "https://unrelated.example/mcp"}}}, True)
