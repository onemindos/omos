/*
 * Hi!
 *
 * Note that this is an EXAMPLE Backstage backend. Please check the README.
 *
 * Happy hacking!
 */

import { createBackend } from '@backstage/backend-defaults';

const backend = createBackend();

backend.add(import('@backstage/plugin-app-backend'));
backend.add(import('@backstage/plugin-proxy-backend'));

// scaffolder plugin
backend.add(import('@backstage/plugin-scaffolder-backend'));
backend.add(import('@backstage/plugin-scaffolder-backend-module-github'));
backend.add(
  import('@backstage/plugin-scaffolder-backend-module-notifications'),
);

// techdocs plugin
backend.add(import('@backstage/plugin-techdocs-backend'));

// auth plugin
backend.add(import('@backstage/plugin-auth-backend'));
// See https://backstage.io/docs/backend-system/building-backends/migrating#the-auth-plugin
backend.add(import('./oidc'));
// See https://backstage.io/docs/auth/guest/provider

// catalog plugin
backend.add(import('@backstage/plugin-catalog-backend'));
backend.add(
  import('@backstage/plugin-catalog-backend-module-scaffolder-entity-model'),
);

// See https://backstage.io/docs/features/software-catalog/configuration#subscribing-to-catalog-errors
backend.add(import('@backstage/plugin-catalog-backend-module-logs'));

// permission plugin
backend.add(import('@backstage/plugin-permission-backend'));
// See https://backstage.io/docs/permissions/getting-started for how to create your own permission policy
backend.add(import('./permissions'));

// search plugin
backend.add(import('@backstage/plugin-search-backend'));

// search engine
// See https://backstage.io/docs/features/search/search-engines
backend.add(import('@backstage/plugin-search-backend-module-pg'));

// search collators
backend.add(import('@backstage/plugin-search-backend-module-catalog'));
backend.add(import('@backstage/plugin-search-backend-module-techdocs'));

// kubernetes plugin
backend.add(import('@backstage/plugin-kubernetes-backend'));

// user settings plugin
backend.add(import('@backstage/plugin-user-settings-backend'));

// Shared database-backed event delivery across backend replicas. No HTTP ingress topics.
backend.add(import('@backstage/plugin-events-backend'));

// notifications and signals plugins
backend.add(import('@backstage/plugin-notifications-backend'));
backend.add(import('@backstage/plugin-signals-backend'));

// mcp actions plugin
backend.add(import('@backstage/plugin-mcp-actions-backend'));

backend.add(import('./operations'));
backend.add(import('./missions'));
backend.add(import('./operationalEvents'));
backend.add(import('./activity'));
backend.add(import('./serviceOperations'));
backend.add(import('./context'));
backend.add(import('./machineReadOnly'));
backend.add(import('./assistant'));
backend.add(import('./offline'));

// Community plugins; their provider credentials are configured separately.
backend.add(import('@backstage/plugin-auth-backend-module-github-provider'));
backend.add(import('@backstage-community/plugin-argocd-backend'));
backend.add(import('./communityReadOnly'));
backend.add(import('@backstage-community/plugin-tech-insights-backend'));
backend.add(import('@backstage-community/plugin-tech-insights-backend-module-jsonfc'));
backend.add(import('@backstage-community/plugin-adr-backend'));
backend.add(import('@backstage-community/search-backend-module-adr'));
backend.add(import('@backstage-community/plugin-announcements-backend'));
backend.add(import('@backstage-community/plugin-search-backend-module-announcements'));

backend.start();
