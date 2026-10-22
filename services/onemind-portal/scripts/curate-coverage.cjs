// Generates the "coverage contract" annotations (flagship F5) for every Kubernetes app from the
// reviewed workload map: Argo CD, Kubernetes resources, Grafana dashboard and alerts, logs, runbook
// and source. Nothing here is typed per app. A new app or workload in catalog/workload-map.yaml gets
// the full set on the next `node scripts/curate-structure.cjs`, and scripts/audit-catalog.cjs fails
// when an app still lacks one without a reviewed entry in catalog/coverage-exceptions.yaml.
//
// Pure: no I/O, no clock. Called from curate-structure.cjs (applyStructure) after the Argo links and
// Kubernetes selectors are set, so it only ever fills what is missing and never replaces a reviewed value.
// Kubernetes label selector grammar: "!key" matches objects without that label, i.e. everything in the namespace.
const NAMESPACE_SCOPE_SELECTOR = '!onemind.io/scope-all';
const NOT_LIVE = new Set(['missing-review-required', 'retirement-review-required']);
// Preferred label keys when workloads share a key but not a value (set-based selector).
const SELECTOR_KEYS = ['app.kubernetes.io/name', 'app', 'k8s-app', 'name', 'app.kubernetes.io/instance', 'app.kubernetes.io/component'];

const annotationsOf = e => e.metadata.annotations || (e.metadata.annotations = {});
const live = workloads => workloads.filter(w => !NOT_LIVE.has(w.metadata.annotations?.['onemind.io/inventory-state']));

/** Equality selector "a=b,c=d" to a Map; null when it uses set-based or negated syntax. */
function parseSelector(text) {
  if (!text) return null;
  const out = new Map();
  for (const part of text.split(',')) {
    const m = part.trim().match(/^([A-Za-z0-9_./-]+)=([A-Za-z0-9_.-]+)$/);
    if (!m) return null;
    out.set(m[1], m[2]);
  }
  return out;
}

/**
 * One label selector for all of an app's live workloads: the equality pairs they all share, else
 * `key in (v1,v2)` over the best key every workload carries, else '' (no honest common selector).
 */
function deriveSelector(workloads) {
  const maps = workloads.map(w => parseSelector(w.metadata.annotations?.['backstage.io/kubernetes-label-selector']));
  if (!maps.length || maps.some(m => !m)) return '';
  const shared = [...maps[0]].filter(([k, v]) => maps.every(m => m.get(k) === v));
  if (shared.length) return shared.map(([k, v]) => `${k}=${v}`).sort().join(',');
  const keys = [...maps[0].keys()].filter(k => maps.every(m => m.has(k)));
  const key = SELECTOR_KEYS.find(k => keys.includes(k)) || keys.sort()[0];
  if (!key) return '';
  const values = [...new Set(maps.map(m => m.get(key)))].sort();
  return values.length === 1 ? `${key}=${values[0]}` : `${key} in (${values.join(',')})`;
}

function namespacesOf(workloads) {
  const counts = new Map();
  for (const w of workloads) {
    const ns = w.metadata.annotations?.['backstage.io/kubernetes-namespace'];
    if (ns) counts.set(ns, (counts.get(ns) || 0) + 1);
  }
  // Most workloads first, then name: the first entry is the primary namespace.
  return [...counts].sort((a, b) => b[1] - a[1] || a[0].localeCompare(b[0])).map(([ns]) => ns);
}

const workloadName = w => w.metadata.annotations?.['onemind.io/kubernetes-workload'];

/** Dashboard-relative path with the app's own variables; the portal adds the Grafana domain. */
function grafanaPath(uid, namespaces, sharedWorkloads) {
  const query = [
    ...namespaces.map(n => `var-namespace=${encodeURIComponent(n)}`),
    ...sharedWorkloads.map(n => `var-workload=${encodeURIComponent(n)}`),
  ];
  return `d/${uid}/${uid}${query.length ? `?${query.join('&')}` : ''}`;
}

/** LogQL for the app. Loki receives logs through OTLP, whose default index labels are k8s_*. */
function logsQuery(coverage, namespaces, sharedWorkloads) {
  const ns = coverage.loki?.namespaceLabel || 'k8s_namespace_name';
  const pod = coverage.loki?.podLabel || 'k8s_pod_name';
  const matchers = [namespaces.length === 1 ? `${ns}="${namespaces[0]}"` : `${ns}=~"${namespaces.join('|')}"`];
  if (sharedWorkloads.length) matchers.push(`${pod}=~"(${sharedWorkloads.join('|')})-.*"`);
  return `{${matchers.join(', ')}}`;
}

