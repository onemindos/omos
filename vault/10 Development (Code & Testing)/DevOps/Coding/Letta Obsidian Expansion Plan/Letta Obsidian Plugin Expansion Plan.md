---
id: ui-12-20260111053926
title: "Letta Obsidian Plugin Expansion Plan"
type: concept
domain: SO
domain: 00-24
folder: "00-24 SO (Sovereign Operations)"
status: active
created: 2026-01-11
modified: 2026-01-11
owner: zeus
tags: [domain/so, folder/12]
todoist_project: "00-24 SO (Sovereign Operations)"
ai_summary: "Letta Obsidian Plugin Expansion Plan"
---


## Current State Summary

The plugin currently provides:

- ✅ Real-time chat with stateful Letta agents
- ✅ Memory block management (CRUD + search/attach)
- ✅ Focus mode with rich note metadata
- ✅ Agent switching across projects
- ✅ File mention system with autocomplete
- ✅ Streaming/non-streaming message support
- ❌ **Disabled**: Custom tool registration (due to API limitations)
- ❌ **Removed**: Vault-wide file syncing

---

## 🎯 Expansion Themes

### **Theme 1: Enhanced Agent Capabilities**

_Make the agent more powerful and useful for knowledge work_

### **Theme 2: Deep Vault Integration**

_Leverage Obsidian's unique features (graph, canvas, dataview, etc.)_

### **Theme 3: Collaboration & Workflows**

_Multi-agent systems, templates, and automation_

### **Theme 4: Knowledge Management**

_Better memory, archiving, and knowledge retrieval_

### **Theme 5: Developer Experience**

_Advanced features for power users_

---

## 📋 Prioritized Feature Roadmap

### 🔥 **HIGH PRIORITY** (High Value, Medium-High Feasibility)

#### 1. **Re-enable Custom Tools with Better Implementation**

**Status**: Framework exists, currently disabled  
**What**: Fix and re-enable the `write_obsidian_note` tool with proper approval handling  
**Value**: Allows agent to create notes, major productivity boost  
**Implementation**:

- Fix approval flow to work with streaming API limitations
- Add more Obsidian-specific tools:
    - `create_note`: Create notes with frontmatter, tags, templates
    - `update_note`: Modify existing notes (with approval)
    - `create_task`: Add tasks to daily notes or task management notes
    - `search_vault`: Semantic search across vault
    - `link_notes`: Create connections between notes
- Implement robust consent/approval UI
- Add tool usage analytics to settings

**Technical Notes**: Lines 1322-1507 have commented implementation to build from

---

#### 2. **Semantic Search & Retrieval**

**Status**: New capability  
**What**: Let the agent semantically search the vault beyond simple file mentions  
**Value**: Agent can find relevant context without user manually mentioning files  
**Implementation**:

- Use Letta's embedding models API (`GET /v1/models/embedding`)
- Build vector index of vault notes (stored locally or in Letta)
- Add `@search:query` syntax for natural language searches
- Display search results as attachable context
- Auto-attach relevant notes based on conversation context
- Command: "Find notes related to..." with UI picker

**Obsidian Integration**: Works with existing metadata cache for fast filtering

---

#### 3. **Canvas Integration**

**Status**: New capability  
**What**: Agent can interact with Obsidian Canvas boards  
**Value**: Visual thinking, project planning, mind mapping with AI assistance  
**Implementation**:

- Read canvas JSON structure
- Tool: `analyze_canvas` - Get overview of canvas structure
- Tool: `add_canvas_card` - Add notes/text to canvas
- Tool: `organize_canvas` - Suggest organization based on content
- Command: "Chat about this canvas" (context-aware)
- Focus mode extension: Include active canvas in focus block

**Use Cases**:

- Project planning assistance
- Brainstorming sessions
- Visual knowledge graphs
- Meeting notes organization

---

#### 4. **Conversation History & Context Management**

**Status**: New capability  
**What**: Persist and manage chat history across sessions  
**Value**: True long-term conversation continuity within Obsidian  
**Implementation**:

- Store conversation history in `.letta/conversations/` folder as markdown
- UI: Load previous conversations from dropdown/command
- Export conversations as notes with metadata
- Search across conversation history
- "Continue last conversation" quick action
- Configurable history retention (last N conversations, time-based)

**Technical**: Use Obsidian's vault API for storage, keep separate from Letta's message history

---

#### 5. **Smart Daily Notes Integration**

