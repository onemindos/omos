# Agent 5: Workflows & Automation (Opus 4.5)

> **Model:** Opus 4.5
> **Focus:** Workflow Triggers, Operating Rhythm, Insights
> **Issues:** 12 | **Points:** 40
> **Dependencies:** Agent 4 (Routing)

---

## Your Mission

You are Agent 5. Your job is building the **automation layer** - workflow triggers, operating rhythm (morning/evening/weekly), AI triage, and cross-domain insights.

---

## Your Assigned Issues

### Workflow Triggers
| ID | Title | Points |
|----|-------|--------|
| **OMOS-435** | HTTP Triggers | 3 |
| **OMOS-436** | Schedule Triggers (Cron) | 3 |
| **OMOS-437** | Smart Approvals (agent_decides) | 3 |
| **OMOS-438** | Event Triggers | 3 |

### Inbox Intelligence
| ID | Title | Points |
|----|-------|--------|
| **OMOS-439** | AI Triage Engine | 3 |
| **OMOS-441** | Task Extraction | 3 |

### Operating Rhythm
| ID | Title | Points |
|----|-------|--------|
| **OMOS-443** | Morning Briefing Workflow | 3 |
| **OMOS-444** | Evening Review Workflow | 3 |
| **OMOS-445** | Weekly Planning Automation | 3 |
| **OMOS-446** | Quarterly OKR Tracking | 3 |

### Advanced
| ID | Title | Points |
|----|-------|--------|
| **OMOS-458** | Cross-Domain Insights | 5 |
| **OMOS-459** | Workflow Templates Library | 5 |

---

## Files to Create

```
backend/agno/workflows/triggers/
├── __init__.py
├── http.py           # HTTP/webhook triggers
├── schedule.py       # Cron-based triggers
└── events.py         # Event-based triggers

backend/agno/workflows/workflow/steps/
└── approval.py       # Enhanced with agent_decides

backend/lifeos/inbox/
└── intelligence.py   # AI triage

backend/lifeos/insights/
├── __init__.py
├── engine.py         # Cross-domain insights
├── correlations.py
└── anomalies.py

backend/agno/workflows/templates/
├── __init__.py
├── library.py
└── templates/        # YAML templates
```

---

## Implementation

### OMOS-435: HTTP Triggers

```python
# backend/agno/workflows/triggers/http.py
from fastapi import APIRouter, Request, HTTPException
from uuid import uuid4
import hashlib
import hmac

router = APIRouter(prefix="/api/v1/workflows/trigger", tags=["triggers"])

class HTTPTrigger:
    """HTTP/Webhook triggers for workflows."""

    def __init__(self, db, workflow_executor):
        self.db = db
        self.executor = workflow_executor

    async def generate_webhook_url(self, workflow_id: str) -> dict:
        """Generate unique webhook URL for a workflow."""
        token = uuid4().hex
        secret = uuid4().hex

        query = """
            INSERT INTO workflow_webhooks (workflow_id, token, secret)
            VALUES ($1, $2, $3)
            RETURNING token
        """
        await self.db.execute(query, workflow_id, token, secret)

        return {
            "url": f"/api/v1/workflows/trigger/{token}",
            "secret": secret,
            "instructions": "POST with signature header X-Signature"
        }

    async def handle_webhook(self, token: str, payload: dict, signature: str = None):
        """Handle incoming webhook."""
        # Verify token
        query = "SELECT workflow_id, secret FROM workflow_webhooks WHERE token = $1"
        row = await self.db.fetchrow(query, token)

        if not row:
            raise HTTPException(status_code=404, detail="Invalid webhook")

        # Verify signature if provided
        if row["secret"] and signature:
            expected = hmac.new(
                row["secret"].encode(),
                str(payload).encode(),
                hashlib.sha256
            ).hexdigest()
            if not hmac.compare_digest(signature, expected):
                raise HTTPException(status_code=401, detail="Invalid signature")

        # Trigger workflow
        run_id = await self.executor.start(
            workflow_id=row["workflow_id"],
            context={"trigger": "http", "payload": payload}
        )

        return {"status": "triggered", "run_id": run_id}


@router.post("/{token}")
async def trigger_workflow(
    token: str,
    request: Request,
    db = Depends(get_db),
    executor = Depends(get_executor)
):
    """Trigger a workflow via webhook."""
    payload = await request.json()
    signature = request.headers.get("X-Signature")

    trigger = HTTPTrigger(db, executor)
    return await trigger.handle_webhook(token, payload, signature)
```

