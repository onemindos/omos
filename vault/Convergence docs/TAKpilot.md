# TAKpilot: AI Chat Copilot for CloudTAK

**Link:** https://corvusintell.com/takpilot/

**Vendor:** Corvus Intelligence (Ukrainian defense tech company)  
**Product Type:** AI-powered chat interface for CloudTAK mission command  
**Architecture:** Agentic AI with auditable tool calls + operator approval gates  
**License:** Commercial (available through Brave1/Delta Marketplace)  
**Target Users:** Tactical operators, commanders, ISR teams, NATO partners, Ukrainian Defense Forces (ZSU)  
**Status:** Production deployment  

---

## What It Is

An AI chat copilot that operates **inside CloudTAK** as an agent interface, enabling tactical operators to control the Common Operating Picture (COP) via natural language instead of menu navigation.

**Core Purpose:** Accelerate tactical decision-making by replacing multi-step menu workflows with single-sentence commands, while maintaining full transparency and operator control over AI actions.

**Tagline:** "AI Chat Copilot for CloudTAK"

---

## What It Can Do: Core Capabilities

### 1. **Chat-Native COP Control**
```
Operator types: "Place a marker at grid 12S DG 34567 89012, label it 'Enemy pos', color red"

TAKpilot:
  1. Parses natural language
  2. Queries CloudTAK for current context
  3. Places marker at specified grid
  4. Labels it correctly
  5. Sets color to red
  6. Displays tool call card for approval

Result: One sentence replaces 5+ menu clicks
```

**Supported operations:**
- Place markers/waypoints
- Draw geometry (polygons, polylines, circles)
- Create/update missions
- Manage overlays/layers
- Query current positions
- Generate reports
- Update COP objects

### 2. **Vision & File Intelligence**
```
Operator uploads SITREP PDF or recon sketch

TAKpilot:
  1. Reads file using built-in vision (PDF text extraction, image OCR)
  2. Extracts key information (grid coordinates, unit positions, threats)
  3. Converts unstructured intel → typed COP objects
  4. Georeferenced automatically (if coordinates provided)

Result: "Unstructured data becomes typed, georeferenced COP objects in one drag-drop"
```

**Supported file types:**
- PDF (SITREP, intelligence reports)
- Images (reconnaissance photos, sketches)
- Text files (hand-written transcripts)
- Overlays (tactical drawings)

### 3. **Streaming Tool Cards**
```
Every action TAKpilot takes appears as a live card:

Card shows:
  ✓ Action type (e.g., "Place marker")
  ✓ Parameters (location, label, color)
  ✓ Status (pending, approved, executing, completed)
  ✓ Timestamp

Operator can:
  - Watch in real-time
  - Inspect before approval
  - Deny action (cancel before execution)
  - See full audit trail
```

**Card types:**
- Read-only (queries, lookups) → Auto-execute
- Write operations (create, modify, delete) → Require approval
- Broadcast (network messages) → Require approval

### 4. **Approve/Deny Control**
```
Destructive operations pause for operator review:

Example: "Delete all old markers"
  → Card appears: "Delete 47 markers from date < 2 days ago?"
  → Shows: Target count, filters, markers to be deleted
  → Operator chooses: [Approve] or [Deny]
  → If approved: Executes with operator's identity
  → If denied: Cancels, no action taken

Result: Zero-surprise writes, full operator authority
```

**Approval gates on:**
- Create/delete markers
- Mission creation/modification
- Layer/overlay changes
- Broadcasts to team
- COP object modifications

### 5. **Model-Agnostic & Edge-Capable**
```
TAKpilot can run:
  - Cloud models (GPT-4, Claude, etc.)
  - Edge-local models (Llama, Mistral on local hardware)
  - Hybrid (cloud for reasoning, edge for decisions)

Operator chooses based on:
  - Network connectivity (no internet → use edge)
  - Operational security (classified data → local only)
  - Threat model (deny cloud access in contested areas)
  - Performance (edge = lower latency)

Result: No model lock-in, flexible deployment
```

**Supported model frameworks:**
- OpenAI API compatible
- Open-source LLMs (via Ollama, vLLM, etc.)
- Proprietary models
- Custom fine-tuned models

