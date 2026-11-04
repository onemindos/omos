import { render, screen, within } from '@testing-library/react';
import { configApiRef } from '@backstage/core-plugin-api';
import { mockApis, TestApiProvider } from '@backstage/test-utils';
import { LivePage } from './LivePage';
import { buildOwnerIndex } from './owners';
import type { LiveResponse } from './types';

let state: { data?: LiveResponse; error: boolean; push: boolean } = { error: false, push: true };
jest.mock('./useLive', () => ({
  useLive: () => ({ ...state, refresh: jest.fn(), refreshedAt: Date.now() }),
}));
jest.mock('./owners', () => ({
  ...jest.requireActual('./owners'),
  useOwnerIndex: () => ({
    loaded: true,
    index: jest.requireActual('./owners').buildOwnerIndex([
      {
        metadata: {
          name: 'matrix',
          title: 'Matrix',
          annotations: { 'argocd/app-name': 'platform-matrix', 'backstage.io/kubernetes-namespace': 'matrix' },
          labels: { 'onemind.io/pillar': 'le' },
        },
      },
    ]),
  }),
}));

const ok = { status: 'ok' as const, observedAt: '2026-10-22T12:00:00Z', sources: [], problems: [] };
const base = (): LiveResponse => ({
  observedAt: '2026-10-22T12:00:00Z',
  cluster: 'do-nyc3-omos-cloud',
  health: { level: 'warn', headline: '1 app out of sync or progressing, 1 volume almost full', reasons: [] },
  counts: { alerts: 1, argo: 1, workloads: 0, backups: 0, backupsNeedSetup: 1, backupsRetired: 2 },
  panels: { alerts: ok, argo: ok, workloads: ok, nodes: ok, storage: ok, backups: ok, uptime: ok, events: ok },
  alerts: [
    { fingerprint: 'a', name: 'KubePodCrashLooping', severity: 'critical', summary: 'Pod is crash looping', namespace: 'matrix', startsAt: '2026-10-22T11:00:00Z', grafanaUrl: 'https://grafana/x', runbookUrl: 'https://runbooks/x' },
    { fingerprint: 'b', name: 'CPUThrottlingHigh', severity: 'warning', summary: 'Throttled', namespace: 'mystery', startsAt: '2026-10-22T11:30:00Z', grafanaUrl: 'https://grafana/y' },
  ],
  argo: { total: 69, items: [{ name: 'platform-matrix', health: 'Healthy', sync: 'OutOfSync', reasons: ['OutOfSync'], severity: 'warning', url: 'https://argocd/x' }] },
  workloads: { issues: [], notReady: [], total: 176 },
  nodes: [{ name: 'node-1', ready: true, pressure: [], cpuPercent: 13.5, memoryPercent: 59.5, diskPercent: 16, cpuRequestedPercent: 40, memoryRequestedPercent: 50, pods: 40 }],
  storage: { pvcs: [{ namespace: 'onetak', name: 'valhalla-data', usedPercent: 100, fillsWithin4Days: false, severity: 'critical' }], volsync: null },
  backups: {
    volsync: { sources: 46, outOfSync: [] },
    items: [
      { name: 'backstage/backstage-backup', lastSuccess: '2026-10-08T07:15:27Z', suspended: true, state: 'retired', health: 'inactive', attention: 'none', supersededBy: 'backup-postgres-backstage' },
      { name: 'ops/backup-offsite-copy', lastSuccess: null, suspended: true, state: 'needs-setup', health: 'suspended', attention: 'setup', pendingReason: 'Credentials not created' },
    ],
  },
  uptime: [{ name: 'Grafana', state: 'up', scope: 'Application health endpoint', link: 'https://g', entity: 'component:default/grafana', observedAt: '2026-10-22T12:00:00Z', latencyMs: 12 }],
  events: { windowMinutes: 30, items: [] },
  attention: [],
  needsSetup: [],
  links: { argocd: 'https://argocd', headlamp: 'https://headlamp', grafana: 'https://grafana.example' },
});

const renderPage = () =>
  render(
    <TestApiProvider apis={[[configApiRef, mockApis.config({ data: {} })]]}>
      <LivePage />
    </TestApiProvider>,
  );

