// The app annotation contract (flagship F7), as code. docs/catalog/app-contract.md is the prose
// version; test-app-contract.cjs fails when the catalog breaks either. Pure: no I/O, no clock.
// The hosting values live in plugin-onemind-common/src/hostingValues.json so the portal UI and this
// validator read one list.
const path = require('node:path');
const hostingValues = require(path.join(__dirname, '../plugins/plugin-onemind-common/src/hostingValues.json'));

const CANONICAL = Object.keys(hostingValues.canonical);
const LEGACY = Object.keys(hostingValues.legacy);
const HOSTING_VALUES = [...CANONICAL, ...LEGACY];

const HOSTING = 'onemind.io/hosting';
const HEALTH_CHECK = 'onemind.io/health-check';
const HEALTH_URL = 'onemind.io/health-url';
// Gatus endpoint names: lower case, digits and hyphens. Also the Prometheus `name` label.
const ENDPOINT_NAME = /^[a-z0-9][a-z0-9-]{0,62}$/;
const ISO_DATE = /^\d{4}-\d{2}-\d{2}$/;
// Facts are text only. These patterns catch a pasted secret, not every possible one; gitleaks is the backstop.
const SECRET_LIKE = [
  /\b(password|passwd|passphrase|secret|api[_-]?key|token)\s*[:=]\s*\S+/i,
  /\b(sk|pk|ghp|gho|xox[abp])[-_][A-Za-z0-9_-]{16,}/,
  /-----BEGIN [A-Z ]*PRIVATE KEY-----/,
  /\b[A-Za-z0-9+/]{40,}={0,2}\b/,
];

const annotationsOf = e => e.metadata?.annotations ?? {};

/** Effective hosting kind; same rules as hostingKind() in plugin-onemind-common/src/hosting.ts. */
function hostingKind(entity) {
  const a = annotationsOf(entity);
  const value = a[HOSTING];
  if (a['onemind.io/documented-status'] === 'Planned') return 'planned';
  const onCluster = !!(a['argocd/app-name'] || a['onemind.io/kubernetes-workloads'] || a['backstage.io/kubernetes-label-selector'] || a['backstage.io/kubernetes-id']);
  switch (value) {
    case 'digitalocean-nyc3': return onCluster ? 'kubernetes' : 'digitalocean';
    case 'external-provider-unverified': return 'saas';
    case 'mesh-endpoint-unverified': return 'home';
    case 'native-application': return 'desktop';
    case 'vehicle-firmware': return 'device';
    case 'external-website': return 'web';
    case 'planned-not-deployed': return 'planned';
    default: return CANONICAL.includes(value) ? value : 'unknown';
  }
}

/** Problems with one entity's contract annotations, as plain strings. Empty means valid. */
function validateEntity(entity, { requireHosting = true } = {}) {
  const a = annotationsOf(entity);
  const problems = [];
  const value = a[HOSTING];
  if (value === undefined || value === '') {
    if (requireHosting) problems.push(`${HOSTING} is missing (use "unknown" when it is not established)`);
  } else if (!HOSTING_VALUES.includes(value)) {
    problems.push(`${HOSTING}: "${value}" is not one of ${HOSTING_VALUES.join(', ')}`);
  }
  const check = a[HEALTH_CHECK];
  if (check !== undefined) {
    if (!ENDPOINT_NAME.test(check)) problems.push(`${HEALTH_CHECK}: "${check}" must be a Gatus endpoint name (lower case letters, digits, hyphens)`);
    if (!a[HEALTH_URL]) problems.push(`${HEALTH_CHECK} is set but ${HEALTH_URL} is missing, so Gatus has nothing to probe`);
  }
  if (a[HEALTH_URL] !== undefined) {
    if (!/^https?:\/\/[^\s/@]+(:\d+)?(\/\S*)?$/.test(a[HEALTH_URL])) problems.push(`${HEALTH_URL}: "${a[HEALTH_URL]}" must be an http(s) URL with no credentials`);
    if (check === undefined) problems.push(`${HEALTH_URL} is set but ${HEALTH_CHECK} (the endpoint name) is missing`);
  }
  if (a['onemind.io/renewal'] !== undefined && !(ISO_DATE.test(a['onemind.io/renewal']) || /^(monthly|yearly|none)$/.test(a['onemind.io/renewal'])))
    problems.push('onemind.io/renewal must be YYYY-MM-DD, monthly, yearly or none');
  for (const key of ['onemind.io/plan', 'onemind.io/access', 'onemind.io/data-export', 'onemind.io/backup-location', 'onemind.io/recovery-runbook', 'onemind.io/renewal']) {
    if (a[key] === undefined) continue;
    if (!String(a[key]).trim()) problems.push(`${key} is empty; remove it so the page shows "Not recorded"`);
    if (SECRET_LIKE.some(p => p.test(String(a[key])))) problems.push(`${key} looks like it contains a credential; record text only`);
  }
  return problems;
}

module.exports = { CANONICAL, LEGACY, HOSTING_VALUES, ENDPOINT_NAME, hostingKind, validateEntity, annotationsOf };
