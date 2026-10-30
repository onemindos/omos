#!/usr/bin/env python3
"""Render the OneMind plugin-ecosystem Backstage catalog.

Reads
  platform/backstage/catalog/plugins/inventory.yaml   reviewed plugin records (plugin.yaml contract fields)
  platform/hermes/mcp/apps.yaml                       one mcp-adapter Component per application entry
  services/onemind-portal/plugins/*/package.json      one backstage-plugin Component per workspace plugin
and writes the same plugins.yaml to both catalog copies:
  platform/backstage/catalog/plugins.yaml             mounted in the cluster (ConfigMap backstage-catalog-plugins)
  services/onemind-portal/catalog/plugins.yaml        development copy

  python scripts/generate_plugin_catalog.py           write both files
  python scripts/generate_plugin_catalog.py --check   fail when a file is stale or a reference does not resolve

The mapping from plugin.yaml fields to catalog fields is documented in
platform/backstage/catalog/plugins/README.md.
"""
from __future__ import annotations

import argparse
import json
import re
import sys
from pathlib import Path

import yaml

ROOT = Path(__file__).resolve().parents[1]
CATALOG = ROOT / "platform/backstage/catalog"
INVENTORY = CATALOG / "plugins/inventory.yaml"
MCP_APPS = ROOT / "platform/hermes/mcp/apps.yaml"
PORTAL_PLUGINS = ROOT / "services/onemind-portal/plugins"
OUTPUTS = [CATALOG / "plugins.yaml", ROOT / "services/onemind-portal/catalog/plugins.yaml"]
EXISTING = [CATALOG / "apps.yaml", CATALOG / "cluster.yaml"]

REPO = "https://github.com/onemindos/omos"
OWNER = "group:default/onemind-admins"
SYSTEM = "onemind-plugin-ecosystem"
VERIFIED_DATE = "2026-10-29"
LIFECYCLE = {"live": "production", "parked": "parked", "planned": "planned", "installed": "production"}
NAME_RE = re.compile(r"^[a-z0-9]([a-z0-9_.-]{0,61}[a-z0-9])?$")
LABEL_RE = re.compile(r"^[A-Za-z0-9]([-A-Za-z0-9_.]{0,61}[A-Za-z0-9])?$")
TAG_RE = re.compile(r"^[a-z0-9:+#]+(-[a-z0-9:+#]+)*$")
REF_RE = re.compile(r"^(component|resource|api|system):default/([a-z0-9][a-z0-9_.-]*)$")

# MCP application ids whose catalog Component has another name.
MCP_ALIASES = {
    "matrix": "k8s-matrix-synapse-c02538fc",
    "tak": "tak-server",
    "home_assistant": "home-assistant",
    "onemind_core": "k8s-omos-omos-tak-projector-ac850041",
}
MCP_HOST = "component:default/k8s-hermes-onemind-mcp-07887746"


def slug(text: str) -> str:
    return re.sub(r"-+", "-", re.sub(r"[^a-z0-9]+", "-", text.lower())).strip("-")


def csv(items) -> str:
    return ",".join(items)


def existing_names() -> dict[str, set[str]]:
    names: dict[str, set[str]] = {"component": set(), "resource": set(), "api": set(), "system": set()}
    for path in EXISTING:
        for doc in yaml.safe_load_all(path.read_text()):
            if doc and doc.get("kind", "").lower() in names:
                names[doc["kind"].lower()].add(doc["metadata"]["name"])
    return names


