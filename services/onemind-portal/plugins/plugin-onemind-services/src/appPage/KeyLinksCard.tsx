import { configApiRef, useApi } from '@backstage/core-plugin-api';
import { useEntity } from '@backstage/plugin-catalog-react';
import { InfoCard, Link } from '@backstage/core-components';
import { List, ListItem } from '@material-ui/core';
import { buildAppLinks } from './links';

/** Key links generated from the coverage annotations. Hidden when an entity has none. */
export function KeyLinksCard() {
  const { entity } = useEntity();
  const grafana = useApi(configApiRef).getOptionalString('grafana.domain');
  const links = buildAppLinks(entity, grafana);
  const items: { label: string; to: string }[] = [
    ...(links.openApp ? [{ label: 'Open app', to: links.openApp }] : []),
    ...links.argo.map(a => ({ label: `Argo CD · ${a.name}`, to: a.url })),
    ...(links.grafanaDashboard ? [{ label: 'Workload dashboard (Grafana)', to: links.grafanaDashboard }] : []),
    ...(links.logsExplore ? [{ label: 'Logs (Loki in Grafana Explore)', to: links.logsExplore }] : []),
    ...(links.runbook ? [{ label: 'Runbook', to: links.runbook }] : []),
    ...(links.opencost ? [{ label: 'Cost (OpenCost)', to: links.opencost }] : []),
    ...(links.source ? [{ label: 'Source repository', to: links.source }] : []),
  ];
  if (!items.length) return null;
  return (
    <InfoCard title="Key links">
      <List dense disablePadding>
        {items.map(item => (
          <ListItem key={item.label} disableGutters>
            <Link to={item.to}>{item.label}</Link>
          </ListItem>
        ))}
      </List>
    </InfoCard>
  );
}
