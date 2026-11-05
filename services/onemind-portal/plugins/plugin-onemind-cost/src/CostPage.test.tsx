import { render, screen, waitFor, within } from '@testing-library/react';
import { discoveryApiRef, fetchApiRef } from '@backstage/core-plugin-api';
import { catalogApiRef } from '@backstage/plugin-catalog-react';
import { TestApiProvider } from '@backstage/test-utils';
import { CostPage, money } from './CostPage';
import fixture from './__fixtures__/allocation-7d.json';

const entity = (name: string, namespace: string, pillar: string) => ({
  metadata: {
    name,
    annotations: {
      'backstage.io/kubernetes-namespace': namespace,
      'onemind.io/pillar-review': 'approved',
    },
    labels: { 'onemind.io/pillar': pillar },
  },
});

const renderPage = (
  respond: () => Promise<Partial<Response>>,
  catalog: () => Promise<{ items: unknown[] }> = async () => ({
    items: [entity('tak', 'takserver', 'so')],
  }),
) =>
  render(
    <TestApiProvider
      apis={[
        [discoveryApiRef, { getBaseUrl: async () => 'http://x/api/proxy' }],
        [fetchApiRef, { fetch: jest.fn(respond) } as any],
        [catalogApiRef, { getEntities: jest.fn(catalog) } as any],
      ]}
    >
      <CostPage />
    </TestApiProvider>,
  );

const json =
  (body: unknown, ok = true, status = 200) =>
  async () => ({
    ok,
    status,
    json: async () => body,
  });

beforeEach(() => {
  jest.useFakeTimers({
    doNotFake: ['nextTick', 'setImmediate', 'queueMicrotask'],
  });
  jest.setSystemTime(Date.parse('2026-11-05T13:00:00Z'));
});
afterEach(() => jest.useRealTimers());

it('formats money without hiding tiny or unknown values', () => {
  expect(money(12.345)).toBe('$12.35');
  expect(money(0.001)).toBe('<$0.01');
  expect(money(0)).toBe('$0.00');
  expect(money(null)).toBe('Unknown');
  expect(money(undefined)).toBe('Unknown');
  expect(money(NaN)).toBe('Unknown');
});

it('always says these are estimates and not the DigitalOcean bill', async () => {
  renderPage(json(fixture));
  const note = await screen.findByRole('note', { name: 'About these numbers' });
  expect(note.textContent).toMatch(/estimates, not your bill/);
  expect(note.textContent).toMatch(/DigitalOcean invoice/);
});

it('shows real namespace cost, the app, the pillar and unattributed cost', async () => {
  renderPage(json(fixture));
  const table = await screen.findByRole('table', { name: 'Cost by namespace' });
  const tak = within(table).getByText('takserver').closest('tr')!;
  expect(
    within(tak).getByRole('link', { name: 'tak' }).getAttribute('href'),
  ).toBe('/catalog/default/component/tak');
  expect(tak.textContent).toMatch(/\$14\.70/);
  expect(tak.textContent).toMatch(/SO · Sovereign Operations/);
  const nats = within(table).getByText('nats').closest('tr')!;
  expect(nats.textContent).toMatch(/No catalog app declares this namespace/);
  const summary = screen.getByLabelText('Cost summary');
  expect(summary.textContent).toMatch(/__unmounted__/);
  expect(summary.textContent).toMatch(/\$26\.97/);
});

it('does not call a normal OpenCost payload incomplete just because it has an unmounted line', async () => {
  renderPage(json(fixture));
  await screen.findByRole('table', { name: 'Cost by namespace' });
  expect(screen.queryByRole('status')).toBeNull();
  expect(screen.getByLabelText('Cost summary').textContent).not.toMatch(
    /at least/,
  );
});

it('marks the total partial and the namespace Unknown when a cost is missing', async () => {
  const payload: any = JSON.parse(JSON.stringify(fixture));
  delete payload.data[0].takserver.totalCost;
  renderPage(json(payload));
  const table = await screen.findByRole('table', { name: 'Cost by namespace' });
  expect(screen.getByRole('status').textContent).toMatch(/incomplete/);
  const tak = within(table).getByText('takserver').closest('tr')!;
  expect(tak.textContent).toMatch(/Unknown/);
  expect(tak.textContent).toMatch(/incomplete data/);
  expect(tak.textContent).not.toMatch(/\$0\.00/);
  const summary = screen.getByLabelText('Cost summary');
  expect(summary.textContent).toMatch(/at least/);
  expect(summary.textContent).toMatch(/1 namespace with unknown cost/);
  const pillars = screen.getByRole('table', { name: 'Cost by pillar' });
  expect(pillars.textContent).toMatch(/1 unknown/);
});

it('shows unknown, not zero, when OpenCost cannot be reached', async () => {
  renderPage(async () => {
    throw new Error('network');
  });
  const alert = await screen.findByRole('alert');
  expect(alert.textContent).toMatch(/could not be reached/);
  expect(screen.getByLabelText('Cost summary').textContent).toMatch(/Unknown/);
  expect(screen.queryByRole('table')).toBeNull();
  expect(screen.queryByText('$0.00')).toBeNull();
});

it('treats a proxy error status as a source error', async () => {
  renderPage(json({ code: 403, error: 'no' }, false, 403));
  const alert = await screen.findByRole('alert');
  expect(alert.textContent).toMatch(/reported an error/);
  expect(screen.queryByRole('table')).toBeNull();
});

it('shows an empty result as empty, not as zero cost', async () => {
  renderPage(json({ code: 200, data: [] }));
  await waitFor(() =>
    expect(screen.getByRole('status').textContent).toMatch(
      /not the same as zero/,
    ),
  );
  expect(screen.queryByRole('table')).toBeNull();
});

it('still shows namespaces when the catalog is unavailable', async () => {
  renderPage(json(fixture), async () => {
    throw new Error('catalog down');
  });
  const table = await screen.findByRole('table', { name: 'Cost by namespace' });
  expect(within(table).getByText('takserver')).toBeTruthy();
  expect(
    screen
      .getAllByRole('status')
      .some(el => /catalog could not be read/.test(el.textContent ?? '')),
  ).toBe(true);
});
