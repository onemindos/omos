"""Tests for scripts/hermes_identity_export.py.

Run with:  python -m unittest scripts/tests/test_hermes_identity_export.py
Needs PyYAML. The fake credentials are assembled at run time so this file never contains a literal
secret-shaped string (the repo's secret scan would rightly flag it).
"""
import importlib.util
import io
import json
import os
import shutil
import tarfile
import tempfile
import unittest

SCRIPT = os.path.join(os.path.dirname(__file__), "..", "hermes_identity_export.py")

try:
    import yaml  # noqa: F401
    HAVE_YAML = True
except ImportError:  # pragma: no cover
    HAVE_YAML = False


def load():
    spec = importlib.util.spec_from_file_location("hermes_identity_export", SCRIPT)
    mod = importlib.util.module_from_spec(spec)
    spec.loader.exec_module(mod)
    return mod


# Fake credentials, split so the source holds no secret-shaped literal.
NOTION = "nt" + "n_" + "A1b2C3d4E5f6G7h8"
GITHUB = "gh" + "p_" + "a" * 30
COMPOSIO = "c" + "k_" + "Zz9Yy8Xx7Ww6"
MATRIX = "m" + "at_" + "q" * 20
BEARER = "Bear" + "er " + "t" * 20


def write(path, text):
    os.makedirs(os.path.dirname(path), exist_ok=True)
    with open(path, "w", encoding="utf-8") as fh:
        fh.write(text)


def make_skill(root, rel, body="# skill\n", author=None, extra=None):
    front = "---\nname: %s\n" % os.path.basename(rel)
    if author:
        front += "author: %s\n" % author
    front += "---\n"
    write(os.path.join(root, rel, "SKILL.md"), front + body)
    for name, text in (extra or {}).items():
        write(os.path.join(root, rel, name), text)


@unittest.skipUnless(HAVE_YAML, "PyYAML not installed")
class ClassifyTests(unittest.TestCase):
    def setUp(self):
        self.tmp = tempfile.mkdtemp()
        self.data = os.path.join(self.tmp, "data")
        self.image = os.path.join(self.tmp, "image")
        self.mod = load()
        skills = os.path.join(self.data, "skills")
        make_skill(skills, "devops/agent-made")
        make_skill(skills, "devops/hub-installed")
        make_skill(skills, "x/upstream-same", body="same\n")
        make_skill(self.image, "x/upstream-same", body="same\n")
        make_skill(skills, "x/upstream-old", body="old\n")
        make_skill(self.image, "x/upstream-old", body="new\n")
        make_skill(skills, "x/upstream-patched", body="patched\n")
        make_skill(self.image, "x/upstream-patched", body="original\n")
        make_skill(skills, "y/ours", author="Legacy (Mother)")
        make_skill(skills, "y/hand-made")
        make_skill(skills, "z/community", author="Someone Else")
        make_skill(skills, ".archive/old-thing")
        write(os.path.join(skills, ".usage.json"), json.dumps({
            "agent-made": {"created_by": "agent"},
            "hub-installed": {"created_by": "installed"},
            "upstream-patched": {"patch_count": 2},
        }))
        write(os.path.join(skills, ".hub", "lock.json"), json.dumps({"installed": {"hub-installed": {}}}))

    def tearDown(self):
        shutil.rmtree(self.tmp)

    def test_rules(self):
        export, review, skipped = self.mod.classify_skills(self.data, [self.image])
        names = {n: why for n, _, why in export}
        self.assertEqual(set(names), {"agent-made", "upstream-patched", "ours", "hand-made"})
        self.assertTrue(names["agent-made"].startswith("R3"))
        self.assertTrue(names["upstream-patched"].startswith("R5"))
        self.assertTrue(names["ours"].startswith("R4"))
        self.assertTrue(names["hand-made"].startswith("R7"))
        self.assertEqual([n for n, _, _ in review], ["community"])
        self.assertEqual(sum(skipped.values()), 3)  # identical, older upstream, hub
        self.assertNotIn("old-thing", names)  # archived directories are never walked


