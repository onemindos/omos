"""Guards for the per-service NATS identities (FEATURES A19).

The identity spec renders into Helm values and the live nats-config ConfigMap. These tests keep
the security invariants (only policy publishes commands, agents never confirm, no reply spoofing,
no stream admin) checkable without a server; permtest/ runs the same config on nats-server.
"""
import copy
import json
import sys
import tempfile
import unittest
from pathlib import Path
from unittest import mock

ROOT = Path(__file__).resolve().parents[2]
sys.path.insert(0, str(ROOT / "scripts"))

import nats_identities as ni  # noqa: E402

# Shape of the live nats-config ConfigMap (read 2026-10-29, values replaced): JSON with bare
# tokens (45Gi, $SERVER_NAME) that nats-server accepts.
LIVE = """{
  "accounts": {
    "ONEMIND": {
      "imports": [{"service": {"account": "SYS", "subject": "$SYS.REQ.>"}}],
      "jetstream": {"max_consumers": -1, "max_file": -1, "max_mem": -1, "max_streams": -1},
      "users": [{"password": "live-hash-onemind", "user": "onemind"}]
    },
    "SYS": {
      "exports": [{"accounts": ["ONEMIND"], "service": "$SYS.REQ.>"}],
      "users": [{"password": "live-hash-sys", "user": "sys"}, {"password": "live-hash-monitor", "user": "monitor"}]
    }
  },
  "http_port": 8222,
  "jetstream": {
    "max_file_store": 45Gi,
    "max_memory_store": 0,
    "store_dir": "/data"
  },
  "leafnodes": {
    "authorization": {"users": [{"account": "ONEMIND", "password": "live-hash-leaf", "user": "leaf_tak"}]},
    "no_advertise": true,
    "port": 7422
  },
  "server_name": $SERVER_NAME,
  "system_account": "SYS"
}
"""


