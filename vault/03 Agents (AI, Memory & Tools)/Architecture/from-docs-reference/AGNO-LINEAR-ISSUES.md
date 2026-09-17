# Agno Gap - Linear Issues

> **Purpose:** Complete list of Linear issues to create for Agno feature gaps
> **Total Issues:** 68 (8 P0 + 22 P1 + 25 P2 + 13 P3)
> **Date:** 2026-01-26

---

## Issue Numbering Convention

Issues will be created as OMOS-203 onwards (continuing from existing OMOS-202).

---

## P0 - Critical Issues (8)

### OMOS-203: Wire Knowledge Auto-Search to Agent Runs
```
Title: Wire knowledge auto-search to agent runs
Priority: Urgent
Labels: backend, agno, knowledge, p0-critical
Estimate: 4 points

Description:
The knowledge search endpoint exists but is not automatically used during agent runs.

**Current State:**
- `/api/knowledge/search` endpoint works
- `search_knowledge` parameter exists but not wired
- Agent doesn't automatically search knowledge base

**Required:**
- When agent has `knowledge` attached and `search_knowledge=True`:
  - Automatically embed user query
  - Search knowledge base for relevant documents
  - Include retrieved documents in agent context

**Files:**
- `core/agent/legacy.py`
- `core/agent/factory.py`
- `core/knowledge/base.py`

**Acceptance Criteria:**
- [ ] Agent with knowledge base automatically searches on each run
- [ ] Retrieved documents appear in agent context
- [ ] `num_documents` parameter respected
- [ ] Test with PDF upload + query
```

### OMOS-204: Enable search_knowledge Parameter in Agent Factory
```
Title: Enable search_knowledge parameter in agent factory
Priority: Urgent
Labels: backend, agno, knowledge, p0-critical
Estimate: 2 points

Description:
Add `search_knowledge` and `num_documents` parameters to agent factory.

**Current State:**
- Parameters not passed to Agno Agent
- Knowledge base attached but not searched

**Required:**
- Pass `search_knowledge=True` by default when knowledge attached
- Pass `num_documents` parameter (default 5)

**Files:**
- `core/agent/factory.py`

**Acceptance Criteria:**
- [ ] search_knowledge parameter works
- [ ] num_documents configurable
```

### OMOS-205: Implement Team Route Mode Properly
```
Title: Implement Team Route coordination mode properly
Priority: Urgent
Labels: backend, agno, teams, p0-critical
Estimate: 4 points

Description:
Team Route mode exists but only basic delegation works. Need full Agno Team.run() behavior.

**Current State:**
- Basic team creation works
- Route mode doesn't properly delegate based on member roles
- Team leader doesn't intelligently choose members

**Required:**
- Team leader analyzes request
- Matches request to appropriate member(s) based on role/expertise
- Returns combined/synthesized response

**Files:**
- `core/teams/coordinator.py`
- `core/orchestrator/main.py`

**Acceptance Criteria:**
- [ ] Team routes requests to appropriate members
- [ ] Member responses combined properly
- [ ] Test with researcher/writer team
```

### OMOS-206: Implement Team Coordinate Mode
```
Title: Implement Team Coordinate coordination mode
Priority: Urgent
Labels: backend, agno, teams, p0-critical
Estimate: 8 points

Description:
Team Coordinate mode not implemented. Leader should explicitly coordinate member contributions.

**Current State:**
- Not implemented
- No way to specify coordinate mode

**Required:**
- Add `mode="coordinate"` parameter to Team
- Leader explicitly assigns tasks to members
- Leader synthesizes all member contributions
- Members can see each other's work via `show_members_responses`

**Files:**
- `core/teams/coordinator.py`
- `core/teams/models.py`

**Acceptance Criteria:**
- [ ] Coordinate mode available
- [ ] Leader explicitly assigns work
- [ ] show_members_responses works
```

### OMOS-207: Implement Team Collaborate Mode
```
Title: Implement Team Collaborate coordination mode
Priority: Urgent
Labels: backend, agno, teams, p0-critical
Estimate: 8 points

Description:
Team Collaborate mode not implemented. Members should discuss and build on each other's work.

**Current State:**
- Not implemented
- No peer discussion capability

**Required:**
- Add `mode="collaborate"` parameter
- Members can see and respond to each other
- Discussion continues until consensus or max rounds
- Leader synthesizes final answer

**Files:**
- `core/teams/coordinator.py`
- `core/teams/models.py`

**Acceptance Criteria:**
- [ ] Collaborate mode available
- [ ] Members discuss with each other
- [ ] Consensus detection works
```

