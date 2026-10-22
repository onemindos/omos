import { tokens, ThemeMode, StatusKey } from './tokens';

const channel = (hex: string, i: number) => {
  const c = parseInt(hex.slice(1 + i * 2, 3 + i * 2), 16) / 255;
  return c <= 0.03928 ? c / 12.92 : ((c + 0.055) / 1.055) ** 2.4;
};

/** WCAG 2.x relative luminance of `#rrggbb`. */
export function luminance(hex: string): number {
  return (
    0.2126 * channel(hex, 0) +
    0.7152 * channel(hex, 1) +
    0.0722 * channel(hex, 2)
  );
}

/** WCAG 2.x contrast ratio between two `#rrggbb` colours (1 to 21). */
export function contrastRatio(a: string, b: string): number {
  const [hi, lo] = [luminance(a), luminance(b)].sort((x, y) => y - x);
  return (hi + 0.05) / (lo + 0.05);
}

export interface ContrastPair {
  name: string;
  fg: string;
  bg: string;
  /** 4.5 for text, 3 for focus rings and graphical objects. */
  min: number;
}

/** Every foreground/background pairing the theme promises is accessible. */
export function contrastPairs(mode: ThemeMode): ContrastPair[] {
  const t = tokens[mode];
  const pairs: ContrastPair[] = [];
  ([0, 1, 2] as const).forEach(i => {
    const bg = t.surface[i];
    const on = `surface ${i}`;
    pairs.push(
      { name: `text 1 on ${on}`, fg: t.text[1], bg, min: 4.5 },
      { name: `text 2 on ${on}`, fg: t.text[2], bg, min: 4.5 },
      { name: `text 3 on ${on}`, fg: t.text[3], bg, min: 4.5 },
      { name: `link (brand text) on ${on}`, fg: t.brand.text, bg, min: 4.5 },
      { name: `focus ring on ${on}`, fg: t.focus, bg, min: 3 },
      { name: `brand fill on ${on}`, fg: t.brand.base, bg, min: 3 },
    );
    (Object.keys(t.status) as StatusKey[]).forEach(k =>
      pairs.push({
        name: `status ${k} on ${on}`,
        fg: t.status[k].fg,
        bg,
        min: 4.5,
      }),
    );
  });
  pairs.push(
    { name: 'text 1 on surface 3', fg: t.text[1], bg: t.surface[3], min: 4.5 },
    { name: 'text 2 on surface 3', fg: t.text[2], bg: t.surface[3], min: 4.5 },
    { name: 'focus ring on surface 3', fg: t.focus, bg: t.surface[3], min: 3 },
    {
      name: 'text on brand button',
      fg: t.brand.on,
      bg: t.brand.base,
      min: 4.5,
    },
    {
      name: 'text on brand button (hover)',
      fg: t.brand.on,
      bg: t.brand.hover,
      min: 4.5,
    },
    {
      name: 'text on brand button (press)',
      fg: t.brand.on,
      bg: t.brand.press,
      min: 4.5,
    },
    { name: 'text on accent', fg: t.accent.on, bg: t.accent.base, min: 4.5 },
    {
      name: 'accent on surface 1',
      fg: t.accent.base,
      bg: t.surface[1],
      min: 4.5,
    },
    { name: 'nav text on nav', fg: t.nav.text, bg: t.nav.bg, min: 4.5 },
    { name: 'nav selected on nav', fg: t.nav.selected, bg: t.nav.bg, min: 4.5 },
    {
      name: 'nav text on nav hover',
      fg: t.nav.text,
      bg: t.nav.hover,
      min: 4.5,
    },
    { name: 'brand indicator on nav', fg: t.brand.base, bg: t.nav.bg, min: 3 },
    {
      name: 'text 1 on brand subtle',
      fg: t.text[1],
      bg: t.brand.subtle,
      min: 4.5,
    },
    {
      name: 'link on brand subtle',
      fg: t.brand.text,
      bg: t.brand.subtle,
      min: 4.5,
    },
  );
  (Object.keys(t.status) as StatusKey[]).forEach(k =>
    pairs.push({
      name: `status ${k} chip (fg on tint)`,
      fg: t.status[k].fg,
      bg: t.status[k].bg,
      min: 4.5,
    }),
  );
  return pairs;
}
