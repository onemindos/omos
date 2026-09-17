"""
Universal Message Protocol
============================

The canonical message format for ALL communication in OneMind OS.
Every message — whether from voice, chat, vision, sensor, or robot —
flows through this protocol. It replaces ad-hoc dicts everywhere.

Design Principles:
    1. EVERY message has the same envelope (NexusMessage)
    2. Modality is explicit — text, audio, image, video, sensor, command
    3. Attachments carry binary/structured data alongside text
    4. Source identifies WHO sent it (human, agent, robot, system)
    5. Routing metadata tells WHERE it should go
    6. Everything serializes to JSON for NATS/WebSocket/REST

Usage:
    from backend.nexus.comms import NexusMessage, Modality, Attachment

    # Text chat from human
    msg = NexusMessage(
        content="Research Q1 revenue",
        modality=Modality.TEXT,
        source_id="human:zeus",
        source_type=SourceType.HUMAN,
    )

    # Voice command (transcribed + raw audio)
    msg = NexusMessage(
        content="Send the research team to analyze this",
        modality=Modality.VOICE,
        source_id="human:zeus",
        source_type=SourceType.HUMAN,
        attachments=[
            Attachment(type=AttachmentType.AUDIO, url="data:audio/wav;base64,...", metadata={"duration_ms": 3200})
        ],
    )

    # Image from camera/vision
    msg = NexusMessage(
        content="What's in this image?",
        modality=Modality.VISION,
        source_id="sensor:garage-cam",
        source_type=SourceType.SENSOR,
        attachments=[
            Attachment(type=AttachmentType.IMAGE, url="http://frigate/snapshot/garage", metadata={"camera": "garage"})
        ],
    )

    # Robot telemetry
    msg = NexusMessage(
        content="Position update",
        modality=Modality.SENSOR,
        source_id="robot:arm-01",
        source_type=SourceType.ROBOT,
        data={"lat": 33.45, "lon": -112.07, "battery": 82},
    )

NATS Subject Convention:
    comms.{modality}.{source_type}.{source_id}
    comms.text.human.zeus
    comms.voice.human.zeus
    comms.vision.sensor.garage-cam
    comms.sensor.robot.arm-01
    comms.command.agent.researcher
"""

import logging
import json
from enum import Enum
from typing import Dict, Any, List, Optional
from dataclasses import dataclass, field, asdict
from datetime import datetime
from uuid import uuid4

logger = logging.getLogger("nexus.comms")


# =============================================================================
# ENUMS
# =============================================================================

class Modality(str, Enum):
    """What kind of communication this is."""
    TEXT = "text"               # Plain text / chat
    VOICE = "voice"             # Voice input (transcribed text + optional audio)
    VISION = "vision"           # Image/video analysis
    AUDIO = "audio"             # Raw audio (not speech — music, ambient, alerts)
    VIDEO = "video"             # Video stream/clip
    SENSOR = "sensor"           # IoT/robot sensor data
    COMMAND = "command"         # Structured command (not natural language)
    MULTIMODAL = "multimodal"   # Multiple modalities combined


class SourceType(str, Enum):
    """Who/what sent this message."""
    HUMAN = "human"
    AGENT = "agent"
    TEAM = "team"
    ROBOT = "robot"
    DRONE = "drone"
    SENSOR = "sensor"
    SERVICE = "service"
    SYSTEM = "system"
    LEGACY = "legacy"          # Legacy super-observer


class AttachmentType(str, Enum):
    """Type of attachment data."""
    IMAGE = "image"            # JPEG, PNG, WebP
    AUDIO = "audio"            # WAV, MP3, OGG, PCM16
    VIDEO = "video"            # MP4, WebM
    DOCUMENT = "document"      # PDF, DOCX, TXT
    SENSOR_DATA = "sensor_data"  # Structured sensor payload
    FILE = "file"              # Generic file


class Priority(str, Enum):
    CRITICAL = "critical"
    HIGH = "high"
    NORMAL = "normal"
    LOW = "low"
    BACKGROUND = "background"


# =============================================================================
# ATTACHMENT — Binary/structured data alongside the message
# =============================================================================

@dataclass
class Attachment:
    """
    A piece of data attached to a message.

    Can be:
    - A URL (http://, data:, file://)
    - Base64-encoded inline data
    - Structured sensor payload in `metadata`
    """
    type: AttachmentType
    url: Optional[str] = None               # URL or data URI
    data_base64: Optional[str] = None       # Inline base64 data
    mime_type: Optional[str] = None         # e.g. "image/jpeg", "audio/wav"
    filename: Optional[str] = None
    size_bytes: Optional[int] = None
    metadata: Dict[str, Any] = field(default_factory=dict)

    def to_dict(self) -> Dict[str, Any]:
        d = asdict(self)
        d["type"] = self.type.value
        return {k: v for k, v in d.items() if v is not None}

    @staticmethod
    def from_dict(data: Dict[str, Any]) -> "Attachment":
        return Attachment(
            type=AttachmentType(data.get("type", "file")),
            url=data.get("url"),
            data_base64=data.get("data_base64"),
            mime_type=data.get("mime_type"),
            filename=data.get("filename"),
            size_bytes=data.get("size_bytes"),
            metadata=data.get("metadata", {}),
        )


