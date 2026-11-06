import { CSSProperties, ReactNode, useEffect, useState } from 'react';
import { Box, makeStyles, Typography } from '@material-ui/core';
import { useHomeTokens } from './homeTokens';

/**
 * Bento rows. Widths come from grid tracks, never fixed heights; rows collapse
 * to one column on narrow viewports.
 */
export const useBentoStyles = makeStyles({
  rowA: {
    display: 'grid',
    gridTemplateColumns: 'minmax(0, 2fr) minmax(0, 1fr)',
    gap: 16,
    '@media (max-width: 720px)': { gridTemplateColumns: 'minmax(0, 1fr)' },
  },
  rowB: {
    display: 'grid',
    gridTemplateColumns: 'repeat(2, minmax(0, 1fr))',
    gap: 16,
    '@media (max-width: 720px)': { gridTemplateColumns: 'minmax(0, 1fr)' },
  },
  tiles: {
    display: 'grid',
    gridTemplateColumns: 'repeat(auto-fit, minmax(min(100%, 120px), 1fr))',
    gap: 12,
  },
});

/** Re-renders on an interval so relative times stay current. */
export function useNow(intervalMs = 30000) {
  const [now, setNow] = useState(() => new Date());
  useEffect(() => {
    const timer = setInterval(() => setNow(new Date()), intervalMs);
    return () => clearInterval(timer);
  }, [intervalMs]);
  return now;
}

/** Small muted label that sits above a card's content. */
export function Eyebrow({
  children,
  color,
}: {
  children: ReactNode;
  color?: string;
}) {
  const t = useHomeTokens();
  return (
    <Typography
      variant="caption"
      component="p"
      style={{
        margin: 0,
        fontSize: 12,
        fontWeight: 600,
        letterSpacing: '0.04em',
        textTransform: 'uppercase',
        color: color ?? t.text[3],
      }}
    >
      {children}
    </Typography>
  );
}

/** Token-styled surface for the bento cards. `hero` gets more padding and a raised tone. */
export function Surface({
  label,
  hero = false,
  children,
  style,
}: {
  label: string;
  hero?: boolean;
  children: ReactNode;
  style?: CSSProperties;
}) {
  const t = useHomeTokens();
  return (
    <Box
      component="section"
      aria-label={label}
      className="onemind-home-card"
      style={{
        minWidth: 0,
        display: 'flex',
        flexDirection: 'column',
        gap: 12,
        padding: hero ? 24 : 20,
        borderRadius: 16,
        border: `1px solid ${t.border.subtle}`,
        background: hero ? t.surface[2] : t.surface[1],
        ...style,
      }}
    >
      {children}
    </Box>
  );
}
