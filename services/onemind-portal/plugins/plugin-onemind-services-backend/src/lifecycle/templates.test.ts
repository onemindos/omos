/** @jest-environment node */
import { mkdtemp, readFile, readdir } from 'node:fs/promises';
import { readFileSync } from 'node:fs';
import os from 'node:os';
import path from 'node:path';
import { PassThrough } from 'node:stream';
import Ajv from 'ajv';
import YAML from 'yaml';
import { mockServices } from '@backstage/backend-test-utils';
import { ScmIntegrations } from '@backstage/integration';
import { createFetchTemplateAction } from '@backstage/plugin-scaffolder-backend';
import { createPublishGithubPullRequestAction } from '@backstage/plugin-scaffolder-backend-module-github';
import { createAddPrepareAction, createRetirePlanAction, normalizeAddInput } from './actions';
import { catalogEntities, localRepo, REPO_ROOT, stacCluster } from './testRepo';

const TEMPLATES = path.join(REPO_ROOT, 'services/onemind-portal/templates');
const template = (dir: string) =>
  YAML.parse(readFileSync(path.join(TEMPLATES, dir, 'template.yaml'), 'utf8'));
const NOW = new Date('2026-11-05T12:00:00Z');
const digest = `ghcr.io/onemindos/demo:1.0.0@sha256:${'b'.repeat(64)}`;

const integrations = ScmIntegrations.fromConfig(
  mockServices.rootConfig({ data: { integrations: { github: [{ host: 'github.com' }] } } }),
);

type PrCall = { files: Record<string, any>; draft: boolean; head: string; body: string; title: string; update: boolean };
function fakeClient() {
  const calls: PrCall[] = [];
  const client = {
    async createPullRequest(options: any) {
      calls.push({
        files: options.changes[0].files,
        draft: options.draft,
        head: options.head,
        body: options.body,
        title: options.title,
        update: options.update,
      });
      return { data: { base: { ref: 'main' }, html_url: 'https://github.com/onemindos/omos/pull/1', number: 1 } };
    },
  };
  return { calls, publish: createPublishGithubPullRequestAction({ integrations, clientFactory: async () => client as any }) };
}

function context(input: any, workspacePath: string, templateDir: string, isDryRun = false) {
  const outputs: Record<string, any> = {};
  const logger: any = { info() {}, warn() {}, error() {}, debug() {}, child: () => logger };
  return {
    outputs,
    ctx: {
      input,
      workspacePath,
      logger,
      logStream: new PassThrough(),
      isDryRun,
      output: (key: string, value: unknown) => {
        outputs[key] = value;
      },
      createTemporaryDirectory: () => mkdtemp(path.join(os.tmpdir(), 'tpl-')),
      checkpoint: async ({ fn }: any) => fn(),
      templateInfo: { baseUrl: `file://${path.join(TEMPLATES, templateDir, 'template.yaml')}`, entityRef: 'template:default/x' },
      task: { id: 'task' },
      getInitiatorCredentials: async () => ({}),
      secrets: {},
    } as any,
  };
}

const workspace = () => mkdtemp(path.join(os.tmpdir(), 'ws-'));

async function list(dir: string, base = dir): Promise<string[]> {
  const out: string[] = [];
  for (const entry of await readdir(dir, { withFileTypes: true })) {
    const full = path.join(dir, entry.name);
    if (entry.isDirectory()) out.push(...(await list(full, base)));
    else out.push(path.relative(base, full).split(path.sep).join('/'));
  }
  return out.sort();
}

const fetchTemplate = createFetchTemplateAction({
  integrations,
  reader: { readUrl: async () => { throw new Error('no network in tests'); }, readTree: async () => { throw new Error('no network in tests'); }, search: async () => { throw new Error('no network in tests'); } } as any,
});

const deps = () => ({
  repo: async () => localRepo(),
  cluster: async () => stacCluster(),
  entities: async () => catalogEntities('services/onemind-portal/catalog/apps.yaml', 'services/onemind-portal/catalog/cluster.yaml'),
  now: () => NOW,
});

