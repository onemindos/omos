import { useEffect, useState } from 'react';
import { Box, Link, Typography, useTheme } from '@material-ui/core';
import { useApi } from '@backstage/core-plugin-api';
import { notificationsApiRef } from '@backstage/plugin-notifications';

export type ChangeItem = {
  id: string;
  title: string;
  topic?: string;
  link?: string | null;
  created: string;
};

// Topics sent by the activity backend and catalog producers that describe a change.
const CHANGE_TOPICS = new Set(['deployments', 'catalog']);

/** Newest first, one row per title (the title carries the resource), capped. */
export function recentChanges(items: ChangeItem[], limit = 10) {
  const seen = new Set<string>();
  return [...items]
    .filter(item => !item.topic || CHANGE_TOPICS.has(item.topic))
    .sort((a, b) => Date.parse(b.created) - Date.parse(a.created))
    .filter(item => {
      const key = item.title.trim().toLowerCase();
      if (seen.has(key)) return false;
      seen.add(key);
      return true;
    })
    .slice(0, limit);
}

export function RecentlyChanged() {
  const api = useApi(notificationsApiRef);
  const theme = useTheme();
  const [items, setItems] = useState<ChangeItem[]>();
  const [error, setError] = useState(false);
  useEffect(() => {
    let active = true;
    const load = async () => {
      try {
        // Reads the same notification feed as the inbox, read and unread alike.
        const result = await api.getNotifications({
          limit: 60,
          sort: 'created',
          sortOrder: 'desc',
        });
        if (!active) return;
        setItems(
          result.notifications.map(item => ({
            id: item.id,
            title: item.payload.title,
            topic: item.payload.topic,
            link: item.payload.link,
            created: String(item.created),
          })),
        );
        setError(false);
      } catch {
        if (active) setError(true);
      }
    };
    load();
    const timer = setInterval(load, 60000);
    return () => {
      active = false;
      clearInterval(timer);
    };
  }, [api]);
  const rows = items ? recentChanges(items) : [];
  return (
    <Box component="section" aria-label="Recently changed">
      <Typography variant="overline" color="textSecondary">
        Recently changed
      </Typography>
      {error && (
        <Typography variant="body2" color="error" role="alert">
          Change feed unavailable; showing nothing rather than guessing.
        </Typography>
      )}
      {!items && !error && <Typography variant="body2">Loading…</Typography>}
      {items && !rows.length && !error && (
        <Typography variant="body2" color="textSecondary">
          No deployments or catalog changes reported yet.
        </Typography>
      )}
      <Box component="ul" m={0} p={0} style={{ listStyle: 'none' }}>
        {rows.map(row => (
          <Box
            component="li"
            key={row.id}
            display="flex"
            justifyContent="space-between"
            style={{
              gap: theme.spacing(2),
              padding: theme.spacing(0.75, 0),
              borderBottom: `1px solid ${theme.palette.divider}`,
            }}
          >
            {row.link ? (
              <Link href={row.link} variant="body2">
                {row.title}
              </Link>
            ) : (
              <Typography variant="body2">{row.title}</Typography>
            )}
            <Typography variant="caption" color="textSecondary" noWrap>
              {new Date(row.created).toLocaleString([], {
                month: 'short',
                day: 'numeric',
                hour: '2-digit',
                minute: '2-digit',
              })}
            </Typography>
          </Box>
        ))}
      </Box>
    </Box>
  );
}
