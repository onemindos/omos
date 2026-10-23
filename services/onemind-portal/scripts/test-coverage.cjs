// Tests for the F5 coverage contract generator (curate-coverage.cjs) and the checked-in catalog.
const { test } = require('node:test');
const assert = require('node:assert/strict');
const fs = require('node:fs');
const path = require('node:path');
const YAML = require('yaml');
const {
  applyCoverage, deriveSelector, grafanaPath, logsQuery, dashboardSelector, renderRunbooks, NAMESPACE_SCOPE_SELECTOR,
} = require('./curate-coverage.cjs');
const { runbooksFor, loadSources, load } = require('./curate-structure.cjs');
const { auditCoverage, loadCoverageExceptions } = require('./audit-catalog.cjs');

const root = path.join(__dirname, '..');
// Repo-level platform/ is absent in the image build context (only services/onemind-portal is copied).
const hasRepo = fs.existsSync(path.join(__dirname, '../../../platform/backstage/catalog'));
const repoOnly = hasRepo ? {} : { skip: 'repo-level platform/ is not in this build context' };
const coverage = YAML.parse(fs.readFileSync(path.join(root, 'catalog/coverage.yaml'), 'utf8'));

const workload = (namespace, name, selector, state = 'observed') => ({
  kind: 'Component',
  metadata: { name: `k8s-${namespace}-${name}`, annotations: {
    'backstage.io/kubernetes-namespace': namespace,
    'backstage.io/kubernetes-label-selector': selector,
    'onemind.io/kubernetes-workload': name,
    'onemind.io/inventory-state': state,
  } },
});
const app = (name, annotations = {}) => ({ kind: 'Component', metadata: { name, title: name, annotations }, spec: { system: 'infrastructure' } });
function model(apps, subs, argo = {}) {
  const appComponent = new Map(apps.map(a => [a.metadata.name, a]));
  const subcomponents = new Map(Object.entries(subs));
  const workloadNamespaces = new Map();
  for (const [n, list] of subcomponents) for (const w of list) {
    const ns = w.metadata.annotations['backstage.io/kubernetes-namespace'];
    if (!workloadNamespaces.has(ns)) workloadNamespaces.set(ns, new Set());
    workloadNamespaces.get(ns).add(n);
  }
  return { appComponent, subcomponents, argoByApp: new Map(Object.entries(argo)), workloadNamespaces };
}

test('selector: shared equality pairs, else a set over the best common key, else none', () => {
  assert.equal(deriveSelector([workload('a', 'x', 'app=x,tier=web'), workload('a', 'y', 'app=y,tier=web')]), 'tier=web');
  assert.equal(deriveSelector([workload('a', 'x', 'app=x'), workload('a', 'y', 'app=y')]), 'app in (x,y)');
  assert.equal(deriveSelector([workload('a', 'x', 'app.kubernetes.io/name=x,app=x'), workload('a', 'y', 'app.kubernetes.io/name=y,app=y')]), 'app.kubernetes.io/name in (x,y)');
  assert.equal(deriveSelector([workload('a', 'x', 'app=x'), workload('a', 'y', 'name=y')]), '');
  assert.equal(deriveSelector([workload('a', 'x', 'app=x'), workload('a', 'y', '')]), '');
  assert.equal(deriveSelector([]), '');
});

test('every annotation of the contract is generated from the workloads, none typed per app', () => {
  const a = app('demo', { 'argocd/app-name': 'platform-demo' });
  const m = model([a], { demo: [workload('demo', 'web', 'app=web'), workload('demo', 'db', 'app=db')] }, { demo: ['platform-demo'] });
  applyCoverage(m, coverage);
  const ann = a.metadata.annotations;
  assert.equal(ann['backstage.io/kubernetes-namespace'], 'demo');
  assert.equal(ann['backstage.io/kubernetes-label-selector'], 'app in (db,web)');
  assert.equal(ann['grafana/dashboard-selector'], "uid == 'onemind-workload'");
  assert.equal(ann['grafana/alert-label-selector'], 'namespace=demo');
  assert.equal(ann['onemind.io/grafana-dashboard-path'], 'd/onemind-workload/onemind-workload?var-namespace=demo');
  assert.equal(ann['onemind.io/logs-query'], '{k8s_namespace_name="demo"}');
  assert.equal(ann['onemind.io/recovery-runbook'], 'apps/demo/');
  assert.equal(ann['github.com/project-slug'], 'onemindos/omos');
  assert.equal(ann['onemind.io/kubernetes-namespaces'], 'demo');
  assert.equal(ann['onemind.io/kubernetes-workloads'], 'demo/db,demo/web');
});

