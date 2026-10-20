import http from 'node:http';
import { list } from './clusterRead';

type RecordValue = Record<string, unknown>;
const object = (v: unknown): RecordValue =>
  v && typeof v === 'object' && !Array.isArray(v) ? (v as RecordValue) : {};
const stamp = () => new Date().toISOString();
const selector =
  'openchoreo.dev/namespace=default,openchoreo.dev/project=onetak,openchoreo.dev/component=stac-browser,openchoreo.dev/environment=development';
const api = 'http://stac-api.onetak.svc.cluster.local:8082';
const gatewayBrowser =
  'http://cilium-gateway-gateway-default.openchoreo-data-plane.svc.cluster.local/stac-browser-http/';
const gatewayHost = 'development-default.apps.openchoreo.100-110-185-30.nip.io';
const prometheus =
  'http://monitoring-kube-prometheus-prometheus.monitoring.svc:9090';
export const stacRefs = {
  browser: 'component:default/stac-browser',
  api: 'component:default/k8s-onetak-stac-api-746ca5f2',
  database: 'component:default/k8s-onetak-onetak-geo-db-a22783b9',
  contract: 'api:default/onetak-stac-readonly',
};

type Deployment = {
  metadata: {
    name: string;
    namespace: string;
    generation?: number;
    labels?: Record<string, string>;
  };
  spec?: {
    replicas?: number;
    template?: { spec?: { containers?: { name: string; image: string }[] } };
  };
  status?: { readyReplicas?: number; observedGeneration?: number };
};
type Application = {
  metadata: { name: string };
  status?: {
    health?: { status?: string };
    sync?: { status?: string; revision?: string };
    reconciledAt?: string;
  };
};
export type HttpObservation = { status: number; body: string };
export type EvidenceDependencies = {
  list: <T>(path: string) => Promise<T[]>;
  get: (url: string) => Promise<HttpObservation>;
};

