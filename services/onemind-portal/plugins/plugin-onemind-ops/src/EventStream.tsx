import { useMemo, useState } from 'react';
import { makeStyles } from '@material-ui/core/styles';
import type { BusEvent } from './contract';
import { matchesEvent } from './filters';
import { SERIES_SECONDS, kindSeries } from './liveState';
import type { Bucket } from './liveState';
import {
  KIND_ORDER,
  Note,
  Panel,
  Sparkline,
  cssVar,
  statusFg,
  useKindColor,
} from './ui';

/** Rows in the DOM. The state keeps more; older rows are simply not drawn. */
export const MAX_ROWS = 150;

const useStyles = makeStyles({
  tools: {
    display: 'flex',
    flexWrap: 'wrap',
    gap: 4,
    padding: '6px 8px',
    borderBottom: `1px solid ${cssVar('border')}`,
  },
  chip: {
    display: 'flex',
    alignItems: 'center',
    gap: 5,
    border: `1px solid ${cssVar('border')}`,
    background: 'transparent',
    color: cssVar('text-1'),
    borderRadius: 4,
    padding: '2px 6px',
    fontSize: 11,
    cursor: 'pointer',
    '&[aria-pressed=false]': { opacity: 0.4 },
    '&:focus-visible': { outline: `2px solid ${cssVar('focus')}` },
  },
  input: {
    flex: '1 1 120px',
    minWidth: 100,
    background: cssVar('surface-0'),
    color: cssVar('text-1'),
    border: `1px solid ${cssVar('border')}`,
    borderRadius: 4,
    padding: '3px 6px',
    fontSize: 12,
  },
  btn: {
    background: cssVar('surface-3'),
    color: cssVar('text-1'),
    border: `1px solid ${cssVar('border-strong')}`,
    borderRadius: 4,
    padding: '2px 10px',
    fontSize: 12,
    cursor: 'pointer',
  },
  row: {
    display: 'grid',
    gridTemplateColumns: '58px 40px 1fr',
    gap: 6,
    padding: '2px 8px',
    fontSize: 12,
    fontFamily: 'ui-monospace, SFMono-Regular, Menlo, monospace',
    borderBottom: `1px solid ${cssVar('surface-2')}`,
  },
  time: { color: cssVar('text-3') },
  subject: {
    color: cssVar('text-1'),
    overflow: 'hidden',
    textOverflow: 'ellipsis',
    whiteSpace: 'nowrap',
  },
  summary: { color: cssVar('text-2') },
  stats: {
    fontSize: 11,
    color: cssVar('text-2'),
    padding: '4px 8px',
    borderTop: `1px solid ${cssVar('border')}`,
  },
});

const clock = (iso: string) => {
  const d = new Date(iso);
  return Number.isNaN(d.getTime())
    ? '--:--:--'
    : d.toLocaleTimeString([], { hour12: false });
};

export function EventStream(props: {
  events: BusEvent[];
  dropped: number;
  buckets: Bucket[];
  now: number;
  /** Why there are no events, when the cause is a source and not simply a quiet bus. */
  unavailable?: string;
}) {
  const c = useStyles();
  const kindColor = useKindColor();
  const [paused, setPaused] = useState<BusEvent[] | undefined>();
  const [hidden, setHidden] = useState<ReadonlySet<string>>(new Set());
  const [text, setText] = useState('');

  const source = paused ?? props.events;
  const rows = useMemo(() => {
    const filter = { hiddenKinds: hidden, text };
    const out: BusEvent[] = [];
    for (const e of source) {
      if (matchesEvent(e, filter)) out.push(e);
      if (out.length >= MAX_ROWS) break;
    }
    return out;
  }, [source, hidden, text]);

  const series = useMemo(
    () => kindSeries(props.buckets, props.now),
    [props.buckets, props.now],
  );
  const kinds = useMemo(() => {
    const extra = Object.keys(series).filter(
      k => !(KIND_ORDER as readonly string[]).includes(k),
    );
    return [...KIND_ORDER, ...extra];
  }, [series]);

  let emptyNote: { text: string; tone?: 'warn' } = {
    text: 'No events match the filter.',
  };
  if (props.unavailable) emptyNote = { text: props.unavailable, tone: 'warn' };
  else if (source.length === 0)
    emptyNote = {
      text: 'No bus events received yet. NATS is connected, but nothing has been published since this page loaded.',
    };

  const toggleKind = (k: string) =>
    setHidden(prev => {
      const next = new Set(prev);
      if (next.has(k)) next.delete(k);
      else next.add(k);
      return next;
    });

  return (
    <Panel
      title="Bus events"
      bodyStyle={{ overflow: 'hidden' }}
      aside={
        <button
          type="button"
          className={c.btn}
          onClick={() => setPaused(paused ? undefined : props.events)}
        >
          {paused
            ? `Resume (${Math.max(
                0,
                props.events.length - paused.length,
              )}+ new)`
            : 'Pause'}
        </button>
      }
    >
      <div style={{ display: 'flex', flexDirection: 'column', height: '100%' }}>
        <div className={c.tools}>
          {kinds.map(k => {
            const values = series[k] ?? new Array(SERIES_SECONDS).fill(0);
            const perSecond = values.slice(-10).reduce((a, b) => a + b, 0) / 10;
            return (
              <button
                key={k}
                type="button"
                className={c.chip}
                aria-pressed={!hidden.has(k)}
                onClick={() => toggleKind(k)}
                title={`${k}: ${perSecond.toFixed(
                  1,
                )}/s over the last 10 s. Click to ${
                  hidden.has(k) ? 'show' : 'hide'
                }.`}
              >
                <span style={{ color: kindColor(k), fontWeight: 700 }}>
                  {k}
                </span>
                <Sparkline
                  values={values}
                  color={kindColor(k)}
                  width={40}
                  height={12}
                />
                <span
                  style={{
                    fontVariantNumeric: 'tabular-nums',
                    color: cssVar('text-2'),
                  }}
                >
                  {perSecond.toFixed(1)}
                </span>
              </button>
            );
          })}
          <input
            className={c.input}
            aria-label="Filter events"
            placeholder="Filter subject or summary"
            value={text}
            onChange={e => setText(e.target.value)}
          />
        </div>
        <div
          style={{ flex: 1, minHeight: 0, overflow: 'auto' }}
          aria-live="off"
        >
          {rows.length === 0 ? (
            <Note tone={emptyNote.tone}>{emptyNote.text}</Note>
          ) : (
            rows.map(e => (
              <div key={e.id} className={c.row} data-kind={e.kind}>
                <span className={c.time}>{clock(e.ts)}</span>
                <span style={{ color: kindColor(e.kind), fontWeight: 700 }}>
                  {e.kind}
                </span>
                <span
                  className={c.subject}
                  title={`${e.subject}\n${e.summary}\n${e.bytes} bytes`}
                >
                  {e.subject} <span className={c.summary}>{e.summary}</span>
                </span>
              </div>
            ))
          )}
        </div>
        {!(props.unavailable && source.length === 0) && (
          <div className={c.stats}>
            {rows.length} shown of {source.length} held ·{' '}
            <span
              style={{
                color: props.dropped > 0 ? statusFg('warn') : undefined,
              }}
            >
              {props.dropped} dropped by the server (rate limited)
            </span>
          </div>
        )}
      </div>
    </Panel>
  );
}
