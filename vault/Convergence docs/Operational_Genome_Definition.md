# Operational Genome: Definition & Tactical System Architecture

**Term:** Operational Genome  
**Domain:** Military/tactical operations, system design  
**Concept:** The complete set of operational behaviors, procedures, and capabilities that define how a military or tactical force operates  
**Relevance:** Blueprint for tactical operations, reproducible procedures, organizational knowledge capture  

---

## What It Means

An **operational genome** is the complete, encoded set of operational behaviors, procedures, tactics, and capabilities that define:
- HOW the force operates
- WHAT it can do
- WHEN to execute procedures
- WHY each procedure matters
- WHO does what

Like a biological genome encodes DNA (blueprint for organisms), an operational genome encodes the tactical blueprint (what the force does and how).

### Core Concept

```
Biological genome:
  DNA → proteins → organism behavior → phenotype

Operational genome:
  Procedures → training → unit behavior → tactical capability
  
Both are:
  - Heritable (pass from generation to generation)
  - Reproducible (express same way each time)
  - Specific (encode particular traits)
  - Evolving (mutate and adapt)
  - Decodable (can be read and understood)
```

---

## Components of an Operational Genome

### 1. **Tactical Procedures (Standard Operating Procedures)**

**What:** Repeatable sequences for common operations

**Examples:**
```
SOP: "Enemy contact drill"
  1. Take cover
  2. Return fire
  3. Report contact to command
  4. Move to next position
  5. Regroup

SOP: "Route clearance"
  1. Scout ahead (visual)
  2. Check for IEDs (detector)
  3. Mark safe route
  4. Report all-clear
  5. Proceed

SOP: "Casualty evacuation"
  1. Treat wound (field medicine)
  2. Call for evacuation
  3. Prepare landing zone
  4. Provide security
  5. Load casualty
```

### 2. **Decision-Making Rules (Rules of Engagement)**

**What:** When to act, what actions are authorized, what's prohibited

**Examples:**
```
ROE: "Engagement authority"
  - Officer can authorize strikes
  - NCO can authorize defensive fire
  - Enlisted can return fire if attacked
  - All can call for support

ROE: "Force protection"
  - Shoot to stop threat (not to kill)
  - Minimum force principle
  - Positive identification required
  - Civilian casualties avoided
```

### 3. **Communication Protocols (How Units Talk)**

**What:** Radio procedures, message formats, reporting requirements

**Examples:**
```
9-line CAS format (documented in manual)
SALUTE spot report format (defined procedure)
SITREP template (standardized fields)
Grid coordinate format (MGRS, UTM, lat/lon)
Status codes ("Contact", "Engage", "Hold")
```

### 4. **Capability Matrix (What We Can Do)**

**What:** Inventory of capabilities and their limitations

**Examples:**
```
Infantry squad:
  - Suppress targets to 300m
  - Move 5km in 1 hour (terrain dependent)
  - Sustain 4 hours continuous ops
  - Carry 40kg per person

Drone (UAV):
  - Loiter 8 hours
  - See 2km away
  - Carry 2kg payload
  - Operate in daylight only

Mortar section:
  - Range: 1-8km
  - Rate: 20 rounds/minute
  - Ammunition: HE, smoke, illumination
```

### 5. **Knowledge Base (Training Records, Lessons Learned)**

**What:** Documented experience and expertise

**Examples:**
```
Lessons learned:
  - Ambush in open terrain requires 360° security
  - Urban operations need building-by-building clearing
  - Night operations improve without moonlight
  - Radio discipline prevents enemy DF (direction finding)

Expertise:
  - Sergeant knows terrain in 12S DG sector
  - Corporal trained in explosives
  - Private speaks local language
  - Lieutenant has CAS coordination experience
```

### 6. **Organizational Structure (Chain of Command)**

**What:** Who reports to whom, authority flow

**Examples:**
```
Company structure:
  Captain (commander)
    ↓
  Platoon leader (lieutenant)
    ↓
  Squad leader (sergeant)
    ↓
  Team leaders (corporal)
    ↓
  Soldiers

Authority:
  - Captain: Operation planning
  - Lieutenant: Unit movement
  - Sergeant: Squad tactics
  - Corporal: Team actions
  - Soldier: Execute orders
```

---

## Operational Genome in Practice

### Scenario: Inherited Tactical Knowledge

**Situation:** New commander takes over unit

**Without operational genome:**
```
New commander learns by:
  - Trial and error (dangerous)
  - Asking people (inconsistent)
  - Reinventing procedures (slow)
  - Making mistakes (costly)

Result: Weeks/months to be effective
        Inconsistent procedures
        Risk of failure
```

