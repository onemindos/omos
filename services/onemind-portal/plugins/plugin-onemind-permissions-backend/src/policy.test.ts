import {
  AuthorizeResult,
  createPermission,
} from '@backstage/plugin-permission-common';
import type { PolicyQueryUser } from '@backstage/plugin-permission-node';
import { createOneMindPolicy } from './policy';

const user = (admin: boolean): PolicyQueryUser => ({
  credentials: {
    $$type: '@backstage/BackstageCredentials',
    principal: { type: 'user', userEntityRef: 'user:default/operator' },
  },
  info: {
    userEntityRef: 'user:default/operator',
    ownershipEntityRefs: admin
      ? ['group:default/onemind-admins']
      : ['user:default/operator'],
  },
});
const query = (name: string, action: 'read' | 'create') => ({
  permission: createPermission({ name, attributes: { action } }),
});

describe('OneMind + OpenChoreo policy intersection', () => {
  test('anonymous denied, including reads', async () => {
    expect(
      await createOneMindPolicy().handle(query('catalog.entity.read', 'read')),
    ).toEqual({ result: AuthorizeResult.DENY });
  });
  test('ordinary user reads allowed for existing host permissions', async () => {
    expect(
      await createOneMindPolicy().handle(
        query('onemind.read', 'read'),
        user(false),
      ),
    ).toEqual({ result: AuthorizeResult.ALLOW });
  });
  test.each([
    'openchoreo.component.create',
    'catalog.entity.create',
    'scaffolder.task.create',
    'onemind.write',
  ])('non-admin write denied: %s', async name => {
    const handle = jest
      .fn()
      .mockResolvedValue({ result: AuthorizeResult.ALLOW });
    expect(
      await createOneMindPolicy({ handle }).handle(
        query(name, 'create'),
        user(false),
      ),
    ).toEqual({ result: AuthorizeResult.DENY });
    expect(handle).not.toHaveBeenCalled();
  });
  test.each([
    'openchoreo.component.create',
    'catalog.entity.read',
    'scaffolder.task.create',
  ])('administrator cannot override upstream denial: %s', async name => {
    const handle = jest
      .fn()
      .mockResolvedValue({ result: AuthorizeResult.DENY });
    expect(
      await createOneMindPolicy({ handle }).handle(
        query(name, 'create'),
        user(true),
      ),
    ).toEqual({ result: AuthorizeResult.DENY });
    expect(handle).toHaveBeenCalled();
  });
  test('conditional decisions are preserved, never promoted to ALLOW', async () => {
    const decision = {
      result: AuthorizeResult.CONDITIONAL,
      pluginId: 'catalog',
      resourceType: 'catalog-entity',
      conditions: {
        rule: 'MATCHES_CAPABILITY',
        resourceType: 'catalog-entity',
        params: {},
      },
    } as const;
    expect(
      await createOneMindPolicy({ handle: async () => decision }).handle(
        query('catalog.entity.read', 'read'),
        user(true),
      ),
    ).toBe(decision);
  });
  test('upstream error fails closed', async () => {
    expect(
      await createOneMindPolicy({
        handle: async () => {
          throw Error('unavailable');
        },
      }).handle(query('openchoreo.component.read', 'read'), user(true)),
    ).toEqual({ result: AuthorizeResult.DENY });
  });
  test('host administrator retains unrelated host operations', async () => {
    const handle = jest.fn();
    expect(
      await createOneMindPolicy({ handle }).handle(
        query('onemind.write', 'create'),
        user(true),
      ),
    ).toEqual({ result: AuthorizeResult.ALLOW });
    expect(handle).not.toHaveBeenCalled();
  });
});
