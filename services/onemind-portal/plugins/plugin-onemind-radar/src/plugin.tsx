import { techRadarApi } from '@backstage-community/plugin-tech-radar/alpha';
import techRadarPlugin from '@backstage-community/plugin-tech-radar/alpha';
import {
  techRadarApiRef,
  TechRadarApi,
} from '@backstage-community/plugin-tech-radar';
import { createApiFactory } from '@backstage/core-plugin-api';
import radar from './radar-data.json';

// Use the plugin's API extension with Git-owned data; never show its demo fallback.
const radarApi: TechRadarApi = {
  async load() {
    return {
      ...radar,
      entries: radar.entries.map(entry => ({
        ...entry,
        timeline: entry.timeline.map(snapshot => ({
          ...snapshot,
          date: new Date(snapshot.date),
        })),
      })),
    };
  },
};

export const oneMindRadar = techRadarPlugin.withOverrides({
  extensions: [
    techRadarApi.override({
      params: define => define(createApiFactory(techRadarApiRef, radarApi)),
    }),
  ],
});

export default oneMindRadar;
