// Applies the reviewed SO / LE / HP / GE structure (catalog/structure.yaml) and the reviewed
// workload and Argo mapping (catalog/workload-map.yaml) to the catalog files. Idempotent and
// pure: it never changes metadata.name, never deletes a non-System entity, and treats nothing
// as a pillar unless the reviewed structure says so. Run last in the refresh chain.
const fs = require('node:fs');
const path = require('node:path');
const YAML = require('yaml');

const root = path.resolve(__dirname, '..');
const OWNER = 'group:default/onemind-admins';
const UNASSIGNED = 'unassigned-workloads';
const PILLARS = ['so', 'le', 'hp', 'ge'];
const ACRONYMS = new Set(['api', 'db', 'ui', 'ldap', 'dns', 'nats', 'tak', 'pg', 'ml', 'id', 'tls', 'ha', 'ws', 'oc', 'ai', 'cpu', 'crd', 'crds', 'tcp', 'url', 'ip', 'sql', 'mcp', 'stac', 'osm']);
const refOf = e => `${e.kind.toLowerCase()}:${(e.metadata.namespace || 'default').toLowerCase()}/${e.metadata.name.toLowerCase()}`;
const load = file => YAML.parseAllDocuments(fs.readFileSync(file, 'utf8')).map(d => d.toJSON()).filter(Boolean);

function loadSources(base = root) {
  const read = f => YAML.parse(fs.readFileSync(path.join(base, 'catalog', f), 'utf8'));
  return { structure: read('structure.yaml'), map: read('workload-map.yaml') };
}

const SPECIAL = { github: 'GitHub', coredns: 'CoreDNS', applicationset: 'ApplicationSet', apiserver: 'API Server', kube: 'Kube' };
function pretty(token) {
  const t = token.toLowerCase();
  if (SPECIAL[t]) return SPECIAL[t];
  return ACRONYMS.has(t) ? t.toUpperCase() : t.charAt(0).toUpperCase() + t.slice(1);
}
const words = name => name.split(/[-_.]+/).filter(t => t && !/^[0-9a-f]{8}$/.test(t)).map(pretty).join(' ');

// Human role of a workload inside its app: drop everything up to and including the app's name.
function workloadRole(workload, appName, namespace = '') {
  const tokens = workload.split('-');
  const proxy = workload.match(/^ts-(.+)-([a-z0-9]{5})$/);
  if (proxy) return { role: `Proxy: ${proxy[1]}`, suffix: proxy[2] };
  const app = appName.split('-');
  for (let i = 0; i + app.length <= tokens.length; i++) {
    if (app.every((a, j) => tokens[i + j] === a)) {
      const rest = tokens.slice(i + app.length);
      return { role: rest.length ? words(rest.join('-')) : 'Core' };
    }
  }
  const ns = namespace.split('-');
  const stripped = namespace && tokens.length > ns.length && ns.every((n, i) => tokens[i] === n) ? tokens.slice(ns.length) : tokens;
  return { role: words(stripped.join('-')) };
}

function workloadIdentity(entity) {
  const a = entity.metadata.annotations || {};
  const namespace = a['backstage.io/kubernetes-namespace'];
  if (!namespace || a['onemind.io/inventory-source'] !== 'kubernetes') return null;
  let name = a['onemind.io/kubernetes-workload'];
  const title = entity.metadata.title || '';
  if (!name && title.startsWith(namespace + ' / ')) name = title.slice(namespace.length + 3);
  if (!name) return null;
  const kind = a['onemind.io/kubernetes-workload-kind'] || (entity.metadata.description || '').match(/(StatefulSet|Deployment|DaemonSet) in /)?.[1] || 'Deployment';
  return { namespace, name, kind };
}

function ruleMatches(pattern, key) {
  const [pn, pw] = pattern.split('/');
  const [kn, kw] = key.split('/');
  if (pn !== kn) return false;
  if (pw === '*') return true;
  if (pw.endsWith('*')) return kw.startsWith(pw.slice(0, -1));
  return pw === kw;
}

// Resolve "<namespace>/<workload>" to an app ref, exact entries first, then "ns/prefix*" and "ns/*".
function appForWorkload(map, namespace, name) {
  const key = `${namespace}/${name}`;
  if (map.workloads[key]) return map.workloads[key];
  const rules = Object.entries(map.workloads).filter(([p]) => p.includes('*'));
  const hit = rules.filter(([p]) => ruleMatches(p, key)).sort((a, b) => b[0].length - a[0].length)[0];
  return hit ? hit[1] : null;
}

