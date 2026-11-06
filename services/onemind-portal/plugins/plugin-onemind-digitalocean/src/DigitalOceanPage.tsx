import { ReactNode, useState } from 'react';
import {
  Box,
  Button,
  Link,
  Table,
  TableBody,
  TableCell,
  TableHead,
  TableRow,
  Typography,
  useTheme,
} from '@material-ui/core';
import { StatusChip } from '@internal/plugin-onemind-theme';
import { useSummary } from './api';
import type { Cluster, Section as Sect, Summary } from './types';

const usd = (value?: string) =>
  value === undefined ? 'Unknown' : `$${Number(value).toFixed(2)}`;
const day = (iso?: string) => (iso ? iso.slice(0, 10) : '');
const gb = (n: number) => `${Math.round(n * 100) / 100} GB`;

function Panel(props: { title: string; hint?: string; children: ReactNode }) {
  const theme = useTheme();
  return (
    <Box
      component="section"
      aria-label={props.title}
      mt={3}
      style={{
        border: `1px solid ${theme.palette.divider}`,
        borderRadius: theme.shape.borderRadius * 2,
        background: theme.palette.background.paper,
        overflowX: 'auto',
      }}
    >
      <Box
        px={2}
        py={1.5}
        style={{ borderBottom: `1px solid ${theme.palette.divider}` }}
      >
        <Typography variant="h6" component="h2" style={{ fontWeight: 600 }}>
          {props.title}
        </Typography>
        {props.hint && (
          <Typography variant="caption" color="textSecondary">
            {props.hint}
          </Typography>
        )}
      </Box>
      {props.children}
    </Box>
  );
}

/** A panel whose rows come from a section; shows the section's reason instead when it is not ok. */
function SectionPanel<T>(props: {
  title: string;
  hint?: string;
  section: Sect<T>;
  empty: string;
  header: string[];
  row: (item: T) => ReactNode[];
  keyOf: (item: T) => string;
  before?: ReactNode;
}) {
  const { section } = props;
  return (
    <Panel title={props.title} hint={props.hint}>
      {props.before}
      {section.state !== 'ok' && (
        <Box p={2} role="status">
          <StatusChip
            status="unknown"
            label={
              section.state === 'not-available'
                ? 'Not available'
                : 'Unavailable'
            }
          />
          <Typography component="span" style={{ marginLeft: 8 }}>
            {section.reason}
          </Typography>
        </Box>
      )}
      {section.state === 'ok' && section.items.length === 0 && (
        <Box p={2}>
          <Typography>{props.empty}</Typography>
        </Box>
      )}
      {section.state === 'ok' && section.items.length > 0 && (
        <Table size="small">
          <TableHead>
            <TableRow>
              {props.header.map(h => (
                <TableCell key={h}>{h}</TableCell>
              ))}
            </TableRow>
          </TableHead>
          <TableBody>
            {section.items.map(item => (
              <TableRow key={props.keyOf(item)}>
                {props.row(item).map((cell, i) => (
                  <TableCell key={props.header[i]}>{cell}</TableCell>
                ))}
              </TableRow>
            ))}
          </TableBody>
        </Table>
      )}
    </Panel>
  );
}

const noApp = (pvc?: unknown) => (pvc ? 'No app known' : '');

function Stat(props: { label: string; value: string; note?: string }) {
  return (
    <Box minWidth={170} mr={4} mb={1}>
      <Typography variant="caption" color="textSecondary" display="block">
        {props.label}
      </Typography>
      <Typography variant="h5" style={{ fontVariantNumeric: 'tabular-nums' }}>
        {props.value}
      </Typography>
      {props.note && (
        <Typography variant="caption" color="textSecondary">
          {props.note}
        </Typography>
      )}
    </Box>
  );
}

