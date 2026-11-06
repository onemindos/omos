// Guard test for the Gatus endpoint file (flagship F7): platform/gatus/endpoints.yaml is generated from the
// catalog, so it must match what the generator produces now, and every endpoint must be sane.
const { test } = require('node:test');
const assert = require('node:assert/strict');
const fs = require('node:fs');
const path = require('node:path');
const YAML = require('yaml');
const { buildEndpoints, render, generate, OUTPUT } = require('./generate-gatus-endpoints.cjs');

const hasRepo = fs.existsSync(path.join(__dirname, '../../../platform/backstage/catalog'));
const repoOnly = hasRepo ? {} : { skip: 'repo-level platform/ is not in this build context' };
const entity = (name, annotations) => ({ kind: 'Component', metadata: { name, annotations: { 'onemind.io/hosting': 'saas', ...annotations } } });

test('the committed endpoints.yaml is exactly what the generator produces from the catalog', repoOnly, () => {
  assert.equal(fs.readFileSync(OUTPUT, 'utf8'), generate());
});

test('the committed endpoints are valid Gatus config with at least ten real endpoints', repoOnly, () => {
  const { endpoints } = YAML.parse(fs.readFileSync(OUTPUT, 'utf8'));
  assert.ok(endpoints.length >= 10, `${endpoints.length} endpoints`);
  const names = new Set();
  for (const e of endpoints) {
    assert.match(e.name, /^[a-z0-9-]+$/);
    assert.ok(!names.has(e.name), `duplicate ${e.name}`);
    names.add(e.name);
    assert.match(e.url, /^https?:\/\//);
    assert.ok(e.conditions.includes('[STATUS] < 400'));
    if (e.url.startsWith('https://')) assert.ok(e.conditions.includes('[CERTIFICATE_EXPIRATION] > 336h'), `${e.name} needs the 14 day certificate condition`);
    assert.notEqual(e.name, 'test-deliberate-failure', 'the deliberate failure lives in test-endpoints.yaml, not in the generated file');
  }
});

test('an entity with health-check and health-url becomes an endpoint; without a url it fails loudly', () => {
  const list = buildEndpoints([
    entity('zoho-one', { 'onemind.io/health-check': 'zoho-one', 'onemind.io/health-url': 'https://www.zoho.com' }),
    entity('plain', {}),
    { kind: 'System', metadata: { name: 'ignored', annotations: { 'onemind.io/health-check': 'ignored', 'onemind.io/health-url': 'https://x.test' } } },
  ]);
  assert.deepEqual(list.map(e => e.name), ['zoho-one']);
  assert.equal(list[0].interval, '5m');
  assert.ok(list[0].conditions.includes('[CERTIFICATE_EXPIRATION] > 336h'));
  assert.throws(() => buildEndpoints([entity('a', { 'onemind.io/health-check': 'a' })]), /without onemind.io\/health-url/);
  assert.throws(() => buildEndpoints([
    entity('a', { 'onemind.io/health-check': 'x', 'onemind.io/health-url': 'http://a.svc' }),
    entity('b', { 'onemind.io/health-check': 'x', 'onemind.io/health-url': 'http://b.svc' }),
  ]), /used by/);
});

test('http endpoints get a response time condition and no certificate condition', () => {
  const [e] = buildEndpoints([entity('app', { 'onemind.io/hosting': 'digitalocean-nyc3', 'argocd/app-name': 'x', 'onemind.io/health-check': 'app', 'onemind.io/health-url': 'http://app.ns.svc.cluster.local/health' })]);
  assert.equal(e.group, 'kubernetes');
  assert.equal(e.interval, '1m');
  assert.deepEqual(e.conditions, ['[STATUS] < 400', '[RESPONSE_TIME] < 3000']);
  assert.match(render([e]), /^endpoints:\n {2}# catalog component: app\n {2}- name: app\n/m);
});
