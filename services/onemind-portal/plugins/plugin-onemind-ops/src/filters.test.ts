import type { BusEvent, Track } from './contract';
import {
  ALL_TRACKS,
  countTracks,
  filterTracks,
  matchesEvent,
  matchesTrack,
  searchTracks,
  toggle,
} from './filters';

const t = (uid: string, type: string, callsign?: string): Track => ({
  uid,
  type,
  callsign,
  lat: 0,
  lon: 0,
  time: '2026-11-07T10:00:00Z',
  stale: '2026-11-07T10:05:00Z',
  subject: 'tak.cot.x',
});
const tracks = [
  t('ANDROID-1', 'a-f-G-U-C', 'ALPHA'),
  t('ae1234', 'a-n-A-C-F', 'UAL123'),
  t('ship-9', 'a-u-S-X-M', 'MAERSK ONE'),
  t('hz-1', 'a-h-G'),
  t('odd', 'a-s-G', 'SUSPECT'),
];

describe('track filters', () => {
  it('shows everything by default', () => {
    expect(filterTracks(tracks, ALL_TRACKS)).toHaveLength(5);
  });

  it('hides an affiliation', () => {
    const f = {
      ...ALL_TRACKS,
      affiliations: toggle(ALL_TRACKS.affiliations, 'hostile' as const),
    };
    expect(filterTracks(tracks, f).map(x => x.uid)).not.toContain('hz-1');
    expect(filterTracks(tracks, f)).toHaveLength(4);
  });

  it('combines affiliation and domain', () => {
    const f = {
      affiliations: new Set(['friendly', 'neutral'] as const),
      domains: new Set(['air'] as const),
    };
    expect(filterTracks(tracks, f).map(x => x.uid)).toEqual(['ae1234']);
  });

  it('puts the "other" affiliation behind its own toggle', () => {
    const f = {
      ...ALL_TRACKS,
      affiliations: toggle(ALL_TRACKS.affiliations, 'other' as const),
    };
    expect(matchesTrack(tracks[4], f)).toBe(false);
  });

  it('toggle does not mutate its input', () => {
    const before = new Set(ALL_TRACKS.affiliations);
    toggle(ALL_TRACKS.affiliations, 'friendly');
    expect(ALL_TRACKS.affiliations).toEqual(before);
  });

  it('counts per category regardless of the filter', () => {
    const c = countTracks(tracks);
    expect(c.total).toBe(5);
    expect(c.affiliation).toEqual({
      friendly: 1,
      hostile: 1,
      neutral: 1,
      unknown: 1,
      other: 1,
    });
    expect(c.domain.ground).toBe(3);
    expect(c.domain.air).toBe(1);
    expect(c.domain.sea).toBe(1);
  });
});

describe('searchTracks', () => {
  it('matches callsign or uid, case-insensitively, best match first', () => {
    expect(searchTracks(tracks, 'ual').map(x => x.uid)).toEqual(['ae1234']);
    expect(searchTracks(tracks, 'ANDROID').map(x => x.uid)).toEqual([
      'ANDROID-1',
    ]);
    expect(searchTracks(tracks, 'one').map(x => x.uid)).toEqual(['ship-9']);
  });

  it('returns nothing for a blank or unmatched query', () => {
    expect(searchTracks(tracks, '  ')).toEqual([]);
    expect(searchTracks(tracks, 'zzz')).toEqual([]);
  });

  it('puts an exact match before a prefix match and respects the limit', () => {
    const many = [
      t('a1', 'a-f-G', 'bravo-2'),
      t('a2', 'a-f-G', 'bravo'),
      t('a3', 'a-f-G', 'bravo-3'),
    ];
    expect(searchTracks(many, 'bravo')[0].uid).toBe('a2');
    expect(searchTracks(many, 'bravo', 2)).toHaveLength(2);
  });
});

describe('event filters', () => {
  const e = (kind: string, subject: string, summary = ''): BusEvent => ({
    id: subject,
    kind,
    subject,
    summary,
    ts: 'x',
    bytes: 1,
  });
  it('hides kinds and filters by text on subject or summary', () => {
    const none = new Set<string>();
    expect(
      matchesEvent(e('tel', 'tel.a.b.c'), { hiddenKinds: none, text: '' }),
    ).toBe(true);
    expect(
      matchesEvent(e('tel', 'tel.a.b.c'), {
        hiddenKinds: new Set(['tel']),
        text: '',
      }),
    ).toBe(false);
    expect(
      matchesEvent(e('tak', 'tak.cot.x', 'UAL123 a-n-A'), {
        hiddenKinds: none,
        text: 'ual',
      }),
    ).toBe(true);
    expect(
      matchesEvent(e('tak', 'tak.cot.x', 'UAL123'), {
        hiddenKinds: none,
        text: 'det.',
      }),
    ).toBe(false);
  });
});
