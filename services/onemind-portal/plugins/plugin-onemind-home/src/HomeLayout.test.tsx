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
function mockStub(label: string) {
  return () => <section aria-label={label} />;
}
jest.mock('./HomeHeader', () => ({
  HomeHeader: ({ action }: any) => (
    <section aria-label="Welcome">{action}</section>
  ),
}));
jest.mock('./CommandBar', () => ({ CommandBar: mockStub('Command bar') }));
jest.mock('./NextUpHero', () => ({ NextUpHero: mockStub('Next up') }));
jest.mock('./PlatformScore', () => ({
  PlatformScore: mockStub('Platform score'),
}));
jest.mock('./LiveTile', () => ({ LiveTile: mockStub('Live') }));
jest.mock('./WaitingOnYou', () => ({
  WaitingOnYou: mockStub('Waiting on you'),
}));
jest.mock('./TasksCard', () => ({ TasksCard: mockStub('Tasks') }));
jest.mock('./ActionTiles', () => ({ ActionTiles: mockStub('Do it now') }));
jest.mock('./JumpBackIn', () => ({ JumpBackIn: mockStub('Jump back in') }));
jest.mock('./PinnedChips', () => ({ PinnedChips: mockStub('Pinned') }));
jest.mock('./BrowseAllApps', () => ({
  BrowseAllApps: ({ children }: any) => (
    <section aria-label="Browse all apps">{children}</section>
  ),
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

const top = [
  'Welcome',
  'Command bar',
  'Next up',
  'Platform score',
  'Live',
  'Waiting on you',
  'Tasks',
  'Do it now',
  'Jump back in',
  'Pinned',
];

test('home sections flow in the specified order with pillars hidden by default', () => {
  pillars = false;
  render(<HomeLayout {...props} />);
  expect(order()).toEqual([
    ...top,
    'Recently changed',
    'Browse all apps',
    'Launcher',
    'Active stack',
  ]);
});

test('pillar tiles appear before recent changes when enabled', () => {
  pillars = true;
  render(<HomeLayout {...props} />);
  expect(order()).toEqual([
    ...top,
    'Pillars',
    'Recently changed',
    'Browse all apps',
    'Launcher',
    'Active stack',
  ]);
});

test('no section sets a fixed height', () => {
  pillars = true;
  const { container } = render(<HomeLayout {...props} />);
  expect(container.innerHTML).not.toMatch(/(?:^|[^-])height:\s*\d/);
});