### 6. **MIL-STD-2525 / APP-6 Symbology**
```
All marker/unit symbols follow NATO standards:

Approved actions show:
  - Correct NATO icon (unit type, status, affiliation)
  - Proper symbology (equipment type, size, readiness)
  - Military standard abbreviations (friendly/enemy/neutral)

Result: COP markers match military standard iconography
        Compatible with allied systems (NATO partners)
        Interoperable with ATAK, WinTAK, other C2 systems
```

---

## Architecture: How TAKpilot Works

### System Overview
```
┌─────────────────────────────────────────────────────────┐
│                   CloudTAK (Web App)                    │
│                  ┌─────────────────┐                   │
│                  │  TAKpilot Chat  │ ← Operator here  │
│                  │    Panel (UI)   │                   │
│                  └────────┬────────┘                   │
└─────────────────────────┼─────────────────────────────┘
                          │
                   (REST API calls)
                          │
┌─────────────────────────┼─────────────────────────────┐
│      TAKpilot Agent (Per-Session Sandbox)             │
│  ┌────────────────────────┬──────────────────────┐   │
│  │   LLM (Cloud or Edge)  │  Vision Module       │   │
│  │   - Reads natural lang │  - OCR PDFs/images   │   │
│  │   - Plans tool calls   │  - Extract data      │   │
│  ├────────────────────────┼──────────────────────┤   │
│  │   Tool Execution       │  Audit Logger        │   │
│  │   - CloudTAK API calls │  - Every action      │   │
│  │   - Operator identity  │  - Operator ID       │   │
│  │   - Session scoped     │  - Timestamp         │   │
│  └────────────────────────┴──────────────────────┘   │
└──────────────────────────┬─────────────────────────┘
                           │
                   (Operator's token)
                           │
┌──────────────────────────┼─────────────────────────┐
│        CloudTAK Backend                            │
│  ├─ REST API                                      │
│  ├─ Mission database                              │
│  ├─ COP state                                     │
│  ├─ Audit logs                                    │
│  └─ User/session management                       │
└──────────────────────────────────────────────────┘
```

### The Four-Step Flow

**Step 1: Operator sends message or file**
```
Input: "Place 3 sniper positions at grid..."
Or: Upload SITREP PDF
Or: Drop tactical sketch image

Storage: Written to private session sandbox
         No cross-session access
         Never stored permanently (ephemeral)
```

**Step 2: Agent reads, reasons, and plans**
```
LLM processes:
  ✓ Natural language command
  ✓ Attached files (vision reading)
  ✓ Current COP state (query CloudTAK)
  ✓ Operator context (who am I? what team?)
  
Output: Plan of minimum tool calls needed

Example:
  Command: "Place 3 snipers at grid 12S DG"
  Plan:
    1. Create entity "Sniper 1" at grid location
    2. Set unit type = sniper (MIL-STD symbol)
    3. Set affiliation = friendly (color)
    4. Create entity "Sniper 2" at offset location
    5. Create entity "Sniper 3" at offset location
```

**Step 3: Tool calls stream as cards; destructive ones pause**
```
Read-only calls (auto-execute):
  - Query: "What markers exist in sector A?"
  - Lookup: "Show me all missions"
  - Analysis: "How many units in grid?"

Write calls (require approval):
  Card shows:
    ┌─────────────────────────────────────┐
    │ Action: Create marker               │
    │ Location: 12S DG 34567 89012        │
    │ Label: Sniper 1                     │
    │ Unit Type: [NATO sniper icon]       │
    │ Affiliation: Friendly               │
    │ Status: Pending operator approval   │
    │                                     │
    │ [Approve]        [Deny]             │
    └─────────────────────────────────────┘
```

**Step 4: Agent acts as operator's identity**
```
Approved calls:
  ✓ Use operator's CloudTAK session token
  ✓ Hit CloudTAK REST API
  ✓ Every write attributed to operator
  ✓ No service account used
  ✓ Audit log shows: "Operator X created marker Y at time Z"

Result: Full accountability
        Operator retains identity
        Audit trail complete
```

---

## Integration with OneMind OS / Tactical Operations

### Scenario 1: Rapid COP Updates Under Time Pressure

```
Tactical situation: New enemy contact, grid 12S DG 34567 89012

Traditional workflow:
  1. Click "Add marker" menu
  2. Enter grid coordinates
  3. Select icon (enemy, unknown)
  4. Enter label
  5. Set color (red = threat)
  6. Confirm
  Time: ~30 seconds, 6 steps

TAKpilot workflow:
  Operator: "Mark enemy contact at 12S DG 34567 89012"
  TAKpilot: Places marker (1 sentence)
  Time: ~3 seconds, 1 step

Result: 10x faster, under pressure
```

