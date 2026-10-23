"""Merge the LiteLLM provider into Hermes without replacing existing config."""

import os
import sys
import tempfile
from pathlib import Path

import yaml


PROVIDER = "litellm"
MODEL = "hermes-agent"
ENDPOINT = "http://litellm.onemind-ai.svc.cluster.local:4000/v1"
KEY_VAR = "HERMES_LITELLM_API_KEY"


def merge(config, key_env):
    provider = {
        "base_url": ENDPOINT,
        "api_mode": "chat_completions",
        "key_env": KEY_VAR,
        "default_model": MODEL,
        "context_length": 8192,
    }
    if not key_env.strip():
        return config
    model = config.setdefault("model", {})
    providers = config.setdefault("providers", {})
    if not isinstance(model, dict) or not isinstance(providers, dict):
        raise ValueError("model and providers must be mappings")
    # Stale explicit endpoints/inline credentials override named-provider routing.
    for field in ("base_url", "api_key", "key_env", "api_key_env", "api_mode"):
        model.pop(field, None)
    model.update({"provider": PROVIDER, "default": MODEL})
    providers[PROVIDER] = provider
    return config


def main(config_path, key_env):
    if not key_env.strip():
        print("LiteLLM key absent; leaving current Hermes provider default")
        return 0
    path = Path(config_path)
    if not path.is_file():
        raise FileNotFoundError("Existing Hermes config is required before provider activation")
    with path.open(encoding="utf-8") as stream:
        config = yaml.safe_load(stream) or {}
    if not isinstance(config, dict):
        raise RuntimeError(f"Expected mapping in {path}, got {type(config).__name__}")
    merge(config, key_env)
    original = path.stat()
    descriptor, temp_name = tempfile.mkstemp(prefix=f"{path.name}.", dir=path.parent)
    try:
        with os.fdopen(descriptor, "w", encoding="utf-8") as stream:
            yaml.safe_dump(config, stream, sort_keys=False, allow_unicode=True)
            stream.flush()
            os.fsync(stream.fileno())
        os.chmod(temp_name, original.st_mode & 0o777)
        os.chown(temp_name, original.st_uid, original.st_gid)
        os.replace(temp_name, path)
    except BaseException:
        try:
            os.unlink(temp_name)
        except FileNotFoundError:
            pass
        raise
    print("LiteLLM provider merged; default is hermes-agent (local-tools fallback at LiteLLM)")
    return 0


if __name__ == "__main__":
    if len(sys.argv) != 2:
        raise SystemExit("Usage: litellm_provider_config.py CONFIG_PATH")
    raise SystemExit(main(sys.argv[1], os.getenv(KEY_VAR, "")))
