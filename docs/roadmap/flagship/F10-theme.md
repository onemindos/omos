# F10 — Theme and visual design system

Work order for one implementing agent. Branch: `flagship/f10-theme`.
Program index: [README](README.md).

## Outcome Zeus sees

- **Brand.** The portal looks like OneMind / Legacy: graphite and near-black
  surfaces with the **Legacy red** as the brand accent, instead of the current
  navy and teal ("it looks blue"). Both dark and light themes.
- **Consistency.** Cards, tables, chips and status colours look the same on
  every page.
- **Status colours.** Green/amber/red/grey mean only health, and stay distinct
  from the brand red. Critical uses a different red tone plus an icon, so the
  brand colour never reads as "error".

## Evidence

- `plugins/plugin-onemind-theme/src/designModule.tsx`: primary `#73d9e4` (teal)
  in dark and `#087e8b` in light; backgrounds `#121b29`, `#101822`, `#192432`
  (navy). `brand.css` overrides plus the `legacy-background-*.webp` image.
- The logo is red and black (`LogoFull.tsx`, `/brand/*`).
- The app mixes MUI v4 components and Backstage UI (`@backstage/ui` 0.18).
  Both need the same tokens.

## Build

1. **Palette.** Extract the exact brand red from the logo assets. Define tokens:
   surface 0–3 (neutral graphite, no blue tint), text 1–3, brand
   (red) plus hover/press, accent (neutral or warm), and status
   ok/warn/critical/unknown/info. Check WCAG AA contrast for text and focus rings
   (record the ratios in the PR).
2. **Apply tokens** to MUI (`createUnifiedTheme`/`UnifiedThemeProvider`) and to
   Backstage UI CSS variables (`--bui-*`), so both systems match. Remove the
   hard-coded hex values scattered in `brand.css` and plugin styles; grep for
   `#[0-9a-f]{6}` in `plugins/**` and replace them with tokens.
3. **Background image.** Make the Legacy background image optional (a user
   setting) and subtle. Default to a flat surface for readability.
4. **Density.** Tighter table rows and card padding (Zeus manages hundreds of
   entities). Keep touch targets at least 40px on mobile.
5. Add a `/design` admin page (or a Storybook-like route) showing all tokens and
   components in both themes for future reviews.

## Files

`plugins/plugin-onemind-theme/**`, style-only changes in other `plugins/**`,
`packages/app/public/brand/**` (assets only), and tests (snapshot of token
values).

## Acceptance

- [ ] No navy-tinted surfaces remain (the test asserts the token values;
      screenshots of Home, Catalog, an app page, Today, dark and light).
- [ ] AA contrast for body text, secondary text, links, focus rings and status
      chips (ratios listed in the PR).
- [ ] No hard-coded hex colours remain in plugin styles, except in the theme
      package (lint/test).
- [ ] Zeus approves the palette from screenshots before merge.
      Required gates pass.

## Status

(Implementing agent: write progress here.)
