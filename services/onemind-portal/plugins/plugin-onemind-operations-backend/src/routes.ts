import express from 'express';
import type { HttpAuthService } from '@backstage/backend-plugin-api';
import { createCollector } from './collect';

type Collector = ReturnType<typeof createCollector>;

/**
 * Operational read routes. Both require a signed-in user (`allow: ['user']`), so an
 * unauthenticated caller gets 401 and a service or plugin credential is refused,
 * exactly like the other /api/onemind routes. They are read-only GETs.
 */
export function createLiveRoutes(options: {
  httpAuth: Pick<HttpAuthService, 'credentials'>;
  collector: Collector;
}) {
  const { httpAuth, collector } = options;
  const routes = express.Router();
  routes.get('/status', async (req, res, next) => {
    try {
      await httpAuth.credentials(req, { allow: ['user'] });
      res.json(await collector.status());
    } catch (e) {
      next(e);
    }
  });
  routes.get('/live', async (req, res, next) => {
    try {
      await httpAuth.credentials(req, { allow: ['user'] });
      res.setHeader('Cache-Control', 'no-store');
      res.json(await collector.live());
    } catch (e) {
      next(e);
    }
  });
  return routes;
}
