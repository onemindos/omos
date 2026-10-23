import { AuthorizeResult } from '@backstage/plugin-permission-common';
import type {
  PermissionPolicy,
  PolicyQueryUser,
} from '@backstage/plugin-permission-node';

export const isOneMindAdmin = (user?: Pick<PolicyQueryUser, 'info'>) =>
  Boolean(
    user?.info.ownershipEntityRefs.includes('group:default/onemind-admins'),
  );

// Administrators (group onemind-admins) are allowed; other signed-in users may only read.
// An optional delegate can narrow (never widen) the result for OpenChoreo-owned permissions.
export function createOneMindPolicy(
  delegate?: PermissionPolicy,
): PermissionPolicy {
  return {
    async handle(request, user) {
      if (
        !user ||
        (!isOneMindAdmin(user) &&
          request.permission.attributes.action !== 'read')
      ) {
        return { result: AuthorizeResult.DENY };
      }
      if (
        delegate &&
        (request.permission.name.startsWith('openchoreo.') ||
          request.permission.name.startsWith('catalog.entity.') ||
          request.permission.name === 'scaffolder.task.create')
      ) {
        try {
          return await delegate.handle(request, user);
        } catch {
          return { result: AuthorizeResult.DENY };
        }
      }
      return { result: AuthorizeResult.ALLOW };
    },
  };
}
