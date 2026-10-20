const { test } = require('node:test');
const assert = require('node:assert/strict');
const { auditCatalog } = require('./audit-catalog.cjs');
const now = new Date('2026-10-20T06:00:00Z');
const entity = (name, extra = {}) => ({
  kind: 'Component',
  metadata: { name },
  ...extra,
});

test('canonical identities, duplicate refs and missing targets are explicit', () => {
  const result = auditCatalog(
    [
      entity('Example', {
        relations: [{ targetRef: 'RESOURCE:DEFAULT/Missing' }],
      }),
      entity('example'),
    ],
    now,
  );
  assert.deepEqual(result.duplicateRefs, ['component:default/example']);
  assert.deepEqual(result.missingRelations, [
    {
      target: 'resource:default/missing',
      sources: ['component:default/example'],
    },
  ]);
});
test('existing relations resolve case-insensitively and default namespaces match', () => {
  const result = auditCatalog(
    [
      entity('a', { relations: [{ targetRef: 'COMPONENT:DEFAULT/B' }] }),
      entity('b'),
    ],
    now,
  );
  assert.deepEqual(result.missingRelations, []);
});
test('freshness is historical metadata, not health; invalid/future dates are not fresh', () => {
  const entries = [
    '2026-10-06T00:00:00Z',
    '2026-10-20T05:00:00Z',
    'bad',
    '2026-10-21T00:00:00Z',
  ].map((date, i) =>
    entity(`a${i}`, {
      metadata: {
        name: `a${i}`,
        annotations: { 'onemind.io/inventory-observed-at': date },
      },
    }),
  );
  const result = auditCatalog(entries, now);
  assert.equal(result.staleInventory.length, 1);
  assert.deepEqual(result.invalidInventoryDates, [
    'component:default/a2',
    'component:default/a3',
  ]);
  assert.match(result.scope, /not service health/);
});
test('shared launchers are candidates; URLs are never emitted and inputs stay unchanged', () => {
  const entries = ['a', 'b'].map(name =>
    entity(name, {
      metadata: {
        name,
        links: [
          {
            title: 'Open application',
            url: 'https://example.test/private?signed=never-output#fragment',
          },
        ],
      },
    }),
  );
  const before = JSON.stringify(entries);
  const result = auditCatalog(entries, now);
  assert.deepEqual(result.sharedLauncherCandidates, [
    ['component:default/a', 'component:default/b'],
  ]);
  assert.ok(!JSON.stringify(result).includes('never-output'));
  assert.equal(JSON.stringify(entries), before);
});
test('orphan and reported errors stay distinct from relation integrity', () => {
  const result = auditCatalog(
    [
      entity('a', {
        metadata: { name: 'a', annotations: { 'backstage.io/orphan': 'true' } },
        status: { items: [{ level: 'error' }] },
      }),
    ],
    now,
  );
  assert.equal(result.orphans.length, 1);
  assert.equal(result.reportedErrors.length, 1);
});
test('rejects partial query envelopes, malformed entities and invalid clocks', () => {
  assert.throws(() => auditCatalog({ items: [] }, now));
  assert.throws(() => auditCatalog([{}], now));
  assert.throws(() => auditCatalog([], new Date('bad')));
});
