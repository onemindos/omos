import { mkdir, writeFile } from 'node:fs/promises';
import path from 'node:path';
import {
  AuthService,
  coreServices,
  createBackendModule,
  resolveSafeChildPath,
  UrlReaderService,
} from '@backstage/backend-plugin-api';
import { CatalogService, catalogServiceRef } from '@backstage/plugin-catalog-node';
import { parseEntityRef, Entity } from '@backstage/catalog-model';
import {
  createTemplateAction,
  scaffolderActionsExtensionPoint,
} from '@backstage/plugin-scaffolder-node';
import { AddParams, prepareAdd } from './add';
import { buildRetirePlan, RetireInputError } from './retire';
import { collectClusterFacts, loadGithubSnapshot } from './sources';
import type { ClusterFacts, RepoReader } from './types';

/**
 * Two read-only scaffolder actions. They read the catalog, Git (GitHub App through the
 * URL reader) and, for retirement, the cluster (portal service account, list calls only),
 * and write files into the scaffolder workspace. The template's own publish step turns
 * the workspace into a draft pull request. Nothing here deletes, applies or syncs.
 */

async function writeWorkspace(workspace: string, files: Map<string, string>) {
  for (const [file, text] of files) {
    const target = resolveSafeChildPath(workspace, file);
    await mkdir(path.dirname(target), { recursive: true });
    await writeFile(target, text, 'utf8');
  }
}

export type LifecycleDeps = {
  repo: () => Promise<RepoReader>;
  cluster: () => Promise<ClusterFacts | undefined>;
  entities: () => Promise<Entity[]>;
  now?: () => Date;
};

export function createRetirePlanAction(deps: LifecycleDeps) {
  return createTemplateAction({
    id: 'onemind:app-retire:plan',
    description:
      'Read-only inventory of what is attached to an app and the file changes for a removal pull request. Deletes nothing.',
    supportsDryRun: true,
    schema: {
      input: {
        entityRef: z => z.string({ description: 'Component entity ref to retire' }),
        retention: z =>
          z.enum(['keep-last-backup', 'export-first', 'no-retention']),
        retentionDays: z => z.number().int().min(0).max(3650).optional(),
        reason: z => z.string().min(10).max(1500),
      },
      output: {
        body: z => z.string(),
        branchName: z => z.string(),
        title: z => z.string(),
        filesToDelete: z => z.array(z.string()),
        scope: z => z.string(),
        inventory: z => z.array(z.any()),
      },
    },
    async handler(ctx) {
      const { entityRef, retention, reason } = ctx.input;
      const ref = parseEntityRef(entityRef, { defaultKind: 'Component', defaultNamespace: 'default' });
      const entities = await deps.entities();
      const target = entities.find(
        e =>
          e.kind.toLowerCase() === ref.kind.toLowerCase() &&
          (e.metadata.namespace ?? 'default') === ref.namespace &&
          e.metadata.name === ref.name,
      );
      if (!target) throw new RetireInputError(`${entityRef} is not in the catalog`);
      const now = (deps.now ?? (() => new Date()))();
      const plan = await buildRetirePlan({
        target,
        catalog: entities,
        repo: await deps.repo(),
        cluster: await deps.cluster(),
        retention: { mode: retention, days: ctx.input.retentionDays ?? 30 },
        reason,
        now,
      });
      await writeWorkspace(ctx.workspacePath, plan.files);
      ctx.logger.info(
        `Retire plan for ${plan.ref}: ${plan.filesToDelete.length} file(s) to delete, ${plan.files.size} written`,
      );
      ctx.output('body', plan.body);
      ctx.output('branchName', `retire/${plan.name}-${now.toISOString().slice(0, 10)}`);
      ctx.output('title', `Retire ${plan.scope === 'app' ? 'app' : 'workload'}: ${plan.title}`);
      ctx.output('filesToDelete', plan.filesToDelete);
      ctx.output('scope', plan.scope);
      ctx.output('inventory', plan.inventory);
    },
  });
}

/** Empty strings from the template engine mean "not provided". */
export function normalizeAddInput(input: Record<string, unknown>): AddParams {
  const out: Record<string, unknown> = {};
  for (const [key, value] of Object.entries(input)) {
    if (value === '' || value === undefined || value === null) continue;
    out[key] = key === 'port' ? Number(value) : value;
  }
  return out as unknown as AddParams;
}

export function createAddPrepareAction(deps: LifecycleDeps) {
  return createTemplateAction({
    id: 'onemind:app-add:prepare',
    description:
      'Read-only: prepares the catalog input edits, the Argo application count and the STACK.md section for a new app, and the pull request body. Applies nothing.',
    supportsDryRun: true,
    schema: {
      input: {
        name: z => z.string(),
        title: z => z.string(),
        description: z => z.string(),
        hosting: z => z.enum(['kubernetes', 'digitalocean', 'saas', 'home-device']),
        pillar: z => z.enum(['so', 'le', 'hp', 'ge']),
        system: z => z.string(),
        owner: z => z.string(),
        // Template parameters that do not apply to a hosting type arrive as empty strings.
        lifecycle: z => z.any().optional(),
        url: z => z.any().optional(),
        image: z => z.any().optional(),
        port: z => z.any().optional(),
        access: z => z.any().optional(),
        dataPlan: z => z.any().optional(),
        healthCheck: z => z.any().optional(),
        healthUrl: z => z.any().optional(),
      },
      output: {
        body: z => z.string(),
        branchName: z => z.string(),
        title: z => z.string(),
      },
    },
    async handler(ctx) {
      const params = normalizeAddInput(ctx.input as Record<string, unknown>);
      const now = (deps.now ?? (() => new Date()))();
      const result = await prepareAdd(params, await deps.repo(), now);
      await writeWorkspace(ctx.workspacePath, result.files);
      for (const w of result.warnings) ctx.logger.warn(w);
      ctx.output('body', result.body);
      ctx.output('branchName', `add-app/${params.name}`);
      ctx.output('title', `Add app: ${params.title}`);
    },
  });
}

export const lifecycleScaffolderModule = createBackendModule({
  pluginId: 'scaffolder',
  moduleId: 'onemind-lifecycle',
  register(env) {
    env.registerInit({
      deps: {
        scaffolder: scaffolderActionsExtensionPoint,
        urlReader: coreServices.urlReader,
        auth: coreServices.auth,
        catalog: catalogServiceRef,
        logger: coreServices.logger,
      },
      async init({ scaffolder, urlReader, auth, catalog, logger }) {
        const deps = lifecycleDeps({ urlReader, auth, catalog, logger });
        scaffolder.addActions(
          createRetirePlanAction(deps),
          createAddPrepareAction(deps),
        );
      },
    });
  },
});

/** Wires the real sources. One repository snapshot per action run. */
export function lifecycleDeps(o: {
  urlReader: UrlReaderService;
  auth: AuthService;
  catalog: CatalogService;
  logger: { warn(message: string): void };
}): LifecycleDeps {
  return {
    repo: () => loadGithubSnapshot(o.urlReader),
    cluster: async () => {
      try {
        return await collectClusterFacts();
      } catch (error) {
        o.logger.warn(`Cluster facts unavailable: ${(error as Error).message}`);
        return undefined;
      }
    },
    entities: async () => {
      const { items } = await o.catalog.getEntities(
        { filter: [{ kind: 'Component' }, { kind: 'Resource' }] },
        { credentials: await o.auth.getOwnServiceCredentials() },
      );
      return items;
    },
  };
}
