import type { Entity } from '@backstage/catalog-model';
import YAML from 'yaml';
import { buildRetirePlan } from './retire';
import { snapshotReader } from './sources';
import { catalogEntities, localRepo, localSnapshot, stacCluster } from './testRepo';

const NOW = new Date('2026-11-05T12:00:00Z');

/** The real catalog has no stac-browser workload entity yet; add the one the generator would make. */
const stacBrowser = (): Entity => ({
  apiVersion: 'backstage.io/v1alpha1',
  kind: 'Component',
  metadata: {
    name: 'k8s-stac-stac-browser-fixture',
    title: 'STAC — Browser',
    annotations: {
      'backstage.io/kubernetes-namespace': 'stac',
      'onemind.io/kubernetes-workload': 'stac-browser',
    },
    labels: { 'onemind.io/catalog-role': 'workload' },
  },
  spec: {
    type: 'service',
    lifecycle: 'production',
    owner: 'group:default/onemind-admins',
    subcomponentOf: 'component:default/stac',
  },
});

describe('Retire app on stac-browser (real repository files)', () => {
  it('lists its workload, Argo app, links and alerts, and removes only its own files', async () => {
    const catalog = [
      ...catalogEntities(
        'services/onemind-portal/catalog/apps.yaml',
        'services/onemind-portal/catalog/cluster.yaml',
      ),
      stacBrowser(),
    ];
    const plan = await buildRetirePlan({
      target: stacBrowser(),
      catalog,
      repo: localRepo(),
      cluster: stacCluster(),
      retention: { mode: 'keep-last-backup', days: 30 },
      reason: 'The STAC browser is replaced by the portal map view',
      now: NOW,
    });

    expect(plan.scope).toBe('workload');
    // Exactly the file that holds only stac-browser objects is deleted.
    expect(plan.filesToDelete).toEqual(['platform/stac/browser.yaml']);
    // The shared ingress file loses only the stac-browser Ingress; the record is new.
    expect([...plan.files.keys()].sort()).toEqual(
      ['platform/stac/ingress.yaml', plan.recordPath].sort(),
    );
    const ingress = plan.files.get('platform/stac/ingress.yaml')!;
    const names = YAML.parseAllDocuments(ingress).map(d => d.toJSON().metadata.name);
    expect(names).toEqual(['stac']);
    // Nothing belonging to the api, writer or sync is touched.
    for (const path of [...plan.filesToDelete, ...plan.files.keys()])
      expect(path).not.toMatch(/api|webodm|migrate|config|stac-db|stac-spaces/);
    // The NetworkPolicy that selects the browser is flagged, not edited.
    expect(plan.warnings.join('\n')).toMatch(/network\.yaml.*NetworkPolicy\/browser/);

    const section = (id: string) => plan.inventory.find(s => s.id === id)!;
    expect(section('workloads').items.map(i => i.label)).toEqual(['Deployment stac/stac-browser']);
    expect(section('argo').items[0].label).toBe('platform-stac');
    expect(section('argo').items[0].detail).toMatch(/resources-finalizer/);
    expect(section('argo').items[0].link).toContain('/applications/argocd/platform-stac');
    expect(section('ingress').items.map(i => i.label)).toEqual(['Ingress stac/stac-browser']);
    expect(section('ingress').items[0].detail).toMatch(/stac-browser/);
    // No DigitalOcean or snapshot data: unavailable, never an empty "none".
    expect(section('snapshots').state).toBe('unavailable');
    expect(section('alerts').title).toBe('Alerts');
    expect(plan.body).toMatch(/Keep the last backup .* until 2026-12-05/);
    expect(plan.body).toMatch(/After merging \(manual, in this order\)/);
    expect(plan.body).toMatch(/The namespace and the other workloads of the app stay/);
    expect(section('argo').items.map(i => i.label)).toEqual(['platform-stac']);
  });
});

