"""
Activity Feed Models
====================

Pydantic models for the unified activity/inbox system.
Handles all items that need attention: events, notifications, approvals, tasks, alerts.
"""

from datetime import datetime, time
from enum import Enum
from typing import Optional, List, Dict, Any
from pydantic import BaseModel, Field
from uuid import UUID


# =============================================================================
# ENUMS
# =============================================================================


class ActivityItemType(str, Enum):
    """Types of activity items."""
    MESSAGE = "message"           # Chat/conversation messages
    NOTIFICATION = "notification" # System notifications
    EVENT = "event"               # External webhook events
    APPROVAL = "approval"         # HITL approval requests
    TASK = "task"                 # Assigned work items
    ALERT = "alert"               # Urgent system alerts
    REMINDER = "reminder"         # Scheduled reminders
    SUGGESTION = "suggestion"     # AI suggestions


class Priority(str, Enum):
    """Priority levels for activity items."""
    CRITICAL = "critical"  # All channels + voice
    HIGH = "high"          # Push + In-App
    NORMAL = "normal"      # In-App
    LOW = "low"            # Silent, batched


class ActivityStatus(str, Enum):
    """Status of an activity item."""
    NEW = "new"            # Unread, unhandled
    READING = "reading"    # Read but not handled
    HANDLED = "handled"    # Completed/resolved
    ARCHIVED = "archived"  # Archived
    SNOOZED = "snoozed"    # Snoozed until later
    FILTERED = "filtered"  # Filtered out by awareness mode


class AssigneeType(str, Enum):
    """Who should handle this item."""
    HUMAN = "human"    # Only Zeus can handle
    AGENT = "agent"    # Legacy AI agent
    TEAM = "team"      # A team of agents
    AUTO = "auto"      # System decides


class ApprovalStatus(str, Enum):
    """Status of an approval request."""
    PENDING = "pending"
    APPROVED = "approved"
    DENIED = "denied"
    EXPIRED = "expired"
    CANCELLED = "cancelled"


class Pillar(str, Enum):
    """Three Paths of life management."""
    HP = "hp"  # Holistic Performance
    LE = "le"  # Legacy Evolution
    GE = "ge"  # Generational Entrepreneurship


# =============================================================================
# ACTIVITY ITEM MODELS
# =============================================================================


class ActivityItem(BaseModel):
    """An item in the unified activity feed."""
    id: UUID
    type: ActivityItemType
    priority: Priority = Priority.NORMAL
    status: ActivityStatus = ActivityStatus.NEW

    # Content
    title: str
    body: Optional[str] = None
    source: Optional[str] = None

    # Assignment
    assignee_type: AssigneeType = AssigneeType.AUTO
    assignee_id: Optional[str] = None
    assigned_at: Optional[datetime] = None
    assigned_by: Optional[str] = None

    # HITL (approval items)
    approval_id: Optional[UUID] = None
    approval_status: Optional[ApprovalStatus] = None

    # Task specific
    task_due_at: Optional[datetime] = None
    task_completed_at: Optional[datetime] = None
    task_pillar: Optional[Pillar] = None
    task_result: Optional[str] = None

    # Scheduling
    snoozed_until: Optional[datetime] = None

    # Metadata
    metadata: Dict[str, Any] = Field(default_factory=dict)
    action_url: Optional[str] = None

    # Related entities
    session_id: Optional[str] = None
    run_id: Optional[str] = None
    workflow_id: Optional[str] = None
    event_id: Optional[str] = None

    # User context
    user_id: str = "zeus"

    # State
    is_read: bool = False
    is_starred: bool = False
    read_at: Optional[datetime] = None

    # Awareness
    awareness_mode: Optional[str] = None

    # Timestamps
    created_at: datetime
    updated_at: datetime

    class Config:
        from_attributes = True


class ActivityItemCreate(BaseModel):
    """Create a new activity item."""
    type: ActivityItemType
    priority: Priority = Priority.NORMAL
    title: str
    body: Optional[str] = None
    source: Optional[str] = None

    assignee_type: AssigneeType = AssigneeType.AUTO
    assignee_id: Optional[str] = None

    # Approval specific
    approval_id: Optional[UUID] = None

    # Task specific
    task_due_at: Optional[datetime] = None
    task_pillar: Optional[Pillar] = None

    # Metadata
    metadata: Dict[str, Any] = Field(default_factory=dict)
    action_url: Optional[str] = None

    # Related entities
    session_id: Optional[str] = None
    run_id: Optional[str] = None
    workflow_id: Optional[str] = None
    event_id: Optional[str] = None

    # User
    user_id: str = "zeus"


