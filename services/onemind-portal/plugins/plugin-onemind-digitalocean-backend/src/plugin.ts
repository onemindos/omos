import {
  coreServices,
  createBackendPlugin,
} from '@backstage/backend-plugin-api';
import { createRouter } from './router';
import { sharedService } from './service';

/** /api/onemind-digitalocean/summary: read-only DigitalOcean account summary for administrators. */
export const digitalOceanPlugin = createBackendPlugin({
  pluginId: 'onemind-digitalocean',
  register(reg) {
    reg.registerInit({
      deps: {
        router: coreServices.httpRouter,
        httpAuth: coreServices.httpAuth,
        userInfo: coreServices.userInfo,
        logger: coreServices.logger,
      },
      async init({ router, httpAuth, userInfo, logger }) {
        router.use(
          createRouter({
            httpAuth,
            userInfo,
            service: sharedService(m => logger.warn(m)),
          }),
        );
      },
    });
  },
});
