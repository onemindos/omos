import {
  createFrontendPlugin,
  PageBlueprint,
} from '@backstage/frontend-plugin-api';

/**
 * Page `/stack-hygiene`. The extension id is `page:onemind-stack-hygiene`, the id the
 * OneMind sidebar looks for under Insights.
 */
export const stackHygienePlugin = createFrontendPlugin({
  pluginId: 'onemind-stack-hygiene',
  extensions: [
    PageBlueprint.make({
      params: {
        path: '/stack-hygiene',
        title: 'Stack hygiene',
        loader: async () => {
          const { StackHygienePage } = await import('./StackHygienePage');
          return <StackHygienePage />;
        },
      },
    }),
  ],
});
