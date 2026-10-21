const { test } = require('node:test');
const assert = require('node:assert/strict');
const { missingDescription } = require('./inventory-status.cjs');
test('makes missing-review state visible without asserting downtime/deletion', () => {
  const result = missingDescription('Deployment in onetak');
  assert.match(result, /^Not observed in the latest inventory; review required\./);
  assert.ok(result.endsWith('Deployment in onetak'));
  assert.ok(!/offline|deleted|healthy/i.test(result));
});
test('repeated missing observations do not grow duplicate warning prefixes', () => {
  const once = missingDescription('Previous snapshot');
  assert.equal(missingDescription(once), once);
});
test('an absent optional description is explicit', () => {
  assert.match(missingDescription(), /No previous description recorded/);
});
