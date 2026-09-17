---
title: "Chain of Thought"
type: reference
domain: SO
status: active
created: 2026-05-27
updated: 2026-05-27
tags: [concept, ai-agents, reasoning, llm, prompting]
---

# Chain of Thought (CoT)

## Definition
A prompting technique where a model reasons step-by-step before producing output. Rather than `input → answer`, the model does `input → internal reasoning → answer`. Improves accuracy on complex, multi-step tasks.

## Variants

| Variant | How It Works |
|---------|-------------|
| **Zero-shot CoT** | Append "think step by step" — no examples needed |
| **Few-shot CoT** | Provide worked reasoning examples in the prompt |
| **Tree of Thought (ToT)** | Branch multiple reasoning paths simultaneously, evaluate and pick best |
| **ReAct** | Interleave reasoning + action (tool calls) in a loop |

## In Multi-Agent Systems
CoT becomes the internal monologue of an agent before it decides:
- Which tool to call
- Whether to delegate to a subagent
- How to decompose a complex task
- Whether to ask for clarification or proceed

Legacy uses CoT-style reasoning internally when `reasoning_effort` is set (low/medium/high/xhigh in Hermes config).

## OneMind OS Angle
- **OM Cortex** would leverage CoT as its decision layer — reason before routing
- Visible CoT = `show_reasoning: true` in Hermes config (currently off for Zeus)
- Hooks can intercept the reasoning step before tool execution fires

## Related Concepts
- [[hook]] — hooks can fire at CoT checkpoints
- [[om-cortex]] — cortex decision-making runs on CoT
- [[chain-of-thought]] and [[hook]] together = observe + intercept reasoning mid-flight

## References
- Wei et al. 2022 — original CoT paper
- Yao et al. 2023 — Tree of Thought
- Hermes config: `agent.reasoning_effort`, `display.show_reasoning`

## Notes
<!-- Add prompting patterns, Hermes-specific CoT config notes here -->
