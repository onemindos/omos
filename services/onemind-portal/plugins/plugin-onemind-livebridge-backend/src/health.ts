import type { CellHealth, NamespaceCell } from './contract';

const HOUR_MS = 3_600_000;

/** The few fields of a Pod the bridge keeps. Slim on purpose: memory stays bounded. */
export interface PodInfo {
  namespace: string;
  name: string;
  phase: string;
  ready: boolean;
  deleting: boolean;
  /** Sum of restartCount over containers. */
  restartCount: number;
  /** Containers whose last termination finished within the last hour, at the time of the update. */
  recentTerminations: number;
  /** Waiting reason that matters (CrashLoopBackOff, ImagePullBackOff, ErrImagePull), if any. */
  waitingReason?: string;
  /** Times at which this bridge saw restartCount grow (one entry per restart). Pruned to 1 h. */
  observedRestarts: number[];
}

const DOWN_REASONS = new Set([
  'CrashLoopBackOff',
  'ImagePullBackOff',
  'ErrImagePull',
]);

export function restartsLastHour(pod: PodInfo, nowMs: number): number {
  const observed = pod.observedRestarts.filter(t => nowMs - t < HOUR_MS).length;
  return Math.max(observed, pod.recentTerminations);
}

/** Short state word used for "pod X <state>" change notes. */
export function podState(pod: PodInfo): string {
  if (pod.waitingReason && DOWN_REASONS.has(pod.waitingReason))
    return pod.waitingReason;
  if (pod.phase === 'Failed') return 'Failed';
  if (pod.phase === 'Pending') return 'Pending';
  if (pod.phase === 'Running' && !pod.ready && !pod.deleting) return 'NotReady';
  return pod.phase;
}

/**
 * Health of one namespace from its pods.
 *  down: any failed pod, or any pod in CrashLoopBackOff / ImagePullBackOff / ErrImagePull
 *  warn: otherwise, any pending pod, any running pod not ready, or any restart in the last hour
 *  ok:   everything else
 * Completed (Succeeded) pods, such as finished Jobs, are not counted.
 */
export function computeNamespaceCell(
  namespace: string,
  pods: PodInfo[],
  nowMs: number,
): Omit<NamespaceCell, 'lastChange' | 'lastChangeAt'> {
  let total = 0;
  let ready = 0;
  let pending = 0;
  let failed = 0;
  let restarts = 0;
  let broken = false;
  let notReady = false;
  for (const pod of pods) {
    if (pod.phase === 'Succeeded') continue;
    total++;
    restarts += restartsLastHour(pod, nowMs);
    if (pod.phase === 'Failed') failed++;
    else if (pod.phase === 'Pending') pending++;
    else if (pod.ready) ready++;
    else if (pod.phase === 'Running' && !pod.deleting) notReady = true;
    if (pod.waitingReason && DOWN_REASONS.has(pod.waitingReason)) broken = true;
  }
  let health: CellHealth = 'ok';
  if (failed > 0 || broken) health = 'down';
  else if (pending > 0 || notReady || restarts > 0) health = 'warn';
  return { namespace, total, ready, pending, failed, restarts, health };
}
