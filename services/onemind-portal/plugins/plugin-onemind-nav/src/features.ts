import { configApiRef, useApi } from '@backstage/frontend-plugin-api';

/** Flags read from `onemind.features.*` in app-config.yaml. Missing means off. */
export type NavFeature = 'pillars' | 'digitalocean' | 'retireApp';

export type NavFlags = Record<NavFeature, boolean>;

export const NAV_FEATURES: NavFeature[] = [
  'pillars',
  'digitalocean',
  'retireApp',
];

export function useNavFlags(): NavFlags {
  const config = useApi(configApiRef);
  return Object.fromEntries(
    NAV_FEATURES.map(name => [
      name,
      config.getOptionalBoolean(`onemind.features.${name}`) ?? false,
    ]),
  ) as NavFlags;
}
