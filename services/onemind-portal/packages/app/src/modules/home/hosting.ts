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