# =============================================================================
# ROUTING — Where this message should go
# =============================================================================

@dataclass
class Routing:
    """
    Optional routing hints for the message.
    If omitted, Legacy decides where to route it.
    """
    target_ids: List[str] = field(default_factory=list)   # Explicit targets
    target_type: Optional[str] = None                     # "agent", "team", "robot", etc.
    capability: Optional[str] = None                      # Route to entities with this capability
    mode: str = "auto"                                    # auto | direct | broadcast
    reply_to: Optional[str] = None                        # NATS subject for response
    callback_url: Optional[str] = None                    # HTTP callback for response

    def to_dict(self) -> Dict[str, Any]:
        return {k: v for k, v in asdict(self).items() if v is not None and v != [] and v != "auto"}

    @staticmethod
    def from_dict(data: Dict[str, Any]) -> "Routing":
        return Routing(
            target_ids=data.get("target_ids", []),
            target_type=data.get("target_type"),
            capability=data.get("capability"),
            mode=data.get("mode", "auto"),
            reply_to=data.get("reply_to"),
            callback_url=data.get("callback_url"),
        )


# =============================================================================
# NEXUS MESSAGE — The universal envelope
# =============================================================================

@dataclass
class NexusMessage:
    """
    The universal message envelope for ALL communication in OneMind OS.

    Every interaction — voice command, chat message, camera snapshot,
    robot telemetry, agent response — uses this format.

    This is the lingua franca of the entire system.
    """

    # ── Core ──
    content: str                                           # The text payload (always present, even if just a description)
    modality: Modality = Modality.TEXT                     # What kind of communication
    id: str = field(default_factory=lambda: str(uuid4()))  # Unique message ID

    # ── Source ──
    source_id: str = "unknown"                             # Who sent it (e.g. "human:zeus", "agent:researcher")
    source_type: SourceType = SourceType.HUMAN             # What kind of sender

    # ── Attachments ──
    attachments: List[Attachment] = field(default_factory=list)  # Images, audio, files, sensor data

    # ── Structured payload ──
    data: Dict[str, Any] = field(default_factory=dict)     # Arbitrary structured data

    # ── Routing ──
    routing: Optional[Routing] = None                     # Where this should go

    # ── Context ──
    session_id: Optional[str] = None                       # Conversation session
    thread_id: Optional[str] = None                        # Thread within a session
    reply_to_id: Optional[str] = None                      # Message this replies to
    correlation_id: Optional[str] = None                   # Links related messages across entities
    priority: Priority = Priority.NORMAL

    # ── Metadata ──
    timestamp: str = field(default_factory=lambda: datetime.utcnow().isoformat())
    metadata: Dict[str, Any] = field(default_factory=dict)

    def to_dict(self) -> Dict[str, Any]:
        """Serialize to dict (safe for JSON/NATS)."""
        d = {
            "id": self.id,
            "content": self.content,
            "modality": self.modality.value,
            "source_id": self.source_id,
            "source_type": self.source_type.value,
            "priority": self.priority.value,
            "timestamp": self.timestamp,
        }
        if self.attachments:
            d["attachments"] = [a.to_dict() for a in self.attachments]
        if self.data:
            d["data"] = self.data
        if self.routing:
            d["routing"] = self.routing.to_dict()
        if self.session_id:
            d["session_id"] = self.session_id
        if self.thread_id:
            d["thread_id"] = self.thread_id
        if self.reply_to_id:
            d["reply_to_id"] = self.reply_to_id
        if self.correlation_id:
            d["correlation_id"] = self.correlation_id
        if self.metadata:
            d["metadata"] = self.metadata
        return d

    def to_json(self) -> str:
        return json.dumps(self.to_dict(), default=str)

    @staticmethod
    def from_dict(data: Dict[str, Any]) -> "NexusMessage":
        """Deserialize from dict."""
        attachments = [Attachment.from_dict(a) for a in data.get("attachments", [])]
        routing = Routing.from_dict(data["routing"]) if data.get("routing") else None

        return NexusMessage(
            id=data.get("id", str(uuid4())),
            content=data.get("content", ""),
            modality=Modality(data.get("modality", "text")),
            source_id=data.get("source_id", "unknown"),
            source_type=SourceType(data.get("source_type", "human")),
            attachments=attachments,
            data=data.get("data", {}),
            routing=routing,
            session_id=data.get("session_id"),
            thread_id=data.get("thread_id"),
            reply_to_id=data.get("reply_to_id"),
            correlation_id=data.get("correlation_id"),
            priority=Priority(data.get("priority", "normal")),
            timestamp=data.get("timestamp", datetime.utcnow().isoformat()),
            metadata=data.get("metadata", {}),
        )

    @staticmethod
    def from_json(text: str) -> "NexusMessage":
        return NexusMessage.from_dict(json.loads(text))

    # ── Convenience Factories ──

    @staticmethod
    def text(content: str, source_id: str = "human:user", **kwargs) -> "NexusMessage":
        """Create a text chat message."""
        return NexusMessage(content=content, modality=Modality.TEXT, source_id=source_id, source_type=SourceType.HUMAN, **kwargs)

    @staticmethod
    def voice(transcript: str, audio_url: str = None, audio_base64: str = None,
              source_id: str = "human:user", **kwargs) -> "NexusMessage":
        """Create a voice message (transcript required, audio optional)."""
        attachments = []
        if audio_url or audio_base64:
            attachments.append(Attachment(
                type=AttachmentType.AUDIO,
                url=audio_url,
                data_base64=audio_base64,
                mime_type="audio/wav",
            ))
        return NexusMessage(
            content=transcript,
            modality=Modality.VOICE,
            source_id=source_id,
            source_type=SourceType.HUMAN,
            attachments=attachments,
            **kwargs,
        )

    @staticmethod
    def vision(description: str, image_url: str = None, image_base64: str = None,
               source_id: str = "sensor:camera", **kwargs) -> "NexusMessage":
        """Create a vision message (image + description/question)."""
        attachments = []
        if image_url or image_base64:
            attachments.append(Attachment(
                type=AttachmentType.IMAGE,
                url=image_url,
                data_base64=image_base64,
                mime_type="image/jpeg",
            ))
        return NexusMessage(
            content=description,
            modality=Modality.VISION,
            source_id=source_id,
            source_type=SourceType.SENSOR,
            attachments=attachments,
            **kwargs,
        )

    @staticmethod
    def sensor(reading: str, data: Dict[str, Any], source_id: str = "sensor:unknown", **kwargs) -> "NexusMessage":
        """Create a sensor telemetry message."""
        return NexusMessage(
            content=reading,
            modality=Modality.SENSOR,
            source_id=source_id,
            source_type=SourceType.SENSOR,
            data=data,
            **kwargs,
        )

    @staticmethod
    def agent_response(content: str, agent_id: str, reply_to_id: str = None, **kwargs) -> "NexusMessage":
        """Create a response from an agent."""
        return NexusMessage(
            content=content,
            modality=Modality.TEXT,
            source_id=f"agent:{agent_id}",
            source_type=SourceType.AGENT,
            reply_to_id=reply_to_id,
            **kwargs,
        )

    # ── Helpers ──

    @property
    def has_images(self) -> bool:
        return any(a.type == AttachmentType.IMAGE for a in self.attachments)

    @property
    def has_audio(self) -> bool:
        return any(a.type == AttachmentType.AUDIO for a in self.attachments)

    @property
    def has_video(self) -> bool:
        return any(a.type == AttachmentType.VIDEO for a in self.attachments)

    @property
    def images(self) -> List[Attachment]:
        return [a for a in self.attachments if a.type == AttachmentType.IMAGE]

    @property
    def audio_clips(self) -> List[Attachment]:
        return [a for a in self.attachments if a.type == AttachmentType.AUDIO]

    @property
    def nats_subject(self) -> str:
        """Generate the canonical NATS subject for this message."""
        src_clean = self.source_id.replace(":", ".").replace(" ", "_")
        return f"comms.{self.modality.value}.{self.source_type.value}.{src_clean}"

    def get_images_for_agno(self) -> Optional[List[Dict[str, Any]]]:
        """
        Convert image attachments to the format Agno agents expect.
        Returns list of {"url": "..."} dicts compatible with Agent.arun(images=...).
        """
        if not self.has_images:
            return None
        result = []
        for img in self.images:
            if img.url:
                result.append({"url": img.url})
            elif img.data_base64:
                mime = img.mime_type or "image/jpeg"
                result.append({"url": f"data:{mime};base64,{img.data_base64}"})
        return result if result else None


# =============================================================================
# MESSAGE BUS HELPERS — Publish/subscribe using NexusMessage
# =============================================================================

async def publish_message(msg: NexusMessage, subject: str = None):
    """
    Publish a NexusMessage to the NATS bus.

    If no subject is provided, the canonical subject is derived:
        comms.{modality}.{source_type}.{source_id}
    """
    from backend.bus import publish
    target_subject = subject or msg.nats_subject
    await publish(target_subject, msg.to_dict())
    logger.debug(f"Published message {msg.id} to {target_subject}")


async def publish_response(
    original: NexusMessage,
    content: str,
    source_id: str,
    source_type: SourceType = SourceType.AGENT,
    **kwargs,
):
    """
    Publish a response to an original message.
    Automatically sets reply_to_id and correlation_id.
    """
    response = NexusMessage(
        content=content,
        modality=Modality.TEXT,
        source_id=source_id,
        source_type=source_type,
        reply_to_id=original.id,
        correlation_id=original.correlation_id or original.id,
        session_id=original.session_id,
        thread_id=original.thread_id,
        **kwargs,
    )
    reply_subject = original.routing.reply_to if original.routing and original.routing.reply_to else None
    await publish_message(response, subject=reply_subject)
    return response
