import { createFrontendPlugin, PageBlueprint } from '@backstage/frontend-plugin-api';
import TimelineIcon from '@material-ui/icons/Timeline';

export const livePlugin = createFrontendPlugin({
  pluginId: 'onemind-live',
  extensions: [
    PageBlueprint.make({
      params: {
        path: '/live',
        title: 'Live Ops',
        icon: <TimelineIcon />,
        loader: async () => {
          const { LivePage } = await import('./LivePage');
          return <LivePage />;
        },
      },
    }),
  ],
});
