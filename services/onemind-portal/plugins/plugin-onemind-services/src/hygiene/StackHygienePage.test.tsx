import { discoveryApiRef, fetchApiRef } from '@backstage/core-plugin-api';
import { renderInTestApp, TestApiProvider } from '@backstage/test-utils';
import { fireEvent, screen, waitFor } from '@testing-library/react';
import { Report, StackHygienePage } from './StackHygienePage';

const report: Report = {
  generatedAt: '2026-11-05T06:00:00.000Z',
  trigger: 'schedule',
  clusterObservedAt: '2026-11-05T05:59:00.000Z',
  checks: [
    { id: 'argo-without-catalog', title: 'Argo apps without a catalog entry', status: 'findings', count: 1 },
    { id: 'no-runbook', title: 'Apps with no runbook', status: 'ok', count: 0 },
    { id: 'old-snapshot', title: 'Old DigitalOcean snapshots', status: 'unavailable', count: 0, reason: 'DigitalOcean snapshot data is not available: F6 is not built.' },
  ],
  findings: [
    {
      check: 'argo-without-catalog',
      severity: 'warn',
      subject: 'platform-unmapped',
      detail: 'Argo application platform-unmapped is named by no catalog app.',
      actions: [{ kind: 'add', label: 'Add to catalog', href: '/create/templates/default/add-app?formData=%7B%7D' }],
    },
  ],
};

function mount(fetch: jest.Mock) {
  return renderInTestApp(
    <TestApiProvider
      apis={[
        [discoveryApiRef, { getBaseUrl: async () => 'http://portal/api/onemind-lifecycle' }],
        [fetchApiRef, { fetch }],
      ]}
    >
      <StackHygienePage />
    </TestApiProvider>,
  );
}

const json = (body: unknown, status = 200) => ({ ok: status < 400, status, json: async () => body });

describe('StackHygienePage', () => {
  it('shows findings with Add links, and unavailable checks with their reason', async () => {
    await mount(jest.fn(async () => json({ report })));
    expect(await screen.findByText('platform-unmapped')).toBeTruthy();
    expect(screen.getByRole('link', { name: 'Add to catalog' }).getAttribute('href')).toContain('/create/templates/default/add-app');
    expect(screen.getByText('Unavailable')).toBeTruthy();
    expect(screen.getByText(/Not evaluated: DigitalOcean snapshot data is not available/)).toBeTruthy();
    expect(screen.getByText('No findings')).toBeTruthy();
  });

  it('says so when no report exists yet', async () => {
    await mount(jest.fn(async () => json({ report: null })));
    expect(await screen.findByText(/No report yet/)).toBeTruthy();
  });

  it('tells a non-administrator that refresh is for administrators', async () => {
    const fetch = jest.fn(async (_url: string, init?: { method?: string }) =>
      init?.method === 'POST' ? json({}, 403) : json({ report }),
    );
    await mount(fetch);
    await screen.findByText('platform-unmapped');
    fireEvent.click(screen.getByRole('button', { name: 'Refresh report' }));
    await waitFor(() => expect(screen.getByRole('status').textContent).toMatch(/Only administrators/));
    expect(fetch).toHaveBeenCalledWith('http://portal/api/onemind-lifecycle/stack-hygiene/refresh', { method: 'POST' });
  });
});
