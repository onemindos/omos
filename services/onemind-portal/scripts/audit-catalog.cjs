// Read-only catalog quality report. Input must be a complete /entities JSON array,
// not one page, a source YAML file, or an inventory snapshot presented as live state.
const fs = require('node:fs');

const ref = entity =>
  `${entity.kind}:${entity.metadata.namespace || 'default'}/${
    entity.metadata.name
  }`.toLowerCase();

function auditCatalog(entities, now = new Date(), staleDays = 7) {
  if (
    !Array.isArray(entities) ||
    !Number.isFinite(now.getTime()) ||
    !(staleDays > 0)
  ) {
    throw new Error(
      'Expected a complete entity array, valid observation time and positive stale-days',
    );
  }
  const identities = new Set();
  const duplicateRefs = new Set();
  const kinds = {};
  for (const entity of entities) {
    if (!entity?.kind || !entity.metadata?.name)
      throw new Error('Invalid catalog entity');
    const id = ref(entity);
    if (identities.has(id)) duplicateRefs.add(id);
    identities.add(id);
    kinds[entity.kind] = (kinds[entity.kind] || 0) + 1;
  }
  const missingRelations = new Map();
  const launcherTargets = new Map();
  const staleInventory = [],
    invalidInventoryDates = [],
    reportedErrors = [],
    orphans = [];
  for (const entity of entities) {
    const id = ref(entity),
      annotations = entity.metadata.annotations || {};
    for (const relation of entity.relations || []) {
      const target = relation.targetRef?.toLowerCase();
      if (target && !identities.has(target)) {
        if (!missingRelations.has(target))
          missingRelations.set(target, new Set());
        missingRelations.get(target).add(id);
      }
    }
    const observedAt =
      annotations['onemind.io/inventory-observed-at'] ||
      annotations['onemind.io/observed-at'];
    if (observedAt) {
      const timestamp = Date.parse(observedAt);
      if (!Number.isFinite(timestamp) || timestamp > now.getTime())
        invalidInventoryDates.push(id);
      else if (now.getTime() - timestamp > staleDays * 86400000)
        staleInventory.push({ ref: id, observedAt });
    }
    if (annotations['backstage.io/orphan'] === 'true') orphans.push(id);
    if (
      entity.status?.items?.some(item =>
        ['error', 'warning'].includes(item.level),
      )
    )
      reportedErrors.push(id);
    if (entity.kind.toLowerCase() !== 'component') continue;
    const launcher =
      annotations['onemind.io/launcher-url'] ||
      entity.metadata.links?.find(link => link.title === 'Open application')
        ?.url;
    if (launcher) {
      if (!launcherTargets.has(launcher))
        launcherTargets.set(launcher, new Set());
      launcherTargets.get(launcher).add(id);
    }
  }
  return {
    auditedAt: now.toISOString(),
    scope:
      'Catalog metadata only; not service health, link reachability, deployment readiness or permission verification',
    entityCount: entities.length,
    kinds,
    duplicateRefs: [...duplicateRefs].sort(),
    missingRelations: [...missingRelations].map(([target, sources]) => ({
      target,
      sources: [...sources].sort(),
    })),
    // Do not print URLs: they may contain personal paths, signed queries or fragments.
    sharedLauncherCandidates: [...launcherTargets.values()]
      .filter(ids => ids.size > 1)
      .map(ids => [...ids].sort()),
    staleAfterDays: staleDays,
    staleInventory,
    invalidInventoryDates,
    reportedErrors,
    orphans,
    note: 'Shared launchers are review candidates, not proof of duplicate services. Missing targets may be permission-filtered; confirm against an authorized complete view. Never fabricate observations or delete entities automatically.',
  };
}

if (require.main === module) {
  try {
    const file = process.argv[2];
    if (!file || process.argv.length !== 3)
      throw new Error(
        'Usage: node scripts/audit-catalog.cjs <complete-entities.json|->',
      );
    const report = auditCatalog(
      JSON.parse(fs.readFileSync(file === '-' ? 0 : file, 'utf8')),
    );
    console.log(JSON.stringify(report, null, 2));
    // Metadata errors gate acceptance; overlap/freshness require explicit review.
    if (
      report.duplicateRefs.length ||
      report.missingRelations.length ||
      report.reportedErrors.length ||
      report.orphans.length ||
      report.invalidInventoryDates.length
    )
      process.exitCode = 1;
  } catch (error) {
    console.error(error.message);
    process.exitCode = 2;
  }
}
module.exports = { auditCatalog };
