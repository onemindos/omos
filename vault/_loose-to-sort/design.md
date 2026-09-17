---
version: alpha
name: OneMind OS
description: Sovereign, minimalist, command-and-control aesthetic. Black ops meets terminal precision.
colors:
  primary: "#cc1a1a"
  black: "#060606"
  panel: "#1a1a1a"
  border: "#333333"
  muted: "#666666"
  white: "#FFFFFF"
  red: "#cc1a1a"
  green: "#00ff00"
  orange: "#ffaa00"
  blue: "#0088cc"
typography:
  h1:
    fontFamily: SF Pro, Inter, -apple-system
    fontSize: 2rem
    fontWeight: 700
    lineHeight: 1.2
    letterSpacing: "-0.01em"
  h2:
    fontFamily: SF Pro, Inter, -apple-system
    fontSize: 1.5rem
    fontWeight: 700
    lineHeight: 1.3
  h3:
    fontFamily: SF Pro, Inter, -apple-system
    fontSize: 1.125rem
    fontWeight: 700
    lineHeight: 1.4
  body-md:
    fontFamily: SF Pro, Inter, -apple-system
    fontSize: 0.875rem
    fontWeight: 400
    lineHeight: 1.5
  body-sm:
    fontFamily: SF Pro, Inter, -apple-system
    fontSize: 0.75rem
    fontWeight: 400
    lineHeight: 1.4
  code:
    fontFamily: SF Mono, Menlo, Consolas
    fontSize: 0.8125rem
    fontWeight: 400
    lineHeight: 1.6
rounded:
  sm: 4px
  md: 8px
  lg: 12px
spacing:
  xs: 4px
  sm: 8px
  md: 16px
  lg: 24px
  xl: 32px
  2xl: 48px
  3xl: 64px
components:
  button-primary:
    backgroundColor: "{colors.red}"
    textColor: "{colors.white}"
    rounded: "{rounded.sm}"
    padding: 12px
  button-primary-hover:
    backgroundColor: "{colors.white}"
    textColor: "{colors.black}"
  button-secondary:
    backgroundColor: "transparent"
    textColor: "{colors.white}"
    rounded: "{rounded.sm}"
    padding: 12px
  button-disabled:
    backgroundColor: "{colors.border}"
    textColor: "{colors.muted}"
  input:
    backgroundColor: "{colors.panel}"
    textColor: "{colors.white}"
    rounded: "{rounded.sm}"
    padding: 12px
  input-focus:
    backgroundColor: "{colors.panel}"
    textColor: "{colors.white}"
  card:
    backgroundColor: "{colors.panel}"
    textColor: "{colors.white}"
    rounded: "{rounded.sm}"
    padding: 24px
  table-header:
    backgroundColor: "{colors.panel}"
    textColor: "{colors.white}"
  table-row:
    backgroundColor: "{colors.black}"
    textColor: "{colors.white}"
  table-row-alt:
    backgroundColor: "{colors.panel}"
    textColor: "{colors.white}"
  status-online:
    backgroundColor: "{colors.green}"
    textColor: "{colors.black}"
  status-warning:
    backgroundColor: "{colors.orange}"
    textColor: "{colors.black}"
  status-error:
    backgroundColor: "{colors.red}"
    textColor: "{colors.white}"
  status-unknown:
    backgroundColor: "{colors.muted}"
    textColor: "{colors.white}"
  status-info:
    backgroundColor: "{colors.blue}"
    textColor: "{colors.white}"
---

## Overview

OneMind OS is an AI-native life operating system. The visual identity reflects its core nature: sovereign, precise, minimal. The aesthetic is black ops meets terminal — every element is intentional, nothing decorative.

Two brand contexts exist under the same system:

- **OneMind OS** (`weareonemind.com`) — command & control, black/red/white strictly enforced
- **Zeus Personal Brand** (`zeusdelacruz.com`) — professional authority, flexible palette, human-first

