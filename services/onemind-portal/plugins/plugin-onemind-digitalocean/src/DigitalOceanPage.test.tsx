import {
  fireEvent,
  render,
  screen,
  waitFor,
  within,
} from '@testing-library/react';
import { discoveryApiRef, fetchApiRef } from '@backstage/core-plugin-api';
import { TestApiProvider } from '@backstage/test-utils';
import { DigitalOceanPage } from './DigitalOceanPage';
import { findDoResources } from './RunsOnCard';
import type { Summary } from './types';

const ok = <T,>(items: T[]) => ({ state: 'ok' as const, items });
const summary = (): Summary => ({
  state: 'connected',
  generatedAt: '2026-11-05T06:00:00Z',
  expiresAt: null,
  account: { state: 'ok' },
  billing: {
    state: 'ok',
    usageBeforeCredits: '776.67',
    amountDue: '0.03',
    accountBalance: '0.03',
    projectedMonth: '800.00',
    daysElapsed: 20,
    daysInMonth: 31,
    invoices: [{ period: '2026-10', amount: '0.03' }],
  },
  kubernetes: ok([
    {
      id: 'c1',
      name: 'omos-cloud',
      region: 'nyc3',
      version: '1.36.5-do.1',
      highAvailability: true,
      upgradeAvailable: '1.36.5-do.3',
      upgradeCheck: 'ok' as const,
      nodePools: [
        {
          id: 'p1',
          name: 'core',
          size: 's-8vcpu-16gb',
          count: 3,
          autoScale: false,
        },
      ],
    },
  ]),
  droplets: ok([
    {
      id: '1',
      name: 'node-a',
      status: 'active',
      size: 's-8vcpu-16gb',
      region: 'nyc3',
      tags: [],
    },
  ]),
  volumes: {
    ...ok([
      {
        id: 'v1',
        name: 'pvc-1',
        sizeGb: 10,
        region: 'nyc3',
        attachedDropletIds: ['1'],
        unattached: false,
        pvc: { namespace: 'n8n', name: 'n8n-data' },
      },
      {
        id: 'v2',
        name: 'pvc-2',
        sizeGb: 5,
        region: 'nyc3',
        attachedDropletIds: [],
        unattached: true,
      },
    ]),
    unattachedCount: 1,
    unattachedSizeGb: 5,
  },
  snapshots: {
    ...ok([]),
    stats: {
      count: 2,
      totalSizeGb: 3,
      olderThanDays: 90,
      olderCount: 1,
      oldest: {
        id: 's1',
        name: 'snap-old',
        resourceType: 'volume',
        sizeGb: 1,
        regions: [],
        ageDays: 120,
        createdAt: '2026-07-07T00:00:00Z',
      },
      largest: {
        id: 's2',
        name: 'snap-big',
        resourceType: 'volume',
        sizeGb: 2,
        regions: [],
        ageDays: 3,
      },
    },
  },
  loadBalancers: ok([]),
  firewalls: ok([]),
  vpcs: ok([]),
  domains: {
    state: 'unavailable',
    reason: 'The token does not have the read scope for this section',
    items: [],
  },
  uptime: ok([
    {
      id: 'u1',
      name: 'gateway',
      type: 'https',
      target: 'https://x',
      enabled: true,
      regions: [],
      state: 'up' as const,
      regionStates: {},
    },
  ]),
  alerts: ok([]),
  spaces: {
    state: 'not-available',
    reason:
      'The DigitalOcean API cannot list Spaces buckets with an API token.',
    items: [],
  },
  changes: [
    {
      at: '2026-11-05T05:00:00Z',
      added: [{ type: 'volume', id: 'v2', name: 'pvc-2' }],
      removed: [],
    },
  ],
});

const renderPage = (response: Partial<Response>) =>
  render(
    <TestApiProvider
      apis={[
        [
          discoveryApiRef,
          { getBaseUrl: async () => 'http://x/api/onemind-digitalocean' },
        ],
        [fetchApiRef, { fetch: jest.fn(async () => response) } as never],
      ]}
    >
      <DigitalOceanPage />
    </TestApiProvider>,
  );

