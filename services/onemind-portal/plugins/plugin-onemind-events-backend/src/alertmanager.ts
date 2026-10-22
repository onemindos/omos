import { createBackendModule, coreServices } from '@backstage/backend-plugin-api';
import { eventsExtensionPoint } from '@backstage/plugin-events-node/alpha';
import type { RequestValidator } from '@backstage/plugin-events-node';
import { createHash, timingSafeEqual } from 'node:crypto';
import { MAX_WEBHOOK_BYTES, parseAlertmanagerWebhook } from '@internal/plugin-onemind-common-backend';

// Published at POST /api/events/http/alertmanager by the native events-backend ingress.
export const ALERTMANAGER_TOPIC = 'alertmanager';
// Optional. Absent or too short: the route rejects every request and the portal starts normally.
export const WEBHOOK_TOKEN_KEY = 'onemind.alerts.webhookToken';
export const MIN_TOKEN_LENGTH = 24;

const digest = (value: string) => createHash('sha256').update(value).digest();
const safeEqual = (a: string, b: string) => timingSafeEqual(digest(a), digest(b));

// The events-backend /http prefix is unauthenticated by design; webhook senders are expected
// to be verified by a topic validator (as the GitHub module does). This one checks a bearer
// token before anything is parsed or published.
export function createAlertmanagerValidator(getToken: () => string | undefined): RequestValidator {
  return async (request, context) => {
    const token = getToken();
    if (!token || token.length < MIN_TOKEN_LENGTH) {
      context.reject({ status: 503, payload: { error: 'Alertmanager webhook is not enabled' } });
      return;
    }
    const header = request.headers.authorization;
    const match = typeof header === 'string' ? /^Bearer (.{1,512})$/.exec(header) : null;
    if (!match || !safeEqual(match[1], token)) {
      context.reject({ status: 401, payload: { error: 'Invalid or missing webhook token' } });
      return;
    }
    if (request.raw.body.length > MAX_WEBHOOK_BYTES) {
      context.reject({ status: 413, payload: { error: 'Payload too large' } });
      return;
    }
    const parsed = parseAlertmanagerWebhook(request.body);
    if (!parsed.ok) {
      context.reject({ status: 400, payload: { error: parsed.error } });
    }
  };
}

export const alertmanagerIngressModule = createBackendModule({
  pluginId: 'events',
  moduleId: 'onemind-alertmanager-ingress',
  register(reg) {
    reg.registerInit({
      deps: { events: eventsExtensionPoint, config: coreServices.rootConfig, logger: coreServices.logger },
      async init({ events, config, logger }) {
        const token = config.getOptionalString(WEBHOOK_TOKEN_KEY);
        if (!token || token.length < MIN_TOKEN_LENGTH) {
          logger.info(`Alertmanager webhook disabled: ${WEBHOOK_TOKEN_KEY} is not set (or too short)`);
        }
        events.addHttpPostIngress({
          topic: ALERTMANAGER_TOPIC,
          validator: createAlertmanagerValidator(() => config.getOptionalString(WEBHOOK_TOKEN_KEY)),
        });
      },
    });
  },
});
