import type { JsonObject } from '@backstage/types';
import { ConfigReader } from '@backstage/config';
import { mockServices } from '@backstage/backend-test-utils';
import { registerLiveInventory } from './module';

const make = (data: JsonObject) => {
  const processing = { addEntityProvider: jest.fn(), addProcessor: jest.fn() };
  const scheduler = { scheduleTask: jest.fn() };
  const catalog = { getEntities: jest.fn() };
  const auth = { getOwnServiceCredentials: jest.fn() };
  const logger = mockServices.logger.mock();
  const registered = registerLiveInventory({
    config: new ConfigReader(data),
    logger,
    scheduler: scheduler as never,
    auth: auth as never,
    catalog: catalog as never,
    processing,
  });
  return { registered, processing, scheduler, catalog, auth, logger };
};

describe('feature flag', () => {
  it.each([
    ['absent', {}],
    ['false', { onemind: { inventory: { live: { enabled: false } } } }],
    [
      'unset with other settings',
      { onemind: { inventory: { live: { frequencySeconds: 30 } } } },
    ],
  ])(
    'registers nothing and makes no cluster or catalog call when the flag is %s',
    (_name, data) => {
      const { registered, processing, scheduler, catalog, auth } = make(data);
      expect(registered).toBe(false);
      expect(processing.addEntityProvider).not.toHaveBeenCalled();
      expect(processing.addProcessor).not.toHaveBeenCalled();
      expect(scheduler.scheduleTask).not.toHaveBeenCalled();
      expect(catalog.getEntities).not.toHaveBeenCalled();
      expect(auth.getOwnServiceCredentials).not.toHaveBeenCalled();
    },
  );

  it('registers one provider, one processor and two tasks (observe locally, publish once) when enabled', () => {
    const { registered, processing, scheduler } = make({
      onemind: { inventory: { live: { enabled: true } } },
    });
    expect(registered).toBe(true);
    expect(processing.addEntityProvider).toHaveBeenCalledTimes(1);
    expect(processing.addProcessor).toHaveBeenCalledTimes(1);
    const tasks = scheduler.scheduleTask.mock.calls.map(c => [
      c[0].id,
      c[0].scope,
      c[0].frequency,
    ]);
    expect(tasks).toEqual([
      ['onemind-live-inventory-observe', 'local', { seconds: 120 }],
      ['onemind-live-inventory-publish', 'global', { seconds: 120 }],
    ]);
  });

  it('rejects nonsensical settings at startup instead of running with them', () => {
    expect(() =>
      make({
        onemind: {
          inventory: { live: { enabled: true, frequencySeconds: 0 } },
        },
      }),
    ).toThrow(/frequencySeconds/);
  });
});
