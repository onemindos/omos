/** @jest-environment node */
import {
  createAlertReader,
  parseAlertSource,
  scopeAlertSource,
} from './prometheusAlerts';

const now = Date.now();
const stamp = new Date(now).toISOString();
const instance = (labels: Record<string, string>, state = 'firing') => ({
  labels,
  annotations: { summary: 'Observed workload alert' },
  state,
  activeAt: stamp,
  value: '1',
});
const rule = (alerts: ReturnType<typeof instance>[] = [], extra = {}) => ({
  type: 'alerting',
  name: 'KubePodNotReady',
  state: alerts.length ? 'firing' : 'inactive',
  health: 'ok',
  lastEvaluation: stamp,
  query: 'up == 0',
  labels: { severity: 'warning' },
  annotations: {},
  alerts,
  ...extra,
});
const source = (rules = [rule()]) => ({
  status: 'success',
  data: { groups: [{ name: 'fixture', rules }] },
});

test('strict workload scope excludes other namespace, database and unrelated pods', () => {
  const data = parseAlertSource(
    source([
      rule([
        instance({ namespace: 'backstage', pod: 'backstage-7c75c78b4-h2hp6' }),
        instance({ namespace: 'backstage', pod: 'backstage-postgres-0' }),
        instance({ namespace: 'other', pod: 'backstage-7c75c78b4-h2hp6' }),
        instance({ namespace: 'backstage', pod: 'unrelated-abcde-abcde' }),
      ]),
    ]),
  );
  expect(
    scopeAlertSource(data, 'backstage', now).data.groups[0].rules[0].alerts,
  ).toHaveLength(1);
  expect(
    scopeAlertSource(data, 'postgres', now).data.groups[0].rules[0].alerts[0]
      .labels.pod,
  ).toBe('backstage-postgres-0');
  expect(
    scopeAlertSource(data, 'nats', now).data.groups[0].rules[0].alerts,
  ).toEqual([]);
});
test('matches controller alerts and pending state without requiring a pod', () => {
  const data = parseAlertSource(
    source([
      rule([
        instance(
          { namespace: 'backstage', deployment: 'backstage' },
          'pending',
        ),
      ]),
    ]),
  );
  expect(
    scopeAlertSource(data, 'backstage', now).data.groups[0].rules[0].alerts[0]
      .state,
  ).toBe('pending');
});
test('NATS custom service rules include explicit bus dependencies but not arbitrary services', () => {
  const data = parseAlertSource(
    source([
      rule([instance({ service: 'nats' })], {
        name: 'TAKCoTFlowStopped',
        labels: { service: 'nats' },
      }),
      rule([instance({ service: 'nats' })], { name: 'ArbitraryRule' }),
    ]),
  );
  expect(
    scopeAlertSource(data, 'nats', now)
      .data.groups.flatMap(g => g.rules)
      .map(r => r.name),
  ).toEqual(['TAKCoTFlowStopped']);
});
test('no-match, no candidate coverage, stale and rule errors remain distinct', () => {
  expect(
    scopeAlertSource(parseAlertSource(source()), 'backstage', now).coverage,
  ).toBe('observed');
  expect(
    scopeAlertSource(parseAlertSource(source([])), 'backstage', now).coverage,
  ).toBe('not-configured');
  const data = parseAlertSource(
    source([
      rule([], {
        health: 'err',
        lastEvaluation: new Date(now - 200000).toISOString(),
      }),
    ]),
  );
  expect(scopeAlertSource(data, 'backstage', now)).toMatchObject({
    coverage: 'partial',
    unhealthyRules: 1,
    staleRules: 1,
  });
  expect(
    scopeAlertSource(
      parseAlertSource(
        source([
          rule([], { lastEvaluation: new Date(now + 60000).toISOString() }),
        ]),
      ),
      'backstage',
      now,
    ).coverage,
  ).toBe('partial');
});
test.each([
  null,
  {},
  { status: 'error' },
  source([rule([], { type: 'recording' })]),
  source([rule([], { alerts: null })]),
  source([rule([], { lastEvaluation: 'invalid' })]),
  { ...source(), warnings: ['partial'] },
])('rejects malformed/partial protocol response %j', value => {
  expect(() => parseAlertSource(value)).toThrow();
});
test('fixed GET shared across scopes, bounded cache and no stale fallback after failure', async () => {
  let time = now;
  const get = jest.fn(async (_url: string) => ({
    status: 200,
    body: JSON.stringify(source()),
  }));
  const reader = createAlertReader(get, () => time);
  const [a, b] = await Promise.all([reader('backstage'), reader('nats')]);
  expect(a.scope).toBe('backstage');
  expect(b.scope).toBe('nats');
  expect(get).toHaveBeenCalledTimes(1);
  expect(get.mock.calls[0][0]).toBe(
    'http://monitoring-kube-prometheus-prometheus.monitoring.svc:9090/api/v1/rules?type=alert',
  );
  time += 16000;
  get.mockResolvedValueOnce({
    status: 503,
    body: 'private error, do not surface',
  });
  await expect(reader('backstage')).rejects.toThrow('unavailable');
  await expect(reader('nats')).resolves.toHaveProperty('scope', 'nats');
});
