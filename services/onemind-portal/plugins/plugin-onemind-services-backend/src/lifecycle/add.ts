/* eslint-disable no-nested-ternary */
import {
  addArgoMapping,
  addPreparationExceptions,
  addPillarOverride,
  addPlacement,
  addStructureApp,
  appendCatalogDocument,
  appendSection,
  bumpCoverageTestCount,
  bumpStructureTestCount,
  PATHS,
  systemNames,
  systemPillar,
} from './catalogEdits';
import type { RepoReader } from './types';

export type Hosting = 'kubernetes' | 'digitalocean' | 'saas' | 'home-device';

export type AddParams = {
  name: string;
  title: string;
  description: string;
  hosting: Hosting;
  pillar: 'so' | 'le' | 'hp' | 'ge';
  /** `system:default/geospatial` or `geospatial`. */
  system: string;
  /** `group:default/onemind-admins`. */
  owner: string;
  lifecycle?: 'experimental' | 'production';
  url?: string;
  image?: string;
  port?: number;
  access?: 'tailnet' | 'cloudflare' | 'internal';
  dataPlan?: 'none' | 'volsync' | 'pg-dump' | 'provider';
  healthCheck?: 'probe' | 'url' | 'none';
  healthUrl?: string;
};

export type AddResult = {
  /** Repo-relative path to full new content, written into the scaffolder workspace. */
  files: Map<string, string>;
  body: string;
  warnings: string[];
};

const NAME = /^[a-z][a-z0-9-]{2,39}$/;
const HOSTING_ANNOTATION: Record<Hosting, string> = {
  kubernetes: 'digitalocean-nyc3',
  digitalocean: 'digitalocean-nyc3',
  saas: 'external-provider',
  'home-device': 'private-mesh',
};

export class LifecycleInputError extends Error {
  readonly name = 'InputError';
}

const bare = (ref: string) => ref.replace(/^system:([^/]+\/)?/, '');
const isoDate = (now: Date) => now.toISOString().slice(0, 10);

function exposure(p: AddParams) {
  if (p.hosting !== 'kubernetes') return 'External or device-local; no cluster ingress';
  if (p.access === 'tailnet') return `Tailnet only (\`https://${p.name}.tail717ba0.ts.net\`)`;
  if (p.access === 'cloudflare')
    return 'Public via Cloudflare (nginx ingress with the Cloudflare allowlist)';
  return 'Internal (ClusterIP only)';
}

function stateText(p: AddParams) {
  switch (p.dataPlan) {
    case 'volsync':
      return 'PVC on `do-block-storage-retain`; daily VolSync snapshot to restic. Restore drill required before real use';
    case 'pg-dump':
      return 'Database; encrypted dump CronJob required before real use';
    case 'provider':
      return 'Data lives with the provider; export path to be recorded';
    default:
      return 'No persistent state';
  }
}

export function stackSection(p: AddParams, now: Date) {
  const system = bare(p.system);
  const row =
    p.hosting === 'kubernetes'
      ? `| 2 | \`platform-${p.name}\` | \`${p.name}\` | ${stateText(p)} | ${exposure(p)} | Manual sync; not live |`
      : `| 2 | \`${p.name}\` (${p.hosting}) | none | ${stateText(p)} | ${exposure(p)} | Proposed; not verified |`;
  return `## ${p.title} preparation — ${isoDate(now)}

Prepared source only, opened from the Backstage "Add app" template (System \`${system}\`, pillar ${p.pillar.toUpperCase()}). ${
    p.hosting === 'kubernetes'
      ? 'The Application is manual-sync: nothing is deployed until a person syncs it after merge.'
      : 'No cluster resources are created by this entry.'
  } Confirm the tier before merge.

| Tier | Application | Namespace | State / recovery | Exposure | Status |
| --- | --- | --- | --- | --- | --- |
${row}`;
}

