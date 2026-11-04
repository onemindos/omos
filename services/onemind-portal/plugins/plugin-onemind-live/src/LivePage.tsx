import { Box, Button, Link, Typography, useTheme } from '@material-ui/core';
import { StatusChip } from '@internal/plugin-onemind-theme';
import { clock } from './format';
import { GrafanaEmbeds } from './GrafanaEmbeds';
import { useOwnerIndex } from './owners';
import {
  AlertsPanel,
  ArgoPanel,
  BackupsPanel,
  EventsPanel,
  NodesPanel,
  StoragePanel,
  UptimePanel,
  WorkloadsPanel,
} from './panels';
import { useLive } from './useLive';
import type { LiveResponse } from './types';

const LEVEL = {
  ok: 'ok',
  warn: 'warn',
  critical: 'critical',
  unknown: 'unknown',
} as const;
const LEVEL_LABEL = { ok: 'Healthy', warn: 'Needs attention', critical: 'Critical', unknown: 'Unknown' } as const;

export function HealthBanner({ data }: { data: LiveResponse }) {
  const theme = useTheme();
  const { health, counts } = data;
  const tone = `var(--onemind-status-${LEVEL[health.level]})`;
  return (
    <Box
      component="section"
      aria-label="Overall health"
      role="status"
      p={2}
      style={{
        border: `1px solid ${theme.palette.divider}`,
        borderLeft: `6px solid ${tone}`,
        borderRadius: theme.shape.borderRadius * 2,
        background: theme.palette.background.paper,
      }}
    >
      <Box display="flex" flexWrap="wrap" alignItems="center" style={{ gap: 12 }}>
        <StatusChip status={LEVEL[health.level]} label={LEVEL_LABEL[health.level]} />
        <Typography variant="h5" component="p" style={{ fontWeight: 700, flex: '1 1 260px' }}>
          {health.headline}
        </Typography>
      </Box>
      {health.reasons.length > 1 && (
        <Typography variant="body2" color="textSecondary" style={{ marginTop: 4 }}>
          {health.reasons.join(' · ')}
        </Typography>
      )}
      <Typography variant="caption" color="textSecondary" display="block" style={{ marginTop: 8 }}>
        Alerts {String(counts.alerts)} · Argo {String(counts.argo)} · Workloads not ready {String(counts.workloads)} ·
        Backups to review {String(counts.backups)} · Needs setup {String(counts.backupsNeedSetup)} · Retired{' '}
        {String(counts.backupsRetired)}
      </Typography>
    </Box>
  );
}

export function LivePage() {
  const theme = useTheme();
  const { data, error, refresh, push, refreshedAt } = useLive();
  const { index } = useOwnerIndex();
  return (
    <Box p={2} className="onemind-live">
      <Box display="flex" flexWrap="wrap" alignItems="baseline" style={{ gap: 12, marginBottom: theme.spacing(2) }}>
        <Box style={{ flex: '1 1 280px' }}>
          <Typography variant="overline">ONE MIND · OPERATIONS</Typography>
          <Typography variant="h4" component="h1" style={{ fontWeight: 700 }}>
            Live Ops
          </Typography>
          <Typography color="textSecondary">What is wrong right now, and where.</Typography>
        </Box>
        <Box textAlign="right">
          <Typography variant="caption" color="textSecondary" display="block" role="status">
            {push ? 'Live: alert Signals plus a 30 second refresh' : 'Refreshing every 30 seconds'}
            {refreshedAt ? ` · updated ${clock(new Date(refreshedAt).toISOString())}` : ''}
          </Typography>
          <Button size="small" variant="outlined" onClick={() => refresh()}>
            Refresh now
          </Button>
        </Box>
      </Box>
      {error && (
        <Typography role="alert" variant="body2" style={{ color: 'var(--onemind-status-critical)', marginBottom: 8 }}>
          {data
            ? `Live refresh failed. Showing results observed at ${clock(data.observedAt)}.`
            : 'Live data is unavailable. This page cannot confirm anything is healthy.'}
        </Typography>
      )}
      {!data && !error && <Typography role="status">Reading the cluster…</Typography>}
      {data && (
        <>
          <HealthBanner data={data} />
          <Box
            mt={2}
            style={{
              display: 'grid',
              gridTemplateColumns: 'repeat(auto-fit, minmax(min(100%, 420px), 1fr))',
              gap: theme.spacing(2),
              alignItems: 'start',
            }}
          >
            <AlertsPanel data={data} index={index} />
            <ArgoPanel data={data} index={index} />
            <WorkloadsPanel data={data} index={index} />
            <NodesPanel data={data} />
            <StoragePanel data={data} index={index} />
            <BackupsPanel data={data} />
            <UptimePanel data={data} />
            <EventsPanel data={data} index={index} />
            <GrafanaEmbeds origin={data.links.grafana} />
          </Box>
          <Typography variant="caption" color="textSecondary" display="block" style={{ marginTop: 16 }}>
            Observed {clock(data.observedAt)} on {data.cluster}. Each panel shows its own observation time; a source that
            did not answer is marked unavailable, never healthy. Links open Argo CD, Headlamp and Grafana;{' '}
            <Link href="/today">Today</Link> shows the same attention list with your inbox.
          </Typography>
        </>
      )}
    </Box>
  );
}