test('scoped Prometheus alert rules are recorded only for the apps that have them', () => {
  const nats = app('nats'), other = app('other');
  applyCoverage(model([nats, other], { nats: [workload('nats', 'nats', 'app=nats')], other: [workload('o', 'o', 'app=o')] }), coverage);
  assert.equal(nats.metadata.annotations['onemind.io/prometheus-scopes'], 'nats');
  assert.equal(other.metadata.annotations['onemind.io/prometheus-scopes'], undefined);
});

test('an app sharing a namespace is narrowed to its own workloads in the dashboard link and the log query', () => {
  const a = app('one'), b = app('two');
  const m = model([a, b], { one: [workload('shared', 'alpha', 'app=alpha')], two: [workload('shared', 'beta', 'app=beta')] });
  applyCoverage(m, coverage);
  assert.equal(a.metadata.annotations['onemind.io/grafana-dashboard-path'], 'd/onemind-workload/onemind-workload?var-namespace=shared&var-workload=alpha');
  assert.equal(a.metadata.annotations['onemind.io/logs-query'], '{k8s_namespace_name="shared", k8s_pod_name=~"(alpha)-.*"}');
  assert.equal(b.metadata.annotations['onemind.io/grafana-dashboard-path'], 'd/onemind-workload/onemind-workload?var-namespace=shared&var-workload=beta');
});

test('reviewed values are never replaced; the generated namespace fallback is re-derived', () => {
  const reviewed = app('r', { 'backstage.io/kubernetes-label-selector': 'custom=1', 'grafana/dashboard-selector': "uid == 'mine'", 'grafana/alert-label-selector': 'team=x' });
  const mixed = app('m');
  const m = model([reviewed, mixed], { r: [workload('r', 'a', 'app=a')], m: [workload('m', 'a', 'app=a'), workload('m', 'b', 'name=b')] });
  applyCoverage(m, coverage);
  assert.equal(reviewed.metadata.annotations['backstage.io/kubernetes-label-selector'], 'custom=1');
  assert.equal(reviewed.metadata.annotations['grafana/dashboard-selector'], "uid == 'mine'");
  assert.equal(reviewed.metadata.annotations['grafana/alert-label-selector'], 'team=x');
  assert.equal(mixed.metadata.annotations['backstage.io/kubernetes-label-selector'], NAMESPACE_SCOPE_SELECTOR);
  assert.equal(mixed.metadata.annotations['onemind.io/kubernetes-scope'], 'namespace');
  // The workloads gain a shared label later: the fallback is replaced and the scope note removed.
  m.subcomponents.set('m', [workload('m', 'a', 'app=a,tier=x'), workload('m', 'b', 'name=b,tier=x')]);
  applyCoverage(m, coverage);
  assert.equal(mixed.metadata.annotations['backstage.io/kubernetes-label-selector'], 'tier=x');
  assert.equal(mixed.metadata.annotations['onemind.io/kubernetes-scope'], undefined);
});

test('no namespace is guessed: an app without live workloads gets the namespace-free annotations only', () => {
  const a = app('daemon', { 'argocd/app-name': 'platform-daemon' });
  const gone = app('gone');
  const m = model([a, gone], { gone: [workload('g', 'x', 'app=x', 'missing-review-required')] }, { daemon: ['platform-daemon'] });
  applyCoverage(m, coverage);
  for (const x of [a, gone]) {
    const ann = x.metadata.annotations;
    assert.equal(ann['backstage.io/kubernetes-namespace'], undefined);
    assert.equal(ann['grafana/alert-label-selector'], undefined);
    assert.equal(ann['onemind.io/logs-query'], undefined);
    assert.equal(ann['onemind.io/recovery-runbook'], `apps/${x.metadata.name}/`);
    assert.equal(ann['grafana/dashboard-selector'], "uid == 'onemind-workload'");
  }
});

test('multi-namespace apps use a regex logs query and one repeated variable per namespace', () => {
  assert.equal(logsQuery(coverage, ['a', 'b'], []), '{k8s_namespace_name=~"a|b"}');
  assert.equal(grafanaPath('onemind-workload', ['a', 'b'], []), 'd/onemind-workload/onemind-workload?var-namespace=a&var-namespace=b');
  assert.equal(dashboardSelector(coverage, 'nats'), "uid == 'onemind-nats' || uid == 'onemind-workload'");
});

test('an app that is not on the cluster (no Argo application, no workloads) is left alone', () => {
  const saas = app('zoho');
  applyCoverage(model([saas], {}), coverage);
  assert.deepEqual(saas.metadata.annotations, {});
});

