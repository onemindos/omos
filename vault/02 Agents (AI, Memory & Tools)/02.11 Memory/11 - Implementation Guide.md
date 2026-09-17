---
id: ui-02-20260111053925
title: "Memory Block Implementation Guide"
type: reference
pillar: UI
domain: 00-24
folder: "00-24 UI (Unified Intelligence)"
status: active
created: 2026-01-11
modified: 2026-01-11
owner: zeus
tags: [domain/ui, folder/02]
todoist_project: "00-24 UI (Unified Intelligence)"
ai_summary: "Memory Block Implementation Guide"
---

# Memory Block Implementation Guide

**Purpose:** Practical guide for implementing these memory blocks in Letta  
**Audience:** Technical implementation team  
**Status:** Ready for implementation

---

## Implementation Sequence

### Phase 1: Core Shared Memory (Week 1)
1. Create the 4 core shared memory blocks in Letta
2. Verify all blocks are readable by test agents
3. Test update permissions (Legacy can write all, specific agents write User Profile)

**Priority Order:**
1. [[00-09 Command (Meta)/01 Legacy-AI (Co-Pilot)/01.02 Letta Agents/Memory Block Architecture/01 - Legacy Codex]] - Foundation for all agent behavior
2. [[00-09 Command (Meta)/01 Legacy-AI (Co-Pilot)/01.02 Letta Agents/Memory Block Architecture/02 - User Profile]] - Critical for personalization
3. [[00-09 Command (Meta)/01 Legacy-AI (Co-Pilot)/01.02 Letta Agents/Memory Block Architecture/03 - System State]] - Required for monitoring
4. [[00-09 Command (Meta)/01 Legacy-AI (Co-Pilot)/01.02 Letta Agents/Memory Block Architecture/04 - Cross-Domain Context]] - Needed for coordination

### Phase 2: Domain Memory (Week 2)
1. Create domain memory blocks for each pillar
2. Assign read access to relevant agents
3. Test cross-domain references

**Implementation Order:**
1. [[00-09 Command (Meta)/01 Legacy-AI (Co-Pilot)/01.02 Letta Agents/Memory Block Architecture/08 - Innovative Technology Domain]] - Infrastructure must be first
2. [[00-09 Command (Meta)/01 Legacy-AI (Co-Pilot)/01.02 Letta Agents/Memory Block Architecture/05 - Holistic Performance Domain]] - Affects all other domains
3. [[00-09 Command (Meta)/01 Legacy-AI (Co-Pilot)/01.02 Letta Agents/Memory Block Architecture/06 - Generational Entrepreneurship Domain]] - Business operations
4. [[00-09 Command (Meta)/01 Legacy-AI (Co-Pilot)/01.02 Letta Agents/Memory Block Architecture/07 - Legacy Evolution Domain]] - Family and homestead

### Phase 3: Agent-Specific Memory (Week 3-4)
1. Create agent-specific blocks one agent at a time
2. Test each agent's ability to read shared memory and write own memory
3. Verify memory isolation (agents can't write to other agent's blocks)

**Agent Rollout Priority:**
1. Infrastructure Agent (must monitor system health first)
2. Intelligence Agent (data pipeline foundation)
3. Communications Agent (external interface)
4. Performance Agent (affects daily operations)
5. Remaining agents as needed

---

## Letta Implementation Code Examples

### Creating a Core Memory Block

```python
from letta import create_block

# Example: Legacy Codex
legacy_codex = create_block(
    label="legacy_codex",
    value="""
# Legacy Codex - Foundational Principles

[Full content from the markdown file]
    """,
    metadata={
        "type": "core_shared",
        "access": "read_all",
        "update_frequency": "quarterly",
        "character_limit": 8000
    }
)
```

### Creating an Agent with Memory Blocks

```python
from letta import create_agent

# Example: Intelligence Agent
intelligence_agent = create_agent(
    name="Intelligence Agent",
    
    # Memory blocks this agent can access
    memory_blocks=[
        legacy_codex,              # Shared core (read-only)
        user_profile,              # Shared core (read-only)
        system_state,              # Shared core (read-only)
        holistic_perf_context,     # Domain (read-only)
        intelligence_agent_memory  # Agent-specific (read/write)
    ],
    
    # Tools this agent uses
    tools=[
        "feedparser",
        "beautifulsoup4",
        "obsidian_api",
        "rag_system_api"
    ],
    
    # System prompt
    system_prompt="""
    You are the Intelligence Agent for Legacy AI.
    
    Your role: Process 100+ information sources daily, distill key insights,
    and deliver actionable intelligence to Zeus and other agents.
    
    Your memory blocks contain:
    - Legacy Codex: Foundational principles (read-only)
    - User Profile: Zeus's patterns and preferences (read-only)
    - Your operational memory: Your own state and history (read/write)
    
    Key behaviors:
    - Deliver morning brief at 6:00 AM daily
    - Deliver evening wrap at 6:00 PM daily
    - Prepare all content for RAG system
    - Track engagement patterns and optimize
    """
)
```

