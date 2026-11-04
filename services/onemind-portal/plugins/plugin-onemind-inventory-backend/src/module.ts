import { resolve } from 'node:path';
import {
  coreServices,
  createBackendModule,
  type AuthService,
  type LoggerService,
  type SchedulerService,
} from '@backstage/backend-plugin-api';
import type { Config } from '@backstage/config';
import {
  catalogProcessingExtensionPoint,
  catalogServiceRef,
  type CatalogProcessingExtensionPoint,
  type CatalogService,
} from '@backstage/plugin-catalog-node';
import { readLiveInventoryConfig } from './config';
import { PROVIDER_NAME } from './entities';
import { inClusterGetJson, listWorkloads, type GetJson } from './kube';
import { ObservationStore } from './observation';
import { LiveInventoryProcessor } from './processor';
import { LiveInventoryProvider, readCurated, toOwnedRecords } from './provider';

export interface RegisterDeps {
  config: Config;
  logger: LoggerService;
  scheduler: SchedulerService;
  auth: AuthService;
  catalog: CatalogService;
  /** Test seam: replaces the in-cluster Kubernetes read. */
  getJson?: GetJson;
  processing: Pick<
    CatalogProcessingExtensionPoint,
    'addEntityProvider' | 'addProcessor'
  >;
}

/**
 * Registers the live inventory. With `onemind.inventory.live.enabled` off (the default) this logs one
 * line and registers nothing: no provider, no processor, no scheduled task, no cluster call.
 */
export function registerLiveInventory({
  config,
  logger,
  scheduler,
  auth,
  catalog,
  processing,
  getJson: injected,
}: RegisterDeps): boolean {
  const settings = readLiveInventoryConfig(config);
  if (!settings.enabled) {
    logger.info(
      'Live inventory is disabled (onemind.inventory.live.enabled is not true); the file-based catalog is authoritative',
    );
    return false;
  }
  const store = new ObservationStore();
  const getJson =
    injected ??
    inClusterGetJson(settings.apiServer, settings.timeoutSeconds * 1000);
  const provider = new LiveInventoryProvider({
    config: settings,
    logger,
    store,
    listWorkloads: () => listWorkloads(getJson),
    readCurated: () => readCurated(resolve(settings.curatedFile)),
    readOwned: async () => {
      const { items } = await catalog.getEntities(
        {
          filter: {
            'metadata.annotations.onemind.io/inventory-owner': PROVIDER_NAME,
          },
          fields: ['kind', 'metadata.name', 'metadata.annotations'],
        },
        { credentials: await auth.getOwnServiceCredentials() },
      );
      return toOwnedRecords(items);
    },
  });
  processing.addEntityProvider(provider);
  processing.addProcessor(new LiveInventoryProcessor(store, settings));

  const every = { seconds: settings.frequencySeconds };
  const timeout = { seconds: settings.timeoutSeconds };
  // Every replica observes locally so its own processor can stamp freshness; only one publishes.
  scheduler.scheduleTask({
    id: `${PROVIDER_NAME}-observe`,
    frequency: every,
    timeout,
    initialDelay: { seconds: settings.initialDelaySeconds },
    scope: 'local',
    fn: async () => {
      await provider.observe();
    },
  });
  scheduler.scheduleTask({
    id: `${PROVIDER_NAME}-publish`,
    frequency: every,
    timeout,
    initialDelay: { seconds: settings.initialDelaySeconds + 5 },
    scope: 'global',
    fn: async () => {
      await provider.refresh();
    },
  });
  logger.info(
    `Live inventory enabled: observing every ${settings.frequencySeconds}s, removed workloads kept ${settings.removedRetentionHours}h`,
  );
  return true;
}

export const catalogModuleOnemindLiveInventory = createBackendModule({
  pluginId: 'catalog',
  moduleId: 'onemind-live-inventory',
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
        registerLiveInventory(deps);
      },
    });
  },
});
