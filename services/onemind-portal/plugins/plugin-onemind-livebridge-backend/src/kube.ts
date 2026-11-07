import {
  CoreV1Api,
  CustomObjectsApi,
  KubeConfig,
  Watch,
} from '@kubernetes/client-node';
import type { LoggerService } from '@backstage/backend-plugin-api';
import type { SourceStatus } from './contract';
import type { ClusterState } from './cluster';
import { RateMeter } from './rate';

interface LoopState {
  everConnected: boolean;
  ok: boolean;
  error?: string;
}

interface LoopSpec {
  name: string;
  state: LoopState;
  list(): Promise<{ items: any[]; resourceVersion?: string }>;
  path: string;
  query: Record<string, string>;
  onList(items: any[]): void;
  onEvent(phase: string, obj: any): void;
}

const MIN_BACKOFF_MS = 1000;
const MAX_BACKOFF_MS = 60_000;

/** Load in-cluster config, falling back to the default kubeconfig (local development). */
export function loadKubeConfig(): KubeConfig | undefined {
  const kc = new KubeConfig();
  try {
    kc.loadFromCluster();
    return kc;
  } catch {
    // not in a pod
  }
  try {
    kc.loadFromDefault();
    return kc.getCurrentCluster() ? kc : undefined;
  } catch {
    return undefined;
  }
}

/**
 * Read-only list+watch of pods, Warning events and Argo CD applications. Every loop re-lists and
 * re-watches when the watch ends or fails (watches end every few minutes; a 410 Gone arrives as an
 * ERROR event and also triggers a re-list), with exponential backoff on errors. Only get, list
 * and watch verbs are used.
 */
export class KubeWatcher {
  private stopped = false;
  private readonly controllers = new Set<AbortController>();
  private readonly timers = new Set<NodeJS.Timeout>();
  private readonly pods: LoopState = { everConnected: false, ok: false };
  private readonly events: LoopState = { everConnected: false, ok: false };
  private readonly applications: LoopState = {
    everConnected: false,
    ok: false,
  };
  private lastKubeAt?: number;
  private lastArgoAt?: number;
  private readonly kubeRate = new RateMeter();
  private readonly argoRate = new RateMeter();

  constructor(
    private readonly kc: KubeConfig,
    private readonly state: ClusterState,
    private readonly logger: LoggerService,
    private readonly argoNamespace: string,
    private readonly now: () => number = Date.now,
  ) {}

  start(): void {
    const core = this.kc.makeApiClient(CoreV1Api);
    const custom = this.kc.makeApiClient(CustomObjectsApi);
    const kube = (obj: any) => {
      this.lastKubeAt = this.now();
      this.kubeRate.hit(this.lastKubeAt);
      return obj;
    };
    void this.loop({
      name: 'pods',
      state: this.pods,
      list: async () => {
        const res = await core.listPodForAllNamespaces();
        return {
          items: res.items,
          resourceVersion: res.metadata?.resourceVersion,
        };
      },
      path: '/api/v1/pods',
      query: {},
      onList: items => {
        kube(items);
        this.state.replacePods(items, this.now());
      },
      onEvent: (phase, obj) =>
        this.state.applyPod(phase, kube(obj), this.now()),
    });
    void this.loop({
      name: 'events',
      state: this.events,
      list: async () => {
        const res = await core.listEventForAllNamespaces({
          fieldSelector: 'type=Warning',
        });
        return {
          items: res.items,
          resourceVersion: res.metadata?.resourceVersion,
        };
      },
      path: '/api/v1/events',
      query: { fieldSelector: 'type=Warning' },
      onList: items => {
        kube(items);
        this.state.replaceWarnings(items);
      },
      onEvent: (phase, obj) => this.state.applyWarning(phase, kube(obj)),
    });
    const argo = (obj: any) => {
      this.lastArgoAt = this.now();
      this.argoRate.hit(this.lastArgoAt);
      return obj;
    };
    void this.loop({
      name: 'applications',
      state: this.applications,
      list: async () => {
        const res: any = await custom.listNamespacedCustomObject({
          group: 'argoproj.io',
          version: 'v1alpha1',
          namespace: this.argoNamespace,
          plural: 'applications',
        });
        return {
          items: res.items ?? [],
          resourceVersion: res.metadata?.resourceVersion,
        };
      },
      path: `/apis/argoproj.io/v1alpha1/namespaces/${this.argoNamespace}/applications`,
      query: {},
      onList: items => {
        argo(items);
        this.state.replaceArgo(items, this.now());
      },
      onEvent: (phase, obj) =>
        this.state.applyArgo(phase, argo(obj), this.now()),
    });
  }

