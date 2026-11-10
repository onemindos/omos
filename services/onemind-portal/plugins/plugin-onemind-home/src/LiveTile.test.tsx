import { render, screen, waitFor } from '@testing-library/react';
import { MemoryRouter } from 'react-router-dom';
import { LiveTile } from './LiveTile';

// If the tile ever pulled in the globe, this would blow up on import.
jest.mock('cesium', () => {
  throw new Error('cesium must not load on the home page');
});

const mockFetch = jest.fn();
jest.mock('@backstage/core-plugin-api', () => {
  const actual = jest.requireActual('@backstage/core-plugin-api');
  return {
    ...actual,
    useApi: (ref: unknown) =>
      ref === actual.discoveryApiRef
        ? { getBaseUrl: async () => 'http://x/api/livebridge' }
        : { fetch: mockFetch },
  };
});

const signals: Record<string, unknown> = {};
jest.mock('@backstage/plugin-signals-react', () => ({
  useSignal: (channel: string) => ({
    lastSignal: signals[channel] ?? null,
    isSignalsAvailable: true,
  }),
}));

const iso = (offsetMs = 0) => new Date(Date.now() + offsetMs).toISOString();

const snapshot = (over: Record<string, unknown> = {}) => ({
  origin: 'pod-a',
  seq: 5,
  generatedAt: iso(),
  status: {
    nats: { state: 'connected', lastMessageAt: iso(-1000) },
    kubernetes: { state: 'connected' },
    argo: { state: 'connected' },
  },
  tracks: [
    {
      uid: 't1',
      type: 'a-f-G-U-C',
      lat: 1,
      lon: 2,
      time: iso(),
      stale: iso(60000),
      subject: 'tak.cot.x',
    },
    {
      uid: 't2',
      type: 'a-n-A-C-F',
      lat: 1,
      lon: 2,
      time: iso(),
      stale: iso(60000),
      subject: 'tak.cot.y',
    },
  ],
  events: [
    {
      id: 'pod-a-1',
      subject: 'evt.mission.p1.completed',
      kind: 'evt',
      ts: iso(-1000),
      bytes: 10,
      summary: 'Patrol 1 completed',
    },
    {
      id: 'pod-a-2',
      subject: 'tel.a.b.c',
      kind: 'tel',
      ts: iso(-2000),
      bytes: 10,
      summary: 'Telemetry sample',
    },
  ],
  cluster: [
    {
      namespace: 'nats',
      total: 5,
      ready: 5,
      pending: 0,
      failed: 0,
      restarts: 0,
      health: 'ok',
    },
    {
      namespace: 'hermes',
      total: 3,
      ready: 2,
      pending: 1,
      failed: 0,
      restarts: 0,
      health: 'warn',
    },
  ],
  argo: [],
  warnings: [],
  ...over,
});

const answer = (body: unknown, status = 200) =>
  mockFetch.mockResolvedValue({
    ok: status >= 200 && status < 300,
    status,
    json: async () => body,
  });

beforeEach(() => {
  mockFetch.mockReset();
});

test('connected: shows counts, latest events, cluster line and the Ops link', async () => {
  answer(snapshot());
  render(<LiveTile />, { wrapper: MemoryRouter });
  expect(await screen.findByText('Connected')).toBeTruthy();
  expect(screen.getByText('tracks').previousSibling?.textContent).toBe('2');
  expect(
    screen.getByText('events, last 5 min').previousSibling?.textContent,
  ).toBe('2');
  expect(screen.getByText('pods ready').previousSibling?.textContent).toBe(
    '7/8',
  );
  expect(screen.getByText('1 namespace needs attention')).toBeTruthy();
  expect(screen.getByText('Patrol 1 completed')).toBeTruthy();
  const link = screen.getByRole('link', { name: 'Open Ops Center' });
  expect(link.getAttribute('href')).toBe('/ops');
});

test('empty bridge shows a calm message and still links to /ops', async () => {
  answer(snapshot({ tracks: [], events: [], cluster: [] }));
  render(<LiveTile />, { wrapper: MemoryRouter });
  expect(await screen.findByText(/Nothing on the bridge yet/)).toBeTruthy();
  expect(screen.getByRole('link', { name: 'Open Ops Center' })).toBeTruthy();
});

test('backend error shows a quiet unavailable state', async () => {
  answer({}, 500);
  render(<LiveTile />, { wrapper: MemoryRouter });
  expect(await screen.findByText('Live bridge unavailable')).toBeTruthy();
  expect(screen.getByRole('link', { name: 'Open Ops Center' })).toBeTruthy();
});

test('livebridge not deployed (404) is also unavailable, not a crash', async () => {
  answer({}, 404);
  render(<LiveTile />, { wrapper: MemoryRouter });
  expect(await screen.findByText('Live bridge unavailable')).toBeTruthy();
});

test('network failure is unavailable', async () => {
  mockFetch.mockRejectedValue(new Error('down'));
  render(<LiveTile />, { wrapper: MemoryRouter });
  await waitFor(() =>
    expect(screen.getByText('Live bridge unavailable')).toBeTruthy(),
  );
});

test('shows loading before the snapshot arrives', () => {
  mockFetch.mockReturnValue(new Promise(() => {}));
  render(<LiveTile />, { wrapper: MemoryRouter });
  expect(screen.getByText(/Loading live picture/)).toBeTruthy();
});

test('NATS down reads as offline', async () => {
  answer(
    snapshot({
      status: {
        nats: { state: 'down' },
        kubernetes: { state: 'connected' },
        argo: { state: 'connected' },
      },
    }),
  );
  render(<LiveTile />, { wrapper: MemoryRouter });
  expect(await screen.findByText('Offline')).toBeTruthy();
});

test('the ops entry point loads without cesium', () => {
  // The cesium mock above throws if anything requires it.
  expect(() =>
    jest.requireActual('@internal/plugin-onemind-ops'),
  ).not.toThrow();
});
