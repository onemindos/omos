import { Box, Link, Typography, useTheme } from '@material-ui/core';
import { configApiRef, useApi } from '@backstage/core-plugin-api';
import { GRAFANA_PANELS, panelEmbedUrl, panelPageUrl } from './grafanaConfig';
import { Panel } from './Panel';

/**
 * Embedded Grafana panels. Embedding follows the same switch the home page uses
 * (`onemind.grafanaEmbedding.enabled`); with it off, each card is a link to the panel.
 */
export function GrafanaEmbeds({ origin }: { origin: string }) {
  const theme = useTheme();
  const embed = useApi(configApiRef).getOptionalBoolean('onemind.grafanaEmbedding.enabled') === true;
  const dark = theme.palette.type === 'dark';
  return (
    <Panel
      title="Grafana"
      wide
      hint={
        embed
          ? 'Cluster CPU and memory, network, NATS and Postgres. Panels load from Grafana and use your Grafana session.'
          : 'Embedded viewing awaits secure Grafana configuration (onemind.grafanaEmbedding.enabled). Panels open in Grafana.'
      }
    >
      <Box
        style={{
          display: 'grid',
          gridTemplateColumns: 'repeat(auto-fit, minmax(min(100%, 320px), 1fr))',
          gap: theme.spacing(1.5),
        }}
      >
        {GRAFANA_PANELS.map(panel => (
          <Box key={`${panel.uid}-${panel.panelId}`} style={{ minWidth: 0 }}>
            <Typography variant="subtitle2" component="h3">
              {panel.title}
            </Typography>
            {embed ? (
              <iframe
                title={`Grafana: ${panel.title}`}
                src={panelEmbedUrl(origin, panel, dark)}
                loading="lazy"
                referrerPolicy="no-referrer"
                sandbox="allow-scripts allow-same-origin"
                style={{
                  width: '100%',
                  height: 200,
                  border: `1px solid ${theme.palette.divider}`,
                  borderRadius: theme.shape.borderRadius,
                  background: theme.palette.background.default,
                }}
              />
            ) : null}
            <Link href={panelPageUrl(origin, panel)} variant="caption" target="_blank" rel="noopener noreferrer">
              Open in Grafana
            </Link>
          </Box>
        ))}
      </Box>
    </Panel>
  );
}