### OMOS-208: Complete Workflow Router Step Execution
```
Title: Complete Workflow Router step execution
Priority: Urgent
Labels: backend, agno, workflows, p0-critical
Estimate: 6 points

Description:
Router step type is scaffolded but doesn't execute properly.

**Current State:**
- Router step defined in step-types
- Selector function not executed
- Dynamic routing not working

**Required:**
- Execute selector function with StepInput
- Return appropriate step(s) based on selector
- Execute returned steps
- Pass output to next step in workflow

**Files:**
- `core/workflows/engine.py`
- `core/workflows/steps.py`

**Acceptance Criteria:**
- [ ] Router selector function executes
- [ ] Correct step chosen based on input
- [ ] Workflow continues after router
```

### OMOS-209: Complete Workflow Parallel Step Execution
```
Title: Complete Workflow Parallel step execution
Priority: Urgent
Labels: backend, agno, workflows, p0-critical
Estimate: 6 points

Description:
Parallel step type is scaffolded but doesn't execute concurrently.

**Current State:**
- Parallel step defined
- Steps execute sequentially
- No concurrent execution

**Required:**
- Execute all parallel steps concurrently (asyncio.gather)
- Collect all outputs
- Combine outputs for next step
- Handle errors in individual steps

**Files:**
- `core/workflows/engine.py`
- `core/workflows/steps.py`

**Acceptance Criteria:**
- [ ] Parallel steps execute concurrently
- [ ] All outputs collected
- [ ] Error handling works
```

### OMOS-210: Enable MCP Auto-Attachment to Agents
```
Title: Enable MCP auto-attachment to agents
Priority: Urgent
Labels: backend, agno, mcp, p0-critical
Estimate: 4 points

Description:
MCP servers are defined but must be manually attached to agents.

**Current State:**
- MCP servers defined in `core/mcp/tools.py`
- Presets available (basic, web, dev, research, full)
- Not auto-attached to agents

**Required:**
- Add `mcp_preset` parameter to agent factory
- Default to "basic" preset (filesystem, memory)
- Auto-connect MCPTools on agent creation
- Proper cleanup on agent destruction

**Files:**
- `core/agent/factory.py`
- `core/mcp/tools.py`

**Acceptance Criteria:**
- [ ] Agents can specify mcp_preset
- [ ] MCPTools auto-attached
- [ ] Tools available in agent.tools
```

---

## P1 - High Priority Issues (22)

### OMOS-211: Add Session State Features to Agent
```
Title: Add session_state features to agent
Priority: High
Labels: backend, agno, agent, p1-high
Estimate: 4 points

Description:
Add session_state, add_session_state_to_context, enable_agentic_state parameters.

**Files:**
- `core/agent/factory.py`
- `core/agent/legacy.py`

**Acceptance Criteria:**
- [ ] session_state persists across runs
- [ ] Agent can update session_state (enable_agentic_state)
- [ ] State appears in context when enabled
```

### OMOS-212: Add Retry/Backoff Configuration to Agent
```
Title: Add retry and backoff configuration to agent
Priority: High
Labels: backend, agno, agent, p1-high
Estimate: 2 points

Description:
Add retries, delay_between_retries, exponential_backoff parameters.

**Files:**
- `core/agent/factory.py`

**Acceptance Criteria:**
- [ ] Agent retries on failure
- [ ] Delay configurable
- [ ] Exponential backoff works
```

### OMOS-213: Add reasoning_model Configuration
```
Title: Add reasoning_model configuration to agent
Priority: High
Labels: backend, agno, agent, p1-high
Estimate: 2 points

Description:
Allow dedicated model for reasoning (e.g., o1-mini for reasoning, gpt-4o for response).

**Files:**
- `core/agent/legacy.py`
- `core/agent/factory.py`

**Acceptance Criteria:**
- [ ] reasoning_model parameter works
- [ ] Reasoning uses dedicated model
- [ ] Response uses main model
```

### OMOS-214: Add Team Model Configuration
```
Title: Add team model configuration
Priority: High
Labels: backend, agno, teams, p1-high
Estimate: 2 points

Description:
Team leader should have configurable model (inherits to members if not set).

**Files:**
- `core/teams/coordinator.py`

**Acceptance Criteria:**
- [ ] Team model configurable
- [ ] Members inherit model if not set
```

