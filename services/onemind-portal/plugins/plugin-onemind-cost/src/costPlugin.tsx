import {
  createFrontendPlugin,
  PageBlueprint,
} from '@backstage/frontend-plugin-api';
import MonetizationOnIcon from '@material-ui/icons/MonetizationOn';

/**
 * The plugin id is `cost-insights` so the page extension is `page:cost-insights`,
 * the id the OneMind sidebar already looks for under Insights.
 */
export const costPlugin = createFrontendPlugin({
  pluginId: 'cost-insights',
  extensions: [
    PageBlueprint.make({
      params: {
        path: '/cost',
        title: 'Cost Insights',
        icon: <MonetizationOnIcon />,
        loader: async () => {
          const { CostPage } = await import('./CostPage');
          return <CostPage />;
        },
      },
    }),
  ],
});
