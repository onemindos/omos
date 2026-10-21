import { useEffect, useState } from 'react';
import {
  discoveryApiRef,
  fetchApiRef,
  useApi,
} from '@backstage/core-plugin-api';
import { Box, Button, Typography } from '@material-ui/core';

type Summary = {
  state: string;
  scope: string;
  observedAt: string | null;
  connections: number | null;
  slowConsumers: number | null;
  messagesIn: number | null;
  messagesOut: number | null;
  memoryBytes: number | null;
  jetstream?: {
    state: string;
    streams: number | null;
    consumers: number | null;
    messages: number | null;
    storageBytes: number | null;
  };
  agentEvents: string;
};
export function FabricPanel() {
  const discovery = useApi(discoveryApiRef);
    const api = useApi(fetchApiRef);
  const [data, setData] = useState<Summary>();
    const [error, setError] = useState(false);
  useEffect(() => {
    let active = true;
      let pending = false;
    const controller = new AbortController();
    const refresh = async () => {
      if (pending) return;
      pending = true;
      try {
        const base = await discovery.getBaseUrl('onemind');
        const response = await api.fetch(`${base  }/fabric`, {
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
    !!data?.observedAt &&
    (error || Date.now() - Date.parse(data.observedAt) > 180000);
  return (
    <Box p={2}>
      <Typography variant="h5">Agent fabric · NATS</Typography>
      {error && (
        <Typography color="error">
          Observation unavailable; previous results may be stale.
        </Typography>
      )}
      {!data && !error && <Typography>Loading monitoring…</Typography>}
      {data && (
        <>
          <Typography variant="caption" display="block">
            {data.scope}
          </Typography>
          <Typography variant="caption" display="block">
            {data.observedAt
              ? `${stale ? 'Stale' : 'Observed'} · ${new Date(
                  data.observedAt,
                ).toLocaleString()}`
              : 'Monitoring unavailable'}{' '}
            · refreshes every minute
          </Typography>
          {data.state === 'observed' && (
            <Box my={2}>
              <Typography>
                Connections: {data.connections ?? 'unknown'} · Slow consumers:{' '}
                {data.slowConsumers ?? 'unknown'}
              </Typography>
              <Typography>
                Messages in/out: {data.messagesIn ?? 'unknown'} /{' '}
                {data.messagesOut ?? 'unknown'} · cumulative server counters
              </Typography>
              <Typography>
                JetStream: {data.jetstream?.state ?? 'unavailable'} · streams{' '}
                {data.jetstream?.streams ?? 'unknown'} · consumers{' '}
                {data.jetstream?.consumers ?? 'unknown'} · stored messages{' '}
                {data.jetstream?.messages ?? 'unknown'}
              </Typography>
            </Box>
          )}
          <Typography variant="body2" color="textSecondary">
            Agent completion events are not connected. Stream traffic does not
            prove a task completed. A producer contract and restricted consumer
            identity are required before enabling that feed.
          </Typography>
        </>
      )}
      <Box mt={2} display="flex" flexWrap="wrap" style={{ gap: 8 }}>
        <Button
          href="https://nats-dashboard-1.tail717ba0.ts.net/"
          target="_blank"
          rel="noopener noreferrer"
        >
          NATS dashboard
        </Button>
        <Button
          href="https://agent-web-ui-1.tail717ba0.ts.net/"
          target="_blank"
          rel="noopener noreferrer"
        >
          Agent workspace
        </Button>
        <Button href="/docs/default/component/onemind-runbooks/nats/">
          Runbook
        </Button>
      </Box>
    </Box>
  );
}
