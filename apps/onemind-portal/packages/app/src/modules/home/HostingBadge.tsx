import { Chip, Tooltip } from '@material-ui/core';
import { hostingInfo } from './hosting';
export function HostingBadge({
  hosting,
  status,
}: {
  hosting?: string;
  status?: string;
}) {
  const info = hostingInfo(hosting, status);
  return (
    <Tooltip title={info.detail}>
      <Chip
        size="small"
        variant="outlined"
        label={info.label}
        style={{ marginTop: 6, maxWidth: '100%' }}
      />
    </Tooltip>
  );
}
