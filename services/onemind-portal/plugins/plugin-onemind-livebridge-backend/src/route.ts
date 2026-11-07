import type { HttpAuthService } from '@backstage/backend-plugin-api';
import express from 'express';
import type { LiveSnapshot } from './contract';

/** GET /snapshot: the current LiveSnapshot for signed-in users. Service and anonymous callers are refused. */
export function createSnapshotRouter(deps: {
  httpAuth: HttpAuthService;
  snapshot: () => LiveSnapshot;
}) {
  const router = express.Router();
  router.get('/snapshot', async (req, res, next) => {
    try {
      await deps.httpAuth.credentials(req, { allow: ['user'] });
      res.json(deps.snapshot());
    } catch (error) {
      next(error);
    }
  });
  return router;
}
