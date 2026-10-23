import {
  createOAuthProviderFactory,
  type AuthProviderRegistrationOptions,
  type AuthResolverContext,
} from '@backstage/plugin-auth-node';
import { oidcAuthenticator } from '@backstage/plugin-auth-backend-module-oidc-provider';
import { signInKnownSubject } from './subject';

// The one sign-in provider id. The frontend (plugin-onemind-auth) must use the same id.
export const ONEMIND_OIDC_PROVIDER_ID = 'oidc';

type ProviderRegistry = {
  registerProvider(options: AuthProviderRegistrationOptions): void;
};

// Identity is the verified, signed OIDC `sub` mapped to the exact `onemind.io/oidc-sub`
// catalog annotation. No email identities, no token-decoded group promotion, no
// auto-registration, and no call to any other service at sign-in.
export async function oneMindSignInResolver(
  info: { result: { fullProfile: { userinfo: { sub?: unknown } } } },
  ctx: Pick<AuthResolverContext, 'signInWithCatalogUser'>,
) {
  const { identity } = await signInKnownSubject(
    info.result.fullProfile.userinfo.sub,
    ctx,
  );
  return identity;
}

// Backstage's own OIDC provider, pointed at Authentik by auth.providers.oidc.* config.
export function registerOneMindProviders(providers: ProviderRegistry) {
  providers.registerProvider({
    providerId: ONEMIND_OIDC_PROVIDER_ID,
    factory: createOAuthProviderFactory({
      authenticator: oidcAuthenticator,
      signInResolver: oneMindSignInResolver,
    }),
  });
}