**With operational genome:**
```
New commander:
  1. Reads documented procedures (1 day)
  2. Reviews lessons learned (1 day)
  3. Understands capability matrix (0.5 day)
  4. Knows chain of command (0.5 day)
  5. Inherits unit's knowledge (immediately)

Result: Effective in 2-3 days
        Consistent procedures
        Known capabilities
        Reduced risk
```

### Example: Encoded Procedure

**TAK Server deployment as operational genome:**

```
Operational genome encoded in:
  - Docker containers (environment captured)
  - Configuration files (settings encoded)
  - Database schemas (data structure encoded)
  - Documentation (procedures encoded)
  - Training materials (knowledge encoded)

Result:
  - Deploy on new machine: Same TAK Server
  - Run same query: Same results
  - Execute same procedure: Same outcome
  - Train new user: Same training, same capability

Operational genome = "here's how TAK works"
```

---

## Operational Genome in OneMind OS

### OneMind as Operational Genome Repository

**OneMind components encode operational knowledge:**

| Component | Encodes | Value |
|-----------|---------|-------|
| **TAK Server** | Position procedures, CoT format, COP display | Units understand procedures |
| **NATS** | Message routing, pub-sub patterns | Teams communicate consistently |
| **PostGIS** | Spatial procedures, geographic queries | Consistent geospatial analysis |
| **DDF** | Intelligence catalog structure | Standardized intel access |
| **Mumble** | Voice comms procedures, encryption | Secure communication |
| **CesiumJS** | 3D visualization, terrain rendering | Common operational view |
| **ATAK** | Tactical procedures, CoT interpretation | Standardized tactical display |

### Capturing Operational Genome in OneMind

```
OneMind operational genome includes:

1. Procedures (documented in TAK)
   - How to mark positions
   - How to send messages
   - How to coordinate fires

2. Capability matrix (coded in components)
   - TAK can handle 10,000 CoT events
   - NATS can route 100,000 messages/sec
   - PostGIS can query 100M geometries

3. Knowledge base (documented in skills, manuals)
   - Training materials
   - Best practices
   - Lessons learned

4. Communication protocols (TAK CoT standard)
   - Format: CoT XML
   - Encoding: deterministic
   - Interpretation: standardized

Result: Unit deploying OneMind inherits
        - Proven procedures
        - Known capabilities
        - Tested knowledge
        - Standardized protocols
```

---

## Operational Genome vs. Doctrine

### Distinction

**Doctrine:**
- High-level strategic principles
- "We believe in rapid maneuver"
- "Minimize civilian casualties"
- General guidance

**Operational Genome:**
- Specific tactical procedures
- "Execute enemy contact drill: take cover → return fire → report → move"
- "9-line CAS format with grid + friendlies + special requests"
- Detailed, repeatable, encoded

### Relationship

```
Doctrine (strategic)
    ↓ (informs)
Operational Genome (tactical)
    ↓ (implements)
Unit Procedures (execution)
    ↓ (results in)
Tactical Effectiveness
```

---

## Building an Operational Genome

### Process

**1. Identify successful procedures**
```
What works?
  - Enemy contact drill (proven effective)
  - Route clearance (low casualties)
  - Casualty evacuation (fast response)
```

**2. Document procedures**
```
Write down:
  - Step-by-step execution
  - Roles and responsibilities
  - Decision points
  - Success criteria
```

**3. Test and validate**
```
Do procedures work?
  - Run drills
  - Measure effectiveness
  - Identify gaps
  - Refine procedures
```

**4. Encode in systems**
```
Make procedures usable:
  - TAK templates
  - NATS message formats
  - Database schemas
  - Software workflows
```

**5. Train personnel**
```
Teach procedures:
  - Classroom training
  - Field exercises
  - Simulations
  - On-the-job training
```

**6. Iterate and improve**
```
Update based on:
  - Lessons learned
  - New technologies
  - Changing threats
  - Feedback from operators
```

---

## Operational Genome in TAK/ATAK

### Example: COP as Operational Genome

**TAK Server encodes:**

1. **Position procedure** — How to report position
   ```
   CoT message = "Here I am, at this grid, with this callsign"
   Format is deterministic, understood by all
   ```

2. **Update procedure** — How often to send updates
   ```
   Every 10 seconds (interval encoded)
   Every 1 kilometer movement (distance encoded)
   Every bearing change > 5 degrees
   ```

3. **Display procedure** — How to show on map
   ```
   CoT type → NATO icon
   Unit type → color (blue=friendly, red=enemy)
   Status → animating (moving) or static
   ```

