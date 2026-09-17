# Open claw memory stack 

**

# 1\. Memory Plugin Explanations

**

**

# memory-core — OpenClaw's built-in structured long-term memory engine. Works like REM sleep: at 0 3 * * * (3AM ET), it runs a&quotdreaming;" consolidation pass over recent sessions, extracts important facts, and writes them to a persistent store. Zero embeddings, zero CPU cost during the day. This is your primary durable memory — the one that actually survives restarts and accumulates over time.

**

**

#   


**

**

# memory-wiki — Organizes and surfaces memories as wiki-style topics. When an agent needs to recall something, it queries this layer like a mini-encyclopedia. It structures what memory-core stores into navigable, topic-indexed entries.

**

**

#   


**

**

# active-memory — Working memory for the current conversation only. Tracks context in-session (what was said earlier this chat), makes it available for context injection. Ephemeral — it doesn't persist between sessions independently. Think RAM vs disk.

**

**

#   


**

**

# Together they form a complete stack: active-memory (RAM) ? memory-core dreaming (writes to disk) ? memory-wiki (indexes the disk for fast lookup).

**