import { readFileSync } from 'node:fs';
import type { JsonObject } from '@backstage/types';
import YAML from 'yaml';
import {
  mockCredentials,
  mockServices,
  startTestBackend,
} from '@backstage/backend-test-utils';
import catalogPlugin from '@backstage/plugin-catalog-backend';
import { catalogModuleOnemindLiveInventory } from './module';
import { clusterYamlPath } from './testing';

jest.setTimeout(120_000);

interface Listed {
  metadata: { uid?: string; etag?: string; name: string };
  kind: string;
}

/** Starts a real catalog on the reviewed file and returns every entity it serves, minus volatile ids. */
type Features = NonNullable<Parameters<typeof startTestBackend>[0]['features']>;
async function serve(extra: Features, onemind?: JsonObject) {
  const backend = await startTestBackend({
    features: [
      catalogPlugin,
      ...extra,
      mockServices.rootConfig.factory({
        data: {
          catalog: {
            rules: [
              {
                allow: [
                  'Component',
                  'System',
                  'API',
                  'Resource',
                  'Location',
                  'Domain',
                ],
              },
            ],
            locations: [{ type: 'file', target: clusterYamlPath }],
          },
          ...(onemind ? { onemind } : {}),
        },
      }),
    ],
  });
  const url = `http://localhost:${backend.server.port()}/api/catalog/entities?filter=kind=component,kind=resource`;
  const headers = { Authorization: mockCredentials.user.header() };
  const expected = YAML.parseAllDocuments(
    readFileSync(clusterYamlPath, 'utf8'),
  ).length;
  const normalise = (items: Listed[]) =>
    items
      .map(e => {
        const c = structuredClone(e) as Listed & {
          relations?: { type: string; targetRef: string }[];
        };
        delete c.metadata.uid;
        delete c.metadata.etag;
        c.relations?.sort((x, y) =>
          `${x.type} ${x.targetRef}`.localeCompare(`${y.type} ${y.targetRef}`),
        );
        return c;
      })
      .sort((x, y) => x.metadata.name.localeCompare(y.metadata.name));
  // Poll until every document is processed and the served set stops changing (relations settle last).
  let previous = '';
  let clean: Listed[] = [];
  for (let i = 0; i < 240; i++) {
    clean = normalise(
      (await (await fetch(url, { headers })).json()) as Listed[],
    );
    const current = JSON.stringify(clean);
    if (clean.length >= expected && current === previous) break;
    previous = current;
    await new Promise(r => setTimeout(r, 500));
  }
  await backend.stop();
  return { clean };
}

describe('portal catalog with the live inventory flag off', () => {
  it('starts, and serves exactly the same entities as a catalog without the module', async () => {
    const baseline = await serve([]);
    const flagAbsent = await serve([catalogModuleOnemindLiveInventory]);
    const flagFalse = await serve([catalogModuleOnemindLiveInventory], {
      inventory: { live: { enabled: false } },
    });
    expect(baseline.clean).toHaveLength(449);
    expect(flagAbsent.clean).toEqual(baseline.clean);
    expect(flagFalse.clean).toEqual(baseline.clean);
    const stamped = baseline.clean.filter(e =>
      JSON.stringify(e).includes('inventory-owner'),
    );
    expect(stamped).toHaveLength(0);
  });
});
