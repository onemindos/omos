"""
AI Triage Engine
================

LLM-powered activity item classification using the Three Paths framework.

Classifies activity items into:
- Pillar (HP, LE, GE)
- Priority (critical, high, normal, low)
- Action type (do, delegate, defer, delete, reference)
- Suggested due date (if applicable)

Usage:
    from backend.consciousness.activity.triage import AITriageEngine

    triage_engine = AITriageEngine()

    # Triage a single item
    result = await triage_engine.triage(activity_item)

    # Batch triage multiple items
    results = await triage_engine.batch_triage(items)
"""

import json
import logging
from dataclasses import dataclass
from datetime import datetime, timezone as tz
from typing import Any, Dict, List, Optional

from .models import ActivityItem, Priority, Pillar

logger = logging.getLogger("consciousness.activity.triage")


# =============================================================================
# TRIAGE RESULT MODEL
# =============================================================================


@dataclass
class TriageResult:
    """Result of AI triage classification."""
    pillar: Optional[Pillar]
    priority: Priority
    action_type: str  # "do", "delegate", "defer", "delete", "reference"
    suggested_due: Optional[datetime]
    suggested_title: Optional[str]
    confidence: float
    reasoning: str


# =============================================================================
# TRIAGE PROMPT
# =============================================================================


TRIAGE_PROMPT = """
Classify this activity item using the Three Paths productivity framework.

**Item Details:**
Title: {title}
Body: {body}
Source: {source}
Current Time: {now}

**The Three Paths (Pillars):**
- **HP (Holistic Performance)**: Personal health, fitness, wellness, self-development, mindfulness, sleep, exercise, nutrition, mental health
- **LE (Legacy Evolution)**: Family, relationships, home, children, parenting, marriage, household management, community
- **GE (Generational Entrepreneurship)**: Business, career, finances, investments, wealth building, professional development, side projects

**Priority Levels:**
- **critical**: Life-threatening, time-critical within 2 hours, major financial impact, safety issues
- **high**: Due within 24 hours, important stakeholders, significant impact
- **normal**: Due within 7 days, standard work items
- **low**: Someday/maybe, nice-to-have, no deadline

**Action Types (GTD Framework):**
- **do**: Can be completed in under 2 minutes, actionable immediately
- **delegate**: Should be assigned to someone else or an agent
- **defer**: Schedule for later, requires planning
- **delete**: Not actionable or relevant, can be dismissed
- **reference**: Useful information to save, not actionable

Return JSON only:
{{
    "pillar": "hp|le|ge|null",
    "priority": "critical|high|normal|low",
    "action_type": "do|delegate|defer|delete|reference",
    "suggested_due": "YYYY-MM-DD or null",
    "suggested_title": "concise action-oriented title or null",
    "confidence": 0.0-1.0,
    "reasoning": "brief explanation of classification"
}}
"""


# =============================================================================
# AI TRIAGE ENGINE
# =============================================================================


