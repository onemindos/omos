import { existsSync } from 'node:fs';
import { resolve } from 'node:path';
import type { JsonObject } from '@backstage/types';
import { ConfigReader, type Config } from '@backstage/config';
import { loadConfig } from '@backstage/config-loader';
import { mockServices } from '@backstage/backend-test-utils';
import type { EntityProviderConnection } from '@backstage/plugin-catalog-node';
import { readAuthentikOrgConfig, reviewedOrgFile } from './config';
import { registerAuthentikOrg } from './module';
import { orgYamlPath, portalPath } from './testing';

const TOKEN = 'fixture-only-token';
/** The reviewed org file's catalog location as app-config.yaml lists it. */
const ORG_LOCATION = {
  type: 'file',
  target: '../../catalog/org.yaml',
  rules: [{ allow: ['User', 'Group'] }],
};

const enabled = (
  orgSync: JsonObject = {},
  authentik: JsonObject = {},
): JsonObject => ({
  catalog: { locations: [ORG_LOCATION] },
  authentik: {
    baseUrl: 'http://authentik.test/',
    token: TOKEN,
    ...authentik,
    orgSync: {
      enabled: true,
      userGroups: ['OneMind Portal Users'],
      ...orgSync,
    },
  },
});

const make = async (
  config: Config,
  scheduleTask = jest.fn().mockResolvedValue(undefined),
) => {
  const processing = { addEntityProvider: jest.fn() };
  const scheduler = { scheduleTask };
  const fetchJson = jest.fn();
  const logger = mockServices.logger.mock();
  const mode = await registerAuthentikOrg({
    config,
    logger,
    scheduler: scheduler as never,
    processing,
    fetchJson,
  });
  return { mode, processing, scheduler, fetchJson, logger };
};
type Made = Awaited<ReturnType<typeof make>>;

/** The default schedule, which the removal always uses. */
const TASK = {
  id: 'onemind-authentik-org-refresh',
  frequency: { seconds: 600 },
  timeout: { seconds: 120 },
  initialDelay: { seconds: 60 },
  scope: 'global',
};
const EMPTY_FULL_MUTATION = [{ type: 'full', entities: [] }];

/** Every state registers one provider under one name and one task id, so replicas share its runs. */
function registeredTask(m: Made) {
  expect(m.processing.addEntityProvider).toHaveBeenCalledTimes(1);
  expect(
    m.processing.addEntityProvider.mock.calls[0][0].getProviderName(),
  ).toBe('onemind-authentik-org');
  expect(m.scheduler.scheduleTask).toHaveBeenCalledTimes(1);
  const { fn, ...task } = m.scheduler.scheduleTask.mock.calls[0][0];
  expect(typeof fn).toBe('function');
  return task;
}

/** Connects the registered provider to a recording catalog, runs its task once, returns the mutations. */
async function connectAndRunOnce(m: Made) {
  const applyMutation = jest.fn();
  const connection = {
    applyMutation,
    refresh: jest.fn(),
  } as unknown as EntityProviderConnection;
  await m.processing.addEntityProvider.mock.calls[0][0].connect(connection);
  await m.scheduler.scheduleTask.mock.calls[0][0].fn(
    new AbortController().signal,
  );
  return applyMutation.mock.calls;
}

