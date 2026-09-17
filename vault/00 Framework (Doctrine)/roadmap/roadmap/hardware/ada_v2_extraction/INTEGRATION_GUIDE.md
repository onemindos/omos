# Integration Guide: ADA V2 → OneMind OS

This guide walks you through integrating the extracted ADA V2 capabilities into your OneMind OS codebase.

---

## Quick Start

### Step 1: Install Dependencies

```bash
# CAD Generation
pip install build123d cadquery-ocp

# 3D Printer Control
pip install zeroconf aiohttp

# Face Authentication
pip install mediapipe opencv-python numpy

# Web Agent (if using Gemini Computer Use)
pip install playwright google-genai
playwright install chromium
```

### Step 2: Copy to Backend

```bash
# From OneMind OS root:
mkdir -p backend/agno/tools/fabrication
mkdir -p backend/agno/multimodal

# Copy agents
cp ada_v2_extraction/backend/cad_agent.py backend/agno/tools/fabrication/
cp ada_v2_extraction/backend/printer_agent.py backend/agno/tools/fabrication/
cp ada_v2_extraction/backend/authenticator.py backend/agno/multimodal/

# Copy tool wrappers
cp ada_v2_extraction/tools/cad_tools.py backend/agno/tools/fabrication/
cp ada_v2_extraction/tools/printer_tools.py backend/agno/tools/fabrication/
cp ada_v2_extraction/tools/face_auth_tools.py backend/agno/multimodal/
```

### Step 3: Register Tools

Edit `backend/agno/tools/__init__.py` or your tool registry:

```python
# Import the new tools
from .fabrication.cad_tools import get_cad_tools
from .fabrication.printer_tools import get_printer_tools
from ..multimodal.face_auth_tools import get_face_auth_tools

# Add to your tool registry
def get_all_tools():
    tools = []
    tools.extend(get_cad_tools())
    tools.extend(get_printer_tools())
    tools.extend(get_face_auth_tools())
    # ... your existing tools
    return tools
```

---

## Feature-Specific Integration

### 🔧 3D CAD Generation

**Location:** `backend/agno/tools/fabrication/`

**Adapt LLM Client (Required):**

The CAD agent uses Gemini. For OneMind's Claude/Bedrock setup:

```python
# In cad_agent.py, modify __init__ or add a factory function:

from backend.agno.ai.bedrock import get_claude_client

def create_onemind_cad_agent():
    """Factory for CAD agent using OneMind's Claude setup."""
    claude_client = get_claude_client()

    agent = CadAgent(llm_client=claude_client)

    # Override the code generation method for Claude
    agent._generate_code_with_custom_llm = _generate_with_claude

    return agent

async def _generate_with_claude(prompt: str) -> str:
    """Generate build123d code using Claude."""
    response = await claude_client.messages.create(
        model="claude-opus-4-0-20250514",
        max_tokens=4096,
        system=CadAgent.SYSTEM_INSTRUCTION,
        messages=[{"role": "user", "content": prompt}]
    )
    return response.content[0].text
```

**Example Usage:**

```python
# In Legacy agent or a dedicated CAD agent
@tool
def design_component(description: str):
    """Design a 3D component from description."""
    from backend.agno.tools.fabrication.cad_tools import generate_3d_model
    return generate_3d_model(description, project_name="robotics")
```

---

### 🖨️ 3D Printer Control

**Location:** `backend/agno/tools/fabrication/`

**No LLM Required** - This is pure Python automation.

**Setup:**

1. Install OrcaSlicer or PrusaSlicer on the system
2. Ensure printers are on the same network

**Example Workflow:**

```python
# Discover printers on network
printers = discover_3d_printers()

# Generate a part
result = generate_3d_model("mounting bracket for camera")
stl_path = result["file_path"]

# Slice and print
print_result = print_model(
    printer="Creality K1",
    stl_path=stl_path
)
```

**NATS Integration:**

For real-time print status updates:

```python
async def publish_print_status(printer_name: str):
    agent = PrinterAgent()
    while True:
        status = await agent.get_print_status(printer_name)
        if status:
            await nats_client.publish(
                f"printer.{printer_name}.status",
                json.dumps(status.to_dict())
            )
        await asyncio.sleep(5)
```

---

### 👤 Face Authentication

**Location:** `backend/agno/multimodal/`

**Use Cases:**

- Confirm identity before payments
- Admin action confirmations
- HITL verification gates

**Example: Protected Operation**

```python
@tool
def transfer_funds(amount: float, to_account: str):
    """Transfer funds with face verification."""

    # Require face verification for amounts over $100
    if amount > 100:
        verify_result = verify_face_identity(timeout_seconds=30)
        result = json.loads(verify_result)

        if not result["verified"]:
            return json.dumps({
                "status": "denied",
                "message": "Face verification failed. Transfer cancelled."
            })

    # Proceed with transfer
    return execute_transfer(amount, to_account)
```