### OMOS-215: Add Team-Level Tools
```
Title: Add team-level tools support
Priority: High
Labels: backend, agno, teams, p1-high
Estimate: 4 points

Description:
Teams should support team-level tools available to all members.

**Files:**
- `core/teams/coordinator.py`

**Acceptance Criteria:**
- [ ] Team tools parameter works
- [ ] All members can use team tools
```

### OMOS-216: Add show_members_responses to Teams
```
Title: Add show_members_responses to teams
Priority: High
Labels: backend, agno, teams, p1-high
Estimate: 2 points

Description:
Show individual member responses in team output.

**Files:**
- `core/teams/coordinator.py`

**Acceptance Criteria:**
- [ ] Member responses visible in output
- [ ] Can toggle on/off
```

### OMOS-217: Implement Team Sequential Mode
```
Title: Implement Team Sequential coordination mode
Priority: High
Labels: backend, agno, teams, p1-high
Estimate: 4 points

Description:
Members execute one after another, passing context.

**Files:**
- `core/teams/coordinator.py`

**Acceptance Criteria:**
- [ ] Sequential mode available
- [ ] Context passes between members
```

### OMOS-218: Implement Team Parallel Mode
```
Title: Implement Team Parallel coordination mode (delegate_to_all_members)
Priority: High
Labels: backend, agno, teams, p1-high
Estimate: 4 points

Description:
All members execute simultaneously via delegate_to_all_members=True.

**Files:**
- `core/teams/coordinator.py`

**Acceptance Criteria:**
- [ ] delegate_to_all_members works
- [ ] All members execute concurrently
- [ ] Results combined
```

### OMOS-219: Implement Workflow Loop Step
```
Title: Implement Workflow Loop step type
Priority: High
Labels: backend, agno, workflows, p1-high
Estimate: 6 points

Description:
Add Loop step for iterative execution until condition met.

**Files:**
- `core/workflows/engine.py`
- `core/workflows/steps.py`

**Acceptance Criteria:**
- [ ] Loop step type available
- [ ] Condition function evaluated each iteration
- [ ] max_iterations enforced
```

### OMOS-220: Enhance Workflow Step Input/Output Piping
```
Title: Enhance workflow step input/output piping
Priority: High
Labels: backend, agno, workflows, p1-high
Estimate: 4 points

Description:
Improve how data flows between workflow steps.

**Files:**
- `core/workflows/engine.py`

**Acceptance Criteria:**
- [ ] Previous step output available to next step
- [ ] StepInput includes previous_step_content
- [ ] StepOutput properly formatted
```

### OMOS-221: Add HITL Support at Workflow Steps
```
Title: Add HITL support at workflow steps
Priority: High
Labels: backend, agno, workflows, hitl, p1-high
Estimate: 6 points

Description:
Workflows should pause for HITL confirmation at specific steps.

**Files:**
- `core/workflows/engine.py`
- `core/approvals/handler.py`

**Acceptance Criteria:**
- [ ] Steps can require confirmation
- [ ] Workflow pauses until approved
- [ ] Workflow continues after approval
```

### OMOS-222: Add ShellTools with HITL Confirmation
```
Title: Add ShellTools with HITL confirmation
Priority: High
Labels: backend, agno, tools, hitl, p1-high
Estimate: 4 points

Description:
Enable ShellTools but require HITL confirmation for all commands.

**Files:**
- `core/tools/builtin/shell.py`
- `core/agent/legacy.py`

**Acceptance Criteria:**
- [ ] ShellTools available
- [ ] All commands require confirmation
- [ ] User can approve/reject
```

### OMOS-223: Add PythonTools with HITL Confirmation
```
Title: Add PythonTools with HITL confirmation
Priority: High
Labels: backend, agno, tools, hitl, p1-high
Estimate: 4 points

Description:
Enable PythonTools but require HITL confirmation for code execution.

**Files:**
- `core/tools/builtin/python.py`
- `core/agent/legacy.py`

**Acceptance Criteria:**
- [ ] PythonTools available
- [ ] Code execution requires confirmation
- [ ] Safe sandbox environment
```

### OMOS-224: Add SQLTools with HITL Confirmation
```
Title: Add SQLTools with HITL confirmation
Priority: High
Labels: backend, agno, tools, hitl, p1-high
Estimate: 4 points

Description:
Enable SQLTools but require HITL for write queries.

**Files:**
- `core/tools/builtin/sql.py`
- `core/agent/legacy.py`

**Acceptance Criteria:**
- [ ] SQLTools available
- [ ] SELECT queries work
- [ ] Write queries require confirmation
```

