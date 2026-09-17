---
title: "TUI Builders — Research"
type: reference
domain: SO
status: active
created: 2026-05-26
updated: 2026-05-26
tags: [tui, terminal-ui, development, tools, python, rust, go, research]
---

# TUI Builders — Research

Terminal User Interface (TUI) frameworks for building interactive CLI applications.

---

## 1. Textual

**Website:** https://textual.textualize.io
**GitHub:** https://github.com/Textualize/textual
**Language:** Python
**By:** Textualize (Will McGugan)
**Stars:** ~26K+

### What it is
Python framework for building full-featured TUI apps. Runs in the terminal, supports mouse, keyboard, CSS-like styling, and reactive data. Same team behind `rich`.

### Key Features
- CSS-based styling system (Textual CSS)
- Reactive state management
- Full widget library: buttons, inputs, tables, trees, data grids, modals
- Mouse support
- Async-first (built on Python asyncio)
- Can run inline in the terminal or fullscreen
- Runs in browser via WebAssembly (Textual Web)
- Hot reload for CSS during development

### Best For
- Python-native TUI apps
- Dashboards, admin panels, interactive CLI tools
- Rapid prototyping with familiar Python patterns

### OneMind Relevance
Hermes uses Textual for its TUI. Directly relevant for building OneMind CLI dashboards and agent interfaces.

---

## 2. Bubble Tea

**Website:** https://github.com/charmbracelet/bubbletea
**Company Site:** https://charm.sh
**Language:** Go
**By:** Charmbracelet
**Stars:** ~29K+

### What it is
Go TUI framework based on The Elm Architecture (TEA) — a functional, message-passing model. Part of the Charm ecosystem which includes `lipgloss` (styling), `bubbles` (reusable components), and `glamour` (markdown rendering).

### Architecture (Elm Model)
```
Model → View → Message → Update → Model (loop)
```
- **Model:** application state
- **View:** renders model to string
- **Update:** handles messages, returns new model

### Key Features
- Functional/immutable state management
- `lipgloss` for layout and styling (borders, colors, alignment)
- `bubbles` component library (spinners, progress bars, text inputs, tables, viewports)
- `glamour` for markdown rendering in terminal
- Mouse support
- Highly composable — nest programs inside programs

### Best For
- Go CLI tools and TUIs
- Complex interactive terminal apps
- Production-grade TUIs with clean architecture

### Charm Ecosystem
- `lipgloss` — CSS-like styling for terminal output
- `bubbles` — reusable Bubble Tea components
- `glamour` — markdown in terminal
- `gum` — shell script TUI components (no code needed)
- `vhs` — record terminal GIFs
- `wishlist` — SSH app directory

---

## 3. Ratatui

**Website:** https://ratatui.rs
**GitHub:** https://github.com/ratatui/ratatui
**Language:** Rust
**Docs:** https://docs.rs/ratatui
**Crates:** https://crates.io/crates/ratatui
**Stars:** ~22K+

### What it is
Rust TUI library. Fork and continuation of `tui-rs`. Immediate-mode rendering — you redraw the entire UI each frame. Extremely fast, zero-cost abstractions.

### Key Features
- Sub-millisecond rendering
- Immediate-mode rendering model
- Rich widget library: charts, sparklines, tables, gauges, scrollable lists, progress bars, calendars, block layouts
- Constraint-based responsive layouts (adapts to any terminal size)
- Works in tiny tmux panes or large terminals
- Backends: crossterm (default), termion, termwiz
- Can embed in web via WASM

### Ecosystem
- `tui-textarea` — text input widget
- `tui-logger` — logging widget
- `color-eyre` / `anyhow` — error handling
- `tokio` / `async-std` for async

### Best For
- Performance-critical TUIs
- Rust-native tooling and system tools
- Dashboards, monitoring, games

### Apps Built with Ratatui
- `csvlens` — CSV viewer (3.8K stars)
- `oxker` — Docker container TUI (1.6K stars)
- Many terminal dashboards and monitoring tools

---

## 4. OpenTUI

**Website:** https://github.com/nicholasgasior/opentui (Go — appears inactive/archived)
**NPM:** https://www.npmjs.com/package/opentui (JavaScript version)
**Language:** Go / JavaScript variants

### What it is
Smaller, lesser-known TUI framework. The Go version (`nicholasgasior/opentui`) appears inactive — not found on pkg.go.dev. There is a JavaScript `opentui` package on npm.

### Status
- Much smaller community vs the other three
- Go version: likely abandoned
- Not recommended for new projects — Bubble Tea dominates the Go TUI space

---

## Comparison

| Framework | Language | Model | Performance | Ecosystem | Best For |
|-----------|----------|-------|-------------|-----------|----------|
| **Textual** | Python | Reactive/CSS | Good | Rich, Textualize | Python devs, rapid prototyping |
| **Bubble Tea** | Go | Elm (TEA) | Great | Charm ecosystem | Go devs, clean architecture |
| **Ratatui** | Rust | Immediate-mode | Fastest | crates.io | Performance, Rust-native |
| **OpenTUI** | Go/JS | — | Unknown | Minimal | Not recommended |

---

## OneMind OS Relevance

- **Hermes TUI** = Textual (Python) — already in use
- **Future Go agents or CLI tools** → Bubble Tea
- **Performance dashboards / fleet monitoring** → Ratatui (Rust)
- **Quick shell scripts with interactive UI** → Charm's `gum` (part of Bubble Tea ecosystem)

---

## Resources

- Textual docs: https://textual.textualize.io/
- Bubble Tea: https://github.com/charmbracelet/bubbletea
- Charm ecosystem: https://charm.sh
- Ratatui: https://ratatui.rs
- Ratatui docs: https://docs.rs/ratatui