describe('feature flag', () => {
  it.each([
    ['absent', {}],
    ['false', { authentik: { orgSync: { enabled: false } } }],
    [
      'false with a token, URL and scope set',
      {
        authentik: {
          baseUrl: 'http://authentik.test',
          token: TOKEN,
          orgSync: { enabled: false, userGroups: ['x'] },
        },
      },
    ],
    [
      'unset with other settings',
      { authentik: { orgSync: { frequencySeconds: 30 } } },
    ],
    [
      'false, with invalid settings it then never reads',
      {
        authentik: {
          baseUrl: 'not a url',
          orgSync: { enabled: false, frequencySeconds: -1 },
        },
      },
    ],
  ])(
    'calls nothing when the flag is %s, and only removes what the sync published, under the same task',
    async (_name, data) => {
      const m = await make(new ConfigReader(data));
      expect(m.mode).toBe('cleanup');
      expect(registeredTask(m)).toEqual(TASK);
      // At startup (connect) and on every run: one full mutation without entities.
      expect(await connectAndRunOnce(m)).toEqual([
        EMPTY_FULL_MUTATION,
        EMPTY_FULL_MUTATION,
      ]);
      expect(m.fetchJson).not.toHaveBeenCalled();
      expect(m.logger.error).not.toHaveBeenCalled();
    },
  );

  it('registers one provider and one global task when enabled, and reads nothing until it runs', async () => {
    const m = await make(new ConfigReader(enabled()));
    expect(m.mode).toBe('sync');
    expect(registeredTask(m)).toEqual(TASK);
    // Connecting the sync (a pod start) never removes what it published.
    const applyMutation = jest.fn();
    await m.processing.addEntityProvider.mock.calls[0][0].connect({
      applyMutation,
      refresh: jest.fn(),
    });
    expect(applyMutation).not.toHaveBeenCalled();
    expect(m.fetchJson).not.toHaveBeenCalled();
    const logged = JSON.stringify(m.logger.info.mock.calls);
    expect(logged).toContain(
      'members of OneMind Portal Users every 600s from http://authentik.test',
    );
    expect(logged).not.toContain(TOKEN);
  });

  it.each([
    ['unset (its secret is missing)', undefined],
    ['blank', '   '],
  ])(
    'does not fail startup when enabled with the token %s: it calls nothing, removes what it published and logs why',
    async (_name, token) => {
      const m = await make(new ConfigReader(enabled({}, { token })));
      expect(m.mode).toBe('cleanup');
      expect(registeredTask(m)).toEqual(TASK);
      expect(await connectAndRunOnce(m)).toEqual([
        EMPTY_FULL_MUTATION,
        EMPTY_FULL_MUTATION,
      ]);
      expect(m.fetchJson).not.toHaveBeenCalled();
      expect(m.logger.error).toHaveBeenCalledTimes(1);
      const [message] = m.logger.error.mock.calls[0];
      expect(message).toContain('authentik.token is not set');
      expect(message).toContain(
        'Secret backstage-authentik-catalog, key AUTHENTIK_CATALOG_TOKEN',
      );
    },
  );

  it('normalises the URL, de-duplicates the scope and reads the org file the catalog reads', () => {
    expect(
      readAuthentikOrgConfig(
        new ConfigReader(
          enabled(
            { userGroups: [' A ', 'A', 'B'], frequencySeconds: 30 },
            { baseUrl: 'https://auth.example.test/prefix//' },
          ),
        ),
      ),
    ).toMatchObject({
      baseUrl: 'https://auth.example.test/prefix',
      token: TOKEN,
      userGroups: ['A', 'B'],
      frequencySeconds: 30,
      curatedFile: '../../catalog/org.yaml',
    });
  });

  it.each([
    [
      'no baseUrl',
      enabled({}, { baseUrl: undefined }),
      'authentik.baseUrl is required',
    ],
    [
      'a baseUrl without scheme',
      enabled({}, { baseUrl: 'authentik.test' }),
      'must be an http(s) URL',
    ],
    [
      'a non-http baseUrl',
      enabled({}, { baseUrl: 'ftp://authentik.test' }),
      'must be an http(s) URL',
    ],
    [
      'credentials in the baseUrl',
      enabled({}, { baseUrl: `http://svc:${TOKEN}@authentik.test` }),
      'without credentials, query or fragment',
    ],
    [
      'no scope',
      enabled({ userGroups: undefined }),
      'authentik.orgSync.userGroups',
    ],
    [
      'an empty scope',
      enabled({ userGroups: [] }),
      'authentik.orgSync.userGroups',
    ],
    [
      'an empty scope while the token is also missing',
      enabled({ userGroups: [] }, { token: undefined }),
      'authentik.orgSync.userGroups',
    ],
    [
      'a blank group name',
      enabled({ userGroups: ['  '] }),
      'authentik.orgSync.userGroups',
    ],
    [
      'a zero frequency',
      enabled({ frequencySeconds: 0 }),
      'frequencySeconds must be a positive number',
    ],
    [
      'a negative timeout',
      enabled({ timeoutSeconds: -5 }),
      'timeoutSeconds must be a positive number',
    ],
    [
      'no org.yaml catalog location and no curatedFile',
      { ...enabled(), catalog: { locations: [] } },
      'authentik.orgSync.curatedFile is required',
    ],
  ])(
    'fails startup on %s, naming the key and never the token',
    async (_name, data, message) => {
      const error = await make(new ConfigReader(data)).then(
        () => undefined,
        (e: unknown) => e,
      );
      expect(error).toBeInstanceOf(Error);
      expect((error as Error).message).toContain(message);
      expect((error as Error).message).not.toContain(TOKEN);
    },
  );

  it.each([
    ['disabled', {}],
    ['enabled', enabled()],
  ])(
    'fails startup when %s and the scheduler cannot store the task, instead of serving without it',
    async (_name, data) => {
      const scheduleTask = jest
        .fn()
        .mockRejectedValue(new Error('Failed to persist task, database blip'));
      await expect(make(new ConfigReader(data), scheduleTask)).rejects.toThrow(
        'Failed to persist task, database blip',
      );
      expect(scheduleTask).toHaveBeenCalledTimes(1);
    },
  );
});

