/**
 * OneMind design tokens: the single source of truth for colour.
 *
 * Pure data, no React, so any plugin (and the tests) can import it through
 * `@internal/plugin-onemind-theme/tokens`. `brand.css` carries the same values
 * as `--onemind-*` CSS variables for the CSS side; `tokens.test.ts` fails if the
 * two drift apart.
 *
 * Brand red: `#ff333a`, sampled from the pixels of the logo
 * (`packages/app/public/brand/onemind-logo-dark.png`; the light JPEG is
 * `#fe343a`). It is the exact brand fill in dark mode. In light mode white text
 * on `#ff333a` is only 3.6:1, so the light fill is a deeper shade of the same
 * hue (`#d4151c`, 5.4:1 with white).
 */
export type ThemeMode = 'dark' | 'light';

export type StatusKey = 'ok' | 'warn' | 'critical' | 'unknown' | 'info';

export interface StatusTokens {
  /** Text, icon and dot colour on any surface. */
  fg: string;
  /** Tinted chip or banner background; `fg` on it is AA. */
  bg: string;
}

export interface ModeTokens {
  /** Neutral graphite, no hue. 0 = app, 1 = card, 2 = raised/hover, 3 = input/selected. */
  surface: { 0: string; 1: string; 2: string; 3: string };
  border: { subtle: string; strong: string };
  text: { 1: string; 2: string; 3: string };
  brand: {
    /** Fill for buttons, indicators, the active nav item. */
    base: string;
    hover: string;
    press: string;
    /** Text and icon colour on a `brand` fill. */
    on: string;
    /** Brand red as text or a link on a surface (AA). */
    text: string;
    /** Faint brand tint for selected rows and highlights. */
    subtle: string;
  };
  /** Neutral warm accent for secondary emphasis. Never a status colour. */
  accent: { base: string; on: string };
  /** Keyboard focus ring. At least 3:1 against every surface. */
  focus: string;
  /** Health only. Critical is a rose red, deliberately not the brand red. */
  status: Record<StatusKey, StatusTokens>;
  /** The sidebar stays graphite in both modes (the logo is white on dark). */
  nav: { bg: string; hover: string; text: string; selected: string };
  shadow: string;
}

export const tokens: Record<ThemeMode, ModeTokens> = {
  dark: {
    surface: { 0: '#0e0e0f', 1: '#151516', 2: '#1c1c1e', 3: '#252527' },
    border: { subtle: '#2e2e31', strong: '#4d4d52' },
    text: { 1: '#f3f3f4', 2: '#b8b8bd', 3: '#93939a' },
    brand: {
      base: '#ff333a',
      hover: '#ff5258',
      press: '#f0272e',
      on: '#120305',
      text: '#ff6b70',
      subtle: '#2a1316',
    },
    accent: { base: '#d9c5a0', on: '#17140d' },
    focus: '#ff6b70',
    status: {
      ok: { fg: '#3ecf72', bg: '#12261a' },
      warn: { fg: '#f5b83d', bg: '#2b2210' },
      critical: { fg: '#ff5fa0', bg: '#2f1224' },
      unknown: { fg: '#9a9aa0', bg: '#232325' },
      info: { fg: '#6cb6ff', bg: '#14212f' },
    },
    nav: {
      bg: '#0a0a0b',
      hover: '#1c1c1e',
      text: '#c4c4c9',
      selected: '#ffffff',
    },
    shadow: '#00000066',
  },
  light: {
    surface: { 0: '#f4f4f5', 1: '#ffffff', 2: '#ececee', 3: '#e4e4e6' },
    border: { subtle: '#d9d9dd', strong: '#77777d' },
    text: { 1: '#19191b', 2: '#4b4b50', 3: '#5f5f65' },
    brand: {
      base: '#d4151c',
      hover: '#b8121a',
      press: '#9c0f16',
      on: '#ffffff',
      text: '#b8121a',
      subtle: '#fdeaeb',
    },
    accent: { base: '#7a5f2e', on: '#ffffff' },
    focus: '#b8121a',
    status: {
      ok: { fg: '#167a3e', bg: '#e4f5ea' },
      warn: { fg: '#8a5a00', bg: '#fbefd2' },
      critical: { fg: '#b0105a', bg: '#fde4ef' },
      unknown: { fg: '#5f5f65', bg: '#ececee' },
      info: { fg: '#1b5fa8', bg: '#e2eefa' },
    },
    nav: {
      bg: '#111112',
      hover: '#222224',
      text: '#c4c4c9',
      selected: '#ffffff',
    },
    shadow: '#00000024',
  },
};

export interface CategoryColor {
  /** Fill / border, readable on light surfaces. */
  base: string;
  /** Text and border on dark surfaces. */
  bright: string;
}

/**
 * Identity hues for app categories on the launcher. Decorative only: none is
 * red or green, so they cannot be mistaken for the brand or for health.
 */
export const categoryPalette = {
  blue: { base: '#1970d4', bright: '#80b8ff' },
  violet: { base: '#8250c4', bright: '#c5a2ff' },
  teal: { base: '#007f76', bright: '#67ddd0' },
  amber: { base: '#a96816', bright: '#ffce83' },
  orange: { base: '#b65b18', bright: '#ffb880' },
  indigo: { base: '#5069c5', bright: '#a7b7ff' },
  cyan: { base: '#237f91', bright: '#77cfe4' },
  moss: { base: '#4f7a2e', bright: '#b5d98a' },
  plum: { base: '#9b4a8a', bright: '#f0a8e0' },
  slate: { base: '#68758b', bright: '#b8c5dd' },
} satisfies Record<string, CategoryColor>;

