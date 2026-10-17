// Reviewed architecture facts from committed manifests, not inferred from co-location.
const fs = require('fs'), YAML = require('yaml');
const docs = YAML.parseAllDocuments(fs.readFileSync('catalog/cluster.yaml', 'utf8')).map(d => d.toJSON());
const portal = docs.find(e => e.kind === 'Component' && e.metadata.title === 'backstage / backstage');
const database = docs.find(e => e.kind === 'Component' && e.metadata.title === 'backstage / backstage-postgres');
if (!portal || !database) throw new Error('Reviewed Backstage workloads missing; inspect inventory before changing dependencies');
portal.spec.dependsOn = [...new Set([...(portal.spec.dependsOn || []), `component:default/${database.metadata.name}`])];
portal.metadata.annotations['argocd/app-name'] = 'platform-backstage';
portal.metadata.annotations['argocd/instance-name'] = 'onemind-cloud';
portal.metadata.annotations['github.com/project-slug'] = 'onemindos/omos';
portal.metadata.annotations['onemind.io/relationship-source'] = 'services/onemind-portal/app-config.production.yaml backend.database.connection.host; platform/backstage/postgres.yaml';
portal.metadata.links = [...(portal.metadata.links || []).filter(l => l.type !== 'runbook'), {title:'Portal recovery runbook',type:'runbook',url:'https://onemind-preview.tail717ba0.ts.net/docs/default/component/onemind-runbooks/resilience/'}];
database.metadata.links = [...(database.metadata.links || []).filter(l => l.type !== 'runbook'), {title:'Database recovery runbook',type:'runbook',url:'https://onemind-preview.tail717ba0.ts.net/docs/default/component/onemind-runbooks/resilience/'}];
const rendered = docs.map(e => YAML.stringify(e)).join('---\n');
for (const base of ['catalog','../../platform/backstage/catalog']) fs.writeFileSync(`${base}/cluster.yaml`, rendered);
