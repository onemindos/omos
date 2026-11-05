/* eslint-disable no-nested-ternary */
import type { Entity } from '@backstage/catalog-model';
import {
  bumpCoverageTestCount,
  bumpStructureTestCount,
  appendSection,
  PATHS,
  removeMapEntries,
  removeCoverageExceptions,
  removePrimaryArgo,
  removeRefBaseline,
  removeStructureApp,
} from './catalogEdits';
import type {
  ClusterFacts,
  InventorySection,
  PvcFact,
  RepoReader,
  Retention,
} from './types';
import { DocRef, listDocuments, removeDocuments } from './yamlText';

/**
 * Read-only retirement planning. Everything here reads (catalog, Git through the GitHub
 * App, the cluster through the portal's read-only service account) and computes. The
 * result is a set of file contents to commit, a list of files to delete and a markdown
 * body. Nothing is deleted, synced or applied from here.
 */

export type RetireInput = {
  target: Entity;
  /** Components and Resources of the catalog. */
  catalog: Entity[];
  repo: RepoReader;
  /** Undefined when the cluster could not be read. */
  cluster?: ClusterFacts;
  retention: Retention;
  reason: string;
  now?: Date;
};

export type RetirePlan = {
  scope: 'app' | 'workload';
  ref: string;
  name: string;
  title: string;
  namespaces: string[];
  inventory: InventorySection[];
  filesToDelete: string[];
  /** Full new content of edited files plus new files (the retirement record). */
  files: Map<string, string>;
  /** Files deliberately left in place, with the reason. */
  kept: { path: string; reason: string }[];
  /** Things a person must do or check that no file edit can express. */
  warnings: string[];
  body: string;
  recordPath: string;
};

export class RetireInputError extends Error {
  readonly name = 'InputError';
}

const ARGO_DIRS = [
  'platform/argocd',
  'platform/ecosystem/applications',
  'platform/ecosystem/foundation-applications',
];
/** Directories shared by many apps: never deleted wholesale. */
const PROTECTED = /^platform\/(argocd|core|ecosystem|security|backups|nats|monitoring|kube|ops|values|parked)(\/|$)/;
const WORKLOAD_KINDS = new Set([
  'Deployment',
  'StatefulSet',
  'DaemonSet',
  'CronJob',
]);
const OWNED_BY_NAME = new Set([
  ...WORKLOAD_KINDS,
  'Service',
  'Ingress',
  'HorizontalPodAutoscaler',
  'PodDisruptionBudget',
  'ServiceAccount',
  'ConfigMap',
  'SealedSecret',
  'ServiceMonitor',
]);

const ann = (e: Entity, key: string) => e.metadata.annotations?.[key];
const csv = (value?: string) =>
  (value ?? '').split(',').map(x => x.trim()).filter(Boolean);
const refOf = (e: Entity) =>
  `${e.kind.toLowerCase()}:${(e.metadata.namespace ?? 'default').toLowerCase()}/${e.metadata.name.toLowerCase()}`;
const unique = <T>(items: T[]) => [...new Set(items)];

type AppFile = { path: string; text: string; docs: { ref: DocRef; doc: any }[] };

async function readApplicationFiles(repo: RepoReader) {
  const out: AppFile[] = [];
  for (const dir of ARGO_DIRS) {
    for (const [path, text] of await repo.readDir(dir)) {
      if (!/\.ya?ml$/.test(path)) continue;
      const docs = listDocuments(text).filter(x => x.ref.kind === 'Application');
      if (docs.length) out.push({ path, text, docs });
    }
  }
  return out;
}

function sourcesOf(application: any): any[] {
  return (
    application.spec?.sources ??
    (application.spec?.source ? [application.spec.source] : [])
  );
}
const localPaths = (application: any) =>
  sourcesOf(application)
    .filter(s => /onemindos\/omos/.test(s.repoURL ?? '') && typeof s.path === 'string')
    .map(s => String(s.path).replace(/\/+$/, ''));
const valueFiles = (application: any) =>
  sourcesOf(application)
    .flatMap(s => s.helm?.valueFiles ?? [])
    .map((v: string) => {
      const i = v.indexOf('platform/');
      return i >= 0 ? v.slice(i) : v;
    });

const related = (a: string, b: string) =>
  a === b || a.startsWith(`${b}/`) || b.startsWith(`${a}/`);

function ingressHosts(doc: any): string[] {
  const hosts = [
    ...(doc.spec?.rules ?? []).map((r: any) => r.host),
    ...(doc.spec?.tls ?? []).flatMap((t: any) => t.hosts ?? []),
    doc.metadata?.annotations?.['tailscale.com/hostname'],
  ].filter(Boolean);
  return unique(hosts as string[]);
}
function ingressBackends(doc: any): string[] {
  const names = [
    doc.spec?.defaultBackend?.service?.name,
    ...(doc.spec?.rules ?? []).flatMap((r: any) =>
      (r.http?.paths ?? []).map((p: any) => p.backend?.service?.name),
    ),
  ].filter(Boolean);
  return unique(names as string[]);
}

