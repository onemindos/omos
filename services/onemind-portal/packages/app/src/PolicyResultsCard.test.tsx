import { render, screen, waitFor } from '@testing-library/react';
import { discoveryApiRef, fetchApiRef } from '@backstage/core-plugin-api';
import { TestApiProvider } from '@backstage/test-utils';
import PolicyResultsCard from './PolicyResultsCard';

let entity: { metadata: { name: string; annotations: Record<string, string> } };
jest.mock('@backstage/plugin-catalog-react', () => ({
  useEntity: () => ({ entity }),
}));

const tally = (pass: number) => ({ pass, fail: 0, warn: 0, error: 0, skip: 1 });
const body = (extra: Record<string, unknown> = {}) => ({
  status: 'ok',
  observedAt: '2026-11-05T12:00:00Z',
  reports: 3,
  truncated: false,
  namespaces: {
    backstage: {
      resources: [
        {
          kind: 'Deployment',
          name: 'backstage',
          workload: 'backstage',
          summary: tally(3),
        },
        {
          kind: 'Pod',
          name: 'backstage-5d9f-abcde',
          workload: 'backstage',
          summary: tally(3),
        },
      ],
      findings: [],
    },
  },
  ...extra,
});

const renderCard = (respond: () => Promise<unknown>) =>
  render(
    <TestApiProvider
      apis={[
        [discoveryApiRef, { getBaseUrl: async () => 'http://x/api/onemind' }],
        [fetchApiRef, { fetch: jest.fn(respond) } as any],
      ]}
    >
      <PolicyResultsCard />
    </TestApiProvider>,
  );
const ok = (value: unknown) => async () => ({
  ok: true,
  json: async () => value,
});

beforeEach(() => {
  entity = {
    metadata: {
      name: 'backstage',
      annotations: {
        'backstage.io/kubernetes-namespace': 'backstage',
        'onemind.io/kubernetes-workloads': 'backstage/backstage',
      },
    },
  };
});

it('shows zero violations with the scanned-resource count and the observation time', async () => {
  renderCard(ok(body()));
  expect(await screen.findByText(/0 policy violations/)).toBeTruthy();
  expect(
    screen.getByText(/6 checks passed, 2 skipped, across 2 scanned/),
  ).toBeTruthy();
  expect(screen.getByText(/Observed 2026-11-05T12:00:00Z/)).toBeTruthy();
  expect(screen.getByText(/not a full security assessment/)).toBeTruthy();
});

it('shows a violation once, with the failing policy listed', async () => {
  const response: any = body();
  response.namespaces.backstage.findings = [
    {
      resource: {
        kind: 'Deployment',
        name: 'backstage',
        workload: 'backstage',
      },
      policy: 'disallow-latest-image-tag',
      rule: 'require-image-tag',
      result: 'fail',
      severity: 'medium',
      message: 'An image tag is required.',
    },
    {
      resource: {
        kind: 'Pod',
        name: 'backstage-5d9f-abcde',
        workload: 'backstage',
      },
      policy: 'disallow-latest-image-tag',
      rule: 'require-image-tag',
      result: 'fail',
      severity: 'medium',
      message: 'An image tag is required.',
    },
  ];
  renderCard(ok(response));
  expect(await screen.findByText(/1 policy violation$/)).toBeTruthy();
  const list = screen.getByRole('list', { name: 'Policy findings' });
  expect(list.textContent).toMatch(/disallow-latest-image-tag/);
});

it('says unavailable, not zero, when the backend call fails', async () => {
  renderCard(async () => {
    throw new Error('network');
  });
  expect((await screen.findByRole('alert')).textContent).toMatch(/unavailable/);
  expect(screen.queryByText(/0 policy violations/)).toBeNull();
});

it('says unavailable when the backend reports a source error', async () => {
  renderCard(
    ok({
      status: 'error',
      observedAt: null,
      error: 'Cluster API returned 403',
    }),
  );
  expect((await screen.findByRole('alert')).textContent).toMatch(/unavailable/);
  expect(screen.queryByText(/policy violation/)).toBeNull();
});

it('warns when the data is stale', async () => {
  renderCard(ok(body({ status: 'stale', error: 'Timed out' })));
  await waitFor(() =>
    expect(screen.getByRole('alert').textContent).toMatch(/last good data/),
  );
});

it('does not call a workload with no report passing', async () => {
  entity.metadata.annotations['onemind.io/kubernetes-workloads'] =
    'backstage/ghost';
  renderCard(ok(body()));
  expect(
    await screen.findByText(/no reports for these workloads yet/),
  ).toBeTruthy();
  expect(screen.getByText(/No report yet for: backstage\/ghost/)).toBeTruthy();
  expect(screen.queryByText(/0 policy violations/)).toBeNull();
});

it('marks counts as lower bounds when the backend truncated findings', async () => {
  renderCard(ok(body({ truncated: true })));
  expect(await screen.findByText(/at least 0 policy violations/)).toBeTruthy();
});
