import { getEvidenceHttp } from '@internal/plugin-onemind-common-backend';

const target =
  'http://monitoring-kube-prometheus-prometheus.monitoring.svc:9090/api/v1/rules?type=alert';
const customNats = new Set([
  'NATSQuorumLost',
  'NATSDown',
  'TAKCoTFlowStopped',
  'NATSConsumerBacklogGrowing',
  'NATSSurveyorDown',
]);
export const alertScopes = ['backstage', 'postgres', 'nats'] as const;
export type AlertScope = (typeof alertScopes)[number];
type Labels = Record<string, string>;
type Row = {
  name: string;
  state: string;
  health: string;
  lastEvaluation: string;
  labels: Labels;
  annotations: Labels;
  query: string;
  alerts: {
    labels: Labels;
    annotations: Labels;
    state: string;
    activeAt: string;
    value: string;
  }[];
};
type Source = {
  status: 'success';
  data: { groups: { name: string; rules: Row[] }[] };
};
export type AlertObservation = Source & {
  schemaVersion: 'onemind.prometheus-alerts.v1';
  observedAt: string;
  scope: AlertScope;
  coverage: 'observed' | 'partial' | 'not-configured';
  candidateRules: number;
  unhealthyRules: number;
  staleRules: number;
};
const record = (v: unknown): v is Record<string, unknown> =>
  !!v && typeof v === 'object' && !Array.isArray(v);
const strings = (v: unknown): v is Labels =>
  record(v) &&
  Object.entries(v).every(
    ([k, x]) => k.length < 256 && typeof x === 'string' && x.length < 20000,
  );
const date = (v: unknown): v is string =>
  typeof v === 'string' && Number.isFinite(Date.parse(v));

// Fail closed on malformed/partial protocol responses, not as an empty healthy list.
export function parseAlertSource(value: unknown): Source {
  if (
    !record(value) ||
    value.status !== 'success' ||
    !record(value.data) ||
    !Array.isArray(value.data.groups) ||
    value.data.groups.length > 2000 ||
    (Array.isArray(value.warnings) && value.warnings.length)
  )
    throw new Error('Invalid alert source');
  let count = 0;
  for (const group of value.data.groups) {
    if (
      !record(group) ||
      typeof group.name !== 'string' ||
      !Array.isArray(group.rules)
    )
      throw new Error('Invalid rule group');
    for (const rule of group.rules) {
      if (
        ++count > 10000 ||
        !record(rule) ||
        rule.type !== 'alerting' ||
        typeof rule.name !== 'string' ||
        typeof rule.query !== 'string' ||
        !['inactive', 'pending', 'firing'].includes(String(rule.state)) ||
        !['ok', 'err', 'unknown'].includes(String(rule.health)) ||
        !date(rule.lastEvaluation) ||
        !strings(rule.labels) ||
        !strings(rule.annotations) ||
        !Array.isArray(rule.alerts) ||
        rule.alerts.length > 10000
      )
        throw new Error('Invalid alert rule');
      for (const alert of rule.alerts) {
        if (
          !record(alert) ||
          !['pending', 'firing'].includes(String(alert.state)) ||
          !strings(alert.labels) ||
          !strings(alert.annotations) ||
          !date(alert.activeAt) ||
          typeof alert.value !== 'string'
        )
          throw new Error('Invalid alert instance');
      }
    }
  }
  // Never expose upstream error details, file paths or undocumented response fields.
  const source = value as unknown as Source;
  return {
    status: 'success',
    data: {
      groups: source.data.groups.map(group => ({
        name: group.name,
        rules: group.rules.map(rule => ({
          name: rule.name,
          state: rule.state,
          health: rule.health,
          lastEvaluation: rule.lastEvaluation,
          labels: rule.labels,
          annotations: rule.annotations,
          query: rule.query,
          alerts: rule.alerts.map(alert => ({
            labels: alert.labels,
            annotations: alert.annotations,
            state: alert.state,
            activeAt: alert.activeAt,
            value: alert.value,
          })),
        })),
      })),
    },
  };
}
function matches(scope: AlertScope, rule: Row, labels: Labels) {
  // Service-level bus rules deliberately include its surveyor and TAK flow dependencies.
  if (
    scope === 'nats' &&
    customNats.has(rule.name) &&
    labels.service === 'nats' &&
    (!labels.namespace || labels.namespace === 'nats')
  )
    return true;
  const ns = scope === 'nats' ? 'nats' : 'backstage';
  if (labels.namespace !== ns) return false;
  if (scope === 'backstage')
    return (
      labels.deployment === 'backstage' ||
      /^backstage-[a-z0-9]+-[a-z0-9]{5}$/.test(labels.pod ?? '')
    );
  const workload = scope === 'nats' ? 'nats' : 'backstage-postgres';
  return (
    labels.statefulset === workload ||
    new RegExp(`^${workload}-[0-9]+$`).test(labels.pod ?? '')
  );
}
function coverageState(
  count: number,
  problems: number,
): AlertObservation['coverage'] {
  if (!count) return 'not-configured';
  return problems ? 'partial' : 'observed';
}
export function scopeAlertSource(
  source: Source,
  scope: AlertScope,
  now = Date.now(),
): AlertObservation {
  let candidateRules = 0;
  let unhealthyRules = 0;
  let staleRules = 0;
  const groups = source.data.groups
    .map(group => {
      const rules = group.rules.flatMap(rule => {
        const alerts = rule.alerts.filter(a =>
          matches(scope, rule, { ...rule.labels, ...a.labels }),
        );
        // Generic Kubernetes families are candidates, not proof of complete service coverage.
        const family =
          /^(KubePod|PodCrashLooping|KubeContainer|KubeDeployment|KubeStatefulSet)/.test(
            rule.name,
          );
        const related =
          alerts.length > 0 || matches(scope, rule, rule.labels) || family;
        if (!related) return [];
        candidateRules++;
        if (rule.health !== 'ok') unhealthyRules++;
        const age = now - Date.parse(rule.lastEvaluation);
        if (age > 180000 || age < -30000) staleRules++;
        return [{ ...rule, alerts }];
      });
      return { name: group.name, rules };
    })
    .filter(g => g.rules.length);
  return {
    status: 'success',
    data: { groups },
    schemaVersion: 'onemind.prometheus-alerts.v1',
    observedAt: new Date(now).toISOString(),
    scope,
    candidateRules,
    unhealthyRules,
    staleRules,
    coverage: coverageState(candidateRules, unhealthyRules + staleRules),
  };
}

// One fixed internal GET, 5-second absolute deadline/2 MB response cap, no redirects
// or incoming credential forwarding. Cache/in-flight state contains only validated reads.
export function createAlertReader(get = getEvidenceHttp, clock = Date.now) {
  let cache: { at: number; source: Source } | undefined;
  let pending: Promise<Source> | undefined;
  return async (scope: AlertScope) => {
    if (!alertScopes.includes(scope)) throw new Error('Unknown alert scope');
    if (!cache || clock() - cache.at >= 15000) {
      if (!pending)
        pending = (async () => {
          const response = await get(target);
          if (response.status !== 200)
            throw new Error('Alert source unavailable');
          const source = parseAlertSource(JSON.parse(response.body));
          cache = { at: clock(), source };
          return source;
        })().finally(() => {
          pending = undefined;
        });
      await pending;
    }
    const observation = scopeAlertSource(cache!.source, scope, clock());
    observation.observedAt = new Date(cache!.at).toISOString();
    return observation;
  };
}
export const readPrometheusAlerts = createAlertReader();
