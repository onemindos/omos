# OneMindOS v3 Enhancement Roadmap

**Date**: 2025-02-17
**Status**: Implementation Planning
**Framework**: Agno AgentOS + Custom Extensions

---

## Overview

This document outlines three critical enhancements for OneMindOS v3 based on architectural analysis of OpenClaw, PicoClaw, and Nanobot. Each feature has been evaluated for Agno AgentOS compatibility and implementation complexity.

**Key Principle**: Don't reinvent what Agno provides. Build only what's missing.

---

## Priority Matrix

| Feature | Priority | Complexity | Timeline | Dependencies |
|---------|----------|-----------|----------|--------------|
| **Browser Automation** | 🔥 High | Medium | 2-3 days | Playwright |
| **Device Event System** | 🔥 High | Low | 1-2 days | NATS (✅ exists) |
| **Voice Pipeline (Full)** | Medium | High | 4-5 days | Porcupine, Groq, WebRTC |

---

## 1. Browser Automation 🌐

### **Objective**
Enable agents to interact with web pages programmatically for:
- Form filling (visa applications, job applications)
- Data extraction (prices, flight status, dashboards)
- Website monitoring (detect changes, capture screenshots)
- Automated testing

### **Agno Status**
❌ **NOT PROVIDED** - Agno has no built-in browser automation

### **Architecture**

```
┌─────────────────────────────────────────────────────────┐
│              Agno Agent (Decision Making)               │
│  - Analyzes user request                                │
│  - Decides which browser tool to use                    │
│  - Interprets results                                   │
└─────────────────────────────────────────────────────────┘
                    ↓ Tool call
┌─────────────────────────────────────────────────────────┐
│         Browser Tools (backend/tools/browser.py)        │
│  @tool navigate_and_extract(url, selector)             │
│  @tool fill_form(url, fields)                          │
│  @tool take_screenshot(url)                            │
│  @tool monitor_page(url, selector, interval)           │
└─────────────────────────────────────────────────────────┘
                    ↓ Browser control
┌─────────────────────────────────────────────────────────┐
│            Playwright Browser Instance                  │
│  - Headless Chromium                                    │
│  - JavaScript execution                                 │
│  - Network interception (optional)                      │
│  - Cookie/session management                            │
└─────────────────────────────────────────────────────────┘
                    ↓ Web interaction
┌─────────────────────────────────────────────────────────┐
│                  Target Websites                        │
└─────────────────────────────────────────────────────────┘
```

### **Implementation Plan**

#### **Phase 1: Core Tools** (Day 1)

**File**: `backend/tools/browser.py` (NEW)

```python
"""
Browser Automation Tools for OneMindOS
Uses Playwright for web page interaction
"""

from agno.tools import tool
from playwright.async_api import async_playwright
import asyncio
from typing import Dict, List, Optional
import logging

logger = logging.getLogger(__name__)


@tool
async def navigate_and_extract(
    url: str,
    selector: str,
    wait_for: Optional[str] = None
) -> str:
    """
    Navigate to a URL and extract text content from CSS selector.

    Args:
        url: Target webpage URL
        selector: CSS selector for element to extract
        wait_for: Optional selector to wait for before extracting

    Returns:
        Extracted text content

    Example:
        navigate_and_extract(
            "https://weather.com",
            ".temperature-value",
            wait_for=".weather-data"
        )
    """
    async with async_playwright() as p:
        browser = await p.chromium.launch(headless=True)
        page = await browser.new_page()

        try:
            await page.goto(url, wait_until="domcontentloaded")

            # Wait for specific element if requested
            if wait_for:
                await page.wait_for_selector(wait_for, timeout=10000)

            # Extract text content
            text = await page.text_content(selector)

            if not text:
                return f"No content found at selector: {selector}"

            logger.info(f"Extracted from {url}: {text[:100]}...")
            return text.strip()

        except Exception as e:
            logger.error(f"Browser navigation error: {e}")
            return f"Error: {str(e)}"
        finally:
            await browser.close()


@tool
async def fill_form(
    url: str,
    fields: Dict[str, str],
    submit_selector: str = "button[type='submit']"
) -> str:
    """
    Navigate to a form and fill in fields.

    Args:
        url: Form URL
        fields: Dictionary of {selector: value} pairs
        submit_selector: CSS selector for submit button

    Returns:
        Success message or error

    Example:
        fill_form(
            "https://example.com/login",
            {
                "input[name='email']": "user@example.com",
                "input[name='password']": "secret123"
            },
            submit_selector="button.login-btn"
        )
    """
    async with async_playwright() as p:
        browser = await p.chromium.launch(headless=True)
        page = await browser.new_page()

        try:
            await page.goto(url, wait_until="domcontentloaded")

            # Fill each field
            for selector, value in fields.items():
                await page.fill(selector, value)
                logger.info(f"Filled field: {selector}")

            # Submit form
            await page.click(submit_selector)
            await page.wait_for_load_state("networkidle")

            # Get final URL to check for redirect
            final_url = page.url

            logger.info(f"Form submitted. Final URL: {final_url}")
            return f"Form submitted successfully. Redirected to: {final_url}"

        except Exception as e:
            logger.error(f"Form filling error: {e}")
            return f"Error: {str(e)}"
        finally:
            await browser.close()


@tool
async def take_screenshot(
    url: str,
    selector: Optional[str] = None,
    wait_seconds: int = 2
) -> str:
    """
    Take a screenshot of a webpage.

    Args:
        url: Target webpage URL
        selector: Optional CSS selector to screenshot (full page if None)
        wait_seconds: Seconds to wait after page load

    Returns:
        Path to saved screenshot

    Example:
        take_screenshot(
            "https://dashboard.example.com",
            selector=".metrics-panel",
            wait_seconds=3
        )
    """
    async with async_playwright() as p:
        browser = await p.chromium.launch(headless=True)
        page = await browser.new_page()

        try:
            await page.goto(url, wait_until="domcontentloaded")
            await asyncio.sleep(wait_seconds)

            # Generate filename
            import time
            filename = f"screenshot_{int(time.time())}.png"
            filepath = f"/tmp/{filename}"

            if selector:
                # Screenshot specific element
                element = await page.query_selector(selector)
                if element:
                    await element.screenshot(path=filepath)
                else:
                    return f"Error: Element not found: {selector}"
            else:
                # Full page screenshot
                await page.screenshot(path=filepath, full_page=True)

            logger.info(f"Screenshot saved: {filepath}")
            return filepath

        except Exception as e:
            logger.error(f"Screenshot error: {e}")
            return f"Error: {str(e)}"
        finally:
            await browser.close()


@tool
async def extract_structured_data(
    url: str,
    extractors: Dict[str, str]
) -> Dict[str, str]:
    """
    Extract multiple pieces of data from a webpage.

    Args:
        url: Target webpage URL
        extractors: Dictionary of {field_name: css_selector}

    Returns:
        Dictionary of extracted data

    Example:
        extract_structured_data(
            "https://product.example.com/item",
            {
                "title": "h1.product-name",
                "price": "span.price",
                "rating": "div.rating",
                "availability": "span.stock-status"
            }
        )
    """
    async with async_playwright() as p:
        browser = await p.chromium.launch(headless=True)
        page = await browser.new_page()

        try:
            await page.goto(url, wait_until="domcontentloaded")

            results = {}
            for field_name, selector in extractors.items():
                try:
                    text = await page.text_content(selector)
                    results[field_name] = text.strip() if text else None
                except Exception as e:
                    results[field_name] = f"Error: {str(e)}"

            logger.info(f"Extracted structured data from {url}")
            return results

        except Exception as e:
            logger.error(f"Data extraction error: {e}")
            return {"error": str(e)}
        finally:
            await browser.close()
```