### Agent Reading from Memory

```python
# Intelligence Agent checking user's information preferences
def get_user_info_preferences(agent):
    # Read from User Profile block
    profile = agent.read_memory_block("user_profile")
    
    # Parse for relevant section
    topics = profile.get_section("Information Diet > High-Interest Topics")
    
    return topics

# Returns: ["AI/ML advancements", "Robotics", "Agriculture tech", ...]
```

### Agent Updating Own Memory

```python
# Intelligence Agent logging successful processing run
def log_processing_run(agent, items_processed):
    # Read current state
    current_stats = agent.read_memory_block(
        "intelligence_agent_memory",
        section="Operational Statistics"
    )
    
    # Update statistics
    agent.update_memory_block(
        "intelligence_agent_memory",
        section="Operational Statistics",
        updates={
            "last_successful_run": datetime.now().isoformat(),
            "items_processed": items_processed,
            "successful_runs": current_stats["successful_runs"] + 1
        }
    )
    
    # Also update System State (shared)
    agent.update_memory_block(
        "system_state",
        section="Agent Status > Intelligence Agent",
        updates={
            "status": "active",
            "last_successful_run": datetime.now().isoformat(),
            "pending_items": 0,
            "errors": []
        }
    )
```

### Cross-Agent Memory Sharing

```python
# Finance Agent checking Performance Agent's biometric data
def get_user_energy_state(finance_agent, legacy_orchestrator):
    # Finance Agent requests through Legacy
    request = {
        "to_agent": "performance",
        "query": "What's Zeus's current energy state?",
        "context": "Need to know if it's a good time for complex financial decisions"
    }
    
    # Legacy routes to Performance Agent
    response = legacy_orchestrator.route_request(request)
    
    # Performance Agent reads its memory and responds
    energy_state = response["energy_state"]  # e.g., "HIGH - optimal for decisions"
    
    return energy_state
```

---

## Memory Access Control Matrix

| Memory Block | Legacy | Intelligence | Performance | Finance | Business Ops | ... |
|--------------|--------|--------------|-------------|---------|--------------|-----|
| Legacy Codex | R/W | R | R | R | R | R |
| User Profile | R/W | R | R/W | R | R | R |
| System State | R/W | R | R | R | R | R |
| Cross-Domain | R/W | R | R | R | R | R |
| Holistic Perf | R/W | R | R/W | R | - | - |
| Gen Ent | R/W | R | - | R | R/W | R |
| Legacy Evo | R/W | - | - | - | - | - |
| Innov Tech | R/W | - | - | - | - | - |
| Agent-Specific | R | Owner: R/W | Owner: R/W | Owner: R/W | Owner: R/W | ... |

**Legend:** R = Read, W = Write, R/W = Read and Write, - = No Access

---

## Memory Update Protocols

### Core Shared Memory
- **Legacy Codex:** Quarterly updates or major strategic shifts (Legacy only)
- **User Profile:** 
  - Patterns: Weekly updates (Legacy + select agents)
  - Current context: Real-time updates (Legacy + relevant agents)
- **System State:** Real-time updates (all agents can update their own status)
- **Cross-Domain Context:** As needed (Legacy primary, agents suggest updates)

### Domain Memory
- Updated by agents within that domain
- Legacy reviews and approves major structural changes
- All updates logged in System State for auditability

### Agent-Specific Memory
- Full autonomy for each agent to update their own memory
- Operational history auto-updates based on activity
- Skills/tools updated when capabilities change

---

## Memory Size Management

**Token/Character Budgets:**
- Core Shared: ~12,000 chars total (always loaded)
- Domain Memory: ~10,000 chars total (loaded on-demand by pillar)
- Agent-Specific: ~5,000 chars per agent (loaded when agent active)

