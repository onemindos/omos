import { useMemo, useState } from 'react';
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
import type { AllocationState, NamespaceAllocation } from './allocation';
import { CostData, DEFAULT_DAYS, useCostData, WINDOW_CHOICES } from './costApi';
import {
  buildCostLines,
  CostLine,
  pillarLabel,
  pillarRollups,
  rollup,
} from './mapping';

/** "$12.34", "<$0.01" for a real but tiny amount, "Unknown" for no usable number. */
export function money(value: number | null | undefined): string {
  if (value === null || value === undefined || !Number.isFinite(value))
    return 'Unknown';
  if (value > 0 && value < 0.005) return '<$0.01';
  return `$${value.toFixed(2)}`;
}

const STATE_TEXT: Record<
  Exclude<AllocationState, 'available'>,
  { level: 'warn' | 'critical' | 'unknown'; text: string }
> = {
  partial: {
    level: 'warn',
    text: 'OpenCost returned incomplete data. Totals below are lower bounds and some figures are marked Unknown.',
  },
  empty: {
    level: 'unknown',
    text: 'OpenCost returned no namespaces for this window. That is not the same as zero cost.',
  },
  unknown: {
    level: 'critical',
    text: 'OpenCost could not be reached through the portal proxy, so no cost is known.',
  },
  'source-error': {
    level: 'critical',
    text: 'OpenCost reported an error, so no cost is known.',
  },
  invalid: {
    level: 'critical',
    text: 'OpenCost returned data in a shape this page does not accept, so no cost is shown.',
  },
};

