// Source-level Home launcher coverage, not rendered UI or provider validation.
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
const structure = YAML.parse(fs.readFileSync(path.join(root, 'catalog/structure.yaml'), 'utf8'));
const curated = JSON.parse(fs.readFileSync(path.join(root, 'plugins/plugin-onemind-common/directory.json'), 'utf8'));
// Only add HTTPS to scheme-less user input; never replace repository links with demos.
const references = {
  osiris: ['https://github.com/simplifaisoul/osiris', ['osint', 'geospatial-intelligence', 'situational-awareness']],
  worldwideview: ['https://github.com/silvertakana/worldwideview', ['geospatial-intelligence', 'situational-awareness']],
  worldmonitor: ['https://github.com/koala73/worldmonitor', ['osint', 'situational-awareness']],
  'argos-atlas': ['https://argosatlas.com', ['osint', 'geospatial-intelligence', 'situational-awareness']],
  'gods-eye-view': ['https://github.com/bilawalsidhu/gods-eye-view', ['osint', 'geospatial-intelligence', 'situational-awareness']],
};
for (const [name, [url, tags]] of Object.entries(references)) {
  test(`${name}: canonical external Intelligence reference with exact Home destination`, () => {
    const matches = docs.filter(e => e.kind === 'Component' && e.metadata.name === name);
    assert.equal(matches.length, 1);
    const [entity] = matches;
    const annotations = entity.metadata.annotations;
    assert.equal(entity.spec.system, 'intelligence');
    assert.equal(structure.placement[name], 'intelligence');
    assert.equal(entity.spec.type, 'website');
    assert.equal(entity.spec.owner, 'group:default/onemind-admins');
    assert.equal(entity.metadata.labels['onemind.io/pillar'], 'so');
    assert.equal(entity.metadata.labels['onemind.io/catalog-role'], 'app');
    assert.equal(annotations['onemind.io/hosting'], 'external-provider-unverified');
    assert.equal(annotations['onemind.io/inventory-review'], 'imported-link');
    assert.equal(annotations['onemind.io/launcher-url'], url);
    assert.equal(annotations['onemind.io/launcher-category'], 'Intelligence');
    assert.deepEqual(entity.metadata.links.map(link => link.url), [url]);
    assert.equal(docs.flatMap(e => e.metadata.links || []).filter(link => link.url === url).length, 1);
    for (const tag of ['reference', ...tags]) assert.ok(entity.metadata.tags.includes(tag), tag);
    assert.match(entity.metadata.description, /^External reference/);
    assert.match(entity.metadata.description, /No OneMind deployment, compatibility or integration verified/);
    // Model useDirectory's actual eligibility and category precedence.
    const launcher = annotations['onemind.io/launcher-url'] ??
      entity.metadata.links.find(link => link.title === 'Open application')?.url;
    assert.equal(launcher, url);
    const destination = new URL(launcher);
    assert.equal(destination.protocol, 'https:');
    assert.equal(destination.username, '');
    assert.equal(destination.password, '');
    const category = annotations['onemind.io/launcher-category'] ??
      curated.find(item => item.name === name)?.groupTitle ?? 'Services';
    assert.equal(category, 'Intelligence');
    for (const key of Object.keys(annotations)) {
      assert.ok(!/^(argocd\/|backstage.io\/kubernetes-|onemind.io\/(verified-at|deployment))/.test(key), key);
    }
  });
}
test('Intelligence reuses the existing SO system and Home category', () => {
  const systems = docs.filter(e => e.kind === 'System' && e.metadata.name === 'intelligence');
  assert.equal(systems.length, 1);
  assert.equal(systems[0].spec.domain, 'domain:default/onemind-so');
  assert.equal(structure.systems.find(s => s.name === 'intelligence').pillar, 'so');
  assert.ok(curated.some(item => item.groupTitle === 'Intelligence'));
});
test("God's Eye View discloses simulated traffic instead of promising all-live data", () => {
  const entity = docs.find(e => e.metadata.name === 'gods-eye-view');
  assert.match(entity.metadata.description, /Traffic is simulated/);
});
test('new reference cards are byte-stable under durable structure regeneration', () => {
  const { applyStructureToText, loadSources } = require('./curate-structure.cjs');
  const cluster = fs.readFileSync(path.join(root, 'catalog/cluster.yaml'), 'utf8');
  const result = applyStructureToText({ apps: text, cluster }, loadSources());
  for (const name of Object.keys(references)) {
    const block = source => source.split('\n---\n')
      .find(part => part.includes(`\n  name: ${name}\n`));
    assert.equal(block(result.apps), block(text), name);
  }
});
const mirror = path.join(root, '../../platform/backstage/catalog/apps.yaml');
test('mounted catalog is byte-identical to source', {
  skip: !fs.existsSync(mirror) && 'repo-level platform/ is absent in image build context',
}, () => {
  assert.equal(fs.readFileSync(mirror, 'utf8'), text);
});
