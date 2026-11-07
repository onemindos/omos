import { Component, Suspense, lazy, useMemo, useState } from 'react';
import type { ReactNode } from 'react';
import { makeStyles } from '@material-ui/core/styles';
import { Progress } from '@backstage/core-components';
import type { SourceStatus, Track } from './contract';
import {
  AFFILIATIONS,
  AFFILIATION_COLOR,
  AFFILIATION_LABEL,
  DOMAINS,
  DOMAIN_LABEL,
  ago,
  decodeCotType,
  toFeet,
  toKnots,
} from './cot';
import type { Affiliation, Domain } from './cot';
import {
  ALL_TRACKS,
  countTracks,
  filterTracks,
  searchTracks,
  toggle,
} from './filters';
import type { TrackFilter } from './filters';
import { Note, Panel, cssVar, statusFg } from './ui';

// The only import of the Cesium module graph. Other portal pages never download this chunk.
const CesiumGlobe = lazy(() => import('./CesiumGlobe'));

const useStyles = makeStyles({
  stage: {
    position: 'relative',
    flex: 1,
    minHeight: 0,
    background: cssVar('surface-0'),
  },
  overlay: {
    position: 'absolute',
    zIndex: 2,
    background: `color-mix(in srgb, ${cssVar('surface-1')} 88%, transparent)`,
    border: `1px solid ${cssVar('border')}`,
    borderRadius: 6,
    fontSize: 12,
    color: cssVar('text-1'),
  },
  toolbar: {
    top: 8,
    left: 8,
    padding: 6,
    maxWidth: 'calc(100% - 340px)',
    display: 'flex',
    flexDirection: 'column',
    gap: 6,
  },
  row: { display: 'flex', flexWrap: 'wrap', gap: 4, alignItems: 'center' },
  toggle: {
    display: 'flex',
    alignItems: 'center',
    gap: 5,
    background: 'transparent',
    border: `1px solid ${cssVar('border')}`,
    color: cssVar('text-1'),
    borderRadius: 4,
    padding: '2px 7px',
    fontSize: 11,
    cursor: 'pointer',
    '&[aria-pressed=false]': { opacity: 0.4 },
    '&:focus-visible': { outline: `2px solid ${cssVar('focus')}` },
  },
  swatch: { width: 9, height: 9, borderRadius: '50%', display: 'inline-block' },
  count: { fontVariantNumeric: 'tabular-nums', color: cssVar('text-2') },
  search: {
    width: 200,
    background: cssVar('surface-0'),
    color: cssVar('text-1'),
    border: `1px solid ${cssVar('border')}`,
    borderRadius: 4,
    padding: '3px 6px',
    fontSize: 12,
  },
  results: {
    listStyle: 'none',
    margin: 0,
    padding: 0,
    maxHeight: 160,
    overflow: 'auto',
  },
  result: {
    display: 'block',
    width: '100%',
    textAlign: 'left',
    background: 'transparent',
    border: 0,
    color: cssVar('text-1'),
    padding: '3px 6px',
    cursor: 'pointer',
    '&:hover, &:focus-visible': { background: cssVar('surface-3') },
  },
  detail: {
    top: 8,
    right: 8,
    width: 320,
    maxHeight: 'calc(100% - 16px)',
    overflow: 'auto',
    padding: 10,
  },
  dl: {
    display: 'grid',
    gridTemplateColumns: '90px 1fr',
    gap: '2px 8px',
    margin: '6px 0 0',
  },
  dt: { color: cssVar('text-2') },
  dd: { margin: 0, wordBreak: 'break-word' },
  empty: {
    position: 'absolute',
    inset: 0,
    zIndex: 1,
    display: 'flex',
    alignItems: 'center',
    justifyContent: 'center',
    pointerEvents: 'none',
    textAlign: 'center',
    padding: 24,
  },
  emptyBox: {
    maxWidth: 380,
    padding: 14,
    background: cssVar('surface-1'),
    border: `1px solid ${cssVar('border')}`,
    borderRadius: 6,
    pointerEvents: 'auto',
  },
  attribution: {
    bottom: 24,
    left: 8,
    padding: '1px 6px',
    fontSize: 10,
    color: cssVar('text-2'),
  },
});

class GlobeBoundary extends Component<
  { children: ReactNode; fallback: (message: string) => ReactNode },
  { error?: string }
