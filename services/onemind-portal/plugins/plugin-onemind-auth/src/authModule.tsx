import {
  SignInPageBlueprint,
  AppRootWrapperBlueprint,
} from '@backstage/plugin-app-react';
import { SignInPage } from '@backstage/core-components';
import {
  createFrontendModule,
  configApiRef,
  discoveryApiRef,
  identityApiRef,
  fetchApiRef,
  ApiBlueprint,
} from '@backstage/frontend-plugin-api';
import { permissionApiRef } from '@backstage/plugin-permission-react';
import {
  openChoreoAuthApiRef,
  OpenChoreoPermissionApi,
} from '@openchoreo/backstage-plugin';
import { OpenChoreoQueryProvider } from '@openchoreo/backstage-plugin-react';
import { createScopedFetch } from './scopedFetch';

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
            apiRef: openChoreoAuthApiRef,
          }}
        />
      ),
  },
});
const fetchApi = ApiBlueprint.make({
  name: fetchApiRef.id,
  params: define =>
    define({
      api: fetchApiRef,
      deps: {
        config: configApiRef,
        identity: identityApiRef,
        oauth: openChoreoAuthApiRef,
        discovery: discoveryApiRef,
      },
      factory: createScopedFetch,
    }),
});
const permissionApi = ApiBlueprint.make({
  name: permissionApiRef.id,
  params: define =>
    define({
      api: permissionApiRef,
      deps: {
        config: configApiRef,
        identity: identityApiRef,
        oauthApi: openChoreoAuthApiRef,
        discovery: discoveryApiRef,
      },
      factory: options => new OpenChoreoPermissionApi(options),
    }),
});
const queries = AppRootWrapperBlueprint.make({
  name: 'openchoreo-query',
  params: { component: OpenChoreoQueryProvider },
});
export const authModule = createFrontendModule({
  pluginId: 'app',
  extensions: [signIn, fetchApi, permissionApi, queries],
});