function domainEntities(structure) {
  const byName = new Map(structure.domains.map(d => [d.name, d]));
  return structure.domains.map(d => ({
    apiVersion: 'backstage.io/v1alpha1',
    kind: 'Domain',
    metadata: {
      name: d.name,
      title: d.title,
      description: d.description,
      labels: { 'onemind.io/catalog-axis': 'mission', ...(PILLARS.includes(d.pillar) ? { 'onemind.io/pillar': d.pillar } : {}) },
      annotations: { 'onemind.io/pillar-review': 'proposed' },
    },
    spec: { owner: OWNER, ...(d.name === 'onemind' ? {} : { subdomainOf: 'domain:default/onemind' }) },
  })).filter(d => byName.has(d.metadata.name));
}

function systemEntities(structure) {
  return structure.systems.map(s => ({
    apiVersion: 'backstage.io/v1alpha1',
    kind: 'System',
    metadata: {
      name: s.name,
      title: s.title,
      description: s.description,
      labels: { 'onemind.io/catalog-axis': 'deployment', ...(PILLARS.includes(s.pillar) ? { 'onemind.io/pillar': s.pillar } : {}) },
      annotations: { 'onemind.io/pillar-review': PILLARS.includes(s.pillar) ? 'proposed' : 'unclassified' },
    },
    spec: { owner: OWNER, domain: s.domain },
  }));
}

function appEntity(app, previous) {
  const base = previous ? structuredClone(previous) : {
    apiVersion: 'backstage.io/v1alpha1',
    kind: 'Component',
    metadata: { name: app.name, tags: ['kubernetes'] },
    spec: { type: 'service', lifecycle: 'production', owner: OWNER },
  };
  base.metadata.title = app.title;
  base.metadata.description = app.description;
  base.metadata.annotations = {
    ...(base.metadata.annotations || {}),
    'onemind.io/hosting': 'digitalocean-nyc3',
    'onemind.io/structure-source': 'catalog/structure.yaml',
  };
  base.spec.system = app.system;
  return base;
}

function pillarFor(structure, system, ref) {
  const override = structure.pillarOverrides?.[ref];
  if (override) return PILLARS.includes(override) ? override : null;
  const p = structure.systems.find(s => s.name === system)?.pillar;
  return PILLARS.includes(p) ? p : null;
}

function setLabel(entity, key, value) {
  const labels = { ...(entity.metadata.labels || {}) };
  if (value) labels[key] = value;
  else delete labels[key];
  if (Object.keys(labels).length) entity.metadata.labels = labels;
  else delete entity.metadata.labels;
}

function commonSelector(workloads) {
  const sets = workloads.map(w => new Set((w.metadata.annotations['backstage.io/kubernetes-label-selector'] || '').split(',').filter(Boolean)));
  if (!sets.length || sets.some(s => !s.size)) return '';
  return [...sets[0]].filter(p => sets.every(s => s.has(p))).sort().join(',');
}

/**
 * @param {{apps: object[], cluster: object[]}} input parsed apps.yaml and cluster.yaml documents
 * @param {{structure: object, map: object}} sources
 */
