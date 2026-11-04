const { test } = require('node:test');
const assert = require('node:assert/strict');
const { spawnSync } = require('node:child_process');
const path = require('node:path');
const { reportLauncherReview, HISTORICAL_NAMES } = require('./report-launcher-review.cjs');
const now = new Date('2026-10-21T08:02:51Z');
const timestamp = '2026-10-20T08:02:51.107324+00:00';
const entity = (name = 'cloudtak', review = 'dns-unresolved-from-operator-mac', date = timestamp, tagged = true) => ({
  kind: 'Component', metadata: { name,
    tags: tagged ? ['launcher-review-required'] : [],
    annotations: { 'onemind.io/launcher-review': review, 'onemind.io/launcher-reviewed-at': date },
  },
});
test('canonical identities, current classifications and all 13 historical closures', () => {
  const input = HISTORICAL_NAMES.map(name => entity(name));
  input[0].kind = 'COMPONENT'; input[0].metadata.name = 'CLOUDTAK';
  input[1].metadata.annotations['onemind.io/launcher-review'] = 'certificate-untrusted-from-operator-mac';
  input[12].metadata.annotations['onemind.io/launcher-review'] = 'timeout-from-operator-mac';
  const before = JSON.stringify(input);
  const result = reportLauncherReview(input, now);
  assert.equal(result.reviewRequiredCount, 13);
  assert.equal(result.historicalFindings.length, 13);
  assert.equal(new Set(result.historicalFindings.map(row => row.ref)).size, 13);
  assert.ok(result.historicalFindings.every(row => row.closureStatus === 'unresolved'));
  assert.equal(result.historicalFindings[0].ref, 'component:default/cloudtak');
  assert.equal(result.historicalFindings[1].classification, 'certificate-untrusted');
  assert.equal(result.historicalFindings[12].classification, 'timeout');
  assert.equal(result.historicalFindings[0].observation.observedAt, timestamp);
  assert.ok(result.findings[0].evidenceRefs.some(ref => ref.endsWith('#onemind.io/launcher-review')));
  assert.equal(JSON.stringify(input), before);
  assert.deepEqual(reportLauncherReview([...input].reverse(), now), result);
});
test('tag removal, absent entities, unknown classifications never silently close findings', () => {
  const result = reportLauncherReview([
    entity('cloudtak', 'reachable-from-operator-mac', timestamp, false),
    entity('new-finding', 'new-unknown-observation'),
    entity('agent-monitor', 'reachable-from-operator-mac', timestamp, false),
  ], now);
  assert.equal(result.historicalFindings[0].closureStatus, 'closure-review-required');
  assert.equal(result.historicalFindings[1].closureStatus, 'unresolved-entity-absent');
  assert.equal(result.findings[0].classification, 'unknown');
  assert.equal(result.findings[0].reviewStatus, 'unresolved');
  assert.equal(result.otherReviewedLaunchers[0].ref, 'component:default/agent-monitor');
});
test('missing, invalid, impossible, future, stale and boundary dates stay distinct', () => {
  const dates = [undefined, 'bad', '2026-02-30T12:00:00Z', '2026-10-22T00:00:00Z',
    '2026-10-01T00:00:00Z', '2026-10-14T08:02:51Z', timestamp];
  const states = ['missing', 'invalid', 'invalid', 'future', 'stale', 'fresh', 'fresh'];
  dates.forEach((date, i) => {
    const input = entity();
    if (date === undefined) delete input.metadata.annotations['onemind.io/launcher-reviewed-at'];
    else input.metadata.annotations['onemind.io/launcher-reviewed-at'] = date;
    assert.equal(reportLauncherReview([input], now).findings[0].observation.freshness, states[i]);
  });
  const conflict = reportLauncherReview([entity('cloudtak', 'reachable-from-operator-mac')], now);
  assert.equal(conflict.findings[0].reviewStatus, 'unresolved');
});
test('reject malformed entities/envelopes, duplicate IDs and invalid report settings', () => {
  for (const input of [{ items: [] }, { entities: [], nextCursor: 'x' }, null, [null], [{}],
    [{ kind: 'Component', metadata: { name: '/private/path' } }],
    [entity(), entity('CLOUDTAK')],
    [{ ...entity(), metadata: { name: 'x', annotations: [] } }],
    [{ ...entity(), metadata: { name: 'x', tags: 'launcher-review-required' } }],
    [{ ...entity(), metadata: { name: 'x', annotations: { bad: {} } } }]]) {
    assert.throws(() => reportLauncherReview(input, now));
  }
  assert.throws(() => reportLauncherReview([], new Date('bad')));
  for (const days of [0, -1, Infinity, NaN, '7']) assert.throws(() => reportLauncherReview([], now, days));
});
test('report and errors do not expose arbitrary URLs, secrets or personal paths', () => {
  const sensitive = 'https://example.invalid/private?token=do-not-print#/Users/person';
  const input = entity('cloudtak', sensitive, sensitive);
  input.metadata.links = [{ url: sensitive }];
  input.metadata.description = sensitive;
  input.metadata.annotations['onemind.io/launcher-url'] = sensitive;
  const result = JSON.stringify(reportLauncherReview([input], now));
  assert.ok(!result.includes(sensitive));
  assert.ok(!result.includes('do-not-print'));
  assert.ok(!result.includes('/Users/'));
  assert.equal(JSON.parse(result).findings[0].observation.sourceValueRedacted, true);
});
const cli = path.join(__dirname, 'report-launcher-review.cjs');
test('CLI accepts complete-array stdin; unresolved findings exit 0 and bad input exits 2', () => {
  const ok = spawnSync(process.execPath, [cli, '-'], { input: JSON.stringify([entity()]), encoding: 'utf8' });
  assert.equal(ok.status, 0); assert.equal(JSON.parse(ok.stdout).reviewRequiredCount, 1);
  for (const input of ['{ "items": [] }', 'secret-invalid-json', '[{}]']) {
    const bad = spawnSync(process.execPath, [cli, '-'], { input, encoding: 'utf8' });
    assert.equal(bad.status, 2); assert.equal(bad.stdout, '');
    assert.ok(!bad.stderr.includes('secret-invalid-json'));
  }
  const missing = spawnSync(process.execPath, [cli, '/Users/private/missing.json'], { encoding: 'utf8' });
  assert.equal(missing.status, 2); assert.ok(!missing.stderr.includes('/Users/private'));
});

