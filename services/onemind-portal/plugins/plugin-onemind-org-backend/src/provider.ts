import type { LoggerService } from '@backstage/backend-plugin-api';
import type {
  EntityProvider,
  EntityProviderConnection,
} from '@backstage/plugin-catalog-node';
import { PERSON_TYPES, type AuthentikSnapshot } from './authentik';
import {
  LOCATION,
  PROVIDER_NAME,
  buildOrgEntities,
  type Curated,
} from './entities';

export interface ProviderDeps {
  userGroups: string[];
  /** Named in log lines only. */
  curatedFile: string;
  logger: Pick<LoggerService, 'info' | 'warn'>;
  readSnapshot: (signal?: AbortSignal) => Promise<AuthentikSnapshot>;
  readCurated: () => Promise<Curated>;
}

export interface SyncSummary {
  users: number;
  groups: number;
  curatedUsers: number;
  curatedGroups: number;
  multiParentGroups: number;
}

/**
 * Publishes Authentik users and groups as catalog Users and Groups: one full mutation per run, so a
 * user who left the configured groups, or was deactivated, is removed at the next run. A failed or
 * implausible read throws before anything is applied, and the catalog keeps the last good state.
 */
export class AuthentikOrgProvider implements EntityProvider {
  private connection?: EntityProviderConnection;

  constructor(private readonly deps: ProviderDeps) {}

  getProviderName() {
    return PROVIDER_NAME;
  }

  async connect(connection: EntityProviderConnection) {
    this.connection = connection;
  }

  async refresh(signal?: AbortSignal): Promise<SyncSummary> {
    if (!this.connection)
      throw new Error('Authentik org provider is not connected');
    const { logger } = this.deps;
    const curated = await this.deps.readCurated();
    const snapshot = await this.deps.readSnapshot(signal);
    // Authentik answers 200 with an empty list when the token lacks a view permission. An identity
    // provider without a single active person or group is that, not an empty organisation.
    if (!snapshot.users.some(u => u.isActive && PERSON_TYPES.includes(u.type)))
      throw new Error(
        'Authentik returned no active users; check that the token may view users. Nothing was changed',
      );
    if (!snapshot.groups.length)
      throw new Error(
        'Authentik returned no groups; check that the token may view groups. Nothing was changed',
      );
    const mapping = buildOrgEntities(snapshot, {
      userGroups: this.deps.userGroups,
      curated,
    });
    for (const warning of mapping.warnings)
      logger.warn(`Authentik org sync: ${warning}`);
    signal?.throwIfAborted();
    await this.connection.applyMutation({
      type: 'full',
      entities: mapping.entities.map(entity => ({
        entity,
        locationKey: LOCATION,
      })),
    });
    logger.info(
      `Authentik org sync: ${mapping.users} users and ${mapping.groups} groups published; ` +
        `${mapping.curatedUsers} users and ${mapping.curatedGroups} groups left to ${this.deps.curatedFile}; ` +
        `${mapping.multiParentGroups} groups published without a parent (several Authentik parents)`,
    );
    return {
      users: mapping.users,
      groups: mapping.groups,
      curatedUsers: mapping.curatedUsers,
      curatedGroups: mapping.curatedGroups,
      multiParentGroups: mapping.multiParentGroups,
    };
  }
}
