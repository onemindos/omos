#!/usr/bin/env python3
"""Bounded Hermes API checks. Metadata only unless --chat-probe is explicit.

No credentials, response bodies, prompts, or raw exception text are logged.
An HTTP 200 SSE stream ending in an error is a failed probe, not success.
"""

import argparse
import json
import os
import sys
import time
import urllib.error
import urllib.parse
import urllib.request

MAX_BYTES = 256 * 1024
MARKER = "ONEMIND_LEGACY_API_OK"


class ProbeError(Exception):
    """Safe, fixed-code error suitable for an operator report."""


class NoRedirect(urllib.request.HTTPRedirectHandler):
    def redirect_request(self, req, fp, code, msg, headers, newurl):
        raise ProbeError("redirect_refused")


def check_stream(lines):
    """Validate a bounded OpenAI-compatible SSE response without exposing text."""
    content = []
    finishes = []
    count = 0
    done = False
    for line in lines:
        count += len(line)
        if count > MAX_BYTES:
            raise ProbeError("response_too_large")
        try:
            line = line.decode("utf-8").strip()
        except UnicodeDecodeError:
            raise ProbeError("invalid_stream_encoding") from None
        if not line or line.startswith(":"):
            continue
        if line == "event: error":
            raise ProbeError("inference_stream_error")
        if not line.startswith("data:"):
            continue
        data = line[5:].strip()
        if data == "[DONE]":
            done = True
            break
        try:
            chunk = json.loads(data)
        except ValueError:
            raise ProbeError("invalid_stream_json") from None
        if not isinstance(chunk, dict) or "error" in chunk:
            raise ProbeError("inference_stream_error")
        choices = chunk.get("choices", [])
        if not isinstance(choices, list):
            raise ProbeError("invalid_stream_shape")
        for choice in choices:
            if not isinstance(choice, dict):
                raise ProbeError("invalid_stream_shape")
            delta = choice.get("delta", {})
            if not isinstance(delta, dict):
                raise ProbeError("invalid_stream_shape")
            if delta.get("tool_calls") or delta.get("function_call"):
                raise ProbeError("unexpected_tool_event")
            text = delta.get("content")
            if text is not None:
                if not isinstance(text, str):
                    raise ProbeError("invalid_stream_shape")
                content.append(text)
            finish = choice.get("finish_reason")
            if finish is not None:
                if finish != "stop":
                    raise ProbeError("inference_stream_error")
                finishes.append(finish)
    if not done or not finishes or "".join(content).strip() != MARKER:
        raise ProbeError("incomplete_or_unexpected_answer")
    return {"answer_matches": True, "terminal_success": True}


