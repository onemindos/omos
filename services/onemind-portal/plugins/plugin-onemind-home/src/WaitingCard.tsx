import { useEffect, useState } from 'react';
import { useApi } from '@backstage/core-plugin-api';
import { notificationsApiRef } from '@backstage/plugin-notifications';
import { Box, Typography, useTheme } from '@material-ui/core';
import { AppLink, HomeCard } from './HomeCard';
import {
  attentionEntries,
  backupNeedsSetup,
  isStale,
  OperationsStatus,
  useOperationsStatus,
} from './status';

export type CountState =
  | { kind: 'loading' }
  | { kind: 'unavailable' }
  | { kind: 'count'; value: number; stale?: boolean; incomplete?: boolean };

export type WaitingRow = {
  key: string;
  label: string;
  href: string;
  external?: boolean;
  hint: string;
  state: CountState | { kind: 'link' };
};

/** Maps the operations status to counts without ever turning an error into 0. */
export function statusCounts(
  data: OperationsStatus | undefined,
  error: boolean,
  now = Date.now(),
): { attention: CountState; setup: CountState } {
  if (!data) {
    const state: CountState = error
      ? { kind: 'unavailable' }
      : { kind: 'loading' };
    return { attention: state, setup: state };
  }
  const stale = error || isStale(data, now);
  const incomplete = !!data.partial;
  const setup =
    data.needsSetup?.length ?? data.backups.filter(backupNeedsSetup).length;
  return {
    attention: {
      kind: 'count',
      value: attentionEntries(data, now).length,
      stale,
      incomplete,
    },
    setup: { kind: 'count', value: setup, stale, incomplete },
  };
}

export function waitingRows(
  unread: CountState,
  status: { attention: CountState; setup: CountState },
): WaitingRow[] {
  return [
    {
      key: 'notifications',
      label: 'Unread notifications',
      href: '/notifications',
      hint: 'Portal inbox',
      state: unread,
    },
    {
      key: 'attention',
      label: 'Attention items',
      href: '/live',
      hint: 'Alerts, workloads, backups',
      state: status.attention,
    },
    {
      key: 'setup',
      label: 'Needs setup',
      href: '/live',
      hint: 'Waiting on a person',
      state: status.setup,
    },
    {
      key: 'reviews',
      label: 'PRs waiting for your review',
      href: 'https://github.com/pulls/review-requested',
      external: true,
      hint: 'Opens GitHub',
      state: { kind: 'link' },
    },
  ];
}

export function useUnreadCount(): CountState {
  const api = useApi(notificationsApiRef);
  const [state, setState] = useState<CountState>({ kind: 'loading' });
  useEffect(() => {
    let active = true;
    const load = async () => {
      try {
        const result = await api.getNotifications({ read: false, limit: 1 });
        if (active) setState({ kind: 'count', value: result.totalCount ?? 0 });
      } catch {
        if (active)
          setState(prev =>
            prev.kind === 'count'
              ? { ...prev, stale: true }
              : { kind: 'unavailable' },
          );
      }
    };
    load();
    const timer = setInterval(load, 30000);
    return () => {
      active = false;
      clearInterval(timer);
    };
  }, [api]);
  return state;
}

function Value({ state }: { state: WaitingRow['state'] }) {
  const theme = useTheme();
  if (state.kind === 'link')
    return (
      <Typography variant="body2" color="textSecondary" aria-hidden>
        ↗
      </Typography>
    );
  if (state.kind === 'loading')
    return (
      <Typography variant="body2" color="textSecondary" role="status">
        …
      </Typography>
    );
  if (state.kind === 'unavailable')
    return (
      <Typography variant="body2" color="error" title="Source unavailable">
        unavailable
      </Typography>
    );
  const flagged = state.stale || state.incomplete;
  let title: string | undefined;
  if (state.stale) title = 'Stale observation';
  else if (state.incomplete)
    title = 'Incomplete view: some sources unavailable';
  let color = theme.palette.text.secondary;
  if (flagged) color = theme.palette.warning.main;
  else if (state.value > 0) color = theme.palette.text.primary;
  return (
    <Typography
      variant="body2"
      title={title}
      style={{ fontWeight: 700, fontVariantNumeric: 'tabular-nums', color }}
    >
      {state.incomplete ? '≥' : ''}
      {state.value}
      {state.stale ? ' (stale)' : ''}
    </Typography>
  );
}

export function WaitingCard() {
  const theme = useTheme();
  const unread = useUnreadCount();
  const { data, error } = useOperationsStatus();
  const rows = waitingRows(unread, statusCounts(data, error));
  return (
    <HomeCard label="Waiting on you" title="Waiting on you">
      <Box component="ul" style={{ listStyle: 'none', margin: 0, padding: 0 }}>
        {rows.map(row => (
          <Box component="li" key={row.key}>
            <AppLink
              href={row.href}
              underline="none"
              color="inherit"
              style={{
                display: 'flex',
                justifyContent: 'space-between',
                alignItems: 'center',
                gap: 8,
                padding: theme.spacing(0.75, 0),
              }}
            >
              <Box>
                <Typography variant="body2">{row.label}</Typography>
                <Typography variant="caption" color="textSecondary">
                  {row.hint}
                </Typography>
              </Box>
              <Value state={row.state} />
            </AppLink>
          </Box>
        ))}
      </Box>
    </HomeCard>
  );
}
