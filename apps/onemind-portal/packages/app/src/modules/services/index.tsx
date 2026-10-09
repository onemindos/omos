import { useEffect, useState } from 'react';
import { createFrontendModule } from '@backstage/frontend-plugin-api';
import { EntityContentBlueprint } from '@backstage/plugin-catalog-react/alpha';
import { useEntity } from '@backstage/plugin-catalog-react';
import { discoveryApiRef, fetchApiRef, useApi } from '@backstage/core-plugin-api';
import { AccessButton } from '../home/AccessDeployments';
import { parseAccessPoints } from '../home/access';
import { safeUrl } from '../home/useDirectory';
import { Box, Link, Typography } from '@material-ui/core';
type Overview = { scope: string; argo: { state: string; observedAt: string | null; applications: { name: string; health: string; sync: string; revision?: string; reconciledAt: string | null; operation: string; finishedAt: string | null; observedResources: number; link: string }[] }; metricsState: string; metrics: { cpuCores: number | null; memoryBytes: number | null; observedAt: string; cpuSampleAt: string | null; memorySampleAt: string | null } | null; grafana: { state: string; link: string } };
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
      {data.argo.observedAt && <Typography variant="caption">Fetched from Kubernetes {new Date(data.argo.observedAt).toLocaleString()} · cached for one minute · fetch time does not prove controller freshness</Typography>}
      {data.argo.state === 'observed' && !data.argo.applications.length && <Typography>No Application destination matches this namespace. This does not prove the service is unmanaged.</Typography>}
      {data.argo.applications.map(app => <Box key={app.name} my={2}><Link href={app.link} target="_blank" rel="noopener noreferrer">{app.name}</Link><Typography>{app.health} · {app.sync} · last operation {app.operation}</Typography>
        <Typography variant="body2" color={!app.reconciledAt || Date.now() - Date.parse(app.reconciledAt) > 600000 ? 'error' : 'textSecondary'}>{!app.reconciledAt ? 'Controller reconciliation time unknown; freshness cannot be established' : `${Date.now() - Date.parse(app.reconciledAt) > 600000 ? 'Stale controller observation' : 'Controller reconciled'} ${new Date(app.reconciledAt).toLocaleString()}`}</Typography><Typography variant="body2">Revision {app.revision ?? 'unknown'} · {app.observedResources} tracked resources</Typography>{app.finishedAt && <Typography variant="caption">Operation finished {new Date(app.finishedAt).toLocaleString()}</Typography>}</Box>)}
      <Typography variant="h5">Metrics and Grafana</Typography>
      <Typography>Metrics: {data.metricsState}. Values describe the namespace, not this individual service.</Typography>
      {data.metrics && <><Typography>CPU {data.metrics.cpuCores === null ? 'unknown' : `${data.metrics.cpuCores.toFixed(2)} cores`} · memory {data.metrics.memoryBytes === null ? 'unknown' : `${(data.metrics.memoryBytes / 1048576).toFixed(0)} MiB`}</Typography><Typography variant="caption">Fetched metrics {new Date(data.metrics.observedAt).toLocaleString()} · CPU average over five minutes · sample times: CPU {data.metrics.cpuSampleAt ?? 'unknown'}, memory {data.metrics.memorySampleAt ?? 'unknown'}</Typography></>}
      {data.grafana.state === 'dashboard-mapping-required' && <Typography variant="body2">A service-specific Grafana dashboard has not been mapped yet.</Typography>}
      <Box mt={2}><Link href={data.grafana.link} target="_blank" rel="noopener noreferrer">Open Grafana</Link> · <Link href="/docs/default/component/onemind-runbooks/monitoring/">Monitoring runbook</Link></Box>
    </>}
  </Box>;
}
export const servicesModule = createFrontendModule({ pluginId: 'catalog', extensions: [EntityContentBlueprint.make({ name: 'onemind-operations', params: { path: '/operations', title: 'Operations', filter: entity => !!entity.metadata.annotations?.['backstage.io/kubernetes-namespace'], loader: async () => <ServiceOperations /> } })] });

function AppAccessPanel() {
  const {entity} = useEntity();
  const a=entity.metadata.annotations ?? {};
  const raw=a['onemind.io/launcher-url'] ?? entity.metadata.links?.find(x=>x.title==='Open application')?.url;
  const url=raw && safeUrl(raw);
  if(!url) return <Typography>No safe access link recorded.</Typography>;
  const app={name:entity.metadata.name,title:entity.metadata.title ?? entity.metadata.name,description:entity.metadata.description ?? '',url,ref:`component:${entity.metadata.namespace ?? 'default'}/${entity.metadata.name}`,group:'services',groupTitle:a['onemind.io/launcher-category'] ?? 'Services',hosting:a['onemind.io/hosting'],status:a['onemind.io/documented-status'],accessPoints:parseAccessPoints(a['onemind.io/access-points']),supportedPlatforms:(a['onemind.io/supported-platforms'] ?? '').split(',').filter(Boolean),dependencies:entity.relations?.filter(x=>x.type==='dependsOn').map(x=>x.targetRef) ?? []};
  return <Box p={3}><Typography variant="h4">Access & Deployments</Typography><Typography>Review recorded installations, platforms, device relationships, and your preferred way to open {app.title}.</Typography><AccessButton app={app}/></Box>;
}
export const accessModule = createFrontendModule({pluginId:'catalog',extensions:[EntityContentBlueprint.make({name:'onemind-access',params:{path:'/access',title:'Access & Deployments',filter:entity=>entity.kind.toLowerCase()==='component' && !!(entity.metadata.annotations?.['onemind.io/launcher-url'] || entity.metadata.links?.some(x=>x.title==='Open application')),loader:async()=> <AppAccessPanel/>}})]});
