// Source-level Home expressions only: not a rendered UI, API, feed or deployment check.
const { test } = require('node:test');
const assert = require('node:assert/strict');
const fs = require('node:fs');
const path = require('node:path');
const YAML = require('yaml');
const root = path.join(__dirname, '..');
const text = fs.readFileSync(path.join(root, 'catalog/apps.yaml'), 'utf8');
const clusterText = fs.readFileSync(path.join(root, 'catalog/cluster.yaml'), 'utf8');
const parse = source => YAML.parseAllDocuments(source).map(d => {
  assert.deepEqual(d.errors, []);
  return d.toJS();
});
const apps = parse(text), cluster = parse(clusterText), inventory = [...apps, ...cluster];
const structure = YAML.parse(fs.readFileSync(path.join(root, 'catalog/structure.yaml'), 'utf8'));
const references = [
  ['adsb-exchange', 'https://www.adsbexchange.com', 'intelligence', 'Intelligence', 'https://adsbexchange.com'],
  ['openstreetmap', 'https://www.openstreetmap.org', 'geospatial', 'Geospatial', 'https://www.openstreetmap.org'],
  ['opensky-network', 'https://opensky-network.org', 'intelligence', 'Intelligence', 'https://opensky-network.org'],
];
const entity = name => apps.find(e => e.kind === 'Component' && e.metadata.name === name);
for (const [name, url, system, category, launcher] of references) {
  test(`${name}: exact supplied URL on one canonical reference with durable Home placement`, () => {
    const holders = inventory.filter(e => e.metadata.links?.some(l => l.url === url));
    assert.equal(holders.length, 1, url);
    assert.equal(inventory.filter(e => e.kind === 'Component' && e.metadata.name === name).length, 1);
    const e = entity(name), a = e.metadata.annotations;
    assert.equal(holders[0], e);
    assert.equal(e.metadata.links.filter(l => l.url === url).length, 1);
    assert.equal(e.spec.system, system);
    assert.equal(structure.placement[name], system);
    assert.equal(e.spec.type, 'website');
    assert.equal(e.spec.owner, 'group:default/onemind-admins');
    assert.equal(e.metadata.labels['onemind.io/pillar'], 'so');
    assert.equal(e.metadata.labels['onemind.io/catalog-role'], 'app');
    assert.ok(e.metadata.tags.includes('reference'));
    assert.equal(a['onemind.io/hosting'], 'external-provider-unverified');
    assert.equal(a['onemind.io/launcher-url'], launcher);
    assert.equal(a['onemind.io/launcher-category'], category);
    for (const destination of [url, launcher]) {
      const u = new URL(destination);
      assert.equal(u.protocol, 'https:');
      assert.equal(u.username + u.password + u.hash + u.search, '');
    }
    for (const key of Object.keys(a)) assert.ok(!/^(argocd\/|backstage.io\/kubernetes-|onemind.io\/(verified-at|deployment|kubernetes-workloads))/.test(key));
    for (const key of ['dependsOn', 'subcomponentOf', 'providesApis', 'consumesApis']) assert.equal(e.spec[key], undefined);
  });
}
test('ADS-B original link, launcher and caveats survive the added www supporting link', () => {
  const e = entity('adsb-exchange');
  assert.deepEqual(e.metadata.links, [
    { url: 'https://adsbexchange.com', title: 'External reference', icon: 'dashboard' },
    { url: 'https://www.adsbexchange.com', title: 'Official website (www reference)', icon: 'dashboard' },
  ]);
  assert.equal(e.metadata.annotations['onemind.io/launcher-url'], 'https://adsbexchange.com');
  assert.match(e.metadata.description, /future CoT bridge requires API, licensing and redistribution review/);
  assert.match(e.metadata.description, /no complete coverage or flight-safety suitability claimed/);
  assert.match(e.metadata.description, /No OneMind installation, hosting or integration verified/);
});
test('new references distinguish collaborative maps from aircraft research and disclaim activation', () => {
  assert.match(entity('openstreetmap').metadata.description, /Collaborative geographic data and maps/);
  assert.match(entity('opensky-network').metadata.description, /Aircraft-tracking research network/);
  assert.match(entity('opensky-network').metadata.description, /not an OpenStreetMap or geocoding service/);
  for (const name of ['openstreetmap', 'opensky-network']) {
    const e = entity(name);
    assert.match(e.metadata.description, /No approved API feed, complete coverage or flight-safety suitability claimed/);
    assert.match(e.metadata.description, /No OneMind deployment, compatibility or integration verified/);
    assert.equal(e.metadata.links.length, 1);
  }
});
test('archived registry Resources remain distinct, not duplicate Home Components', () => {
  for (const [name, archived, url] of [
    ['openstreetmap', 'fabric-openstreetmap-91ae68a8', 'https://openstreetmap.org'],
    ['opensky-network', 'fabric-opensky-20188a97', 'https://opensky-network.org'],
  ]) {
    const a = cluster.find(e => e.metadata.name === archived);
    assert.equal(a.kind, 'Resource');
    assert.equal(a.spec.type, 'documented-capability');
    assert.equal(a.metadata.annotations['onemind.io/documented-address'], url);
    assert.equal(a.metadata.annotations['onemind.io/launcher-url'], undefined);
    assert.notEqual(entity(name).metadata.title, a.metadata.title);
  }
});
test('actual Home source expressions preserve original launcher and select requested categories', () => {
  const hook = fs.readFileSync(path.join(root, 'plugins/plugin-onemind-common/src/useDirectory.ts'), 'utf8');
  const launch = new Function('entity', `return (${hook.match(/const raw =([\s\S]*?);/)[1]});`);
  const group = new Function('entity', 'known', `return (${hook.match(/const groupTitle =([\s\S]*?);/)[1]});`);
  const directory = references.map(([name, , , category, launcher]) => {
    assert.equal(launch(entity(name)), launcher);
    assert.equal(group(entity(name), { groupTitle: 'Wrong fallback' }), category);
    assert.equal(group(entity(name), undefined), category);
    return { name, groupTitle: category };
  });
  const home = fs.readFileSync(path.join(root, 'plugins/plugin-onemind-home/src/homeModule.tsx'), 'utf8');
  const groups = new Function('directory', `return (${home.match(/const groups = (.*);/)[1]});`)(directory);
  assert.deepEqual(groups, ['Intelligence', 'Geospatial']);
});
test('canonical reference blocks are byte-stable under durable structure regeneration', () => {
  const { applyStructureToText, loadSources } = require('./curate-structure.cjs');
  const out = applyStructureToText({ apps: text, cluster: clusterText }, loadSources());
  for (const [name] of references) {
    const block = source => source.split('---\n').find(b => b.includes(`  name: ${name}\n`));
    assert.ok(block(text), name);
    assert.equal(block(out.apps), block(text), name);
  }
});
const mirror = path.join(root, '../../platform/backstage/catalog/apps.yaml');
test('source and mounted catalogs are identical', { skip: !fs.existsSync(mirror) }, () => {
  assert.equal(fs.readFileSync(mirror, 'utf8'), text);
});
