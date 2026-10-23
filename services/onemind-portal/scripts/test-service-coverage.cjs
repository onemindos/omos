const { test } = require('node:test');
const assert = require('node:assert/strict');
const { spawnSync } = require('node:child_process');
const path = require('node:path');
const { reportServiceCoverage: report, FIELDS } = require('./report-service-coverage.cjs');
const entity = (name, annotations = {}, spec = {}, extra = {}) => ({
  kind: 'Component', metadata: { name, annotations }, spec, ...extra,
});
const row = input => report(input).components[0].fields;

test('canonical default namespace and noncomponents; no mutation, order or clock dependence', () => {
  const input = [entity('Z'), entity('a'), { kind: 'Resource', metadata: { name: 'r' } }];
  const before = JSON.stringify(input);
  const result = report(input);
  assert.deepEqual(result.components.map(e => e.ref), ['component:default/a', 'component:default/z']);
  assert.equal(result.ignoredNonComponents, 1);
  assert.equal(result.componentCount, 2);
  assert.deepEqual(result, report([...input].reverse()));
  assert.equal(JSON.stringify(input), before);
  assert.equal(report([]).componentCount, 0);
});

test('rejects envelopes, snapshots, duplicate canonical IDs and malformed structures', () => {
  for (const input of [null, {}, { items: [] }, [null], [{}], [entity('a'), entity('A')],
    [entity('a', {}, {}, { metadata: { name: 'a', namespace: '' } })],
    [entity('https://private.test')], [entity('a', { x: 3 })],
    [entity('a', {}, { dependsOn: 'resource:default/r' })],
    [entity('a', {}, { providesApis: [42] })], [entity('a', {}, { owner: 'https://private.test' })],
    [entity('a', {}, {}, { relations: [{ type: 'dependsOn', targetRef: 'r' }] })],
    [entity('a', {}, {}, { metadata: { name: 'a', links: [{ title: 'runbook' }] } })],
    [entity('a', {}, {}, { spec: [] })], new Array(10001).fill(entity('a'))]) {
    assert.throws(() => report(input), /Invalid complete/);
  }
});

test('available fields identify evidence kinds, never verified capabilities', () => {
  const fields = row([entity('k8s-nats-nats-0dee54ac', {
    'backstage.io/source-location': 'url:https://example.test/source',
    'backstage.io/kubernetes-label-selector': 'app=nats',
    'backstage.io/kubernetes-namespace': 'nats',
    'argocd/app-name': 'platform-nats',
    'grafana/dashboard-selector': 'uid == example',
    'prometheus.io/service-name': 'nats',
    'onemind.io/verification': 'live-verified',
  }, { owner: 'team', dependsOn: ['resource:default/cluster'], providesApis: ['contract'] }, {
    metadata: { name: 'k8s-nats-nats-0dee54ac', annotations: {
      'backstage.io/source-location': 'url:https://example.test/source',
      'backstage.io/kubernetes-label-selector': 'app=nats',
      'backstage.io/kubernetes-namespace': 'nats', 'argocd/app-name': 'platform-nats',
      'grafana/dashboard-selector': 'uid == example', 'prometheus.io/service-name': 'nats',
      'onemind.io/verification': 'live-verified',
    }, links: [{ type: 'runbook', url: 'https://example.test/recovery' }] },
  }), { kind: 'Resource', metadata: { name: 'cluster' }, spec: { type: 'kubernetes-cluster' } },
  { kind: 'API', metadata: { name: 'contract' }, spec: { definition: 'openapi: 3.0.3' } }]);
  for (const field of FIELDS) {
    assert.equal(fields[field].status, 'available', field);
    assert.equal(fields[field].verification, 'not-tested');
  }
  assert.equal(fields.aiApi.authenticatedAccess, 'not-tested');
  assert.equal(fields.aiApi.contracts.withDefinition, 1);
  assert.ok(fields.prometheus.evidence.some(e => e.mode === 'configured-mapping'));
});