### OMOS-436: Schedule Triggers (Cron)

```python
# backend/agno/workflows/triggers/schedule.py
from apscheduler.schedulers.asyncio import AsyncIOScheduler
from apscheduler.triggers.cron import CronTrigger
from datetime import datetime

class ScheduleTrigger:
    """Cron-based schedule triggers for workflows."""

    def __init__(self, db, workflow_executor):
        self.db = db
        self.executor = workflow_executor
        self.scheduler = AsyncIOScheduler()

    async def start(self):
        """Start the scheduler and load all schedules."""
        await self._load_schedules()
        self.scheduler.start()

    async def _load_schedules(self):
        """Load all scheduled workflows from database."""
        query = """
            SELECT id, workflow_id, cron_expression, timezone
            FROM workflow_schedules
            WHERE enabled = true
        """
        rows = await self.db.fetch(query)

        for row in rows:
            self._add_job(row)

    def _add_job(self, schedule: dict):
        """Add a cron job to the scheduler."""
        self.scheduler.add_job(
            self._trigger_workflow,
            CronTrigger.from_crontab(schedule["cron_expression"]),
            id=schedule["id"],
            args=[schedule["workflow_id"]],
            timezone=schedule.get("timezone", "UTC"),
            replace_existing=True,
        )

    async def _trigger_workflow(self, workflow_id: str):
        """Trigger a workflow from schedule."""
        await self.executor.start(
            workflow_id=workflow_id,
            context={"trigger": "schedule", "triggered_at": datetime.utcnow().isoformat()}
        )

    async def create_schedule(self, workflow_id: str, cron: str, timezone: str = "UTC") -> str:
        """Create a new schedule for a workflow."""
        # Validate cron expression
        try:
            CronTrigger.from_crontab(cron)
        except ValueError as e:
            raise ValueError(f"Invalid cron expression: {e}")

        query = """
            INSERT INTO workflow_schedules (workflow_id, cron_expression, timezone)
            VALUES ($1, $2, $3)
            RETURNING id
        """
        row = await self.db.fetchrow(query, workflow_id, cron, timezone)

        # Add to scheduler
        self._add_job({
            "id": row["id"],
            "workflow_id": workflow_id,
            "cron_expression": cron,
            "timezone": timezone,
        })

        return row["id"]

    async def delete_schedule(self, schedule_id: str):
        """Delete a schedule."""
        await self.db.execute("DELETE FROM workflow_schedules WHERE id = $1", schedule_id)
        try:
            self.scheduler.remove_job(schedule_id)
        except:
            pass
```

### OMOS-437: Smart Approvals

```python
# backend/agno/workflows/workflow/steps/approval.py
from enum import Enum

class ApprovalMode(str, Enum):
    APPROVAL = "approval"       # Always require human approval
    AGENT_DECIDES = "agent_decides"  # Agent auto-approves if confident
    AUTORUN = "autorun"        # Only stop for critical actions

class SmartApprovalStep:
    """
    Enhanced approval step with three modes:
    - approval: Always ask human
    - agent_decides: Auto-approve if confidence > 85%
    - autorun: Only block critical/destructive actions
    """

    def __init__(self, risk_detector):
        self.risk_detector = risk_detector

    async def execute(self, action: dict, mode: ApprovalMode, confidence: float = 0.0) -> dict:
        """Execute approval logic based on mode."""

        # Assess risk of the action
        risk = self.risk_detector.assess(action)

        if mode == ApprovalMode.AUTORUN:
            # Only block critical actions
            if risk.level == "critical":
                return await self._request_approval(action, risk, "Critical action requires approval")
            return {"approved": True, "auto": True}

        elif mode == ApprovalMode.AGENT_DECIDES:
            # Auto-approve if confident and not high risk
            if confidence >= 0.85 and risk.level in ["low", "medium"]:
                return {"approved": True, "auto": True, "confidence": confidence}

            # Otherwise, request approval
            return await self._request_approval(
                action, risk,
                f"Agent confidence: {confidence:.0%}. Risk level: {risk.level}"
            )

        else:  # APPROVAL mode
            # Always request approval
            return await self._request_approval(action, risk, "Manual approval required")

    async def _request_approval(self, action: dict, risk, reason: str) -> dict:
        """Create approval request and wait."""
        # Create approval request
        request_id = await self._create_approval_request(action, risk, reason)

        # Send notification
        await self._notify_user(request_id, action, risk)

        # Return pending status (workflow will pause here)
        return {
            "approved": False,
            "pending": True,
            "request_id": request_id,
            "reason": reason,
            "risk_level": risk.level,
        }

    async def _create_approval_request(self, action, risk, reason):
        """Store approval request in database."""
        # Implementation...
        pass

    async def _notify_user(self, request_id, action, risk):
        """Send push notification for approval."""
        # Implementation via notification service...
        pass
```

