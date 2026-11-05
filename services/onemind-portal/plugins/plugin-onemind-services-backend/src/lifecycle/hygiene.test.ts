/** @jest-environment node */
import type { Entity } from '@backstage/catalog-model';
import express from 'express';
import type { Server } from 'node:http';
import { computeHygiene } from './hygiene';
import { hygieneRouter, memoryStore, MIN_REFRESH_MS } from './hygienePlugin';
import type { ClusterFacts, WorkloadFact } from './types';
import { volsyncPvcsFrom } from './hygiene';

const NOW = new Date('2026-11-05T12:00:00Z');
const pinned = `ghcr.io/x/y:1@sha256:${'c'.repeat(64)}`;

const app = (name: string, annotations: Record<string, string> = {}): Entity => ({
  apiVersion: 'backstage.io/v1alpha1',
  kind: 'Component',
  metadata: {
    name,
    annotations: {
      'argocd/app-name': `platform-${name}`,
      'onemind.io/kubernetes-workloads': `${name}/${name}`,
      'onemind.io/recovery-runbook': `apps/${name}/`,
      ...annotations,
    },
    labels: { 'onemind.io/catalog-role': 'app' },
  },
  spec: { type: 'service', owner: 'group:default/onemind-admins' },
});

const workload = (ns: string, name: string, over: Partial<WorkloadFact> = {}): WorkloadFact => ({
  kind: 'Deployment',
  namespace: ns,
  name,
  images: [pinned],
  hasProbe: true,
  claims: [],
  ...over,
});

// Seeded fixtures: each defect below must be flagged by exactly the check that owns it.
const cluster = (): ClusterFacts => ({
  observedAt: '2026-11-05T11:59:00Z',
  argoApps: [
    { name: 'platform-good', finalizers: [], automated: true, prune: true, sourcePaths: [], charts: [], valueFiles: [], resources: [] },
    { name: 'platform-unmapped', finalizers: [], automated: false, prune: false, sourcePaths: [], charts: [], valueFiles: [], resources: [], sync: 'Synced', health: 'Healthy' },
  ],
  workloads: [
    workload('good', 'good'),
    workload('stray', 'stray'),
    workload('unpinned', 'unpinned', { images: ['ghcr.io/x/y:latest'] }),
    workload('noprobe', 'noprobe', { hasProbe: false }),
    workload('data', 'data', { claims: ['data'] }),
    workload('safe', 'safe', { claims: ['data'] }),
  ],
  pvcs: [],
  ingresses: [],
  cronJobs: [],
});

const entities = (): Entity[] => [
  app('good'),
  app('ghost'),
  app('unpinned'),
  app('noprobe'),
  app('data'),
  app('safe'),
  app('norunbook', { 'onemind.io/recovery-runbook': '' }),
].map(e => {
  if (e.metadata.name === 'norunbook') delete e.metadata.annotations!['onemind.io/recovery-runbook'];
  return e;
});

const run = (over: Partial<Parameters<typeof computeHygiene>[0]> = {}) =>
  computeHygiene({
    now: NOW,
    trigger: 'manual',
    cluster: cluster(),
    entities: entities(),
    volsyncPvcs: new Set(['safe/data']),
    ...over,
  });

const subjects = (report: ReturnType<typeof run>, check: string) =>
  report.findings.filter(f => f.check === check).map(f => f.subject).sort();

