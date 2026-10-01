#!/usr/bin/env python3
"""
NATS → Matrix Bridge Bot
Subscribes to NATS subjects, publishes formatted messages to Matrix rooms.
Part of OneMind OS fabric integration.
"""

import asyncio
import json
import os
import signal
import sys
from dataclasses import dataclass
from typing import Optional
from urllib.parse import urljoin

import nats
from nats.aio.msg import Msg
from matrix_nio import AsyncClient, MatrixRoom, RoomMessageText
from matrix_nio.crypto import OlmDevice
import structlog

logger = structlog.get_logger()

# =============================================================================
# Configuration
# =============================================================================

@dataclass
class Config:
    # NATS
    nats_url: str = os.getenv("NATS_URL", "nats://localhost:4222")
    
    # Matrix
    matrix_homeserver: str = os.getenv("MATRIX_HOMESERVER", "https://matrix.onemindos.com")
    matrix_user_id: str = os.getenv("MATRIX_USER_ID", "@legacy-bot:onemindos.com")
    matrix_access_token: str = os.getenv("MATRIX_ACCESS_TOKEN", "")
    matrix_device_id: Optional[str] = os.getenv("MATRIX_DEVICE_ID")
    
    # Room mappings (NATS subject pattern → Matrix room ID)
    room_mappings: dict = None
    
    # Bot behavior
    max_message_length: int = 4000
    enable_encryption: bool = True
    
    def __post_init__(self):
        if self.room_mappings is None:
            self.room_mappings = {
                "alerts.critical": "!alerts:onemindos.com",
                "alerts.warning": "!alerts:onemindos.com",
                "alerts.info": "!alerts:onemindos.com",
                "tak.cot.mission": "!missions:onemindos.com",
                "tak.cot.emergency": "!alerts:onemindos.com",
                "tak.cot.alert": "!alerts:onemindos.com",
                "infra.health.degraded": "!infra:onemindos.com",
                "infra.health.critical": "!alerts:onemindos.com",
                "infra.health.recovered": "!infra:onemindos.com",
                "legacy.voice.transcript": "!voice:onemindos.com",
                "legacy.task.complete": "!legacy:onemindos.com",
                "legacy.task.failed": "!legacy:onemindos.com",
                "geo.sensor.alert": "!sensors:onemindos.com",
                "fabric.agent.status": "!agents:onemindos.com",
            }

CONFIG = Config()

# =============================================================================
# Message Formatters
# =============================================================================

