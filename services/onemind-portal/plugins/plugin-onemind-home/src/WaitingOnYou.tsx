import { Box, Typography } from '@material-ui/core';
import { AppLink } from './HomeCard';
import { Eyebrow, Surface } from './bento';
import { useHomeTokens } from './homeTokens';
import { attentionEntries, useOperationsStatus } from './status';
import {
  CountState,
  statusCounts,
  useUnreadCount,
  waitingRows,
  WaitingRow,
} from './WaitingCard';

export type BigTone = 'plain' | 'muted' | 'warn' | 'danger';
export type BigValue = { text: string; tone: BigTone; note?: string };

/**
 * What a big number shows. A failed or unknown source is never 0: it renders
 * `—` with a note. Colour is only used to signal: danger for a critical item,
 * warning for a stale or incomplete count.
 */
export function bigValue(
  state: CountState,
  options: { critical?: boolean } = {},
): BigValue {
  if (state.kind === 'loading')
    return { text: '…', tone: 'muted', note: 'loading' };
  if (state.kind === 'unavailable')
    return { text: '—', tone: 'warn', note: 'unavailable' };
  const flagged = state.stale || state.incomplete;
  let note: string | undefined;
  if (state.stale) note = 'stale';
  else if (state.incomplete) note = 'incomplete';
  let tone: BigTone = state.value > 0 ? 'plain' : 'muted';
  if (flagged) tone = 'warn';
  else if (options.critical && state.value > 0) tone = 'danger';
  return {
    text: `${state.incomplete ? '≥' : ''}${state.value}`,
    tone,
    note,
  };
}

const COLUMNS = ['notifications', 'attention', 'setup'];

/** Row B: the same three counts as the Waiting card, as big numbers. */
export function WaitingOnYou() {
  const t = useHomeTokens();
  const unread = useUnreadCount();
  const { data, error } = useOperationsStatus();
  const rows = waitingRows(unread, statusCounts(data, error));
  const reviews = rows.find(row => row.key === 'reviews') as WaitingRow;
  const critical =
    !!data && attentionEntries(data).some(e => e.severity === 'critical');
  const toneColor: Record<BigTone, string> = {
    plain: t.text[1],
    muted: t.text[3],
    warn: t.status.warn.fg,
    danger: t.status.critical.fg,
  };
  return (
    <Surface label="Waiting on you">
      <Eyebrow>Waiting on you</Eyebrow>
      <Box
        component="ul"
        style={{
          listStyle: 'none',
          margin: 0,
          padding: 0,
          display: 'grid',
          gridTemplateColumns: 'repeat(3, minmax(0, 1fr))',
          gap: 12,
        }}
      >
        {COLUMNS.map(key => {
          const row = rows.find(r => r.key === key) as WaitingRow;
          const big = bigValue(row.state as CountState, {
            critical: key === 'attention' && critical,
          });
          return (
            <li key={key}>
              <AppLink
                href={row.href}
                underline="none"
                color="inherit"
                title={row.hint}
                aria-label={`${row.label}: ${big.text}${
                  big.note ? `, ${big.note}` : ''
                }`}
                style={{ display: 'block' }}
              >
                <Typography
                  component="span"
                  display="block"
                  style={{
                    fontSize: 36,
                    lineHeight: 1.1,
                    fontWeight: 700,
                    fontVariantNumeric: 'tabular-nums',
                    color: toneColor[big.tone],
                  }}
                >
                  {big.text}
                </Typography>
                <Typography
                  component="span"
                  display="block"
                  variant="caption"
                  style={{ color: t.text[2] }}
                >
                  {row.label}
                </Typography>
                {big.note && (
                  <Typography
                    component="span"
                    display="block"
                    variant="caption"
                    style={{ color: toneColor[big.tone] }}
                  >
                    {big.note}
                  </Typography>
                )}
              </AppLink>
            </li>
          );
        })}
      </Box>
      <AppLink href={reviews.href} variant="caption">
        PRs waiting for your review ↗
      </AppLink>
    </Surface>
  );
}
