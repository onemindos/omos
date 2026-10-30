"""Gate core: per-call pending approvals, owner-only approval by typed word, single use, TTL."""

from __future__ import annotations

import json
import re
import secrets
import threading
import time
from collections import OrderedDict

from .audit import Audit
from .canonical import call_hash, parse_arguments, summarise
from .policy import PHYSICAL, WORD_FOR_CLASS, Policy

TTL_SECONDS = 600
_CMD = re.compile(r"^\s*(execute|override)\s+([0-9a-fA-F]{4})\s*$", re.IGNORECASE)
_HREC = "h."
_IREC = "i."


def _enc(rec: dict) -> bytes:
    return json.dumps(rec, sort_keys=True).encode()


def _dec(raw: bytes) -> dict:
    return json.loads(raw.decode())


class Gate:
    def __init__(self, store, policy: Policy, owner: str, approvals_room: str, bot_user: str,
                 notifier, audit: Audit, ttl: int = TTL_SECONDS, clock=time.time):
        self.store = store
        self.policy = policy
        self.owner = owner
        self.room = approvals_room
        self.bot_user = bot_user
        self.notifier = notifier  # object with send(text) -> None, raising on failure
        self.audit = audit
        self.ttl = ttl
        self.clock = clock
        self._lock = threading.Lock()
        self._seen = OrderedDict()

    # ---------------------------------------------------------------- AuthZEN side

    def evaluate(self, subject: str, resource: str, arguments) -> tuple:
        """Return (decision, reason, request_id). Fails closed on every error."""
        try:
            with self._lock:
                return self._evaluate(subject, resource, arguments)
        except Exception as exc:  # any error is a deny
            self.audit.emit("error", resource=resource, error=type(exc).__name__)
            return False, "error", None

    def _evaluate(self, subject, resource, arguments):
        if not isinstance(resource, str) or not resource or not isinstance(subject, str):
            self.audit.emit("denied", reason="bad_request")
            return False, "bad_request", None
        if not self.policy.subject_allowed(subject):
            self.audit.emit("denied", reason="subject_not_allowed", subject=subject, resource=resource)
            return False, "subject_not_allowed", None
        try:
            args = parse_arguments(arguments)
            h = call_hash(resource, args)
        except (ValueError, TypeError):
            self.audit.emit("denied", reason="bad_arguments", resource=resource)
            return False, "bad_arguments", None

        cls = self.policy.classify(resource)
        now = self.clock()
        hkey = _HREC + h
        entry = self.store.get(hkey)
        if entry is not None:
            rec, rev = _dec(entry[0]), entry[1]
            if now >= rec["expires"]:
                self._drop(rec, rev)
                self.audit.emit("expired", id=rec["id"], resource=resource, hash=h[:12])
            elif rec["state"] == "approved":
                if self.store.delete(hkey, rev):
                    self._del_index(rec["id"])
                    self.audit.emit("consumed", id=rec["id"], resource=resource, cls=rec["cls"], hash=h[:12])
                    return True, "approved", rec["id"]
                self.audit.emit("denied", reason="consume_race", id=rec["id"], resource=resource)
                return False, "consume_race", rec["id"]
            else:
                if not rec.get("announced"):
                    self._announce(rec, hkey, rev)
                self.audit.emit("pending", id=rec["id"], resource=resource, cls=rec["cls"], hash=h[:12])
                return False, "pending_approval", rec["id"]

        return self._open_request(resource, cls, h, args, now)

    def _open_request(self, resource, cls, h, args, now):
        rid = None
        for _ in range(16):
            cand = secrets.token_hex(2).upper()
            if self.store.create(_IREC + cand, h.encode()):
                rid = cand
                break
        if rid is None:
            self.audit.emit("denied", reason="no_free_id", resource=resource)
            return False, "no_free_id", None
        rec = {"id": rid, "hash": h, "resource": resource, "cls": cls, "state": "pending",
               "created": now, "expires": now + self.ttl, "announced": False,
               "summary": summarise(args)}
        if not self.store.create(_HREC + h, _enc(rec)):
            self._del_index(rid)
            self.audit.emit("denied", reason="create_race", resource=resource)
            return False, "create_race", None
        entry = self.store.get(_HREC + h)
        if entry:
            self._announce(rec, _HREC + h, entry[1])
        self.audit.emit("requested", id=rid, resource=resource, cls=cls, hash=h[:12])
        return False, "approval_required", rid

    def _announce(self, rec, hkey, rev) -> None:
        word = WORD_FOR_CLASS[rec["cls"]]
        mins = max(1, int(self.ttl // 60))
        text = (
            f"Approval needed [{rec['id']}]\n"
            f"Tool: {rec['resource']} ({rec['cls']})\n"
            f"Arguments: {rec['summary']}\n"
            f"Reply `{word} {rec['id']}` within {mins} minutes to allow this one call once."
        )
        try:
            self.notifier.send(text)
        except Exception as exc:
            self.audit.emit("announce_failed", id=rec["id"], error=type(exc).__name__)
            return
        rec["announced"] = True
        self.store.update(hkey, _enc(rec), rev)

    def _drop(self, rec, rev) -> None:
        self.store.delete(_HREC + rec["hash"], rev)
        self._del_index(rec["id"])

    def _del_index(self, rid) -> None:
        e = self.store.get(_IREC + rid)
        if e:
            self.store.delete(_IREC + rid, e[1])

    # ---------------------------------------------------------------- Matrix side

    def handle_message(self, room_id, sender, event_id, body, msgtype="m.text", edited=False) -> str:
        """Process one room message. Returns an outcome label (for tests and logs)."""
        try:
            with self._lock:
                return self._handle(room_id, sender, event_id, body, msgtype, edited)
        except Exception as exc:
            self.audit.emit("error", where="handle_message", error=type(exc).__name__)
            return "error"

    def _reply(self, text) -> None:
        try:
            self.notifier.send(text)
        except Exception:
            pass

    def _handle(self, room_id, sender, event_id, body, msgtype, edited):
        if room_id != self.room or sender == self.bot_user:
            return "ignored"
        m = _CMD.match(body) if isinstance(body, str) else None
        if not m:
            return "ignored"
        if event_id:
            if event_id in self._seen:
                return "replay"
            self._seen[event_id] = True
            while len(self._seen) > 2048:
                self._seen.popitem(last=False)
        word, rid = m.group(1).lower(), m.group(2).upper()
        if sender != self.owner:
            self.audit.emit("denied", reason="not_owner", sender=sender, id=rid)
            return "not_owner"
        if msgtype != "m.text" or edited:
            self.audit.emit("denied", reason="not_plain_message", id=rid)
            return "not_plain_message"
        idx = self.store.get(_IREC + rid)
        if idx is None:
            self.audit.emit("denied", reason="unknown_id", id=rid)
            self._reply(f"No open request {rid}.")
            return "unknown_id"
        hkey = _HREC + idx[0].decode()
        entry = self.store.get(hkey)
        if entry is None:
            self.audit.emit("denied", reason="unknown_id", id=rid)
            self._reply(f"No open request {rid}.")
            return "unknown_id"
        rec, rev = _dec(entry[0]), entry[1]
        if self.clock() >= rec["expires"]:
            self._drop(rec, rev)
            self.audit.emit("denied", reason="expired", id=rid)
            self._reply(f"Request {rid} expired.")
            return "expired"
        if rec["state"] != "pending":
            self.audit.emit("denied", reason="already_approved", id=rid)
            return "already_approved"
        if WORD_FOR_CLASS[rec["cls"]] != word:
            self.audit.emit("denied", reason="wrong_word", id=rid, cls=rec["cls"], word=word)
            self._reply(f"{rid} is {rec['cls']}: type `{WORD_FOR_CLASS[rec['cls']]} {rid}`.")
            return "wrong_word"
        rec["state"] = "approved"
        rec["approved_by"] = sender
        rec["approved_at"] = self.clock()
        if not self.store.update(hkey, _enc(rec), rev):
            self.audit.emit("denied", reason="update_race", id=rid)
            return "update_race"
        self.audit.emit("granted", id=rid, resource=rec["resource"], cls=rec["cls"], approver=sender,
                        hash=rec["hash"][:12])
        self._reply(f"Approved {rid}. Legacy may retry the call once.")
        return "granted"


__all__ = ["Gate", "TTL_SECONDS", "PHYSICAL"]
