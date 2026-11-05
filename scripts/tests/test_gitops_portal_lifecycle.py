"""Offline F8 contract: the Add app and Retire app templates, their registration and their limits.

No GitHub, cluster or portal credentials. The behavioural tests (rendering the templates, the retire
file list for stac-browser, the hygiene fixtures) live in the portal workspace under
plugins/plugin-onemind-services-backend/src/lifecycle; this file guards the wiring around them.
"""
from pathlib import Path
import re
import unittest
import yaml

ROOT = Path(__file__).resolve().parents[2]
PORTAL = ROOT / "services/onemind-portal"
TEMPLATES = PORTAL / "templates"
ALLOWED_ACTIONS = {
    "onemind:app-add:prepare",
    "onemind:app-retire:plan",
    "fetch:template",
    "publish:github:pull-request",
}
MUTATING = re.compile(r"kubectl|argocd|doctl|digitalocean|helm |terraform|\bsync\b", re.I)


def template(name):
    return yaml.safe_load((TEMPLATES / name / "template.yaml").read_text())


class LifecycleTemplateTests(unittest.TestCase):
    def test_templates_are_registered_everywhere_the_others_are(self):
        for config, prefix in (
            (PORTAL / "app-config.yaml", "../../templates/"),
            (PORTAL / "app-config.production.yaml", "/app/templates/"),
            (ROOT / "platform/backstage/app-config.cluster.yaml", "/app/templates/"),
        ):
            targets = {
                loc["target"]
                for loc in yaml.safe_load(config.read_text())["catalog"]["locations"]
                if loc.get("type") == "file"
            }
            for directory in ("app-add", "app-retire", "recovery-preparation"):
                self.assertIn(f"{prefix}{directory}/template.yaml", targets, f"{config.name}: {directory}")

    def test_names_match_the_sidebar_route_and_the_flag_is_on_in_the_image(self):
        self.assertEqual(template("app-add")["metadata"]["name"], "add-app")
        self.assertEqual(template("app-retire")["metadata"]["name"], "retire-app")
        nav = (PORTAL / "plugins/plugin-onemind-nav/src/navModel.ts").read_text()
        self.assertIn("/create/templates/default/retire-app", nav)
        config = yaml.safe_load((PORTAL / "app-config.yaml").read_text())
        self.assertTrue(config["onemind"]["features"]["retireApp"])

    def test_steps_use_only_read_only_actions_and_open_drafts_with_the_app_credentials(self):
        for name in ("app-add", "app-retire"):
            tpl = template(name)
            text = (TEMPLATES / name / "template.yaml").read_text()
            self.assertNotIn("githubToken", text)
            self.assertNotIn("ui:field: Secret", text)
            for step in tpl["spec"]["steps"]:
                self.assertIn(step["action"], ALLOWED_ACTIONS, f"{name}/{step['id']}")
                self.assertIsNone(
                    MUTATING.search(yaml.safe_dump({"action": step["action"], "input": step.get("input")})),
                    f"{name}/{step['id']} names a mutating tool",
                )
            publish = tpl["spec"]["steps"][-1]
            self.assertEqual(publish["action"], "publish:github:pull-request")
            self.assertIs(publish["input"]["draft"], True)
            self.assertNotIn("token", publish["input"])
            self.assertEqual(publish["input"]["targetBranchName"], "main")

    def test_only_the_retire_template_deletes_files_and_only_through_the_pr(self):
        add = template("app-add")["spec"]["steps"][-1]["input"]
        retire = template("app-retire")["spec"]["steps"][-1]["input"]
        self.assertNotIn("filesToDelete", add)
        self.assertEqual(retire["filesToDelete"], "${{ steps.plan.output.filesToDelete }}")

    def test_the_kubernetes_skeleton_follows_the_new_service_checklist(self):
        skeleton = TEMPLATES / "app-add/skeleton"
        text = {p.name: p.read_text() for p in skeleton.rglob("*") if p.is_file()}
        self.assertIn("do-block-storage-retain", text["data-pvc.yaml"])
        backup = text["backup.yaml"]
        self.assertRegex(backup, r"trigger:\n    schedule: \$\{\{ values\.backupSchedule \}\}")
        self.assertRegex(backup, r"\n    storageClassName: do-block-storage\n")
        self.assertRegex(backup, r"\n    cacheStorageClassName: do-block-storage\n")
        deployment = text["deployment.yaml"]
        self.assertIn("resources:", deployment)
        self.assertIn("requests:", deployment)
        self.assertIn("limits:", deployment)
        self.assertIn("image: ${{ values.image }}", deployment)
        self.assertNotIn(":latest", "".join(text.values()))
        application = yaml.safe_load(text["platform-${{ values.name }}.yaml"])
        self.assertNotIn("finalizers", application["metadata"])
        self.assertNotIn("automated", application["spec"]["syncPolicy"])
        self.assertIn("whitelist-source-range", text["ingress.yaml"] + "".join(
            p.read_text() for p in (skeleton / "k8s-cloudflare").rglob("ingress.yaml")))
        self.assertIn("ingressClassName: tailscale", "".join(
            p.read_text() for p in (skeleton / "k8s-tailnet").rglob("ingress.yaml")))
        # The placeholder lists key names only, never a value.
        secret = next(v for k, v in text.items() if k.endswith(".secret.example"))
        self.assertNotRegex(secret, r"^[A-Z_]+=\S", "a value in the secret placeholder")

    def test_the_hard_coded_argo_application_count_is_consistent(self):
        structure = (PORTAL / "scripts/test-catalog-structure.cjs").read_text()
        coverage = (PORTAL / "scripts/test-coverage.cjs").read_text()
        a = re.search(r"inGit\.size, (\d+), 'platform/ should define the (\d+) source", structure)
        b = re.search(r"argoApplicationsLinked, (\d+)\)", coverage)
        c = re.search(r"All (\d+) Argo applications are linked", coverage)
        self.assertTrue(a and b and c)
        self.assertEqual({a.group(1), a.group(2), b.group(1), c.group(1)}, {a.group(1)})

    def test_no_new_rbac_and_the_observer_stays_read_only(self):
        docs = list(yaml.safe_load_all((ROOT / "platform/backstage/rbac.yaml").read_text()))
        roles = [d for d in docs if d and d["kind"] in ("ClusterRole", "Role")]
        for role in roles:
            for rule in role["rules"]:
                self.assertLessEqual(set(rule["verbs"]), {"get", "list", "watch"}, role["metadata"]["name"])
                self.assertNotIn("secrets", rule["resources"])

    def test_the_lifecycle_backend_adds_no_secret_and_no_permission(self):
        package = (PORTAL / "plugins/plugin-onemind-services-backend/package.json").read_text()
        self.assertNotRegex(package, r"octokit|@octokit|kubernetes-client|@kubernetes")
        source = "".join(
            p.read_text()
            for p in (PORTAL / "plugins/plugin-onemind-services-backend/src/lifecycle").glob("*.ts")
            if ".test." not in p.name and p.name != "testRepo.ts"
        )
        self.assertNotRegex(source, r"process\.env|githubToken|GITHUB_TOKEN")
        # Refresh reuses scaffolder.task.create, which the policy limits to administrators.
        self.assertIn("taskCreatePermission", source)
        self.assertNotRegex(source, r"createPermission\(")


if __name__ == "__main__":
    unittest.main()
