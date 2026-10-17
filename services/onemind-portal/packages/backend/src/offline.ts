import { createBackendPlugin, coreServices } from '@backstage/backend-plugin-api';
import express from 'express';
import { readFileSync } from 'node:fs';
import directory from '../../app/src/modules/home/directory.json';
import { makePack, names } from './offlinePack';

export default createBackendPlugin({ pluginId: 'onemind-offline', register(reg) {
  reg.registerInit({ deps: { router: coreServices.httpRouter, httpAuth: coreServices.httpAuth,
    userInfo: coreServices.userInfo }, async init({ router, httpAuth, userInfo }) {
    const routes = express.Router();
    routes.get('/pack', async (req, res, next) => {
      try {
        const credentials = await httpAuth.credentials(req, { allow: ['user'] });
        const user = await userInfo.getUserInfo(credentials);
        if (!user.ownershipEntityRefs.includes('group:default/onemind-admins')) { res.status(403).json({ error: 'Administrator access required' }); return; }
        const documents = names.map(name => ({ name, content: readFileSync(`/app/runbooks/docs/${name}.md`, 'utf8') }));
        const content = makePack(directory, documents, new Date().toISOString());
        res.setHeader('Cache-Control', 'no-store');
        res.setHeader('Content-Disposition', 'attachment; filename="onemind-offline.html"');
        res.type('html').send(content);
      } catch (error) { next(error); }
    });
    router.use(routes);
  } });
} });
