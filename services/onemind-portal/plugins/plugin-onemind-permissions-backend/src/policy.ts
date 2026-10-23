import { AuthorizeResult } from '@backstage/plugin-permission-common';
import type {
  PermissionPolicy,
  PolicyQueryUser,
} from '@backstage/plugin-permission-node';

export const ADMIN_GROUP_REF = 'group:default/onemind-admins';

export const isOneMindAdmin = (user?: Pick<PolicyQueryUser, 'info'>) =>
  Boolean(user?.info.ownershipEntityRefs.includes(ADMIN_GROUP_REF));

/**
 * The one OneMind permission policy. It applies to every Backstage permission,
 * including `catalog.entity.*` and `scaffolder.task.create`:
 *
 *   anonymous / no user identity                 DENY, including reads
 *   member of group:default/onemind-admins       ALLOW everything
 *   any other signed-in user, action `read`      ALLOW
 *   any other signed-in user, create/update/delete   DENY
 *
 * Decisions are never CONDITIONAL, so there is no per-entity filtering: a user who may
 * read the catalog may read every entity. Machine callers (MCP token, plugin-to-plugin)
 * are service principals; they have no user identity here and are limited by the
 * backend `externalAccess` restrictions, not by this policy.
 */
export function createOneMindPolicy(): PermissionPolicy {
  return {
    async handle(request, user) {
      if (!user) return { result: AuthorizeResult.DENY };
      if (
        isOneMindAdmin(user) ||
        request.permission.attributes.action === 'read'
      ) {
        return { result: AuthorizeResult.ALLOW };
      }
      return { result: AuthorizeResult.DENY };
    },
  };
}