/** Runs the add-app template the way the scaffolder would for the given parameters. */
async function runAdd(parameters: Record<string, any>, dryRun = false) {
  const tpl = template('app-add');
  const ws = await workspace();
  const prepare = createAddPrepareAction(deps());
  const { ctx, outputs } = context(
    Object.fromEntries(Object.keys(tpl.spec.steps[0].input).map(k => [k, parameters[k] ?? ''])),
    ws,
    'app-add',
    dryRun,
  );
  await prepare.handler(ctx);
  const render = async (dir: string, values: Record<string, any>) => {
    const c = context({ url: `./skeleton/${dir}`, values }, ws, 'app-add', dryRun);
    await fetchTemplate.handler(c.ctx);
  };
  const k8s = parameters.hosting === 'kubernetes';
  if (k8s) await render('k8s', parameters);
  if (k8s && parameters.access === 'tailnet') await render('k8s-tailnet', parameters);
  if (k8s && parameters.access === 'cloudflare') await render('k8s-cloudflare', parameters);
  if (k8s && parameters.dataPlan === 'volsync') await render('k8s-volsync', parameters);
  const { calls, publish } = fakeClient();
  const publishStep = tpl.spec.steps.find((s: any) => s.id === 'publish');
  const pc = context(
    {
      repoUrl: publishStep.input.repoUrl,
      branchName: outputs.branchName,
      targetBranchName: 'main',
      title: outputs.title,
      description: outputs.body,
      sourcePath: '.',
      draft: true,
      update: true,
      forceEmptyGitAuthor: true,
    },
    ws,
    'app-add',
    dryRun,
  );
  await publish.handler(pc.ctx);
  return { ws, calls, outputs, files: await list(ws), publishOutputs: pc.outputs };
}

const K8S = {
  name: 'demo-app', title: 'Demo app', description: 'A demonstration application for the lifecycle template.',
  hosting: 'kubernetes', pillar: 'so', system: 'system:default/geospatial', owner: 'group:default/onemind-admins',
  lifecycle: 'experimental', image: digest, port: 8080, layer: 'make-sense', access: 'tailnet', dataPlan: 'volsync',
  storage: '5Gi', backupSchedule: '30 3 * * *', healthCheck: 'probe', healthPath: '/health',
};

describe('Add app template', () => {
  it('Kubernetes dry run renders the expected files and opens a draft PR with no token', async () => {
    const run = await runAdd(K8S);
    expect(run.files).toEqual([
      'docs/STACK.md',
      'platform/argocd/platform-demo-app.yaml',
      'platform/demo-app/README.md',
      'platform/demo-app/backup.yaml',
      'platform/demo-app/data-pvc.yaml',
      'platform/demo-app/demo-app.secret.example',
      'platform/demo-app/deployment.yaml',
      'platform/demo-app/ingress.yaml',
      'platform/demo-app/namespace.yaml',
      'platform/demo-app/network.yaml',
      'platform/demo-app/service.yaml',
      'services/onemind-portal/catalog/coverage-exceptions.yaml',
      'services/onemind-portal/catalog/structure.yaml',
      'services/onemind-portal/catalog/workload-map.yaml',
      'services/onemind-portal/scripts/test-catalog-structure.cjs',
      'services/onemind-portal/scripts/test-coverage.cjs',
    ]);
    expect(run.calls).toHaveLength(1);
    expect(run.calls[0]).toMatchObject({ draft: true, update: true, head: 'add-app/demo-app' });
    expect(Object.keys(run.calls[0].files).sort()).toEqual(run.files);
    expect(run.calls[0].body).toMatch(/Before merging/);
  });

  it('renders manifests that meet the new-service checklist', async () => {
    const { ws } = await runAdd(K8S);
    const docs = async (file: string) =>
      YAML.parseAllDocuments(await readFile(path.join(ws, file), 'utf8')).map(d => {
        expect(d.errors).toEqual([]);
        return d.toJSON();
      });
    const [deployment] = await docs('platform/demo-app/deployment.yaml');
    const container = deployment.spec.template.spec.containers[0];
    expect(container.image).toMatch(/@sha256:[a-f0-9]{64}$/);
    expect(container.image).not.toMatch(/:latest/);
    expect(container.resources.requests).toEqual({ cpu: '50m', memory: '128Mi' });
    expect(container.resources.limits).toBeDefined();
    expect(container.readinessProbe.httpGet.path).toBe('/health');
    expect(deployment.spec.template.spec.nodeSelector).toEqual({ 'doks.digitalocean.com/node-pool': 'onemind-ecosystem-cpu' });
    const [ns] = await docs('platform/demo-app/namespace.yaml');
    expect(ns.metadata.labels['onemind.io/layer']).toBe('make-sense');
    const [pvc] = await docs('platform/demo-app/data-pvc.yaml');
    expect(pvc.spec.storageClassName).toBe('do-block-storage-retain');
    const backup = (await docs('platform/demo-app/backup.yaml')).find(d => d.kind === 'ReplicationSource');
    // A ReplicationSource with an empty trigger is never scheduled; the temp classes are do-block-storage.
    expect(backup.spec.trigger.schedule).toBe('30 3 * * *');
    expect(backup.spec.restic.storageClassName).toBe('do-block-storage');
    expect(backup.spec.restic.cacheStorageClassName).toBe('do-block-storage');
    const [ingress] = await docs('platform/demo-app/ingress.yaml');
    expect(ingress.spec.ingressClassName).toBe('tailscale');
    const [application] = await docs('platform/argocd/platform-demo-app.yaml');
    expect(application.metadata.finalizers).toBeUndefined();
    expect(application.spec.syncPolicy.automated).toBeUndefined();
    expect(application.spec.source.path).toBe('platform/demo-app');
    // The placeholder carries key names only.
    const secret = await readFile(path.join(ws, 'platform/demo-app/demo-app.secret.example'), 'utf8');
    expect(secret).not.toMatch(/^[A-Z_]+=\S+/m);
  });

  it('Cloudflare access uses nginx with the Cloudflare allowlist and no VolSync files without data', async () => {
    const run = await runAdd({ ...K8S, access: 'cloudflare', dataPlan: 'none', healthCheck: 'none' });
    expect(run.files).not.toContain('platform/demo-app/backup.yaml');
    expect(run.files).not.toContain('platform/demo-app/data-pvc.yaml');
    const ingress = YAML.parse(await readFile(path.join(run.ws, 'platform/demo-app/ingress.yaml'), 'utf8'));
    expect(ingress.spec.ingressClassName).toBe('nginx');
    expect(ingress.metadata.annotations['nginx.ingress.kubernetes.io/whitelist-source-range']).toMatch(/173\.245\.48\.0\/20/);
    const deployment = YAML.parse(await readFile(path.join(run.ws, 'platform/demo-app/deployment.yaml'), 'utf8'));
    expect(deployment.spec.template.spec.containers[0].readinessProbe).toBeUndefined();
    expect(deployment.spec.template.spec.volumes.map((v: any) => v.name)).toEqual(['tmp']);
  });

  it('SaaS variant produces only catalog changes', async () => {
    const run = await runAdd({
      name: 'acme-crm', title: 'Acme CRM', description: 'A hosted CRM used for partner tracking.', hosting: 'saas',
      pillar: 'ge', system: 'productivity-saas', owner: 'group:default/onemind-admins', url: 'https://crm.example.com', dataPlan: 'provider',
    });
    expect(run.files).toEqual([
      'docs/STACK.md',
      'platform/backstage/catalog/apps.yaml',
      'services/onemind-portal/catalog/apps.yaml',
      'services/onemind-portal/catalog/structure.yaml',
    ]);
    expect(run.calls[0].draft).toBe(true);
  });

  it('a scaffolder dry run reads and renders but never calls GitHub', async () => {
    const run = await runAdd(K8S, true);
    expect(run.calls).toHaveLength(0);
    expect(run.publishOutputs.remoteUrl).toBeDefined();
  });

  it('has no token or secret field and only the allowed actions', () => {
    const tpl = template('app-add');
    expect(JSON.stringify(tpl)).not.toMatch(/githubToken|ui:field: Secret|secrets\./);
    expect(tpl.spec.steps.map((s: any) => s.action)).toEqual([
      'onemind:app-add:prepare', 'fetch:template', 'fetch:template', 'fetch:template', 'fetch:template', 'publish:github:pull-request',
    ]);
    const publish = tpl.spec.steps.at(-1).input;
    expect(publish.draft).toBe(true);
    expect(publish.targetBranchName).toBe('main');
    expect(publish.token).toBeUndefined();
    expect(normalizeAddInput({ port: '' })).toEqual({});
  });
});