> {
  static getDerivedStateFromError(e: unknown) {
    return { error: e instanceof Error ? e.message : String(e) };
  }

  state: { error?: string } = {};

  render() {
    return this.state.error === undefined
      ? this.props.children
      : this.props.fallback(this.state.error);
  }
}

export function TrackDetail(props: {
  track?: Track;
  selectedUid: string;
  now: number;
  onClose: () => void;
}) {
  const c = useStyles();
  const { track, now } = props;
  if (!track) {
    return (
      <div
        className={`${c.overlay} ${c.detail}`}
        role="complementary"
        aria-label="Track detail"
      >
        <strong>{props.selectedUid}</strong>
        <Note tone="warn">
          This track is no longer reported (stale or removed).
        </Note>
        <button type="button" className={c.toggle} onClick={props.onClose}>
          Close
        </button>
      </div>
    );
  }
  const seen = Date.parse(track.time);
  const rows: Array<[string, string]> = [
    ['Type', `${decodeCotType(track.type)} (${track.type})`],
    [
      'Speed',
      track.speed === undefined
        ? 'unknown'
        : `${Math.round(toKnots(track.speed))} kt`,
    ],
    [
      'Altitude',
      track.hae === undefined
        ? 'unknown'
        : `${Math.round(toFeet(track.hae)).toLocaleString()} ft HAE`,
    ],
    [
      'Course',
      track.course === undefined ? 'unknown' : `${Math.round(track.course)}°`,
    ],
    ['Position', `${track.lat.toFixed(4)}, ${track.lon.toFixed(4)}`],
    ['Last seen', Number.isFinite(seen) ? ago(seen, now) : 'unknown'],
    ['Subject', track.subject],
    ['UID', track.uid],
  ];
  return (
    <div
      className={`${c.overlay} ${c.detail}`}
      role="complementary"
      aria-label="Track detail"
    >
      <div style={{ display: 'flex', alignItems: 'center', gap: 6 }}>
        <strong style={{ flex: 1, fontSize: 14 }}>
          {track.callsign || track.uid}
        </strong>
        <button
          type="button"
          className={c.toggle}
          onClick={props.onClose}
          aria-label="Close track detail"
        >
          Close
        </button>
      </div>
      <dl className={c.dl}>
        {rows.map(([k, v]) => (
          <div key={k} style={{ display: 'contents' }}>
            <dt className={c.dt}>{k}</dt>
            <dd className={c.dd}>{v}</dd>
          </div>
        ))}
      </dl>
    </div>
  );
}

/** Why there are no tracks, said honestly: a down or disabled source is not the same as a quiet one. */
export function emptyTracksMessage(nats?: SourceStatus): string {
  if (!nats) return 'Waiting for the first status from the backend.';
  if (nats.state === 'disabled')
    return 'Tracks unavailable: the NATS source is disabled on the backend, so no positions are collected.';
  if (nats.state === 'down')
    return 'Tracks unavailable: NATS is down. No positions are arriving.';
  if (nats.state === 'connecting')
    return 'NATS is still connecting. No positions yet.';
  return 'NATS is connected but no position has been received yet.';
}

