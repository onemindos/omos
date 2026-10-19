/* eslint-disable */
import { createBackendPlugin, coreServices } from '@backstage/backend-plugin-api';
import { readFileSync } from 'node:fs';
import express from 'express';

const core = 'http://onemind-core.onemind-core.svc.cluster.local:8080';
const fga = 'http://openfga.onemind-core.svc.cluster.local:8080';
const reads = new Set(['missions', 'objects', 'links', 'agents', 'capabilities']);

async function check(user: string, groups: string[], capability: string) {
  const store = readFileSync('/etc/onemind-fga/store-id', 'utf8').trim();
  const model = readFileSync('/etc/onemind-fga/model-id', 'utf8').trim();
  if (!store || !model || !process.env.ONEMIND_OPENFGA_TOKEN) throw new Error('Authorization is initializing');
  const response = await fetch(`${fga}/stores/${encodeURIComponent(store)}/check`, {
    method: 'POST', signal: AbortSignal.timeout(8000),
    headers: { 'Content-Type': 'application/json', Authorization: `Bearer ${process.env.ONEMIND_OPENFGA_TOKEN}` },
    body: JSON.stringify({ authorization_model_id: model,
      tuple_key: { user, relation: 'runner', object: `capability:${capability}` },
      contextual_tuples: { tuple_keys: groups.filter(group => group.startsWith('group:')).map(group => ({ user, relation: 'member', object: group })) },
    }),
  });
  if (!response.ok) throw new Error('Authorization service unavailable');
  return (await response.json()).allowed === true;
}

export default createBackendPlugin({ pluginId: 'onemind-missions', register(reg) {
  reg.registerInit({ deps: { router: coreServices.httpRouter, httpAuth: coreServices.httpAuth, userInfo: coreServices.userInfo },
    async init({ router, httpAuth, userInfo }) {
      const routes = express.Router();
      routes.use(express.json({ limit: '80kb' }));
      async function identity(req: express.Request) {
        const credentials = await httpAuth.credentials(req, { allow: ['user'] });
        return userInfo.getUserInfo(credentials);
      }
      async function forward(req: express.Request, res: express.Response, path: string, method = 'GET', capability?: string) {
        const user = await identity(req);
        if (capability) {
          if (!user.ownershipEntityRefs.includes('group:default/onemind-admins')) { res.status(403).json({ error: 'An operator role is required' }); return; }
          if (!await check(user.userEntityRef, user.ownershipEntityRefs, capability)) { res.status(403).json({ error: 'This capability is not authorized' }); return; }
        }
        const key = process.env.ONEMIND_CORE_GATEWAY_KEY;
        if (!key) { res.status(503).json({ error: 'Mission service is initializing' }); return; }
        const upstream = await fetch(core + path, { method, signal: AbortSignal.timeout(35000),
          headers: { 'Content-Type': 'application/json', Authorization: `Bearer ${key}`, 'X-OneMind-Actor': user.userEntityRef },
          ...(method === 'POST' ? { body: JSON.stringify(req.body ?? {}) } : {}),
        });
        res.status(upstream.status).json(await upstream.json());
      }
      routes.get('/viewer', async (req, res, next) => { try {
        const user = await identity(req);
        res.json({ actor: user.userEntityRef, canOperate: user.ownershipEntityRefs.includes('group:default/onemind-admins') });
      } catch (error) { next(error); } });
      routes.get('/:collection', async (req, res, next) => { try {
        if (!reads.has(req.params.collection)) { res.sendStatus(404); return; }
        await forward(req, res, `/${req.params.collection}${req.params.collection === 'objects' ? '?limit=500' : ''}`);
      } catch (error) { next(error); } });
      routes.post('/missions', async (req, res, next) => { try {
        if (typeof req.body?.capability !== 'string' || !['inspect-system', 'summarize-site', 'delegate-agent', 'restart-worker'].includes(req.body.capability)) {
          res.status(400).json({ error: 'Choose a supported capability' }); return;
        }
        await forward(req, res, '/missions', 'POST', req.body.capability);
      } catch (error) { next(error); } });
      routes.post('/missions/:id/cancel', async (req, res, next) => { try {
        if (!/^[a-f0-9-]{36}$/.test(req.params.id)) { res.sendStatus(400); return; }
        await forward(req, res, `/missions/${encodeURIComponent(req.params.id)}/cancel`, 'POST', 'cancel-mission');
      } catch (error) { next(error); } });
      for (const collection of ['objects', 'links']) routes.post(`/${collection}`, async (req, res, next) => { try {
        await forward(req, res, `/${collection}`, 'POST', 'edit-world');
      } catch (error) { next(error); } });
      router.use(routes);
    },
  });
} });
