---
id: so-01-20260412000001
title: "JSON Context Profile"
type: reference
path: SO
domain: 00-24
folder: "01"
status: active
created: 2026-04-12
modified: 2026-04-12
owner: zeus
tags: [path/so, path/so/01, type/reference, status/active, prompt-engineering, json, llm, context]
ai_summary: "Structured JSON format used to give LLMs zero-noise, machine-readable context for consistent brand, image, or knowledge outputs"
---
# JSON Context Profile

Owner: Zeus
Tags: Prompt-Engineering, AI, JSON, LLM, Context

## What It Is

A **JSON Context Profile** is a structured, noise-free data format used to codify context — such as brand identity, user preferences, or image attributes — into a machine-readable blueprint that LLMs consume for consistent, high-fidelity results.

Unlike unstructured text (which introduces ambiguity), JSON profiles provide zero-noise access to key information, allowing AI to reason, write, or generate images with greater precision and repeatability.

---

## Core Use Cases

| Use Case | What It Does |
|----------|-------------|
| **Brand Voice** | Captures tone, keywords, and narrative structures for consistent marketing copy across teams |
| **Image Generation (AJCP)** | "Advanced JSON Context Profiles" — detailed blueprints for replicating specific styles, lighting, and compositions in AI art |
| **AI Knowledge Optimization** | Reformulates documents (team bios, service descriptions) into structured data for faster, more accurate LLM responses |

---

## How to Generate One

1. Prompt an AI to analyze a source (website, image, document)
2. Request output as a structured JSON object
3. Paste the JSON into future prompts to enforce specific constraints or styles
4. No need to re-describe the full context every time — the profile carries it

---

## Example Structure (Brand Voice)

```json
{
  "brand": {
    "name": "One Mind",
    "tone": ["direct", "empowering", "systems-focused"],
    "avoid": ["corporate jargon", "passive voice", "vague claims"],
    "keywords": ["AI-native", "human enhancement", "sovereign", "codex"],
    "narrative_arc": "problem → framework → transformation"
  }
}
```

---

## Example Structure (Image Generation — AJCP)

```json
{
  "style": "cinematic realism",
  "lighting": "golden hour, side-lit",
  "composition": "rule of thirds, subject left",
  "color_palette": ["#1a1a2e", "#e94560", "#0f3460"],
  "mood": "intense, focused",
  "negative_prompts": ["cartoon", "flat", "blurry"]
}
```

---

## One Mind Applications

- **Brand consistency** — JSON profile for One Mind voice used in marketing, content, and community posts
- **Agent context injection** — feed JSON profile to Legacy/Empire agents as system context
- **AJCP for content** — standardize visual style for One Mind brand assets

---

## Related

- [[prompt-engineering]]
- [[_Legacy AI]]
