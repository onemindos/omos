"""
Temporal Core Models
====================

Data models for the Temporal Core system - permanent time infrastructure.

This is not an agent - it's foundational infrastructure that all components use.
Time awareness is baked into the OS itself.
"""

from datetime import datetime, time, timedelta
from enum import Enum, IntEnum
from typing import Optional, List, Dict, Any
from pydantic import BaseModel, Field
from uuid import UUID, uuid4


# =============================================================================
# ENUMS
# =============================================================================


class TimeOfDay(str, Enum):
    """Time of day classification."""
    EARLY_MORNING = "early_morning"  # 5am-7am
    MORNING = "morning"              # 7am-12pm
    AFTERNOON = "afternoon"          # 12pm-5pm
    EVENING = "evening"              # 5pm-9pm
    NIGHT = "night"                  # 9pm-12am
    LATE_NIGHT = "late_night"        # 12am-5am


class DayType(str, Enum):
    """Type of day."""
    WORKDAY = "workday"
    WEEKEND = "weekend"
    HOLIDAY = "holiday"
    VACATION = "vacation"


class DeadlineUrgency(IntEnum):
    """Urgency levels for deadlines based on time remaining."""
    OVERDUE = 0       # Past due
    CRITICAL = 1      # < 1 hour
    URGENT = 2        # < 4 hours
    SOON = 3          # < 24 hours
    UPCOMING = 4      # < 1 week
    LATER = 5         # > 1 week


class PatternType(str, Enum):
    """Types of temporal patterns detected."""
    PRODUCTIVITY = "productivity"     # Work output patterns
    ENERGY = "energy"                 # Energy level patterns
    FOCUS = "focus"                   # Deep work patterns
    STRESS = "stress"                 # Stress patterns
    SLEEP = "sleep"                   # Sleep patterns
    ACTIVITY = "activity"             # General activity patterns
    DEADLINE = "deadline"             # Deadline behavior patterns
    COMMUNICATION = "communication"   # Response time patterns


class ScheduleEventType(str, Enum):
    """Types of scheduled events."""
    REMINDER = "reminder"
    TASK_DUE = "task_due"
    MEETING = "meeting"
    DEADLINE = "deadline"
    RECURRING = "recurring"
    TRIGGER = "trigger"      # Triggers an action/workflow
    DIGEST = "digest"        # Scheduled digest/summary


class ScheduleStatus(str, Enum):
    """Status of a scheduled event."""
    PENDING = "pending"
    TRIGGERED = "triggered"
    COMPLETED = "completed"
    FAILED = "failed"
    CANCELLED = "cancelled"
    SNOOZED = "snoozed"


class AggregationPeriod(str, Enum):
    """Time aggregation periods."""
    MINUTE = "minute"
    HOUR = "hour"
    DAY = "day"
    WEEK = "week"
    MONTH = "month"


# =============================================================================
# CORE TIME MODELS
# =============================================================================


class TimeRange(BaseModel):
    """A range of time."""
    start: datetime
    end: datetime

    @property
    def duration(self) -> timedelta:
        return self.end - self.start

    @property
    def duration_minutes(self) -> int:
        return int(self.duration.total_seconds() / 60)

    def contains(self, dt: datetime) -> bool:
        return self.start <= dt <= self.end

    def overlaps(self, other: "TimeRange") -> bool:
        return self.start < other.end and other.start < self.end


class TimeWindow(BaseModel):
    """A recurring time window (e.g., 9am-5pm on weekdays)."""
    start_hour: int = Field(ge=0, le=23)
    start_minute: int = Field(default=0, ge=0, le=59)
    end_hour: int = Field(ge=0, le=23)
    end_minute: int = Field(default=59, ge=0, le=59)
    days: Optional[List[int]] = None  # 0=Monday, 6=Sunday. None = all days
    timezone: str = "UTC"

    def is_active(self, dt: datetime) -> bool:
        """Check if the given time falls within this window."""
        if self.days is not None and dt.weekday() not in self.days:
            return False

        start = time(self.start_hour, self.start_minute)
        end = time(self.end_hour, self.end_minute)
        current = dt.time()

        if start <= end:
            return start <= current <= end
        else:
            # Overnight window (e.g., 22:00 - 06:00)
            return current >= start or current <= end


# =============================================================================
# CALENDAR MODELS
# =============================================================================


