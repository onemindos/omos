"""Offline regression tests. Never call a provider or read a real credential."""

import importlib.util
import json
from pathlib import Path
import unittest
from unittest.mock import Mock

SPEC = importlib.util.spec_from_file_location(
    "legacy_api_preflight", Path(__file__).resolve().parents[1] / "ai-operator/legacy_api_preflight.py")
probe = importlib.util.module_from_spec(SPEC)
SPEC.loader.exec_module(probe)


def event(delta=None, finish=None):
    return ("data: " + json.dumps({"choices": [{"delta": delta or {}, "finish_reason": finish}]}) + "\n").encode()


class StreamTests(unittest.TestCase):
    def test_success_requires_answer_terminal_event_and_done(self):
        result = probe.check_stream([event({"content": probe.MARKER}), event(finish="stop"), b"data: [DONE]\n"])
        self.assertEqual(result, {"answer_matches": True, "terminal_success": True})

    def test_http_200_error_finish_is_failure(self):
        with self.assertRaisesRegex(probe.ProbeError, "inference_stream_error"):
            probe.check_stream([event(finish="error"), b"data: [DONE]\n"])

    def test_error_event_is_failure(self):
        with self.assertRaisesRegex(probe.ProbeError, "inference_stream_error"):
            probe.check_stream([b"event: error\n", b'data: {"error":"private upstream detail"}\n'])

    def test_top_level_error_body_not_reported(self):
        with self.assertRaises(probe.ProbeError) as cm:
            probe.check_stream([b'data: {"error":"private upstream detail"}\n'])
        self.assertNotIn("private", str(cm.exception))

    def test_early_eof_and_empty_success_are_failures(self):
        for lines in ([event({"content": probe.MARKER})], [event(finish="stop"), b"data: [DONE]\n"]):
            with self.subTest(lines=lines), self.assertRaises(probe.ProbeError):
                probe.check_stream(lines)

    def test_tool_events_fail(self):
        with self.assertRaisesRegex(probe.ProbeError, "unexpected_tool_event"):
            probe.check_stream([event({"tool_calls": [{"function": {"name": "fixture"}}]})])

    def test_limits_and_malformed_stream(self):
        for lines in ([b"data: " + b"x" * probe.MAX_BYTES], [b"data: not-json\n"],
                      [b"data: []\n"], [b'data: {"choices": [null]}\n'], [b"data: \xff\n"]):
            with self.subTest(), self.assertRaises(probe.ProbeError):
                probe.check_stream(lines)


class MetadataTests(unittest.TestCase):
    def client(self):
        client = Mock()

        def get(path, authenticated=True):
            if path == "/health":
                return {"status": "ok", "platform": "hermes-agent"}
            if not authenticated:
                raise probe.ProbeError("http_401")
            return {"data": [{"id": "hermes-agent"}]} if path == "/v1/models" else {"capabilities": {}}

        client.get_json.side_effect = get
        return client

    def test_default_never_calls_inference(self):
        client = self.client()
        result = probe.run(client)
        self.assertEqual(result["inference"], "not_tested")
        client.chat.assert_not_called()

    def test_auth_failure_stops_before_chat(self):
        client = self.client()
        client.get_json.side_effect = probe.ProbeError("http_401")
        with self.assertRaises(probe.ProbeError):
            probe.run(client, chat_probe=True)
        client.chat.assert_not_called()

    def test_successful_health_does_not_hide_chat_failure_or_retry(self):
        client = self.client()
        client.chat.side_effect = probe.ProbeError("inference_stream_error")
        with self.assertRaisesRegex(probe.ProbeError, "inference_stream_error"):
            probe.run(client, chat_probe=True)
        client.chat.assert_called_once()

    def test_no_redirects_with_credentials(self):
        with self.assertRaisesRegex(probe.ProbeError, "redirect_refused"):
            probe.NoRedirect().redirect_request(None, None, 302, "", {}, "https://untrusted.invalid")

    def test_no_plaintext_remote_or_embedded_credentials(self):
        for url in ("http://untrusted.invalid", "https://user:pass@host.invalid", "https://host.invalid/?key=private"):
            with self.subTest(url=url), self.assertRaises(probe.ProbeError):
                probe.Client(url, "fixture", 1)


if __name__ == "__main__":
    unittest.main()
