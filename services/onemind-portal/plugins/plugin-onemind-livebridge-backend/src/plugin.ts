import {
  createBackendPlugin,
  coreServices,
} from '@backstage/backend-plugin-api';
import { actionsRegistryServiceRef } from '@backstage/backend-plugin-api/alpha';
import { signalsServiceRef } from '@backstage/plugin-signals-node';
import { randomBytes } from 'node:crypto';
import { hostname } from 'node:os';
import { LiveBridge } from './bridge';
import { ClusterState } from './cluster';
import type { SourceStatus } from './contract';
import { KubeWatcher, loadKubeConfig } from './kube';
import { NatsSource } from './natsSource';
import { createSnapshotRouter } from './route';
import { summarizeSnapshot } from './summary';

const DEFAULT_NATS_URL = 'nats://nats.nats.svc.cluster.local:4222';

/**
 * Live bridge: read-only NATS and Kubernetes watchers, pushed to browsers over Signals
 * (broadcast channels onemind:live:*) and served as a snapshot at GET /api/livebridge/snapshot.
 */
export const livebridgePlugin = createBackendPlugin({
  pluginId: 'livebridge',
  register(env) {
    env.registerInit({
      deps: {
        logger: coreServices.logger,
        config: coreServices.rootConfig,
        lifecycle: coreServices.rootLifecycle,
        httpRouter: coreServices.httpRouter,
        httpAuth: coreServices.httpAuth,
        signals: signalsServiceRef,
        actions: actionsRegistryServiceRef,
      },
      async init({
        logger,
        config,
        lifecycle,
        httpRouter,
        httpAuth,
        signals,
        actions,
      }) {
        const origin = `${hostname()}-${randomBytes(3).toString('hex')}`;
        const base = config.getOptionalConfig('onemind.livebridge');

        const natsUser =
          base?.getOptionalString('nats.user') ?? process.env.NATS_USER;
        const natsPassword =
          base?.getOptionalString('nats.password') ?? process.env.NATS_PASSWORD;
        const cluster = new ClusterState();

        const ref: { bridge?: LiveBridge } = {};
        const nats = new NatsSource({
          url: base?.getOptionalString('nats.url') ?? DEFAULT_NATS_URL,
          user: natsUser || undefined,
          password: natsPassword || undefined,
          caFile: base?.getOptionalString('nats.caFile'),
          logger,
          onMessage: (subject, data) => ref.bridge?.ingest(subject, data),
        });

        const kubeEnabled =
          base?.getOptionalBoolean('kubernetes.enabled') ?? true;
        const kc = kubeEnabled ? loadKubeConfig() : undefined;
        const watcher = kc
          ? new KubeWatcher(
              kc,
              cluster,
              logger,
              base?.getOptionalString('kubernetes.argoNamespace') ?? 'argocd',
            )
          : undefined;
        const off = (detail: string): SourceStatus => ({
          state: 'disabled',
          detail,
        });

        const bridge = new LiveBridge({
          origin,
          logger,
          cluster,
          providers: {
            nats: () => nats.status(),
            kubernetes: () =>
              watcher?.status().kubernetes ??
              off(
                kubeEnabled
                  ? 'no Kubernetes configuration'
                  : 'disabled by config',
              ),
            argo: () =>
              watcher?.status().argo ??
              off(
                kubeEnabled
                  ? 'no Kubernetes configuration'
                  : 'disabled by config',
              ),
          },
          publish: (channel, message) =>
            signals.publish({
              recipients: { type: 'broadcast' },
              channel,
              message: message as never,
            }),
        });

        ref.bridge = bridge;
        httpRouter.use(
          createSnapshotRouter({ httpAuth, snapshot: () => bridge.snapshot() }),
        );

        actions.register({
          name: 'snapshot_summary',
          title: 'OneMind live bridge summary',
          description:
            'Digest of the live bridge: namespace health counts and unhealthy namespaces, Argo CD applications that are not Synced and Healthy, CoT track counts by affiliation, bus message rates and the last 20 bus events. Live observation by this portal replica; the bridge is read-only.',
          attributes: { readOnly: true, destructive: false, idempotent: true },
          schema: {
            input: z => z.object({}),
            output: z => z.object({ summary: z.unknown() }),
          },
          async action() {
            return {
              output: { summary: summarizeSnapshot(bridge.snapshot()) },
            };
          },
        });

        lifecycle.addStartupHook(() => {
          bridge.start();
          nats.start();
          watcher?.start();
          logger.info(
            `livebridge started as ${origin} (nats ${
              natsUser && natsPassword ? 'enabled' : 'disabled'
            }, kubernetes ${watcher ? 'enabled' : 'disabled'})`,
          );
        });
        lifecycle.addShutdownHook(async () => {
          bridge.stop();
          watcher?.stop();
          await nats.stop();
        });
      },
    });
  },
});
