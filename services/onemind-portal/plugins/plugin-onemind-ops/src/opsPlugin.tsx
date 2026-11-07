import {
  createFrontendPlugin,
  PageBlueprint,
} from '@backstage/frontend-plugin-api';
import RadarIcon from '@material-ui/icons/TrackChanges';

export const opsPlugin = createFrontendPlugin({
  pluginId: 'onemind-ops',
  extensions: [
    PageBlueprint.make({
      params: {
        path: '/ops',
        title: 'Ops Center',
        icon: <RadarIcon />,
        loader: async () => {
          const { OpsPage } = await import('./OpsPage');
          return <OpsPage />;
        },
      },
    }),
  ],
});
