import { makeStyles } from '@material-ui/core/styles';
import { Tooltip } from '@material-ui/core';
import type { SourceStatus } from './contract';
import { ago } from './cot';
import type { FeedState, Freshness } from './liveState';
import { CONNECTION_STATUS, Dot, cssVar, statusBg, statusFg } from './ui';
import type { StatusVar } from './ui';

const useStyles = makeStyles({
  bar: {
    display: 'flex',
    alignItems: 'stretch',
    flexWrap: 'wrap',
    gap: 1,
    background: cssVar('border'),
    border: `1px solid ${cssVar('border')}`,
    borderRadius: 6,
    overflow: 'hidden',
  },
  cell: {
    display: 'flex',
    alignItems: 'center',
    gap: 8,
    padding: '8px 14px',
    background: cssVar('surface-1'),
    flex: '1 1 200px',
    minWidth: 0,
  },
  name: { fontWeight: 700, fontSize: 13, color: cssVar('text-1') },
  meta: {
    fontSize: 12,
    color: cssVar('text-2'),
    whiteSpace: 'nowrap',
    overflow: 'hidden',
    textOverflow: 'ellipsis',
  },
  num: { fontVariantNumeric: 'tabular-nums' },
  badge: {
    display: 'flex',
    alignItems: 'center',
    gap: 8,
    padding: '8px 16px',
    fontWeight: 800,
    fontSize: 13,
    letterSpacing: '0.1em',
    flex: 'none',
  },
});

const STATE_WORD: Record<SourceStatus['state'], string> = {
  connected: 'connected',
  connecting: 'connecting',
  down: 'down',
  disabled: 'disabled',
};

const FRESH: Record<Freshness, { label: string; tone: StatusVar }> = {
  live: { label: 'LIVE', tone: 'ok' },
  stale: { label: 'STALE', tone: 'warn' },
  waiting: { label: 'WAITING', tone: 'unknown' },
};

function Source(props: { name: string; source?: SourceStatus; now: number }) {
  const c = useStyles();
  const { name, source, now } = props;
  const tone: StatusVar = source ? CONNECTION_STATUS[source.state] : 'unknown';
  const last = source?.lastMessageAt
    ? Date.parse(source.lastMessageAt)
    : undefined;
  const running = source?.state === 'connected';
  return (
    <Tooltip
      title={
        source?.detail ??
        (source ? STATE_WORD[source.state] : 'No status received yet')
      }
      arrow
    >
      <div className={c.cell} data-source={name}>
        <Dot status={tone} pulse={running} />
        <span className={c.name}>{name}</span>
        <span className={c.meta} style={{ color: statusFg(tone) }}>
          {source ? STATE_WORD[source.state] : 'waiting'}
        </span>
        {source && source.state !== 'disabled' && (
          <>
            <span className={`${c.meta} ${c.num}`}>
              {source.rate === undefined
                ? 'rate unknown'
                : `${source.rate.toFixed(source.rate < 10 ? 1 : 0)} msg/s`}
            </span>
            <span className={`${c.meta} ${c.num}`}>
              {last === undefined
                ? 'no message yet'
                : `last message ${ago(last, now)}`}
            </span>
          </>
        )}
      </div>
    </Tooltip>
  );
}

export function ConnectionBar(props: {
  state: FeedState;
  freshness: Freshness;
  now: number;
  signalsConnected: boolean;
}) {
  const c = useStyles();
  const { state, freshness, now, signalsConnected } = props;
  const f = FRESH[freshness];
  return (
    <div className={c.bar} role="region" aria-label="Connections">
      <Source name="NATS" source={state.status?.nats} now={now} />
      <Source name="Kubernetes" source={state.status?.kubernetes} now={now} />
      <Source name="Argo" source={state.status?.argo} now={now} />
      <Tooltip
        title={
          signalsConnected
            ? 'Deltas arrive through Signals. STALE means no status message for 15 seconds.'
            : 'The Signals connection is not available, so only snapshots refresh this page.'
        }
        arrow
      >
        <div
          className={c.badge}
          role="status"
          aria-label={`Feed ${f.label}`}
          style={{ background: statusBg(f.tone), color: statusFg(f.tone) }}
        >
          <Dot status={f.tone} pulse={freshness === 'live'} />
          {f.label}
          {!signalsConnected && <span className={c.meta}> (no push)</span>}
        </div>
      </Tooltip>
    </div>
  );
}