class InvariantTests(unittest.TestCase):
    def setUp(self):
        self.spec = ni.load_spec()

    def user(self, spec, name):
        return next(u for u in spec["users"] if u["user"] == name)

    def test_spec_holds(self):
        self.assertEqual(ni.check_invariants(self.spec), [])

    def test_values_files_current(self):
        args = mock.Mock(check=True)
        self.assertEqual(ni.cmd_values(args), 0)

    def test_only_policy_publishes_commands(self):
        for name, subj in (("ingest", "cmd.>"), ("zeus", "cmd.fleet.>"), ("hermes", "cmd.*.*.*.goto"), ("agent-fabric", "cmd.ha.>")):
            spec = copy.deepcopy(self.spec)
            self.user(spec, name)["publish"].append(subj)
            self.assertTrue(any("only policy publishes" in e for e in ni.check_invariants(spec)), (name, subj))

    def test_only_projector_injects_into_tak(self):
        spec = copy.deepcopy(self.spec)
        self.user(spec, "mesh-gateway")["publish"].append("tak.cot.inject")
        self.assertTrue(any("only projector, projector-gated" in e for e in ni.check_invariants(spec)))

    def test_only_entity_publishes_entities(self):
        spec = copy.deepcopy(self.spec)
        self.user(spec, "zeus")["publish"].append("ent.>")
        self.assertTrue(any("(only entity)" in e for e in ni.check_invariants(spec)))

    def test_no_inbox_publish_and_no_stream_admin(self):
        spec = copy.deepcopy(self.spec)
        self.user(spec, "layer")["publish"].append("_INBOX.>")
        self.user(spec, "entity")["publish"].append("$JS.API.STREAM.CREATE.X")
        errs = ni.check_invariants(spec)
        self.assertTrue(any("spoof replies" in e for e in errs))
        self.assertTrue(any("stream admin" in e for e in errs))

    def test_agents_never_confirm(self):
        spec = copy.deepcopy(self.spec)
        spec["accounts"]["AGENTS"]["imports"].append({"service": {"account": "ONEMIND", "subject": "svc.policy.command.confirm"}})
        self.assertTrue(any("never confirm" in e for e in ni.check_invariants(spec)))
        spec = copy.deepcopy(self.spec)
        self.user(spec, "hermes")["publish"].append("svc.policy.>")
        self.assertTrue(any("only humans (CLIENT)" in e for e in ni.check_invariants(spec)))

    def test_single_writers(self):
        # #569 D1-P6, #568 D2-P5, #554: one writer each.
        for name, subj, owner in (("policy", "$KV.geofences.>", "geofence"), ("policy", "$KV.mission-state.>", "mission"),
                                  ("agentops", "evt.mission.>", "mission"), ("zeus", "svc.mission.apply", "policy"),
                                  ("agentops", "tak.cot.inject", "projector")):
            spec = copy.deepcopy(self.spec)
            self.user(spec, name)["publish"].append(subj)
            self.assertTrue(any(f"may publish" in e and owner in e for e in ni.check_invariants(spec)), (name, subj))

    def test_services_never_confirm(self):
        # Only humans (CLIENT users) confirm; an ONEMIND service such as agentops cannot.
        spec = copy.deepcopy(self.spec)
        self.user(spec, "agentops")["publish"].append("svc.policy.*.confirm")
        self.assertTrue(any("only humans (CLIENT)" in e for e in ni.check_invariants(spec)))

    def test_callout_usernames_are_reserved(self):
        # #565 J2-P2: callout users are `ak:<user>;<groups>`.
        for bad in ("ak:zeus", "zeus;onemind-admins"):
            spec = copy.deepcopy(self.spec)
            self.user(spec, "zeus")["user"] = bad
            self.assertTrue(any("reserved for auth-callout" in e for e in ni.check_invariants(spec)), bad)

    def test_auth_account_has_no_jetstream(self):
        spec = copy.deepcopy(self.spec)
        spec["accounts"]["AUTH"]["jetstream"] = {"max_mem": -1}
        self.assertTrue(any("account AUTH" in e for e in ni.check_invariants(spec)))

    def test_no_credentials_in_spec(self):
        spec = copy.deepcopy(self.spec)
        self.user(spec, "policy")["password"] = "x"
        self.assertTrue(any("password field" in e for e in ni.check_invariants(spec)))
        for stage in ("A1", "A2", "B"):
            text = ni.VALUES[stage].read_text()
            for u in self.spec["users"]:
                if stage == "A1" and u["user"] == "leaf_tak":
                    continue  # stays in leafnodes.authorization, untouched by these values
                if u.get("test_only"):
                    self.assertNotIn(ni.placeholder(u["user"]), text)  # never rendered for live
                    continue
                self.assertIn(ni.placeholder(u["user"]), text)

    def test_client_and_agents_have_no_jetstream(self):
        spec = copy.deepcopy(self.spec)
        spec["accounts"]["CLIENT"]["jetstream"] = {"max_mem": -1}
        self.assertTrue(any("must not enable JetStream" in e for e in ni.check_invariants(spec)))


