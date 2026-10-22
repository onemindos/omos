const { test } = require('node:test');
const assert = require('node:assert/strict');
const fs = require('node:fs');
const path = require('node:path');
const YAML = require('yaml');
const { curateTopology, graphUrl, names, storage, pvcRef, renderCluster } = require('./curate-topology.cjs');
// Repo-level platform/ is absent in the image build context (only services/onemind-portal is copied).
const hasRepo = fs.existsSync(path.join(__dirname, '../../../platform/backstage/catalog'));
const repoOnly = hasRepo ? {} : { skip: 'repo-level platform/ is not in this build context' };
const input = () => YAML.parseAllDocuments(fs.readFileSync(path.join(__dirname, '../catalog/cluster.yaml'), 'utf8')).map(d => d.toJSON());
const key = e => `${e.kind.toLowerCase()}:${e.metadata.namespace || 'default'}/${e.metadata.name}`;

test('idempotent, immutable curation preserves identities and unrelated metadata', () => {
  const docs = input(), original = structuredClone(docs);
  const first = curateTopology(docs);
  assert.deepEqual(docs, original);
  assert.deepEqual(curateTopology(first), first);
  assert.deepEqual(first.map(key).sort(), docs.map(key).sort());
  const portal = first.find(e => e.metadata.name === names.portal);
  assert.ok(portal.spec.dependsOn.includes(`component:default/${names.database}`));
  assert.ok(portal.spec.dependsOn.includes(`component:default/${names.nats}`));
  assert.ok(!portal.metadata.annotations['onemind.io/topology-scope'].includes('verified healthy'));
});
test('missing/duplicate/relocated or changed target fails closed', () => {
  let docs = input().filter(e => e.metadata.name !== names.identity);
  assert.throws(() => curateTopology(docs), /missing\/ambiguous/);
  docs = input(); docs.push(structuredClone(docs.find(e => e.metadata.name === names.grafana)));
  assert.throws(() => curateTopology(docs), /ambiguous/);
  docs = input(); docs.find(e => e.metadata.name === names.argo).metadata.annotations['backstage.io/kubernetes-label-selector'] = 'app=unreviewed';
  assert.throws(() => curateTopology(docs), /identity changed/);
  docs = input(); docs.find(e => e.metadata.name === names.nats).spec.dependsOn = [];
  assert.throws(() => curateTopology(docs), /location changed/);
});
test('five existing PVC identities only; collision rejected and custom metadata retained', () => {
  let docs = input().filter(e => !storage.some(s => key(e) === pvcRef(s)));
  const out = curateTopology(docs);
  assert.equal(out.length - docs.length, 5);
  for (const s of storage) assert.ok(out.find(e => e.metadata.name === s.component).spec.dependsOn.includes(pvcRef(s)));
  docs = input(); const pvc = docs.find(e => key(e) === pvcRef(storage[0]));
  pvc.metadata.links.push({ title: 'Custom link', url: 'https://example.test/' });
  assert.ok(curateTopology(docs).find(e => key(e) === key(pvc)).metadata.links.some(l => l.title === 'Custom link'));
  pvc.metadata.annotations['onemind.io/inventory-source'] = 'manual';
  assert.throws(() => curateTopology(docs), /Unreviewed storage identity collision/);
});
test('native graph queries bound direction/depth and preserve exact root refs', () => {
  for (const reverse of [false, true]) {
    const u = new URL(graphUrl([names.portal, names.nats], reverse));
    assert.equal(u.pathname, '/catalog-graph');
    assert.deepEqual(u.searchParams.getAll('rootEntityRefs[]'), [names.portal, names.nats].map(n => `component:default/${n}`));
    assert.deepEqual(u.searchParams.getAll('selectedRelations[]'), [reverse ? 'dependencyOf' : 'dependsOn']);
    assert.equal(u.searchParams.get('unidirectional'), 'true');
    assert.equal(u.searchParams.get('maxDepth'), reverse ? '1' : '2');
  }
});
test('quoted long URLs round-trip and catalog pairs match without duplicate edges/links', repoOnly, () => {
  const out = curateTopology(input()), text = renderCluster(out);
  assert.deepEqual(YAML.parseAllDocuments(text).map(d => d.toJSON()), out);
  assert.ok(text.split('\n').filter(line => line.length > 200).every(line => line.includes('backstage.io/kubernetes-label-selector')));
  assert.equal(fs.readFileSync(path.join(__dirname, '../catalog/cluster.yaml'), 'utf8'), fs.readFileSync(path.join(__dirname, '../../../platform/backstage/catalog/cluster.yaml'), 'utf8'));
  for (const name of Object.values(names)) {
    const e = out.find(x => x.metadata.name === name);
    assert.equal(new Set(e.spec.dependsOn).size, e.spec.dependsOn.length);
    assert.equal(e.metadata.links.filter(l => l.title === 'Dependency topology (documented)').length, 1);
  }
});
// F2 structure guards run with the topology gate so the existing pipeline step covers them.
require('./test-catalog-structure.cjs');
require('./test-snapshot-cluster.cjs');
