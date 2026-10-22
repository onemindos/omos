// Offline metadata report. Supply the complete, authorized /entities JSON array.
// No requests, source edits, URL output or workload-health inference.
const fs = require('node:fs');
const REVIEW = 'onemind.io/launcher-review';
const REVIEWED_AT = 'onemind.io/launcher-reviewed-at';
const TAG = 'launcher-review-required';
const EVIDENCE = 'docs/operations/portal-acceptance-followthrough-2026-10-20.md#launcher-review-and-repairs';
// OpenClaw was retired on 2026-10-22, which resolves its historical finding.
const HISTORICAL_NAMES = Object.freeze([
  'cloudtak', 'tak-server', 'onetak', 'fedhub', 'clickhouse', 'minio',
  'postgis', 'longhorn', 'prometheus', 'nats-surveyor', 'tileserver-gl',
  'valhalla', 'maputnik', 'webodm', 'home-assistant',
]);
const CLASSIFICATIONS = Object.freeze({
  'dns-unresolved-from-operator-mac': 'dns-unresolved',
  'certificate-untrusted-from-operator-mac': 'certificate-untrusted',
  'timeout-from-operator-mac': 'timeout',
  'reachable-from-operator-mac': 'reachable-observation',
});
const record = value => value !== null && typeof value === 'object' && !Array.isArray(value);
const part = value => typeof value === 'string' && /^[a-z0-9][a-z0-9_.-]{0,62}$/i.test(value);
// Deliberately accept only timestamp-shaped strings for output; arbitrary malformed
// annotation text might contain a URL, credential or personal filesystem path.
const dateShape = /^\d{4}-\d{2}-\d{2}T\d{2}:\d{2}:\d{2}(?:\.\d{1,9})?(?:Z|[+-]\d{2}:\d{2})$/;
function observation(value, now, staleDays) {
  if (value === undefined) return { observedAt: null, freshness: 'missing' };
  if (!dateShape.test(value)) return { observedAt: null, freshness: 'invalid', sourceValueRedacted: true };
  const stamp = Date.parse(value);
  // Date.parse normalizes impossible dates such as February 30: reject those too.
  const [year, month, day] = value.slice(0, 10).split('-').map(Number);
  const timeValid = Number(value.slice(11, 13)) < 24 && Number(value.slice(14, 16)) < 60 &&
    Number(value.slice(17, 19)) < 60 && (value.endsWith('Z') ||
      (Number(value.slice(-5, -3)) < 24 && Number(value.slice(-2)) < 60));
  const calendarValid = timeValid && month >= 1 && month <= 12 && day >= 1 &&
    day <= new Date(Date.UTC(year, month, 0)).getUTCDate();
  const freshness = !Number.isFinite(stamp) || !calendarValid ? 'invalid' :
    stamp > now.getTime() ? 'future' :
    now.getTime() - stamp > staleDays * 86400000 ? 'stale' : 'fresh';
  return { observedAt: value, freshness };
}
function reportLauncherReview(entities, now = new Date(), staleDays = 7) {
  if (!Array.isArray(entities) || !(now instanceof Date) ||
      !Number.isFinite(now.getTime()) || !Number.isFinite(staleDays) || staleDays <= 0) {
    throw new Error('Expected complete entity array, valid clock and positive finite stale-days');
  }
  const byRef = new Map();
  for (const entity of entities) {
    if (!record(entity) || !part(entity.kind) || !record(entity.metadata) ||
        !part(entity.metadata.name) || (entity.metadata.namespace !== undefined && !part(entity.metadata.namespace)) ||
        (entity.metadata.tags !== undefined && (!Array.isArray(entity.metadata.tags) ||
          !entity.metadata.tags.every(tag => typeof tag === 'string'))) ||
        (entity.metadata.annotations !== undefined && (!record(entity.metadata.annotations) ||
          !Object.values(entity.metadata.annotations).every(value => typeof value === 'string')))) {
      throw new Error('Invalid catalog entity metadata');
    }
    const ref = `${entity.kind}:${entity.metadata.namespace || 'default'}/${entity.metadata.name}`.toLowerCase();
    if (byRef.has(ref)) throw new Error('Duplicate canonical catalog identity');
    const annotations = entity.metadata.annotations || {};
    const required = (entity.metadata.tags || []).includes(TAG);
    const rawReview = annotations[REVIEW];
    byRef.set(ref, {
      ref,
      reviewRequired: required,
      classification: Object.hasOwn(CLASSIFICATIONS, rawReview) ? CLASSIFICATIONS[rawReview] : 'unknown',
      reviewStatus: required ? 'unresolved' : 'closure-review-required',
      observation: observation(annotations[REVIEWED_AT], now, staleDays),
      evidenceRefs: [
        ...(rawReview !== undefined ? [`${ref}#${REVIEW}`] : []),
        ...(annotations[REVIEWED_AT] !== undefined ? [`${ref}#${REVIEWED_AT}`] : []),
      ],
      hasReviewMetadata: rawReview !== undefined || annotations[REVIEWED_AT] !== undefined,
    });
  }
  const publicRow = row => {
    const { hasReviewMetadata, ...safe } = row;
    return safe;
  };
  const findings = [...byRef.values()].filter(row => row.reviewRequired)
    .sort((a, b) => a.ref.localeCompare(b.ref)).map(publicRow);
  const historicalFindings = HISTORICAL_NAMES.map(name => {
    const ref = `component:default/${name}`;
    const row = byRef.get(ref);
    return {
      ...(row ? publicRow(row) : { ref, reviewRequired: null, classification: 'unknown',
        reviewStatus: 'unresolved', observation: { observedAt: null, freshness: 'missing' }, evidenceRefs: [] }),
      closureStatus: !row ? 'unresolved-entity-absent' : row.reviewRequired ? 'unresolved' : 'closure-review-required',
      evidenceRefs: [...(row?.evidenceRefs || []), EVIDENCE],
    };
  });
  return {
    schemaVersion: 1,
    reportedAt: now.toISOString(),
    scope: 'Offline catalog metadata only; not workload health, endpoint reachability, authorization or deployment acceptance',
    completeness: 'Caller must supply a complete authorized /entities array; array completeness cannot be proven offline',
    entityCount: entities.length,
    staleAfterDays: staleDays,
    reviewRequiredCount: findings.length,
    findings,
    historicalFindings,
    otherReviewedLaunchers: [...byRef.values()]
      .filter(row => row.hasReviewMetadata && !row.reviewRequired && !HISTORICAL_NAMES.some(name => row.ref === `component:default/${name}`))
      .sort((a, b) => a.ref.localeCompare(b.ref)).map(publicRow),
    note: 'Tag removal, missing entities and reachable metadata do not certify repair. Preserve identities and source observation times; obtain reviewed closure evidence separately. No automatic changes.',
  };
}
if (require.main === module) {
  try {
    if (process.argv.length !== 3) throw new Error('usage');
    const file = process.argv[2];
    const result = reportLauncherReview(JSON.parse(fs.readFileSync(file === '-' ? 0 : file, 'utf8')));
    console.log(JSON.stringify(result, null, 2));
    // Unresolved review findings are a successful report, not workload failures.
  } catch {
    // Do not echo parser errors or filenames, which can expose input values/paths.
    console.error('Invalid input. Usage: node scripts/report-launcher-review.cjs <complete-entities.json|->');
    process.exitCode = 2;
  }
}
module.exports = { reportLauncherReview, HISTORICAL_NAMES };
