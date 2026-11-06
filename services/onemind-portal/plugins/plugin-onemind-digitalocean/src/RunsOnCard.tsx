import { useEffect, useState } from 'react';
import { InfoCard } from '@backstage/core-components';
import { Link, Typography } from '@material-ui/core';
import { useEntity, catalogApiRef } from '@backstage/plugin-catalog-react';
import { useApi } from '@backstage/core-plugin-api';
import { parseEntityRef, stringifyEntityRef } from '@backstage/catalog-model';
import type { Entity } from '@backstage/catalog-model';

const isDo = (e: Entity) =>
  e.kind === 'Resource' && String(e.spec?.type ?? '').startsWith('do-');

const dependsOn = (e: Entity) =>
  (e.relations ?? []).filter(r => r.type === 'dependsOn').map(r => r.targetRef);

/**
 * "Runs on DigitalOcean": the DigitalOcean resources an app depends on in the catalog, found by
 * following dependsOn (the app, its parts, then the cluster Resource to node pools).
 */
export async function findDoResources(
  start: Entity,
  load: (refs: string[]) => Promise<(Entity | undefined)[]>,
  depth = 3,
): Promise<Entity[]> {
  const found = new Map<string, Entity>();
  const seen = new Set<string>([stringifyEntityRef(start)]);
  const parts = (start.relations ?? [])
    .filter(r => r.type === 'hasPart')
    .map(r => r.targetRef);
  let frontier = [...dependsOn(start), ...parts];
  for (let level = 0; level < depth && frontier.length; level++) {
    const refs = [...new Set(frontier)].filter(r => !seen.has(r));
    refs.forEach(r => seen.add(r));
    const entities = (await load(refs)).filter((e): e is Entity => !!e);
    frontier = [];
    for (const e of entities) {
      if (isDo(e)) found.set(stringifyEntityRef(e), e);
      frontier.push(...dependsOn(e));
    }
  }
  return [...found.values()];
}

const LABEL: Record<string, string> = {
  'do-volume': 'Data volumes',
  'do-node-pool': 'Node pools',
  'do-kubernetes-cluster': 'Kubernetes cluster',
};

export function RunsOnCard() {
  const { entity } = useEntity();
  const catalog = useApi(catalogApiRef);
  const [items, setItems] = useState<Entity[] | undefined>();
  useEffect(() => {
    let active = true;
    findDoResources(
      entity,
      async refs =>
        (await catalog.getEntitiesByRefs({ entityRefs: refs })).items,
    )
      .then(r => active && setItems(r))
      .catch(() => active && setItems([]));
    return () => {
      active = false;
    };
  }, [entity, catalog]);

  if (!items?.length) return null;
  const groups = Object.keys(LABEL).map(type => ({
    type,
    list: items.filter(e => e.spec?.type === type),
  }));
  return (
    <InfoCard title="Runs on DigitalOcean">
      {groups
        .filter(g => g.list.length)
        .map(g => (
          <div key={g.type}>
            <Typography variant="subtitle2">{LABEL[g.type]}</Typography>
            {g.list.map(e => {
              const ref = parseEntityRef(stringifyEntityRef(e));
              return (
                <Typography key={e.metadata.name} variant="body2">
                  <Link
                    href={`/catalog/${
                      ref.namespace
                    }/${ref.kind.toLowerCase()}/${ref.name}`}
                  >
                    {e.metadata.title ?? e.metadata.name}
                  </Link>
                </Typography>
              );
            })}
          </div>
        ))}
    </InfoCard>
  );
}
