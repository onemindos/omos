/* eslint-disable no-nested-ternary */
import type { Entity } from '@backstage/catalog-model';
import type { ClusterFacts, WorkloadFact } from './types';
import { listDocuments } from './yamlText';

/**
 * Stack hygiene: pure comparison of catalog, Argo CD, live workloads and Git. Read-only.
 * A check that cannot be evaluated is reported as `unavailable` with its reason, never
 * as an empty list, so a missing data source cannot look like a clean stack.
 */

export type FindingAction = {
  kind: 'add' | 'retire' | 'review';
  label: string;
  href: string;
};

export type Finding = {
  check: string;
  severity: 'info' | 'warn';
  subject: string;
  detail: string;
  actions: FindingAction[];
};

export type CheckResult = {
  id: string;
  title: string;
  status: 'ok' | 'findings' | 'unavailable';
  count: number;
  reason?: string;
};

export type HygieneReport = {
  generatedAt: string;
  trigger: 'schedule' | 'manual';
  clusterObservedAt: string | null;
  checks: CheckResult[];
  findings: Finding[];
};

export type HygieneInput = {
  now: Date;
  trigger: 'schedule' | 'manual';
  /** Undefined when the cluster could not be read. */
  cluster?: ClusterFacts;
  /** Components of the catalog (apps and workloads). */
  entities: Entity[];
  /** `namespace/pvc` for every VolSync ReplicationSource in Git; undefined when Git could not be read. */
  volsyncPvcs?: Set<string>;
  /** DigitalOcean volumes and snapshots (F6). Absent until F6 exists. */
  digitalocean?: undefined;
};

export const CHECKS: { id: string; title: string }[] = [
  { id: 'argo-without-catalog', title: 'Argo apps without a catalog entry' },
  { id: 'catalog-without-workload', title: 'Catalog apps without a running workload' },
  { id: 'workload-without-owner', title: 'Workloads with no owner app' },
  { id: 'no-health-check', title: 'Apps with no health check' },
  { id: 'no-backup', title: 'Apps with a volume and no backup' },
  { id: 'no-runbook', title: 'Apps with no runbook' },
  { id: 'image-without-digest', title: 'Images without a digest pin' },
  { id: 'orphan-volume', title: 'Orphan DigitalOcean volumes' },
  { id: 'old-snapshot', title: 'Old DigitalOcean snapshots' },
];

const ann = (e: Entity, key: string) => e.metadata.annotations?.[key];
const csv = (value?: string) =>
  (value ?? '').split(',').map(x => x.trim()).filter(Boolean);
const role = (e: Entity) => e.metadata.labels?.['onemind.io/catalog-role'];
const createUrl = (template: string, formData: object) =>
  `/create/templates/default/${template}?formData=${encodeURIComponent(JSON.stringify(formData))}`;
const entityUrl = (e: Entity) => `/catalog/default/component/${e.metadata.name}`;

const ELIGIBLE_KINDS = new Set(['Deployment', 'StatefulSet', 'DaemonSet', 'CronJob']);
const BACKUP_JOB = /backup|dump/;

