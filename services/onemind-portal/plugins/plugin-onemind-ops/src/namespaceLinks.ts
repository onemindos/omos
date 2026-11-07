import { useEffect, useMemo, useState } from 'react';
import { configApiRef, useApi } from '@backstage/core-plugin-api';
import { catalogApiRef } from '@backstage/plugin-catalog-react';

const NAMESPACE_ANNOTATION = 'backstage.io/kubernetes-namespace';

interface Entry {
  metadata: {
    name: string;
    namespace?: string;
    title?: string;
    annotations?: Record<string, string>;
  };
}

export interface NamespaceTarget {
  href: string;
  title: string;
  /** `catalog` opens the component's Kubernetes tab; `headlamp` opens the Headlamp dashboard. */
  via: 'catalog' | 'headlamp';
}

export type NamespaceLinks = (namespace: string) => NamespaceTarget | undefined;

/**
 * Pure: a namespace links to the Kubernetes tab of the first catalog Component annotated with that namespace
 * (sorted by name, so the choice is stable). Without one it falls back to the Headlamp URL, if configured.
 */
export function buildNamespaceLinks(
  components: Entry[],
  headlampUrl?: string,
): NamespaceLinks {
  const byNamespace = new Map<string, NamespaceTarget>();
  const sorted = [...components].sort((a, b) =>
    a.metadata.name.localeCompare(b.metadata.name),
  );
  for (const e of sorted) {
    const ns = e.metadata.annotations?.[NAMESPACE_ANNOTATION];
    if (!ns || byNamespace.has(ns)) continue;
    byNamespace.set(ns, {
      href: `/catalog/${e.metadata.namespace ?? 'default'}/component/${
        e.metadata.name
      }/kubernetes`,
      title: e.metadata.title ?? e.metadata.name,
      via: 'catalog',
    });
  }
  return namespace =>
    byNamespace.get(namespace) ??
    (headlampUrl
      ? { href: headlampUrl, title: 'Headlamp', via: 'headlamp' }
      : undefined);
}

/** First `dashboardUrl` under kubernetes.clusterLocatorMethods[].clusters[], if the config exposes it. */
export function headlampUrlFrom(read: () => unknown): string | undefined {
  try {
    const methods = read();
    if (!Array.isArray(methods)) return undefined;
    for (const m of methods) {
      const clusters = (m as { clusters?: Array<{ dashboardUrl?: unknown }> })
        ?.clusters;
      for (const cl of clusters ?? [])
        if (typeof cl?.dashboardUrl === 'string' && cl.dashboardUrl)
          return cl.dashboardUrl;
    }
  } catch {
    // Not configured for the frontend.
  }
  return undefined;
}

export function useNamespaceLinks(): NamespaceLinks {
  const catalog = useApi(catalogApiRef);
  const config = useApi(configApiRef);
  const [components, setComponents] = useState<Entry[]>([]);
  useEffect(() => {
    let active = true;
    catalog
      .getEntities({
        filter: { kind: 'Component' },
        fields: [
          'metadata.name',
          'metadata.namespace',
          'metadata.title',
          'metadata.annotations',
        ],
      })
      .then(r => active && setComponents(r.items as Entry[]))
      .catch(() => {
        // Cells stay unlinked when the catalog is unreachable.
      });
    return () => {
      active = false;
    };
  }, [catalog]);
  const headlamp = useMemo(
    () =>
      headlampUrlFrom(() =>
        config.getOptional('kubernetes.clusterLocatorMethods'),
      ),
    [config],
  );
  return useMemo(
    () => buildNamespaceLinks(components, headlamp),
    [components, headlamp],
  );
}
