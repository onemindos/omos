import githubActionsPlugin from '@backstage-community/plugin-github-actions/alpha';
import githubPullRequestsPlugin from '@roadiehq/backstage-plugin-github-pull-requests/alpha';
import grafanaPlugin from '@backstage-community/plugin-grafana/alpha';
import argocdPlugin, {
  argocdTranslationsModule,
} from '@backstage-community/plugin-argocd';
import techInsightsPlugin from '@backstage-community/plugin-tech-insights/alpha';
import adrPlugin from '@backstage-community/plugin-adr/alpha';
import announcementsPlugin from '@backstage-community/plugin-announcements/alpha';
import { assistantPlugin } from '@internal/plugin-onemind-assistant';
import { authModule } from '@internal/plugin-onemind-auth';
import { costPlugin } from '@internal/plugin-onemind-cost';
import { missionsPlugin } from '@internal/plugin-onemind-missions';
import { navModule } from '@internal/plugin-onemind-nav';
import offlinePlugin from '@internal/plugin-onemind-offline';
import { oneMindRadar } from '@internal/plugin-onemind-radar';
import {
  servicesModule,
  accessModule,
  standardAppPageModule,
  githubGateModules,
} from '@internal/plugin-onemind-services';
import { systemsPlugin } from '@internal/plugin-onemind-systems';
import { designModule } from '@internal/plugin-onemind-theme';
import '@internal/plugin-onemind-theme/brand.css';
import { devtoolsFeature } from './devtoolsFeature';
import { unprocessedEntitiesFeature } from './unprocessedEntitiesFeature';
import { shortcutsFeature } from './shortcutsFeature';
import { createApp } from '@backstage/frontend-defaults';
import catalogPlugin from '@backstage/plugin-catalog/alpha';
import { homeModule, todayPlugin } from '@internal/plugin-onemind-home';
import { livePlugin } from '@internal/plugin-onemind-live';
import { prometheusPlugin, prometheusApiModule } from './prometheusPlugin';
import { catalogCompositionModule } from './catalogCompositionModules';
import { kyvernoPolicyModule } from './kyvernoPolicy';
import { systemRedirectsModule } from './systemRedirects';

export default createApp({
  features: [
    catalogCompositionModule,
    kyvernoPolicyModule,
    systemRedirectsModule,
    githubActionsPlugin,
    githubPullRequestsPlugin,
    grafanaPlugin,
    prometheusPlugin,
    prometheusApiModule,
    argocdPlugin,
    argocdTranslationsModule,
    techInsightsPlugin,
    adrPlugin,
    oneMindRadar,
    announcementsPlugin,
    designModule,
    catalogPlugin,
    navModule,
    homeModule,
    todayPlugin,
    livePlugin,
    costPlugin,
    shortcutsFeature,
    devtoolsFeature,
    unprocessedEntitiesFeature,
    authModule,
    systemsPlugin,
    servicesModule,
    accessModule,
    standardAppPageModule,
    ...githubGateModules,
    assistantPlugin,
    offlinePlugin,
    missionsPlugin,
  ],
});