class RenderTests(unittest.TestCase):
    def setUp(self):
        self.spec = ni.load_spec()

    def users(self, accounts):
        return {u["user"]: (a, u) for a, v in accounts.items() for u in v["users"]}

    def test_every_new_user_is_restricted(self):
        accounts = ni.render_accounts(self.spec, "B", ni.placeholder)
        for name, (_, u) in self.users(accounts).items():
            if name in ("sys", "monitor"):
                self.assertNotIn("permissions", u)
                continue
            p = u["permissions"]
            if name in ("onemind", "leaf_tak"):
                self.assertIn("deny", p["publish"])
                self.assertIn("cmd.fleet.>", p["publish"]["deny"])
                self.assertIn("svc.policy.>", p["subscribe"]["deny"])
            else:
                self.assertIn("allow", p["publish"], name)
                self.assertIn("allow", p["subscribe"], name)

    def test_policy_receives_zone_results(self):
        # onemind-xr E2E run 1: without this grant zone.create/update/remove end no_ack.
        sub = self.users(ni.render_accounts(self.spec, "B", ni.placeholder))["policy"][1]["permissions"]["subscribe"]["allow"]
        self.assertIn("evt.map.geofence.>", sub)

    def test_test_only_humans_never_reach_live_renders(self):
        for stage in ni.STAGES:
            self.assertNotIn("alice", self.users(ni.render_accounts(self.spec, stage, ni.placeholder)))
        for stage in ni.VALUES:
            self.assertNotIn("alice", ni.render_values(self.spec, stage))
        test = self.users(ni.render_accounts(self.spec, "B", ni.placeholder, test_users=True))
        (acc_a, alice), (acc_z, zeus) = test["alice"], test["zeus"]
        self.assertEqual((acc_a, acc_z), ("CLIENT", "CLIENT"))
        self.assertEqual(alice["permissions"], zeus["permissions"])
        self.assertEqual(alice["allowed_connection_types"], zeus["allowed_connection_types"])
        humans = [n for n, (a, _) in test.items() if a == "CLIENT"]
        self.assertGreaterEqual(len(humans), 2, "A4 actions need two distinct humans in test configs")

    def test_test_only_rules(self):
        spec = copy.deepcopy(self.spec)
        spec["users"].append({"user": "svc-x", "account": "ONEMIND", "stage": "A", "test_only": True,
                              "secret": {"namespace": "a", "name": "b", "user_key": "u", "password_key": "p"}})
        errs = ni.check_invariants(spec)
        self.assertTrue(any("only for CLIENT humans" in e for e in errs), errs)
        self.assertTrue(any("never sealed" in e for e in errs), errs)

    def test_conf_test_users_flag(self):
        with tempfile.TemporaryDirectory() as d:
            pw = Path(d) / "pw.json"
            pw.write_text(json.dumps({u["user"]: "x" for u in self.spec["users"]}))
            out = Path(d) / "conf.json"
            for flag, want in (([], False), (["--test-users"], True)):
                argv = ["nats_identities.py", "conf", "--passwords", str(pw), "--out", str(out)] + flag
                with mock.patch.object(sys, "argv", argv), mock.patch("sys.stdout"):
                    self.assertEqual(ni.main(), 0)
                users = [u["user"] for u in json.loads(out.read_text())["accounts"]["CLIENT"]["users"]]
                self.assertEqual("alice" in users, want)

    def test_standin_leaf_conf(self):
        conf = ni.render_leaf_conf(self.spec, "hermes", "pw", remote="127.0.0.1:17422", plain=True, port=14222)
        (remote,) = conf["leafnodes"]["remotes"]
        self.assertEqual(remote["urls"], ["nats://hermes:pw@127.0.0.1:17422"])
        self.assertNotIn("tls", remote)
        self.assertEqual(conf["port"], 14222)
        self.assertNotIn("http", conf)
        self.assertEqual(conf["no_auth_user"], "local")

    def test_stage_a1_is_additive(self):
        accounts = ni.render_accounts(self.spec, "A1", ni.placeholder)
        users = self.users(accounts)
        self.assertNotIn("leaf_tak", users)
        self.assertNotIn("permissions", users["onemind"][1])
        values = ni.render_values(self.spec, "A1")
        self.assertNotIn("leafnodes", values.split("config:", 1)[1].split("merge:", 1)[0])

    def test_stage_a2_moves_the_leaf_user(self):
        users = self.users(ni.render_accounts(self.spec, "A2", ni.placeholder))
        acc, leaf = users["leaf_tak"]
        self.assertEqual(acc, "ONEMIND")
        self.assertEqual(leaf["allowed_connection_types"], ["LEAFNODE"])
        self.assertNotIn("permissions", leaf)
        self.assertIn("authorization: null", ni.render_values(self.spec, "A2"))

    def test_legacy_connection_types_from_a2(self):
        a1 = self.users(ni.render_accounts(self.spec, "A1", ni.placeholder))
        for name in ("onemind", "sys", "monitor"):
            self.assertNotIn("allowed_connection_types", a1[name][1], name)
        for stage in ("A2", "B"):
            users = self.users(ni.render_accounts(self.spec, stage, ni.placeholder))
            self.assertEqual(users["onemind"][1]["allowed_connection_types"], ["STANDARD", "WEBSOCKET"])
            self.assertEqual(users["sys"][1]["allowed_connection_types"], ["STANDARD"])
            self.assertEqual(users["monitor"][1]["allowed_connection_types"], ["STANDARD"])

    def test_hermes_sidecar_carries_the_hub_allow_lists(self):
        hub = self.users(ni.render_accounts(self.spec, "B", ni.placeholder))["hermes"][1]
        frag = ni.sidecar_local(self.spec, "hermes")
        self.assertEqual(frag["no_auth_user"], "local")
        (local,) = frag["accounts"]["OMOS"]["users"]
        self.assertEqual(local["permissions"], hub["permissions"])
        self.assertNotIn("password", local)
        self.assertNotIn("cmd.>", local["permissions"]["publish"]["allow"])
        conf = ni.render_leaf_conf(self.spec, "hermes", "p@ss/w:rd")
        (remote,) = conf["leafnodes"]["remotes"]
        self.assertEqual(remote["urls"], ["tls://hermes:p%40ss%2Fw%3Ard@nats.nats.svc.cluster.local:7422"])
        self.assertEqual(remote["account"], "OMOS")
        self.assertEqual(conf["server_name"], "hermes-legacy-leaf")

    def test_sidecar_only_on_leaf_users(self):
        spec = copy.deepcopy(self.spec)
        self.user_of(spec, "policy")["sidecar"] = self.user_of(spec, "hermes")["sidecar"]
        self.assertTrue(any("leaf-only user" in e for e in ni.check_invariants(spec)))

    def user_of(self, spec, name):
        return next(u for u in spec["users"] if u["user"] == name)

    def test_stage_c_drops_shared_users_only(self):
        users = self.users(ni.render_accounts(self.spec, "C", ni.placeholder))
        self.assertNotIn("onemind", users)
        self.assertNotIn("leaf_tak", users)
        self.assertIn("sys", users)
        self.assertIn("policy", users)

    def test_jetstream_expansion_is_per_stream(self):
        js = ni.expand_jetstream({"consume": ["ENT"], "kv_write": ["geofences"]})
        self.assertIn("$JS.API.CONSUMER.MSG.NEXT.ENT.*", js)
        self.assertIn("$KV.geofences.>", js)
        self.assertNotIn("$JS.API.>", js)
        self.assertFalse(any(s.startswith("$JS.API.STREAM.CREATE") for s in js))