function applyStructure(input, sources) {
  const { structure, map } = sources;
  const apps = structuredClone(input.apps);
  const cluster = structuredClone(input.cluster);
  const systemNames = new Set(structure.systems.map(s => s.name));
  const newApps = new Map((structure.apps || []).map(a => [a.name, a]));
  const appSystem = new Map([...newApps].map(([n, a]) => [n, a.system]));
  for (const [n, s] of Object.entries(structure.placement)) appSystem.set(n, s);
  const domainRefs = new Set(structure.domains.map(d => `domain:default/${d.name}`));
  const keep = e => e.kind !== 'System' && !(e.kind === 'Domain' && domainRefs.has(refOf(e)));
  const oldApps = apps.filter(keep);
  const oldCluster = cluster.filter(keep);

  const existingByName = new Map(oldApps.filter(e => e.kind === 'Component').map(e => [e.metadata.name, e]));
  const appDocs = [];
  for (const app of newApps.values()) appDocs.push(appEntity(app, existingByName.get(app.name)));
  const newAppNames = new Set(newApps.keys());
  const appsOut = [...oldApps.filter(e => !(e.kind === 'Component' && newAppNames.has(e.metadata.name))), ...appDocs];
  const clusterOut = oldCluster;

  const appComponent = new Map(appsOut.filter(e => e.kind === 'Component').map(e => [e.metadata.name, e]));
  const argoByApp = new Map();
  for (const [argo, ref] of Object.entries(map.argoApplications)) {
    const name = ref.split('/')[1];
    argoByApp.set(name, [...(argoByApp.get(name) || []), argo]);
  }

  // --- placement
  const subcomponents = new Map();
  const workloadTitles = [];
  for (const e of clusterOut) {
    if (e.kind === 'Component') {
      const id = workloadIdentity(e);
      if (!id) continue;
      e.metadata.annotations['onemind.io/kubernetes-workload'] = id.name;
      e.metadata.annotations['onemind.io/kubernetes-workload-kind'] = id.kind;
      const appRef = appForWorkload(map, id.namespace, id.name);
      const appName = appRef?.split('/')[1];
      if (appName && appComponent.has(appName) && appSystem.has(appName)) {
        e.spec.system = appSystem.get(appName);
        e.spec.subcomponentOf = appRef;
        subcomponents.set(appName, [...(subcomponents.get(appName) || []), e]);
        const app = appComponent.get(appName);
        const { role, suffix } = workloadRole(id.name, appName, id.namespace);
        workloadTitles.push({ entity: e, base: `${app.metadata.title} — ${role}`, suffix, namespace: id.namespace });
      } else {
        e.spec.system = UNASSIGNED;
        delete e.spec.subcomponentOf;
        workloadTitles.push({ entity: e, base: `${words(id.namespace)} — ${words(id.name)}`, namespace: id.namespace });
      }
    }
  }
  const resolveSystem = e => {
    const old = e.spec?.system;
    if (systemNames.has(old) && !structure.fallbackSystems?.[old]) return old;
    if (structure.fallbackSystems?.[old]) return structure.fallbackSystems[old];
    return old && systemNames.has(old) ? old : null;
  };
  for (const e of [...appsOut, ...clusterOut]) {
    if (e.kind === 'Domain') continue;
    if (e.kind === 'Component' && workloadIdentity(e)) continue;
    const placed = e.kind === 'Component' || e.kind === 'API' ? structure.placement[e.metadata.name] : null;
    e.spec.system = placed || resolveSystem(e) || (e.spec.type === 'kubernetes-cluster' ? 'infrastructure' : null);
    if (!e.spec.system) throw new Error(`No reviewed System for ${refOf(e)}`);
  }

  // --- Argo links and Kubernetes selectors on the app entities
  for (const [name, list] of argoByApp) {
    const app = appComponent.get(name);
    if (!app) throw new Error(`Argo application target ${name} has no app Component`);
    const exact = list.find(a => a.replace(/^(platform|ecosystem|edge|monitoring)-/, '') === name);
    const primary = map.primaryArgoApplication?.[name] || exact || [...list].sort()[0];
    const ann = app.metadata.annotations || (app.metadata.annotations = {});
    ann['argocd/app-name'] = primary;
    ann['argocd/instance-name'] = 'onemind-cloud';
    if (list.length > 1) ann['onemind.io/argocd-apps'] = [...list].sort().join(',');
    else delete ann['onemind.io/argocd-apps'];
  }
  for (const [name, list] of subcomponents) {
    const app = appComponent.get(name);
    const live = list.filter(w => w.metadata.annotations['onemind.io/inventory-state'] !== 'missing-review-required' && w.metadata.annotations['onemind.io/inventory-state'] !== 'retirement-review-required');
    const namespaces = new Set(live.map(w => w.metadata.annotations['backstage.io/kubernetes-namespace']));
    const selector = commonSelector(live);
    const ann = app.metadata.annotations || (app.metadata.annotations = {});
    // Only fill what is missing; a reviewed selector or kubernetes-id is never replaced.
    if (namespaces.size === 1 && selector && !ann['backstage.io/kubernetes-id'] && !ann['backstage.io/kubernetes-label-selector']) {
      ann['backstage.io/kubernetes-namespace'] = [...namespaces][0];
      ann['backstage.io/kubernetes-label-selector'] = selector;
    }
  }

  // --- titles: workloads first, then guarantee uniqueness inside each System
  for (const w of workloadTitles) w.entity.metadata.title = w.base;
  for (const e of [...appsOut, ...clusterOut]) {
    if (!e.metadata.title) e.metadata.title = e.metadata.name;
    for (const [name, title] of Object.entries(structure.titles || {})) if (refOf(e) === name) e.metadata.title = title;
  }
  const seen = new Map();
  const bySystem = [...appsOut, ...clusterOut].filter(e => e.kind !== 'Domain');
  for (const e of bySystem) {
    const key = `${e.spec.system}\u0000${e.metadata.title}`;
    if (!seen.has(key)) { seen.set(key, e); continue; }
    const w = workloadTitles.find(x => x.entity === e);
    const tag = w ? (w.suffix || w.namespace) : e.kind === 'Resource' ? 'registry' : e.spec.type || e.kind.toLowerCase();
    let candidate = `${e.metadata.title} (${tag})`;
    if (seen.has(`${e.spec.system}\u0000${candidate}`)) candidate = `${e.metadata.title} (${tag} ${e.metadata.name.slice(-6)})`;
    e.metadata.title = candidate;
    seen.set(`${e.spec.system}\u0000${candidate}`, e);
  }

  // --- pillar labels follow the reviewed System; nothing is guessed
  const finalDocs = [...appsOut, ...clusterOut];
  for (const e of finalDocs) {
    if (e.kind === 'Domain') continue;
    setLabel(e, 'onemind.io/pillar', pillarFor(structure, e.spec.system, refOf(e)));
  }
  for (const e of finalDocs.filter(x => x.kind === 'Component')) {
    const w = workloadIdentity(e);
    if (w && !e.spec.subcomponentOf) delete e.spec.subcomponentOf;
  }

  const head = [...domainEntities(structure), ...systemEntities(structure)];
  return { apps: [...head, ...appsOut], cluster: clusterOut };
}