#### **Phase 2: Tool Registration** (Day 1)

**File**: `backend/tools/registry.py` (MODIFY)

```python
# Add browser tools import
from backend.tools.browser import (
    navigate_and_extract,
    fill_form,
    take_screenshot,
    extract_structured_data,
)

# Add to TOOL_REGISTRY
TOOL_REGISTRY = {
    # ... existing tools ...

    # Browser automation
    "browser_navigate": navigate_and_extract,
    "browser_form": fill_form,
    "browser_screenshot": take_screenshot,
    "browser_extract": extract_structured_data,
}

# Add browser preset
TOOL_PRESETS = {
    # ... existing presets ...

    "browser": [
        navigate_and_extract,
        fill_form,
        take_screenshot,
        extract_structured_data,
    ],
}
```

#### **Phase 3: Agent Configuration** (Day 2)

**File**: `backend/agents/presets.py` (MODIFY)

```python
AGENT_PRESETS = {
    # ... existing agents ...

    "web_automator": {
        "name": "Web Automator",
        "model_id": "gpt-4o",
        "tools": ["browser", "duckduckgo"],  # Browser preset
        "instructions": [
            "You are a web automation specialist.",
            "You can navigate websites, extract data, fill forms, and take screenshots.",
            "Always verify URLs before interacting with them.",
            "For forms, confirm field selectors are correct before filling.",
            "When taking screenshots, save them and provide the file path.",
        ],
        "enable_user_memories": True,
        "enable_reasoning": True,
        "reasoning_max_steps": 8,
        "run_parallel_tools": False,  # Sequential for browser actions
    },
}
```

#### **Phase 4: Advanced Features** (Day 3)

**File**: `backend/tools/browser_advanced.py` (NEW)

```python
"""
Advanced browser automation features
"""

from agno.tools import tool
from playwright.async_api import async_playwright
import asyncio
from typing import List, Dict


@tool
async def monitor_page_for_changes(
    url: str,
    selector: str,
    interval_seconds: int = 60,
    duration_minutes: int = 60
) -> List[Dict]:
    """
    Monitor a webpage for changes over time.

    Args:
        url: Target webpage URL
        selector: CSS selector to monitor
        interval_seconds: Check interval in seconds
        duration_minutes: Total monitoring duration

    Returns:
        List of change events with timestamps

    Example:
        monitor_page_for_changes(
            "https://stocks.example.com/AAPL",
            ".stock-price",
            interval_seconds=30,
            duration_minutes=60
        )
    """
    changes = []
    last_value = None
    iterations = (duration_minutes * 60) // interval_seconds

    async with async_playwright() as p:
        browser = await p.chromium.launch(headless=True)
        page = await browser.new_page()

        try:
            for i in range(iterations):
                await page.goto(url, wait_until="domcontentloaded")
                current_value = await page.text_content(selector)

                if current_value and current_value != last_value:
                    changes.append({
                        "timestamp": asyncio.get_event_loop().time(),
                        "old_value": last_value,
                        "new_value": current_value
                    })
                    last_value = current_value

                await asyncio.sleep(interval_seconds)

            return changes

        finally:
            await browser.close()


@tool
async def execute_javascript(
    url: str,
    javascript_code: str
) -> str:
    """
    Execute custom JavaScript on a webpage.

    Args:
        url: Target webpage URL
        javascript_code: JavaScript code to execute

    Returns:
        Result of JavaScript execution

    Example:
        execute_javascript(
            "https://example.com",
            "return document.querySelectorAll('a').length"
        )
    """
    async with async_playwright() as p:
        browser = await p.chromium.launch(headless=True)
        page = await browser.new_page()

        try:
            await page.goto(url, wait_until="domcontentloaded")
            result = await page.evaluate(javascript_code)
            return str(result)

        except Exception as e:
            return f"Error executing JavaScript: {str(e)}"
        finally:
            await browser.close()
```

### **Testing Plan**

**File**: `tests/test_browser_tools.py` (NEW)

```python
import pytest
from backend.tools.browser import (
    navigate_and_extract,
    fill_form,
    take_screenshot,
)


@pytest.mark.asyncio
async def test_navigate_and_extract():
    """Test basic page navigation and extraction"""
    result = await navigate_and_extract(
        "https://example.com",
        "h1"
    )
    assert "Example Domain" in result


@pytest.mark.asyncio
async def test_screenshot():
    """Test screenshot capture"""
    filepath = await take_screenshot("https://example.com")
    assert filepath.endswith(".png")
    # Verify file exists
    import os
    assert os.path.exists(filepath)
```