### OMOS-225: Add SlackTools Integration
```
Title: Add SlackTools integration
Priority: High
Labels: backend, agno, tools, integrations, p1-high
Estimate: 4 points

Description:
Add Agno's built-in SlackTools for messaging.

**Files:**
- `core/tools/builtin/slack.py`
- `core/agent/legacy.py`

**Acceptance Criteria:**
- [ ] SlackTools available
- [ ] Can send messages
- [ ] Can list channels
```

### OMOS-226: Add EmailTools Integration
```
Title: Add EmailTools integration
Priority: High
Labels: backend, agno, tools, integrations, p1-high
Estimate: 4 points

Description:
Add Agno's built-in EmailTools (via Resend or SMTP).

**Files:**
- `core/tools/builtin/email.py`
- `core/agent/legacy.py`

**Acceptance Criteria:**
- [ ] EmailTools available
- [ ] Can send emails
- [ ] Template support
```

### OMOS-227: Add FileTools with HITL
```
Title: Add FileTools with HITL confirmation
Priority: High
Labels: backend, agno, tools, hitl, p1-high
Estimate: 4 points

Description:
Enable FileTools but require confirmation for write operations.

**Files:**
- `core/tools/builtin/file.py`

**Acceptance Criteria:**
- [ ] FileTools available
- [ ] Read operations work
- [ ] Write operations require confirmation
```

### OMOS-228: Add FileGenerationTools
```
Title: Add FileGenerationTools
Priority: High
Labels: backend, agno, tools, p1-high
Estimate: 4 points

Description:
Add tools for generating JSON, CSV, PDF, TXT files.

**Files:**
- `core/tools/builtin/file_generation.py`

**Acceptance Criteria:**
- [ ] Can generate JSON files
- [ ] Can generate CSV files
- [ ] Can generate PDF files
- [ ] Can generate TXT files
```

### OMOS-229: Add ReasoningTools
```
Title: Add ReasoningTools
Priority: High
Labels: backend, agno, tools, p1-high
Estimate: 4 points

Description:
Add Agno's ReasoningTools for structured thinking.

**Files:**
- `core/tools/builtin/reasoning.py`
- `core/agent/legacy.py`

**Acceptance Criteria:**
- [ ] ReasoningTools available
- [ ] think() function works
- [ ] analyze() function works
- [ ] plan() function works
```

### OMOS-230: Add Semantic Chunking to Knowledge
```
Title: Add semantic chunking to knowledge base
Priority: High
Labels: backend, agno, knowledge, p1-high
Estimate: 4 points

Description:
Use semantic chunking for better document retrieval.

**Files:**
- `core/knowledge/chunking.py`
- `core/knowledge/base.py`

**Acceptance Criteria:**
- [ ] SemanticChunking strategy available
- [ ] similarity_threshold configurable
- [ ] Better retrieval quality
```

### OMOS-231: Enable MCPTools in Agent by Default
```
Title: Enable MCPTools in agent by default
Priority: High
Labels: backend, agno, mcp, p1-high
Estimate: 2 points

Description:
Make MCPTools available in agent by default (basic preset).

**Files:**
- `core/agent/factory.py`
- `core/mcp/tools.py`

**Acceptance Criteria:**
- [ ] MCPTools auto-available
- [ ] Can override preset
```

### OMOS-232: Implement AccuracyEval
```
Title: Implement AccuracyEval evaluation
Priority: High
Labels: backend, agno, evals, p1-high
Estimate: 8 points

Description:
Complete AccuracyEval implementation for testing agent correctness.

**Files:**
- `core/evals/accuracy.py`
- `core/orchestrator/main.py`

**Acceptance Criteria:**
- [ ] POST /eval-runs with type=accuracy works
- [ ] Results stored in database
- [ ] GET /eval-runs returns results
```

---

## P2 - Medium Priority Issues (25)

### OMOS-233: Add expected_output Parameter
```
Title: Add expected_output parameter to agent
Priority: Medium
Labels: backend, agno, agent, p2-medium
Estimate: 1 point
```

### OMOS-234: Add timezone_identifier Parameter
```
Title: Add timezone_identifier parameter to agent
Priority: Medium
Labels: backend, agno, agent, p2-medium
Estimate: 1 point
```

### OMOS-235: Add tool_choice and tool_call_limit Parameters
```
Title: Add tool_choice and tool_call_limit parameters
Priority: Medium
Labels: backend, agno, agent, p2-medium
Estimate: 2 points
```

