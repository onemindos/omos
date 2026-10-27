import { existsSync } from 'node:fs';
import type { JsonObject } from '@backstage/types';
import { ConfigReader, type Config } from '@backstage/config';
import { loadConfig } from '@backstage/config-loader';
import { mockServices } from '@backstage/backend-test-utils';
import { readAuthentikOrgConfig } from './config';
import { registerAuthentikOrg } from './module';
import { portalPath } from './testing';

const TOKEN = 'fixture-only-token';

const enabled = (
  orgSync: JsonObject = {},
  authentik: JsonObject = {},
): JsonObject => ({
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

const make = (config: Config) => {
  const processing = { addEntityProvider: jest.fn() };
  const scheduler = { scheduleTask: jest.fn() };
  const fetchJson = jest.fn();
  const logger = mockServices.logger.mock();
  const registered = registerAuthentikOrg({
    config,
    logger,
    scheduler: scheduler as never,
    processing,
    fetchJson,
  });
  return { registered, processing, scheduler, fetchJson, logger };
};

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
    'registers nothing and calls nothing when the flag is %s',
    (_name, data) => {
      const m = make(new ConfigReader(data));
      expect(m.registered).toBe(false);
      expect(m.processing.addEntityProvider).not.toHaveBeenCalled();
      expect(m.scheduler.scheduleTask).not.toHaveBeenCalled();
      expect(m.fetchJson).not.toHaveBeenCalled();
    },
  );

  it('registers one provider and one global task when enabled, and reads nothing until it runs', () => {
    const m = make(new ConfigReader(enabled()));
    expect(m.registered).toBe(true);
    expect(m.processing.addEntityProvider).toHaveBeenCalledTimes(1);
    expect(
      m.processing.addEntityProvider.mock.calls[0][0].getProviderName(),
    ).toBe('onemind-authentik-org');
    expect(m.scheduler.scheduleTask).toHaveBeenCalledTimes(1);
    const { fn, ...task } = m.scheduler.scheduleTask.mock.calls[0][0];
    expect(typeof fn).toBe('function');
    expect(task).toEqual({
      id: 'onemind-authentik-org-refresh',
      frequency: { seconds: 600 },
      timeout: { seconds: 120 },
      initialDelay: { seconds: 60 },
      scope: 'global',
    });
    expect(m.fetchJson).not.toHaveBeenCalled();
    const logged = JSON.stringify(m.logger.info.mock.calls);
    expect(logged).toContain(
      'members of OneMind Portal Users every 600s from http://authentik.test',
    );
    expect(logged).not.toContain(TOKEN);
  });

  it('normalises the URL and de-duplicates the scope', () => {
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
      userGroups: ['A', 'B'],
      frequencySeconds: 30,
      curatedFile: 'catalog/org.yaml',
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
      'no token (env var absent)',
      enabled({}, { token: undefined }),
      'authentik.token is required',
    ],
    [
      'a blank token',
      enabled({}, { token: '   ' }),
      'authentik.token is required',
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
  ])(
    'fails startup on %s, naming the key and never the token',
    (_name, data, message) => {
      let error: unknown;
      try {
        make(new ConfigReader(data));
      } catch (e) {
        error = e;
      }
      expect(error).toBeInstanceOf(Error);
      expect((error as Error).message).toContain(message);
      expect((error as Error).message).not.toContain(TOKEN);
    },
  );
});

const overlay = portalPath('../../platform/backstage/app-config.cluster.yaml');
// The image build context holds only services/onemind-portal.
const describeRepo = existsSync(overlay) ? describe : describe.skip;

describeRepo('the committed cluster overlay', () => {
  const load = async (env: Record<string, string>) => {
    const result = await loadConfig({
      configRoot: portalPath('.'),
      configTargets: [
        { path: portalPath('app-config.yaml') },
        { path: portalPath('app-config.production.yaml') },
        { path: overlay },
      ],
      experimentalEnvFunc: async name => env[name],
    });
    return ConfigReader.fromConfigs(result.appConfigs);
  };

  it('ships the sync disabled: merging registers nothing, with or without the token', async () => {
    const envs: Record<string, string>[] = [
      {},
      { AUTHENTIK_CATALOG_TOKEN: TOKEN },
    ];
    for (const env of envs) {
      const m = make(await load(env));
      expect(m.registered).toBe(false);
      expect(m.processing.addEntityProvider).not.toHaveBeenCalled();
      expect(m.scheduler.scheduleTask).not.toHaveBeenCalled();
    }
  });

  it('takes the token only from AUTHENTIK_CATALOG_TOKEN and reads the in-cluster Authentik service', async () => {
    expect(
      (await load({})).getOptionalString('authentik.token'),
    ).toBeUndefined();
    const config = await load({ AUTHENTIK_CATALOG_TOKEN: TOKEN });
    expect(config.getString('authentik.token')).toBe(TOKEN);
    expect(config.getString('authentik.baseUrl')).toBe(
      'http://platform-auth-helm-authentik-server.auth.svc.cluster.local',
    );
  });
});
