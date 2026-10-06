import { createBackendModule } from '@backstage/backend-plugin-api';
import { authProvidersExtensionPoint, createOAuthProviderFactory } from '@backstage/plugin-auth-node';
import { oidcAuthenticator } from '@backstage/plugin-auth-backend-module-oidc-provider';

// Map the signed subject from this configured issuer to an explicitly registered catalog user.
// No email trust override or automatic administrator registration is permitted.
export default createBackendModule({
  pluginId: 'auth',
  moduleId: 'onemind-oidc',
  register(reg) {
    reg.registerInit({
      deps: { providers: authProvidersExtensionPoint },
      async init({ providers }) {
        providers.registerProvider({
          providerId: 'oidc',
          factory: createOAuthProviderFactory({
            authenticator: oidcAuthenticator,
            async signInResolver(info, ctx) {
              const subject = info.result.fullProfile.userinfo.sub;
              if (!subject) throw new Error('OIDC subject is missing');
              return ctx.signInWithCatalogUser({ annotations: { 'onemind.io/oidc-sub': subject } });
            },
          }),
        });
      },
    });
  },
});
