import { InfoCard } from '@backstage/core-components';
import { Box, List, ListItem, Typography } from '@material-ui/core';
import { InfoTooltip, SourceLine, StatusDot, useEntityStatus } from '@internal/plugin-onemind-common';

const NAME: Record<string, string> = { alerts: 'Alerts', argo: 'Argo CD', workloads: 'Replicas ready', uptime: 'Uptime check' };

/** Overview summary: the signals behind the dot. One tooltip carries the caveat; no repeated disclaimers. */
export function StatusSummaryCard() {
  const { status, signals, observedAt, loading, error } = useEntityStatus();
  return (
    <InfoCard
      title="Live status"
      subheader={<SourceLine source="Argo CD, Kubernetes, uptime checks, scoped alerts" observedAt={observedAt} />}
      action={
        <InfoTooltip label="live status">
          Combines the signals that are connected for this app. Grey means unknown, never healthy: a signal that
          cannot be read, or an observation older than three minutes, is not counted as good. A connected signal that
          shows no problem is not proof the app works end to end.
        </InfoTooltip>
      }
    >
      <Box display="flex" alignItems="center" mb={1}>
        <StatusDot status={status} showLabel size={12} />
      </Box>
      {loading && <Typography role="status" variant="body2">Reading live signals…</Typography>}
      {error && !loading && (
        <Typography role="alert" variant="body2" color="error">Live update unavailable; showing the last reading.</Typography>
      )}
      {!loading && !signals.length && (
        <Typography variant="body2">No live signal is connected for this app yet.</Typography>
      )}
      <List dense disablePadding>
        {signals.map(signal => (
          <ListItem key={signal.name} disableGutters style={{ alignItems: 'flex-start', gap: 8 }}>
            <Box pt={0.5}><StatusDot status={signal.status} /></Box>
            <Box>
              <Typography variant="subtitle2" component="span">{NAME[signal.name]}</Typography>
              <Typography variant="body2" color="textSecondary">{signal.detail}</Typography>
            </Box>
          </ListItem>
        ))}
      </List>
    </InfoCard>
  );
}
