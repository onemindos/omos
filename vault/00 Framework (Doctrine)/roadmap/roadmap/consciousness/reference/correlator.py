"""
Legacy Consciousness - Event Correlation Engine

Links related events across services to build complete stories.

Example:
- 10:00 GitHub issue #100 created
- 10:05 Todoist task created mentioning "#100"
- 10:30 PR #42 created fixing issue #100
- 11:00 PR merged
- 11:05 Todoist task completed

The correlator links these as ONE workflow.

Features:
- Entity extraction (issue #, PR #, task ID)
- Relationship graph building
- Timeline queries
- Cross-service correlation
"""

import asyncio
import json
import logging
import re
from dataclasses import dataclass, field
from datetime import datetime, timedelta
from typing import Any, Dict, List, Optional, Set

from backend.consciousness.core.bus import Event
from backend.consciousness.core.kv_store import KVStore, get_kv_store

logger = logging.getLogger("consciousness.correlator")


@dataclass
class Entity:
    """An entity extracted from an event."""
    type: str           # issue, pr, task, repo, user
    id: str             # The entity identifier
    service: str        # github, todoist, slack
    display: str        # Human-readable display

    def __hash__(self):
        return hash((self.type, self.id, self.service))

    def __eq__(self, other):
        if not isinstance(other, Entity):
            return False
        return (
            self.type == other.type and
            self.id == other.id and
            self.service == other.service
        )


@dataclass
class Correlation:
    """A correlation between events."""
    id: str
    entities: Set[Entity]
    events: List[Dict[str, Any]]  # Event summaries
    created_at: datetime
    updated_at: datetime
    status: str = "active"  # active, resolved, archived

    def add_event(self, event: Event, entities: Set[Entity]) -> None:
        """Add an event to this correlation."""
        self.events.append({
            "seq": event.sequence,
            "type": str(event.type),
            "source": event.source,
            "timestamp": event.timestamp,
            "summary": self._summarize_event(event),
        })
        self.entities.update(entities)
        self.updated_at = datetime.utcnow()

    def _summarize_event(self, event: Event) -> str:
        """Create a one-line summary of the event."""
        data = event.data
        event_type = str(event.type)

        if "github" in event_type:
            if "issue" in event_type:
                return f"Issue #{data.get('issue_number', '?')}: {data.get('title', 'unknown')}"
            elif "pr" in event_type:
                return f"PR #{data.get('pr_number', '?')}: {data.get('title', 'unknown')}"
            elif "commit" in event_type:
                return f"Commit to {data.get('repo', 'unknown')}"
        elif "todoist" in event_type:
            return f"Task: {data.get('title', data.get('content', 'unknown'))}"
        elif "action" in event_type:
            return f"Action: {data.get('action', 'unknown')}"

        return f"{event.source}: {event_type}"