class Client:
    def __init__(self, base_url, key, timeout):
        parsed = urllib.parse.urlsplit(base_url)
        if (parsed.scheme not in ("http", "https") or not parsed.hostname
                or parsed.username or parsed.password or parsed.query or parsed.fragment
                or parsed.path not in ("", "/")):
            raise ProbeError("invalid_base_url")
        host = parsed.hostname
        private_http = host in ("localhost", "127.0.0.1", "::1") or host.endswith(".svc.cluster.local")
        if parsed.scheme == "http" and not private_http:
            raise ProbeError("plaintext_remote_url_refused")
        self.base_url = base_url.rstrip("/")
        self.key = key
        self.timeout = timeout
        self.opener = urllib.request.build_opener(NoRedirect())

    def request(self, path, authenticated=True, body=None):
        headers = {"Accept": "application/json"}
        if authenticated:
            headers["Authorization"] = "Bearer " + self.key
        data = None
        if body is not None:
            headers.update({"Content-Type": "application/json", "Accept": "text/event-stream"})
            data = json.dumps(body).encode()
        req = urllib.request.Request(self.base_url + path, data=data, headers=headers)
        try:
            return self.opener.open(req, timeout=self.timeout)
        except urllib.error.HTTPError as exc:
            if exc.code in (301, 302, 303, 307, 308):
                exc.close()
                raise ProbeError("redirect_refused") from None
            # Return only the status, never the upstream error body.
            code = exc.code
            exc.close()
            raise ProbeError("http_" + str(code)) from None
        except (urllib.error.URLError, OSError, ValueError):
            raise ProbeError("transport_error") from None

    def get_json(self, path, authenticated=True):
        with self.request(path, authenticated=authenticated) as response:
            if response.status != 200:
                raise ProbeError("unexpected_http_status")
            if response.headers.get_content_type() != "application/json":
                raise ProbeError("expected_json_not_html")
            data = response.read(MAX_BYTES + 1)
        if len(data) > MAX_BYTES:
            raise ProbeError("response_too_large")
        try:
            data = json.loads(data)
        except (ValueError, UnicodeDecodeError):
            raise ProbeError("invalid_json") from None
        if not isinstance(data, dict):
            raise ProbeError("invalid_json_shape")
        return data

    def chat(self):
        # This is a real inference request; the model's instruction is not a
        # permission boundary for the harness. Never run without authorization.
        body = {"model": "hermes-agent", "messages": [{"role": "user", "content":
                "Connectivity verification only. Do not use tools, read files, delegate, "
                "run commands or change anything. Reply exactly: " + MARKER}],
                "stream": True, "max_tokens": 32, "tools": [], "tool_choice": "none"}
        deadline = time.monotonic() + self.timeout
        with self.request("/v1/chat/completions", body=body) as response:
            if response.status != 200 or response.headers.get_content_type() != "text/event-stream":
                raise ProbeError("expected_sse")

            def lines():
                remaining = MAX_BYTES
                while remaining >= 0:
                    if time.monotonic() >= deadline:
                        raise ProbeError("stream_deadline")
                    line = response.readline(min(remaining + 1, MAX_BYTES + 1))
                    if not line:
                        return
                    remaining -= len(line)
                    if remaining < 0:
                        raise ProbeError("response_too_large")
                    yield line

            return check_stream(lines())


def run(client, chat_probe=False):
    health = client.get_json("/health", authenticated=False)
    if health.get("status") != "ok" or health.get("platform") != "hermes-agent":
        raise ProbeError("unexpected_platform_health")
    for path in ("/v1/models", "/v1/capabilities"):
        try:
            client.get_json(path, authenticated=False)
        except ProbeError as exc:
            if str(exc) not in ("http_401", "http_403"):
                raise
        else:
            raise ProbeError("protected_endpoint_allows_anonymous")
    models = client.get_json("/v1/models")
    model_data = models.get("data")
    if not isinstance(model_data, list) or not any(
            isinstance(m, dict) and m.get("id") == "hermes-agent" for m in model_data):
        raise ProbeError("agent_model_missing")
    client.get_json("/v1/capabilities")
    report = {"metadata_passed": True, "capabilities": "advertised_not_executed",
              "inference": "not_tested", "provider_changed": False}
    if chat_probe:
        report.update({"inference": "passed", **client.chat()})
    return report


def main():
    parser = argparse.ArgumentParser(description=__doc__)
    parser.add_argument("--base-url", default="http://127.0.0.1:8642")
    parser.add_argument("--key-env", default="HERMES_API_KEY")
    parser.add_argument("--timeout", type=float, default=30)
    parser.add_argument("--chat-probe", action="store_true",
                        help="Run one REAL inference call; no retry or provider switch")
    args = parser.parse_args()
    try:
        key = os.environ.get(args.key_env, "")
        if not key:
            raise ProbeError("missing_key_environment")
        if not 0 < args.timeout <= 60:
            raise ProbeError("timeout_out_of_range")
        print(json.dumps(run(Client(args.base_url, key, args.timeout), args.chat_probe), sort_keys=True))
        return 0
    except ProbeError as exc:
        print(json.dumps({"passed": False, "reason": str(exc)}))
    except (OSError, ValueError):
        print(json.dumps({"passed": False, "reason": "transport_or_decode_error"}))
    return 1


if __name__ == "__main__":
    sys.exit(main())