### **Dependencies**

```bash
# Add to requirements.txt or pyproject.toml
playwright==1.41.0
```

```bash
# Install Playwright browsers
playwright install chromium
```

### **Use Cases**

1. **Flight Check-in**: "Check in for my United flight tomorrow"
2. **Price Monitoring**: "Monitor Amazon product price and alert me when it drops"
3. **Form Submission**: "Fill out this visa application form for me"
4. **Dashboard Screenshots**: "Take a screenshot of my server metrics dashboard"
5. **Data Extraction**: "Extract all product prices from this competitor website"

---

## 2. Device Event System 🔌

### **Objective**
Treat hardware devices as event sources that publish to NATS bus, enabling:
- Loose coupling between devices and agents
- Multiple agents reacting to same device event
- Easy addition of new sensors without agent changes
- Full event logging and debugging

### **Agno Status**
⚠️ **PARTIAL** - Agno doesn't provide device abstraction, but OneMindOS already has:
- ✅ NATS bus ([backend/bus.py](backend/bus.py))
- ✅ Event gateway ([backend/events/gateway.py](backend/events/gateway.py))
- ✅ WebSocket broadcasting ([backend/websocket_handler.py](backend/websocket_handler.py))

### **Architecture**

```
┌─────────────────────────────────────────────────────────┐
│                 Device Event Sources                    │
│  ┌─────────────┐  ┌─────────────┐  ┌─────────────┐    │
│  │  MaixCAM    │  │ USB Sensor  │  │ Smart Home  │    │
│  │  (Vision)   │  │ (I2C/SPI)   │  │ (Home Asst) │    │
│  └─────────────┘  └─────────────┘  └─────────────┘    │
└─────────────────────────────────────────────────────────┘
         ↓                  ↓                  ↓
         └──────────────────┴──────────────────┘
                           ↓
┌─────────────────────────────────────────────────────────┐
│              Device Event Adapters                      │
│  - Normalize device-specific data                       │
│  - Add metadata (timestamp, confidence, location)       │
│  - Validate event structure                             │
└─────────────────────────────────────────────────────────┘
                           ↓ Publish
┌─────────────────────────────────────────────────────────┐
│            NATS Event Bus (Already Exists!)             │
│  Topics:                                                │
│    onemind.devices.camera.person_detected              │
│    onemind.devices.sensor.temperature_high             │
│    onemind.devices.sensor.motion_detected              │
│    onemind.devices.usb.hotplug                         │
│    onemind.devices.home.door_opened                    │
└─────────────────────────────────────────────────────────┘
         ↓ Subscribe            ↓ Subscribe            ↓
┌──────────────┐      ┌──────────────┐      ┌──────────────┐
│ Security     │      │ HVAC Agent   │      │ Monitor      │
│ Agent        │      │              │      │ Agent        │
└──────────────┘      └──────────────┘      └──────────────┘
```

### **Event Schema**

```python
# Standard device event structure
{
    "event_type": "person_detected",
    "device_id": "living_room_camera_01",
    "device_type": "camera",
    "timestamp": "2025-02-17T14:35:22.123Z",
    "confidence": 0.87,
    "data": {
        "class": "person",
        "bounding_box": [100, 150, 300, 450],
        "snapshot_url": "http://camera/snapshot.jpg"
    },
    "metadata": {
        "location": "living_room",
        "zone": "indoor",
        "floor": 1
    }
}
```

### **Implementation Plan**

#### **Phase 1: Core Device Framework** (Day 1, Morning)

**File**: `backend/devices/__init__.py` (NEW)

```python
"""
Device Event System for OneMindOS
Treats hardware as event sources that publish to NATS
"""

from abc import ABC, abstractmethod
from typing import Dict, Any, Optional
from datetime import datetime
import logging

logger = logging.getLogger(__name__)


class DeviceEvent:
    """Standard device event structure"""

    def __init__(
        self,
        event_type: str,
        device_id: str,
        device_type: str,
        data: Dict[str, Any],
        confidence: Optional[float] = None,
        metadata: Optional[Dict[str, Any]] = None,
    ):
        self.event_type = event_type
        self.device_id = device_id
        self.device_type = device_type
        self.data = data
        self.confidence = confidence
        self.metadata = metadata or {}
        self.timestamp = datetime.utcnow().isoformat()

    def to_dict(self) -> Dict[str, Any]:
        """Convert to dictionary for NATS publishing"""
        return {
            "event_type": self.event_type,
            "device_id": self.device_id,
            "device_type": self.device_type,
            "timestamp": self.timestamp,
            "confidence": self.confidence,
            "data": self.data,
            "metadata": self.metadata,
        }

    def to_nats_topic(self) -> str:
        """Generate NATS topic from event"""
        return f"onemind.devices.{self.device_type}.{self.event_type}"


class DeviceSource(ABC):
    """
    Abstract base class for device event sources.
    All device integrations should inherit from this.
    """

    def __init__(self, device_id: str, device_type: str):
        self.device_id = device_id
        self.device_type = device_type
        self.running = False
        logger.info(f"Device source initialized: {device_id} ({device_type})")

    @abstractmethod
    async def start(self):
        """Start monitoring device for events"""
        pass

    @abstractmethod
    async def stop(self):
        """Stop monitoring device"""
        pass

    @abstractmethod
    async def get_status(self) -> Dict[str, Any]:
        """Get device status/health"""
        pass


class DeviceRegistry:
    """
    Registry for managing all device sources
    """

    def __init__(self):
        self.sources: Dict[str, DeviceSource] = {}
        logger.info("Device registry initialized")

    def register(self, source: DeviceSource):
        """Register a device source"""
        self.sources[source.device_id] = source
        logger.info(f"Registered device: {source.device_id}")

    def unregister(self, device_id: str):
        """Unregister a device source"""
        if device_id in self.sources:
            del self.sources[device_id]
            logger.info(f"Unregistered device: {device_id}")

    async def start_all(self):
        """Start all registered devices"""
        for source in self.sources.values():
            try:
                await source.start()
            except Exception as e:
                logger.error(f"Failed to start {source.device_id}: {e}")

    async def stop_all(self):
        """Stop all registered devices"""
        for source in self.sources.values():
            try:
                await source.stop()
            except Exception as e:
                logger.error(f"Failed to stop {source.device_id}: {e}")

    def get_all_status(self) -> Dict[str, Any]:
        """Get status of all devices"""
        return {
            device_id: source.get_status()
            for device_id, source in self.sources.items()
        }


# Global registry singleton
_registry: Optional[DeviceRegistry] = None

def get_device_registry() -> DeviceRegistry:
    """Get the global device registry"""
    global _registry
    if _registry is None:
        _registry = DeviceRegistry()
    return _registry
```

