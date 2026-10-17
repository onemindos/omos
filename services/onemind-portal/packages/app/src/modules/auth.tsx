import { OpenIdConnectApi, ProfileInfoApi, BackstageIdentityApi, SessionApi } from '@backstage/core-plugin-api';
import { OAuth2 } from '@backstage/core-app-api';
import { SignInPageBlueprint } from '@backstage/plugin-app-react';
import { SignInPage } from '@backstage/core-components';
import { createApiRef, createFrontendModule, configApiRef, discoveryApiRef, oauthRequestApiRef, ApiBlueprint } from '@backstage/frontend-plugin-api';

const authApiRef = createApiRef<OpenIdConnectApi & ProfileInfoApi & BackstageIdentityApi & SessionApi>().with({ id: 'auth.onemind' });
const authApi = ApiBlueprint.make({
  name: 'onemind',
  params: defineParams => defineParams({
    api: authApiRef,
    deps: { discoveryApi: discoveryApiRef, oauthRequestApi: oauthRequestApiRef, configApi: configApiRef },
    factory: ({ discoveryApi, oauthRequestApi, configApi }) => OAuth2.create({
      configApi, discoveryApi, oauthRequestApi,
      environment: configApi.getOptionalString('auth.environment'),
      provider: { id: 'oidc', title: 'OneMind', icon: () => null },
      defaultScopes: ['openid', 'profile', 'email'],
    }),
  }),
});
const signIn = SignInPageBlueprint.make({
  params: { loader: async () => props => <SignInPage {...props} provider={{
    id: 'onemind-auth', title: 'OneMind', message: 'Sign in with your OneMind account', apiRef: authApiRef,
  }} /> },
});
export const authModule = createFrontendModule({ pluginId: 'app', extensions: [authApi, signIn] });
