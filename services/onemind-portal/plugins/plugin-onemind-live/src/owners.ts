import { useEffect, useMemo, useState } from 'react';
import { useApi } from '@backstage/core-plugin-api';
import { catalogApiRef } from '@backstage/plugin-catalog-react';

export type Owner = {
  name: string;
  title: string;
  href: string;
  pillar?: string;
};

type Entry = {
  metadata: {
    name: string;
    namespace?: string;
    title?: string;
    annotations?: Record<string, string>;
    labels?: Record<string, string>;
  };
};

const ARGO = 'argocd/app-name';
const NAMESPACE = 'backstage.io/kubernetes-namespace';
const PILLAR = 'onemind.io/pillar';

export type OwnerIndex = {
  byArgoApp(name?: string): Owner | undefined;
  byNamespace(namespace?: string): Owner | undefined;
  /** Argo application first (exact), then the Kubernetes namespace. */
  resolve(input: { argoApp?: string; namespace?: string }): Owner | undefined;
};

export const emptyIndex: OwnerIndex = {
  byArgoApp: () => undefined,
  byNamespace: () => undefined,
  resolve: () => undefined,
};

/** Pure: builds lookups from catalog Components (F2 annotations and pillar labels). */
export function buildOwnerIndex(components: Entry[]): OwnerIndex {
  const argo = new Map<string, Owner>();
  const namespaces = new Map<string, Owner>();
  for (const entity of components) {
    const { name, annotations = {}, labels = {} } = entity.metadata;
    const owner: Owner = {
      name,
      title: entity.metadata.title ?? name,
      href: `/catalog/${entity.metadata.namespace ?? 'default'}/component/${name}`,
      pillar: labels[PILLAR],
    };
    const app = annotations[ARGO];
    if (app && !argo.has(app)) argo.set(app, owner);
    const namespace = annotations[NAMESPACE];
    // A namespace can hold several components; the first (by name) is only a fallback.
    if (namespace && !namespaces.has(namespace)) namespaces.set(namespace, owner);
  }
  const byArgoApp = (name?: string) => (name ? argo.get(name) : undefined);
  const byNamespace = (namespace?: string) => (namespace ? namespaces.get(namespace) : undefined);
  return {
    byArgoApp,
    byNamespace,
    resolve: ({ argoApp, namespace }) => byArgoApp(argoApp) ?? byNamespace(namespace),
  };
}

/** Loads the Component index once. Unassigned items stay visible when the catalog is unreachable. */
export function useOwnerIndex() {
  const catalog = useApi(catalogApiRef);
  const [components, setComponents] = useState<Entry[]>([]);
  const [loaded, setLoaded] = useState(false);
  useEffect(() => {
    let active = true;
    (async () => {
      try {
        const result = await catalog.getEntities({
          filter: { kind: 'Component' },
          fields: [
            'metadata.name',
            'metadata.namespace',
            'metadata.title',
            'metadata.annotations',
            'metadata.labels',
          ],
        });
        if (active) setComponents(result.items as Entry[]);
      } catch {
        // Items fall back to Unassigned; the page still renders.
      } finally {
        if (active) setLoaded(true);
      }
    })();
    return () => {
      active = false;
    };
  }, [catalog]);
  return { index: useMemo(() => buildOwnerIndex(components), [components]), loaded };
}
