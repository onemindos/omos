import type { NamespaceRow } from './allocation';

/** A catalog Component that declares which Kubernetes namespace it runs in. */
export type CatalogApp = {
  name: string;
  title: string;
  href: string;
  namespaces: string[];
  pillar?: string;
  pillarApproved: boolean;
};

type EntityLike = {
  metadata: {
    name: string;
    namespace?: string;
    title?: string;
    annotations?: Record<string, string>;
    labels?: Record<string, string>;
  };
};

const NAMESPACE = 'backstage.io/kubernetes-namespace';
const NAMESPACES = 'onemind.io/kubernetes-namespaces';
const PILLAR = 'onemind.io/pillar';
const PILLAR_REVIEW = 'onemind.io/pillar-review';

export const PILLAR_TITLES: Record<string, string> = {
  so: 'SO · Sovereign Operations',
  le: 'LE · Legacy Evolution',
  hp: 'HP · Holistic Performance',
  ge: 'GE · Generational Empire',
  shared: 'Shared platform',
};

/** Pure: reads the namespace and pillar annotations the catalog already carries. */
export function appsFromEntities(entities: EntityLike[]): CatalogApp[] {
  const apps: CatalogApp[] = [];
  for (const entity of entities) {
    const { name, annotations = {}, labels = {} } = entity.metadata;
    const namespaces = [
      ...(annotations[NAMESPACE] ? [annotations[NAMESPACE]] : []),
      ...(annotations[NAMESPACES] ?? '').split(','),
    ]
      .map(value => value.trim())
      .filter(Boolean);
    if (!namespaces.length) continue;
    apps.push({
      name,
      title: entity.metadata.title ?? name,
      href: `/catalog/${
        entity.metadata.namespace ?? 'default'
      }/component/${name}`,
      namespaces: [...new Set(namespaces)],
      pillar: labels[PILLAR],
      pillarApproved: annotations[PILLAR_REVIEW] === 'approved',
    });
  }
  return apps.sort((a, b) => a.name.localeCompare(b.name));
}

/**
 * How a namespace's cost relates to catalog apps:
 *  - app:    exactly one catalog app declares this namespace, so the cost is its cost
 *  - shared: several apps share the namespace; OpenCost's namespace total is NOT split
 *  - none:   no catalog app declares the namespace
 */
export type Attribution = 'app' | 'shared' | 'none';

export type PillarKey = string; // a pillar slug, or one of the values below
export const MIXED = 'mixed';
export const UNSET = 'unset';
export const UNMAPPED = 'unmapped';

export type CostLine = {
  namespace: string;
  /** null means OpenCost did not give a usable number. It is never shown as zero. */
  cost: number | null;
  row: NamespaceRow;
  apps: CatalogApp[];
  attribution: Attribution;
  pillar: PillarKey;
  /** True when the pillar comes from an assignment nobody has approved yet. */
  pillarProposed: boolean;
};

function pillarOf(apps: CatalogApp[]): {
  pillar: PillarKey;
  proposed: boolean;
} {
  if (!apps.length) return { pillar: UNMAPPED, proposed: false };
  const pillars = new Set(apps.map(app => app.pillar ?? UNSET));
  if (pillars.size > 1) return { pillar: MIXED, proposed: false };
  const [pillar] = [...pillars];
  return {
    pillar,
    proposed: pillar !== UNSET && apps.some(app => !app.pillarApproved),
  };
}

function attributionOf(count: number): Attribution {
  if (count === 0) return 'none';
  return count === 1 ? 'app' : 'shared';
}

export function buildCostLines(
  rows: NamespaceRow[],
  apps: CatalogApp[],
): CostLine[] {
  return rows
    .map(row => {
      const matched = apps.filter(app =>
        app.namespaces.includes(row.namespace),
      );
      const { pillar, proposed } = pillarOf(matched);
      return {
        namespace: row.namespace,
        cost: row.totalCost,
        row,
        apps: matched,
        attribution: attributionOf(matched.length),
        pillar,
        pillarProposed: proposed,
      };
    })
    .sort(
      (a, b) =>
        (b.cost ?? -1) - (a.cost ?? -1) ||
        a.namespace.localeCompare(b.namespace),
    );
}

export type Rollup = {
  /** Sum of the namespaces whose cost is known. A lower bound when `unknown > 0`. */
  known: number;
  /** Namespaces in the group whose cost is unknown. */
  unknown: number;
  namespaces: number;
};

export function rollup(lines: CostLine[]): Rollup {
  return {
    known: lines.reduce((sum, line) => sum + (line.cost ?? 0), 0),
    unknown: lines.filter(line => line.cost === null).length,
    namespaces: lines.length,
  };
}

export type PillarRollup = Rollup & {
  pillar: PillarKey;
  proposed: boolean;
};

const ORDER = ['so', 'le', 'hp', 'ge', 'shared', MIXED, UNSET, UNMAPPED];

export function pillarRollups(lines: CostLine[]): PillarRollup[] {
  const groups = new Map<PillarKey, CostLine[]>();
  for (const line of lines)
    groups.set(line.pillar, [...(groups.get(line.pillar) ?? []), line]);
  return [...groups]
    .map(([pillar, members]) => ({
      pillar,
      proposed: members.some(line => line.pillarProposed),
      ...rollup(members),
    }))
    .sort(
      (a, b) =>
        (ORDER.indexOf(a.pillar) + 1 || 99) -
        (ORDER.indexOf(b.pillar) + 1 || 99),
    );
}

export function pillarLabel(pillar: PillarKey): string {
  if (pillar === MIXED) return 'Mixed pillars (shared namespace)';
  if (pillar === UNSET) return 'App has no pillar yet';
  if (pillar === UNMAPPED) return 'No catalog app for this namespace';
  return PILLAR_TITLES[pillar] ?? pillar;
}
