import { useEffect, useRef } from 'react';
import {
  Card,
  CardContent,
  MenuItem,
  MenuList,
  Typography,
} from '@material-ui/core';
import {
  PILLAR_LABEL_KEY,
  ROLE_LABEL_KEY,
} from '@internal/plugin-onemind-common';
import { CatalogFilterBlueprint } from '@backstage/plugin-catalog-react/alpha';
import {
  DefaultEntityFilters,
  EntityFilter,
  EntityKindFilter,
  EntityTypeFilter,
  useEntityList,
} from '@backstage/plugin-catalog-react';

// Pillars are a reviewed classification, not a guess (docs/roadmap/onemind-catalog-pillars.md).
// Backstage filters can only ask for values that exist, so every classified entity carries the
// label `onemind.io/pillar` and every other entity carries the annotation
// `onemind.io/pillar-review: unclassified`; both come from catalog/structure.yaml.
export const PILLAR_LABEL = `metadata.labels.${PILLAR_LABEL_KEY}`;
export const PILLAR_REVIEW = 'metadata.annotations.onemind.io/pillar-review';
export const ROLE_LABEL = `metadata.labels.${ROLE_LABEL_KEY}`;

export type PillarOption =
  | 'all'
  | 'so'
  | 'le'
  | 'hp'
  | 'ge'
  | 'shared'
  | 'unclassified';
export const pillarOptions: { id: PillarOption; label: string }[] = [
  { id: 'all', label: 'All' },
  { id: 'so', label: 'SO — Sovereign Operations' },
  { id: 'le', label: 'LE — Legacy Evolution' },
  { id: 'hp', label: 'HP — Holistic Performance' },
  { id: 'ge', label: 'GE — Generational Empire' },
  { id: 'shared', label: 'Shared' },
  { id: 'unclassified', label: 'Unclassified' },
];

export class EntityPillarFilter implements EntityFilter {
  constructor(readonly value: PillarOption) {}
  getCatalogFilters(): Record<string, string | symbol> {
    switch (this.value) {
      case 'so':
      case 'le':
      case 'hp':
      case 'ge':
        return { [PILLAR_LABEL]: this.value };
      case 'shared':
        // Shared is a reviewed view (label value), not a fifth Domain.
        return { [PILLAR_LABEL]: 'shared' };
      case 'unclassified':
        return { [PILLAR_REVIEW]: 'unclassified' };
      default:
        return {};
    }
  }
  toQueryValue() {
    return this.value;
  }
}

export type ViewOption = 'apps' | 'infrastructure' | 'repos' | 'everything';
export const viewOptions: { id: ViewOption; label: string }[] = [
  { id: 'apps', label: 'Apps' },
  { id: 'infrastructure', label: 'Infrastructure' },
  { id: 'repos', label: 'Repos' },
  { id: 'everything', label: 'Everything' },
];

/** Separates the primary app entity from the Kubernetes workloads that hang under it. */
export class EntityRoleFilter implements EntityFilter {
  constructor(readonly value: 'app' | 'workload') {}
  getCatalogFilters(): Record<string, string> {
    return { [ROLE_LABEL]: this.value };
  }
  toQueryValue() {
    return this.value;
  }
}

type OneMindFilters = DefaultEntityFilters & {
  pillar?: EntityPillarFilter;
  role?: EntityRoleFilter;
};

/** Filters a "View" preset sets: kind, optional type and optional role. */
export function viewFilters(view: ViewOption): Partial<OneMindFilters> {
  switch (view) {
    case 'apps':
      return {
        kind: new EntityKindFilter('component', 'Component'),
        type: undefined,
        role: new EntityRoleFilter('app'),
      };
    case 'infrastructure':
      return {
        kind: new EntityKindFilter('component', 'Component'),
        type: undefined,
        role: new EntityRoleFilter('workload'),
      };
    case 'repos':
      return {
        kind: new EntityKindFilter('resource', 'Resource'),
        type: new EntityTypeFilter('code-repository'),
        role: undefined,
      };
    default:
      return {
        kind: new EntityKindFilter('component', 'Component'),
        type: undefined,
        role: undefined,
      };
  }
}

export function currentView(filters: OneMindFilters): ViewOption {
  const kind = filters.kind?.value?.toLowerCase();
  if (kind === 'resource' && filters.type) return 'repos';
  if (kind === 'component' && filters.role?.value === 'app') return 'apps';
  if (kind === 'component' && filters.role?.value === 'workload')
    return 'infrastructure';
  return 'everything';
}

function Picker<T extends string>(props: {
  title: string;
  options: { id: T; label: string }[];
  selected: T;
  onSelect: (id: T) => void;
  testId: string;
}) {
  return (
    <Card data-testid={props.testId} variant="outlined">
      <CardContent>
        <Typography variant="subtitle2" gutterBottom>
          {props.title}
        </Typography>
        <MenuList dense disablePadding>
          {props.options.map(option => (
            <MenuItem
              key={option.id}
              selected={option.id === props.selected}
              onClick={() => props.onSelect(option.id)}
            >
              {option.label}
            </MenuItem>
          ))}
        </MenuList>
      </CardContent>
    </Card>
  );
}

export function PillarPicker() {
  const { filters, updateFilters } = useEntityList<OneMindFilters>();
  return (
    <Picker
      testId="pillar-picker"
      title="Pillar"
      options={pillarOptions}
      selected={filters.pillar?.value ?? 'all'}
      onSelect={id =>
        updateFilters({
          pillar: id === 'all' ? undefined : new EntityPillarFilter(id),
        })
      }
    />
  );
}

/** Saved views. The catalog opens on Apps so workloads, repositories and mesh endpoints are one click away, not in the way. */
export function ViewPicker() {
  const { filters, updateFilters } = useEntityList<OneMindFilters>();
  const applied = useRef(false);
  useEffect(() => {
    if (applied.current) return;
    applied.current = true;
    // Start on Apps unless something else (a deep link) already narrowed the list.
    if (
      !filters.role &&
      !filters.type &&
      (filters.kind?.value ?? 'component').toLowerCase() === 'component'
    )
      updateFilters(viewFilters('apps'));
  }, [filters, updateFilters]);
  return (
    <Picker
      testId="view-picker"
      title="View"
      options={viewOptions}
      selected={currentView(filters)}
      onSelect={id => updateFilters(viewFilters(id))}
    />
  );
}

// Catalog index sidebar: saved views (Apps / Infrastructure / Repos / Everything) and the
// pillar filter (All / SO / LE / HP / GE / Shared / Unclassified).
export const catalogViewFilterExtension = CatalogFilterBlueprint.make({
  name: 'view',
  params: { loader: async () => <ViewPicker /> },
});
export const catalogPillarFilterExtension = CatalogFilterBlueprint.make({
  name: 'pillar',
  params: { loader: async () => <PillarPicker /> },
});