When building for OneMind OS, apply this spec literally. When building for Zeus's personal brand, use the voice & tone guidelines but allow palette flexibility.

## Colors

- **Black (#060606):** The foundation. Backgrounds, main UI surfaces, full-bleed sections.
- **Panel (#1a1a1a):** Cards, sidebars, elevated surfaces — one step above background.
- **Border (#333333):** Dividers, input borders, structural lines. Never decorative.
- **Muted (#666666):** Secondary text, placeholders, metadata. Minimum AA contrast floor.
- **White (#FFFFFF):** Primary text, contrast. 21:1 on black — maximum legibility.
- **Red (#cc1a1a):** The only accent. Alerts, primary actions, focus rings, hover states. Do not introduce a second accent color.
- **Green (#00ff00):** Success and operational states only. Terminal aesthetic intentional.
- **Orange (#ffaa00):** Warning and pending states only.
- **Blue (#0088cc):** Informational/neutral states only.

Semantic mapping: Red = error/critical, Orange = warning/pending, Green = success/operational, Blue = info/neutral.

## Typography

Two font stacks only:

- **Sans-serif:** SF Pro → Inter → -apple-system (UI, headings, body)
- **Monospace:** SF Mono → Menlo → Consolas (terminals, code, data, metrics)

Use monospace for any live data: uptime counters, task counts, NATS subjects, log output. Use sans-serif for all editorial and UI copy.

Type scale is 8px-based. Line heights are generous for readability in dense data interfaces.

## Layout

- **12-column grid** for desktop layouts
- **8px base unit** — all spacing is a multiple of 8
- **Columns snap to grid** — no arbitrary pixel offsets
- Layouts are information-dense by design. Whitespace is functional, not decorative.

Standard spacing scale: 4 / 8 / 16 / 24 / 32 / 48 / 64px.

## Elevation & Depth

No shadows. No blur. No frosted glass. Depth is expressed through background color steps only:

- Layer 0: `#060606` — base
- Layer 1: `#1a1a1a` — cards, panels
- Layer 2: `#333333` — borders, active states

If a UI element needs to feel elevated, increase its background one step and add a 1px `#333333` border. Never use `box-shadow`.

## Shapes

Border radius is minimal — this is not a rounded, friendly UI:

- `4px` (sm) — buttons, inputs, tags
- `8px` (md) — cards, panels
- `12px` (lg) — modals, large containers

Never use `border-radius > 12px` or `border-radius: 50%` on anything except avatar circles.

## Components

### Buttons
- **Primary:** Red background, white text, 4px radius. Single call to action per view.
- **Secondary:** Transparent, red border, white text. Supporting actions.
- **Disabled:** #333 background, #666 text. No pointer events.

### Inputs
- Background #1a1a1a, border #333, white text, placeholder #666.
- Focus: border switches to red (#cc1a1a), no glow, no shadow.

### Tables
- Header: #1a1a1a background, red text.
- Rows: alternating #060606 / #1a1a1a.
- Hover: 2px left border in red.

### Cards
- #1a1a1a background, 1px #333 border, 24px padding, 4px radius.

### Status Indicators
- 🟢 Green — Online, operational, healthy
- 🔴 Red — Offline, critical, failed
- 🟡 Orange — Warning, degraded, pending
- ⚪ Gray — Unknown, disabled, dormant

### Command Palette
- 95% opacity black overlay.
- Red highlight on selected item.
- Terminal green for active/confirmed commands.
- Fuzzy search with instant feedback.

## Do's and Don'ts

**Do:**
- Use a single accent (red) per screen. One call to action.
- Use monospace for all live data and metrics.
- Keep animations under 200ms. Fast = confident.
- Use status indicators consistently (same emoji/color system everywhere).
- Follow the 8px spacing grid.

**Don't:**
- Add gradients, shadows, or blur effects.
- Introduce new accent colors.
- Use rounded corners > 12px.
- Add decorative animation or transitions.
- Use light backgrounds. This is a dark-only system.
- Mix sans-serif and monospace in the same data cell.
