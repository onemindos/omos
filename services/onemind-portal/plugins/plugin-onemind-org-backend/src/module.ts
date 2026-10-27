import { resolve } from 'node:path';
import {
  coreServices,
  createBackendModule,
  type LoggerService,
  type SchedulerService,
} from '@backstage/backend-plugin-api';
import type { Config } from '@backstage/config';
import {
  catalogProcessingExtensionPoint,
  type CatalogProcessingExtensionPoint,
} from '@backstage/plugin-catalog-node';
import {
  authentikFetchJson,
  readAuthentikOrg,
  type FetchJson,
} from './authentik';
import { readAuthentikOrgConfig } from './config';
import { PROVIDER_NAME, readCurated } from './entities';
import { AuthentikOrgProvider } from './provider';

export interface RegisterDeps {
  config: Config;
  logger: LoggerService;
  scheduler: SchedulerService;
  processing: Pick<CatalogProcessingExtensionPoint, 'addEntityProvider'>;
  /** Test seam: replaces the HTTP read of the Authentik API. */
  fetchJson?: FetchJson;
}

/**
 * Registers the Authentik org sync. With `authentik.orgSync.enabled` off (the default) this logs one
 * line and registers nothing: no provider, no scheduled task, no Authentik call.
 */
export function registerAuthentikOrg({
  config,
  logger,
  scheduler,
  processing,
  fetchJson,
}: RegisterDeps): boolean {
  const settings = readAuthentikOrgConfig(config);
  if (!settings.enabled) {
    logger.info(
      'Authentik org sync is disabled (authentik.orgSync.enabled is not true); users and groups come only from the reviewed org file',
    );
    return false;
  }
  const read = fetchJson ?? authentikFetchJson(settings.token);
  const provider = new AuthentikOrgProvider({
    userGroups: settings.userGroups,
    curatedFile: settings.curatedFile,
    logger,
    readSnapshot: signal => readAuthentikOrg(read, settings.baseUrl, signal),
    readCurated: () => readCurated(resolve(settings.curatedFile)),
  });
  processing.addEntityProvider(provider);
  // Both portal replicas register the provider; the global scope runs each sync on one of them.
  scheduler.scheduleTask({
    id: `${PROVIDER_NAME}-refresh`,
    frequency: { seconds: settings.frequencySeconds },
    timeout: { seconds: settings.timeoutSeconds },
    initialDelay: { seconds: settings.initialDelaySeconds },
    scope: 'global',
    fn: async (signal: AbortSignal) => {
      await provider.refresh(signal);
    },
  });
  logger.info(
    `Authentik org sync enabled: members of ${settings.userGroups.join(
      ', ',
    )} every ${settings.frequencySeconds}s from ${settings.baseUrl}`,
  );
  return true;
}

export const catalogModuleOnemindAuthentikOrg = createBackendModule({
  pluginId: 'catalog',
  moduleId: 'onemind-authentik-org',
  register(reg) {
    reg.registerInit({
      deps: {
        config: coreServices.rootConfig,
        logger: coreServices.logger,
        scheduler: coreServices.scheduler,
        processing: catalogProcessingExtensionPoint,
      },
      async init(deps) {
        registerAuthentikOrg(deps);
      },
    });
  },
});
