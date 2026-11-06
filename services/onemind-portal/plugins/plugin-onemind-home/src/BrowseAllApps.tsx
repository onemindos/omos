import { ReactNode, useState } from 'react';
import { storageApiRef, useApi } from '@backstage/core-plugin-api';
import { Box, Button, Typography, useTheme } from '@material-ui/core';
import ExpandMoreIcon from '@material-ui/icons/ExpandMore';
import ChevronRightIcon from '@material-ui/icons/ChevronRight';
import { useDirectory } from '@internal/plugin-onemind-common';

const KEY = 'browseAllOpen';

/** Collapsed-by-default wrapper; the 75-app launcher is not mounted until opened. */
export function BrowseAllApps({ children }: { children: ReactNode }) {
  const theme = useTheme();
  const bucket = useApi(storageApiRef).forBucket('onemind-home');
  const { directory } = useDirectory();
  const [open, setOpen] = useState(
    () => bucket.snapshot<boolean>(KEY).value === true,
  );
  const toggle = () => {
    const next = !open;
    setOpen(next);
    bucket.set(KEY, next).catch(() => undefined);
  };
  return (
    <Box component="section" aria-label="Browse all apps">
      <Button
        onClick={toggle}
        aria-expanded={open}
        startIcon={open ? <ExpandMoreIcon /> : <ChevronRightIcon />}
        style={{
          textTransform: 'none',
          border: `1px solid ${theme.palette.divider}`,
          borderRadius: 14,
          justifyContent: 'flex-start',
          width: '100%',
          padding: theme.spacing(1, 2),
          background: theme.palette.background.paper,
        }}
      >
        <Typography
          variant="subtitle1"
          component="span"
          style={{ fontWeight: 700 }}
        >
          Browse all apps ({directory.length})
        </Typography>
      </Button>
      {open && <Box mt={2}>{children}</Box>}
    </Box>
  );
}