function Section(props: {
  title: string;
  hint?: string;
  children: React.ReactNode;
}) {
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

function Amount({ cost, partial }: { cost: number | null; partial?: boolean }) {
  return (
    <span
      style={{
        fontVariantNumeric: 'tabular-nums',
        opacity: cost === null ? 0.7 : 1,
      }}
    >
      {partial && cost !== null ? 'at least ' : ''}
      {money(cost)}
    </span>
  );
}

function activeNote(line: CostLine, windowHours: number): string {
  const { row } = line;
  if (row.minutes === null) return 'Unknown';
  const hours = row.minutes / 60;
  if (Math.abs(hours - windowHours) < 0.02) return 'Whole window';
  return `${hours.toFixed(hours < 10 ? 1 : 0)} h of ${windowHours} h`;
}

function Attribution({ line }: { line: CostLine }) {
  if (line.attribution === 'none')
    return <span>No catalog app declares this namespace</span>;
  if (line.attribution === 'app') {
    const [app] = line.apps;
    return <Link href={app.href}>{app.title}</Link>;
  }
  return (
    <span title={line.apps.map(app => app.title).join(', ')}>
      Shared by {line.apps.length} apps (cost not split)
    </span>
  );
}

function pillarCell(line: CostLine) {
  return (
    <>
      {pillarLabel(line.pillar)}
      {line.pillarProposed && (
        <Typography variant="caption" color="textSecondary" display="block">
          proposed, not approved
        </Typography>
      )}
    </>
  );
}

export function CostView({ data }: { data: CostData }) {
  const allocation: NamespaceAllocation | undefined = data.allocation;
  const lines = useMemo(
    () => (allocation ? buildCostLines(allocation.rows, data.apps) : []),
    [allocation, data.apps],
  );
  const pillars = useMemo(() => pillarRollups(lines), [lines]);
  if (!allocation) return null;

  const total = rollup(lines);
  // OpenCost always adds `__unmounted__`-style lines. They are listed on their own
  // below, so on their own they do not make the namespace figures partial.
  const onlySpecialLines =
    allocation.state === 'partial' &&
    lines.length > 0 &&
    allocation.rows.every(row => row.state === 'available') &&
    allocation.issues.every(issue => issue === 'special-allocation-excluded');
  const known = allocation.state === 'available' || onlySpecialLines;
  const partial = allocation.state === 'partial' && !onlySpecialLines;
  const hasNumbers = known || (partial && lines.length > 0);
  const notice = known
    ? undefined
    : STATE_TEXT[allocation.state as Exclude<AllocationState, 'available'>];
  const unattributed = data.specials;
  const windowDays = allocation.window.hours / 24;

  return (
    <>
      {notice && (
        <Box mt={2} role={notice.level === 'critical' ? 'alert' : 'status'}>
          <StatusChip
            status={notice.level}
            label={allocation.state.replace('-', ' ')}
          />
          <Typography component="span" style={{ marginLeft: 8 }}>
            {notice.text}
          </Typography>
        </Box>
      )}
      {data.catalogError && (
        <Box mt={2} role="status">
          <Typography>
            The catalog could not be read, so namespaces are shown without apps
            or pillars.
          </Typography>
        </Box>
      )}
      {!known && allocation.issues.length > 0 && hasNumbers && (
        <Typography variant="caption" color="textSecondary" display="block">
          Validation notes: {[...new Set(allocation.issues)].join(', ')}
        </Typography>
      )}

      <Box
        mt={3}
        display="flex"
        flexWrap="wrap"
        style={{ gap: 24 }}
        aria-label="Cost summary"
      >
        <Box>
          <Typography variant="overline">
            Namespaces · {windowDays} day{windowDays === 1 ? '' : 's'}
          </Typography>
          <Typography variant="h4" component="p">
            {hasNumbers ? (
              <Amount
                cost={total.known}
                partial={!known || total.unknown > 0}
              />
            ) : (
              'Unknown'
            )}
          </Typography>
          {hasNumbers && total.unknown > 0 && (
            <Typography variant="caption" color="textSecondary">
              {total.unknown} namespace{total.unknown === 1 ? '' : 's'} with
              unknown cost are not included
            </Typography>
          )}
        </Box>
        <Box>
          <Typography variant="overline">Not owned by a namespace</Typography>
          <Typography variant="h4" component="p">
            {unattributed.length ? (
              <Amount
                cost={
                  unattributed.some(item => item.totalCost === null)
                    ? null
                    : unattributed.reduce(
                        (sum, item) => sum + (item.totalCost ?? 0),
                        0,
                      )
                }
              />
            ) : (
              'None reported'
            )}
          </Typography>
          <Typography variant="caption" color="textSecondary">
            {unattributed.length
              ? `OpenCost lines: ${unattributed
                  .map(item => item.name)
                  .join(', ')} (for example load balancers)`
              : 'OpenCost reported no unattributed lines'}
          </Typography>
        </Box>
      </Box>

      {hasNumbers && (
        <Section
          title="By pillar"
          hint="Namespaces are grouped by the pillar of the catalog app that declares them."
        >
          <Table size="small" aria-label="Cost by pillar">
            <TableHead>
              <TableRow>
                <TableCell>Pillar</TableCell>
                <TableCell align="right">Namespaces</TableCell>
                <TableCell align="right">Cost</TableCell>
              </TableRow>
            </TableHead>
            <TableBody>
              {pillars.map(item => (
                <TableRow key={item.pillar}>
                  <TableCell>
                    {pillarLabel(item.pillar)}
                    {item.proposed && (
                      <Typography
                        variant="caption"
                        color="textSecondary"
                        display="block"
                      >
                        includes proposed, not approved assignments
                      </Typography>
                    )}
                  </TableCell>
                  <TableCell align="right">{item.namespaces}</TableCell>
                  <TableCell align="right">
                    <Amount cost={item.known} partial={item.unknown > 0} />
                    {item.unknown > 0 && (
                      <Typography
                        variant="caption"
                        color="textSecondary"
                        display="block"
                      >
                        {item.unknown} unknown
                      </Typography>
                    )}
                  </TableCell>
                </TableRow>
              ))}
            </TableBody>
          </Table>
        </Section>
      )}

      {hasNumbers && (
        <Section
          title="By namespace and app"
          hint="One catalog app per namespace is shown by name. Where several apps share a namespace the namespace total is not divided between them."
        >
          <Table size="small" aria-label="Cost by namespace">
            <TableHead>
              <TableRow>
                <TableCell>Namespace</TableCell>
                <TableCell>App</TableCell>
                <TableCell>Pillar</TableCell>
                <TableCell align="right">Cost</TableCell>
                <TableCell align="right">Per hour</TableCell>
                <TableCell>Running</TableCell>
              </TableRow>
            </TableHead>
            <TableBody>
              {lines.map(line => (
                <TableRow key={line.namespace}>
                  <TableCell>{line.namespace}</TableCell>
                  <TableCell>
                    <Attribution line={line} />
                  </TableCell>
                  <TableCell>{pillarCell(line)}</TableCell>
                  <TableCell align="right">
                    <Amount cost={line.cost} />
                    {line.row.state === 'partial' && (
                      <Typography
                        variant="caption"
                        color="textSecondary"
                        display="block"
                      >
                        incomplete data
                      </Typography>
                    )}
                  </TableCell>
                  <TableCell align="right">
                    <Amount cost={line.row.averageHourlyCost} />
                  </TableCell>
                  <TableCell>
                    {activeNote(line, allocation.window.hours)}
                  </TableCell>
                </TableRow>
              ))}
            </TableBody>
          </Table>
        </Section>
      )}
    </>
  );
}

export function CostPage() {
  const [days, setDays] = useState<number>(DEFAULT_DAYS);
  const { data, refresh } = useCostData(days);
  const allocation = data.allocation;
  return (
    <Box p={3}>
      <Typography variant="overline">ONE MIND · INSIGHTS</Typography>
      <Typography variant="h4" component="h1">
        Cost
      </Typography>
      <Typography color="textSecondary">
        Estimated cost of what runs in the cluster, per namespace, app and
        pillar.
      </Typography>

      <Box
        mt={2}
        p={2}
        role="note"
        aria-label="About these numbers"
        style={{ border: '1px solid currentColor', borderRadius: 12 }}
      >
        <Typography variant="body2">
          <strong>These are estimates, not your bill.</strong> OpenCost prices
          each namespace from Kubernetes resource use and list prices. It will
          not match the DigitalOcean invoice exactly: the invoice has rounding,
          credits, unused capacity, snapshots, Spaces and other items that
          OpenCost does not attribute. OpenCost does not report a currency;
          amounts are shown in USD, which is how DigitalOcean bills. Use the
          DigitalOcean billing page for the real total.
        </Typography>
      </Box>

      <Box
        mt={2}
        display="flex"
        alignItems="center"
        flexWrap="wrap"
        style={{ gap: 8 }}
      >
        <Typography component="span" variant="body2">
          Window (complete UTC days):
        </Typography>
        {WINDOW_CHOICES.map(choice => (
          <Button
            key={choice}
            size="small"
            variant={choice === days ? 'contained' : 'outlined'}
            color="primary"
            aria-pressed={choice === days}
            onClick={() => setDays(choice)}
          >
            {choice === 1 ? '1 day' : `${choice} days`}
          </Button>
        ))}
        <Button size="small" onClick={refresh}>
          Refresh
        </Button>
        {allocation && (
          <Typography variant="caption" color="textSecondary">
            {allocation.window.start} to {allocation.window.end}
            {data.observedAt ? ` · fetched ${data.observedAt}` : ''}
          </Typography>
        )}
      </Box>

      {data.status === 'loading' && !allocation && (
        <Typography role="status" style={{ marginTop: 16 }}>
          Loading cost from OpenCost…
        </Typography>
      )}
      {data.status === 'loading' && allocation && (
        <Typography role="status" variant="caption" display="block">
          Refreshing…
        </Typography>
      )}
      <CostView data={data} />
    </Box>
  );
}
