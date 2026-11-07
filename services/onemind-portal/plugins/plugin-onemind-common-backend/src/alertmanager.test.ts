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

describe('Falco events from falcosidekick', () => {
  const falcoLabels = (over: Record<string, string> = {}) => ({
    source: 'falco',
    rule: 'Drop and execute new binary in container',
    priority: 'Critical',
    severity: 'critical',
    k8s_ns_name: 'github-runner',
    k8s_pod_name: 'github-runner-5f5c5d6bfb-sqrl9',
    container_name: 'runner',
    container_image_repository: 'docker.io/myoung34/github-runner',
    proc_name: 'native-context',
    user_name: 'root',
    evt_time: '1794067956465324958',
    ...over,
  });
  const falco = (labels: Record<string, string>, fingerprint: string, status = 'firing') => {
    const parsed = parseAlertmanagerWebhook(
      payload(alert({ status, labels, annotations: {}, fingerprint })),
    );
    if (!parsed.ok) throw new Error(parsed.error);
    return alertNotification(parsed.alerts[0]);
  };

  test('an event without alertname becomes a readable alert', () => {
    const result = falco(falcoLabels(), 'aaaaaaaaaaaaaaa1');
    expect(result?.severity).toBe('critical');
    expect(result?.payload.title).toBe(
      'Alert: Falco: Drop and execute new binary in container (github-runner)',
    );
    expect(result?.payload.description).toBe(
      'Drop and execute new binary in container, in github-runner/github-runner-5f5c5d6bfb-sqrl9, process native-context, user root, image docker.io/myoung34/github-runner',
    );
  });

  test('repeats of one rule on one workload share a fingerprint and scope', () => {
    const first = falco(falcoLabels(), 'aaaaaaaaaaaaaaa1');
    const again = falco(
      falcoLabels({ k8s_pod_name: 'github-runner-other', evt_time: '1794067999999999999' }),
      'bbbbbbbbbbbbbbb2',
    );
    expect(again?.fingerprint).toBe(first?.fingerprint);
    expect(again?.payload.scope).toBe(first?.payload.scope);
    const otherRule = falco(falcoLabels({ rule: 'Terminal shell in container' }), 'ccccccccccccccc3');
    expect(otherRule?.fingerprint).not.toBe(first?.fingerprint);
    const otherNamespace = falco(falcoLabels({ k8s_ns_name: 'hermes' }), 'ddddddddddddddd4');
    expect(otherNamespace?.fingerprint).not.toBe(first?.fingerprint);
  });

  test('resolved Falco events keep the same scope', () => {
    const firing = falco(falcoLabels(), 'aaaaaaaaaaaaaaa1');
    const resolved = falco(falcoLabels(), 'eeeeeeeeeeeeeee5', 'resolved');
    expect(resolved?.payload.scope).toBe(firing?.payload.scope);
    expect(resolved?.severity).toBe('low');
  });

  test('labels without source=falco and a rule are still ignored without alertname', () => {
    expect(falco(falcoLabels({ source: 'k8saudit', rule: '' }), 'fffffffffffffff6')).toBeUndefined();
  });
});
