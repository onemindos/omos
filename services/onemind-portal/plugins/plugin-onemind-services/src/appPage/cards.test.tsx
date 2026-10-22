import { configApiRef } from '@backstage/core-plugin-api';
import type { Entity } from '@backstage/catalog-model';
import { catalogApiRef, entityRouteRef, EntityProvider } from '@backstage/plugin-catalog-react';
import { mockApis, renderInTestApp, TestApiProvider } from '@backstage/test-utils';
import { appStatusApiRef, StatusSnapshot } from '@internal/plugin-onemind-common';
import { screen, waitFor } from '@testing-library/react';
import { KeyLinksCard } from './KeyLinksCard';
import { MembersStatusCard } from './MembersStatusCard';
import { StatusSummaryCard } from './StatusSummaryCard';

const NOW = () => new Date().toISOString();
const snapshot = (): StatusSnapshot => {
  const at = NOW();
  return {
    observedAt: at,
    argo: { state: 'observed', observedAt: at, value: [{ name: 'platform-broken', health: 'Degraded', sync: 'Synced' }] },
    notReady: { state: 'observed', observedAt: at, value: [] },
  };
};
const source = (read: () => Promise<StatusSnapshot> = async () => snapshot()) => ({ read });

const component = (name: string, annotations: Record<string, string>, extra: Record<string, unknown> = {}) => ({
  apiVersion: 'backstage.io/v1alpha1',
  kind: 'Component',
  metadata: { name, namespace: 'default', annotations, ...extra },
  spec: { type: 'service', lifecycle: 'production', owner: 'group:default/onemind-admins', system: 'messaging-nats' },
});

const mountedRoutes = { '/catalog/:namespace/:kind/:name': entityRouteRef };
const config = mockApis.config({ data: { grafana: { domain: 'https://grafana.example.test' } } });

describe('StatusSummaryCard', () => {
  it('shows the live signals behind the dot, and never calls an app healthy that has a failing Argo app', async () => {
    await renderInTestApp(
      <TestApiProvider apis={[[appStatusApiRef, source()]]}>
        <EntityProvider entity={component('broken', { 'argocd/app-name': 'platform-broken' })}>
          <StatusSummaryCard />
        </EntityProvider>
      </TestApiProvider>,
      { mountedRoutes },
    );
    expect(screen.getByText('Live status')).toBeTruthy();
    await waitFor(() => expect(screen.getByText('Argo CD')).toBeTruthy());
    expect(screen.getByText(/Degraded/)).toBeTruthy();
  });

  it('says so when no signal is connected instead of showing a green dot', async () => {
    await renderInTestApp(
      <TestApiProvider apis={[[appStatusApiRef, source()]]}>
        <EntityProvider entity={component('bare', {})}>
          <StatusSummaryCard />
        </EntityProvider>
      </TestApiProvider>,
      { mountedRoutes },
    );
    await waitFor(() => expect(screen.getByText(/No live signal is connected/)).toBeTruthy());
  });

  it('keeps the card and reports it when the status source fails', async () => {
    await renderInTestApp(
      <TestApiProvider apis={[[appStatusApiRef, source(async () => { throw new Error('down'); })]]}>
        <EntityProvider entity={component('x', { 'argocd/app-name': 'platform-x' })}>
          <StatusSummaryCard />
        </EntityProvider>
      </TestApiProvider>,
      { mountedRoutes },
    );
    await waitFor(() => expect(screen.getByRole('alert')).toBeTruthy());
  });
});

describe('KeyLinksCard', () => {
  const annotated = component('nats', {
    'argocd/app-name': 'platform-nats',
    'onemind.io/grafana-dashboard-path': 'd/onemind-workload/onemind-workload?var-namespace=nats',
    'onemind.io/logs-query': '{k8s_namespace_name="nats"}',
    'onemind.io/recovery-runbook': 'apps/nats/',
    'onemind.io/kubernetes-namespaces': 'nats',
  });

  it('lists every generated link', async () => {
    await renderInTestApp(
      <TestApiProvider apis={[[configApiRef, config]]}>
        <EntityProvider entity={annotated}>
          <KeyLinksCard />
        </EntityProvider>
      </TestApiProvider>,
      { mountedRoutes },
    );
    expect(screen.getByText('Key links')).toBeTruthy();
    expect(screen.getByText(/Argo CD · platform-nats/)).toBeTruthy();
    expect(screen.getByText(/Workload dashboard/)).toBeTruthy();
    expect(screen.getByText(/Logs \(Loki/)).toBeTruthy();
    expect(screen.getByText('Runbook')).toBeTruthy();
    expect(screen.getByText(/Cost \(OpenCost\)/)).toBeTruthy();
  });

  it('renders nothing for an entity with no links, so no empty card appears', async () => {
    await renderInTestApp(
      <TestApiProvider apis={[[configApiRef, config]]}>
        <EntityProvider entity={component('bare', {})}>
          <KeyLinksCard />
        </EntityProvider>
      </TestApiProvider>,
      { mountedRoutes },
    );
    expect(screen.queryByText('Key links')).toBeNull();
  });
});

describe('MembersStatusCard', () => {
  const system = {
    apiVersion: 'backstage.io/v1alpha1',
    kind: 'System',
    metadata: { name: 'messaging-nats', namespace: 'default' },
    spec: { owner: 'group:default/onemind-admins' },
  };
  const members = [
    component('good', { 'argocd/app-name': 'platform-good' }, { title: 'Good app' }),
    component('broken', { 'argocd/app-name': 'platform-broken' }, { title: 'Broken app' }),
  ];
  const catalog = (items: unknown[], fail = false) => ({
    getEntities: jest.fn(async () => {
      if (fail) throw new Error('catalog down');
      return { items: items as Entity[] };
    }),
  });

  it('summarises member apps, puts the problem first and keeps unknown separate from healthy', async () => {
    await renderInTestApp(
      <TestApiProvider apis={[[appStatusApiRef, source()], [catalogApiRef, catalog(members)]]}>
        <EntityProvider entity={system}>
          <MembersStatusCard />
        </EntityProvider>
      </TestApiProvider>,
      { mountedRoutes },
    );
    await waitFor(() => expect(screen.getByText('Broken app')).toBeTruthy());
    expect(screen.getByTestId('count-down').textContent).toMatch(/1\s+Down/);
    expect(screen.getByTestId('count-ok').textContent).toMatch(/1\s+Healthy/);
    expect(screen.queryByText(/Every member app reports healthy/)).toBeNull();
  });

  it('says when a System has no member apps', async () => {
    await renderInTestApp(
      <TestApiProvider apis={[[appStatusApiRef, source()], [catalogApiRef, catalog([])]]}>
        <EntityProvider entity={system}>
          <MembersStatusCard />
        </EntityProvider>
      </TestApiProvider>,
      { mountedRoutes },
    );
    await waitFor(() => expect(screen.getByText(/No member apps are recorded here/)).toBeTruthy());
  });

  it('reports a catalog failure instead of an empty or healthy summary', async () => {
    await renderInTestApp(
      <TestApiProvider apis={[[appStatusApiRef, source()], [catalogApiRef, catalog([], true)]]}>
        <EntityProvider entity={system}>
          <MembersStatusCard />
        </EntityProvider>
      </TestApiProvider>,
      { mountedRoutes },
    );
    await waitFor(() => expect(screen.getByRole('alert').textContent).toMatch(/could not be loaded/));
  });
});
