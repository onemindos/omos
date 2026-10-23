import {
  createBackendModule,
  coreServices,
} from '@backstage/backend-plugin-api';
import { policyExtensionPoint } from '@backstage/plugin-permission-node/alpha';
import { catalogServiceRef } from '@backstage/plugin-catalog-node';
import {
  AuthzProfileService,
  AuthzProfileCache,
} from '@openchoreo/backstage-plugin-permission-backend-module-openchoreo-policy';
import { openchoreoPermissions } from '@openchoreo/backstage-plugin-common';
import { openChoreoTokenServiceRef } from '@openchoreo/openchoreo-auth';
import express from 'express';
import { createOneMindPolicy, isOneMindAdmin } from './policy';
import { createUpstreamRouter } from './upstreamRouter';
import { createSharedAuthzCache } from './sharedCache';

export default createBackendModule({
  pluginId: 'permission',
  moduleId: 'onemind-policy',
  register(reg) {
    reg.registerInit({
      deps: {
        policy: policyExtensionPoint,
        config: coreServices.rootConfig,
        logger: coreServices.logger,
        database: coreServices.database,
        httpRouter: coreServices.httpRouter,
        auth: coreServices.auth,
        httpAuth: coreServices.httpAuth,
        userInfo: coreServices.userInfo,
        catalogService: catalogServiceRef,
        tokenService: openChoreoTokenServiceRef,
      },
      async init({
        policy,
        config,
        logger,
        database,
        httpRouter,
        auth,
        httpAuth,
        userInfo,
        catalogService,
        tokenService,
      }) {
        // Production must never fall through to upstream's allow-all mode.
        if (
          config.getBoolean('openchoreo.features.auth.enabled') !== true ||
          config.getBoolean('openchoreo.features.authz.enabled') !== true
        ) {
          throw new Error(
            'OneMind requires OpenChoreo authentication and authorization',
          );
        }
        const authzService = new AuthzProfileService({
          baseUrl: config.getString('openchoreo.baseUrl'),
          logger,
          cache: new AuthzProfileCache(
            await createSharedAuthzCache(await database.getClient()),
          ),
        });
        // Sign-in no longer obtains an OpenChoreo token or caches OpenChoreo capabilities, so
        // the upstream OpenChoreo policy has nothing to evaluate. The OneMind rule applies
        // alone: administrators allowed, other signed-in users read-only, anonymous denied.
        policy.setPolicy(createOneMindPolicy());
        httpRouter.use(express.json());
        httpRouter.use(async (req, res, next) => {
          try {
            if (req.path === '/cache-capabilities') {
              const credentials = await httpAuth.credentials(req, {
                allow: ['service'],
              });
              if (credentials.principal.subject !== 'plugin:auth') {
                res.status(403).json({
                  error: 'Only the sign-in backend may cache capabilities',
                });
                return;
              }
            }
            if (req.path === '/evaluate-with-context') {
              const credentials = await httpAuth.credentials(req, {
                allow: ['user'],
              });
              const info = await userInfo.getUserInfo(credentials);
              const permission = openchoreoPermissions.find(
                p => p.name === req.body?.permissionName,
              );
              if (
                !permission ||
                (!isOneMindAdmin({ info }) &&
                  permission.attributes.action !== 'read')
              ) {
                res.json({ allowed: false });
                return;
              }
            }
            next();
          } catch (error) {
            next(error);
          }
        });
        httpRouter.use(
          await createUpstreamRouter({
            authzService,
            logger,
            auth,
            httpAuth,
            catalogService,
            tokenService,
          }),
        );
      },
    });
  },
});
