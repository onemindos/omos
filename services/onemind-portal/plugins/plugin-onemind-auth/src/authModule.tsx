import { SignInPageBlueprint } from '@backstage/plugin-app-react';
import { SignInPage } from '@backstage/core-components';
import { OAuth2 } from '@backstage/core-app-api';
import {
  createFrontendModule,
  configApiRef,
  discoveryApiRef,
  oauthRequestApiRef,
  ApiBlueprint,
} from '@backstage/frontend-plugin-api';
import { oneMindAuthApiRef, ONEMIND_AUTH_PROVIDER_ID } from './authApi';

const signIn = SignInPageBlueprint.make({
  params: {
    loader: async () => props =>
      (
        <SignInPage
          {...props}
          provider={{
            id: 'onemind-auth',
            title: 'OneMind',
            message: 'Sign in with your OneMind account',
            apiRef: oneMindAuthApiRef,
          }}
        />
      ),
  },
});
// Standard Backstage OAuth2 session against the backend `oidc` provider. The same
// default scopes as before; the backend adds offline_access (auth.providers.oidc).
const oidcAuthApi = ApiBlueprint.make({
  name: 'onemind-oidc',
  params: define =>
    define({
      api: oneMindAuthApiRef,
      deps: {
        discoveryApi: discoveryApiRef,
        oauthRequestApi: oauthRequestApiRef,
        configApi: configApiRef,
      },
      factory: ({ discoveryApi, oauthRequestApi, configApi }) =>
        OAuth2.create({
          discoveryApi,
          oauthRequestApi,
          configApi,
          provider: {
            id: ONEMIND_AUTH_PROVIDER_ID,
            title: 'OneMind',
            icon: () => null,
          },
          environment: configApi.getOptionalString('auth.environment'),
          defaultScopes: ['openid', 'profile', 'email'],
        }),
    }),
});
export const authModule = createFrontendModule({
  pluginId: 'app',
  extensions: [signIn, oidcAuthApi],
});