class CalendarEvent(BaseModel):
    """A calendar event."""
    id: str
    title: str
    start: datetime
    end: datetime

    # Optional fields
    description: Optional[str] = None
    location: Optional[str] = None
    attendees: List[str] = Field(default_factory=list)
    source: str = "calendar"  # google, apple, outlook, etc.
    source_id: Optional[str] = None

    # Metadata
    is_all_day: bool = False
    is_recurring: bool = False
    is_busy: bool = True
    color: Optional[str] = None

    @property
    def duration_minutes(self) -> int:
        return int((self.end - self.start).total_seconds() / 60)

    @property
    def is_happening_now(self) -> bool:
        now = datetime.utcnow()
        return self.start <= now <= self.end

    def minutes_until(self, from_time: Optional[datetime] = None) -> int:
        """Minutes until this event starts."""
        now = from_time or datetime.utcnow()
        if now >= self.start:
            return 0
        return int((self.start - now).total_seconds() / 60)


class Deadline(BaseModel):
    """A deadline from any source."""
    id: str
    title: str
    due: datetime
    source: str  # todoist, github, calendar, etc.
    source_id: Optional[str] = None

    # Calculated
    urgency: Optional[DeadlineUrgency] = None

    # Metadata
    priority: int = Field(default=1, ge=0, le=4)
    metadata: Dict[str, Any] = Field(default_factory=dict)

    def calculate_urgency(self, from_time: Optional[datetime] = None) -> DeadlineUrgency:
        """Calculate urgency based on time remaining."""
        now = from_time or datetime.utcnow()
        remaining = self.due - now

        if remaining.total_seconds() <= 0:
            return DeadlineUrgency.OVERDUE
        elif remaining <= timedelta(hours=1):
            return DeadlineUrgency.CRITICAL
        elif remaining <= timedelta(hours=4):
            return DeadlineUrgency.URGENT
        elif remaining <= timedelta(days=1):
            return DeadlineUrgency.SOON
        elif remaining <= timedelta(weeks=1):
            return DeadlineUrgency.UPCOMING
        else:
            return DeadlineUrgency.LATER


# =============================================================================
# PATTERN MODELS
# =============================================================================


class TemporalPattern(BaseModel):
    """A detected temporal pattern."""
    id: UUID = Field(default_factory=uuid4)
    type: PatternType
    name: str
    description: str

    # Pattern data
    hour_distribution: Dict[int, float] = Field(default_factory=dict)  # Hour -> score
    day_distribution: Dict[int, float] = Field(default_factory=dict)   # Day of week -> score

    # Statistics
    confidence: float = Field(default=0.0, ge=0.0, le=1.0)
    sample_size: int = 0
    first_observed: datetime
    last_observed: datetime

    # Insights
    peak_hours: List[int] = Field(default_factory=list)
    peak_days: List[int] = Field(default_factory=list)
    insights: List[str] = Field(default_factory=list)


class ProductivityPattern(TemporalPattern):
    """Productivity-specific pattern with additional fields."""
    type: PatternType = PatternType.PRODUCTIVITY

    # Productivity specific
    avg_tasks_per_hour: Dict[int, float] = Field(default_factory=dict)
    completion_rate_by_hour: Dict[int, float] = Field(default_factory=dict)
    focus_score_by_hour: Dict[int, float] = Field(default_factory=dict)

    # Recommendations
    optimal_deep_work_hours: List[int] = Field(default_factory=list)
    optimal_meeting_hours: List[int] = Field(default_factory=list)
    recommended_break_times: List[int] = Field(default_factory=list)


# =============================================================================
# SCHEDULED EVENT MODELS
# =============================================================================


class ScheduledEvent(BaseModel):
    """A scheduled future event/action."""
    id: UUID = Field(default_factory=uuid4)
    type: ScheduleEventType
    status: ScheduleStatus = ScheduleStatus.PENDING

    # Timing
    scheduled_at: datetime
    timezone: str = "UTC"

    # Recurrence
    is_recurring: bool = False
    recurrence_rule: Optional[str] = None  # RRULE format
    next_occurrence: Optional[datetime] = None

    # Payload
    title: str
    description: Optional[str] = None
    payload: Dict[str, Any] = Field(default_factory=dict)

    # Action to trigger
    action_type: Optional[str] = None  # workflow, agent, tool, notification
    action_target: Optional[str] = None  # workflow_id, agent_id, tool_name
    action_params: Dict[str, Any] = Field(default_factory=dict)

    # Context
    user_id: str = "zeus"
    source: Optional[str] = None
    source_id: Optional[str] = None

    # State
    triggered_at: Optional[datetime] = None
    completed_at: Optional[datetime] = None
    error: Optional[str] = None

    # Metadata
    metadata: Dict[str, Any] = Field(default_factory=dict)
    created_at: datetime = Field(default_factory=datetime.utcnow)
    updated_at: datetime = Field(default_factory=datetime.utcnow)


