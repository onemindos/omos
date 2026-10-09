export const platforms = [
  'Web',
  'iOS',
  'macOS',
  'Android',
  'Linux',
  'Windows',
  'Other',
];
export const methods = [
  'Browser',
  'Native app',
  'Docker',
  'Kubernetes',
  'Other',
];
export type AccessPoint = {
  id: string;
  label: string;
  platform: string;
  method: string;
  location: string;
  targetRef?: string;
  url?: string;
  instructions?: string;
  observedAt?: string;
  evidence:
    | 'Catalog snapshot'
    | 'Installation inspected'
    | 'User recorded'
    | 'Link only';
  version?: string;
};
export function safeAccessUrl(value?: string) {
  if (!value) return undefined;
  try {
    const u = new URL(value);
    return u.protocol === 'https:' && !u.username && !u.password
      ? u.href
      : undefined;
  } catch {
    return undefined;
  }
}
export function parseAccessPoints(raw?: string): AccessPoint[] {
  try {
    const value = JSON.parse(raw ?? '[]');
    if (!Array.isArray(value)) return [];
    return value
      .slice(0, 50)
      .filter(
        x =>
          x &&
          typeof x.id === 'string' &&
          typeof x.label === 'string' &&
          typeof x.platform === 'string' &&
          typeof x.method === 'string' &&
          typeof x.location === 'string' &&
          [
            'Catalog snapshot',
            'Installation inspected',
            'User recorded',
            'Link only',
          ].includes(x.evidence),
      )
      .map(x => ({
        id: x.id,
        label: x.label,
        platform: x.platform,
        method: x.method,
        location: x.location,
        evidence: x.evidence,
        url: typeof x.url === 'string' ? safeAccessUrl(x.url) : undefined,
        ...Object.fromEntries(
          ['targetRef', 'instructions', 'observedAt', 'version']
            .filter(k => typeof x[k] === 'string')
            .map(k => [k, x[k]]),
        ),
      }));
  } catch {
    return [];
  }
}
export function catalogPath(ref?: string) {
  const m = ref?.match(
    /^(component|resource|system):([a-z0-9_-]+)\/([a-z0-9_.-]+)$/i,
  );
  return m ? `/catalog/${m[2]}/${m[1].toLowerCase()}/${m[3]}` : undefined;
}
