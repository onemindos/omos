// Transactional discovery: partial failures never replace the reviewed inventory.
const fs = require('fs');
const cp = require('child_process');
const path = require('path');
const YAML = require('yaml');
const { missingDescription } = require('./inventory-status.cjs');
process.chdir(path.resolve(__dirname, '..'));
const files = ['catalog/cluster.yaml', 'catalog/apps.yaml', 'catalog/repositories-source.json', 'catalog/runtime-source.json', '../../platform/backstage/catalog/cluster.yaml', '../../platform/backstage/catalog/apps.yaml'];
const original = new Map(files.map(file => [file, fs.existsSync(file) ? fs.readFileSync(file) : null]));
const documents = bytes => YAML.parseAllDocuments(bytes.toString()).map(doc => doc.toJSON());
const identity = entity => `${entity.kind}:${entity.metadata.namespace || 'default'}/${entity.metadata.name}`;
try {
  // enrich-inventory enforces the exact cloud context and requires all three sources.
  cp.execFileSync(process.execPath, ['scripts/enrich-inventory.cjs'], { stdio: 'inherit', timeout: 180000 });
  const previous = documents(original.get('catalog/cluster.yaml'));
  let priorProposal = [];
  try { priorProposal = documents(cp.execFileSync('git', ['show', 'origin/automation/onemind-inventory-review:services/onemind-portal/catalog/cluster.yaml'], { maxBuffer: 20 * 1024 * 1024 })); } catch {}
  const priorByIdentity = new Map(priorProposal.map(entity => [identity(entity), entity]));
  const previousIdentities = new Set(previous.map(identity));
  for (const entity of priorProposal) if (!previousIdentities.has(identity(entity)) && ['kubernetes', 'tailscale'].includes(entity.metadata.annotations?.['onemind.io/inventory-source'])) previous.push(entity);
  const current = documents(fs.readFileSync('catalog/cluster.yaml'));
  const observedAt = new Date().toISOString();
  const present = new Set(current.map(identity));
  for (const entity of current) {
    if (!['kubernetes', 'tailscale'].includes(entity.metadata.annotations?.['onemind.io/inventory-source'])) continue;
    entity.metadata.annotations['onemind.io/inventory-state'] = 'observed';
    entity.metadata.annotations['onemind.io/inventory-missing-observations'] = '0';
  }
  for (const entity of previous) {
    if (present.has(identity(entity))) continue;
    const annotations = entity.metadata.annotations || {};
    if (!['kubernetes', 'tailscale'].includes(annotations['onemind.io/inventory-source'])) {
      // Retain relationship targets and reviewed documentation even when discovery changes.
      current.push(entity); continue;
    }
    const prior = priorByIdentity.get(identity(entity))?.metadata.annotations || annotations;
    const misses = Math.max(0, Number(prior['onemind.io/inventory-missing-observations']) || 0) + 1;
    entity.metadata.annotations = { ...annotations, 'onemind.io/inventory-state': misses >= 3 ? 'retirement-review-required' : 'missing-review-required', 'onemind.io/inventory-missing-observations': String(misses), 'onemind.io/inventory-last-checked-at': observedAt };
    entity.metadata.description = missingDescription(entity.metadata.description);
    // Keep last successful observed-at: disappearance is not evidence of deletion.
    current.push(entity);
  }
  const rendered = current.map(entity => YAML.stringify(entity)).join('---\n');
  fs.writeFileSync('catalog/cluster.yaml', rendered);
  fs.writeFileSync('../../platform/backstage/catalog/cluster.yaml', rendered);
  const repositoryFile = 'catalog/repositories-source.json';
  const repositories = JSON.parse(fs.readFileSync(repositoryFile));
  const repositoryNames = new Set(repositories.map(repository => repository.name));
  for (const repository of repositories) {
    repository.inventoryState = repository.isArchived ? 'archived' : 'observed';
    repository.observedAt = observedAt;
  }
  for (const repository of JSON.parse(original.get(repositoryFile) || '[]')) {
    if (!repositoryNames.has(repository.name)) repositories.push({ ...repository, inventoryState: 'missing-review-required', lastCheckedAt: observedAt });
  }
  fs.writeFileSync(repositoryFile, JSON.stringify(repositories, null, 2) + '\n');
  cp.execFileSync(process.execPath, ['scripts/verify-inventory.cjs'], { stdio: 'inherit' });
  // Last: re-parent every entity into the reviewed SO/LE/HP/GE structure (verify and enrich re-add legacy Systems).
  cp.execFileSync(process.execPath, ['scripts/curate-structure.cjs'], { stdio: 'inherit' });
  console.log('Complete discovery staged for review. No entities deleted; no deployment or catalog API writes.');
} catch (error) {
  for (const [file, bytes] of original) {
    if (bytes === null) fs.rmSync(file, { force: true }); else fs.writeFileSync(file, bytes);
  }
  console.error('Discovery failed; previous inventory restored. Inspect source access without publishing partial results.');
  process.exitCode = 1;
}
