/*
 * OneMind backend wiring.
 * All custom modules are loaded from @internal plugin packages.
 */

import { createBackend } from '@backstage/backend-defaults';

const backend = createBackend();
// OneMind is the sole permission policy and the sole sign-in resolver.

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
// Failed and pending catalog entities. Administrator-only through the OneMind policy.
backend.add(import('@backstage/plugin-catalog-backend-module-unprocessed'));
// GitHub discovery through the read-only GitHub App; discovers nothing without catalog.providers.github.
// Also lets Catalog Import analyze GitHub repository URLs (catalog.location.create stays admin-only).
backend.add(import('@backstage/plugin-catalog-backend-module-github'));
backend.add(import('@internal/plugin-onemind-machine-backend'));
// Live cluster inventory. Disabled unless onemind.inventory.live.enabled is true; see docs/catalog/live-inventory.md.
backend.add(import('@internal/plugin-onemind-inventory-backend'));
// F6 DigitalOcean: read-only summary route plus catalog Resources. Inert without DIGITALOCEAN_TOKEN; see docs/operations/digitalocean.md.
backend.add(import('@internal/plugin-onemind-digitalocean-backend'));
backend.add(
  import('@internal/plugin-onemind-digitalocean-backend').then(m => ({
    default: m.catalogModuleOnemindDigitalOcean,
  })),
);
// Authentik users and groups. Disabled unless authentik.orgSync.enabled is true; see docs/catalog/authentik-org-sync.md.
backend.add(import('@internal/plugin-onemind-org-backend'));

backend.add(import('@backstage/plugin-permission-backend'));
backend.add(import('@internal/plugin-onemind-permissions-backend'));

backend.add(import('@backstage/plugin-search-backend'));
backend.add(import('@backstage/plugin-search-backend-module-pg'));
backend.add(import('@backstage/plugin-search-backend-module-catalog'));
backend.add(import('@backstage/plugin-search-backend-module-techdocs'));

// Admin diagnostics. The OneMind policy makes it administrator-only and closes /config.
backend.add(import('@backstage/plugin-devtools-backend'));
backend.add(import('@backstage/plugin-kubernetes-backend'));
backend.add(import('@backstage/plugin-user-settings-backend'));
backend.add(import('@backstage/plugin-events-backend'));
backend.add(import('@backstage/plugin-notifications-backend'));
backend.add(import('@backstage/plugin-signals-backend'));
backend.add(import('@backstage/plugin-mcp-actions-backend'));

backend.add(import('@internal/plugin-onemind-operations-backend'));
backend.add(import('@internal/plugin-onemind-missions-backend'));
backend.add(import('@internal/plugin-onemind-events-backend'));
// Live bridge: read-only NATS and Kubernetes watchers pushed to browsers over Signals; see plugins/plugin-onemind-livebridge-backend/README.md.
backend.add(import('@internal/plugin-onemind-livebridge-backend'));
backend.add(import('@internal/plugin-onemind-activity-backend'));
backend.add(import('@internal/plugin-onemind-services-backend'));
// F8 app lifecycle: read-only scaffolder actions (Add/Retire app) and the weekly stack hygiene report.
backend.add(
  import('@internal/plugin-onemind-services-backend').then(m => ({
    default: m.lifecycleScaffolderModule,
  })),
);
backend.add(
  import('@internal/plugin-onemind-services-backend').then(m => ({
    default: m.lifecyclePlugin,
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
