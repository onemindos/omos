import {
  createFrontendPlugin,
  PageBlueprint,
} from '@backstage/frontend-plugin-api';
import { EntityCardBlueprint } from '@backstage/plugin-catalog-react/alpha';
import CloudIcon from '@material-ui/icons/Cloud';

/** /digitalocean (read-only account summary, administrators) and the "Runs on DigitalOcean" app card. */
export const digitalOceanPlugin = createFrontendPlugin({
  pluginId: 'onemind-digitalocean',
  extensions: [
    PageBlueprint.make({
      params: {
        path: '/digitalocean',
        title: 'DigitalOcean',
        icon: <CloudIcon />,
        loader: async () => {
          const { DigitalOceanPage } = await import('./DigitalOceanPage');
          return <DigitalOceanPage />;
        },
      },
    }),
    EntityCardBlueprint.make({
      name: 'runs-on-digitalocean',
      params: {
        type: 'content',
        filter: e => e.kind.toLowerCase() === 'component',
        loader: async () => {
          const { RunsOnCard } = await import('./RunsOnCard');
          return <RunsOnCard />;
        },
      },
    }),
  ],
});
