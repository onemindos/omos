import type { Entity } from '@backstage/catalog-model';
import type { Change, ResourceType, Snapshots } from './types';

export const PROVIDER_NAME = 'onemind-digitalocean';
export const OWNER = 'group:default/onemind-admins';
const CONSOLE = 'https://cloud.digitalocean.com';

export const A = {
  owner: 'onemind.io/digitalocean-owner',
  id: 'onemind.io/digitalocean-id',
  pvc: 'onemind.io/digitalocean-pvc',
  unattached: 'onemind.io/digitalocean-unattached',
} as const;

/** One DigitalOcean thing the catalog should know about, before it becomes an Entity. */
export interface DoResource {
  type: ResourceType;
  id: string;
  title: string;
  description: string;
  url: string;
  region?: string;
  dependsOn?: string[];
  dependencyOf?: string[];
  annotations?: Record<string, string>;
}

export const entityNameFor = (type: ResourceType, id: string) =>
  `do-${type}-${id}`.toLowerCase();

const tag = (value: string) =>
  value
    .toLowerCase()
    .replace(/[^a-z0-9+#]+/g, '-')
    .replace(/^-+|-+$/g, '')
    .slice(0, 63);

const ref = (type: ResourceType, id: string) =>
  `resource:default/${entityNameFor(type, id)}`;

/**
 * The resources of every section that read successfully, grouped by whether the section is
 * complete. A failed section is left out of `complete`, so a failed read is never mistaken
 * for "everything was deleted".
 */
export function collectResources(s: Snapshots): {
  resources: DoResource[];
  complete: Set<ResourceType>;
} {
  const resources: DoResource[] = [];
  const complete = new Set<ResourceType>();

  if (s.kubernetes.state === 'ok') {
    complete.add('kubernetes-cluster');
    complete.add('node-pool');
    for (const c of s.kubernetes.items) {
      resources.push({
        type: 'kubernetes-cluster',
        id: c.id,
        title: `${c.name} (DOKS ${c.version})`,
        description: `DigitalOcean Kubernetes cluster in ${c.region}`,
        url: `${CONSOLE}/kubernetes/clusters/${c.id}`,
        region: c.region,
        dependsOn: c.nodePools.map(p => ref('node-pool', p.id)),
        // The reviewed catalog already has one Resource for this cluster (do-<region>-<name>).
        dependencyOf: [`resource:default/do-${c.region}-${c.name}`],
      });
      for (const p of c.nodePools) {
        resources.push({
          type: 'node-pool',
          id: p.id,
          title: `${c.name} / ${p.name} (${p.count} x ${p.size})`,
          description: `Node pool of ${c.name}: ${p.count} x ${p.size}${
            p.autoScale ? `, autoscale ${p.minNodes} to ${p.maxNodes}` : ''
          }`,
          url: `${CONSOLE}/kubernetes/clusters/${c.id}`,
          region: c.region,
        });
      }
    }
  }
  if (s.droplets.state === 'ok') {
    complete.add('droplet');
    for (const d of s.droplets.items)
      resources.push({
        type: 'droplet',
        id: d.id,
        title: d.name,
        description: `Droplet ${d.size} in ${d.region} (${d.status})`,
        url: `${CONSOLE}/droplets/${d.id}`,
        region: d.region,
      });
  }
  if (s.volumes.state === 'ok') {
    complete.add('volume');
    for (const v of s.volumes.items) {
      const annotations: Record<string, string> = {};
      if (v.pvc) annotations[A.pvc] = `${v.pvc.namespace}/${v.pvc.name}`;
      if (v.unattached) annotations[A.unattached] = 'true';
      resources.push({
        type: 'volume',
        id: v.id,
        title: v.pvc
          ? `${v.pvc.namespace}/${v.pvc.name} (${v.sizeGb} GB)`
          : `${v.name} (${v.sizeGb} GB)`,
        description: `Block storage volume ${v.name}, ${v.sizeGb} GB in ${
          v.region
        }${v.unattached ? ', not attached to a droplet' : ''}`,
        url: `${CONSOLE}/volumes/${v.id}`,
        region: v.region,
        dependencyOf: v.app ? [v.app.ref] : undefined,
        annotations,
      });
    }
  }
  if (s.snapshots.state === 'ok') {
    complete.add('snapshot');
    for (const x of s.snapshots.items)
      resources.push({
        type: 'snapshot',
        id: x.id,
        title: x.name,
        description: `${x.resourceType} snapshot, ${x.sizeGb} GB`,
        url: `${CONSOLE}/images/snapshots/${
          x.resourceType === 'volume' ? 'volume' : 'droplet'
        }`,
        region: x.regions[0],
      });
  }
  if (s.loadBalancers.state === 'ok') {
    complete.add('load-balancer');
    for (const l of s.loadBalancers.items)
      resources.push({
        type: 'load-balancer',
        id: l.id,
        title: l.name,
        description: `Load balancer in ${l.region}, ${l.forwardingRules} forwarding rules`,
        url: `${CONSOLE}/networking/load_balancers/${l.id}`,
        region: l.region,
      });
  }
  if (s.firewalls.state === 'ok') {
    complete.add('firewall');
    for (const f of s.firewalls.items)
      resources.push({
        type: 'firewall',
        id: f.id,
        title: f.name,
        description: `Cloud firewall, ${f.inboundRules} inbound and ${f.outboundRules} outbound rules`,
        url: `${CONSOLE}/networking/firewalls/${f.id}`,
      });
  }
  if (s.vpcs.state === 'ok') {
    complete.add('vpc');
    for (const v of s.vpcs.items)
      resources.push({
        type: 'vpc',
        id: v.id,
        title: v.name,
        description: `VPC ${v.ipRange} in ${v.region}${
          v.isDefault ? ' (default)' : ''
        }`,
        url: `${CONSOLE}/networking/vpc/${v.id}`,
        region: v.region,
      });
  }
  return { resources, complete };
}

export const toChange = (r: DoResource): Change => ({
  type: r.type,
  id: r.id,
  name: r.title,
});

/** Catalog Resource for one DigitalOcean thing. No timestamps: the same input is the same entity. */
export function buildEntity(r: DoResource): Entity {
  const name = entityNameFor(r.type, r.id);
  const tags = ['digitalocean', r.type, ...(r.region ? [tag(r.region)] : [])];
  return {
    apiVersion: 'backstage.io/v1alpha1',
    kind: 'Resource',
    metadata: {
      name,
      title: r.title.slice(0, 200),
      description: r.description,
      annotations: {
        'backstage.io/managed-by-location': `provider:${PROVIDER_NAME}`,
        'backstage.io/managed-by-origin-location': `provider:${PROVIDER_NAME}`,
        [A.owner]: PROVIDER_NAME,
        [A.id]: r.id,
        'onemind.io/hosting': 'digitalocean',
        ...r.annotations,
      },
      tags: [...new Set(tags.filter(Boolean))],
      links: [{ url: r.url, title: 'Open in DigitalOcean', icon: 'dashboard' }],
    },
    spec: {
      type: `do-${r.type}`,
      owner: OWNER,
      system: 'infrastructure',
      ...(r.dependsOn?.length ? { dependsOn: r.dependsOn } : {}),
      ...(r.dependencyOf?.length ? { dependencyOf: r.dependencyOf } : {}),
    },
  };
}

/** All entities for a snapshot, in a fixed order. */
export function buildEntities(s: Snapshots): Entity[] {
  return collectResources(s)
    .resources.map(buildEntity)
    .sort((a, b) => a.metadata.name.localeCompare(b.metadata.name));
}
