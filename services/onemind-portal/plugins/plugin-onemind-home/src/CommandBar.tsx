import { useEffect, useId, useMemo, useRef, useState } from 'react';
import { useNavigate } from 'react-router-dom';
import { Box, Link, Typography, useTheme } from '@material-ui/core';
import { Link as RouterLink } from '@backstage/core-components';
import SearchIcon from '@material-ui/icons/Search';
import { AppEntry, useDirectory } from '@internal/plugin-onemind-common';

/** Instant directory matches: title hits rank above description or category hits. */
export function matchApps(apps: AppEntry[], query: string, limit = 6) {
  const q = query.trim().toLowerCase();
  if (!q) return [];
  const scored = apps.flatMap(app => {
    const title = app.title.toLowerCase();
    let rank = -1;
    if (title.startsWith(q)) rank = 0;
    else if (title.includes(q)) rank = 1;
    else if (`${app.description} ${app.groupTitle}`.toLowerCase().includes(q))
      rank = 2;
    return rank < 0 ? [] : [{ app, rank }];
  });
  return scored
    .sort((a, b) => a.rank - b.rank || a.app.title.localeCompare(b.app.title))
    .slice(0, limit)
    .map(item => item.app);
}

export const searchUrl = (query: string) =>
  `/search?query=${encodeURIComponent(query.trim())}`;

const isTypingTarget = (target: EventTarget | null) => {
  const node = target as HTMLElement | null;
  if (!node?.tagName) return false;
  return (
    ['INPUT', 'TEXTAREA', 'SELECT'].includes(node.tagName) ||
    node.isContentEditable
  );
};

export function CommandBar({
  navigate: navigateOverride,
}: {
  navigate?: (url: string) => void;
}) {
  const routerNavigate = useNavigate();
  const navigate = navigateOverride ?? routerNavigate;
  const theme = useTheme();
  const inputId = useId();
  const listId = useId();
  const { directory } = useDirectory();
  const [query, setQuery] = useState('');
  const inputRef = useRef<HTMLInputElement | null>(null);
  const resultsRef = useRef<HTMLUListElement | null>(null);
  const matches = useMemo(
    () => matchApps(directory, query),
    [directory, query],
  );
  const trimmed = query.trim();

  useEffect(() => {
    const onKey = (event: KeyboardEvent) => {
      if (
        event.key !== '/' ||
        event.ctrlKey ||
        event.metaKey ||
        event.altKey ||
        isTypingTarget(event.target)
      )
        return;
      event.preventDefault();
      inputRef.current?.focus();
    };
    document.addEventListener('keydown', onKey);
    return () => document.removeEventListener('keydown', onKey);
  }, []);

  const rowStyle = {
    display: 'flex',
    alignItems: 'baseline',
    gap: 8,
    padding: theme.spacing(0.75, 1),
    borderRadius: 8,
    color: theme.palette.text.primary,
    textDecoration: 'none',
  } as const;

  return (
    <Box
      component="section"
      aria-label="Command bar"
      style={{
        border: `1px solid ${theme.palette.divider}`,
        borderRadius: 14,
        background: theme.palette.background.paper,
        padding: theme.spacing(1, 1.5),
      }}
    >
      <Box display="flex" alignItems="center" style={{ gap: 8 }}>
        <SearchIcon color="action" fontSize="small" aria-hidden />
        <input
          ref={inputRef}
          id={inputId}
          type="search"
          value={query}
          autoComplete="off"
          aria-label="Search apps, runbooks, services or ask Legacy"
          aria-controls={trimmed ? listId : undefined}
          placeholder="Search apps, runbooks, services — or ask Legacy"
          onChange={event => setQuery(event.target.value)}
          onKeyDown={event => {
            if (event.key === 'Enter' && trimmed) {
              event.preventDefault();
              navigate(searchUrl(trimmed));
            } else if (event.key === 'Escape') {
              setQuery('');
            } else if (event.key === 'ArrowDown') {
              const first = resultsRef.current?.querySelector<HTMLElement>('a');
              if (first) {
                event.preventDefault();
                first.focus();
              }
            }
          }}
          style={{
            flex: 1,
            minWidth: 0,
            border: 'none',
            outline: 'none',
            background: 'transparent',
            color: theme.palette.text.primary,
            font: 'inherit',
            fontSize: 15,
            padding: theme.spacing(0.75, 0),
          }}
        />
        <Typography variant="caption" color="textSecondary" aria-hidden>
          Press / to focus
        </Typography>
      </Box>
      {trimmed && (
        <ul
          id={listId}
          ref={resultsRef}
          aria-label="Matches"
          style={{
            listStyle: 'none',
            margin: 0,
            marginTop: theme.spacing(0.75),
            borderTop: `1px solid ${theme.palette.divider}`,
            paddingTop: theme.spacing(0.75),
          }}
        >
          {matches.map(app => (
            <li key={app.ref}>
              <Link
                href={app.url}
                target="_blank"
                rel="noopener noreferrer"
                underline="none"
                style={rowStyle}
              >
                <Typography variant="body2" style={{ fontWeight: 600 }}>
                  {app.title} ↗
                </Typography>
                <Typography
                  variant="caption"
                  color="textSecondary"
                  noWrap
                  style={{ minWidth: 0 }}
                >
                  {app.groupTitle}
                  {app.description ? ` · ${app.description}` : ''}
                </Typography>
              </Link>
            </li>
          ))}
          {!matches.length && (
            <li>
              <Typography
                variant="caption"
                color="textSecondary"
                style={{ padding: theme.spacing(0.75, 1) }}
                display="block"
              >
                No apps match. Search everything or ask Legacy.
              </Typography>
            </li>
          )}
          <li>
            <RouterLink
              to={searchUrl(trimmed)}
              underline="none"
              style={rowStyle}
            >
              <Typography variant="body2" style={{ fontWeight: 600 }}>
                Search everything for &lsquo;{trimmed}&rsquo;
              </Typography>
              <Typography variant="caption" color="textSecondary">
                Enter
              </Typography>
            </RouterLink>
          </li>
          <li>
            <RouterLink to="/assistant" underline="none" style={rowStyle}>
              <Typography variant="body2" style={{ fontWeight: 600 }}>
                Ask Legacy
              </Typography>
              <Typography variant="caption" color="textSecondary">
                Opens the assistant; type your question there
              </Typography>
            </RouterLink>
          </li>
        </ul>
      )}
    </Box>
  );
}
