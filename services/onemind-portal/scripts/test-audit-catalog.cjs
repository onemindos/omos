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

// Native visibility metadata is reapplied by the reviewed inventory refresh.
// Keep these in the existing catalog gate (local pipeline and Linux Dockerfile).
const fs = require('node:fs');
const path = require('node:path');
const YAML = require('yaml');
const {
  curateDependencies,
  reviewedWorkloads,
} = require('./curate-dependencies.cjs');
const root = path.resolve(__dirname, '..');
const load = file =>
  YAML.parseAllDocuments(fs.readFileSync(file, 'utf8')).map(d => d.toJSON());
const catalog = () => load(path.join(root, 'catalog/cluster.yaml'));

test('curation is pure, idempotent and restores all reviewed native mappings after refresh', () => {
  const entries = catalog();
  for (const e of entries) {
    for (const key of [
      'argocd/app-name',
      'argocd/instance-name',
      'backstage.io/source-location',
      'grafana/dashboard-selector',
    ]) {
      delete e.metadata.annotations?.[key];
    }
  }
  const before = JSON.stringify(entries);
  const curated = curateDependencies(entries);
  assert.equal(JSON.stringify(entries), before);
  assert.deepEqual(curateDependencies(curated), curated);
  assert.deepEqual(
    curated.map(e => [e.kind, e.metadata.name]),
    entries.map(e => [e.kind, e.metadata.name]),
  );
  for (const target of reviewedWorkloads) {
    const e = curated.find(
      x => x.kind === 'Component' && x.metadata.name === target.name,
    );
    assert.equal(e.metadata.annotations['argocd/app-name'], target.application);
    assert.equal(
      e.metadata.annotations['grafana/dashboard-selector'],
      `uid == '${target.grafanaUid}'`,
    );
    assert.ok(
      e.metadata.links.some(
        link =>
          link.type === 'grafana' &&
          link.url ===
            `https://grafana.tail717ba0.ts.net/d/${target.grafanaUid}`,
      ),
    );
    assert.equal(
      e.metadata.annotations['argocd/instance-name'],
      'onemind-cloud',
    );
    for (const suffix of ['kubernetes', 'deployment-summary']) {
      assert.ok(
        e.metadata.links.some(link =>
          link.url.endsWith(`/${target.name}/${suffix}`),
        ),
      );
    }
    assert.ok(
      e.metadata.links.some(
        link =>
          link.type === 'runbook' && link.url.endsWith(`/${target.runbook}/`),
      ),
    );
  }
});
test('native mappings preserve unrelated metadata, dependencies, owners and same-namespace services', () => {
  const entries = catalog();
  const e = entries.find(x => x.metadata.name === reviewedWorkloads[0].name);
  e.metadata.annotations['example.test/keep'] = 'reviewed';
  e.metadata.links.push({
    title: 'Other runbook',
    type: 'runbook',
    url: 'https://example.test/runbook',
  });
  e.spec.dependsOn.push('resource:default/another-reviewed-dependency');
  const result = curateDependencies(entries);
  for (let i = 0; i < entries.length; i++) {
    assert.deepEqual(result[i].spec, entries[i].spec);
    if (!reviewedWorkloads.some(w => w.name === entries[i].metadata.name)) {
      assert.deepEqual(result[i], entries[i]);
    }
  }
  const curated = result.find(x => x.metadata.name === e.metadata.name);
  assert.equal(curated.metadata.annotations['example.test/keep'], 'reviewed');
  assert.ok(
    curated.metadata.links.some(link => link.title === 'Other runbook'),
  );
});
test('curation rejects missing, duplicate, wrong namespace or changed selectors without partial input writes', () => {
  const mutations = [
    entries =>
      entries.filter(e => e.metadata.name !== reviewedWorkloads[0].name),
    entries => [
      ...entries,
      structuredClone(
        entries.find(e => e.metadata.name === reviewedWorkloads[0].name),
      ),
    ],
    entries => {
      entries.find(
        e => e.metadata.name === reviewedWorkloads[0].name,
      ).metadata.namespace = 'other';
      return entries;
    },
    entries => {
      entries.find(
        e => e.metadata.name === reviewedWorkloads[0].name,
      ).metadata.annotations['backstage.io/kubernetes-namespace'] = 'nats';
      return entries;
    },
    entries => {
      entries.find(
        e => e.metadata.name === reviewedWorkloads[0].name,
      ).metadata.annotations['backstage.io/kubernetes-label-selector'] =
        'app=other';
      return entries;
    },
  ];
  for (const mutate of mutations) {
    const input = mutate(catalog());
    const before = JSON.stringify(input);
    assert.throws(() => curateDependencies(input), /review ownership/);
    assert.equal(JSON.stringify(input), before);
  }
});
test('label selector order does not change reviewed workload identity', () => {
  const entries = catalog();
  const nats = entries.find(e => e.metadata.name === reviewedWorkloads[2].name);
  nats.metadata.annotations['backstage.io/kubernetes-label-selector'] =
    reviewedWorkloads[2].selector.split(',').reverse().join(',');
  assert.doesNotThrow(() => curateDependencies(entries));
});
test('portal dependency is re-created exactly once after inventory refresh', () => {
  const entries = catalog();
  const portal = entries.find(
    e => e.metadata.name === reviewedWorkloads[0].name,
  );
  const databaseRef = `component:default/${reviewedWorkloads[1].name}`;
  portal.spec.dependsOn = portal.spec.dependsOn.filter(
    ref => ref !== databaseRef,
  );
  const result = curateDependencies(curateDependencies(entries));
  assert.equal(
    result
      .find(e => e.metadata.name === portal.metadata.name)
      .spec.dependsOn.filter(ref => ref === databaseRef).length,
    1,
  );
});

