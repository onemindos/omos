"""Reviewed integrations; adding an app instance is catalog-only."""
from . import nats, prometheus, applications, tile38

ADAPTERS = {"nats": nats.register, "prometheus": prometheus.register, "tile38": tile38.tile38,
            **{name: getattr(applications, name) for name in (
                "stac", "ollama", "alertmanager", "loki", "tempo", "opencost",
                "mediamtx", "karakeep", "n8n", "paperclip", "grafana", "litellm", "argocd", "backstage")}}
# Adapters that accept a list of upstreams, one per node (`<url_env>S`); others get one.
MULTI_NODE = {"nats"}
