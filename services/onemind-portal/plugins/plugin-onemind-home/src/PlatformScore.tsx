import { Box, Typography } from '@material-ui/core';
import { AppLink } from './HomeCard';
import { Eyebrow, Surface } from './bento';
import { healthScore, HealthScore, scoreBand } from './health';
import { useHomeTokens } from './homeTokens';
import { summarizeStatus, useOperationsStatus } from './status';
import { TILES } from './StatusStrip';

const BAND_LABEL = {
  healthy: 'Healthy',
  degraded: 'Degraded',
  'at-risk': 'At risk',
} as const;

/** Text for the number: `—` when unknown, `≤` when sources are missing. */
export function scoreText(health: HealthScore | undefined) {
  if (!health) return '—';
  return `${health.partial ? '≤' : ''}${health.score}`;
}

function countText(value: number | 'unknown' | undefined) {
  if (value === undefined) return '—';
  return value === 'unknown' ? '?' : String(value);
}

type Tokens = ReturnType<typeof useHomeTokens>;

/** Number colour: warning whenever the view is stale or incomplete. */
function scoreColor(health: HealthScore | undefined, t: Tokens) {
  if (!health) return t.text[3];
  if (health.stale || health.partial) return t.status.warn.fg;
  const band = scoreBand(health.score);
  if (band === 'healthy') return t.status.ok.fg;
  return band === 'degraded' ? t.status.warn.fg : t.status.critical.fg;
}

/** Row A, narrow: a 0 to 100 score from the live operations status. */
export function PlatformScore() {
  const t = useHomeTokens();
  const { data, error } = useOperationsStatus();
  const health = healthScore(data, error);
  const counts = data ? summarizeStatus(data) : undefined;
  const band = health ? scoreBand(health.score) : undefined;
  const color = scoreColor(health, t);
  const flags = [
    health?.stale ? 'stale' : '',
    health?.partial ? 'incomplete' : '',
  ].filter(Boolean);
  let summary: string;
  if (!health)
    summary = error ? 'Live status unavailable' : 'Waiting for observation';
  else if (health.top) summary = health.top.title;
  else summary = 'All checks passing';
  return (
    <Surface label="Platform score">
      <Eyebrow>Platform score</Eyebrow>
      <Box
        display="flex"
        alignItems="baseline"
        style={{ gap: 8, flexWrap: 'wrap' }}
      >
        <Typography
          component="p"
          aria-label={
            health
              ? `Platform score ${health.score} out of 100${
                  health.partial ? ', at most, incomplete view' : ''
                }`
              : 'Platform score unavailable'
          }
          style={{
            margin: 0,
            fontSize: 56,
            lineHeight: 1,
            fontWeight: 700,
            fontVariantNumeric: 'tabular-nums',
            color,
          }}
        >
          {scoreText(health)}
        </Typography>
        <Typography variant="body2" style={{ color, fontWeight: 600 }}>
          {health && band
            ? `${BAND_LABEL[band]}${
                flags.length ? ` · ${flags.join(', ')}` : ''
              }`
            : 'unavailable'}
        </Typography>
      </Box>
      <AppLink href="/live" underline="none" color="inherit">
        <Typography
          variant="body2"
          style={{ color: t.text[1], overflowWrap: 'anywhere' }}
        >
          {summary}
          {health && health.attentionCount > 1
            ? ` (+${health.attentionCount - 1} more)`
            : ''}
        </Typography>
      </AppLink>
      <Box
        component="ul"
        aria-label="Score inputs"
        style={{
          listStyle: 'none',
          margin: 0,
          padding: 0,
          display: 'flex',
          flexWrap: 'wrap',
          gap: '4px 14px',
        }}
      >
        {TILES.map(tile => (
          <li key={tile.key}>
            <AppLink
              href={tile.href}
              underline="none"
              title={tile.hint}
              style={{ color: t.text[2], fontSize: 12 }}
            >
              {tile.label.split(' ')[0]}{' '}
              <span
                style={{
                  fontWeight: 700,
                  fontVariantNumeric: 'tabular-nums',
                  color: t.text[1],
                }}
              >
                {countText(counts?.[tile.key])}
              </span>
            </AppLink>
          </li>
        ))}
      </Box>
      <Typography variant="caption" style={{ color: t.text[3] }}>
        100 minus 20 per critical and 8 per warning item. No trend: the status
        APIs return a current snapshot only.
      </Typography>
    </Surface>
  );
}