class MessageFormatter:
    @staticmethod
    def format_alert(subject: str, data: dict) -> str:
        severity = "critical" if "critical" in subject else "warning" if "warning" in subject else "info"
        emoji = {"critical": "🚨", "warning": "⚠️", "info": "ℹ️"}[severity]
        
        title = data.get("title", subject)
        message = data.get("message", json.dumps(data, indent=2)[:500])
        source = data.get("source", "unknown")
        timestamp = data.get("timestamp", "")
        
        lines = [
            f"{emoji} **ALERT [{severity.upper()}]**",
            f"**{title}**",
            f"",
            message,
            f"",
            f"📍 Source: `{source}`",
        ]
        if timestamp:
            lines.append(f"🕐 Time: `{timestamp}`")
        
        return "\n".join(lines)
    
    @staticmethod
    def format_tak_cot(subject: str, data: dict) -> str:
        cot_type = data.get("type", "unknown")
        callsign = data.get("callsign", "UNKNOWN")
        lat = data.get("lat")
        lon = data.get("lon")
        remarks = data.get("remarks", "")
        
        emoji_map = {
            "a-f-G-U-C": "🚨",  # Emergency
            "a-f-G-U-C-R": "🔴",  # Rescue
            "a-h-G-U-C": "🏥",  # Medical
            "a-u-G": "👤",  # Friendly
            "a-u-S": "❓",  # Unknown
        }
        emoji = emoji_map.get(cot_type, "📍")
        
        lines = [
            f"{emoji} **TAK CoT: {cot_type}**",
            f"**Callsign:** `{callsign}`",
        ]
        
        if lat and lon:
            lines.append(f"**Location:** `{lat:.6f}, {lon:.6f}`")
            lines.append(f"🗺️ <https://maps.onemindos.com/?lat={lat}&lon={lon}|Map>")
        
        if remarks:
            lines.extend(["", f"**Remarks:**", remarks])
        
        return "\n".join(lines)
    
    @staticmethod
    def format_infra_health(subject: str, data: dict) -> str:
        status = data.get("status", "unknown")
        service = data.get("service", "unknown")
        details = data.get("details", {})
        
        emoji_map = {
            "healthy": "✅",
            "degraded": "⚠️",
            "critical": "🔴",
            "recovered": "🟢",
            "unknown": "❓",
        }
        emoji = emoji_map.get(status, "❓")
        
        lines = [
            f"{emoji} **INFRA: {service.upper()} → {status.upper()}**",
        ]
        
        if details:
            lines.append("")
            for key, value in details.items():
                lines.append(f"• `{key}`: `{value}`")
        
        return "\n".join(lines)
    
    @staticmethod
    def format_voice_transcript(data: dict) -> str:
        speaker = data.get("speaker", "Unknown")
        transcript = data.get("transcript", "")
        channel = data.get("channel", "unknown")
        confidence = data.get("confidence", 0)
        
        conf_emoji = "🟢" if confidence > 0.9 else "🟡" if confidence > 0.7 else "🔴"
        
        return f"🎙️ **VOICE [{channel}]** {conf_emoji}\n**{speaker}:**\n{transcript}"
    
    @staticmethod
    def format_task_result(subject: str, data: dict) -> str:
        task_id = data.get("task_id", "unknown")
        status = data.get("status", "unknown")
        result = data.get("result", "")
        error = data.get("error", "")
        
        emoji = "✅" if status == "complete" else "❌" if status == "failed" else "⏳"
        
        lines = [
            f"{emoji} **TASK {status.upper()}**",
            f"**Task ID:** `{task_id}`",
        ]
        
        if result:
            lines.extend(["", "**Result:**", f"```\n{result[:2000]}\n```"])
        
        if error:
            lines.extend(["", "**Error:**", f"```\n{error[:1000]}\n```"])
        
        return "\n".join(lines)
    
    @staticmethod
    def format_generic(subject: str, data: dict) -> str:
        return f"📡 **{subject}**\n```json\n{json.dumps(data, indent=2)[:3000]}\n```"
    
    @classmethod
    def format(cls, subject: str, data: dict) -> str:
        if subject.startswith("alerts."):
            return cls.format_alert(subject, data)
        elif subject.startswith("tak.cot."):
            return cls.format_tak_cot(subject, data)
        elif subject.startswith("infra.health."):
            return cls.format_infra_health(subject, data)
        elif subject.startswith("legacy.voice."):
            return cls.format_voice_transcript(data)
        elif subject.startswith("legacy.task."):
            return cls.format_task_result(subject, data)
        else:
            return cls.format_generic(subject, data)

# =============================================================================
# Matrix Client Wrapper
# =============================================================================

