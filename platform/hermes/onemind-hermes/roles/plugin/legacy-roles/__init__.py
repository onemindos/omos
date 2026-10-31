"""Delegation backend `legacy-roles` (delegation.backend: legacy-roles).

Adds one optional per-task parameter to delegate_task, `worker`, an enum of Legacy's worker roles. A task
with a worker runs as that Hermes profile (its model, toolsets and brief; credentials come from the default
profile's pool). A task without one takes the `delegation` section's route, exactly like the bundled
`default` backend. The role list is the DEFAULT_WORKERS tuple, overridable with
delegation.backends.legacy-roles.workers in config.yaml. Contract: agent/delegation_provider.py.
"""

from __future__ import annotations

from typing import Any, Mapping

from agent.delegation_provider import DelegationProvider, DelegationRouteError

DEFAULT_WORKERS = ("ops-cluster", "tak-mission", "builder", "researcher", "checker")


class LegacyRolesProvider(DelegationProvider):
    @property
    def name(self) -> str:
        return "legacy-roles"

    @property
    def display_name(self) -> str:
        return "Legacy roles"

    @property
    def summary(self) -> str:
        return "delegate_task(worker=<role>) runs the child as that role's profile; no worker = the delegation section."

    def _workers(self, config: Mapping[str, Any]) -> list[str]:
        configured = self.settings(config).get("workers")
        if isinstance(configured, (list, tuple)) and all(isinstance(w, str) and w for w in configured) and configured:
            return list(configured)
        return list(DEFAULT_WORKERS)

    def task_parameters(self, config: Mapping[str, Any]) -> dict[str, dict[str, Any]]:
        return {
            "worker": {
                "type": "string",
                "enum": self._workers(config),
                "description": "Optional worker role for this task. Omit for a generic subagent.",
            }
        }

    def description(self, config: Mapping[str, Any]) -> str:
        return (
            "Worker roles (set `worker` per task): ops-cluster = cluster/Argo/K8s ops; tak-mission = TAK, missions, "
            "CoT, video; builder = code, manifests, PRs; researcher = web and docs research; checker = verifies "
            "other workers' output. A worker sees only the goal and context you give it, so make each task self-contained."
        )

    def route_task(self, task: Mapping[str, Any], config: Mapping[str, Any]) -> dict[str, Any]:
        worker = task.get("worker")
        if worker:
            if worker not in self._workers(config):
                raise DelegationRouteError(f"Unknown worker {worker!r}. Use one of: {', '.join(self._workers(config))}.")
            return {"profile": worker, "label": worker}
        return {k: v for k, v in config.items() if k not in ("service_tier", "backend", "backends")}


def register(ctx) -> None:
    ctx.register_delegation_provider(LegacyRolesProvider())
