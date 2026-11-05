import express from 'express';
import type { HttpAuthService } from '@backstage/backend-plugin-api';
import { createSource, Source } from './source';
import { KubePolicyReport, PolicySummary, summarizePolicyReports } from './policyReports';

/** The only cluster path this reads. It is fixed here: no caller supplies a path or query. */
export const POLICY_REPORTS_PATH = '/apis/wgpolicyk8s.io/v1alpha2/policyreports';

/** One cluster list (about 1 MB today) at most once a minute, shared by every caller. */
export const POLICY_TTL_MS = 60000;

export type PolicyReportsResponse =
  | ({ status: 'ok' | 'stale'; observedAt: string; error?: string } & PolicySummary)
  | { status: 'error'; observedAt: null; error: string };

export function createPolicySource(
  list: (path: string) => Promise<KubePolicyReport[]>,
  now: () => number = Date.now,
): Source<PolicySummary> {
  return createSource<PolicySummary>(
    'policyreports',
    async () => summarizePolicyReports(await list(POLICY_REPORTS_PATH)),
    { timeoutMs: 10000, ttlMs: POLICY_TTL_MS, now },
  );
}

/**
 * GET /policy-reports. Signed-in users only (`allow: ['user']`): an unauthenticated caller
 * gets 401 and a service or plugin credential is refused, like the other /api/onemind
 * routes. A failed source answers 200 with status "error" so the card can say "unavailable"
 * rather than showing zero violations.
 */
export function createPolicyRoutes(options: {
  httpAuth: Pick<HttpAuthService, 'credentials'>;
  source: Source<PolicySummary>;
}) {
  const routes = express.Router();
  routes.get('/policy-reports', async (req, res, next) => {
    try {
      await options.httpAuth.credentials(req, { allow: ['user'] });
      const result = await options.source.read();
      res.setHeader('Cache-Control', 'no-store');
      const body: PolicyReportsResponse =
        result.status === 'error' || !result.data
          ? { status: 'error', observedAt: null, error: result.error ?? 'Unavailable' }
          : { status: result.status, observedAt: result.observedAt!, error: result.error, ...result.data };
      res.json(body);
    } catch (e) {
      next(e);
    }
  });
  return routes;
}
