import { createBackendModule } from '@backstage/backend-plugin-api';
import { actionsRegistryServiceRef } from '@backstage/backend-plugin-api/alpha';
import { health } from './health';
import { readFileSync } from 'node:fs';

const runbooks = ['index', 'backstage', 'identity', 'nats', 'tak', 'monitoring', 'resilience', 'ai', 'cloud', 'workflows', 'inventory-refresh', 'assistant', 'offline'];
export default createBackendModule({
  pluginId: 'onemind', moduleId: 'operational-context',
  register(reg) {
    reg.registerInit({ deps: { actions: actionsRegistryServiceRef }, async init({ actions }) {
      actions.register({
        name: 'health', title: 'OneMind observed service health',
        description: 'Read configured internal HTTP health observations and namespace metrics. Reachability is not application correctness.',
        attributes: { readOnly: true, destructive: false, idempotent: true },
        schema: { input: z => z.object({}), output: z => z.object({ observedAt: z.string(), source: z.string(), observations: z.unknown() }) },
        async action() { return { output: { observedAt: new Date().toISOString(), source: 'do-nyc3-omos-cloud; configured internal probes and Prometheus', observations: await health() } }; },
      });
      actions.register({
        name: 'runbook', title: 'Read OneMind runbook',
        description: 'Read Git-owned operational documentation. Documented instructions are not live observations.',
        attributes: { readOnly: true, destructive: false, idempotent: true },
        schema: { input: z => z.object({ name: z.enum(['index', 'backstage', 'identity', 'nats', 'tak', 'monitoring', 'resilience', 'ai', 'cloud', 'workflows', 'inventory-refresh', 'assistant', 'offline']) }), output: z => z.object({ source: z.string(), freshness: z.string(), content: z.string() }) },
        async action({ input }) {
          if (!runbooks.includes(input.name)) throw new Error('Unknown runbook');
          return { output: { source: `Git: services/onemind-portal/runbooks/docs/${input.name}.md`, freshness: 'Documentation packaged in the deployed image; not live state', content: readFileSync(`/app/runbooks/docs/${input.name}.md`, 'utf8') } };
        },
      });
    } });
  },
});
