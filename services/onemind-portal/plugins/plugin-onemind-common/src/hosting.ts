import hostingValues from './hostingValues.json';

/** Every value `onemind.io/hosting` may carry: the canonical set plus the legacy values that stay valid. */
export const HOSTING_VALUES: readonly string[] = [
  ...Object.keys(hostingValues.canonical),
  ...Object.keys(hostingValues.legacy),
];

/** Where an app runs, as one word per kind. Used by the hosting badge, the facts card and the Not on Kubernetes view. */
export type HostingKind =
  | 'kubernetes'
  | 'digitalocean'
  | 'saas'
  | 'home'
  | 'device'
  | 'desktop'
  | 'mobile'
  | 'web'
  | 'planned'
  | 'unknown';

export const HOSTING_KIND_LABEL: Record<HostingKind, string> = {
  kubernetes: 'Kubernetes',
  digitalocean: 'DigitalOcean (not Kubernetes)',
  saas: 'SaaS',
  home: 'Home / edge',
  device: 'Device',
  desktop: 'Desktop app',
  mobile: 'Mobile app',
  web: 'Reference website',
  planned: 'Planned',
  unknown: 'Unknown',
};

type HostingEntity = { metadata: { annotations?: Record<string, string> } };
const onCluster = (a: Record<string, string>) =>
  !!(
    a['argocd/app-name'] ||
    a['onemind.io/kubernetes-workloads'] ||
    a['backstage.io/kubernetes-label-selector'] ||
    a['backstage.io/kubernetes-id']
  );

/**
 * The effective hosting kind of an entity. Canonical values map to themselves; the legacy values
 * map as documented in docs/catalog/app-contract.md. A missing or unrecognised value is `unknown`.
 */
export function hostingKind(entity: HostingEntity): HostingKind {
  const a = entity.metadata.annotations ?? {};
  const value = a['onemind.io/hosting'];
  if (a['onemind.io/documented-status'] === 'Planned') return 'planned';
  switch (value) {
    case 'digitalocean-nyc3':
      return onCluster(a) ? 'kubernetes' : 'digitalocean';
    case 'external-provider-unverified':
      return 'saas';
    case 'mesh-endpoint-unverified':
      return 'home';
    case 'native-application':
      return 'desktop';
    case 'vehicle-firmware':
      return 'device';
    case 'external-website':
      return 'web';
    case 'planned-not-deployed':
      return 'planned';
    default:
      return value && value in hostingValues.canonical
        ? (value as HostingKind)
        : 'unknown';
  }
}

export type HostingGroup =
  | 'Self-hosted'
  | 'Cloud service'
  | 'Local / native'
  | 'Reference'
  | 'Planned'
  | 'Unknown';
export function hostingInfo(
  hosting?: string,
  status?: string,
): { group: HostingGroup; label: string; detail: string } {
  if (status === 'Planned' || hosting === 'planned' || hosting === 'planned-not-deployed')
    return {
      group: 'Planned',
      label: 'Planned',
      detail:
        'Recorded for future use; deployment location has not been verified.',
    };
  if (hosting === 'kubernetes' || hosting === 'digitalocean')
    return {
      group: 'Self-hosted',
      label:
        hosting === 'kubernetes'
          ? 'Self-hosted · Kubernetes'
          : 'Self-hosted · DigitalOcean',
      detail:
        'OneMind-managed deployment on DigitalOcean. Hosting inventory does not establish current availability.',
    };
  if (hosting === 'saas')
    return {
      group: 'Cloud service',
      label: 'SaaS',
      detail:
        'Provider-hosted service. Health is a public reachability check, not proof the account works.',
    };
  if (hosting === 'home')
    return {
      group: 'Self-hosted',
      label: 'Self-hosted · home / edge',
      detail: 'Runs at home or at the edge; reachability depends on the tailnet.',
    };
  if (hosting === 'device' || hosting === 'desktop' || hosting === 'mobile') {
    const labels = { device: 'Device', desktop: 'Desktop app', mobile: 'Mobile app' };
    return {
      group: 'Local / native',
      label: labels[hosting],
      detail:
        'Recorded application type; installation and active use have not been verified.',
    };
  }
  if (hosting === 'web')
    return {
      group: 'Reference',
      label: 'Website reference',
      detail: 'External resource link; not a deployed OneMind application.',
    };
  if (hosting === 'digitalocean-nyc3')
    return {
      group: 'Self-hosted',
      label: 'Self-hosted · cloud',
      detail:
        'OneMind-managed deployment on DigitalOcean NYC3. Hosting inventory does not establish current availability.',
    };
  if (hosting === 'external-provider-unverified')
    return {
      group: 'Cloud service',
      label: 'Cloud service',
      detail:
        'External provider. Account access and availability have not been verified.',
    };
  if (hosting === 'native-application' || hosting === 'vehicle-firmware')
    return {
      group: 'Local / native',
      label: hosting === 'vehicle-firmware' ? 'Vehicle firmware' : 'Native app',
      detail:
        'Recorded application type; installation and active use have not been verified.',
    };
  if (hosting === 'external-website')
    return {
      group: 'Reference',
      label: 'Website reference',
      detail: 'External resource link; not a deployed OneMind application.',
    };
  if (hosting === 'mesh-endpoint-unverified')
    return {
      group: 'Unknown',
      label: 'Mesh · unverified',
      detail: 'Mesh endpoint recorded; hosting location has not been verified.',
    };
  return {
    group: 'Unknown',
    label: 'Hosting unknown',
    detail: 'Hosting has not been recorded or verified.',
  };
}

export type HostingBadge = 'K8s' | 'DO' | 'SaaS' | 'Home' | 'Device' | 'Web' | 'Planned' | 'Unknown';
type Annotated = { metadata: { annotations?: Record<string, string> } };

/**
 * The hosting badge in the standard app header: one short word, the same set on every app page.
 * K8s means it runs on the cluster (Argo application, workloads or a Kubernetes selector);
 * DO is other DigitalOcean infrastructure. Nothing is promoted: unverified entries keep "Unknown".
 */
export function hostingBadge(entity: Annotated): { badge: HostingBadge; detail: string } {
  const a = entity.metadata.annotations ?? {};
  const info = hostingInfo(a['onemind.io/hosting'], a['onemind.io/documented-status']);
  if (info.group === 'Planned') return { badge: 'Planned', detail: info.detail };
  const kind = hostingKind(entity);
  if (kind === 'planned') return { badge: 'Planned', detail: hostingInfo('planned').detail };
  const badges: Partial<Record<HostingKind, HostingBadge>> = {
    kubernetes: 'K8s',
    digitalocean: 'DO',
    saas: 'SaaS',
    home: 'Home',
    device: 'Device',
    desktop: 'Device',
    mobile: 'Device',
    web: 'Web',
  };
  return { badge: badges[kind] ?? 'Unknown', detail: info.detail };
}
