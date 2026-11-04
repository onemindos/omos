// Offline metadata report; no catalog writes, network calls or credential access.
const fs = require('node:fs');
const MAX_BYTES = 32 * 1024 * 1024;
const MAX_ENTITIES = 10000;
const FIELDS = ['source', 'owner', 'kubernetesSelector', 'kubernetesNamespace',
  'location', 'argo', 'grafana', 'prometheus', 'operations', 'dependency', 'runbook', 'aiApi'];
const object = value => value !== null && typeof value === 'object' && !Array.isArray(value);
const nonempty = value => typeof value === 'string' && value.trim().length > 0;
const part = value => typeof value === 'string' && /^[a-z0-9][a-z0-9_.-]{0,62}$/i.test(value);
const fail = () => { throw new Error('Invalid complete catalog array or reviewed exclusions'); };
function reference(value, kind, namespace) {
  if (!nonempty(value)) fail();
  const match = /^(?:([a-z][a-z0-9]{0,62}):)?(?:([^/:]+)\/)?([^/:]+)$/i.exec(value);
  if (!match || !part(match[2] || namespace) || !part(match[3])) fail();
  return `${match[1] || kind}:${match[2] || namespace}/${match[3]}`.toLowerCase();
}
const ref = e => reference(`${e.kind}:${e.metadata.namespace ?? 'default'}/${e.metadata.name}`);
function validate(e) {
  if (!object(e) || typeof e.kind !== 'string' || !/^[a-z][a-z0-9]{0,62}$/i.test(e.kind) ||
      !object(e.metadata) || !part(e.metadata.name) ||
      !part(e.metadata.namespace ?? 'default')) fail();
  for (const field of ['annotations', 'labels']) {
    const map = e.metadata[field];
    if (map !== undefined && (!object(map) || Object.values(map).some(v => typeof v !== 'string'))) fail();
  }
  if (e.metadata.links !== undefined && (!Array.isArray(e.metadata.links) ||
      e.metadata.links.some(l => !object(l) || !nonempty(l.url) ||
        ['type', 'title'].some(k => l[k] !== undefined && typeof l[k] !== 'string')))) fail();
  if (e.spec !== undefined && !object(e.spec)) fail();
  const spec = e.spec || {};
  if (spec.owner !== undefined) reference(spec.owner, 'group', e.metadata.namespace ?? 'default');
  for (const key of ['dependsOn', 'providesApis', 'consumesApis']) {
    if (spec[key] !== undefined && (!Array.isArray(spec[key]) || spec[key].length > MAX_ENTITIES)) fail();
    for (const value of spec[key] || []) reference(value, key === 'dependsOn' ? 'component' : 'api', e.metadata.namespace ?? 'default');
  }
  if (e.relations !== undefined && (!Array.isArray(e.relations) || e.relations.length > MAX_ENTITIES)) fail();
  for (const r of e.relations || []) {
    if (!object(r) || !nonempty(r.type)) fail();
    if (!/^[a-z][a-z0-9]*:[^/]+\/[^/]+$/i.test(r.targetRef || '')) fail();
    reference(r.targetRef);
  }
}
function reportServiceCoverage(entities, reviewedNotApplicable = []) {
  if (!Array.isArray(entities) || entities.length > MAX_ENTITIES ||
      Buffer.byteLength(JSON.stringify(entities)) > MAX_BYTES || !Array.isArray(reviewedNotApplicable) ||
      reviewedNotApplicable.length > MAX_ENTITIES) fail();
  const byRef = new Map();
  for (const e of entities) {
    validate(e);
    const id = ref(e);
    if (byRef.has(id)) fail();
    byRef.set(id, e);
  }
  const exclusions = new Map();
  for (const item of reviewedNotApplicable) {
    if (!object(item) || !nonempty(item.ref) || !Array.isArray(item.fields) ||
        item.fields.length === 0 || item.fields.some(f => !FIELDS.includes(f))) fail();
    const id = reference(item.ref);
    if (!id.startsWith('component:') || !byRef.has(id) || exclusions.has(id)) fail();
    exclusions.set(id, new Set(item.fields));
  }
  const rows = [];
  for (const [id, e] of byRef) {
    if (e.kind.toLowerCase() !== 'component') continue;
    const a = e.metadata.annotations || {}, links = e.metadata.links || [], spec = e.spec || {};
    const namespace = e.metadata.namespace ?? 'default';
    const fields = {};
    for (const field of FIELDS) fields[field] = { status: 'missing', evidence: [], verification: 'not-tested' };
    const add = (field, mode, key) => fields[field].evidence.push({ mode, key });
    const annotations = (field, keys) => {
      for (const key of keys) if (nonempty(a[key])) add(field, 'annotation', key);
    };
    const typedLinks = (field, types, titles = []) => {
      if (links.some(l => types.includes(l.type) || titles.includes(l.title))) add(field, 'link', `metadata.links:${field}`);
    };
    annotations('source', ['backstage.io/source-location', 'github.com/project-slug', 'onemind.io/deployment-source']);
    if (nonempty(spec.owner)) add('owner', 'spec', 'spec.owner');
    annotations('kubernetesSelector', ['backstage.io/kubernetes-label-selector', 'backstage.io/kubernetes-id']);
    annotations('kubernetesNamespace', ['backstage.io/kubernetes-namespace']);
    typedLinks('location', ['location']);
    annotations('argo', ['argocd/app-name', 'argocd/app-selector']);
    typedLinks('argo', ['argocd']);
    annotations('grafana', ['grafana/dashboard-selector', 'grafana/overview-dashboard', 'grafana/alert-label-selector']);
    typedLinks('grafana', ['grafana']);
    annotations('prometheus', ['prometheus.io/service-name', 'prometheus.io/alert', 'prometheus.io/rule', 'prometheus.io/labels']);
    typedLinks('prometheus', ['prometheus']);
    // Source mapping mirrors prometheusApi.ts, not evidence of a working backend.
    const scopes = { backstage: 'k8s-backstage-backstage-21b99354', postgres: 'k8s-backstage-backstage-postgres-17d2f351', nats: 'k8s-nats-nats-0dee54ac' };
    if (namespace === 'default' && Object.hasOwn(scopes, a['prometheus.io/service-name']) &&
        scopes[a['prometheus.io/service-name']] === e.metadata.name)
      add('prometheus', 'configured-mapping', 'prometheusApi.scopeForEntity');
    if (nonempty(a['backstage.io/kubernetes-namespace']))
      add('operations', 'configured-mapping', 'ServicesModule.namespace-filter');
    typedLinks('operations', ['operations']);
    typedLinks('runbook', ['runbook'], ['Operations runbook', 'Recovery runbook']);
    const targets = (key, relation, kind) => [...new Set([
      ...(spec[key] || []).map(v => reference(v, kind, namespace)),
      ...(e.relations || []).filter(r => r.type === relation).map(r => reference(r.targetRef)),
    ])].sort();
    const dependencies = targets('dependsOn', 'dependsOn', 'component');
    if (dependencies.length) add('dependency', 'spec-or-relation', 'dependsOn');
    fields.dependency.targets = { recorded: dependencies.length, present: dependencies.filter(t => byRef.has(t)).length };
    if (dependencies.some(t => byRef.get(t)?.spec?.type === 'kubernetes-cluster'))
      add('location', 'configured-mapping', 'dependsOn:kubernetes-cluster');
    if ((e.relations || []).some(r => r.type === 'ownedBy')) add('owner', 'relation', 'ownedBy');
    const apis = targets('providesApis', 'providesApi', 'api');
    if (apis.length) add('aiApi', 'spec-or-relation', 'providesApi');
    typedLinks('aiApi', ['api', 'openapi', 'mcp']);
    fields.aiApi.contracts = { recorded: apis.length,
      present: apis.filter(t => byRef.get(t)?.kind.toLowerCase() === 'api').length,
      withDefinition: apis.filter(t => byRef.get(t)?.kind.toLowerCase() === 'api' && nonempty(byRef.get(t)?.spec?.definition)).length };
    fields.aiApi.authenticatedAccess = 'not-tested';
    for (const field of FIELDS) {
      if (fields[field].evidence.length) fields[field].status = 'available';
      if (exclusions.get(id)?.has(field)) {
        if (fields[field].evidence.length) fail();
        fields[field].status = 'not-applicable';
        add(field, 'reviewed-exclusion', 'reviewedNotApplicable');
      }
    }
    rows.push({ ref: id, fields });
  }
  rows.sort((a, b) => a.ref < b.ref ? -1 : a.ref > b.ref ? 1 : 0);
  const summary = {};
  for (const field of FIELDS) {
    summary[field] = { available: 0, missing: 0, 'not-applicable': 0 };
    for (const row of rows) summary[field][row.fields[field].status]++;
  }
  return { schemaVersion: 'onemind.service-coverage.v1',
    scope: 'Offline catalog metadata only; no runtime health, reachability, authenticated AI access or live acceptance tested. Complete array required; completeness cannot be inferred from JSON.',
    entityCount: entities.length, componentCount: rows.length,
    ignoredNonComponents: entities.length - rows.length, summary, components: rows };
}
function readJson(file) {
  // Read at most the bound, including stdin; reject oversized input before parsing.
  const fd = file === '-' ? 0 : fs.openSync(file, 'r');
  const chunks = [];
  let total = 0;
  try {
    for (;;) {
      const buffer = Buffer.alloc(65536);
      const size = fs.readSync(fd, buffer, 0, buffer.length, null);
      if (!size) break;
      total += size;
      if (total > MAX_BYTES) fail();
      chunks.push(buffer.subarray(0, size));
    }
    return JSON.parse(Buffer.concat(chunks).toString('utf8'));
  } finally { if (file !== '-') fs.closeSync(fd); }
}
if (require.main === module) {
  try {
    const args = process.argv.slice(2);
    if (args.length < 1 || args.length > 2 || (args.length === 2 && args[1] === '-'))
      throw new Error('usage');
    console.log(JSON.stringify(reportServiceCoverage(readJson(args[0]), args[1] ? readJson(args[1]) : []), null, 2));
  } catch {
    // Never echo parser errors, file paths, annotation values or URL queries.
    console.error('Coverage report rejected: expected a complete entity JSON array and optional reviewed exclusions JSON. Usage: node scripts/report-service-coverage.cjs <complete-entities.json|-> [reviewed-exclusions.json]');
    process.exitCode = 2;
  }
}
module.exports = { reportServiceCoverage, FIELDS };
