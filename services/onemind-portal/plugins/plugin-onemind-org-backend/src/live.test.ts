import {
  coreServices,
  createBackendModule,
  type SchedulerService,
} from '@backstage/backend-plugin-api';
import { catalogProcessingExtensionPoint } from '@backstage/plugin-catalog-node';
import {
  TestDatabases,
  mockCredentials,
  mockServices,
  startTestBackend,
} from '@backstage/backend-test-utils';
import catalogPlugin from '@backstage/plugin-catalog-backend';
import type { Entity } from '@backstage/catalog-model';
import type { JsonObject } from '@backstage/types';
import { readCurated } from './entities';
import { TASK_ID, registerAuthentikOrg } from './module';
import { SCOPE, directory, fakeAuthentik, orgYamlPath, uidOf } from './testing';

jest.setTimeout(120_000);

type Database = Awaited<ReturnType<TestDatabases['init']>>;

const sleep = (ms: number) => new Promise(r => setTimeout(r, ms));
const memberOf = (e: Entity | undefined) =>
  (e?.relations ?? [])
    .filter(r => r.type === 'memberOf')
    .map(r => r.targetRef)
    .sort();
/** The catalog processes and stitches asynchronously, so assertions on its API poll. */
async function eventually(what: string, check: () => Promise<boolean>) {
  for (let i = 0; i < 240; i++) {
    if (await check()) return;
    await sleep(250);
  }
  throw new Error(`timed out waiting until ${what}`);
}
/** The people the fake directory must publish: [Authentik pk, entity name]. */
const SYNCED = [
  [2, 'daniel-example.test'],
  [3, 'night.owl'],
  [8, 'two-parents-member'],
  [9, 'external.partner'],
] as const;

/**
 * One portal replica: the catalog with the reviewed org file and the real module registration, with
 * the Authentik HTTP read replaced by the fake directory. Replicas given the same database share it,
 * as the production replicas share PostgreSQL; each replica's task runs only when the test makes that
 * replica claim it.
 */
async function startReplica(orgSync: JsonObject, database?: Database) {
  let scheduler: SchedulerService | undefined;
  const module = createBackendModule({
    pluginId: 'catalog',
    moduleId: 'test-authentik-org',
    register(reg) {
      reg.registerInit({
        deps: {
          config: coreServices.rootConfig,
          logger: coreServices.logger,
          scheduler: coreServices.scheduler,
          processing: catalogProcessingExtensionPoint,
        },
        async init(deps) {
          scheduler = deps.scheduler;
          await registerAuthentikOrg({
            ...deps,
            fetchJson: fakeAuthentik(directory).fetchJson,
          });
        },
      });
    },
  });
  const backend = await startTestBackend({
    features: [
      catalogPlugin,
      module,
      mockServices.scheduler.factory({ skipTaskRunOnStartup: true }),
      ...(database ? [mockServices.database.factory({ knex: database })] : []),
      mockServices.rootConfig.factory({
        data: {
          catalog: {
            processingInterval: { seconds: 2 },
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
            // The module derives its reviewed org file from this location.
            locations: [
              {
                type: 'file',
                target: orgYamlPath,
                rules: [{ allow: ['User', 'Group'] }],
              },
            ],
          },
          authentik: {
            baseUrl: 'http://authentik.test',
            token: 'fixture-only-token',
            orgSync: { userGroups: SCOPE, ...orgSync },
          },
        },
      }),
    ],
  });
  const base = `http://localhost:${backend.server.port()}/api/catalog`;
  const headers = { Authorization: mockCredentials.user.header() };
  const list = async (filter: string) =>
    (await (
      await fetch(`${base}/entities?${new URLSearchParams({ filter })}`, {
        headers,
      })
    ).json()) as Entity[];
  let stopped = false;
  return {
    list,
    /** The lookup the OneMind sign-in resolver makes for a verified subject. */
    signIn: (subject: string) =>
      list(`kind=user,metadata.annotations.onemind.io/oidc-sub=${subject}`),
    /** Everything the module published: the check docs/catalog/authentik-org-sync.md asks for. */
    published: () =>
      list(
        'metadata.annotations.backstage.io/managed-by-location=provider:onemind-authentik-org',
      ),
    /** This replica claims the next run of the shared global task. */
    runTask: async () => {
      if (!scheduler) throw new Error('the module was not initialised');
      await scheduler.triggerTask(TASK_ID);
    },
    stop: async () => {
      if (stopped) return;
      stopped = true;
      await backend.stop();
    },
  };
}

