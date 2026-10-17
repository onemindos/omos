import { createBackendModule, coreServices } from '@backstage/backend-plugin-api';

// The upstream Argo backend also exposes create-resource routes and actions.
// This portal uses it for observation; block mutations independently of UI.
export default createBackendModule({
  pluginId: 'backstage-community-argocd',
  moduleId: 'onemind-readonly',
  register(reg) {
    reg.registerInit({
      deps: { router: coreServices.httpRouter },
      async init({ router }) {
        router.use((req, res, next) => {
          if (!['GET', 'HEAD'].includes(req.method)) {
            res.status(403).json({ error: 'OneMind Argo integration is read-only' });
            return;
          }
          next();
        });
      },
    });
  },
});
