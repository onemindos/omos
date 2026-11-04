import { render, screen } from '@testing-library/react';
import { HomeLayout } from './HomeLayout';

let pillars = false;
jest.mock('./features', () => ({ useHomeFeature: () => pillars }));
jest.mock('@backstage/core-components', () => ({
  Content: ({ children }: any) => <div>{children}</div>,
}));
jest.mock('@backstage/plugin-home', () => ({ CustomHomepageGrid: () => null }));
jest.mock('@backstage/plugin-home-react/alpha', () => ({}));
jest.mock('./homeModule', () => ({
  Directory: () => <section aria-label="Launcher" />,
}));
jest.mock('./Launcher', () => ({
  QuickAccess: () => <section aria-label="Favorites" />,
}));
jest.mock('./StatusStrip', () => ({
  StatusStrip: () => <section aria-label="Live status" />,
}));
jest.mock('./PillarTiles', () => ({
  PillarTiles: () => <section aria-label="Pillars" />,
}));
jest.mock('./RecentlyChanged', () => ({
  RecentlyChanged: () => <section aria-label="Recently changed" />,
}));
jest.mock('./ActiveStack', () => ({
  ActiveStack: () => <section aria-label="Active stack" />,
}));

const props = { widgets: [], defaultConfig: [] } as any;
const order = () =>
  screen.getAllByRole('region').map(node => node.getAttribute('aria-label'));

test('home sections flow in the specified order with pillars hidden by default', () => {
  pillars = false;
  render(<HomeLayout {...props} />);
  expect(order()).toEqual([
    'Live status',
    'Favorites',
    'Recently changed',
    'Launcher',
    'Active stack',
  ]);
});

test('pillar tiles appear between favorites and recent changes when enabled', () => {
  pillars = true;
  render(<HomeLayout {...props} />);
  expect(order()).toEqual([
    'Live status',
    'Favorites',
    'Pillars',
    'Recently changed',
    'Launcher',
    'Active stack',
  ]);
});

test('no section sets a fixed height', () => {
  pillars = true;
  const { container } = render(<HomeLayout {...props} />);
  expect(container.innerHTML).not.toMatch(/(?:^|[^-])height:\s*\d/);
});
