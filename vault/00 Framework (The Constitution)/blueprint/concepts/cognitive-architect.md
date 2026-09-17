---
title: "Cognitive Architect"
type: reference
domain: SO
status: active
created: 2026-05-27
updated: 2026-05-27
tags: [ai, agent, architecture, role, concept, cognition, cognitive-science]
---

# Cognitive Architect

## What It Is
Cognitive Architect is not a standardized technical term with a formal definition in AI literature — it is an emerging role/discipline at the intersection of cognitive science and AI system design. It describes someone who designs the internal structure of intelligent systems: how they remember, reason, plan, and act — based on principles drawn from cognitive science.

The established academic field it draws from is **Cognitive Architecture** — computational models of the structure and processes of the human mind.

## Foundational Cognitive Architecture Frameworks

| Framework | Origin | Core Mechanism |
|-----------|--------|----------------|
| ACT-R (Adaptive Control of Thought-Rational) | Carnegie Mellon University | Production rules + declarative memory modules; models human problem-solving and learning |
| SOAR (State, Operator, And Result) | University of Michigan | Problem spaces and operators; decision-making via goal decomposition; learns via chunking |
| CLARION | UMASS Amherst | Dual-process: integrates implicit (neural) and explicit (symbolic) learning |
| EPIC (Executive Process/Interactive Control) | University of Michigan | Models human multitasking and attention |
| LIDA (Learning Intelligent Distribution Agent) | U of Memphis | Based on Global Workspace Theory; integrates perception, attention, and action selection |

## What a Cognitive Architect Designs

| Layer | What It Shapes |
|-------|---------------|
| Memory Systems | Working memory (short-term context), long-term memory (knowledge bases, vector stores), episodic memory (experience replay) |
| Reasoning Mechanisms | Symbolic (rule-based), subsymbolic (neural), hybrid (neuro-symbolic) |
| Planning & Goal Management | Goal decomposition, plan generation, meta-cognitive monitoring |
| Learning & Adaptation | Reinforcement learning, transfer learning, meta-learning |
| Perception-Action Cycles | Sense-think-act loops, reactive vs deliberative processing, action selection |

## Modern Implementations
Classic academic frameworks (ACT-R, SOAR) are primarily research tools. In modern AI agent systems, cognitive architecture principles appear in:
- **LangChain / LangGraph** — agent orchestration with memory and tool use
- **AutoGPT / BabyAGI** — autonomous goal decomposition agents
- **ReAct** (Reasoning + Acting) — interleaved reasoning and action prompting pattern
- **BDI (Belief-Desire-Intention)** architectures in multi-agent systems

## OneMind OS Angle
In the OneMind fleet, this role maps to the function of designing how agents think — their memory schemas, skill structures, decision frameworks, and context injection patterns. Currently absorbed informally by Legacy. A dedicated Cognitive Architect function would own vault schema design (SCHEMA.md), agent profile templates (_roster/), memory bank structure, and cross-agent knowledge graph design.

## Related Concepts
- [[onemind-librarian]] — downstream executor of the schemas this role designs
- [[SO/02 Intelligence & Agents/fleet/]] — the fleet whose architecture this role shapes

## Notes
<!-- Decide: standalone agent role vs Legacy sub-function vs design discipline only -->
