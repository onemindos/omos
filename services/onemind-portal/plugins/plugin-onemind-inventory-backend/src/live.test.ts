import { readFileSync } from 'node:fs';
import {
  coreServices,
  createBackendModule,
  type SchedulerService,
} from '@backstage/backend-plugin-api';
import {
  catalogProcessingExtensionPoint,
  catalogServiceRef,
} from '@backstage/plugin-catalog-node';
import {
  mockCredentials,
  mockServices,
  startTestBackend,
} from '@backstage/backend-test-utils';
import catalogPlugin from '@backstage/plugin-catalog-backend';
import type { Entity } from '@backstage/catalog-model';
import { registerLiveInventory } from './module';
import { A, entityName } from './entities';
import { clusterYamlPath, deployment, recordedApi } from './testing';

jest.setTimeout(120_000);

let scheduler: SchedulerService;

/** The inventory module with the Kubernetes read replaced by the recorded cluster plus a new Deployment. */
const withRecordedCluster = createBackendModule({
  pluginId: 'catalog',
  moduleId: 'test-live-inventory',
  register(reg) {
    reg.registerInit({
      deps: {
        config: coreServices.rootConfig,
        logger: coreServices.logger,
        scheduler: coreServices.scheduler,
        auth: coreServices.auth,
        catalog: catalogServiceRef,
        processing: catalogProcessingExtensionPoint,
      },
      async init(deps) {
        scheduler = deps.scheduler;
        registerLiveInventory({
          ...deps,
          getJson: recordedApi([
            deployment('team-a', 'worker'),
            deployment('dp-default-shop-development-1234abcd', 'shop-api'),
          ]),
        });
      },
    });
  },
});

describe('live inventory enabled, against a real catalog', () => {
  it('adds the new workload once, stamps freshness on reviewed ones, flags the vanished, and never touches dp-*', async () => {
    const backend = await startTestBackend({
      features: [
        catalogPlugin,
        withRecordedCluster,
        mockServices.scheduler.factory({ skipTaskRunOnStartup: true }),
        mockServices.rootConfig.factory({
          data: {
            catalog: {
              processingInterval: { seconds: 3 },
              rules: [
                {
                  allow: [
                    'Component',
                    'System',
                    'API',
                    'Resource',
                    'Location',
                    'Domain',
                  ],
                },
              ],
              locations: [{ type: 'file', target: clusterYamlPath }],
            },
            onemind: {
              inventory: {
                live: {
                  enabled: true,
                  initialDelaySeconds: 1,
                  curatedFile: clusterYamlPath,
                },
              },
            },
          },
        }),
      ],
    });
    const url = `http://localhost:${backend.server.port()}/api/catalog/entities?filter=kind=component`;
    const headers = { Authorization: mockCredentials.user.header() };
    // The scheduler is a test double that only runs a task when asked; trigger the real tasks by id.
    const settle = async () => {
      for (let i = 0; i < 120; i++) {
        const count = (
          (await (await fetch(url, { headers })).json()) as Entity[]
        ).length;
        if (count >= 154) return;
        await new Promise(r => setTimeout(r, 500));
      }
    };
    await settle();
    await scheduler.triggerTask('onemind-live-inventory-observe');
    await scheduler.triggerTask('onemind-live-inventory-publish');
    const worker = entityName('team-a', 'worker');
    const dp = entityName('dp-default-shop-development-1234abcd', 'shop-api');
    let items: Entity[] = [];
    for (let i = 0; i < 120; i++) {
      items = (await (await fetch(url, { headers })).json()) as Entity[];
      const reviewedDone = items.filter(
        e => e.metadata.annotations?.[A.checkedAt],
      );
      const flagged = items.find(
        e => e.metadata.name === entityName('nats', 'agenttap'),
      )?.spec?.lifecycle;
      // Processing is a loop: wait until the new entity exists and every reviewed one was re-stamped.
      if (
        items.some(e => e.metadata.name === worker) &&
        flagged === 'removed' &&
        reviewedDone.length >= 143
      )
        break;
      await new Promise(r => setTimeout(r, 500));
    }
    await backend.stop();

    const byName = new Map(items.map(e => [e.metadata.name, e]));
    const added = byName.get(worker)!;
    expect(added.spec).toMatchObject({
      system: 'unassigned-workloads',
      lifecycle: 'production',
    });
    expect(added.metadata.annotations![A.owner]).toBe('onemind-live-inventory');
    expect(byName.has(dp)).toBe(false);
    expect(items.filter(e => e.metadata.name === worker)).toHaveLength(1);

    const source = readFileSync(clusterYamlPath, 'utf8');
    const reviewed = items.filter(
      e =>
        e.metadata.annotations?.[A.source] === 'kubernetes' &&
        !e.metadata.annotations[A.owner],
    );
    expect(reviewed).toHaveLength(
      (source.match(/onemind.io\/inventory-source: kubernetes/g) ?? []).length,
    );
    const portal = byName.get('k8s-backstage-backstage-21b99354')!;
    expect(portal.metadata.annotations![A.state]).toBe('observed');
    expect(portal.spec!.subcomponentOf).toBe('component:default/backstage'); // reviewed placement untouched
    const gone = byName.get(entityName('nats', 'agenttap'))!;
    expect(gone.metadata.annotations![A.state]).toBe('missing-review-required');
    expect(gone.spec!.lifecycle).toBe('removed');
  });
});
