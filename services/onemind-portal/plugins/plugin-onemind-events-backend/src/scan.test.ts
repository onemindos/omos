import { destroyMemoryDbs, memoryDb } from './testDb';
import {
  classifyApplications,
  classifyBackups,
  ensureDeliveryTables,
  processConditions,
  processObservations,
  runDigest,
  type Observation,
  type Sinks,
} from './scan';
import { createStore } from './store';

const T0 = new Date('2026-10-21T10:00:00Z');
let clock = T0;
afterAll(destroyMemoryDbs);
const now = () => clock;
const at = (minutes: number) => new Date(T0.getTime() + minutes * 60_000);

async function setup() {
  const db = memoryDb();
  await ensureDeliveryTables(db);
  const store = await createStore(db);
  const notify = jest.fn(async () => {});
  const signal = jest.fn(async () => {});
  const sinks: Sinks = { notify, signal };
  clock = T0;
  return { db, store, sinks, notify, signal };
}

const op = (kind: Observation['kind'], minutes: number, key = `${kind}:${minutes}`): Observation => ({
  key,
  resource: 'platform-backstage',
  kind,
  occurredAt: at(minutes).toISOString(),
  link: 'https://argocd.example/applications/platform-backstage',
});

describe('Argo operations', () => {
  test('a synthetic success produces zero notifications and lands on the timeline', async () => {
    const { db, store, sinks, notify, signal } = await setup();
    await processObservations({ db, store, source: 'argocd', values: [], sinks, now });
    clock = at(5);
    await processObservations({ db, store, source: 'argocd', values: [op('deployment.succeeded', 3)], sinks, now });
    expect(notify).not.toHaveBeenCalled();
    expect(signal).not.toHaveBeenCalled();
    expect((await store.timeline.recent(10)).map(e => e.kind)).toEqual(['deployment.succeeded']);
  });

  test('a synthetic failure produces one notification and repeating it adds none', async () => {
    const { db, store, sinks, notify, signal } = await setup();
    await processObservations({ db, store, source: 'argocd', values: [], sinks, now });
    clock = at(5);
    const failed = op('deployment.failed', 3);
    await processObservations({ db, store, source: 'argocd', values: [failed], sinks, now });
    await processObservations({ db, store, source: 'argocd', values: [failed], sinks, now });
    expect(notify).toHaveBeenCalledTimes(1);
    expect(notify).toHaveBeenCalledWith(expect.objectContaining({ severity: 'high', topic: 'deployments', scope: 'argo:platform-backstage' }));
    expect(signal).toHaveBeenCalledTimes(1);
  });

  test('a later failure of the same app reuses the scope so the inbox shows one item', async () => {
    const { db, store, sinks, notify } = await setup();
    await processObservations({ db, store, source: 'argocd', values: [], sinks, now });
    clock = at(30);
    await processObservations({ db, store, source: 'argocd', values: [op('deployment.failed', 3, 'a'), op('deployment.failed', 20, 'b')], sinks, now });
    const scopes = notify.mock.calls.map(call => (call as unknown as [{ scope: string }])[0].scope);
    expect(new Set(scopes)).toEqual(new Set(['argo:platform-backstage']));
  });

  test('history before the baseline is never replayed', async () => {
    const { db, store, sinks, notify } = await setup();
    await processObservations({ db, store, source: 'argocd', values: [op('deployment.failed', -600)], sinks, now });
    expect(notify).not.toHaveBeenCalled();
    clock = at(1);
    await processObservations({ db, store, source: 'argocd', values: [op('deployment.failed', -600)], sinks, now });
    expect(notify).not.toHaveBeenCalled();
  });

  test('a success after a failure announces recovery on the same scope', async () => {
    const { db, store, sinks, notify } = await setup();
    await processObservations({ db, store, source: 'argocd', values: [], sinks, now });
    clock = at(5);
    await processObservations({ db, store, source: 'argocd', values: [op('deployment.failed', 3)], sinks, now });
    clock = at(10);
    await processObservations({ db, store, source: 'argocd', values: [op('deployment.succeeded', 8)], sinks, now });
    expect(notify).toHaveBeenCalledTimes(2);
    expect(notify).toHaveBeenLastCalledWith(expect.objectContaining({ severity: 'low', scope: 'argo:platform-backstage', title: 'Recovered: platform-backstage' }));
  });
});