test('installed predicate schema admits only explicit Grafana alert selector metadata', () => {
  const {
    parseFilterPredicate,
    filterPredicateToFilterFunction,
  } = require('@backstage/filter-predicates');
  const predicate = parseFilterPredicate({
    'metadata.annotations.grafana/alert-label-selector': { $exists: true },
  });
  const accepts = filterPredicateToFilterFunction(predicate);
  assert.equal(
    accepts({
      metadata: { annotations: { 'grafana/dashboard-selector': 'nats' } },
    }),
    false,
  );
  assert.equal(
    accepts({
      metadata: {
        annotations: { 'grafana/alert-label-selector': 'namespace=nats' },
      },
    }),
    true,
  );
  assert.equal(accepts({ metadata: {} }), false);
});


test('reviewed service Prometheus scopes survive inventory curation', () => {
  const { curateDependencies, reviewedWorkloads } = require('./curate-dependencies.cjs');
  const YAML = require('yaml');
  const fs = require('node:fs');
  const path = require('node:path');
  const docs = YAML.parseAllDocuments(fs.readFileSync(path.join(__dirname, '../catalog/cluster.yaml'), 'utf8')).map(d => d.toJSON());
  for (const item of docs) {
    delete item.metadata.annotations?.['prometheus.io/service-name'];
    delete item.metadata.annotations?.['prometheus.io/alert'];
  }
  const curated = curateDependencies(docs);
  for (const reviewed of reviewedWorkloads) {
    const entry = curated.find(e => e.kind === 'Component' && e.metadata.name === reviewed.name);
    assert.equal(entry.metadata.annotations['prometheus.io/service-name'], reviewed.alertScope);
    assert.equal(entry.metadata.annotations['prometheus.io/alert'], 'all');
    assert.equal(entry.metadata.annotations['prometheus.io/labels'], undefined);
  }
});

// --- F5 coverage contract guard -------------------------------------------------------------
const { auditCoverage, kubernetesApps, COVERAGE } = require('./audit-catalog.cjs');
const fullAnnotations = () => ({
  'argocd/app-name': 'platform-x',
  'backstage.io/kubernetes-label-selector': 'app=x',
  'grafana/dashboard-selector': "uid == 'onemind-workload'",
  'grafana/alert-label-selector': 'namespace=x',
  'onemind.io/logs-query': '{k8s_namespace_name="x"}',
  'onemind.io/recovery-runbook': 'apps/x/',
  'github.com/project-slug': 'onemindos/omos',
});
const app = (name, annotations = fullAnnotations()) => ({
  kind: 'Component',
  metadata: { name, labels: { 'onemind.io/catalog-role': 'app' }, annotations },
  spec: {},
});
// A workload that belongs to the app keeps it in scope even when its Argo annotation is the missing one.
const owned = name => ({ kind: 'Component', metadata: { name: `${name}-workload` }, spec: { subcomponentOf: `component:default/${name}` } });
const without = (...keys) => {
  const a = fullAnnotations();
  for (const key of keys) delete a[key];
  return a;
};
const reason = 'Deployed by a shared Argo application that the reviewed map attributes to another app.';

