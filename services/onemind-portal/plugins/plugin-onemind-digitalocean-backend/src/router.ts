import express from 'express';
import type {
  HttpAuthService,
  UserInfoService,
} from '@backstage/backend-plugin-api';
import type { DigitalOceanService } from './service';

export const ADMIN_GROUP_REF = 'group:default/onemind-admins';

/**
 * GET /summary. A signed-in user who belongs to onemind-admins, the same membership test the
 * other administrator-only routes of this portal use (/api/onemind/costs, /calendar, /cloud).
 * Service and plugin tokens are refused. Responses are never cached by the browser.
 */
export function createRouter(deps: {
  httpAuth: Pick<HttpAuthService, 'credentials'>;
  userInfo: Pick<UserInfoService, 'getUserInfo'>;
  service: Pick<DigitalOceanService, 'get'>;
}) {
  const router = express.Router();
  router.get('/summary', async (req, res, next) => {
    try {
      const credentials = await deps.httpAuth.credentials(req, {
        allow: ['user'],
      });
      const info = await deps.userInfo.getUserInfo(credentials);
      if (!info.ownershipEntityRefs.includes(ADMIN_GROUP_REF)) {
        res.status(403).json({
          error: 'DigitalOcean requires OneMind administrator membership',
        });
        return;
      }
      res.setHeader('Cache-Control', 'no-store');
      res.json(await deps.service.get());
    } catch (error) {
      next(error);
    }
  });
  return router;
}
