import { act, render, screen } from '@testing-library/react';
import { discoveryApiRef, useApi } from '@backstage/core-plugin-api';
import { StacOperations, isStacEntity } from './StacOperations';

jest.mock('@backstage/core-plugin-api', () => ({
  ...jest.requireActual('@backstage/core-plugin-api'),
  useApi: jest.fn(),
}));
const fetch = jest.fn();
const discovery = {
  getBaseUrl: async () => 'https://portal.example.test/api/onemind-services',
};
const fetchApi = { fetch };
afterEach(() => jest.useRealTimers());
beforeEach(() => {
  fetch.mockReset();
  (useApi as jest.Mock).mockImplementation(ref =>
    ref === discoveryApiRef ? discovery : fetchApi,
  );
});
const fixture = () => ({
  schemaVersion: 'onemind.stac-evidence.v1',
  observedAt: new Date().toISOString(),
  scope: 'Read-only Development pilot',
  partial: true,
  workloads: [],
  probes: {
    apiLanding: {
      state: 'responding',
      httpStatus: 200,
      observedAt: new Date().toISOString(),
    },
    collections: {
      state: 'error',
      httpStatus: 500,
      observedAt: new Date().toISOString(),
    },
  },
  data: {
    state: 'unknown',
    observedAt: new Date().toISOString(),
    collectionsReturned: null,
    collectionsTruncated: null,
    sampleItemsReturned: null,
    latestItemDatetime: null,
    freshness: 'No freshness conclusion can be drawn.',
  },
  argo: {
    state: 'observed',
    observedAt: new Date().toISOString(),
    applications: [],
  },
  recovery: { state: 'not-certified', detail: 'Not a restore test' },
});
test('shows real query failure and unknown freshness instead of a green empty catalog', async () => {
  fetch.mockResolvedValue({ ok: true, json: async () => fixture() });
  render(<StacOperations />);
  expect(await screen.findByText(/Incomplete or failing checks/)).toBeTruthy();
  expect(screen.getByText(/collections: error/).textContent).toContain('500');
  expect(screen.getByText(/Collections returned:/).textContent).toContain(
    'Unknown',
  );
  expect(
    screen.getByText(/Newest queried item acquisition time:/).textContent,
  ).toContain('Unknown');
  expect(
    screen.getByRole('link', { name: 'Live pilot logs' }).getAttribute('href'),
  ).toContain('runtime-logs');
  expect(
    screen.queryByRole('button', {
      name: /restart|promote|migrate|initialize/i,
    }),
  ).toBeNull();
});
test('does not replace a failed observation with healthy zeroes', async () => {
  fetch.mockRejectedValue(new Error('Denied'));
  render(<StacOperations />);
  expect(await screen.findByText(/Observation unavailable/)).toBeTruthy();
  expect(screen.queryByText(/Collections returned:/)).toBeNull();
});
test('a hung request becomes unavailable instead of remaining loading forever', async () => {
  jest.useFakeTimers();
  fetch.mockImplementation(() => new Promise(() => {}));
  render(<StacOperations />);
  await act(async () => {
    await Promise.resolve();
  });
  await act(async () => {
    await jest.advanceTimersByTimeAsync(15_001);
  });
  expect(screen.getByText(/Observation unavailable/)).toBeTruthy();
  expect(screen.queryByText(/Loading service evidence/)).toBeNull();
});
test('a timed-out refresh retains previous evidence only as stale and can retry', async () => {
  jest.useFakeTimers();
  fetch.mockResolvedValueOnce({ ok: true, json: async () => fixture() });
  render(<StacOperations />);
  await act(async () => {
    await Promise.resolve();
  });
  expect(screen.getByText(/Collections returned:/)).toBeTruthy();
  fetch.mockImplementation(() => new Promise(() => {}));
  await act(async () => {
    await jest.advanceTimersByTimeAsync(60_000);
  });
  await act(async () => {
    await jest.advanceTimersByTimeAsync(15_001);
  });
  expect(screen.getByText(/Stale observation/)).toBeTruthy();
  expect(screen.getByText(/Collections returned:/)).toBeTruthy();
  fetch.mockResolvedValue({ ok: true, json: async () => fixture() });
  await act(async () => {
    await jest.advanceTimersByTimeAsync(60_000);
  });
  expect(screen.queryByText(/Observation unavailable/)).toBeNull();
  expect(screen.queryByText(/Stale observation/)).toBeNull();
});
test('rejects an unknown observation contract', async () => {
  fetch.mockResolvedValue({
    ok: true,
    json: async () => ({ ...fixture(), schemaVersion: 'unknown' }),
  });
  render(<StacOperations />);
  expect(await screen.findByText(/Observation unavailable/)).toBeTruthy();
});
test.each([
  ['Component', 'stac-browser', true],
  ['System', 'onetak', true],
  ['Component', 'onetak', false],
  ['System', 'other', false],
  ['API', 'stac-browser', false],
])(
  'selects only reviewed managed pilot pages: %s/%s',
  (kind, name, expected) => {
    expect(
      isStacEntity({
        kind,
        metadata: { name, labels: { 'openchoreo.io/managed': 'true' } },
      }),
    ).toBe(expected);
  },
);
test('unmanaged names and other namespaces cannot opt into the pilot page', () => {
  expect(
    isStacEntity({ kind: 'Component', metadata: { name: 'stac-browser' } }),
  ).toBe(false);
  expect(
    isStacEntity({
      kind: 'Component',
      metadata: {
        name: 'stac-browser',
        namespace: 'other',
        labels: { 'openchoreo.io/managed': 'true' },
      },
    }),
  ).toBe(false);
});
