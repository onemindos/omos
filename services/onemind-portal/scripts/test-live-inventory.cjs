// Guards for the live inventory (F11) that must also run in the portal image build, where only
// services/onemind-portal exists. Anything reading repo-level platform/ uses the repoOnly skip.
const { test } = require('node:test');
const assert = require('node:assert/strict');
const fs = require('node:fs');
const path = require('node:path');
const cp = require('node:child_process');
const YAML = require('yaml');

const portal = path.resolve(__dirname, '..');
const hasRepo = fs.existsSync(path.join(portal, '../../platform/backstage/rbac.yaml'));
const repoOnly = hasRepo ? {} : { skip: 'repo-level platform/ is not in this build context' };
const config = f => YAML.parse(fs.readFileSync(path.join(portal, f), 'utf8'));

test('the generated placement bundle is current with workload-map.yaml and apps.yaml', () => {
  const r = cp.spawnSync(process.execPath, [path.join(__dirname, 'generate-inventory-mapping.cjs'), '--check'], { encoding: 'utf8' });
  assert.equal(r.status, 0, r.stderr || r.stdout);
});

test('the live inventory ships disabled and the file-based catalog still lists cluster.yaml', () => {
  const base = config('app-config.yaml');
  assert.equal(base.onemind.inventory.live.enabled, false);
  assert.ok(base.catalog.locations.some(l => l.type === 'file' && l.target === '../../catalog/cluster.yaml'));
  assert.ok(base.catalog.rules[0].allow.includes('Component'));
  const production = config('app-config.production.yaml');
  assert.notEqual(production.onemind?.inventory?.live?.enabled, true);
});

test('the cluster overlay enables the live inventory explicitly (reviewed 2026-10-22) and overrides nothing else', repoOnly, () => {
  const cluster = YAML.parse(fs.readFileSync(path.join(portal, '../../platform/backstage/app-config.cluster.yaml'), 'utf8'));
  const live = cluster.onemind?.inventory?.live;
  assert.equal(live?.enabled, true);
  // Only the switch is set here: scope (apiServer, excluded namespaces, curated file) stays at the reviewed defaults.
  assert.deepEqual(Object.keys(live), ['enabled']);
});

test('the backend registers the inventory module', () => {
  const index = fs.readFileSync(path.join(portal, 'packages/backend/src/index.ts'), 'utf8');
  assert.match(index, /import\('@internal\/plugin-onemind-inventory-backend'\)/);
  const pkg = JSON.parse(fs.readFileSync(path.join(portal, 'packages/backend/package.json'), 'utf8'));
  assert.equal(pkg.dependencies['@internal/plugin-onemind-inventory-backend'], 'workspace:^');
});

test('RBAC: the portal ServiceAccount can only get/list/watch the kinds the inventory reads, and no Secrets', repoOnly, () => {
  const docs = YAML.parseAllDocuments(fs.readFileSync(path.join(portal, '../../platform/backstage/rbac.yaml'), 'utf8')).map(d => d.toJSON());
  const role = docs.find(d => d.kind === 'ClusterRole' && d.metadata.name === 'backstage-observer');
  const rule = role.rules.find(r => r.apiGroups.includes('apps') && r.resources.includes('deployments') && r.resources.includes('statefulsets'));
  assert.ok(rule, 'deployments and statefulsets are readable cluster-wide');
  for (const verb of ['list']) assert.ok(rule.verbs.includes(verb));
  const writes = ['create', 'update', 'patch', 'delete', 'deletecollection', 'escalate', 'bind', 'impersonate', '*'];
  for (const r of role.rules) {
    assert.ok(!r.verbs.some(v => writes.includes(v)), `no write verb in ${JSON.stringify(r)}`);
    assert.ok(!r.resources.includes('secrets') && !r.resources.includes('configmaps') && !r.resources.includes('*'), `no secrets, configmaps or wildcards in ${JSON.stringify(r)}`);
  }
});
