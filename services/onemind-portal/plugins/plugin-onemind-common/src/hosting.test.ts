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

import { hostingBadge } from './hosting';
const e = (annotations: Record<string, string>) => ({ metadata: { annotations } });
test('hosting badge: cluster apps are K8s, other DigitalOcean entries DO, and the rest keep their own word', () => {
  expect(hostingBadge(e({ 'onemind.io/hosting': 'digitalocean-nyc3', 'argocd/app-name': 'platform-x' })).badge).toBe('K8s');
  expect(hostingBadge(e({ 'onemind.io/hosting': 'digitalocean-nyc3', 'onemind.io/kubernetes-workloads': 'a/b' })).badge).toBe('K8s');
  expect(hostingBadge(e({ 'onemind.io/hosting': 'digitalocean-nyc3' })).badge).toBe('DO');
  expect(hostingBadge(e({ 'onemind.io/hosting': 'external-provider-unverified' })).badge).toBe('SaaS');
  expect(hostingBadge(e({ 'onemind.io/hosting': 'mesh-endpoint-unverified' })).badge).toBe('Home');
  expect(hostingBadge(e({ 'onemind.io/hosting': 'native-application' })).badge).toBe('Device');
  expect(hostingBadge(e({ 'onemind.io/hosting': 'vehicle-firmware' })).badge).toBe('Device');
  expect(hostingBadge(e({ 'onemind.io/hosting': 'external-website' })).badge).toBe('Web');
});
test('hosting badge: planned and unrecorded entries are never promoted', () => {
  expect(hostingBadge(e({ 'onemind.io/hosting': 'digitalocean-nyc3', 'onemind.io/documented-status': 'Planned' })).badge).toBe('Planned');
  expect(hostingBadge(e({})).badge).toBe('Unknown');
});
