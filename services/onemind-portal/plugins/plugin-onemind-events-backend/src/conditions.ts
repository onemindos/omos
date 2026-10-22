import { OUT_OF_SYNC_GRACE_MS } from '@internal/plugin-onemind-common-backend';
import type { ConditionRow } from './store';

export type Reason = ConditionRow['reason'];

// A bad object seen in this scan.
export type ActiveCondition = { scope: string; reason: Reason; resource: string };

export type ConditionPlan = {
  // Conditions that just became notifiable (grace elapsed or immediate reasons).
  raise: { row: ConditionRow; resource: string }[];
  // Previously notified conditions that are healthy again (send a quiet update).
  recover: ConditionRow[];
  upsert: ConditionRow[];
  remove: string[];
};

// Healthy must be seen this long after the condition began before we call it recovered,
// so a failure and the same minute's healthy snapshot do not cancel each other out.
export const RECOVERY_MIN_AGE_MS = 2 * 60 * 1000;

const ready = (row: ConditionRow, now: number) =>
  row.reason === 'out-of-sync' ? now - row.since >= OUT_OF_SYNC_GRACE_MS : true;

// Pure planner: given what is wrong now, what is known healthy and what we stored,
// decide what to notify, what to mark recovered and what to persist.
export function planConditions(input: {
  now: number;
  rows: ConditionRow[];
  active: ActiveCondition[];
  healthy: string[];
}): ConditionPlan {
  const { now } = input;
  const rows = new Map(input.rows.map(row => [row.scope, row]));
  const plan: ConditionPlan = { raise: [], recover: [], upsert: [], remove: [] };
  for (const item of input.active) {
    const stored = rows.get(item.scope);
    // A reason change keeps the clock and notified flag: one incident, one notification.
    const row: ConditionRow = stored
      ? { ...stored, reason: item.reason }
      : { scope: item.scope, reason: item.reason, since: now, notified: false };
    if (!row.notified && ready(row, now)) {
      row.notified = true;
      plan.raise.push({ row, resource: item.resource });
    }
    if (!stored || stored.reason !== row.reason || stored.notified !== row.notified) {
      plan.upsert.push(row);
    }
  }
  for (const scope of input.healthy) {
    const stored = rows.get(scope);
    if (!stored || now - stored.since < RECOVERY_MIN_AGE_MS) continue;
    if (stored.notified) plan.recover.push(stored);
    plan.remove.push(scope);
  }
  return plan;
}