### Scenario 2: Unstructured Intel → COP Objects

```
Reconnaissance team returns with sketch:
  - Hand-drawn map showing:
    * Enemy positions (3 marked)
    * Supply line (hand-annotated)
    * Possible ambush points (question marks)

Traditional workflow:
  1. Redraw sketch manually
  2. Estimate grid coordinates
  3. Create markers one-by-one
  4. Time: ~20 minutes

TAKpilot workflow:
  1. Operator uploads sketch image
  2. TAKpilot reads image (vision)
  3. Extracts enemy positions, supply lines, suspects
  4. Converts to COP markers (georeferenced)
  5. Creates layer group "Recon sketch"
  
  Result: Full COP layer in ~30 seconds
```

### Scenario 3: Integration with DDF Intelligence

```
Intelligence analyst queries DDF (Distributed Data Framework):
  "Show me all known enemy camps in sector A"

DDF returns:
  - Historical positions
  - Timing patterns
  - Confidence levels
  - Related intelligence

TAKpilot processes:
  1. Receives DDF query results
  2. Converts to COP markers (color-coded by confidence)
  3. Creates layer "DDF Known Threats"
  4. Displays on COP with time slider (replay over time)

Result: Historical intel integrated into live COP
```

### Scenario 4: Real-Time NATS Position Stream

```
NATS stream: agents.tactical.position (friendly forces)

TAKpilot integration:
  1. Subscribe to NATS position stream
  2. Parse position updates
  3. Update COP markers in real-time
  4. Smooth animation between updates
  5. Keep operator's identity (position updates show as system broadcast)

Result: Live COP feed without operator intervention
```

### Scenario 5: Operator Commands via Voice/Chat

```
Scenario: Operator using voice comms (radio)

Workflow:
  1. Operator speaks to NCO: "Mark grid 12S DG..., set sniper overwatch"
  2. NCO types into TAKpilot chat
  3. TAKpilot places markers on COP
  4. Shows to operator for approval
  5. Operator says "confirm"
  6. NCO clicks "Approve"
  7. Markers go live on COP

Result: Voice commands translate to COP actions
        Low cognitive load (NCO types, TAKpilot acts)
```

---

## Key Features Deep Dive

### Full Transparency & Auditability
```
Every action is:
  ✓ Visible (as tool card)
  ✓ Auditable (timestamp, operator ID, action type)
  ✓ Logged (CloudTAK audit trail)
  ✓ Identifiable (operator's token, not service account)

Result: Impossible for AI action to be invisible
        Commanders can inspect what TAKpilot did
        Compliance with audit requirements
```

### Operator Control Over Destructive Actions
```
Destructive operations:
  - Delete markers/missions
  - Overwrite COP objects
  - Broadcast to team
  - Change layer visibility
  - Modify access controls

All require:
  1. TAKpilot plans action (card shows full details)
  2. Operator reviews (sees exactly what will happen)
  3. Operator decides (approve/deny before execution)
  4. No surprise changes to COP

Result: Operator retains final authority
        No "AI did this without asking"
        Appropriate for combat/tactical use
```

### Model Flexibility (No Lock-in)
```
Deployment options:

Option A: Cloud-connected
  - Use GPT-4, Claude, Llama-API
  - No local hardware needed
  - Requires internet connection
  - Good for rear-area operations

Option B: Edge-local
  - Run local Llama/Mistral on laptop/server
  - No internet required
  - Works in GPS-denied/comms-denied areas
  - Lower latency

Option C: Hybrid
  - Cloud for complex reasoning
  - Edge for decisions/approvals
  - Best of both worlds

Result: Operator chooses based on situation
        Not locked into one vendor
```

### MIL-STD-2525 Compliance
```
All markers follow NATO standards:

Unit symbols:
  ✓ Correct icon shape (infantry, armor, air, etc.)
  ✓ Correct color (friendly=blue, enemy=red, neutral=green)
  ✓ Correct modifiers (size, condition, type)
  ✓ Correct orientation (direction of movement/facing)

Result: Compatible with:
  - ATAK (TAK Android)
  - WinTAK (TAK Windows)
  - Other NATO C2 systems
  - Printed maps (symbols match standards)
```

---

## Target Audience