**Mobile/Remote Verification:**

```python
# For Flutter app integration
@app.post("/api/verify-face")
async def verify_face_endpoint(image_b64: str):
    result = verify_face_from_image(image_b64)
    return json.loads(result)
```

---

### 🌐 Web Agent (Optional)

**Consideration:** The Web Agent uses Gemini's Computer Use model, which is unique to Google's offering. You have two options:

**Option A: Keep Gemini for Computer Use Only**

```python
# backend/agno/tools/browser_visual.py
# Use for visual web tasks, keep Claude for everything else

from ada_v2_extraction.backend.web_agent import WebAgent

@tool
def visual_web_task(task: str, start_url: str = "https://google.com"):
    """Perform visual web tasks (clicking, navigating visually)."""
    agent = WebAgent(headless=True)
    result = _run_async(agent.run_task(task, start_url=start_url))
    return result
```

**Option B: Skip It**

Your existing `browser.py` handles extraction well. The Computer Use model is primarily for tasks where you need to:

- Navigate sites without stable selectors
- Handle dynamically rendered content
- Perform multi-step visual workflows

If your current browser tool meets your needs, skip this integration.

---

## Architecture Diagram

```
┌─────────────────────────────────────────────────────────────────┐
│                         OneMind OS                               │
├─────────────────────────────────────────────────────────────────┤
│  Flutter UI                                                      │
│  └── 3D Viewer Component (for STL preview)                      │
├─────────────────────────────────────────────────────────────────┤
│  NATS JetStream                                                  │
│  ├── agent.cad.request                                          │
│  ├── agent.cad.status                                           │
│  ├── printer.{name}.status                                      │
│  └── auth.face.verify                                           │
├─────────────────────────────────────────────────────────────────┤
│  Python Backend (FastAPI + Agno)                                │
│  ├── Legacy Agent                                               │
│  │   └── Tools                                                  │
│  │       ├── fabrication/                                       │
│  │       │   ├── cad_agent.py ──────► build123d                │
│  │       │   ├── cad_tools.py                                   │
│  │       │   ├── printer_agent.py ──► OrcaSlicer CLI           │
│  │       │   └── printer_tools.py     └──► Moonraker/OctoPrint │
│  │       └── browser.py (existing)                              │
│  │                                                              │
│  └── Multimodal                                                 │
│      ├── authenticator.py ──────────► MediaPipe                 │
│      └── face_auth_tools.py                                     │
└─────────────────────────────────────────────────────────────────┘
```

---

## Testing the Integration

### Test CAD Generation

```bash
cd backend
python -c "
from ada_v2_extraction.backend.cad_agent import CadAgent
import asyncio

async def test():
    agent = CadAgent()
    result = await agent.generate_prototype('a simple cube 20mm on each side')
    print(f'Result: {result}')

asyncio.run(test())
"
```

### Test Printer Discovery

```bash
cd backend
python -c "
from ada_v2_extraction.backend.printer_agent import PrinterAgent
import asyncio

async def test():
    agent = PrinterAgent()
    printers = await agent.discover_printers()
    print(f'Found printers: {printers}')

asyncio.run(test())
"
```

### Test Face Auth

```bash
cd backend
python ada_v2_extraction/backend/authenticator.py
# Choose option 1 to enroll, option 2 to test
```

---

## What NOT to Integrate

1. **Socket.IO patterns** - Replace with NATS publishing
2. **Project Manager** - Use OneMind's existing session/memory systems
3. **Kasa Smart Home** - You have Home Assistant which is more comprehensive
4. **Electron/React patterns** - You have Flutter

---

## Files Reference

| File                 | Purpose                            | Priority       |
| -------------------- | ---------------------------------- | -------------- |
| `cad_agent.py`       | 3D model generation                | High           |
| `cad_tools.py`       | Agno tool wrappers for CAD         | High           |
| `printer_agent.py`   | Printer discovery/slicing/printing | Medium         |
| `printer_tools.py`   | Agno tool wrappers for printing    | Medium         |
| `authenticator.py`   | Face recognition                   | Medium         |
| `face_auth_tools.py` | Agno tool wrappers for face auth   | Medium         |
| `web_agent.py`       | Visual browser automation          | Low (optional) |

---

## Next Steps

1. [ ] Adapt CadAgent to use Claude instead of Gemini
2. [ ] Add 3D model preview to Flutter UI
3. [ ] Create NATS topics for CAD/printer status
4. [ ] Test end-to-end: "Legacy, design and print a phone stand"
5. [ ] Consider gesture control for future integration
