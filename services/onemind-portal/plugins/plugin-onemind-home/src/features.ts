import { configApiRef, useApi } from '@backstage/core-plugin-api';

/**
 * Flags live in app-config.yaml under `onemind.features` (schema in
 * plugin-onemind-nav/config.d.ts). All default to off until the work order
 * that builds the target (F2 Domains, F6 DigitalOcean, F8 Retire) merges.
 */
export type HomeFeature = 'pillars';

export function useHomeFeature(name: HomeFeature): boolean {
  const config = useApi(configApiRef);
  return config.getOptionalBoolean(`onemind.features.${name}`) ?? false;
}
