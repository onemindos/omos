// Source-level reference and Home-expression guards, not rendered UI or integration tests.
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
// Scheme-less owner input: add ONLY https://, preserving domain/path even after redirects.
const references = [
  ['opencage', 'opencagedata.com'],
  ['graphhopper', 'graphhopper.com'],
  ['overpass-api', 'overpass-api.de'],
  ['overpass-turbo', 'github.com/tyrasd/overpass-turbo', 'github.com/cipher387/awesome-overpass-turbo'],
  ['valhalla', 'github.com/valhalla/valhalla', 'valhalla.openstreetmap.de'],
  ['nominatim-upstream', 'nominatim.org'],
  ['organic-maps', 'organicmaps.app'],
  ['maputnik', 'maputnik.github.io'],
  ['suncalc', 'suncalc.org'],
  ['komoot', 'komoot.com'],
  ['adsb-exchange', 'adsbexchange.com'],
];
const entity = name => apps.find(e => e.kind === 'Component' && e.metadata.name === name);
for (const [name, ...addresses] of references) {
  test(`${name}: exact references, safe launcher, durable Home classification, no runtime claims`, () => {
    const matches = inventory.filter(e => e.kind === 'Component' && e.metadata.name === name);
    assert.equal(matches.length, 1);
    const [e] = matches, a = e.metadata.annotations;
    const system = name === 'adsb-exchange' ? 'intelligence' : 'geospatial';
    assert.equal(e.spec.system, system);
    assert.equal(structure.placement[name], system);
    assert.equal(e.spec.type, 'website');
    assert.equal(e.spec.owner, 'group:default/onemind-admins');
    assert.equal(e.metadata.labels['onemind.io/pillar'], 'so');
    assert.equal(e.metadata.labels['onemind.io/catalog-role'], 'app');
    assert.ok(e.metadata.tags.includes('reference'));
    assert.equal(a['onemind.io/hosting'], 'external-provider-unverified');
    assert.equal(a['onemind.io/inventory-review'], 'imported-link');
    assert.equal(a['onemind.io/launcher-url'], `https://${addresses[0]}`);
    assert.equal(a['onemind.io/launcher-category'], system === 'intelligence' ? 'Intelligence' : 'Geospatial');
    assert.deepEqual(e.metadata.links.map(l => l.url), addresses.map(u => `https://${u}`));
    for (const url of e.metadata.links.map(l => l.url)) {
      assert.equal(inventory.flatMap(x => x.metadata.links || []).filter(l => l.url === url).length, 1, url);
      const parsed = new URL(url);
      assert.equal(parsed.protocol, 'https:');
      assert.equal(parsed.username, '');
      assert.equal(parsed.password, '');
      assert.equal(parsed.hash, '');
      assert.equal(parsed.search, '');
    }
    assert.equal(inventory.filter(x => x.metadata.annotations?.['onemind.io/launcher-url'] === `https://${addresses[0]}`).length, 1);
    assert.match(e.metadata.description, /^External reference\./);
    assert.match(e.metadata.description, /No OneMind installation, hosting or integration verified\./);
    for (const key of Object.keys(a)) assert.ok(!/^(argocd\/|backstage.io\/kubernetes-|onemind.io\/(verified-at|deployment|kubernetes-workloads))/.test(key), key);
    for (const key of ['subcomponentOf', 'dependsOn', 'dependencyOf', 'providesApis', 'consumesApis']) assert.equal(e.spec[key], undefined);
  });
}
test('all thirteen exact owner links are covered on eleven cards; examples/demo are not separate launchers', () => {
  const requested = references.flatMap(([, ...urls]) => urls.map(u => `https://${u}`));
  assert.equal(requested.length, 13);
  assert.equal(new Set(requested).size, 13);
  for (const url of requested) {
    const holders = inventory.filter(e => e.metadata.links?.some(l => l.url === url));
    assert.equal(holders.length, 1, url);
    assert.ok(references.some(([name]) => holders[0] === entity(name)));
  }
  for (const url of ['https://github.com/cipher387/awesome-overpass-turbo', 'https://valhalla.openstreetmap.de']) {
    assert.equal(inventory.some(e => e.metadata.annotations?.['onemind.io/launcher-url'] === url), false);
  }
  assert.equal(entity('overpass-turbo').metadata.links[1].title, 'Community examples and learning resources');
  assert.match(entity('valhalla').metadata.links[1].title, /FOSSGIS public demo/);
});
test('actual Home source expressions select explicit launchers and categories', () => {
  const hook = fs.readFileSync(path.join(root, 'plugins/plugin-onemind-common/src/useDirectory.ts'), 'utf8');
  const launch = new Function('entity', `return (${hook.match(/const raw =([\s\S]*?);/)[1]});`);
  const group = new Function('entity', 'known', `return (${hook.match(/const groupTitle =([\s\S]*?);/)[1]});`);
  const directory = references.map(([name, url]) => {
    const e = entity(name), category = name === 'adsb-exchange' ? 'Intelligence' : 'Geospatial';
    assert.equal(launch(e), `https://${url}`);
    assert.equal(group(e, { groupTitle: 'Wrong fallback' }), category);
    assert.equal(group(e, undefined), category);
    return { name, groupTitle: category };
  });
  const home = fs.readFileSync(path.join(root, 'plugins/plugin-onemind-home/src/homeModule.tsx'), 'utf8');
  const groups = new Function('directory', `return (${home.match(/const groups = (.*);/)[1]});`)(directory);
  assert.deepEqual(groups, ['Geospatial', 'Intelligence']);
});
test('reference caveats distinguish upstream projects, hosted products and public usage limits', () => {
  const desc = name => entity(name).metadata.description;
  assert.match(desc('opencage'), /Hosted geocoding.*not a sovereign self-hosted backend/);
  assert.match(desc('graphhopper'), /open-source routing engine.*separate hosted Directions API/);
  assert.match(desc('maputnik'), /MapLibre Style Specification, not a tile generator or proof of ATAK style compatibility/);
  assert.match(desc('suncalc'), /not a complete terrain shade assessment/);
  assert.match(desc('komoot'), /Hosted outdoor route planning/);
  assert.match(desc('adsb-exchange'), /future CoT bridge requires API, licensing and redistribution review/);
  assert.match(desc('adsb-exchange'), /no complete coverage or flight-safety suitability claimed/);
  for (const name of ['overpass-api', 'nominatim-upstream', 'valhalla']) assert.match(desc(name), /not an approved automated operational dependency/);
  const archive = cluster.find(e => e.metadata.name === 'fabric-nominatim-817008f1');
  assert.equal(archive.kind, 'Resource');
  assert.equal(archive.metadata.annotations['onemind.io/documented-address'], 'https://nominatim.tail717ba0.ts.net');
  assert.match(desc('nominatim-upstream'), /Distinct from the archived OneMind Nominatim registry record/);
  assert.notEqual(entity('nominatim-upstream').metadata.title, archive.metadata.title);
});
test('new reference blocks survive durable structure regeneration byte-for-byte', () => {
  const { applyStructureToText, loadSources } = require('./curate-structure.cjs');
  const out = applyStructureToText({ apps: text, cluster: clusterText }, loadSources());
  for (const [name] of references) {
    const block = source => source.split('---\n').find(b => b.includes(`  name: ${name}\n`));
    assert.ok(block(text), name);
    assert.equal(block(out.apps), block(text), name);
  }
});
const mirror = path.join(root, '../../platform/backstage/catalog/apps.yaml');
test('mapping source and mounted app catalogs stay byte-identical', { skip: !fs.existsSync(mirror) }, () => {
  assert.equal(fs.readFileSync(mirror, 'utf8'), text);
});
