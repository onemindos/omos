import type {
  FetchApi,
  IdentityApi,
  OAuthApi,
  ConfigApi,
  DiscoveryApi,
} from '@backstage/core-plugin-api';

export function isWithinUrl(candidate: string, base: string): boolean {
  const url = new URL(candidate),
    target = new URL(base);
  const path = target.pathname.replace(/\/$/, '');
  return (
    url.origin === target.origin &&
    (url.pathname === path || url.pathname.startsWith(`${path}/`))
  );
}

// Construct each Request only once: some upstream fetch middleware consumes a
// Request body and then forwards the now-used original. Keep tokens URL-scoped.
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
  return {
    async fetch(input, init) {
      const inputUrl = input instanceof Request ? input.url : String(input);
      const parsed = new URL(inputUrl, backend);
      const plugin = parsed.protocol === 'plugin:';
      let target = parsed.href;
      if (plugin) {
        if (parsed.username || parsed.password)
          throw new Error('Plugin URL must not contain credentials');
        target = `${(await discovery.getBaseUrl(parsed.hostname)).replace(
          /\/$/,
          '',
        )}${parsed.pathname}${parsed.search}${parsed.hash}`;
      }
      const sourceHeaders = new Headers(
        init?.headers ?? (input instanceof Request ? input.headers : undefined),
      );
      const direct = sourceHeaders.has('x-openchoreo-direct');
      const host = isWithinUrl(target, backend);
      if (direct && !directUrls.some(base => isWithinUrl(target, base))) {
        throw new Error('Direct OpenChoreo request outside configured API');
      }
      const request = new Request(
        input instanceof Request ? input : parsed.href,
        {
          ...init,
          ...(direct || host ? { redirect: 'error' } : {}),
          ...(direct ? { credentials: 'omit' } : {}),
        },
      );
      request.headers.delete('x-openchoreo-direct');
      if (direct || host) {
        const token = await oauth.getAccessToken();
        if (!token) throw new Error('Missing OpenChoreo identity token');
        if (direct) {
          request.headers.delete('x-openchoreo-token');
          request.headers.set('Authorization', `Bearer ${token}`);
        } else {
          request.headers.set('x-openchoreo-token', token);
          if (!request.headers.has('Authorization')) {
            const { token: backstage } = await identity.getCredentials();
            if (backstage)
              request.headers.set('Authorization', `Bearer ${backstage}`);
          }
        }
      } else {
        request.headers.delete('x-openchoreo-token');
      }
      return plugin ? fetch(target, request) : fetch(request);
    },
  };
}
