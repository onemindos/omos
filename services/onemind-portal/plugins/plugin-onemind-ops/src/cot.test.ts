import { tokens } from '@internal/plugin-onemind-theme';
import {
  AFFILIATION_COLOR,
  DOMAIN_SIZE,
  affiliationOf,
  ago,
  decodeCotType,
  domainOf,
  fadeAlpha,
  hasCourse,
  toFeet,
  toKnots,
} from './cot';

describe('affiliation and domain', () => {
  it.each([
    ['a-f-G-U-C', 'friendly', 'ground'],
    ['a-h-A-M-F', 'hostile', 'air'],
    ['a-n-A-C-F', 'neutral', 'air'],
    ['a-u-S-X', 'unknown', 'sea'],
    ['a-s-G', 'other', 'ground'],
    ['a-a-U', 'other', 'subsurface'],
    ['a-p-P', 'other', 'space'],
    ['a-f-X', 'friendly', 'other'],
    ['b-m-p-s-p-loc', 'other', 'other'],
  ])('%s -> %s %s', (type, affiliation, domain) => {
    expect(affiliationOf(type)).toBe(affiliation);
    expect(domainOf(type)).toBe(domain);
  });
});

describe('decodeCotType', () => {
  it.each([
    ['a-f-G-U-C', 'friendly ground unit (combat)'],
    ['a-n-A-C-F', 'neutral aircraft (civil fixed wing)'],
    ['a-h-A-M-H', 'hostile aircraft (military helicopter)'],
    ['a-u-S-X-M', 'unknown surface vessel (merchant)'],
    ['a-f-G', 'friendly ground unit'],
    ['a-s-G-E-V', 'suspect ground unit (vehicle)'],
    ['b-m-p-s-p-loc', 'sensor or detail event'],
    ['t-x-c-t', 'tasking request'],
    ['z-weird', 'CoT z-weird'],
  ])('%s -> %s', (type, words) => {
    expect(decodeCotType(type)).toBe(words);
  });
});

describe('colours and sizes', () => {
  it('uses the dark theme tokens for the five affiliations', () => {
    expect(AFFILIATION_COLOR.friendly).toBe(tokens.dark.status.info.fg);
    expect(AFFILIATION_COLOR.hostile).toBe(tokens.dark.brand.base);
    expect(AFFILIATION_COLOR.neutral).toBe(tokens.dark.status.ok.fg);
    expect(AFFILIATION_COLOR.unknown).toBe(tokens.dark.status.warn.fg);
    expect(AFFILIATION_COLOR.other).toBe(tokens.dark.status.unknown.fg);
    expect(new Set(Object.values(AFFILIATION_COLOR)).size).toBe(5);
  });

  it('draws aircraft biggest', () => {
    expect(DOMAIN_SIZE.air).toBeGreaterThan(DOMAIN_SIZE.ground);
    expect(DOMAIN_SIZE.air).toBeGreaterThan(DOMAIN_SIZE.sea);
    expect(DOMAIN_SIZE.air).toBeGreaterThan(DOMAIN_SIZE.other);
  });
});

describe('units and course', () => {
  it('converts m/s to knots and metres to feet', () => {
    expect(Math.round(toKnots(10))).toBe(19);
    expect(Math.round(toFeet(1000))).toBe(3281);
  });

  it('knows when a course is usable', () => {
    expect(hasCourse({ course: 0 })).toBe(true);
    expect(hasCourse({ course: 359.9 })).toBe(true);
    expect(hasCourse({ course: undefined })).toBe(false);
    expect(hasCourse({ course: NaN })).toBe(false);
    expect(hasCourse({ course: -1 })).toBe(false);
  });
});

describe('fadeAlpha', () => {
  const t = { time: '2026-11-07T10:00:00Z', stale: '2026-11-07T10:10:00Z' };
  const at = (s: string) => Date.parse(s);
  it('is opaque for the first half, fades to 0.25 at stale and is dim after', () => {
    expect(fadeAlpha(t, at('2026-11-07T10:02:00Z'))).toBe(1);
    expect(fadeAlpha(t, at('2026-11-07T10:05:00Z'))).toBe(1);
    expect(fadeAlpha(t, at('2026-11-07T10:07:30Z'))).toBeCloseTo(0.625, 3);
    expect(fadeAlpha(t, at('2026-11-07T10:09:59Z'))).toBeLessThan(0.3);
    expect(fadeAlpha(t, at('2026-11-07T10:11:00Z'))).toBe(0.15);
  });
  it('stays opaque when the times are unusable', () => {
    expect(fadeAlpha({ time: 'x', stale: 'y' }, 0)).toBe(1);
  });
});

describe('ago', () => {
  it('formats seconds, minutes, hours and days', () => {
    expect(ago(1000, 3000)).toBe('2 s ago');
    expect(ago(0, 125_000)).toBe('2 min ago');
    expect(ago(0, 7_300_000)).toBe('2 h ago');
    expect(ago(0, 200_000_000)).toBe('2 d ago');
    expect(ago(undefined, 0)).toBe('never');
  });
});