- **Tactical operators on CloudTAK** — Primary users, direct COP control
- **Mission commanders & S3 staff** — Manage operations, coordinate fires
- **ISR & reconnaissance teams** — Convert field observations to COP
- **Intelligence analysts (S2)** — Integrate intel layers, analyze trends
- **Ukrainian Defense Forces** — Primary deployment (available via Brave1)
- **NATO partner units on TAK** — Interoperable deployment
- **Public safety & emergency management** — Disaster response, incident command
- **Defense integrators & CloudTAK deployers** — Integration partners

---

## Deployment & Availability

### Distribution Channel
- **Brave1 / Delta Marketplace** (Ukrainian defense innovation platform)
- Available to Ukrainian armed forces + NATO partners
- Commercial licensing model

### Deployment Modes
```
Option 1: CloudTAK SaaS
  - TAKpilot runs as chat panel in CloudTAK web app
  - No installation needed
  - Works on any device with browser + CloudTAK access

Option 2: Self-hosted
  - Deploy TAKpilot agent server on-premises
  - Connect to local CloudTAK instance
  - Full control, no cloud dependency
```

---

## Problems It Solves

### Challenge 1: Menu Navigation Under Time Pressure
```
Problem:
  Operators lose time navigating menus under fire
  One marker takes 30 seconds (6 clicks)
  When every second matters, this is unacceptable

Solution:
  TAKpilot: One sentence = one marker placed instantly
  Result: 10x faster updates
```

### Challenge 2: Unstructured Intel Has No Fast Path
```
Problem:
  Reconnaissance sketches, handwritten notes, PDF reports
  Cannot be quickly added to COP
  Must be manually redrawn by S2 analyst
  Can take hours

Solution:
  TAKpilot reads images/PDFs with vision
  Extracts georeferenced data automatically
  Converts to COP objects in seconds
  Result: Intel on COP in < 1 minute
```

### Challenge 3: AI Actions Are Invisible
```
Problem:
  Traditional AI systems execute actions without showing work
  Operators cannot audit what happened
  Impossible to debug failures
  Creates distrust

Solution:
  TAKpilot streams every tool call as visible card
  Operator sees exactly what will happen before approval
  Full audit trail in CloudTAK logs
  Result: Transparent, auditable AI
```

### Challenge 4: Destructive AI Actions Execute Without Review
```
Problem:
  AI system deletes/overwrites COP without asking
  No chance to review before damage done
  Operator loses control

Solution:
  TAKpilot requires explicit approval for all destructive actions
  Shows full context (what exactly will be deleted?)
  Operator can deny (cancels before execution)
  Result: Operator retains final authority
```

### Challenge 5: Model Lock-in
```
Problem:
  System requires specific model provider
  If connectivity lost or provider unavailable, system fails
  Network security policies may forbid cloud models

Solution:
  TAKpilot model-agnostic
  Deploy cloud OR edge models based on situation
  Switch between models without reconfiguration
  Result: Flexibility, resilience, operational security
```

---

## Related Services (Corvus Intelligence)

Corvus Intelligence offers complementary services:

| Service | What It Does | Integration with TAKpilot |
|---------|------------|--------------------------|
| **C2 Dashboard Development** | Custom command & control interfaces | Could build custom TAKpilot variants |
| **Data Fusion & Integration** | Combine data from multiple sources | TAKpilot consumes fused data |
| **Secure Cloud & Infrastructure** | Deployment & hosting | Runs TAKpilot infrastructure |

---

## Competitive Advantages

| Aspect | TAKpilot | Traditional COP Update |
|--------|----------|----------------------|
| **Speed** | 1 sentence | 5+ menu clicks |
| **Audit** | Visible cards + logs | Implicit operator action |
| **Control** | Operator approves | No approval gate |
| **Unstructured data** | Auto-convert (vision) | Manual re-entry |
| **Model flexibility** | Cloud or edge | Lock-in |
| **NATO compliance** | MIL-STD-2525 guaranteed | User responsible |
| **Integration** | CloudTAK native | Same as other apps |

---

## Technical Architecture Details

### Security Model
```
Session-scoped sandbox:
  ✓ Each operator gets isolated session
  ✓ Files written to ephemeral storage
  ✓ No cross-session data access
  ✓ Operator's CloudTAK token used for API calls
  ✓ No service account (no credential sharing)
  ✓ Per-session encryption

Result: Operator data confined, no data leakage
```

