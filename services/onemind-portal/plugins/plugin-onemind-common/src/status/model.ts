/**
 * The single definition of an app's status dot (flagship F5).
 *
 * Four signals are combined: alerts, Argo CD, workload readiness and uptime checks. A signal that is
 * not wired for an app is left out (it says nothing). A signal whose source failed or is stale is
 * `unknown`. The dot is `ok` only when at least one signal is `ok` and none is `unknown`, `warn` or
 * `down`: absence of data is never green.
 */
export type AppStatus = 'ok' | 'warn' | 'down' | 'unknown';
export type SignalName = 'alerts' | 'argo' | 'workloads' | 'uptime';

export type Signal = {
  name: SignalName;
  status: AppStatus;
  /** One short sentence for the tooltip. */
  detail: string;
  observedAt?: string;
};

export type AppStatusResult = {
  status: AppStatus;
  signals: Signal[];
  /** Oldest observation among the signals; undefined when nothing was observed. */
  observedAt?: string;
};

/** What a page or an app entity tells us to look for in a snapshot. */
export type AppTarget = {
  /** Argo CD application names. */
  argoApps: string[];
  /** `<namespace>/<workload>` of every live workload. */
  workloads: string[];
  /** Catalog refs whose uptime checks count (the app and its workloads). */
  refs: string[];
  /** Scoped Prometheus alert scopes (services backend allow-list). */
  alertScopes: string[];
};

export type Reading<T> = { state: 'observed'; observedAt: string; value: T } | { state: 'unavailable'; reason?: string };

/**
 * Normalized observation, independent of where it came from. Today `createOperationsStatusSource`
 * fills it from the existing status and Prometheus routes; the Live Ops aggregator (F4) can fill the
 * same shape later without touching the model, the hook or any component.
 */
export type StatusSnapshot = {
  observedAt: string;
  /** Argo CD applications that are not Healthy or not Synced, by name. */
  argo?: Reading<{ name: string; health: string; sync: string }[]>;
  /** `<namespace>/<workload>` with fewer ready replicas than desired. */
  notReady?: Reading<string[]>;
  /** Uptime checks by catalog ref. */
  uptime?: Reading<{ entity?: string; name: string; state: string }[]>;
  /** Firing and pending alerts by scoped Prometheus scope. */
  alerts?: Record<string, Reading<{ firing: number; critical: number; pending: number }>>;
};

export const STATUS_STALE_MS = 180000;

const annotation = (entity: EntityLike, key: string) => entity.metadata.annotations?.[key] ?? '';
const list = (value: string) => value.split(',').map(x => x.trim()).filter(Boolean);

export type EntityLike = {
  kind: string;
  metadata: { name: string; namespace?: string; annotations?: Record<string, string> };
  relations?: { type: string; targetRef: string }[];
};

export const refOfEntity = (entity: EntityLike) =>
  `${entity.kind}:${entity.metadata.namespace ?? 'default'}/${entity.metadata.name}`.toLowerCase();

/** Reads the generated coverage annotations (scripts/curate-coverage.cjs); nothing is guessed. */
export function targetForEntity(entity: EntityLike): AppTarget {
  const argo = [annotation(entity, 'argocd/app-name'), ...list(annotation(entity, 'onemind.io/argocd-apps'))];
  const workloads = list(annotation(entity, 'onemind.io/kubernetes-workloads'));
  const ns = annotation(entity, 'backstage.io/kubernetes-namespace');
  const own = annotation(entity, 'onemind.io/kubernetes-workload');
  // A single workload entity (not an app) carries its own namespace and workload name.
  if (!workloads.length && ns && own) workloads.push(`${ns}/${own}`);
  const scopes = [...list(annotation(entity, 'onemind.io/prometheus-scopes')), annotation(entity, 'prometheus.io/service-name')];
  return {
    argoApps: [...new Set(argo.filter(Boolean))],
    workloads: [...new Set(workloads)],
    refs: [
      refOfEntity(entity),
      ...(entity.relations ?? []).filter(r => r.type === 'hasPart').map(r => r.targetRef.toLowerCase()),
    ],
    alertScopes: [...new Set(scopes.filter(Boolean))],
  };
}

const RANK: Record<AppStatus, number> = { ok: 0, unknown: 1, warn: 2, down: 3 };
const worst = (a: AppStatus, b: AppStatus) => (RANK[b] > RANK[a] ? b : a);

function argoSignal(target: AppTarget, reading: NonNullable<StatusSnapshot['argo']>): Signal {
  if (reading.state === 'unavailable') return { name: 'argo', status: 'unknown', detail: 'Argo CD could not be read.' };
  let status: AppStatus = 'ok';
  const notes: string[] = [];
  for (const name of target.argoApps) {
    const hit = reading.value.find(app => app.name === name);
    if (!hit) continue;
    const health = hit.health;
    const level: AppStatus = health === 'Degraded' || health === 'Missing' ? 'down' : health === 'Healthy' || health === 'Progressing' || health === 'Suspended' ? 'warn' : 'unknown';
    status = worst(status, level);
    notes.push(`${name}: ${health}, ${hit.sync}`);
  }
  return {
    name: 'argo',
    status,
    detail: notes.length ? notes.join('; ') : `${target.argoApps.length} Argo CD application${target.argoApps.length === 1 ? '' : 's'} healthy and synced.`,
    observedAt: reading.observedAt,
  };
}

