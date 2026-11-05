import {
  countsFor,
  matchPolicy,
  targetsForEntity,
  type PolicyFinding,
  type PolicyResource,
  type PolicyResponse,
} from './policyMatch';

const entity = (annotations: Record<string, string>) => ({
  metadata: { annotations },
});

const tally = (over: Partial<PolicyResource['summary']> = {}) => ({
  pass: 0,
  fail: 0,
  warn: 0,
  error: 0,
  skip: 0,
  ...over,
});
const resource = (
  kind: string,
  name: string,
  workload: string,
  summary = tally({ pass: 3 }),
): PolicyResource => ({ kind, name, workload, summary });
const finding = (
  kind: string,
  name: string,
  workload: string,
  policy: string,
  result: PolicyFinding['result'] = 'fail',
): PolicyFinding => ({
  resource: { kind, name, workload },
  policy,
  rule: 'rule',
  result,
  severity: null,
  message: '',
});

type Data = Extract<PolicyResponse, { status: 'ok' | 'stale' }>;
const response = (overrides: Partial<Data> = {}): Data => ({
  status: 'ok',
  observedAt: '2026-11-05T12:00:00Z',
  reports: 6,
  truncated: false,
  namespaces: {
    backstage: {
      resources: [
        resource('Deployment', 'backstage', 'backstage'),
        resource('Pod', 'backstage-5d9f-abcde', 'backstage'),
        resource('StatefulSet', 'backstage-postgres', 'backstage-postgres'),
        resource('Pod', 'backstage-postgres-0', 'backstage-postgres'),
      ],
      findings: [
        finding('Deployment', 'backstage', 'backstage', 'p1'),
        finding('Pod', 'backstage-5d9f-abcde', 'backstage', 'p1'),
        finding('Pod', 'backstage-postgres-0', 'backstage-postgres', 'p2'),
        finding('Deployment', 'backstage', 'backstage', 'p3', 'warn'),
      ],
    },
  },
  ...overrides,
});

describe('targetsForEntity', () => {
  it('reads namespaces and named workloads', () => {
    expect(
      targetsForEntity(
        entity({
          'backstage.io/kubernetes-namespace': 'backstage',
          'onemind.io/kubernetes-workloads':
            'backstage/backstage,backstage/backstage-postgres,nats/nats',
          'onemind.io/kubernetes-namespaces': 'backstage,nats',
        }),
      ),
    ).toEqual([
      {
        namespace: 'backstage',
        workloads: ['backstage', 'backstage-postgres'],
      },
      { namespace: 'nats', workloads: ['nats'] },
    ]);
  });

  it('reads the single-workload annotation of a cluster workload entity', () => {
    expect(
      targetsForEntity(
        entity({
          'backstage.io/kubernetes-namespace': 'backstage',
          'onemind.io/kubernetes-workload': 'backstage',
          'onemind.io/kubernetes-workload-kind': 'Deployment',
        }),
      ),
    ).toEqual([{ namespace: 'backstage', workloads: ['backstage'] }]);
  });

  it('returns whole-namespace targets when no workload is named, and none without a namespace', () => {
    expect(
      targetsForEntity(entity({ 'backstage.io/kubernetes-namespace': 'ops' })),
    ).toEqual([{ namespace: 'ops', workloads: null }]);
    expect(targetsForEntity(entity({}))).toEqual([]);
  });
});

describe('matchPolicy', () => {
  it('keeps backstage-postgres results out of the backstage app', () => {
    const matched = matchPolicy(response(), [
      { namespace: 'backstage', workloads: ['backstage'] },
    ]);
    expect(matched.scope).toBe('workloads');
    expect(matched.resources.map(r => r.name)).toEqual([
      'backstage',
      'backstage-5d9f-abcde',
    ]);
    expect(matched.findings.map(f => f.policy)).toEqual(['p1', 'p1', 'p3']);
  });

  it('shows the whole namespace only when no workload is named', () => {
    const matched = matchPolicy(response(), [
      { namespace: 'backstage', workloads: null },
    ]);
    expect(matched.scope).toBe('namespace');
    expect(matched.resources).toHaveLength(4);
  });

  it('lists named workloads that have no report instead of calling them clean', () => {
    const matched = matchPolicy(response(), [
      { namespace: 'backstage', workloads: ['backstage', 'ghost'] },
      { namespace: 'nowhere', workloads: ['x'] },
    ]);
    expect(matched.unreported).toEqual(['backstage/ghost', 'nowhere/x']);
  });
});

describe('countsFor', () => {
  it('counts a violation once per workload, policy and rule, not once per Pod and ReplicaSet', () => {
    const matched = matchPolicy(response(), [
      { namespace: 'backstage', workloads: ['backstage'] },
    ]);
    const counts = countsFor(matched);
    expect(counts.violations).toBe(1);
    expect(counts.warnings).toBe(1);
    expect(counts.passed).toBe(6);
    expect(counts.resources).toBe(2);
  });

  it('counts errors as violations', () => {
    const matched = matchPolicy(
      response({
        namespaces: {
          a: {
            resources: [resource('Deployment', 'x', 'x', tally({ error: 1 }))],
            findings: [finding('Deployment', 'x', 'x', 'p', 'error')],
          },
        },
      }),
      [{ namespace: 'a', workloads: null }],
    );
    expect(countsFor(matched).violations).toBe(1);
  });
});
