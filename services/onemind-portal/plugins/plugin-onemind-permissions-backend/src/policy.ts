import { AuthorizeResult } from '@backstage/plugin-permission-common';
import type {
  PermissionPolicy,
  PolicyQueryUser,
} from '@backstage/plugin-permission-node';

export const ADMIN_GROUP_REF = 'group:default/onemind-admins';

export const isOneMindAdmin = (user?: Pick<PolicyQueryUser, 'info'>) =>
  Boolean(user?.info.ownershipEntityRefs.includes(ADMIN_GROUP_REF));

/**
 * Permissions whose action is `read` but whose answer exposes operator internals
 * (host details, package inventory, ingestion failures). Without this list the rule
 * "any signed-in user may read" would let every user open the admin tools. Only
 * administrators are allowed these.
 */
export const ADMIN_ONLY_READ_PERMISSIONS: readonly string[] = [
  'devtools.administer',
  'devtools.info',
  'devtools.external-dependencies',
  // Failed and pending catalog entities, with their raw source documents and errors.
  'catalog.entities.unprocessed.read',
];

/**
 * Permissions nobody gets, administrators included.
 *  - devtools.config: the DevTools config view returns the loaded configuration, and it
 *    hides only values that a plugin's config schema marks `@visibility secret`. A secret
 *    in a key without that mark would be shown, so the endpoint stays closed.
 *  - devtools.scheduler.*: the scheduled-task view is not part of this portal, and its
 *    trigger would let a click run background jobs.
 */
export const NEVER_ALLOWED_PERMISSIONS: readonly string[] = [
  'devtools.config',
  'devtools.scheduler.read',
  'devtools.scheduler.trigger',
];

/**
 * The one OneMind permission policy. It applies to every Backstage permission,
 * including `catalog.entity.*` and `scaffolder.task.create`:
 *
 *   anonymous / no user identity                 DENY, including reads
 *   NEVER_ALLOWED_PERMISSIONS                    DENY for everyone, administrators included
 *   member of group:default/onemind-admins       ALLOW everything else
 *   any other signed-in user, action `read`      ALLOW, except ADMIN_ONLY_READ_PERMISSIONS
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
      const { name, attributes } = request.permission;
      if (NEVER_ALLOWED_PERMISSIONS.includes(name)) {
        return { result: AuthorizeResult.DENY };
      }
      if (isOneMindAdmin(user)) return { result: AuthorizeResult.ALLOW };
      if (
        attributes.action === 'read' &&
        !ADMIN_ONLY_READ_PERMISSIONS.includes(name)
      ) {
        return { result: AuthorizeResult.ALLOW };
      }
      return { result: AuthorizeResult.DENY };
    },
  };
}
