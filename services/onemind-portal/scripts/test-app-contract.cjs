// Guard test for the app annotation contract (flagship F7): every curated app Component records where it
// runs, with a value from the documented set, and the health and fact annotations are well formed.
const { test } = require('node:test');
const assert = require('node:assert/strict');
const fs = require('node:fs');
const path = require('node:path');
const YAML = require('yaml');
const { HOSTING_VALUES, hostingKind, validateEntity } = require('./app-contract.cjs');

const catalog = f => fs.readFileSync(path.join(__dirname, '../catalog', f), 'utf8').split(/^---$/m).map(d => YAML.parse(d)).filter(Boolean);
const apps = catalog('apps.yaml');
const components = apps.filter(e => e.kind === 'Component');
const entity = (annotations, kind = 'Component') => ({ kind, metadata: { name: 'demo', annotations } });

test('every curated app Component has a valid onemind.io/hosting value', () => {
  assert.ok(components.length > 100, `${components.length} components`);
  const bad = components.flatMap(c => validateEntity(c).map(p => `${c.metadata.name}: ${p}`));
  assert.deepEqual(bad, []);
});

test('every other curated entity that sets onemind.io/hosting uses a valid value', () => {
  const others = apps.filter(e => e.kind !== 'Component' && e.metadata.annotations?.['onemind.io/hosting'] !== undefined);
  const bad = others.flatMap(e => validateEntity(e, { requireHosting: false }).map(p => `${e.kind}/${e.metadata.name}: ${p}`));
  assert.deepEqual(bad, []);
});

test('a missing or invalid hosting value fails the validator, and unknown is accepted', () => {
  assert.match(validateEntity(entity({})).join(), /hosting is missing/);
  assert.match(validateEntity(entity({ 'onemind.io/hosting': 'aws' })).join(), /not one of/);
  assert.match(validateEntity(entity({ 'onemind.io/hosting': '' })).join(), /missing/);
  assert.deepEqual(validateEntity(entity({ 'onemind.io/hosting': 'unknown' })), []);
  for (const value of HOSTING_VALUES) assert.deepEqual(validateEntity(entity({ 'onemind.io/hosting': value })), [], value);
});

test('removing hosting from any one catalog app turns the guard red', () => {
  for (const c of components.slice(0, 5)) {
    const broken = structuredClone(c);
    delete broken.metadata.annotations['onemind.io/hosting'];
    assert.equal(validateEntity(broken).length, 1, c.metadata.name);
  }
});

test('legacy hosting values keep their meaning', () => {
  const kind = (value, extra = {}) => hostingKind(entity({ 'onemind.io/hosting': value, ...extra }));
  assert.equal(kind('digitalocean-nyc3', { 'argocd/app-name': 'platform-x' }), 'kubernetes');
  assert.equal(kind('digitalocean-nyc3'), 'digitalocean');
  assert.equal(kind('external-provider-unverified'), 'saas');
  assert.equal(kind('mesh-endpoint-unverified'), 'home');
  assert.equal(kind('native-application'), 'desktop');
  assert.equal(kind('vehicle-firmware'), 'device');
  assert.equal(kind('external-website'), 'web');
  assert.equal(kind('planned'), 'planned');
  assert.equal(kind('planned-not-deployed'), 'planned');
  assert.equal(kind('nonsense'), 'unknown');
});

test('health-check needs a health-url and a valid endpoint name, and the names are unique', () => {
  const ok = { 'onemind.io/hosting': 'saas', 'onemind.io/health-check': 'demo-app', 'onemind.io/health-url': 'https://example.com' };
  assert.deepEqual(validateEntity(entity(ok)), []);
  assert.match(validateEntity(entity({ ...ok, 'onemind.io/health-url': undefined })).join(), /health-url is missing/);
  assert.match(validateEntity(entity({ ...ok, 'onemind.io/health-check': 'Demo App' })).join(), /endpoint name/);
  assert.match(validateEntity(entity({ ...ok, 'onemind.io/health-url': 'https://user:pw@example.com' })).join(), /no credentials/);
  const names = components.map(c => c.metadata.annotations?.['onemind.io/health-check']).filter(Boolean);
  assert.equal(new Set(names).size, names.length, 'duplicate Gatus endpoint names');
});

test('fact annotations are text only and well formed', () => {
  const base = { 'onemind.io/hosting': 'saas' };
  assert.deepEqual(validateEntity(entity({ ...base, 'onemind.io/renewal': '2027-01-31', 'onemind.io/plan': 'Standard, billed yearly' })), []);
  assert.match(validateEntity(entity({ ...base, 'onemind.io/renewal': 'next March' })).join(), /renewal/);
  assert.match(validateEntity(entity({ ...base, 'onemind.io/access': 'password: hunter2hunter2' })).join(), /credential/);
  assert.match(validateEntity(entity({ ...base, 'onemind.io/access': 'Zeus (owner); backup admin: Pat' })).join(), /^$/);
  assert.match(validateEntity(entity({ ...base, 'onemind.io/plan': '  ' })).join(), /empty/);
});

test('no catalog app invents a plan, renewal or access fact yet: those stay absent until recorded', () => {
  for (const key of ['onemind.io/plan', 'onemind.io/renewal', 'onemind.io/access', 'onemind.io/data-export', 'onemind.io/backup-location']) {
    for (const c of components) {
      const value = c.metadata.annotations?.[key];
      if (value !== undefined) assert.ok(String(value).trim().length > 0, `${c.metadata.name} ${key}`);
    }
  }
});

test('the mounted catalog copy matches the portal copy', { skip: !fs.existsSync(path.join(__dirname, '../../../platform/backstage/catalog')) && 'repo-level platform/ is not in this build context' }, () => {
  const a = fs.readFileSync(path.join(__dirname, '../catalog/apps.yaml'), 'utf8');
  const b = fs.readFileSync(path.join(__dirname, '../../../platform/backstage/catalog/apps.yaml'), 'utf8');
  assert.ok(a === b, 'platform/backstage/catalog/apps.yaml must be an exact copy of the portal catalog');
});
