import {
  FieldFormatEntityPolicy,
  KubernetesValidatorFunctions,
  NoForeignRootFieldsEntityPolicy,
  SchemaValidEntityPolicy,
  groupEntityV1alpha1Validator,
  userEntityV1alpha1Validator,
  type Entity,
} from '@backstage/catalog-model';
import { parseGroup, parseUser, type AuthentikSnapshot } from './authentik';
import {
  A,
  LOCATION,
  buildOrgEntities,
  entityName,
  readCurated,
  type Curated,
} from './entities';
import {
  SCOPE,
  apiGroup,
  apiUser,
  directory,
  orgYamlPath,
  uidOf,
} from './testing';

const snapshot = (
  data: { users: unknown[]; groups: unknown[] } = directory,
): AuthentikSnapshot => ({
  users: data.users.map(u => parseUser(u as Record<string, unknown>)),
  groups: data.groups.map(g => parseGroup(g as Record<string, unknown>)),
});
const none: Curated = {
  users: new Set(),
  groups: new Set(),
  subjects: new Set(),
};
const names = (entities: Entity[], kind: string) =>
  entities.filter(e => e.kind === kind).map(e => e.metadata.name);
const find = (entities: Entity[], kind: string, name: string) =>
  entities.find(e => e.kind === kind && e.metadata.name === name)!;
/** The catalog lookup the sign-in resolver makes: kind user, exact annotation value. */
const signInMatches = (entities: Entity[], subject: string) =>
  entities.filter(
    e => e.kind === 'User' && e.metadata.annotations?.[A.oidcSub] === subject,
  );

describe('entityName', () => {
  it.each([
    ['akadmin', 'akadmin'],
    ['daniel@example.test', 'daniel-example.test'],
    ['Daniel@Example.test', 'daniel-example.test'],
    ['OneMind Admins', 'onemind-admins'],
    ['authentik Admins', 'authentik-admins'],
    ['Field Ops / Night Shift', 'field-ops-night-shift'],
    ['  -José Núñez_ ', 'jose-nunez'],
    ['first.last', 'first.last'],
    ['a__b', 'a__b'],
  ])('%p -> %p', (raw, expected) => {
    expect(entityName(raw)).toBe(expected);
  });

  it('cuts long names to 63 characters with a stable hash suffix', () => {
    const long = `${'x'.repeat(70)}@example.test`;
    const name = entityName(long);
    expect(name).toHaveLength(63);
    expect(name).toMatch(/^x{54}-[0-9a-f]{8}$/);
    expect(entityName(long)).toBe(name);
    expect(entityName(`${'x'.repeat(70)}@example.org`)).not.toBe(name);
  });

  it('gives names without a usable character a hash name instead of nothing', () => {
    expect(entityName('日本語')).toMatch(/^authentik-[0-9a-f]{8}$/);
    expect(entityName('---')).toMatch(/^authentik-[0-9a-f]{8}$/);
    expect(entityName('日本語')).not.toBe(entityName('中文'));
  });

  it('always returns a valid catalog entity name', () => {
    const samples = [
      'Ünïcödé Group (Readers)',
      'x.'.repeat(40),
      `${'a-'.repeat(27)}-.-tail`,
      '@@@user@@@',
      'UPPER_lower.123',
      ' trailing space ',
      '🙂 smile',
      ...directory.groups.map(g => g.name),
      ...directory.users.map(u => u.username),
    ];
    for (const raw of samples)
      expect([
        raw,
        KubernetesValidatorFunctions.isValidObjectName(entityName(raw)),
      ]).toEqual([raw, true]);
  });
});

describe('readCurated', () => {
  it('reads the reviewed users, groups and subjects of catalog/org.yaml', async () => {
    const curated = await readCurated(orgYamlPath);
    expect([...curated.users]).toEqual(['akadmin']);
    expect([...curated.groups]).toEqual(['onemind-admins']);
    expect(curated.subjects.size).toBe(1);
    expect([...curated.subjects][0]).toMatch(/^[0-9a-f]{64}$/);
  });
});

