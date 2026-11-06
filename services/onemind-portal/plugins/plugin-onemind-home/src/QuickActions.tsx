import { ReactNode } from 'react';
import { Box, Typography, useTheme } from '@material-ui/core';
import BuildIcon from '@material-ui/icons/Build';
import AddCircleOutlineIcon from '@material-ui/icons/AddCircleOutline';
import AssignmentIcon from '@material-ui/icons/Assignment';
import ChatIcon from '@material-ui/icons/Chat';
import MenuBookIcon from '@material-ui/icons/MenuBook';
import OpenInNewIcon from '@material-ui/icons/OpenInNew';
import { AppLink, HomeCard } from './HomeCard';
import { useCloudWork } from './cloudWork';

export type QuickAction = {
  key: string;
  label: string;
  hint: string;
  href: string;
  external?: boolean;
  icon: ReactNode;
};

/** Every route here is registered by an existing plugin or template. */
export function quickActions(zohoProjectsUrl?: string): QuickAction[] {
  return [
    {
      key: 'change',
      label: 'Prepare a change',
      hint: 'Reviewed pull request',
      href: '/create',
      icon: <BuildIcon fontSize="small" />,
    },
    {
      key: 'add-app',
      label: 'Add an app',
      hint: 'Catalog + deploy draft',
      href: '/create/templates/default/add-app',
      icon: <AddCircleOutlineIcon fontSize="small" />,
    },
    {
      key: 'mission',
      label: 'Start a mission',
      hint: 'Plan and track work',
      href: '/missions',
      icon: <AssignmentIcon fontSize="small" />,
    },
    {
      key: 'legacy',
      label: 'Ask Legacy',
      hint: 'Answers with sources',
      href: '/assistant',
      icon: <ChatIcon fontSize="small" />,
    },
    {
      key: 'runbook',
      label: 'Find a runbook',
      hint: 'Docs and procedures',
      href: '/docs',
      icon: <MenuBookIcon fontSize="small" />,
    },
    {
      key: 'zoho',
      label: 'Open Zoho Projects',
      hint: 'Opens Zoho',
      href: zohoProjectsUrl ?? 'https://projects.zoho.com',
      external: true,
      icon: <OpenInNewIcon fontSize="small" />,
    },
  ];
}

export function QuickActions() {
  const theme = useTheme();
  const { data } = useCloudWork();
  const actions = quickActions(data?.link);
  return (
    <HomeCard label="Do it now" title="Do it now">
      <Box
        style={{
          display: 'grid',
          gridTemplateColumns: 'repeat(2, minmax(0, 1fr))',
          gap: 8,
        }}
      >
        {actions.map(action => (
          <AppLink
            key={action.key}
            href={action.href}
            underline="none"
            color="inherit"
            style={{
              display: 'flex',
              alignItems: 'center',
              gap: 8,
              padding: theme.spacing(0.75, 1),
              border: `1px solid ${theme.palette.divider}`,
              borderRadius: 10,
              minWidth: 0,
            }}
          >
            <Box color="primary.main" display="flex">
              {action.icon}
            </Box>
            <Box style={{ minWidth: 0 }}>
              <Typography
                variant="body2"
                style={{ fontWeight: 600, lineHeight: 1.25 }}
              >
                {action.label}
              </Typography>
              <Typography
                variant="caption"
                color="textSecondary"
                display="block"
                style={{ lineHeight: 1.2 }}
              >
                {action.hint}
              </Typography>
            </Box>
          </AppLink>
        ))}
      </Box>
    </HomeCard>
  );
}