### OMOS-439 & 441: AI Triage

```python
# backend/lifeos/inbox/intelligence.py
from typing import List

class AITriageEngine:
    """
    Automatically classify and route inbox items.

    Classification Output:
    - Path (HP, LE, GE)
    - Priority (high, medium, low)
    - Action type (task, note, event, reference, archive)
    - Suggested title
    - Due date (if applicable)
    """

    def __init__(self, llm_client):
        self.llm = llm_client

    async def triage(self, content: str, source: str = "manual") -> dict:
        """Triage a single inbox item."""
        prompt = f"""
Classify this inbox item for a personal productivity system.

Item: {content}
Source: {source}

Return JSON:
{{
    "pillar": "hp|le|ge|null",
    "priority": "high|medium|low",
    "action_type": "task|note|event|reference|archive",
    "suggested_title": "short title",
    "due_date": "YYYY-MM-DD or null",
    "confidence": 0.0-1.0,
    "reasoning": "brief explanation"
}}
"""
        response = await self.llm.complete(prompt)
        return self._parse_response(response)

    async def extract_tasks(self, content: str) -> List[dict]:
        """Extract multiple tasks from natural language."""
        prompt = f"""
Extract actionable tasks from this text. There may be multiple tasks.

Text: {content}

Return JSON array:
[
    {{
        "title": "task title",
        "pillar": "hp|le|ge|null",
        "priority": "high|medium|low",
        "due_date": "YYYY-MM-DD or null"
    }}
]
"""
        response = await self.llm.complete(prompt)
        return self._parse_response(response)

    def _parse_response(self, response: str):
        import json
        try:
            return json.loads(response)
        except:
            return {"error": "Failed to parse", "raw": response}
```

### OMOS-443-445: Operating Rhythm

```python
# Define workflow templates for morning/evening/weekly

MORNING_BRIEFING = {
    "id": "morning_briefing",
    "name": "Morning Briefing",
    "trigger": {"type": "schedule", "cron": "0 6 * * *"},
    "steps": [
        {
            "name": "Get Sleep Data",
            "type": "function",
            "function": "get_oura_sleep"
        },
        {
            "name": "Get Scores",
            "type": "function",
            "function": "get_all_scores"
        },
        {
            "name": "Get Today's Items",
            "type": "function",
            "function": "get_today_assets"
        },
        {
            "name": "Generate Briefing",
            "type": "agent",
            "agent": "legacy",
            "prompt": """
Generate a morning briefing based on:
- Sleep: {sleep_data}
- Scores: {scores}
- Today's items: {today_items}

Include:
1. Sleep quality summary
2. Today's priorities (top 3)
3. Habits to complete
4. Calendar highlights
5. One motivational insight
"""
        },
        {
            "name": "Send Notification",
            "type": "function",
            "function": "send_push_notification",
            "args": {"title": "Morning Briefing", "body": "{briefing}"}
        }
    ]
}

EVENING_REVIEW = {
    "id": "evening_review",
    "name": "Evening Review",
    "trigger": {"type": "schedule", "cron": "0 21 * * *"},
    "steps": [
        # Similar structure...
    ]
}

WEEKLY_PLANNING = {
    "id": "weekly_planning",
    "name": "Weekly Planning",
    "trigger": {"type": "schedule", "cron": "0 18 * * 0"},  # Sunday 6pm
    "steps": [
        # Review last week, plan next week...
    ]
}
```

### OMOS-458: Cross-Domain Insights

