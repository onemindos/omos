import { useMemo } from 'react';
import { makeStyles } from '@material-ui/core/styles';
import { Tooltip } from '@material-ui/core';
import { Link } from '@backstage/core-components';
import type {
  ArgoCell,
  NamespaceCell,
  SourceStatus,
  WarningEvent,
} from './contract';
import { ago } from './cot';
import { sortArgo, sortCluster } from './liveState';
import type { NamespaceLinks } from './namespaceLinks';
import {
  CELL_STATUS,
  Dot,
  Note,
  Panel,
  cssVar,
  statusBg,
  statusFg,
} from './ui';
import type { StatusVar } from './ui';

/** A cell shows the pulse for this long after a delta changed it. */
export const PULSE_MS = 3000;

const useStyles = makeStyles({
  grid: {
    display: 'grid',
    gridTemplateColumns: 'repeat(auto-fill, minmax(112px, 1fr))',
    gap: 4,
    padding: 8,
  },
  cell: {
    display: 'block',
    padding: '5px 7px',
    borderRadius: 4,
    border: '1px solid transparent',
    fontSize: 12,
    lineHeight: 1.25,
    color: 'inherit',
    textDecoration: 'none',
    overflow: 'hidden',
    '&:hover, &:focus-visible': {
      borderColor: cssVar('text-2'),
      textDecoration: 'none',
    },
  },
  cellName: {
    fontWeight: 700,
    whiteSpace: 'nowrap',
    overflow: 'hidden',
    textOverflow: 'ellipsis',
  },
  cellMeta: { fontSize: 11, opacity: 0.85, fontVariantNumeric: 'tabular-nums' },
  '@keyframes opsPulse': {
    '0%': {
      boxShadow: `0 0 0 0 ${cssVar('text-1')}`,
      transform: 'scale(1.04)',
    },
    '100%': { boxShadow: '0 0 0 8px transparent', transform: 'scale(1)' },
  },
  pulse: { animation: '$opsPulse 1.2s ease-out 2' },
  strip: { display: 'flex', flexWrap: 'wrap', gap: 4, padding: 8 },
  chip: {
    display: 'flex',
    alignItems: 'center',
    gap: 5,
    padding: '2px 7px',
    borderRadius: 999,
    fontSize: 11,
    border: `1px solid ${cssVar('border')}`,
    whiteSpace: 'nowrap',
  },
  warnRow: {
    display: 'grid',
    gridTemplateColumns: '58px 1fr',
    gap: 8,
    padding: '4px 10px',
    fontSize: 12,
    borderBottom: `1px solid ${cssVar('surface-2')}`,
  },
  muted: { color: cssVar('text-2') },
});

function cellTitle(
  cell: NamespaceCell,
  now: number,
  target?: { title: string; via: string },
) {
  const changed = cell.lastChangeAt ? Date.parse(cell.lastChangeAt) : undefined;
  return (
    <div>
      <strong>{cell.namespace}</strong>
      <div>
        {cell.ready} of {cell.total} pods ready · {cell.pending} pending ·{' '}
        {cell.failed} failed
      </div>
      <div>{cell.restarts} restarts in the last hour</div>
      <div>
        {cell.lastChange
          ? `${cell.lastChange}${
              changed !== undefined ? `, ${ago(changed, now)}` : ''
            }`
          : 'No recent change recorded'}
      </div>
      <div>
        {target
          ? `Opens ${target.title} (${target.via})`
          : 'No catalog system or Headlamp link for this namespace'}
      </div>
    </div>
  );
}

function unavailableNote(
  label: string,
  source?: SourceStatus,
): string | undefined {
  if (!source)
    return `${label}: waiting for the first status from the backend.`;
  if (source.state === 'disabled')
    return `${label} is disabled on the backend${
      source.detail ? ` (${source.detail})` : ''
    }. Nothing is being collected.`;
  if (source.state === 'down')
    return `${label} is down${
      source.detail ? `: ${source.detail}` : ''
    }. Values below are the last ones received, if any.`;
  if (source.state === 'connecting') return `${label} is still connecting.`;
  return undefined;
}

const sourceTone = (source?: SourceStatus): StatusVar => {
  if (source?.state === 'down') return 'critical';
  if (source?.state === 'connecting') return 'warn';
  return 'unknown';
};