export function GlobePanel(props: {
  tracks: ReadonlyMap<string, Track>;
  now: number;
  nats?: SourceStatus;
}) {
  const c = useStyles();
  const { tracks, now } = props;
  const [filter, setFilter] = useState<TrackFilter>(ALL_TRACKS);
  const [query, setQuery] = useState('');
  const [selectedUid, setSelectedUid] = useState<string>();
  const [flyTo, setFlyTo] = useState<{
    uid: string;
    lat: number;
    lon: number;
    nonce: number;
  }>();
  const [globeError, setGlobeError] = useState<string>();

  const counts = useMemo(() => countTracks(tracks.values()), [tracks]);
  const visible = useMemo(
    () => filterTracks(tracks.values(), filter),
    [tracks, filter],
  );
  const results = useMemo(
    () => searchTracks(tracks.values(), query),
    [tracks, query],
  );
  // The fade only needs a coarse clock.
  const fadeNowMs = Math.floor(now / 5000) * 5000;
  const selected = selectedUid ? tracks.get(selectedUid) : undefined;

  const pick = (t: Track) => {
    setSelectedUid(t.uid);
    setFlyTo({ uid: t.uid, lat: t.lat, lon: t.lon, nonce: Date.now() });
    setQuery('');
  };

  const natsDown =
    props.nats &&
    (props.nats.state === 'down' || props.nats.state === 'disabled');
  const emptyMessage =
    tracks.size > 0 ? undefined : emptyTracksMessage(props.nats);

  return (
    <Panel
      title="Tracks"
      aside={
        <span style={{ textTransform: 'none', letterSpacing: 0 }}>
          {natsDown
            ? 'source unavailable'
            : `${visible.length} shown of ${counts.total}`}
        </span>
      }
      bodyStyle={{
        display: 'flex',
        flexDirection: 'column',
        overflow: 'hidden',
      }}
    >
      <div className={c.stage}>
        <GlobeBoundary
          fallback={message => (
            <div className={c.empty}>
              <div className={c.emptyBox}>
                <Note tone="critical">
                  The globe could not start: {message}. The track list and
                  filters still work; check WebGL support in this browser.
                </Note>
              </div>
            </div>
          )}
        >
          <Suspense fallback={<Progress />}>
            <CesiumGlobe
              tracks={visible}
              fadeNowMs={fadeNowMs}
              selectedUid={selectedUid}
              onSelect={setSelectedUid}
              flyTo={flyTo}
              onError={setGlobeError}
            />
          </Suspense>
        </GlobeBoundary>

        {globeError && (
          <div className={c.empty}>
            <div className={c.emptyBox}>
              <Note tone="critical">
                The globe could not start: {globeError}. Check WebGL support,
                and that /cesium/ assets are served.
              </Note>
            </div>
          </div>
        )}
        {emptyMessage && !globeError && (
          <div className={c.empty}>
            <div className={c.emptyBox}>
              <Note tone={natsDown ? 'critical' : undefined}>
                {emptyMessage}
              </Note>
            </div>
          </div>
        )}

        <div
          className={`${c.overlay} ${c.toolbar}`}
          role="group"
          aria-label="Track filters"
        >
          <div className={c.row} aria-label="Affiliation">
            {AFFILIATIONS.map((a: Affiliation) => (
              <button
                key={a}
                type="button"
                className={c.toggle}
                aria-pressed={filter.affiliations.has(a)}
                onClick={() =>
                  setFilter(f => ({
                    ...f,
                    affiliations: toggle(f.affiliations, a),
                  }))
                }
              >
                <span
                  className={c.swatch}
                  style={{ background: AFFILIATION_COLOR[a] }}
                />
                {AFFILIATION_LABEL[a]}
                <span className={c.count}>{counts.affiliation[a]}</span>
              </button>
            ))}
          </div>
          <div className={c.row} aria-label="Domain">
            {DOMAINS.map((d: Domain) => (
              <button
                key={d}
                type="button"
                className={c.toggle}
                aria-pressed={filter.domains.has(d)}
                onClick={() =>
                  setFilter(f => ({ ...f, domains: toggle(f.domains, d) }))
                }
              >
                {DOMAIN_LABEL[d]}
                <span className={c.count}>{counts.domain[d]}</span>
              </button>
            ))}
          </div>
          <div>
            <input
              className={c.search}
              aria-label="Search callsign or uid"
              placeholder="Search callsign or uid"
              value={query}
              onChange={e => setQuery(e.target.value)}
            />
            {query.trim() !== '' && (
              <ul className={c.results} aria-label="Search results">
                {results.length === 0 ? (
                  <li
                    style={{ padding: '3px 6px', color: statusFg('unknown') }}
                  >
                    No track matches
                  </li>
                ) : (
                  results.map(t => (
                    <li key={t.uid}>
                      <button
                        type="button"
                        className={c.result}
                        onClick={() => pick(t)}
                      >
                        {t.callsign || t.uid}{' '}
                        <span className={c.count}>{decodeCotType(t.type)}</span>
                      </button>
                    </li>
                  ))
                )}
              </ul>
            )}
          </div>
        </div>

        {selectedUid && (
          <TrackDetail
            track={selected}
            selectedUid={selectedUid}
            now={now}
            onClose={() => setSelectedUid(undefined)}
          />
        )}
        <div className={`${c.overlay} ${c.attribution}`}>
          Map © OpenStreetMap contributors
        </div>
      </div>
    </Panel>
  );
}