class MatrixBot:
    def __init__(self, config: Config):
        self.config = config
        self.client = AsyncClient(
            config.matrix_homeserver,
            config.matrix_user_id,
            device_id=config.matrix_device_id
        )
        self.client.access_token = config.matrix_access_token
        self.formatter = MessageFormatter()
        self._synced = False
    
    async def connect(self):
        logger.info("Connecting to Matrix", homeserver=self.config.matrix_homeserver)
        
        if self.config.enable_encryption:
            await self.client.keys_upload()
            await self.client.key_query()
        
        # Initial sync
        await self.client.sync(full_state=True)
        self._synced = True
        
        logger.info("Matrix bot connected", user_id=self.config.matrix_user_id)
    
    async def send_message(self, room_id: str, message: str, thread_root: Optional[str] = None):
        if not self._synced:
            await self.connect()
        
        # Truncate if needed
        if len(message) > self.config.max_message_length:
            message = message[:self.config.max_message_length - 100] + "\n\n[Truncated]"
        
        content = {
            "msgtype": "m.text",
            "body": message,
            "format": "org.matrix.custom.html",
            "formatted_body": message.replace("\n", "<br>")
        }
        
        if thread_root:
            content["m.relates_to"] = {
                "rel_type": "m.thread",
                "event_id": thread_root
            }
        
        response = await self.client.room_send(
            room_id=room_id,
            message_type="m.room.message",
            content=content
        )
        
        if hasattr(response, "event_id"):
            return response.event_id
        return None
    
    async def close(self):
        await self.client.close()

# =============================================================================
# NATS Subscriber
# =============================================================================

class NATSMatrixBridge:
    def __init__(self, config: Config):
        self.config = config
        self.nc = None
        self.matrix = MatrixBot(config)
        self.subscriptions = []
        self.running = False
    
    async def connect_nats(self):
        self.nc = await nats.connect(self.config.nats_url)
        logger.info("Connected to NATS", url=self.config.nats_url)
    
    async def subscribe_subject(self, subject: str, room_id: str):
        async def handler(msg: Msg):
            try:
                data = json.loads(msg.data.decode())
                formatted = self.formatter.format(subject, data)
                await self.matrix.send_message(room_id, formatted)
                logger.debug("Forwarded to Matrix", subject=subject, room=room_id)
            except Exception as e:
                logger.error("Failed to forward message", subject=subject, error=str(e))
        
        sub = await self.nc.subscribe(subject, cb=handler)
        self.subscriptions.append((subject, sub))
        logger.info("Subscribed", subject=subject, room=room_id)
    
    async def start(self):
        await self.connect_nats()
        await self.matrix.connect()
        
        # Subscribe to all mapped subjects
        for subject, room_id in self.config.room_mappings.items():
            await self.subscribe_subject(subject, room_id)
        
        # Wildcard subscriptions for patterns
        wildcard_patterns = [
            ("alerts.>", "!alerts:onemindos.com"),
            ("tak.cot.>", "!missions:onemindos.com"),
            ("infra.health.>", "!infra:onemindos.com"),
            ("legacy.>", "!legacy:onemindos.com"),
            ("geo.>", "!sensors:onemindos.com"),
            ("fabric.agent.>", "!agents:onemindos.com"),
        ]
        
        for pattern, room_id in wildcard_patterns:
            await self.subscribe_subject(pattern, room_id)
        
        self.running = True
        logger.info("NATS-Matrix bridge started", subscriptions=len(self.subscriptions))
    
    async def stop(self):
        self.running = False
        for subject, sub in self.subscriptions:
            await sub.unsubscribe()
        if self.nc:
            await self.nc.close()
        await self.matrix.close()
        logger.info("NATS-Matrix bridge stopped")

# =============================================================================
# Main
# =============================================================================

async def main():
    structlog.configure(
        processors=[
            structlog.processors.TimeStamper(fmt="iso"),
            structlog.processors.add_log_level,
            structlog.processors.JSONRenderer()
        ]
    )
    
    bridge = NATSMatrixBridge(CONFIG)
    
    # Signal handling
    loop = asyncio.get_event_loop()
    for sig in (signal.SIGTERM, signal.SIGINT):
        loop.add_signal_handler(sig, lambda: asyncio.create_task(bridge.stop()))
    
    try:
        await bridge.start()
        # Keep running
        while bridge.running:
            await asyncio.sleep(1)
    except KeyboardInterrupt:
        pass
    finally:
        await bridge.stop()

if __name__ == "__main__":
    asyncio.run(main())