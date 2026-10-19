import knex from 'knex';
import type { Knex } from 'knex';
import { AuthzProfileCache } from '@openchoreo/backstage-plugin-permission-backend-module-openchoreo-policy';
import { createSharedAuthzCache } from './sharedCache';

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
    const a = await createSharedAuthzCache(db),
      b = await createSharedAuthzCache(db);
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
