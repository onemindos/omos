import knex from 'knex';
import type { Knex } from 'knex';
import {
  AuthzProfileCache,
  AuthzProfileService,
  OpenChoreoPermissionPolicy,
} from '@openchoreo/backstage-plugin-permission-backend-module-openchoreo-policy';
import type { LoggerService } from '@backstage/backend-plugin-api';
import {
  AuthorizeResult,
  createPermission,
} from '@backstage/plugin-permission-common';
import { createOneMindPolicy } from './policy';
import { createSharedAuthzCache } from './sharedCache';

type ProfileApiBoundary = {
  getCapabilitiesFromApi(
    ...args: Parameters<AuthzProfileService['getCapabilities']>
  ): ReturnType<AuthzProfileService['getCapabilities']>;
};

describe('shared authorization cache', () => {
  let db: Knex;
  beforeEach(() => {
    db = knex({
      client: 'better-sqlite3',
      connection: { filename: ':memory:' },
      useNullAsDefault: true,
    });
  });
  afterEach(async () => {
    jest.restoreAllMocks();
    await db.destroy();
  });

  test('a separate replica reads capability profiles warmed at sign-in', async () => {
    const replicaA = new AuthzProfileCache(await createSharedAuthzCache(db));
    const replicaB = new AuthzProfileCache(await createSharedAuthzCache(db));
    const profile = { capabilities: {} };
    await replicaA.setByUser(
      'user:default/operator',
      profile,
      60_000,
      'global',
    );
    expect(await replicaB.getByUser('user:default/operator', 'global')).toEqual(
      profile,
    );
    expect(
      await replicaB.getByUser('user:default/other', 'global'),
    ).toBeUndefined();
  });
  test('expired entries never authorize and are reclaimed on a subsequent write', async () => {
    const now = jest.spyOn(Date, 'now').mockReturnValue(1000);
    const cache = await createSharedAuthzCache(db);
    await cache.set('old', { allowed: true }, { ttl: 10 });
    now.mockReturnValue(1010);
    expect(await cache.get('old')).toBeUndefined();
    await cache.set('fresh', false);
    expect(await db('onemind_authz_cache').count('* as n').first()).toEqual({
      n: 1,
    });
  });
  test('real upstream cache-miss warning fails closed without fetching, even for an administrator', async () => {
    const now = jest.spyOn(Date, 'now').mockReturnValue(1000);
    const cache = new AuthzProfileCache(await createSharedAuthzCache(db));
    const logger: LoggerService = {
      debug: jest.fn(),
      info: jest.fn(),
      warn: jest.fn(),
      error: jest.fn(),
      child: () => logger,
    };
    const service = new AuthzProfileService({
      baseUrl: 'https://unused.example.test',
      logger,
      cache,
    });
    const fetchCapabilities = jest
      .spyOn(service, 'getCapabilities')
      .mockRejectedValue(new Error('Unexpected fetch'));
    const id = 'user:default/operator';
    const profile = { capabilities: {} };
    await cache.setByUser(id, profile, 1000, 'global');
    expect(await service.getCapabilitiesForUser(id)).toEqual(profile);
    expect(logger.warn).not.toHaveBeenCalled();
    now.mockReturnValue(2000);
    const policy = createOneMindPolicy(
      new OpenChoreoPermissionPolicy({ authzService: service, logger }),
    );
    const decision = await policy.handle(
      {
        permission: createPermission({
          name: 'openchoreo.component.create',
          attributes: { action: 'create' },
        }),
      },
      {
        credentials: {
          $$type: '@backstage/BackstageCredentials',
          principal: { type: 'user', userEntityRef: id },
        },
        info: {
          userEntityRef: id,
          ownershipEntityRefs: ['group:default/onemind-admins'],
        },
      },
    );
    expect(decision).toEqual({ result: AuthorizeResult.DENY });
    expect(logger.warn).toHaveBeenCalledWith(
      `No token and no cached capabilities for ${id}`,
    );
    expect(fetchCapabilities).not.toHaveBeenCalled();
  });
  test('a renewed profile restores cross-replica reads only until its new token expiry', async () => {
    const now = jest.spyOn(Date, 'now').mockReturnValue(1000);
    const logger: LoggerService = {
      debug: jest.fn(),
      info: jest.fn(),
      warn: jest.fn(),
      error: jest.fn(),
      child: () => logger,
    };
    const a = new AuthzProfileCache(await createSharedAuthzCache(db));
    const b = new AuthzProfileCache(await createSharedAuthzCache(db));
    const service = new AuthzProfileService({
      baseUrl: 'https://unused.example.test',
      logger,
      cache: a,
    });
    const profile = { capabilities: {} };
    // Only the upstream network boundary is mocked. Actual cache/TTL logic runs.
    const fetchProfile = jest
      .spyOn(service as unknown as ProfileApiBoundary, 'getCapabilitiesFromApi')
      .mockResolvedValue(profile);
    const token = (exp: number) =>
      `fixture.${Buffer.from(JSON.stringify({ exp })).toString(
        'base64url',
      )}.fixture`;
    const actor = 'user:default/operator';
    await service.getCapabilitiesForUser(actor, token(2));
    expect(await b.getByUser(actor, 'global')).toEqual(profile);
    now.mockReturnValue(2000);
    expect(await b.getByUser(actor, 'global')).toBeUndefined();
    expect(await service.getCapabilitiesForUser(actor)).toEqual({
      capabilities: {},
    });
    expect(fetchProfile).toHaveBeenCalledTimes(1);
    await service.getCapabilitiesForUser(actor, token(9));
    expect(await b.getByUser(actor, 'global')).toEqual(profile);
    expect(await b.getByUser('user:default/other', 'global')).toBeUndefined();
    expect(await b.getByUser(actor, 'other-namespace')).toBeUndefined();
    expect(fetchProfile).toHaveBeenCalledTimes(2);
    now.mockReturnValue(9000);
    expect(await b.getByUser(actor, 'global')).toBeUndefined();
    expect(
      JSON.stringify(await db('onemind_authz_cache').select()),
    ).not.toContain('fixture.');
  });
  test('refresh failure does not resurrect an expired profile', async () => {
    const now = jest.spyOn(Date, 'now').mockReturnValue(1000);
    const cache = new AuthzProfileCache(await createSharedAuthzCache(db));
    const logger: LoggerService = {
      debug: jest.fn(),
      info: jest.fn(),
      warn: jest.fn(),
      error: jest.fn(),
      child: () => logger,
    };
    const service = new AuthzProfileService({
      baseUrl: 'https://unused.example.test',
      logger,
      cache,
    });
    await cache.setByUser(
      'user:default/operator',
      { capabilities: {} },
      1,
      'global',
    );
    now.mockReturnValue(1001);
    jest
      .spyOn(service as unknown as ProfileApiBoundary, 'getCapabilitiesFromApi')
      .mockRejectedValue(new Error('Provider unavailable'));
    await expect(
      service.getCapabilitiesForUser('user:default/operator', 'fixture'),
    ).rejects.toThrow('Provider unavailable');
    expect(
      await cache.getByUser('user:default/operator', 'global'),
    ).toBeUndefined();
  });
  test('preserves token TTL, stores only hashed keys, and isolates scopes', async () => {
    const cache = await createSharedAuthzCache(db);
    jest.spyOn(Date, 'now').mockReturnValue(1000);
    await cache.set(
      'user:default/operator:global',
      { allowed: true },
      { ttl: 1234 },
    );
    const row = await db('onemind_authz_cache').first();
    expect(row.cache_key).toMatch(/^[a-f0-9]{64}$/);
    expect(row.cache_key).not.toContain('operator');
    expect(Number(row.expires_at)).toBe(2234);
    expect(await cache.get('user:default/operator:other')).toBeUndefined();
  });
  test('concurrent writes upsert and deletion is visible across replicas', async () => {
    const a = await createSharedAuthzCache(db);
    const b = await createSharedAuthzCache(db);
    await Promise.all([a.set('same', true), b.set('same', false)]);
    expect(await db('onemind_authz_cache').count('* as n').first()).toEqual({
      n: 1,
    });
    await b.delete('same');
    expect(await a.get('same')).toBeUndefined();
  });
  test('supports duration defaults, explicit overrides and nonpositive TTL invalidation', async () => {
    jest.spyOn(Date, 'now').mockReturnValue(1000);
    const a = (await createSharedAuthzCache(db)).withOptions({
      defaultTtl: { seconds: 2 },
    });
    await a.set('a', true);
    expect(Number((await db('onemind_authz_cache').first()).expires_at)).toBe(
      3000,
    );
    await a.set('a', false, { ttl: 0 });
    expect(await a.get('a')).toBeUndefined();
    await expect(a.set('bad', true, { ttl: Number.NaN })).rejects.toThrow(
      'Invalid',
    );
  });
  test('database failures propagate rather than falling back to stale memory', async () => {
    const cache = await createSharedAuthzCache(db);
    await cache.set('a', true);
    await db.schema.dropTable('onemind_authz_cache');
    await expect(cache.get('a')).rejects.toThrow();
    await expect(cache.set('a', true)).rejects.toThrow();
  });
  test('initialization is idempotent and preserves unrelated plugin data', async () => {
    await db.schema.createTable('existing_data', t => t.string('value'));
    await db('existing_data').insert({ value: 'preserved' });
    await createSharedAuthzCache(db);
    await createSharedAuthzCache(db);
    expect(await db('existing_data').first()).toEqual({ value: 'preserved' });
  });
});
