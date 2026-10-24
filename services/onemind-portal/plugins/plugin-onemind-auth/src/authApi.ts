import {
  createApiRef,
  type ApiRef,
  type BackstageIdentityApi,
  type OAuthApi,
  type ProfileInfoApi,
  type SessionApi,
} from '@backstage/core-plugin-api';

// Backstage ships no ApiRef for the generic OIDC provider, so the portal declares its own.
// It talks to the backend `oidc` provider (auth.providers.oidc), which points at Authentik.
export const ONEMIND_AUTH_PROVIDER_ID = 'oidc';

export const oneMindAuthApiRef: ApiRef<
  OAuthApi & ProfileInfoApi & BackstageIdentityApi & SessionApi
> = createApiRef({ id: 'internal.auth.onemind-oidc' });