#### **Phase 2: Camera Device Source** (Day 1, Afternoon)

**File**: `backend/devices/sources/camera.py` (NEW)

```python
"""
Camera device source for person/object detection
Supports MaixCAM and generic USB cameras
"""

import asyncio
from typing import Dict, Any
import logging
from backend.devices import DeviceSource, DeviceEvent
from backend.bus import publish

logger = logging.getLogger(__name__)


class CameraSource(DeviceSource):
    """
    Generic camera device source for vision-based detection
    """

    def __init__(
        self,
        device_id: str,
        camera_url: str,
        detection_classes: list[str] = None,
        confidence_threshold: float = 0.7,
        fps: int = 10,
    ):
        super().__init__(device_id, "camera")
        self.camera_url = camera_url
        self.detection_classes = detection_classes or ["person"]
        self.confidence_threshold = confidence_threshold
        self.fps = fps
        self.running = False
        self._task = None

    async def start(self):
        """Start camera monitoring"""
        self.running = True
        self._task = asyncio.create_task(self._monitoring_loop())
        logger.info(f"Camera started: {self.device_id}")

    async def stop(self):
        """Stop camera monitoring"""
        self.running = False
        if self._task:
            self._task.cancel()
        logger.info(f"Camera stopped: {self.device_id}")

    async def get_status(self) -> Dict[str, Any]:
        """Get camera status"""
        return {
            "device_id": self.device_id,
            "device_type": "camera",
            "running": self.running,
            "camera_url": self.camera_url,
            "fps": self.fps,
            "detection_classes": self.detection_classes,
        }

    async def _monitoring_loop(self):
        """Main detection loop"""
        interval = 1.0 / self.fps

        while self.running:
            try:
                # In real implementation, this would:
                # 1. Capture frame from camera
                # 2. Run inference (YOLO, etc.)
                # 3. Filter detections by confidence
                # 4. Publish events to NATS

                # Placeholder for actual detection logic
                detections = await self._detect_objects()

                for detection in detections:
                    if detection["confidence"] >= self.confidence_threshold:
                        await self._publish_detection(detection)

                await asyncio.sleep(interval)

            except Exception as e:
                logger.error(f"Camera monitoring error: {e}")
                await asyncio.sleep(interval)

    async def _detect_objects(self) -> list[Dict[str, Any]]:
        """
        Placeholder for actual object detection.
        In production, this would call MaixCAM API or run local inference.
        """
        # TODO: Integrate with actual camera/inference API
        return []

    async def _publish_detection(self, detection: Dict[str, Any]):
        """Publish detection event to NATS"""
        event = DeviceEvent(
            event_type=f"{detection['class']}_detected",
            device_id=self.device_id,
            device_type="camera",
            confidence=detection["confidence"],
            data={
                "class": detection["class"],
                "bounding_box": detection.get("bbox"),
                "snapshot_url": detection.get("snapshot_url"),
            },
            metadata={
                "location": self.metadata.get("location", "unknown"),
                "zone": self.metadata.get("zone", "unknown"),
            }
        )

        # Publish to NATS
        topic = event.to_nats_topic()
        await publish(topic, event.to_dict())

        logger.info(
            f"Detection published: {detection['class']} "
            f"({detection['confidence']:.2f}) on {self.device_id}"
        )


class MaixCAMSource(CameraSource):
    """
    Specialized source for Sipeed MaixCAM ($50 AI camera)
    """

    def __init__(
        self,
        device_id: str,
        maixcam_ip: str,
        detection_classes: list[str] = None,
    ):
        camera_url = f"http://{maixcam_ip}:8080/stream"
        super().__init__(
            device_id=device_id,
            camera_url=camera_url,
            detection_classes=detection_classes,
            confidence_threshold=0.7,
            fps=30,  # MaixCAM supports 30fps
        )
        self.maixcam_ip = maixcam_ip

    async def _detect_objects(self) -> list[Dict[str, Any]]:
        """
        Use MaixCAM's built-in detection API
        """
        import httpx

        try:
            async with httpx.AsyncClient() as client:
                response = await client.get(
                    f"http://{self.maixcam_ip}:8080/detect",
                    timeout=2.0
                )
                if response.status_code == 200:
                    return response.json().get("detections", [])
        except Exception as e:
            logger.error(f"MaixCAM API error: {e}")

        return []
```

#### **Phase 3: Agent Subscriptions** (Day 2, Morning)

**File**: `backend/devices/listeners.py` (NEW)

