// Catalog-only projection of reviewed source wiring. Never a traffic or health map.
const YAML = require('yaml');
const portalUrl = 'https://onemind-preview.tail717ba0.ts.net';
const cluster = 'resource:default/do-nyc3-omos-cloud';
const names = {
  portal: 'k8s-backstage-backstage-21b99354',
  database: 'k8s-backstage-backstage-postgres-17d2f351',
  nats: 'k8s-nats-nats-0dee54ac',
  grafana: 'k8s-monitoring-monitoring-grafana-8ed6cac6',
  identity: 'k8s-auth-platform-auth-helm-authentik-server-78f74c27',
  argo: 'k8s-argocd-argocd-server-4414022f',
};
const expectedNamespaces = { portal: 'backstage', database: 'backstage', nats: 'nats', grafana: 'monitoring', identity: 'auth', argo: 'argocd' };
const expectedSelectors = {
  grafana: 'app.kubernetes.io/instance=monitoring,app.kubernetes.io/name=grafana',
  identity: 'app.kubernetes.io/component=server,app.kubernetes.io/instance=platform-auth-helm,app.kubernetes.io/name=authentik',
  argo: 'app.kubernetes.io/instance=argocd,app.kubernetes.io/name=argocd-server',
};
const normalize = s => (s || '').split(',').sort().join(',');
const storage = [
  { namespace: 'backstage', claim: 'data-backstage-postgres-0', component: names.database, source: 'platform/backstage/postgres.yaml volumeClaimTemplates.data' },
  { namespace: 'monitoring', claim: 'grafana-data', component: names.grafana, source: 'platform/values/kube-prometheus-stack.yaml grafana.persistence.existingClaim; platform/monitoring/grafana-state.yaml' },
  ...[0, 1, 2].map(n => ({ namespace: 'nats', claim: `nats-js-nats-${n}`, component: names.nats, source: 'platform/nats/nats-statefulset.yaml volumeClaimTemplates.nats-js and replicas=3' })),
];
const pvcRef = s => `resource:default/pvc-${s.namespace}-${s.claim}`;
function graphUrl(roots, reverse = false) {
  // Installed native graph parses bracket-array query keys. No default-cluster
  // reverse traversal: that would expand every workload sharing this location.
  const query = new URLSearchParams();
  for (const name of roots) query.append('rootEntityRefs[]', `component:default/${name}`);
  query.append('selectedRelations[]', reverse ? 'dependencyOf' : 'dependsOn');
  query.set('maxDepth', reverse ? '1' : '2');
  query.set('unidirectional', 'true');
  query.set('mergeRelations', 'true');
  return `${portalUrl}/catalog-graph?${query}`;
}
function exactly(docs, kind, name) {
  const found = docs.filter(e => e.kind === kind && (e.metadata.namespace || 'default') === 'default' && e.metadata.name === name);
  if (found.length !== 1) throw new Error(`Topology target missing/ambiguous: ${kind}/${name}`);
  return found[0];
}
function curateTopology(entities) {
  const docs = structuredClone(entities);
  const location = exactly(docs, 'Resource', cluster.split('/')[1]);
  if (location.spec.type !== 'kubernetes-cluster') throw new Error('Topology cluster type changed');
  for (const [role, name] of Object.entries(names)) {
    const e = exactly(docs, 'Component', name);
    if (e.metadata.annotations?.['backstage.io/kubernetes-namespace'] !== expectedNamespaces[role] || (expectedSelectors[role] && normalize(e.metadata.annotations['backstage.io/kubernetes-label-selector']) !== normalize(expectedSelectors[role]))) throw new Error(`Topology workload identity changed for ${name}`);
    if (!e.spec.dependsOn?.includes(cluster)) throw new Error(`Topology location changed for ${name}`);
    e.metadata.annotations['onemind.io/topology-scope'] = 'Reviewed source dependencies, not observed traffic, live health or complete failure impact';
    e.metadata.annotations['onemind.io/topology-evidence'] = 'docs/operations/native-service-topology-2026-10-20.md';
    const links = [
      { title: 'Dependency topology (documented)', type: 'topology', url: graphUrl([name]) },
      { title: 'Dependent services (documented)', type: 'topology', url: graphUrl([name], true) },
      { title: 'Deployment location — NYC3', type: 'location', url: `${portalUrl}/catalog/default/resource/do-nyc3-omos-cloud` },
    ];
    if (role === 'portal') links.push({ title: 'Core service topology (documented)', type: 'topology', url: graphUrl([names.portal, names.nats, names.database]) });
    const titles = new Set(links.map(l => l.title));
    e.metadata.links = [...(e.metadata.links || []).filter(l => !titles.has(l.title)), ...links];
  }
  // These are real feature dependencies, not a claim that every outage kills the
  // portal. The PostgreSQL runtime dependency is curated by the existing registry.
  const portal = exactly(docs, 'Component', names.portal);
  portal.spec.dependsOn = [...new Set([...portal.spec.dependsOn, ...['identity', 'nats', 'grafana', 'argo'].map(role => `component:default/${names[role]}`)])];
  portal.metadata.annotations['onemind.io/dependency-scope'] = 'PostgreSQL: runtime; Authentik: sign-in; NATS: bus summaries; Grafana: dashboards; Argo CD: delivery visibility. Feature failure can degrade one view without stopping the portal.';

  for (const s of storage) {
    const name = pvcRef(s).split('/')[1];
    const existing = docs.filter(e => e.kind === 'Resource' && (e.metadata.namespace || 'default') === 'default' && e.metadata.name === name);
    if (existing.length > 1 || existing.some(e => e.metadata.annotations?.['onemind.io/inventory-source'] !== 'reviewed-topology')) throw new Error(`Unreviewed storage identity collision: ${name}`);
    const resource = {
      apiVersion: 'backstage.io/v1alpha1', kind: 'Resource',
      metadata: {
        name, title: `${s.namespace} / ${s.claim} (PVC)`,
        description: 'Existing Kubernetes PVC identity from reviewed workload wiring. Catalog metadata only, not a newly provisioned volume or live storage/backup health.',
        annotations: { 'onemind.io/inventory-source': 'reviewed-topology', 'onemind.io/kubernetes-namespace': s.namespace, 'onemind.io/kubernetes-pvc': s.claim, 'onemind.io/relationship-source': s.source },
        links: [{ title: 'Reviewed storage source', type: 'source', url: `https://github.com/onemindos/omos/blob/main/${s.source.split(' ')[0]}` }],
      },
      spec: { type: 'kubernetes-persistent-volume-claim', owner: 'group:default/onemind-admins', system: `cluster-${s.namespace}`, dependsOn: [cluster] },
    };
    if (existing.length) {
      const old = existing[0];
      if (old.spec.type !== resource.spec.type || old.metadata.annotations['onemind.io/kubernetes-namespace'] !== s.namespace || old.metadata.annotations['onemind.io/kubernetes-pvc'] !== s.claim) throw new Error(`Storage identity changed: ${name}`);
      const titles = new Set(resource.metadata.links.map(l => l.title));
      resource.metadata.annotations = { ...old.metadata.annotations, ...resource.metadata.annotations };
      resource.metadata.links = [...(old.metadata.links || []).filter(l => !titles.has(l.title)), ...resource.metadata.links];
      resource.spec.dependsOn = [...new Set([...(old.spec.dependsOn || []), ...resource.spec.dependsOn])];
      // Keep the System the structure review assigned; a new volume gets cluster-<namespace>, which curate-structure maps to the app's System.
      resource.spec.system = old.spec.system || resource.spec.system;
      docs[docs.indexOf(old)] = { ...old, metadata: { ...old.metadata, ...resource.metadata }, spec: { ...old.spec, ...resource.spec } };
    }
    else docs.push(resource);
    const component = exactly(docs, 'Component', s.component);
    component.spec.dependsOn = [...new Set([...component.spec.dependsOn, pvcRef(s)])];
  }
  // Snapshot-generated Systems group namespaces; their old "Cluster /" titles
  // misleadingly suggested dozens of clusters. Preserve IDs and system membership.
  for (const e of docs.filter(x => x.kind === 'System' && x.metadata.name.startsWith('cluster-'))) {
    const ns = e.metadata.name.slice('cluster-'.length);
    const members = docs.filter(x => x.kind === 'Component' && x.spec.system === e.metadata.name && x.metadata.annotations?.['backstage.io/kubernetes-namespace'] === ns);
    if (!members.length || members.some(x => !x.spec.dependsOn?.includes(cluster))) continue;
    e.metadata.title = `Namespace / ${ns}`;
    e.metadata.description = `Kubernetes namespace ${ns} on OneMind Cloud — NYC3. Inventory grouping, not an independent cluster or live health.`;
  }
  location.metadata.description = 'Current DigitalOcean NYC3 Kubernetes deployment location. Home/Mac/edge clusters are not claimed deployed by this graph. Dependency edges are reviewed catalog metadata, not a live traffic or health map.';
  location.metadata.links = [...(location.metadata.links || []).filter(l => l.title !== 'Core service topology (documented)'), { title: 'Core service topology (documented)', type: 'topology', url: graphUrl([names.portal, names.nats, names.database]) }];
  return docs;
}
function renderCluster(entities) {
  return entities.map(e => {
    const doc = new YAML.Document(e);
    // Long graph URLs need quoted escaping to respect changed-YAML line limits.
    // Keep their parsed value identical; leave unrelated scalar styles unchanged.
    for (const link of doc.getIn(['metadata', 'links'], true)?.items || []) {
      const url = link.get('url', true);
      if (typeof url?.value === 'string' && url.value.length > 180) url.type = YAML.Scalar.QUOTE_DOUBLE;
    }
    return doc.toString();
  }).join('---\n');
}
module.exports = { curateTopology, graphUrl, names, storage, pvcRef, renderCluster };
