import { createBackendFeatureLoader } from '@backstage/backend-plugin-api';
import { alertmanagerIngressModule } from './alertmanager';
import { alertRoutingModule } from './alertRouting';
import { operationalEventsModule } from './plugin';

// One default export keeps packages/backend untouched: the loader returns every module.
export default createBackendFeatureLoader({
  loader() {
    return [operationalEventsModule, alertmanagerIngressModule, alertRoutingModule];
  },
});

export { operationalEventsModule, alertmanagerIngressModule, alertRoutingModule };
export { createAlertmanagerValidator } from './alertmanager';
export { routeAlerts } from './alertRouting';