export function catalogEntity(p: AddParams) {
  const annotations: Record<string, string> = {
    'onemind.io/hosting': HOSTING_ANNOTATION[p.hosting],
    'onemind.io/inventory-review': 'proposed',
    'onemind.io/pillar-review': 'proposed',
  };
  if (p.healthCheck === 'url' && p.healthUrl)
    annotations['onemind.io/health-check'] = p.healthUrl;
  const metadata: Record<string, unknown> = {
    name: p.name,
    title: p.title,
    description: p.description,
    tags: ['proposed', p.hosting],
    annotations,
  };
  if (p.url) metadata.links = [{ url: p.url, title: 'Open application' }];
  return {
    apiVersion: 'backstage.io/v1alpha1',
    kind: 'Component',
    metadata,
    spec: {
      type: 'website',
      lifecycle: p.lifecycle ?? 'experimental',
      owner: p.owner,
      system: bare(p.system),
    },
  };
}

function checklist(p: AddParams, argoCount?: { from: number; to: number }) {
  const k8s = p.hosting === 'kubernetes';
  const lines: string[] = [];
  lines.push(
    '- [ ] **Regenerate the catalog** from `services/onemind-portal`: `node scripts/curate-structure.cjs`, then `node scripts/test-catalog-structure.cjs` and `node scripts/test-coverage.cjs`. Commit both catalog copies and any generated runbook page. The generated files are not hand-edited by this template.',
  );
  if (k8s) {
    lines.push(
      argoCount
        ? `- [ ] The portal tests hard-code the Argo application count; this PR moved it from ${argoCount.from} to ${argoCount.to}. Confirm it equals the number of Application files in Git (\`platform/argocd\` plus the ecosystem directories).`
        : '- [ ] The portal tests hard-code the Argo application count (`scripts/test-catalog-structure.cjs`, `scripts/test-coverage.cjs`). The template could not update them; do it by hand.',
      `- [ ] Image \`${p.image ?? '(missing)'}\` exists in the registry for the node architecture and is pinned by digest (never \`:latest\`).`,
      `- [ ] Three preparation exceptions for \`${p.name}\` were added to \`catalog/coverage-exceptions.yaml\` (no live workload yet). After the first sync, refresh the inventory (\`node scripts/enrich-inventory.cjs\`) and delete them: the audit fails on a stale exception.`,
      '- [ ] Requests and limits fit the ecosystem pool; the hub-core pool is memory-saturated. Do not raise requests without a capacity check.',
      `- [ ] Secrets: seal real values with \`kubeseal\` (see CLAUDE.md, Secrets) into \`platform/${p.name}/${p.name}-secrets.sealedsecret.yaml\` and delete \`${p.name}.secret.example\`. Never commit plaintext, not even test values.`,
      `- [ ] Application \`platform-${p.name}\` is manual-sync with no \`resources-finalizer\`. Before enabling automated sync, decide the finalizer and prune policy against \`docs/gitops/APP-FINALIZERS.md\`; a finalizer on a stateful app cascades on delete.`,
      '- [ ] Namespace carries `onemind.io/layer` and Pod Security `restricted`; the Namespace manifest has `Prune=false,Delete=false`.',
    );
    if (p.dataPlan === 'volsync') {
      lines.push(
        `- [ ] Backup: \`platform/${p.name}/backup.yaml\` has a non-empty \`trigger.schedule\` (a ReplicationSource with an empty \`trigger\` is never scheduled). The claim uses \`do-block-storage-retain\`; the source and cache use \`do-block-storage\`. Seal the restic repository Secret \`${p.name}-data-restic\` (pattern: \`platform/karakeep/backup-data.sealedsecret.yaml\`) and run an isolated restore before real data.`,
      );
    } else if (p.dataPlan === 'pg-dump') {
      lines.push(
        '- [ ] Backup: add a dump CronJob under `platform/backups` or use CNPG; verify one restore.',
      );
    }
    if (p.access === 'cloudflare') {
      lines.push(
        '- [ ] Public hostname: create the Cloudflare-proxied A record to the ingress-nginx LB, and confirm the source-range annotation lists the current Cloudflare ranges.',
      );
    } else if (p.access === 'tailnet') {
      lines.push(
        `- [ ] Tailnet hostname \`${p.name}\` is not claimed by another Tailscale device or Ingress.`,
      );
    }
    lines.push(
      '- [ ] After merge a person syncs the Application (the portal never syncs Argo) and checks pods, the health endpoint and one backup run.',
    );
  } else {
    lines.push(
      '- [ ] Record the owner account and billing context, SSO and who holds admin. A URL alone does not prove hosting.',
      '- [ ] Record the data export path and where the export is stored. Provider-side retention is not a backup.',
    );
  }
  lines.push(
    `- [ ] \`docs/STACK.md\`: confirm the tier of the appended section for ${p.title}.`,
    '- [ ] If the app publishes to NATS, subjects follow the locked grammar in `docs/NATS-DOCTRINE.md`.',
  );
  return lines;
}

