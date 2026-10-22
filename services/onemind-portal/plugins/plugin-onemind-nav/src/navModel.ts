import type { NavContentNavItem } from '@backstage/plugin-app-react';
import type { NavFeature, NavFlags } from './features';

/**
 * Pure description of the sidebar. Rendering lives in Sidebar.tsx so this file
 * can be unit-tested without a Backstage app.
 */

export type IconKey =
  | 'home'
  | 'liveOps'
  | 'pillar'
  | 'catalog'
  | 'api'
  | 'graph'
  | 'cloud'
  | 'cost'
  | 'delivery'
  | 'insights'
  | 'radar'
  | 'docs'
  | 'runbook'
  | 'announce'
  | 'adr'
  | 'create'
  | 'register'
  | 'retire'
  | 'audit'
  | 'visualizer'
  | 'devtools'
  | 'more'
  | 'search';

export type EntrySpec = {
  key: string;
  title: string;
  icon: IconKey;
  /** Extension ids of a plugin page to reuse (first one present wins). */
  ids?: string[];
  /** Fallback match on the page's title when its id is not known. */
  titles?: string[];
  /** Always-available destination, also used when no plugin page is found. */
  href?: string;
  flag?: NavFeature;
};

export type GroupSpec = {
  key: string;
  label: string;
  icon: IconKey;
  entries: EntrySpec[];
  flag?: NavFeature;
  adminOnly?: boolean;
  collapsed?: boolean;
};

// F2 Domains: domain:default/onemind-so|le|hp|ge (parent onemind).
const pillar = (slug: string, title: string): EntrySpec => ({
  key: `pillar-${slug}`,
  title,
  icon: 'pillar',
  href: `/catalog/default/domain/onemind-${slug}`,
});

export const NAV_GROUPS: GroupSpec[] = [
  {
    key: 'now',
    label: 'Now',
    icon: 'home',
    entries: [
      {
        key: 'home',
        title: 'Home',
        icon: 'home',
        ids: ['page:home'],
        href: '/',
      },
      {
        key: 'live-ops',
        title: 'Live Ops',
        icon: 'liveOps',
        ids: ['page:onemind-live'],
        href: '/live',
      },
      // Today (inbox, calendar, work) keeps its own entry now that Live Ops has the page.
      { key: 'today', title: 'Today', icon: 'insights', href: '/today' },
      // Notifications render as NotificationsSidebarItem (live unread badge).
    ],
  },
  {
    key: 'pillars',
    label: 'Pillars',
    icon: 'pillar',
    flag: 'pillars',
    entries: [
      pillar('so', 'SO · Sovereign Operations'),
      pillar('le', 'LE · Legacy Evolution'),
      pillar('hp', 'HP · Holistic Performance'),
      pillar('ge', 'GE · Generational Empire'),
    ],
  },
  {
    key: 'services',
    label: 'Services',
    icon: 'catalog',
    entries: [
      {
        key: 'catalog',
        title: 'Catalog',
        icon: 'catalog',
        ids: ['page:catalog'],
        href: '/catalog',
      },
      {
        key: 'apis',
        title: 'APIs',
        icon: 'api',
        ids: ['page:api-docs'],
        titles: ['APIs'],
      },
      {
        key: 'graph',
        title: 'Catalog Graph',
        icon: 'graph',
        ids: ['page:catalog-graph'],
        titles: ['Catalog Graph'],
      },
      {
        key: 'digitalocean',
        title: 'DigitalOcean',
        icon: 'cloud',
        href: '/digitalocean',
        flag: 'digitalocean',
      },
    ],
  },
  {
    key: 'insights',
    label: 'Insights',
    icon: 'insights',
    entries: [
      {
        key: 'cost',
        title: 'Cost',
        icon: 'cost',
        ids: ['page:cost-insights'],
        titles: ['Cost Insights'],
      },
      {
        key: 'delivery',
        title: 'Delivery',
        icon: 'delivery',
        ids: ['page:delivery-insights'],
        titles: ['Delivery Insights'],
      },
      {
        key: 'tech-insights',
        title: 'Tech Insights',
        icon: 'insights',
        ids: ['page:tech-insights'],
        titles: ['Tech Insights'],
      },
      {
        key: 'tech-radar',
        title: 'Tech Radar',
        icon: 'radar',
        ids: ['page:tech-radar'],
        titles: ['Tech Radar'],
      },
    ],
  },
  {
    key: 'knowledge',
    label: 'Knowledge',
    icon: 'docs',
    entries: [
      {
        key: 'docs',
        title: 'Docs',
        icon: 'docs',
        ids: ['page:techdocs'],
        href: '/docs',
      },
      {
        key: 'runbooks',
        title: 'Runbooks',
        icon: 'runbook',
        href: '/docs/default/component/onemind-runbooks',
      },
      {
        key: 'announcements',
        title: 'Announcements',
        icon: 'announce',
        ids: ['page:announcements'],
        titles: ['Announcements'],
      },
      {
        key: 'adrs',
        title: 'ADRs',
        icon: 'adr',
        ids: ['page:adr'],
        titles: ['ADRs', 'ADR'],
      },
    ],
  },
  {
    key: 'change',
    label: 'Change',
    icon: 'create',
    entries: [
      {
        key: 'create',
        title: 'Create',
        icon: 'create',
        ids: ['page:scaffolder'],
        href: '/create',
      },
      {
        key: 'register',
        title: 'Register existing',
        icon: 'register',
        ids: ['page:catalog-import'],
        titles: ['Register Existing Component'],
      },
      // F8 ships the retire template; confirm this route there.
      {
        key: 'retire',
        title: 'Retire app',
        icon: 'retire',
        href: '/create/templates/default/retire-app',
        flag: 'retireApp',
      },
    ],
  },
  {
    key: 'admin',
    label: 'Admin',
    icon: 'audit',
    adminOnly: true,
    collapsed: true,
    entries: [
      {
        key: 'audit',
        title: 'Audit Logs',
        icon: 'audit',
        ids: ['page:audit-logs'],
        titles: ['Audit Logs'],
      },
      {
        key: 'visualizer',
        title: 'Visualizer',
        icon: 'visualizer',
        ids: ['page:app-visualizer'],
        titles: ['Visualizer'],
      },
      // F9 adds DevTools; hidden until its page exists.
      {
        key: 'devtools',
        title: 'DevTools',
        icon: 'devtools',
        ids: ['page:devtools'],
        titles: ['DevTools'],
      },
    ],
  },
];

