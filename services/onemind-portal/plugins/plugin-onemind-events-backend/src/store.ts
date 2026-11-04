import type { DatabaseService } from '@backstage/backend-plugin-api';

export type Db = Awaited<ReturnType<DatabaseService['getClient']>>;

export type TimelineEntry = {
  id: string;
  kind: string;
  resource: string;
  occurredAt: string;
  source: string;
};

export type ConditionRow = {
  scope: string;
  reason: 'failed' | 'degraded' | 'out-of-sync' | 'suspended';
  since: number;
  notified: boolean;
};

export type AlertRow = {
  fingerprint: string;
  status: 'firing' | 'resolved';
  severity: string;
};

export type DigestSummary = {
  total: number;
  byResource: { resource: string; count: number }[];
};

// Creates a table only if absent. Two modules of the same plugin may initialise at once, so a
// failed create is retried briefly and tolerated when the table now exists.
async function ensureTable(
  db: Db,
  name: string,
  build: Parameters<Db['schema']['createTable']>[1],
) {
  for (let attempt = 0; ; attempt += 1) {
    if (await db.schema.hasTable(name)) return;
    try {
      await db.schema.createTable(name, build);
      return;
    } catch (error) {
      if (attempt >= 3) throw error;
      await new Promise(resolve => setTimeout(resolve, 250 * (attempt + 1)));
    }
  }
}

const stores = new WeakMap<object, Promise<Store>>();

// One store per database client: the observer and alert routing modules share it, so the
// tables are created once.
export function createStore(db: Db): Promise<Store> {
  let store = stores.get(db);
  if (!store) {
    store = buildStore(db);
    stores.set(db, store);
  }
  return store;
}

async function buildStore(db: Db) {
  await ensureTable(db, 'onemind_activity_timeline', table => {
    table.string('id', 64).primary();
    table.string('kind', 40).notNullable();
    table.string('resource', 200).notNullable();
    table.string('occurred_at', 40).notNullable().index();
    table.string('source', 40).notNullable();
  });
  await ensureTable(db, 'onemind_activity_digest', table => {
    table.string('day', 10).primary();
    table.string('sent_at', 40).notNullable();
    table.integer('total').notNullable();
  });
  await ensureTable(db, 'onemind_condition_state', table => {
    table.string('scope', 200).primary();
    table.string('reason', 24).notNullable();
    table.string('since', 40).notNullable();
    table.integer('notified').notNullable();
  });
  await ensureTable(db, 'onemind_alert_state', table => {
    table.string('fingerprint', 64).primary();
    table.string('status', 16).notNullable();
    table.string('severity', 16).notNullable();
    table.string('updated_at', 40).notNullable();
  });

  return {
    timeline: {
      async record(entry: TimelineEntry) {
        await db('onemind_activity_timeline')
          .insert({
            id: entry.id,
            kind: entry.kind,
            resource: entry.resource.slice(0, 200),
            occurred_at: entry.occurredAt,
            source: entry.source,
          })
          .onConflict('id')
          .ignore();
      },
      async recent(limit: number): Promise<TimelineEntry[]> {
        const rows = await db('onemind_activity_timeline')
          .orderBy('occurred_at', 'desc')
          .limit(Math.min(Math.max(limit, 1), 200));
        return rows.map(
          (row: Record<string, string>): TimelineEntry => ({
            id: row.id,
            kind: row.kind,
            resource: row.resource,
            occurredAt: row.occurred_at,
            source: row.source,
          }),
        );
      },
      async summary(sinceIso: string, kind: string): Promise<DigestSummary> {
        const rows = await db('onemind_activity_timeline')
          .where('kind', kind)
          .andWhere('occurred_at', '>', sinceIso)
          .select('resource')
          .count({ count: '*' })
          .groupBy('resource')
          .orderBy('count', 'desc');
        const byResource = rows.map((row: Record<string, unknown>) => ({
          resource: String(row.resource),
          count: Number(row.count),
        }));
        return {
          total: byResource.reduce(
            (sum: number, row: { count: number }) => sum + row.count,
            0,
          ),
          byResource,
        };
      },
      async prune(beforeIso: string) {
        await db('onemind_activity_timeline')
          .where('occurred_at', '<', beforeIso)
          .delete();
      },
    },
    digest: {
      async sent(day: string) {
        return Boolean(await db('onemind_activity_digest').where({ day }).first());
      },
      async markSent(day: string, total: number) {
        await db('onemind_activity_digest')
          .insert({ day, sent_at: new Date().toISOString(), total })
          .onConflict('day')
          .ignore();
      },
    },
    conditions: {
      async all(): Promise<ConditionRow[]> {
        const rows = await db('onemind_condition_state');
        return rows.map((row: Record<string, unknown>) => ({
          scope: String(row.scope),
          reason: row.reason as ConditionRow['reason'],
          since: Date.parse(String(row.since)),
          notified: Number(row.notified) === 1,
        }));
      },
      async put(row: ConditionRow) {
        await db('onemind_condition_state')
          .insert({
            scope: row.scope,
            reason: row.reason,
            since: new Date(row.since).toISOString(),
            notified: row.notified ? 1 : 0,
          })
          .onConflict('scope')
          .merge();
      },
      async remove(scope: string) {
        await db('onemind_condition_state').where({ scope }).delete();
      },
    },
    alerts: {
      async get(fingerprint: string): Promise<AlertRow | undefined> {
        const row = await db('onemind_alert_state').where({ fingerprint }).first();
        return row
          ? {
              fingerprint,
              status: row.status as AlertRow['status'],
              severity: String(row.severity),
            }
          : undefined;
      },
      async put(row: AlertRow) {
        await db('onemind_alert_state')
          .insert({
            fingerprint: row.fingerprint,
            status: row.status,
            severity: row.severity,
            updated_at: new Date().toISOString(),
          })
          .onConflict('fingerprint')
          .merge();
      },
      async pruneResolved(beforeIso: string) {
        await db('onemind_alert_state')
          .where('status', 'resolved')
          .andWhere('updated_at', '<', beforeIso)
          .delete();
      },
    },
  };
}

export type Store = Awaited<ReturnType<typeof buildStore>>;
