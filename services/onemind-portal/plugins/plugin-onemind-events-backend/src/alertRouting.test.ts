import { routeAlerts } from './alertRouting';
import { createStore } from './store';
import { destroyMemoryDbs, memoryDb } from './testDb';
import type { Sinks } from './scan';

afterAll(destroyMemoryDbs);

const alert = (status: 'firing' | 'resolved', over: Record<string, unknown> = {}) => ({
  status,
  labels: { alertname: 'KubePodCrashLooping', severity: 'critical', namespace: 'backstage', pod: 'backstage-1' },
  annotations: { summary: 'Pod is crash looping' },
  startsAt: '2026-10-21T10:00:00Z',
  endsAt: status === 'resolved' ? '2026-10-21T10:30:00Z' : '0001-01-01T00:00:00Z',
  fingerprint: 'abc123def4567890',
  ...over,
});
const body = (...alerts: unknown[]) => ({ version: '4', status: 'firing', alerts });

async function setup() {
  const store = await createStore(memoryDb());
  const notify = jest.fn(async () => {});
  const signal = jest.fn(async () => {});
  const sinks: Sinks = { notify, signal };
  return { store, sinks, notify, signal };
}

describe('alert routing', () => {
  test('firing then resolved creates one item and resolves it on the same scope, with signals', async () => {
    const { store, sinks, notify, signal } = await setup();
    const fired = await routeAlerts({ body: body(alert('firing')), store, sinks });
    expect(fired).toMatchObject({ notified: 1 });
    expect(notify).toHaveBeenCalledWith(expect.objectContaining({ severity: 'critical', topic: 'alerts', scope: 'alert:abc123def4567890' }));
    expect(signal).toHaveBeenCalledWith(expect.objectContaining({ action: 'alert.firing', severity: 'critical', source: 'alertmanager' }));

    const resolved = await routeAlerts({ body: body(alert('resolved')), store, sinks });
    expect(resolved).toMatchObject({ resolved: 1 });
    expect(notify).toHaveBeenCalledTimes(2);
    expect(notify).toHaveBeenLastCalledWith(expect.objectContaining({ severity: 'low', scope: 'alert:abc123def4567890', title: expect.stringMatching(/^Resolved:/) }));
    expect(signal).toHaveBeenLastCalledWith(expect.objectContaining({ action: 'alert.resolved' }));
  });

  test('Alertmanager repeat notifications for a still-firing alert add nothing', async () => {
    const { store, sinks, notify, signal } = await setup();
    await routeAlerts({ body: body(alert('firing')), store, sinks });
    const repeat = await routeAlerts({ body: body(alert('firing')), store, sinks });
    expect(repeat).toMatchObject({ repeats: 1, notified: 0 });
    expect(notify).toHaveBeenCalledTimes(1);
    expect(signal).toHaveBeenCalledTimes(1);
  });

  test('re-firing after resolution notifies again on the same scope', async () => {
    const { store, sinks, notify } = await setup();
    await routeAlerts({ body: body(alert('firing')), store, sinks });
    await routeAlerts({ body: body(alert('resolved')), store, sinks });
    await routeAlerts({ body: body(alert('firing')), store, sinks });
    expect(notify).toHaveBeenCalledTimes(3);
    const scopes = notify.mock.calls.map(call => (call as unknown as [{ scope: string }])[0].scope);
    expect(new Set(scopes).size).toBe(1);
  });

  test('a resolved alert we never announced is ignored', async () => {
    const { store, sinks, notify } = await setup();
    const result = await routeAlerts({ body: body(alert('resolved')), store, sinks });
    expect(result).toMatchObject({ ignored: 1, resolved: 0 });
    expect(notify).not.toHaveBeenCalled();
  });

  test('warning is high with a signal; info is normal with none; Watchdog is dropped', async () => {
    const { store, sinks, notify, signal } = await setup();
    await routeAlerts({
      body: body(
        alert('firing', { fingerprint: 'w1', labels: { alertname: 'HighLatency', severity: 'warning' } }),
        alert('firing', { fingerprint: 'i1', labels: { alertname: 'FYI', severity: 'info' } }),
        alert('firing', { fingerprint: 'd1', labels: { alertname: 'Watchdog', severity: 'none' } }),
      ),
      store,
      sinks,
    });
    expect(notify.mock.calls.map(call => (call as unknown as [{ severity: string }])[0].severity)).toEqual(['high', 'normal']);
    expect(signal).toHaveBeenCalledTimes(1);
  });

  test('a malformed payload is dropped without side effects', async () => {
    const { store, sinks, notify } = await setup();
    const result = await routeAlerts({ body: { alerts: 'nope' }, store, sinks });
    expect(result).toEqual({ notified: 0, resolved: 0, repeats: 0, ignored: 0 });
    expect(notify).not.toHaveBeenCalled();
  });

});
