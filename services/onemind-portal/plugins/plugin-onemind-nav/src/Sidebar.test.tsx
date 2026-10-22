import { screen } from '@testing-library/react';
import { renderInTestApp, TestApiProvider } from '@backstage/test-utils';
import { configApiRef, identityApiRef } from '@backstage/frontend-plugin-api';
import type { NavContentNavItem } from '@backstage/plugin-app-react';
import { OneMindSidebar } from './Sidebar';

jest.mock('@backstage/plugin-notifications', () => ({
  NotificationsSidebarItem: () => <a href="/notifications">Notifications</a>,
}));
jest.mock('@backstage/plugin-user-settings', () => ({
  UserSettingsSignInAvatar: () => <span>avatar</span>,
}));
jest.mock('@backstage/plugin-search', () => ({
  SearchModalProvider: ({ children }: { children: unknown }) => <>{children}</>,
  useSearchModal: () => ({
    state: { open: false, hidden: true },
    toggleModal: jest.fn(),
  }),
  SearchModal: () => null,
}));
jest.mock('./SidebarLogo', () => ({ SidebarLogo: () => null }));

const item = (id: string, title: string, href: string) =>
  ({
    node: { spec: { id } },
    title,
    href,
    icon: null,
    routeRef: null,
  } as unknown as NavContentNavItem);
const navItems = (list: NavContentNavItem[]) =>
  ({
    rest: () => list,
    take: () => undefined,
    clone: () => navItems(list),
    withComponent: () => ({}),
  } as any);

const pages = [
  item('page:home', 'Home', '/'),
  item('page:catalog', 'Catalog', '/catalog'),
  item('page:app-visualizer', 'Visualizer', '/visualizer'),
  item('page:onemind-systems', 'Systems', '/systems'),
];

const render = (flags: Record<string, boolean>, admin: boolean) =>
  renderInTestApp(
    <TestApiProvider
      apis={[
        [
          configApiRef,
          { getOptionalBoolean: (key: string) => flags[key.split('.').pop()!] },
        ],
        [
          identityApiRef,
          {
            getBackstageIdentity: async () => ({
              type: 'user' as const,
              userEntityRef: 'user:default/zeus',
              ownershipEntityRefs: admin
                ? ['group:default/onemind-admins']
                : [],
            }),
          },
        ],
      ]}
    >
      <OneMindSidebar navItems={navItems(pages)} />
    </TestApiProvider>,
  );

describe('OneMindSidebar', () => {
  it('renders the groups in order with accessible headings and live-ops link', async () => {
    await render({}, false);
    const headings = screen
      .getAllByRole('heading', { level: 2 })
      .map(h => h.textContent);
    expect(headings).toEqual(['Now', 'Services', 'Knowledge', 'Change']);
    // Collapsed groups are one button each; unplaced pages stay reachable in More.
    expect(screen.getByRole('button', { name: /More/ })).toBeTruthy();
    expect(
      screen.getByRole('link', { name: /Live Ops/ }).getAttribute('href'),
    ).toBe('/today');
    expect(screen.getByRole('link', { name: 'Notifications' })).toBeTruthy();
    expect(
      screen.getByRole('link', { name: /Catalog/ }).getAttribute('href'),
    ).toBe('/catalog');
    expect(screen.queryByText('DigitalOcean')).toBeNull();
    expect(screen.queryByText('Pillars')).toBeNull();
    expect(screen.queryByRole('link', { name: /Visualizer/ })).toBeNull();
  });

  it('shows flagged groups when their flags are on', async () => {
    await render({ pillars: true, digitalocean: true }, false);
    expect(screen.getByRole('heading', { name: 'Pillars' })).toBeTruthy();
    expect(
      screen.getByRole('link', { name: /SO · Sovereign Operations/ }),
    ).toBeTruthy();
    expect(
      screen.getByRole('link', { name: /DigitalOcean/ }).getAttribute('href'),
    ).toBe('/digitalocean');
  });

  it('puts admin pages behind a button for admins only', async () => {
    await render({}, true);
    expect(await screen.findByRole('button', { name: /Admin/ })).toBeTruthy();
  });
});