describe('Add app form schema', () => {
  const schema = () => template('app-add').spec.parameters[0];
  const validate = (data: any) => {
    const ajv = new Ajv({ strict: false, allErrors: true });
    const fn = ajv.compile(schema());
    return { ok: fn(data), errors: fn.errors };
  };
  const base = { name: 'demo-app', title: 'Demo app', description: 'A demonstration application.' };

  it('accepts each hosting variant with its own fields', () => {
    expect(validate({ ...base, hosting: 'kubernetes', image: digest, port: 8080, access: 'tailnet', dataPlan: 'volsync' }).ok).toBe(true);
    expect(validate({ ...base, hosting: 'saas', url: 'https://crm.example.com', dataPlan: 'provider' }).ok).toBe(true);
    expect(validate({ ...base, hosting: 'digitalocean', url: 'https://app.example.com', dataPlan: 'provider' }).ok).toBe(true);
    expect(validate({ ...base, hosting: 'home-device', dataPlan: 'none' }).ok).toBe(true);
  });

  it('rejects an image that is not pinned by digest and a SaaS app without a URL', () => {
    expect(validate({ ...base, hosting: 'kubernetes', image: 'ghcr.io/x/y:latest', port: 8080, access: 'tailnet', dataPlan: 'none' }).ok).toBe(false);
    expect(validate({ ...base, hosting: 'kubernetes', port: 8080, access: 'tailnet', dataPlan: 'none' }).ok).toBe(false);
    expect(validate({ ...base, hosting: 'saas', dataPlan: 'provider' }).ok).toBe(false);
    expect(validate({ ...base, hosting: 'kubernetes', image: digest, port: 8080, access: 'tailnet', dataPlan: 'volsync', backupSchedule: '' }).ok).toBe(false);
  });
});