**Status**: New capability  
**What**: Agent awareness and interaction with daily notes workflow  
**Value**: Huge for users who live in daily notes (journaling, logging, task management)  
**Implementation**:

- Detect daily note format from settings or popular plugins (Periodic Notes, Calendar)
- Tool: `append_to_daily_note` - Add content to today's note
- Tool: `create_daily_summary` - Summarize today's activities
- Command: "Add to daily notes" - Quick capture
- Auto-attach today's note as context in conversations
- Template support for structured entries

**Integrations**: Works with Templater, Periodic Notes, Calendar plugins

---

### 🟡 **MEDIUM PRIORITY** (Medium-High Value, Variable Feasibility)

#### 6. **Multi-Agent Collaboration**

**Status**: New capability  
**What**: Multiple agents working together on tasks  
**Value**: Specialized agents for different domains (research, writing, coding, etc.)  
**Implementation**:

- Create "agent team" configurations
- UI: Agent panel showing active agents
- Tool: `consult_agent` - One agent asks another for help
- Workflow: User orchestrates multi-agent conversations
- Templates: Pre-configured teams (Research Team, Writing Team, etc.)
- Message routing: "@AgentName" mentions to direct messages

**Use Cases**:

- Research agent finds sources → Writing agent synthesizes
- Code agent reviews → Documentation agent explains
- Brainstorm agent generates → Critic agent evaluates

---

#### 7. **Advanced File Sync (Reimplemented)**

**Status**: Deprecated/removed, needs redesign  
**What**: Selective, intelligent vault syncing to Letta  
**Value**: Agent has broader vault awareness  
**Implementation**:

- **Selective sync**: Folder/tag-based rules (e.g., "sync #project/active")
- **Change detection**: Only sync modified files (hash-based)
- **Sync strategies**:
    - Manual: User-triggered
    - Auto: On file save (debounced)
    - Scheduled: Background sync every N minutes
- **Sync UI**: Show sync status, file count, last sync time
- **Conflict handling**: Detect external changes
- **Exclude patterns**: `.obsidian/`, templates, etc.

**Alternative Approach**: Instead of "syncing," use on-demand retrieval with caching

---

#### 8. **Knowledge Graph Analysis**

**Status**: New capability  
**What**: Agent understands and navigates vault's link structure  
**Value**: Helps users discover connections and insights  
**Implementation**:

- Tool: `analyze_graph` - Get graph metrics (clusters, hubs, orphans)
- Tool: `find_path` - Find connections between notes
- Tool: `suggest_links` - Recommend links based on content
- Command: "What's connected to this note?"
- Visualization: Show subgraph of current context
- Focus mode: Include graph context (already partially implemented)

**Obsidian API**: Uses `app.metadataCache.resolvedLinks` and `app.metadataCache.getBacklinksForFile()`

---

#### 9. **Template System for Prompts & Workflows**

**Status**: New capability  
**What**: User-defined templates for common agent interactions  
**Value**: Streamlines repetitive workflows  
**Implementation**:

- Slash commands in chat: `/summarize`, `/brainstorm`, `/review`
- Template marketplace or sharing
- Template editor UI in settings
- Variables: `{{title}}`, `{{selection}}`, `{{daily_note}}`
- Chain templates for complex workflows
- Community templates repository

**Examples**:

- `/summarize` → "Summarize this note in 3 bullet points"
- `/meeting` → "Convert these notes into structured meeting minutes"
- `/expand` → "Expand this outline into full paragraphs"

---

#### 10. **Dataview Integration**

**Status**: New capability  
**What**: Agent can query vault using Dataview-like queries  
**Value**: Structured data access for power users  
**Implementation**:

- Tool: `query_vault` - Execute Dataview-style queries
- Parse frontmatter across notes
- Aggregate data (e.g., "all books I've read this year")
- Generate reports based on metadata
- Create dynamic dashboards

**Requires**: Dataview plugin API or custom implementation

---

#### 11. **Memory Block Enhancements**

**Status**: Enhancement to existing  
**What**: More powerful memory management  
**Implementation**:

- **Versioning**: Track changes to memory blocks over time
- **Diff view**: Show what changed in block updates
- **Templates**: Pre-configured memory structures
- **Import/Export**: Share memory configurations
- **Auto-summarization**: Condense old memories to save space
- **Memory analytics**: Show which blocks are most referenced
- **Memory graph**: Visualize relationships between memory blocks

---

#### 12. **Voice Input Integration**

