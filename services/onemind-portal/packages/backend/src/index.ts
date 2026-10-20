import { rootHttpRouterServiceFactory } from '@backstage/backend-defaults/rootHttpRouter';
import {
  immediateCatalogServiceFactory,
  annotationStoreFactory,
} from '@openchoreo/backstage-plugin-catalog-backend-module';
import { createIdpTokenHeaderMiddleware } from '@openchoreo/openchoreo-auth';
/*
 * OneMind backend wiring.
 * All custom modules are loaded from @internal plugin packages.
 */

import { createBackend } from '@backstage/backend-defaults';

const backend = createBackend();
backend.add(immediateCatalogServiceFactory);
backend.add(annotationStoreFactory);
backend.add(
  rootHttpRouterServiceFactory({
    configure: ({ app, applyDefaults }) => {
      // 1.55.3 applyDefaults installs the security/logging stack once.
      app.use(createIdpTokenHeaderMiddleware());
      applyDefaults();
    },
  }),
);
backend.add(import('@openchoreo/backstage-plugin-catalog-backend-module'));
backend.add(import('@openchoreo/backstage-plugin-backend'));
backend.add(import('@openchoreo/backstage-plugin-scaffolder-backend-module'));
backend.add(
  import('@openchoreo/backstage-plugin-openchoreo-observability-backend'),
);
// OneMind remains the sole permission policy and the sole sign-in resolver.

backend.add(import('@backstage/plugin-app-backend'));
backend.add(import('@backstage/plugin-proxy-backend'));

backend.add(import('@backstage/plugin-scaffolder-backend'));
backend.add(import('@backstage/plugin-scaffolder-backend-module-github'));
backend.add(
  import('@backstage/plugin-scaffolder-backend-module-notifications'),
);

backend.add(import('@backstage/plugin-techdocs-backend'));

backend.add(import('@backstage/plugin-auth-backend'));
backend.add(import('@internal/plugin-onemind-auth-backend'));

backend.add(import('@backstage/plugin-catalog-backend'));
backend.add(
  import('@backstage/plugin-catalog-backend-module-scaffolder-entity-model'),
);
backend.add(import('@backstage/plugin-catalog-backend-module-logs'));
backend.add(import('@internal/plugin-onemind-machine-backend'));

backend.add(import('@backstage/plugin-permission-backend'));
backend.add(import('@internal/plugin-onemind-permissions-backend'));
backend.add(
  import('@internal/plugin-onemind-permissions-backend').then(m => ({
    default: m.openChoreoHostGuard,
  })),
);

backend.add(import('@backstage/plugin-search-backend'));
backend.add(import('@backstage/plugin-search-backend-module-pg'));
backend.add(import('@backstage/plugin-search-backend-module-catalog'));
backend.add(import('@backstage/plugin-search-backend-module-techdocs'));

backend.add(import('@backstage/plugin-kubernetes-backend'));
backend.add(import('@backstage/plugin-user-settings-backend'));
backend.add(import('@backstage/plugin-events-backend'));
backend.add(import('@backstage/plugin-notifications-backend'));
backend.add(import('@backstage/plugin-signals-backend'));
backend.add(import('@backstage/plugin-mcp-actions-backend'));

backend.add(import('@internal/plugin-onemind-operations-backend'));
backend.add(import('@internal/plugin-onemind-missions-backend'));
backend.add(import('@internal/plugin-onemind-events-backend'));
backend.add(import('@internal/plugin-onemind-activity-backend'));
backend.add(import('@internal/plugin-onemind-services-backend'));
backend.add(
  import('@internal/plugin-onemind-services-backend').then(m => ({
    default: m.stacCatalogModule,
  })),
);

backend.add(import('@internal/plugin-onemind-context-backend'));
backend.add(import('@internal/plugin-onemind-assistant-backend'));
backend.add(import('@internal/plugin-onemind-offline-backend'));

backend.add(import('@backstage/plugin-auth-backend-module-github-provider'));
backend.add(import('@backstage-community/plugin-argocd-backend'));
backend.add(import('@internal/plugin-onemind-community-backend'));
backend.add(import('@backstage-community/plugin-tech-insights-backend'));
backend.add(
  import('@backstage-community/plugin-tech-insights-backend-module-jsonfc'),
);
backend.add(import('@backstage-community/plugin-adr-backend'));
backend.add(import('@backstage-community/search-backend-module-adr'));
backend.add(import('@backstage-community/plugin-announcements-backend'));
backend.add(
  import('@backstage-community/plugin-search-backend-module-announcements'),
);

backend.start();
