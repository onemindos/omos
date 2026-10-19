import { useEffect, useState } from 'react';
import {
  discoveryApiRef,
  fetchApiRef,
  useApi,
} from '@backstage/core-plugin-api';
import { Box, Link, Typography } from '@material-ui/core';

type Costs = {
  state: string;
  scope: string;
  currency: string;
  observedAt: string | null;
  generatedAt: string | null;
  expiresAt: string | null;
  monthToDateUsage?: string;
  monthToDateBalance?: string;
  accountBalance?: string;
};
export function CostPanel() {
  const discovery = useApi(discoveryApiRef);
    const api = useApi(fetchApiRef);
    const [data, setData] = useState<Costs>();
    const [error, setError] = useState('');
  useEffect(() => {
    let active = true;
      let pending = false;
    const controller = new AbortController();
    const refresh = async () => {
      if (pending) return;
      pending = true;
      try {
        const base = await discovery.getBaseUrl('onemind');
        const response = await api.fetch(`${base  }/costs`, {
          signal: controller.signal,
        });
        if (!response.ok)
          throw new Error(
            response.status === 403
              ? 'Administrator access required'
              : 'Billing observation unavailable',
          );
        const value = await response.json();
        if (active) {
          setData(value);
          setError('');
        }
      } catch (e) {
        if (active)
          setError(e instanceof Error ? e.message : 'Billing unavailable');
      } finally {
        pending = false;
      }
    };
    refresh();
    const timer = setInterval(refresh, 900000);
    return () => {
      active = false;
      controller.abort();
      clearInterval(timer);
    };
  }, [discovery, api]);
  const stale =
    !!data &&
    (!!error ||
      data.state === 'stale' ||
      (!!data.expiresAt && Date.parse(data.expiresAt) < Date.now()));
  return (
    <Box p={2}>
      <Typography variant="h5">DigitalOcean account costs</Typography>
      <Typography variant="body2" color="textSecondary">
        Entire connected account · USD · no workload allocation
      </Typography>
      {error && (
        <Typography color="error" role="alert">
          {error}
        </Typography>
      )}
      {!data && !error && <Typography>Loading billing…</Typography>}
      {data && (
        <>
          <Typography variant="caption" display="block">
            {data.scope}
          </Typography>
          {data.state === 'configuration-required' && (
            <Typography>
              A dedicated billing:read token is required. Account balances are
              not connected yet.
            </Typography>
          )}
          {data.state === 'unavailable' && (
            <Typography color="error">
              Billing source unavailable; no totals assumed.
            </Typography>
          )}
          {data.observedAt && (
            <>
              <Typography>
                Month-to-date usage: {data.monthToDateUsage ?? 'unknown'} USD
              </Typography>
              <Typography>
                Month-to-date balance: {data.monthToDateBalance ?? 'unknown'}{' '}
                USD
              </Typography>
              <Typography>
                Account balance: {data.accountBalance ?? 'unknown'} USD
              </Typography>
              <Typography variant="caption" display="block">
                {stale ? 'Stale · ' : ''}Fetched{' '}
                {new Date(data.observedAt).toLocaleString()} · source generated{' '}
                {data.generatedAt
                  ? new Date(data.generatedAt).toLocaleString()
                  : 'unknown'}{' '}
                · cached 15 minutes
              </Typography>
            </>
          )}
        </>
      )}
      <Box mt={2}>
        <Link
          href="https://cloud.digitalocean.com/account/billing"
          target="_blank"
          rel="noopener noreferrer"
        >
          Open DigitalOcean billing
        </Link>
      </Box>
    </Box>
  );
}