test('runbook pages: one per covered app plus an index, facts only from the catalog', () => {
  const a = app('demo', { 'argocd/app-name': 'platform-demo' });
  a.metadata.description = 'Demo service.';
  const m = model([a], { demo: [workload('demo', 'web', 'app=web')] }, { demo: ['platform-demo'] });
  applyCoverage(m, coverage);
  const files = renderRunbooks(m, coverage);
  assert.deepEqual([...files.keys()].sort(), ['apps/demo.md', 'apps/index.md']);
  const page = files.get('apps/demo.md');
  assert.match(page, /platform-demo/);
  assert.match(page, /`web` \| Deployment \| `demo`/);
  assert.match(page, /Do not `kubectl apply`/);
  assert.doesNotMatch(page, /is healthy|are healthy/i);
  assert.match(files.get('apps/index.md'), /\[demo\]\(demo\.md\)/);
});

// --- the checked-in catalog ----------------------------------------------------------------
const entities = ['apps.yaml', 'cluster.yaml'].flatMap(f => load(path.join(root, 'catalog', f)));

test('catalog coverage: every Kubernetes app has the full contract or a reviewed exception', () => {
  const report = auditCoverage(entities, loadCoverageExceptions());
  assert.deepEqual(report.missing, []);
  assert.deepEqual(report.invalidExceptions, []);
  assert.deepEqual(report.staleExceptions, []);
  assert.ok(report.kubernetesApps >= 55, `${report.kubernetesApps} Kubernetes apps`);
  // All 55 Argo applications are linked to an app (F2); the per-app coverage is reported, not hidden.
  assert.equal(report.argoApplicationsLinked, 55);
  assert.equal(report.byCapability['grafana-dashboard'].covered, report.kubernetesApps);
  assert.equal(report.byCapability.runbook.covered, report.kubernetesApps);
});

test('catalog coverage: removing a generated annotation from any one app turns the guard red', () => {
  const exceptions = loadCoverageExceptions();
  const candidates = entities.filter(e => e.metadata.annotations?.['onemind.io/logs-query']);
  assert.ok(candidates.length > 10);
  for (const victim of candidates.slice(0, 5)) {
    const copy = structuredClone(entities);
    const target = copy.find(e => e.metadata.name === victim.metadata.name && e.kind === 'Component');
    delete target.metadata.annotations['onemind.io/logs-query'];
    const report = auditCoverage(copy, exceptions);
    assert.deepEqual(report.missing, [{ app: `component:default/${victim.metadata.name}`, capability: 'logs' }]);
  }
});

test('catalog coverage: each exception has a reason naming a verifiable fact', () => {
  for (const item of loadCoverageExceptions()) assert.ok(item.reason.length >= 60, `${item.app} ${item.capability}`);
});

test('every runbook annotation resolves to a generated page that matches the catalog byte for byte', () => {
  const sources = loadSources();
  const expected = runbooksFor({ apps: load(path.join(root, 'catalog/apps.yaml')), cluster: load(path.join(root, 'catalog/cluster.yaml')) }, sources);
  assert.ok(expected.size > 55);
  for (const [rel, text] of expected) assert.equal(fs.readFileSync(path.join(root, 'runbooks/docs', rel), 'utf8'), text, rel);
  for (const e of entities) {
    const value = e.metadata.annotations?.['onemind.io/recovery-runbook'];
    if (value) assert.ok(fs.existsSync(path.join(root, 'runbooks/docs', value.replace(/\/$/, '.md'))), value);
  }
  const files = fs.readdirSync(path.join(root, 'runbooks/docs/apps'));
  assert.equal(files.length, expected.size, 'no page for an app that no longer exists');
});

test('the app runbooks are part of the runbook site navigation', () => {
  assert.match(fs.readFileSync(path.join(root, 'runbooks/mkdocs.yml'), 'utf8'), /apps\/index\.md/);
});

test('the generic dashboard named in coverage.yaml is the one GitOps provisions', repoOnly, () => {
  const cm = YAML.parse(fs.readFileSync(path.join(__dirname, '../../../platform/monitoring/grafana-dashboard-onemind-workload.yaml'), 'utf8'));
  const board = JSON.parse(cm.data['onemind-workload.json']);
  assert.equal(board.uid, coverage.genericDashboard.uid);
  assert.deepEqual(board.templating.list.map(v => v.name), ['namespace', 'workload']);
  for (const e of entities) {
    const link = e.metadata.annotations?.['onemind.io/grafana-dashboard-path'];
    if (link) assert.match(link, new RegExp(`^d/${board.uid}/${board.uid}\\?var-namespace=`));
  }
});

test('the mounted catalog copy matches the portal copy', repoOnly, () => {
  const mounted = f => fs.readFileSync(path.join(__dirname, '../../../platform/backstage/catalog', f), 'utf8');
  assert.equal(mounted('apps.yaml'), fs.readFileSync(path.join(root, 'catalog/apps.yaml'), 'utf8'));
  assert.equal(mounted('cluster.yaml'), fs.readFileSync(path.join(root, 'catalog/cluster.yaml'), 'utf8'));
});
