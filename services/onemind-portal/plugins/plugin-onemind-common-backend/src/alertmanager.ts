import { createHash } from 'node:crypto';
import {
  alertScope,
  alertSeverity,
  type NotificationSeverity,
} from './notificationPolicy';

// Alertmanager webhook payload, version 4. Only the fields we use.
export type AlertmanagerAlert = {
  status: 'firing' | 'resolved';
  labels: Record<string, string>;
  annotations: Record<string, string>;
  startsAt: string;
  endsAt?: string;
  fingerprint: string;
};

export type ParsedWebhook =
  | { ok: true; alerts: AlertmanagerAlert[] }
  | { ok: false; error: string };

export const MAX_WEBHOOK_BYTES = 512 * 1024;
const MAX_ALERTS = 100;
const MAX_ENTRIES = 64;

const isRecord = (value: unknown): value is Record<string, unknown> =>
  typeof value === 'object' && value !== null && !Array.isArray(value);

function stringMap(value: unknown): Record<string, string> | undefined {
  if (value === undefined) return {};
  if (!isRecord(value)) return undefined;
  const entries = Object.entries(value);
  if (entries.length > MAX_ENTRIES) return undefined;
  const result: Record<string, string> = {};
  for (const [key, item] of entries) {
    if (typeof item !== 'string' || key.length > 128) return undefined;
    result[key] = item.slice(0, 2000);
  }
  return result;
}

export function parseAlertmanagerWebhook(body: unknown): ParsedWebhook {
  if (!isRecord(body) || !Array.isArray(body.alerts)) {
    return { ok: false, error: 'Not an Alertmanager webhook payload' };
  }
  if (body.version !== undefined && body.version !== '4') {
    return { ok: false, error: 'Unsupported Alertmanager payload version' };
  }
  if (body.alerts.length > MAX_ALERTS) {
    return { ok: false, error: 'Too many alerts in one payload' };
  }
  const alerts: AlertmanagerAlert[] = [];
  for (const raw of body.alerts) {
    if (!isRecord(raw) || (raw.status !== 'firing' && raw.status !== 'resolved')) {
      return { ok: false, error: 'Invalid alert status' };
    }
    const labels = stringMap(raw.labels);
    const annotations = stringMap(raw.annotations);
    if (!labels || !annotations || typeof raw.startsAt !== 'string') {
      return { ok: false, error: 'Invalid alert labels or annotations' };
    }
    if (raw.endsAt !== undefined && typeof raw.endsAt !== 'string') {
      return { ok: false, error: 'Invalid alert end time' };
    }
    // Older Alertmanagers omit the fingerprint; derive a stable one from the labels.
    const fingerprint =
      typeof raw.fingerprint === 'string' && /^[a-zA-Z0-9]{1,64}$/.test(raw.fingerprint)
        ? raw.fingerprint
        : createHash('sha256')
            .update(JSON.stringify(Object.entries(labels).sort()))
            .digest('hex')
            .slice(0, 16);
    alerts.push({
      status: raw.status,
      labels,
      annotations,
      startsAt: raw.startsAt,
      endsAt: raw.endsAt as string | undefined,
      fingerprint,
    });
  }
  return { ok: true, alerts };
}

export type AlertNotification = {
  fingerprint: string;
  status: 'firing' | 'resolved';
  severity: NotificationSeverity;
  // Severity of the item as it was while firing; resolved updates are quiet (low).
  firingSeverity: NotificationSeverity;
  payload: {
    title: string;
    description: string;
    severity: NotificationSeverity;
    topic: 'alerts';
    scope: string;
    link: string;
  };
};

// eslint-disable-next-line no-control-regex
const control = /[\u0000-\u001f\u007f]/g;
const clean = (value: string, max: number) =>
  value.replace(control, ' ').trim().slice(0, max);

const targetLabels = ['deployment', 'statefulset', 'daemonset', 'pod', 'job', 'service', 'instance', 'namespace'];

// Only https runbooks are linked; the inbox otherwise falls back to the Today page.
function alertLink(alert: AlertmanagerAlert): string {
  const runbook = alert.annotations.runbook_url;
  if (runbook && /^https:\/\/[^\s]{1,400}$/.test(runbook)) return runbook;
  return '/today';
}

// Returns undefined for alerts that are not meant for a human (severity none, Watchdog).
export function alertNotification(
  alert: AlertmanagerAlert,
): AlertNotification | undefined {
  const firingSeverity = alertSeverity(alert.labels.severity);
  const name = alert.labels.alertname;
  if (!firingSeverity || !name) return undefined;
  const label = targetLabels.find(item => alert.labels[item]);
  const target = label ? ` (${clean(alert.labels[label], 80)})` : '';
  const heading = `${clean(name, 80)}${target}`;
  const resolved = alert.status === 'resolved';
  const text = alert.annotations.summary ?? alert.annotations.description ?? '';
  return {
    fingerprint: alert.fingerprint,
    status: alert.status,
    severity: resolved ? 'low' : firingSeverity,
    firingSeverity,
    payload: {
      title: resolved ? `Resolved: ${heading}` : `Alert: ${heading}`,
      description: resolved
        ? `Alertmanager reports this alert resolved${alert.endsAt ? ` at ${clean(alert.endsAt, 40)}` : ''}.`
        : clean(text, 500),
      severity: resolved ? 'low' : firingSeverity,
      topic: 'alerts',
      scope: alertScope(alert.fingerprint),
      link: alertLink(alert),
    },
  };
}