/** Prefix of every CSS custom property this package defines. */
export const CSS_PREFIX = '--onemind-';

const statusCss: Record<StatusKey, string> = {
  ok: 'ok',
  warn: 'warn',
  critical: 'critical',
  unknown: 'unknown',
  info: 'info',
};

/** The `--onemind-*` variables for one mode, in a stable order. */
export function cssVariables(mode: ThemeMode): Record<string, string> {
  const t = tokens[mode];
  const v: Record<string, string> = {};
  ([0, 1, 2, 3] as const).forEach(i => {
    v[`${CSS_PREFIX}surface-${i}`] = t.surface[i];
  });
  v[`${CSS_PREFIX}border`] = t.border.subtle;
  v[`${CSS_PREFIX}border-strong`] = t.border.strong;
  ([1, 2, 3] as const).forEach(i => {
    v[`${CSS_PREFIX}text-${i}`] = t.text[i];
  });
  v[`${CSS_PREFIX}brand`] = t.brand.base;
  v[`${CSS_PREFIX}brand-hover`] = t.brand.hover;
  v[`${CSS_PREFIX}brand-press`] = t.brand.press;
  v[`${CSS_PREFIX}on-brand`] = t.brand.on;
  v[`${CSS_PREFIX}brand-text`] = t.brand.text;
  v[`${CSS_PREFIX}brand-subtle`] = t.brand.subtle;
  v[`${CSS_PREFIX}accent`] = t.accent.base;
  v[`${CSS_PREFIX}on-accent`] = t.accent.on;
  v[`${CSS_PREFIX}focus`] = t.focus;
  (Object.keys(statusCss) as StatusKey[]).forEach(k => {
    v[`${CSS_PREFIX}status-${statusCss[k]}`] = t.status[k].fg;
    v[`${CSS_PREFIX}status-${statusCss[k]}-bg`] = t.status[k].bg;
  });
  v[`${CSS_PREFIX}nav-bg`] = t.nav.bg;
  v[`${CSS_PREFIX}nav-hover`] = t.nav.hover;
  v[`${CSS_PREFIX}nav-text`] = t.nav.text;
  v[`${CSS_PREFIX}nav-selected`] = t.nav.selected;
  v[`${CSS_PREFIX}shadow`] = t.shadow;
  return v;
}

/**
 * Backstage UI (`--bui-*`) variables mapped onto the same tokens, so BUI
 * components and MUI components match. Brand red is the BUI "accent"; critical
 * is the BUI "negative"/"danger" family (the deprecated `bg-danger` etc. are kept
 * only because `<Box bg="danger">` still reads them), so errors never use the brand red.
 */
export function buiVariables(mode: ThemeMode): Record<string, string> {
  const t = tokens[mode];
  const s = t.status;
  const v: Record<string, string> = {
    '--bui-ring': t.focus,
    '--bui-bg-app': t.surface[0],
    '--bui-bg-neutral-1': t.surface[1],
    '--bui-bg-neutral-2': t.surface[2],
    '--bui-bg-neutral-3': t.surface[3],
    '--bui-bg-neutral-4': t.surface[2],
    '--bui-fg-primary': t.text[1],
    '--bui-fg-secondary': t.text[2],
    '--bui-fg-disabled': t.text[3],
    '--bui-fg-positive': s.ok.fg,
    '--bui-fg-negative': s.critical.fg,
    '--bui-fg-warning': s.warn.fg,
    '--bui-fg-announcement': s.info.fg,
    '--bui-border-1': t.border.subtle,
    '--bui-border-2': t.border.strong,
    '--bui-accent-bg': t.brand.base,
    '--bui-accent-bg-hover': t.brand.hover,
    '--bui-accent-bg-disabled': t.surface[3],
    '--bui-accent-fg': t.brand.on,
    '--bui-accent-fg-disabled': t.text[3],
    // Deprecated in BUI 0.18 but still read by Button, Checkbox, Slider etc.
    '--bui-bg-solid': t.brand.base,
    '--bui-bg-solid-hover': t.brand.hover,
    '--bui-bg-solid-pressed': t.brand.press,
    '--bui-bg-solid-disabled': t.surface[3],
    '--bui-fg-solid': t.brand.on,
    '--bui-fg-solid-disabled': t.text[3],
    '--bui-bg-danger': s.critical.bg,
    '--bui-bg-info': s.info.bg,
    '--bui-fg-danger-on-bg': s.critical.fg,
    '--bui-fg-warning-on-bg': s.warn.fg,
    '--bui-fg-success-on-bg': s.ok.fg,
    '--bui-fg-info-on-bg': s.info.fg,
    '--bui-bg-warning': s.warn.bg,
    '--bui-bg-success': s.ok.bg,
  };
  (
    [
      ['negative', 'critical'],
      ['warning', 'warn'],
      ['positive', 'ok'],
      ['announcement', 'info'],
    ] as const
  ).forEach(([bui, key]) => {
    v[`--bui-${bui}-bg-subdued`] = s[key].bg;
    v[`--bui-${bui}-fg-subdued`] = s[key].fg;
    v[`--bui-${bui}-border`] = s[key].fg;
  });
  return v;
}

/** Hero image behind page headers; off unless the user opts in. */
export const heroImage = (mode: ThemeMode) =>
  `url(/brand/legacy-background-${mode}.webp)`;