test('coverage: a fully annotated Kubernetes app passes', () => {
  const result = auditCoverage([app('x')]);
  assert.equal(result.kubernetesApps, 1);
  assert.equal(result.fullyCovered, 1);
  assert.deepEqual(result.missing, []);
});
test('coverage: a missing annotation is reported with its capability, never silently skipped', () => {
  for (const capability of COVERAGE) {
    const result = auditCoverage([app('x', without(...capability.keys)), owned('x')]);
    assert.deepEqual(result.missing, [{ app: 'component:default/x', capability: capability.id }], capability.id);
    assert.equal(result.byCapability[capability.id].missing, 1);
  }
  const report = auditCatalog([app('x', without('onemind.io/recovery-runbook')), owned('x')], now);
  assert.equal(report.coverage.missing.length, 1);
});
test('coverage: either Kubernetes selector or kubernetes-id satisfies the Kubernetes capability', () => {
  const a = without('backstage.io/kubernetes-label-selector');
  a['backstage.io/kubernetes-id'] = 'x';
  assert.deepEqual(auditCoverage([app('x', a)]).missing, []);
});
test('coverage: a reviewed exception with a reason allows exactly one missing capability', () => {
  const entities = [app('x', without('argocd/app-name', 'onemind.io/logs-query')), owned('x')];
  const result = auditCoverage(entities, [{ app: 'component:default/x', capability: 'argocd', reason }]);
  assert.deepEqual(result.missing, [{ app: 'component:default/x', capability: 'logs' }]);
  assert.deepEqual(result.excepted, [{ app: 'component:default/x', capability: 'argocd' }]);
});
test('coverage: exceptions without a real reason, unknown capability or no app are invalid', () => {
  const entities = [app('x', without('argocd/app-name')), owned('x')];
  for (const bad of [
    { app: 'component:default/x', capability: 'argocd', reason: 'todo' },
    { app: 'component:default/x', capability: 'argocd' },
    { app: 'component:default/x', capability: 'colour', reason },
    { capability: 'argocd', reason },
  ]) {
    const result = auditCoverage(entities, [bad]);
    assert.equal(result.invalidExceptions.length, 1);
    assert.equal(result.missing.length, 1, 'an invalid exception grants nothing');
  }
});
test('coverage: stale exceptions (annotation present, or app gone) are reported for deletion', () => {
  const result = auditCoverage([app('x')], [
    { app: 'component:default/x', capability: 'argocd', reason },
    { app: 'component:default/gone', capability: 'logs', reason },
  ]);
  assert.deepEqual(result.staleExceptions.map(s => s.app), ['component:default/x', 'component:default/gone']);
});
test('coverage: only cluster apps are in scope, found by Argo name or by owning workloads', () => {
  const launcher = app('saas', {});
  const workload = { kind: 'Component', metadata: { name: 'w' }, spec: { subcomponentOf: 'component:default/owner' } };
  const owner = app('owner', {});
  const entities = [launcher, owner, workload, { kind: 'Resource', metadata: { name: 'r' } }, app('argo-app', { 'argocd/app-name': 'a' })];
  assert.deepEqual(kubernetesApps(entities).map(e => e.metadata.name).sort(), ['argo-app', 'owner']);
});
test('coverage: counts distinct linked Argo applications, including the extra ones on an app', () => {
  const a = app('x', { ...fullAnnotations(), 'onemind.io/argocd-apps': 'platform-x,platform-y' });
  assert.equal(auditCoverage([a]).argoApplicationsLinked, 2);
});
test('coverage: CLI exits 1 on a missing annotation and 0 when everything is present', () => {
  const { spawnSync } = require('node:child_process');
  const run = entities => spawnSync(process.execPath, [require.resolve('./audit-catalog.cjs'), '-'], { input: JSON.stringify(entities), encoding: 'utf8', env: { ...process.env, ONEMIND_COVERAGE_EXCEPTIONS: '/nonexistent/coverage-exceptions.yaml' } });
  assert.equal(run([app('x')]).status, 0);
  const failed = run([app('x', without('grafana/dashboard-selector')), owned('x')]);
  assert.equal(failed.status, 1);
  assert.match(failed.stdout, /grafana-dashboard/);
});
