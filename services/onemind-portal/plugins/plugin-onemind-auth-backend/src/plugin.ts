import { createBackendModule } from '@backstage/backend-plugin-api';
import { authProvidersExtensionPoint } from '@backstage/plugin-auth-node';
import { registerOneMindProviders } from './providers';

// Sign-in is Backstage's standard OIDC provider against Authentik (auth.providers.oidc).
export default createBackendModule({
  pluginId: 'auth',
  moduleId: 'onemind-oidc',
  register(reg) {
    reg.registerInit({
      deps: { providers: authProvidersExtensionPoint },
      async init({ providers }) {
        registerOneMindProviders(providers);
      },
    });
  },
});
