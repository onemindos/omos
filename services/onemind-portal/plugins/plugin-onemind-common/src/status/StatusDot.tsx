import { CSSProperties } from 'react';
import type { AppStatus } from './model';

const LABEL: Record<AppStatus, string> = { ok: 'Healthy', warn: 'Degraded', down: 'Down', unknown: 'Unknown' };
// The shape carries the meaning as well as the colour (WCAG 1.4.1): circle, diamond, square, hollow ring.
const TOKEN: Record<AppStatus, string> = { ok: 'ok', warn: 'warn', down: 'critical', unknown: 'unknown' };

export type StatusDotProps = {
  status: AppStatus;
  /** Visible text next to the dot. */
  showLabel?: boolean;
  /** Extra detail for assistive technology and the hover title. */
  title?: string;
  size?: number;
};

/** Grey (hollow) for unknown, never green. Colours are the portal's `--onemind-status-*` tokens. */
export function StatusDot({ status, showLabel, title, size = 10 }: StatusDotProps) {
  const color = `var(--onemind-status-${TOKEN[status]})`;
  const shape: CSSProperties =
    status === 'ok'
      ? { borderRadius: '50%', background: color }
      : status === 'warn'
      ? { transform: 'rotate(45deg) scale(0.85)', background: color }
      : status === 'down'
      ? { background: color }
      : { borderRadius: '50%', border: `2px solid ${color}`, background: 'transparent', boxSizing: 'border-box' };
  const text = `Status: ${LABEL[status]}${title ? `. ${title}` : ''}`;
  return (
    <span
      role="img"
      aria-label={text}
      title={text}
      data-status={status}
      style={{ display: 'inline-flex', alignItems: 'center', gap: 6, fontSize: 13, color: 'inherit' }}
    >
      <span aria-hidden style={{ display: 'inline-block', width: size, height: size, flex: 'none', ...shape }} />
      {showLabel && <span>{LABEL[status]}</span>}
    </span>
  );
}
