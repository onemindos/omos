// Guard tests for the SO / LE / HP / GE catalog structure (F2). Reads the checked-in catalog
// files, so a hand edit or a regenerated inventory that breaks the model fails here.
const { test } = require('node:test');
const assert = require('node:assert/strict');
const fs = require('node:fs');
const path = require('node:path');
const YAML = require('yaml');
const { applyStructure, loadSources, load, refOf, UNASSIGNED } = require('./curate-structure.cjs');

const root = path.join(__dirname, '..');
const read = f => load(path.join(root, 'catalog', f));
const apps = read('apps.yaml');
const cluster = read('cluster.yaml');
const all = [...apps, ...cluster];
const sources = loadSources();
const org = read('org.yaml');
const byRef = new Map([...all, ...org].map(e => [refOf(e), e]));
const PILLARS = ['so', 'le', 'hp', 'ge'];
const systems = all.filter(e => e.kind === 'System');
const domains = all.filter(e => e.kind === 'Domain');
const members = all.filter(e => !['System', 'Domain'].includes(e.kind));
const workloads = cluster.filter(e => e.kind === 'Component' && e.metadata.annotations?.['onemind.io/inventory-source'] === 'kubernetes');
const ann = (e, k) => e.metadata.annotations?.[k];

// Applications defined in Git (parked and chart templates are not deployed).
function argoApplicationsInGit() {
  const names = new Set();
  const skip = new Set(['node_modules', 'charts', 'chart', 'parked', 'templates']);
  const walk = dir => {
    for (const entry of fs.readdirSync(dir, { withFileTypes: true })) {
      const p = path.join(dir, entry.name);
      if (entry.isDirectory()) { if (!skip.has(entry.name)) walk(p); continue; }
      if (!/\.ya?ml$/.test(entry.name)) continue;
      let docs;
      try { docs = YAML.parseAllDocuments(fs.readFileSync(p, 'utf8'), { logLevel: 'silent' }).map(d => d.toJS()); } catch { continue; }
      for (const d of docs) if (d?.kind === 'Application' && d.apiVersion === 'argoproj.io/v1alpha1') names.add(d.metadata.name);
    }
  };
  walk(path.join(root, '../../platform'));
  return names;
}

test('five Domains: OneMind parent and SO, LE, HP, GE children owned by the admins group', () => {
  const byName = Object.fromEntries(domains.map(d => [d.metadata.name, d]));
  assert.deepEqual(Object.keys(byName).sort(), ['onemind', 'onemind-ge', 'onemind-hp', 'onemind-le', 'onemind-so']);
  assert.equal(byName['onemind-so'].metadata.title, 'SO — Sovereign Operations');
  assert.equal(byName['onemind-le'].metadata.title, 'LE — Legacy Evolution');
  assert.equal(byName['onemind-hp'].metadata.title, 'HP — Holistic Performance');
  assert.equal(byName['onemind-ge'].metadata.title, 'GE — Generational Empire');
  assert.equal(byName.onemind.spec.subdomainOf, undefined);
  for (const p of PILLARS) {
    const d = byName[`onemind-${p}`];
    assert.equal(d.spec.subdomainOf, 'domain:default/onemind');
    assert.equal(d.metadata.labels['onemind.io/pillar'], p);
  }
  for (const d of domains) {
    assert.equal(d.spec.owner, 'group:default/onemind-admins');
    assert.ok(d.metadata.description);
  }
});

test('at most 25 Systems, each with title, description, owner and an existing Domain', () => {
  assert.ok(systems.length <= 25, `${systems.length} Systems`);
  assert.equal(new Set(systems.map(refOf)).size, systems.length);
  for (const s of systems) {
    assert.ok(s.metadata.title && s.metadata.description, `${s.metadata.name} needs title and description`);
    assert.ok(s.spec.owner);
    assert.ok(byRef.get(s.spec.domain)?.kind === 'Domain', `${s.metadata.name} domain ${s.spec.domain}`);
  }
  assert.equal(cluster.filter(e => e.kind === 'System').length, 0, 'Systems are defined once, in apps.yaml');
});

