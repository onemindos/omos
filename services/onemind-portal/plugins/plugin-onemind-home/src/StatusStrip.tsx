import { Box, Link, Theme, Typography, useTheme } from '@material-ui/core';
import {
  isStale,
  summarizeStatus,
  StripCounts,
  useOperationsStatus,
} from './status';

type Tile = {
  key: keyof StripCounts;
  label: string;
  href: string;
  external?: boolean;
  hint: string;
};

// Links reuse the destinations `/today` already offers for each source.
const TILES: Tile[] = [
  {
    key: 'alerts',
    label: 'Alerts firing',
    href: '/live',
    hint: 'Firing Alertmanager alerts and service checks that are not up',
  },
  {
    key: 'argo',
    label: 'Argo apps to review',
    href: 'https://argocd.tail717ba0.ts.net',
    external: true,
    hint: 'Out of sync or not healthy',
  },
  {
    key: 'workloads',
    label: 'Workloads not ready',
    href: 'https://headlamp.tail717ba0.ts.net/c/main/workloads',
    external: true,
    hint: 'Fewer ready replicas than desired',
  },
  {
    key: 'backups',
    label: 'Backups to review',
    href: '/live',
    hint: 'Failing, overdue, or suspended without a reason. Retired jobs are not counted.',
  },
];

const observedTime = (iso?: string) =>
  iso
    ? new Date(iso).toLocaleTimeString([], {
        hour: '2-digit',
        minute: '2-digit',
      })
    : '';

type Value = number | 'unknown' | undefined;

function accentFor(value: Value, palette: Theme['palette']) {
  if (value === 'unknown') return palette.warning.main;
  if (typeof value === 'number' && value > 0) return palette.error.main;
  return palette.success.main;
}

function valueLabel(value: Value) {
  if (value === undefined) return '—';
  return value === 'unknown' ? 'Unknown' : String(value);
}

// Retired and needs-setup backups are shown beside the count, never inside it.
function backupNote(counts: StripCounts | undefined) {
  if (!counts) return '';
  const parts: string[] = [];
  const setup = counts.backupsNeedSetup;
  const retired = counts.backupsRetired;
  if (setup === 'unknown' || retired === 'unknown') return '';
  if (setup) parts.push(`${setup} needs setup`);
  if (retired) parts.push(`${retired} retired`);
  return parts.join(' · ');
}

function observedLabel(
  observed: string | undefined,
  stale: boolean,
  error: boolean,
) {
  if (observed)
    return `${stale ? 'Stale · ' : ''}observed ${observedTime(observed)}`;
  return error ? 'Live update unavailable' : 'Waiting for observation';
}

export function StatusStrip() {
  const theme = useTheme();
  const { data, error } = useOperationsStatus();
  const stale = !!data && (error || isStale(data));
  const counts = data ? summarizeStatus(data) : undefined;
  return (
    <Box
      component="section"
      aria-label="Live status"
      className="onemind-status-strip"
      style={{
        display: 'grid',
        gridTemplateColumns: 'repeat(auto-fit,minmax(min(100%,200px),1fr))',
        gap: theme.spacing(1.5),
      }}
    >
      {TILES.map(tile => {
        const value = counts?.[tile.key];
        const accent = accentFor(value, theme.palette);
        const observed =
          tile.key === 'alerts'
            ? data?.serviceHealth?.observedAt ?? data?.observedAt
            : data?.observedAt;
        return (
          <Link
            key={tile.key}
            href={tile.href}
            {...(tile.external
              ? { target: '_blank', rel: 'noopener noreferrer' }
              : {})}
            underline="none"
            color="inherit"
            aria-label={`${tile.label}: ${value ?? 'loading'}. ${tile.hint}.`}
            style={{
              display: 'block',
              padding: theme.spacing(1.5, 2),
              minHeight: 44,
              border: `1px solid ${theme.palette.divider}`,
              borderLeft: `4px solid ${accent}`,
              borderRadius: theme.shape.borderRadius * 2,
              background: theme.palette.background.paper,
            }}
          >
            <Typography variant="caption" color="textSecondary" display="block">
              {tile.label}
            </Typography>
            <Typography
              variant="h4"
              component="p"
              style={{ fontWeight: 700, lineHeight: 1.2 }}
            >
              {valueLabel(value)}
            </Typography>
            {tile.key === 'backups' && backupNote(counts) && (
              <Typography variant="caption" display="block">
                {backupNote(counts)}
              </Typography>
            )}
            <Typography variant="caption" color="textSecondary" display="block">
              {observedLabel(observed, stale, error)}
            </Typography>
          </Link>
        );
      })}
      {data?.partial && (
        <Typography
          variant="caption"
          color="error"
          role="status"
          style={{ gridColumn: '1 / -1' }}
        >
          Incomplete view:{' '}
          {data.unavailableSources?.join(', ') || 'some sources'} unavailable.
          Unknown is not healthy.
        </Typography>
      )}
      {error && (
        <Typography
          variant="caption"
          color="error"
          role="alert"
          style={{ gridColumn: '1 / -1' }}
        >
          Live refresh unavailable; previous results may be stale.
        </Typography>
      )}
    </Box>
  );
}