export function renderAddBody(p: AddParams, argoCount?: { from: number; to: number }) {
  const rows = [
    ['Name', `\`${p.name}\``],
    ['Hosting', p.hosting],
    ['Pillar / System', `${p.pillar.toUpperCase()} / \`${bare(p.system)}\``],
    ['Owner', `\`${p.owner}\``],
    ['Access', exposure(p)],
    ['Data / backup', stateText(p)],
    [
      'Health check',
      p.healthCheck === 'url' && p.healthUrl
        ? p.healthUrl
        : p.healthCheck === 'probe'
        ? 'Readiness and liveness probes on the Deployment'
        : 'None declared (shown in the stack hygiene report)',
    ],
  ];
  return `## Add app: ${p.title}

Opened by the Backstage "Add app" template as a draft. Nothing changes live until this is reviewed and merged${
    p.hosting === 'kubernetes' ? ', and a person then syncs the Application' : ''
  }. The portal reads the repository and the cluster read-only and runs no kubectl, Argo sync or DigitalOcean call.

${p.description}

| | |
| --- | --- |
${rows.map(([k, v]) => `| ${k} | ${v} |`).join('\n')}

### Before merging

${checklist(p, argoCount).join('\n')}
`;
}

export async function prepareAdd(
  p: AddParams,
  repo: RepoReader,
  now: Date = new Date(),
): Promise<AddResult> {
  if (!NAME.test(p.name))
    throw new LifecycleInputError(`Invalid app name: ${p.name}`);
  const files = new Map<string, string>();
  const warnings: string[] = [];
  const read = async (path: string) => {
    const text = await repo.readFile(path);
    if (text === undefined) throw new Error(`${path} could not be read from the repository`);
    return text;
  };

  let structure = await read(PATHS.structure);
  const system = bare(p.system);
  if (!systemNames(structure).includes(system))
    throw new LifecycleInputError(
      `System ${system} is not defined in catalog/structure.yaml`,
    );
  const ref = `component:default/${p.name}`;

  if (p.hosting === 'kubernetes') {
    structure = addStructureApp(structure, {
      name: p.name,
      title: p.title,
      system,
      description: p.description,
    });
    const map = addArgoMapping(
      await read(PATHS.workloadMap),
      `platform-${p.name}`,
      ref,
    );
    files.set(PATHS.workloadMap, map);
    files.set(
      PATHS.coverageExceptions,
      addPreparationExceptions(await read(PATHS.coverageExceptions), p.name, p.title),
    );
  } else {
    structure = addPlacement(structure, p.name, system);
    for (const path of PATHS.apps) {
      files.set(path, appendCatalogDocument(await read(path), catalogEntity(p)));
    }
  }
  if (systemPillar(structure, system) !== p.pillar) {
    structure = addPillarOverride(structure, ref, p.pillar);
  }
  files.set(PATHS.structure, structure);

  let argoCount: { from: number; to: number } | undefined;
  if (p.hosting === 'kubernetes') {
    const structureTest = await read(PATHS.structureTest);
    const coverageTest = await read(PATHS.coverageTest);
    const a = bumpStructureTestCount(structureTest, 1);
    const b = bumpCoverageTestCount(coverageTest, 1);
    if (a && b) {
      files.set(PATHS.structureTest, a);
      files.set(PATHS.coverageTest, b);
      const from = Number(/inGit\.size, (\d+),/.exec(structureTest)![1]);
      argoCount = { from, to: from + 1 };
    } else {
      warnings.push('Argo application count in the portal tests was not updated');
    }
  }
  files.set(PATHS.stack, appendSection(await read(PATHS.stack), stackSection(p, now)));
  return { files, body: renderAddBody(p, argoCount), warnings };
}
