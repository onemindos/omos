import { tokens } from '@internal/plugin-onemind-theme';
import type { Track } from './contract';

/**
 * Pure helpers for Cursor-on-Target types. A CoT atom type reads `a-<affiliation>-<dimension>-<function...>`,
 * e.g. `a-f-G-U-C` is a friendly ground unit (combat), `a-n-A-C-F` a neutral civil fixed-wing aircraft.
 */

export type Affiliation =
  | 'friendly'
  | 'hostile'
  | 'neutral'
  | 'unknown'
  | 'other';
export type Domain =
  | 'air'
  | 'ground'
  | 'sea'
  | 'subsurface'
  | 'space'
  | 'other';

export const AFFILIATIONS: Affiliation[] = [
  'friendly',
  'hostile',
  'neutral',
  'unknown',
  'other',
];
export const DOMAINS: Domain[] = [
  'air',
  'ground',
  'sea',
  'subsurface',
  'space',
  'other',
];

export const AFFILIATION_LABEL: Record<Affiliation, string> = {
  friendly: 'Friendly',
  hostile: 'Hostile',
  neutral: 'Neutral',
  unknown: 'Unknown',
  other: 'Other',
};

export const DOMAIN_LABEL: Record<Domain, string> = {
  air: 'Air',
  ground: 'Ground',
  sea: 'Sea surface',
  subsurface: 'Subsurface',
  space: 'Space',
  other: 'Other',
};

/** Second atom of the type. Anything that is not f, h, n or u (assumed friend, suspect, pending...) is "other". */
export function affiliationOf(type: string): Affiliation {
  const [root, second] = type.split('-');
  if (root?.toLowerCase() !== 'a') return 'other';
  const atom = second?.toLowerCase();
  if (atom === 'f') return 'friendly';
  if (atom === 'h') return 'hostile';
  if (atom === 'n') return 'neutral';
  if (atom === 'u') return 'unknown';
  return 'other';
}

/** Third atom of the type. */
export function domainOf(type: string): Domain {
  const [root, , third] = type.split('-');
  if (root?.toLowerCase() !== 'a') return 'other';
  const atom = third?.toUpperCase();
  if (atom === 'A') return 'air';
  if (atom === 'G' || atom === 'F') return 'ground';
  if (atom === 'S') return 'sea';
  if (atom === 'U') return 'subsurface';
  if (atom === 'P') return 'space';
  return 'other';
}

const AFFILIATION_WORD: Record<string, string> = {
  f: 'friendly',
  h: 'hostile',
  n: 'neutral',
  u: 'unknown',
  a: 'assumed friendly',
  s: 'suspect',
  j: 'joker',
  k: 'faker',
  p: 'pending',
  o: 'unaffiliated',
  g: 'exercise',
  w: 'exercise',
  d: 'exercise',
  l: 'exercise',
  m: 'exercise',
};

const DOMAIN_WORD: Record<string, string> = {
  A: 'aircraft',
  G: 'ground unit',
  F: 'special forces unit',
  S: 'surface vessel',
  U: 'subsurface vessel',
  P: 'space object',
  X: 'other object',
};

/** Function atoms we can name. Keyed by `<domain>-<rest of the type>`; the longest matching prefix wins. */
const FUNCTION_WORD: Record<string, string> = {
  'A-C-F': 'civil fixed wing',
  'A-C-H': 'civil helicopter',
  'A-C': 'civil',
  'A-M-F': 'military fixed wing',
  'A-M-H': 'military helicopter',
  'A-M': 'military',
  'A-W-M': 'weapon (missile)',
  'G-U-C': 'combat',
  'G-U-I': 'infantry',
  'G-U-S': 'combat support',
  'G-U-U': 'combat service support',
  'G-E-V': 'vehicle',
  'G-E-S': 'equipment',
  'G-I': 'installation',
  'S-X-M': 'merchant',
  'S-C': 'civilian vessel',
  'S-F': 'fishing vessel',
  'S-X': 'vessel',
};

