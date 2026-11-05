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
import catalogPlugin from '@backstage/plugin-catalog-backend';
import unprocessedModule from '@backstage/plugin-catalog-backend-module-unprocessed';
import { createOneMindPolicy, ADMIN_GROUP_REF } from './policy';

jest.setTimeout(120_000);

/**
 * The real unprocessed-entities routes in the real catalog backend, with the real
 * OneMind policy answering their permission checks. Only identity lookup is faked.
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

describe('catalog unprocessed-entities routes are admin only', () => {
  let call: (
    path: string,
    authorization: string | undefined,
    method?: string,
  ) => Promise<Response>;
  const admin = mockCredentials.user.header('user:default/admin');
  const reader = mockCredentials.user.header('user:default/reader');

  let stopBackend: () => Promise<void>;
  afterAll(async () => {
    await stopBackend();
  });

  beforeAll(async () => {
    const backend = await startTestBackend({
      features: [
        catalogPlugin,
        unprocessedModule,
        permissionsThroughPolicy,
        mockServices.rootConfig.factory({
          data: { backend: { baseUrl: 'http://localhost' } },
        }),
      ],
    });
    stopBackend = () => backend.stop();
    call = (path, authorization, method = 'GET') =>
      fetch(`http://localhost:${backend.server.port()}/api/catalog${path}`, {
        method,
        headers: authorization ? { Authorization: authorization } : {},
      });
  });

  it.each(['failed', 'pending'])(
    'refuses a signed-in non-admin on %s with 403',
    async which => {
      const response = await call(`/entities/unprocessed/${which}`, reader);
      expect(response.status).toBe(403);
      expect(await response.text()).not.toMatch(/unprocessed_entity/);
    },
  );

  it('refuses a non-admin delete with 403', async () => {
    const response = await call(
      '/entities/unprocessed/delete/some-id',
      reader,
      'DELETE',
    );
    expect(response.status).toBe(403);
  });

  it('refuses a service credential such as the MCP machine token with 403', async () => {
    const service = mockCredentials.service.header();
    for (const which of ['failed', 'pending']) {
      expect(
        (await call(`/entities/unprocessed/${which}`, service)).status,
      ).toBe(403);
    }
    expect(
      (await call('/entities/unprocessed/delete/some-id', service, 'DELETE'))
        .status,
    ).toBe(403);
  });

  it('refuses a request with no credentials', async () => {
    // The test backend treats a missing token as the "none" principal, which the policy
    // denies; the production http auth service answers 401 before the policy.
    const response = await call('/entities/unprocessed/failed', undefined);
    expect([401, 403]).toContain(response.status);
  });

  it('lets an administrator read both lists', async () => {
    for (const which of ['failed', 'pending']) {
      const response = await call(`/entities/unprocessed/${which}`, admin);
      expect(response.status).toBe(200);
      expect(await response.json()).toMatchObject({
        type: which,
        entities: expect.any(Array),
      });
    }
  });
});
