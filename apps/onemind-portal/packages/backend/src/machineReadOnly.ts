import { createBackendModule, coreServices } from '@backstage/backend-plugin-api';

// Backstage cannot delegate permission-restricted external service credentials.
// Keep plugin-level delegation, and deny direct machine catalog mutations here.
export default createBackendModule({
  pluginId: 'catalog', moduleId: 'onemind-machine-readonly',
  register(reg) {
    reg.registerInit({
      deps: { router: coreServices.httpRouter, httpAuth: coreServices.httpAuth, auth: coreServices.auth },
      async init({ router, httpAuth, auth }) {
        router.use(async (req, res, next) => {
          try {
            const credentials = await httpAuth.credentials(req);
            if (auth.isPrincipal(credentials, 'service') && credentials.principal.subject === 'onemind-readonly-ai' && !['GET', 'HEAD'].includes(req.method)) {
              res.status(403).json({ error: 'Machine identity has read-only catalog access' });
              return;
            }
            next();
          } catch (error) { next(error); }
        });
      },
    });
  },
});
