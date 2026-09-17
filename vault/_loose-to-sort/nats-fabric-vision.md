---
title: "NATS Fabric Vision (Interactive HTML)"
type: artifact-card
status: active
created: 2026-06-11
updated: 2026-06-11
owner: zeus
category: infrastructure
tags: [nats, fabric, visualization, html, interactive, architecture, ui-reference]
ai_summary: "Self-contained interactive HTML visualization of the OneMind NATS fabric. Six animated views (Fabric, Anatomy, Languages, Voice, Topology, SaaS) built with canvas particle simulation. Opens in any browser, zero dependencies, works offline. Serves as the visual reference for understanding the architecture and as a design seed for the OneMind Command UI."
artifact_path: "/Users/onemind/onemind/nats-fabric-vision.html"
---

# NATS Fabric Vision — Interactive HTML

> The OneMind NATS fabric made visual and tangible. One self-contained HTML file, six explorable views, real particle-simulation animation. The "see everything and understand what we're building" artifact.

## What It Is

A standalone, dependency-free HTML page that visualizes NATS as the **nervous system** of OneMind OS — every machine (agents, robotics, satellites, drones, edge devices, the hub, cloud, human biometrics, nature) as a node on one fabric. Built to make the architecture *understandable at a glance* and to seed the design language for the OneMind Command UI.

- **File:** `nats-fabric-vision.html` (vault root, ~46 KB)
- **Dependencies:** none (Google Fonts only, degrades gracefully offline)
- **Tech:** pure HTML/CSS/JS, HTML5 Canvas particle simulation, no framework
- **Render note:** Obsidian shows it in the file tree but **will not preview it** — Obsidian only renders Markdown. Open it in a browser to see it live.

## How To Open

```bash
# Brave (Zeus's primary browser)
open -a "Brave Browser" ~/onemind/nats-fabric-vision.html

# Or any browser via file URL
file:///Users/onemind/onemind/nats-fabric-vision.html
```

## The Six Views

| View | What it shows |
|---|---|
| **◉ The Fabric** | Animated core with 9 node-classes orbiting. Real messages flow node → core → fan out (M:N delivery). Live counters, streaming subject feed, click any node for its real subject namespace. |
| **⊞ Anatomy** | How one message lives: Publisher → Subject → Core fan-out → 3 subscribers, with JetStream capture, KV, and Object Store branching. Explains Core vs JetStream vs KV vs Object Store + the `$SRV` self-discovery point. |
| **⟨⟩ Languages** | Go / Python / TypeScript / Dart side by side with real code, why each is chosen, what each runs. Edge=Go, Brain=Python, Cockpit=TypeScript, App=Dart. |
| **🎙 Voice** | The 5-stage realtime pipeline: Mic → LiveKit Room → Python Bridge (dual NATS+LiveKit) → Intent on bus → TTS back. Shows subjects fired at each stage. |
| **⊹ Topology** | Physical map: Homestead ⇄ Hub ⇄ Office ⇄ Cloud with live traffic across the Tailscale mesh. Leaves survive offline; cloud is the public doorway. |
| **⇄ SaaS** | The 5-stage SaaS flow (Skool webhook → Wombat → bus event → agent → ClickHouse) proving "no API calls in your code, ever." |

## Why It Matters

1. **Comprehension** — turns the abstract "nats is the nervous system" doctrine into something Zeus can see move.
2. **UI seed** — the visual language (dark substrate, cyan/teal core, color-coded message classes, node treatment) is the starting palette for OneMind Command's NATS management pages.
3. **Communication** — a single file that explains the whole architecture to anyone, hostable on GitHub Pages later (see below).

## Hosting Live (Future)

This file is **GitHub-Pages-ready** — fully static, zero build step. The vault repo (`One-Mind`) is private, so do **not** enable Pages on it (would expose the whole vault + needs paid plan). Instead:

1. Create a separate **public** repo (e.g. `onemind-visuals`)
2. Drop HTML artifacts there → Settings → Pages → enable
3. Live URL: `https://onemindos.github.io/onemind-visuals/nats-fabric-vision.html`

## Relationship to the Real System

The visualization is **illustrative**, not wired to live NATS. The subjects shown (`tel.barn.soil`, `cmd.go2.patrol.start`, `agents.prompt.legacy.zeus.main`, etc.) follow the real Fabric subject law but the animation is a simulation. The *next* step — turning the Topology/Fabric view into a real management page wired to `$SRV.PING` / `$SYS` / KV watch — is tracked as the OneMind Command "NATS management page."

## Related

- [[ONEMIND-OS-STACK]] — the full infrastructure stack this visualizes
- [[NATS-INFRASTRUCTURE]] — the live NATS server documentation
- `~/onemindos-core/docs/NATS-DOCTRINE.md` — binding subject law
