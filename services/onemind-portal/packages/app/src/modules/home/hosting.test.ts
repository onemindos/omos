import { hostingInfo } from './hosting';
test('distinguishes self-managed cloud deployments from provider services', () => {
  expect(hostingInfo('digitalocean-nyc3').label).toBe('Self-hosted · cloud');
  expect(hostingInfo('external-provider-unverified').group).toBe(
    'Cloud service',
  );
});
test('does not promote planned, mesh or unrecorded entries to deployed apps', () => {
  expect(hostingInfo('vehicle-firmware', 'Planned').group).toBe('Planned');
  expect(hostingInfo('planned').group).toBe('Planned');
  expect(hostingInfo('mesh-endpoint-unverified').group).toBe('Unknown');
  expect(hostingInfo().group).toBe('Unknown');
  expect(hostingInfo('external-website').group).toBe('Reference');
});