```python
"""
Agent listeners for device events
Agents subscribe to NATS topics and react to events
"""

import asyncio
import logging
from backend.bus import subscribe
from backend.agents.factory import create_agent_from_config
from backend.agents.presets import AGENT_PRESETS

logger = logging.getLogger(__name__)


class DeviceEventListener:
    """
    Listens to device events and triggers agent responses
    """

    def __init__(self, agent_id: str, topic_pattern: str):
        self.agent_id = agent_id
        self.topic_pattern = topic_pattern
        self.agent = None
        self.running = False
        self._task = None

    async def start(self):
        """Start listening for device events"""
        # Create agent from preset
        preset = AGENT_PRESETS.get(self.agent_id)
        if not preset:
            raise ValueError(f"Unknown agent: {self.agent_id}")

        self.agent = await create_agent_from_config(preset)
        self.running = True
        self._task = asyncio.create_task(self._listen_loop())
        logger.info(f"Device listener started: {self.agent_id} on {self.topic_pattern}")

    async def stop(self):
        """Stop listening"""
        self.running = False
        if self._task:
            self._task.cancel()
        logger.info(f"Device listener stopped: {self.agent_id}")

    async def _listen_loop(self):
        """Main event listening loop"""
        async for event in subscribe(self.topic_pattern):
            try:
                await self._handle_event(event)
            except Exception as e:
                logger.error(f"Event handling error: {e}")

    async def _handle_event(self, event: dict):
        """Process device event with agent"""
        # Build prompt from event
        prompt = self._build_prompt(event)

        # Run agent
        response = await self.agent.run(prompt)

        # Log agent decision
        logger.info(
            f"Agent {self.agent_id} handled {event['event_type']}: "
            f"{response.message[:100]}"
        )

        # Agent response can trigger actions via tools
        # (e.g., send notification, adjust thermostat, etc.)

    def _build_prompt(self, event: dict) -> str:
        """Build agent prompt from device event"""
        return f"""
Device Event Received:

Type: {event['event_type']}
Device: {event['device_id']} ({event['device_type']})
Timestamp: {event['timestamp']}
Confidence: {event.get('confidence', 'N/A')}

Data: {event.get('data', {})}
Location: {event.get('metadata', {}).get('location', 'unknown')}

What action should be taken? Consider:
- Is this normal or unexpected?
- Does the user need to be notified?
- Should any automation be triggered?
- Is this part of a pattern worth tracking?
"""


# Pre-configured listeners
DEVICE_LISTENERS = [
    # Security monitoring
    DeviceEventListener(
        agent_id="assistant",
        topic_pattern="onemind.devices.camera.person_detected"
    ),

    # Temperature monitoring
    DeviceEventListener(
        agent_id="assistant",
        topic_pattern="onemind.devices.sensor.temperature_*"
    ),

    # Motion detection
    DeviceEventListener(
        agent_id="assistant",
        topic_pattern="onemind.devices.sensor.motion_detected"
    ),
]


async def start_all_listeners():
    """Start all device event listeners"""
    for listener in DEVICE_LISTENERS:
        await listener.start()


async def stop_all_listeners():
    """Stop all device event listeners"""
    for listener in DEVICE_LISTENERS:
        await listener.stop()
```

#### **Phase 4: Integration with main.py** (Day 2, Afternoon)

**File**: `backend/main.py` (MODIFY)

```python
# Add to startup_event()

# ─── 9. Start Device Event System ───
try:
    from backend.devices import get_device_registry
    from backend.devices.listeners import start_all_listeners

    # Initialize device registry
    registry = get_device_registry()

    # Register devices (example)
    # from backend.devices.sources.camera import MaixCAMSource
    # maixcam = MaixCAMSource(
    #     device_id="living_room_camera",
    #     maixcam_ip="192.168.1.100"
    # )
    # registry.register(maixcam)

    # Start all devices
    await registry.start_all()
    logger.info("✓ Device registry started")

    # Start event listeners
    await start_all_listeners()
    logger.info("✓ Device event listeners started")

    # Register health
    from backend.system.status import register_subsystem
    register_subsystem("devices", "healthy")

except Exception as e:
    logger.warning(f"Failed to start device system: {e}")
```

### **NATS Topics**

```
onemind.devices.*                          # All device events
onemind.devices.camera.*                   # All camera events
onemind.devices.camera.person_detected     # Person detection
onemind.devices.camera.motion_detected     # Motion detection
onemind.devices.sensor.*                   # All sensor events
onemind.devices.sensor.temperature_high    # Temperature alert
onemind.devices.sensor.humidity_low        # Humidity alert
onemind.devices.usb.*                      # USB device events
onemind.devices.usb.hotplug                # USB device connected
onemind.devices.home.*                     # Smart home events
onemind.devices.home.door_opened           # Door sensor
onemind.devices.home.light_turned_on       # Light switch
```

### **Use Cases**

1. **Security Camera**: MaixCAM detects person → Agent decides if alert needed
2. **Smart Home**: Door opens when no one home → Agent sends notification
3. **Temperature**: Sensor detects high temp → Agent adjusts thermostat
4. **OM1 Robot**: Vision system detects obstacle → Agent stops motors
5. **USB Devices**: New camera plugged in → Agent auto-registers it

### **OM1 Integration**

This architecture is **exactly** what OM1 needs:

```python
# backend/openmind/devices/om1_sensors.py

class OM1VisionSource(DeviceSource):
    """OM1 robot vision system as device source"""

    async def _monitoring_loop(self):
        while self.running:
            # Get vision data from ROS2
            detections = await self.ros2_client.get_detections()

            # Publish to NATS
            for det in detections:
                event = DeviceEvent(
                    event_type="obstacle_detected",
                    device_id="om1_camera",
                    device_type="camera",
                    data=det
                )
                await publish(event.to_nats_topic(), event.to_dict())
```

---

## 3. Voice Pipeline (Full) 🎙️

### **Objective**
Enable complete voice interaction loop:
- Wake word detection ("Hey OneMind")
- Speech-to-text (user input)
- Agent processing (with voice-optimized prompts)
- Text-to-speech (agent response)
- Audio streaming (WebSocket for remote devices)

### **Agno Status**
⚠️ **PARTIAL**
- ✅ Text-to-speech: Gemini native audio (5 voices)
- ❌ Speech-to-text: Not provided
- ❌ Wake word detection: Not provided
- ❌ Audio streaming: Not provided

### **Architecture**

