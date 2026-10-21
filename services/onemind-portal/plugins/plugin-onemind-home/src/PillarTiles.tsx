import { useEffect, useState } from 'react';
import { Box, Link, Typography, useTheme } from '@material-ui/core';
import { useApi } from '@backstage/core-plugin-api';
import { catalogApiRef } from '@backstage/plugin-catalog-react';

// The four pillar Domains created by F2. Order is the sidebar order.
export const PILLARS = [
  { slug: 'so', short: 'SO', title: 'Sovereign Operations' },
  { slug: 'le', short: 'LE', title: 'Legacy Evolution' },
  { slug: 'hp', short: 'HP', title: 'Holistic Performance' },
  { slug: 'ge', short: 'GE', title: 'Generational Empire' },
] as const;

const PILLAR_LABEL = 'onemind.io/pillar';

/** Counts Components by their `onemind.io/pillar` label (F2 copies it from the System). */
export function countAppsByPillar(
  components: { labels?: Record<string, string> }[],
) {
  const counts: Record<string, number> = {};
  for (const pillar of PILLARS) counts[pillar.slug] = 0;
  for (const component of components) {
    const value = component.labels?.[PILLAR_LABEL];
    if (value && value in counts) counts[value] += 1;
  }
  return counts;
}

function countLabel(count: number | undefined, error: boolean) {
  if (count !== undefined) return `${count} apps`;
  return error ? 'Unavailable' : '—';
}

export function PillarTiles() {
  const catalog = useApi(catalogApiRef);
  const theme = useTheme();
  const [counts, setCounts] = useState<Record<string, number>>();
  const [error, setError] = useState(false);
  useEffect(() => {
    let active = true;
    (async () => {
      try {
        const components = await catalog.getEntities({
          filter: {
            kind: 'Component',
            [`metadata.labels.${PILLAR_LABEL}`]: PILLARS.map(p => p.slug),
          },
          fields: ['metadata.labels'],
        });
        if (!active) return;
        setCounts(
          countAppsByPillar(
            components.items.map(item => ({ labels: item.metadata.labels })),
          ),
        );
      } catch {
        if (active) setError(true);
      }
    })();
    return () => {
      active = false;
    };
  }, [catalog]);
  return (
    <Box component="section" aria-label="Pillars">
      <Typography variant="overline" color="textSecondary">
        Pillars
      </Typography>
      <Box
        style={{
          display: 'grid',
          gridTemplateColumns: 'repeat(auto-fit,minmax(min(100%,150px),1fr))',
          gap: theme.spacing(1.5),
        }}
      >
        {PILLARS.map(pillar => (
          <Link
            key={pillar.slug}
            href={`/catalog/default/domain/onemind-${pillar.slug}`}
            underline="none"
            color="inherit"
            aria-label={`${pillar.title}: ${
              counts ? `${counts[pillar.slug]} apps` : 'apps unknown'
            }`}
            style={{
              display: 'block',
              minHeight: 44,
              padding: theme.spacing(1.5, 2),
              border: `1px solid ${theme.palette.divider}`,
              borderTop: `3px solid ${theme.palette.primary.main}`,
              borderRadius: theme.shape.borderRadius * 2,
              background: theme.palette.background.paper,
            }}
          >
            <Typography variant="subtitle2" style={{ fontWeight: 700 }}>
              {pillar.short}
            </Typography>
            <Typography variant="caption" color="textSecondary" display="block">
              {pillar.title}
            </Typography>
            <Typography variant="body2">
              {countLabel(counts?.[pillar.slug], error)}
            </Typography>
          </Link>
        ))}
      </Box>
    </Box>
  );
}