class EntityExtractor:
    """Extracts entities from events."""

    # Patterns for entity extraction
    PATTERNS = {
        "github_issue": re.compile(r"(?:issue\s*)?#(\d+)", re.IGNORECASE),
        "github_pr": re.compile(r"(?:pr|pull\s*request)\s*#?(\d+)", re.IGNORECASE),
        "github_repo": re.compile(r"([a-zA-Z0-9_-]+/[a-zA-Z0-9_-]+)"),
        "todoist_task": re.compile(r"task[_\s]?(?:id)?[:\s]?(\d+)", re.IGNORECASE),
        "slack_channel": re.compile(r"<#([A-Z0-9]+)\|([^>]+)>"),
        "slack_user": re.compile(r"<@([A-Z0-9]+)>"),
    }

    def extract(self, event: Event) -> Set[Entity]:
        """Extract all entities from an event."""
        entities = set()

        # Get text to search
        text = self._get_searchable_text(event)

        # Extract based on source
        if event.source == "github" or "github" in str(event.type):
            entities.update(self._extract_github(event, text))
        elif event.source == "todoist" or "todoist" in str(event.type):
            entities.update(self._extract_todoist(event, text))
        elif event.source == "slack" or "slack" in str(event.type):
            entities.update(self._extract_slack(event, text))

        # Always look for cross-references
        entities.update(self._extract_references(text))

        return entities

    def _get_searchable_text(self, event: Event) -> str:
        """Get all text from event that might contain entities."""
        parts = []
        data = event.data

        # Common fields
        for key in ["title", "content", "description", "body", "message", "text"]:
            if key in data and isinstance(data[key], str):
                parts.append(data[key])

        # Nested fields
        if "params" in data and isinstance(data["params"], dict):
            for v in data["params"].values():
                if isinstance(v, str):
                    parts.append(v)

        return " ".join(parts)

    def _extract_github(self, event: Event, text: str) -> Set[Entity]:
        """Extract GitHub entities."""
        entities = set()
        data = event.data
        event_type = str(event.type)

        # Direct from event data
        if "issue_number" in data:
            entities.add(Entity(
                type="issue",
                id=str(data["issue_number"]),
                service="github",
                display=f"GitHub Issue #{data['issue_number']}",
            ))

        if "pr_number" in data:
            entities.add(Entity(
                type="pr",
                id=str(data["pr_number"]),
                service="github",
                display=f"GitHub PR #{data['pr_number']}",
            ))

        if "repo" in data or "repository" in data:
            repo = data.get("repo") or data.get("repository", {}).get("full_name", "")
            if repo:
                entities.add(Entity(
                    type="repo",
                    id=repo,
                    service="github",
                    display=f"Repo: {repo}",
                ))

        return entities

    def _extract_todoist(self, event: Event, text: str) -> Set[Entity]:
        """Extract Todoist entities."""
        entities = set()
        data = event.data

        # Task ID
        task_id = data.get("task_id") or data.get("id")
        if task_id:
            entities.add(Entity(
                type="task",
                id=str(task_id),
                service="todoist",
                display=f"Todoist Task #{task_id}",
            ))

        # Project
        project_id = data.get("project_id")
        if project_id:
            entities.add(Entity(
                type="project",
                id=str(project_id),
                service="todoist",
                display=f"Project #{project_id}",
            ))

        return entities

    def _extract_slack(self, event: Event, text: str) -> Set[Entity]:
        """Extract Slack entities."""
        entities = set()

        # Channel mentions
        for match in self.PATTERNS["slack_channel"].finditer(text):
            entities.add(Entity(
                type="channel",
                id=match.group(1),
                service="slack",
                display=f"#{match.group(2)}",
            ))

        # User mentions
        for match in self.PATTERNS["slack_user"].finditer(text):
            entities.add(Entity(
                type="user",
                id=match.group(1),
                service="slack",
                display=f"@{match.group(1)}",
            ))

        return entities

    def _extract_references(self, text: str) -> Set[Entity]:
        """Extract cross-service references from text."""
        entities = set()

        # GitHub issue/PR references in other services
        for match in self.PATTERNS["github_issue"].finditer(text):
            entities.add(Entity(
                type="issue_ref",
                id=match.group(1),
                service="github",
                display=f"Issue #{match.group(1)}",
            ))

        return entities