**Letta Character Limit:** 20,000 per block (we're well within limits)

**Optimization Strategies:**
1. **Compression:** Older operational history summarized, details archived to Obsidian
2. **Relevance Pruning:** Low-value data periodically removed
3. **Tiered Loading:** Load only relevant memory blocks per agent interaction
4. **Archive System:** Move historical data to Obsidian vault (remains RAG-accessible)

---

## Versioning & Backup

**Git-Style Version Control:**
```python
# Every memory block update creates a version
def update_memory_with_versioning(block_name, updates):
    # Get current version
    current = get_memory_block(block_name)
    
    # Create version snapshot
    version = {
        "timestamp": datetime.now().isoformat(),
        "previous_content": current,
        "changes": updates,
        "changed_by": agent_name
    }
    
    # Store version in PostgreSQL
    store_version(block_name, version)
    
    # Apply updates
    apply_updates(block_name, updates)
```

**Backup Schedule:**
- Critical blocks (Codex, Profile): Hourly snapshots
- Domain blocks: Daily snapshots
- Agent blocks: Daily snapshots
- All backups: 30-day retention in PostgreSQL
- Weekly full system backup to S3

**Rollback Capability:**
```python
# Rollback to previous version if needed
def rollback_memory_block(block_name, version_id):
    version = get_version(block_name, version_id)
    restore_memory_block(block_name, version["previous_content"])
    log_rollback(block_name, version_id, reason="Manual rollback")
```

---

## Testing & Validation

### Phase 1 Tests: Core Memory
```python
# Test 1: All agents can read Legacy Codex
for agent in all_agents:
    codex = agent.read_memory_block("legacy_codex")
    assert codex is not None
    assert "Four Pillars" in codex

# Test 2: Only authorized agents can write User Profile
try:
    intelligence_agent.write_memory_block("user_profile", {...})
    assert False, "Intelligence Agent should not write to User Profile"
except PermissionError:
    pass  # Expected

# Test 3: Legacy can write to all core blocks
for block in core_blocks:
    legacy.write_memory_block(block, test_data)
    assert legacy.read_memory_block(block) == test_data
```

### Phase 2 Tests: Domain Memory
```python
# Test 1: Domain isolation (agents can't access irrelevant domains)
try:
    intelligence_agent.read_memory_block("legacy_evolution_domain")
    assert False, "Intelligence Agent should not access Legacy Evolution"
except PermissionError:
    pass  # Expected

# Test 2: Domain agents can read their domain
perf_context = performance_agent.read_memory_block("holistic_performance_domain")
assert "Performance Baselines" in perf_context
```

### Phase 3 Tests: Agent Memory
```python
# Test 1: Agents can write to own memory only
intelligence_agent.write_memory_block("intelligence_agent_memory", updates)
assert intelligence_agent.read_memory_block("intelligence_agent_memory") == updates

try:
    intelligence_agent.write_memory_block("infrastructure_agent_memory", updates)
    assert False, "Should not write to other agent's memory"
except PermissionError:
    pass  # Expected

# Test 2: Memory persists across agent restarts
agent_id = intelligence_agent.id
intelligence_agent.write_memory_block("intelligence_agent_memory", test_data)

# Restart agent
new_agent = load_agent(agent_id)
assert new_agent.read_memory_block("intelligence_agent_memory") == test_data
```

---

## Monitoring & Observability

**Memory Health Dashboard:**
- Block size tracking (ensure under character limits)
- Update frequency per block
- Version history and rollback events
- Access patterns (which agents read which blocks)
- Performance metrics (read/write latency)

**Alerts:**
- Block approaching character limit (>80% of 20,000)
- Excessive update frequency (potential bug)
- Permission violations (unauthorized access attempts)
- Version conflicts (concurrent update attempts)

---

## Troubleshooting Guide

**Problem:** Agent can't read a memory block
- Check: Access control matrix - does agent have read permission?
- Check: Block exists in Letta (may not be created yet)
- Check: Block name matches exactly (case-sensitive)

**Problem:** Agent's updates not persisting
- Check: Agent has write permission to that block
- Check: Update method being called correctly
- Check: No concurrent updates causing conflicts
- Check: Block hasn't reached character limit

**Problem:** Memory blocks too large
- Solution: Move historical data to Obsidian archives
- Solution: Compress verbose sections
- Solution: Split into multiple smaller blocks if necessary

**Problem:** Slow memory access
- Check: Are we loading unnecessary blocks?
- Solution: Implement lazy loading (load on-demand)
- Solution: Cache frequently accessed blocks in Redis

---

## Next Steps After Implementation

1. **Monitor Performance:** Track memory access patterns for first 2 weeks
2. **Optimize Patterns:** Adjust block structure based on usage
3. **Agent Feedback:** Gather feedback from agents about memory usefulness
4. **Iterate:** Refine memory content based on actual operations
5. **Scale:** Add remaining agent-specific blocks as agents are deployed

---

## Related Documents
- [[00-09 Command (Meta)/01 Legacy-AI (Co-Pilot)/01.02 Letta Agents/Memory Block Architecture/00 - Memory Architecture Overview]] - High-level architecture
- [[Agent Roster]] - Complete list of agents needing memory blocks
- [[Letta Documentation]] - Official Letta memory system docs

*Created: 2025-11-08*
*Status: Ready for implementation*
