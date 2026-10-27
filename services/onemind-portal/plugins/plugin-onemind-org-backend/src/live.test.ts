import {
  coreServices,
  createBackendModule,
  type SchedulerService,
} from '@backstage/backend-plugin-api';
import { catalogProcessingExtensionPoint } from '@backstage/plugin-catalog-node';
import {
  mockCredentials,
  mockServices,
  startTestBackend,
} from '@backstage/backend-test-utils';
import catalogPlugin from '@backstage/plugin-catalog-backend';
import type { Entity } from '@backstage/catalog-model';
import { readCurated } from './entities';
import { registerAuthentikOrg } from './module';
import { SCOPE, directory, fakeAuthentik, orgYamlPath, uidOf } from './testing';

jest.setTimeout(120_000);

let scheduler: SchedulerService;

/** The real module registration, with the Authentik HTTP read replaced by the fake directory. */
const withFakeAuthentik = createBackendModule({
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
        registerAuthentikOrg({
          ...deps,
          fetchJson: fakeAuthentik(directory).fetchJson,
        });
      },
    });
  },
});

const sleep = (ms: number) => new Promise(r => setTimeout(r, ms));
const memberOf = (e: Entity | undefined) =>
  (e?.relations ?? [])
    .filter(r => r.type === 'memberOf')
    .map(r => r.targetRef)
    .sort();

describe('Authentik org sync enabled, against a real catalog with the reviewed org file', () => {
  it('serves synced users the sign-in lookup finds exactly once, and leaves admin to the reviewed file', async () => {
    const [reviewedSubject] = [...(await readCurated(orgYamlPath)).subjects];
    const backend = await startTestBackend({
      features: [
        catalogPlugin,
        withFakeAuthentik,
        mockServices.scheduler.factory({ skipTaskRunOnStartup: true }),
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
              orgSync: {
                enabled: true,
                userGroups: SCOPE,
                curatedFile: orgYamlPath,
              },
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
    /** The lookup the OneMind sign-in resolver makes for a verified subject. */
    const signIn = (subject: string) =>
      list(`kind=user,metadata.annotations.onemind.io/oidc-sub=${subject}`);

    try {
      // The reviewed file lands first, as it does at startup in production.
      for (let i = 0; i < 120; i++) {
        if (memberOf((await list('kind=user'))[0]).length) break;
        await sleep(250);
      }
      await scheduler.triggerTask('onemind-authentik-org-refresh');
      let users: Entity[] = [];
      let groups: Entity[] = [];
      for (let i = 0; i < 240; i++) {
        users = await list('kind=user');
        groups = await list('kind=group');
        const stitched =
          users.length === 5 &&
          groups.length === 7 &&
          users.every(u => memberOf(u).length > 0) &&
          groups.some(g => (g.relations ?? []).some(r => r.type === 'childOf'));
        if (stitched) break;
        await sleep(250);
      }
      const byName = new Map(
        [...users, ...groups].map(e => [`${e.kind}:${e.metadata.name}`, e]),
      );

      // Each synced person resolves to exactly one catalog user.
      for (const [pk, name] of [
        [2, 'daniel-example.test'],
        [3, 'night.owl'],
        [8, 'two-parents-member'],
        [9, 'external.partner'],
      ] as const) {
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
      await backend.stop();
    }
  });
});