```
┌─────────────────────────────────────────────────────────┐
│              Device Microphone (Always On)              │
│  - Continuous audio capture                             │
│  - Local buffering (circular buffer)                    │
└─────────────────────────────────────────────────────────┘
                    ↓ Audio stream
┌─────────────────────────────────────────────────────────┐
│         Wake Word Detector (Porcupine/Vosk)             │
│  - Offline processing (privacy-preserving)              │
│  - Custom wake words: "Hey OneMind", "Hey OM1"         │
│  - Low latency (~100ms)                                 │
└─────────────────────────────────────────────────────────┘
                    ↓ Wake word detected
┌─────────────────────────────────────────────────────────┐
│           Audio Capture (3-5 second window)             │
│  - VAD (Voice Activity Detection)                       │
│  - Stop on silence                                      │
└─────────────────────────────────────────────────────────┘
                    ↓ Audio buffer
┌─────────────────────────────────────────────────────────┐
│        Speech-to-Text (Groq Whisper / Google STT)       │
│  - Transcription: "What's the weather today?"           │
│  - Latency: <500ms                                      │
│  - Accuracy: 95%+                                       │
└─────────────────────────────────────────────────────────┘
                    ↓ Text query
┌─────────────────────────────────────────────────────────┐
│         Agno Agent (Voice-Optimized)                    │
│  - Voice-specific instructions                          │
│  - Short, conversational responses                      │
│  - Tool execution (weather, etc.)                       │
└─────────────────────────────────────────────────────────┘
                    ↓ Text response
┌─────────────────────────────────────────────────────────┐
│      Text-to-Speech (Gemini Native Audio)               │
│  - Agno already supports this!                          │
│  - 5 voice options (Kore, Charon, Aoede, etc.)         │
│  - Latency: <1s                                         │
└─────────────────────────────────────────────────────────┘
                    ↓ Audio stream (PCM16)
┌─────────────────────────────────────────────────────────┐
│         Audio Playback / WebSocket Streaming            │
│  - Local: pyaudio playback                              │
│  - Remote: WebSocket streaming to device                │
└─────────────────────────────────────────────────────────┘
```

### **Implementation Plan**

#### **Phase 1: Wake Word Detection** (Day 1)

**File**: `backend/voice/wake_word.py` (NEW)

```python
"""
Wake word detection using Porcupine
Offline, privacy-preserving, low latency
"""

import pvporcupine
import pyaudio
import struct
import logging
from typing import Callable, Optional

logger = logging.getLogger(__name__)


class WakeWordDetector:
    """
    Offline wake word detection using Porcupine
    """

    def __init__(
        self,
        access_key: str,
        keywords: list[str] = None,
        sensitivities: list[float] = None,
    ):
        """
        Initialize wake word detector

        Args:
            access_key: Picovoice access key (free tier available)
            keywords: List of wake words (default: ["computer"])
            sensitivities: Detection sensitivity 0-1 (default: 0.5)
        """
        self.keywords = keywords or ["computer"]
        self.sensitivities = sensitivities or [0.5] * len(self.keywords)

        # Initialize Porcupine
        self.porcupine = pvporcupine.create(
            access_key=access_key,
            keywords=self.keywords,
            sensitivities=self.sensitivities
        )

        # Audio setup
        self.audio = pyaudio.PyAudio()
        self.stream = None
        self.running = False

        logger.info(f"Wake word detector initialized: {self.keywords}")

    async def start(self, on_wake: Callable):
        """
        Start listening for wake word

        Args:
            on_wake: Callback function when wake word detected
        """
        self.running = True

        self.stream = self.audio.open(
            rate=self.porcupine.sample_rate,
            channels=1,
            format=pyaudio.paInt16,
            input=True,
            frames_per_buffer=self.porcupine.frame_length,
            stream_callback=self._audio_callback
        )

        self.stream.start_stream()
        logger.info("Listening for wake word...")

        # Keep running
        while self.running:
            await asyncio.sleep(0.1)

    def stop(self):
        """Stop listening"""
        self.running = False
        if self.stream:
            self.stream.stop_stream()
            self.stream.close()
        self.porcupine.delete()
        self.audio.terminate()
        logger.info("Wake word detector stopped")

    def _audio_callback(self, in_data, frame_count, time_info, status):
        """Process audio frames"""
        pcm = struct.unpack_from("h" * self.porcupine.frame_length, in_data)
        keyword_index = self.porcupine.process(pcm)

        if keyword_index >= 0:
            logger.info(f"Wake word detected: {self.keywords[keyword_index]}")
            # Trigger callback in main thread
            asyncio.create_task(self.on_wake_callback())

        return (in_data, pyaudio.paContinue)
```

#### **Phase 2: Speech-to-Text** (Day 2)

**File**: `backend/voice/stt.py` (NEW)

