import type { BusEvent, Track } from './contract';
import { AFFILIATIONS, DOMAINS, affiliationOf, domainOf } from './cot';
import type { Affiliation, Domain } from './cot';

export interface TrackFilter {
  affiliations: ReadonlySet<Affiliation>;
  domains: ReadonlySet<Domain>;
}

export const ALL_TRACKS: TrackFilter = {
  affiliations: new Set(AFFILIATIONS),
  domains: new Set(DOMAINS),
};

export function matchesTrack(track: Track, filter: TrackFilter): boolean {
  return (
    filter.affiliations.has(affiliationOf(track.type)) &&
    filter.domains.has(domainOf(track.type))
  );
}

export function filterTracks(
  tracks: Iterable<Track>,
  filter: TrackFilter,
): Track[] {
  const out: Track[] = [];
  for (const t of tracks) if (matchesTrack(t, filter)) out.push(t);
  return out;
}

export function toggle<T>(set: ReadonlySet<T>, value: T): Set<T> {
  const next = new Set(set);
  if (next.has(value)) next.delete(value);
  else next.add(value);
  return next;
}

export interface TrackCounts {
  total: number;
  affiliation: Record<Affiliation, number>;
  domain: Record<Domain, number>;
}

/** Counts every track, ignoring the filter, so the legend shows what a toggle would reveal. */
export function countTracks(tracks: Iterable<Track>): TrackCounts {
  const affiliation = Object.fromEntries(
    AFFILIATIONS.map(a => [a, 0]),
  ) as Record<Affiliation, number>;
  const domain = Object.fromEntries(DOMAINS.map(d => [d, 0])) as Record<
    Domain,
    number
  >;
  let total = 0;
  for (const t of tracks) {
    total += 1;
    affiliation[affiliationOf(t.type)] += 1;
    domain[domainOf(t.type)] += 1;
  }
  return { total, affiliation, domain };
}

/** Case-insensitive substring match on callsign or uid. Exact and prefix matches come first. */
export function searchTracks(
  tracks: Iterable<Track>,
  query: string,
  limit = 8,
): Track[] {
  const q = query.trim().toLowerCase();
  if (!q) return [];
  const scored: Array<{ t: Track; score: number }> = [];
  for (const t of tracks) {
    const fields = [t.callsign?.toLowerCase(), t.uid.toLowerCase()];
    let score = -1;
    for (const f of fields) {
      if (!f) continue;
      if (f === q) score = Math.max(score, 3);
      else if (f.startsWith(q)) score = Math.max(score, 2);
      else if (f.includes(q)) score = Math.max(score, 1);
    }
    if (score > 0) scored.push({ t, score });
  }
  scored.sort(
    (a, b) =>
      b.score - a.score ||
      (a.t.callsign ?? a.t.uid).localeCompare(b.t.callsign ?? b.t.uid),
  );
  return scored.slice(0, limit).map(s => s.t);
}

export interface EventFilter {
  /** Kinds that are hidden. Empty hides nothing. */
  hiddenKinds: ReadonlySet<string>;
  text: string;
}

export function matchesEvent(event: BusEvent, filter: EventFilter): boolean {
  if (filter.hiddenKinds.has(event.kind)) return false;
  const q = filter.text.trim().toLowerCase();
  if (!q) return true;
  return (
    event.subject.toLowerCase().includes(q) ||
    event.summary.toLowerCase().includes(q)
  );
}
