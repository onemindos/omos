import type { AuthService, DiscoveryService } from '@backstage/backend-plugin-api';

const GROUP = 'onemind-admins';
export const ADMIN_GROUP_REF = `group:default/${GROUP}`;

type Entity = {
  spec?: { members?: unknown };
  relations?: { type: string; targetRef: string }[];
};

// Resolves the intended recipients from the catalog: the onemind-admins group and its
// explicit members. Anything unexpected means "not ready", so nothing is sent blindly.
export function createAdminRecipients(deps: {
  auth: AuthService;
  discovery: DiscoveryService;
}) {
  const entity = async (kind: string, name: string): Promise<Entity> => {
    const { token } = await deps.auth.getPluginRequestToken({
      onBehalfOf: await deps.auth.getOwnServiceCredentials(),
      targetPluginId: 'catalog',
    });
    const base = await deps.discovery.getBaseUrl('catalog');
    const response = await fetch(
      `${base}/entities/by-name/${kind}/default/${encodeURIComponent(name)}`,
      {
        headers: { Authorization: `Bearer ${token}` },
        signal: AbortSignal.timeout(8000),
        redirect: 'error',
      },
    );
    if (!response.ok) throw new Error('Recipient catalog unavailable');
    return response.json();
  };

  // Verified user entity refs of the admin group, or undefined when not ready.
  const userRefs = async (): Promise<string[] | undefined> => {
    const group = await entity('group', GROUP);
    const members = group.spec?.members;
    if (!Array.isArray(members) || members.length < 1 || members.length > 10) {
      return undefined;
    }
    const refs: string[] = [];
    for (const member of members) {
      if (typeof member !== 'string' || !/^[a-z0-9_-]+$/.test(member)) {
        return undefined;
      }
      const user = await entity('user', member);
      if (
        !user.relations?.some(
          relation =>
            relation.type === 'memberOf' && relation.targetRef === ADMIN_GROUP_REF,
        )
      ) {
        return undefined;
      }
      refs.push(`user:default/${member}`);
    }
    return refs;
  };

  return { userRefs, ready: async () => Boolean(await userRefs()) };
}
