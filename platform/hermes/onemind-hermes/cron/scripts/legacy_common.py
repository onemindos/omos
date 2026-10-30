"""Shared helpers for the Legacy cron scripts (platform/hermes/onemind-hermes/cron/).

Standard library only. Runs inside the onemind-hermes pod under Hermes' cron runner, which
passes ONEMIND_LITELLM_MCP_KEY (verified: not scrubbed) and has no Kubernetes token.

The pod's NetworkPolicy only reaches LiteLLM in onemind-ai on TCP 4000, so every cluster
fact comes from the read-only onemind MCP tools behind LiteLLM (Prometheus, Alertmanager,
Argo CD, NATS, OpenCost). Nothing here writes anything.
"""

import json
import os
import sys
import urllib.error
import urllib.request

MCP_URL = os.environ.get("ONEMIND_MCP_URL", "http://litellm.onemind-ai.svc.cluster.local:4000/mcp/")


class McpError(RuntimeError):
    pass


def _parse_body(raw):
    """Streamable-HTTP replies are either plain JSON or one SSE `data:` event."""
    raw = raw.strip()
    if raw.startswith("{"):
        return json.loads(raw)
    for line in raw.splitlines():
        if line.startswith("data:"):
            return json.loads(line[5:].strip())
    raise McpError("unparseable MCP reply")


def call(tool, arguments=None, timeout=25):
    """Call one onemind MCP tool through LiteLLM; return the list of text content items."""
    key = os.environ.get("ONEMIND_LITELLM_MCP_KEY", "")
    if not key:
        raise McpError("ONEMIND_LITELLM_MCP_KEY is not set")
    body = json.dumps({
        "jsonrpc": "2.0", "id": 1, "method": "tools/call",
        "params": {"name": tool, "arguments": arguments or {}},
    }).encode()
    req = urllib.request.Request(MCP_URL, data=body, method="POST", headers={
        "x-litellm-api-key": "Bearer " + key,
        "Content-Type": "application/json",
        "Accept": "application/json, text/event-stream",
    })
    last = None
    for _ in range(2):  # one retry for a transient gateway blip
        try:
            with urllib.request.urlopen(req, timeout=timeout) as resp:
                reply = _parse_body(resp.read().decode("utf-8", "replace"))
            break
        except urllib.error.HTTPError as exc:
            last = "HTTP %s" % exc.code
        except (urllib.error.URLError, TimeoutError, OSError, ValueError) as exc:
            last = type(exc).__name__
    else:
        raise McpError("MCP gateway unreachable (%s)" % last)
    if "error" in reply:
        raise McpError("MCP error: %s" % str(reply["error"].get("message", "unknown"))[:120])
    result = reply.get("result", {})
    texts = [c.get("text", "") for c in result.get("content", []) if c.get("type") == "text"]
    if result.get("isError"):
        raise McpError("tool %s failed: %s" % (tool, (texts[0] if texts else "")[:120]))
    return texts


def call_json(tool, arguments=None, timeout=25):
    """Like call(), with each text item decoded as JSON (one item -> that value)."""
    items = [json.loads(t) for t in call(tool, arguments, timeout) if t.strip()]
    return items[0] if len(items) == 1 else items


def prom(query):
    """Instant PromQL; returns [(labels_dict, float_value), ...]."""
    data = call_json("onemind-prometheus_query", {"query": query})
    if not isinstance(data, dict) or data.get("status") != "success":
        raise McpError("prometheus query failed")
    out = []
    for series in data.get("data", {}).get("result", []):
        try:
            out.append((series.get("metric", {}), float(series["value"][1])))
        except (KeyError, ValueError, IndexError):
            continue
    return out


def alerts():
    """Active Alertmanager alerts as dicts (Watchdog and InfoInhibitor dropped)."""
    data = call_json("onemind-alertmanager_alerts")
    if isinstance(data, dict):
        data = [data]
    out = []
    for a in data or []:
        name = alert_name(a.get("labels", {}))
        if name in ("Watchdog", "InfoInhibitor"):
            continue
        if a.get("status", {}).get("state", "active") != "active":
            continue
        out.append(a)
    return out


def alert_name(labels):
    """alertname, or `falco:<rule>` for Falco events (which carry a rule, not an alertname)."""
    if labels.get("alertname"):
        return labels["alertname"]
    if labels.get("rule"):
        return "falco:" + labels["rule"]
    return "unnamed-alert"


def alert_ns(labels):
    return labels.get("namespace") or labels.get("obj_namespace") or labels.get("k8s_ns_name") or ""


def argo_apps():
    data = call_json("onemind-argocd_applications", {"limit": 100})
    return data.get("applications", []) if isinstance(data, dict) else []


def safe(fn, default, problems=None, label=""):
    """Run fn(); on failure return default and note it (never raise out of a cron script)."""
    try:
        return fn()
    except Exception as exc:  # noqa: BLE001 - a broken probe must not hide the others
        if problems is not None:
            problems.append("%s: %s" % (label or "probe", str(exc)[:100]))
        return default


def ref(labels):
    ns = labels.get("namespace", "")
    name = labels.get("pod") or labels.get("deployment") or labels.get("persistentvolumeclaim") or ""
    return ("%s/%s" % (ns, name)).strip("/")


def emit(lines):
    sys.stdout.write("\n".join(lines) + "\n")
