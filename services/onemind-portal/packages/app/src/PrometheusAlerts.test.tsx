import { render, screen, waitFor } from '@testing-library/react';
import PrometheusAlerts from './PrometheusAlerts';

const getAlerts = jest.fn();
const api = { getAlerts };
const entity = {
  kind: 'Component',
  metadata: {
    name: 'k8s-nats-nats-0dee54ac',
    annotations: { 'prometheus.io/service-name': 'nats' },
  },
};
jest.mock('@backstage/core-plugin-api', () => ({
  ...jest.requireActual('@backstage/core-plugin-api'),
  useApi: () => api,
}));
jest.mock('@backstage/plugin-catalog-react', () => ({
  useEntity: () => ({ entity }),
}));
jest.mock('@roadiehq/backstage-plugin-prometheus', () => ({
  PrometheusAlertStatus: () => <div>Native Prometheus table</div>,
}));
const value = (coverage = 'observed', alerts: unknown[] = []) => ({
  scope: 'nats',
  observedAt: new Date().toISOString(),
  coverage,
  candidateRules: 12,
  unhealthyRules: coverage === 'partial' ? 1 : 0,
  staleRules: 0,
  data: { groups: [{ rules: [{ alerts }] }] },
});
beforeEach(() => getAlerts.mockReset());
test('no-match is not an all-clear; real links and timestamp stay visible', async () => {
  getAlerts.mockResolvedValue(value());
  render(<PrometheusAlerts />);
  await waitFor(() =>
    expect(screen.getByText(/No firing or pending/)).toBeInTheDocument(),
  );
  expect(
    screen.getByText(/not proof that the service is healthy/),
  ).toBeInTheDocument();
  expect(screen.getByText(/Observed .*refresh every 30/)).toBeInTheDocument();
  expect(
    screen.getByRole('link', { name: 'Recovery runbook' }),
  ).toHaveAttribute('href', '/docs/default/component/onemind-runbooks/nats/');
  expect(screen.getByText('Native Prometheus table')).toBeInTheDocument();
});
test('unavailable does not keep a table or emit healthy', async () => {
  getAlerts.mockRejectedValue(new Error('private upstream error'));
  render(<PrometheusAlerts />);
  expect(await screen.findByRole('alert')).toHaveTextContent('Unknown');
  expect(screen.queryByText('Native Prometheus table')).not.toBeInTheDocument();
  expect(screen.queryByText(/private upstream/)).not.toBeInTheDocument();
});
test('partial coverage is explicit while active alerts remain visible', async () => {
  getAlerts.mockResolvedValue(value('partial', [{ state: 'firing' }]));
  render(<PrometheusAlerts />);
  expect(await screen.findByRole('alert')).toHaveTextContent(
    'Coverage partial',
  );
  expect(screen.getByText('Native Prometheus table')).toBeInTheDocument();
  expect(screen.queryByText(/No firing or pending/)).not.toBeInTheDocument();
});
