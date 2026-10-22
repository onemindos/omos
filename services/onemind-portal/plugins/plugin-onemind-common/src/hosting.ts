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
  if (status === 'Planned' || hosting === 'planned')
    return {
      group: 'Planned',
      label: 'Planned',
      detail:
        'Recorded for future use; deployment location has not been verified.',
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
  const onCluster = !!(a['argocd/app-name'] || a['onemind.io/kubernetes-workloads'] || a['backstage.io/kubernetes-label-selector'] || a['backstage.io/kubernetes-id']);
  if (a['onemind.io/hosting'] === 'digitalocean-nyc3') return { badge: onCluster ? 'K8s' : 'DO', detail: info.detail };
  if (a['onemind.io/hosting'] === 'external-provider-unverified') return { badge: 'SaaS', detail: info.detail };
  if (a['onemind.io/hosting'] === 'mesh-endpoint-unverified') return { badge: 'Home', detail: info.detail };
  if (a['onemind.io/hosting'] === 'native-application' || a['onemind.io/hosting'] === 'vehicle-firmware') return { badge: 'Device', detail: info.detail };
  if (a['onemind.io/hosting'] === 'external-website') return { badge: 'Web', detail: info.detail };
  return { badge: 'Unknown', detail: info.detail };
}
