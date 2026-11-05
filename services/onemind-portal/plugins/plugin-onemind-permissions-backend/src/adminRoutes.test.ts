import {
  coreServices,
  createServiceFactory,
} from '@backstage/backend-plugin-api';
import type { BackstageCredentials } from '@backstage/backend-plugin-api';
import {
  mockCredentials,
  mockServices,
  startTestBackend,
} from '@backstage/backend-test-utils';
import devtoolsPlugin from '@backstage/plugin-devtools-backend';
import { createOneMindPolicy, ADMIN_GROUP_REF } from './policy';

jest.setTimeout(60_000);

/**
 * The real devtools-backend router, with the real OneMind policy answering its
 * permission checks. Only the identity lookup is faked: the admin user is a member of
 * group:default/onemind-admins, the reader is not.
 */
const MEMBERSHIPS: Record<string, string[]> = {
  'user:default/admin': ['user:default/admin', ADMIN_GROUP_REF],
  'user:default/reader': ['user:default/reader', 'group:default/family'],
};

const permissionsThroughPolicy = createServiceFactory({
  service: coreServices.permissions,
  deps: {},
  factory: () => {
    const policy = createOneMindPolicy();
    const authorize = async (
      requests: { permission: any }[],
      options: { credentials: BackstageCredentials },
    ) => {
      const { principal } = options.credentials as any;
      const user =
        principal.type === 'user'
          ? {
              credentials: options.credentials,
              info: {
                userEntityRef: principal.userEntityRef,
                ownershipEntityRefs: MEMBERSHIPS[principal.userEntityRef] ?? [],
              },
            }
          : undefined;
      return Promise.all(requests.map(r => policy.handle(r, user as any)));
    };
    return { authorize, authorizeConditional: authorize } as any;
  },
});

async function start() {
  const backend = await startTestBackend({
    features: [
      devtoolsPlugin,
      permissionsThroughPolicy,
      mockServices.rootConfig.factory({
        data: { backend: { baseUrl: 'http://localhost' } },
      }),
    ],
  });
  const call = (path: string, authorization?: string) =>
    fetch(`http://localhost:${backend.server.port()}/api/devtools${path}`, {
      headers: authorization ? { Authorization: authorization } : {},
    });
  return {
    stop: () => backend.stop(),
    call,
    asAdmin: (path: string) =>
      call(path, mockCredentials.user.header('user:default/admin')),
    asReader: (path: string) =>
      call(path, mockCredentials.user.header('user:default/reader')),
    asService: (path: string) => call(path, mockCredentials.service.header()),
  };
}

describe('DevTools routes are admin only', () => {
  let backend: Awaited<ReturnType<typeof start>>;
  beforeAll(async () => {
    backend = await start();
  });
  afterAll(async () => {
    await backend.stop();
  });

  it('refuses a signed-in non-admin on /info with 403', async () => {
    const response = await backend.asReader('/info');
    expect(response.status).toBe(403);
    expect(JSON.stringify(await response.json())).not.toMatch(/dependencies/);
  });

  it('refuses a signed-in non-admin on /config with 403', async () => {
    expect((await backend.asReader('/config')).status).toBe(403);
  });

  it('refuses a signed-in non-admin on /external-dependencies with 403', async () => {
    // Not reachable from the portal UI, but the route exists in the backend.
    const response = await backend.asReader('/external-dependencies');
    expect(response.status).toBe(403);
  });

  it('refuses a service credential such as the MCP machine token with 403', async () => {
    for (const path of ['/info', '/config', '/external-dependencies']) {
      expect((await backend.asService(path)).status).toBe(403);
    }
  });

  it('refuses a request with no credentials', async () => {
    // The test backend turns a missing token into the "none" principal, which the
    // policy denies; the production http auth service answers 401 before the policy.
    for (const path of ['/info', '/config']) {
      expect([401, 403]).toContain((await backend.call(path)).status);
    }
  });

  it('never serves the raw configuration, not even to an administrator', async () => {
    const response = await backend.asAdmin('/config');
    expect(response.status).toBe(403);
    expect(await response.text()).not.toMatch(/"config"/);
  });

  it('lets an administrator read /info', async () => {
    const response = await backend.asAdmin('/info');
    expect(response.status).toBe(200);
    const body = await response.json();
    expect(body).toHaveProperty('nodeJsVersion');
    expect(body).toHaveProperty('dependencies');
  });
});
