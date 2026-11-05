// Pure helpers for the Kyverno "Policy results" card. The wire shape is produced by
// plugins/plugin-onemind-operations-backend/src/policyReports.ts (GET /api/onemind/policy-reports).

export type Tally = {
  pass: number;
  fail: number;
  warn: number;
  error: number;
  skip: number;
};
export type PolicyResource = {
  kind: string;
  name: string;
  workload: string;
  summary: Tally;
};
export type PolicyFinding = {
  resource: { kind: string; name: string; workload: string };
  policy: string;
  rule: string;
  result: 'fail' | 'warn' | 'error';
  severity: string | null;
  message: string;
};
export type PolicyResponse =
  | {
      status: 'ok' | 'stale';
      observedAt: string;
      error?: string;
      reports: number;
      truncated: boolean;
      namespaces: Record<
        string,
        { resources: PolicyResource[]; findings: PolicyFinding[] }
      >;
    }
  | { status: 'error'; observedAt: null; error: string };

type EntityLike = {
  metadata: { annotations?: Record<string, string> };
};

const NAMESPACE = 'backstage.io/kubernetes-namespace';
const NAMESPACES = 'onemind.io/kubernetes-namespaces';
const WORKLOADS = 'onemind.io/kubernetes-workloads';
const WORKLOAD = 'onemind.io/kubernetes-workload';

/** `workloads: null` means the entity names no workloads, so the whole namespace is shown. */
export type Target = { namespace: string; workloads: string[] | null };

const split = (value?: string) =>
  (value ?? '')
    .split(',')
    .map(item => item.trim())
    .filter(Boolean);

/** What the catalog says this entity runs: its namespaces and, when known, its workloads. */
export function targetsForEntity(entity: EntityLike): Target[] {
  const annotations = entity.metadata.annotations ?? {};
  const namespaces = [
    ...new Set([
      ...split(annotations[NAMESPACE]),
      ...split(annotations[NAMESPACES]),
    ]),
  ];
  const named = new Map<string, string[]>();
  for (const item of split(annotations[WORKLOADS])) {
    const [namespace, ...rest] = item.split('/');
    if (namespace && rest.length)
      named.set(namespace, [...(named.get(namespace) ?? []), rest.join('/')]);
  }
  const single = annotations[WORKLOAD];
  if (single && annotations[NAMESPACE])
    named.set(annotations[NAMESPACE], [
      ...(named.get(annotations[NAMESPACE]) ?? []),
      single,
    ]);
  return namespaces.map(namespace => ({
    namespace,
    workloads: named.has(namespace) ? [...new Set(named.get(namespace))] : null,
  }));
}

export type Matched = {
  scope: 'workloads' | 'namespace';
  resources: PolicyResource[];
  findings: PolicyFinding[];
  /** Workloads the catalog names that Kyverno has not reported on. */
  unreported: string[];
};

/** Workloads are matched exactly on the `workload` the backend resolved for each resource. */
export function matchPolicy(
  response: Extract<PolicyResponse, { status: 'ok' | 'stale' }>,
  targets: Target[],
): Matched {
  const resources: PolicyResource[] = [];
  const findings: PolicyFinding[] = [];
  const unreported: string[] = [];
  for (const { namespace, workloads } of targets) {
    const entry = response.namespaces[namespace];
    const wanted = workloads ? new Set(workloads) : undefined;
    const keep = (workload: string) => !wanted || wanted.has(workload);
    for (const resource of entry?.resources ?? [])
      if (keep(resource.workload)) resources.push(resource);
    for (const finding of entry?.findings ?? [])
      if (keep(finding.resource.workload)) findings.push(finding);
    for (const workload of workloads ?? []) {
      if (!entry?.resources.some(resource => resource.workload === workload))
        unreported.push(`${namespace}/${workload}`);
    }
  }
  return {
    scope: targets.some(target => target.workloads) ? 'workloads' : 'namespace',
    resources,
    findings,
    unreported,
  };
}

export type Counts = {
  /** Distinct failing or erroring policy rules, so one bad Deployment is not counted
   * again for its ReplicaSet and Pods (Kyverno reports each of them). */
  violations: number;
  warnings: number;
  passed: number;
  skipped: number;
  resources: number;
};

export function countsFor(matched: Matched): Counts {
  const distinct = (states: string[]) =>
    new Set(
      matched.findings
        .filter(finding => states.includes(finding.result))
        .map(
          finding =>
            `${finding.resource.workload}|${finding.policy}|${finding.rule}`,
        ),
    ).size;
  const sum = (field: keyof Tally) =>
    matched.resources.reduce((n, resource) => n + resource.summary[field], 0);
  return {
    violations: distinct(['fail', 'error']),
    warnings: distinct(['warn']),
    passed: sum('pass'),
    skipped: sum('skip'),
    resources: matched.resources.length,
  };
}
