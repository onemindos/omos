import { ReactNode, useEffect, useState } from 'react';
import { identityApiRef, useApi } from '@backstage/core-plugin-api';
import { Box, Typography, useTheme } from '@material-ui/core';

export function greetingFor(hour: number) {
  if (hour < 12) return 'Good morning';
  if (hour < 18) return 'Good afternoon';
  return 'Good evening';
}

export function firstNameOf(displayName?: string) {
  const name = displayName?.trim();
  if (!name || name.includes('@')) return undefined;
  return name.split(/\s+/)[0];
}

export function formatClock(date: Date, timeZone: string) {
  return new Intl.DateTimeFormat('en-GB', {
    timeZone,
    hour: '2-digit',
    minute: '2-digit',
    hourCycle: 'h23',
  }).format(date);
}

export function HomeHeader({ action }: { action?: ReactNode }) {
  const theme = useTheme();
  const identity = useApi(identityApiRef);
  const [name, setName] = useState<string>();
  const [now, setNow] = useState(() => new Date());
  useEffect(() => {
    let active = true;
    identity
      .getProfileInfo()
      .then(profile => {
        if (active) setName(firstNameOf(profile.displayName));
      })
      .catch(() => undefined);
    return () => {
      active = false;
    };
  }, [identity]);
  useEffect(() => {
    const timer = setInterval(() => setNow(new Date()), 30000);
    return () => clearInterval(timer);
  }, []);
  const clock = (label: string, zone: string) => (
    <Typography
      variant="caption"
      color="textSecondary"
      component="span"
      style={{ fontVariantNumeric: 'tabular-nums' }}
    >
      {label} <time dateTime={now.toISOString()}>{formatClock(now, zone)}</time>
    </Typography>
  );
  return (
    <Box
      component="section"
      aria-label="Welcome"
      className="onemind-home-toolbar"
      display="flex"
      justifyContent="space-between"
      alignItems="center"
      flexWrap="wrap"
      style={{ gap: theme.spacing(1) }}
    >
      <Box>
        <Typography variant="h5" component="h1" style={{ fontWeight: 700 }}>
          {greetingFor(now.getHours())}
          {name ? `, ${name}` : ''}
        </Typography>
        <Box display="flex" flexWrap="wrap" style={{ gap: 12 }}>
          <Typography variant="caption" color="textSecondary">
            {now.toLocaleDateString(undefined, {
              weekday: 'long',
              month: 'long',
              day: 'numeric',
              year: 'numeric',
            })}
          </Typography>
          {clock('New York', 'America/New_York')}
          {clock('UTC', 'UTC')}
        </Box>
      </Box>
      {action}
    </Box>
  );
}
