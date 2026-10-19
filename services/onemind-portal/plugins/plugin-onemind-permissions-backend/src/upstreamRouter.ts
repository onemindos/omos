import type {
  AuthService,
  HttpAuthService,
  LoggerService,
} from '@backstage/backend-plugin-api';
import type { CatalogService } from '@backstage/plugin-catalog-node';
import type { AuthzProfileService } from '@openchoreo/backstage-plugin-permission-backend-module-openchoreo-policy';
import type { OpenChoreoTokenService } from '@openchoreo/openchoreo-auth';
import type { Router } from 'express';

// v1.3.1 exports policy/services but not its ABAC router. Keep this private API
// boundary isolated and version-locked rather than copying the CEL evaluator or
// installing a second policy. Revalidate/remove this adapter on any upgrade.
export function createUpstreamRouter(options: {
  authzService: AuthzProfileService;
  logger: LoggerService;
  auth: AuthService;
  httpAuth: HttpAuthService;
  catalogService: CatalogService;
  tokenService: OpenChoreoTokenService;
}): Promise<Router> {
  const pkg =
    '@openchoreo/backstage-plugin-permission-backend-module-openchoreo-policy';
  if (require(`${pkg}/package.json`).version !== '1.3.1') {
    throw new Error(
      'OpenChoreo permission router requires reviewed version 1.3.1',
    );
  }
  return require(`${pkg}/dist/router.cjs.js`).createRouter(options);
}
