import { ReactNode } from 'react';
import { Box, Typography } from '@material-ui/core';
import SearchIcon from '@material-ui/icons/Search';
import { categoryPalette } from '@internal/plugin-onemind-theme/tokens';
import { AppLink } from './HomeCard';
import { useBentoStyles } from './bento';
import { useCloudWork } from './cloudWork';
import { useHomeTokens, withAlpha } from './homeTokens';
import { quickActions } from './QuickActions';

type Hue = keyof typeof categoryPalette;

/** One identity hue per destination, from the shared category palette. */
export const ACTION_HUES: Record<string, Hue> = {
  change: 'orange',
  'add-app': 'teal',
  mission: 'violet',
  legacy: 'blue',
  runbook: 'moss',
  zoho: 'plum',
  search: 'slate',
};

/** Same effect as pressing `/`: the command bar listens for it on the document. */
export function focusCommandBar() {
  document.dispatchEvent(
    new KeyboardEvent('keydown', { key: '/', bubbles: true }),
  );
}

/** Row C: the Do-it-now destinations as colored tiles, plus Search. */
export function ActionTiles() {
  const t = useHomeTokens();
  const classes = useBentoStyles();
  const { data } = useCloudWork();
  const tiles: {
    key: string;
    label: string;
    icon: ReactNode;
    href?: string;
  }[] = [
    ...quickActions(data?.link).map(a => ({
      key: a.key,
      label: a.label,
      icon: a.icon,
      href: a.href,
    })),
    {
      key: 'search',
      label: 'Search  /',
      icon: <SearchIcon fontSize="small" />,
    },
  ];
  return (
    <Box component="section" aria-label="Do it now">
      <ul
        className={classes.tiles}
        style={{ listStyle: 'none', margin: 0, padding: 0 }}
      >
        {tiles.map(tile => {
          const hue = categoryPalette[ACTION_HUES[tile.key] ?? 'slate'];
          const color = t.dark ? hue.bright : hue.base;
          const style = {
            display: 'flex',
            flexDirection: 'column' as const,
            alignItems: 'center',
            justifyContent: 'center',
            gap: 8,
            padding: '16px 8px',
            minHeight: 44,
            borderRadius: 14,
            border: `1px solid ${withAlpha(hue.base, t.dark ? '66' : '55')}`,
            background: withAlpha(hue.base, t.dark ? '26' : '14'),
            color: t.text[1],
            textAlign: 'center' as const,
            textDecoration: 'none',
            font: 'inherit',
            cursor: 'pointer',
            width: '100%',
          };
          const body = (
            <>
              <span style={{ color, display: 'flex', fontSize: 24 }}>
                {tile.icon}
              </span>
              <Typography
                variant="body2"
                style={{ fontWeight: 600, lineHeight: 1.2 }}
              >
                {tile.label}
              </Typography>
            </>
          );
          return (
            <li key={tile.key}>
              {tile.href ? (
                <AppLink
                  href={tile.href}
                  underline="none"
                  color="inherit"
                  style={style}
                >
                  {body}
                </AppLink>
              ) : (
                <button type="button" onClick={focusCommandBar} style={style}>
                  {body}
                </button>
              )}
            </li>
          );
        })}
      </ul>
    </Box>
  );
}