describe('buildOrgEntities', () => {
  let curated: Curated;
  let result: ReturnType<typeof buildOrgEntities>;
  beforeAll(async () => {
    curated = await readCurated(orgYamlPath);
    result = buildOrgEntities(snapshot(), { userGroups: SCOPE, curated });
  });

  it('publishes active people in the scope group and every group below it, nothing else', () => {
    expect(names(result.entities, 'User')).toEqual([
      'daniel-example.test',
      'external.partner',
      'night.owl',
      'two-parents-member',
    ]);
    // Not published: akadmin (reviewed), a service account, a deactivated user, a member of an
    // unrelated group only, and a username that maps to an already used entity name.
    expect(result.users).toBe(4);
    expect(result.curatedUsers).toBe(1);
  });

  it('sets onemind.io/oidc-sub to the exact Authentik uid, so the sign-in resolver finds exactly one user', () => {
    for (const [pk, name] of [
      [2, 'daniel-example.test'],
      [3, 'night.owl'],
      [8, 'two-parents-member'],
      [9, 'external.partner'],
    ] as const) {
      const matches = signInMatches(result.entities, uidOf(pk));
      expect(matches.map(e => e.metadata.name)).toEqual([name]);
    }
    // Excluded people have no catalog identity, so their sign-in fails closed.
    for (const pk of [1, 4, 5, 6, 7])
      expect(signInMatches(result.entities, uidOf(pk))).toHaveLength(0);
  });

  it('maps profile and provenance like the reviewed org file does', () => {
    expect(find(result.entities, 'User', 'daniel-example.test')).toEqual({
      apiVersion: 'backstage.io/v1alpha1',
      kind: 'User',
      metadata: {
        name: 'daniel-example.test',
        annotations: {
          'backstage.io/managed-by-location': LOCATION,
          'backstage.io/managed-by-origin-location': LOCATION,
          'onemind.io/oidc-sub': uidOf(2),
          'onemind.io/authentik-name': 'daniel@example.test',
          'onemind.io/authentik-uuid': '10000000-0000-4000-8000-000000000002',
        },
      },
      spec: {
        profile: {
          displayName: 'Daniel Example',
          email: 'person2@example.test',
        },
        memberOf: ['onemind-portal-users'],
      },
    });
    // An empty name falls back to the username; an empty email is left out.
    expect(
      find(result.entities, 'User', 'external.partner').spec?.profile,
    ).toEqual({
      displayName: 'external.partner',
      email: 'person9@example.test',
    });
    expect(
      find(result.entities, 'User', 'two-parents-member').spec?.profile,
    ).toEqual({ displayName: 'Person 8' });
  });

  it('publishes every Authentik group except one the reviewed file defines', () => {
    expect(names(result.entities, 'Group')).toEqual([
      'authentik-admins',
      'field-ops',
      'field-ops-night-shift',
      'onemind-portal-users',
      'tak-users',
      'two-parents',
    ]);
    expect(result.curatedGroups).toBe(1);
    expect(find(result.entities, 'Group', 'field-ops-night-shift')).toEqual({
      apiVersion: 'backstage.io/v1alpha1',
      kind: 'Group',
      metadata: {
        name: 'field-ops-night-shift',
        annotations: {
          'backstage.io/managed-by-location': LOCATION,
          'backstage.io/managed-by-origin-location': LOCATION,
          'onemind.io/authentik-name': 'Field Ops / Night Shift',
          'onemind.io/authentik-uuid': '00000000-0000-4000-8000-000000000003',
        },
      },
      spec: {
        type: 'team',
        profile: { displayName: 'Field Ops / Night Shift' },
        parent: 'field-ops',
        children: [],
        members: ['night.owl'],
      },
    });
  });

  it('never grants or links reviewed identities: admin membership stays in catalog/org.yaml', () => {
    const text = JSON.stringify(result.entities);
    expect(text).not.toContain('onemind-admins');
    expect(text).not.toMatch(/"akadmin"/);
    for (const e of result.entities.filter(x => x.kind === 'User'))
      expect(e.spec?.memberOf).not.toContain('onemind-admins');
    // Even an Authentik group that sanitises to the reviewed admin group name is not published.
    expect(names(result.entities, 'Group')).not.toContain('onemind-admins');
  });

  it('copies the hierarchy only where a group has exactly one Authentik parent', () => {
    expect(
      find(result.entities, 'Group', 'onemind-portal-users').spec,
    ).toMatchObject({
      children: ['field-ops'],
      members: ['daniel-example.test'],
    });
    expect(find(result.entities, 'Group', 'field-ops').spec).toMatchObject({
      parent: 'onemind-portal-users',
      children: ['field-ops-night-shift'],
      members: ['external.partner'],
    });
    const two = find(result.entities, 'Group', 'two-parents');
    expect(two.spec?.parent).toBeUndefined();
    expect(two.spec?.members).toEqual(['two-parents-member']);
    expect(result.multiParentGroups).toBe(1);
    // Members of groups outside the scope, and reviewed users, are not listed.
    expect(find(result.entities, 'Group', 'tak-users').spec?.members).toEqual(
      [],
    );
    expect(
      find(result.entities, 'Group', 'authentik-admins').spec?.members,
    ).toEqual([]);
  });

  it('every reference points at a published entity, and both sides of each link agree', () => {
    const users = names(result.entities, 'User');
    const groups = names(result.entities, 'Group');
    const userEntities = result.entities.filter(e => e.kind === 'User');
    const groupEntities = result.entities.filter(e => e.kind === 'Group');
    const list = (e: Entity, key: string) => (e.spec?.[key] ?? []) as string[];
    const parentOf = (g: Entity) =>
      g.spec?.parent ? [g.spec.parent as string] : [];
    const pairs = (p: string[][]) => p.map(x => x.join(' -> ')).sort();

    expect(
      userEntities
        .flatMap(u => list(u, 'memberOf'))
        .filter(g => !groups.includes(g)),
    ).toEqual([]);
    expect(
      groupEntities
        .flatMap(g => list(g, 'members'))
        .filter(u => !users.includes(u)),
    ).toEqual([]);
    expect(
      groupEntities
        .flatMap(g => [...list(g, 'children'), ...parentOf(g)])
        .filter(g => !groups.includes(g)),
    ).toEqual([]);
    // user.memberOf matches group.members, and group.parent matches the parent's children.
    expect(
      pairs(
        userEntities.flatMap(u =>
          list(u, 'memberOf').map(g => [u.metadata.name, g]),
        ),
      ),
    ).toEqual(
      pairs(
        groupEntities.flatMap(g =>
          list(g, 'members').map(u => [u, g.metadata.name]),
        ),
      ),
    );
    expect(
      pairs(
        groupEntities.flatMap(g => parentOf(g).map(p => [g.metadata.name, p])),
      ),
    ).toEqual(
      pairs(
        groupEntities.flatMap(g =>
          list(g, 'children').map(c => [c, g.metadata.name]),
        ),
      ),
    );
  });

  it('passes the catalog validators and entity policies', async () => {
    const policies = [
      new SchemaValidEntityPolicy(),
      new FieldFormatEntityPolicy(),
      new NoForeignRootFieldsEntityPolicy(),
    ];
    for (const e of result.entities) {
      const validator =
        e.kind === 'User'
          ? userEntityV1alpha1Validator
          : groupEntityV1alpha1Validator;
      await expect(validator.check(e)).resolves.toBe(true);
      for (const policy of policies)
        await expect(policy.enforce(e)).resolves.toBe(e);
    }
  });

  it('skips the newer of two users that map to one entity name, with a warning', () => {
    expect(
      find(result.entities, 'User', 'daniel-example.test').metadata
        .annotations?.[A.oidcSub],
    ).toBe(uidOf(2));
    expect(result.warnings).toEqual([
      'user "Daniel@Example.test" maps to the entity name daniel-example.test, already used by user "daniel@example.test"; skipped',
    ]);
  });

  it('never publishes a user whose subject the reviewed file already holds, under any name', () => {
    const [reviewed] = [...curated.subjects];
    const data = {
      groups: directory.groups,
      users: [
        apiUser(20, 'renamed-admin', [1], { uid: reviewed }),
        apiUser(21, 'fresh', [1]),
      ],
    };
    const out = buildOrgEntities(snapshot(data), {
      userGroups: SCOPE,
      curated,
    });
    expect(names(out.entities, 'User')).toEqual(['fresh']);
    expect(signInMatches(out.entities, reviewed)).toHaveLength(0);
    expect(out.curatedUsers).toBe(1);
  });

  it('keeps the first group by name when two group names map to one entity name', () => {
    const data = {
      groups: [
        apiGroup(1, 'Portal'),
        apiGroup(2, 'portal'),
        apiGroup(3, 'PORTAL'),
      ],
      users: [apiUser(2, 'member', [2])],
    };
    const out = buildOrgEntities(snapshot(data), {
      userGroups: ['portal'],
      curated: none,
    });
    expect(out.entities.filter(e => e.kind === 'Group')).toHaveLength(1);
    expect(
      find(out.entities, 'Group', 'portal').metadata.annotations?.[A.name],
    ).toBe('PORTAL');
    // The member is in scope through the skipped group "portal", but only published groups are referenced.
    expect(find(out.entities, 'User', 'member').spec?.memberOf).toEqual([]);
    expect(out.warnings).toHaveLength(2);
  });

  it('is deterministic whatever order Authentik lists records in', () => {
    const reversed = {
      users: [...directory.users].reverse(),
      groups: [...directory.groups].reverse(),
    };
    expect(
      buildOrgEntities(snapshot(reversed), { userGroups: SCOPE, curated }),
    ).toEqual(result);
  });

  it('refuses to publish when no configured scope group exists, and warns about a missing one', () => {
    expect(() =>
      buildOrgEntities(snapshot(), {
        userGroups: ['Portal Users (typo)'],
        curated,
      }),
    ).toThrow(
      /None of the groups in authentik\.orgSync\.userGroups .* exists in Authentik/,
    );
    const out = buildOrgEntities(snapshot(), {
      userGroups: ['Portal Users (typo)', 'Field Ops'],
      curated,
    });
    expect(names(out.entities, 'User')).toEqual([
      'external.partner',
      'night.owl',
    ]);
    expect(out.warnings).toContain(
      'authentik.orgSync.userGroups names groups that do not exist in Authentik: Portal Users (typo)',
    );
  });

  it('never publishes a parent cycle', () => {
    const data = {
      groups: [
        apiGroup(1, 'A', [2]),
        apiGroup(2, 'B', [1]),
        apiGroup(3, 'Self', [3]),
      ],
      users: [apiUser(2, 'member', [1])],
    };
    const out = buildOrgEntities(snapshot(data), {
      userGroups: ['A'],
      curated: none,
    });
    expect(find(out.entities, 'Group', 'a').spec?.parent).toBeUndefined();
    expect(find(out.entities, 'Group', 'b').spec?.parent).toBe('a');
    expect(find(out.entities, 'Group', 'a').spec?.children).toEqual(['b']);
    expect(find(out.entities, 'Group', 'self').spec?.parent).toBeUndefined();
    expect(out.warnings).toHaveLength(2);
    // Scope traversal terminates on the cycle too.
    expect(names(out.entities, 'User')).toEqual(['member']);
  });
});
