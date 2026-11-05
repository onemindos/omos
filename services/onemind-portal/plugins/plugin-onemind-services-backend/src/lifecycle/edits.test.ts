import YAML from 'yaml';
import {
  addPillarOverride,
  addPreparationExceptions,
  removeCoverageExceptions,
  bumpCoverageTestCount,
  bumpStructureTestCount,
  removeMapEntries,
  removeStructureApp,
} from './catalogEdits';
import { collectClusterFacts } from './sources';
import { joinChunks, removeDocuments, splitChunks } from './yamlText';
import { localSnapshot } from './testRepo';

describe('yaml text helpers', () => {
  it('splitting and joining reproduces every manifest in platform/ byte for byte', () => {
    let checked = 0;
    for (const [path, text] of localSnapshot()) {
      if (!path.startsWith('platform/') || !/\.ya?ml$/.test(path)) continue;
      expect(joinChunks(splitChunks(text))).toBe(text);
      checked += 1;
    }
    expect(checked).toBeGreaterThan(100);
  });

  it('removes only the matching documents and reports an empty file', () => {
    const text = 'kind: A\nmetadata:\n  name: a\n---\nkind: B\nmetadata:\n  name: b\n';
    const one = removeDocuments(text, r => r.name === 'a');
    expect(one.text).toBe('---\nkind: B\nmetadata:\n  name: b\n');
    expect(one.removed.map(r => r.name)).toEqual(['a']);
    expect(removeDocuments(text, () => true).text).toBeUndefined();
    expect(removeDocuments(text, () => false).text).toBe(text);
  });
});

describe('catalog input edits', () => {
  it('removes a mapping only when it points at the app', () => {
    const text = 'workloads:\n  a/x: component:default/one\n  a/y: component:default/two\n';
    const res = removeMapEntries(text, { ref: 'component:default/one', keys: ['a/x', 'a/y'] });
    expect(res.text).toBe('workloads:\n  a/y: component:default/two\n');
    expect(res.removed).toEqual(['a/x']);
  });

  it('removes an app from structure.yaml without touching other sections', () => {
    const text = 'apps:\n  - name: one\n    title: One\n  - name: two\n    title: Two\nplacement:\n  one: s\n  two: s\nfallbackSystems:\n  one: one\n';
    const out = removeStructureApp(text, 'one');
    expect(YAML.parse(out)).toEqual({ apps: [{ name: 'two', title: 'Two' }], placement: { two: 's' }, fallbackSystems: { one: 'one' } });
  });

  it('adds a pillar override to an empty or non-empty section', () => {
    expect(YAML.parse(addPillarOverride('pillarOverrides: {}\ntitles: {}\n', 'component:default/a', 'ge')).pillarOverrides).toEqual({ 'component:default/a': 'ge' });
    expect(YAML.parse(addPillarOverride('pillarOverrides:\n  component:default/b: so\ntitles: {}\n', 'component:default/a', 'ge')).pillarOverrides).toEqual({ 'component:default/b': 'so', 'component:default/a': 'ge' });
  });

  it('bumps the hard-coded Argo count in both tests, and refuses a changed layout', () => {
    const real = (path: string) => localSnapshot().get(`services/onemind-portal/scripts/${path}`)!;
    const a = bumpStructureTestCount(real('test-catalog-structure.cjs'), 1)!;
    const b = bumpCoverageTestCount(real('test-coverage.cjs'), 1)!;
    const n = Number(/inGit\.size, (\d+),/.exec(real('test-catalog-structure.cjs'))![1]);
    expect(a).toContain(`inGit.size, ${n + 1}, 'platform/ should define the ${n + 1} source`);
    expect(b).toContain(`argoApplicationsLinked, ${n + 1}`);
    expect(bumpStructureTestCount('nothing here', 1)).toBeUndefined();
  });
});

describe('coverage preparation exceptions', () => {
  it('adds the three exceptions and removes exactly them again', () => {
    const base = 'exceptions:\n  - app: component:default/other\n    capability: logs\n    reason: >-\n      keep this one around\n';
    const added = addPreparationExceptions(base, 'demo', 'Demo');
    const parsed = YAML.parse(added).exceptions;
    expect(parsed.filter((e: any) => e.app === 'component:default/demo').map((e: any) => e.capability)).toEqual(['kubernetes', 'grafana-alerts', 'logs']);
    expect(parsed.filter((e: any) => e.app.endsWith('/demo')).every((e: any) => e.reason.length >= 60)).toBe(true);
    expect(removeCoverageExceptions(added, 'component:default/demo')).toBe(base);
  });
});

describe('collectClusterFacts', () => {
  it('only lists resources the portal service account already reads', async () => {
    const seen: string[] = [];
    const facts = await collectClusterFacts(async path => {
      seen.push(path);
      if (path.endsWith('/statefulsets'))
        return [{ metadata: { namespace: 'n', name: 'db' }, spec: { replicas: 2, volumeClaimTemplates: [{ metadata: { name: 'data' } }], template: { spec: { containers: [{ image: `x@sha256:${'a'.repeat(64)}` }] } } }, status: { readyReplicas: 2 } }];
      return [];
    }, () => new Date('2026-11-05T00:00:00Z'));
    expect(seen.sort()).toEqual([
      '/api/v1/persistentvolumeclaims',
      '/apis/apps/v1/daemonsets',
      '/apis/apps/v1/deployments',
      '/apis/apps/v1/statefulsets',
      '/apis/argoproj.io/v1alpha1/namespaces/argocd/applications',
      '/apis/batch/v1/cronjobs',
      '/apis/networking.k8s.io/v1/ingresses',
    ]);
    expect(facts.workloads[0].claims).toEqual(['data-db-0', 'data-db-1']);
    expect(facts.observedAt).toBe('2026-11-05T00:00:00.000Z');
  });
});
