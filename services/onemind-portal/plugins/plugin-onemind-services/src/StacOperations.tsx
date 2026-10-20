import { useEffect, useState } from 'react';
import {
  discoveryApiRef,
  fetchApiRef,
  useApi,
} from '@backstage/core-plugin-api';
import { Box, Button, Link, Typography } from '@material-ui/core';

type Metric = {
  value: number | null;
  sampledAt: string | null;
  state: string;
} | null;
type Workload = {
  role: string;
  entityRef: string;
  state: string;
  observedAt: string;
  resource: string | null;
  desired: number | null;
  ready: number | null;
  images: { name: string; image: string }[];
  renderedRelease?: string | null;
  cpu: Metric;
  memory: Metric;
};
type Evidence = {
  schemaVersion: string;
  observedAt: string;
  scope: string;
  partial: boolean;
  workloads: Workload[];
  probes: Record<
    string,
    { state: string; httpStatus: number | null; observedAt: string }
  >;
  data: {
    state: string;
    observedAt: string;
    collectionsReturned: number | null;
    collectionsTruncated: boolean | null;
    sampleItemsReturned: number | null;
    latestItemDatetime: string | null;
    freshness: string;
  };
  argo: {
    state: string;
    observedAt: string;
    applications: {
      name: string;
      health: string;
      sync: string;
      revision: string | null;
      reconciledAt: string | null;
    }[];
  };
  recovery: { state: string; detail: string };
};
export function isStacEntity(entity: {
  kind: string;
  metadata: {
    name: string;
    namespace?: string;
    labels?: Record<string, string>;
  };
}) {
  return (
    entity.metadata.labels?.['openchoreo.io/managed'] === 'true' &&
    (entity.metadata.namespace ?? 'default') === 'default' &&
    ((entity.kind.toLowerCase() === 'component' &&
      entity.metadata.name === 'stac-browser') ||
      (entity.kind.toLowerCase() === 'system' &&
        entity.metadata.name === 'onetak'))
  );
}
const time = (v: string | null) =>
  v && Number.isFinite(Date.parse(v))
    ? new Date(v).toLocaleString()
    : 'Unknown';