def component(name, title, description, ctype, status, origin, *, runs_kind, runs_on, connects=(), publishes=(),
              subscribes=(), subtype=None, version=None, source=None, build=None, verified=None,
              subcomponent_of=None, depends_on=(), provides=(), consumes=(), extra_tags=(), extra_ann=None,
              subjects=None):
    ann = {
        "onemind.io/plugin-type": ctype,
        "onemind.io/plugin-status": status,
        "onemind.io/plugin-origin": origin,
        "onemind.io/runs-on-kind": runs_kind,
        "onemind.io/runs-on": runs_on,
    }
    if subtype:
        ann["onemind.io/plugin-subtype"] = subtype
    if connects:
        ann["onemind.io/connects-to"] = csv(connects)
    if publishes:
        ann["onemind.io/nats-publishes"] = csv(subjects[s]["subject"] for s in publishes)
    if subscribes:
        ann["onemind.io/nats-subscribes"] = csv(subjects[s]["subject"] for s in subscribes)
    if version:
        ann["onemind.io/plugin-version"] = str(version)
    if source:
        ann["onemind.io/plugin-source"] = source
        ann["backstage.io/source-location"] = f"url:{REPO}/tree/main/{source}"
    else:
        ann["onemind.io/plugin-source"] = "not in this repository"
    if build:
        ann["onemind.io/plugin-build"] = build
    if verified:
        ann["onemind.io/plugin-verified"] = verified
    ann.update(extra_ann or {})
    tags = ["plugin", ctype, status, runs_kind]
    if origin == "third-party":
        tags.append("third-party")
    if publishes or subscribes:
        tags.append("nats")
    tags += list(extra_tags)
    spec = {"type": ctype, "lifecycle": LIFECYCLE[status], "owner": OWNER, "system": SYSTEM}
    if subcomponent_of:
        spec["subcomponentOf"] = subcomponent_of
    deps = list(depends_on) + [f"resource:default/nats-subject-{s}" for s in dict.fromkeys([*publishes, *subscribes])]
    if deps:
        spec["dependsOn"] = list(dict.fromkeys(deps))
    if provides:
        spec["providesApis"] = [f"api:default/{a}" for a in provides]
    if consumes:
        spec["consumesApis"] = [f"api:default/{a}" for a in consumes]
    return {
        "apiVersion": "backstage.io/v1alpha1",
        "kind": "Component",
        "metadata": {
            "name": name,
            "title": title,
            "description": description,
            "tags": list(dict.fromkeys(tags)),
            "labels": {
                "onemind.io/plugin-type": ctype,
                "onemind.io/plugin-status": status,
                "onemind.io/plugin-origin": origin,
                "onemind.io/runs-on-kind": runs_kind,
                "onemind.io/pillar": "so",
            },
            "annotations": ann,
        },
        "spec": spec,
    }


def resource(name, title, description, rtype, depends_on=(), annotations=None, tags=(), labels=None):
    doc = {
        "apiVersion": "backstage.io/v1alpha1",
        "kind": "Resource",
        "metadata": {
            "name": name,
            "title": title,
            "description": description,
            "tags": list(dict.fromkeys(["plugin-ecosystem", rtype, *tags])),
            "labels": {"onemind.io/pillar": "so", **(labels or {})},
            "annotations": annotations or {},
        },
        "spec": {"type": rtype, "owner": OWNER, "system": SYSTEM},
    }
    if depends_on:
        doc["spec"]["dependsOn"] = list(depends_on)
    return doc


