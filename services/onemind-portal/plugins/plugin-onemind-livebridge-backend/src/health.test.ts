import { ClusterState, toPodInfo } from './cluster';
import { computeNamespaceCell, type PodInfo } from './health';

const NOW = Date.parse('2026-11-07T12:00:00Z');
const pod = (over: Partial<PodInfo> = {}): PodInfo => ({
  namespace: 'ns',
  name: 'p',
  phase: 'Running',
  ready: true,
  deleting: false,
  restartCount: 0,
  recentTerminations: 0,
  observedRestarts: [],
  ...over,
});

describe('computeNamespaceCell', () => {
  it('is ok when everything runs and is ready', () => {
    expect(
      computeNamespaceCell('ns', [pod(), pod({ name: 'q' })], NOW),
    ).toMatchObject({
      health: 'ok',
      total: 2,
      ready: 2,
      pending: 0,
      failed: 0,
      restarts: 0,
    });
  });

  it('is down for a failed pod', () => {
    expect(
      computeNamespaceCell(
        'ns',
        [pod(), pod({ phase: 'Failed', ready: false })],
        NOW,
      ),
    ).toMatchObject({ health: 'down', failed: 1 });
  });

  it.each(['CrashLoopBackOff', 'ImagePullBackOff', 'ErrImagePull'])(
    'is down for %s',
    reason => {
      expect(
        computeNamespaceCell(
          'ns',
          [pod({ ready: false, waitingReason: reason })],
          NOW,
        ).health,
      ).toBe('down');
    },
  );

  it('is warn for pending, not ready or a recent restart', () => {
    expect(
      computeNamespaceCell(
        'ns',
        [pod({ phase: 'Pending', ready: false })],
        NOW,
      ),
    ).toMatchObject({ health: 'warn', pending: 1 });
    expect(
      computeNamespaceCell('ns', [pod({ ready: false })], NOW).health,
    ).toBe('warn');
    expect(
      computeNamespaceCell('ns', [pod({ recentTerminations: 1 })], NOW),
    ).toMatchObject({ health: 'warn', restarts: 1 });
    expect(
      computeNamespaceCell(
        'ns',
        [pod({ observedRestarts: [NOW - 60_000, NOW - 120_000] })],
        NOW,
      ),
    ).toMatchObject({ health: 'warn', restarts: 2 });
  });

  it('ignores restarts seen more than an hour ago and completed pods', () => {
    const cell = computeNamespaceCell(
      'ns',
      [
        pod({ observedRestarts: [NOW - 3_700_000] }),
        pod({ name: 'job', phase: 'Succeeded', ready: false }),
      ],
      NOW,
    );
    expect(cell).toMatchObject({ health: 'ok', total: 1, restarts: 0 });
  });

  it('does not warn for a terminating pod that is not ready', () => {
    expect(
      computeNamespaceCell('ns', [pod({ ready: false, deleting: true })], NOW)
        .health,
    ).toBe('ok');
  });
});

const raw = (restartCount: number, extra: object = {}) => ({
  metadata: { namespace: 'forgejo', name: 'forgejo-0' },
  status: {
    phase: 'Running',
    conditions: [{ type: 'Ready', status: 'True' }],
    containerStatuses: [{ restartCount, state: { running: {} }, ...extra }],
  },
});

describe('ClusterState', () => {
  it('derives restarts from restartCount deltas and lastState.terminated', () => {
    const first = toPodInfo(raw(2), NOW);
    expect(first.recentTerminations).toBe(0);
    const next = toPodInfo(raw(4), NOW + 1000, first);
    expect(next.observedRestarts).toHaveLength(2);
    const recent = toPodInfo(
      raw(1, {
        lastState: { terminated: { finishedAt: '2026-11-07T11:30:00Z' } },
      }),
      NOW,
    );
    expect(recent.recentTerminations).toBe(1);
  });

  it('records the last change and sends only changed cells', () => {
    const s = new ClusterState();
    s.replacePods([raw(0)], NOW);
    const first = s.delta(NOW)!;
    expect(first.namespaces).toHaveLength(1);
    expect(s.delta(NOW + 1000)).toBeUndefined();
    s.applyPod(
      'MODIFIED',
      raw(1, { state: { waiting: { reason: 'CrashLoopBackOff' } } }),
      NOW + 2000,
    );
    const d = s.delta(NOW + 2000)!;
    expect(d.namespaces[0]).toMatchObject({
      namespace: 'forgejo',
      health: 'down',
      lastChange: 'pod forgejo-0 CrashLoopBackOff',
    });
    s.applyPod('DELETED', raw(1), NOW + 3000);
    expect(s.delta(NOW + 3000)!.removedNamespaces).toEqual(['forgejo']);
  });

  it('builds argo cells and warnings', () => {
    const s = new ClusterState();
    s.replaceArgo(
      [
        {
          metadata: { name: 'platform-nats' },
          status: {
            sync: { status: 'OutOfSync', revision: 'abcdef' },
            health: { status: 'Healthy' },
          },
        },
      ],
      NOW,
    );
    s.replaceWarnings([
      {
        metadata: { uid: '1', namespace: 'a' },
        involvedObject: { kind: 'Pod', name: 'x' },
        reason: 'BackOff',
        message: 'back-off',
        lastTimestamp: '2026-11-07T11:00:00Z',
        count: 3,
      },
      {
        metadata: { uid: '2', namespace: 'b' },
        involvedObject: { kind: 'Pod', name: 'y' },
        reason: 'Failed',
        message: 'm',
        lastTimestamp: '2026-11-07T11:59:00Z',
      },
    ]);
    expect(s.argoCells()[0]).toMatchObject({
      name: 'platform-nats',
      sync: 'OutOfSync',
      health: 'Healthy',
      revision: 'abcdef',
    });
    expect(s.warningList().map(w => w.reason)).toEqual(['Failed', 'BackOff']);
    const d = s.delta(NOW)!;
    expect(d.argo).toHaveLength(1);
    expect(d.warnings).toHaveLength(2);
    expect(s.delta(NOW)).toBeUndefined();
  });
});
