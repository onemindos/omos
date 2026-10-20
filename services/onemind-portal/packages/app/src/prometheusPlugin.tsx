import { convertLegacyPlugin } from '@backstage/core-compat-api';
import {
  ApiBlueprint,
  createFrontendModule,
} from '@backstage/frontend-plugin-api';
import { backstagePluginPrometheusPlugin } from '@roadiehq/backstage-plugin-prometheus';
import { EntityCardBlueprint } from '@backstage/plugin-catalog-react/alpha';
import { discoveryApiRef, fetchApiRef } from '@backstage/core-plugin-api';
import {
  createScopedPrometheusApi,
  scopeForEntity,
  scopedPrometheusApiRef,
} from './prometheusApi';

export const prometheusCard = EntityCardBlueprint.make({
  name: 'service-alerts',
  params: {
    filter: entity => !!scopeForEntity(entity),
    loader: async () => {
      const { default: PrometheusAlerts } = await import('./PrometheusAlerts');
      return <PrometheusAlerts />;
    },
  },
});
export const prometheusPlugin = convertLegacyPlugin(
  backstagePluginPrometheusPlugin,
  {
    extensions: [prometheusCard],
  },
);
// Modules override an existing plugin extension by the exact public API ID.
export const prometheusApiExtension = ApiBlueprint.make({
  name: scopedPrometheusApiRef.id,
  params: define =>
    define({
      api: scopedPrometheusApiRef,
      deps: { discovery: discoveryApiRef, fetchApi: fetchApiRef },
      factory: ({ discovery, fetchApi }) =>
        createScopedPrometheusApi(discovery, fetchApi),
    }),
});
export const prometheusApiModule = createFrontendModule({
  pluginId: 'backstage-plugin-prometheus',
  extensions: [prometheusApiExtension],
});
