import { MemoryRouter } from 'react-router-dom';
import { render, screen, waitFor } from '@testing-library/react';
import { statusCounts, waitingRows, WaitingCard } from './WaitingCard';

let unreadResult: any;
let statusResult: any;
jest.mock('@backstage/core-plugin-api', () => ({
  ...jest.requireActual('@backstage/core-plugin-api'),
  useApi: () => ({ getNotifications: () => unreadResult() }),
}));
jest.mock('@backstage/plugin-notifications', () => ({
  notificationsApiRef: {},
}));
jest.mock('./status', () => ({
  ...jest.requireActual('./status'),
  useOperationsStatus: () => statusResult,
}));

const data = (extra: any = {}) => ({
  observedAt: new Date().toISOString(),
  partial: false,
  totalWorkloads: 3,
  unavailable: [],
  deployments: [],
  backups: [],
  attention: [{ key: 'a' }, { key: 'b' }],
  needsSetup: [{ key: 's' }],
  ...extra,
});

test('statusCounts never turns a failed source into zero', () => {
  expect(statusCounts(undefined, true).attention).toEqual({
    kind: 'unavailable',
  });
  expect(statusCounts(undefined, false).attention).toEqual({ kind: 'loading' });
  expect(statusCounts(data() as any, false).attention).toMatchObject({
    kind: 'count',
    value: 2,
  });
  expect(statusCounts(data() as any, true).setup).toMatchObject({
    value: 1,
    stale: true,
  });
  expect(
    statusCounts(data({ partial: true }) as any, false).attention,
  ).toMatchObject({ incomplete: true });
});

test('rows link to the right places', () => {
  const rows = waitingRows({ kind: 'loading' }, statusCounts(undefined, false));
  expect(rows.map(r => r.href)).toEqual([
    '/notifications',
    '/live',
    '/live',
    'https://github.com/pulls/review-requested',
  ]);
});

test('renders unavailable, not 0, when sources fail', async () => {
  unreadResult = () => Promise.reject(new Error('down'));
  statusResult = { data: undefined, error: true };
  render(
    <MemoryRouter>
      <WaitingCard />
    </MemoryRouter>,
  );
  await waitFor(() =>
    expect(screen.getAllByText('unavailable')).toHaveLength(3),
  );
  expect(screen.queryByText('0')).toBeNull();
});

test('renders counts when sources succeed', async () => {
  unreadResult = () => Promise.resolve({ totalCount: 4, notifications: [] });
  statusResult = { data: data(), error: false };
  render(
    <MemoryRouter>
      <WaitingCard />
    </MemoryRouter>,
  );
  expect(await screen.findByText('4')).not.toBeNull();
  expect(screen.getByText('2')).not.toBeNull();
  expect(screen.getByText('1')).not.toBeNull();
});
