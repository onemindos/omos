import {
  coreServices,
  createBackendModule,
} from '@backstage/backend-plugin-api';
import { isOneMindAdmin } from './policy';

export const isReadOnlyOpenChoreoRequest = (method: string, path: string) =>
  ['GET', 'HEAD'].includes(method) ||
  (method === 'POST' && path === '/dashboard/bindings-count');

// The upstream REST proxy checks the IDP token, not the host permission policy.
// Enforce the host upper bound on actual mutations as well as permission queries.
export const openChoreoHostGuard = createBackendModule({
  pluginId: 'openchoreo',
  moduleId: 'onemind-host-boundary',
  register(reg) {
    reg.registerInit({
      deps: {
        router: coreServices.httpRouter,
        httpAuth: coreServices.httpAuth,
        userInfo: coreServices.userInfo,
      },
      async init({ router, httpAuth, userInfo }) {
        router.use(async (req, res, next) => {
          try {
            const credentials = await httpAuth.credentials(req, {
              allow: ['user'],
            });
            if (!isReadOnlyOpenChoreoRequest(req.method, req.path)) {
              const info = await userInfo.getUserInfo(credentials);
              if (!isOneMindAdmin({ info })) {
                res
                  .status(403)
                  .json({
                    error:
                      'OneMind administrator required for platform mutations',
                  });
                return;
              }
            }
            next();
          } catch (error) {
            next(error);
          }
        });
      },
    });
  },
});
