import {
  activityKinds,
  activityPolicy,
  alertSeverity,
  argoScope,
  isActivityKind,
  isDeploymentSuccessNotification,
  isSignalSeverity,
  OUT_OF_SYNC_GRACE_MS,
  scopeFor,
} from './notificationPolicy';

describe('notification severity policy', () => {
  test('successes are timeline-only and never notify', () => {
    expect(activityPolicy('deployment.succeeded').delivery).toBe('timeline');
    expect(activityPolicy('backup.succeeded').delivery).toBe('timeline');
  });

  test.each([
    'deployment.failed',
    'deployment.degraded',
    'deployment.out-of-sync',
    'backup.failed',
    'backup.suspended',
  ] as const)('%s is a high inbox item', kind => {
    expect(activityPolicy(kind)).toMatchObject({
      delivery: 'notify',
      severity: 'high',
    });
  });

  test('every kind has a policy and only successes skip the inbox', () => {
    const skipped = activityKinds.filter(
      kind => activityPolicy(kind).delivery === 'timeline',
    );
    expect(skipped.sort()).toEqual(['backup.succeeded', 'deployment.succeeded']);
  });

  test('out-of-sync grace is fifteen minutes', () => {
    expect(OUT_OF_SYNC_GRACE_MS).toBe(15 * 60_000);
  });

  test('alert severities map critical, warning and info; others are ignored', () => {
    expect(alertSeverity('critical')).toBe('critical');
    expect(alertSeverity('Warning')).toBe('high');
    expect(alertSeverity('info')).toBe('normal');
    expect(alertSeverity('none')).toBeUndefined();
    expect(alertSeverity(undefined)).toBeUndefined();
  });

  test('only high and critical items push signals', () => {
    expect(isSignalSeverity('critical')).toBe(true);
    expect(isSignalSeverity('high')).toBe(true);
    expect(isSignalSeverity('normal')).toBe(false);
    expect(isSignalSeverity('low')).toBe(false);
  });

  test('scope is stable per source object, not per occurrence', () => {
    expect(scopeFor('deployment.failed', 'platform-backstage', 'op-1')).toBe(
      scopeFor('deployment.failed', 'platform-backstage', 'op-2'),
    );
    expect(scopeFor('deployment.failed', 'platform-backstage', 'op-1')).toBe(
      argoScope('platform-backstage'),
    );
    expect(scopeFor('backup.failed', 'backup-postgres', 'x')).toBe(
      'backup:backup-postgres',
    );
  });

  test('kind guard rejects prototype keys and unknown kinds', () => {
    expect(isActivityKind('deployment.failed')).toBe(true);
    expect(isActivityKind('toString')).toBe(false);
    expect(isActivityKind('deployment.exploded')).toBe(false);
  });

  test('backlog predicate matches only old low-severity success notifications', () => {
    const old = {
      title: 'Deployment completed: platform-backstage',
      topic: 'deployments',
      severity: 'low',
    };
    expect(isDeploymentSuccessNotification(old)).toBe(true);
    expect(
      isDeploymentSuccessNotification({ ...old, title: 'Deployment needs attention: x', severity: 'high' }),
    ).toBe(false);
    expect(isDeploymentSuccessNotification({ ...old, topic: 'backups' })).toBe(false);
    expect(isDeploymentSuccessNotification({ ...old, severity: 'high' })).toBe(false);
  });
});
