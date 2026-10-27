"""Register new reviewed adapters here; catalog configuration chooses instances."""
from . import nats, prometheus

ADAPTERS = {"nats": nats.register, "prometheus": prometheus.register}
# Adapters that accept a list of upstreams, one per node (`<url_env>S`); others get one.
MULTI_NODE = {"nats"}
