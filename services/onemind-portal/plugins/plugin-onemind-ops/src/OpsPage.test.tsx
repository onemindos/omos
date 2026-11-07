import '@testing-library/jest-dom';
import { fireEvent, screen, waitFor, within } from '@testing-library/react';
import {
  configApiRef,
  discoveryApiRef,
  fetchApiRef,
} from '@backstage/core-plugin-api';
import { catalogApiRef } from '@backstage/plugin-catalog-react';
import {
  mockApis,
  renderInTestApp,
  TestApiProvider,
} from '@backstage/test-utils';
import { LIVE_CHANNELS } from './contract';
import type { LiveSnapshot } from './contract';
import { OpsPage } from './OpsPage';

// Cesium needs WebGL; the page test replaces the whole globe module.
jest.mock('./CesiumGlobe', () => {
  const React = require('react');
  return {
    __esModule: true,
    default: (p: {
      tracks: Array<{ uid: string }>;
      selectedUid?: string;
      onSelect: (uid?: string) => void;
    }) =>
      React.createElement(
        'div',
        { 'data-testid': 'globe-mock', 'data-count': p.tracks.length },
        p.tracks.map((t: { uid: string }) =>
          React.createElement(
            'button',
            { key: t.uid, type: 'button', onClick: () => p.onSelect(t.uid) },
            `pick ${t.uid}`,
          ),
        ),
      ),
  };
});

const signals: Record<string, unknown> = {};
jest.mock('@backstage/plugin-signals-react', () => ({
  useSignal: (channel: string) => ({
    lastSignal: signals[channel] ?? null,
    isSignalsAvailable: true,
  }),
}));

const NOW = Date.now();
const iso = (offsetMs = 0) => new Date(NOW + offsetMs).toISOString();

const snapshot = (over: Partial<LiveSnapshot> = {}): LiveSnapshot => ({
  origin: 'pod-a',
  seq: 5,
  generatedAt: iso(),
  status: {
    nats: { state: 'connected', rate: 42, lastMessageAt: iso(-2000) },
    kubernetes: { state: 'connected', rate: 1.5, lastMessageAt: iso(-1000) },
    argo: { state: 'connected', lastMessageAt: iso(-3000) },
  },
  tracks: [
    {
      uid: 'ae1',
      type: 'a-n-A-C-F',
      callsign: 'UAL123',
      lat: 40,
      lon: -73,
      hae: 10000,
      course: 90,
      speed: 200,
      time: iso(-5000),
      stale: iso(60000),
      subject: 'tak.cot.adsb',
    },
    {
      uid: 'and-1',
      type: 'a-f-G-U-C',
      callsign: 'ALPHA',
      lat: 41,
      lon: -74,
      time: iso(-5000),
      stale: iso(60000),
      subject: 'tak.cot.eud',
    },
    {
      uid: 'hz-1',
      type: 'a-h-G',
      lat: 42,
      lon: -75,
      time: iso(-5000),
      stale: iso(60000),
      subject: 'tak.cot.x',
    },
  ],
  events: [
    {
      id: 'pod-a:4',
      subject: 'tak.cot.adsb.ae1',
      kind: 'tak',
      ts: iso(-1000),
      bytes: 300,
      summary: 'UAL123 a-n-A-C-F',
    },
    {
      id: 'pod-a:3',
      subject: 'tel.node.x.cpu',
      kind: 'tel',
      ts: iso(-2000),
      bytes: 30,
      summary: 'cpu 0.4',
    },
  ],
  cluster: [
    {
      namespace: 'nats',
      total: 3,
      ready: 3,
      pending: 0,
      failed: 0,
      restarts: 0,
      health: 'ok',
    },
    {
      namespace: 'matrix',
      total: 2,
      ready: 1,
      pending: 0,
      failed: 1,
      restarts: 7,
      health: 'down',
      lastChange: 'pod matrix-0 CrashLoopBackOff',
      lastChangeAt: iso(-60000),
    },
  ],
  argo: [
    { name: 'platform-nats', sync: 'Synced', health: 'Healthy' },
    { name: 'platform-matrix', sync: 'OutOfSync', health: 'Degraded' },
  ],
  warnings: [
    {
      namespace: 'matrix',
      object: 'pod/matrix-0',
      reason: 'BackOff',
      message: 'Back-off restarting failed container',
      ts: iso(-4000),
      count: 3,
    },
  ],
  ...over,
});

