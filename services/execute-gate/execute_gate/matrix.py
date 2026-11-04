"""Minimal Matrix client-server API over urllib: send to the approvals room and long-poll sync.

The room must be unencrypted (the bot does not do E2EE). The sender of each event is the MXID
authenticated by Synapse; the gate only trusts that field.
"""

from __future__ import annotations

import json
import secrets
import time
import urllib.error
import urllib.parse
import urllib.request


class MatrixClient:
    def __init__(self, homeserver: str, token: str, room_id: str, timeout: float = 40.0):
        self.base = homeserver.rstrip("/")
        self._token = token
        self.room_id = room_id
        self.timeout = timeout

    def _req(self, method, path, body=None, timeout=None):
        data = None if body is None else json.dumps(body).encode()
        req = urllib.request.Request(self.base + path, data=data, method=method)
        req.add_header("Authorization", "Bearer " + self._token)
        if data is not None:
            req.add_header("Content-Type", "application/json")
        with urllib.request.urlopen(req, timeout=timeout or self.timeout) as resp:
            return json.loads(resp.read().decode() or "{}")

    def send(self, text: str) -> None:
        txn = secrets.token_hex(8)
        room = urllib.parse.quote(self.room_id, safe="")
        self._req("PUT", f"/_matrix/client/v3/rooms/{room}/send/m.room.message/{txn}",
                  {"msgtype": "m.text", "body": text})

    def whoami(self) -> str:
        return self._req("GET", "/_matrix/client/v3/account/whoami")["user_id"]

    def sync(self, since: str | None, timeout_ms: int = 30000) -> dict:
        flt = json.dumps({"room": {"rooms": [self.room_id], "timeline": {"limit": 20},
                                   "state": {"types": []}, "ephemeral": {"types": []}},
                          "presence": {"types": []}, "account_data": {"types": []}})
        q = {"timeout": str(timeout_ms if since else 0), "filter": flt}
        if since:
            q["since"] = since
        return self._req("GET", "/_matrix/client/v3/sync?" + urllib.parse.urlencode(q),
                         timeout=timeout_ms / 1000 + 15)


def events_from_sync(resp: dict, room_id: str):
    """Yield (event_id, sender, msgtype, body, edited) for text messages in the room."""
    room = (resp.get("rooms", {}).get("join", {}) or {}).get(room_id)
    if not room:
        return
    for ev in room.get("timeline", {}).get("events", []):
        if ev.get("type") != "m.room.message":
            continue
        content = ev.get("content", {}) or {}
        rel = content.get("m.relates_to") or {}
        yield (ev.get("event_id"), ev.get("sender"), content.get("msgtype"), content.get("body"),
               rel.get("rel_type") == "m.replace" or "m.new_content" in content)


def sync_loop(client: MatrixClient, gate, stop, log=print) -> None:
    """Run until stop is set. First sync only fetches a token, so history is never replayed."""
    since = None
    backoff = 1.0
    while not stop.is_set():
        try:
            resp = client.sync(since)
            first = since is None
            since = resp.get("next_batch", since)
            backoff = 1.0
            if first:
                continue
            for ev_id, sender, msgtype, body, edited in events_from_sync(resp, client.room_id):
                gate.handle_message(client.room_id, sender, ev_id, body, msgtype or "m.text", edited)
        except (urllib.error.URLError, OSError, ValueError, KeyError) as exc:
            log(json.dumps({"svc": "execute-gate", "event": "matrix_sync_error", "error": type(exc).__name__}))
            stop.wait(backoff)
            backoff = min(backoff * 2, 30.0)
        time.sleep(0)
