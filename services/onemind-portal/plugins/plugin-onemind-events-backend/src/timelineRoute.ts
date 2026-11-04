import express from 'express';
import type { HttpAuthService, UserInfoService } from '@backstage/backend-plugin-api';
import { ADMIN_GROUP_REF } from './recipients';
import type { Store } from './store';

// Read-only activity timeline (successful deployments and backups are recorded here instead
// of the inbox). Administrators only; service and anonymous callers are refused.
export function createTimelineRouter(deps: {
  httpAuth: HttpAuthService;
  userInfo: UserInfoService;
  store: Store;
}) {
  const routes = express.Router();
  routes.get('/activity/timeline', async (req, res, next) => {
    try {
      const credentials = await deps.httpAuth.credentials(req, { allow: ['user'] });
      const info = await deps.userInfo.getUserInfo(credentials);
      if (!info.ownershipEntityRefs.includes(ADMIN_GROUP_REF)) {
        res.sendStatus(403);
        return;
      }
      const limit = Number.parseInt(String(req.query.limit ?? '50'), 10);
      res.json({
        entries: await deps.store.timeline.recent(Number.isFinite(limit) ? limit : 50),
        note: 'Successful operations are recorded here and summarised in the daily digest. They are not notifications.',
      });
    } catch (error) {
      next(error);
    }
  });
  return routes;
}
