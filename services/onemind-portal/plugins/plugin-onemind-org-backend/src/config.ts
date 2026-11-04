import { basename } from 'node:path';
import type { Config } from '@backstage/config';

/** Default schedule of the sync, and the fixed schedule of the removal that replaces it when off. */
export const DEFAULT_SCHEDULE = {
  frequencySeconds: 600,
  timeoutSeconds: 120,
  initialDelaySeconds: 60,
} as const;

/** Settings of an enabled sync, from `authentik.*` and `authentik.orgSync.*`. */
export interface AuthentikOrgSettings {
  enabled: true;
  /** Without a trailing slash. */
  baseUrl: string;
  /**
   * Absent when `authentik.token` is unset or blank, for example because its secret is missing. The
   * module then removes what it published instead of failing startup.
   */
  token?: string;
  /** Exact Authentik group names whose members become catalog Users. */
  userGroups: string[];
  frequencySeconds: number;
  timeoutSeconds: number;
  initialDelaySeconds: number;
  /** Reviewed catalog file whose Users and Groups always win (never published or linked). */
  curatedFile: string;
}

/**
 * The reviewed org file: `authentik.orgSync.curatedFile` when set, otherwise the one `file` location
 * in `catalog.locations` named org.yaml. A relative path resolves against the working directory, as
 * the catalog resolves that location, so the module reads the file the catalog reads: in the image
 * /app/catalog/org.yaml, under `yarn start` ../../catalog/org.yaml from packages/backend.
 */
export function reviewedOrgFile(config: Config): string {
  const explicit = config
    .getOptionalString('authentik.orgSync.curatedFile')
    ?.trim();
  if (explicit) return explicit;
  const targets = (config.getOptionalConfigArray('catalog.locations') ?? [])
    .filter(location => location.getOptionalString('type') === 'file')
    .map(location => location.getOptionalString('target') ?? '')
    .filter(target => basename(target) === 'org.yaml');
  if (targets.length !== 1)
    throw new Error(
      'authentik.orgSync.curatedFile is required when catalog.locations does not have exactly one file location named org.yaml',
    );
  return targets[0];
}

/**
 * The sync is off unless `authentik.orgSync.enabled` is true. While it is off nothing else is read,
 * so an absent token or URL cannot affect startup. When it is on, missing or nonsensical reviewed
 * settings throw at startup instead of running with them; messages name keys, never values. The
 * token is the exception: it comes from a secret that can disappear at runtime, so an absent token
 * is returned as such and never fails startup.
 */
export function readAuthentikOrgConfig(
  config: Config,
): { enabled: false } | AuthentikOrgSettings {
  const sync = config.getOptionalConfig('authentik.orgSync');
  if (!sync || !(sync.getOptionalBoolean('enabled') ?? false))
    return { enabled: false };

  const required = 'is required when authentik.orgSync.enabled is true';
  const rawUrl = config.getOptionalString('authentik.baseUrl')?.trim();
  if (!rawUrl) throw new Error(`authentik.baseUrl ${required}`);
  let url: URL;
  try {
    url = new URL(rawUrl);
  } catch {
    throw new Error('authentik.baseUrl must be an http(s) URL');
  }
  if (
    !['http:', 'https:'].includes(url.protocol) ||
    url.username ||
    url.password ||
    url.search ||
    url.hash
  )
    throw new Error(
      'authentik.baseUrl must be an http(s) URL without credentials, query or fragment',
    );

  const token = config.getOptionalString('authentik.token')?.trim();

  const userGroups = [
    ...new Set(
      (sync.getOptionalStringArray('userGroups') ?? []).map(g => g.trim()),
    ),
  ];
  if (!userGroups.length || userGroups.some(g => !g))
    throw new Error(
      `authentik.orgSync.userGroups (non-empty Authentik group names) ${required}`,
    );

  const num = (key: string, fallback: number) => {
    const value = sync.getOptionalNumber(key) ?? fallback;
    if (!Number.isFinite(value) || value <= 0)
      throw new Error(`authentik.orgSync.${key} must be a positive number`);
    return value;
  };
  return {
    enabled: true,
    baseUrl: rawUrl.replace(/\/+$/, ''),
    ...(token ? { token } : {}),
    userGroups,
    frequencySeconds: num(
      'frequencySeconds',
      DEFAULT_SCHEDULE.frequencySeconds,
    ),
    timeoutSeconds: num('timeoutSeconds', DEFAULT_SCHEDULE.timeoutSeconds),
    initialDelaySeconds: num(
      'initialDelaySeconds',
      DEFAULT_SCHEDULE.initialDelaySeconds,
    ),
    curatedFile: reviewedOrgFile(config),
  };
}
