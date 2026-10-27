import type { Entity } from '@backstage/catalog-model';
import type { EntityProviderConnection } from '@backstage/plugin-catalog-node';
import { readAuthentikOrg } from './authentik';
import { readCurated } from './entities';
import {
  AuthentikOrgCleanup,
  AuthentikOrgProvider,
  type ProviderDeps,
} from './provider';
import {
  SCOPE,
  apiUser,
  directory,
  fakeAuthentik,
  orgYamlPath,
} from './testing';

type Mutation = {
  type: string;
  entities: { entity: Entity; locationKey?: string }[];
};

/** A provider wired to a fake Authentik and a "catalog" that records every mutation. */
function harness(
  data: { users: unknown[]; groups: unknown[] } = directory,
  overrides: Partial<ProviderDeps> = {},
) {
  const mutations: Mutation[] = [];
  const connection = {
    applyMutation: async (m: Mutation) => {
      mutations.push(m);
    },
    refresh: async () => {},
  } as unknown as EntityProviderConnection;
  const logger = { info: jest.fn(), warn: jest.fn() };
  let api = fakeAuthentik(data);
  const provider = new AuthentikOrgProvider({
    userGroups: SCOPE,
    curatedFile: 'catalog/org.yaml',
    logger,
    readSnapshot: signal =>
      readAuthentikOrg(api.fetchJson, 'http://authentik.test', signal),
    readCurated: () => readCurated(orgYamlPath),
    ...overrides,
  });
  return {
    provider,
    mutations,
    logger,
    connect: () => provider.connect(connection),
    setData: (next: typeof data) => (api = fakeAuthentik(next)),
  };
}

const refs = (m: Mutation) =>
  m.entities.map(d => `${d.entity.kind}:${d.entity.metadata.name}`);

describe('AuthentikOrgProvider', () => {
  it('applies one full mutation of the mapped entities under its own location key', async () => {
    const h = harness();
    await h.connect();
    await expect(h.provider.refresh()).resolves.toEqual({
      users: 4,
      groups: 4,
      curatedUsers: 1,
      curatedGroups: 1,
      superuserGroups: 2,
      multiParentGroups: 1,
    });
    expect(h.provider.getProviderName()).toBe('onemind-authentik-org');
    expect(h.mutations).toHaveLength(1);
    expect(h.mutations[0].type).toBe('full');
    expect(refs(h.mutations[0])).toEqual([
      'User:daniel-example.test',
      'User:external.partner',
      'User:night.owl',
      'User:two-parents-member',
      'Group:field-ops',
      'Group:field-ops-night-shift',
      'Group:onemind-portal-users',
      'Group:two-parents',
    ]);
    expect(
      h.mutations[0].entities.every(
        d => d.locationKey === 'provider:onemind-authentik-org',
      ),
    ).toBe(true);
    expect(h.logger.warn).toHaveBeenCalledTimes(1);
    expect(h.logger.info).toHaveBeenCalledWith(
      'Authentik org sync: 4 users and 4 groups published; 1 users and 1 groups left to catalog/org.yaml; ' +
        '2 superuser groups not published; 1 groups published without a parent (several Authentik parents)',
    );
  });

  it('drops a user at the next run once they leave the scope group', async () => {
    const h = harness();
    await h.connect();
    await h.provider.refresh();
    h.setData({
      groups: directory.groups,
      users: directory.users.map(u => (u.pk === 3 ? { ...u, groups: [] } : u)),
    });
    await h.provider.refresh();
    expect(refs(h.mutations[1])).not.toContain('User:night.owl');
    expect(refs(h.mutations[1])).toContain('User:daniel-example.test');
  });

  it.each([
    [
      'no users (the token cannot view users)',
      { users: [], groups: directory.groups },
      'Authentik returned no active users; check that the token may view users. Nothing was changed',
    ],
    [
      'service accounts only',
      {
        users: [apiUser(4, 'ak-outpost', [1], { type: 'service_account' })],
        groups: directory.groups,
      },
      'Authentik returned no active users',
    ],
    [
      'no groups (the token cannot view groups)',
      { users: directory.users, groups: [] },
      'Authentik returned no groups; check that the token may view groups. Nothing was changed',
    ],
  ])(
    'changes nothing when Authentik returns %s',
    async (_name, data, message) => {
      const h = harness(data);
      await h.connect();
      await expect(h.provider.refresh()).rejects.toThrow(message);
      expect(h.mutations).toHaveLength(0);
    },
  );

  it('changes nothing when Authentik fails', async () => {
    const h = harness(directory, {
      readSnapshot: async () => {
        throw new Error('Authentik API answered 503 for /api/v3/core/users/');
      },
    });
    await h.connect();
    await expect(h.provider.refresh()).rejects.toThrow('answered 503');
    expect(h.mutations).toHaveLength(0);
  });

  it('changes nothing, and does not call Authentik, when the reviewed file is unreadable', async () => {
    const readSnapshot = jest.fn();
    const h = harness(directory, {
      readSnapshot,
      readCurated: () => readCurated('/nonexistent/catalog/org.yaml'),
    });
    await h.connect();
    await expect(h.provider.refresh()).rejects.toThrow(/ENOENT/);
    expect(readSnapshot).not.toHaveBeenCalled();
    expect(h.mutations).toHaveLength(0);
  });

  it('changes nothing after its task was aborted', async () => {
    const h = harness();
    await h.connect();
    const controller = new AbortController();
    controller.abort();
    await expect(h.provider.refresh(controller.signal)).rejects.toThrow();
    expect(h.mutations).toHaveLength(0);
  });

  it('refuses to run before the catalog connected it', async () => {
    await expect(harness().provider.refresh()).rejects.toThrow(
      'Authentik org provider is not connected',
    );
  });
});

describe('AuthentikOrgCleanup', () => {
  const EMPTY = { type: 'full', entities: [] };

  it('removes everything the sync published when connected and on every run, under the same name', async () => {
    const mutations: Mutation[] = [];
    const logger = { warn: jest.fn() };
    const cleanup = new AuthentikOrgCleanup(logger);
    expect(cleanup.getProviderName()).toBe(
      harness().provider.getProviderName(),
    );
    await cleanup.connect({
      applyMutation: async (m: Mutation) => {
        mutations.push(m);
      },
      refresh: async () => {},
    } as unknown as EntityProviderConnection);
    await cleanup.clear();
    expect(mutations).toEqual([EMPTY, EMPTY]);
    expect(logger.warn).not.toHaveBeenCalled();
  });

  it('lets the catalog start when the removal at startup fails, and the next run retries it', async () => {
    const logger = { warn: jest.fn() };
    const cleanup = new AuthentikOrgCleanup(logger);
    const applyMutation = jest
      .fn()
      .mockRejectedValueOnce(new Error('database is locked'))
      .mockResolvedValue(undefined);
    await expect(
      cleanup.connect({
        applyMutation,
        refresh: jest.fn(),
      } as unknown as EntityProviderConnection),
    ).resolves.toBeUndefined();
    expect(logger.warn).toHaveBeenCalledWith(
      expect.stringContaining('database is locked'),
    );
    await cleanup.clear();
    expect(applyMutation.mock.calls).toEqual([[EMPTY], [EMPTY]]);
  });

  it('refuses to run before the catalog connected it', async () => {
    await expect(
      new AuthentikOrgCleanup({ warn: jest.fn() }).clear(),
    ).rejects.toThrow('Authentik org provider is not connected');
  });
});
