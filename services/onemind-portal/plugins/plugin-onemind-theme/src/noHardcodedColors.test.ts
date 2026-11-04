import { readdirSync, readFileSync, statSync } from 'node:fs';
import { join, relative, sep } from 'node:path';

const pluginsDir = join(__dirname, '..', '..');
const portalDir = join(pluginsDir, '..');

/**
 * Plugins (and the app) must take colours from the theme, not hard-code them.
 * Exempt:
 * - plugin-onemind-theme: it defines the tokens.
 * - plugin-onemind-nav and plugin-onemind-home: being rewritten by F1, which
 *   adopts the tokens. REMOVE these two entries when F1 merges.
 * - plugin-onemind-offline-backend: generates a standalone HTML file that must
 *   work without the portal's CSS, and a backend cannot import the frontend
 *   theme.
 */
const EXEMPT = [
  'plugin-onemind-theme',
  'plugin-onemind-nav',
  'plugin-onemind-home',
  'plugin-onemind-offline-backend',
];

const HEX = /#[0-9a-fA-F]{8}\b|#[0-9a-fA-F]{6}\b|#[0-9a-fA-F]{3,4}\b/;
const FUNC = /\brgba?\(|\bhsla?\(/;
const SOURCE = /\.(ts|tsx|css|scss)$/;
const TEST = /\.test\.(ts|tsx)$/;

function* walk(dir: string): Generator<string> {
  for (const name of readdirSync(dir)) {
    if (name === 'node_modules' || name === 'dist' || name === 'dist-types')
      continue;
    const full = join(dir, name);
    if (statSync(full).isDirectory()) yield* walk(full);
    else if (SOURCE.test(name) && !TEST.test(name)) yield full;
  }
}

function offenders(root: string, skip: (rel: string) => boolean) {
  const found: string[] = [];
  for (const file of walk(root)) {
    const rel = relative(portalDir, file);
    if (skip(rel)) continue;
    readFileSync(file, 'utf8')
      .split('\n')
      .forEach((line, i) => {
        // `#` followed by hex inside a URL fragment or an issue reference is not a colour.
        const code = line.replace(/https?:\/\/\S+/g, '');
        if (HEX.test(code) || FUNC.test(code))
          found.push(`${rel}:${i + 1}: ${line.trim().slice(0, 80)}`);
      });
  }
  return found;
}

describe('no hard-coded colours outside the theme package', () => {
  it('plugins use tokens, not hex/rgb literals', () => {
    const skip = (rel: string) =>
      EXEMPT.some(e => rel.split(sep).slice(0, 2).join('/') === `plugins/${e}`);
    expect(offenders(pluginsDir, skip)).toEqual([]);
  });

  it('the app package uses tokens, not hex/rgb literals', () => {
    expect(
      offenders(join(portalDir, 'packages', 'app', 'src'), () => false),
    ).toEqual([]);
  });
});