/** Selector for the Grafana plugin: bespoke dashboards of the app first, then the generic one. */
function dashboardSelector(coverage, appName) {
  const uids = [...(coverage.bespokeDashboards?.[appName] || []), coverage.genericDashboard.uid];
  return uids.map(uid => `uid == '${uid}'`).join(' || ');
}

/**
 * @param {{appComponent: Map<string, object>, subcomponents: Map<string, object[]>,
 *   argoByApp: Map<string, string[]>, workloadNamespaces: Map<string, Set<string>>}} model
 * @param {object} coverage parsed catalog/coverage.yaml
 */
function applyCoverage(model, coverage) {
  if (!coverage) return;
  const { appComponent, subcomponents, argoByApp, workloadNamespaces } = model;
  for (const [name, app] of appComponent) {
    if (!argoByApp.has(name) && !subcomponents.has(name)) continue;
    const ann = annotationsOf(app);
    const members = live(subcomponents.get(name) || []);
    const namespaces = namespacesOf(members);
    // A workload name narrows the dashboard and logs only when the namespace is shared with other apps.
    const shared = namespaces.some(ns => (workloadNamespaces.get(ns) || new Set()).size > 1);
    const sharedWorkloads = shared ? [...new Set(members.map(workloadName).filter(Boolean))].sort() : [];

    if (namespaces.length) {
      ann['onemind.io/kubernetes-namespaces'] = namespaces.join(',');
      if (namespaces.length === 1 && !ann['backstage.io/kubernetes-namespace']) ann['backstage.io/kubernetes-namespace'] = namespaces[0];
      const current = ann['backstage.io/kubernetes-label-selector'];
      // A reviewed selector or kubernetes-id is never replaced; the generated namespace fallback is re-derived.
      if (!ann['backstage.io/kubernetes-id'] && (!current || current === NAMESPACE_SCOPE_SELECTOR)) {
        const selector = deriveSelector(members);
        if (selector) {
          ann['backstage.io/kubernetes-label-selector'] = selector;
          delete ann['onemind.io/kubernetes-scope'];
        }
        // No label is common to the workloads: fall back to the whole namespace through a selector that
        // every object satisfies (a label that does not exist), and say so. Multi-namespace apps cannot.
        else if (namespaces.length === 1) {
          ann['backstage.io/kubernetes-label-selector'] = NAMESPACE_SCOPE_SELECTOR;
          ann['onemind.io/kubernetes-scope'] = 'namespace';
        }
      }
      if (!ann['grafana/alert-label-selector']) ann['grafana/alert-label-selector'] = `namespace=${namespaces[0]}`;
      ann['onemind.io/grafana-dashboard-path'] = grafanaPath(coverage.genericDashboard.uid, namespaces, sharedWorkloads);
      ann['onemind.io/logs-query'] = logsQuery(coverage, namespaces, sharedWorkloads);
    } else {
      // Namespace unknown (no Deployment or StatefulSet observed): never guess one.
      for (const key of ['onemind.io/kubernetes-namespaces', 'onemind.io/grafana-dashboard-path', 'onemind.io/logs-query']) delete ann[key];
    }
    if (!ann['grafana/dashboard-selector']) ann['grafana/dashboard-selector'] = dashboardSelector(coverage, name);
    ann['onemind.io/recovery-runbook'] = `${coverage.runbooks.base}${name}/`;
    if (!ann['github.com/project-slug']) ann['github.com/project-slug'] = coverage.repo;
  }
}

const RUNBOOK_DISCLAIMER = 'Generated from the reviewed catalog (workload-map.yaml, structure.yaml). It lists what is recorded, not what is currently healthy; check the Live tab for current state.';

function kindOf(w) {
  return w.metadata.annotations?.['onemind.io/kubernetes-workload-kind'] || 'Deployment';
}