class ActivityItemUpdate(BaseModel):
    """Update an activity item."""
    status: Optional[ActivityStatus] = None
    priority: Optional[Priority] = None
    is_read: Optional[bool] = None
    is_starred: Optional[bool] = None
    assignee_type: Optional[AssigneeType] = None
    assignee_id: Optional[str] = None
    snoozed_until: Optional[datetime] = None
    task_result: Optional[str] = None
    metadata: Optional[Dict[str, Any]] = None


# =============================================================================
# STATS & PREFERENCES
# =============================================================================


class ActivityStats(BaseModel):
    """Statistics about the activity feed."""
    total: int = 0
    unread: int = 0
    starred: int = 0
    critical: int = 0
    today_count: int = Field(default=0)
    by_type: Dict[str, int] = Field(default_factory=dict)
    by_priority: Dict[str, int] = Field(default_factory=dict)
    by_status: Dict[str, int] = Field(default_factory=dict)
    pending_approvals: int = 0
    pending_tasks: int = 0
    agent_tasks: int = 0

    class Config:
        populate_by_name = True


class ActivityPreferences(BaseModel):
    """User preferences for activity notifications."""
    user_id: str

    # Channels
    push_enabled: bool = True
    email_enabled: bool = False
    sms_enabled: bool = False

    # Type preferences
    type_preferences: Dict[str, Dict[str, bool]] = Field(default_factory=dict)

    # Source preferences
    source_preferences: Dict[str, Dict[str, Any]] = Field(default_factory=dict)

    # Digest
    digest_enabled: bool = True
    digest_frequency: str = "daily"
    digest_time: time = time(9, 0)

    # Quiet hours
    quiet_hours_enabled: bool = False
    quiet_hours_start: time = time(22, 0)
    quiet_hours_end: time = time(7, 0)

    created_at: datetime
    updated_at: datetime

    class Config:
        from_attributes = True


class ActivityPreferencesUpdate(BaseModel):
    """Update activity preferences."""
    push_enabled: Optional[bool] = None
    email_enabled: Optional[bool] = None
    sms_enabled: Optional[bool] = None
    type_preferences: Optional[Dict[str, Dict[str, bool]]] = None
    source_preferences: Optional[Dict[str, Dict[str, Any]]] = None
    digest_enabled: Optional[bool] = None
    digest_frequency: Optional[str] = None
    digest_time: Optional[time] = None
    quiet_hours_enabled: Optional[bool] = None
    quiet_hours_start: Optional[time] = None
    quiet_hours_end: Optional[time] = None


# =============================================================================
# HISTORY
# =============================================================================


class ActivityHistoryEntry(BaseModel):
    """History entry for activity item changes."""
    id: UUID
    activity_item_id: UUID
    action: str
    actor: str
    old_value: Optional[Dict[str, Any]] = None
    new_value: Optional[Dict[str, Any]] = None
    metadata: Dict[str, Any] = Field(default_factory=dict)
    created_at: datetime

    class Config:
        from_attributes = True


# =============================================================================
# LIST RESPONSE
# =============================================================================


class ActivityListResponse(BaseModel):
    """Response for listing activity items."""
    items: List[ActivityItem]
    total: int
    unread_count: int
    critical_count: int
    has_more: bool


# =============================================================================
# RUN TRACKING (for Activity tab)
# =============================================================================


class RunStatus(str, Enum):
    """Status of an agent/team/workflow run."""
    QUEUED = "queued"
    RUNNING = "running"
    PAUSED = "paused"        # Waiting for HITL
    COMPLETED = "completed"
    FAILED = "failed"
    CANCELLED = "cancelled"


class ExecutionType(str, Enum):
    """Type of execution."""
    AGENT = "agent"
    TEAM = "team"
    WORKFLOW = "workflow"


class ActiveRun(BaseModel):
    """An active agent/team/workflow run."""
    id: str
    execution_type: ExecutionType
    target_id: str
    target_name: str
    status: RunStatus

    # Task
    task: Optional[str] = None
    progress: float = 0.0
    current_step: Optional[str] = None

    # HITL
    is_paused: bool = False
    pending_approvals: List[str] = Field(default_factory=list)

    # Result
    result: Optional[str] = None
    error: Optional[str] = None

    # Context
    session_id: Optional[str] = None
    activity_item_id: Optional[UUID] = None
    user_id: str = "zeus"

    # Timestamps
    created_at: datetime
    started_at: Optional[datetime] = None
    ended_at: Optional[datetime] = None

    class Config:
        from_attributes = True