class LiveMergeTests(unittest.TestCase):
    def run_live(self, stage):
        with tempfile.TemporaryDirectory() as d:
            live = Path(d) / "live.conf"
            live.write_text(LIVE)
            pwf = Path(d) / "pw.json"
            spec = ni.load_spec()
            pwf.write_text(json.dumps({u["user"]: "test-only" for u in spec["users"]}))
            out = Path(d) / "new.conf"
            args = mock.Mock(live_conf=str(live), passwords=str(pwf), stage=stage, out=str(out))
            with mock.patch.object(ni, "bcrypt_hash", side_effect=lambda p: "hashed"), mock.patch("builtins.print"):
                self.assertEqual(ni.cmd_live(args), 0)
            return out.read_text()

    def test_live_round_trip_keeps_bare_tokens_and_live_hashes(self):
        text = self.run_live("A1")
        self.assertIn('"max_file_store": 45Gi', text)
        self.assertIn('"server_name": $SERVER_NAME', text)
        conf = ni.load_natsconf(text)
        users = {u["user"]: u for v in conf["accounts"].values() for u in v["users"]}
        self.assertEqual(users["onemind"]["password"], "live-hash-onemind")
        self.assertEqual(users["sys"]["password"], "live-hash-sys")
        self.assertEqual(users["policy"]["password"], "hashed")
        self.assertEqual(conf["leafnodes"]["authorization"]["users"][0]["password"], "live-hash-leaf")
        self.assertEqual(conf["http_port"], 8222)

    def test_live_a2_removes_leaf_authorization(self):
        conf = ni.load_natsconf(self.run_live("A2"))
        self.assertNotIn("authorization", conf["leafnodes"])
        users = {u["user"]: u for u in conf["accounts"]["ONEMIND"]["users"]}
        self.assertEqual(users["leaf_tak"]["password"], "live-hash-leaf")
        self.assertEqual(conf["leafnodes"]["port"], 7422)

    def test_refuses_credentials_inside_the_repo(self):
        with self.assertRaises(SystemExit):
            ni.write_private(str(ROOT / "platform/nats/x.json"), "{}")
        with self.assertRaises(SystemExit):
            ni.read_passwords(str(ROOT / "platform/nats/identities/identities.yaml"))