class AITriageEngine:
    """
    LLM-powered activity item classification.

    Uses Claude Haiku for fast, cost-effective classification.
    """

    def __init__(self, model_id: str = "claude-haiku"):
        """
        Initialize the triage engine.

        Args:
            model_id: Model to use for classification
        """
        self.model_id = model_id
        self._model = None

    async def triage(self, item: ActivityItem) -> TriageResult:
        """
        Triage a single activity item.

        Args:
            item: ActivityItem to classify

        Returns:
            TriageResult with classification
        """
        logger.info(f"Triaging activity item: {item.id}")

        try:
            prompt = TRIAGE_PROMPT.format(
                title=item.title or "No title",
                body=(item.body or "No body")[:1000],  # Limit body length
                source=item.source or "Unknown",
                now=datetime.now(tz.utc).strftime("%Y-%m-%d %H:%M"),
            )

            response = await self._call_llm(prompt)
            result = self._parse_response(response)

            logger.info(
                f"Triage complete for {item.id}: "
                f"pillar={result.pillar}, priority={result.priority}, "
                f"confidence={result.confidence:.0%}"
            )

            return result

        except Exception as e:
            logger.error(f"Triage failed for {item.id}: {e}")
            # Return default classification on error
            return TriageResult(
                pillar=None,
                priority=Priority.NORMAL,
                action_type="defer",
                suggested_due=None,
                suggested_title=None,
                confidence=0.0,
                reasoning=f"Triage failed: {e}",
            )

    async def batch_triage(
        self,
        items: List[ActivityItem],
        max_concurrent: int = 5,
    ) -> Dict[str, TriageResult]:
        """
        Triage multiple items concurrently.

        Args:
            items: List of ActivityItems to classify
            max_concurrent: Maximum concurrent LLM calls

        Returns:
            Dict mapping item ID to TriageResult
        """
        import asyncio

        results = {}
        semaphore = asyncio.Semaphore(max_concurrent)

        async def triage_with_semaphore(item):
            async with semaphore:
                return item.id, await self.triage(item)

        tasks = [triage_with_semaphore(item) for item in items]
        completed = await asyncio.gather(*tasks, return_exceptions=True)

        for result in completed:
            if isinstance(result, Exception):
                logger.error(f"Batch triage error: {result}")
            else:
                item_id, triage_result = result
                results[str(item_id)] = triage_result

        return results

    async def suggest_priority(self, item: ActivityItem) -> Priority:
        """
        Quick priority suggestion without full triage.

        Uses heuristics first, falls back to LLM if uncertain.
        """
        # Quick heuristic checks
        title_lower = (item.title or "").lower()
        body_lower = (item.body or "").lower()
        combined = f"{title_lower} {body_lower}"

        # Critical keywords
        critical_keywords = ["urgent", "emergency", "asap", "critical", "immediately", "911"]
        if any(kw in combined for kw in critical_keywords):
            return Priority.CRITICAL

        # High priority keywords
        high_keywords = ["important", "deadline", "today", "tomorrow", "eod", "end of day"]
        if any(kw in combined for kw in high_keywords):
            return Priority.HIGH

        # Low priority keywords
        low_keywords = ["fyi", "when you can", "no rush", "someday", "maybe", "eventually"]
        if any(kw in combined for kw in low_keywords):
            return Priority.LOW

        # Default to normal
        return Priority.NORMAL

    async def suggest_pillar(self, item: ActivityItem) -> Optional[Pillar]:
        """
        Quick pillar suggestion using keyword matching.

        Uses heuristics first for common cases.
        """
        combined = f"{item.title or ''} {item.body or ''}".lower()

        # HP keywords
        hp_keywords = [
            "health", "fitness", "workout", "exercise", "gym", "sleep",
            "meditation", "wellness", "doctor", "medical", "nutrition",
            "diet", "weight", "running", "yoga", "mental health"
        ]
        if any(kw in combined for kw in hp_keywords):
            return Pillar.HP

        # LE keywords
        le_keywords = [
            "family", "kids", "children", "wife", "husband", "spouse",
            "home", "house", "repair", "school", "parent", "birthday",
            "holiday", "vacation", "relationship", "community"
        ]
        if any(kw in combined for kw in le_keywords):
            return Pillar.LE

        # GE keywords
        ge_keywords = [
            "business", "work", "project", "meeting", "client", "customer",
            "revenue", "profit", "investment", "finance", "money", "budget",
            "career", "job", "startup", "entrepreneur", "sales", "marketing"
        ]
        if any(kw in combined for kw in ge_keywords):
            return Pillar.GE

        # No clear match - would need LLM for accurate classification
        return None

    async def _call_llm(self, prompt: str) -> str:
        """Call the LLM for classification."""
        try:
            from backend.agno.models import get_model
            from agno.agent import Agent

            if self._model is None:
                self._model = get_model(self.model_id)

            agent = Agent(
                name="activity-triage",
                model=self._model,
                instructions="You classify activity items. Return only valid JSON.",
            )

            response = await agent.arun(prompt)

            if response and response.content:
                return response.content

            return "{}"

        except ImportError as e:
            logger.warning(f"Model not available: {e}")
            raise
        except Exception as e:
            logger.error(f"LLM call failed: {e}")
            raise

    def _parse_response(self, response: str) -> TriageResult:
        """Parse LLM response into TriageResult."""
        try:
            # Extract JSON from response
            content = response
            if "```json" in content:
                content = content.split("```json")[1].split("```")[0]
            elif "```" in content:
                content = content.split("```")[1].split("```")[0]

            data = json.loads(content.strip())

            # Parse pillar
            pillar = None
            if data.get("pillar"):
                pillar_str = data["pillar"].lower()
                if pillar_str in ("hp", "le", "ge"):
                    pillar = Pillar(pillar_str)

            # Parse priority
            priority_str = data.get("priority", "normal").lower()
            try:
                priority = Priority(priority_str)
            except ValueError:
                priority = Priority.NORMAL

            # Parse suggested due date
            suggested_due = None
            if data.get("suggested_due"):
                try:
                    suggested_due = datetime.strptime(
                        data["suggested_due"], "%Y-%m-%d"
                    ).replace(tzinfo=tz.utc)
                except (ValueError, TypeError):
                    pass

            return TriageResult(
                pillar=pillar,
                priority=priority,
                action_type=data.get("action_type", "defer"),
                suggested_due=suggested_due,
                suggested_title=data.get("suggested_title"),
                confidence=float(data.get("confidence", 0.5)),
                reasoning=data.get("reasoning", ""),
            )

        except json.JSONDecodeError as e:
            logger.warning(f"Failed to parse triage response: {e}")
            return TriageResult(
                pillar=None,
                priority=Priority.NORMAL,
                action_type="defer",
                suggested_due=None,
                suggested_title=None,
                confidence=0.0,
                reasoning=f"Parse error: {e}",
            )


# =============================================================================
# SINGLETON ACCESS
# =============================================================================


_triage_engine: Optional[AITriageEngine] = None


def get_triage_engine() -> AITriageEngine:
    """Get the global triage engine instance."""
    global _triage_engine
    if _triage_engine is None:
        _triage_engine = AITriageEngine()
    return _triage_engine


async def auto_triage_item(item: ActivityItem) -> TriageResult:
    """Convenience function to triage a single item."""
    engine = get_triage_engine()
    return await engine.triage(item)