const isObject = v => v && typeof v === 'object' && !Array.isArray(v);
const same = (a, b) => JSON.stringify(a) === JSON.stringify(b);

// Apply only the changed leaves to the original YAML document so untouched entities keep their
// exact bytes (scalar quoting, folding, ordering) and the diff shows real changes only.
function patchDocument(doc, before, after, trail = []) {
  if (isObject(before) && isObject(after)) {
    for (const key of Object.keys(after)) {
      if (!(key in before)) doc.setIn([...trail, key], doc.createNode(after[key]));
      else patchDocument(doc, before[key], after[key], [...trail, key]);
    }
    for (const key of Object.keys(before)) if (!(key in after)) doc.deleteIn([...trail, key]);
  } else if (!same(before, after)) doc.setIn(trail, doc.createNode(after));
}

function ownedRefs(structure) {
  return new Set([
    ...structure.domains.map(d => `domain:default/${d.name}`),
    ...structure.systems.map(s => `system:default/${s.name}`),
    ...(structure.apps || []).map(a => `component:default/${a.name}`),
  ]);
}

/** Text in, text out: the form the catalog files are rewritten in. */
function applyStructureToText(text, sources) {
  const parse = t => YAML.parseAllDocuments(t).filter(d => d.contents);
  const sides = { apps: parse(text.apps), cluster: parse(text.cluster) };
  const before = Object.fromEntries(Object.entries(sides).map(([k, docs]) => [k, docs.map(d => d.toJS())]));
  const out = applyStructure(before, sources);
  const owned = ownedRefs(sources.structure);
  const render = key => {
    const existing = new Map(sides[key].map((d, i) => [refOf(before[key][i]), { doc: d, js: before[key][i] }]));
    return out[key].map(entity => {
      const ref = refOf(entity);
      const hit = existing.get(ref);
      if (owned.has(ref) || !hit) return YAML.stringify(entity);
      patchDocument(hit.doc, hit.js, entity);
      hit.doc.directives.docStart = null;
      return hit.doc.toString();
    }).join('---\n');
  };
  return { apps: render('apps'), cluster: render('cluster') };
}

function writeCatalog(result, base = root) {
  for (const dir of [path.join(base, 'catalog'), path.join(base, '../../platform/backstage/catalog')]) {
    fs.writeFileSync(path.join(dir, 'apps.yaml'), result.apps);
    fs.writeFileSync(path.join(dir, 'cluster.yaml'), result.cluster);
  }
}

if (require.main === module) {
  const read = f => fs.readFileSync(path.join(root, 'catalog', f), 'utf8');
  const result = applyStructureToText({ apps: read('apps.yaml'), cluster: read('cluster.yaml') }, loadSources());
  writeCatalog(result);
  const count = (text, kind) => (text.match(new RegExp(`^kind: ${kind}$`, 'gm')) || []).length;
  console.log(`Structure applied: ${count(result.apps, 'Domain')} domains, ${count(result.apps, 'System')} systems, ${(result.cluster.match(/^ {2}subcomponentOf: /gm) || []).length} subcomponents.`);
}
module.exports = { applyStructure, applyStructureToText, loadSources, load, refOf, workloadIdentity, appForWorkload, workloadRole, writeCatalog, UNASSIGNED };
