import type { ArgoCell, NamespaceCell, WarningEvent } from './contract';
import { computeNamespaceCell, podState, type PodInfo } from './health';

const HOUR_MS = 3_600_000;
const MAX_WARNINGS = 50;
const DOWN_REASONS = new Set([
  'CrashLoopBackOff',
  'ImagePullBackOff',
  'ErrImagePull',
]);

type Loose = any; // Kubernetes objects arrive as untyped JSON; only a few fields are read.

/** Reduce a Pod to the fields the bridge keeps. */
export function toPodInfo(
  raw: Loose,
  nowMs: number,
  previous?: PodInfo,
): PodInfo {
  const status = raw.status ?? {};
  const containers: Loose[] = [
    ...(status.containerStatuses ?? []),
    ...(status.initContainerStatuses ?? []),
  ];
  let restartCount = 0;
  let recentTerminations = 0;
  let waitingReason: string | undefined;
  for (const cs of containers) {
    restartCount += cs.restartCount ?? 0;
    const finished = Date.parse(cs.lastState?.terminated?.finishedAt ?? '');
    if (Number.isFinite(finished) && nowMs - finished < HOUR_MS)
      recentTerminations++;
    const reason = cs.state?.waiting?.reason;
    if (!waitingReason && reason && DOWN_REASONS.has(reason))
      waitingReason = reason;
  }
  const observedRestarts = (previous?.observedRestarts ?? []).filter(
    t => nowMs - t < HOUR_MS,
  );
  if (previous && restartCount > previous.restartCount) {
    for (let i = 0; i < restartCount - previous.restartCount; i++)
      observedRestarts.push(nowMs);
  }
  return {
    namespace: raw.metadata?.namespace ?? 'default',
    name: raw.metadata?.name ?? '',
    phase: status.phase ?? 'Unknown',
    ready: (status.conditions ?? []).some(
      (c: Loose) => c.type === 'Ready' && c.status === 'True',
    ),
    deleting: Boolean(raw.metadata?.deletionTimestamp),
    restartCount,
    recentTerminations,
    waitingReason,
    observedRestarts,
  };
}

export function toArgoCell(raw: Loose, nowMs: number): ArgoCell {
  const status = raw.status ?? {};
  const revision = status.sync?.revision;
  const changed = Date.parse(
    status.health?.lastTransitionTime ?? status.reconciledAt ?? '',
  );
  return {
    name: raw.metadata?.name ?? '',
    sync: status.sync?.status ?? 'Unknown',
    health: status.health?.status ?? 'Unknown',
    revision: typeof revision === 'string' ? revision.slice(0, 40) : undefined,
    changedAt: new Date(
      Number.isFinite(changed) ? changed : nowMs,
    ).toISOString(),
  };
}

export function toWarning(
  raw: Loose,
): (WarningEvent & { uid: string }) | undefined {
  const uid = raw.metadata?.uid;
  if (!uid) return undefined;
  const ts =
    raw.lastTimestamp ??
    raw.series?.lastObservedTime ??
    raw.eventTime ??
    raw.metadata?.creationTimestamp;
  const t = Date.parse(ts ?? '');
  const obj = raw.involvedObject ?? raw.regarding ?? {};
  return {
    uid,
    namespace: raw.metadata?.namespace ?? obj.namespace ?? '',
    object: `${obj.kind ?? 'Object'}/${obj.name ?? ''}`,
    reason: String(raw.reason ?? ''),
    message: String(raw.message ?? raw.note ?? '').slice(0, 300),
    ts: new Date(Number.isFinite(t) ? t : Date.now()).toISOString(),
    count: raw.count ?? raw.series?.count,
  };
}

/**
 * Cluster picture built from the pod, Warning event and Argo application streams, plus the
 * delta bookkeeping. Pure state: the watchers (kube.ts) feed it, the bridge reads it.
 */
export class ClusterState {
  private pods = new Map<string, PodInfo>();
  private readonly lastChange = new Map<string, { text: string; at: string }>();
  private argo = new Map<string, ArgoCell>();
  private warnings = new Map<string, WarningEvent & { uid: string }>();
  private readonly sentCells = new Map<string, string>();
  private readonly sentArgo = new Map<string, string>();
  private readonly sentWarnings = new Map<string, string>();
  private warmed = false;

  /** Full re-list. Keeps restart history for pods that survive. */
  replacePods(items: Loose[], nowMs: number): void {
    const next = new Map<string, PodInfo>();
    for (const raw of items) {
      const key = `${raw.metadata?.namespace}/${raw.metadata?.name}`;
      next.set(key, toPodInfo(raw, nowMs, this.pods.get(key)));
    }
    this.pods = next;
    this.warmed = true;
  }

  applyPod(phase: string, raw: Loose, nowMs: number): void {
    const key = `${raw.metadata?.namespace}/${raw.metadata?.name}`;
    const previous = this.pods.get(key);
    if (phase === 'DELETED') {
      if (previous) {
        this.pods.delete(key);
        this.note(previous.namespace, `pod ${previous.name} deleted`, nowMs);
      }
      return;
    }
    const info = toPodInfo(raw, nowMs, previous);
    this.pods.set(key, info);
    const before = previous ? podState(previous) : undefined;
    const after = podState(info);
    if (before !== after && (this.warmed || previous))
      this.note(info.namespace, `pod ${info.name} ${after}`, nowMs);
  }