def build_entities(inv: dict) -> list[dict]:
    subjects = {s["slug"]: s for s in inv["subjects"]}
    entities: list[dict] = []

    entities.append({
        "apiVersion": "backstage.io/v1alpha1",
        "kind": "System",
        "metadata": {
            "name": SYSTEM,
            "title": "OneMind Plugin Ecosystem",
            "description": (
                "Every plugin, adapter and bridge OneMind builds or runs, what it connects to and where it runs: "
                "TAK Server plugins, ATAK plugins, Hermes agent plugins, MCP adapters, omos-bridge protocol "
                "adapters and the Backstage portal plugins. Filter the catalog by Type (tak-server-plugin, "
                "atak-plugin, hermes-plugin, mcp-adapter, bridge, backstage-plugin), Lifecycle (production, "
                "parked, planned) or tag. Source of truth: platform/backstage/catalog/plugins/inventory.yaml."
            ),
            "tags": ["plugin-ecosystem", "plugins"],
            "labels": {"onemind.io/pillar": "so"},
            "annotations": {
                "onemind.io/pillar-review": "proposed",
                "backstage.io/source-location": f"url:{REPO}/tree/main/platform/backstage/catalog/plugins",
            },
            "links": [
                {"url": f"{REPO}/blob/main/platform/backstage/catalog/plugins/README.md", "title": "Catalog model and generator"},
            ],
        },
        "spec": {"owner": OWNER, "domain": "domain:default/onemind-so"},
    })

    for h in inv["hosts"]:
        entities.append(resource(
            h["name"], h["title"], h["description"], "plugin-host", depends_on=h["dependsOn"],
            tags=[h["kind"]],
            annotations={
                "onemind.io/runs-on-kind": h["kind"],
                "onemind.io/plugin-source": h["source"],
                "backstage.io/source-location": f"url:{REPO}/tree/main/{h['source']}",
                "onemind.io/plugin-verified": h["verified"],
            },
            labels={"onemind.io/runs-on-kind": h["kind"]},
        ))

    for s in inv["streams"]:
        entities.append(resource(
            f"nats-stream-{s['name'].lower()}", f"JetStream {s['name']}",
            f"JetStream stream {s['name']} capturing {s['subjects']}. Defined in platform/nats/streams/{s['file']} (NACK).",
            "nats-stream", depends_on=["component:default/nats"], tags=["nats"],
            annotations={
                "onemind.io/nats-subject": s["subjects"],
                "onemind.io/plugin-source": f"platform/nats/streams/{s['file']}",
                "backstage.io/source-location": f"url:{REPO}/blob/main/platform/nats/streams/{s['file']}",
            },
        ))

    for s in inv["subjects"]:
        entities.append(resource(
            f"nats-subject-{s['slug']}", s["subject"], s["description"], "nats-subject",
            depends_on=[f"resource:default/nats-stream-{x.lower()}" for x in s["streams"]] or ["component:default/nats"],
            tags=["nats", s["grammar"]],
            annotations={
                "onemind.io/nats-subject": s["subject"],
                "onemind.io/subject-grammar": s["grammar"],
            },
            labels={"onemind.io/subject-grammar": s["grammar"]},
        ))

    for p in inv["plugins"]:
        runs = p["runsOn"]
        host = runs.get("host")
        deps = list(p.get("dependsOn", []))
        if host:
            deps.append(host)
        entities.append(component(
            p["name"], p["title"], p["description"], p["type"], p["status"], p["origin"],
            runs_kind=runs["kind"], runs_on=runs["location"], connects=p.get("connectsTo", []),
            publishes=p.get("publishes", []), subscribes=p.get("subscribes", []), subtype=p.get("subtype"),
            version=p.get("version"), source=p.get("source"), build=p.get("build"), verified=p.get("verified"),
            subcomponent_of=p.get("subcomponentOf"), depends_on=deps, provides=p.get("providesApis", []),
            consumes=p.get("consumesApis", []), subjects=subjects,
        ))

    # ATAK plugins installed on the S24 (third party).
    defaults = inv["atak_connects"]
    overrides = inv.get("atak_overrides", {})
    seen: set[str] = set()
    for pkg, name, category, note in inv["atak_third_party"]:
        ov = overrides.get(pkg, {})
        cname = f"atak-plugin-{slug(name)}"
        assert cname not in seen, cname
        seen.add(cname)
        publishes = ov.get("publishes", [])
        subscribes = ov.get("subscribes", [])
        text = f"{note}." if not note.endswith(".") else note
        if ov.get("note"):
            text += " " + ov["note"]
        entities.append(component(
            cname, name, f"{text} Package {pkg}, installed in ATAK on the S24 Ultra; not built by OneMind.",
            "atak-plugin", "installed", "third-party",
            runs_kind="device", runs_on="Android EUD: ATAK-CIV 5.7.0.10 on the Samsung S24 Ultra",
            connects=ov.get("connectsTo", defaults[category]), publishes=publishes, subscribes=subscribes,
            subtype=category, verified="2026-07 (adb pm list packages; not re-verified)",
            depends_on=[*ov.get("dependsOn", []), "resource:default/host-atak-android-eud"],
            extra_tags=[category] + (["flagged"] if category == "suspicious" else []),
            extra_ann={"onemind.io/android-package": pkg}, subjects=subjects,
        ))

    # MCP adapters: one per application entry in platform/hermes/mcp/apps.yaml.
    names = existing_names()
    apps = yaml.safe_load(MCP_APPS.read_text())["apps"]
    for app in apps:
        enabled = app.get("enabled", True)
        has_adapter = bool(app.get("adapter"))
        status = "live" if enabled and has_adapter else "parked" if has_adapter else "planned"
        aid = app["id"]
        target = MCP_ALIASES.get(aid) or next((c for c in (aid, aid.replace("_", "-")) if c in names["component"]), None)
        extra = {"onemind.io/mcp-app-id": aid}
        if app.get("url_env"):
            extra["onemind.io/mcp-url-env"] = app["url_env"]
        if app.get("blocked_reason"):
            extra["onemind.io/blocked-reason"] = " ".join(app["blocked_reason"].split())
        deps = [f"component:default/{target}"] if target else []
        entities.append(component(
            f"mcp-adapter-{slug(aid)}", f"MCP adapter: {aid}",
            f"{app['description']}. " + (
                "Read-only adapter served by onemind-mcp and reached by Hermes through the LiteLLM MCP hub." if status == "live"
                else "Adapter implemented but disabled in apps.yaml until the blocker is cleared." if status == "parked"
                else "No adapter implemented yet; the entry reserves the application in apps.yaml."
            ),
            "mcp-adapter", status, "onemind", runs_kind="cluster",
            runs_on="ns hermes / deploy onemind-mcp (MCP server, streamable HTTP :8000/mcp)",
            connects=[aid.replace("_", "-"), "onemind-mcp", "litellm"],
            subtype=app.get("adapter"), source="platform/hermes/mcp/apps.yaml",
            build="services/onemind-mcp (Python MCP server image, onemind-mcp.yml); adapters in onemind_mcp/adapters/",
            verified=f"{VERIFIED_DATE} onemind-mcp pod Running; enabled flag from apps.yaml, upstream health not probed",
            subcomponent_of="component:default/onemind-mcp", depends_on=[MCP_HOST, *deps],
            extra_ann=extra, extra_tags=["mcp"], subjects=subjects,
        ))

    # Backstage portal plugins.
    for pkg in sorted(PORTAL_PLUGINS.glob("plugin-onemind-*/package.json")):
        data = json.loads(pkg.read_text())
        role = data.get("backstage", {}).get("role", "plugin")
        short = data["name"].split("/")[-1].removeprefix("plugin-")
        directory = pkg.parent.relative_to(ROOT).as_posix()
        entities.append(component(
            f"backstage-plugin-{short}", short, f"Backstage {role} {data['name']} in the OneMind portal workspace.",
            "backstage-plugin", "live", "onemind", runs_kind="cluster",
            runs_on="ns backstage / deploy backstage (compiled into the portal image)",
            connects=["backstage", "postgres-backstage"], subtype=role, version=data.get("version"),
            source=directory, build="services/onemind-portal image (build-portal.yml); digest pinned in platform/backstage/deployment.yaml",
            verified=f"{VERIFIED_DATE} workspace package; deployed image digest not matched to this commit",
            subcomponent_of="component:default/backstage", depends_on=["component:default/k8s-backstage-backstage-21b99354"],
            extra_tags=["backstage"], subjects=subjects,
        ))
    return entities