function NotConnected() {
  return (
    <Box
      mt={3}
      p={2}
      role="status"
      aria-label="Not connected"
      style={{ border: '1px solid currentColor', borderRadius: 12 }}
    >
      <StatusChip status="unknown" label="Not connected" />
      <Typography component="span" style={{ marginLeft: 8 }}>
        <strong>Not connected: needs the read-only DigitalOcean token.</strong>{' '}
        Seal it with{' '}
        <code>services/onemind-portal/scripts/seal-digitalocean.py</code> (steps
        in docs/operations/digitalocean.md). Nothing is read from DigitalOcean
        until then.
      </Typography>
    </Box>
  );
}

function Billing({ summary }: { summary: Summary }) {
  const b = summary.billing;
  return (
    <Panel
      title="Spend"
      hint="Credits can hide real usage, so usage before credits is shown next to the amount due."
    >
      {b.state !== 'ok' ? (
        <Box p={2} role="status">
          <StatusChip status="unknown" label="Unavailable" />
          <Typography component="span" style={{ marginLeft: 8 }}>
            {b.reason}
          </Typography>
        </Box>
      ) : (
        <Box p={2} display="flex" flexWrap="wrap">
          <Stat
            label="Usage before credits (this month)"
            value={usd(b.usageBeforeCredits)}
          />
          <Stat
            label="Amount due (this month)"
            value={usd(b.amountDue)}
            note="after credits"
          />
          <Stat
            label="Projected month"
            value={usd(b.projectedMonth)}
            note={
              b.daysElapsed !== undefined
                ? `usage before credits, scaled from ${b.daysElapsed} of ${b.daysInMonth} days`
                : undefined
            }
          />
          <Stat label="Account balance" value={usd(b.accountBalance)} />
        </Box>
      )}
      {b.invoices.length > 0 && (
        <Table size="small" aria-label="Invoices">
          <TableHead>
            <TableRow>
              <TableCell>Invoice</TableCell>
              <TableCell align="right">Amount</TableCell>
            </TableRow>
          </TableHead>
          <TableBody>
            {b.invoicePreview && (
              <TableRow>
                <TableCell>
                  {b.invoicePreview.period} (open, not final)
                </TableCell>
                <TableCell align="right">
                  {usd(b.invoicePreview.amount)}
                </TableCell>
              </TableRow>
            )}
            {b.invoices.map(i => (
              <TableRow key={i.period}>
                <TableCell>{i.period}</TableCell>
                <TableCell align="right">{usd(i.amount)}</TableCell>
              </TableRow>
            ))}
          </TableBody>
        </Table>
      )}
    </Panel>
  );
}

function Volumes({ summary }: { summary: Summary }) {
  const [only, setOnly] = useState(false);
  const v = summary.volumes;
  const shown = {
    ...v,
    items: only ? v.items.filter(x => x.unattached) : v.items,
  };
  return (
    <SectionPanel
      title="Volumes"
      hint="Reports only. Nothing here deletes or changes a volume."
      section={shown}
      empty={only ? 'No unattached volumes.' : 'No volumes.'}
      header={['Name', 'Size', 'Used by', 'Status', 'Created']}
      keyOf={x => x.id}
      before={
        <Box
          px={2}
          pt={1.5}
          display="flex"
          alignItems="center"
          style={{ gap: 12 }}
        >
          <Typography variant="body2">
            {v.items.length} volumes,{' '}
            {gb(v.items.reduce((n, x) => n + x.sizeGb, 0))};{' '}
            <strong>{v.unattachedCount} unattached</strong> (
            {gb(v.unattachedSizeGb)})
          </Typography>
          <Button
            size="small"
            variant={only ? 'contained' : 'outlined'}
            color="primary"
            aria-pressed={only}
            onClick={() => setOnly(!only)}
          >
            Unattached only
          </Button>
        </Box>
      }
      row={x => [
        x.pvc ? `${x.pvc.namespace}/${x.pvc.name}` : x.name,
        gb(x.sizeGb),
        x.app ? (
          <Link
            href={`/catalog/${x.app.ref
              .replace('component:', '')
              .replace('/', '/component/')}`}
          >
            {x.app.title ?? x.app.ref}
          </Link>
        ) : (
          noApp(x.pvc)
        ),
        x.unattached ? (
          <StatusChip status="warn" label="Unattached" />
        ) : (
          'Attached'
        ),
        day(x.createdAt),
      ]}
    />
  );
}

