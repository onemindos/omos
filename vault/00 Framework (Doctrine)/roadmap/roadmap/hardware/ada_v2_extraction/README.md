# ADA V2 → OneMind OS Extraction

## 🔗 Source

| Item                | Link                                                                 |
| ------------------- | -------------------------------------------------------------------- |
| **Repository**      | [github.com/nazirlouis/ada_v2](https://github.com/nazirlouis/ada_v2) |
| **License**         | MIT                                                                  |
| **Extraction Date** | 2026-01-30                                                           |
| **Creator**         | Nazir Louis                                                          |

---

## 📋 Overview

This folder contains code and documentation extracted from **A.D.A V2** (Advanced Design Assistant) for potential integration into OneMind OS. ADA V2 is a desktop AI assistant focused on 3D CAD design, printing, and maker workflows.

---

## 🔍 Feature Comparison Matrix

| Feature                  | ADA V2                           | OneMind OS                    | Action                          |
| ------------------------ | -------------------------------- | ----------------------------- | ------------------------------- |
| **Voice**                | Gemini Native Audio              | LiveKit + Deepgram + Cartesia | ✅ OneMind is more robust       |
| **Smart Home**           | TP-Link Kasa only                | Home Assistant (15+ devices)  | ✅ OneMind wins                 |
| **3D CAD**               | build123d + OrcaSlicer           | ❌ Not present                | 🆕 **Port from ADA**            |
| **Gesture Control**      | MediaPipe Hands                  | ❌ Not present                | 🆕 **Port from ADA**            |
| **Web Automation**       | Playwright + Gemini Computer Use | Playwright (browser.py)       | 🔄 **Compare implementations**  |
| **Face Auth**            | MediaPipe Face Landmarker        | ❌ Not present                | ⚖️ **Consider for local auth**  |
| **3D Printer Discovery** | mDNS + OctoPrint/Moonraker       | ❌ Not present                | 🆕 **Port from ADA**            |
| **Memory**               | Project-based                    | Persistent Vector Memory      | ✅ OneMind wins                 |
| **Multi-Agent**          | Single agent                     | Multi-agent teams             | ✅ OneMind wins                 |
| **Platform**             | Electron desktop only            | Flutter (6 platforms)         | ✅ OneMind wins                 |
| **Messaging**            | Socket.IO                        | NATS JetStream                | ✅ OneMind wins (lower latency) |

---

## 📁 Extracted Code Structure

```
ada_v2_extraction/
├── README.md                     # This file
├── INTEGRATION_GUIDE.md          # How to integrate into OneMind OS
├── backend/
│   ├── cad_agent.py              # 3D CAD generation with build123d
│   ├── printer_agent.py          # 3D printer discovery & slicing
│   ├── web_agent.py              # Browser automation with Gemini Computer Use
│   ├── authenticator.py          # Face authentication with MediaPipe
│   ├── kasa_agent.py             # TP-Link Kasa smart device control ✨ NEW
│   ├── project_manager.py        # Project workspace management ✨ NEW
│   ├── hand_gesture_test.py      # Reference gesture detection ✨ NEW
│   ├── capture_face.py           # Utility to capture reference face
│   └── tools.py                  # Tool definitions for Gemini
├── tools/                        # OneMind-adapted tool wrappers
│   ├── cad_tools.py              # Agno-compatible CAD tools
│   ├── printer_tools.py          # Agno-compatible printer tools
│   └── face_auth_tools.py        # Agno-compatible face auth tools
└── frontend/
    └── gesture_notes.md          # Notes on gesture implementation
```

---

## 🎯 Priority Integration Order

### **Priority 1: 3D CAD Generation (cad_agent.py)**

Generate 3D models from natural language prompts using build123d.

**Why:** Perfect for hardware/robotics projects in your `hardware/` directory.

**Example Use Cases:**

- "Legacy, create a mounting bracket for my Raspberry Pi camera"
- "Design a case for this sensor with ventilation holes"
- "Make a gear with 24 teeth, module 1.5"

**Dependencies:**

```bash
pip install build123d cadquery-ocp
```

---

### **Priority 2: 3D Printer Control (printer_agent.py)**

Discover printers on network, slice STL files, and submit print jobs.

**Why:** End-to-end fabrication workflow from design to physical part.

**Features:**

- mDNS printer discovery (OctoPrint, Moonraker, Klipper)
- OrcaSlicer/PrusaSlicer CLI integration
- Auto-profile detection based on printer name
- Upload and start prints via REST API
- Print status monitoring

**Dependencies:**

```bash
pip install zeroconf aiohttp
# Also requires OrcaSlicer or PrusaSlicer installed
```

---

### **Priority 3: Face Authentication (authenticator.py)**

Local biometric authentication using MediaPipe Face Landmarker.

**Why:** Adds local identity verification for sensitive operations without cloud dependencies.

**Features:**

- Cosine similarity-based face matching
- Reference face enrollment
- Real-time webcam processing
- Async-compatible with asyncio

**Dependencies:**

```bash
pip install mediapipe opencv-python numpy
```

---

### **Priority 4: Browser Agent Comparison (web_agent.py)**

Enhanced browser automation using Gemini Computer Use model.

**What ADA Has That OneMind Could Use:**

1. **Gemini Computer Use Model** (`gemini-2.5-computer-use-preview-10-2025`)
2. **Normalized coordinate system** (1000x1000 → actual viewport)
3. **Built-in safety decision handling**
4. **Streaming screenshot feedback**

**OneMind Already Has:**

- Playwright integration in `browser.py`
- Extract types (links, text, tables, etc.)
- Search engine integration

**Recommendation:** Consider adding Computer Use capabilities for visual web tasks while keeping your extraction-focused approach for data tasks.

---

## 🏗️ Architecture Differences

### **ADA V2 Architecture**

```
┌─────────────────────────────────────────────────────────┐
│  Electron + React Frontend                              │
│  ├── Three.js 3D Viewer                                │
│  ├── Socket.IO Client                                  │
│  └── MediaPipe Gesture Detection                       │
└─────────────────────────────────────────────────────────┘
                          │
                          ▼
┌─────────────────────────────────────────────────────────┐
│  Python Backend (FastAPI + Socket.IO)                   │
│  ├── ada.py (Gemini Live API)                          │
│  ├── cad_agent.py (build123d)                          │
│  ├── printer_agent.py (OrcaSlicer + Moonraker)         │
│  ├── web_agent.py (Playwright + Computer Use)          │
│  └── authenticator.py (MediaPipe Face)                 │
└─────────────────────────────────────────────────────────┘
```

### **OneMind OS Architecture**

```
┌─────────────────────────────────────────────────────────┐
│  Flutter Multi-Platform UI                              │
│  ├── iOS, Android, Web, macOS, Windows, Linux          │
│  └── WebSocket + REST API Client                       │
└─────────────────────────────────────────────────────────┘
                          │
                          ▼
┌─────────────────────────────────────────────────────────┐
│  NATS JetStream (Nervous System)                        │
│  ├── Event Streaming (200-400μs)                       │
│  ├── MQTT Bridge (IoT)                                 │
│  └── WebSocket Bridge (UI)                             │
└─────────────────────────────────────────────────────────┘
                          │
                          ▼
┌─────────────────────────────────────────────────────────┐
│  Python Backend (FastAPI + Agno Framework)              │
│  ├── Multi-Agent Teams (Legacy, Researcher, etc.)      │
│  ├── Persistent Memory (PostgreSQL + pgvector)         │
│  ├── Knowledge Base (LanceDB + RAG)                    │
│  ├── 40+ Tools (Todoist, GitHub, Home Assistant, etc.) │
│  └── Voice (LiveKit + Deepgram + Cartesia)             │
└─────────────────────────────────────────────────────────┘
```

---

## 📦 Dependencies to Add

Add these to your requirements when integrating:

```txt
# 3D CAD Generation
build123d>=0.1.0
cadquery-ocp>=7.7.0

# 3D Printer Control
zeroconf>=0.131.0
aiohttp>=3.9.0

# Face Authentication
mediapipe>=0.10.0
opencv-python>=4.8.0
numpy>=1.24.0
```

---

## 🚀 Quick Integration Steps

### Step 1: Copy Backend Agents

```bash
# From onemind-os root:
cp ada_v2_extraction/backend/cad_agent.py backend/agno/tools/cad/
cp ada_v2_extraction/backend/printer_agent.py backend/agno/tools/hardware/
cp ada_v2_extraction/backend/authenticator.py backend/agno/multimodal/
```

### Step 2: Create Agno Tool Wrappers

Use the files in `ada_v2_extraction/tools/` as reference for creating Agno-compatible tool decorators.

### Step 3: Register Tools

Add to your tool registry in `backend/agno/tools/registry.py`

### Step 4: Install Dependencies

```bash
pip install build123d zeroconf aiohttp mediapipe opencv-python
```

---

## ⚠️ Adaptation Notes

### Port Considerations

1. **Gemini → Claude/Bedrock**: ADA uses Gemini API. You'll need to adapt:
   - `cad_agent.py` uses Gemini for code generation → adapt to use your existing Claude Opus setup
   - `web_agent.py` uses Gemini Computer Use → may need to keep Gemini for this specific use case or adapt

2. **Socket.IO → NATS**: ADA uses Socket.IO for real-time updates. Replace with NATS publishing:

   ```python
   # ADA style
   await sio.emit('cad_status', {'status': 'generating'})

   # OneMind style
   await nats_client.publish('agent.cad.status', json.dumps({'status': 'generating'}))
   ```

3. **Project Manager → Session/Memory**: ADA has simple project-based context. OneMind has full persistent memory - leverage it.

4. **Async Patterns**: ADA uses `asyncio.to_thread()` for blocking operations - same pattern works in OneMind.

---

## 📝 License

All extracted code is from [nazirlouis/ada_v2](https://github.com/nazirlouis/ada_v2) under MIT License.
This allows free use, modification, and distribution with attribution.

---

## 🔗 Related OneMind OS Locations

- **Backend Tools**: `/backend/agno/tools/`
- **Hardware Agents**: `/hardware/`
- **Multimodal Processing**: `/backend/agno/multimodal/`
- **Voice Interface**: `/backend/agno/voice/`