export function computeHygiene(input: HygieneInput): HygieneReport {
  const { cluster, entities } = input;
  const findings: Finding[] = [];
  const unavailable = new Map<string, string>();

  const components = entities.filter(e => e.kind === 'Component');
  const apps = components.filter(e => role(e) !== 'workload');
  const workloadEntities = components.filter(e => role(e) === 'workload');
  const k8sApps = apps.filter(
    e => ann(e, 'onemind.io/kubernetes-workloads') || ann(e, 'argocd/app-name'),
  );

  if (!cluster) {
    for (const id of [
      'argo-without-catalog',
      'catalog-without-workload',
      'workload-without-owner',
      'no-health-check',
      'no-backup',
      'image-without-digest',
    ])
      unavailable.set(id, 'The cluster could not be read.');
  }
  unavailable.set(
    'orphan-volume',
    'DigitalOcean volume data is not available: F6 (the DigitalOcean page) is not built.',
  );
  unavailable.set(
    'old-snapshot',
    'DigitalOcean snapshot data is not available: F6 (the DigitalOcean page) is not built.',
  );
  if (!input.volsyncPvcs && cluster)
    unavailable.set('no-backup', 'Git could not be read, so VolSync sources are unknown.');

  // Index what the catalog declares.
  const declaredArgo = new Set(
    apps.flatMap(e => [ann(e, 'argocd/app-name') ?? '', ...csv(ann(e, 'onemind.io/argocd-apps'))]).filter(Boolean),
  );
  const ownedWorkloads = new Map<string, Entity>();
  for (const app of apps)
    for (const w of csv(ann(app, 'onemind.io/kubernetes-workloads'))) ownedWorkloads.set(w, app);
  for (const w of workloadEntities) {
    const ns = ann(w, 'backstage.io/kubernetes-namespace');
    const name = ann(w, 'onemind.io/kubernetes-workload');
    const parent = (w.spec?.subcomponentOf as string | undefined)?.split('/')[1];
    const owner = apps.find(a => a.metadata.name === parent);
    if (ns && name && owner) ownedWorkloads.set(`${ns}/${name}`, owner);
  }

  const liveWorkloads = (cluster?.workloads ?? []).filter(w => ELIGIBLE_KINDS.has(w.kind));
  const key = (w: WorkloadFact) => `${w.namespace}/${w.name}`;

  if (cluster) {
    for (const argo of cluster.argoApps) {
      if (declaredArgo.has(argo.name)) continue;
      findings.push({
        check: 'argo-without-catalog',
        severity: 'warn',
        subject: argo.name,
        detail: `Argo application ${argo.name} (${argo.sync ?? '?'}/${argo.health ?? '?'}) is named by no catalog app.`,
        actions: [
          {
            kind: 'add',
            label: 'Add to catalog',
            href: createUrl('add-app', {
              name: argo.name.replace(/^(platform|ecosystem|edge|monitoring)-/, ''),
              hosting: 'kubernetes',
            }),
          },
        ],
      });
    }
    const liveKeys = new Set(liveWorkloads.map(key));
    for (const app of k8sApps) {
      const listed = csv(ann(app, 'onemind.io/kubernetes-workloads'));
      if (!listed.length || listed.some(w => liveKeys.has(w))) continue;
      findings.push({
        check: 'catalog-without-workload',
        severity: 'warn',
        subject: app.metadata.name,
        detail: `None of its ${listed.length} catalogued workload(s) is running: ${listed.slice(0, 4).join(', ')}.`,
        actions: [
          {
            kind: 'retire',
            label: 'Retire app',
            href: createUrl('retire-app', { app: `component:default/${app.metadata.name}` }),
          },
          { kind: 'review', label: 'Open entity', href: entityUrl(app) },
        ],
      });
    }
    for (const w of liveWorkloads) {
      if (ownedWorkloads.has(key(w))) continue;
      findings.push({
        check: 'workload-without-owner',
        severity: 'warn',
        subject: key(w),
        detail: `${w.kind} ${key(w)} is running and no catalog app lists it.`,
        actions: [
          {
            kind: 'add',
            label: 'Add app',
            href: createUrl('add-app', { name: w.name, hosting: 'kubernetes' }),
          },
        ],
      });
    }
    for (const w of liveWorkloads) {
      const bare = w.images.filter(image => !/@sha256:[a-f0-9]{64}/.test(image));
      if (!bare.length) continue;
      const owner = ownedWorkloads.get(key(w));
      findings.push({
        check: 'image-without-digest',
        severity: 'warn',
        subject: key(w),
        detail: `Not pinned by digest: ${bare.map(i => i.slice(0, 90)).join(', ')}.`,
        actions: owner ? [{ kind: 'review', label: 'Open app', href: entityUrl(owner) }] : [],
      });
    }
    for (const app of k8sApps) {
      const mine = liveWorkloads.filter(w => ownedWorkloads.get(key(w)) === app);
      if (!mine.length) continue;
      if (!ann(app, 'onemind.io/health-check') && !mine.some(w => w.hasProbe)) {
        findings.push({
          check: 'no-health-check',
          severity: 'warn',
          subject: app.metadata.name,
          detail: 'No readiness or liveness probe on any of its workloads and no `onemind.io/health-check` annotation.',
          actions: [{ kind: 'review', label: 'Open app', href: entityUrl(app) }],
        });
      }
      if (input.volsyncPvcs) {
        const namespaces = new Set(mine.map(w => w.namespace));
        const claims = [...new Set(mine.flatMap(w => w.claims.map(c => `${w.namespace}/${c}`)))];
        const uncovered = claims.filter(c => !input.volsyncPvcs!.has(c));
        const dumpJob = cluster.cronJobs.some(c => namespaces.has(c.namespace) && BACKUP_JOB.test(c.name));
        if (uncovered.length && !dumpJob) {
          findings.push({
            check: 'no-backup',
            severity: 'warn',
            subject: app.metadata.name,
            detail: `No VolSync source and no backup CronJob for ${uncovered.slice(0, 4).join(', ')}.`,
            actions: [{ kind: 'review', label: 'Open app', href: entityUrl(app) }],
          });
        }
      }
    }
  }
  for (const app of k8sApps) {
    if (!ann(app, 'onemind.io/recovery-runbook')) {
      findings.push({
        check: 'no-runbook',
        severity: 'info',
        subject: app.metadata.name,
        detail: 'No `onemind.io/recovery-runbook` annotation.',
        actions: [{ kind: 'review', label: 'Open app', href: entityUrl(app) }],
      });
    }
  }

  const checks: CheckResult[] = CHECKS.map(check => {
    const reason = unavailable.get(check.id);
    if (reason)
      return { ...check, status: 'unavailable', count: 0, reason };
    const count = findings.filter(f => f.check === check.id).length;
    return { ...check, status: count ? 'findings' : 'ok', count };
  });
  const order = new Map(CHECKS.map((c, i) => [c.id, i]));
  findings.sort(
    (a, b) =>
      (a.severity === b.severity ? 0 : a.severity === 'warn' ? -1 : 1) ||
      (order.get(a.check)! - order.get(b.check)!) ||
      a.subject.localeCompare(b.subject),
  );
  return {
    generatedAt: input.now.toISOString(),
    trigger: input.trigger,
    clusterObservedAt: cluster?.observedAt ?? null,
    checks,
    findings,
  };
}

/** `namespace/pvc` for every ReplicationSource document in the VolSync sources file. */
export function volsyncPvcsFrom(text: string) {
  return new Set(
    listDocuments(text)
      .filter(d => d.ref.kind === 'ReplicationSource' && d.doc.spec?.sourcePVC)
      .map(d => `${d.ref.namespace}/${d.doc.spec.sourcePVC}`),
  );
}
