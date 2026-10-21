import { createBackendModule, coreServices } from '@backstage/backend-plugin-api';
import { notificationService } from '@backstage/plugin-notifications-node';
import { eventsServiceRef } from '@backstage/plugin-events-node';
import { signalsServiceRef } from '@backstage/plugin-signals-node';
import express from 'express';
import {
  ACTIVITY_EVENT_TOPIC,
  ALERT_SIGNAL_CHANNEL,
  isActivityKind,
  list,
  type ActivityEvent,
} from '@internal/plugin-onemind-common-backend';
import { ADMIN_GROUP_REF, createAdminRecipients } from './recipients';
import {
  classifyApplications,
  classifyBackups,
  ensureDeliveryTables,
  processConditions,
  processObservations,
  runDigest,
  type AppView,
  type Observation,
  type Sinks,
} from './scan';
import { createStore } from './store';
import { createTimelineRouter } from './timelineRoute';

type Application = {
  metadata: { name: string; uid: string };
  status?: {
    health?: { status?: string };
    sync?: { status?: string };
    operationState?: { phase?: string; finishedAt?: string };
  };
};
type Job = {
  metadata: { name: string; uid: string; ownerReferences?: { kind: string; name: string }[] };
  status?: { conditions?: { type: string; status: string; lastTransitionTime?: string }[] };
};
type CronJob = { metadata: { name: string }; spec?: { suspend?: boolean } };

const argoLink = (name: string) =>
  `https://argocd.tail717ba0.ts.net/applications/argocd/${encodeURIComponent(name)}`;
const RUNBOOK = '/docs/default/component/onemind-runbooks/resilience/';
// Superseded by backup-postgres-backstage; kept suspended on purpose, so not watched.
const DEFAULT_BACKUP_CRONJOBS = ['backup-postgres-backstage'];