describe('Authentik org sync enabled, against a real catalog with the reviewed org file', () => {
  it('serves synced users the sign-in lookup finds exactly once, and leaves admin to the reviewed file', async () => {
    const [reviewedSubject] = [...(await readCurated(orgYamlPath)).subjects];
    const replica = await startReplica({ enabled: true });
    const { list, signIn } = replica;

    try {
      // The reviewed file lands first, as it does at startup in production.
      for (let i = 0; i < 120; i++) {
        if (memberOf((await list('kind=user'))[0]).length) break;
        await sleep(250);
      }
      await replica.runTask();
      let users: Entity[] = [];
      let groups: Entity[] = [];
      for (let i = 0; i < 240; i++) {
        users = await list('kind=user');
        groups = await list('kind=group');
        const stitched =
          users.length === 5 &&
          groups.length === 5 &&
          users.every(u => memberOf(u).length > 0) &&
          groups.some(g => (g.relations ?? []).some(r => r.type === 'childOf'));
        if (stitched) break;
        await sleep(250);
      }
      const byName = new Map(
        [...users, ...groups].map(e => [`${e.kind}:${e.metadata.name}`, e]),
      );

      // Each synced person resolves to exactly one catalog user.
      for (const [pk, name] of SYNCED) {
        const found = await signIn(uidOf(pk));
        expect(found.map(e => e.metadata.name)).toEqual([name]);
        expect(
          found[0].metadata.annotations?.['backstage.io/managed-by-location'],
        ).toBe('provider:onemind-authentik-org');
      }
      // Excluded people (service account, deactivated, out of scope) cannot sign in.
      for (const pk of [4, 5, 6, 7])
        expect(await signIn(uidOf(pk))).toHaveLength(0);

      // The reviewed administrator is still the file's entity, found exactly once, still admin.
      const admin = await signIn(reviewedSubject);
      expect(admin.map(e => e.metadata.name)).toEqual(['akadmin']);
      expect(
        admin[0].metadata.annotations?.['backstage.io/managed-by-location'],
      ).toBe(`file:${orgYamlPath}`);
      expect(memberOf(admin[0])).toEqual(['group:default/onemind-admins']);

      // Nobody else becomes a member of the admin group.
      for (const u of users.filter(e => e.metadata.name !== 'akadmin'))
        expect(memberOf(u)).not.toContain('group:default/onemind-admins');
      const adminGroup = byName.get('Group:onemind-admins');
      expect(
        adminGroup?.metadata.annotations?.['backstage.io/managed-by-location'],
      ).toBe(`file:${orgYamlPath}`);
      expect(
        (adminGroup?.relations ?? [])
          .filter(r => r.type === 'hasMember')
          .map(r => r.targetRef),
      ).toEqual(['user:default/akadmin']);

      // Only the scope's groups are in the catalog: no group outside it (authentik Admins,
      // argocd-admins, TAK Users) and no superuser group, so nobody can read who is in them.
      expect(groups.map(g => g.metadata.name).sort()).toEqual([
        'field-ops',
        'field-ops-night-shift',
        'onemind-admins',
        'onemind-portal-users',
        'two-parents',
      ]);
      expect(memberOf(byName.get('User:daniel-example.test'))).toEqual([
        'group:default/onemind-portal-users',
      ]);

      // Membership and hierarchy relations are stitched like the reviewed file's.
      expect(memberOf(byName.get('User:night.owl'))).toEqual([
        'group:default/field-ops-night-shift',
      ]);
      expect(
        (byName.get('Group:field-ops-night-shift')?.relations ?? [])
          .filter(r => r.type === 'childOf')
          .map(r => r.targetRef),
      ).toEqual(['group:default/field-ops']);
    } finally {
      await replica.stop();
    }
  });
});

describe('turning the sync off during a rolling update, against one catalog database shared by an old and a new replica', () => {
  const databases = TestDatabases.create({ ids: ['SQLITE_3'] });

  it('revokes every synced user at the first run after the last enabled replica is gone, even if it ran the task after the new replica started', async () => {
    const database = await databases.init('SQLITE_3');
    const [reviewedSubject] = [...(await readCurated(orgYamlPath)).subjects];
    const old = await startReplica({ enabled: true }, database);
    let fresh: Awaited<ReturnType<typeof startReplica>> | undefined;
    try {
      await eventually(
        'the reviewed org file is stitched',
        async () => memberOf((await old.list('kind=user'))[0]).length > 0,
      );
      await old.runTask();
      await eventually(
        'the enabled replica published every synced person',
        async () =>
          (await old.published()).filter(e => e.kind === 'User').length ===
          SYNCED.length,
      );

      // The new, disabled replica removes them as soon as its catalog connects the provider.
      fresh = await startReplica({ enabled: false }, database);
      const replica = fresh;
      await eventually(
        'the new replica removed them at startup',
        async () => (await replica.published()).length === 0,
      );

      // The old replica claims the shared task once more before it stops: the rollout race.
      await old.runTask();
      await eventually(
        'the old replica republished them',
        async () => (await replica.signIn(uidOf(3))).length === 1,
      );
      await old.stop();

      // The next run of the same task can only happen on a new replica now, and removes them.
      await replica.runTask();
      await eventually(
        'the next run removed them again',
        async () => (await replica.published()).length === 0,
      );
      for (const [pk] of SYNCED)
        expect(await replica.signIn(uidOf(pk))).toHaveLength(0);

      // The reviewed administrator is untouched throughout.
      const admin = await replica.signIn(reviewedSubject);
      expect(admin.map(e => e.metadata.name)).toEqual(['akadmin']);
      expect(memberOf(admin[0])).toEqual(['group:default/onemind-admins']);
    } finally {
      await fresh?.stop();
      await old.stop();
    }
  }, 240_000);
});