```python
"""
Speech-to-Text using Groq Whisper
Fast, accurate, cost-effective
"""

import asyncio
import pyaudio
import wave
import tempfile
import logging
from groq import Groq

logger = logging.getLogger(__name__)


class SpeechToText:
    """
    Speech-to-Text using Groq Whisper API
    """

    def __init__(self, api_key: str):
        self.client = Groq(api_key=api_key)
        self.audio = pyaudio.PyAudio()
        logger.info("Speech-to-text initialized (Groq Whisper)")

    async def capture_and_transcribe(
        self,
        duration_seconds: int = 5,
        silence_threshold: int = 500,
        silence_duration: float = 1.5
    ) -> str:
        """
        Capture audio from microphone and transcribe

        Args:
            duration_seconds: Maximum recording duration
            silence_threshold: Audio level below which is silence
            silence_duration: Seconds of silence to stop recording

        Returns:
            Transcribed text
        """
        # Capture audio
        audio_file = await self._capture_audio(
            duration_seconds,
            silence_threshold,
            silence_duration
        )

        # Transcribe
        text = await self._transcribe(audio_file)

        return text

    async def _capture_audio(
        self,
        max_duration: int,
        silence_threshold: int,
        silence_duration: float
    ) -> str:
        """Capture audio from microphone"""
        CHUNK = 1024
        FORMAT = pyaudio.paInt16
        CHANNELS = 1
        RATE = 16000

        stream = self.audio.open(
            format=FORMAT,
            channels=CHANNELS,
            rate=RATE,
            input=True,
            frames_per_buffer=CHUNK
        )

        logger.info("Recording... (speak now)")

        frames = []
        silent_chunks = 0
        max_silent_chunks = int((silence_duration * RATE) / CHUNK)

        for i in range(0, int(RATE / CHUNK * max_duration)):
            data = stream.read(CHUNK)
            frames.append(data)

            # Check for silence
            audio_data = struct.unpack(str(CHUNK) + 'h', data)
            volume = max(audio_data)

            if volume < silence_threshold:
                silent_chunks += 1
                if silent_chunks > max_silent_chunks:
                    logger.info("Silence detected, stopping...")
                    break
            else:
                silent_chunks = 0

        stream.stop_stream()
        stream.close()

        # Save to temporary file
        temp_file = tempfile.NamedTemporaryFile(suffix=".wav", delete=False)
        wf = wave.open(temp_file.name, 'wb')
        wf.setnchannels(CHANNELS)
        wf.setsampwidth(self.audio.get_sample_size(FORMAT))
        wf.setframerate(RATE)
        wf.writeframes(b''.join(frames))
        wf.close()

        logger.info(f"Audio captured: {temp_file.name}")
        return temp_file.name

    async def _transcribe(self, audio_file: str) -> str:
        """Transcribe audio file using Groq Whisper"""
        try:
            with open(audio_file, "rb") as file:
                transcription = self.client.audio.transcriptions.create(
                    file=(audio_file, file.read()),
                    model="whisper-large-v3",
                    response_format="text",
                    language="en",
                    temperature=0.0
                )

            text = transcription.strip()
            logger.info(f"Transcribed: {text}")
            return text

        except Exception as e:
            logger.error(f"Transcription error: {e}")
            return ""
```

#### **Phase 3: Voice Gateway** (Day 3)

**File**: `backend/voice/gateway.py` (NEW)

```python
"""
Voice Gateway - Complete voice interaction loop
Integrates wake word, STT, agent, and TTS
"""

import asyncio
import logging
from backend.voice.wake_word import WakeWordDetector
from backend.voice.stt import SpeechToText
from backend.agents.factory import create_agent_from_config
from backend.agents.presets import AGENT_PRESETS

logger = logging.getLogger(__name__)


class VoiceGateway:
    """
    Complete voice interaction pipeline
    """

    def __init__(
        self,
        porcupine_key: str,
        groq_api_key: str,
        agent_id: str = "assistant",
        wake_words: list[str] = None,
    ):
        # Wake word detection
        self.wake_detector = WakeWordDetector(
            access_key=porcupine_key,
            keywords=wake_words or ["hey onemind"],
            sensitivities=[0.5]
        )

        # Speech-to-text
        self.stt = SpeechToText(api_key=groq_api_key)

        # Create voice-enabled agent
        agent_preset = AGENT_PRESETS[agent_id].copy()
        self.agent = create_agent_from_config(
            config=agent_preset,
            voice_enabled=True,      # Gemini native audio
            voice_id="kore",         # Female, warm voice
            streaming=True,          # Enable audio streaming
        )

        logger.info("Voice gateway initialized")

    async def start(self):
        """Start voice interaction loop"""
        logger.info("Voice gateway started. Say wake word to begin...")

        # Start listening for wake word
        await self.wake_detector.start(on_wake=self._on_wake_word)

    async def _on_wake_word(self):
        """Handle wake word detection"""
        logger.info("Wake word detected! Listening...")

        # Capture user speech
        text = await self.stt.capture_and_transcribe()

        if not text:
            logger.info("No speech detected")
            return

        logger.info(f"User said: {text}")

        # Send to agent
        response = await self.agent.run(text)

        # Play audio response (Gemini native audio)
        if hasattr(response, 'audio_data') and response.audio_data:
            await self._play_audio(response.audio_data)
        else:
            logger.warning("No audio in response")

    async def _play_audio(self, audio_data: bytes):
        """Play audio through speakers"""
        import pyaudio

        p = pyaudio.PyAudio()
        stream = p.open(
            format=pyaudio.paInt16,
            channels=1,
            rate=24000,  # Gemini audio sample rate
            output=True
        )

        stream.write(audio_data)
        stream.stop_stream()
        stream.close()
        p.terminate()

        logger.info("Audio playback complete")
```

#### **Phase 4: WebSocket Audio Streaming** (Day 4)

**File**: `backend/voice/websocket.py` (NEW)

```python
"""
WebSocket audio streaming for remote devices
Enables voice interaction from browser/mobile
"""

from fastapi import WebSocket
import asyncio
import logging

logger = logging.getLogger(__name__)


class VoiceWebSocketHandler:
    """
    Handle bidirectional audio streaming over WebSocket
    """

    def __init__(self, websocket: WebSocket):
        self.websocket = websocket
        self.stt = SpeechToText(api_key=os.getenv("GROQ_API_KEY"))
        self.agent = None  # Initialize on connection

    async def handle(self):
        """Main WebSocket handler"""
        await self.websocket.accept()
        logger.info("Voice WebSocket connected")

        try:
            while True:
                # Receive audio from client
                audio_data = await self.websocket.receive_bytes()

                # Save to temp file for transcription
                temp_file = self._save_audio(audio_data)

                # Transcribe
                text = await self.stt._transcribe(temp_file)

                if text:
                    # Send transcription to client
                    await self.websocket.send_json({
                        "type": "transcription",
                        "text": text
                    })

                    # Get agent response
                    response = await self.agent.run(text)

                    # Stream audio response back to client
                    if hasattr(response, 'audio_data'):
                        await self.websocket.send_bytes(response.audio_data)

        except Exception as e:
            logger.error(f"WebSocket error: {e}")
        finally:
            await self.websocket.close()


# Add to main.py
@app.websocket("/voice")
async def voice_websocket(websocket: WebSocket):
    """Voice interaction via WebSocket"""
    handler = VoiceWebSocketHandler(websocket)
    await handler.handle()
```

