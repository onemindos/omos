import { createFrontendPlugin, PageBlueprint } from '@backstage/frontend-plugin-api';

/**
 * Page `/non-kubernetes`. The extension id is `page:onemind-non-kubernetes`, the id the OneMind sidebar
 * looks for under Services.
 */
export const nonKubernetesPlugin = createFrontendPlugin({
  pluginId: 'onemind-non-kubernetes',
  extensions: [
    PageBlueprint.make({
      params: {
        path: '/non-kubernetes',
        title: 'Not on Kubernetes',
        loader: async () => {
          const { NonKubernetesPage } = await import('./NonKubernetesPage');
          return <NonKubernetesPage />;
        },
      },
    }),
  ],
});