const entityPath = (ref: string) => {
  const [kind, rest] = ref.split(':');
  const [namespace, name] = rest.split('/');
  return `/catalog/${namespace}/${kind}/${name}`;
};
export function StacOperations() {
  const discovery = useApi(discoveryApiRef);
  const api = useApi(fetchApiRef);
  const [data, setData] = useState<Evidence>();
  const [error, setError] = useState(false);
  const [refresh, setRefresh] = useState(0);
  useEffect(() => {
    let active = true;
    let pending = false;
    const controller = new AbortController();
    const load = async () => {
      if (pending) return;
      pending = true;
      try {
        const base = await discovery.getBaseUrl('onemind-services');
        const response = await api.fetch(`${base}/stac`, {
          signal: controller.signal,
        });
        if (!response.ok) throw new Error('Observation unavailable');
        const value = await response.json();
        if (
          value.schemaVersion !== 'onemind.stac-evidence.v1' ||
          !Array.isArray(value.workloads)
        )
          throw new Error('Invalid observation contract');
        if (active) {
          setData(value);
          setError(false);
        }
      } catch {
        if (active) setError(true);
      } finally {
        pending = false;
      }
    };
    load();
    const timer = setInterval(load, 60_000);
    return () => {
      active = false;
      controller.abort();
      clearInterval(timer);
    };
  }, [discovery, api, refresh]);
  const stale =
    error ||
    (!!data &&
      (!Number.isFinite(Date.parse(data.observedAt)) ||
        Date.now() - Date.parse(data.observedAt) > 180_000));
  return (
    <Box p={3}>
      <Typography variant="h4">OneTAK · STAC operational picture</Typography>
      <Typography>
        Existing Development pilot only. STAC API and its database remain
        Argo-managed in onetak; this page does not migrate or control them.
      </Typography>
      <Box my={2}>
        <Button variant="outlined" onClick={() => setRefresh(x => x + 1)}>
          Refresh observations
        </Button>
      </Box>
      {error && (
        <Typography color="error">
          Observation unavailable. Previous results, if shown, are stale—not
          proof of current health.
        </Typography>
      )}
      {!data && !error && <Typography>Loading service evidence…</Typography>}
      {data && (
        <>
          <Typography color={stale ? 'error' : 'textSecondary'}>
            {stale ? 'Stale observation' : 'Observed'} · {time(data.observedAt)}{' '}
            · refreshes every minute; server cache up to 60 seconds
          </Typography>
          <Typography>{data.scope}</Typography>
          {data.partial && (
            <Typography color="error">
              Incomplete or failing checks. A ready deployment or HTTP-200
              landing page does not prove catalog searches work.
            </Typography>
          )}
          <Box my={3}>
            <Typography variant="h5">
              Identity, dependencies and source
            </Typography>
            <Typography>
              <Link href="/catalog/default/component/stac-browser">
                STAC Browser
              </Link>{' '}
              →{' '}
              <Link href="/catalog/default/component/k8s-onetak-stac-api-746ca5f2">
                STAC API
              </Link>{' '}
              →{' '}
              <Link href="/catalog/default/component/k8s-onetak-onetak-geo-db-a22783b9">
                OneTAK geo database
              </Link>
            </Typography>
            <Typography>
              <Link href="/catalog/default/api/onetak-stac-readonly">
                Read-only API contract
              </Link>{' '}
              ·{' '}
              <Link href="/catalog/default/component/stac-browser/environments">
                Deployment and active release
              </Link>{' '}
              ·{' '}
              <Link href="/catalog/default/component/stac-browser/definition">
                OpenChoreo definition
              </Link>
            </Typography>
            <Typography>
              <Link href="https://github.com/onemindos/omos/tree/main/platform/openchoreo/projects/onetak/stac-browser">
                Reviewed pilot source
              </Link>{' '}
              ·{' '}
              <Link href="https://github.com/onemindos/omos/blob/main/platform/onetak/stac.yaml">
                Backing API source
              </Link>
            </Typography>
          </Box>
          <Typography variant="h5">
            Observed workloads and resource use
          </Typography>
          <Typography variant="caption">
            Kubernetes readiness is separate from API behavior. Metrics cover
            each selected Deployment's pods, not the entire namespace.
            Prometheus query times are not proof of ingestion or scrape
            freshness.
          </Typography>
          {data.workloads.map(w => (
            <Box
              key={w.role}
              my={2}
              p={2}
              border={1}
              borderColor="divider"
              borderRadius={8}
              style={{ overflowWrap: 'anywhere' }}
            >
              <Typography variant="h6">
                {w.role} · {w.state}
              </Typography>
              <Typography>
                {w.resource ?? 'No unique workload observed'} · ready{' '}
                {w.ready ?? '?'}/{w.desired ?? '?'} · {time(w.observedAt)}
              </Typography>
              {w.renderedRelease && (
                <Typography>
                  RenderedRelease: {w.renderedRelease} (not the ComponentRelease
                  version)
                </Typography>
              )}
              {w.images.map(i => (
                <Typography variant="body2" key={i.name}>
                  {i.name}: {i.image}
                </Typography>
              ))}
              <Typography>
                CPU:{' '}
                {w.cpu?.value === null || w.cpu?.value === undefined
                  ? 'Unknown'
                  : `${(w.cpu.value * 1000).toFixed(2)} mCPU`}{' '}
                · memory:{' '}
                {w.memory?.value === null || w.memory?.value === undefined
                  ? 'Unknown'
                  : `${(w.memory.value / 1048576).toFixed(1)} MiB`}
              </Typography>
              <Typography variant="caption">
                CPU query: {time(w.cpu?.sampledAt ?? null)} · memory query:{' '}
                {time(w.memory?.sampledAt ?? null)}
              </Typography>
              <Typography>
                <Link href={entityPath(w.entityRef)}>Catalog record</Link>
              </Typography>
            </Box>
          ))}
          <Box my={3}>
            <Typography variant="h5">Real application checks</Typography>
            {Object.entries(data.probes).map(([name, p]) => (
              <Typography key={name}>
                {name}: {p.state} · HTTP {p.httpStatus ?? 'unknown'} ·{' '}
                {time(p.observedAt)}
              </Typography>
            ))}
            <Typography variant="caption">
              Fixed cluster-side read-only GET requests. These do not test every
              client, login path or STAC operation.
            </Typography>
          </Box>
          <Box my={3}>
            <Typography variant="h5">Catalog data and freshness</Typography>
            <Typography>
              Query state: {data.data.state} · {time(data.data.observedAt)}
            </Typography>
            <Typography>
              Collections returned: {data.data.collectionsReturned ?? 'Unknown'}
              {data.data.collectionsTruncated ? ' (partial page)' : ''} · sample
              items returned: {data.data.sampleItemsReturned ?? 'Unknown'}
            </Typography>
            <Typography>
              Newest queried item acquisition time:{' '}
              {time(data.data.latestItemDatetime)}
            </Typography>
            <Typography>{data.data.freshness}</Typography>
          </Box>
          <Box my={3}>
            <Typography variant="h5">
              GitOps configuration—not exclusive workload health
            </Typography>
            <Typography>
              Source state: {data.argo.state} · {time(data.argo.observedAt)}
            </Typography>
            {data.argo.applications.map(a => (
              <Box key={a.name} my={1} style={{ overflowWrap: 'anywhere' }}>
                <Typography>
                  <Link
                    href={`https://argocd.tail717ba0.ts.net/applications/argocd/${encodeURIComponent(
                      a.name,
                    )}`}
                  >
                    {a.name}
                  </Link>
                  : {a.health} / {a.sync}
                </Typography>
                <Typography variant="caption">
                  Git revision: {a.revision ?? 'Unknown'} · Argo reconciled:{' '}
                  {time(a.reconciledAt)}
                </Typography>
              </Box>
            ))}
          </Box>
          <Box my={3}>
            <Typography variant="h5">Investigate and recover</Typography>
            <Typography>
              <Link href="/catalog/default/component/stac-browser/runtime-logs">
                Live pilot logs
              </Link>{' '}
              ·{' '}
              <Link href="/catalog/default/component/stac-browser/metrics">
                Native OpenChoreo metrics
              </Link>{' '}
              ·{' '}
              <Link href="/docs/default/component/onemind-runbooks/stac/">
                STAC diagnosis and recovery runbook
              </Link>
            </Typography>
            <Typography>{data.recovery.detail}</Typography>
            <Typography>
              No restart, migration, schema initialization, write API or
              production promotion is performed here.
            </Typography>
          </Box>
        </>
      )}
    </Box>
  );
}
