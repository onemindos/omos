import {
  ApiBlueprint,
  createFrontendModule,
  discoveryApiRef,
  fetchApiRef,
} from '@backstage/frontend-plugin-api';
import {
  EntityCardBlueprint,
  EntityContentBlueprint,
  EntityHeaderLayoutBlueprint,
} from '@backstage/plugin-catalog-react/alpha';
import {
  appStatusApiRef,
  createOperationsStatusSource,
  targetForEntity,
  EntityLike,
} from '@internal/plugin-onemind-common';
import { tabFilters } from './links';

type Filterable = Parameters<typeof targetForEntity>[0];
const kindIs = (...kinds: string[]) => (e: { kind: string }) => kinds.includes(e.kind.toLowerCase());

/** Entities that have at least one status signal to show (generated annotations). */
export function hasStatusTarget(entity: Filterable) {
  const t = targetForEntity(entity);
  return !!(t.argoApps.length || t.workloads.length || t.alertScopes.length);
}

/** The status source is registered once; the Live Ops aggregator can replace it with one line. */
export const appStatusApiExtension = ApiBlueprint.make({
  name: 'app-status',
  params: define =>
    define({
      api: appStatusApiRef,
      deps: { discovery: discoveryApiRef, fetchApi: fetchApiRef },
      factory: ({ discovery, fetchApi }) => createOperationsStatusSource(discovery, fetchApi),
    }),
});

export const appHeaderExtension = EntityHeaderLayoutBlueprint.make({
  name: 'onemind-app',
  params: {
    filter: e => kindIs('component', 'system', 'domain')(e),
    loader: async () => (await import('./AppHeader')).AppHeader,
  },
});

export const statusSummaryCard = EntityCardBlueprint.make({
  name: 'onemind-live-status',
  params: {
    type: 'content',
    filter: e => kindIs('component')(e) && hasStatusTarget(e as unknown as EntityLike),
    loader: async () => {
      const { StatusSummaryCard } = await import('./StatusSummaryCard');
      return <StatusSummaryCard />;
    },
  },
});

export const keyLinksCard = EntityCardBlueprint.make({
  name: 'onemind-key-links',
  params: {
    type: 'info',
    filter: e => kindIs('component')(e),
    loader: async () => {
      const { KeyLinksCard } = await import('./KeyLinksCard');
      return <KeyLinksCard />;
    },
  },
});

export const membersStatusCard = EntityCardBlueprint.make({
  name: 'onemind-members-status',
  params: {
    type: 'content',
    filter: e => kindIs('system', 'domain')(e),
    loader: async () => {
      const { MembersStatusCard } = await import('./MembersStatusCard');
      return <MembersStatusCard />;
    },
  },
});

export const logsContent = EntityContentBlueprint.make({
  name: 'onemind-logs',
  params: {
    path: '/logs',
    title: 'Logs',
    group: 'logs',
    filter: e => kindIs('component')(e) && tabFilters.logs(e),
    loader: async () => {
      const { LogsContent } = await import('./contents');
      return <LogsContent />;
    },
  },
});

export const costContent = EntityContentBlueprint.make({
  name: 'onemind-cost',
  params: {
    path: '/cost',
    title: 'Cost',
    group: 'cost',
    filter: e => kindIs('component')(e) && tabFilters.cost(e),
    loader: async () => {
      const { CostContent } = await import('./contents');
      return <CostContent />;
    },
  },
});

export const runbookContent = EntityContentBlueprint.make({
  name: 'onemind-runbook',
  params: {
    path: '/runbook',
    title: 'Runbook',
    group: 'docs',
    filter: e => kindIs('component')(e) && tabFilters.runbook(e),
    loader: async () => {
      const { RunbookContent } = await import('./contents');
      return <RunbookContent />;
    },
  },
});

/** The standard app page (F5): header, summary cards and the Logs, Cost and Docs tabs. */
export const standardAppPageModule = createFrontendModule({
  pluginId: 'catalog',
  extensions: [
    appStatusApiExtension,
    appHeaderExtension,
    statusSummaryCard,
    keyLinksCard,
    membersStatusCard,
    logsContent,
    costContent,
    runbookContent,
  ],
});