### OMOS-236: Add input_schema Validation
```
Title: Add input_schema validation to agent
Priority: Medium
Labels: backend, agno, agent, p2-medium
Estimate: 2 points
```

### OMOS-237: Add reasoning_agent and Step Config
```
Title: Add reasoning_agent and step configuration
Priority: Medium
Labels: backend, agno, agent, p2-medium
Estimate: 4 points
```

### OMOS-238: Add search_session_history Feature
```
Title: Add search_session_history feature
Priority: Medium
Labels: backend, agno, agent, p2-medium
Estimate: 4 points
```

### OMOS-239: Add cache_session Feature
```
Title: Add cache_session feature
Priority: Medium
Labels: backend, agno, agent, p2-medium
Estimate: 2 points
```

### OMOS-240: Add Team enable_agentic_memory
```
Title: Add enable_agentic_memory to teams
Priority: Medium
Labels: backend, agno, teams, p2-medium
Estimate: 2 points
```

### OMOS-241: Add Workflow Function Step
```
Title: Add workflow Function step type
Priority: Medium
Labels: backend, agno, workflows, p2-medium
Estimate: 4 points
```

### OMOS-242: Add Workflow Input Caching
```
Title: Add workflow input caching
Priority: Medium
Labels: backend, agno, workflows, p2-medium
Estimate: 4 points
```

### OMOS-243: Add ExaTools
```
Title: Add ExaTools (neural search)
Priority: Medium
Labels: backend, agno, tools, p2-medium
Estimate: 2 points
```

### OMOS-244: Add TavilyTools
```
Title: Add TavilyTools (research search)
Priority: Medium
Labels: backend, agno, tools, p2-medium
Estimate: 2 points
```

### OMOS-245: Add FirecrawlTools
```
Title: Add FirecrawlTools (web crawling)
Priority: Medium
Labels: backend, agno, tools, p2-medium
Estimate: 4 points
```

### OMOS-246: Add Newspaper4kTools
```
Title: Add Newspaper4kTools (article extraction)
Priority: Medium
Labels: backend, agno, tools, p2-medium
Estimate: 2 points
```

### OMOS-247: Add E2BTools
```
Title: Add E2BTools (sandboxed execution)
Priority: Medium
Labels: backend, agno, tools, p2-medium
Estimate: 4 points
```

### OMOS-248: Add TwilioTools
```
Title: Add TwilioTools (SMS/Voice)
Priority: Medium
Labels: backend, agno, tools, p2-medium
Estimate: 4 points
```

### OMOS-249: Add GoogleCalendarTools
```
Title: Add GoogleCalendarTools
Priority: Medium
Labels: backend, agno, tools, p2-medium
Estimate: 4 points
```

### OMOS-250: Add Memory Topic Classification
```
Title: Add memory topic classification
Priority: Medium
Labels: backend, agno, memory, p2-medium
Estimate: 4 points
```

### OMOS-251: Add Knowledge DELETE Endpoint
```
Title: Add DELETE /knowledge/{id} endpoint
Priority: Medium
Labels: backend, agno, knowledge, p2-medium
Estimate: 2 points
```

### OMOS-252: Add Agentic Chunking
```
Title: Add agentic chunking to knowledge
Priority: Medium
Labels: backend, agno, knowledge, p2-medium
Estimate: 6 points
```

### OMOS-253: Add MCP SSE Transport Testing
```
Title: Test and enable MCP SSE transport
Priority: Medium
Labels: backend, agno, mcp, p2-medium
Estimate: 2 points
```

### OMOS-254: Add Additional MCP Servers
```
Title: Add slack, google-maps, sequential-thinking MCP servers
Priority: Medium
Labels: backend, agno, mcp, p2-medium
Estimate: 4 points
```

### OMOS-255: Add Approval Webhooks
```
Title: Add approval webhooks for external notifications
Priority: Medium
Labels: backend, agno, hitl, p2-medium
Estimate: 4 points
```

### OMOS-256: Add Approval Timeout
```
Title: Add approval timeout/expiry
Priority: Medium
Labels: backend, agno, hitl, p2-medium
Estimate: 2 points
```

### OMOS-257: Add Agent-as-Judge Eval
```
Title: Add agent-as-judge evaluation type
Priority: Medium
Labels: backend, agno, evals, p2-medium
Estimate: 6 points
```

---

## P3 - Low Priority Issues (13)

### OMOS-258: Add system_message_role Parameter
```
Title: Add system_message_role parameter
Priority: Low
Labels: backend, agno, agent, p3-low
Estimate: 1 point
```