it('shows the health banner with its reason', () => {
  state = { data: base(), error: false, push: true };
  renderPage();
  const banner = screen.getByRole('status', { name: 'Overall health' });
  expect(banner.textContent).toMatch(/1 volume almost full/);
  expect(screen.getByText(/Live: alert Signals plus a 30 second refresh/)).toBeTruthy();
});

it('groups alerts by pillar and app, and marks unmatched ones Unassigned', () => {
  state = { data: base(), error: false, push: true };
  renderPage();
  const alerts = screen.getByRole('region', { name: 'Firing alerts' });
  expect(within(alerts).getByText('LE · Legacy Evolution')).toBeTruthy();
  expect(within(alerts).getByRole('link', { name: 'Owning app Matrix' }).getAttribute('href')).toBe('/catalog/default/component/matrix');
  expect(within(alerts).getAllByText('Unassigned').length).toBeGreaterThan(0);
  expect(within(alerts).getByRole('link', { name: 'Runbook' }).getAttribute('href')).toBe('https://runbooks/x');
});

it('shows retired backups and the needs-setup job apart from attention', () => {
  state = { data: base(), error: false, push: true };
  renderPage();
  const backups = screen.getByRole('region', { name: 'Backups' });
  expect(within(backups).getByText('Retired')).toBeTruthy();
  expect(within(backups).getByText(/Replaced by backup-postgres-backstage/)).toBeTruthy();
  expect(within(backups).getByText('Needs setup')).toBeTruthy();
  expect(within(backups).getByText(/Credentials not created/)).toBeTruthy();
  expect(within(backups).getByText(/0 to review · 1 needs setup · 2 retired/)).toBeTruthy();
});

it('marks a failed source Unavailable and still renders the other panels', () => {
  const data = base();
  data.panels.storage = { status: 'error', observedAt: null, sources: ['prometheus'], problems: [{ source: 'prometheus', status: 'error', error: 'Timed out' }] };
  data.storage = null;
  state = { data, error: false, push: false };
  renderPage();
  const storage = screen.getByRole('region', { name: 'Storage' });
  expect(within(storage).getByText('Unavailable')).toBeTruthy();
  expect(within(storage).getByRole('alert').textContent).toMatch(/Timed out/);
  expect(within(screen.getByRole('region', { name: 'Argo CD' })).getByText('platform-matrix')).toBeTruthy();
  expect(screen.getByText(/Refreshing every 30 seconds/)).toBeTruthy();
});

it('says it cannot confirm health when nothing has loaded and the request failed', () => {
  state = { data: undefined, error: true, push: false };
  renderPage();
  expect(screen.getByRole('alert').textContent).toMatch(/cannot confirm anything is healthy/);
});

it('keeps the last results visible and says the refresh failed', () => {
  state = { data: base(), error: true, push: true };
  renderPage();
  expect(screen.getByText(/Live refresh failed\. Showing results observed at/)).toBeTruthy();
  expect(buildOwnerIndex([]).resolve({})).toBeUndefined();
});

it('opens Grafana panels as links while embedding is off', () => {
  state = { data: base(), error: false, push: true };
  renderPage();
  expect(document.querySelectorAll('iframe')).toHaveLength(0);
  expect(screen.getAllByRole('link', { name: 'Open in Grafana' })[0].getAttribute('href')).toMatch(
    /^https:\/\/grafana\.example\/d\/efa86fd1d0c121a26444b636a3f509a8\//,
  );
});

it('embeds the five panels, sandboxed, when embedding is enabled', () => {
  state = { data: base(), error: false, push: true };
  render(
    <TestApiProvider apis={[[configApiRef, mockApis.config({ data: { onemind: { grafanaEmbedding: { enabled: true } } } })]]}>
      <LivePage />
    </TestApiProvider>,
  );
  const frames = document.querySelectorAll('iframe');
  expect(frames).toHaveLength(5);
  expect(frames[0].getAttribute('src')).toMatch(/^https:\/\/grafana\.example\/d-solo\/.*panelId=7.*theme=light/);
  expect(frames[0].getAttribute('sandbox')).toBe('allow-scripts allow-same-origin');
});
