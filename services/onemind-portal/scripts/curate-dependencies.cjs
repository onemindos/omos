// Reviewed manifest ownership, never namespace-only inference or live health.
const fs = require('node:fs');
const path = require('node:path');
const YAML = require('yaml');

const portalUrl = 'https://onemind-preview.tail717ba0.ts.net';
const reviewedWorkloads = [
  {
    name: 'k8s-backstage-backstage-21b99354',
    namespace: 'backstage',
    selector: 'app=backstage',
    application: 'platform-backstage',
    source: 'platform/backstage',
    runbook: 'resilience',
    runbookTitle: 'Portal recovery runbook',
    grafanaUid: 'onemind-backstage',
  },
  {
    name: 'k8s-backstage-backstage-postgres-17d2f351',
    namespace: 'backstage',
    selector: 'app=backstage-postgres',
    application: 'platform-backstage',
    source: 'platform/backstage',
    runbook: 'resilience',
    runbookTitle: 'Database recovery runbook',
    grafanaUid: 'onemind-backstage-postgres',
  },
  {
    name: 'k8s-nats-nats-0dee54ac',
    namespace: 'nats',
    selector:
      'app.kubernetes.io/component=nats,app.kubernetes.io/instance=nats,app.kubernetes.io/name=nats',
    application: 'platform-nats',
    source: 'platform/nats',
    runbook: 'nats',
    runbookTitle: 'NATS operations runbook',
    grafanaUid: 'onemind-nats',
  },
];
const normalizeSelector = value => (value || '').split(',').sort().join(',');

function curateDependencies(entities) {
  const docs = structuredClone(entities);
  for (const reviewed of reviewedWorkloads) {
    const matches = docs.filter(
      e =>
        e.kind === 'Component' &&
        (e.metadata.namespace || 'default') === 'default' &&
        e.metadata.name === reviewed.name,
    );
    const entity = matches[0];
    const annotations = entity?.metadata.annotations;
    if (
      matches.length !== 1 ||
      annotations?.['backstage.io/kubernetes-namespace'] !==
        reviewed.namespace ||
      normalizeSelector(
        annotations?.['backstage.io/kubernetes-label-selector'],
      ) !== normalizeSelector(reviewed.selector)
    ) {
      throw new Error(
        `Reviewed workload ${reviewed.name} missing, ambiguous or changed; review ownership before remapping`,
      );
    }
    Object.assign(annotations, {
      'grafana/dashboard-selector': `uid == '${reviewed.grafanaUid}'`,
      'argocd/app-name': reviewed.application,
      'argocd/instance-name': 'onemind-cloud',
      'github.com/project-slug': 'onemindos/omos',
      'backstage.io/source-location': `url:https://github.com/onemindos/omos/tree/main/${reviewed.source}`,
      'onemind.io/deployment-source': `${reviewed.source}; platform/argocd/${reviewed.application}.yaml`,
    });
    const entityUrl = `${portalUrl}/catalog/default/component/${reviewed.name}`;
    const links = [
      {
        title: 'Grafana workload dashboard',
        type: 'grafana',
        url: `https://grafana.tail717ba0.ts.net/d/${reviewed.grafanaUid}`,
      },
      {
        title: 'Live Kubernetes resources',
        type: 'kubernetes',
        url: `${entityUrl}/kubernetes`,
      },
      {
        title: 'Argo CD deployment summary',
        type: 'argocd',
        url: `${entityUrl}/deployment-summary`,
      },
      {
        title: reviewed.runbookTitle,
        type: 'runbook',
        url: `${portalUrl}/docs/default/component/onemind-runbooks/${reviewed.runbook}/`,
      },
    ];
    const managedTitles = new Set(links.map(link => link.title));
    entity.metadata.links = [
      ...(entity.metadata.links || []).filter(
        link => !managedTitles.has(link.title),
      ),
      ...links,
    ];
  }
  const portal = docs.find(
    e =>
      e.kind === 'Component' &&
      (e.metadata.namespace || 'default') === 'default' &&
      e.metadata.name === reviewedWorkloads[0].name,
  );
  const database = docs.find(
    e =>
      e.kind === 'Component' &&
      (e.metadata.namespace || 'default') === 'default' &&
      e.metadata.name === reviewedWorkloads[1].name,
  );
  portal.spec.dependsOn = [
    ...new Set([
      ...(portal.spec.dependsOn || []),
      `component:default/${database.metadata.name}`,
    ]),
  ];
  portal.metadata.annotations['onemind.io/relationship-source'] =
    'services/onemind-portal/app-config.production.yaml backend.database.connection.host; platform/backstage/postgres.yaml';
  return docs;
}

if (require.main === module) {
  const root = path.resolve(__dirname, '..');
  const docs = YAML.parseAllDocuments(
    fs.readFileSync(path.join(root, 'catalog/cluster.yaml'), 'utf8'),
  ).map(doc => doc.toJSON());
  const rendered = curateDependencies(docs)
    .map(e => YAML.stringify(e))
    .join('---\n');
  for (const base of ['catalog', '../../platform/backstage/catalog']) {
    fs.writeFileSync(path.join(root, base, 'cluster.yaml'), rendered);
  }
}
module.exports = { curateDependencies, reviewedWorkloads };