4. **Query procedure** — How to ask "where are you?"
   ```
   SQL query on CoT database
   Returns position + time + status
   Deterministic (same query = same result)
   ```

**Result:** All units understand procedures consistently

---

## Operational Genome for Resilience

### Why Operational Genome Matters for OneMind

**Problem:** If ORP (central processor) fails, units lose operational knowledge

**Solution:** Encode operational genome in distributed systems

```
Normal operation:
  Units report → TAK Server (ORP knows procedures)

ORP fails:
  Units still have:
    - TAK (knows CoT procedures)
    - NATS (knows message procedures)
    - Mumble (knows voice procedures)
    - Local knowledge

Result: Operational genome survives ORP failure
        Units can continue with reduced capabilities
        Graceful degradation, not collapse
```

---

## Capturing OneMind's Operational Genome

### What OneMind Encodes

**Procedural knowledge:**
- "Here's how to deploy TAK Server"
- "Here's how to write a NATS converter"
- "Here's how to query PostGIS for threats"

**Tactical knowledge:**
- "Here's how to mark positions"
- "Here's how to coordinate with allies"
- "Here's how to request air support (9-line)"

**Capability knowledge:**
- "TAK handles 10,000 CoT events"
- "NATS routes 100k messages/sec"
- "PostGIS queries million positions instantly"

**Communication knowledge:**
- "CoT XML format for positions"
- "9-line format for CAS"
- "SALUTE format for spot reports"

### Encoding the Genome

```
OneMind operational genome is documented in:

1. Code repositories
   - TAK plugins
   - NATS subjects
   - PostGIS queries

2. Configuration files
   - TAK Server settings
   - NATS JetStream
   - Mumble channels

3. Documentation
   - Deployment guides
   - Procedure manuals
   - Capability reference

4. Skills (Hermes)
   - Training materials
   - Best practices
   - Troubleshooting

5. Tests
   - Validation procedures
   - Performance benchmarks
   - Failure scenarios
```

---

## Operational Genome Evolution

### How Genomes Change

**Natural evolution:**
```
DNA mutates → New traits → Selection → Population changes
```

**Operational evolution:**
```
New lessons learned → Procedures update → Training evolves → Capability changes
```

### Versioning Operational Genome

```
Version 1.0: Initial deployment
  - Basic procedures
  - Known limitations

Version 1.1: After-action review
  - Refined procedures
  - Added lessons learned

Version 2.0: Major upgrade
  - New capabilities added
  - Obsolete procedures removed
  - Efficiency improvements

Result: Operational genome evolves
        Captured in version control
        Can revert if needed
        Improvements documented
```

---

## Summary: Operational Genome

**Operational genome = complete blueprint of how force operates**

- **Encodes:** Procedures, capabilities, knowledge, communication
- **Enables:** Consistency, training, resilience, inheritance
- **Captured in:** TAK, NATS, PostGIS, documentation, code
- **Evolves:** Through lessons learned, updates, improvements
- **Ensures:** New commanders/units inherit proven methods
- **Survives:** Equipment failures, personnel turnover, technology changes

**For OneMind:**
```
OneMind's operational genome includes:
  - TAK procedures (how to use COP)
  - NATS protocols (how to route messages)
  - CoT standards (how to report positions)
  - Spatial procedures (how to analyze geography)
  - Voice procedures (how to communicate securely)

Deploy OneMind = Deploy operational genome
  Units inherit:
    - Proven procedures
    - Known capabilities
    - Tested knowledge
    - Standardized protocols

Result: Tactically effective unit from day 1
```

---

**Captured:** 2026-07-09  
**Domain:** Military/tactical operations, system design  
**Relevance:** Encodes how tactical forces operate, ensures reproducibility and inheritance  
**Concept:** DNA of operations (procedures, capabilities, knowledge)  
**Implementation:** TAK, NATS, PostGIS, documentation, training  
**Value:** Consistency, resilience, training efficiency, knowledge preservation  
**Evolution:** Versioned, updated with lessons learned, continuously improving  

---

## Related Resources in Convergence

- **Evolution_to_C4ISR.md** — Operational procedures for command and control
- **Protocol_to_CoT_Converters.md** — Standardized data conversion procedures
- **TAKpilot.md** — AI chat interface for operational procedures
- **NSTAC_Family.md** — Standardized protocol procedures
- **ATAK_GIS_Data_Display.md** — Spatial analysis procedures
- **Deterministic_Definition.md** — Reproducible procedures guarantee
