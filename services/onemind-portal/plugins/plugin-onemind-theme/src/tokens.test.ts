import { readFileSync } from 'node:fs';
import { join } from 'node:path';
import { contrastPairs, contrastRatio } from './contrast';
import { createOneMindTheme } from './muiTheme';
import { buiVariables, cssVariables, tokens, ThemeMode } from './tokens';

const modes: ThemeMode[] = ['dark', 'light'];
const css = readFileSync(join(__dirname, 'brand.css'), 'utf8');

const rgb = (hex: string) =>
  [1, 3, 5].map(i => parseInt(hex.slice(i, i + 2), 16)) as [
    number,
    number,
    number,
  ];
/** Spread between the strongest and weakest channel; 0 is pure grey. */
const chroma = (hex: string) => {
  const c = rgb(hex);
  return Math.max(...c) - Math.min(...c);
};
const blueTinted = (hex: string) => {
  const [r, , b] = rgb(hex);
  return b - r > 4;
};

const hue = (hex: string) => {
  const [r, g, b] = rgb(hex).map(x => x / 255);
  const max = Math.max(r, g, b);
  const d = max - Math.min(r, g, b);
  if (d === 0) return 0;
  let h = (r - g) / d + 4;
  if (max === r) h = ((g - b) / d) % 6;
  else if (max === g) h = (b - r) / d + 2;
  return (h * 60 + 360) % 360;
};

describe('token values', () => {
  it('pins the brand red sampled from the logo', () => {
    expect(tokens.dark.brand.base).toBe('#ff333a');
    expect(tokens.light.brand.base).toBe('#d4151c');
  });

  it('keeps surfaces neutral graphite with no blue tint', () => {
    modes.forEach(mode => {
      Object.values(tokens[mode].surface).forEach(hex => {
        expect(chroma(hex)).toBeLessThanOrEqual(2);
        expect(blueTinted(hex)).toBe(false);
      });
      expect(blueTinted(tokens[mode].nav.bg)).toBe(false);
    });
  });

  it('matches the token snapshot', () => {
    expect(tokens).toMatchSnapshot();
  });

  it('keeps critical distinct from the brand red', () => {
    modes.forEach(mode => {
      const { brand, status } = tokens[mode];
      expect(status.critical.fg).not.toBe(brand.base);
      expect(status.critical.fg).not.toBe(brand.text);
      // A clearly different hue (degrees on the colour wheel), not just a shade.
      const gap = Math.abs(hue(status.critical.fg) - hue(brand.text));
      expect(Math.min(gap, 360 - gap)).toBeGreaterThanOrEqual(20);
    });
  });
});

describe('contrast (WCAG 2.x)', () => {
  modes.forEach(mode => {
    it(`every ${mode} pair meets its minimum`, () => {
      const failures = contrastPairs(mode)
        .map(p => ({ ...p, ratio: contrastRatio(p.fg, p.bg) }))
        .filter(p => p.ratio < p.min)
        .map(p => `${p.name}: ${p.ratio.toFixed(2)} < ${p.min}`);
      expect(failures).toEqual([]);
    });
  });

  it('computes known ratios', () => {
    expect(contrastRatio('#000000', '#ffffff')).toBeCloseTo(21, 5);
    expect(contrastRatio('#ffffff', '#ffffff')).toBeCloseTo(1, 5);
  });
});

describe('brand.css stays in sync with tokens.ts', () => {
  const declared = (name: string) => {
    const m = [...css.matchAll(new RegExp(`${name}:\\s*([^;]+);`, 'g'))];
    return m.map(x => x[1].trim());
  };

  modes.forEach(mode => {
    it(`declares every ${mode} variable`, () => {
      const all = { ...cssVariables(mode), ...buiVariables(mode) };
      const section =
        mode === 'dark'
          ? css.slice(
              css.indexOf("body[data-theme-mode='dark']"),
              css.indexOf('/* tokens:end */'),
            )
          : css.slice(
              css.indexOf(':root'),
              css.indexOf("body[data-theme-mode='dark']"),
            );
      const missing = Object.entries(all).filter(
        ([k, v]) => !new RegExp(`${k}:\\s*${v};`).test(section),
      );
      expect(missing).toEqual([]);
    });
  });

  it('defines each --onemind variable in both modes', () => {
    Object.keys(cssVariables('dark')).forEach(k => {
      expect(declared(k).length).toBeGreaterThanOrEqual(2);
    });
  });
});

describe('base background in brand.css', () => {
  it('paints html, body and #root with the surface token and fills the viewport', () => {
    expect(css).toMatch(
      /html,\s*body,\s*#root\s*\{[^}]*background-color:\s*var\(--bui-bg-app\);[^}]*min-height:\s*100%/,
    );
    expect(css).toMatch(/body\s*\{[^}]*min-height:\s*100vh/);
    expect(css).toMatch(
      /html:has\(body\[data-theme-mode='dark'\]\)\s*\{[^}]*var\(--onemind-surface-0\)/,
    );
    // The variable it reads is surface 0 in both modes.
    modes.forEach(mode =>
      expect(buiVariables(mode)['--bui-bg-app']).toBe(tokens[mode].surface[0]),
    );
  });
});

describe('MUI theme', () => {
  modes.forEach(mode => {
    it(`${mode} theme uses the tokens`, () => {
      const t = tokens[mode];
      const theme = createOneMindTheme(mode).getTheme('v5') as any;
      expect(theme.palette.mode).toBe(mode);
      expect(theme.palette.primary.main).toBe(t.brand.base);
      expect(theme.palette.background.default).toBe(t.surface[0]);
      expect(theme.palette.background.paper).toBe(t.surface[1]);
      expect(theme.palette.error.main).toBe(t.status.critical.fg);
      expect(theme.palette.navigation.indicator).toBe(t.brand.base);
    });
  });

  modes.forEach(mode => {
    it(`${mode} theme paints html, body and #root with surface 0`, () => {
      const base = (createOneMindTheme(mode).getTheme('v5') as any).components
        .MuiCssBaseline.styleOverrides['html, body, #root'];
      expect(base.backgroundColor).toBe(tokens[mode].surface[0]);
      expect(base.minHeight).toBe('100%');
    });
  });

  it('shows no header image by default and adds it when enabled', () => {
    const header = (opts?: { backgroundImage?: boolean }) =>
      (createOneMindTheme('dark', opts).getTheme('v5') as any).components
        .BackstageHeader.styleOverrides.header.backgroundImage as string;
    expect(header()).toBe('none');
    expect(header({ backgroundImage: true })).toContain(
      'legacy-background-dark.webp',
    );
  });
});