// A small synthetic repository for the app-wide case, so every edit can be asserted exactly.
const files = () =>
  new Map<string, string>([
    ['platform/argocd/platform-demo.yaml', 'apiVersion: argoproj.io/v1alpha1\nkind: Application\nmetadata:\n  name: platform-demo\n  namespace: argocd\n  finalizers:\n    - resources-finalizer.argocd.argoproj.io\nspec:\n  source:\n    repoURL: https://github.com/onemindos/omos.git\n    path: platform/demo\n  destination:\n    namespace: demo\n'],
    ['platform/argocd/platform-other.yaml', 'apiVersion: argoproj.io/v1alpha1\nkind: Application\nmetadata:\n  name: platform-other\n  namespace: argocd\nspec:\n  source:\n    repoURL: https://github.com/onemindos/omos.git\n    path: platform/other\n'],
    ['platform/demo/deployment.yaml', 'apiVersion: apps/v1\nkind: Deployment\nmetadata:\n  name: demo\n  namespace: demo\n'],
    ['platform/demo/data-pvc.yaml', 'apiVersion: v1\nkind: PersistentVolumeClaim\nmetadata:\n  name: data\n  namespace: demo\n'],
    ['platform/demo/demo.sealedsecret.yaml', 'apiVersion: bitnami.com/v1alpha1\nkind: SealedSecret\nmetadata:\n  name: demo-secrets\n  namespace: demo\n'],
    ['platform/demo/backup-restic.sealedsecret.yaml', 'apiVersion: bitnami.com/v1alpha1\nkind: SealedSecret\nmetadata:\n  name: demo-data-restic\n  namespace: demo\n'],
    ['platform/demo/README.md', '# demo\n'],
    ['platform/other/deployment.yaml', 'apiVersion: apps/v1\nkind: Deployment\nmetadata:\n  name: other\n  namespace: other\n'],
    ['platform/backups/volsync/sources.yaml', 'apiVersion: volsync.backube/v1alpha1\nkind: ReplicationSource\nmetadata:\n  name: volume-aaa\n  namespace: demo\nspec:\n  sourcePVC: data\n---\napiVersion: volsync.backube/v1alpha1\nkind: ReplicationSource\nmetadata:\n  name: volume-bbb\n  namespace: other\nspec:\n  sourcePVC: data\n'],
    ['platform/monitoring/alerts.yaml', 'groups:\n- rules:\n  - alert: DemoDown\n    expr: up{namespace="demo"} == 0\n'],
    ['services/onemind-portal/catalog/apps.yaml', 'apiVersion: backstage.io/v1alpha1\nkind: Component\nmetadata:\n  name: demo\n---\napiVersion: backstage.io/v1alpha1\nkind: Component\nmetadata:\n  name: other\n'],
    ['platform/backstage/catalog/apps.yaml', 'apiVersion: backstage.io/v1alpha1\nkind: Component\nmetadata:\n  name: demo\n---\napiVersion: backstage.io/v1alpha1\nkind: Component\nmetadata:\n  name: other\n'],
    ['services/onemind-portal/catalog/cluster.yaml', 'apiVersion: backstage.io/v1alpha1\nkind: Component\nmetadata:\n  name: k8s-demo-demo-1\nspec:\n  subcomponentOf: component:default/demo\n---\napiVersion: backstage.io/v1alpha1\nkind: Resource\nmetadata:\n  name: argocd-app-platform-demo\n  annotations:\n    argocd/app-name: platform-demo\nspec:\n  type: argocd-application\n'],
    ['platform/backstage/catalog/cluster.yaml', 'apiVersion: backstage.io/v1alpha1\nkind: Component\nmetadata:\n  name: k8s-demo-demo-1\nspec:\n  subcomponentOf: component:default/demo\n---\napiVersion: backstage.io/v1alpha1\nkind: Resource\nmetadata:\n  name: argocd-app-platform-demo\n  annotations:\n    argocd/app-name: platform-demo\nspec:\n  type: argocd-application\n'],
    ['services/onemind-portal/catalog/workload-map.yaml', 'workloads:\n  demo/demo: component:default/demo\n  other/other: component:default/other\nargoApplications:\n  platform-demo: component:default/demo\n  platform-other: component:default/other\nprimaryArgoApplication:\n  demo: platform-demo\n'],
    ['services/onemind-portal/catalog/structure.yaml', 'apps:\n  - name: demo\n    title: Demo\n    system: s\n    description: d\n  - name: other\n    title: Other\n    system: s\n    description: d\nplacement:\n  demo: s\ntitles: {}\n'],
    ['services/onemind-portal/catalog/ref-baseline.json', `${JSON.stringify({ refs: ['component:default/demo', 'component:default/other'] }, null, 2)}\n`],
    ['services/onemind-portal/scripts/test-catalog-structure.cjs', "  assert.equal(inGit.size, 65, 'platform/ should define the 65 source Argo applications');\n"],
    ['services/onemind-portal/scripts/test-coverage.cjs', '  // All 65 Argo applications are linked to an app (F2); the per-app coverage is reported, not hidden.\n  assert.equal(report.argoApplicationsLinked, 65);\n'],
    ['docs/STACK.md', '# stack\n'],
  ]);

const demo = (): Entity => ({
  apiVersion: 'backstage.io/v1alpha1',
  kind: 'Component',
  metadata: {
    name: 'demo',
    title: 'Demo',
    annotations: {
      'argocd/app-name': 'platform-demo',
      'onemind.io/kubernetes-workloads': 'demo/demo',
      'onemind.io/kubernetes-namespaces': 'demo',
    },
    labels: { 'onemind.io/catalog-role': 'app' },
  },
  spec: { type: 'service', owner: 'group:default/onemind-admins' },
});

