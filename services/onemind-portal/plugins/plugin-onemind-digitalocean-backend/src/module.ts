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
import { inClusterGetJson, type GetJson } from './kube';
import { readToken } from './client';
import { linkVolumes, listPvcs, type WorkloadComponent } from './links';
import { DigitalOceanProvider } from './provider';
import { sharedService, type DigitalOceanService } from './service';

export interface RegisterDeps {
  config: Config;
  logger: LoggerService;
  scheduler: SchedulerService;
  auth: AuthService;
  catalog: CatalogService;
  processing: Pick<CatalogProcessingExtensionPoint, 'addEntityProvider'>;
  /** Test seams. */
  service?: DigitalOceanService;
  getJson?: GetJson;
  env?: NodeJS.ProcessEnv;
}

export const REFRESH_SECONDS = 600;

/**
 * Registers the DigitalOcean provider. With no DIGITALOCEAN_TOKEN it logs one line and registers
 * nothing: no provider, no scheduled task, no DigitalOcean or cluster call.
 */
export function registerDigitalOcean({
  config,
  logger,
  scheduler,
  auth,
  catalog,
  processing,
  service: injected,
  getJson: injectedGetJson,
  env,
}: RegisterDeps): boolean {
  if (!readToken(env)) {
    logger.info(
      'DigitalOcean catalog is off (DIGITALOCEAN_TOKEN is not set); no DigitalOcean resources are published',
    );
    return false;
  }
  const service = injected ?? sharedService(m => logger.warn(m));
  const apiServer =
    config.getOptionalString('onemind.inventory.live.apiServer') ??
    'https://kubernetes.default.svc';
  const getJson = injectedGetJson ?? inClusterGetJson(apiServer, 60_000);

  service.setLinker(async volumes => {
    const [pvcs, { items }] = await Promise.all([
      listPvcs(getJson),
      (async () =>
        catalog.getEntities(
          {
            filter: [{ kind: 'Component' }],
            fields: [
              'kind',
              'metadata.name',
              'metadata.title',
              'metadata.namespace',
              'metadata.annotations',
            ],
          },
          { credentials: await auth.getOwnServiceCredentials() },
        ))(),
    ]);
    const workloads: WorkloadComponent[] = [];
    for (const e of items) {
      const a = e.metadata.annotations ?? {};
      const workload = a['onemind.io/kubernetes-workload'];
      const namespace = a['backstage.io/kubernetes-namespace'];
      if (!workload || !namespace) continue;
      workloads.push({
        ref: `component:${e.metadata.namespace ?? 'default'}/${
          e.metadata.name
        }`,
        title: e.metadata.title,
        namespace,
        workload,
        selector: a['backstage.io/kubernetes-label-selector'] ?? '',
      });
    }
    return linkVolumes(volumes, pvcs, workloads);
  });

  const provider = new DigitalOceanProvider({ service, logger });
  processing.addEntityProvider(provider);
  scheduler.scheduleTask({
    id: 'onemind-digitalocean-refresh',
    frequency: { seconds: REFRESH_SECONDS },
    timeout: { seconds: 300 },
    initialDelay: { seconds: 30 },
    scope: 'global',
    fn: async () => {
      await provider.refresh();
    },
  });
  logger.info(
    `DigitalOcean catalog enabled: refreshing every ${REFRESH_SECONDS}s`,
  );
  return true;
}

export const catalogModuleOnemindDigitalOcean = createBackendModule({
  pluginId: 'catalog',
  moduleId: 'onemind-digitalocean',
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
        registerDigitalOcean(deps);
      },
    });
  },
});
