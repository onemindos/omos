import { ReactNode } from 'react';
import { Box, Link, Typography, useTheme } from '@material-ui/core';
import { StatusChip } from '@internal/plugin-onemind-theme';
import { age, clock } from './format';
import type { PanelState } from './types';
import type { Owner } from './owners';

export function Panel(props: {
  title: string;
  state?: PanelState;
  hint?: string;
  wide?: boolean;
  children: ReactNode;
}) {
  const theme = useTheme();
  const { state } = props;
  const unavailable = state?.status === 'error';
  return (
    <Box
      component="section"
      aria-label={props.title}
      style={{
        border: `1px solid ${theme.palette.divider}`,
        borderRadius: theme.shape.borderRadius * 2,
        background: theme.palette.background.paper,
        minWidth: 0,
        gridColumn: props.wide ? '1 / -1' : undefined,
      }}
    >
      <Box
        px={2}
        py={1.5}
        display="flex"
        flexWrap="wrap"
        alignItems="center"
        style={{ gap: 8, borderBottom: `1px solid ${theme.palette.divider}` }}
      >
        <Typography variant="h6" component="h2" style={{ flex: '1 1 auto', fontWeight: 600 }}>
          {props.title}
        </Typography>
        {state?.status === 'stale' && <StatusChip status="warn" label="Stale" />}
        {unavailable && <StatusChip status="critical" label="Unavailable" />}
        {state?.observedAt && (
          <Typography variant="caption" color="textSecondary">
            observed {clock(state.observedAt)}
          </Typography>
        )}
      </Box>
      {props.hint && (
        <Typography variant="caption" color="textSecondary" display="block" style={{ padding: theme.spacing(1, 2, 0) }}>
          {props.hint}
        </Typography>
      )}
      {state && state.status !== 'ok' && (
        <Typography
          variant="body2"
          role={unavailable ? 'alert' : 'status'}
          style={{ padding: theme.spacing(1, 2, 0), color: 'var(--onemind-status-warn)' }}
        >
          {unavailable
            ? 'This source did not answer, so nothing is shown here. Unknown is not healthy.'
            : `Showing results from ${clock(state.observedAt)} (${age(state.observedAt)} ago): the latest refresh failed.`}{' '}
          {state.problems.map(p => `${p.source}: ${p.error ?? p.status}`).join('; ')}
        </Typography>
      )}
      <Box p={2} pt={1}>
        {props.children}
      </Box>
    </Box>
  );
}

export function Row(props: { children: ReactNode; tone?: 'critical' | 'warn' | 'ok' | 'none' }) {
  const theme = useTheme();
  const tone = props.tone && props.tone !== 'none' ? `var(--onemind-status-${props.tone})` : 'transparent';
  return (
    <Box
      py={1}
      pl={1.5}
      display="flex"
      flexWrap="wrap"
      alignItems="baseline"
      style={{
        gap: theme.spacing(0.5, 1.5),
        borderTop: `1px solid ${theme.palette.divider}`,
        borderLeft: `3px solid ${tone}`,
        minHeight: 40,
      }}
    >
      {props.children}
    </Box>
  );
}

export function Empty({ children }: { children: ReactNode }) {
  return (
    <Typography variant="body2" color="textSecondary" style={{ paddingTop: 8 }}>
      {children}
    </Typography>
  );
}

export function OwnerLink({ owner }: { owner?: Owner }) {
  return owner ? (
    <Link href={owner.href} variant="caption" aria-label={`Owning app ${owner.title}`}>
      {owner.title}
    </Link>
  ) : (
    <Typography variant="caption" color="textSecondary">
      Unassigned
    </Typography>
  );
}

function meterTone(value: number | null | undefined, warnAt: number, criticalAt: number) {
  if (value === null || value === undefined) return 'unknown';
  if (value >= criticalAt) return 'critical';
  return value >= warnAt ? 'warn' : 'ok';
}

/** A value with a bar. The number is always shown; colour is only a second cue. */
export function Meter(props: { label: string; value: number | null | undefined; warnAt?: number; criticalAt?: number }) {
  const { value, warnAt = 80, criticalAt = 90 } = props;
  const tone = meterTone(value, warnAt, criticalAt);
  return (
    <Box style={{ minWidth: 92 }}>
      <Typography variant="caption" color="textSecondary" display="block">
        {props.label} {value === null || value === undefined ? 'n/a' : `${value}%`}
      </Typography>
      <Box
        aria-hidden
        style={{ height: 4, borderRadius: 2, background: 'var(--onemind-border)', overflow: 'hidden' }}
      >
        <Box style={{ height: '100%', width: `${Math.min(100, Math.max(0, value ?? 0))}%`, background: `var(--onemind-status-${tone})` }} />
      </Box>
    </Box>
  );
}
