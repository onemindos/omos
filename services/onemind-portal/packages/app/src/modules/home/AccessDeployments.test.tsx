import { fireEvent, render, screen, waitFor } from '@testing-library/react';
import { AccessDeployments } from './AccessDeployments';
let saved: any = {};
const set = jest.fn();
const bucket = {
  snapshot: () => ({ presence: 'present', value: saved }),
  observe$: () => ({ subscribe: () => ({ unsubscribe: jest.fn() }) }),
  set,
};
const catalog = {
  queryEntities: jest.fn(async () => ({
    items: [
      {
        kind: 'Resource',
        metadata: {
          name: 'mac',
          title: 'My Mac',
          annotations: { 'onemind.io/mesh-dns': 'mac.example.ts.net' },
        },
        spec: { type: 'mesh-endpoint' },
      },
    ],
    pageInfo: {},
  })),
  getEntityByRef: jest.fn(),
};
jest.mock('@backstage/core-plugin-api', () => ({
  storageApiRef: { kind: 'storage' },
  useApi: (ref: any) =>
    ref.kind === 'catalog' ? catalog : { forBucket: () => bucket },
}));
jest.mock('@backstage/plugin-catalog-react', () => ({
  catalogApiRef: { kind: 'catalog' },
}));
const app = {
  name: 'mail',
  title: 'Mail',
  description: 'Mail',
  url: 'https://example.com',
  group: 'services',
  groupTitle: 'Communication',
  ref: 'component:default/mail',
};
beforeEach(() => {
  saved = {};
  set.mockReset();
  set.mockResolvedValue(undefined);
});
test('saves an additional native access point and preferred route while preserving another app', async () => {
  saved = { 'component:default/other': { preferred: 'web', points: [] } };
  render(<AccessDeployments app={app} open onClose={() => {}} />);
  await screen.findByText('No deployment relationship recorded.');
  fireEvent.click(screen.getByRole('button', { name: 'Add access point' }));
  fireEvent.change(screen.getByLabelText('Access name'), {
    target: { value: 'My Mac installation' },
  });
  fireEvent.change(screen.getByLabelText('Actual location'), {
    target: { value: 'My Mac' },
  });
  fireEvent.mouseDown(screen.getByLabelText('Platform'));
  fireEvent.click(await screen.findByRole('option', { name: 'macOS' }));
  fireEvent.mouseDown(screen.getByLabelText('Installation method'));
  fireEvent.click(await screen.findByRole('option', { name: 'Native app' }));
  fireEvent.mouseDown(
    screen.getByLabelText('Preferred access on this account'),
  );
  fireEvent.click(
    await screen.findByRole('option', { name: 'My Mac installation · macOS' }),
  );
  fireEvent.click(
    screen.getByRole('button', { name: 'Save access preferences' }),
  );
  await waitFor(() => expect(set).toHaveBeenCalled());
  const next = set.mock.calls[0][1];
  expect(next['component:default/other']).toEqual(
    saved['component:default/other'],
  );
  expect(next[app.ref].points[0]).toMatchObject({
    platform: 'macOS',
    method: 'Native app',
    evidence: 'User recorded',
  });
  expect(next[app.ref].preferred).toBe(next[app.ref].points[0].id);
});
test('does not lose stored access preferences when inventory lookup fails', async () => {
  catalog.queryEntities.mockRejectedValueOnce(new Error('unavailable'));
  saved = {
    [app.ref]: {
      preferred: 'custom',
      points: [
        {
          id: 'custom',
          label: 'Home server',
          platform: 'Linux',
          method: 'Docker',
          location: 'Home',
          evidence: 'User recorded',
        },
      ],
    },
  };
  render(<AccessDeployments app={app} open onClose={() => {}} />);
  expect(await screen.findByRole('alert')).toHaveProperty(
    'textContent',
    'Inventory lookup failed. Previous device results may be stale.',
  );
  expect(screen.getByLabelText('Access name')).toHaveProperty(
    'value',
    'Home server',
  );
});
