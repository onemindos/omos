---
title: "AI Pointer — Research & Open Source Landscape"
type: reference
domain: SO
status: active
created: 2026-05-26
updated: 2026-05-26
tags: [ai, human-computer-interaction, computer-use, context-aware, pointer, deepmind, gemini, agents, open-source, macos]
---

# AI Pointer — Research & Open Source Landscape

## 1. DeepMind AI Pointer (Research / Experimental)

**Source:** https://deepmind.google/blog/ai-pointer/
**Published:** May 12, 2026
**Authors:** Adrien Baranes, Rob Marchant — Google DeepMind

### Summary

DeepMind is reimagining the mouse pointer for the AI era. The AI-enabled pointer, powered by Gemini, understands not just *where* you're pointing but *what* you're pointing at and *why it matters* — turning pixels into structured, actionable entities without forcing users into separate AI windows.

Goal: AI that meets users in their existing tools, not the reverse.

### Four Interaction Principles

1. **Maintain the Flow** — AI available across all apps, no "AI detours." Point at a PDF → get summary. Hover over table → get chart.
2. **Show and Tell** — Captures visual and semantic context around the pointer automatically. No verbose prompts.
3. **Embrace "This" and "That"** — Natural shorthand like "Fix this" / "Move that here" via pointer + speech + context.
4. **Pixels → Actionable Entities** — Understands *what* you're pointing at: places, dates, objects. Photo of a note → to-do list. Paused video frame → booking link.

**Status:** Experimental. Demos in Google AI Studio. Not a released product. Gemini-proprietary.

---

## 2. AIPointer — Open Source (Best Match)

**Repo:** https://github.com/gonemedia/aipointer
**Site:** https://aipointer.app
**Stars:** 125 | **Language:** TypeScript (Electron) | **License:** BSL-1.1
**Created:** May 13, 2026 — built in direct response to DeepMind's announcement

### What It Does

Open-source desktop overlay for macOS, Windows, Linux. Hold a hotkey (Right-Cmd on macOS), a glassmorphism box pops up next to the cursor, ask a question, vision-capable LLM answers about whatever's on screen around the cursor. BYO API key, local-first, no telemetry.

### Key Features

- **Multi-provider:** OpenRouter (recommended), Anthropic, OpenAI, Gemini — fallback chain handles outages
- **Vision-native:** Screenshots cursor region automatically, sends to LLM
- **Agentic tools:** fetch URL, open URL, copy text, save answer, launch app — all with approval gates
- **Finder auto-attach (v1.1.5):** Select files in Finder, press hotkey → auto-attaches up to 5 files
- **Voice in/out:** Mic input + TTS readback, hands-free loop mode
- **Region capture:** Hold trigger + drag to capture exact screen region
- **Chat-only mode:** Pure chat without screenshot auto-attach
- **Child Mode:** PIN-locked kid-friendly layer
- **Mouse-wiggle activation:** Wiggle cursor left-right-left to summon
- **Templates:** /summary, /brief, /translate, /explain, /code, /improve, /define, /solve, /reply, /identify

### Downloads

- macOS Apple Silicon: https://github.com/gonemedia/aipointer/releases/latest/download/AIPointer-arm64.dmg
- macOS Intel: https://github.com/gonemedia/aipointer/releases/latest/download/AIPointer-x64.dmg
- Windows: https://github.com/gonemedia/aipointer/releases/latest/download/AIPointer-Setup.exe
- Linux: https://github.com/gonemedia/aipointer/releases/latest/download/AIPointer-x86_64.AppImage

---

## 3. Other Open Source Analogues

| Project | Repo | Stars | What It Does |
|---------|------|-------|-------------|
| **OmniParser** | https://github.com/microsoft/OmniParser | ~10k | Screen region parsing → structured UI entities (Microsoft) |
| **screenpipe** | https://github.com/mediar-ai/screenpipe | ~12k | Continuous screen + audio capture with local AI context |
| **OpenAdapt** | https://github.com/OpenAdaptAI/OpenAdapt | ~1k | Open-source AI-powered RPA / computer use recording + replay |
| **UI-TARS** | https://github.com/bytedance-research/UI-TARS | ~5k | ByteDance GUI agent model for screen understanding (HF weights available) |
| **SoM (Set-of-Mark)** | https://github.com/microsoft/SoM | ~2k | Visual prompting — overlays marks on screen regions for LLM action |

---

## Relevance to OneMind

The pointer-as-context-bridge is the missing UI layer for agentic OS control. DeepMind's principles map directly to OneMind: AI that lives in your environment, not the reverse.

**AIPointer** is the most immediately deployable — works today on M2 Max, BYO key, supports Anthropic/OpenAI/Gemini/OpenRouter. Good candidate for Legacy integration evaluation.

**For deeper integration:** UI-TARS + OmniParser combo enables self-hosted screen understanding without cloud dependency.

## Task

- [ ] Legacy to evaluate AIPointer on Zeuss-MBP — install, test with OpenRouter key, assess fit for OneMind workflow

## Notes