class RoleTests(unittest.TestCase):
    """check-roles: the #565 role file against the CLIENT account."""

    FIXTURE = ROOT / "platform/nats/identities/permtest/testdata/client-roles.json"

    def setUp(self):
        self.spec = ni.load_spec()
        self.roles = json.loads(self.FIXTURE.read_text())

    def role(self, roles, name):
        return next(r for r in roles["roles"] if r["name"] == name)

    def test_fixture_fits(self):
        self.assertEqual(ni.check_roles(self.spec, self.roles), [])

    def test_no_role_commands_or_forges(self):
        for subj in ("cmd.fleet.>", "ent.*.*.state", "tak.cot.inject", "cmd.>", "evt.so.action.>", "$JS.API.>"):
            roles = copy.deepcopy(self.roles)
            self.role(roles, "admin")["publish"].append(subj)
            self.assertTrue(any("may publish" in e for e in ni.check_roles(self.spec, roles)), subj)

    def test_newcomer_and_viewer_never_confirm(self):
        for name in ("newcomer", "viewer"):
            for subj in ("svc.policy.command.confirm", "svc.policy.>", "svc.policy.*.confirm"):
                roles = copy.deepcopy(self.roles)
                self.role(roles, name)["publish"].append(subj)
                self.assertTrue(any("never confirm" in e for e in ni.check_roles(self.spec, roles)), (name, subj))
        self.assertIn("svc.policy.command.confirm", self.role(self.roles, "operator")["publish"])

    def test_subjects_need_a_client_import(self):
        roles = copy.deepcopy(self.roles)
        self.role(roles, "viewer")["proposed_subscribe"].append("tel.fleet.>")
        self.role(roles, "viewer")["proposed_publish"].append("svc.mission.apply")
        errs = ni.check_roles(self.spec, roles)
        self.assertTrue(any("tel.fleet.> is not imported" in e for e in errs))
        self.assertTrue(any("svc.mission.apply has no CLIENT service import" in e for e in errs))

    def test_within(self):
        self.assertTrue(ni._within("evt.map.geofence.*.*", ["evt.map.geofence.>"]))
        self.assertTrue(ni._within("svc.tak.mission.>", ["svc.tak.mission.>"]))
        self.assertFalse(ni._within("tel.feed.>", ["tel.feed.*.status"]))
        self.assertFalse(ni._within("svc.layer.>", ["svc.layer.list"]))
        self.assertTrue(ni._within("tel.feed.nws.status", ["tel.feed.*.status"]))


if __name__ == "__main__":
    unittest.main()