function matchesWorkload(ref: DocRef, doc: any, workload: string) {
  if (!OWNED_BY_NAME.has(ref.kind)) return false;
  if (ref.kind === 'Ingress') {
    const backends = ingressBackends(doc);
    return (
      ref.name === workload ||
      (backends.length > 0 && backends.every(b => b === workload))
    );
  }
  return ref.name === workload;
}

function mentions(doc: any, word: string) {
  return JSON.stringify(doc).includes(`"${word}"`);
}

/** True when a `matchLabels` in the document only uses labels the workload's pods carry. */
function selectsPods(doc: any, podLabels: Record<string, string>[]): boolean {
  const found: Record<string, string>[] = [];
  const walk = (node: any) => {
    if (!node || typeof node !== 'object') return;
    if (node.matchLabels && typeof node.matchLabels === 'object') found.push(node.matchLabels);
    for (const value of Object.values(node)) walk(value);
  };
  walk(doc);
  return found.some(
    labels =>
      Object.keys(labels).length > 0 &&
      podLabels.some(pod => Object.entries(labels).every(([k, v]) => pod[k] === v)),
  );
}

function matchingLines(text: string, words: string[], limit = 4) {
  const out: string[] = [];
  const lines = text.split('\n');
  for (let i = 0; i < lines.length && out.length < limit; i++) {
    if (words.some(w => w && new RegExp(`(^|[^a-z0-9-])${w.replace(/[.*+?^${}()|[\]\\]/g, '\\$&')}([^a-z0-9-]|$)`).test(lines[i])))
      out.push(`${i + 1}: ${lines[i].trim().slice(0, 100)}`);
  }
  return out;
}

