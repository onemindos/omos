import YAML from 'yaml';
import { AddParams, prepareAdd } from './add';
import { normalizeAddInput } from './actions';
import { localRepo, localSnapshot } from './testRepo';

const NOW = new Date('2026-11-05T12:00:00Z');
const digest = `ghcr.io/onemindos/demo:1.0.0@sha256:${'b'.repeat(64)}`;

const k8s = (over: Partial<AddParams> = {}): AddParams => ({
  name: 'demo-app',
  title: 'Demo app',
  description: 'A demonstration application for the lifecycle template.',
  hosting: 'kubernetes',
  pillar: 'so',
  system: 'system:default/geospatial',
  owner: 'group:default/onemind-admins',
  image: digest,
  port: 8080,
  access: 'tailnet',
  dataPlan: 'volsync',
  healthCheck: 'probe',
  ...over,
});

const count = (text: string) => Number(/inGit\.size, (\d+),/.exec(text)![1]);

describe('prepareAdd', () => {
  it('Kubernetes: adds the mapping, the structure entry, the count bump and the STACK section', async () => {
    const result = await prepareAdd(k8s(), localRepo(), NOW);
    const paths = [...result.files.keys()].sort();
    expect(paths).toEqual([
      'docs/STACK.md',
      'services/onemind-portal/catalog/coverage-exceptions.yaml',
      'services/onemind-portal/catalog/structure.yaml',
      'services/onemind-portal/catalog/workload-map.yaml',
      'services/onemind-portal/scripts/test-catalog-structure.cjs',
      'services/onemind-portal/scripts/test-coverage.cjs',
    ]);
    const map = YAML.parse(result.files.get('services/onemind-portal/catalog/workload-map.yaml')!);
    expect(map.argoApplications['platform-demo-app']).toBe('component:default/demo-app');
    const structure = YAML.parse(result.files.get('services/onemind-portal/catalog/structure.yaml')!);
    expect(structure.apps.find((a: any) => a.name === 'demo-app')).toMatchObject({ system: 'geospatial', title: 'Demo app' });

    // The portal tests hard-code the Argo application count: it moves up by exactly one.
    const before = localSnapshot().get('services/onemind-portal/scripts/test-catalog-structure.cjs')!;
    const after = result.files.get('services/onemind-portal/scripts/test-catalog-structure.cjs')!;
    expect(count(after)).toBe(count(before) + 1);
    expect(result.files.get('services/onemind-portal/scripts/test-coverage.cjs')).toContain(
      `argoApplicationsLinked, ${count(before) + 1}`,
    );
    expect(result.files.get('docs/STACK.md')).toMatch(/## Demo app preparation — 2026-11-05[\s\S]*platform-demo-app/);
    // The review checklist covers the new-service checklist and the lessons.
    expect(result.body).toMatch(/moved it from \d+ to \d+/);
    expect(result.body).toMatch(/never scheduled/);
    expect(result.body).toMatch(/do-block-storage-retain/);
    expect(result.body).toMatch(/kubeseal/);
    expect(result.body).toMatch(/curate-structure\.cjs/);
    expect(result.body).toMatch(/portal never syncs Argo/);
  });

  it('Kubernetes: a pillar that differs from the System gets an override', async () => {
    const result = await prepareAdd(k8s({ pillar: 'ge' }), localRepo(), NOW);
    const structure = YAML.parse(result.files.get('services/onemind-portal/catalog/structure.yaml')!);
    expect(structure.pillarOverrides['component:default/demo-app']).toBe('ge');
  });

  it('SaaS: appends one catalog document to both copies and no Argo application', async () => {
    const result = await prepareAdd(
      { name: 'acme-crm', title: 'Acme CRM', description: 'A hosted CRM used for partner tracking.', hosting: 'saas', pillar: 'ge', system: 'productivity-saas', owner: 'group:default/onemind-admins', url: 'https://crm.example.com', dataPlan: 'provider' },
      localRepo(),
      NOW,
    );
    const apps = ['services/onemind-portal/catalog/apps.yaml', 'platform/backstage/catalog/apps.yaml'];
    for (const path of apps) {
      const docs = YAML.parseAllDocuments(result.files.get(path)!).map(d => d.toJSON());
      const added = docs[docs.length - 1];
      expect(added).toMatchObject({
        kind: 'Component',
        metadata: { name: 'acme-crm', links: [{ url: 'https://crm.example.com', title: 'Open application' }] },
        spec: { owner: 'group:default/onemind-admins', system: 'productivity-saas' },
      });
      expect(added.metadata.annotations['onemind.io/hosting']).toBe('external-provider');
    }
    expect(result.files.get(apps[0])).toBe(result.files.get(apps[1]));
    expect([...result.files.keys()]).not.toContain('services/onemind-portal/catalog/workload-map.yaml');
    expect(result.files.has('services/onemind-portal/scripts/test-catalog-structure.cjs')).toBe(false);
    const structure = YAML.parse(result.files.get('services/onemind-portal/catalog/structure.yaml')!);
    expect(structure.placement['acme-crm']).toBe('productivity-saas');
    expect(result.body).toMatch(/Record the data export path/);
    expect(result.body).not.toMatch(/kubeseal|ReplicationSource/);
  });

  it('rejects an unknown System and a bad name', async () => {
    await expect(prepareAdd(k8s({ system: 'nope' }), localRepo(), NOW)).rejects.toThrow(/not defined/);
    await expect(prepareAdd(k8s({ name: 'Bad Name' }), localRepo(), NOW)).rejects.toThrow(/Invalid app name/);
  });

  it('treats empty template parameters as not provided', () => {
    expect(normalizeAddInput({ name: 'x', url: '', port: '8080', image: undefined })).toEqual({ name: 'x', port: 8080 });
  });
});