function Snapshots({ summary }: { summary: Summary }) {
  const s = summary.snapshots;
  const { stats } = s;
  return (
    <Panel title="Snapshots" hint="Summary only; reports, no deletes.">
      {s.state !== 'ok' ? (
        <Box p={2} role="status">
          <StatusChip status="unknown" label="Unavailable" />
          <Typography component="span" style={{ marginLeft: 8 }}>
            {s.reason}
          </Typography>
        </Box>
      ) : (
        <Box p={2} display="flex" flexWrap="wrap">
          <Stat
            label="Snapshots"
            value={String(stats.count)}
            note={`${gb(stats.totalSizeGb)} in total`}
          />
          <Stat
            label="Oldest"
            value={stats.oldest ? `${stats.oldest.ageDays} days` : 'None'}
            note={
              stats.oldest
                ? `${stats.oldest.name}, ${day(stats.oldest.createdAt)}`
                : undefined
            }
          />
          <Stat
            label="Largest"
            value={stats.largest ? gb(stats.largest.sizeGb) : 'None'}
            note={stats.largest?.name}
          />
          <Stat
            label={`Older than ${stats.olderThanDays} days`}
            value={String(stats.olderCount)}
          />
        </Box>
      )}
    </Panel>
  );
}

function Changes({ summary }: { summary: Summary }) {
  const list = summary.changes;
  return (
    <Panel
      title="What changed"
      hint="Resources added or removed between refreshes since the portal started."
    >
      {list.length === 0 ? (
        <Box p={2}>
          <Typography>No changes seen since the portal started.</Typography>
        </Box>
      ) : (
        <Box p={2}>
          {list.map(c => (
            <Box key={c.at} mb={1}>
              <Typography variant="caption" color="textSecondary">
                {c.at}
              </Typography>
              {c.added.map(a => (
                <Typography key={`+${a.type}${a.id}`}>
                  + added {a.type}: {a.name}
                </Typography>
              ))}
              {c.removed.map(a => (
                <Typography key={`-${a.type}${a.id}`}>
                  - removed {a.type}: {a.name}
                </Typography>
              ))}
            </Box>
          ))}
        </Box>
      )}
    </Panel>
  );
}

function upgradeCell(c: Cluster) {
  if (c.upgradeCheck === 'unavailable') return 'Unknown';
  if (!c.upgradeAvailable) return 'Up to date';
  return (
    <StatusChip status="warn" label={`Available: ${c.upgradeAvailable}`} />
  );
}

const UP_TEXT = { up: 'Up', down: 'Down', unknown: 'Unknown' } as const;
const upLabel = { up: 'ok', down: 'critical', unknown: 'unknown' } as const;

