import { Chip, Tooltip, useTheme } from '@material-ui/core';
import { hostingInfo } from './hosting';
export function HostingBadge({
  hosting,
  status,
}: {
  hosting?: string;
  status?: string;
}) {
  const info = hostingInfo(hosting, status);
  const dark = useTheme().palette.type === 'dark';
  const colors = {
    'Self-hosted': ['#166534', '#dcfce7', '#86efac', '#143d28'],
    'Cloud service': ['#1e40af', '#dbeafe', '#93c5fd', '#172f50'],
    'Local / native': ['#115e59', '#ccfbf1', '#5eead4', '#123c3a'],
    Reference: ['#6b21a8', '#f3e8ff', '#d8b4fe', '#38234f'],
    Planned: ['#92400e', '#fef3c7', '#fcd34d', '#453416'],
    Unknown: ['#475569', '#e2e8f0', '#cbd5e1', '#293342'],
  }[info.group];
  const foreground = colors[dark ? 2 : 0];
  const background = colors[dark ? 3 : 1];
  return (
    <Tooltip title={info.detail}>
      <Chip
        size="small"
        variant="outlined"
        label={info.label}
        style={{ marginTop: 6, maxWidth: '100%', color: foreground, backgroundColor: background, borderColor: foreground, fontWeight: 600 }}
      />
    </Tooltip>
  );
}
