import { mockServices, startTestBackend } from '@backstage/backend-test-utils';
import { policyExtensionPoint } from '@backstage/plugin-permission-node/alpha';
import { AuthorizeResult } from '@backstage/plugin-permission-common';
import { catalogEntityReadPermission } from '@backstage/plugin-catalog-common/alpha';
import module from './plugin';

test('installs the OneMind policy without any OpenChoreo configuration or services', async () => {
  const setPolicy = jest.fn();
  await startTestBackend({
    extensionPoints: [[policyExtensionPoint, { setPolicy }]],
    features: [
      module,
      // No `openchoreo.*` block at all: the module must not require one.
      mockServices.rootConfig.factory({
        data: { permission: { enabled: true } },
      }),
    ],
  });
  expect(setPolicy).toHaveBeenCalledTimes(1);
  const policy = setPolicy.mock.calls[0][0];
  expect(
    await policy.handle({ permission: catalogEntityReadPermission }, undefined),
  ).toEqual({ result: AuthorizeResult.DENY });
});