export function ClusterGrid(props: {
  cells: ReadonlyMap<string, NamespaceCell>;
  changedMs: ReadonlyMap<string, number>;
  now: number;
  links: NamespaceLinks;
  source?: SourceStatus;
}) {
  const c = useStyles();
  const sorted = useMemo(
    () => sortCluster(props.cells.values()),
    [props.cells],
  );
  const note = unavailableNote('Kubernetes', props.source);
  const counts = { ok: 0, warn: 0, down: 0 };
  for (const cell of sorted) counts[cell.health] += 1;
  return (
    <Panel
      title="Namespaces"
      aside={
        sorted.length > 0 ? (
          <span style={{ textTransform: 'none', letterSpacing: 0 }}>
            {counts.down} down · {counts.warn} warn · {counts.ok} ok
          </span>
        ) : undefined
      }
    >
      {note && <Note tone={sourceTone(props.source)}>{note}</Note>}
      {sorted.length === 0 ? (
        !note && (
          <Note>Kubernetes is connected but reported no namespaces.</Note>
        )
      ) : (
        <div className={c.grid} role="list">
          {sorted.map(cell => {
            const tone: StatusVar = CELL_STATUS[cell.health];
            const target = props.links(cell.namespace);
            const changed = props.changedMs.get(cell.namespace);
            const pulsing =
              changed !== undefined && props.now - changed < PULSE_MS;
            const body = (
              <>
                <div className={c.cellName}>{cell.namespace}</div>
                <div className={c.cellMeta}>
                  {cell.ready}/{cell.total} ready
                  {cell.restarts > 0 ? ` · ${cell.restarts}↻` : ''}
                </div>
              </>
            );
            const common = {
              className: `${c.cell} ${pulsing ? c.pulse : ''}`,
              style: { background: statusBg(tone), color: statusFg(tone) },
              'data-health': cell.health,
              'data-namespace': cell.namespace,
            };
            return (
              <Tooltip
                key={`${cell.namespace}:${pulsing ? changed : 0}`}
                title={cellTitle(cell, props.now, target)}
                arrow
              >
                <div role="listitem">
                  {target ? (
                    <Link
                      to={target.href}
                      {...common}
                      aria-label={`${cell.namespace}, ${cell.health}`}
                    >
                      {body}
                    </Link>
                  ) : (
                    <div
                      {...common}
                      aria-label={`${cell.namespace}, ${cell.health}`}
                    >
                      {body}
                    </div>
                  )}
                </div>
              </Tooltip>
            );
          })}
        </div>
      )}
    </Panel>
  );
}

const argoTone = (a: ArgoCell): StatusVar => {
  if (a.sync === 'Synced' && a.health === 'Healthy') return 'ok';
  if (a.health === 'Degraded' || a.health === 'Missing') return 'critical';
  if (a.health === 'Unknown') return 'unknown';
  return 'warn';
};

export function ArgoStrip(props: {
  apps: ReadonlyMap<string, ArgoCell>;
  source?: SourceStatus;
  now: number;
}) {
  const c = useStyles();
  const sorted = useMemo(() => sortArgo(props.apps.values()), [props.apps]);
  const note = unavailableNote('Argo CD', props.source);
  const good = sorted.filter(a => argoTone(a) === 'ok').length;
  return (
    <Panel
      title="Argo applications"
      aside={
        sorted.length > 0 ? (
          <span style={{ textTransform: 'none', letterSpacing: 0 }}>
            {good} of {sorted.length} Synced and Healthy
          </span>
        ) : undefined
      }
    >
      {note && <Note tone={sourceTone(props.source)}>{note}</Note>}
      {sorted.length === 0 ? (
        !note && <Note>Argo CD is connected but reported no applications.</Note>
      ) : (
        <div className={c.strip}>
          {sorted.map(a => {
            const tone = argoTone(a);
            const changed = a.changedAt ? Date.parse(a.changedAt) : undefined;
            return (
              <Tooltip
                key={a.name}
                arrow
                title={`${a.name}: ${a.sync}, ${a.health}${
                  a.revision ? `, revision ${a.revision.slice(0, 8)}` : ''
                }${
                  changed !== undefined
                    ? `, changed ${ago(changed, props.now)}`
                    : ''
                }`}
              >
                <span
                  className={c.chip}
                  data-tone={tone}
                  style={
                    tone === 'ok'
                      ? undefined
                      : {
                          background: statusBg(tone),
                          color: statusFg(tone),
                          borderColor: statusFg(tone),
                        }
                  }
                >
                  <Dot status={tone} size={7} />
                  {a.name}
                  {tone !== 'ok' && (
                    <span>
                      {a.sync} / {a.health}
                    </span>
                  )}
                </span>
              </Tooltip>
            );
          })}
        </div>
      )}
    </Panel>
  );
}

export function WarningsFeed(props: {
  warnings: WarningEvent[];
  now: number;
  source?: SourceStatus;
}) {
  const c = useStyles();
  const note = unavailableNote('Kubernetes', props.source);
  return (
    <Panel
      title="Kubernetes warnings"
      aside={
        props.warnings.length > 0 ? (
          <span>{props.warnings.length}</span>
        ) : undefined
      }
    >
      {note && <Note tone={sourceTone(props.source)}>{note}</Note>}
      {props.warnings.length === 0
        ? !note && (
            <Note>
              No Warning events held. The cluster has not reported one recently.
            </Note>
          )
        : props.warnings.map(w => {
            const t = Date.parse(w.ts);
            return (
              <div
                key={`${w.namespace}/${w.object}/${w.reason}/${w.ts}`}
                className={c.warnRow}
              >
                <span
                  className={c.muted}
                  style={{ fontVariantNumeric: 'tabular-nums' }}
                >
                  {Number.isFinite(t) ? ago(t, props.now) : '?'}
                </span>
                <span>
                  <strong style={{ color: statusFg('warn') }}>
                    {w.reason}
                  </strong>{' '}
                  {w.namespace}/{w.object}
                  {w.count && w.count > 1 ? ` ×${w.count}` : ''}
                  <div className={c.muted}>{w.message}</div>
                </span>
              </div>
            );
          })}
    </Panel>
  );
}
