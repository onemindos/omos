/** @jest-environment node */
import { createScopedFetch, isWithinUrl } from './scopedFetch';
import type {
  ConfigApi,
  IdentityApi,
  OAuthApi,
  DiscoveryApi,
} from '@backstage/core-plugin-api';

const backend = 'https://portal.example.test';
const api = 'https://control.example.test/api/v1';
const observer = 'https://observer.example.test';
const config = {
  getString: (key: string) => (key === 'backend.baseUrl' ? backend : api),
  getOptionalStringArray: () => [api, observer],
} as unknown as ConfigApi;
const oauth = {
  getAccessToken: jest.fn(async () => 'idp-fixture'),
} as unknown as OAuthApi;
const identity = {
  getCredentials: jest.fn(async () => ({ token: 'backstage-fixture' })),
} as unknown as IdentityApi;
const discovery = {
  getBaseUrl: jest.fn(async () => `${backend}/api/catalog`),
} as unknown as DiscoveryApi;
const originalFetch = global.fetch;
let outgoing: Request;
beforeEach(() => {
  jest.clearAllMocks();
  global.fetch = jest.fn(async input => {
    outgoing = input instanceof Request ? input : new Request(input);
    return new Response('{}', { status: 200 });
  });
});
afterAll(() => {
  global.fetch = originalFetch;
});
const client = () => createScopedFetch({ config, identity, oauth, discovery });

test.each([
  'https://portal.example.test.evil.test/api',
  'https://portal.example.test@evil.test/api',
  'http://portal.example.test/api',
])('reject origin confusion: %s', url => {
  expect(isWithinUrl(url, backend)).toBe(false);
});
test('path allowlist uses segment boundaries', () => {
  expect(isWithinUrl(`${api}-evil`, api)).toBe(false);
  expect(isWithinUrl(`${api}/components`, api)).toBe(true);
});
test('external fetch receives neither injected token nor direct signal', async () => {
  await client().fetch('https://external.example.test/data', {
    headers: { 'x-openchoreo-token': 'discard' },
  });
  expect(outgoing.headers.has('Authorization')).toBe(false);
  expect(outgoing.headers.has('x-openchoreo-token')).toBe(false);
  expect(oauth.getAccessToken).not.toHaveBeenCalled();
});
test('host request receives both token types and cannot redirect them', async () => {
  await client().fetch(`${backend}/api/catalog/entities`);
  expect(outgoing.headers.get('Authorization')).toBe(
    'Bearer backstage-fixture',
  );
  expect(outgoing.headers.get('x-openchoreo-token')).toBe('idp-fixture');
  expect(outgoing.redirect).toBe('error');
});
test.each([`${api}/components`, `${observer}/api/v1/metrics/query`])(
  'trusted direct request has only IDP auth: %s',
  async url => {
    await client().fetch(url, { headers: { 'x-openchoreo-direct': 'true' } });
    expect(outgoing.headers.get('Authorization')).toBe('Bearer idp-fixture');
    expect(outgoing.headers.has('x-openchoreo-token')).toBe(false);
    expect(outgoing.headers.has('x-openchoreo-direct')).toBe(false);
    expect(outgoing.redirect).toBe('error');
    expect(identity.getCredentials).not.toHaveBeenCalled();
  },
);
test('untrusted direct request is blocked before reading tokens or fetching', async () => {
  await expect(
    client().fetch('https://evil.example.test', {
      headers: { 'x-openchoreo-direct': 'true' },
    }),
  ).rejects.toThrow('outside configured API');
  expect(global.fetch).not.toHaveBeenCalled();
  expect(oauth.getAccessToken).not.toHaveBeenCalled();
});
