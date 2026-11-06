import { useState } from 'react';
import { Box, Button, Chip, Typography } from '@material-ui/core';
import { FavoritesPicker, useAppFavorites } from './Launcher';
import { HomeCard } from './HomeCard';

export function PinnedChips() {
  const controls = useAppFavorites();
  const [open, setOpen] = useState(false);
  return (
    <HomeCard
      label="Pinned"
      title="Pinned"
      action={
        <Button size="small" color="primary" onClick={() => setOpen(true)}>
          Manage
        </Button>
      }
    >
      {controls.error && (
        <Typography
          variant="caption"
          color="error"
          role="alert"
          display="block"
        >
          {controls.error}
        </Typography>
      )}
      {!controls.loaded && (
        <Typography variant="caption" color="textSecondary" role="status">
          Loading saved favorites…
        </Typography>
      )}
      {controls.loaded && !controls.favorites.length && (
        <Typography variant="caption" color="textSecondary">
          No pinned apps yet. Use Manage to choose some.
        </Typography>
      )}
      {controls.loaded && !!controls.favorites.length && (
        <Box display="flex" flexWrap="wrap" style={{ gap: 6 }}>
          {controls.favorites.map(app => (
            <Chip
              key={app.name}
              size="small"
              clickable
              component="a"
              href={app.url}
              target="_blank"
              rel="noopener noreferrer"
              label={app.title}
              variant="outlined"
            />
          ))}
        </Box>
      )}
      <FavoritesPicker
        open={open}
        onClose={() => setOpen(false)}
        controls={controls}
      />
    </HomeCard>
  );
}