function functionWord(type: string): string | undefined {
  const atoms = type
    .split('-')
    .slice(2)
    .map(a => a.toUpperCase());
  for (let n = atoms.length; n >= 2; n -= 1) {
    const hit = FUNCTION_WORD[atoms.slice(0, n).join('-')];
    if (hit) return hit;
  }
  return undefined;
}

/**
 * Plain words for a CoT type: "friendly ground unit (combat)", "neutral aircraft (civil fixed wing)".
 * Types outside the `a-` atom family (bits, tasking, alerts) are described by their root.
 */
export function decodeCotType(type: string): string {
  const atoms = type.split('-');
  if (atoms[0]?.toLowerCase() !== 'a') {
    const root: Record<string, string> = {
      b: 'sensor or detail event',
      t: 'tasking request',
      r: 'reply',
      c: 'capability',
      u: 'user-defined event',
    };
    return root[atoms[0]?.toLowerCase() ?? ''] ?? `CoT ${type}`;
  }
  const affiliation =
    AFFILIATION_WORD[atoms[1]?.toLowerCase() ?? ''] ?? 'unspecified';
  const noun = DOMAIN_WORD[atoms[2]?.toUpperCase() ?? ''] ?? 'object';
  const detail = functionWord(type);
  return `${affiliation} ${noun}${detail ? ` (${detail})` : ''}`;
}

/**
 * Affiliation colours come from the dark theme tokens (the globe is always dark):
 * friendly = info blue, hostile = brand red, neutral = ok green, unknown = warn yellow, other = grey.
 */
export const AFFILIATION_COLOR: Record<Affiliation, string> = {
  friendly: tokens.dark.status.info.fg,
  hostile: tokens.dark.brand.base,
  neutral: tokens.dark.status.ok.fg,
  unknown: tokens.dark.status.warn.fg,
  other: tokens.dark.status.unknown.fg,
};

/** Marker size in pixels. Aircraft are the biggest because they move fastest and fly over everything. */
export const DOMAIN_SIZE: Record<Domain, number> = {
  air: 12,
  sea: 9,
  ground: 8,
  subsurface: 7,
  space: 7,
  other: 6,
};

export const MS_TO_KNOTS = 1.943844;
export const M_TO_FEET = 3.28084;

export const toKnots = (metresPerSecond: number) =>
  metresPerSecond * MS_TO_KNOTS;
export const toFeet = (metres: number) => metres * M_TO_FEET;

/** True for a usable course in degrees. */
export const hasCourse = (t: Pick<Track, 'course'>): t is { course: number } =>
  typeof t.course === 'number' &&
  Number.isFinite(t.course) &&
  t.course >= 0 &&
  t.course <= 360;

/**
 * Opacity of a track by age. Full for the first half of its life (time -> stale), fading to 0.25 at stale,
 * and 0.15 once past it (the backend removes it soon after, and a removal arrives as a delta).
 */
export function fadeAlpha(
  track: Pick<Track, 'time' | 'stale'>,
  nowMs: number,
): number {
  const start = Date.parse(track.time);
  const end = Date.parse(track.stale);
  if (!Number.isFinite(start) || !Number.isFinite(end) || end <= start)
    return 1;
  if (nowMs >= end) return 0.15;
  const fraction = (nowMs - start) / (end - start);
  if (fraction <= 0.5) return 1;
  return 1 - ((fraction - 0.5) / 0.5) * 0.75;
}

/** "5 s ago", "3 min ago", "2 h ago". */
export function ago(fromMs: number | undefined, nowMs: number): string {
  if (fromMs === undefined || !Number.isFinite(fromMs)) return 'never';
  const s = Math.max(0, Math.round((nowMs - fromMs) / 1000));
  if (s < 60) return `${s} s ago`;
  if (s < 3600) return `${Math.floor(s / 60)} min ago`;
  if (s < 86400) return `${Math.floor(s / 3600)} h ago`;
  return `${Math.floor(s / 86400)} d ago`;
}