describe('the reviewed org file', () => {
  const files = (...targets: string[]) =>
    targets.map(target => ({ type: 'file', target }));

  it('is the one file location named org.yaml, unless curatedFile names one', () => {
    const locations = [
      { type: 'url', target: 'https://example.test/catalog/org.yaml' },
      ...files('/app/catalog/cluster.yaml', '/app/catalog/org.yaml'),
    ];
    expect(reviewedOrgFile(new ConfigReader({ catalog: { locations } }))).toBe(
      '/app/catalog/org.yaml',
    );
    expect(
      reviewedOrgFile(
        new ConfigReader({
          catalog: { locations },
          authentik: { orgSync: { curatedFile: ' /elsewhere/org.yaml ' } },
        }),
      ),
    ).toBe('/elsewhere/org.yaml');
  });

  it.each([
    ['no catalog locations', {}],
    [
      'no org.yaml file location',
      { catalog: { locations: files('/app/catalog/cluster.yaml') } },
    ],
    [
      'two org.yaml file locations',
      { catalog: { locations: files('/a/org.yaml', '/b/org.yaml') } },
    ],
  ])('must be named explicitly with %s', (_name, data) => {
    expect(() => reviewedOrgFile(new ConfigReader(data))).toThrow(
      'authentik.orgSync.curatedFile is required',
    );
  });
});

const load = async (files: string[], env: Record<string, string> = {}) => {
  const result = await loadConfig({
    configRoot: portalPath('.'),
    configTargets: files.map(path => ({ path })),
    experimentalEnvFunc: async name => env[name],
  });
  return ConfigReader.fromConfigs(result.appConfigs);
};
const base = portalPath('app-config.yaml');
const production = portalPath('app-config.production.yaml');

describe('the committed portal configs', () => {
  it('read the reviewed org file the catalog reads, under yarn start and in the image', async () => {
    // `yarn start` runs the backend from packages/backend.
    const local = reviewedOrgFile(await load([base]));
    expect(local).toBe('../../catalog/org.yaml');
    expect(resolve(portalPath('packages/backend'), local)).toBe(orgYamlPath);
    // The image runs from /app, where the backstage-catalog ConfigMap is mounted at /app/catalog.
    expect(reviewedOrgFile(await load([base, production]))).toBe(
      '/app/catalog/org.yaml',
    );
  });
});

const overlay = portalPath('../../platform/backstage/app-config.cluster.yaml');
// The image build context holds only services/onemind-portal.
const describeRepo = existsSync(overlay) ? describe : describe.skip;

describeRepo('the committed cluster overlay', () => {
  const cluster = (env: Record<string, string>) =>
    load([base, production, overlay], env);

  it('ships the sync disabled: merging only registers the removal, which never calls Authentik, with or without the token', async () => {
    const envs: Record<string, string>[] = [
      {},
      { AUTHENTIK_CATALOG_TOKEN: TOKEN },
    ];
    for (const env of envs) {
      const m = await make(await cluster(env));
      expect(m.mode).toBe('cleanup');
      expect(registeredTask(m)).toEqual(TASK);
      expect(await connectAndRunOnce(m)).toEqual([
        EMPTY_FULL_MUTATION,
        EMPTY_FULL_MUTATION,
      ]);
      expect(m.fetchJson).not.toHaveBeenCalled();
      expect(m.logger.error).not.toHaveBeenCalled();
    }
  });

  it('takes the token only from AUTHENTIK_CATALOG_TOKEN and reads the in-cluster Authentik service', async () => {
    expect(
      (await cluster({})).getOptionalString('authentik.token'),
    ).toBeUndefined();
    const config = await cluster({ AUTHENTIK_CATALOG_TOKEN: TOKEN });
    expect(config.getString('authentik.token')).toBe(TOKEN);
    expect(config.getString('authentik.baseUrl')).toBe(
      'http://platform-auth-helm-authentik-server.auth.svc.cluster.local',
    );
  });

  it('keeps the production catalog locations, so the reviewed org file is the mounted one', async () => {
    expect(reviewedOrgFile(await cluster({}))).toBe('/app/catalog/org.yaml');
  });
});
