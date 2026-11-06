import { ReactNode } from 'react';
import { Link as RouterLink } from '@backstage/core-components';
import {
  Box,
  Link as MuiLink,
  LinkProps,
  Typography,
  useTheme,
} from '@material-ui/core';

/** Compact bordered card shared by the home command-center sections. */
export function HomeCard({
  label,
  title,
  action,
  children,
}: {
  label: string;
  title: string;
  action?: ReactNode;
  children: ReactNode;
}) {
  const theme = useTheme();
  return (
    <Box
      component="section"
      aria-label={label}
      className="onemind-home-card"
      style={{
        border: `1px solid ${theme.palette.divider}`,
        borderRadius: 14,
        background: theme.palette.background.paper,
        padding: theme.spacing(1.5, 2),
        minWidth: 0,
      }}
    >
      <Box
        display="flex"
        justifyContent="space-between"
        alignItems="center"
        style={{ gap: 8, marginBottom: theme.spacing(1) }}
      >
        <Typography
          variant="subtitle1"
          component="h2"
          style={{ fontWeight: 700 }}
        >
          {title}
        </Typography>
        {action}
      </Box>
      {children}
    </Box>
  );
}

export const cardGridStyle = (minWidth = 300) => ({
  display: 'grid',
  gridTemplateColumns: `repeat(auto-fit, minmax(min(100%, ${minWidth}px), 1fr))`,
  gap: 16,
  alignItems: 'start' as const,
});

/**
 * Internal paths navigate inside the app (Backstage Link); external URLs open
 * in a new tab with `noopener noreferrer`.
 */
export function AppLink({
  href,
  children,
  ...rest
}: Omit<LinkProps, 'href'> & { href: string }) {
  if (href.startsWith('/'))
    return (
      <RouterLink to={href} {...(rest as object)}>
        {children}
      </RouterLink>
    );
  return (
    <MuiLink href={href} target="_blank" rel="noopener noreferrer" {...rest}>
      {children}
    </MuiLink>
  );
}
