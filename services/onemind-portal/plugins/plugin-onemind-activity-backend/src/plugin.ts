/* eslint-disable */
import { createBackendPlugin, coreServices } from '@backstage/backend-plugin-api';
import { notificationService } from '@backstage/plugin-notifications-node';
import express from 'express';

// Producers cannot choose recipients, arbitrary links, severity, or message text.
const kinds = {
  'deployment.failed': { severity: 'high' as const, title: 'Deployment needs attention', topic: 'deployments' },
  'deployment.succeeded': { severity: 'low' as const, title: 'Deployment completed', topic: 'deployments' },
  'backup.failed': { severity: 'high' as const, title: 'Backup failed', topic: 'backups' },
  'backup.succeeded': { severity: 'low' as const, title: 'Backup completed', topic: 'backups' },
  'agent.completed': { severity: 'normal' as const, title: 'Agent work completed', topic: 'agents' },
};
export default createBackendPlugin({
  pluginId: 'onemind-events',
  register(reg) {
    reg.registerInit({ deps: { router: coreServices.httpRouter, httpAuth: coreServices.httpAuth, userInfo: coreServices.userInfo, notifications: notificationService },
      async init({ router, httpAuth, userInfo, notifications }) {
        const routes = express.Router();
        const windows = new Map<string, { time: number; count: number }>();
        routes.use(express.json({ limit: '8kb', strict: true }));
        routes.post('/events', async (req, res, next) => {
          try {
            const credentials = await httpAuth.credentials(req, { allow: ['user', 'service'] });
            let actor: string;
            if (credentials.principal.type === 'service') {
              if (credentials.principal.subject !== 'onemind-event-ingestor') { res.sendStatus(403); return; }
              actor = credentials.principal.subject;
            } else {
              const info = await userInfo.getUserInfo(credentials);
              if (!info.ownershipEntityRefs.includes('group:default/onemind-admins')) { res.sendStatus(403); return; }
              actor = info.userEntityRef;
            }
            const now = Date.now();
            for (const [key, value] of windows) if (now - value.time > 60000) windows.delete(key);
            const window = windows.get(actor) ?? { time: now, count: 0 };
            window.count += 1; windows.set(actor, window);
            if (window.count > 30) { res.setHeader('Retry-After', '60'); res.sendStatus(429); return; }
            const body = req.body;
            if (!body || Object.keys(body).some(key => !['id', 'kind', 'resource', 'occurredAt'].includes(key)) ||
              typeof body.id !== 'string' || !/^[a-zA-Z0-9_.:-]{1,120}$/.test(body.id) ||
              typeof body.kind !== 'string' || !Object.hasOwn(kinds, body.kind) ||
              typeof body.resource !== 'string' || !/^[a-z0-9][a-z0-9./_-]{0,120}$/.test(body.resource) ||
              typeof body.occurredAt !== 'string' || !Number.isFinite(Date.parse(body.occurredAt)) ||
              Math.abs(now - Date.parse(body.occurredAt)) > 24 * 3600000) {
              res.status(400).json({ error: 'Invalid event; use a supported kind, bounded identifier/resource, and a timestamp within 24 hours' }); return;
            }
            const kind = kinds[body.kind as keyof typeof kinds];
            const link = body.kind.startsWith('deployment.')
              ? `https://argocd.tail717ba0.ts.net/applications/argocd/${encodeURIComponent(body.resource)}`
              : body.kind.startsWith('backup.') ? '/docs/default/component/onemind-runbooks/resilience/' : '/catalog';
            await notifications.send({ recipients: { type: 'entity', entityRef: 'group:default/onemind-admins' }, payload: {
              title: `${kind.title}: ${body.resource}`, severity: kind.severity, topic: kind.topic,
              scope: `${body.kind}:${body.id}`, link,
              description: `Source reported ${body.kind} at ${new Date(body.occurredAt).toISOString()}. This event is not a live health check.`,
            } });
            // notificationService acknowledges processing, not resolved recipient count or delivery.
            res.status(202).json({ accepted: true, deliveryConfirmed: false, recipientReadiness: 'not-verified', recipientGroup: 'group:default/onemind-admins' });
          } catch (error) { next(error); }
        });
        router.use(routes);
      },
    });
  },
});