test('every Component, API and Resource has an existing System and a human title', () => {
  for (const e of members) {
    assert.ok(byRef.get(`system:default/${e.spec.system}`), `${refOf(e)} has no System (${e.spec.system})`);
    const title = e.metadata.title;
    assert.ok(title && title.trim(), `${refOf(e)} has no title`);
    assert.ok(!/-[0-9a-f]{8}$/.test(title), `${refOf(e)} title looks generated: ${title}`);
  }
});

test('no duplicate titles inside a System', () => {
  const seen = new Map();
  for (const e of members) {
    const key = `${e.spec.system}\u0000${e.metadata.title}`;
    assert.ok(!seen.has(key), `${e.spec.system}: "${e.metadata.title}" is ${seen.get(key)} and ${refOf(e)}`);
    seen.set(key, refOf(e));
  }
});

test('every workload is a subcomponent of its mapped app or sits in Unassigned workloads', () => {
  assert.ok(workloads.length >= 150);
  for (const w of workloads) {
    const ref = w.spec.subcomponentOf;
    if (ref) {
      const app = byRef.get(ref);
      assert.equal(app?.kind, 'Component', `${w.metadata.name} parent ${ref}`);
      assert.equal(w.spec.system, app.spec.system, `${w.metadata.name} follows its app's System`);
      assert.ok(!app.spec.subcomponentOf, 'apps are top-level');
    } else assert.equal(w.spec.system, UNASSIGNED, `${w.metadata.name} is neither mapped nor unassigned`);
  }
  const argocd = workloads.filter(w => w.spec.subcomponentOf === 'component:default/argocd').map(w => ann(w, 'onemind.io/kubernetes-workload'));
  for (const name of ['argocd-server', 'argocd-repo-server', 'argocd-dex-server', 'argocd-redis']) assert.ok(argocd.includes(name), `argocd needs ${name}`);
});

test('a workload no mapping covers lands in Unassigned workloads and nothing is dropped', () => {
  const input = structuredClone({ apps, cluster });
  const sample = structuredClone(workloads[0]);
  sample.metadata.name = 'k8s-new-ns-thing-deadbeef';
  sample.metadata.title = 'new-ns / thing';
  sample.metadata.annotations['backstage.io/kubernetes-namespace'] = 'new-ns';
  delete sample.metadata.annotations['onemind.io/kubernetes-workload'];
  delete sample.spec.subcomponentOf;
  input.cluster.push(sample);
  const out = applyStructure(input, sources);
  assert.equal(out.cluster.length, cluster.length + 1);
  const added = out.cluster.find(e => e.metadata.name === sample.metadata.name);
  assert.equal(added.spec.system, UNASSIGNED);
  assert.equal(added.spec.subcomponentOf, undefined);
  assert.equal(added.metadata.title, 'New Ns — Thing');
});

test('applying the structure is idempotent and never renames or drops a non-System entity', () => {
  const once = applyStructure({ apps, cluster }, sources);
  assert.deepEqual(once.apps, apps);
  assert.deepEqual(once.cluster, cluster);
  const baseline = new Set(JSON.parse(fs.readFileSync(path.join(root, 'catalog/ref-baseline.json'), 'utf8')).refs);
  const names = new Set(byRef.keys());
  for (const ref of baseline) if (!ref.startsWith('system:')) assert.ok(names.has(ref), `${ref} disappeared`);
});

test('colliding titles are disambiguated instead of dropped', () => {
  const input = structuredClone({ apps, cluster });
  const mesh = input.cluster.filter(e => e.spec?.system === 'onemind-mesh');
  mesh[1].metadata.title = mesh[0].metadata.title;
  const out = applyStructure(input, sources);
  const titles = out.cluster.filter(e => e.spec.system === 'onemind-mesh').map(e => e.metadata.title);
  assert.equal(new Set(titles).size, titles.length);
  assert.equal(out.cluster.length, cluster.length);
});

