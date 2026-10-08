import { useEffect, useState } from 'react';
import {
  discoveryApiRef,
  fetchApiRef,
  useApi,
} from '@backstage/core-plugin-api';
import { Box, Button, Typography, useTheme } from '@material-ui/core';
import { ActivityInbox } from './ActivityInbox';
type Observation = {
  observedAt: string;
  partial: boolean;
  totalWorkloads: number;
  unavailable: string[];
  unavailableSources?: string[];
  deployments: { name: string; health: string; sync: string }[];
  backups: { name: string; lastSuccess: string | null; suspended: boolean }[];
  serviceHealth?: { checks: { state: string }[] };
};
export function LiveOverview() {
  const discovery = useApi(discoveryApiRef),
    api = useApi(fetchApiRef),
    theme = useTheme();
  const [data, setData] = useState<Observation>(),
    [error, setError] = useState(false),
    [inbox, setInbox] = useState(false);
  useEffect(() => {
    let active = true,
      pending = false;
    const controller = new AbortController();
    const refresh = async () => {
      if (pending) return;
      pending = true;
      try {
        const base = await discovery.getBaseUrl('onemind');
        const response = await api.fetch(base + '/status', {
          signal: controller.signal,
        });
        if (!response.ok) throw new Error();
        const value = await response.json();
        if (active) {
          setData(value);
          setError(false);
        }
      } catch {
        if (active) setError(true);
      } finally {
        pending = false;
      }
    };
    refresh();
    const timer = setInterval(refresh, 60000);
    return () => {
      active = false;
      controller.abort();
      clearInterval(timer);
    };
  }, [discovery, api]);
  const stale =
    !!data && (error || Date.now() - Date.parse(data.observedAt) > 180000);
  const unavailable = new Set(data?.unavailableSources ?? []);
  const backup = data?.backups.filter(
    item =>
      item.suspended ||
      !item.lastSuccess ||
      (item.name === 'backstage/backstage-backup' &&
        Date.now() - Date.parse(item.lastSuccess) > 129600000),
  ).length;
  const cards = [
    {
      title: 'Workloads not ready',
      value:
        unavailable.has('Deployments') || unavailable.has('StatefulSets')
          ? 'Unknown'
          : data?.unavailable.length,
      detail: data
        ? `of ${data.totalWorkloads} observed workloads`
        : 'Waiting for observation',
      color: '#237f91',
    },
    {
      title: 'Deployments to review',
      value: unavailable.has('Argo CD') ? 'Unknown' : data?.deployments.length,
      detail: 'Argo health or sync needs review',
      color: '#7858b5',
    },
    {
      title: 'Backup schedules to review',
      value: unavailable.has('Backup schedules') ? 'Unknown' : backup,
      detail: 'Activity only; not proof of restoration',
      color: '#a76d17',
    },
    {
      title: 'Service checks to review',
      value: data?.serviceHealth?.checks.filter(item => item.state !== 'up')
        .length,
      detail: 'Point-in-time checks; not uptime history',
      color: '#287756',
    },
  ];
  return (
    <Box
      mb={3}
      p={2}
      style={{ border: `1px solid ${theme.palette.divider}`, borderRadius: 20 }}
    >
      <Box
        display="flex"
        justifyContent="space-between"
        flexWrap="wrap"
        alignItems="center"
        style={{ gap: 8 }}
      >
        <Typography variant="h5" component="h2">
          Live overview
        </Typography>
        <Button color="primary" href="/today">
          Open Today
        </Button>
      </Box>
      <Typography variant="caption" display="block">
        {data
          ? `${stale ? 'Stale observation' : 'Observed'} · ${new Date(
              data.observedAt,
            ).toLocaleString()} · refreshes every minute`
          : 'Loading live sources…'}
      </Typography>
      {error && (
        <Typography role="alert" color="error">
          Live refresh unavailable; previous results may be stale.
        </Typography>
      )}
      {data?.partial && (
        <Typography role="status" color="error">
          Incomplete view: {data.unavailableSources?.join(', ')}. Missing data
          is unknown.
        </Typography>
      )}
      <Box
        mt={2}
        style={{
          display: 'grid',
          gridTemplateColumns: 'repeat(auto-fit,minmax(min(100%,190px),1fr))',
          gap: 12,
        }}
      >
        {cards.map(card => (
          <Box
            key={card.title}
            p={2}
            style={{
              borderTop: `3px solid ${card.color}`,
              borderRadius: 12,
              background: `${card.color}18`,
            }}
          >
            <Typography variant="subtitle2">{card.title}</Typography>
            <Typography variant="h4">{card.value ?? '—'}</Typography>
            <Typography variant="caption">{card.detail}</Typography>
          </Box>
        ))}
      </Box>
      <Box mt={2} display="flex" flexWrap="wrap" style={{ gap: 8 }}>
        <Button variant="outlined" href="/notifications">
          Notifications
        </Button>
        <Button
          variant="outlined"
          href="https://github.com/onemindos/omos/actions"
          target="_blank"
          rel="noopener noreferrer"
        >
          Builds and deployments
        </Button>
        <Button onClick={() => setInbox(value => !value)}>
          {inbox ? 'Hide inbox' : 'Show inbox'}
        </Button>
      </Box>
      {inbox && <ActivityInbox />}
    </Box>
  );
}
