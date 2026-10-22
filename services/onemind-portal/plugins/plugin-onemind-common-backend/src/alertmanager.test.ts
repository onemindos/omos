import { alertNotification, parseAlertmanagerWebhook } from './alertmanager';

const alert = (over: Record<string, unknown> = {}) => ({
  status: 'firing',
  labels: { alertname: 'KubePodCrashLooping', severity: 'critical', namespace: 'backstage', pod: 'backstage-1' },
  annotations: { summary: 'Pod is crash looping' },
  startsAt: '2026-10-21T10:00:00Z',
  endsAt: '0001-01-01T00:00:00Z',
  fingerprint: 'abc123def4567890',
  ...over,
});
const payload = (...alerts: unknown[]) => ({ version: '4', status: 'firing', alerts });

describe('Alertmanager webhook parsing', () => {
  test('accepts a version 4 payload', () => {
    const parsed = parseAlertmanagerWebhook(payload(alert()));
    expect(parsed.ok && parsed.alerts).toHaveLength(1);
  });

  test.each([
    ['not an object', 'text'],
    ['no alerts array', { version: '4' }],
    ['wrong version', { version: '3', alerts: [] }],
    ['bad status', payload(alert({ status: 'pending' }))],
    ['non-string label', payload(alert({ labels: { alertname: 5 } }))],
    ['array labels', payload(alert({ labels: ['x'] }))],
    ['missing start', payload(alert({ startsAt: undefined }))],
  ])('rejects %s', (_name, body) => {
    expect(parseAlertmanagerWebhook(body).ok).toBe(false);
  });

  test('rejects oversized batches', () => {
    expect(
      parseAlertmanagerWebhook(payload(...Array.from({ length: 101 }, () => alert()))).ok,
    ).toBe(false);
  });

  test('derives a stable fingerprint when Alertmanager omits it', () => {
    const a = parseAlertmanagerWebhook(payload(alert({ fingerprint: undefined })));
    const b = parseAlertmanagerWebhook(payload(alert({ fingerprint: undefined })));
    expect(a.ok && b.ok && a.alerts[0].fingerprint).toBe(b.ok && b.alerts[0].fingerprint);
    expect(a.ok && a.alerts[0].fingerprint).toMatch(/^[0-9a-f]{16}$/);
  });
});

describe('alert to notification mapping', () => {
  const parse = (over?: Record<string, unknown>) => {
    const parsed = parseAlertmanagerWebhook(payload(alert(over)));
    if (!parsed.ok) throw new Error(parsed.error);
    return parsed.alerts[0];
  };

  test('critical firing alert is critical with a stable per-fingerprint scope', () => {
    const result = alertNotification(parse());
    expect(result).toMatchObject({
      status: 'firing',
      severity: 'critical',
      payload: { severity: 'critical', topic: 'alerts', scope: 'alert:abc123def4567890', link: '/today' },
    });
    expect(result?.payload.title).toBe('Alert: KubePodCrashLooping (backstage-1)');
  });

  test('warning maps to high and info to normal', () => {
    expect(alertNotification(parse({ labels: { alertname: 'A', severity: 'warning' } }))?.severity).toBe('high');
    expect(alertNotification(parse({ labels: { alertname: 'A', severity: 'info' } }))?.severity).toBe('normal');
  });

  test('Watchdog and unlabelled alerts are ignored', () => {
    expect(alertNotification(parse({ labels: { alertname: 'Watchdog', severity: 'none' } }))).toBeUndefined();
    expect(alertNotification(parse({ labels: { alertname: 'X' } }))).toBeUndefined();
  });

  test('resolved update keeps the scope and goes quiet', () => {
    const firing = alertNotification(parse());
    const resolved = alertNotification(parse({ status: 'resolved', endsAt: '2026-10-21T10:30:00Z' }));
    expect(resolved?.payload.scope).toBe(firing?.payload.scope);
    expect(resolved).toMatchObject({ status: 'resolved', severity: 'low', firingSeverity: 'critical' });
    expect(resolved?.payload.title).toBe('Resolved: KubePodCrashLooping (backstage-1)');
  });

  test('only https runbook links are used', () => {
    expect(alertNotification(parse({ annotations: { runbook_url: 'https://runbooks.example/x' } }))?.payload.link).toBe('https://runbooks.example/x');
    expect(alertNotification(parse({ annotations: { runbook_url: 'javascript:alert(1)' } }))?.payload.link).toBe('/today');
    expect(alertNotification(parse({ annotations: { runbook_url: 'http://insecure.example' } }))?.payload.link).toBe('/today');
  });

  test('control characters are stripped and text is bounded', () => {
    const result = alertNotification(parse({ annotations: { summary: `a\nb${'x'.repeat(900)}` } }));
    expect(result?.payload.description).not.toMatch(/[\n\r]/);
    expect(result?.payload.description.length).toBeLessThanOrEqual(500);
  });
});