describe('computeHygiene', () => {
  it('flags every seeded fixture', () => {
    const report = run();
    expect(subjects(report, 'argo-without-catalog')).toEqual(['platform-unmapped']);
    expect(subjects(report, 'catalog-without-workload')).toEqual(['ghost', 'norunbook']);
    expect(subjects(report, 'workload-without-owner')).toEqual(['stray/stray']);
    expect(subjects(report, 'image-without-digest')).toEqual(['unpinned/unpinned']);
    expect(subjects(report, 'no-health-check')).toEqual(['noprobe']);
    expect(subjects(report, 'no-backup')).toEqual(['data']);
    expect(subjects(report, 'no-runbook')).toEqual(['norunbook']);
  });

  it('links each finding to Add or Retire with pre-filled parameters', () => {
    const report = run();
    const form = (href: string) => JSON.parse(decodeURIComponent(href.split('formData=')[1]));
    const unmapped = report.findings.find(f => f.check === 'argo-without-catalog')!;
    expect(unmapped.actions[0].href).toMatch(/^\/create\/templates\/default\/add-app\?formData=/);
    expect(form(unmapped.actions[0].href)).toEqual({ name: 'unmapped', hosting: 'kubernetes' });
    const ghost = report.findings.find(f => f.check === 'catalog-without-workload' && f.subject === 'ghost')!;
    expect(ghost.actions[0].href).toMatch(/\/retire-app\?formData=/);
    expect(form(ghost.actions[0].href)).toEqual({ app: 'component:default/ghost' });
  });

  it('marks DigitalOcean checks unavailable with a reason, never as an empty pass', () => {
    const report = run();
    for (const id of ['orphan-volume', 'old-snapshot']) {
      const check = report.checks.find(c => c.id === id)!;
      expect(check.status).toBe('unavailable');
      expect(check.reason).toMatch(/F6/);
      expect(report.findings.some(f => f.check === id)).toBe(false);
    }
    expect(report.checks.find(c => c.id === 'argo-without-catalog')!.status).toBe('findings');
  });

  it('marks cluster-dependent checks unavailable when the cluster is unreadable', () => {
    const report = run({ cluster: undefined });
    for (const id of ['argo-without-catalog', 'workload-without-owner', 'no-health-check', 'no-backup', 'image-without-digest'])
      expect(report.checks.find(c => c.id === id)!.status).toBe('unavailable');
    expect(report.clusterObservedAt).toBeNull();
    // Catalog-only checks still run.
    expect(subjects(report, 'no-runbook')).toEqual(['norunbook']);
  });

  it('marks the backup check unavailable when Git could not be read', () => {
    const report = run({ volsyncPvcs: undefined });
    expect(report.checks.find(c => c.id === 'no-backup')).toMatchObject({ status: 'unavailable' });
    expect(subjects(report, 'no-backup')).toEqual([]);
  });

  it('reads VolSync sources from manifests', () => {
    const text = 'apiVersion: volsync.backube/v1alpha1\nkind: ReplicationSource\nmetadata:\n  name: a\n  namespace: ns\nspec:\n  sourcePVC: pvc\n';
    expect([...volsyncPvcsFrom(text)]).toEqual(['ns/pvc']);
  });
});

describe('stack hygiene routes', () => {
  let server: Server;
  let base: string;
  let clock = Date.parse('2026-11-05T12:00:00Z');
  const store = memoryStore();
  const run1 = jest.fn(async (trigger: 'schedule' | 'manual') => {
    const report = computeHygiene({ now: new Date(clock), trigger, cluster: cluster(), entities: entities(), volsyncPvcs: new Set() });
    await store.save(report);
    return report;
  });
  beforeAll(async () => {
    const server2 = express();
    server2.use(
      hygieneRouter({
        store,
        run: run1,
        now: () => clock,
        user: async req => {
          if (!req.get('Authorization')) throw Object.assign(new Error('no'), { status: 401 });
        },
        admin: async req => {
          if (req.get('Authorization') !== 'Bearer admin') throw Object.assign(new Error('no'), { status: 403 });
        },
      }),
    );
    server2.use((error: any, _req: express.Request, res: express.Response, _next: express.NextFunction) => {
      res.status(error.status ?? 500).json({ error: error.message });
    });
    await new Promise<void>(resolve => {
      server = server2.listen(0, resolve);
    });
    base = `http://localhost:${(server.address() as any).port}`;
  });
  afterAll(() => new Promise<void>(resolve => server.close(() => resolve())));

  const call = (path: string, token?: string, method = 'GET') =>
    fetch(`${base}${path}`, { method, headers: token ? { Authorization: `Bearer ${token}` } : {} });

  it('needs a signed-in user to read, and has no report before the first run', async () => {
    expect((await call('/stack-hygiene')).status).toBe(401);
    const response = await call('/stack-hygiene', 'user');
    expect(await response.json()).toEqual({ report: null });
  });

  it('refreshes for administrators only, and rate-limits repeat refreshes', async () => {
    expect((await call('/stack-hygiene/refresh', 'user', 'POST')).status).toBe(403);
    expect(run1).not.toHaveBeenCalled();
    const first = await (await call('/stack-hygiene/refresh', 'admin', 'POST')).json();
    expect(first.cached).toBe(false);
    expect(run1).toHaveBeenCalledWith('manual');
    const second = await (await call('/stack-hygiene/refresh', 'admin', 'POST')).json();
    expect(second.cached).toBe(true);
    expect(run1).toHaveBeenCalledTimes(1);
    clock += MIN_REFRESH_MS + 1;
    expect((await (await call('/stack-hygiene/refresh', 'admin', 'POST')).json()).cached).toBe(false);
    const read = await (await call('/stack-hygiene', 'user')).json();
    expect(read.report.findings.length).toBeGreaterThan(0);
  });
});
