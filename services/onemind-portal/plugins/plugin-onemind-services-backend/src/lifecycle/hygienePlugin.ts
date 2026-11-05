import express from 'express';
import {
  coreServices,
  DatabaseService,
  createBackendPlugin,
} from '@backstage/backend-plugin-api';
import { catalogServiceRef } from '@backstage/plugin-catalog-node';
import { AuthorizeResult } from '@backstage/plugin-permission-common';
import { taskCreatePermission } from '@backstage/plugin-scaffolder-common/alpha';
import { computeHygiene, HygieneReport, volsyncPvcsFrom } from './hygiene';
import { collectClusterFacts, loadGithubSnapshot } from './sources';

/**
 * Weekly stack hygiene report. A backend scheduled task computes it and stores the JSON;
 * the page reads the latest one. Reading needs a signed-in user (the same bar as the other
 * portal APIs). Refreshing on demand reuses `scaffolder.task.create`, which the portal
 * policy gives to administrators only, so no new permission exists.
 */

export interface ReportStore {
  latest(): Promise<HygieneReport | undefined>;
  save(report: HygieneReport): Promise<void>;
}

const TABLE = 'stack_hygiene_reports';
const KEEP = 12;

type Knex = Awaited<ReturnType<DatabaseService['getClient']>>;

export async function createKnexStore(knex: Knex): Promise<ReportStore> {
  if (!(await knex.schema.hasTable(TABLE))) {
    await knex.schema.createTable(TABLE, table => {
      table.increments('id');
      table.string('generated_at').notNullable();
      table.text('report').notNullable();
    });
  }
  return {
    async latest() {
      const row = await knex(TABLE).orderBy('id', 'desc').first();
      return row ? (JSON.parse(row.report) as HygieneReport) : undefined;
    },
    async save(report) {
      await knex(TABLE).insert({
        generated_at: report.generatedAt,
        report: JSON.stringify(report),
      });
      const old = await knex(TABLE).orderBy('id', 'desc').offset(KEEP).select('id');
      if (old.length) await knex(TABLE).whereIn('id', old.map(r => r.id)).delete();
    },
  };
}

export function memoryStore(): ReportStore {
  let current: HygieneReport | undefined;
  return {
    async latest() {
      return current;
    },
    async save(report) {
      current = report;
    },
  };
}

export const MIN_REFRESH_MS = 60_000;

export function hygieneRouter(o: {
  store: ReportStore;
  run: (trigger: 'schedule' | 'manual') => Promise<HygieneReport>;
  user: (req: express.Request) => Promise<unknown>;
  admin: (req: express.Request) => Promise<void>;
  now?: () => number;
}) {
  const router = express.Router();
  const now = o.now ?? Date.now;
  router.get('/stack-hygiene', async (req, res, next) => {
    try {
      await o.user(req);
      res.setHeader('Cache-Control', 'no-store');
      res.json({ report: (await o.store.latest()) ?? null });
    } catch (error) {
      next(error);
    }
  });
  router.post('/stack-hygiene/refresh', async (req, res, next) => {
    try {
      await o.admin(req);
      const latest = await o.store.latest();
      if (latest && now() - Date.parse(latest.generatedAt) < MIN_REFRESH_MS) {
        res.json({ report: latest, cached: true });
        return;
      }
      res.json({ report: await o.run('manual'), cached: false });
    } catch (error) {
      next(error);
    }
  });
  return router;
}

export const lifecyclePlugin = createBackendPlugin({
  pluginId: 'onemind-lifecycle',
  register(env) {
    env.registerInit({
      deps: {
        router: coreServices.httpRouter,
        httpAuth: coreServices.httpAuth,
        auth: coreServices.auth,
        permissions: coreServices.permissions,
        scheduler: coreServices.scheduler,
        database: coreServices.database,
        urlReader: coreServices.urlReader,
        logger: coreServices.logger,
        catalog: catalogServiceRef,
      },
      async init({ router, httpAuth, auth, permissions, scheduler, database, urlReader, logger, catalog }) {
        const store = await createKnexStore(await database.getClient());
        const run = async (trigger: 'schedule' | 'manual') => {
          const [cluster, entities, repo] = await Promise.all([
            collectClusterFacts().catch(error => {
              logger.warn(`Hygiene: cluster unreadable: ${(error as Error).message}`);
              return undefined;
            }),
            catalog
              .getEntities(
                { filter: [{ kind: 'Component' }] },
                { credentials: await auth.getOwnServiceCredentials() },
              )
              .then(r => r.items),
            loadGithubSnapshot(urlReader).catch(error => {
              logger.warn(`Hygiene: Git unreadable: ${(error as Error).message}`);
              return undefined;
            }),
          ]);
          let volsyncPvcs: Set<string> | undefined;
          if (repo) {
            volsyncPvcs = new Set();
            for (const [file, text] of await repo.readDir('platform')) {
              if (/\.ya?ml$/.test(file) && text.includes('kind: ReplicationSource'))
                for (const pvc of volsyncPvcsFrom(text)) volsyncPvcs.add(pvc);
            }
          }
          const report = computeHygiene({ now: new Date(), trigger, cluster, entities, volsyncPvcs });
          await store.save(report);
          return report;
        };
        await scheduler.scheduleTask({
          id: 'onemind-stack-hygiene-weekly',
          frequency: { cron: '0 6 * * 1' },
          timeout: { minutes: 5 },
          scope: 'global',
          fn: async () => {
            await run('schedule');
          },
        });
        router.use(
          hygieneRouter({
            store,
            run,
            user: req => httpAuth.credentials(req, { allow: ['user'] }),
            admin: async req => {
              const credentials = await httpAuth.credentials(req, { allow: ['user'] });
              const [decision] = await permissions.authorize(
                [{ permission: taskCreatePermission }],
                { credentials },
              );
              if (decision.result !== AuthorizeResult.ALLOW) {
                throw Object.assign(new Error('Administrators only'), {
                  name: 'NotAllowedError',
                });
              }
            },
          }),
        );
      },
    });
  },
});
