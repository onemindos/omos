import { Navigate } from 'react-router-dom';
import { PageBlueprint, createFrontendModule } from '@backstage/frontend-plugin-api';
import redirects from './systemRedirects.json';

// Systems removed by the F2 restructure keep working URLs. Backstage has no entity redirect, so
// each removed name gets its own static page route. A static path outranks the catalog's
// `/catalog/:namespace/:kind/:name` route in React Router, and no entity exists under these
// names any more, so live entities are never intercepted. Source: catalog/redirects.yaml.
export const systemRedirectPath = (name: string) =>
  `/catalog/default/system/${name}`;

export const systemRedirectExtensions = Object.entries(
  redirects as Record<string, string>,
).map(([from, to]) =>
  PageBlueprint.make({
    name: `system-redirect-${from}`,
    params: {
      path: systemRedirectPath(from),
      loader: async () => <Navigate replace to={systemRedirectPath(to)} />,
    },
  }),
);

export const systemRedirectsModule = createFrontendModule({
  pluginId: 'catalog',
  extensions: systemRedirectExtensions,
});
