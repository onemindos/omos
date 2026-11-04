import { createBackendModule } from '@backstage/backend-plugin-api';
import { policyExtensionPoint } from '@backstage/plugin-permission-node/alpha';
import { createOneMindPolicy } from './policy';

// Installs the single OneMind policy. It needs no other service and no OpenChoreo
// configuration; the permission backend stays enabled through `permission.enabled`.
export default createBackendModule({
  pluginId: 'permission',
  moduleId: 'onemind-policy',
  register(reg) {
    reg.registerInit({
      deps: { policy: policyExtensionPoint },
      async init({ policy }) {
        policy.setPolicy(createOneMindPolicy());
      },
    });
  },
});
