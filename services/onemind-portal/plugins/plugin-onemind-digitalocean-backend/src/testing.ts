/* Test-only helpers (never imported by index.ts): the recorded, redacted DigitalOcean answers. */
import { readFileSync } from 'node:fs';
import { resolve } from 'node:path';
import { resolvePackagePath } from '@backstage/backend-plugin-api';
import { listPvcs, type WorkloadComponent, type Pvc } from './links';

const read = (name: string) =>
  JSON.parse(
    readFileSync(
      resolve(
        resolvePackagePath(
          '@internal/plugin-onemind-digitalocean-backend',
          'src',
        ),
        '__fixtures__',
        name,
      ),
      'utf8',
    ),
  );

export const api: Record<string, unknown> & {
  _meta: { counts: Record<string, number> };
} = read('do-api.json');
export const counts = api._meta.counts;
/** The claims as listPvcs reads them from the recorded claim and pod lists. */
export const pvcs = (): Promise<Pvc[]> =>
  listPvcs(async path =>
    path.startsWith('/api/v1/pods')
      ? read('k8s-pods.json')
      : read('k8s-pvcs.json'),
  );
export const workloads: WorkloadComponent[] = read(
  'catalog-workloads.json',
).map(
  (w: {
    name: string;
    namespace: string;
    workload: string;
    selector: string;
  }) => ({
    ref: `component:default/${w.name}`,
    namespace: w.namespace,
    workload: w.workload,
    selector: w.selector,
  }),
);

/** A made-up value (assembled at run time so secret scanners have nothing to match). */
export const TOKEN = ['fake', 'do', 'token', 'for', 'tests', '0123456789'].join(
  '-',
);

export interface FakeOptions {
  /** Paths (no query) that answer with this status instead of the recording. */
  status?: Record<string, number>;
  /** Replaces a recorded answer. */
  override?: Record<string, unknown>;
  /** Paths whose items are dropped (by id) before answering. */
  drop?: Record<string, string[]>;
}

/** A fetch that answers from the recording, by path, and remembers every request. */
export function fakeFetch(options: FakeOptions = {}) {
  const calls: { url: string; init?: RequestInit }[] = [];
  const impl = (async (input: RequestInfo | URL, init?: RequestInit) => {
    const url = new URL(String(input));
    calls.push({ url: String(input), init });
    const path = url.pathname;
    const status = options.status?.[path];
    if (status) return new Response('{"id":"x","message":"no"}', { status });
    let body = (options.override?.[path] ?? api[path]) as
      | Record<string, unknown>
      | undefined;
    if (!body) return new Response('{"id":"not_found"}', { status: 404 });
    const drop = options.drop?.[path];
    if (drop) {
      body = JSON.parse(JSON.stringify(body));
      for (const [k, v] of Object.entries(body!))
        if (Array.isArray(v))
          (body as Record<string, unknown>)[k] = v.filter(
            (x: { id?: unknown }) => !drop.includes(String(x.id)),
          );
    }
    return new Response(JSON.stringify(body), { status: 200 });
  }) as typeof fetch;
  return { impl, calls };
}