class ScheduledEventCreate(BaseModel):
    """Create a scheduled event."""
    type: ScheduleEventType
    scheduled_at: datetime
    timezone: str = "UTC"

    title: str
    description: Optional[str] = None
    payload: Dict[str, Any] = Field(default_factory=dict)

    # Recurrence
    is_recurring: bool = False
    recurrence_rule: Optional[str] = None

    # Action
    action_type: Optional[str] = None
    action_target: Optional[str] = None
    action_params: Dict[str, Any] = Field(default_factory=dict)

    # Context
    user_id: str = "zeus"
    source: Optional[str] = None
    source_id: Optional[str] = None


# =============================================================================
# TEMPORAL CONTEXT MODEL
# =============================================================================


class TemporalContext(BaseModel):
    """
    The complete temporal context - injected into every agent run.

    This is what makes time "permanent" in the OS - every component
    automatically knows the time context.
    """
    # Current moment
    now: datetime
    timezone: str
    time_of_day: TimeOfDay
    day_type: DayType

    # Calendar awareness
    current_event: Optional[CalendarEvent] = None
    next_event: Optional[CalendarEvent] = None
    minutes_until_next: Optional[int] = None
    today_events: List[CalendarEvent] = Field(default_factory=list)
    events_count_today: int = 0

    # Deadlines
    urgent_deadlines: List[Deadline] = Field(default_factory=list)
    deadline_pressure: float = Field(default=0.0, ge=0.0, le=1.0)
    overdue_count: int = 0
    due_today_count: int = 0

    # Historical patterns (from TimescaleDB)
    productivity_score: float = Field(default=0.5, ge=0.0, le=1.0)
    typical_load: str = "normal"  # lighter, normal, heavier
    energy_prediction: float = Field(default=0.5, ge=0.0, le=1.0)

    # Focus windows
    is_focus_time: bool = False
    focus_windows_today: List[TimeRange] = Field(default_factory=list)
    next_focus_window: Optional[TimeRange] = None

    # Activity context
    events_last_hour: int = 0
    events_last_24h: int = 0
    active_since: Optional[datetime] = None

    # Anomalies
    is_anomalous: bool = False
    anomaly_description: Optional[str] = None

    def to_prompt_context(self) -> str:
        """Generate a string suitable for system prompts."""
        lines = [
            f"Current time: {self.now.strftime('%Y-%m-%d %H:%M')} ({self.timezone})",
            f"Time of day: {self.time_of_day.value} on a {self.day_type.value}",
        ]

        if self.current_event:
            lines.append(f"Currently in: {self.current_event.title}")

        if self.next_event and self.minutes_until_next:
            lines.append(f"Next event: {self.next_event.title} in {self.minutes_until_next}m")

        if self.urgent_deadlines:
            urgent = self.urgent_deadlines[0]
            lines.append(f"Urgent deadline: {urgent.title} ({urgent.urgency.name})")

        if self.deadline_pressure > 0.7:
            lines.append(f"Deadline pressure: HIGH ({self.deadline_pressure:.0%})")

        if self.is_focus_time:
            lines.append("Currently in a FOCUS WINDOW - minimize interruptions")

        if self.productivity_score > 0.7:
            lines.append(f"Productivity typically PEAKS at this time")
        elif self.productivity_score < 0.3:
            lines.append(f"Productivity typically DIPS at this time")

        return "\n".join(lines)


# =============================================================================
# METRICS MODELS (for TimescaleDB)
# =============================================================================


class EventMetric(BaseModel):
    """A single event metric for time-series storage."""
    time: datetime
    event_type: str
    source: str
    count: int = 1
    value: Optional[float] = None
    metadata: Dict[str, Any] = Field(default_factory=dict)


class AggregatedMetric(BaseModel):
    """Aggregated metrics over a time period."""
    period_start: datetime
    period_end: datetime
    period_type: AggregationPeriod

    # Counts
    total_events: int
    events_by_type: Dict[str, int] = Field(default_factory=dict)
    events_by_source: Dict[str, int] = Field(default_factory=dict)

    # Productivity
    tasks_completed: int = 0
    tasks_created: int = 0
    completion_rate: float = 0.0

    # Activity
    active_minutes: int = 0
    peak_activity_hour: Optional[int] = None

    # Biometric (if available)
    avg_heart_rate: Optional[float] = None
    avg_stress: Optional[float] = None
    avg_energy: Optional[float] = None


class TimeSeriesQuery(BaseModel):
    """Query parameters for time-series data."""
    start: datetime
    end: datetime
    period: AggregationPeriod = AggregationPeriod.HOUR

    # Filters
    event_types: Optional[List[str]] = None
    sources: Optional[List[str]] = None

    # Options
    include_raw: bool = False
    limit: int = 1000


class TimeSeriesResult(BaseModel):
    """Result of a time-series query."""
    query: TimeSeriesQuery
    data: List[AggregatedMetric]

    # Summary
    total_events: int
    period_count: int

    # Patterns detected
    patterns: List[TemporalPattern] = Field(default_factory=list)
