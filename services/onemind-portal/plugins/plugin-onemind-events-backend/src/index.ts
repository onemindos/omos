import { createBackendFeatureLoader } from '@backstage/backend-plugin-api';
import { operationalEventsModule } from './plugin';

// One default export keeps packages/backend untouched: the loader returns every module.
export default createBackendFeatureLoader({
  loader() {
    return [operationalEventsModule];
  },
});

export { operationalEventsModule };
