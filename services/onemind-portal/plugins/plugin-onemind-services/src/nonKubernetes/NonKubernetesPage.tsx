import { useEffect, useState } from 'react';
import { useApi } from '@backstage/core-plugin-api';
import { catalogApiRef, EntityRefLink } from '@backstage/plugin-catalog-react';
import { Content, Header, Page } from '@backstage/core-components';
import { Box, Table, TableBody, TableCell, TableHead, TableRow, Typography } from '@material-ui/core';
import { StatusDot } from '@internal/plugin-onemind-common';
import { GatusEndpoint, summariseEndpoint, useGatusReader } from '../appPage/gatus';
import { AppRow, groupByHosting, nonKubernetesRows } from './groups';

function CheckCell({ row, statuses }: { row: AppRow; statuses?: GatusEndpoint[] }) {
  if (!row.check) return <Typography variant="body2" color="textSecondary">No check</Typography>;
  if (!statuses) return <StatusDot status="unknown" showLabel title="Gatus is not reachable" />;
  const summary = summariseEndpoint(statuses, row.check);
  return summary.kind === 'checked' ? (
    <StatusDot status={summary.status} showLabel />
  ) : (
    <StatusDot status="unknown" showLabel title="Gatus has no result for this check yet" />
  );
}

/** Apps that are not on Kubernetes, grouped by hosting type. Unknown stays visible; nothing is promoted. */
export function NonKubernetesPage() {
  const catalog = useApi(catalogApiRef);
  const read = useGatusReader();
  const [rows, setRows] = useState<AppRow[]>();
  const [failed, setFailed] = useState(false);
  const [statuses, setStatuses] = useState<GatusEndpoint[]>();
  useEffect(() => {
    let active = true;
    catalog
      .getEntities({
        filter: { kind: 'Component' },
        fields: ['kind', 'metadata.name', 'metadata.namespace', 'metadata.title', 'metadata.annotations', 'spec.owner'],
      })
      .then(result => active && setRows(nonKubernetesRows(result.items as never)))
      .catch(() => active && setFailed(true));
    read().then(value => active && setStatuses(value)).catch(() => undefined);
    return () => {
      active = false;
    };
  }, [catalog, read]);
  const groups = rows ? groupByHosting(rows) : [];
  return (
    <Page themeId="home">
      <Header title="Not on Kubernetes" subtitle="Apps that run elsewhere, grouped by where they run" />
      <Content>
        {failed && <Typography role="alert" color="error">The catalog could not be loaded.</Typography>}
        {!rows && !failed && <Typography role="status">Reading the catalog…</Typography>}
        {rows && (
          <Typography variant="body2" color="textSecondary" gutterBottom>
            {rows.length} apps. “Unknown” means hosting has not been established; it is shown, not hidden. Uptime comes from Gatus
            where an app has a check.
          </Typography>
        )}
        {groups.map(group => (
          <Box key={group.kind} mt={3} data-testid={`group-${group.kind}`}>
            <Typography variant="h6">
              {group.label} ({group.rows.length})
            </Typography>
            <Table size="small" aria-label={group.label}>
              <TableHead>
                <TableRow>
                  <TableCell>App</TableCell>
                  <TableCell>Owner</TableCell>
                  <TableCell>Uptime check</TableCell>
                </TableRow>
              </TableHead>
              <TableBody>
                {group.rows.map(row => (
                  <TableRow key={row.ref}>
                    <TableCell>
                      <EntityRefLink entityRef={row.ref} title={row.title} />
                    </TableCell>
                    <TableCell>{row.owner ?? 'Not recorded'}</TableCell>
                    <TableCell>
                      <CheckCell row={row} statuses={statuses} />
                    </TableCell>
                  </TableRow>
                ))}
              </TableBody>
            </Table>
          </Box>
        ))}
      </Content>
    </Page>
  );
}