describe('Argo conditions', () => {
  const app = (name: string, health: string, sync: string) => ({ name, health, sync, link: '/x' });

  test('degraded notifies once; out-of-sync waits fifteen minutes; healthy recovers', async () => {
    const { store, sinks, notify } = await setup();
    const run = async (minutes: number, apps: ReturnType<typeof app>[]) => {
      clock = at(minutes);
      const { active, healthy } = classifyApplications(apps);
      await processConditions({ store, active, healthy, topic: 'deployments', sinks, now });
    };
    await run(0, [app('a', 'Degraded', 'Synced'), app('b', 'Healthy', 'OutOfSync')]);
    expect(notify).toHaveBeenCalledTimes(1);
    expect(notify).toHaveBeenCalledWith(expect.objectContaining({ scope: 'argo:a', severity: 'high' }));
    await run(10, [app('a', 'Degraded', 'Synced'), app('b', 'Healthy', 'OutOfSync')]);
    expect(notify).toHaveBeenCalledTimes(1);
    await run(15, [app('a', 'Degraded', 'Synced'), app('b', 'Healthy', 'OutOfSync')]);
    expect(notify).toHaveBeenCalledTimes(2);
    expect(notify).toHaveBeenLastCalledWith(expect.objectContaining({ scope: 'argo:b', title: expect.stringContaining('out of sync') }));
    await run(20, [app('a', 'Healthy', 'Synced'), app('b', 'Healthy', 'Synced')]);
    expect(notify).toHaveBeenCalledTimes(4);
    expect(notify).toHaveBeenCalledWith(expect.objectContaining({ title: 'Recovered: a', severity: 'low' }));
  });

  test('progressing apps are neither raised nor recovered', () => {
    expect(classifyApplications([app('c', 'Progressing', 'Synced')])).toEqual({ active: [], healthy: [] });
  });
});

describe('Backups', () => {
  test('a suspended watched CronJob is one high notification, and unsuspending recovers it', async () => {
    const { store, sinks, notify } = await setup();
    const run = async (minutes: number, suspended: boolean) => {
      clock = at(minutes);
      const { active, healthy } = classifyBackups([{ name: 'backup-postgres-backstage', suspended, link: '/runbook' }]);
      await processConditions({ store, active, healthy, topic: 'backups', sinks, now });
    };
    await run(0, true);
    await run(5, true);
    expect(notify).toHaveBeenCalledTimes(1);
    expect(notify).toHaveBeenCalledWith(expect.objectContaining({ topic: 'backups', severity: 'high' }));
    await run(10, false);
    expect(notify).toHaveBeenCalledTimes(2);
  });

  test('a failed job is not auto-recovered by the CronJob merely being active', async () => {
    const { db, store, sinks, notify } = await setup();
    const job = (kind: Observation['kind'], minutes: number): Observation => ({ key: `${kind}${minutes}`, resource: 'backstage/backup-postgres-backstage', kind, occurredAt: at(minutes).toISOString(), link: '/runbook' });
    await processObservations({ db, store, source: 'backstage-backup', values: [], sinks, now });
    clock = at(5);
    await processObservations({ db, store, source: 'backstage-backup', values: [job('backup.failed', 3)], sinks, now });
    clock = at(30);
    const { active, healthy } = classifyBackups([{ name: 'backup-postgres-backstage', suspended: false, link: '/runbook' }]);
    await processConditions({ store, active, healthy, topic: 'backups', sinks, now });
    expect(notify).toHaveBeenCalledTimes(1);
    clock = at(40);
    await processObservations({ db, store, source: 'backstage-backup', values: [job('backup.succeeded', 35)], sinks, now });
    expect(notify).toHaveBeenCalledTimes(2);
    expect(notify).toHaveBeenLastCalledWith(expect.objectContaining({ severity: 'low', topic: 'backups' }));
  });
});

describe('Daily digest', () => {
  test('summarises the last 24 hours once per day, and sends nothing when empty', async () => {
    const { db, store, sinks, notify } = await setup();
    await processObservations({ db, store, source: 'argocd', values: [], sinks, now });
    clock = new Date('2026-10-21T09:00:00Z');
    const successes = [op('deployment.succeeded', 1, 's1'), op('deployment.succeeded', 2, 's2'), { ...op('deployment.succeeded', 3, 's3'), resource: 'other-app' }];
    await processObservations({ db, store, source: 'argocd', values: successes, sinks, now });
    expect(notify).not.toHaveBeenCalled();

    expect(await runDigest({ store, sinks, now: () => new Date('2026-10-21T12:00:00Z'), hourUtc: 13 })).toBe('too-early');
    const sendTime = () => new Date('2026-10-21T13:30:00Z');
    expect(await runDigest({ store, sinks, now: sendTime, hourUtc: 13 })).toBe('sent');
    expect(await runDigest({ store, sinks, now: sendTime, hourUtc: 13 })).toBe('already-sent');
    expect(notify).toHaveBeenCalledTimes(1);
    expect(notify).toHaveBeenCalledWith(expect.objectContaining({
      title: 'Daily digest: 3 deployments succeeded',
      severity: 'normal',
      topic: 'digest',
      scope: 'digest:deployments:2026-10-21',
    }));

    const nextDay = () => new Date('2026-10-23T13:30:00Z');
    expect(await runDigest({ store, sinks, now: nextDay, hourUtc: 13 })).toBe('empty');
    expect(notify).toHaveBeenCalledTimes(1);
  });
});
