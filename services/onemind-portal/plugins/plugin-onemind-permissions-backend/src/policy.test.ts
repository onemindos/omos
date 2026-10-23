import {
  AuthorizeResult,
  createPermission,
  type Permission,
} from '@backstage/plugin-permission-common';
import type { PolicyQueryUser } from '@backstage/plugin-permission-node';
import {
  catalogEntityCreatePermission,
  catalogEntityDeletePermission,
  catalogEntityReadPermission,
  catalogEntityRefreshPermission,
  catalogLocationCreatePermission,
  catalogLocationDeletePermission,
  catalogLocationReadPermission,
} from '@backstage/plugin-catalog-common/alpha';
import {
  actionExecutePermission,
  taskCancelPermission,
  taskCreatePermission,
  taskReadPermission,
  templateParameterReadPermission,
  templateStepReadPermission,
} from '@backstage/plugin-scaffolder-common/alpha';
import { createOneMindPolicy, isOneMindAdmin } from './policy';

const user = (admin: boolean): PolicyQueryUser => ({
  credentials: {
    $$type: '@backstage/BackstageCredentials',
    principal: { type: 'user', userEntityRef: 'user:default/operator' },
  },
  info: {
    userEntityRef: 'user:default/operator',
    ownershipEntityRefs: admin
      ? ['user:default/operator', 'group:default/onemind-admins']
      : ['user:default/operator', 'group:default/some-other-group'],
  },
});
const decide = (permission: Permission, who?: PolicyQueryUser) =>
  createOneMindPolicy().handle({ permission }, who);

// Every permission the portal's plugins can ask for that touches the catalog or the
// scaffolder, with the action Backstage declares for it.
const catalogAndScaffolder: [string, Permission][] = [
  ['catalog.entity.read', catalogEntityReadPermission],
  ['catalog.entity.create', catalogEntityCreatePermission],
  ['catalog.entity.refresh', catalogEntityRefreshPermission],
  ['catalog.entity.delete', catalogEntityDeletePermission],
  ['catalog.location.read', catalogLocationReadPermission],
  ['catalog.location.create', catalogLocationCreatePermission],
  ['catalog.location.delete', catalogLocationDeletePermission],
  ['scaffolder.task.create', taskCreatePermission],
  ['scaffolder.task.read', taskReadPermission],
  ['scaffolder.task.cancel', taskCancelPermission],
  ['scaffolder.action.execute', actionExecutePermission],
  ['scaffolder.template.parameter.read', templateParameterReadPermission],
  ['scaffolder.template.step.read', templateStepReadPermission],
];

describe('OneMind permission policy', () => {
  test.each(catalogAndScaffolder)(
    'administrator is allowed: %s',
    async (_name, permission) => {
      expect(await decide(permission, user(true))).toEqual({
        result: AuthorizeResult.ALLOW,
      });
    },
  );

  test.each(catalogAndScaffolder)(
    'other signed-in user: read allowed, everything else denied: %s',
    async (_name, permission) => {
      expect(await decide(permission, user(false))).toEqual({
        result:
          permission.attributes.action === 'read'
            ? AuthorizeResult.ALLOW
            : AuthorizeResult.DENY,
      });
    },
  );

  test.each(catalogAndScaffolder)(
    'anonymous is denied, including reads: %s',
    async (_name, permission) => {
      expect(await decide(permission, undefined)).toEqual({
        result: AuthorizeResult.DENY,
      });
    },
  );

  test('the effective table the portal relies on', async () => {
    const table = async (admin: boolean) =>
      Object.fromEntries(
        await Promise.all(
          catalogAndScaffolder.map(async ([name, p]) => [
            name,
            (await decide(p, user(admin))).result,
          ]),
        ),
      );
    const admin = await table(true);
    expect(Object.values(admin).every(r => r === AuthorizeResult.ALLOW)).toBe(
      true,
    );
    const reader = await table(false);
    expect(
      Object.entries(reader)
        .filter(([, r]) => r === AuthorizeResult.ALLOW)
        .map(([n]) => n)
        .sort(),
    ).toEqual(
      [
        'catalog.entity.read',
        'catalog.location.read',
        'scaffolder.task.read',
        'scaffolder.template.parameter.read',
        'scaffolder.template.step.read',
      ].sort(),
    );
  });

  test('decisions are never conditional, so nothing is filtered per entity', async () => {
    for (const [, permission] of catalogAndScaffolder) {
      for (const who of [user(true), user(false), undefined]) {
        expect((await decide(permission, who)).result).not.toBe(
          AuthorizeResult.CONDITIONAL,
        );
      }
    }
  });

  test('unrelated host permissions follow the same rule', async () => {
    const read = createPermission({
      name: 'onemind.read',
      attributes: { action: 'read' },
    });
    const write = createPermission({
      name: 'onemind.write',
      attributes: { action: 'create' },
    });
    expect((await decide(read, user(false))).result).toBe(
      AuthorizeResult.ALLOW,
    );
    expect((await decide(write, user(false))).result).toBe(
      AuthorizeResult.DENY,
    );
    expect((await decide(write, user(true))).result).toBe(
      AuthorizeResult.ALLOW,
    );
  });

  test('admin is exactly membership of group:default/onemind-admins', () => {
    expect(isOneMindAdmin(user(true))).toBe(true);
    expect(isOneMindAdmin(user(false))).toBe(false);
    expect(isOneMindAdmin(undefined)).toBe(false);
    expect(
      isOneMindAdmin({
        info: {
          userEntityRef: 'user:default/x',
          ownershipEntityRefs: ['group:default/onemind-admins-extra'],
        },
      }),
    ).toBe(false);
  });
});