export async function buildRetirePlan(input: RetireInput): Promise<RetirePlan> {
  const { target, catalog, repo, cluster, retention } = input;
  const now = input.now ?? new Date();
  const today = now.toISOString().slice(0, 10);
  if (target.kind !== 'Component')
    throw new RetireInputError('Only Components can be retired');
  const role = target.metadata.labels?.['onemind.io/catalog-role'];
  const scope: 'app' | 'workload' = role === 'workload' ? 'workload' : 'app';

  const byRef = new Map(catalog.map(e => [refOf(e), e]));
  let app = target;
  let workloadName: string | undefined;
  let workloadNamespace: string | undefined;
  if (scope === 'workload') {
    workloadName = ann(target, 'onemind.io/kubernetes-workload');
    workloadNamespace = ann(target, 'backstage.io/kubernetes-namespace');
    const parentRef = target.spec?.subcomponentOf as string | undefined;
    const parent = parentRef ? byRef.get(parentRef.toLowerCase()) : undefined;
    if (!workloadName || !parent)
      throw new RetireInputError(
        'This workload entity has no workload name or parent app in the catalog',
      );
    app = parent;
  }
  const appRef = refOf(app);
  const subcomponents = catalog.filter(
    e => (e.spec?.subcomponentOf as string | undefined)?.toLowerCase() === appRef,
  );
  const argoNames = unique([
    ...csv(ann(app, 'onemind.io/argocd-apps')),
    ...csv(ann(app, 'argocd/app-name')),
  ]);
  const argoResources = catalog.filter(
    e =>
      e.kind === 'Resource' &&
      e.spec?.type === 'argocd-application' &&
      argoNames.includes(ann(e, 'argocd/app-name') ?? ''),
  );
  const workloadRefs = scope === 'workload'
    ? [`${workloadNamespace}/${workloadName}`]
    : csv(ann(app, 'onemind.io/kubernetes-workloads'));
  const namespaces = unique(
    scope === 'workload'
      ? [workloadNamespace ?? '']
      : [
          ...csv(ann(app, 'onemind.io/kubernetes-namespaces')),
          ann(app, 'backstage.io/kubernetes-namespace') ?? '',
          ...workloadRefs.map(w => w.split('/')[0]),
        ],
  ).filter(Boolean);

  const warnings: string[] = [];
  const kept: RetirePlan['kept'] = [];
  const filesToDelete: string[] = [];
  const files = new Map<string, string>();
  const textOf = new Map<string, string>();
  const readOnce = async (path: string) => {
    if (files.has(path)) return files.get(path);
    if (!textOf.has(path)) {
      const text = await repo.readFile(path);
      if (text !== undefined) textOf.set(path, text);
      else return undefined;
    }
    return textOf.get(path);
  };
  const edit = (path: string, text: string) => {
    files.set(path, text);
  };

  // --- Argo applications (Git side and live side) -----------------------------------
  const appFiles = await readApplicationFiles(repo);
  const targetApps = appFiles.flatMap(f =>
    f.docs.filter(d => argoNames.includes(d.ref.name)).map(d => ({ file: f, ...d })),
  );
  const otherApps = appFiles.flatMap(f =>
    f.docs.filter(d => !argoNames.includes(d.ref.name)).map(d => ({ file: f, ...d })),
  );
  const liveApp = (name: string) => cluster?.argoApps.find(a => a.name === name);

  const argo: InventorySection = {
    id: 'argo',
    title: 'Argo CD application',
    state: argoNames.length ? 'found' : 'none',
    items: [],
  };
  for (const name of argoNames) {
    const found = targetApps.find(a => a.ref.name === name);
    const lv = liveApp(name);
    const finalizers = lv?.finalizers ?? found?.doc.metadata?.finalizers ?? [];
    const cascades = finalizers.includes('resources-finalizer.argocd.argoproj.io');
    const parentPrunes = found ? found.file.path.startsWith('platform/argocd/') : undefined;
    argo.items.push({
      label: name,
      detail: [
        found ? `defined in \`${found.file.path}\`` : 'no Application file found in Git',
        lv ? `live ${lv.sync ?? '?'}/${lv.health ?? '?'}` : cluster ? 'not live' : 'live state unavailable',
        cascades ? 'has resources-finalizer (deleting it cascades to its resources)' : 'no resources-finalizer (deleting it orphans its resources)',
        parentPrunes === false ? 'parent does not prune: delete the Application by hand' : parentPrunes ? 'parent app-of-apps prunes: merge deletes the Application CR' : '',
      ].filter(Boolean).join('; '),
      link: `https://argocd.tail717ba0.ts.net/applications/argocd/${encodeURIComponent(name)}`,
    });
    if (!found)
      warnings.push(`Application ${name} has no file under ${ARGO_DIRS.join(', ')}; check how it is created before merging.`);
  }

  // --- repository files that belong to the app ---------------------------------------
  const sourcePaths = unique(targetApps.flatMap(a => localPaths(a.doc)));
  const sourceFiles = new Map<string, string>();
  for (const path of sourcePaths) {
    for (const [file, text] of await repo.readDir(path)) sourceFiles.set(file, text);
  }
  const sharedPaths = sourcePaths.filter(p =>
    otherApps.some(o => localPaths(o.doc).some(q => related(p, q))),
  );
  if (scope === 'app') {
    // Application files.
    for (const f of unique(targetApps.map(a => a.file))) {
      const allTarget = f.docs.every(d => argoNames.includes(d.ref.name));
      if (allTarget) filesToDelete.push(f.path);
      else {
        const res = removeDocuments(f.text, r => r.kind === 'Application' && argoNames.includes(r.name));
        if (res.text !== undefined) edit(f.path, res.text);
      }
    }
    for (const path of sourcePaths) {
      if (sharedPaths.includes(path)) {
        kept.push({ path: `${path}/`, reason: 'Another Argo application also sources this directory' });
        warnings.push(`Source directory ${path}/ is shared with another Application; it was not deleted.`);
        continue;
      }
      if (PROTECTED.test(path)) {
        kept.push({ path: `${path}/`, reason: 'Protected platform directory: remove this app\'s files by hand' });
        warnings.push(`Source directory ${path}/ is a protected platform directory; remove only this app's files by hand.`);
        continue;
      }
      for (const [file] of [...sourceFiles].filter(([f]) => f.startsWith(`${path}/`))) {
        if (/\.md$/i.test(file)) kept.push({ path: file, reason: 'Documentation kept for the record; delete it in review if not wanted' });
        else if (/restic|backup[^/]*sealedsecret/i.test(file)) kept.push({ path: file, reason: 'Restic repository credentials: keep until the retention period ends or the backup repository becomes unreadable' });
        else filesToDelete.push(file);
      }
    }
    // Value files used only by this app.
    for (const a of targetApps) {
      for (const vf of valueFiles(a.doc)) {
        const shared = otherApps.some(o => valueFiles(o.doc).includes(vf));
        const text = await readOnce(vf);
        if (text === undefined) continue;
        if (shared) kept.push({ path: vf, reason: 'Also used by another Application' });
        else filesToDelete.push(vf);
      }
    }
  } else {
    // Workload scope: only documents named after the workload, in the parent's source dirs.
    const parentApps = appFiles.flatMap(f =>
      f.docs.filter(d => argoNames.includes(d.ref.name)),
    );
    const dirs = unique(parentApps.flatMap(d => localPaths(d.doc)));
    for (const dir of dirs) {
      for (const [file, text] of await repo.readDir(dir)) {
        sourceFiles.set(file, text);
      }
    }
    const podLabels: Record<string, string>[] = [];
    for (const [path, text] of sourceFiles) {
      if (!/\.ya?ml$/.test(path)) continue;
      const res = removeDocuments(text, (ref, doc) => {
        const hit = matchesWorkload(ref, doc, workloadName!);
        if (hit && WORKLOAD_KINDS.has(ref.kind) && doc.spec?.selector?.matchLabels)
          podLabels.push(doc.spec.selector.matchLabels);
        return hit;
      });
      if (!res.removed.length) continue;
      if (res.text === undefined) filesToDelete.push(path);
      else edit(path, res.text);
    }
    for (const [path, text] of sourceFiles) {
      if (!/\.ya?ml$/.test(path)) continue;
      const removedFile = filesToDelete.includes(path);
      if (removedFile) continue;
      const remaining = listDocuments(files.get(path) ?? text);
      for (const d of remaining) {
        if (matchesWorkload(d.ref, d.doc, workloadName!)) continue;
        if (d.ref.kind === 'PersistentVolumeClaim' && d.ref.name.includes(workloadName!)) {
          warnings.push(`PVC manifest ${d.ref.name} in ${path} looks like this workload's data; it is not auto-removed (data).`);
        } else if (mentions(d.doc, workloadName!) || selectsPods(d.doc, podLabels)) {
          warnings.push(`${path}: ${d.ref.kind}/${d.ref.name} still references or selects \`${workloadName}\`; review it.`);
        }
      }
    }
  }

  // --- cluster facts ---------------------------------------------------------------
  const clusterNote = cluster ? undefined : 'The cluster could not be read, so live state is unknown.';
  const wlSelected = (cluster?.workloads ?? []).filter(w =>
    scope === 'workload'
      ? w.namespace === workloadNamespace && w.name === workloadName
      : workloadRefs.includes(`${w.namespace}/${w.name}`),
  );
  const workloads: InventorySection = {
    id: 'workloads',
    title: 'Workloads',
    state: !cluster ? 'unavailable' : wlSelected.length ? 'found' : 'none',
    note: clusterNote ?? (wlSelected.length ? undefined : 'No matching live workload; the catalog may be stale.'),
    items: wlSelected.map(w => ({
      label: `${w.kind} ${w.namespace}/${w.name}`,
      detail: `${w.ready ?? '?'}/${w.replicas ?? '?'} ready; images: ${w.images.map(i => i.split('@')[0]).join(', ')}`,
    })),
  };
  const catalogWorkloads = workloadRefs.filter(
    ref => !wlSelected.some(w => `${w.namespace}/${w.name}` === ref),
  );
  if (cluster && catalogWorkloads.length)
    workloads.items.push(...catalogWorkloads.map(ref => ({ label: `${ref}`, detail: 'in the catalog, not running' })));

  // PVCs: claimed by the selected workloads, plus PVC manifests in the owned files.
  const claimNames = new Set<string>();
  for (const w of wlSelected) for (const c of w.claims) claimNames.add(`${w.namespace}/${c}`);
  const manifestPvcs = [...sourceFiles].flatMap(([path, text]) =>
    /\.ya?ml$/.test(path)
      ? listDocuments(text).filter(d => d.ref.kind === 'PersistentVolumeClaim').map(d => ({ path, ref: d.ref, doc: d.doc }))
      : [],
  );
  if (scope === 'app') for (const m of manifestPvcs) claimNames.add(`${m.ref.namespace ?? namespaces[0]}/${m.ref.name}`);
  const pvcFacts: PvcFact[] = (cluster?.pvcs ?? []).filter(p => claimNames.has(`${p.namespace}/${p.name}`));
  const pvcs: InventorySection = {
    id: 'volumes',
    title: 'Volumes (PVCs and bound volumes)',
    state: !cluster ? 'unavailable' : pvcFacts.length ? 'found' : 'none',
    note: clusterNote ?? 'Reclaim policy of the bound PVs is not readable by the portal; check it before deleting a claim.',
    items: pvcFacts.map(p => ({
      label: `PVC ${p.namespace}/${p.name}`,
      detail: `${p.size ?? '?'} on \`${p.storageClass ?? '?'}\`${p.volumeName ? `, volume \`${p.volumeName}\`` : ''}${p.storageClass === 'do-block-storage-retain' ? '' : '; NOT the retain class'}`,
    })),
  };
  if (!cluster && claimNames.size)
    pvcs.items.push(...[...claimNames].map(c => ({ label: `PVC ${c}`, detail: 'from manifests; live state unknown' })));

  const snapshots: InventorySection = {
    id: 'snapshots',
    title: 'Snapshots',
    state: 'unavailable',
    note: 'DigitalOcean snapshots (F6) and VolumeSnapshots are not readable by the portal. List them by hand before teardown and keep the final one for the retention period.',
    items: [],
  };

  // Ingress, DNS, tailnet.
  const ingressDocs = [...sourceFiles].flatMap(([path, text]) =>
    /\.ya?ml$/.test(path)
      ? listDocuments(text).filter(d => d.ref.kind === 'Ingress' && (scope === 'app' || matchesWorkload(d.ref, d.doc, workloadName!))).map(d => ({ path, ...d }))
      : [],
  );
  const ingress: InventorySection = {
    id: 'ingress',
    title: 'Ingress, DNS and tailnet names',
    state: ingressDocs.length ? 'found' : 'none',
    items: ingressDocs.map(d => {
      const cls = d.doc.spec?.ingressClassName;
      const hosts = ingressHosts(d.doc).join(', ');
      const isLive = cluster?.ingresses.some(i => i.namespace === (d.ref.namespace ?? '') && i.name === d.ref.name);
      return {
        label: `Ingress ${d.ref.namespace ?? ''}/${d.ref.name}`,
        detail: `${cls === 'tailscale' ? 'tailnet host' : 'public host'} ${hosts}; ${cluster ? (isLive ? 'live' : 'not live') : 'live state unknown'}`,
      };
    }),
    note: 'Public hosts need their Cloudflare record removed; tailnet hosts leave a stale Tailscale device to delete.',
  };

  // Secrets.
  const secretFiles = [...sourceFiles.keys()].filter(f => /sealedsecret|\.secret\./.test(f) && filesToDelete.includes(f));
  const argoSecrets = argoNames.flatMap(n => liveApp(n)?.resources.filter(r => r.kind === 'SealedSecret' || r.kind === 'Secret').map(r => `${r.kind} ${r.namespace ?? ''}/${r.name}`) ?? []);
  const secrets: InventorySection = {
    id: 'secrets',
    title: 'Secrets',
    state: secretFiles.length || argoSecrets.length ? 'found' : cluster ? 'none' : 'unavailable',
    note: 'The portal cannot list Secrets (no RBAC). Secrets created by hand outside Git survive a merge: check the namespace by hand. Rotate credentials that other systems still hold.',
    items: unique([...secretFiles.map(f => `SealedSecret file ${f}`), ...argoSecrets]).map(label => ({ label })),
  };

  // Backup sources: ReplicationSources in the shared VolSync file, local backup files.
  const pvcNamesInNs = new Set(pvcFacts.map(p => `${p.namespace}/${p.name}`));
  for (const k of claimNames) pvcNamesInNs.add(k);
  const volsyncText = await readOnce('platform/backups/volsync/sources.yaml');
  const backupItems: { label: string; detail?: string }[] = [];
  if (volsyncText !== undefined) {
    const res = removeDocuments(volsyncText, (ref, doc) =>
      ref.kind === 'ReplicationSource' && pvcNamesInNs.has(`${ref.namespace}/${doc.spec?.sourcePVC}`));
    if (res.removed.length) {
      for (const r of res.removed) backupItems.push({ label: `ReplicationSource ${r.namespace}/${r.name}`, detail: 'removed from platform/backups/volsync/sources.yaml by this PR (stops new backups; the restic repository and its sealed credentials are kept)' });
      if (res.text !== undefined) edit('platform/backups/volsync/sources.yaml', res.text);
    }
  }
  for (const path of filesToDelete.filter(f => sourceFiles.has(f))) {
    for (const d of listDocuments(sourceFiles.get(path) ?? '')) {
      if (d.ref.kind === 'ReplicationSource') backupItems.push({ label: `ReplicationSource ${d.ref.namespace}/${d.ref.name}`, detail: `defined in ${path}, deleted with the app` });
      if (d.ref.kind === 'CronJob' && /backup|dump/.test(d.ref.name)) backupItems.push({ label: `CronJob ${d.ref.namespace}/${d.ref.name}`, detail: `defined in ${path}, deleted with the app` });
    }
  }
  const cronHits = (cluster?.cronJobs ?? []).filter(c => namespaces.includes(c.namespace) && /backup|dump/.test(c.name) && !backupItems.some(b => b.label.endsWith(`${c.namespace}/${c.name}`)));
  for (const c of cronHits) backupItems.push({ label: `CronJob ${c.namespace}/${c.name}`, detail: 'live; defined outside this app (platform/backups): remove it there' });
  const sealedRestic = [...(await repo.readDir('platform/backups/volsync/sealed')).keys()].filter(f => namespaces.some(ns => f.split('/').pop()!.startsWith(`${ns}-volume-`)));
  for (const f of sealedRestic) kept.push({ path: f, reason: 'Restic repository credentials: keep until the retention period ends or the repository becomes unreadable' });
  const backups: InventorySection = {
    id: 'backups',
    title: 'Backup sources',
    state: backupItems.length ? 'found' : 'none',
    note: sealedRestic.length ? `Sealed restic credentials (${sealedRestic.length}) are kept for the retention period.` : 'Read from Git and live CronJobs; ReplicationSources themselves are not readable by the portal.',
    items: backupItems,
  };

  // Dashboards and alerts: text search of platform/monitoring.
  const monitoring = await repo.readDir('platform/monitoring');
  const words = scope === 'workload' ? [workloadName!] : unique([app.metadata.name, ...namespaces]);
  const hits = [...monitoring].filter(([p]) => /\.ya?ml$/.test(p)).map(([path, text]) => ({ path, lines: matchingLines(text, words) })).filter(h => h.lines.length);
  const alertHits = hits.filter(h => /alert|rule/i.test(h.path));
  const dashHits = hits.filter(h => !/alert|rule/i.test(h.path));
  const toItems = (hs: typeof hits) => hs.map(h => ({ label: h.path, detail: h.lines.join(' | ') }));
  const dashboards: InventorySection = {
    id: 'dashboards',
    title: 'Dashboards',
    state: dashHits.length ? 'found' : 'none',
    note: `Text search of platform/monitoring for ${words.join(', ')}. The generic workload dashboard is variable-driven and needs no change.`,
    items: toItems(dashHits),
  };
  const alerts: InventorySection = {
    id: 'alerts',
    title: 'Alerts',
    state: alertHits.length ? 'found' : 'none',
    note: 'Text search of alert rule files only. Live alerts, silences and VolSync series for removed sources are not readable here: check Alertmanager after teardown.',
    items: toItems(alertHits),
  };

  // Catalog.
  const catalogItems = [
    { label: refOf(target), detail: 'catalog entity' },
    ...(scope === 'app'
      ? [...subcomponents.map(e => ({ label: refOf(e), detail: 'workload entity' })), ...argoResources.map(e => ({ label: refOf(e), detail: 'Argo application resource' }))]
      : []),
  ];
  const catalogSection: InventorySection = {
    id: 'catalog',
    title: 'Catalog entries',
    state: 'found',
    items: catalogItems,
    note: 'Entity documents are removed from both catalog copies; the mapping inputs are edited; regenerate with `node scripts/curate-structure.cjs`.',
  };

  // --- catalog and test edits --------------------------------------------------------
  const removeNames = new Set(
    scope === 'app'
      ? [target.metadata.name, ...subcomponents.map(e => e.metadata.name), ...argoResources.map(e => e.metadata.name)]
      : [target.metadata.name],
  );
  const removeRefs = catalogItems.map(i => i.label);
  for (const path of [...PATHS.apps, ...PATHS.cluster]) {
    const text = await readOnce(path);
    if (text === undefined) { warnings.push(`${path} could not be read`); continue; }
    const res = removeDocuments(text, ref => ['Component', 'Resource'].includes(ref.kind) && removeNames.has(ref.name) && !(ref.kind === 'Component' && ref.name === app.metadata.name && scope === 'workload'));
    if (res.removed.length && res.text !== undefined) edit(path, res.text);
  }
  const mapText = await readOnce(PATHS.workloadMap);
  if (mapText !== undefined) {
    const keys = scope === 'app' ? [...argoNames, ...workloadRefs] : workloadRefs;
    const res = removeMapEntries(mapText, { ref: appRef, keys });
    let next = res.text;
    if (scope === 'app') next = removePrimaryArgo(next, app.metadata.name);
    if (next !== mapText) edit(PATHS.workloadMap, next);
  }
  if (scope === 'app') {
    const structure = await readOnce(PATHS.structure);
    if (structure !== undefined) {
      const next = removeStructureApp(structure, app.metadata.name);
      if (next !== structure) edit(PATHS.structure, next);
    }
    const exceptions = await readOnce(PATHS.coverageExceptions);
    if (exceptions !== undefined) {
      const next = removeCoverageExceptions(exceptions, appRef);
      if (next !== exceptions) edit(PATHS.coverageExceptions, next);
    }
    const baseline = await readOnce(PATHS.refBaseline);
    if (baseline !== undefined) {
      const next = removeRefBaseline(baseline, removeRefs);
      if (next === undefined) warnings.push(`${PATHS.refBaseline} is not in canonical JSON format; remove ${removeRefs.length} refs by hand.`);
      else if (next !== baseline) edit(PATHS.refBaseline, next);
    }
    const removedApplications = unique(targetApps.map(a => a.ref.name).filter(n => filesToDelete.some(f => targetApps.find(t => t.ref.name === n && t.file.path === f))));
    if (removedApplications.length) {
      const st = await readOnce(PATHS.structureTest);
      const ct = await readOnce(PATHS.coverageTest);
      const a = st && bumpStructureTestCount(st, -removedApplications.length);
      const b = ct && bumpCoverageTestCount(ct, -removedApplications.length);
      if (a && b) { edit(PATHS.structureTest, a); edit(PATHS.coverageTest, b); }
      else warnings.push(`Update the hard-coded Argo application count (-${removedApplications.length}) in scripts/test-catalog-structure.cjs and scripts/test-coverage.cjs by hand.`);
    }
  }

  // --- references that will turn CI red ----------------------------------------------
  const referenceHits: string[] = [];
  for (const path of ['platform/security/shared-namespaces.yaml', 'platform/ecosystem/chart-lock.json', 'platform/ecosystem/foundation-lock.json', 'platform/nats/network-policy.yaml', 'platform/kube/ai-operator.yaml', 'platform/hermes/mcp/coverage.yaml']) {
    const text = await readOnce(path);
    if (text === undefined) continue;
    const m = matchingLines(text, words, 3);
    if (m.length) referenceHits.push(`\`${path}\`: ${m.join(' | ')}`);
  }

  // --- retention and body ------------------------------------------------------------
  const retainUntil = new Date(now.getTime() + retention.days * 86400000).toISOString().slice(0, 10);
  if (scope === 'workload') {
    const changed = [...filesToDelete, ...files.keys()];
    const used = new Set(
      targetApps
        .filter(a => localPaths(a.doc).some(p => changed.some(c => c.startsWith(`${p}/`))))
        .map(a => a.ref.name),
    );
    if (used.size) argo.items = argo.items.filter(i => used.has(i.label));
  }
  const inventory = [workloads, argo, pvcs, snapshots, secrets, ingress, backups, dashboards, alerts, catalogSection];
  const title = target.metadata.title ?? target.metadata.name;
  const recordPath = `docs/operations/retirements/${target.metadata.name}-${today}.md`;
  const sharedNamespace = scope === 'app' && namespaces.some(ns => (cluster?.workloads ?? []).some(w => w.namespace === ns && !workloadRefs.includes(`${w.namespace}/${w.name}`)));
  const cascade = argoNames.some(n => (liveApp(n)?.finalizers ?? targetApps.find(t => t.ref.name === n)?.doc.metadata?.finalizers ?? []).includes('resources-finalizer.argocd.argoproj.io'));
  const parentNoPrune = targetApps.some(a => !a.file.path.startsWith('platform/argocd/'));

  const retentionText = retention.mode === 'keep-last-backup'
    ? `Keep the last backup and the restic repository until ${retainUntil} (${retention.days} days).`
    : retention.mode === 'export-first'
    ? `Export the data first, store the export off-cluster, then keep the last backup until ${retainUntil} (${retention.days} days).`
    : 'No retention: the owner confirmed the data is not needed. Take a final snapshot anyway if the volume is a retain-class volume you cannot recreate.';

  const before = [
    `- [ ] Retention: ${retentionText} Record the final backup id and checksum in \`${recordPath}\`. Never commit backup contents.`,
    '- [ ] **Regenerate the catalog** from `services/onemind-portal`: `node scripts/curate-structure.cjs`, then `node scripts/test-catalog-structure.cjs` and `node scripts/test-coverage.cjs`. Commit both catalog copies. The generated mapping Resources are rebuilt by the generator.',
    ...(referenceHits.length ? [`- [ ] These files still mention the app and will fail review or tests until cleaned: ${referenceHits.join('; ')}`] : []),
    '- [ ] `rg` the repository for the app name and each namespace; docs that are dated audits stay as written.',
  ];
  const mergeEffect = scope === 'workload'
    ? 'Merging removes the workload\'s manifests from the parent Application. With prune enabled Argo deletes the live objects on the next sync; with manual sync nothing changes until a person syncs.'
    : parentNoPrune
    ? 'The Application sits under a parent that does **not** prune: merging deletes nothing. The Application and its objects stay until deleted by hand.'
    : cascade
    ? 'The Application has `resources-finalizer`: once the parent app-of-apps prunes the Application CR, Argo cascades and deletes every object it tracks (Deployment, Service, PVC objects). StatefulSet claim templates survive.'
    : 'The Application has no `resources-finalizer`: merging deletes only the Application CR. Workloads, PVCs, Secrets and the namespace keep running unmanaged until the teardown below.';
  const workloadAfter = [
    '1. **Stop traffic first:** delete the Ingress (if any) and check nothing else calls the workload (`rg` for its Service DNS name).',
    '2. **Remove the live objects:** a person syncs the Application (never the portal) or deletes the objects by hand. Objects the repository no longer lists are removed only when prune is on.',
    '3. **Its claims:** a PVC the workload used stays until deleted by hand; check the reclaim policy of its volume first and keep it until the retention date.',
    `4. **Volumes and snapshots:** after ${retainUntil}, and with explicit owner approval, delete Released PVs, their DigitalOcean volumes and final snapshots.`,
    '5. **Names:** remove the stale Tailscale device and the Cloudflare DNS record for each host above.',
    '6. **Alerts and dashboards:** resolve stale alerts and silences; remove the rules and boards listed above.',
    `7. Tick this list in \`${recordPath}\`. The namespace and the other workloads of the app stay.`,
  ];
  const appAfter = [
    '1. **Stop traffic and writers:** delete the Ingress and scale writers to zero. Check no other app depends on it (`rg` for its Service DNS name).',
    '2. **Remove the backup source before the namespace:** the ReplicationSource is removed from Git by this PR; confirm it is gone live. Deleting the namespace first leaves VolSync series stale and the out-of-sync alert firing. Silence or resolve that alert, and expect the series to go absent after staleness.',
    '3. **Application:** if it still exists, delete it respecting the finalizer: cascade deletes its objects, `--cascade=orphan` leaves them. A `prune: false` parent never removes it.',
    `4. **Namespace:** ${sharedNamespace ? '**shared with other workloads: do not delete the namespace.**' : 'delete it only after step 1-3; PVs of the retain class go Released and keep their volumes.'}`,
    '5. **CRDs last:** only if this app owned them and nothing else uses them; deleting a CRD deletes every custom resource of that kind cluster-wide.',
    `6. **Volumes and snapshots:** after ${retainUntil}, and with explicit owner approval, delete Released PVs, their DigitalOcean volumes, the final snapshots and the restic repository, then the sealed restic credentials.`,
    '7. **Names:** remove the stale Tailscale device and the Cloudflare DNS record for each host above.',
    '8. **Alerts and dashboards:** resolve stale alerts and silences; remove the app\'s rules and boards listed above.',
    `9. Update \`docs/STACK.md\` counts and tick this list in \`${recordPath}\`.`,
  ];
  const after = scope === 'workload' ? workloadAfter : appAfter;

  const section = (s: InventorySection) => {
    const head = `**${s.title}:** ${s.state === 'unavailable' ? '_unavailable_' : s.state === 'none' ? '_none found_' : `${s.items.length} found`}`;
    const items = s.items.map(i => `  - ${i.label}${i.detail ? `: ${i.detail}` : ''}`).join('\n');
    return `${head}${s.note ? `\n  - _${s.note}_` : ''}${items ? `\n${items}` : ''}`;
  };
  const body = `## Retire ${scope === 'app' ? 'app' : 'workload'}: ${title}

Opened by the Backstage "Retire app" template as a draft. The portal read the catalog, Git and the cluster read-only and deleted nothing. ${scope === 'workload' ? `Scope is the single workload \`${workloadNamespace}/${workloadName}\` of \`${app.metadata.name}\`.` : ''}

**Reason:** ${input.reason}

**Retention:** ${retentionText}

**What merging does:** ${mergeEffect}

### Attached resources

${inventory.map(section).join('\n\n')}

### Files in this PR

- Deleted (${filesToDelete.length}): ${filesToDelete.length ? filesToDelete.map(f => `\`${f}\``).join(', ') : 'none'}
- Edited (${[...files.keys()].length}): ${[...files.keys()].map(f => `\`${f}\``).join(', ') || 'none'}
- Kept on purpose (${kept.length}): ${kept.map(k => `\`${k.path}\` (${k.reason})`).join('; ') || 'none'}
${warnings.length ? `\n### Needs a person\n\n${warnings.map(w => `- ${w}`).join('\n')}\n` : ''}
### Before merging

${before.join('\n')}

### After merging (manual, in this order)

${after.join('\n')}
`;
  files.set(recordPath, `# ${title}: retirement record, ${today}\n\n${body}`);
  if (scope === 'app') {
    const stack = await readOnce(PATHS.stack);
    if (stack !== undefined)
      edit(PATHS.stack, appendSection(stack, `## ${title}: retirement proposed — ${today}\n\nSource removal opened by the Backstage "Retire app" template; live teardown and data deletion follow [${recordPath}](../${recordPath.replace(/^docs\//, '')}). Not retired until that checklist is complete.`));
  }
  return {
    scope,
    ref: refOf(target),
    name: target.metadata.name,
    title,
    namespaces,
    inventory,
    filesToDelete: unique(filesToDelete).sort(),
    files,
    kept,
    warnings,
    body,
    recordPath,
  };
}