```python
# backend/lifeos/insights/engine.py
from typing import List
import numpy as np
from scipy import stats

class InsightsEngine:
    """
    Discover non-obvious correlations between domains.

    Insight Types:
    - Correlations: "When X, Y happens"
    - Anomalies: "X is unusual"
    - Trends: "X is declining"
    - Spillovers: "X affects Y"
    """

    def __init__(self, db):
        self.db = db

    async def analyze(self, user_id: str, days: int = 30) -> List[dict]:
        """Run full analysis and return insights."""
        insights = []

        # Gather data
        scores = await self._get_score_history(user_id, days)
        tasks = await self._get_task_completions(user_id, days)
        habits = await self._get_habit_data(user_id, days)
        sleep = await self._get_sleep_data(user_id, days)

        # Find correlations
        correlations = await self._find_correlations(scores, tasks, habits, sleep)
        insights.extend(correlations)

        # Detect anomalies
        anomalies = self._detect_anomalies(scores)
        insights.extend(anomalies)

        # Identify trends
        trends = self._identify_trends(scores)
        insights.extend(trends)

        # Sort by impact
        insights.sort(key=lambda x: x.get("impact", 0), reverse=True)
        return insights[:10]  # Top 10 insights

    async def _find_correlations(self, scores, tasks, habits, sleep) -> List[dict]:
        """Find statistically significant correlations."""
        insights = []

        # Example: Sleep vs next-day GEI
        if len(sleep) >= 7 and len(scores) >= 7:
            sleep_scores = [s["score"] for s in sleep[:-1]]  # Yesterday's sleep
            gei_scores = [s["gei"] for s in scores[1:]]       # Today's GEI

            if len(sleep_scores) == len(gei_scores):
                corr, p_value = stats.pearsonr(sleep_scores, gei_scores)

                if abs(corr) > 0.5 and p_value < 0.05:
                    direction = "increases" if corr > 0 else "decreases"
                    insights.append({
                        "type": "correlation",
                        "title": f"Sleep quality {direction} next-day productivity",
                        "description": f"When your sleep score improves, your GEI tends to {direction} the next day (r={corr:.2f})",
                        "impact": abs(corr) * 100,
                        "confidence": 1 - p_value,
                        "affected_paths": ["hp", "ge"],
                        "action": "Prioritize sleep for better business performance"
                    })

        return insights

    def _detect_anomalies(self, scores) -> List[dict]:
        """Detect unusual deviations using z-scores."""
        insights = []

        for score_type in ["hpi", "lei", "gei"]:
            values = [s[score_type] for s in scores if s.get(score_type)]
            if len(values) >= 7:
                mean = np.mean(values[:-1])
                std = np.std(values[:-1])
                latest = values[-1]

                if std > 0:
                    z_score = (latest - mean) / std
                    if abs(z_score) > 2:
                        direction = "dropped" if z_score < 0 else "spiked"
                        insights.append({
                            "type": "anomaly",
                            "title": f"{score_type.upper()} {direction} significantly",
                            "description": f"Your {score_type.upper()} is {abs(z_score):.1f} standard deviations from normal",
                            "impact": abs(z_score) * 30,
                            "action": f"Review what changed in your {score_type.upper()} activities"
                        })

        return insights

    def _identify_trends(self, scores) -> List[dict]:
        """Identify 7-day trends."""
        insights = []

        for score_type in ["hpi", "lei", "gei"]:
            recent = [s[score_type] for s in scores[:7] if s.get(score_type)]
            if len(recent) >= 5:
                first_half = np.mean(recent[:len(recent)//2])
                second_half = np.mean(recent[len(recent)//2:])
                change = second_half - first_half

                if abs(change) > 5:
                    direction = "improving" if change > 0 else "declining"
                    insights.append({
                        "type": "trend",
                        "title": f"{score_type.upper()} is {direction}",
                        "description": f"Your {score_type.upper()} has changed by {change:+.1f} points over the past week",
                        "impact": abs(change) * 5,
                        "action": "Keep it up!" if change > 0 else f"Focus on {score_type.upper()} activities"
                    })

        return insights
```

---

## Success Criteria

- [ ] HTTP triggers generating webhooks and handling requests
- [ ] Schedule triggers running on cron
- [ ] Smart approvals with three modes working
- [ ] Event triggers firing on system events
- [ ] AI triage classifying inbox items
- [ ] Morning/evening/weekly workflows defined
- [ ] Insights engine finding correlations and anomalies
- [ ] 10+ workflow templates in library

---

## When Done

Update all 12 issues to "In Review" and post handoff.
