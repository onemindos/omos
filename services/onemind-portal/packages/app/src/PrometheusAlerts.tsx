import { useEffect, useState } from 'react';
import { useApi } from '@backstage/core-plugin-api';
import { useEntity } from '@backstage/plugin-catalog-react';
import { PrometheusAlertStatus } from '@roadiehq/backstage-plugin-prometheus';
import { Box, Button, Link, Typography } from '@material-ui/core';
import {
  scopeForEntity,
  scopedPrometheusApiRef,
  Observation,
} from './prometheusApi';

export default function PrometheusAlerts() {
  const { entity } = useEntity();
  const scope = scopeForEntity(entity);
  const api = useApi(scopedPrometheusApiRef);
  const [value, setValue] = useState<Observation>();
  const [error, setError] = useState(false);
  const [refresh, setRefresh] = useState(0);
  useEffect(() => {
    let alive = true;
    setValue(undefined);
    setError(false);
    const load = async () => {
      try {
        const next = await api.getAlerts({ serviceName: scope });
        if (alive) {
          setValue(next);
          setError(false);
        }
      } catch {
        if (alive) {
          setValue(undefined);
          setError(true);
        }
      }
    };
    void load();
    const interval = setInterval(() => {
      void load();
    }, 30000);
    return () => {
      alive = false;
      clearInterval(interval);
    };
  }, [api, scope, refresh]);
  if (!scope) return null;
  const count = value?.data.groups.reduce(
    (n, g) => n + g.rules.reduce((total, r) => total + r.alerts.length, 0),
    0,
  );
  return (
    <Box p={2}>
      <Typography variant="h5" component="h3">
        Prometheus service alerts
      </Typography>
      <Typography variant="body2">
        Scope: {scope}. Kubernetes workload labels and reviewed service rules
        only; cluster-wide alerts and complete application correctness are not
        covered.
      </Typography>
      <Box my={1}>
        <Link
          href={`/catalog/default/component/${entity.metadata.name}/kubernetes`}
        >
          Live Kubernetes resources
        </Link>
        {' · '}
        <Link
          href={`https://grafana.tail717ba0.ts.net/d/${
            scope === 'postgres'
              ? 'onemind-backstage-postgres'
              : `onemind-${scope}`
          }`}
          target="_blank"
          rel="noopener noreferrer"
        >
          Grafana workload dashboard
        </Link>
        {' · '}
        <Link
          href={`/docs/default/component/onemind-runbooks/${
            scope === 'nats' ? 'nats' : 'resilience'
          }/`}
        >
          Recovery runbook
        </Link>
      </Box>
      <Button size="small" onClick={() => setRefresh(n => n + 1)}>
        Refresh alerts
      </Button>
      {error && (
        <Typography role="alert" color="error">
          Unknown: Prometheus alerts are unavailable. No all-clear can be
          inferred.
        </Typography>
      )}
      {!error && !value && (
        <Typography role="status">Loading Prometheus observation…</Typography>
      )}
      {value && (
        <>
          <Typography variant="caption" component="p">
            Observed {value.observedAt}; refresh every 30 seconds. Relevant rule
            families: {value.candidateRules} (not a complete coverage
            guarantee).
          </Typography>
          {value.coverage !== 'observed' && (
            <Typography role="alert" color="error">
              Coverage {value.coverage}: {value.unhealthyRules} rule evaluation
              errors; {value.staleRules} stale evaluations. Monitoring is not
              fully verified.
            </Typography>
          )}
          {count === 0 && (
            <Typography>
              No firing or pending alerts matched this scope. This is not proof
              that the service is healthy.
            </Typography>
          )}
          <PrometheusAlertStatus
            key={`${scope}-${value.observedAt}`}
            alerts="all"
            showInactiveAlerts={false}
            showAnnotations={false}
            showLabels={false}
            extraColumns={[
              {
                title: 'Severity',
                render: row => row.labels.severity ?? 'unknown',
              },
              {
                title: 'Affected target',
                render: row =>
                  row.labels.pod ??
                  row.labels.deployment ??
                  row.labels.statefulset ??
                  row.labels.service ??
                  'Unknown target',
              },
              {
                title: 'Summary',
                render: row => row.annotations.summary ?? 'No summary supplied',
              },
            ]}
          />
        </>
      )}
    </Box>
  );
}
