import { createFrontendModule } from '@backstage/frontend-plugin-api';
import { EntityCardBlueprint } from '@backstage/plugin-catalog-react/alpha';
import { targetsForEntity } from './policyMatch';

/**
 * Kyverno results per app. The community Policy Reporter and TeraSky plugins were
 * checked and not used (see platform/backstage/README.md): one needs another in-cluster
 * workload, the other needs core packages newer than Backstage 1.55.3. This card reads
 * the PolicyReports Kyverno already writes, through the portal backend's existing
 * read-only cluster access (GET /api/onemind/policy-reports).
 */
export const policyResultsCard = EntityCardBlueprint.make({
  name: 'policy-results',
  params: {
    filter: entity => targetsForEntity(entity).length > 0,
    loader: async () => {
      const { default: PolicyResultsCard } = await import(
        './PolicyResultsCard'
      );
      return <PolicyResultsCard />;
    },
  },
});

export const kyvernoPolicyModule = createFrontendModule({
  pluginId: 'catalog',
  extensions: [policyResultsCard],
});
