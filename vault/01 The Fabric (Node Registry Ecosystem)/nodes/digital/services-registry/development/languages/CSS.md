# CSS

> **Cascading Style Sheets — visual presentation layer of the web**

## Overview

CSS styles all web UIs in the One Mind stack. onemind-web uses Tailwind CSS v4 with a `--neutral-tint` single-knob theme system and custom Solarpunk palette tokens (`--sp-*`, `--td-*`). codex-live (Quartz) has its own CSS overrides.

## Status

| Attribute   | Value                                  |
| ----------- | -------------------------------------- |
| Status      | ✅ Active                               |
| Category    | Development / Web                      |
| Type        | Styling Language                       |
| Integration | Tailwind CSS, PostCSS, CSS Modules     |

## Use Cases in OneMind

1. **onemind-web UI** — Tailwind classes throughout React components
2. **Brand Theming** — Solarpunk palette tokens defined in `globals.css`
3. **codex-live** — Custom CSS overrides for the Quartz-generated site
4. **Animations** — Framer Motion + CSS transitions for UI interactions
5. **Responsive Design** — Mobile-first breakpoints for all views

## Theme Tokens (onemind-web)

```css
/* globals.css — Solarpunk Palette */
--sp-primary: ...;      /* Main brand color */
--sp-path-*: ...;       /* Domain path colors */
--td-*: ...;            /* Tactical Dark (Zeus Red + Cyan) */
--ts-*: ...;            /* Tactical Solarpunk (Amber + Moss) */
--cat-*: ...;           /* Agent category colors */
--neutral-tint: ...;    /* Single-knob theme control */
```

## Tailwind (onemind-web)

```bash
# Tailwind v4 — config-less, CSS-first
# postcss.config.mjs handles compilation
npm run dev   # Hot reload includes CSS changes
```

## Links

- MDN: [developer.mozilla.org/en-US/docs/Web/CSS](https://developer.mozilla.org/en-US/docs/Web/CSS)
- Tailwind CSS: [tailwindcss.com](https://tailwindcss.com)


---

## Related

- [[00-24 SO (Sovereign Operations)/02 Infrastructure (VPS, Servers & Networking) (The Body)/02 Infrastructure|02 Infrastructure]]
