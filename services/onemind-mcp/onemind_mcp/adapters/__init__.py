"""Reviewed integrations; adding an app instance is catalog-only."""
from . import nats, prometheus, applications

ADAPTERS = {"nats": nats.register, "prometheus": prometheus.register,
            **{name: getattr(applications, name) for name in (
                "stac", "ollama", "alertmanager", "loki", "tempo", "opencost",
                "mediamtx", "karakeep", "n8n", "paperclip", "grafana")}}