@unittest.skipUnless(HAVE_YAML, "PyYAML not installed")
class ConfigTests(unittest.TestCase):
    def setUp(self):
        self.mod = load()

    def config(self):
        return {
            "model": {"provider": "digitalocean", "key_env": "DIGITALOCEAN_AI_KEY", "context_length": 131072},
            "secrets": {"onepassword": {"env": {"GITHUB_TOKEN": "op://Vault/Item/field"}}},
            "mcp": {"servers": {"notion": {"args": '["-y", "--env", "H={\\"Authorization\\": \\"%s\\"}"]' % BEARER}}},
            "mcp_servers": {
                "notion": {"command": "npx", "args": ["-y", "--args", "x", "--env",
                                                       'H={"Authorization": "%s"}' % BEARER]},
                "github": {"url": "https://example.test/mcp", "headers": {"Authorization": BEARER}},
                "composio": {"headers": {"x-consumer-api-key": COMPOSIO}},
                "ghl-pro": {"env": {"GHL_API_KEY": "plainvalue123456"}},
                "miro": {"auth": "oauth"},
            },
            "platforms": {"matrix": {"enabled": True, "user_id": "@someone:example.test", "token": MATRIX}},
            "dashboard": {"basic_auth": {"username": "admin", "password_hash": "x" * 60, "secret": "y" * 40}},
            "TELEGRAM_ALLOWED_USERS": 123456789,
            "compression": {"proactive_prune_tokens": 3000},
        }

    def test_redacts_and_preserves(self):
        red, reds, problems = self.mod.redact_config(self.config())
        self.assertEqual(problems, [])
        text = self.mod.dump_yaml(red)
        for secret in (NOTION, GITHUB, COMPOSIO, MATRIX, BEARER.split(" ")[1], "123456789", "@someone",
                       "plainvalue123456"):
            self.assertNotIn(secret, text)
        for kept in ("DIGITALOCEAN_AI_KEY", "op://Vault/Item/field", "oauth", "131072", "npx", "digitalocean", "3000"):
            self.assertIn(kept, text)
        self.assertIn("${PLATFORMS_MATRIX_TOKEN}", text)
        self.assertIn("${MCP_SERVERS_GITHUB_HEADERS_AUTHORIZATION}", text)
        self.assertIn("${TELEGRAM_ALLOWED_USERS}", text)

    def test_list_flags_survive_while_secret_goes(self):
        red, _, _ = self.mod.redact_config(self.config())
        args = red["mcp_servers"]["notion"]["args"]
        self.assertEqual(args[:2], ["-y", "--args"])
        self.assertNotIn(BEARER.split(" ")[1], " ".join(args))

    def test_unknown_sensitive_key_fails_closed(self):
        _, _, problems = self.mod.redact_config({"foo": {"client_secret": "abcdefghijklmnop"}})
        self.assertEqual(len(problems), 1)
        self.assertIn("foo.client_secret", problems[0])

    def test_safe_forms_are_not_problems(self):
        cfg = {"a": {"api_key": "SOME_ENV_VAR", "token": "op://v/i/f", "secret": "", "password": "${X_Y}"}}
        _, _, problems = self.mod.redact_config(cfg)
        self.assertEqual(problems, [])

    def test_plain_word_is_not_treated_as_safe(self):
        _, _, problems = self.mod.redact_config({"a": {"password": "hunter2hunter2"}})
        self.assertEqual(len(problems), 1)

    def test_idempotent(self):
        red, _, _ = self.mod.redact_config(self.config())
        again, _, problems = self.mod.redact_config(red)
        self.assertEqual(problems, [])
        self.assertEqual(red, again)


@unittest.skipUnless(HAVE_YAML, "PyYAML not installed")
class ScanTests(unittest.TestCase):
    def setUp(self):
        self.mod = load()

    def test_flags_real_looking_values(self):
        self.assertTrue(self.mod.scan_text("token = %s" % GITHUB))
        self.assertTrue(self.mod.scan_text("api-key: %s" % COMPOSIO))
        self.assertTrue(self.mod.scan_text("password = 'hunter2hunter2'"))

    def test_allows_placeholders_and_env_references(self):
        for line in ("password=NATS_PASSWORD", "token=${TOKEN_VALUE}", "password: <your password>",
                     "api-key: YOUR_API_KEY", "echo hello"):
            self.assertEqual(self.mod.scan_text(line), [], line)

    def test_reports_line_numbers_not_values(self):
        findings = self.mod.scan_text("ok\nok\ntoken = %s\n" % GITHUB)
        self.assertEqual(findings[0][0], 3)
        self.assertNotIn(GITHUB, repr(findings))


@unittest.skipUnless(HAVE_YAML, "PyYAML not installed")
class ExportTests(unittest.TestCase):
    def setUp(self):
        self.tmp = tempfile.mkdtemp()
        self.data = os.path.join(self.tmp, "data")
        self.mod = load()
        make_skill(os.path.join(self.data, "skills"), "devops/agent-made")
        write(os.path.join(self.data, "skills", ".usage.json"), json.dumps({"agent-made": {"created_by": "agent"}}))
        write(os.path.join(self.data, "SOUL.md"), "You are Legacy.\n")
        write(os.path.join(self.data, "cron", "jobs.json"), "[]\n")
        write(os.path.join(self.data, "config.yaml"), "model:\n  provider: digitalocean\n  key_env: DIGITALOCEAN_AI_KEY\n")
        write(os.path.join(self.data, "profiles", "so", "config.yaml"), "model:\n  provider: digitalocean\n")

    def tearDown(self):
        shutil.rmtree(self.tmp)

    def test_clean_export_layout(self):
        files, ok, report = self.mod.build_export(self.data, [])
        self.assertTrue(ok, report)
        self.assertEqual(set(files), {
            "skills/devops/agent-made/SKILL.md", "live/SOUL.md", "cron/jobs.json",
            "config/config.template.yaml", "profiles/so/config.template.yaml", "EXPORT-REPORT.md"})
        self.assertIn("Status: **OK**", report)

    def test_secret_in_a_skill_refuses_the_export_and_stays_out_of_the_report(self):
        write(os.path.join(self.data, "skills", "devops", "agent-made", "notes.md"), "token = %s\n" % GITHUB)
        files, ok, report = self.mod.build_export(self.data, [])
        self.assertFalse(ok)
        self.assertEqual(set(files), {"EXPORT-REPORT.md"})
        self.assertIn("FAILED", report)
        self.assertIn("notes.md:1", report)
        self.assertNotIn(GITHUB, report)

    def test_tar_stream_is_deterministic(self):
        files, _, _ = self.mod.build_export(self.data, [])
        a, b = io.BytesIO(), io.BytesIO()
        self.mod.write_tar(files, a)
        self.mod.write_tar(files, b)
        self.assertEqual(a.getvalue(), b.getvalue())
        with tarfile.open(fileobj=io.BytesIO(a.getvalue())) as tar:
            self.assertIn("EXPORT-REPORT.md", tar.getnames())


if __name__ == "__main__":
    unittest.main()
