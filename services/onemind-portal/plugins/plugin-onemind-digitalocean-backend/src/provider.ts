import type {
  EntityProvider,
  EntityProviderConnection,
} from '@backstage/plugin-catalog-node';
import type { LoggerService } from '@backstage/backend-plugin-api';
import { buildEntities, collectResources, PROVIDER_NAME } from './resources';
import type { DigitalOceanService } from './service';

export interface ProviderDeps {
  service: Pick<DigitalOceanService, 'refresh'>;
  logger: Pick<LoggerService, 'info' | 'warn'>;
}

/**
 * Publishes DigitalOcean resources as catalog Resources. Without a token it publishes nothing.
 * It replaces its whole set only when every resource section read successfully, so a partial
 * outage or a missing scope never deletes entities that are still real.
 */
export class DigitalOceanProvider implements EntityProvider {
  private connection?: EntityProviderConnection;

  constructor(private readonly deps: ProviderDeps) {}

  getProviderName() {
    return PROVIDER_NAME;
  }

  async connect(connection: EntityProviderConnection) {
    this.connection = connection;
  }

  async refresh(): Promise<{ published: number } | undefined> {
    if (!this.connection)
      throw new Error('DigitalOcean provider is not connected');
    const summary = await this.deps.service.refresh();
    if (summary.state === 'not-connected') return undefined;
    const { complete } = collectResources(summary);
    const required = [
      'kubernetes-cluster',
      'droplet',
      'volume',
      'snapshot',
      'load-balancer',
      'firewall',
      'vpc',
    ] as const;
    const missing = required.filter(t => !complete.has(t));
    if (missing.length) {
      this.deps.logger.warn(
        `DigitalOcean catalog refresh skipped: ${missing.join(
          ', ',
        )} could not be read`,
      );
      return undefined;
    }
    const entities = buildEntities(summary);
    await this.connection.applyMutation({
      type: 'full',
      entities: entities.map(entity => ({
        entity,
        locationKey: `provider:${PROVIDER_NAME}`,
      })),
    });
    this.deps.logger.info(
      `DigitalOcean catalog: ${entities.length} resources published`,
    );
    return { published: entities.length };
  }
}
