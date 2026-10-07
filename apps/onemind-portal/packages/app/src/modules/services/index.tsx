import { useEffect, useState } from 'react';
import { createFrontendModule } from '@backstage/frontend-plugin-api';
import { EntityContentBlueprint } from '@backstage/plugin-catalog-react/alpha';
import { useEntity } from '@backstage/plugin-catalog-react';
import { discoveryApiRef, fetchApiRef, useApi } from '@backstage/core-plugin-api';
import { Box, Link, Typography } from '@material-ui/core';
type Overview = { scope: string; argo: { state: string; observedAt: string | null; applications: { name: string; health: string; sync: string; revision?: string; operation: string; finishedAt: string | null; observedResources: number; link: string }[] }; metricsState: string; metrics: { cpuCores: number | null; memoryBytes: number | null; observedAt: string; cpuSampleAt: string | null; memorySampleAt: string | null } | null; grafana: { state: string; link: string } };
function ServiceOperations() {
  const { entity } = useEntity(); const namespace = entity.metadata.annotations?.['backstage.io/kubernetes-namespace'];
  const discovery = useApi(discoveryApiRef), api = useApi(fetchApiRef);
  const [data, setData] = useState<Overview>(); const [error, setError] = useState(false);
  useEffect(() => { let active = true; const controller = new AbortController(); setData(undefined);
    const load = async () => { try { const base = await discovery.getBaseUrl('onemind-services'); const response = await api.fetch(`${base}/overview?namespace=${encodeURIComponent(namespace ?? '')}`, { signal: controller.signal }); if (!response.ok) throw new Error(); const value = await response.json(); if (active) { setData(value); setError(false); } } catch { if (active) setError(true); } };
    load(); const timer = setInterval(load, 60000); return () => { active = false; controller.abort(); clearInterval(timer); };
  }, [namespace, discovery, api]);
  return <Box p={3}><Typography variant="h4">Operations · {namespace}</Typography>
    {error && <Typography color="error">Observation unavailable; previous results may be stale.</Typography>}
    {!data && !error && <Typography>Loading observations…</Typography>}
    {data && <><Typography>{data.scope}. A shared namespace may contain multiple services and deployments.</Typography>
      <Typography variant="h5">Argo CD deployments</Typography>
      {data.argo.state === 'unavailable' && <Typography color="error">Argo observation unavailable. Missing results do not mean healthy.</Typography>}
      {data.argo.observedAt && <Typography variant="caption">Observed {new Date(data.argo.observedAt).toLocaleString()} · cached for one minute</Typography>}
      {data.argo.state === 'observed' && !data.argo.applications.length && <Typography>No Application destination matches this namespace. This does not prove the service is unmanaged.</Typography>}
      {data.argo.applications.map(app => <Box key={app.name} my={2}><Link href={app.link} target="_blank" rel="noopener noreferrer">{app.name}</Link><Typography>{app.health} · {app.sync} · last operation {app.operation}</Typography><Typography variant="body2">Revision {app.revision ?? 'unknown'} · {app.observedResources} tracked resources</Typography>{app.finishedAt && <Typography variant="caption">Operation finished {new Date(app.finishedAt).toLocaleString()}</Typography>}</Box>)}
      <Typography variant="h5">Metrics and Grafana</Typography>
      <Typography>Metrics: {data.metricsState}. Values describe the namespace, not this individual service.</Typography>
      {data.metrics && <><Typography>CPU {data.metrics.cpuCores === null ? 'unknown' : `${data.metrics.cpuCores.toFixed(2)} cores`} · memory {data.metrics.memoryBytes === null ? 'unknown' : `${(data.metrics.memoryBytes / 1048576).toFixed(0)} MiB`}</Typography><Typography variant="caption">Observed {new Date(data.metrics.observedAt).toLocaleString()} · CPU average over five minutes</Typography></>}
      {data.grafana.state === 'dashboard-mapping-required' && <Typography variant="body2">A service-specific Grafana dashboard has not been mapped yet.</Typography>}
      <Box mt={2}><Link href={data.grafana.link} target="_blank" rel="noopener noreferrer">Open Grafana</Link> · <Link href="/docs/default/component/onemind-runbooks/monitoring/">Monitoring runbook</Link></Box>
    </>}
  </Box>;
}
export const servicesModule = createFrontendModule({ pluginId: 'catalog', extensions: [EntityContentBlueprint.make({ name: 'onemind-operations', params: { path: '/operations', title: 'Operations', filter: entity => !!entity.metadata.annotations?.['backstage.io/kubernetes-namespace'], loader: async () => <ServiceOperations /> } })] });
