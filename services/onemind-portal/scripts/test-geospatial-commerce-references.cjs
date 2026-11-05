// Source-level catalog and Home expression coverage, not rendered UI or integration validation.
const { test } = require('node:test');
const assert = require('node:assert/strict');
const fs = require('node:fs');
const path = require('node:path');
const YAML = require('yaml');
const root = path.join(__dirname, '..');
const text = fs.readFileSync(path.join(root, 'catalog/apps.yaml'), 'utf8');
const parse = source => YAML.parseAllDocuments(source).map(d => { assert.deepEqual(d.errors, []); return d.toJS(); });
const docs = parse(text);
const cluster = fs.readFileSync(path.join(root, 'catalog/cluster.yaml'), 'utf8');
const inventory = [...docs, ...parse(cluster)];
const structure = YAML.parse(fs.readFileSync(path.join(root, 'catalog/structure.yaml'), 'utf8'));
const references = [
  ['google-earth', 'https://earth.google.com', 'geospatial', 'Geospatial', ['geospatial'], 'so'],
  ['cesium-native', 'https://github.com/CesiumGS/cesium-native', 'geospatial', 'Geospatial', ['cesium', '3d-tiles'], 'so'],
  ['3d-tiles', 'https://github.com/CesiumGS/3d-tiles', 'geospatial', 'Geospatial', ['cesium', '3d-tiles'], 'so'],
  ['cesium-unreal', 'https://github.com/CesiumGS/cesium-unreal', 'geospatial', 'Geospatial', ['cesium', '3d-tiles', 'unreal'], 'so'],
  ['cesium-unity', 'https://github.com/CesiumGS/cesium-unity', 'geospatial', 'Geospatial', ['cesium', '3d-tiles', 'unity'], 'so'],
  ['cesium-omniverse', 'https://github.com/CesiumGS/cesium-omniverse', 'geospatial', 'Geospatial', ['cesium', '3d-tiles', 'omniverse'], 'so'],
  ['cesium-ion-for-tak', 'https://github.com/CesiumGS/cesium-ion-for-TAK-plugins', 'tak-ecosystem', 'TAK Ecosystem', ['tak-ecosystem'], 'so'],
  ['whop', 'https://whop.com', 'business-operations', 'Business & Commerce', ['business', 'commerce'], 'ge'],
];
for (const [name, url, system, category, tags, pillar] of references) {
  test(`${name}: one canonical reference with exact URL, category and durable system`, () => {
    const matches = inventory.filter(e => e.kind === 'Component' && e.metadata.name === name);
    assert.equal(matches.length, 1);
    const [e] = matches, a = e.metadata.annotations;
    assert.equal(e.spec.system, system);
    assert.equal(structure.placement[name], system);
    assert.equal(e.spec.type, 'website');
    assert.equal(e.spec.owner, 'group:default/onemind-admins');
    assert.equal(e.metadata.labels['onemind.io/pillar'], pillar);
    assert.equal(e.metadata.labels['onemind.io/catalog-role'], 'app');
    assert.equal(a['onemind.io/hosting'], 'external-provider-unverified');
    assert.equal(a['onemind.io/inventory-review'], 'imported-link');
    assert.equal(a['onemind.io/launcher-url'], url);
    assert.equal(a['onemind.io/launcher-category'], category);
    assert.deepEqual(e.metadata.links, [{ url, title: 'External reference', icon: 'dashboard' }]);
    assert.equal(inventory.flatMap(x => x.metadata.links || []).filter(l => l.url === url).length, 1);
    assert.equal(inventory.filter(x => x.metadata.annotations?.['onemind.io/launcher-url'] === url).length, 1);
    for (const tag of ['reference', ...tags]) assert.ok(e.metadata.tags.includes(tag), tag);
    assert.match(e.metadata.description, /^External reference/);
    if (name !== 'cesium-ion-for-tak') {
      assert.match(e.metadata.description, /No OneMind installation, hosting or integration verified\./);
    }
    const destination = new URL(url);
    assert.equal(destination.protocol, 'https:');
    assert.equal(destination.username, '');
    assert.equal(destination.password, '');
    for (const key of Object.keys(a)) assert.ok(!/^(argocd\/|backstage.io\/kubernetes-|onemind.io\/(verified-at|deployment))/.test(key), key);
  });
}
test('actual Home expressions admit all eight references and dynamically group Business & Commerce', () => {
  const hook = fs.readFileSync(path.join(root, 'plugins/plugin-onemind-common/src/useDirectory.ts'), 'utf8');
  const launch = new Function('entity', `return (${hook.match(/const raw =([\s\S]*?);/)[1]});`);
  const group = new Function('entity', 'known', `return (${hook.match(/const groupTitle =([\s\S]*?);/)[1]});`);
  const directory = references.map(([name, url, , category]) => {
    const entity = docs.find(e => e.kind === 'Component' && e.metadata.name === name);
    assert.equal(launch(entity), url);
    assert.equal(group(entity, { groupTitle: 'Wrong fallback' }), category);
    return { name, groupTitle: group(entity, undefined) };
  });
  const home = fs.readFileSync(path.join(root, 'plugins/plugin-onemind-home/src/homeModule.tsx'), 'utf8');
  const groups = new Function('directory', `return (${home.match(/const groups = (.*);/)[1]});`)(directory);
  assert.deepEqual(groups, ['Geospatial', 'TAK Ecosystem', 'Business & Commerce']);
  assert.deepEqual(directory.filter(e => e.groupTitle === 'Business & Commerce').map(e => e.name), ['whop']);
});
test('presentation categories reuse canonical systems; no new category Systems', () => {
  for (const system of ['geospatial', 'tak-ecosystem', 'business-operations']) {
    assert.equal(docs.filter(e => e.kind === 'System' && e.metadata.name === system).length, 1);
    assert.equal(structure.systems.filter(e => e.name === system).length, 1);
  }
  assert.equal(docs.some(e => e.kind === 'System' && e.metadata.name === 'business-commerce'), false);
});
test('all requested blocks remain byte-stable under structure regeneration', () => {
  const { applyStructureToText, loadSources } = require('./curate-structure.cjs');
  const out = applyStructureToText({ apps: text, cluster }, loadSources());
  for (const [name] of references) {
    const block = s => s.split('\n---\n').find(b => b.includes(`\n  name: ${name}\n`));
    assert.ok(block(text), name);
    assert.equal(block(out.apps), block(text), name);
  }
});
const mirror = path.join(root, '../../platform/backstage/catalog/apps.yaml');
test('source and mounted apps catalogs are byte-identical', { skip: !fs.existsSync(mirror) }, () => {
  assert.equal(fs.readFileSync(mirror, 'utf8'), text);
});