/** Markdown for one app runbook. Every statement comes from the catalog or a fixed procedure. */
function renderRunbook(app, members, coverage, argoApps) {
  const a = app.metadata.annotations || {};
  const lines = [];
  lines.push(`# ${app.metadata.title || app.metadata.name}`, '', `> ${RUNBOOK_DISCLAIMER}`, '');
  if (app.metadata.description) lines.push(app.metadata.description, '');
  lines.push('## Where it runs', '');
  lines.push(`- Catalog entity: \`component:default/${app.metadata.name}\` in System \`${app.spec.system}\``);
  if (argoApps.length) lines.push(`- Argo CD application${argoApps.length > 1 ? 's' : ''}: ${argoApps.map(n => `\`${n}\``).join(', ')}`);
  else lines.push('- Argo CD application: none recorded for this app');
  const namespaces = (a['onemind.io/kubernetes-namespaces'] || '').split(',').filter(Boolean);
  lines.push(`- Namespace${namespaces.length === 1 ? '' : 's'}: ${namespaces.length ? namespaces.map(n => `\`${n}\``).join(', ') : 'not recorded (no Deployment or StatefulSet observed)'}`);
  lines.push('');
  const liveMembers = live(members);
  if (liveMembers.length) {
    lines.push('## Workloads', '', '| Workload | Kind | Namespace |', '| --- | --- | --- |');
    for (const w of liveMembers.slice().sort((x, y) => workloadName(x).localeCompare(workloadName(y))).slice(0, 40)) {
      lines.push(`| \`${workloadName(w)}\` | ${kindOf(w)} | \`${w.metadata.annotations['backstage.io/kubernetes-namespace']}\` |`);
    }
    if (liveMembers.length > 40) lines.push('', `${liveMembers.length - 40} more workloads are listed in the catalog.`);
    lines.push('');
  }
  lines.push('## When it looks wrong', '',
    '1. Open the app page. The status dot combines Argo CD, workload readiness, alerts and uptime; grey means unknown, not healthy.',
    '2. **Live** tab: confirm pods and replicas. **Logs** tab: open the generated Loki query. **Metrics** tab: open the workload dashboard with this app\'s variables.',
    '3. Find the cause before acting. A restart that hides the cause is not a fix.', '',
    '## Safe recovery', '',
    '- Change the cluster only through Git and Argo CD (pull request, merge, sync). Do not `kubectl apply`, patch or delete by hand, and do not run Helm directly.',
    '- A Git change is reversible with `git revert`; Argo CD then restores the previous state.',
    '- Data stores: confirm the latest backup in the [backups runbook](../backups.md) before any change that touches storage.', '',
    '## Generated links', '');
  if (a['onemind.io/logs-query']) lines.push(`- Loki query: \`${a['onemind.io/logs-query']}\``);
  if (a['onemind.io/grafana-dashboard-path']) lines.push(`- Grafana dashboard path: \`${a['onemind.io/grafana-dashboard-path']}\``);
  lines.push(`- Source: \`${a['github.com/project-slug'] || coverage.repo}\``, '',
    '## Escalate', '', 'If this app is unreachable and the cause is not clear after the Live and Logs tabs, treat it as an incident: record what was observed (time, tab, message) and follow the [resilience runbook](../resilience.md).', '');
  return lines.join('\n');
}

function renderRunbookIndex(entries) {
  const lines = ['# App runbooks', '', `> ${RUNBOOK_DISCLAIMER}`, '',
    'One page per Kubernetes app, generated from the catalog so a new app gets a runbook automatically.', '', '| App | System |', '| --- | --- |'];
  for (const { name, title, system } of entries) lines.push(`| [${title}](${name}.md) | ${system} |`);
  lines.push('');
  return lines.join('\n');
}

/** Map of relative path (under runbooks/docs) to Markdown for every covered app. */
function renderRunbooks(model, coverage) {
  const { appComponent, subcomponents, argoByApp } = model;
  const files = new Map();
  const entries = [];
  for (const [name, app] of [...appComponent].sort(([x], [y]) => x.localeCompare(y))) {
    if (!app.metadata.annotations?.['onemind.io/recovery-runbook']) continue;
    files.set(`${coverage.runbooks.base}${name}.md`, renderRunbook(app, subcomponents.get(name) || [], coverage, argoByApp.get(name) || []));
    entries.push({ name, title: app.metadata.title || name, system: app.spec.system });
  }
  files.set(`${coverage.runbooks.base}index.md`, renderRunbookIndex(entries));
  return files;
}

module.exports = { NAMESPACE_SCOPE_SELECTOR, applyCoverage, deriveSelector, parseSelector, namespacesOf, grafanaPath, logsQuery, dashboardSelector, renderRunbooks, renderRunbook, NOT_LIVE };