export function DigitalOceanView({ summary }: { summary: Summary }) {
  if (summary.state === 'not-connected') return <NotConnected />;
  return (
    <>
      <Billing summary={summary} />
      <SectionPanel
        title="Kubernetes cluster and node pools"
        section={summary.kubernetes}
        empty="No clusters."
        header={['Cluster', 'Version', 'Upgrade', 'Node pools']}
        keyOf={c => c.id}
        row={c => [
          `${c.name} (${c.region}${c.highAvailability ? ', HA' : ''})`,
          c.version,
          upgradeCell(c),
          c.nodePools
            .map(
              p =>
                `${p.name}: ${p.count} x ${p.size}${
                  p.autoScale
                    ? ` (autoscale ${p.minNodes} to ${p.maxNodes})`
                    : ''
                }`,
            )
            .join('; '),
        ]}
      />
      <SectionPanel
        title="Droplets"
        section={summary.droplets}
        empty="No droplets."
        header={['Name', 'Status', 'Size', 'Region', 'Tags']}
        keyOf={d => d.id}
        row={d => [d.name, d.status, d.size, d.region, d.tags.join(', ')]}
      />
      <Volumes summary={summary} />
      <Snapshots summary={summary} />
      <SectionPanel
        title="Load balancers"
        section={summary.loadBalancers}
        empty="No load balancers."
        header={['Name', 'Region', 'Rules', 'Status']}
        keyOf={l => l.id}
        row={l => [
          l.name,
          l.region,
          String(l.forwardingRules),
          l.idle ? <StatusChip status="warn" label="Idle" /> : l.status,
        ]}
      />
      <SectionPanel
        title="Firewalls"
        section={summary.firewalls}
        empty="No firewalls."
        header={['Name', 'Status', 'Inbound rules', 'Outbound rules']}
        keyOf={f => f.id}
        row={f => [
          f.name,
          f.status,
          String(f.inboundRules),
          String(f.outboundRules),
        ]}
      />
      <SectionPanel
        title="VPCs"
        section={summary.vpcs}
        empty="No VPCs."
        header={['Name', 'Region', 'IP range', 'Default']}
        keyOf={v => v.id}
        row={v => [v.name, v.region, v.ipRange, v.isDefault ? 'Yes' : '']}
      />
      <SectionPanel
        title="Spaces buckets"
        section={summary.spaces}
        empty=""
        header={['Bucket']}
        keyOf={() => 'none'}
        row={() => ['']}
      />
      <SectionPanel
        title="Domains"
        section={summary.domains}
        empty="No domains in DigitalOcean DNS."
        header={['Domain', 'DNS records']}
        keyOf={d => d.name}
        row={d => [
          d.name,
          d.recordCount === null ? 'Unknown' : String(d.recordCount),
        ]}
      />
      <SectionPanel
        title="Uptime checks"
        section={summary.uptime}
        empty="No uptime checks."
        header={['Name', 'Target', 'State']}
        keyOf={u => u.id}
        row={u => [
          u.name,
          u.target,
          <StatusChip
            key="s"
            status={upLabel[u.state]}
            label={UP_TEXT[u.state]}
          />,
        ]}
      />
      <SectionPanel
        title="Alert policies"
        section={summary.alerts}
        empty="No monitoring alert policies are set."
        header={['Description', 'Type', 'Enabled']}
        keyOf={a => a.id}
        row={a => [a.description, a.type, a.enabled ? 'Yes' : 'No']}
      />
      <Changes summary={summary} />
    </>
  );
}

export function DigitalOceanPage() {
  const { state, refresh } = useSummary();
  return (
    <Box p={3}>
      <Typography variant="overline">ONE MIND · SERVICES</Typography>
      <Typography variant="h4" component="h1">
        DigitalOcean
      </Typography>
      <Typography color="textSecondary">
        What the account runs and costs, read-only. Refreshed about every 10
        minutes.
      </Typography>
      <Box mt={1}>
        <Button size="small" onClick={refresh}>
          Reload
        </Button>
        {state.status === 'done' && state.summary.state === 'connected' && (
          <Typography variant="caption" color="textSecondary">
            {' '}
            Data from {state.summary.generatedAt}
          </Typography>
        )}
      </Box>
      {state.status === 'loading' && (
        <Typography role="status">Loading DigitalOcean…</Typography>
      )}
      {state.status === 'forbidden' && (
        <Typography role="alert">
          Only OneMind administrators can see this page.
        </Typography>
      )}
      {state.status === 'error' && (
        <Typography role="alert">
          The portal could not load the DigitalOcean summary.
        </Typography>
      )}
      {state.status === 'done' && <DigitalOceanView summary={state.summary} />}
    </Box>
  );
}
