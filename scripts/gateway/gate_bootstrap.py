"""Create the execute-gate's Matrix bot and approvals room. Runs INSIDE the Synapse pod (stdlib only), started by
scripts/gateway/keys-session.sh through `kubectl exec -i deploy/synapse -- python3 -c "$(cat gate_bootstrap.py)"`.

stdin (one JSON object, never in argv):
  {"user": "gate", "password": "<generated>", "owner": "@zeus:onemindos.dev", "room_name": "approvals",
   "base": "http://localhost:8008", "register_cmd": [...optional, tests only...]}
stdout: one JSON object {"token", "room_id", "mxid"}. Progress goes to stderr and never contains the password or token.

Steps: check the owner exists, register `@<user>` WITHOUT admin with Synapse's register_new_matrix_user (the shared secret is read
from the pod's own /config/homeserver.yaml, so it does not depend on the matrix-secrets Secret), log in for a token, create a private
invite-only room without encryption (the gate cannot read encrypted events), invite the owner, and check the room is unencrypted.
Exit codes: 2 owner missing, 3 bot user exists, 4 register failed, 5 login failed, 6 room failed, 7 room is encrypted.
"""
import json
import subprocess
import sys
import urllib.error
import urllib.parse
import urllib.request


def log(msg):
    print(msg, file=sys.stderr, flush=True)


def call(base, method, path, body=None, token=None):
    req = urllib.request.Request(base + path, method=method, data=None if body is None else json.dumps(body).encode(),
                                 headers={"Content-Type": "application/json"})
    if token:
        req.add_header("Authorization", "Bearer " + token)
    try:
        with urllib.request.urlopen(req, timeout=30) as r:
            return r.status, json.loads(r.read() or b"{}")
    except urllib.error.HTTPError as e:
        try:
            return e.code, json.loads(e.read() or b"{}")
        except ValueError:
            return e.code, {}


def main():
    cfg = json.load(sys.stdin)
    user, pw, owner = cfg["user"], cfg["password"], cfg["owner"]
    base = cfg.get("base", "http://localhost:8008")
    server = owner.split(":", 1)[1]
    mxid = "@%s:%s" % (user, server)
    if mxid == owner:
        log("bot and owner must differ")
        return 2

    code, _ = call(base, "GET", "/_matrix/client/v3/profile/" + urllib.parse.quote(owner, safe=""))
    if code == 404:
        log("owner %s does not exist on this server" % owner)
        return 2
    log("owner %s: %s" % (owner, "exists" if code == 200 else "not checked (profile lookup needs auth, HTTP %s)" % code))

    cmd = cfg.get("register_cmd") or ["register_new_matrix_user", "-c", "/config/homeserver.yaml", "-u", user, "--no-admin", base]
    # The tool asks for the password twice on stdin when no -p is given (getpass falls back to stdin without a tty),
    # so the password is never in an argument list.
    r = subprocess.run(cmd, input=pw + "\n" + pw + "\n", capture_output=True, text=True, timeout=60)
    if r.returncode != 0:
        text = (r.stdout + r.stderr).replace(pw, "<password>")
        log("register_new_matrix_user failed (exit %d): %s" % (r.returncode, text.strip()[-400:]))
        return 3 if "taken" in text.lower() or "in use" in text.lower() or "exists" in text.lower() else 4
    log("registered %s (not admin)" % mxid)

    code, res = call(base, "POST", "/_matrix/client/v3/login", {
        "type": "m.login.password", "identifier": {"type": "m.id.user", "user": user}, "password": pw,
        "initial_device_display_name": "execute-gate"})
    if code != 200 or "access_token" not in res:
        log("login failed: HTTP %s" % code)
        return 5
    token = res["access_token"]
    log("logged in")

    code, res = call(base, "POST", "/_matrix/client/v3/createRoom", {
        "name": cfg.get("room_name", "approvals"), "topic": "Execute gate approvals: reply `execute <id>` or `override <id>`",
        "preset": "private_chat", "visibility": "private", "is_direct": False, "invite": [owner]}, token)
    if code != 200 or "room_id" not in res:
        log("create room failed: HTTP %s %s" % (code, res.get("error", "")))
        return 6
    room = res["room_id"]
    log("created private room %s and invited %s" % (room, owner))

    code, _ = call(base, "GET", "/_matrix/client/v3/rooms/%s/state/m.room.encryption" % urllib.parse.quote(room, safe=""), token=token)
    if code != 404:
        log("room is encrypted (HTTP %s): the gate cannot read it" % code)
        return 7
    log("room is unencrypted")
    print(json.dumps({"token": token, "room_id": room, "mxid": mxid}))
    return 0


if __name__ == "__main__":
    sys.exit(main())
