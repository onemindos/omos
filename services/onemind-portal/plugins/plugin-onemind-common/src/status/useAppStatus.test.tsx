import { render, screen, act } from '@testing-library/react';
import { useAppStatus, useAppStatuses } from './useAppStatus';
import { StatusSnapshot } from './model';

let source: { read: jest.Mock };
jest.mock('@backstage/core-plugin-api', () => ({
  ...jest.requireActual('@backstage/core-plugin-api'),
  useApi: () => source,
}));
jest.mock('@backstage/plugin-catalog-react', () => ({ useEntity: jest.fn() }));

const fresh = (): StatusSnapshot => {
  const observedAt = new Date().toISOString();
  return {
    observedAt,
    argo: { state: 'observed', observedAt, value: [{ name: 'platform-x', health: 'Degraded', sync: 'Synced' }] },
    notReady: { state: 'observed', observedAt, value: [] },
  };
};
const entity = (name: string, argo: string) => ({ kind: 'Component', metadata: { name, annotations: { 'argocd/app-name': argo } } });

function One({ e }: { e: ReturnType<typeof entity> }) {
  const s = useAppStatus(e);
  return <div data-testid={e.metadata.name}>{s.loading ? 'loading' : s.status}</div>;
}
function Many({ list }: { list: ReturnType<typeof entity>[] }) {
  const { results } = useAppStatuses(list);
  return <div data-testid="many">{results.map(r => `${r.entity.metadata.name}=${r.result.status}`).join(',')}</div>;
}

beforeEach(() => { source = { read: jest.fn(async () => fresh()) }; });

test('loading first, then the status from the shared snapshot; many dots cost one request', async () => {
  render(<><One e={entity('x', 'platform-x')} /><One e={entity('y', 'platform-y')} /><Many list={[entity('x', 'platform-x'), entity('y', 'platform-y')]} /></>);
  expect(screen.getByTestId('x').textContent).toBe('loading');
  await act(async () => { await Promise.resolve(); await Promise.resolve(); });
  expect(screen.getByTestId('x').textContent).toBe('down');
  expect(screen.getByTestId('y').textContent).toBe('ok');
  expect(screen.getByTestId('many').textContent).toBe('x=down,y=ok');
  expect(source.read).toHaveBeenCalledTimes(1);
});

test('a source that never answers leaves the dot unknown', async () => {
  source = { read: jest.fn(async () => { throw new Error('down'); }) };
  render(<One e={entity('z', 'platform-z')} />);
  await act(async () => { await Promise.resolve(); await Promise.resolve(); });
  expect(screen.getByTestId('z').textContent).toBe('unknown');
});
