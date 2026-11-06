import { MemoryRouter } from 'react-router-dom';
import { fireEvent, render, screen } from '@testing-library/react';
import { CommandBar, matchApps, searchUrl } from './CommandBar';

const apps = [
  {
    ref: 'component:default/grafana',
    name: 'grafana',
    title: 'Grafana',
    description: 'Dashboards',
    groupTitle: 'Observability',
    url: 'https://g.example',
  },
  {
    ref: 'component:default/argocd',
    name: 'argocd',
    title: 'Argo CD',
    description: 'GitOps deploys',
    groupTitle: 'Platform',
    url: 'https://a.example',
  },
  {
    ref: 'component:default/loki',
    name: 'loki',
    title: 'Loki',
    description: 'Logs for Grafana',
    groupTitle: 'Observability',
    url: 'https://l.example',
  },
] as any[];
jest.mock('@internal/plugin-onemind-common', () => ({
  useDirectory: () => ({ directory: apps }),
}));

test('ranks title matches before description matches and caps results', () => {
  expect(matchApps(apps, 'graf').map(a => a.name)).toEqual(['grafana', 'loki']);
  expect(matchApps(apps, 'gitops').map(a => a.name)).toEqual(['argocd']);
  expect(matchApps(apps, '  ')).toEqual([]);
  expect(matchApps(apps, 'o', 1)).toHaveLength(1);
});

test('shows instant matches with safe external links and the two actions', () => {
  render(
    <MemoryRouter>
      <CommandBar navigate={jest.fn()} />
    </MemoryRouter>,
  );
  fireEvent.change(screen.getByRole('searchbox'), {
    target: { value: 'graf' },
  });
  const link = screen.getByRole('link', { name: /Grafana ↗/ });
  expect(link.getAttribute('href')).toBe('https://g.example');
  expect(link.getAttribute('rel')).toBe('noopener noreferrer');
  expect(
    screen
      .getByRole('link', { name: /Search everything for/ })
      .getAttribute('href'),
  ).toBe('/search?query=graf');
  expect(
    screen.getByRole('link', { name: /Ask Legacy/ }).getAttribute('href'),
  ).toBe('/assistant');
});

test('Enter searches everything with an encoded query', () => {
  const navigate = jest.fn();
  render(
    <MemoryRouter>
      <CommandBar navigate={navigate} />
    </MemoryRouter>,
  );
  const input = screen.getByRole('searchbox');
  fireEvent.keyDown(input, { key: 'Enter' });
  expect(navigate).not.toHaveBeenCalled();
  fireEvent.change(input, { target: { value: 'a b&c' } });
  fireEvent.keyDown(input, { key: 'Enter' });
  expect(navigate).toHaveBeenCalledWith(searchUrl('a b&c'));
  expect(searchUrl('a b&c')).toBe('/search?query=a%20b%26c');
});

test('pressing / focuses the bar unless already typing in a field', () => {
  render(
    <MemoryRouter>
      <input aria-label="other" />
      <CommandBar />
    </MemoryRouter>,
  );
  fireEvent.keyDown(document.body, { key: '/' });
  expect(document.activeElement).toBe(screen.getByRole('searchbox'));
  const other = screen.getByLabelText('other');
  other.focus();
  fireEvent.keyDown(other, { key: '/' });
  expect(document.activeElement).toBe(other);
});
