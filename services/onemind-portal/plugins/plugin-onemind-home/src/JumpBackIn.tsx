import { ComponentType, useEffect, useState } from 'react';
import { useApiHolder } from '@backstage/core-plugin-api';
import { Visit, visitsApiRef } from '@backstage/plugin-home';
import { Link } from '@backstage/core-components';
import { Box, Chip, Typography } from '@material-ui/core';
import { HomeCard } from './HomeCard';

/** Drops the home page itself and repeats, keeping first occurrence order. */
export function dedupeVisits(
  visits: Visit[],
  exclude: Set<string> = new Set(),
) {
  const seen = new Set(exclude);
  return visits.filter(visit => {
    if (!visit.pathname || visit.pathname === '/' || seen.has(visit.pathname))
      return false;
    seen.add(visit.pathname);
    return true;
  });
}

// MUI v4's Chip typing cannot express a router Link as `component`.
const LinkChip = Chip as unknown as ComponentType<Record<string, unknown>>;

function Chips({ label, visits }: { label: string; visits: Visit[] }) {
  if (!visits.length) return null;
  return (
    <Box mb={1}>
      <Typography variant="caption" color="textSecondary" display="block">
        {label}
      </Typography>
      <Box display="flex" flexWrap="wrap" style={{ gap: 6, marginTop: 4 }}>
        {visits.map(visit => (
          <LinkChip
            key={visit.id ?? visit.pathname}
            size="small"
            clickable
            component={Link}
            to={visit.pathname}
            label={visit.name}
            variant="outlined"
            style={{ maxWidth: '100%' }}
          />
        ))}
      </Box>
    </Box>
  );
}

export function JumpBackIn() {
  const visitsApi = useApiHolder().get(visitsApiRef);
  const [state, setState] = useState<
    | { kind: 'loading' }
    | { kind: 'unavailable' }
    | { kind: 'ready'; recent: Visit[]; top: Visit[] }
  >({ kind: 'loading' });
  useEffect(() => {
    if (!visitsApi) {
      setState({ kind: 'unavailable' });
      return undefined;
    }
    let active = true;
    Promise.all([
      visitsApi.list({
        limit: 12,
        orderBy: [{ field: 'timestamp', direction: 'desc' }],
      }),
      visitsApi.list({
        limit: 12,
        orderBy: [{ field: 'hits', direction: 'desc' }],
      }),
    ])
      .then(([recent, top]) => {
        if (!active) return;
        const recentShown = dedupeVisits(recent).slice(0, 6);
        const topShown = dedupeVisits(
          top,
          new Set(recentShown.map(visit => visit.pathname)),
        ).slice(0, 6);
        setState({ kind: 'ready', recent: recentShown, top: topShown });
      })
      .catch(() => active && setState({ kind: 'unavailable' }));
    return () => {
      active = false;
    };
  }, [visitsApi]);
  return (
    <HomeCard label="Jump back in" title="Jump back in">
      {state.kind === 'loading' && (
        <Typography variant="caption" color="textSecondary" role="status">
          Loading your visit history…
        </Typography>
      )}
      {state.kind === 'unavailable' && (
        <Typography variant="caption" color="error" role="alert">
          Visit history unavailable.
        </Typography>
      )}
      {state.kind === 'ready' && (
        <>
          <Chips label="Recent" visits={state.recent} />
          <Chips label="Most visited" visits={state.top} />
          {!state.recent.length && !state.top.length && (
            <Typography variant="caption" color="textSecondary">
              Pages you open will show up here.
            </Typography>
          )}
        </>
      )}
    </HomeCard>
  );
}