test('an arbitrary Prometheus annotation or link cannot create a fixed native scope mapping', () => {
  const fields = row([entity('other', { 'prometheus.io/service-name': 'nats' })]);
  assert.equal(fields.prometheus.status, 'available');
  assert.equal(fields.prometheus.evidence.some(e => e.mode === 'configured-mapping'), false);
  assert.equal(fields.operations.status, 'missing');
  assert.equal(fields.kubernetesNamespace.status, 'missing'); // catalog default is not workload default
});

test('STAC Operations predicate requires the exact STAC API entity; no invented mappings', () => {
  assert.equal(row([entity('stac-browser')]).operations.status, 'missing');
  const e = entity('k8s-onetak-stac-api-746ca5f2');
  assert.equal(row([e]).operations.status, 'available');
  e.metadata.namespace = 'other';
  assert.equal(row([e]).operations.status, 'missing');
});

test('dependency and API relations deduplicate, resolve defaults and retain missing target counts', () => {
  const fields = row([entity('a', {}, { dependsOn: ['resource:r'], providesApis: ['x'] }, {
    relations: [{ type: 'dependsOn', targetRef: 'RESOURCE:DEFAULT/R' },
      { type: 'providesApi', targetRef: 'api:default/x' },
      { type: 'providesApi', targetRef: 'api:default/z' }],
  }), { kind: 'Resource', metadata: { name: 'r' } }, { kind: 'API', metadata: { name: 'x' } }]);
  assert.deepEqual(fields.dependency.targets, { recorded: 1, present: 1 });
  assert.deepEqual(fields.aiApi.contracts, { recorded: 2, present: 1, withDefinition: 0 });
  assert.equal(fields.location.status, 'missing');
});

test('consuming an API and generic TechDocs do not prove a provided AI interface or runbook', () => {
  const fields = row([entity('a', { 'backstage.io/techdocs-ref': 'dir:.' }, { consumesApis: ['x'] })]);
  assert.equal(fields.aiApi.status, 'missing');
  assert.equal(fields.runbook.status, 'missing');
});

test('not applicable requires separate reviewed exclusion; conflicts and unknown refs/fields fail', () => {
  const input = [entity('a')];
  const scope = [{ ref: 'component:default/a', fields: ['kubernetesSelector', 'argo'] }];
  assert.equal(report(input, scope).components[0].fields.argo.status, 'not-applicable');
  assert.equal(report(input).components[0].fields.argo.status, 'missing');
  for (const exclusions of [null, [{ ref: 'component:default/missing', fields: ['argo'] }],
    [{ ref: 'component:default/a', fields: ['invented'] }], [...scope, ...scope]])
    assert.throws(() => report(input, exclusions));
  assert.throws(() => report([entity('a', { 'argocd/app-name': 'app' })], scope));
});

test('reports and CLI failures never echo URLs, queries, metadata or parser snippets', () => {
  const privateText = 'https://private.test/path?token=do-not-print#private-fragment';
  const e = entity('a', { 'backstage.io/source-location': privateText,
    'argocd/app-name': privateText, 'grafana/dashboard-selector': privateText,
    'unknown-secret': privateText });
  e.metadata.title = privateText;
  e.metadata.description = privateText;
  e.metadata.links = [{ type: 'runbook', title: privateText, url: privateText }];
  const output = JSON.stringify(report([e]));
  for (const marker of ['https:', 'private.test', 'do-not-print', 'private-fragment', 'unknown-secret'])
    assert.equal(output.includes(marker), false);
  const cli = path.join(__dirname, 'report-service-coverage.cjs');
  for (const input of ['{"items":[]}', `[{"metadata": "${privateText}"}`, JSON.stringify([entity(privateText)])]) {
    const result = spawnSync(process.execPath, [cli, '-'], { input, encoding: 'utf8' });
    assert.equal(result.status, 2);
    assert.equal(result.stdout, '');
    assert.equal(result.stderr.includes('private.test'), false);
  }
  const success = spawnSync(process.execPath, [cli, '-'], { input: JSON.stringify([e]), encoding: 'utf8' });
  assert.equal(success.status, 0);
  assert.deepEqual(JSON.parse(success.stdout), report([e]));
});