test('current paired source catalogs retain all 13 historical findings and four other observations', () => {
  const fs = require('node:fs');
  const root = path.resolve(__dirname, '../../..');
  const source = fs.readFileSync(path.join(root, 'services/onemind-portal/catalog/apps.yaml'), 'utf8');
  const paired = fs.readFileSync(path.join(root, 'platform/backstage/catalog/apps.yaml'), 'utf8');
  assert.equal(source, paired, 'paired current source must agree');
  // Fixture extraction only, intentionally limited to the current source layout.
  // The reporting CLI accepts JSON and never attempts to parse YAML.
  const reviewed = source.split(/^---\s*$/m).filter(doc => doc.includes('onemind.io/launcher-review:'));
  assert.equal(reviewed.length, 17);
  const entries = reviewed.map(doc => {
    const name = doc.match(/^  name: (\S+)$/m)?.[1];
    const review = doc.match(/^    onemind.io\/launcher-review: (\S+)$/m)?.[1];
    const date = doc.match(/^    onemind.io\/launcher-reviewed-at: "([^"]+)"$/m)?.[1];
    assert.ok(name && review && date, 'source layout changed: update reviewed fixture extraction');
    return entity(name, review, date, /^    - launcher-review-required$/m.test(doc));
  });
  const result = reportLauncherReview(entries, now);
  assert.equal(result.reviewRequiredCount, 13);
  assert.ok(result.historicalFindings.every(row => row.closureStatus === 'unresolved'));
  assert.deepEqual(result.historicalFindings.map(row => row.ref).sort(),
    result.findings.map(row => row.ref).sort());
  assert.equal(result.findings.filter(row => row.classification === 'dns-unresolved').length, 11);
  assert.equal(result.findings.filter(row => row.classification === 'certificate-untrusted').length, 1);
  assert.equal(result.findings.filter(row => row.classification === 'timeout').length, 1);
  assert.equal(result.otherReviewedLaunchers.length, 4);
  assert.ok(result.otherReviewedLaunchers.every(row => row.classification === 'reachable-observation'));
});

test('timestamp normalization cannot turn invalid clock fields into fresh observations', () => {
  for (const date of ['2026-10-20T24:00:00Z', '2026-10-20T08:60:00Z', '2026-10-20T08:00:60Z',
    '2026-10-20T08:00:00+00:60', '2026-10-20T08:00:00+24:00']) {
    assert.equal(reportLauncherReview([entity('x', undefined, date)], now).findings[0].observation.freshness, 'invalid');
  }
});
