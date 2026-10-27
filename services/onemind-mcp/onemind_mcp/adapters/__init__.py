"""Register new reviewed adapters here; catalog configuration chooses instances."""
from . import nats, prometheus

ADAPTERS = {"nats": nats.register, "prometheus": prometheus.register}