class CorrelationEngine:
    """
    Correlates events across services.

    Builds a relationship graph of related events and entities.
    """

    def __init__(self, kv_store: Optional[KVStore] = None):
        """Initialize correlation engine."""
        self._kv = kv_store or get_kv_store()
        self._extractor = EntityExtractor()

        # In-memory correlation cache (also persisted to KV)
        self._correlations: Dict[str, Correlation] = {}

        # Entity to correlation mapping
        self._entity_index: Dict[Entity, str] = {}

    async def process_event(self, event: Event) -> Optional[Correlation]:
        """
        Process an event and correlate it with related events.

        Returns the correlation if event was correlated, None otherwise.
        """
        # Extract entities
        entities = self._extractor.extract(event)

        if not entities:
            return None

        # Find existing correlation
        correlation = self._find_correlation(entities)

        if correlation:
            # Add to existing correlation
            correlation.add_event(event, entities)
            logger.debug(f"Added event to correlation {correlation.id}")
        else:
            # Create new correlation
            correlation = Correlation(
                id=f"corr_{event.sequence}_{datetime.utcnow().strftime('%H%M%S')}",
                entities=entities,
                events=[],
                created_at=datetime.utcnow(),
                updated_at=datetime.utcnow(),
            )
            correlation.add_event(event, entities)
            self._correlations[correlation.id] = correlation
            logger.info(f"Created new correlation: {correlation.id}")

        # Update entity index
        for entity in entities:
            self._entity_index[entity] = correlation.id

        # Persist to KV
        await self._persist_correlation(correlation)

        return correlation

    def _find_correlation(self, entities: Set[Entity]) -> Optional[Correlation]:
        """Find existing correlation that shares entities."""
        for entity in entities:
            if entity in self._entity_index:
                corr_id = self._entity_index[entity]
                if corr_id in self._correlations:
                    return self._correlations[corr_id]
        return None

    async def _persist_correlation(self, correlation: Correlation) -> None:
        """Persist correlation to KV store."""
        if not self._kv:
            return

        data = {
            "id": correlation.id,
            "entities": [
                {"type": e.type, "id": e.id, "service": e.service, "display": e.display}
                for e in correlation.entities
            ],
            "events": correlation.events,
            "created_at": correlation.created_at.isoformat(),
            "updated_at": correlation.updated_at.isoformat(),
            "status": correlation.status,
        }

        await self._kv.cache_set(f"correlation.{correlation.id}", data)

    async def get_correlation(self, corr_id: str) -> Optional[Correlation]:
        """Get a correlation by ID."""
        if corr_id in self._correlations:
            return self._correlations[corr_id]

        # Try loading from KV
        if self._kv:
            data = await self._kv.cache_get(f"correlation.{corr_id}")
            if data:
                return self._deserialize_correlation(data)

        return None

    async def get_related_events(self, entity: Entity) -> List[Dict]:
        """Get all events related to an entity."""
        if entity in self._entity_index:
            corr_id = self._entity_index[entity]
            correlation = await self.get_correlation(corr_id)
            if correlation:
                return correlation.events
        return []

    async def get_story(self, entity: Entity) -> Optional[str]:
        """
        Get the complete story for an entity.

        Returns a human-readable narrative of related events.
        """
        correlation = None
        if entity in self._entity_index:
            correlation = await self.get_correlation(self._entity_index[entity])

        if not correlation:
            return None

        # Build narrative
        lines = [f"Story for {entity.display}:", ""]

        for evt in sorted(correlation.events, key=lambda x: x["timestamp"]):
            timestamp = evt["timestamp"][:16].replace("T", " ")
            lines.append(f"  {timestamp} | {evt['summary']}")

        return "\n".join(lines)

    async def query(
        self,
        entity_type: Optional[str] = None,
        service: Optional[str] = None,
        since: Optional[datetime] = None,
    ) -> List[Correlation]:
        """
        Query correlations.

        Args:
            entity_type: Filter by entity type (issue, pr, task)
            service: Filter by service (github, todoist)
            since: Only correlations updated since this time
        """
        results = []

        for correlation in self._correlations.values():
            # Check time filter
            if since and correlation.updated_at < since:
                continue

            # Check entity filters
            if entity_type or service:
                match = False
                for entity in correlation.entities:
                    if entity_type and entity.type != entity_type:
                        continue
                    if service and entity.service != service:
                        continue
                    match = True
                    break
                if not match:
                    continue

            results.append(correlation)

        return results

    def _deserialize_correlation(self, data: Dict) -> Correlation:
        """Deserialize correlation from stored data."""
        entities = set()
        for e in data.get("entities", []):
            entities.add(Entity(
                type=e["type"],
                id=e["id"],
                service=e["service"],
                display=e["display"],
            ))

        return Correlation(
            id=data["id"],
            entities=entities,
            events=data.get("events", []),
            created_at=datetime.fromisoformat(data["created_at"]),
            updated_at=datetime.fromisoformat(data["updated_at"]),
            status=data.get("status", "active"),
        )


# Global instance
_correlator: Optional[CorrelationEngine] = None


def create_correlator(kv_store: Optional[KVStore] = None) -> CorrelationEngine:
    """Create correlation engine."""
    global _correlator
    _correlator = CorrelationEngine(kv_store)
    return _correlator


def get_correlator() -> Optional[CorrelationEngine]:
    """Get the global correlator instance."""
    return _correlator


__all__ = [
    "Entity",
    "Correlation",
    "EntityExtractor",
    "CorrelationEngine",
    "create_correlator",
    "get_correlator",
]
