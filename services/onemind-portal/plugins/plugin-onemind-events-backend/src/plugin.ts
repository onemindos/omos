/* eslint-disable */
import { createBackendModule, coreServices } from '@backstage/backend-plugin-api';
import { notificationService } from '@backstage/plugin-notifications-node';
import express from 'express';
import { createHash } from 'node:crypto';
import { list } from '@internal/plugin-onemind-common-backend';

type Observation = { key: string; resource: string; kind: 'deployment.failed' | 'deployment.succeeded' | 'backup.failed' | 'backup.succeeded'; occurredAt: string; link: string };
type Application = { metadata: { name: string; uid: string }; status?: { operationState?: { phase?: string; finishedAt?: string } } };
type Job = { metadata: { name: string; uid: string; ownerReferences?: { kind: string; name: string }[] }; status?: { conditions?: { type: string; status: string; lastTransitionTime?: string }[] } };
const hash = (value: string) => createHash('sha256').update(value).digest('hex');

// Internal observer reuses existing read-only cluster permissions. No external webhook
// credential, Kubernetes writes, arbitrary recipients or raw provider payloads.
export default createBackendModule({ pluginId: 'onemind', moduleId: 'operational-events', register(reg) {
  reg.registerInit({ deps: { config: coreServices.rootConfig, scheduler: coreServices.scheduler, database: coreServices.database, logger: coreServices.logger, notifications: notificationService, auth: coreServices.auth, discovery: coreServices.discovery, router: coreServices.httpRouter, httpAuth: coreServices.httpAuth },
    async init({ config, scheduler, database, logger, notifications, auth, discovery, router, httpAuth }) {
      const enabled = config.getOptionalBoolean('onemind.operationalEvents.enabled') ?? false;
      const state = { enabled, recipientReadiness: 'not-verified', lastAttemptAt: null as string | null, lastSuccessAt: null as string | null, sources: {} as Record<string, string> };
      let readState = async () => state;
      const routes = express.Router();
      routes.get('/event-sources', async (req, res, next) => { try { await httpAuth.credentials(req, { allow: ['user'] }); res.json(await readState()); } catch (e) { next(e); } });
      router.use(routes);
      if (!enabled) return;
      const db = await database.getClient();
      // Global scheduler registration happens only after additive schema creation.
      if (!await db.schema.hasTable('onemind_operation_delivery')) await db.schema.createTable('onemind_operation_delivery', table => {
        table.string('id', 64).primary(); table.string('source', 32).notNullable(); table.string('recorded_at', 40).notNullable();
      }).catch(async error => { if (!await db.schema.hasTable('onemind_operation_delivery')) throw error; });
      if (!await db.schema.hasTable('onemind_operation_observer')) await db.schema.createTable('onemind_operation_observer', table => { table.string('id', 32).primary(); table.text('payload').notNullable(); }).catch(async error => { if (!await db.schema.hasTable('onemind_operation_observer')) throw error; });
      readState = async () => { const saved = await db('onemind_operation_observer').where({ id: 'status' }).first(); return saved ? { ...JSON.parse(saved.payload), enabled } : state; };
      const catalogEntity = async (kind: string, name: string) => {
        const { token } = await auth.getPluginRequestToken({ onBehalfOf: await auth.getOwnServiceCredentials(), targetPluginId: 'catalog' });
        const base = await discovery.getBaseUrl('catalog');
        const response = await fetch(`${base}/entities/by-name/${kind}/default/${encodeURIComponent(name)}`, { headers: { Authorization: `Bearer ${token}` }, signal: AbortSignal.timeout(8000), redirect: 'error' });
        if (!response.ok) throw new Error('Recipient catalog unavailable');
        return response.json();
      };
      const readyRecipients = async () => {
        const group = await catalogEntity('group', 'onemind-admins');
        const members: unknown = group.spec?.members;
        if (!Array.isArray(members) || members.length < 1 || members.length > 10) return false;
        for (const member of members) {
          if (typeof member !== 'string' || !/^[a-z0-9_-]+$/.test(member)) return false;
          const user = await catalogEntity('user', member);
          if (!user.relations?.some((relation: { type: string; targetRef: string }) => relation.type === 'memberOf' && relation.targetRef === 'group:default/onemind-admins')) return false;
        }
        return true;
      };
      const observeApplications = async (): Promise<Observation[]> => (await list<Application>('/apis/argoproj.io/v1alpha1/namespaces/argocd/applications')).flatMap(app => {
        const operation = app.status?.operationState;
        if (!operation?.finishedAt || !['Succeeded', 'Failed', 'Error'].includes(operation.phase ?? '') || !Number.isFinite(Date.parse(operation.finishedAt))) return [];
        return [{ key: `${app.metadata.uid}:${operation.finishedAt}:${operation.phase}`, resource: app.metadata.name, kind: operation.phase === 'Succeeded' ? 'deployment.succeeded' : 'deployment.failed', occurredAt: operation.finishedAt, link: `https://argocd.tail717ba0.ts.net/applications/argocd/${encodeURIComponent(app.metadata.name)}` }];
      });
      const observeBackups = async (): Promise<Observation[]> => (await list<Job>('/apis/batch/v1/namespaces/backstage/jobs')).flatMap(job => {
        if (!job.metadata.ownerReferences?.some(owner => owner.kind === 'CronJob' && owner.name === 'backstage-backup')) return [];
        const condition = job.status?.conditions?.find(item => ['Complete', 'Failed'].includes(item.type) && item.status === 'True');
        if (!condition?.lastTransitionTime || !Number.isFinite(Date.parse(condition.lastTransitionTime))) return [];
        return [{ key: `${job.metadata.uid}:${condition.type}`, resource: `backstage/${job.metadata.name}`, kind: condition.type === 'Complete' ? 'backup.succeeded' : 'backup.failed', occurredAt: condition.lastTransitionTime, link: '/docs/default/component/onemind-runbooks/resilience/' }];
      });
      await scheduler.scheduleTask({ id: 'onemind-operational-events', scope: 'global', frequency: { minutes: 1 }, timeout: { seconds: 50 }, initialDelay: { seconds: 45 }, fn: async () => {
        state.lastAttemptAt = new Date().toISOString();
        try {
          if (!await readyRecipients()) { state.recipientReadiness = 'not-ready'; logger.warn('Operational events paused: intended catalog recipients are not ready'); return; }
          state.recipientReadiness = 'catalog-membership-verified';
          for (const [source, observe] of [['argocd', observeApplications], ['backstage-backup', observeBackups]] as const) {
            try {
              const values = await observe();
              const baseline = await db('onemind_operation_delivery').where({ id: hash(`baseline:${source}`) }).first();
              if (!baseline) {
                await db.transaction(async tx => {
                  for (const value of values) await tx('onemind_operation_delivery').insert({ id: hash(`${source}:${value.key}`), source, recorded_at: new Date().toISOString() }).onConflict('id').ignore();
                  await tx('onemind_operation_delivery').insert({ id: hash(`baseline:${source}`), source, recorded_at: new Date().toISOString() }).onConflict('id').ignore();
                });
                state.sources[source] = 'baseline-established'; continue;
              }
              for (const value of values) {
                const id = hash(`${source}:${value.key}`);
                if (await db('onemind_operation_delivery').where({ id }).first()) continue;
                // Avoid replaying historical jobs newly discovered after retention changes.
                if (Date.parse(value.occurredAt) < Date.parse(baseline.recorded_at)) {
                  await db('onemind_operation_delivery').insert({ id, source, recorded_at: new Date().toISOString() }).onConflict('id').ignore(); continue;
                }
                const failed = value.kind.endsWith('.failed');
                await notifications.send({ recipients: { type: 'entity', entityRef: 'group:default/onemind-admins' }, payload: {
                  title: `${source === 'argocd' ? 'Deployment' : 'Backup job'} ${failed ? 'failed' : 'completed'}: ${value.resource}`,
                  description: `Observed ${value.kind} at ${value.occurredAt}. ${source === 'argocd' ? 'Operation completion is not proof of current application health.' : 'Job success is not separate proof of archive restoration; inspect the backup report.'}`,
                  severity: failed ? 'high' : 'low', topic: source === 'argocd' ? 'deployments' : 'backups', scope: `onemind-observer:${id}`, link: value.link,
                } });
                // send acknowledges processing, not user delivery. A crash before this
                // insert may retry; notification scope collapses duplicate occurrences.
                await db('onemind_operation_delivery').insert({ id, source, recorded_at: new Date().toISOString() }).onConflict('id').ignore();
              }
              state.sources[source] = 'observing';
            } catch { state.sources[source] = 'unavailable'; logger.warn(`Operational event source unavailable: ${source}`); }
          }
          if (Object.values(state.sources).every(value => value !== 'unavailable')) state.lastSuccessAt = new Date().toISOString();
        } catch { state.recipientReadiness = 'unavailable'; logger.warn('Operational event recipient verification unavailable'); }
        finally { await db('onemind_operation_observer').insert({ id: 'status', payload: JSON.stringify(state) }).onConflict('id').merge(); }
      } });
    },
  });
} });
