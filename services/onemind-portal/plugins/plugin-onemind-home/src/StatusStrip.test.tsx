import { render, screen } from '@testing-library/react';
import { StatusStrip } from './StatusStrip';

let state: { data?: unknown; error: boolean } = { error: false };
jest.mock('./status', () => ({
  ...jest.requireActual('./status'),
  useOperationsStatus: () => state,
}));

const data = {
  observedAt: new Date().toISOString(),
  partial: false,
  totalWorkloads: 10,
  unavailable: ['n/one'],
  deployments: [],
  backups: [{ name: 'x/b', lastSuccess: null, suspended: true }],
  serviceHealth: {
    observedAt: new Date().toISOString(),
    checks: [{ name: 'a', state: 'down' }],
  },
};

test('shows four linked counts with observed times', () => {
  state = { data, error: false };
  render(<StatusStrip />);
  const alerts = screen.getByRole('link', { name: /Alerts firing: 1/ });
  expect(alerts.getAttribute('href')).toBe('/today');
  expect(
    screen.getByRole('link', { name: /Argo apps to review: 0/ }),
  ).toBeTruthy();
  expect(
    screen.getByRole('link', { name: /Workloads not ready: 1/ }),
  ).toBeTruthy();
  expect(
    screen.getByRole('link', { name: /Backups to review: 1/ }),
  ).toBeTruthy();
  expect(screen.getAllByText(/observed \d/)).toHaveLength(4);
});

test('marks stale data and a failed refresh', () => {
  state = {
    data: { ...data, observedAt: new Date(Date.now() - 600000).toISOString() },
    error: true,
  };
  render(<StatusStrip />);
  expect(screen.getAllByText(/Stale · observed/).length).toBeGreaterThan(0);
  expect(screen.getByRole('alert').textContent).toMatch(
    /Live refresh unavailable/,
  );
});

test('shows unknown rather than zero when a source failed', () => {
  state = {
    data: { ...data, partial: true, unavailableSources: ['Argo CD'] },
    error: false,
  };
  render(<StatusStrip />);
  expect(
    screen.getByRole('link', { name: /Argo apps to review: unknown/ }),
  ).toBeTruthy();
  expect(screen.getByRole('status').textContent).toMatch(/Argo CD unavailable/);
});

test('does not claim zero while loading', () => {
  state = { error: false };
  render(<StatusStrip />);
  expect(screen.getAllByText('Waiting for observation')).toHaveLength(4);
});
