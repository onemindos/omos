import {
  MAX_FINDINGS,
  ownerOf,
  summarizePolicyReports,
  type KubePolicyReport,
} from './policyReports';
import real from './__fixtures__/policyreports-real.json';

const report = (
  namespace: string,
  kind: string,
  name: string,
  results: {
    result: string;
    policy?: string;
    rule?: string;
    message?: string;
    severity?: string;
  }[],
): KubePolicyReport => ({
  metadata: { name: `${name}-uid`, namespace },
  scope: { kind, name, namespace },
  results,
});

describe('summarizePolicyReports on real Kyverno reports (backstage and nats namespaces)', () => {
  const summary = summarizePolicyReports(real as KubePolicyReport[]);

  it('counts one resource per report and keeps namespaces apart', () => {
    expect(summary.reports).toBe(real.length);
    expect(Object.keys(summary.namespaces).sort()).toEqual([
      'backstage',
      'nats',
    ]);
  });

  it('keeps backstage-postgres out of the backstage workload', () => {
    const resources = summary.namespaces.backstage.resources;
    const postgresPod = resources.find(r => r.name === 'backstage-postgres-0');
    expect(postgresPod?.workload).toBe('backstage-postgres');
    expect(
      resources.find(r => r.kind === 'Deployment' && r.name === 'backstage')
        ?.workload,
    ).toBe('backstage');
  });

  it('tallies passes and skips, and has no findings when nothing failed', () => {
    const total = Object.values(summary.namespaces)
      .flatMap(ns => ns.resources)
      .reduce((sum, r) => sum + r.summary.pass + r.summary.skip, 0);
    expect(total).toBeGreaterThan(0);
    expect(
      Object.values(summary.namespaces).every(ns => ns.findings.length === 0),
    ).toBe(true);
    expect(summary.truncated).toBe(false);
  });
});

describe('summarizePolicyReports findings', () => {
  const items = [
    report('app', 'Deployment', 'web', [
      { result: 'pass', policy: 'p1', rule: 'r1' },
      { result: 'warn', policy: 'p2', rule: 'r2', message: 'careful' },
    ]),
    report('app', 'Pod', 'web-7d9f8c6b5-abcde', [
      {
        result: 'fail',
        policy: 'p3',
        rule: 'r3',
        severity: 'high',
        message: 'x'.repeat(1000),
      },
      { result: 'error', policy: 'p4', rule: 'r4' },
      { result: 'skip', policy: 'p5', rule: 'r5' },
    ]),
  ];
  const { namespaces } = summarizePolicyReports(items);

  it('lists errors and failures before warnings, with the owning workload', () => {
    const findings = namespaces.app.findings;
    expect(findings.map(f => f.result)).toEqual(['error', 'fail', 'warn']);
    expect(findings.every(f => f.resource.workload === 'web')).toBe(true);
  });

  it('truncates long messages and keeps severity', () => {
    const fail = namespaces.app.findings.find(f => f.result === 'fail')!;
    expect(fail.message).toHaveLength(300);
    expect(fail.severity).toBe('high');
  });

  it('tallies every state on the resource', () => {
    const pod = namespaces.app.resources.find(r => r.kind === 'Pod')!;
    expect(pod.summary).toEqual({
      pass: 0,
      fail: 1,
      warn: 0,
      error: 1,
      skip: 1,
    });
  });
});

describe('summarizePolicyReports is conservative', () => {
  it('does not count an unknown result state as a pass', () => {
    const { namespaces } = summarizePolicyReports([
      report('a', 'Deployment', 'x', [
        { result: 'mystery' },
        { result: 'pass' },
      ]),
    ]);
    expect(namespaces.a.resources[0].summary).toEqual({
      pass: 1,
      fail: 0,
      warn: 0,
      error: 0,
      skip: 0,
    });
  });

  it('skips reports that name no namespace or resource', () => {
    const out = summarizePolicyReports([
      { results: [{ result: 'fail' }] },
      { metadata: { namespace: 'a' }, results: [{ result: 'fail' }] },
    ]);
    expect(out.reports).toBe(0);
    expect(out.namespaces).toEqual({});
  });

  it('caps findings but never the tallies, and says so', () => {
    const many = Array.from({ length: MAX_FINDINGS + 5 }, (_, i) =>
      report('a', 'Deployment', `d${String(i).padStart(4, '0')}`, [
        { result: 'fail', policy: 'p', rule: 'r' },
      ]),
    );
    const out = summarizePolicyReports(many);
    expect(out.namespaces.a.findings).toHaveLength(MAX_FINDINGS);
    expect(out.truncated).toBe(true);
    expect(
      out.namespaces.a.resources.reduce((n, r) => n + r.summary.fail, 0),
    ).toBe(MAX_FINDINGS + 5);
  });
});

describe('ownerOf', () => {
  it('uses the longest controller prefix and the resource itself otherwise', () => {
    const controllers = ['backstage', 'backstage-postgres'];
    expect(ownerOf('backstage-5d9f-abcde', 'Pod', controllers)).toBe(
      'backstage',
    );
    expect(ownerOf('backstage-postgres-0', 'Pod', controllers)).toBe(
      'backstage-postgres',
    );
    expect(ownerOf('lonely', 'Pod', controllers)).toBe('lonely');
    expect(ownerOf('backstage', 'Deployment', controllers)).toBe('backstage');
    // A name that merely starts with the same letters is not a child.
    expect(ownerOf('backstagefoo-1', 'Pod', controllers)).toBe(
      'backstagefoo-1',
    );
  });
});