  stop(): void {
    this.stopped = true;
    for (const c of this.controllers) c.abort();
    for (const t of this.timers) clearTimeout(t);
    this.timers.clear();
  }

  status(): { kubernetes: SourceStatus; argo: SourceStatus } {
    const at = this.now();
    const make = (
      states: Array<[string, LoopState]>,
      last: number | undefined,
      meter: RateMeter,
    ): SourceStatus => {
      const failing = states.filter(([, s]) => !s.ok && s.everConnected);
      const notYet = states.filter(([, s]) => !s.everConnected);
      let state: SourceStatus['state'] = 'connected';
      let detail: string | undefined;
      if (failing.length) {
        state = 'down';
        detail = failing.map(([n, s]) => `${n}: ${s.error}`).join('; ');
      } else if (notYet.length) {
        state = 'connecting';
        detail =
          notYet
            .map(([n, s]) => (s.error ? `${n}: ${s.error}` : undefined))
            .filter(Boolean)
            .join('; ') || undefined;
        if (notYet.some(([, s]) => s.error)) state = 'down';
      }
      return {
        state,
        lastMessageAt: last ? new Date(last).toISOString() : undefined,
        rate: meter.rate(at),
        detail,
      };
    };
    return {
      kubernetes: make(
        [
          ['pods', this.pods],
          ['events', this.events],
        ],
        this.lastKubeAt,
        this.kubeRate,
      ),
      argo: make(
        [['applications', this.applications]],
        this.lastArgoAt,
        this.argoRate,
      ),
    };
  }

  private sleep(ms: number): Promise<void> {
    return new Promise(resolve => {
      const t = setTimeout(() => {
        this.timers.delete(t);
        resolve();
      }, ms);
      t.unref?.();
      this.timers.add(t);
    });
  }

  private async loop(spec: LoopSpec): Promise<void> {
    let backoff = MIN_BACKOFF_MS;
    const watch = new Watch(this.kc);
    while (!this.stopped) {
      let watchError: string | undefined;
      const startedAt = this.now();
      try {
        const { items, resourceVersion } = await spec.list();
        spec.onList(items);
        spec.state.ok = true;
        spec.state.everConnected = true;
        spec.state.error = undefined;
        await new Promise<void>(resolve => {
          let finished = false;
          let controller: AbortController | undefined;
          const end = (error?: string) => {
            if (finished) return;
            finished = true;
            if (error) watchError = error;
            controller?.abort();
            if (controller) this.controllers.delete(controller);
            resolve();
          };
          watch
            .watch(
              spec.path,
              {
                ...spec.query,
                resourceVersion,
                allowWatchBookmarks: true,
                timeoutSeconds: 240 + Math.floor(Math.random() * 60),
              },
              (phase, obj) => {
                if (phase === 'ERROR') {
                  // 410 Gone (resourceVersion too old) and other watch errors: re-list.
                  end(
                    obj?.code === 410
                      ? undefined
                      : `watch error ${obj?.code ?? ''} ${
                          obj?.reason ?? ''
                        }`.trim(),
                  );
                  return;
                }
                if (phase === 'BOOKMARK') return;
                try {
                  spec.onEvent(phase, obj);
                } catch (e) {
                  this.logger.warn(
                    `livebridge ${spec.name} event failed: ${
                      (e as Error).message
                    }`,
                  );
                }
              },
              err =>
                end(
                  err && err.name !== 'AbortError'
                    ? String(err.message ?? err)
                    : undefined,
                ),
            )
            .then(c => {
              controller = c;
              this.controllers.add(c);
              if (this.stopped || finished) end();
            })
            .catch(e => end(String(e?.message ?? e)));
        });
      } catch (e) {
        watchError = String((e as Error)?.message ?? e);
        spec.state.ok = false;
      }
      if (this.stopped) return;
      if (watchError) {
        spec.state.ok = false;
        spec.state.error = watchError.slice(0, 200);
        this.logger.warn(
          `livebridge ${spec.name} watch: ${spec.state.error}; retrying in ${backoff} ms`,
        );
        await this.sleep(backoff + Math.floor(Math.random() * 250));
        backoff =
          this.now() - startedAt > 30_000
            ? MIN_BACKOFF_MS
            : Math.min(backoff * 2, MAX_BACKOFF_MS);
      } else {
        backoff = MIN_BACKOFF_MS;
        await this.sleep(250 + Math.floor(Math.random() * 250));
      }
    }
  }
}
