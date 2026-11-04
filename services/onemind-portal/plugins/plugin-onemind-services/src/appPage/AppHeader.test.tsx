import { configApiRef } from '@backstage/core-plugin-api';
import { entityRouteRef, EntityProvider, starredEntitiesApiRef } from '@backstage/plugin-catalog-react';
import { mockApis, renderInTestApp, TestApiProvider } from '@backstage/test-utils';
import { appStatusApiRef, StatusSnapshot } from '@internal/plugin-onemind-common';
import { screen, waitFor } from '@testing-library/react';
import { AppHeader } from './AppHeader';

const at = () => new Date().toISOString();
const snapshot = (): StatusSnapshot => ({
  observedAt: at(),
  argo: { state: 'observed', observedAt: at(), value: [] },
  notReady: { state: 'observed', observedAt: at(), value: [] },
});
const mountedRoutes = { '/catalog/:namespace/:kind/:name': entityRouteRef };
type Observer = ((value: Set<string>) => void) | { next: (value: Set<string>) => void };
const starred = {
  starredEntitie$: () => ({
    subscribe: (observer: Observer) => {
      const next = typeof observer === 'function' ? observer : observer.next.bind(observer);
      next(new Set());
      return { unsubscribe() {}, closed: false };
    },
    [Symbol.observable]() {
      return this;
    },
  }),
  toggleStarred: async () => {},
};

const entity = (annotations: Record<string, string>, labels: Record<string, string> = {}) => ({
  apiVersion: 'backstage.io/v1alpha1',
  kind: 'Component',
  metadata: { name: 'nats', title: 'NATS', namespace: 'default', annotations, labels },
  spec: { type: 'service', lifecycle: 'production', owner: 'group:default/onemind-admins', system: 'messaging-nats' },
  relations: [{ type: 'ownedBy', targetRef: 'group:default/onemind-admins' }],
});

async function renderHeader(e: ReturnType<typeof entity>) {
  return renderInTestApp(
    <TestApiProvider
      apis={[
        [appStatusApiRef, { read: async () => snapshot() }],
        [starredEntitiesApiRef, starred as never],
        [configApiRef, mockApis.config({ data: {} })],
      ]}
    >
      <EntityProvider entity={e}>
        <AppHeader tabs={[{ id: 'overview', label: 'Overview', href: '/catalog/default/component/nats' }]} activeTabId="overview" />
      </EntityProvider>
    </TestApiProvider>,
    { mountedRoutes },
  );
}

describe('AppHeader', () => {
  it('shows the name, kind, pillar and hosting badges, status and an Open app button', async () => {
    await renderHeader(
      entity(
        { 'argocd/app-name': 'platform-nats', 'onemind.io/hosting': 'digitalocean-nyc3', 'onemind.io/launcher-url': 'https://nats.example.test' },
        { 'onemind.io/pillar': 'so' },
      ),
    );
    expect(screen.getByText('NATS')).toBeTruthy();
    expect(screen.getByText('Component')).toBeTruthy();
    expect(screen.getByText('SO')).toBeTruthy();
    expect(screen.getByText('K8s')).toBeTruthy();
    expect(screen.getByRole('link', { name: 'Open app' }).getAttribute('href')).toBe('https://nats.example.test/');
    await waitFor(() => expect(screen.getByText('Status')).toBeTruthy());
  });

  it('never promotes an unrecorded app: the hosting badge says Unknown and there is no Open app button', async () => {
    await renderHeader(entity({}));
    // Two honest unknowns: the hosting badge, and the status dot (no live signal is connected).
    await waitFor(() => expect(screen.getAllByText('Unknown')).toHaveLength(2));
    expect(screen.queryByRole('link', { name: 'Open app' })).toBeNull();
  });

  it('keeps an unsafe launcher URL off the button', async () => {
    await renderHeader(entity({ 'onemind.io/launcher-url': 'javascript:alert(1)' }));
    expect(screen.queryByRole('link', { name: 'Open app' })).toBeNull();
  });
});
