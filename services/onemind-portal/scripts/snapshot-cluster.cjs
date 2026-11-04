// Run after reviewing the active kubectl context. Captures identities, never credentials or runtime specs.
// Workloads are attached to their reviewed app (catalog/workload-map.yaml); unmapped ones land in
// system:default/unassigned-workloads. Entity names keep the original hash scheme.
const path = require('path');
const crypto = require('crypto');
const fs = require('fs');
const cp = require('child_process');
const YAML = require('yaml');
const { applyStructure, loadSources, load } = require('./curate-structure.cjs');

const base = path.resolve(__dirname, '..');

function buildSnapshot(items, context, observedAt, apps) {
  const docs = [{ apiVersion: 'backstage.io/v1alpha1', kind: 'Resource', metadata: { name: context, title: 'OneMind Cloud — NYC3' }, spec: { type: 'kubernetes-cluster', owner: 'group:default/onemind-admins' } }];
  for (const x of items) {
    const ns = x.metadata.namespace;
    const name = x.metadata.name;
    const labels = Object.entries(x.spec.selector.matchLabels).map(([k, v]) => k + '=' + v).join(',');
    docs.push({
      apiVersion: 'backstage.io/v1alpha1', kind: 'Component',
      metadata: {
        name: ('k8s-' + ns + '-' + name).slice(0, 54) + '-' + crypto.createHash('sha256').update(ns + '/' + name).digest('hex').slice(0, 8),
        title: ns + ' / ' + name, description: x.kind + ' in ' + ns,
        annotations: {
          'backstage.io/kubernetes-label-selector': labels, 'backstage.io/kubernetes-namespace': ns,
          'onemind.io/inventory-source': 'kubernetes', 'onemind.io/inventory-observed-at': observedAt,
          'onemind.io/kubernetes-workload': name, 'onemind.io/kubernetes-workload-kind': x.kind,
        },
        tags: ['kubernetes', ns],
      },
      spec: { type: 'service', lifecycle: 'production', owner: 'group:default/onemind-admins', system: 'unassigned-workloads', dependsOn: ['resource:default/' + context] },
    });
  }
  return applyStructure({ apps, cluster: docs }, loadSources(base)).cluster;
}

if (require.main === module) {
  process.chdir(base);
  const context = cp.execFileSync('kubectl', ['config', 'current-context'], { encoding: 'utf8' }).trim();
  if (context !== 'do-nyc3-omos-cloud') throw new Error('Unexpected Kubernetes context');
  const data = JSON.parse(cp.execFileSync('kubectl', ['get', 'deployments,statefulsets', '-A', '-o', 'json'], { encoding: 'utf8', maxBuffer: 20 * 1024 * 1024 }));
  const docs = buildSnapshot(data.items, context, new Date().toISOString(), load('catalog/apps.yaml'));
  const rendered = docs.map(x => YAML.stringify(x)).join('---\n');
  fs.writeFileSync('catalog/cluster.yaml', rendered);
  fs.mkdirSync('../../platform/backstage/catalog', { recursive: true });
  fs.writeFileSync('../../platform/backstage/catalog/cluster.yaml', rendered);
  console.log('Captured ' + data.items.length + ' workload identities in ' + new Set(data.items.map(x => x.metadata.namespace)).size + ' namespaces');
}
module.exports = { buildSnapshot };
