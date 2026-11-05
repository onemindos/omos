import { parseNamespaceAllocation } from './allocation';
import {
  appsFromEntities,
  buildCostLines,
  MIXED,
  pillarRollups,
  rollup,
  UNMAPPED,
  UNSET,
} from './mapping';
import fixture from './__fixtures__/allocation-7d.json';

const window = { start: '2026-10-29T00:00:00Z', end: '2026-11-05T00:00:00Z' };
const entity = (
  name: string,
  namespace: string,
  pillar?: string,
  review = 'approved',
  extra: Record<string, string> = {},
) => ({
  metadata: {
    name,
    annotations: {
      'backstage.io/kubernetes-namespace': namespace,
      'onemind.io/pillar-review': review,
      ...extra,
    },
    labels: (pillar ? { 'onemind.io/pillar': pillar } : {}) as Record<
      string,
      string
    >,
  },
});

const rows = () => parseNamespaceAllocation(fixture, { window }).rows;

describe('appsFromEntities', () => {
  it('reads singular and plural namespace annotations, skipping components without any', () => {
    const apps = appsFromEntities([
      entity('a', 'nats', 'so', 'approved', {
        'onemind.io/kubernetes-namespaces': 'nats, hermes',
      }),
      { metadata: { name: 'no-namespace' } },
    ]);
    expect(apps).toHaveLength(1);
    expect(apps[0].namespaces).toEqual(['nats', 'hermes']);
    expect(apps[0].href).toBe('/catalog/default/component/a');
  });
});

describe('buildCostLines', () => {
  const apps = appsFromEntities([
    entity('hermes-agent', 'hermes', 'le', 'proposed'),
    entity('tak', 'takserver', 'so'),
    entity('nats-a', 'nats', 'so'),
    entity('nats-b', 'nats', 'ge'),
    entity('prom', 'monitoring'),
  ]);
  const lines = buildCostLines(rows(), apps);
  const byNs = (name: string) => lines.find(line => line.namespace === name)!;

  it('attributes a single-app namespace to that app', () => {
    expect(byNs('takserver')).toMatchObject({
      attribution: 'app',
      pillar: 'so',
      pillarProposed: false,
    });
  });
  it('does not split a shared namespace and marks mixed pillars', () => {
    expect(byNs('nats')).toMatchObject({
      attribution: 'shared',
      pillar: MIXED,
    });
    expect(byNs('nats').cost).toBeCloseTo(6.14, 1);
  });
  it('flags proposed pillars and apps without a pillar', () => {
    expect(byNs('hermes').pillarProposed).toBe(true);
    expect(byNs('monitoring').pillar).toBe(UNSET);
  });
  it('keeps namespaces with no catalog app visible', () => {
    expect(byNs('nextcloud')).toMatchObject({
      attribution: 'none',
      pillar: UNMAPPED,
    });
  });
  it('sorts by cost, highest first', () => {
    const costs = lines.map(line => line.cost ?? -1);
    expect(costs).toEqual([...costs].sort((a, b) => b - a));
  });
});

describe('rollups never count unknown cost as zero', () => {
  it('reports unknown namespaces next to a lower-bound sum', () => {
    const payload: any = JSON.parse(JSON.stringify(fixture));
    delete payload.data[0].takserver.totalCost;
    const parsed = parseNamespaceAllocation(payload, { window });
    const tak = parsed.rows.find(row => row.namespace === 'takserver')!;
    expect(tak.totalCost).toBeNull();
    expect(tak.state).toBe('partial');
    const apps = appsFromEntities([entity('tak', 'takserver', 'so')]);
    const lines = buildCostLines(parsed.rows, apps);
    expect(lines.find(l => l.namespace === 'takserver')!.cost).toBeNull();
    expect(rollup(lines).unknown).toBe(1);
    const so = pillarRollups(lines).find(item => item.pillar === 'so')!;
    expect(so.unknown).toBe(1);
    expect(so.known).toBe(0);
  });
  it('orders pillars with the four pillars first', () => {
    const apps = appsFromEntities([
      entity('x', 'nats', 'ge'),
      entity('y', 'takserver', 'so'),
    ]);
    const order = pillarRollups(buildCostLines(rows(), apps)).map(
      item => item.pillar,
    );
    expect(order.slice(0, 2)).toEqual(['so', 'ge']);
    expect(order[order.length - 1]).toBe(UNMAPPED);
  });
});
