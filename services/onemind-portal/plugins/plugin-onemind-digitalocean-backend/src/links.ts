import type { GetJson } from './kube';
import type { AppLink, PvcLink, Volume } from './types';

export interface Pvc extends PvcLink {
  volumeName: string;
  labels: Record<string, string>;
  /** Names of the Deployments/StatefulSets/DaemonSets whose pods mount the claim right now. */
  mountedBy: string[];
}

/** A catalog Component that stands for one Deployment or StatefulSet. */
export interface WorkloadComponent {
  ref: string;
  title?: string;
  namespace: string;
  workload: string;
  /** "k=v,k=v" from the catalog annotation backstage.io/kubernetes-label-selector. */
  selector: string;
}

/**
 * Every bound PersistentVolumeClaim, read-only (the portal's ServiceAccount already lists claims).
 * The CSI driver names the DigitalOcean volume after its PersistentVolume, pvc-<uid>, which is the
 * claim's spec.volumeName, so no PersistentVolume read (and no extra RBAC) is needed.
 */
export async function listPvcs(getJson: GetJson): Promise<Pvc[]> {
  const mounts = await listMounts(getJson);
  const out: Pvc[] = [];
  let token: string | undefined;
  do {
    const page = (await getJson(
      `/api/v1/persistentvolumeclaims?limit=500${
        token ? `&continue=${encodeURIComponent(token)}` : ''
      }`,
    )) as {
      items?: {
        metadata?: {
          name?: string;
          namespace?: string;
          labels?: Record<string, string>;
        };
        spec?: { volumeName?: string };
      }[];
      metadata?: { continue?: string };
    };
    if (!page || !Array.isArray(page.items))
      throw new Error('Unexpected PersistentVolumeClaim list response');
    for (const i of page.items) {
      const { name, namespace } = i.metadata ?? {};
      const volumeName = i.spec?.volumeName;
      if (name && namespace && volumeName)
        out.push({
          name,
          namespace,
          volumeName,
          labels: i.metadata?.labels ?? {},
          mountedBy: [...(mounts.get(`${namespace}/${name}`) ?? [])],
        });
    }
    token = page.metadata?.continue || undefined;
  } while (token);
  return out;
}

/** Workload that owns a pod: StatefulSet/DaemonSet by name, ReplicaSet minus its pod-template hash. */
function workloadOf(owner?: { kind?: string; name?: string }) {
  if (!owner?.name) return undefined;
  if (owner.kind === 'StatefulSet' || owner.kind === 'DaemonSet')
    return owner.name;
  if (owner.kind === 'ReplicaSet')
    return owner.name.replace(/-[a-z0-9]{5,10}$/, '');
  return undefined;
}

/** namespace/claim -> workloads whose pods mount it (read-only pod list; only claim names are kept). */
async function listMounts(getJson: GetJson): Promise<Map<string, Set<string>>> {
  const out = new Map<string, Set<string>>();
  let token: string | undefined;
  do {
    const page = (await getJson(
      `/api/v1/pods?limit=500${
        token ? `&continue=${encodeURIComponent(token)}` : ''
      }`,
    )) as {
      items?: {
        metadata?: {
          namespace?: string;
          ownerReferences?: { kind?: string; name?: string }[];
        };
        spec?: {
          volumes?: { persistentVolumeClaim?: { claimName?: string } }[];
        };
      }[];
      metadata?: { continue?: string };
    };
    if (!page || !Array.isArray(page.items))
      throw new Error('Unexpected Pod list response');
    for (const pod of page.items) {
      const namespace = pod.metadata?.namespace;
      const workload = workloadOf(pod.metadata?.ownerReferences?.[0]);
      if (!namespace || !workload) continue;
      for (const v of pod.spec?.volumes ?? []) {
        const claim = v.persistentVolumeClaim?.claimName;
        if (!claim) continue;
        const key = `${namespace}/${claim}`;
        out.set(key, (out.get(key) ?? new Set()).add(workload));
      }
    }
    token = page.metadata?.continue || undefined;
  } while (token);
  return out;
}

const escape = (s: string) => s.replace(/[.*+?^${}()|[\]\\]/g, '\\$&');

const parseSelector = (selector: string): [string, string][] =>
  selector
    .split(',')
    .map(part => part.split('=') as [string, string])
    .filter(([k, v]) => k && v !== undefined);

/**
 * The workload that owns a claim, or undefined when the mapping is not certain:
 *  0. the one workload whose pods mount the claim (certain)
 *  1. a workload whose selector labels are all on the claim (most labels wins; a tie is no answer)
 *  2. a StatefulSet claim name, <template>-<workload>-<ordinal>, or a claim named <workload>[-...]
 */
export function ownerOf(
  pvc: Pvc,
  workloads: WorkloadComponent[],
): WorkloadComponent | undefined {
  const here = workloads.filter(w => w.namespace === pvc.namespace);
  const mounted = here.filter(w => pvc.mountedBy.includes(w.workload));
  if (mounted.length === 1) return mounted[0];
  if (mounted.length > 1) return undefined;
  let best: WorkloadComponent[] = [];
  let score = 0;
  for (const w of here) {
    const pairs = parseSelector(w.selector);
    if (!pairs.length || !pairs.every(([k, v]) => pvc.labels[k] === v))
      continue;
    if (pairs.length > score) {
      score = pairs.length;
      best = [w];
    } else if (pairs.length === score) best.push(w);
  }
  if (best.length === 1) return best[0];
  if (best.length > 1) return undefined;
  const byName = here
    .filter(w => new RegExp(`(^|-)${escape(w.workload)}-\\d+$`).test(pvc.name))
    .sort((a, b) => b.workload.length - a.workload.length);
  if (byName[0]) return byName[0];
  const byPrefix = here
    .filter(
      w => pvc.name === w.workload || pvc.name.startsWith(`${w.workload}-`),
    )
    .sort((a, b) => b.workload.length - a.workload.length);
  return byPrefix[0];
}

/** Adds pvc and app to every volume that a claim names. Volumes without a claim are left as they are. */
export function linkVolumes(
  volumes: Volume[],
  pvcs: Pvc[],
  workloads: WorkloadComponent[],
): Volume[] {
  const byVolume = new Map(pvcs.map(p => [p.volumeName, p]));
  return volumes.map(v => {
    const pvc = byVolume.get(v.name);
    if (!pvc) return v;
    const owner = ownerOf(pvc, workloads);
    const app: AppLink | undefined = owner
      ? { ref: owner.ref, title: owner.title }
      : undefined;
    return { ...v, pvc: { namespace: pvc.namespace, name: pvc.name }, app };
  });
}
