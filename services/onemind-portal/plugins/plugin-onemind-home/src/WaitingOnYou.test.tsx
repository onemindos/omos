import { MemoryRouter } from 'react-router-dom';
import { render, screen, waitFor } from '@testing-library/react';
import { bigValue, WaitingOnYou } from './WaitingOnYou';

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
  attention: [
    { key: 'a', severity: 'warning' },
    { key: 'b', severity: 'critical' },
  ],
  needsSetup: [{ key: 's' }],
  ...extra,
});

describe('bigValue', () => {
  test('failure and loading never read as 0', () => {
    expect(bigValue({ kind: 'unavailable' })).toEqual({
      text: '—',
      tone: 'warn',
      note: 'unavailable',
    });
    expect(bigValue({ kind: 'loading' }).text).toBe('…');
  });

  test('colour only signals: critical is danger, stale or partial is warning', () => {
    expect(bigValue({ kind: 'count', value: 0 }).tone).toBe('muted');
    expect(bigValue({ kind: 'count', value: 3 }).tone).toBe('plain');
    expect(bigValue({ kind: 'count', value: 3 }, { critical: true }).tone).toBe(
      'danger',
    );
    expect(
      bigValue({ kind: 'count', value: 3, stale: true }, { critical: true }),
    ).toMatchObject({ tone: 'warn', note: 'stale' });
    expect(
      bigValue({ kind: 'count', value: 3, incomplete: true }),
    ).toMatchObject({ text: '≥3', tone: 'warn', note: 'incomplete' });
  });
});

test('renders unavailable, not 0, when every source fails', async () => {
  unreadResult = () => Promise.reject(new Error('down'));
  statusResult = { data: undefined, error: true };
  render(
    <MemoryRouter>
      <WaitingOnYou />
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
      <WaitingOnYou />
    </MemoryRouter>,
  );
  expect(await screen.findByText('4')).not.toBeNull();
  expect(screen.getByText('2')).not.toBeNull();
  expect(screen.getByText('1')).not.toBeNull();
  expect(screen.getByText(/PRs waiting for your review/)).not.toBeNull();
});