// Internal observer reuses existing read-only cluster permissions. No external webhook
// credential, Kubernetes writes, arbitrary recipients or raw provider payloads.
export const operationalEventsModule = createBackendModule({
  pluginId: 'onemind',
  moduleId: 'operational-events',
  register(reg) {
    reg.registerInit({
      deps: {
        config: coreServices.rootConfig,
        scheduler: coreServices.scheduler,
        database: coreServices.database,
        logger: coreServices.logger,
        notifications: notificationService,
        signals: signalsServiceRef,
        events: eventsServiceRef,
        auth: coreServices.auth,
        discovery: coreServices.discovery,
        router: coreServices.httpRouter,
        httpAuth: coreServices.httpAuth,
        userInfo: coreServices.userInfo,
      },
      async init({ config, scheduler, database, logger, notifications, signals, events, auth, discovery, router, httpAuth, userInfo }) {
        const enabled = config.getOptionalBoolean('onemind.operationalEvents.enabled') ?? false;
        const state = {
          enabled,
          recipientReadiness: 'not-verified',
          lastAttemptAt: null as string | null,
          lastSuccessAt: null as string | null,
          sources: {} as Record<string, string>,
        };
        let readState = async () => state;
        const routes = express.Router();
        routes.get('/event-sources', async (req, res, next) => {
          try {
            await httpAuth.credentials(req, { allow: ['user'] });
            res.json(await readState());
          } catch (e) {
            next(e);
          }
        });
        router.use(routes);
        if (!enabled) return;

        const db = await database.getClient();
        await ensureDeliveryTables(db);
        const store = await createStore(db);
        routes.use(createTimelineRouter({ httpAuth, userInfo, store }));
        readState = async () => {
          const saved = await db('onemind_operation_observer').where({ id: 'status' }).first();
          return saved ? { ...JSON.parse(saved.payload), enabled } : state;
        };
        const recipients = createAdminRecipients({ auth, discovery });
        const backupCronJobs =
          config.getOptionalStringArray('onemind.operationalEvents.backupCronJobs') ?? DEFAULT_BACKUP_CRONJOBS;
        const digestHourUtc = config.getOptionalNumber('onemind.operationalEvents.digestHourUtc') ?? 13;

        // External producers (the activity bridge) post successes here; they are recorded, never notified.
        await events.subscribe({
          id: 'onemind-activity-timeline',
          topics: [ACTIVITY_EVENT_TOPIC],
          onEvent: async ({ eventPayload }) => {
            const event = eventPayload as Partial<ActivityEvent>;
            if (!isActivityKind(event?.kind) || typeof event.id !== 'string' || typeof event.resource !== 'string' || typeof event.occurredAt !== 'string') return;
            await store.timeline.record({ id: `ext:${event.id}`.slice(0, 64), kind: event.kind, resource: event.resource, occurredAt: event.occurredAt, source: 'external' });
          },
        });

        const makeSinks = (userRefs: string[]): Sinks => ({
          notify: async payload => {
            await notifications.send({ recipients: { type: 'entity', entityRef: ADMIN_GROUP_REF }, payload });
          },
          signal: async message => {
            await signals.publish({ recipients: { type: 'user', entityRef: userRefs }, channel: ALERT_SIGNAL_CHANNEL, message });
          },
        });

        const observeApplications = async () => {
          const apps = await list<Application>('/apis/argoproj.io/v1alpha1/namespaces/argocd/applications');
          const observations: Observation[] = apps.flatMap(app => {
            const operation = app.status?.operationState;
            if (!operation?.finishedAt || !['Succeeded', 'Failed', 'Error'].includes(operation.phase ?? '') || !Number.isFinite(Date.parse(operation.finishedAt))) return [];
            return [{ key: `${app.metadata.uid}:${operation.finishedAt}:${operation.phase}`, resource: app.metadata.name, kind: operation.phase === 'Succeeded' ? 'deployment.succeeded' : 'deployment.failed', occurredAt: operation.finishedAt, link: argoLink(app.metadata.name) } as Observation];
          });
          const views: AppView[] = apps.map(app => ({ name: app.metadata.name, health: app.status?.health?.status, sync: app.status?.sync?.status, link: argoLink(app.metadata.name) }));
          return { observations, views };
        };
        const observeBackups = async () => {
          const jobs = await list<Job>('/apis/batch/v1/namespaces/backstage/jobs');
          const cronJobs = await list<CronJob>('/apis/batch/v1/namespaces/backstage/cronjobs');
          const observations: Observation[] = jobs.flatMap(job => {
            const owner = job.metadata.ownerReferences?.find(item => item.kind === 'CronJob' && backupCronJobs.includes(item.name));
            if (!owner) return [];
            const condition = job.status?.conditions?.find(item => ['Complete', 'Failed'].includes(item.type) && item.status === 'True');
            if (!condition?.lastTransitionTime || !Number.isFinite(Date.parse(condition.lastTransitionTime))) return [];
            return [{ key: `${job.metadata.uid}:${condition.type}`, resource: `backstage/${owner.name}`, kind: condition.type === 'Complete' ? 'backup.succeeded' : 'backup.failed', occurredAt: condition.lastTransitionTime, link: RUNBOOK } as Observation];
          });
          const watched = cronJobs
            .filter(item => backupCronJobs.includes(item.metadata.name))
            .map(item => ({ name: item.metadata.name, suspended: item.spec?.suspend === true, link: RUNBOOK }));
          return { observations, watched };
        };

        await scheduler.scheduleTask({
          id: 'onemind-operational-events',
          scope: 'global',
          frequency: { minutes: 1 },
          timeout: { seconds: 50 },
          initialDelay: { seconds: 45 },
          fn: async () => {
            state.lastAttemptAt = new Date().toISOString();
            try {
              const userRefs = await recipients.userRefs();
              if (!userRefs) {
                state.recipientReadiness = 'not-ready';
                logger.warn('Operational events paused: intended catalog recipients are not ready');
                return;
              }
              state.recipientReadiness = 'catalog-membership-verified';
              const sinks = makeSinks(userRefs);
              try {
                const { observations, views } = await observeApplications();
                const result = await processObservations({ db, store, source: 'argocd', values: observations, sinks });
                const { active, healthy } = classifyApplications(views);
                await processConditions({ store, active, healthy, topic: 'deployments', sinks });
                state.sources.argocd = result;
              } catch {
                state.sources.argocd = 'unavailable';
                logger.warn('Operational event source unavailable: argocd');
              }
              try {
                const { observations, watched } = await observeBackups();
                const result = await processObservations({ db, store, source: 'backstage-backup', values: observations, sinks });
                const { active, healthy } = classifyBackups(watched);
                await processConditions({ store, active, healthy, topic: 'backups', sinks });
                state.sources['backstage-backup'] = result;
              } catch {
                state.sources['backstage-backup'] = 'unavailable';
                logger.warn('Operational event source unavailable: backstage-backup');
              }
              if (Object.values(state.sources).every(value => value !== 'unavailable')) state.lastSuccessAt = new Date().toISOString();
            } catch {
              state.recipientReadiness = 'unavailable';
              logger.warn('Operational event recipient verification unavailable');
            } finally {
              await db('onemind_operation_observer').insert({ id: 'status', payload: JSON.stringify(state) }).onConflict('id').merge();
            }
          },
        });

        // Checked every 30 minutes and idempotent per UTC day, so downtime at the digest hour only delays it.
        await scheduler.scheduleTask({
          id: 'onemind-activity-digest',
          scope: 'global',
          frequency: { minutes: 30 },
          timeout: { seconds: 50 },
          initialDelay: { seconds: 90 },
          fn: async () => {
            try {
              const userRefs = await recipients.userRefs();
              if (!userRefs) return;
              await runDigest({ store, sinks: makeSinks(userRefs), hourUtc: digestHourUtc });
              await store.timeline.prune(new Date(Date.now() - 30 * 24 * 3600_000).toISOString());
              await store.alerts.pruneResolved(new Date(Date.now() - 7 * 24 * 3600_000).toISOString());
            } catch {
              logger.warn('Daily activity digest unavailable; will retry');
            }
          },
        });
      },
    });
  },
});