// Internal, fixed targets only. No client URL, credential forwarding or redirects.
export function getEvidenceHttp(url: string): Promise<HttpObservation> {
  return new Promise((resolve, reject) => {
    const req = http.get(
      url,
      { headers: url === gatewayBrowser ? { Host: gatewayHost } : {} },
      response => {
        const chunks: Buffer[] = [];
        let bytes = 0;
        response.on('data', (chunk: Buffer) => {
          bytes += chunk.length;
          if (bytes > 2_000_000) {
            req.destroy(new Error('Observation response limit'));
            return;
          }
          chunks.push(chunk);
        });
        response.on('error', reject);
        response.on('end', () =>
          resolve({
            status: response.statusCode ?? 0,
            body: Buffer.concat(chunks).toString('utf8'),
          }),
        );
      },
    );
    const deadline = setTimeout(
      () => req.destroy(new Error('Observation deadline')),
      5000,
    );
    req.on('close', () => clearTimeout(deadline));
    req.on('error', reject);
  });
}
async function probe(get: EvidenceDependencies['get'], url: string) {
  try {
    const r = await get(url);
    let state = 'error';
    if (r.status === 200) state = 'responding';
    else if (r.status === 401 || r.status === 403)
      state = 'authentication-required';
    else if (r.status >= 300 && r.status < 400) state = 'redirect-not-followed';
    return { state, httpStatus: r.status, observedAt: stamp() };
  } catch {
    return { state: 'unknown', httpStatus: null, observedAt: stamp() };
  }
}
async function json(get: EvidenceDependencies['get'], url: string) {
  const r = await get(url);
  if (r.status !== 200) throw new Error('Source not successful');
  const value = JSON.parse(r.body);
  if (!value || typeof value !== 'object')
    throw new Error('Invalid JSON source');
  return object(value);
}
async function metric(
  get: EvidenceDependencies['get'],
  namespace: string,
  deployment: string,
  memory: boolean,
) {
  const labels = `namespace="${namespace}",pod=~"${deployment}-.*",container!="",container!="POD"`;
  const query = memory
    ? `sum(container_memory_working_set_bytes{${labels}})`
    : `sum(rate(container_cpu_usage_seconds_total{${labels}}[5m]))`;
  try {
    const x = await json(
      get,
      `${prometheus}/api/v1/query?query=${encodeURIComponent(query)}`,
    );
    const data = object(x.data);
    const rows = data.result;
    if (
      x.status !== 'success' ||
      data.resultType !== 'vector' ||
      !Array.isArray(rows) ||
      rows.length !== 1
    )
      throw new Error('No unique metric');
    const pair = object(rows[0]).value;
    if (
      !Array.isArray(pair) ||
      pair.length !== 2 ||
      typeof pair[0] !== 'number' ||
      typeof pair[1] !== 'string' ||
      pair[1].trim().length === 0 ||
      !Number.isFinite(Number(pair[0])) ||
      !Number.isFinite(Number(pair[1])) ||
      Number(pair[1]) < 0
    )
      throw new Error('Invalid metric');
    const time = Number(pair[0]) * 1000;
    if (time > Date.now() + 60_000 || Date.now() - time > 180_000)
      return {
        value: null,
        sampledAt: new Date(time).toISOString(),
        state: 'stale',
      };
    return {
      value: Number(pair[1]),
      sampledAt: new Date(time).toISOString(),
      state: 'observed',
    };
  } catch {
    return { value: null, sampledAt: null, state: 'unknown' };
  }
}
async function workload(
  d: Deployment | undefined,
  role: keyof typeof stacRefs,
  deps: EvidenceDependencies,
) {
  if (
    !d ||
    !/^[a-z0-9]([a-z0-9-]*[a-z0-9])?$/.test(d.metadata.name) ||
    !/^[a-z0-9]([a-z0-9-]*[a-z0-9])?$/.test(d.metadata.namespace)
  ) {
    return {
      role,
      entityRef: stacRefs[role],
      state: 'unknown',
      observedAt: stamp(),
      resource: null,
      desired: null,
      ready: null,
      images: [],
      cpu: null,
      memory: null,
    };
  }
  const desired = d.spec?.replicas ?? 1;
  const ready = d.status?.readyReplicas ?? 0;
  const current =
    Number.isInteger(d.metadata.generation) &&
    (d.status?.observedGeneration ?? -1) >= d.metadata.generation!;
  const [cpu, memory] = await Promise.all([
    metric(deps.get, d.metadata.namespace, d.metadata.name, false),
    metric(deps.get, d.metadata.namespace, d.metadata.name, true),
  ]);
  let state = 'not-ready';
  if (!current) state = 'reconciling';
  else if (desired === 0) state = 'scaled-to-zero';
  else if (ready >= desired) state = 'ready';
  return {
    role,
    entityRef: stacRefs[role],
    state,
    observedAt: stamp(),
    resource: `${d.metadata.namespace}/${d.metadata.name}`,
    desired,
    ready,
    images: (d.spec?.template?.spec?.containers ?? []).map(c => ({
      name: c.name,
      image: c.image,
    })),
    renderedRelease:
      d.metadata.labels?.['openchoreo.dev/rendered-release-name'] ?? null,
    cpu,
    memory,
  };
}
async function catalog(get: EvidenceDependencies['get']) {
  const observedAt = stamp();
  try {
    const [collections, search] = await Promise.all([
      json(get, `${api}/collections`),
      json(get, `${api}/search?limit=1&sortby=-datetime`),
    ]);
    if (
      !Array.isArray(collections.collections) ||
      search.type !== 'FeatureCollection' ||
      !Array.isArray(search.features)
    )
      throw new Error('Invalid STAC response');
    const item = object(search.features[0]);
    const properties = object(item.properties);
    const value = properties.datetime ?? properties.end_datetime;
    const latestItemDatetime =
      typeof value === 'string' &&
      Number.isFinite(Date.parse(value)) &&
      Date.parse(value) <= Date.now() + 60_000
        ? value
        : null;
    const truncated =
      Array.isArray(collections.links) &&
      collections.links.some(l => object(l).rel === 'next');
    return {
      state: 'observed',
      observedAt,
      collectionsReturned: collections.collections.length,
      collectionsTruncated: truncated,
      sampleItemsReturned: search.features.length,
      latestItemDatetime,
      freshness:
        'No ingestion SLA configured. Item acquisition time is not ingestion time; this is a bounded sorted query, not a full catalog audit.',
    };
  } catch {
    return {
      state: 'unknown',
      observedAt,
      collectionsReturned: null,
      collectionsTruncated: null,
      sampleItemsReturned: null,
      latestItemDatetime: null,
      freshness:
        'Collection/search query unavailable or invalid. No freshness or empty-catalog conclusion can be drawn.',
    };
  }
}

