// Source-level Home coverage; not a rendered UI, installation or deployment test.
const { test } = require('node:test');
const assert = require('node:assert/strict');
const fs = require('node:fs');
const path = require('node:path');
const YAML = require('yaml');
const root = path.join(__dirname, '..');
const text = fs.readFileSync(path.join(root, 'catalog/apps.yaml'), 'utf8');
const docs = YAML.parseAllDocuments(text).map(d => { assert.deepEqual(d.errors, []); return d.toJS(); });
const structure = YAML.parse(fs.readFileSync(path.join(root, 'catalog/structure.yaml'), 'utf8'));
const curated = JSON.parse(fs.readFileSync(path.join(root, 'plugins/plugin-onemind-common/directory.json'), 'utf8'));
const references = [
  [
    "hermes-atlas",
    "Hermes Atlas",
    "https://hermesatlas.com",
    "Community-curated Hermes ecosystem directory with guides, tools, skills and plugins.",
    "community"
  ],
  [
    "hermes-bible",
    "Hermes Bible",
    "https://www.hermesbible.com",
    "Unofficial community documentation and workflows for Hermes Agent; not affiliated with Nous Research.",
    "community"
  ],
  [
    "oh-my-hermes",
    "Oh My Hermes",
    "https://omher.app",
    "Community coding harness for Hermes Agent with IDE workflows.",
    "community"
  ],
  [
    "hermes-official-plugins",
    "Hermes Official Plugins",
    "https://github.com/NousResearch/hermes-official-plugins/tree/main",
    "Official Hermes Agent plugins built and maintained by Nous Research.",
    "official"
  ],
  [
    "hermes-skills",
    "Hermes Skills",
    "https://hermes-agent.nousresearch.com/docs/skills",
    "Official Nous Research Hermes skills documentation entry point, redirecting to the skills marketplace.",
    "official"
  ],
  [
    "hermes-documentation",
    "Hermes Documentation",
    "https://hermes-agent.nousresearch.com/docs/",
    "Official Nous Research documentation for Hermes Agent.",
    "official"
  ]
];
for (const [name, title, url, description, source] of references) {
  test(`${name}: exact unique external Home reference in Hermes Ecosystem`, () => {
    const matches = docs.filter(e => e.kind === 'Component' && e.metadata.name === name);
    assert.equal(matches.length, 1);
    const [e] = matches, a = e.metadata.annotations;
    assert.equal(e.metadata.title, title);
    assert.equal(e.metadata.description, `External reference: ${description} No OneMind installation or integration verified.`);
    assert.equal(e.spec.system, 'ai-agents');
    assert.equal(structure.placement[name], 'ai-agents');
    assert.equal(e.spec.type, 'website');
    assert.equal(e.spec.owner, 'group:default/onemind-admins');
    assert.equal(e.metadata.labels['onemind.io/pillar'], 'so');
    assert.equal(e.metadata.labels['onemind.io/catalog-role'], 'app');
    assert.ok(e.metadata.tags.includes(source));
    assert.ok(e.metadata.tags.includes('reference'));
    assert.ok(!e.metadata.tags.includes(source === 'official' ? 'community' : 'official'));
    assert.equal(a['onemind.io/hosting'], 'external-provider-unverified');
    assert.equal(a['onemind.io/inventory-review'], 'imported-link');
    assert.equal(a['onemind.io/launcher-url'], url);
    assert.equal(a['onemind.io/launcher-category'], 'Hermes Ecosystem');
    assert.deepEqual(e.metadata.links, [{ url, title: 'External reference', icon: 'dashboard' }]);
    assert.equal(docs.flatMap(x => x.metadata.links || []).filter(l => l.url === url).length, 1);
    assert.equal(docs.filter(x => x.metadata.annotations?.['onemind.io/launcher-url'] === url).length, 1);
    const destination = new URL(url);
    assert.equal(destination.protocol, 'https:');
    assert.equal(destination.username, '');
    assert.equal(destination.password, '');
    for (const key of Object.keys(a)) assert.ok(!/^(argocd\/|backstage.io\/kubernetes-|onemind.io\/(verified-at|deployment))/.test(key), key);
  });
}
test('actual directory hook expressions admit annotations and prefer the explicit category', () => {
  const hook = fs.readFileSync(path.join(root, 'plugins/plugin-onemind-common/src/useDirectory.ts'), 'utf8');
  // Execute the current hook expressions, not a second implementation of them.
  const raw = hook.match(/const raw =([\s\S]*?);/)[1];
  const category = hook.match(/const groupTitle =([\s\S]*?);/)[1];
  const launch = new Function('entity', `return (${raw});`);
  const group = new Function('entity', 'known', `return (${category});`);
  for (const [name, , url] of references) {
    const e = docs.find(x => x.metadata.name === name);
    assert.equal(launch(e), url);
    assert.equal(group(e, curated.find(x => x.name === name)), 'Hermes Ecosystem');
    assert.equal(group(e, { groupTitle: 'AI & Agents' }), 'Hermes Ecosystem');
  }
  assert.equal(launch({ metadata: { links: [{ title: 'External reference', url: 'https://example.test' }] } }), undefined);
  assert.equal(launch({ metadata: { links: [{ title: 'Open application', url: 'https://example.test' }] } }), 'https://example.test');
  assert.equal(group({ metadata: {} }, { groupTitle: 'AI & Agents' }), 'AI & Agents');
  assert.equal(group({ metadata: {} }, undefined), 'Services');
  const home = fs.readFileSync(path.join(root, 'plugins/plugin-onemind-home/src/homeModule.tsx'), 'utf8');
  const groupsExpression = home.match(/const groups = (.*);/)[1];
  const directory = references.map(([name]) => ({ groupTitle: group(docs.find(e => e.metadata.name === name), undefined) }));
  assert.deepEqual(new Function('directory', `return (${groupsExpression});`)(directory), ['Hermes Ecosystem']);
});
test('Legacy stays canonical in the broader AI system, without a new Hermes System', () => {
  assert.equal(docs.some(e => e.kind === 'System' && e.metadata.name === 'hermes-ecosystem'), false);
  assert.equal(docs.filter(e => e.kind === 'System' && e.metadata.name === 'ai-agents').length, 1);
  const legacy = docs.filter(e => e.kind === 'Component' && e.metadata.name === 'onemind-hermes');
  assert.equal(legacy.length, 1);
  assert.equal(legacy[0].spec.system, 'ai-agents');
});
test('Hermes blocks survive durable structure regeneration byte-for-byte', () => {
  const { applyStructureToText, loadSources } = require('./curate-structure.cjs');
  const out = applyStructureToText({ apps: text, cluster: fs.readFileSync(path.join(root, 'catalog/cluster.yaml'), 'utf8') }, loadSources());
  for (const [name] of references) {
    const block = s => s.split('\n---\n').find(b => b.includes(`\n  name: ${name}\n`));
    assert.ok(block(text), name);
    assert.equal(block(out.apps), block(text), name);
  }
});
const mirror = path.join(root, '../../platform/backstage/catalog/apps.yaml');
test('Hermes catalog mirrors are byte-identical', { skip: !fs.existsSync(mirror) }, () => {
  assert.equal(fs.readFileSync(mirror, 'utf8'), text);
});