### Latency & Performance
```
Expected performance:
  - Read-only tool calls: <100ms (auto-execute)
  - Write operations: <500ms (display card, wait for approval)
  - Approval decision: <1s (operator clicks)
  - API execution: <100ms (CloudTAK response)
  
Total for typical command: ~1-2 seconds
```

### Scalability
```
Per-operator:
  - TAKpilot agent sandbox: ~50-100MB RAM
  - Per-session storage: ~10-100MB
  - Network bandwidth: <1Mbps (mostly CloudTAK API)

For 100 concurrent operators:
  - Server RAM: ~5-10GB
  - Network: ~100Mbps uplink
  - Storage: ~1TB (all sessions)
```

---

## Use Cases & Scenarios

### Tactical Operations Center (TOC)
```
S3 operations staff managing live operation:
  - UAV feeds incoming
  - Enemy contacts reported
  - Friendly positions changing
  
Workflow:
  1. Intelligence officer: "Mark those three enemy positions"
  2. TAKpilot: Places markers (3 seconds)
  3. S3 approves (sees markers)
  4. COP updates live
  5. Everyone sees updated picture
```

### Forward Operating Base (FOB)
```
Small team with TAKpilot on laptop:
  - Limited connectivity
  - Deploy edge model (local inference)
  - Generate local COP without internet
  
Workflow:
  1. Operator: "Reconn photos uploaded, extract positions"
  2. TAKpilot (local): Reads photos, extracts grids
  3. Creates local layer "Recon intel"
  4. Later syncs to TOC when comms available
```

### Intelligence Analysis
```
S2 analyst building threat picture:
  - Combine multiple data sources
  - Historical + current intel
  - Pattern recognition
  
Workflow:
  1. Analyst: "Show me all enemy camps from DDF + recent reports"
  2. TAKpilot: Queries DDF, uploads PDFs, combines
  3. Creates layer "Threat picture"
  4. Shows confidence levels (color-coded)
```

---

## Known Limitations & Considerations

### Limitations
- Requires CloudTAK (cannot work standalone)
- Operator must have CloudTAK access + proper permissions
- Vision module works best with clear PDFs/images
- LLM reasoning limited by model capability
- Edge models require adequate local hardware

### Operational Considerations
- Train operators on TAKpilot interface
- Establish approval procedures (who approves what?)
- Define command vocabulary (standardize phrasing)
- Set up audit review process
- Plan model switching strategy (cloud vs. edge)

---

## Integration Roadmap (Potential)

Future TAKpilot enhancements for OneMind OS:

```
Phase 1: Current
  ✓ CloudTAK integration
  ✓ Manual approval gates
  ✓ File upload (vision)
  
Phase 2: Planned
  - [ ] NATS stream subscriber (real-time position updates)
  - [ ] DDF query integration
  - [ ] Automated rules engine (approve certain actions silently)
  - [ ] Voice input (microphone → command)
  
Phase 3: Future
  - [ ] Multi-agent coordination (TAKpilot agents collaborate)
  - [ ] Predictive suggestions (TAKpilot suggests likely next actions)
  - [ ] Swarm drone control (via TAKpilot chat)
```

---

## Contact & Deployment

**Vendor:** Corvus Intelligence  
**Email:** contact@corvusintell.com  
**Website:** https://corvusintell.com/  
**Demo:** https://corvusintell.com/takpilot/ (request demo button)  

**Procurement:** Through Brave1 / Delta Marketplace (Ukrainian Defense Procurement)

---

**Captured:** 2026-07-09  
**Status:** Production-ready, active deployment  
**Relevance:** Critical for rapid COP updates, unstructured intel processing, operator workload reduction  
**Integration Priority:** High (complements TAK Server, NATS position streams, DDF intel)  
**Learning curve:** Low (natural language interface, familiar to operators)  
**Cost model:** Commercial license (pricing via Brave1 marketplace)  
**Competitive advantage:** Model-agnostic, operator-controlled, auditable, NATO-compliant  
**Key innovation:** Agentic AI with transparent tool calls + operator approval gates (not autonomous)  

---

## Related Resources in Convergence

- **TAK Server** — Mission command backend (TAKpilot client)
- **DDF (Distributed Data Framework)** — Intelligence source (TAKpilot queries)
- **NATS** — Event bus (TAKpilot could subscribe to position streams)
- **CesiumJS** — 3D visualization (TAKpilot could drive CesiumJS via CloudTAK)
- **Anduril EagleEye** — Warfighter AR platform (complement to TAKpilot for ground ops)
