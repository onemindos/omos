import { useEffect, useMemo, useState } from 'react';
import {
  discoveryApiRef,
  fetchApiRef,
  useApi,
} from '@backstage/core-plugin-api';
import { useEntity } from '@backstage/plugin-catalog-react';
import { Box, Button, Typography } from '@material-ui/core';
import {
  countsFor,
  matchPolicy,
  PolicyResponse,
  targetsForEntity,
} from './policyMatch';

type State =
  | { kind: 'loading' }
  | { kind: 'unavailable' }
  | {
      kind: 'data';
      response: Extract<PolicyResponse, { status: 'ok' | 'stale' }>;
    };

/** Reads Kyverno PolicyReports through the portal backend. A failed read is "unavailable", never zero. */
export function usePolicyReports(refresh: number): State {
  const discovery = useApi(discoveryApiRef);
  const fetchApi = useApi(fetchApiRef);
  const [state, setState] = useState<State>({ kind: 'loading' });
  useEffect(() => {
    let alive = true;
    (async () => {
      try {
        const base = await discovery.getBaseUrl('onemind');
        const response = await fetchApi.fetch(`${base}/policy-reports`);
        if (!response.ok) throw new Error('unavailable');
        const body = (await response.json()) as PolicyResponse;
        if (!alive) return;
        setState(
          body.status === 'ok' || body.status === 'stale'
            ? { kind: 'data', response: body }
            : { kind: 'unavailable' },
        );
      } catch {
        if (alive) setState({ kind: 'unavailable' });
      }
    })();
    return () => {
      alive = false;
    };
  }, [discovery, fetchApi, refresh]);
  return state;
}

export default function PolicyResultsCard() {
  const { entity } = useEntity();
  const targets = useMemo(() => targetsForEntity(entity), [entity]);
  const [refresh, setRefresh] = useState(0);
  const state = usePolicyReports(refresh);

  return (
    <Box p={2}>
      <Typography variant="h5" component="h3">
        Policy results (Kyverno)
      </Typography>
      <Typography variant="body2" color="textSecondary">
        What Kyverno's reports say about this app's workloads. They come from
        Kyverno's own scans of resources already running; they are not a full
        security assessment.
      </Typography>
      {state.kind === 'loading' && (
        <Typography role="status">Loading policy results…</Typography>
      )}
      {state.kind === 'unavailable' && (
        <Typography role="alert" color="error">
          Unknown: policy results are unavailable. No all-clear can be inferred.
        </Typography>
      )}
      {state.kind === 'data' && <Results state={state} targets={targets} />}
      <Box mt={1}>
        <Button size="small" onClick={() => setRefresh(n => n + 1)}>
          Refresh policy results
        </Button>
      </Box>
    </Box>
  );
}

function Results({
  state,
  targets,
}: {
  state: Extract<State, { kind: 'data' }>;
  targets: ReturnType<typeof targetsForEntity>;
}) {
  const { response } = state;
  const matched = useMemo(
    () => matchPolicy(response, targets),
    [response, targets],
  );
  const counts = countsFor(matched);
  const atLeast = response.truncated ? 'at least ' : '';
  const shown = matched.findings.slice(0, 10);
  return (
    <>
      <Typography variant="caption" component="p">
        Observed {response.observedAt}; scope:{' '}
        {matched.scope === 'workloads'
          ? "this app's named workloads"
          : 'the whole namespace (the catalog names no workloads for this app)'}
        .
      </Typography>
      {response.status === 'stale' && (
        <Typography role="alert" color="error">
          The latest read failed ({response.error ?? 'unavailable'}); showing
          the last good data from {response.observedAt}.
        </Typography>
      )}
      {!targets.length && (
        <Typography>
          This entity has no Kubernetes namespace, so there is nothing to match.
        </Typography>
      )}
      {targets.length > 0 && counts.resources === 0 && (
        <Typography>
          Kyverno has no reports for{' '}
          {matched.scope === 'workloads' ? 'these workloads' : 'this namespace'}{' '}
          yet. That is not the same as passing.
        </Typography>
      )}
      {counts.resources > 0 && (
        <>
          <Typography variant="h6" component="p">
            {atLeast}
            {counts.violations} policy violation
            {counts.violations === 1 && !atLeast ? '' : 's'}
            {counts.warnings > 0 &&
              `, ${atLeast}${counts.warnings} warning${
                counts.warnings === 1 ? '' : 's'
              }`}
          </Typography>
          <Typography variant="body2">
            {counts.passed} checks passed, {counts.skipped} skipped, across{' '}
            {counts.resources} scanned resource
            {counts.resources === 1 ? '' : 's'} (Deployments, ReplicaSets, Pods
            and similar; one problem can show on several of them).
          </Typography>
        </>
      )}
      {matched.unreported.length > 0 && (
        <Typography variant="body2" role="status">
          No report yet for: {matched.unreported.join(', ')}.
        </Typography>
      )}
      {shown.length > 0 && (
        <Box component="ul" pl={2} mt={1} aria-label="Policy findings">
          {shown.map(finding => (
            <li
              key={`${finding.resource.kind}/${finding.resource.name}/${finding.policy}/${finding.rule}`}
            >
              <strong>{finding.result}</strong> · {finding.policy} /{' '}
              {finding.rule} · {finding.resource.kind} {finding.resource.name}
              {finding.severity ? ` · ${finding.severity}` : ''}
              {finding.message ? ` — ${finding.message}` : ''}
            </li>
          ))}
        </Box>
      )}
      {matched.findings.length > shown.length && (
        <Typography variant="caption">
          Showing {shown.length} of {matched.findings.length} findings.
        </Typography>
      )}
      {response.truncated && (
        <Typography variant="caption" component="p">
          Kyverno reported more findings than the portal lists; counts are lower
          bounds.
        </Typography>
      )}
    </>
  );
}
