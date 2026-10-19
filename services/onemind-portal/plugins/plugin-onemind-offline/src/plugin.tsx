import GetAppIcon from '@material-ui/icons/GetApp';
import { createFrontendPlugin, PageBlueprint } from '@backstage/frontend-plugin-api';
import { OfflinePage } from './components/OfflinePage';

export const offlinePlugin = createFrontendPlugin({
  pluginId: 'onemind-offline',
  extensions: [
    PageBlueprint.make({
      params: {
        path: '/offline',
        title: 'Offline knowledge',
        icon: <GetAppIcon />,
        loader: async () => <OfflinePage />,
      },
    }),
  ],
});
