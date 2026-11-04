import { ReactElement } from 'react';
import CheckCircleIcon from '@material-ui/icons/CheckCircle';
import ErrorIcon from '@material-ui/icons/Error';
import HelpIcon from '@material-ui/icons/Help';
import InfoIcon from '@material-ui/icons/Info';
import WarningIcon from '@material-ui/icons/Warning';
import { StatusKey, CSS_PREFIX } from './tokens';

const icons: Record<StatusKey, ReactElement> = {
  ok: <CheckCircleIcon fontSize="inherit" />,
  warn: <WarningIcon fontSize="inherit" />,
  critical: <ErrorIcon fontSize="inherit" />,
  unknown: <HelpIcon fontSize="inherit" />,
  info: <InfoIcon fontSize="inherit" />,
};

const defaultLabels: Record<StatusKey, string> = {
  ok: 'Healthy',
  warn: 'Degraded',
  critical: 'Critical',
  unknown: 'Unknown',
  info: 'Info',
};

export interface StatusChipProps {
  status: StatusKey;
  label?: string;
}

/**
 * Health chip. The icon shape carries the meaning as well as the colour, so
 * critical (rose red with an alert icon) is never confused with the brand red.
 * Colours come from the `--onemind-status-*` variables, so it follows the theme.
 */
export function StatusChip({ status, label }: StatusChipProps) {
  const fg = `var(${CSS_PREFIX}status-${status})`;
  return (
    <span
      data-status={status}
      style={{
        display: 'inline-flex',
        alignItems: 'center',
        gap: 4,
        padding: '1px 8px 1px 6px',
        borderRadius: 6,
        fontSize: 12,
        fontWeight: 600,
        lineHeight: '20px',
        color: fg,
        background: `var(${CSS_PREFIX}status-${status}-bg)`,
        border: `1px solid ${fg}`,
      }}
    >
      {icons[status]}
      {label ?? defaultLabels[status]}
    </span>
  );
}
