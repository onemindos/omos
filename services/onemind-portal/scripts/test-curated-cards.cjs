const { test } = require('node:test');
const assert = require('node:assert/strict');
const fs = require('fs');
const path = require('path');
const { curatedBlocks, restoreCuratedCards, dockerfileTestArgs } = require('./curated-cards.cjs');

const card = (name, tags, extra = '') => `apiVersion: backstage.io/v1alpha1
kind: Component
metadata:
  name: ${name}
  description: >-
    Curated text for ${name}.
  tags:
${tags.map(t => `    - ${t}`).join('\n')}
  annotations:
    onemind.io/hosting: external-provider-unverified
${extra}spec:
  type: website
  owner: group:default/onemind-admins
`;

test('reference and planned cards are curated; others are not', () => {
  const text = [card('ref', ['reference']), card('plan', ['planned']), card('live', ['documented'])].join('---\n');
  assert.deepEqual([...curatedBlocks(text).keys()].sort(), ['plan', 'ref']);
});

test('a refresh that rewrites curated cards is put back byte-for-byte; other cards keep the refresh', () => {
  const original = [card('ref', ['reference']), card('live', ['documented'])].join('---\n');
  const refreshed = [
    card('ref', ['reference', 'geo'], '    onemind.io/documented-status: Live\n').replace('external-provider-unverified', 'digitalocean-nyc3'),
    card('live', ['documented'], '    onemind.io/verified-at: now\n'),
  ].join('---\n');
  const restored = restoreCuratedCards(original, refreshed);
  const blocks = restored.split(/^---\n/m);
  assert.equal(blocks[0], card('ref', ['reference']));
  assert.match(blocks[1], /verified-at: now/);
});

test('a curated card the refresh dropped comes back', () => {
  const original = [card('ref', ['reference']), card('live', ['documented'])].join('---\n');
  const restored = restoreCuratedCards(original, card('live', ['documented']));
  assert.deepEqual([...curatedBlocks(restored).keys()], ['ref']);
});

test('restoring an untouched file is a no-op', () => {
  const text = [card('ref', ['reference']), card('plan', ['planned'])].join('---\n');
  assert.equal(restoreCuratedCards(text, text), text);
});

test('refresh runs the exact catalog tests the Dockerfile runs', () => {
  const args = dockerfileTestArgs(path.join(__dirname, '..', 'packages/backend/Dockerfile'));
  assert.equal(args[0], '--test');
  assert.ok(args.includes('scripts/test-audit-catalog.cjs') && args.includes('scripts/test-topology.cjs'));
  assert.match(fs.readFileSync(path.join(__dirname, 'refresh-inventory.cjs'), 'utf8'), /runDockerfileCatalogTests\(\)/);
});
