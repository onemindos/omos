import githubActionsPlugin from '@backstage-community/plugin-github-actions/alpha';
import githubPullRequestsPlugin from '@roadiehq/backstage-plugin-github-pull-requests/alpha';
import grafanaPlugin from '@backstage-community/plugin-grafana/alpha';
import argocdPlugin, { argocdTranslationsModule } from '@backstage-community/plugin-argocd';
import techInsightsPlugin from '@backstage-community/plugin-tech-insights/alpha';
import adrPlugin from '@backstage-community/plugin-adr/alpha';
import techRadarPlugin, { techRadarApi } from '@backstage-community/plugin-tech-radar/alpha';
import announcementsPlugin from '@backstage-community/plugin-announcements/alpha';
import { techRadarApiRef, TechRadarApi } from '@backstage-community/plugin-tech-radar';
import { createApiFactory } from '@backstage/core-plugin-api';
import radar from './modules/radar/data.json';

// Use the plugin's API extension with Git-owned data; never show its demo fallback.
const radarApi: TechRadarApi = {
  async load() {
    return { ...radar, entries: radar.entries.map(entry => ({ ...entry,
      timeline: entry.timeline.map(snapshot => ({ ...snapshot, date: new Date(snapshot.date) })),
    })) };
  },
};
const oneMindRadar = techRadarPlugin.withOverrides({
  extensions: [techRadarApi.override({
    params: define => define(createApiFactory(techRadarApiRef, radarApi)),
  })],
});

import { designModule } from './theme/designModule';
import './theme/brand.css';
import { createApp } from '@backstage/frontend-defaults';
import catalogPlugin from '@backstage/plugin-catalog/alpha';
import { navModule } from './modules/nav';
import { homeModule } from './modules/home';
import { todayPlugin } from './modules/home/Today';
import { systemsPlugin } from './modules/systems';
import { authModule } from './modules/auth';
import { servicesModule } from './modules/services';
import { assistantPlugin } from './modules/assistant';
import { offlinePlugin } from './modules/offline';

export default createApp({
  features: [githubActionsPlugin, githubPullRequestsPlugin, grafanaPlugin, argocdPlugin, argocdTranslationsModule, techInsightsPlugin, adrPlugin, oneMindRadar, announcementsPlugin, designModule, catalogPlugin, navModule, homeModule, todayPlugin, authModule, systemsPlugin, servicesModule, assistantPlugin, offlinePlugin],
});