function renderPage(response: { status: number; body?: unknown } | Error) {
  const fetchMock = jest.fn(async () => {
    if (response instanceof Error) throw response;
    return {
      ok: response.status < 400,
      status: response.status,
      json: async () => response.body,
    } as Response;
  });
  return renderInTestApp(
    <TestApiProvider
      apis={[
        [
          discoveryApiRef,
          mockApis.discovery({ baseUrl: 'http://localhost:7007' }),
        ],
        [fetchApiRef, { fetch: fetchMock } as never],
        [
          catalogApiRef,
          {
            getEntities: async () => ({
              items: [
                {
                  metadata: {
                    name: 'matrix-app',
                    annotations: {
                      'backstage.io/kubernetes-namespace': 'matrix',
                    },
                  },
                },
              ],
            }),
          } as never,
        ],
        [
          configApiRef,
          mockApis.config({
            data: { app: { baseUrl: 'http://localhost:3000' } },
          }),
        ],
      ]}
    >
      <OpsPage />
    </TestApiProvider>,
  ).then(result => ({ ...result, fetchMock }));
}

beforeEach(() => {
  for (const k of Object.keys(signals)) delete signals[k];
});

describe('OpsPage', () => {
  it('renders the connection bar, globe, bus stream, cluster grid, Argo strip and warnings from a snapshot', async () => {
    const { fetchMock } = await renderPage({ status: 200, body: snapshot() });
    await waitFor(() =>
      expect(screen.getByTestId('globe-mock')).toHaveAttribute(
        'data-count',
        '3',
      ),
    );
    expect(fetchMock).toHaveBeenCalledWith(
      'http://localhost:7007/api/livebridge/snapshot',
    );

    const bar = screen.getByRole('region', { name: 'Connections' });
    expect(within(bar).getByText('NATS')).toBeInTheDocument();
    expect(within(bar).getByText('Kubernetes')).toBeInTheDocument();
    expect(within(bar).getByText('Argo')).toBeInTheDocument();
    expect(within(bar).getByText('42 msg/s')).toBeInTheDocument();
    expect(within(bar).getAllByText(/last message \d+ s ago/)).toHaveLength(3);
    expect(
      within(bar).getByRole('status', { name: 'Feed LIVE' }),
    ).toBeInTheDocument();

    expect(screen.getByText('tak.cot.adsb.ae1')).toBeInTheDocument();
    // Unhealthy namespace first, linked to its component's Kubernetes tab.
    const cells = screen.getAllByRole('listitem');
    expect(cells[0]).toHaveTextContent('matrix');
    expect(within(cells[0]).getByRole('link')).toHaveAttribute(
      'href',
      '/catalog/default/component/matrix-app/kubernetes',
    );
    expect(within(cells[1]).queryByRole('link')).toBeNull();

    const argo = screen.getByRole('region', { name: 'Argo applications' });
    expect(argo).toHaveTextContent('1 of 2 Synced and Healthy');
    const chips = argo.querySelectorAll('[data-tone]');
    expect(chips[0]).toHaveTextContent('platform-matrix');
    expect(screen.getByText('BackOff')).toBeInTheDocument();
  });

  it('shows track detail in plain words when a track is picked', async () => {
    await renderPage({ status: 200, body: snapshot() });
    fireEvent.click(await screen.findByText('pick ae1'));
    const detail = screen.getByRole('complementary', { name: 'Track detail' });
    expect(detail).toHaveTextContent('UAL123');
    expect(detail).toHaveTextContent('neutral aircraft (civil fixed wing)');
    expect(detail).toHaveTextContent('389 kt');
    expect(detail).toHaveTextContent('32,808 ft HAE');
    expect(detail).toHaveTextContent('90°');
    expect(detail).toHaveTextContent('tak.cot.adsb');
  });

  it('filters tracks by affiliation and finds a track by callsign', async () => {
    await renderPage({ status: 200, body: snapshot() });
    await screen.findByTestId('globe-mock');
    const filters = screen.getByRole('group', { name: 'Track filters' });
    expect(filters).toHaveTextContent('Friendly1');
    fireEvent.click(within(filters).getByRole('button', { name: /^Hostile/ }));
    await waitFor(() =>
      expect(screen.getByTestId('globe-mock')).toHaveAttribute(
        'data-count',
        '2',
      ),
    );
    fireEvent.change(within(filters).getByLabelText('Search callsign or uid'), {
      target: { value: 'alph' },
    });
    fireEvent.click(within(filters).getByRole('button', { name: /ALPHA/ }));
    expect(
      screen.getByRole('complementary', { name: 'Track detail' }),
    ).toHaveTextContent('friendly ground unit (combat)');
  });

  it('applies a Signals delta that is newer than the snapshot', async () => {
    signals[LIVE_CHANNELS.tracks] = {
      origin: 'pod-a',
      seq: 6,
      ts: iso(),
      upserts: [
        {
          uid: 'new',
          type: 'a-f-A',
          lat: 1,
          lon: 1,
          time: iso(),
          stale: iso(60000),
          subject: 's',
        },
      ],
      removes: ['hz-1'],
    };
    await renderPage({ status: 200, body: snapshot() });
    // 3 in the snapshot, one removed and one added by the delta.
    await waitFor(() =>
      expect(screen.getByTestId('globe-mock')).toHaveAttribute(
        'data-count',
        '3',
      ),
    );
    expect(screen.getByText('pick new')).toBeInTheDocument();
    expect(screen.queryByText('pick hz-1')).toBeNull();
  });

  it('says so when the backend route does not exist', async () => {
    await renderPage({ status: 404 });
    expect(
      await screen.findByText(/The live bridge is not available/),
    ).toBeInTheDocument();
    expect(screen.getByText(/answered 404/)).toBeInTheDocument();
    expect(screen.queryByTestId('globe-mock')).toBeNull();
  });

  it('says so when the snapshot fails', async () => {
    await renderPage(new Error('network down'));
    expect(
      await screen.findByText(/Could not load the live snapshot/),
    ).toBeInTheDocument();
    expect(screen.getByText(/network down/)).toBeInTheDocument();
  });

  it('never shows an unavailable source as zero', async () => {
    await renderPage({
      status: 200,
      body: snapshot({
        status: {
          nats: { state: 'disabled', detail: 'NATS_URL not set' },
          kubernetes: { state: 'down', detail: 'watch failed' },
          argo: { state: 'disabled' },
        },
        tracks: [],
        events: [],
        cluster: [],
        argo: [],
        warnings: [],
      }),
    });
    await screen.findByTestId('globe-mock');
    expect(
      screen.getAllByText(/Tracks unavailable: the NATS source is disabled/)
        .length,
    ).toBeGreaterThan(0);
    expect(
      screen.getByText(
        /Bus events unavailable: NATS is disabled \(NATS_URL not set\)/,
      ),
    ).toBeInTheDocument();
    expect(
      screen.getAllByText(/Kubernetes is down: watch failed/).length,
    ).toBeGreaterThan(0);
    expect(
      screen.getByText(/Argo CD is disabled on the backend/),
    ).toBeInTheDocument();
    expect(screen.queryByText(/0 shown of 0/)).toBeNull();
    expect(screen.queryByText(/Synced and Healthy/)).toBeNull();
    expect(screen.queryByText(/0 down · 0 warn · 0 ok/)).toBeNull();
  });
});
