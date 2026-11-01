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
  });
}
test('external OmniTAK is explicitly unrelated to the retired OneMind name', () => {
  const entity = docs.find(e => e.metadata.name === 'omnitak-external');
  assert.match(entity.metadata.title, /external/i);
  assert.match(entity.metadata.description, /third-party.*unrelated to OneMind.*retired OmniTAK name/);
});
test('TAK Product Center retains its existing repository-list link', () => {
  const entity = docs.find(e => e.metadata.name === 'tak-product-center');
  assert.ok(entity.metadata.links.some(link => link.url === 'https://github.com/orgs/TAK-Product-Center/repositories'));
});