def render(entities: list[dict]) -> str:
    header = (
        "# GENERATED by scripts/generate_plugin_catalog.py from platform/backstage/catalog/plugins/inventory.yaml,\n"
        "# platform/hermes/mcp/apps.yaml and services/onemind-portal/plugins/*/package.json. Do not edit by hand.\n"
        "# Both copies (platform/backstage/catalog and services/onemind-portal/catalog) are identical.\n"
    )
    body = "---\n".join(
        yaml.safe_dump(e, sort_keys=False, allow_unicode=True, width=110) for e in entities
    )
    return header + "---\n" + body


def check(entities: list[dict], text: str) -> list[str]:
    errors: list[str] = []
    names = existing_names()
    own = {(e["kind"].lower(), e["metadata"]["name"]) for e in entities}
    if len(own) != len(entities):
        errors.append("duplicate entity name within the generated file")
    for kind, name in own:
        if name in names.get(kind, set()):
            errors.append(f"{kind}/{name} collides with an entity in apps.yaml or cluster.yaml")
    for e in entities:
        n = e["metadata"]["name"]
        if not NAME_RE.match(n):
            errors.append(f"invalid name {n}")
        for k, v in e["metadata"].get("labels", {}).items():
            if not LABEL_RE.match(str(v)):
                errors.append(f"{n}: label {k}={v!r} is not a valid label value")
        for t in e["metadata"].get("tags", []):
            if not TAG_RE.match(t) or len(t) > 63:
                errors.append(f"{n}: invalid tag {t}")
        spec = e["spec"]
        refs = [spec.get("subcomponentOf"), spec.get("system"), *spec.get("dependsOn", []),
                *spec.get("providesApis", []), *spec.get("consumesApis", [])]
        for ref in filter(None, refs):
            m = REF_RE.match(ref) or REF_RE.match(f"system:default/{ref}")
            if not m:
                errors.append(f"{n}: unparseable ref {ref}")
                continue
            kind, target = m.groups()
            if target not in names[kind] and (kind, target) not in own:
                errors.append(f"{n}: dangling ref {ref}")
        src = e["metadata"]["annotations"].get("onemind.io/plugin-source")
        if src and src != "not in this repository" and not (ROOT / src).exists():
            errors.append(f"{n}: plugin-source {src} does not exist")
    for out in OUTPUTS:
        if not out.exists() or out.read_text() != text:
            errors.append(f"{out.relative_to(ROOT)} is stale; run scripts/generate_plugin_catalog.py")
    return errors


def main() -> int:
    ap = argparse.ArgumentParser(description=__doc__.split("\n")[0])
    ap.add_argument("--check", action="store_true")
    args = ap.parse_args()
    entities = build_entities(yaml.safe_load(INVENTORY.read_text()))
    text = render(entities)
    if args.check:
        errors = check(entities, text)
        for err in errors:
            print(err, file=sys.stderr)
        print(f"{len(entities)} entities, {len(errors)} problems")
        return 1 if errors else 0
    for out in OUTPUTS:
        out.write_text(text)
    print(f"wrote {len(entities)} entities to {len(OUTPUTS)} files")
    return 0


if __name__ == "__main__":
    sys.exit(main())
