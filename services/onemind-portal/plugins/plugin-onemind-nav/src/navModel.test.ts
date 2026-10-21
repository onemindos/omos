import type { NavContentNavItem } from '@backstage/plugin-app-react';
import { buildNavModel, NAV_GROUPS } from './navModel';

const item = (id: string, title: string, href: string) =>
  ({
    node: { spec: { id } },
    title,
    href,
    icon: null,
    routeRef: null,
  } as unknown as NavContentNavItem);

const pages = [
  item('page:home', 'Home', '/'),
  item('page:catalog', 'Catalog', '/catalog'),
  item('page:api-docs', 'APIs', '/api-docs'),
  item('page:catalog-graph', 'Catalog Graph', '/catalog-graph'),
  item('page:cost-insights', 'Cost Insights', '/cost-insights'),
  item('page:tech-insights', 'Tech Insights', '/tech-insights'),
  item('page:tech-radar', 'Tech Radar', '/tech-radar'),
  item('page:techdocs', 'Docs', '/docs'),
  item('page:announcements', 'Announcements', '/announcements'),
  item('page:scaffolder', 'Create...', '/create'),
  item('page:catalog-import', 'Register Existing Component', '/catalog-import'),
  item('page:app-visualizer', 'Visualizer', '/visualizer'),
  item('x:audit', 'Audit Logs', '/audit-logs'),
  item('page:search', 'Search', '/search'),
  item('page:notifications', 'Notifications', '/notifications'),
  item('page:onemind-today', 'Today', '/today'),
  item('page:user-settings', 'Settings', '/settings'),
  item('page:onemind-systems', 'Systems', '/systems'),
  item('page:onemind-assistant', 'Assistant', '/assistant'),
];
const off = { pillars: false, digitalocean: false, retireApp: false };
const labels = (groups: ReturnType<typeof buildNavModel>) =>
  groups.map(g => g.label);
const titles = (groups: ReturnType<typeof buildNavModel>, key: string) =>
  groups.find(g => g.key === key)?.entries.map(e => e.title);

describe('buildNavModel', () => {
  it('orders the groups as the work order specifies', () => {
    expect(NAV_GROUPS.map(g => g.label)).toEqual([
      'Now',
      'Pillars',
      'Services',
      'Insights',
      'Knowledge',
      'Change',
      'Admin',
    ]);
    const groups = buildNavModel({
      items: pages,
      flags: { ...off, pillars: true },
      isAdmin: true,
    });
    expect(labels(groups)).toEqual([
      'Now',
      'Pillars',
      'Services',
      'Insights',
      'Knowledge',
      'Change',
      'Admin',
      'More',
    ]);
  });

  it('puts Live Ops on /today and keeps Today reachable', () => {
    const groups = buildNavModel({ items: pages, flags: off, isAdmin: false });
    expect(groups[0].entries.map(e => [e.title, e.href])).toEqual([
      ['Home', '/'],
      ['Live Ops', '/today'],
    ]);
  });

  it('hides flagged entries by default and shows them when enabled', () => {
    const hidden = buildNavModel({ items: pages, flags: off, isAdmin: true });
    expect(labels(hidden)).not.toContain('Pillars');
    expect(titles(hidden, 'services')).not.toContain('DigitalOcean');
    expect(titles(hidden, 'change')).not.toContain('Retire app');

    const shown = buildNavModel({
      items: pages,
      flags: { pillars: true, digitalocean: true, retireApp: true },
      isAdmin: true,
    });
    expect(titles(shown, 'pillars')).toEqual([
      'SO · Sovereign Operations',
      'LE · Legacy Evolution',
      'HP · Holistic Performance',
      'GE · Generational Empire',
    ]);
    expect(
      shown.find(g => g.key === 'pillars')?.entries.map(e => e.href),
    ).toEqual([
      '/catalog/default/domain/onemind-so',
      '/catalog/default/domain/onemind-le',
      '/catalog/default/domain/onemind-hp',
      '/catalog/default/domain/onemind-ge',
    ]);
    expect(titles(shown, 'services')).toContain('DigitalOcean');
    expect(titles(shown, 'change')).toContain('Retire app');
  });

  it('shows Admin only to admins and does not leak its pages into More', () => {
    const user = buildNavModel({ items: pages, flags: off, isAdmin: false });
    expect(labels(user)).not.toContain('Admin');
    expect(titles(user, 'more')).not.toContain('Visualizer');
    expect(titles(user, 'more')).not.toContain('Audit Logs');
    const admin = buildNavModel({ items: pages, flags: off, isAdmin: true });
    expect(titles(admin, 'admin')).toEqual(['Audit Logs', 'Visualizer']);
    expect(admin.find(g => g.key === 'admin')?.collapsed).toBe(true);
  });

  it('matches by title when the extension id is unknown', () => {
    const groups = buildNavModel({ items: pages, flags: off, isAdmin: true });
    expect(titles(groups, 'admin')).toContain('Audit Logs');
    expect(titles(groups, 'change')).toContain('Register existing');
  });

  it('keeps every other page reachable under More, once', () => {
    const groups = buildNavModel({ items: pages, flags: off, isAdmin: true });
    expect(titles(groups, 'more')).toEqual(['Assistant', 'Systems']);
    const placed = groups.flatMap(g => g.entries.map(e => e.href));
    for (const page of pages.filter(
      p => !['/search', '/notifications', '/settings'].includes(p.href),
    )) {
      expect(placed.filter(href => href === page.href)).toHaveLength(1);
    }
  });

  it('omits an entry whose page is absent and has no fallback route', () => {
    const groups = buildNavModel({ items: [], flags: off, isAdmin: true });
    expect(titles(groups, 'knowledge')).toEqual(['Docs', 'Runbooks']);
    expect(titles(groups, 'admin')).toBeUndefined();
    expect(titles(groups, 'insights')).toBeUndefined();
  });
});