**Status**: New capability  
**What**: Voice-to-text for chat messages  
**Value**: Hands-free interaction, accessibility  
**Implementation**:

- Browser Speech Recognition API (Web Speech API)
- Whisper API integration (if user provides OpenAI key)
- Push-to-talk button in chat interface
- Auto-transcription to message input
- Optional: Voice responses using TTS

**Mobile consideration**: Especially valuable on mobile Obsidian

---

### 🟢 **LOWER PRIORITY / FUTURE** (Nice-to-Have, Lower Urgency)

#### 13. **Scheduled Agent Tasks**

**Status**: New capability  
**What**: Agent performs tasks on schedule  
**Examples**:

- Daily summary of changes at 6 PM
- Weekly review preparation
- Monthly archive organization
- Automated tagging suggestions

**Implementation**: Obsidian's background timers + task queue

---

#### 14. **Agent Training & Fine-Tuning UI**

**Status**: New capability  
**What**: Configure agent personality, examples, and behavior  
**Implementation**:

- System prompt templates
- Few-shot example management
- Response style configuration
- Custom persona creation

**Note**: Depends on Letta's fine-tuning capabilities

---

#### 15. **Conflict Resolution for Collaborative Editing**

**Status**: New capability  
**What**: Handle simultaneous edits by user and agent  
**Implementation**: Git-style merge with 3-way diff

---

#### 16. **Mobile-Optimized Interface**

**Status**: Enhancement  
**What**: Better mobile UX for chat and memory management  
**Implementation**: Responsive design, touch-optimized controls

---

#### 17. **Export & Sharing**

**Status**: New capability  
**What**: Share conversations, export agent configurations  
**Implementation**:

- Export conversations as markdown
- Share memory blocks
- Agent configuration export/import (JSON)

---

#### 18. **Analytics & Insights**

**Status**: New capability  
**What**: Usage stats, most-discussed topics, interaction patterns  
**Implementation**: Local analytics dashboard in settings

---

## 🎨 **Experimental / Research Ideas**

#### 19. **Real-Time Collaborative Notes**

**What**: Agent and user edit the same note simultaneously  
**How**: Operational transforms or CRDTs for conflict-free editing

#### 20. **Obsidian Publish Integration**

**What**: Agent helps prepare notes for publishing  
**How**: Suggest improvements, check links, optimize for readers

#### 21. **Plugin Ecosystem Integration**

**What**: Integrate with popular plugins (Tasks, Kanban, Excalidraw)  
**How**: Tool APIs for each plugin

#### 22. **Vault Health Monitoring**

**What**: Agent suggests vault improvements (broken links, orphans, duplicates)  
**How**: Background analysis + recommendations

---

## 🛠️ **Implementation Strategy**

### Phase 1: Foundation (Months 1-2)

1. ✅ Re-enable custom tools with fixed approval flow
2. ✅ Implement conversation history
3. ✅ Add semantic search capability

### Phase 2: Deep Integration (Months 3-4)

4. ✅ Canvas integration
5. ✅ Daily notes integration
6. ✅ Knowledge graph analysis

### Phase 3: Advanced Features (Months 5-6)

7. ✅ Multi-agent collaboration
8. ✅ Template system
9. ✅ Memory enhancements

### Phase 4: Polish & Expansion (Ongoing)

10. ✅ Voice input
11. ✅ Mobile optimization
12. ✅ Plugin integrations
13. ✅ Community features

---

## 🎯 **Quick Wins** (Implement First for Maximum Impact)

1. **Re-enable `write_obsidian_note` tool** - Framework exists, big value
2. **Conversation history** - Use Obsidian's vault, simple storage
3. **Daily notes append** - Single tool, high usage
4. **Canvas reading** - JSON parsing, easy to start
5. **Slash commands in chat** - UI sugar, great UX improvement

---

## 💡 **Key Decision Points**

### Should we prioritize:

1. **Breadth** (many small features) vs **Depth** (few amazing features)?
2. **Power users** (advanced features) vs **Beginners** (simplicity)?
3. **Letta-first** (leverage Letta's strengths) vs **Obsidian-first** (deep integration)?
4. **Stability** (polish existing) vs **Innovation** (experimental features)?

### Technical questions:

- **Storage**: Where to store plugin data? (Vault vs Letta vs local DB)
- **Performance**: How to handle large vaults? (Indexing, caching strategies)
- **Privacy**: What data should stay local vs sync to Letta?