function workloadSignal(target: AppTarget, reading: NonNullable<StatusSnapshot['notReady']>): Signal {
  if (reading.state === 'unavailable') return { name: 'workloads', status: 'unknown', detail: 'Workload readiness could not be read.' };
  const mine = new Set(target.workloads);
  const down = reading.value.filter(name => mine.has(name));
  if (!down.length) return { name: 'workloads', status: 'ok', detail: `All ${mine.size} workload${mine.size === 1 ? '' : 's'} have their desired replicas ready.`, observedAt: reading.observedAt };
  const status: AppStatus = down.length >= mine.size ? 'down' : 'warn';
  return { name: 'workloads', status, detail: `${down.length} of ${mine.size} workloads below desired replicas: ${down.slice(0, 3).join(', ')}${down.length > 3 ? '…' : ''}.`, observedAt: reading.observedAt };
}

function uptimeSignal(target: AppTarget, reading: NonNullable<StatusSnapshot['uptime']>): Signal | undefined {
  if (reading.state === 'unavailable') return { name: 'uptime', status: 'unknown', detail: 'Uptime checks could not be read.' };
  const refs = new Set(target.refs);
  const checks = reading.value.filter(check => check.entity && refs.has(check.entity.toLowerCase()));
  if (!checks.length) return undefined;
  let status: AppStatus = 'ok';
  const notes: string[] = [];
  for (const check of checks) {
    // A page that asks for sign-in answered: reachable. Anything else that is not "up" is not green.
    const level: AppStatus = check.state === 'up' || check.state === 'authentication-required' ? 'ok' : check.state === 'down' ? 'down' : 'unknown';
    status = worst(status, level);
    if (level !== 'ok') notes.push(`${check.name}: ${check.state}`);
  }
  return { name: 'uptime', status, detail: notes.length ? notes.join('; ') : `${checks.length} uptime check${checks.length === 1 ? '' : 's'} passing.`, observedAt: reading.observedAt };
}

function alertSignal(target: AppTarget, alerts: NonNullable<StatusSnapshot['alerts']>): Signal | undefined {
  const readings = target.alertScopes.map(scope => ({ scope, reading: alerts[scope] })).filter(x => x.reading);
  if (!readings.length) return undefined;
  let status: AppStatus = 'ok';
  const notes: string[] = [];
  let observedAt: string | undefined;
  for (const { scope, reading } of readings) {
    if (reading!.state === 'unavailable') { status = worst(status, 'unknown'); notes.push(`${scope}: alerts unavailable`); continue; }
    const { firing, critical, pending } = reading!.value;
    observedAt = observedAt && observedAt < reading!.observedAt ? observedAt : reading!.observedAt;
    if (critical > 0) { status = worst(status, 'down'); notes.push(`${scope}: ${critical} critical firing`); }
    else if (firing > 0 || pending > 0) { status = worst(status, 'warn'); notes.push(`${scope}: ${firing} firing, ${pending} pending`); }
  }
  return { name: 'alerts', status, detail: notes.length ? notes.join('; ') : 'No firing or pending alerts in the scoped rules.', observedAt };
}

/** Combines the four signals for one app. `snapshot` undefined means nothing has been read yet. */
export function combineStatus(target: AppTarget, snapshot: StatusSnapshot | undefined, now = Date.now()): AppStatusResult {
  if (!snapshot) return { status: 'unknown', signals: [] };
  const stale = now - Date.parse(snapshot.observedAt) > STATUS_STALE_MS || !Number.isFinite(Date.parse(snapshot.observedAt));
  const signals: Signal[] = [];
  if (target.argoApps.length && snapshot.argo) signals.push(argoSignal(target, snapshot.argo));
  if (target.workloads.length && snapshot.notReady) signals.push(workloadSignal(target, snapshot.notReady));
  if (snapshot.uptime) { const s = uptimeSignal(target, snapshot.uptime); if (s) signals.push(s); }
  if (snapshot.alerts) { const s = alertSignal(target, snapshot.alerts); if (s) signals.push(s); }
  // A stale snapshot may no longer describe the cluster: a problem stays visible, health does not.
  const adjusted = stale
    ? signals.map(s => (s.status === 'ok' ? { ...s, status: 'unknown' as const, detail: `${s.detail} (observation is older than 3 minutes)` } : s))
    : signals;
  const observed = adjusted.map(s => s.observedAt).filter((x): x is string => !!x).sort();
  if (!adjusted.length) return { status: 'unknown', signals: [] };
  return { status: adjusted.reduce<AppStatus>((acc, s) => worst(acc, s.status), 'ok'), signals: adjusted, observedAt: observed[0] };
}

/** Counts for a System or Domain summary. */
export function summarize(results: AppStatusResult[]): Record<AppStatus, number> {
  const out: Record<AppStatus, number> = { ok: 0, warn: 0, down: 0, unknown: 0 };
  for (const r of results) out[r.status]++;
  return out;
}

export const statusRank = (status: AppStatus) => RANK[status];
