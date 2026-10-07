import { createApp } from '@backstage/frontend-defaults';
import catalogPlugin from '@backstage/plugin-catalog/alpha';
import { navModule } from './modules/nav';
import { homeModule } from './modules/home';
import { systemsPlugin } from './modules/systems';
import { authModule } from './modules/auth';
import { servicesModule } from './modules/services';
import { assistantPlugin } from './modules/assistant';
import { offlinePlugin } from './modules/offline';

export default createApp({
  features: [catalogPlugin, navModule, homeModule, authModule, systemsPlugin, servicesModule, assistantPlugin, offlinePlugin],
});