export async function observeStacEvidence(
  dependencies: EvidenceDependencies = { list, get: getEvidenceHttp },
) {
  const requests = new Map<string, Promise<HttpObservation>>();
  const deps: EvidenceDependencies = {
    ...dependencies,
    get: url => {
      let request = requests.get(url);
      if (!request) {
        request = dependencies.get(url);
        requests.set(url, request);
      }
      return request;
    },
  };
  const observedAt = stamp();
  const sources = await Promise.allSettled([
    deps.list<Deployment>(
      `/apis/apps/v1/deployments?labelSelector=${encodeURIComponent(selector)}`,
    ),
    deps.list<Deployment>('/apis/apps/v1/namespaces/onetak/deployments'),
    deps.list<Application>(
      '/apis/argoproj.io/v1alpha1/namespaces/argocd/applications',
    ),
  ]);
  const deployments = (i: number) =>
    sources[i].status === 'fulfilled'
      ? (sources[i] as PromiseFulfilledResult<Deployment[]>).value
      : [];
  const browsers = deployments(0).filter(
    d =>
      d.metadata.labels?.['openchoreo.dev/component'] === 'stac-browser' &&
      d.metadata.labels?.['openchoreo.dev/project'] === 'onetak' &&
      d.metadata.labels?.['openchoreo.dev/environment'] === 'development' &&
      d.metadata.labels?.['openchoreo.dev/namespace'] === 'default',
  );
  const browser = browsers.length === 1 ? browsers[0] : undefined;
  const onetak = deployments(1);
  const choose = (name: string) => {
    const matches = onetak.filter(
      d => d.metadata.namespace === 'onetak' && d.metadata.name === name,
    );
    return matches.length === 1 ? matches[0] : undefined;
  };
  const browserUrl =
    browser &&
    /^dp-default-onetak-development-[a-f0-9]+$/.test(browser.metadata.namespace)
      ? gatewayBrowser
      : undefined;
  const [workloads, browserHttp, apiHttp, collectionHttp, searchHttp, data] =
    await Promise.all([
      Promise.all([
        workload(browser, 'browser', deps),
        workload(choose('stac-api'), 'api', deps),
        workload(choose('onetak-geo-db'), 'database', deps),
      ]),
      browserUrl
        ? probe(deps.get, browserUrl)
        : Promise.resolve({
            state: 'unknown',
            httpStatus: null,
            observedAt: stamp(),
          }),
      probe(deps.get, `${api}/`),
      probe(deps.get, `${api}/collections`),
      probe(deps.get, `${api}/search?limit=1&sortby=-datetime`),
      catalog(deps.get),
    ]);
  const argoResult = sources[2];
  const names = ['platform-openchoreo-projects', 'platform-onetak'];
  const apps =
    argoResult.status === 'fulfilled'
      ? (argoResult.value as Application[])
          .filter(a => names.includes(a.metadata.name))
          .map(a => ({
            name: a.metadata.name,
            health: a.status?.health?.status ?? 'Unknown',
            sync: a.status?.sync?.status ?? 'Unknown',
            revision: a.status?.sync?.revision ?? null,
            reconciledAt: a.status?.reconciledAt ?? null,
          }))
      : [];
  return {
    schemaVersion: 'onemind.stac-evidence.v1',
    serviceRef: stacRefs.browser,
    systemRef: 'system:default/onetak',
    observedAt,
    scope:
      'Read-only development pilot observations. Argo application state is shared configuration scope, not exclusive app health. No mutation capability is exposed.',
    partial:
      sources.some(r => r.status === 'rejected') ||
      workloads.some(w => w.state === 'unknown') ||
      [browserHttp, apiHttp, collectionHttp, searchHttp].some(
        p => p.state !== 'responding',
      ) ||
      data.state !== 'observed' ||
      apps.length !== names.length ||
      workloads.some(
        w => w.cpu?.state !== 'observed' || w.memory?.state !== 'observed',
      ),
    workloads,
    probes: {
      browser: browserHttp,
      apiLanding: apiHttp,
      collections: collectionHttp,
      search: searchHttp,
    },
    data,
    argo: {
      state:
        argoResult.status === 'fulfilled' && apps.length === names.length
          ? 'observed'
          : 'partial-or-unknown',
      observedAt: stamp(),
      applications: apps,
    },
    recovery: {
      state: 'not-certified-by-this-check',
      detail:
        'A database pod being ready is not backup/restore evidence. Consult the STAC runbook and volume-backup evidence before any data/schema change.',
    },
  };
}
export type StacEvidence = Awaited<ReturnType<typeof observeStacEvidence>>;
let cached: StacEvidence | undefined;
let pending: Promise<StacEvidence> | undefined;
export function stacEvidence(): Promise<StacEvidence> {
  if (cached && Date.now() - Date.parse(cached.observedAt) < 60_000)
    return Promise.resolve(cached);
  if (!pending)
    pending = observeStacEvidence()
      .then(value => {
        cached = value;
        return value;
      })
      .finally(() => {
        pending = undefined;
      });
  return pending;
}
