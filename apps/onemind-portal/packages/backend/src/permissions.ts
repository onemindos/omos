import { createBackendModule } from '@backstage/backend-plugin-api';
import { AuthorizeResult } from '@backstage/plugin-permission-common';
import { policyExtensionPoint } from '@backstage/plugin-permission-node/alpha';

// Permissions remain enforced on the backend, independent of visible UI controls.
export default createBackendModule({
  pluginId: 'permission',
  moduleId: 'onemind-policy',
  register(reg) {
    reg.registerInit({
      deps: { policy: policyExtensionPoint },
      async init({ policy }) {
        policy.setPolicy({
          async handle(request, user) {
            if (!user) return { result: AuthorizeResult.DENY };
            const admin = user.info.ownershipEntityRefs.includes('group:default/onemind-admins');
            return { result: admin || request.permission.attributes.action === 'read'
              ? AuthorizeResult.ALLOW : AuthorizeResult.DENY };
          },
        });
      },
    });
  },
});
