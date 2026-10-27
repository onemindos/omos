import type { Config } from '@backstage/config';

/** Settings of an enabled sync, from `authentik.*` and `authentik.orgSync.*`. */
export interface AuthentikOrgSettings {
  enabled: true;
  /** Without a trailing slash. */
  baseUrl: string;
  token: string;
  /** Exact Authentik group names whose members become catalog Users. */
  userGroups: string[];
  frequencySeconds: number;
  timeoutSeconds: number;
  initialDelaySeconds: number;
  /** Reviewed catalog file whose Users and Groups always win (never published or linked). */
  curatedFile: string;
}

/**
 * The sync is off unless `authentik.orgSync.enabled` is true. While it is off nothing else is read,
 * so an absent token or URL cannot affect startup. When it is on, missing or nonsensical settings
 * throw at startup instead of running with them; messages name keys, never values.
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
  if (!token) throw new Error(`authentik.token ${required}`);

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
    token,
    userGroups,
    frequencySeconds: num('frequencySeconds', 600),
    timeoutSeconds: num('timeoutSeconds', 120),
    initialDelaySeconds: num('initialDelaySeconds', 60),
    curatedFile: sync.getOptionalString('curatedFile') ?? 'catalog/org.yaml',
  };
}
