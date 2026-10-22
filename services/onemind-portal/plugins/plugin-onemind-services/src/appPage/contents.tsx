import { configApiRef, useApi } from '@backstage/core-plugin-api';
import { useEntity } from '@backstage/plugin-catalog-react';
import { InfoCard, Link } from '@backstage/core-components';
import { Box, Button, Typography } from '@material-ui/core';
import { InfoTooltip } from '@internal/plugin-onemind-common';
import { buildAppLinks } from './links';

const code: React.CSSProperties = {
  display: 'block',
  padding: 8,
  overflowX: 'auto',
  fontFamily: 'monospace',
  fontSize: 13,
  borderRadius: 6,
  border: '1px solid var(--onemind-border, currentColor)',
  background: 'var(--onemind-surface-2, transparent)',
};

function useLinks() {
  const { entity } = useEntity();
  const grafana = useApi(configApiRef).getOptionalString('grafana.domain');
  return { entity, links: buildAppLinks(entity, grafana) };
}

/** Logs tab: no new log store. The generated LogQL opens in Grafana Explore against the existing Loki. */
export function LogsContent() {
  const { links } = useLinks();
  return (
    <InfoCard
      title="Logs"
      action={
        <InfoTooltip label="logs">
          Logs are read from the existing Loki through Grafana Explore; this portal stores none. The query uses the
          namespace and workload names recorded in the catalog. An empty result can mean the app writes nothing, or that
          log collection is not covering it; it does not prove the app is quiet.
        </InfoTooltip>
      }
    >
      <Typography variant="body2" gutterBottom>Loki query for this app</Typography>
      <code style={code} data-testid="logs-query">{links.logsQuery}</code>
      <Box mt={2} display="flex" style={{ gap: 12 }} flexWrap="wrap">
        {links.logsExplore && (
          <Button variant="contained" color="primary" href={links.logsExplore} target="_blank" rel="noopener noreferrer">
            Open in Grafana Explore
          </Button>
        )}
        <Typography variant="body2" component="span" style={{ alignSelf: 'center' }}>
          Live pod logs: open the <b>Live</b> tab, then a pod.
        </Typography>
      </Box>
    </InfoCard>
  );
}

/** Cost tab: links to OpenCost until the per-app allocation card (F9) lands. Shown only for cluster apps. */
export function CostContent() {
  const { entity, links } = useLinks();
  const namespaces = (entity.metadata.annotations?.['onemind.io/kubernetes-namespaces'] ?? '').split(',').filter(Boolean);
  return (
    <InfoCard
      title="Cost"
      action={
        <InfoTooltip label="cost">
          Cost is allocated by OpenCost per namespace. A namespace can hold several apps, so a namespace figure is an
          upper bound for any one of them. Numbers are not shown here until the OpenCost card is added (work order F9).
        </InfoTooltip>
      }
    >
      <Typography variant="body2" gutterBottom>
        Namespace{namespaces.length === 1 ? '' : 's'}: {namespaces.join(', ')}
      </Typography>
      {links.opencost && (
        <Button variant="contained" color="primary" href={links.opencost} target="_blank" rel="noopener noreferrer">
          Open OpenCost
        </Button>
      )}
    </InfoCard>
  );
}

/** Docs tab entry: the generated runbook page in the existing TechDocs runbook site. */
export function RunbookContent() {
  const { entity, links } = useLinks();
  return (
    <InfoCard
      title="Runbook"
      action={
        <InfoTooltip label="runbook">
          Generated from the catalog: where the app runs, its workloads, and the safe recovery path (Git and Argo CD
          only). It records what is configured, not what is healthy.
        </InfoTooltip>
      }
    >
      <Typography variant="body2" gutterBottom>
        Recovery guide for {entity.metadata.title ?? entity.metadata.name}.
      </Typography>
      {links.runbook && <Link to={links.runbook}>Open the runbook</Link>}
    </InfoCard>
  );
}
