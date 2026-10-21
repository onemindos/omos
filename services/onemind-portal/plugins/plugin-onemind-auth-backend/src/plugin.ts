import {
  createBackendModule,
  coreServices,
} from '@backstage/backend-plugin-api';
import {
  authProvidersExtensionPoint,
  createOAuthProviderFactory,
} from '@backstage/plugin-auth-node';
import { oidcAuthenticator } from '@backstage/plugin-auth-backend-module-oidc-provider';
import { stringifyEntityRef } from '@backstage/catalog-model';
import { signInKnownSubject } from './subject';

// Retain verified OIDC and explicit catalog-subject mapping. No email identities
// or token-decoded group promotion. The OC frontend requires only the provider ID.
export default createBackendModule({
  pluginId: 'auth',
  moduleId: 'onemind-oidc',
  register(reg) {
    reg.registerInit({
      deps: {
        providers: authProvidersExtensionPoint,
        discovery: coreServices.discovery,
        auth: coreServices.auth,
        logger: coreServices.logger,
      },
      async init({ providers, discovery, auth, logger }) {
        for (const providerId of ['oidc', 'openchoreo-auth']) {
          providers.registerProvider({
            providerId,
            factory: createOAuthProviderFactory({
              authenticator: oidcAuthenticator,
              async signInResolver(info, ctx) {
                const subject = info.result.fullProfile.userinfo.sub;
                const { query, identity } = await signInKnownSubject(
                  subject,
                  ctx,
                );
                if (providerId === 'openchoreo-auth') {
                  const { entity } = await ctx.findCatalogUser(query);
                  const { token } = await auth.getPluginRequestToken({
                    onBehalfOf: await auth.getOwnServiceCredentials(),
                    targetPluginId: 'permission',
                  });
                  const response = await fetch(
                    `${await discovery.getBaseUrl(
                      'permission',
                    )}/cache-capabilities`,
                    {
                      method: 'POST',
                      headers: {
                        'Content-Type': 'application/json',
                        Authorization: `Bearer ${token}`,
                      },
                      body: JSON.stringify({
                        userEntityRef: stringifyEntityRef(entity),
                        accessToken: info.result.session.accessToken,
                      }),
                    },
                  );
                  if (!response.ok) {
                    logger.warn(
                      `OpenChoreo capability initialization failed: HTTP ${response.status}`,
                    );
                    throw new Error(
                      'Unable to initialize platform permissions',
                    );
                  }
                }
                return identity;
              },
            }),
          });
        }
      },
    });
  },
});
