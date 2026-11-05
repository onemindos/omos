import { useCallback, useEffect, useState } from 'react';
import {
  discoveryApiRef,
  fetchApiRef,
  useApi,
} from '@backstage/core-plugin-api';
import { Link } from '@backstage/core-components';
import {
  Box,
  Button,
  Chip,
  Table,
  TableBody,
  TableCell,
  TableHead,
  TableRow,
  Typography,
} from '@material-ui/core';

export type Action = { kind: string; label: string; href: string };
export type Finding = {
  check: string;
  severity: 'info' | 'warn';
  subject: string;
  detail: string;
  actions: Action[];
};
export type Check = {
  id: string;
  title: string;
  status: 'ok' | 'findings' | 'unavailable';
  count: number;
  reason?: string;
};
export type Report = {
  generatedAt: string;
  trigger: 'schedule' | 'manual';
  clusterObservedAt: string | null;
  checks: Check[];
  findings: Finding[];
};

export function CheckStatus({ check }: { check: Check }) {
  if (check.status === 'unavailable')
    return <Chip size="small" label="Unavailable" color="default" />;
  if (check.status === 'ok')
    return <Chip size="small" label="No findings" color="primary" />;
  return (
    <Chip size="small" label={`${check.count} to review`} color="secondary" />
  );
}

export function HygieneReportView({ report }: { report: Report }) {
  const titles = new Map(report.checks.map(c => [c.id, c.title]));
  return (
    <>
      <Typography variant="body2" color="textSecondary">
        Generated {new Date(report.generatedAt).toLocaleString()} (
        {report.trigger === 'schedule' ? 'weekly run' : 'on demand'}); cluster
        observed{' '}
        {report.clusterObservedAt
          ? new Date(report.clusterObservedAt).toLocaleString()
          : 'never (cluster unreadable)'}
        . Read-only: this page changes nothing.
      </Typography>
      <Box mt={2}>
        <Table size="small" aria-label="Checks">
          <TableHead>
            <TableRow>
              <TableCell>Check</TableCell>
              <TableCell>Result</TableCell>
              <TableCell>Note</TableCell>
            </TableRow>
          </TableHead>
          <TableBody>
            {report.checks.map(check => (
              <TableRow key={check.id}>
                <TableCell>{check.title}</TableCell>
                <TableCell>
                  <CheckStatus check={check} />
                </TableCell>
                <TableCell>
                  {check.status === 'unavailable'
                    ? `Not evaluated: ${check.reason}`
                    : ''}
                </TableCell>
              </TableRow>
            ))}
          </TableBody>
        </Table>
      </Box>
      <Box mt={3}>
        <Typography variant="h6">Findings ({report.findings.length})</Typography>
        {report.findings.length === 0 ? (
          <Typography>
            No findings in the checks that could be evaluated.
          </Typography>
        ) : (
          <Table size="small" aria-label="Findings">
            <TableHead>
              <TableRow>
                <TableCell>Check</TableCell>
                <TableCell>Subject</TableCell>
                <TableCell>Detail</TableCell>
                <TableCell>Next step</TableCell>
              </TableRow>
            </TableHead>
            <TableBody>
              {report.findings.map(f => (
                <TableRow key={`${f.check}:${f.subject}`}>
                  <TableCell>{titles.get(f.check) ?? f.check}</TableCell>
                  <TableCell>
                    <code>{f.subject}</code>
                  </TableCell>
                  <TableCell>{f.detail}</TableCell>
                  <TableCell>
                    {f.actions.map(a => (
                      <Box key={a.href} component="span" mr={1}>
                        <Link to={a.href}>{a.label}</Link>
                      </Box>
                    ))}
                  </TableCell>
                </TableRow>
              ))}
            </TableBody>
          </Table>
        )}
      </Box>
    </>
  );
}

export function StackHygienePage() {
  const discovery = useApi(discoveryApiRef);
  const fetchApi = useApi(fetchApiRef);
  const [report, setReport] = useState<Report | null | undefined>();
  const [message, setMessage] = useState<string>();
  const [busy, setBusy] = useState(false);

  const call = useCallback(
    async (path: string, init?: RequestInit) => {
      const base = await discovery.getBaseUrl('onemind-lifecycle');
      return fetchApi.fetch(`${base}${path}`, init);
    },
    [discovery, fetchApi],
  );

  useEffect(() => {
    let active = true;
    call('/stack-hygiene')
      .then(async response => {
        if (!response.ok) throw new Error(String(response.status));
        const body = await response.json();
        if (active) setReport(body.report);
      })
      .catch(() => {
        if (active) setMessage('The hygiene report could not be loaded.');
      });
    return () => {
      active = false;
    };
  }, [call]);

  const refresh = async () => {
    setBusy(true);
    setMessage(undefined);
    try {
      const response = await call('/stack-hygiene/refresh', { method: 'POST' });
      if (response.status === 403) {
        setMessage('Only administrators can refresh the report on demand.');
      } else if (!response.ok) {
        setMessage('The refresh failed. The previous report is still shown.');
      } else {
        setReport((await response.json()).report);
      }
    } catch {
      setMessage('The refresh failed. The previous report is still shown.');
    } finally {
      setBusy(false);
    }
  };

  return (
    <Box p={3}>
      <Box display="flex" alignItems="center" justifyContent="space-between">
        <Typography variant="h4" component="h1">
          Stack hygiene
        </Typography>
        <Button
          variant="outlined"
          onClick={refresh}
          disabled={busy}
          aria-label="Refresh report"
        >
          {busy ? 'Refreshing' : 'Refresh now'}
        </Button>
      </Box>
      <Typography paragraph color="textSecondary">
        What to add and what to remove. Computed weekly from the catalog, Argo
        CD, the cluster and Git, all read-only.
      </Typography>
      {message && (
        <Typography role="status" color="error">
          {message}
        </Typography>
      )}
      {report === undefined && !message && <Typography>Loading</Typography>}
      {report === null && (
        <Typography>
          No report yet. The weekly run has not happened; an administrator can
          refresh now.
        </Typography>
      )}
      {report && <HygieneReportView report={report} />}
    </Box>
  );
}