### OMOS-259: Add introduction Parameter
```
Title: Add introduction parameter
Priority: Low
Labels: backend, agno, agent, p3-low
Estimate: 1 point
```

### OMOS-260: Add parser_model Support
```
Title: Add parser_model support
Priority: Low
Labels: backend, agno, agent, p3-low
Estimate: 4 points
```

### OMOS-261: Add output_model Support
```
Title: Add output_model support
Priority: Low
Labels: backend, agno, agent, p3-low
Estimate: 4 points
```

### OMOS-262: Add show_tool_calls Parameter
```
Title: Add show_tool_calls parameter
Priority: Low
Labels: backend, agno, agent, p3-low
Estimate: 1 point
```

### OMOS-263: Add BraveSearchTools
```
Title: Add BraveSearchTools
Priority: Low
Labels: backend, agno, tools, p3-low
Estimate: 2 points
```

### OMOS-264: Add CrawlbaseTools
```
Title: Add CrawlbaseTools
Priority: Low
Labels: backend, agno, tools, p3-low
Estimate: 2 points
```

### OMOS-265: Add NotionTools
```
Title: Add NotionTools
Priority: Low
Labels: backend, agno, tools, p3-low
Estimate: 4 points
```

### OMOS-266: Add Pinecone Vector DB
```
Title: Add Pinecone vector database support
Priority: Low
Labels: backend, agno, knowledge, p3-low
Estimate: 4 points
```

### OMOS-267: Add Qdrant Vector DB
```
Title: Add Qdrant vector database support
Priority: Low
Labels: backend, agno, knowledge, p3-low
Estimate: 4 points
```

### OMOS-268: Add Mistral Model Provider
```
Title: Add Mistral model provider configuration
Priority: Low
Labels: backend, agno, models, p3-low
Estimate: 2 points
```

### OMOS-269: Add Cohere Model Provider
```
Title: Add Cohere model provider configuration
Priority: Low
Labels: backend, agno, models, p3-low
Estimate: 2 points
```

### OMOS-270: Add Custom Table Name Support
```
Title: Add custom table name support for storage
Priority: Low
Labels: backend, agno, storage, p3-low
Estimate: 2 points
```

---

## Summary

| Priority | Count | Total Points |
|----------|-------|--------------|
| P0 (Critical) | 8 | 42 |
| P1 (High) | 22 | 86 |
| P2 (Medium) | 25 | 78 |
| P3 (Low) | 13 | 33 |
| **Total** | **68** | **239** |

## Implementation Order

### Sprint 1 (Week 1-2): P0 Critical
- OMOS-203 to OMOS-210
- Focus: Knowledge auto-search, Team coordination, Workflow steps, MCP auto-attach
- **42 points**

### Sprint 2 (Week 3-4): P1 High (Part 1)
- OMOS-211 to OMOS-221
- Focus: Agent features, Team modes, Workflow improvements
- **46 points**

### Sprint 3 (Week 5-6): P1 High (Part 2)
- OMOS-222 to OMOS-232
- Focus: Tools with HITL, Evals
- **40 points**

### Sprint 4 (Week 7-8): P2 Medium
- OMOS-233 to OMOS-257
- Focus: Additional features, tools, integrations
- **78 points**

### Sprint 5 (Future): P3 Low
- OMOS-258 to OMOS-270
- Focus: Nice-to-have features
- **33 points**

---

## Script to Create Issues

A Python script to create these issues via Linear API is available at:
`scripts/linear/create_agno_gap_issues.py`

### Usage

```bash
# First, ensure LINEAR_API_KEY is set in .env or environment

# Dry run (see what would be created)
python scripts/linear/create_agno_gap_issues.py --dry-run

# Create P0 critical issues only
python scripts/linear/create_agno_gap_issues.py --priority p0

# Create P1 high priority issues
python scripts/linear/create_agno_gap_issues.py --priority p1

# Create ALL issues (68 total)
python scripts/linear/create_agno_gap_issues.py --priority all
```

### What the Script Does

1. **Creates Project:** "Agno Feature Gap Sprint" in Linear
2. **Creates Labels:** 18 labels (backend, agno, p0-critical, p1-high, etc.)
3. **Creates Issues:** All issues with proper priority, labels, and estimates
4. **Organizes:** All issues linked to the project for tracking

### Requirements

```bash
pip install requests python-dotenv
```

---

*Generated 2026-01-26*
*Reference: [AGNO-GAP-ANALYSIS.md](./AGNO-GAP-ANALYSIS.md)*