describe('Retire an app with its own namespace (synthetic repository)', () => {
  it('removes the Application, its files, its catalog entries and its backup source, and keeps the rest', async () => {
    const plan = await buildRetirePlan({
      target: demo(),
      catalog: [demo(), ...catalogEntities('services/onemind-portal/catalog/apps.yaml').slice(0, 0)],
      repo: snapshotReader(files()),
      cluster: undefined,
      retention: { mode: 'export-first', days: 14 },
      reason: 'Replaced and no longer used by anyone',
      now: NOW,
    });
    expect(plan.filesToDelete).toEqual([
      'platform/argocd/platform-demo.yaml',
      'platform/demo/data-pvc.yaml',
      'platform/demo/demo.sealedsecret.yaml',
      'platform/demo/deployment.yaml',
    ]);
    // The restic credentials and the README survive.
    expect(plan.kept.map(k => k.path).sort()).toEqual([
      'platform/demo/README.md',
      'platform/demo/backup-restic.sealedsecret.yaml',
    ]);
    // The other app is untouched everywhere.
    for (const text of plan.files.values()) expect(text).not.toMatch(/platform\/other\/x/);
    const edited = (path: string) => plan.files.get(path)!;
    expect(edited('platform/backups/volsync/sources.yaml')).not.toMatch(/volume-aaa/);
    expect(edited('platform/backups/volsync/sources.yaml')).toMatch(/volume-bbb/);
    expect(edited('services/onemind-portal/catalog/workload-map.yaml')).toBe(
      'workloads:\n  other/other: component:default/other\nargoApplications:\n  platform-other: component:default/other\nprimaryArgoApplication:\n',
    );
    expect(edited('services/onemind-portal/catalog/structure.yaml')).not.toMatch(/demo/);
    expect(edited('services/onemind-portal/catalog/structure.yaml')).toMatch(/name: other/);
    for (const apps of ['services/onemind-portal/catalog/apps.yaml', 'platform/backstage/catalog/apps.yaml'])
      expect(YAML.parseAllDocuments(edited(apps)).map(d => d.toJSON().metadata.name)).toEqual(['other']);
    // cluster.yaml: both the workload entity and the Argo resource go; nothing is left to write.
    expect(plan.files.has('services/onemind-portal/catalog/cluster.yaml')).toBe(false);
    expect(edited('services/onemind-portal/catalog/ref-baseline.json')).toBe(
      `${JSON.stringify({ refs: ['component:default/other'] }, null, 2)}\n`,
    );
    // The portal tests hard-code the Argo application count: 65 goes to 64.
    expect(edited('services/onemind-portal/scripts/test-catalog-structure.cjs')).toMatch(/inGit\.size, 64, .*the 64 source/);
    expect(edited('services/onemind-portal/scripts/test-coverage.cjs')).toMatch(/All 64 Argo applications[\s\S]*argoApplicationsLinked, 64/);
    expect(edited('docs/STACK.md')).toMatch(/Demo: retirement proposed/);
    // Lessons in the body: finalizer cascade, deletion order, VolSync metrics, retention date.
    expect(plan.body).toMatch(/has `resources-finalizer`/);
    expect(plan.body).toMatch(/Deleting the namespace first leaves VolSync series stale/);
    expect(plan.body).toMatch(/CRDs last/);
    expect(plan.body).toMatch(/until 2026-11-19/);
    expect(plan.inventory.find(s => s.id === 'workloads')!.state).toBe('unavailable');
  });

  it('never deletes a source directory another Application also uses', async () => {
    const repo = files();
    repo.set(
      'platform/argocd/platform-other.yaml',
      repo.get('platform/argocd/platform-other.yaml')!.replace('platform/other', 'platform/demo'),
    );
    const plan = await buildRetirePlan({
      target: demo(),
      catalog: [demo()],
      repo: snapshotReader(repo),
      retention: { mode: 'keep-last-backup', days: 30 },
      reason: 'Replaced and no longer used by anyone',
      now: NOW,
    });
    expect(plan.filesToDelete).toEqual(['platform/argocd/platform-demo.yaml']);
    expect(plan.warnings.join('\n')).toMatch(/shared with another Application/);
  });

  it('refuses to delete protected platform directories wholesale', async () => {
    const repo = files();
    repo.set('platform/argocd/platform-demo.yaml', repo.get('platform/argocd/platform-demo.yaml')!.replace('platform/demo', 'platform/security'));
    repo.set('platform/security/policy.yaml', 'apiVersion: v1\nkind: ConfigMap\nmetadata:\n  name: p\n');
    const plan = await buildRetirePlan({
      target: demo(),
      catalog: [demo()],
      repo: snapshotReader(repo),
      retention: { mode: 'keep-last-backup', days: 30 },
      reason: 'Replaced and no longer used by anyone',
      now: NOW,
    });
    expect(plan.filesToDelete).not.toContain('platform/security/policy.yaml');
    expect(plan.warnings.join('\n')).toMatch(/protected platform directory/);
  });

  it('reads the real repository without throwing for a real app', async () => {
    const catalog = catalogEntities('services/onemind-portal/catalog/apps.yaml', 'services/onemind-portal/catalog/cluster.yaml');
    const target = catalog.find(e => e.metadata.name === 'karakeep')!;
    expect(target).toBeDefined();
    const plan = await buildRetirePlan({
      target,
      catalog,
      repo: localRepo(),
      retention: { mode: 'keep-last-backup', days: 30 },
      reason: 'A dry run of the retirement plan',
      now: NOW,
    });
    expect(plan.filesToDelete).toContain('platform/argocd/platform-karakeep.yaml');
    expect(plan.filesToDelete.some(f => f.startsWith('platform/karakeep/'))).toBe(true);
    expect(localSnapshot().size).toBeGreaterThan(100);
  });
});