test('every Argo application in Git maps to exactly one app Component, and a new one fails', () => {
  const map = sources.map.argoApplications;
  const inGit = argoApplicationsInGit();
  assert.equal(inGit.size, 69, 'platform/ should define the 69 live Argo applications');
  const unmapped = [...inGit].filter(a => !map[a]);
  assert.deepEqual(unmapped, [], `Argo applications without a mapping in catalog/workload-map.yaml: ${unmapped}`);
  assert.deepEqual(Object.keys(map).filter(a => !inGit.has(a)), [], 'mapping names an Argo application that Git does not define');
  // The guard itself: an extra application is reported.
  assert.deepEqual([...new Set([...inGit, 'platform-new-thing'])].filter(a => !map[a]), ['platform-new-thing']);
  for (const [argo, ref] of Object.entries(map)) {
    const holders = all.filter(e => e.kind === 'Component' && !e.spec.subcomponentOf && [ann(e, 'argocd/app-name'), ...(ann(e, 'onemind.io/argocd-apps') || '').split(',')].includes(argo));
    assert.deepEqual(holders.map(refOf), [ref], `${argo} must resolve to exactly ${ref}`);
    assert.ok(ann(byRef.get(ref), 'argocd/instance-name'));
  }
});

test('no previously valid ref disappears without a redirect, and redirects resolve', () => {
  const { refs } = JSON.parse(fs.readFileSync(path.join(root, 'catalog/ref-baseline.json'), 'utf8'));
  const { redirects } = YAML.parse(fs.readFileSync(path.join(root, 'catalog/redirects.yaml'), 'utf8'));
  const from = new Map(redirects.map(r => [r.from, r]));
  assert.equal(from.size, redirects.length, 'duplicate redirect source');
  for (const ref of refs) assert.ok(byRef.has(ref) || from.has(ref), `${ref} was removed without a redirect`);
  for (const r of redirects) {
    assert.ok(!byRef.has(r.from), `${r.from} still exists, so it needs no redirect`);
    assert.ok(byRef.has(r.to), `${r.from} redirects to missing ${r.to}`);
    assert.ok(r.reason);
  }
});

test('pillars are reviewed, never guessed: enum only, inherited from the System, no served pillars yet', () => {
  const { structure } = sources;
  for (const e of [...systems, ...members]) {
    const p = e.metadata.labels?.['onemind.io/pillar'];
    assert.ok(p === undefined || PILLARS.includes(p), `${refOf(e)} pillar ${p}`);
    assert.equal(ann(e, 'onemind.io/pillars-served'), undefined);
  }
  for (const s of systems) {
    const declared = structure.systems.find(x => x.name === s.metadata.name).pillar;
    assert.equal(s.metadata.labels?.['onemind.io/pillar'], PILLARS.includes(declared) ? declared : undefined);
    assert.equal(s.spec.domain, PILLARS.includes(declared) ? `domain:default/onemind-${declared}` : 'domain:default/onemind');
  }
  for (const e of members) {
    const override = structure.pillarOverrides?.[refOf(e)];
    const expected = override ?? s(e).metadata.labels?.['onemind.io/pillar'];
    assert.equal(e.metadata.labels?.['onemind.io/pillar'], PILLARS.includes(expected) ? expected : undefined, refOf(e));
  }
  function s(e) { return byRef.get(`system:default/${e.spec.system}`); }
  const unclassified = members.filter(e => !e.metadata.labels?.['onemind.io/pillar']);
  assert.ok(unclassified.length > 0 && unclassified.some(e => e.metadata.name === 'mercury'));
});

test('mounted copies equal the source copies', () => {
  for (const f of ['apps.yaml', 'cluster.yaml', 'org.yaml']) {
    assert.equal(fs.readFileSync(path.join(root, 'catalog', f), 'utf8'), fs.readFileSync(path.join(root, '../../platform/backstage/catalog', f), 'utf8'), f);
  }
});
