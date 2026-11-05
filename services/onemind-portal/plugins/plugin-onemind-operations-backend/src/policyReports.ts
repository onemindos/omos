// Kyverno policy results, read from the PolicyReport objects Kyverno's reports controller
// already writes (wgpolicyk8s.io/v1alpha2), through the same read-only cluster access the
// Live Operations page uses. Nothing here installs or calls another workload.
//
// One namespaced PolicyReport exists per scanned resource (a Deployment, ReplicaSet, Pod,
// CronJob, ...). This module turns the list into a small, per-namespace summary the portal
// can match against a catalog app's workloads.

export type PolicyResultState = 'pass' | 'fail' | 'warn' | 'error' | 'skip';

export type KubePolicyReport = {
  metadata?: { name?: string; namespace?: string };
  scope?: { kind?: string; name?: string; namespace?: string };
  results?: {
    policy?: string;
    rule?: string;
    result?: string;
    severity?: string;
    category?: string;
    message?: string;
  }[];
};

export type Tally = Record<PolicyResultState, number>;

export type PolicyResource = {
  kind: string;
  name: string;
  /**
   * The controller this resource belongs to: itself for a Deployment, StatefulSet,
   * DaemonSet or CronJob; otherwise the controller of the same namespace whose name is the
   * longest prefix of this resource's name (a ReplicaSet or Pod). Matching on the longest
   * prefix keeps `backstage-postgres-0` with `backstage-postgres`, not with `backstage`.
   * Falls back to the resource's own name when no controller matches.
   */
  workload: string;
  summary: Tally;
};

export type PolicyFinding = {
  resource: { kind: string; name: string; workload: string };
  policy: string;
  rule: string;
  result: Exclude<PolicyResultState, 'pass' | 'skip'>;
  severity: string | null;
  message: string;
};

export type NamespacePolicy = {
  resources: PolicyResource[];
  findings: PolicyFinding[];
};

export type PolicySummary = {
  namespaces: Record<string, NamespacePolicy>;
  reports: number;
  /** True when findings were cut at MAX_FINDINGS; the tallies are always complete. */
  truncated: boolean;
};

export const MAX_FINDINGS = 500;
const MAX_MESSAGE = 300;
const STATES: PolicyResultState[] = ['pass', 'fail', 'warn', 'error', 'skip'];
const RANK: Record<string, number> = { error: 0, fail: 1, warn: 2 };

const blank = (): Tally => ({ pass: 0, fail: 0, warn: 0, error: 0, skip: 0 });
const isState = (value: unknown): value is PolicyResultState =>
  typeof value === 'string' && (STATES as string[]).includes(value);

const CONTROLLERS = ['Deployment', 'StatefulSet', 'DaemonSet', 'CronJob'];

export function ownerOf(name: string, kind: string, controllers: string[]): string {
  if (CONTROLLERS.includes(kind)) return name;
  let best = name;
  let bestLength = 0;
  for (const candidate of controllers) {
    if (name.startsWith(`${candidate}-`) && candidate.length > bestLength) {
      best = candidate;
      bestLength = candidate.length;
    }
  }
  return best;
}

/**
 * Pure. A result with an unrecognised state is not counted as a pass: it is left out of
 * every tally, so a Kyverno version that adds a state cannot inflate "passed".
 * A report that names no namespace or resource is skipped for the same reason.
 */
export function summarizePolicyReports(items: KubePolicyReport[]): PolicySummary {
  type Scanned = { namespace: string; kind: string; name: string; summary: Tally; results: NonNullable<KubePolicyReport['results']> };
  const scanned: Scanned[] = [];
  for (const report of items) {
    const namespace = report.scope?.namespace ?? report.metadata?.namespace;
    const kind = report.scope?.kind;
    const name = report.scope?.name;
    if (!namespace || !kind || !name) continue;
    const summary = blank();
    const results = (report.results ?? []).filter(item => isState(item.result));
    for (const item of results) summary[item.result as PolicyResultState] += 1;
    scanned.push({ namespace, kind, name, summary, results });
  }
  const controllers = new Map<string, string[]>();
  for (const { namespace, kind, name } of scanned) {
    if (CONTROLLERS.includes(kind)) controllers.set(namespace, [...(controllers.get(namespace) ?? []), name]);
  }
  const namespaces: Record<string, NamespacePolicy> = {};
  const findings: (PolicyFinding & { namespace: string })[] = [];
  for (const { namespace, kind, name, summary, results } of scanned) {
    const workload = ownerOf(name, kind, controllers.get(namespace) ?? []);
    (namespaces[namespace] ??= { resources: [], findings: [] }).resources.push({ kind, name, workload, summary });
    for (const item of results) {
      if (item.result === 'fail' || item.result === 'warn' || item.result === 'error') {
        findings.push({
          namespace,
          resource: { kind, name, workload },
          policy: item.policy ?? 'unknown policy',
          rule: item.rule ?? 'unknown rule',
          result: item.result,
          severity: item.severity ?? null,
          message: (item.message ?? '').slice(0, MAX_MESSAGE),
        });
      }
    }
  }
  findings.sort(
    (a, b) =>
      RANK[a.result] - RANK[b.result] ||
      a.namespace.localeCompare(b.namespace) ||
      a.resource.name.localeCompare(b.resource.name) ||
      a.policy.localeCompare(b.policy) ||
      a.rule.localeCompare(b.rule),
  );
  for (const { namespace, ...finding } of findings.slice(0, MAX_FINDINGS)) {
    namespaces[namespace].findings.push(finding);
  }
  for (const entry of Object.values(namespaces)) {
    entry.resources.sort((a, b) => a.kind.localeCompare(b.kind) || a.name.localeCompare(b.name));
  }
  return { namespaces, reports: scanned.length, truncated: findings.length > MAX_FINDINGS };
}
