import type {
  FetchApi,
  IdentityApi,
  OAuthApi,
  ConfigApi,
  DiscoveryApi,
} from '@backstage/core-plugin-api';
import { createFetchApi, FetchMiddlewares } from '@backstage/core-app-api';

export function isWithinUrl(candidate: string, base: string): boolean {
  const url = new URL(candidate),
    target = new URL(base);
  const path = target.pathname.replace(/\/$/, '');
  return (
    url.origin === target.origin &&
    (url.pathname === path || url.pathname.startsWith(`${path}/`))
  );
}

// Never attach platform tokens to arbitrary external URLs used by existing plugins.
export function createScopedFetch(options: {
  config: ConfigApi;
  identity: IdentityApi;
  oauth: OAuthApi;
  discovery: DiscoveryApi;
}): FetchApi {
  const { config, identity, oauth, discovery } = options;
  const backend = config.getString('backend.baseUrl');
  const directUrls =
    config.getOptionalStringArray('onemind.openchoreo.directApiAllowlist') ??
    [];
  const normal = createFetchApi({
    middleware: [
      FetchMiddlewares.resolvePluginProtocol({ discoveryApi: discovery }),
      FetchMiddlewares.injectIdentityAuth({
        identityApi: identity,
        allowUrl: url => isWithinUrl(url, backend),
      }),
    ],
  });
  return {
    async fetch(input, init) {
      const request = new Request(input, init);
      const headers = new Headers(request.headers);
      const direct = headers.has('x-openchoreo-direct');
      headers.delete('x-openchoreo-direct');
      if (direct) {
        if (!directUrls.some(base => isWithinUrl(request.url, base)))
          throw new Error('Direct OpenChoreo request outside configured API');
        const token = await oauth.getAccessToken();
        if (!token) throw new Error('Missing OpenChoreo identity token');
        headers.delete('x-openchoreo-token');
        headers.set('Authorization', `Bearer ${token}`);
        return fetch(
          new Request(request, {
            headers,
            credentials: 'omit',
            redirect: 'error',
          }),
        );
      }
      if (isWithinUrl(request.url, backend)) {
        const token = await oauth.getAccessToken();
        if (!token) throw new Error('Missing OpenChoreo identity token');
        headers.set('x-openchoreo-token', token);
      } else {
        headers.delete('x-openchoreo-token');
      }
      return normal.fetch(
        new Request(request, {
          headers,
          ...(isWithinUrl(request.url, backend)
            ? { redirect: 'error' as const }
            : {}),
        }),
      );
    },
  };
}
