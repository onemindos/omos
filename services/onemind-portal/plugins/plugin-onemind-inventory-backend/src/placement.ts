import type { InventoryMapping } from './mapping.generated';

/*
 * Placement rules, ported from scripts/curate-structure.cjs (pretty, words, workloadRole, ruleMatches,
 * appForWorkload). Kept as a port, not an import, because the script is not part of the backend bundle;
 * placement.test.ts runs both over every recorded workload so they cannot drift.
 */
const ACRONYMS = new Set(['api', 'db', 'ui', 'ldap', 'dns', 'nats', 'tak', 'pg', 'ml', 'id', 'tls', 'ha', 'ws', 'oc', 'ai', 'cpu', 'crd', 'crds', 'tcp', 'url', 'ip', 'sql', 'mcp', 'stac', 'osm']); // prettier-ignore
const SPECIAL: Record<string, string> = {
  github: 'GitHub',
  coredns: 'CoreDNS',
  applicationset: 'ApplicationSet',
  apiserver: 'API Server',
  kube: 'Kube',
};

function pretty(token: string): string {
  const t = token.toLowerCase();
  if (SPECIAL[t]) return SPECIAL[t];
  return ACRONYMS.has(t)
    ? t.toUpperCase()
    : t.charAt(0).toUpperCase() + t.slice(1);
}

export const words = (name: string): string =>
  name
    .split(/[-_.]+/)
    .filter(t => t && !/^[0-9a-f]{8}$/.test(t))
    .map(pretty)
    .join(' ');

/** Human role of a workload inside its app: drop everything up to and including the app's name. */
export function workloadRole(
  workload: string,
  appName: string,
  namespace = '',
): { role: string; suffix?: string } {
  const tokens = workload.split('-');
  const proxy = workload.match(/^ts-(.+)-([a-z0-9]{5})$/);
  if (proxy) return { role: `Proxy: ${proxy[1]}`, suffix: proxy[2] };
  const app = appName.split('-');
  for (let i = 0; i + app.length <= tokens.length; i++) {
    if (app.every((a, j) => tokens[i + j] === a)) {
      const rest = tokens.slice(i + app.length);
      return { role: rest.length ? words(rest.join('-')) : 'Core' };
    }
  }
  const ns = namespace.split('-');
  const stripped =
    namespace &&
    tokens.length > ns.length &&
    ns.every((n, i) => tokens[i] === n)
      ? tokens.slice(ns.length)
      : tokens;
  return { role: words(stripped.join('-')) };
}

function ruleMatches(pattern: string, key: string): boolean {
  const [pn, pw] = pattern.split('/');
  const [kn, kw] = key.split('/');
  if (pn !== kn) return false;
  if (pw === '*') return true;
  if (pw.endsWith('*')) return kw.startsWith(pw.slice(0, -1));
  return pw === kw;
}

/** Resolve "<namespace>/<workload>" to an app ref: exact entry first, then the longest "ns/prefix*" or "ns/*". */
export function appForWorkload(
  map: Record<string, string>,
  namespace: string,
  name: string,
): string | null {
  const key = `${namespace}/${name}`;
  if (map[key]) return map[key];
  const hit = Object.entries(map)
    .filter(([p]) => p.includes('*') && ruleMatches(p, key))
    .sort((a, b) => b[0].length - a[0].length)[0];
  return hit ? hit[1] : null;
}

export const UNASSIGNED_SYSTEM = 'unassigned-workloads';

export interface Placement {
  system: string;
  subcomponentOf?: string;
  title: string;
  /** Distinguishes duplicate titles inside one System (same rule as the curator). */
  suffix?: string;
  pillar: string | null;
  review: string;
}

/** Where a workload belongs. Unmapped workloads land in the "Unassigned workloads" System. */
export function placeWorkload(
  namespace: string,
  name: string,
  mapping: InventoryMapping,
): Placement {
  const appRef = appForWorkload(mapping.workloads, namespace, name);
  const appName = appRef?.split('/')[1];
  const app = appName ? mapping.apps[appName] : undefined;
  if (appRef && app) {
    const { role, suffix } = workloadRole(name, appName!, namespace);
    return {
      system: app.system,
      subcomponentOf: appRef,
      title: `${app.title} — ${role}`,
      suffix,
      pillar: app.pillar,
      review: app.review,
    };
  }
  return {
    system: UNASSIGNED_SYSTEM,
    title: `${words(namespace)} — ${words(name)}`,
    pillar: null,
    review: 'unclassified',
  };
}
