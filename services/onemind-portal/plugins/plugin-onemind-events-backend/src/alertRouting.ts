import { createBackendModule, coreServices } from '@backstage/backend-plugin-api';
import { notificationService } from '@backstage/plugin-notifications-node';
import { eventsServiceRef } from '@backstage/plugin-events-node';
import { signalsServiceRef } from '@backstage/plugin-signals-node';
import {
  ALERT_SIGNAL_CHANNEL,
  alertNotification,
  isSignalSeverity,
  parseAlertmanagerWebhook,
  type AlertSignal,
} from '@internal/plugin-onemind-common-backend';
import { ALERTMANAGER_TOPIC } from './alertmanager';
import { ADMIN_GROUP_REF, createAdminRecipients } from './recipients';
import type { NotificationPayload, Sinks } from './scan';
import { createStore, type Store } from './store';

export type AlertRoutingResult = { notified: number; resolved: number; repeats: number; ignored: number };

// Turns one Alertmanager webhook payload into inbox items and signals. Alertmanager repeats
// firing alerts every few hours; the stored state per fingerprint turns those into no-ops,
// and the per-fingerprint scope makes the resolved update replace the firing item.
export async function routeAlerts(args: {
  body: unknown;
  store: Store;
  sinks: Sinks;
  now?: () => Date;
}): Promise<AlertRoutingResult> {
  const result: AlertRoutingResult = { notified: 0, resolved: 0, repeats: 0, ignored: 0 };
  const parsed = parseAlertmanagerWebhook(args.body);
  if (!parsed.ok) return result;
  const now = args.now ?? (() => new Date());
  for (const alert of parsed.alerts) {
    const item = alertNotification(alert);
    if (!item) {
      result.ignored += 1;
      continue;
    }
    const known = await args.store.alerts.get(item.fingerprint);
    const signal = (action: AlertSignal['action']): AlertSignal => ({
      action,
      source: 'alertmanager',
      key: item.payload.scope,
      severity: item.firingSeverity,
      title: item.payload.title,
      at: now().toISOString(),
    });
    if (item.status === 'firing') {
      if (known?.status === 'firing') {
        result.repeats += 1;
        continue;
      }
      await args.sinks.notify(item.payload as NotificationPayload);
      if (isSignalSeverity(item.firingSeverity)) await args.sinks.signal(signal('alert.firing'));
      await args.store.alerts.put({ fingerprint: item.fingerprint, status: 'firing', severity: item.firingSeverity });
      result.notified += 1;
    } else {
      // Resolved for something we never announced (or already resolved): nothing to update.
      if (known?.status !== 'firing') {
        result.ignored += 1;
        continue;
      }
      await args.sinks.notify(item.payload as NotificationPayload);
      if (isSignalSeverity(item.firingSeverity)) await args.sinks.signal(signal('alert.resolved'));
      await args.store.alerts.put({ fingerprint: item.fingerprint, status: 'resolved', severity: item.firingSeverity });
      result.resolved += 1;
    }
  }
  return result;
}

export const alertRoutingModule = createBackendModule({
  pluginId: 'onemind',
  moduleId: 'alert-routing',
  register(reg) {
    reg.registerInit({
      deps: {
        events: eventsServiceRef,
        notifications: notificationService,
        signals: signalsServiceRef,
        database: coreServices.database,
        logger: coreServices.logger,
        auth: coreServices.auth,
        discovery: coreServices.discovery,
      },
      async init({ events, notifications, signals, database, logger, auth, discovery }) {
        const store = await createStore(await database.getClient());
        const recipients = createAdminRecipients({ auth, discovery });
        await events.subscribe({
          id: 'onemind-alertmanager-routing',
          topics: [ALERTMANAGER_TOPIC],
          onEvent: async ({ eventPayload }) => {
            try {
              // Signals need explicit users; the inbox item itself resolves the group.
              const userRefs = await recipients.userRefs().catch(() => undefined);
              const result = await routeAlerts({
                body: eventPayload,
                store,
                sinks: {
                  notify: async payload => {
                    await notifications.send({ recipients: { type: 'entity', entityRef: ADMIN_GROUP_REF }, payload });
                  },
                  signal: async message => {
                    if (!userRefs) {
                      logger.warn('Alert signal skipped: admin recipients not verified');
                      return;
                    }
                    await signals.publish({ recipients: { type: 'user', entityRef: userRefs }, channel: ALERT_SIGNAL_CHANNEL, message });
                  },
                },
              });
              logger.info(`Alertmanager webhook processed: ${JSON.stringify(result)}`);
            } catch (error) {
              logger.error(`Alertmanager webhook processing failed: ${(error as Error).message}`);
            }
          },
        });
      },
    });
  },
});
