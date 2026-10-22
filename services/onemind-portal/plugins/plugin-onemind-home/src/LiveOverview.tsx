import { useState } from 'react';
import { Box, Button, Typography, useTheme } from '@material-ui/core';
import { alpha } from '@material-ui/core/styles';
import { ActivityInbox } from './ActivityInbox';
import { backupNeedsReview, isStale, useOperationsStatus } from './status';

export function LiveOverview() {
  const theme = useTheme();
  const { data, error } = useOperationsStatus();
  const [inbox, setInbox] = useState(false);
  const stale = !!data && (error || isStale(data));
  const unavailable = new Set(data?.unavailableSources ?? []);
  const backup = data?.backups.filter(item => backupNeedsReview(item)).length;
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
      color: theme.palette.primary.main,
    },
    {
      title: 'Deployments to review',
      value: unavailable.has('Argo CD') ? 'Unknown' : data?.deployments.length,
      detail: 'Argo health or sync needs review',
      color: theme.palette.secondary.main,
    },
    {
      title: 'Backup schedules to review',
      value: unavailable.has('Backup schedules') ? 'Unknown' : backup,
      detail: 'Activity only; not proof of restoration',
      color: theme.palette.warning.main,
    },
    {
      title: 'Service checks to review',
      value: data?.serviceHealth?.checks.filter(item => item.state !== 'up')
        .length,
      detail: 'Point-in-time checks; not uptime history',
      color: theme.palette.success.main,
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
              background: alpha(card.color, 0.09),
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
