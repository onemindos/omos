/* eslint-disable */
import { createBackendPlugin, coreServices } from '@backstage/backend-plugin-api';
import express from 'express';
import https from 'node:https';
import { alertScopes, AlertScope, readPrometheusAlerts } from './prometheusAlerts';
import { readFileSync } from 'node:fs';
import { health, stacEvidence } from '@internal/plugin-onemind-common-backend';

type Application = { metadata: { name: string }; spec?: { destination?: { namespace?: string } }; status?: { reconciledAt?: string; health?: { status?: string }; sync?: { status?: string; revision?: string }; operationState?: { phase?: string; finishedAt?: string }; resources?: { health?: { status?: string } }[] } };
let cache: { observedAt: string; apps: Application[] } | undefined;
let pending: Promise<NonNullable<typeof cache>> | undefined;
async function applications() {
  if (cache && Date.now() - Date.parse(cache.observedAt) < 60000) return cache;
  if (!pending) pending = new Promise<NonNullable<typeof cache>>((resolve, reject) => {
    const root = '/var/run/secrets/kubernetes.io/serviceaccount/';
    const request = https.get({ hostname: 'kubernetes.default.svc', port: 443,
      path: '/apis/argoproj.io/v1alpha1/namespaces/argocd/applications',
      ca: readFileSync(`${root}ca.crt`), headers: { Authorization: `Bearer ${readFileSync(`${root}token`, 'utf8').trim()}` },
    }, response => {
      let body = ''; response.setEncoding('utf8');
      response.on('data', chunk => { body += chunk; if (body.length > 20000000) request.destroy(new Error('Response limit')); });
      response.on('error', reject);
      response.on('end', () => { if (response.statusCode !== 200) { reject(new Error('Argo observation unavailable')); return; }
        try { const apps = JSON.parse(body).items; if (!Array.isArray(apps)) throw new Error(); resolve({ observedAt: new Date().toISOString(), apps }); } catch { reject(new Error('Invalid Argo observation')); } });
    });
    const deadline = setTimeout(() => request.destroy(new Error('Observation deadline')), 8000);
    request.on('close', () => clearTimeout(deadline)); request.on('error', reject);
  }).then(value => (cache = value)).finally(() => { pending = undefined; });
  return pending;
}
export default createBackendPlugin({ pluginId: 'onemind-services', register(reg) {
  reg.registerInit({ deps: { router: coreServices.httpRouter, httpAuth: coreServices.httpAuth }, async init({ router, httpAuth }) {
    const routes = express.Router();
    routes.get('/overview', async (req, res, next) => { try {
      await httpAuth.credentials(req, { allow: ['user'] });
      const namespace = req.query.namespace;
      if (typeof namespace !== 'string' || !/^[a-z0-9](?:[a-z0-9-]{0,61}[a-z0-9])?$/.test(namespace)) { res.status(400).json({ error: 'A valid namespace is required' }); return; }
      const [argo, observations] = await Promise.allSettled([applications(), health()]);
      const apps = argo.status === 'fulfilled' ? argo.value.apps.filter(app => app.spec?.destination?.namespace === namespace).map(app => ({
        name: app.metadata.name, health: app.status?.health?.status ?? 'Unknown', sync: app.status?.sync?.status ?? 'Unknown',
        reconciledAt: typeof app.status?.reconciledAt === 'string' && Number.isFinite(Date.parse(app.status.reconciledAt)) ? app.status.reconciledAt : null,
        revision: app.status?.sync?.revision?.slice(0, 120), operation: app.status?.operationState?.phase ?? 'Unknown', finishedAt: app.status?.operationState?.finishedAt ?? null,
        observedResources: app.status?.resources?.length ?? 0,
        link: `https://argocd.tail717ba0.ts.net/applications/argocd/${encodeURIComponent(app.metadata.name)}`,
      })) : [];
      const metric = observations.status === 'fulfilled' ? observations.value.metrics.find(value => value.namespace === namespace) : undefined;
      res.json({ namespace, scope: 'Namespace-level associations, not exclusive service ownership', argo: { state: argo.status === 'fulfilled' ? 'observed' : 'unavailable', observedAt: argo.status === 'fulfilled' ? argo.value.observedAt : null, applications: apps }, metrics: metric ?? null,
        metricsState: observations.status === 'rejected' ? 'unavailable' : !metric ? 'not-configured' : metric.cpuCores === null && metric.memoryBytes === null ? 'no-data' : observations.value.metricsUnavailable || metric.cpuCores === null || metric.memoryBytes === null ? 'partial' : 'observed',
        grafana: { state: namespace === 'nats' ? 'linked' : 'dashboard-mapping-required', link: namespace === 'nats' ? 'https://grafana.tail717ba0.ts.net/d/omos-nats-observ' : 'https://grafana.tail717ba0.ts.net/dashboards' } });
    } catch (error) { next(error); } });
    routes.get('/stac', async (req,res,next)=>{try{await httpAuth.credentials(req,{allow:['user']});res.json(await stacEvidence());}catch(error){next(error);}});
    routes.get('/prometheus/:scope', async (req, res, next) => {
      try {
        await httpAuth.credentials(req, { allow: ['user'] });
        res.setHeader('Cache-Control', 'no-store');
        if (!alertScopes.includes(req.params.scope as AlertScope) || Object.keys(req.query).length) {
          res.status(400).json({ error: 'A supported alert scope without query parameters is required' }); return;
        }
        try { res.json(await readPrometheusAlerts(req.params.scope as AlertScope)); }
        catch { res.status(503).json({ error: 'Prometheus alert observation unavailable', observedAt: new Date().toISOString() }); }
      } catch (error) { next(error); }
    });
    router.use(routes);
  } });
} });
