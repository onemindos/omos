import { ReactElement, useEffect, useState } from 'react';
import {
  Sidebar,
  SidebarDivider,
  SidebarGroup,
  SidebarItem,
  SidebarSpace,
  SidebarSubmenu,
  SidebarSubmenuItem,
  useSidebarOpenState,
  useSidebarPinState,
} from '@backstage/core-components';
import { identityApiRef, useApi } from '@backstage/frontend-plugin-api';
import {
  NavContentBlueprint,
  NavContentNavItems,
} from '@backstage/plugin-app-react';
import { SearchModalProvider } from '@backstage/plugin-search';
import { UserSettingsSignInAvatar } from '@backstage/plugin-user-settings';
import { NotificationsSidebarItem } from '@backstage/plugin-notifications';
import { Typography, useTheme } from '@material-ui/core';
import { makeStyles } from '@material-ui/core/styles';
import { SidebarLogo } from './SidebarLogo';
import { CommandPalette, PaletteTrigger } from './CommandPalette';
import { useNavFlags } from './features';
import {
  buildNavModel,
  IconKey,
  ResolvedEntry,
  ResolvedGroup,
} from './navModel';
import { ICONS } from './icons';

// Eight bottom-bar groups would overflow a 390px phone at Backstage's default
// 80px minimum width, so let them share the row evenly.
const useBarStyles = makeStyles({
  '@global': {
    '.MuiBottomNavigation-root .MuiBottomNavigationAction-root': {
      flex: '1 1 0',
      minWidth: 0,
      margin: 0,
      padding: '6px 0',
    },
  },
});

export const ADMIN_GROUP = 'group:default/onemind-admins';

/** UI convenience only: hides admin links. Backend permissions are unchanged. */
export function useIsAdmin() {
  const identity = useApi(identityApiRef);
  const [admin, setAdmin] = useState(false);
  useEffect(() => {
    let active = true;
    identity
      .getBackstageIdentity()
      .then(
        value =>
          active && setAdmin(value.ownershipEntityRefs.includes(ADMIN_GROUP)),
      )
      .catch(() => active && setAdmin(false));
    return () => {
      active = false;
    };
  }, [identity]);
  return admin;
}

const iconFor = (key: IconKey) => ICONS[key];

function GroupHeading({ label }: { label: string }) {
  const theme = useTheme();
  const { isOpen } = useSidebarOpenState();
  const { isMobile } = useSidebarPinState();
  if (isMobile) return null;
  // Collapsed rail: no room for text, so a divider keeps the groups separate.
  if (!isOpen) return <SidebarDivider />;
  return (
    <Typography
      variant="overline"
      component="h2"
      style={{
        display: 'block',
        padding: theme.spacing(1.5, 3, 0.5),
        lineHeight: 1.6,
        opacity: 0.75,
        color: (theme.palette as any).navigation?.color,
      }}
    >
      {label}
    </Typography>
  );
}

function FlatItems({ entries }: { entries: ResolvedEntry[] }) {
  return (
    <>
      {entries.map(entry => (
        <SidebarItem
          key={entry.key}
          icon={iconFor(entry.icon)}
          to={entry.href}
          text={entry.title}
        />
      ))}
    </>
  );
}

function GroupBody({
  group,
  extra,
}: {
  group: ResolvedGroup;
  extra?: ReactElement;
}) {
  const { isMobile } = useSidebarPinState();
  // Collapsed groups sit behind one button that opens a keyboard-operable flyout.
  if (group.collapsed && !isMobile) {
    return (
      <SidebarItem icon={iconFor(group.icon)} text={group.label}>
        <SidebarSubmenu title={group.label}>
          {group.entries.map(entry => (
            <SidebarSubmenuItem
              key={entry.key}
              title={entry.title}
              to={entry.href}
              icon={iconFor(entry.icon)}
            />
          ))}
        </SidebarSubmenu>
      </SidebarItem>
    );
  }
  return (
    <>
      <GroupHeading label={group.label} />
      <FlatItems entries={group.entries} />
      {extra}
    </>
  );
}

// Returns a SidebarGroup element directly (not a wrapper component): the mobile
// bottom bar finds groups by element type among the Sidebar's direct children.
function renderGroup(
  group: ResolvedGroup,
  priority: number,
  extra?: ReactElement,
) {
  const Icon = iconFor(group.icon);
  return (
    <SidebarGroup
      key={group.key}
      label={group.label}
      icon={<Icon />}
      priority={priority}
    >
      <GroupBody group={group} extra={extra} />
    </SidebarGroup>
  );
}

export function OneMindSidebar({ navItems }: { navItems: NavContentNavItems }) {
  useBarStyles();
  const flags = useNavFlags();
  const isAdmin = useIsAdmin();
  const groups = buildNavModel({ items: [...navItems.rest()], flags, isAdmin });
  const SearchGlyph = iconFor('search');
  return (
    <SearchModalProvider>
      <Sidebar>
        <SidebarLogo />
        <SidebarGroup label="Search" icon={<SearchGlyph />} priority={100}>
          <PaletteTrigger />
        </SidebarGroup>
        <SidebarDivider />
        {groups.map((group, index) =>
          renderGroup(
            group,
            90 - index,
            // Notifications keep their live unread badge instead of a plain link.
            group.key === 'now' ? <NotificationsSidebarItem /> : undefined,
          ),
        )}
        <SidebarSpace />
        <SidebarDivider />
        <SidebarGroup
          label="Settings"
          icon={<UserSettingsSignInAvatar />}
          to="/settings"
          priority={1}
        >
          <SidebarItem
            icon={() => <UserSettingsSignInAvatar />}
            to="/settings"
            text="Settings"
          />
        </SidebarGroup>
      </Sidebar>
      <CommandPalette />
    </SearchModalProvider>
  );
}

export const SidebarContent = NavContentBlueprint.make({
  params: {
    component: ({ navItems }) => <OneMindSidebar navItems={navItems} />,
  },
});
