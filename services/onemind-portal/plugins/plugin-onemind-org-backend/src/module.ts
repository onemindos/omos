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
import { DEFAULT_SCHEDULE, readAuthentikOrgConfig } from './config';
import { PROVIDER_NAME, readCurated } from './entities';
import { AuthentikOrgCleanup, AuthentikOrgProvider } from './provider';

export interface RegisterDeps {
  config: Config;
  logger: LoggerService;
  scheduler: SchedulerService;
  processing: Pick<CatalogProcessingExtensionPoint, 'addEntityProvider'>;
  /** Test seam: replaces the HTTP read of the Authentik API. */
  fetchJson?: FetchJson;
}

/** What the module runs: the sync, or only the removal of entities an earlier sync published. */
export type AuthentikOrgMode = 'sync' | 'cleanup';

/** One task id in both modes, so that during a rollout old and new replicas share its runs. */
export const TASK_ID = `${PROVIDER_NAME}-refresh`;

/**
 * Registers the Authentik org sync. A provider named `onemind-authentik-org` and the global task
 * `onemind-authentik-org-refresh` are registered in every state:
 * - enabled with a token: each run reads Authentik and publishes the result;
 * - disabled (the default), or enabled while `authentik.token` is unset: nothing calls Authentik; the
 *   provider removes what the sync published when the catalog connects it and on every run of the
 *   task, so turning the sync off, or losing its secret, revokes synced users even if a replica of
 *   the enabled rollout runs the task once more.
 * Invalid reviewed settings still fail startup; an absent token never does. So does a task that
 * could not be stored (a database error): registration waits for the scheduler, so a pod never
 * serves without the task that the rollback relies on.
 */
export async function registerAuthentikOrg({
  config,
  logger,
  scheduler,
  processing,
  fetchJson,
}: RegisterDeps): Promise<AuthentikOrgMode> {
  const settings = readAuthentikOrgConfig(config);
  // Both portal replicas register the task; the global scope runs each run on one of them.
  const schedule = (
    timing: {
      frequencySeconds: number;
      timeoutSeconds: number;
      initialDelaySeconds: number;
    },
    fn: (signal: AbortSignal) => Promise<unknown>,
  ) =>
    scheduler.scheduleTask({
      id: TASK_ID,
      frequency: { seconds: timing.frequencySeconds },
      timeout: { seconds: timing.timeoutSeconds },
      initialDelay: { seconds: timing.initialDelaySeconds },
      scope: 'global',
      fn: async (signal: AbortSignal) => {
        await fn(signal);
      },
    });

  if (!settings.enabled || !settings.token) {
    if (settings.enabled)
      logger.error(
        'Authentik org sync is enabled but authentik.token is not set (in the cluster: Secret backstage-authentik-catalog, key AUTHENTIK_CATALOG_TOKEN, read when the pod starts); it calls nothing and removes the users and groups it published',
      );
    else
      logger.info(
        'Authentik org sync is disabled (authentik.orgSync.enabled is not true); it calls nothing and only removes entities it published earlier, so users and groups come only from the reviewed org file',
      );
    const cleanup = new AuthentikOrgCleanup(logger);
    processing.addEntityProvider(cleanup);
    await schedule(DEFAULT_SCHEDULE, () => cleanup.clear());
    return 'cleanup';
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
  await schedule(settings, signal => provider.refresh(signal));
  logger.info(
    `Authentik org sync enabled: members of ${settings.userGroups.join(
      ', ',
    )} every ${settings.frequencySeconds}s from ${settings.baseUrl}`,
  );
  return 'sync';
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
        await registerAuthentikOrg(deps);
      },
    });
  },
});