describe('Retire app template', () => {
  it('on stac-browser lists its workload and opens a draft PR that removes only its files', async () => {
    const tpl = template('app-retire');
    const catalog = await deps().entities();
    const target = {
      apiVersion: 'backstage.io/v1alpha1', kind: 'Component',
      metadata: { name: 'stac-browser', title: 'STAC Browser', annotations: { 'backstage.io/kubernetes-namespace': 'stac', 'onemind.io/kubernetes-workload': 'stac-browser' }, labels: { 'onemind.io/catalog-role': 'workload' } },
      spec: { subcomponentOf: 'component:default/stac' },
    };
    const ws = await workspace();
    const plan = createRetirePlanAction({ ...deps(), entities: async () => [...catalog, target as any] });
    const { ctx, outputs } = context(
      { entityRef: 'component:default/stac-browser', retention: 'keep-last-backup', retentionDays: 30, reason: 'The STAC browser is replaced by the portal map view' },
      ws, 'app-retire',
    );
    await plan.handler(ctx);
    expect(outputs.scope).toBe('workload');
    expect(outputs.filesToDelete).toEqual(['platform/stac/browser.yaml']);
    expect(outputs.inventory.find((s: any) => s.id === 'workloads').items[0].label).toBe('Deployment stac/stac-browser');

    const { calls, publish } = fakeClient();
    const publishInput = tpl.spec.steps.find((s: any) => s.id === 'publish').input;
    expect(publishInput.filesToDelete).toBe('${{ steps.plan.output.filesToDelete }}');
    const pc = context(
      { repoUrl: publishInput.repoUrl, branchName: outputs.branchName, targetBranchName: 'main', title: outputs.title, description: outputs.body, sourcePath: '.', filesToDelete: outputs.filesToDelete, draft: true, update: true, forceEmptyGitAuthor: true },
      ws, 'app-retire',
    );
    await publish.handler(pc.ctx);
    expect(calls).toHaveLength(1);
    const { files } = calls[0];
    const deleted = Object.keys(files).filter(p => typeof files[p] === 'symbol');
    expect(deleted).toEqual(['platform/stac/browser.yaml']);
    const written = Object.keys(files).filter(p => !deleted.includes(p)).sort();
    expect(written).toEqual(['docs/operations/retirements/stac-browser-2026-11-05.md', 'platform/stac/ingress.yaml']);
    expect(calls[0].draft).toBe(true);
    expect(calls[0].body).toMatch(/After merging \(manual, in this order\)/);
    // The api, writer, sync, secrets and namespace files are never touched.
    for (const p of Object.keys(files)) expect(p).not.toMatch(/api|webodm|migrate|stac-db|namespace\.yaml/);
  });

  it('is an EntityPicker plus read-only steps, with no token field', () => {
    const tpl = template('app-retire');
    expect(tpl.spec.parameters[0].properties.app['ui:field']).toBe('EntityPicker');
    expect(tpl.spec.steps.map((s: any) => s.action)).toEqual(['onemind:app-retire:plan', 'publish:github:pull-request']);
    expect(JSON.stringify(tpl)).not.toMatch(/githubToken|ui:field: Secret/);
    expect(tpl.spec.steps[1].input.draft).toBe(true);
    expect(tpl.spec.steps[1].input.token).toBeUndefined();
  });
});

describe('No template step mutates the cluster, Argo CD or DigitalOcean', () => {
  const FORBIDDEN = /kubectl|argocd|argo-?cd app|doctl|digitalocean|helm |terraform|\bsync\b/i;
  const ALLOWED = new Set([
    'onemind:app-add:prepare', 'onemind:app-retire:plan', 'fetch:template', 'publish:github:pull-request',
  ]);
  it.each(['app-add', 'app-retire'])('%s steps use only allow-listed actions and name no mutating tool', dir => {
    const steps = template(dir).spec.steps as any[];
    for (const step of steps) {
      expect(ALLOWED.has(step.action)).toBe(true);
      // The step definitions, not the human-facing PR text, must not name a mutating tool.
      expect(JSON.stringify({ action: step.action, input: step.input })).not.toMatch(FORBIDDEN);
    }
  });

  it('the lifecycle code spawns no process and makes no write call to a cluster API', () => {
    for (const file of ['actions', 'add', 'retire', 'sources', 'hygiene', 'catalogEdits', 'yamlText']) {
      const text = readFileSync(path.join(__dirname, `${file}.ts`), 'utf8');
      expect(text).not.toMatch(/child_process|execSync|spawn\(|(?<![.\w])exec\(|kubectl |doctl |https\.request|method:\s*'(POST|PUT|PATCH|DELETE)'/);
    }
  });
});