  private note(namespace: string, text: string, nowMs: number): void {
    this.lastChange.set(namespace, { text, at: new Date(nowMs).toISOString() });
  }

  replaceArgo(items: Loose[], nowMs: number): void {
    const next = new Map<string, ArgoCell>();
    for (const raw of items) {
      const cell = toArgoCell(raw, nowMs);
      const old = this.argo.get(cell.name);
      if (
        old &&
        old.sync === cell.sync &&
        old.health === cell.health &&
        old.revision === cell.revision
      )
        cell.changedAt = old.changedAt;
      next.set(cell.name, cell);
    }
    this.argo = next;
  }

  applyArgo(phase: string, raw: Loose, nowMs: number): void {
    const cell = toArgoCell(raw, nowMs);
    if (phase === 'DELETED') {
      this.argo.delete(cell.name);
      return;
    }
    const old = this.argo.get(cell.name);
    if (
      old &&
      old.sync === cell.sync &&
      old.health === cell.health &&
      old.revision === cell.revision
    )
      cell.changedAt = old.changedAt;
    else if (old) cell.changedAt = new Date(nowMs).toISOString();
    this.argo.set(cell.name, cell);
  }

  replaceWarnings(items: Loose[]): void {
    this.warnings = new Map();
    for (const raw of items) this.addWarning(raw);
  }

  applyWarning(phase: string, raw: Loose): void {
    if (phase === 'DELETED') {
      this.warnings.delete(raw.metadata?.uid);
      return;
    }
    this.addWarning(raw);
  }

  private addWarning(raw: Loose): void {
    const w = toWarning(raw);
    if (!w) return;
    this.warnings.set(w.uid, w);
    if (this.warnings.size > MAX_WARNINGS * 2) {
      const keep = [...this.warnings.values()]
        .sort((a, b) => b.ts.localeCompare(a.ts))
        .slice(0, MAX_WARNINGS);
      this.warnings = new Map(keep.map(x => [x.uid, x]));
    }
  }

  cells(nowMs: number): NamespaceCell[] {
    const byNs = new Map<string, PodInfo[]>();
    for (const pod of this.pods.values()) {
      const list = byNs.get(pod.namespace);
      if (list) list.push(pod);
      else byNs.set(pod.namespace, [pod]);
    }
    return [...byNs.entries()]
      .map(([ns, pods]) => {
        const cell: NamespaceCell = computeNamespaceCell(ns, pods, nowMs);
        const change = this.lastChange.get(ns);
        if (change) {
          cell.lastChange = change.text;
          cell.lastChangeAt = change.at;
        }
        return cell;
      })
      .sort((a, b) => a.namespace.localeCompare(b.namespace));
  }

  argoCells(): ArgoCell[] {
    return [...this.argo.values()].sort((a, b) => a.name.localeCompare(b.name));
  }

  /** Kubernetes Warning events, newest first, at most 50. */
  warningList(): WarningEvent[] {
    return [...this.warnings.values()]
      .sort((a, b) => b.ts.localeCompare(a.ts))
      .slice(0, MAX_WARNINGS)
      .map(({ uid: _uid, ...w }) => w);
  }

  /** Only what changed since the last call; undefined when nothing did. */
  delta(
    nowMs: number,
  ):
    | {
        namespaces: NamespaceCell[];
        removedNamespaces: string[];
        argo: ArgoCell[];
        warnings: WarningEvent[];
      }
    | undefined {
    const cells = this.cells(nowMs);
    const namespaces: NamespaceCell[] = [];
    const seen = new Set<string>();
    for (const cell of cells) {
      seen.add(cell.namespace);
      const key = JSON.stringify(cell);
      if (this.sentCells.get(cell.namespace) !== key) {
        this.sentCells.set(cell.namespace, key);
        namespaces.push(cell);
      }
    }
    const removedNamespaces = [...this.sentCells.keys()].filter(
      ns => !seen.has(ns),
    );
    for (const ns of removedNamespaces) {
      this.sentCells.delete(ns);
      this.lastChange.delete(ns);
    }
    const argo: ArgoCell[] = [];
    for (const cell of this.argoCells()) {
      const key = JSON.stringify(cell);
      if (this.sentArgo.get(cell.name) !== key) {
        this.sentArgo.set(cell.name, key);
        argo.push(cell);
      }
    }
    const warnings: WarningEvent[] = [];
    const live = new Set<string>();
    for (const w of [...this.warnings.values()]
      .sort((a, b) => b.ts.localeCompare(a.ts))
      .slice(0, MAX_WARNINGS)) {
      live.add(w.uid);
      const key = `${w.ts}|${w.count ?? ''}|${w.message}`;
      if (this.sentWarnings.get(w.uid) !== key) {
        this.sentWarnings.set(w.uid, key);
        const { uid: _uid, ...rest } = w;
        warnings.push(rest);
      }
    }
    for (const uid of [...this.sentWarnings.keys()])
      if (!live.has(uid)) this.sentWarnings.delete(uid);
    if (
      !namespaces.length &&
      !removedNamespaces.length &&
      !argo.length &&
      !warnings.length
    )
      return undefined;
    return { namespaces, removedNamespaces, argo, warnings };
  }
}
