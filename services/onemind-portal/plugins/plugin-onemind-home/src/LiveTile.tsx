import { useMemo } from 'react';
import { Box, Typography } from '@material-ui/core';
import { summarizeFeed, useLiveFeed } from '@internal/plugin-onemind-ops';
import type { Connection, LiveSummary } from '@internal/plugin-onemind-ops';
import { AppLink } from './HomeCard';
import { Eyebrow, Surface } from './bento';
import { useHomeTokens } from './homeTokens';

// Light by design: only the ops plugin's hook entry points are imported. The
// Ops Center page, and CesiumJS with it, stay behind that plugin's lazy route.

const CONNECTION_LABEL: Record<Connection, string> = {
  connected: 'Connected',
  reconnecting: 'Reconnecting',
  offline: 'Offline',
};

function ago(fromMs: number | undefined, nowMs: number) {
  if (fromMs === undefined) return 'no update yet';
  const s = Math.max(0, Math.round((nowMs - fromMs) / 1000));
  if (s < 5) return 'just now';
  if (s < 60) return `${s} s ago`;
  if (s < 3600) return `${Math.floor(s / 60)} min ago`;
  return `${Math.floor(s / 3600)} h ago`;
}

function Stat({ value, label }: { value: string; label: string }) {
  const t = useHomeTokens();
  return (
    <Box>
      <Typography
        component="p"
        style={{
          margin: 0,
          fontSize: 28,
          fontWeight: 700,
          lineHeight: 1.1,
          fontVariantNumeric: 'tabular-nums',
          color: t.text[1],
        }}
      >
        {value}
      </Typography>
      <Typography variant="caption" style={{ color: t.text[3] }}>
        {label}
      </Typography>
    </Box>
  );
}

/** Presentational part, kept separate so it is easy to test. */
export function LiveTileView({
  summary,
  now,
}: {
  summary: LiveSummary;
  now: number;
}) {
  const t = useHomeTokens();
  const tone = {
    connected: t.status.ok.fg,
    reconnecting: t.status.warn.fg,
    offline: t.status.critical.fg,
  }[summary.connection];
  const empty =
    summary.tracks === 0 &&
    summary.latestEvents.length === 0 &&
    summary.podsTotal === 0;
  const attention = summary.namespacesNeedingAttention;
  return (
    <>
      <Box
        display="flex"
        alignItems="baseline"
        style={{ gap: 8, flexWrap: 'wrap' }}
      >
        <Typography variant="body2" style={{ color: tone, fontWeight: 700 }}>
          {CONNECTION_LABEL[summary.connection]}
        </Typography>
        <Typography variant="caption" style={{ color: t.text[3] }}>
          Updated {ago(summary.lastUpdateMs, now)}
        </Typography>
      </Box>
      {empty ? (
        <Typography variant="body2" style={{ color: t.text[2] }}>
          Nothing on the bridge yet. Tracks, events and cluster health will show
          here as they arrive.
        </Typography>
      ) : (
        <>
          <Box display="flex" style={{ gap: 28, flexWrap: 'wrap' }}>
            <Stat value={String(summary.tracks)} label="tracks" />
            <Stat
              value={String(summary.recentEvents)}
              label="events, last 5 min"
            />
            <Stat
              value={`${summary.podsReady}/${summary.podsTotal}`}
              label="pods ready"
            />
          </Box>
          <Typography
            variant="body2"
            style={{
              color:
                attention || summary.warnings ? t.status.warn.fg : t.text[2],
            }}
          >
            {attention
              ? `${attention} namespace${
                  attention === 1 ? ' needs' : 's need'
                } attention`
              : 'All namespaces healthy'}
            {summary.warnings
              ? ` · ${summary.warnings} warning${
                  summary.warnings === 1 ? '' : 's'
                }`
              : ''}
          </Typography>
          {summary.latestEvents.length > 0 && (
            <Box
              component="ul"
              aria-label="Latest events"
              style={{ listStyle: 'none', margin: 0, padding: 0 }}
            >
              {summary.latestEvents.map(e => (
                <li key={e.id}>
                  <Typography
                    variant="caption"
                    style={{ color: t.text[2], overflowWrap: 'anywhere' }}
                  >
                    {e.summary || e.subject}
                  </Typography>
                </li>
              ))}
            </Box>
          )}
        </>
      )}
    </>
  );
}

/** Home tile: the live bridge at a glance, linking to the Ops Center. */
export function LiveTile() {
  const t = useHomeTokens();
  const feed = useLiveFeed();
  const summary = useMemo(
    () => summarizeFeed(feed.state, feed.now),
    [feed.state, feed.now],
  );
  let body;
  if (feed.load.kind === 'loading')
    body = (
      <Typography variant="body2" style={{ color: t.text[3] }}>
        Loading live picture…
      </Typography>
    );
  else if (feed.load.kind === 'ready')
    body = <LiveTileView summary={summary} now={feed.now} />;
  else
    body = (
      <Typography variant="body2" style={{ color: t.text[2] }}>
        Live bridge unavailable
      </Typography>
    );
  return (
    <Surface label="Live">
      <Eyebrow>Live</Eyebrow>
      {body}
      <AppLink href="/ops" underline="none" color="primary">
        <Typography variant="body2" style={{ fontWeight: 600 }}>
          Open Ops Center
        </Typography>
      </AppLink>
    </Surface>
  );
}
