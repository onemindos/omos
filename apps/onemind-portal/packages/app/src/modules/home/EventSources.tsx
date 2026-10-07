import { useEffect, useState } from 'react';
import { discoveryApiRef, fetchApiRef, useApi } from '@backstage/core-plugin-api';
import { Box, Typography } from '@material-ui/core';
type SourceState = { enabled: boolean; recipientReadiness: string; lastAttemptAt: string | null; lastSuccessAt: string | null; sources: Record<string, string> };
export function EventSources() {
  const discovery = useApi(discoveryApiRef), api = useApi(fetchApiRef);
  const [data, setData] = useState<SourceState>();
  const [error, setError] = useState(false);
  useEffect(() => {
    let active = true; const controller = new AbortController();
    const refresh = async () => { try {
      const base = await discovery.getBaseUrl('onemind');
      const response = await api.fetch(`${base}/event-sources`, { signal: controller.signal });
      if (!response.ok) throw new Error();
      const value = await response.json(); if (active) { setData(value); setError(false); }
    } catch { if (active) setError(true); } };
    refresh(); const timer = setInterval(refresh, 60000);
    return () => { active = false; controller.abort(); clearInterval(timer); };
  }, [discovery, api]);
  const stale = !!data?.lastSuccessAt && Date.now() - Date.parse(data.lastSuccessAt) > 180000;
  return <Box p={2} borderTop={1} borderColor="divider">
    <Typography variant="subtitle2">Infrastructure event sources</Typography>
    {error && <Typography color="error">Source status unavailable; previous observations may be stale.</Typography>}
    {!data && !error && <Typography variant="body2">Checking source connections…</Typography>}
    {data && <>
      <Typography variant="body2">{data.enabled ? 'Read-only observer enabled' : 'Observer not configured'} · {data.recipientReadiness}</Typography>
      {Object.entries(data.sources).map(([name, value]) => <Typography variant="caption" display="block" key={name}>{name}: {value}</Typography>)}
      {data.enabled && !data.lastSuccessAt && <Typography variant="caption">Waiting for the first successful observation.</Typography>}
      {data.lastSuccessAt && <Typography variant="caption" display="block" color={stale ? 'error' : 'textSecondary'}>{stale ? 'Stale · ' : 'Last checked · '}{new Date(data.lastSuccessAt).toLocaleString()}</Typography>}
      <Typography variant="caption" display="block">Only new terminal Argo operations and scheduled Backstage backup jobs are observed. Historical results are baselined. Agent events are not connected. Recipient membership is checked; accepted processing is not delivery confirmation.</Typography>
    </>}
  </Box>;
}
