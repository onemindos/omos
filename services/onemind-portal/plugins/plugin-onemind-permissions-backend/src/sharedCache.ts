import { createHash } from 'node:crypto';
import type {
  CacheService,
  CacheServiceOptions,
  CacheServiceSetOptions,
  DatabaseService,
} from '@backstage/backend-plugin-api';
import { durationToMilliseconds } from '@backstage/types';
import type { JsonValue } from '@backstage/types';

type Client = Awaited<ReturnType<DatabaseService['getClient']>>;
const TABLE = 'onemind_authz_cache';
const digest = (key: string) => createHash('sha256').update(key).digest('hex');
const milliseconds = (ttl: NonNullable<CacheServiceSetOptions['ttl']>) =>
  typeof ttl === 'number' ? ttl : durationToMilliseconds(ttl);

// Permission calls can cross replicas without the IDP token. Share the upstream
// cache in the permission plugin DB, storing only hashes/authorization results.
export async function createSharedAuthzCache(
  db: Client,
): Promise<CacheService> {
  if (!(await db.schema.hasTable(TABLE))) {
    try {
      await db.schema.createTable(TABLE, table => {
        table.string('cache_key', 64).primary();
        table.text('value').notNullable();
        table.bigInteger('expires_at').notNullable();
      });
    } catch (error) {
      const code = (error as { code?: string }).code;
      // PostgreSQL concurrent CREATE can report duplicate table or duplicate type.
      if (
        !['42P07', '23505', 'SQLITE_ERROR'].includes(code ?? '') ||
        !(await db.schema.hasTable(TABLE))
      )
        throw error;
    }
  }
  return new SharedAuthzCache(db);
}

class SharedAuthzCache implements CacheService {
  constructor(
    private readonly db: Client,
    private readonly options: CacheServiceOptions = {},
  ) {}

  async get<T extends JsonValue>(key: string): Promise<T | undefined> {
    const row = await this.db(TABLE)
      .where({ cache_key: digest(key) })
      .where('expires_at', '>', Date.now())
      .first('value');
    return row ? (JSON.parse(row.value) as T) : undefined;
  }

  async set(key: string, value: JsonValue, options?: CacheServiceSetOptions) {
    const ttl = milliseconds(
      options?.ttl ?? this.options.defaultTtl ?? 300_000,
    );
    if (!Number.isFinite(ttl))
      throw new Error('Invalid authorization cache TTL');
    if (ttl <= 0) return this.delete(key);
    const now = Date.now();
    await this.db(TABLE).where('expires_at', '<=', now).delete();
    await this.db(TABLE)
      .insert({
        cache_key: digest(key),
        value: JSON.stringify(value),
        expires_at: now + ttl,
      })
      .onConflict('cache_key')
      .merge();
  }

  async delete(key: string) {
    await this.db(TABLE)
      .where({ cache_key: digest(key) })
      .delete();
  }

  withOptions(options: CacheServiceOptions): CacheService {
    return new SharedAuthzCache(this.db, { ...this.options, ...options });
  }
}
