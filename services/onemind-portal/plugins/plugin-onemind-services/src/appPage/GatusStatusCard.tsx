import { useEntity } from '@backstage/plugin-catalog-react';
import { InfoCard } from '@backstage/core-components';
import { Box, Typography } from '@material-ui/core';
import { InfoTooltip, SourceLine, StatusDot } from '@internal/plugin-onemind-common';
import { useGatusStatus } from './gatus';

const ago = (iso: string) => {
  const seconds = Math.max(0, Math.round((Date.now() - new Date(iso).getTime()) / 1000));
  if (seconds < 90) return `${seconds}s ago`;
  if (seconds < 5400) return `${Math.round(seconds / 60)} min ago`;
  return `${Math.round(seconds / 3600)} h ago`;
};

/**
 * Uptime check from Gatus for the entity's `onemind.io/health-check` endpoint: up or down, last check,
 * response time and the 14 day certificate condition. Never green without a fresh passing check.
 */
export function GatusStatusCard() {
  const { entity } = useEntity();
  const name = entity.metadata.annotations?.['onemind.io/health-check'] ?? '';
  const { loading, summary, readAt } = useGatusStatus(name);
  return (
    <InfoCard
      title="Uptime check"
      subheader={<SourceLine source={`Gatus · ${name}`} observedAt={readAt} />}
      action={
        <InfoTooltip label="uptime check">
          Gatus probes the app's health URL from inside the cluster, so it shows reachability, not that your account or
          data work. Check history lives in Gatus memory and resets when it restarts. The certificate line reports the
          14 day condition; the exact expiry date is in Prometheus (gatus_results_certificate_expiration_seconds).
        </InfoTooltip>
      }
    >
      {loading && <Typography role="status" variant="body2">Reading Gatus…</Typography>}
      {summary?.kind === 'unavailable' && (
        <Box display="flex" alignItems="center" style={{ gap: 8 }}>
          <StatusDot status="unknown" showLabel />
          <Typography role="alert" variant="body2" color="error">{summary.reason}</Typography>
        </Box>
      )}
      {summary?.kind === 'not-checked' && (
        <Box display="flex" alignItems="center" style={{ gap: 8 }}>
          <StatusDot status="unknown" showLabel />
          <Typography variant="body2">Not checked yet: Gatus has no result for “{name}”.</Typography>
        </Box>
      )}
      {summary?.kind === 'checked' && (
        <>
          <Box display="flex" alignItems="center" mb={1} style={{ gap: 8 }}>
            <StatusDot status={summary.status} showLabel size={12} />
            {summary.stale && <Typography variant="body2" color="textSecondary">last check is too old to trust</Typography>}
          </Box>
          <Typography variant="body2" data-testid="gatus-last-check">Last check: {ago(summary.lastCheck)}</Typography>
          <Typography variant="body2" data-testid="gatus-response">
            Response time: {summary.responseMs !== undefined ? `${summary.responseMs} ms` : 'Not recorded'}
            {summary.httpStatus ? ` · HTTP ${summary.httpStatus}` : ''}
          </Typography>
          <Typography variant="body2" data-testid="gatus-certificate">
            Certificate:{' '}
            {summary.certificate === 'valid' && 'more than 14 days left'}
            {summary.certificate === 'expiring' && 'under 14 days left, or could not be read'}
            {summary.certificate === 'not-checked' && 'not checked (plain HTTP inside the cluster)'}
          </Typography>
          {summary.status === 'down' && summary.failing.length > 0 && (
            <Typography variant="body2" color="error" data-testid="gatus-failing">Failed: {summary.failing.join('; ')}</Typography>
          )}
          {summary.status === 'down' && summary.errors.length > 0 && (
            <Typography variant="body2" color="textSecondary">{summary.errors[0]}</Typography>
          )}
        </>
      )}
    </InfoCard>
  );
}