/**
 * Pages rendered elsewhere in the sidebar (search modal, notifications item,
 * Settings group) or replaced by a group entry. They never fall into "More".
 */
export const HANDLED_IDS = [
  'page:search',
  'page:notifications',
  'page:onemind-today',
  'page:user-settings',
];

export type ResolvedEntry = {
  key: string;
  title: string;
  icon: IconKey;
  href: string;
};
export type ResolvedGroup = {
  key: string;
  label: string;
  icon: IconKey;
  collapsed: boolean;
  entries: ResolvedEntry[];
};

export const itemId = (item: NavContentNavItem) => item.node.spec.id;

export type BuildOptions = {
  items: NavContentNavItem[];
  flags: NavFlags;
  isAdmin: boolean;
};

/**
 * Resolves the spec against the pages the app actually has. Every plugin page
 * is placed exactly once: in its group, hidden by design (admin or flag), or in
 * the trailing "More" group, so no existing page becomes unreachable.
 */
export function buildNavModel({
  items,
  flags,
  isAdmin,
}: BuildOptions): ResolvedGroup[] {
  const claimed = new Set<string>(HANDLED_IDS);
  const byId = new Map(items.map(item => [itemId(item), item] as const));
  const find = (entry: EntrySpec) => {
    for (const id of entry.ids ?? []) {
      const hit = byId.get(id);
      if (hit && !claimed.has(id)) return hit;
    }
    for (const title of entry.titles ?? []) {
      const hit = items.find(
        item => item.title === title && !claimed.has(itemId(item)),
      );
      if (hit) return hit;
    }
    return undefined;
  };

  const groups: ResolvedGroup[] = [];
  for (const group of NAV_GROUPS) {
    const entries: ResolvedEntry[] = [];
    for (const entry of group.entries) {
      const page = find(entry);
      // A matched plugin page is claimed even when its entry is hidden, so it
      // does not reappear under "More".
      if (page) claimed.add(itemId(page));
      const visible =
        (!entry.flag || flags[entry.flag]) && (!group.adminOnly || isAdmin);
      const href = page?.href ?? entry.href;
      if (visible && href) {
        entries.push({
          key: entry.key,
          title: entry.title,
          icon: entry.icon,
          href,
        });
      }
    }
    if ((group.flag && !flags[group.flag]) || !entries.length) continue;
    groups.push({
      key: group.key,
      label: group.label,
      icon: group.icon,
      collapsed: !!group.collapsed,
      entries,
    });
  }

  const more = items
    .filter(item => !claimed.has(itemId(item)))
    .sort((a, b) => a.title.localeCompare(b.title))
    .map<ResolvedEntry>(item => ({
      key: itemId(item),
      title: item.title,
      icon: 'more',
      href: item.href,
    }));
  if (more.length) {
    groups.push({
      key: 'more',
      label: 'More',
      icon: 'more',
      collapsed: true,
      entries: more,
    });
  }
  return groups;
}
