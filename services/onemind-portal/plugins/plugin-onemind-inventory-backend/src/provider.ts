import { readFile } from 'node:fs/promises';
import type { Entity } from '@backstage/catalog-model';
import type {
  EntityProvider,
  EntityProviderConnection,
} from '@backstage/plugin-catalog-node';
import type { LoggerService } from '@backstage/backend-plugin-api';
import YAML from 'yaml';
import { isExcludedNamespace, type LiveInventoryConfig } from './config';
import {
  A,
  PROVIDER_NAME,
  buildWorkloadEntity,
  entityName,
  workloadKey,
  type Workload,
} from './entities';
import { reconcile, type OwnedRecord } from './ledger';
import type { ObservationStore } from './observation';
import type { InventoryMapping } from './mapping.generated';

/** What the reviewed catalog file already defines, so a workload is never published twice. */
export interface Curated {
  names: Set<string>;
  titlesBySystem: Map<string, Set<string>>;
}

/** Reads the workload entities of the reviewed catalog file (catalog/cluster.yaml). */
export async function readCurated(file: string): Promise<Curated> {
  const text = await readFile(file, 'utf8');
  const curated: Curated = { names: new Set(), titlesBySystem: new Map() };
  for (const doc of YAML.parseAllDocuments(text)) {
    const e = doc.toJSON() as Entity | null;
    if (!e || e.kind !== 'Component') continue;
    if (e.metadata.annotations?.[A.source] !== 'kubernetes') continue;
    curated.names.add(e.metadata.name);
    const system = String(e.spec?.system ?? '');
    const titles = curated.titlesBySystem.get(system) ?? new Set<string>();
    if (e.metadata.title) titles.add(e.metadata.title);
    curated.titlesBySystem.set(system, titles);
  }
  return curated;
}

/** Turns catalog entities this provider published earlier back into ledger records. */
export function toOwnedRecords(entities: Entity[]): OwnedRecord[] {
  const out: OwnedRecord[] = [];
  for (const e of entities) {
    const a = e.metadata.annotations ?? {};
    if (a[A.owner] !== PROVIDER_NAME || !a[A.namespace] || !a[A.workload])
      continue;
    out.push({
      kind: a[A.workloadKind] ?? 'Deployment',
      namespace: a[A.namespace],
      name: a[A.workload],
      selector: a[A.selector] ?? '',
      observedAt: a[A.observedAt] ?? new Date(0).toISOString(),
      removedAt: a[A.removedAt],
    });
  }
  return out;
}

export interface ProviderDeps {
  config: LiveInventoryConfig;
  logger: Pick<LoggerService, 'info' | 'warn' | 'error'>;
  store: ObservationStore;
  listWorkloads: () => Promise<Workload[]>;
  readOwned: () => Promise<OwnedRecord[]>;
  readCurated: () => Promise<Curated>;
  mapping?: InventoryMapping;
  now?: () => Date;
}

export interface RefreshSummary {
  observed: number;
  tracked: number;
  published: number;
  removed: number;
  added: string[];
  dropped: string[];
  skippedCurated: number;
}

/**
 * Publishes the cluster workloads that the reviewed catalog file does not define yet, and keeps
 * removed ones visible (lifecycle "removed") for the retention window before dropping them.
 * One full mutation per run, so anything not emitted is removed from the catalog.
 */
export class LiveInventoryProvider implements EntityProvider {
  private connection?: EntityProviderConnection;

  constructor(private readonly deps: ProviderDeps) {}

  getProviderName() {
    return PROVIDER_NAME;
  }

  async connect(connection: EntityProviderConnection) {
    this.connection = connection;
  }

  /** Observe the cluster only: updates the freshness the processor stamps, publishes nothing. */
  async observe(): Promise<Workload[]> {
    const { config, store } = this.deps;
    const now = (this.deps.now ?? (() => new Date()))();
    const tracked = (await this.deps.listWorkloads()).filter(
      w => !isExcludedNamespace(w.namespace, config),
    );
    // A cluster with no workloads at all means the read failed or is partial; never act on it.
    if (!tracked.length)
      throw new Error(
        'Cluster observation returned no workloads; refusing to act on it',
      );
    store.set(now.toISOString(), tracked.map(workloadKey));
    return tracked;
  }

  async refresh(): Promise<RefreshSummary> {
    if (!this.connection)
      throw new Error('Live inventory provider is not connected');
    const { config, logger } = this.deps;
    const now = (this.deps.now ?? (() => new Date()))();
    const tracked = await this.observe();
    const [curated, previous] = await Promise.all([
      this.deps.readCurated(),
      this.deps.readOwned(),
    ]);

    const seen = new Set<string>();
    const candidates: Workload[] = [];
    let skippedCurated = 0;
    for (const w of tracked) {
      const name = entityName(w.namespace, w.name);
      if (curated.names.has(name)) {
        skippedCurated++;
        continue;
      }
      if (seen.has(name)) {
        logger.warn(
          `Live inventory: ${w.kind} ${workloadKey(
            w,
          )} shares its entity name with another workload; skipped`,
        );
        continue;
      }
      seen.add(name);
      candidates.push(w);
    }

    const ledger = reconcile({
      candidates,
      previous: previous.filter(
        p => !curated.names.has(entityName(p.namespace, p.name)),
      ),
      now,
      retentionMs: config.removedRetentionHours * 3_600_000,
    });
    const takenTitles = new Map(
      [...curated.titlesBySystem].map(([k, v]) => [k, new Set(v)]),
    );
    const entities = ledger.entries.map(entry =>
      buildWorkloadEntity(entry.workload, {
        clusterName: config.clusterName,
        observedAt: entry.observedAt,
        state: entry.state,
        removedAt: entry.removedAt,
        mapping: this.deps.mapping,
        takenTitles,
      }),
    );
    await this.connection.applyMutation({
      type: 'full',
      entities: entities.map(entity => ({
        entity,
        locationKey: `provider:${PROVIDER_NAME}`,
      })),
    });
    const removed = ledger.entries.filter(e => e.state === 'removed').length;
    logger.info(
      `Live inventory: ${tracked.length} workloads tracked, ${skippedCurated} defined in ${config.curatedFile}, ` +
        `${entities.length} published (${ledger.added.length} new, ${removed} removed, ${ledger.dropped.length} dropped)`,
    );
    return {
      observed: tracked.length,
      tracked: tracked.length,
      published: entities.length,
      removed,
      added: ledger.added,
      dropped: ledger.dropped,
      skippedCurated,
    };
  }
}
