import knexFactory from 'knex';
import type { Db } from './store';

const open: Db[] = [];

// In-memory SQLite, the same engine the portal uses locally. Call destroyMemoryDbs() in afterAll.
export function memoryDb(): Db {
  const db = knexFactory({
    client: 'better-sqlite3',
    connection: { filename: ':memory:' },
    useNullAsDefault: true,
  }) as unknown as Db;
  open.push(db);
  return db;
}

export async function destroyMemoryDbs() {
  await Promise.all(open.splice(0).map(db => db.destroy()));
}
