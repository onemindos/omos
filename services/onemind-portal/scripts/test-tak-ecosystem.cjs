// Exact owner-supplied external references; no deployment or affiliation implied.
const { test } = require('node:test');
const assert = require('node:assert/strict');
const fs = require('node:fs');
const path = require('node:path');
const YAML = require('yaml');
const root = path.join(__dirname, '..');
const text = fs.readFileSync(path.join(root, 'catalog/apps.yaml'), 'utf8');
const docs = YAML.parseAllDocuments(text).map(doc => {
  assert.deepEqual(doc.errors, []);
  return doc.toJS();
});
const placement = YAML.parse(fs.readFileSync(path.join(root, 'catalog/structure.yaml'), 'utf8')).placement;
const links = {
  district47: 'https://district47.us',
  cotak: 'https://cotak.gov/pages/guest-home',
  gotak: 'https://getgotak.com',
  takgrid: 'https://takgrid.com',
  'tak-product-center': 'https://github.com/TAK-Product-Center',
  pytak: 'https://pytak.readthedocs.io/en/latest/',
  eolianvr: 'https://www.eolianvr.com',
  'omnitak-external': 'https://omnitak.engindearing.soy',
  'engindearing-projects': 'https://github.com/engindearing-projects',
};
for (const [name, url] of Object.entries(links)) {
  test(`${name}: exact URL once, canonical external entity in TAK ecosystem`, () => {
    const matches = docs.filter(e => e.kind === 'Component' && e.metadata.name === name);
    assert.equal(matches.length, 1);
    const [entity] = matches;
    assert.equal(entity.spec.system, 'tak-ecosystem');
    assert.equal(placement[name], 'tak-ecosystem');
    assert.equal(entity.metadata.labels['onemind.io/pillar'], 'so');
    assert.equal(entity.metadata.annotations['onemind.io/hosting'], 'external-provider-unverified');
    assert.equal(docs.flatMap(e => e.metadata.links || []).filter(link => link.url === url).length, 1);
    assert.ok(entity.metadata.links.some(link => link.url === url));
    // Match useDirectory's launcher eligibility, not merely catalog ingestion.
    const launcherUrl = entity.metadata.annotations['onemind.io/launcher-url'] ??
      entity.metadata.links.find(link => link.title === 'Open application')?.url;
    assert.ok(launcherUrl, `${name} is ingested but excluded from the launcher`);
    if (name === 'district47') {
      assert.equal(entity.metadata.annotations['onemind.io/launcher-url'], url);
      assert.equal(entity.metadata.annotations['onemind.io/launcher-category'], 'TAK Ecosystem');
    }
    const parsed = new URL(launcherUrl);
    assert.equal(parsed.protocol, 'https:');
    assert.equal(parsed.username, '');
    assert.equal(parsed.password, '');
    const curated = JSON.parse(fs.readFileSync(path.join(root,
      'plugins/plugin-onemind-common/directory.json'), 'utf8'));
    const known = curated.find(item => item.name === name);
    const category = entity.metadata.annotations['onemind.io/launcher-category'] ??
      known?.groupTitle ?? 'Services';
    assert.equal(category, 'TAK Ecosystem', `${name} must match the home category dropdown`);
  });
}
test('XMPP for TAK is a deferred Home planning card, not a running endpoint', () => {
  const matches = docs.filter(e => e.kind === 'Component' && e.metadata.name === 'xmpp-for-tak');
  assert.equal(matches.length, 1);
  const [entity] = matches;
  assert.equal(entity.metadata.title, 'XMPP for TAK');
  assert.equal(entity.metadata.description, 'Self-hosted messaging for ATAK’s TAK Chat plugin and standard desktop/mobile XMPP clients, with persistent rooms, direct messages and a future Hermes bridge. Deferred; plugin compatibility and integration remain untested.');
  assert.equal(entity.spec.system, 'tak-ecosystem');
  assert.equal(placement['xmpp-for-tak'], 'tak-ecosystem');
  assert.equal(entity.metadata.labels['onemind.io/pillar'], 'so');
  assert.equal(entity.spec.lifecycle, 'experimental');
  const annotations = entity.metadata.annotations;
  assert.equal(annotations['onemind.io/documented-status'], 'Planned — deferred, not deployed; plugin compatibility and integration untested');
  assert.equal(annotations['onemind.io/hosting'], 'planned-not-deployed');
  assert.equal(annotations['onemind.io/launcher-category'], 'TAK Ecosystem');
  // Production app.baseUrl comes from PORTAL_URL, not the localhost dev config.
  const production = YAML.parse(fs.readFileSync(path.join(root, 'app-config.production.yaml'), 'utf8'));
  assert.equal(production.app.baseUrl, '${PORTAL_URL}');
  const deployment = path.join(root, '../../platform/backstage/deployment.yaml');
  const baseUrl = fs.existsSync(deployment)
    ? YAML.parseAllDocuments(fs.readFileSync(deployment, 'utf8')).map(doc => doc.toJS())
      .find(doc => doc.kind === 'Deployment' && doc.metadata.name === 'backstage').spec.template.spec.containers
      .flatMap(container => container.env || []).find(env => env.name === 'PORTAL_URL').value
    : 'https://onemind-preview.tail717ba0.ts.net'; // Repo absent in image build context.
  assert.equal(new URL(baseUrl).protocol, 'https:');
  const destination = `${baseUrl}/catalog/default/component/xmpp-for-tak`;
  assert.equal(annotations['onemind.io/launcher-url'], destination);
  assert.deepEqual(entity.metadata.links.map(link => link.url), [destination]);
  for (const key of Object.keys(annotations)) {
    assert.ok(!/^(argocd\/|backstage.io\/kubernetes-|onemind.io\/(verified-at|deployment))/.test(key), key);
  }
});
test('District47 is a reference, with no OneMind affiliation or integration claimed', () => {
  const entity = docs.find(e => e.metadata.name === 'district47');
  assert.equal(entity.metadata.title, 'District47');
  assert.ok(entity.metadata.tags.includes('reference'));
  assert.match(entity.metadata.description, /External reference.*ATAK plugins and KMZ data packages/);
  assert.match(entity.metadata.description, /No OneMind affiliation or integration implied/);
});
test('external OmniTAK is explicitly unrelated to the retired OneMind name', () => {
  const entity = docs.find(e => e.metadata.name === 'omnitak-external');
  assert.match(entity.metadata.title, /external/i);
  assert.match(entity.metadata.description, /third-party.*unrelated to OneMind.*retired OmniTAK name/);
});
test('TAK Product Center retains its existing repository-list link', () => {
  const entity = docs.find(e => e.metadata.name === 'tak-product-center');
  assert.ok(entity.metadata.links.some(link => link.url === 'https://github.com/orgs/TAK-Product-Center/repositories'));
});
