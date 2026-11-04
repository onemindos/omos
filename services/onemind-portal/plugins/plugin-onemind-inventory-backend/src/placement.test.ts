import { A, buildWorkloadEntity } from './entities';
import { appForWorkload, placeWorkload, workloadRole } from './placement';
import { mapping } from './mapping.generated';
import {
  curatedWorkloads,
  curator,
  generator,
  liveWorkloads,
  readPortalFile,
} from './testing';
import YAML from 'yaml';

const observedAt = '2026-10-22T00:00:00.000Z';

describe('mapping through F2', () => {
  it('uses the same app and role rules as scripts/curate-structure.cjs for every live workload', async () => {
    const map = YAML.parse(readPortalFile('catalog/workload-map.yaml'));
    const script = curator();
    const live = await liveWorkloads();
    const ours = live.map(w =>
      appForWorkload(mapping.workloads, w.namespace, w.name),
    );
    expect(ours).toEqual(
      live.map(w => script.appForWorkload(map, w.namespace, w.name)),
    );
    const roles = live.map((w, i) =>
      workloadRole(w.name, ours[i]?.split('/')[1] ?? 'unmapped', w.namespace),
    );
    expect(roles).toEqual(
      live.map((w, i) =>
        script.workloadRole(
          w.name,
          ours[i]?.split('/')[1] ?? 'unmapped',
          w.namespace,
        ),
      ),
    );
  });

  it('sets system, subcomponentOf, pillar and review for every curated workload exactly as the committed file does', async () => {
    const curated = new Map(
      curatedWorkloads().map(e => [
        `${e.metadata.annotations![A.namespace]}/${
          e.metadata.annotations![A.workload]
        }`,
        e,
      ]),
    );
    const taken = new Map<string, Set<string>>();
    let titleMismatches = 0;
    for (const w of await liveWorkloads()) {
      const want = curated.get(`${w.namespace}/${w.name}`)!;
      const got = buildWorkloadEntity(w, {
        clusterName: 'do-nyc3-omos-cloud',
        observedAt,
        state: 'observed',
        takenTitles: taken,
      });
      expect(got.spec).toMatchObject({
        system: want.spec!.system,
        subcomponentOf: want.spec!.subcomponentOf,
        owner: want.spec!.owner,
        type: 'service',
      });
      expect(got.metadata.labels).toEqual(want.metadata.labels);
      expect(got.metadata.annotations!['onemind.io/pillar-review']).toBe(
        want.metadata.annotations!['onemind.io/pillar-review'],
      );
      expect(got.metadata.annotations![A.selector]).toBe(
        want.metadata.annotations![A.selector],
      );
      expect(got.metadata.annotations![A.namespace]).toBe(
        want.metadata.annotations![A.namespace],
      );
      expect(got.metadata.annotations![A.workloadKind]).toBe(
        want.metadata.annotations![A.workloadKind],
      );
      expect(got.metadata.tags).toEqual(want.metadata.tags);
      expect(got.metadata.description).toBe(want.metadata.description);
      if (got.metadata.title !== want.metadata.title) titleMismatches++;
    }
    expect(titleMismatches).toBe(0);
  });

  it('the generated bundle is current and every mapped app exists', () => {
    expect(generator().buildMapping()).toEqual(mapping);
    for (const ref of Object.values(mapping.workloads))
      expect(mapping.apps[ref.split('/')[1]]).toBeDefined();
  });

  it('puts an unmapped workload in Unassigned workloads and a mapped new one under its app', () => {
    const lone = placeWorkload('brand-new-ns', 'thing-api', mapping);
    expect(lone).toMatchObject({
      system: 'unassigned-workloads',
      pillar: null,
      review: 'unclassified',
      title: 'Brand New Ns — Thing API',
    });
    expect(lone.subcomponentOf).toBeUndefined();
    const argo = placeWorkload('argocd', 'argocd-server', mapping);
    expect(argo).toMatchObject({
      subcomponentOf: 'component:default/argocd',
      pillar: 'so',
    });
  });
});