#### **Phase 5: Integration** (Day 5)

**File**: `backend/main.py` (MODIFY)

```python
# Add to startup_event()

# ─── 10. Start Voice Gateway (Optional) ───
if os.getenv("VOICE_ENABLED") == "true":
    try:
        from backend.voice.gateway import VoiceGateway

        voice_gateway = VoiceGateway(
            porcupine_key=os.getenv("PORCUPINE_ACCESS_KEY"),
            groq_api_key=os.getenv("GROQ_API_KEY"),
            agent_id="assistant",
            wake_words=["hey onemind", "hey om one"]
        )

        # Start in background
        asyncio.create_task(voice_gateway.start())
        logger.info("✓ Voice gateway started (wake word active)")

        # Register health
        from backend.system.status import register_subsystem
        register_subsystem("voice", "healthy")

    except Exception as e:
        logger.warning(f"Failed to start voice gateway: {e}")
```

### **Environment Configuration**

```bash
# .env additions
VOICE_ENABLED=true
PORCUPINE_ACCESS_KEY=your_picovoice_key  # Free tier: https://picovoice.ai
GROQ_API_KEY=your_groq_key               # Free tier: https://groq.com
```

### **Dependencies**

```
# requirements.txt additions
pvporcupine==3.0.2         # Wake word detection
pyaudio==0.2.14            # Audio I/O
groq==0.4.2                # Fast Whisper STT
```

### **Use Cases**

1. **Desktop**: "Hey OneMind, what's the weather?" → Voice response
2. **OM1 Robot**: "Hey OM1, move forward" → Motor control
3. **Smart Home**: "Hey OneMind, turn off the lights" → Home Assistant
4. **Remote**: Browser WebSocket → Voice interaction from phone
5. **Continuous**: Always-listening mode for hands-free operation

### **Latency Breakdown**

| Stage | Latency | Notes |
|-------|---------|-------|
| Wake Word | ~100ms | Offline (Porcupine) |
| Audio Capture | 1-3s | Variable (until silence) |
| STT (Groq) | ~300ms | Fast Whisper inference |
| Agent Processing | 1-2s | Depends on tool calls |
| TTS (Gemini) | ~800ms | Native audio generation |
| **Total** | **3-6s** | End-to-end |

---

## Testing Strategy

### Browser Automation Tests
```bash
pytest tests/test_browser_tools.py -v
```

### Device Event Tests
```bash
pytest tests/test_device_events.py -v
```

### Voice Pipeline Tests
```bash
pytest tests/test_voice_pipeline.py -v
```

---

## Deployment Checklist

### Browser Automation
- [ ] Install Playwright: `pip install playwright`
- [ ] Install browsers: `playwright install chromium`
- [ ] Add browser tools to registry
- [ ] Configure web_automator agent
- [ ] Test with example.com
- [ ] Document use cases

### Device Event System
- [ ] Create device framework
- [ ] Implement camera source
- [ ] Set up NATS topics
- [ ] Create event listeners
- [ ] Integrate with main.py
- [ ] Test with mock device
- [ ] Document event schema

### Voice Pipeline
- [ ] Get Picovoice key (free tier)
- [ ] Get Groq API key (free tier)
- [ ] Install audio dependencies
- [ ] Configure wake word detector
- [ ] Test STT accuracy
- [ ] Verify Gemini audio output
- [ ] Test WebSocket streaming
- [ ] Document voice commands

---

## Success Metrics

### Browser Automation
- ✅ Can extract data from 10 different websites
- ✅ Form filling works for 5 common use cases
- ✅ Screenshot capture <2s latency
- ✅ No browser crashes after 100 operations

### Device Event System
- ✅ Events published to NATS <100ms
- ✅ Multiple agents can subscribe to same event
- ✅ Device failures don't crash agents
- ✅ Full event logging for debugging

### Voice Pipeline
- ✅ Wake word detection <100ms latency
- ✅ STT accuracy >95%
- ✅ End-to-end latency <6s
- ✅ TTS sounds natural (user subjective test)
- ✅ WebSocket streaming works from browser

---

## Future Enhancements

### Browser Automation (Phase 2)
- Cookie persistence across sessions
- Browser automation recording/replay
- Multi-page workflows
- Network request interception

### Device Event System (Phase 2)
- Device discovery (auto-detect cameras on network)
- Event aggregation (combine multiple sensors)
- Event replay (debugging)
- Device firmware updates via NATS

### Voice Pipeline (Phase 2)
- Multi-language support
- Voice cloning (custom voices)
- Emotion detection in speech
- Background noise cancellation
- Conversation memory (context across sessions)

---

## Resources

### Documentation
- **Agno AgentOS**: https://github.com/agno-oss/agno
- **Playwright**: https://playwright.dev/python/
- **Porcupine**: https://picovoice.ai/docs/porcupine/
- **Groq Whisper**: https://groq.com/docs/

### Code References
- **OpenClaw**: https://github.com/openclaw/openclaw
- **PicoClaw**: https://github.com/sipeed/picoclaw
- **Nanobot**: https://github.com/HKUDS/nanobot

---

## Questions & Answers

**Q: Why not use OpenClaw's gateway for chat platforms?**
A: OneMindOS has a Flutter UI, which provides better UX than chat for complex workflows.

**Q: Why Groq instead of OpenAI Whisper?**
A: Groq provides faster inference (300ms vs 1-2s) and has generous free tier.

**Q: Can device events trigger workflows?**
A: Yes! Nexus can listen to device events and trigger Agno workflows automatically.

**Q: Is voice pipeline production-ready?**
A: Not yet. Needs extensive testing for edge cases (accents, noise, latency).

---

**Document Version**: 1.0
**Last Updated**: 2025-02-17
**Maintained By**: OneMindOS Core Team

---

## Future Canvas Ideas

- **NATS Pulse Visualization** — OpenClaw canvas screen that renders live NATS message flow as animated particle effects, showing agent activity and system health as a living visual (think: each message = a spark traveling between nodes).
