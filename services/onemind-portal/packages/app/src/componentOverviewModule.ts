import { createFrontendModule } from '@backstage/frontend-plugin-api';
import openchoreoPlugin from '@openchoreo/backstage-plugin/alpha';

export const componentOverviewExtension = openchoreoPlugin
  .getExtension('entity-content-layout:openchoreo/component-overview')
  .override({
    params: {
      loader: () => import('./ComponentOverview').then(m => m.default),
    },
  });

export const componentOverviewModule = createFrontendModule({
  pluginId: 'openchoreo',
  extensions: [componentOverviewExtension],
});
