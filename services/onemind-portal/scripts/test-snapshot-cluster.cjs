const { test } = require('node:test');
const assert = require('node:assert/strict');
const path = require('node:path');
const { load } = require('./curate-structure.cjs');
const { buildSnapshot } = require('./snapshot-cluster.cjs');

const apps = load(path.join(__dirname, '../catalog/apps.yaml'));
const item = (kind, namespace, name, labels) => ({ kind, metadata: { namespace, name }, spec: { selector: { matchLabels: labels } } });

test('snapshot keeps the hash scheme, attaches mapped workloads and leaves new ones visible', () => {
  const docs = buildSnapshot([
    item('Deployment', 'argocd', 'argocd-server', { 'app.kubernetes.io/instance': 'argocd', 'app.kubernetes.io/name': 'argocd-server' }),
    item('Deployment', 'brand-new', 'thing', { app: 'thing' }),
  ], 'do-nyc3-omos-cloud', '2026-10-21T00:00:00.000Z', apps);
  assert.ok(!docs.some(d => d.kind === 'System'), 'namespace Systems are no longer generated');
  const server = docs.find(d => d.metadata.name === 'k8s-argocd-argocd-server-4414022f');
  assert.equal(server.spec.subcomponentOf, 'component:default/argocd');
  assert.equal(server.spec.system, 'gitops-delivery');
  assert.equal(server.metadata.title, 'ArgoCD — Server');
  const fresh = docs.find(d => d.metadata.annotations?.['onemind.io/kubernetes-workload'] === 'thing');
  assert.equal(fresh.spec.system, 'unassigned-workloads');
  assert.equal(fresh.spec.subcomponentOf, undefined);
  assert.equal(docs.find(d => d.kind === 'Resource').spec.system, 'infrastructure');
});