describe('DigitalOceanPage', () => {
  it('shows usage before credits apart from the amount due, the projection and invoices', async () => {
    renderPage({ ok: true, status: 200, json: async () => summary() });
    await screen.findByText('Usage before credits (this month)');
    expect(screen.getByText('$776.67')).toBeTruthy();
    expect(screen.getByText('Amount due (this month)')).toBeTruthy();
    expect(screen.getAllByText('$0.03').length).toBeGreaterThan(0);
    expect(screen.getByText('$800.00')).toBeTruthy();
    expect(screen.getByText('Available: 1.36.5-do.3')).toBeTruthy();
    expect(screen.getByText(/Not available/)).toBeTruthy(); // Spaces
    expect(screen.getByText(/does not have the read scope/)).toBeTruthy(); // domains, on its own
    expect(screen.getByText(/added volume: pvc-2/)).toBeTruthy();
    expect(screen.queryByRole('button', { name: /delete/i })).toBeNull();
  });

  it('filters volumes to the unattached ones', async () => {
    renderPage({ ok: true, status: 200, json: async () => summary() });
    const volumes = await screen.findByRole('region', { name: 'Volumes' });
    expect(within(volumes).getByText('n8n/n8n-data')).toBeTruthy();
    fireEvent.click(
      within(volumes).getByRole('button', { name: 'Unattached only' }),
    );
    expect(within(volumes).queryByText('n8n/n8n-data')).toBeNull();
    expect(within(volumes).getByText('pvc-2')).toBeTruthy();
  });

  it('shows the not connected state with the token instruction', async () => {
    renderPage({
      ok: true,
      status: 200,
      json: async () => ({ ...summary(), state: 'not-connected' }),
    });
    expect(
      await screen.findByText(/needs the read-only DigitalOcean token/),
    ).toBeTruthy();
  });

  it('tells a non-administrator the page is for administrators', async () => {
    renderPage({ ok: false, status: 403, json: async () => ({}) });
    await waitFor(() =>
      expect(screen.getByText(/Only OneMind administrators/)).toBeTruthy(),
    );
  });
});

describe('findDoResources', () => {
  const res = (name: string, type: string, depends: string[] = []) => ({
    apiVersion: 'backstage.io/v1alpha1',
    kind: 'Resource',
    metadata: { name },
    spec: { type },
    relations: depends.map(d => ({ type: 'dependsOn', targetRef: d })),
  });
  it('follows dependsOn from an app through the cluster Resource to the node pools', async () => {
    const all: Record<string, ReturnType<typeof res>> = {
      'resource:default/do-volume-1': res('do-volume-1', 'do-volume'),
      'resource:default/do-nyc3-omos-cloud': res(
        'do-nyc3-omos-cloud',
        'kubernetes-cluster',
        ['resource:default/do-kubernetes-cluster-1'],
      ),
      'resource:default/do-kubernetes-cluster-1': res(
        'do-kubernetes-cluster-1',
        'do-kubernetes-cluster',
        ['resource:default/do-node-pool-1'],
      ),
      'resource:default/do-node-pool-1': res('do-node-pool-1', 'do-node-pool'),
    };
    const app = {
      apiVersion: 'backstage.io/v1alpha1',
      kind: 'Component',
      metadata: { name: 'x' },
      spec: {},
      relations: [
        { type: 'dependsOn', targetRef: 'resource:default/do-volume-1' },
        { type: 'dependsOn', targetRef: 'resource:default/do-nyc3-omos-cloud' },
      ],
    };
    const found = await findDoResources(app, async refs =>
      refs.map(r => all[r] as never),
    );
    expect(found.map(e => e.metadata.name).sort()).toEqual([
      'do-kubernetes-cluster-1',
      'do-node-pool-1',
      'do-volume-1',
    ]);
  });
});
