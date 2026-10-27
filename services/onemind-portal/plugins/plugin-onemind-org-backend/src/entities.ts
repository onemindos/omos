import { createHash } from 'node:crypto';
import { readFile } from 'node:fs/promises';
import type { Entity } from '@backstage/catalog-model';
import YAML from 'yaml';
import {
  PERSON_TYPES,
  type AuthentikGroup,
  type AuthentikSnapshot,
  type AuthentikUser,
} from './authentik';

export const PROVIDER_NAME = 'onemind-authentik-org';
export const LOCATION = `provider:${PROVIDER_NAME}`;
export const A = {
  /**
   * The annotation the OneMind sign-in resolver (plugin-onemind-auth-backend, subject.ts) looks up
   * with the verified OIDC `sub`, compared as an exact value.
   */
  oidcSub: 'onemind.io/oidc-sub',
  /** The Authentik username or group name, before it was made a valid entity name. */
  name: 'onemind.io/authentik-name',
  uuid: 'onemind.io/authentik-uuid',
} as const;

const MAX_NAME = 63;
const digest = (raw: string) =>
  createHash('sha256').update(raw).digest('hex').slice(0, 8);

/**
 * A valid, stable catalog entity name for an Authentik username or group name: lower case, accents
 * dropped, every run of other characters turned into one "-", starting and ending alphanumeric.
 * Names longer than 63 characters are cut to 54 plus "-" and 8 hex of sha256(raw); a name with no
 * usable character at all becomes "authentik-" plus that hash.
 */
export function entityName(raw: string): string {
  const cleaned = raw
    .normalize('NFKD')
    .replace(/[̀-ͯ]/g, '')
    .toLowerCase()
    .replace(/[^a-z0-9._-]+/g, '-')
    .replace(/-{2,}/g, '-')
    .replace(/^[^a-z0-9]+|[^a-z0-9]+$/g, '');
  if (!cleaned) return `authentik-${digest(raw)}`;
  if (cleaned.length <= MAX_NAME) return cleaned;
  return `${cleaned.slice(0, 54).replace(/[^a-z0-9]+$/, '')}-${digest(raw)}`;
}

/** What the reviewed org file (catalog/org.yaml) defines; those entities stay Git-owned. */
export interface Curated {
  /** Lower-case User names in the default namespace. */
  users: Set<string>;
  /** Lower-case Group names in the default namespace. */
  groups: Set<string>;
  /** Lower-case `onemind.io/oidc-sub` values of every User, in any namespace. */
  subjects: Set<string>;
}

/** Reads the reviewed org file. An unreadable or invalid file throws, so nothing is published blind. */
export async function readCurated(file: string): Promise<Curated> {
  const text = await readFile(file, 'utf8');
  const curated: Curated = {
    users: new Set(),
    groups: new Set(),
    subjects: new Set(),
  };
  for (const doc of YAML.parseAllDocuments(text)) {
    if (doc.errors.length) throw new Error(`${file} is not valid YAML`);
    const e = doc.toJSON() as Partial<Entity> | null;
    if (!e || typeof e !== 'object') continue;
    const name = String(e.metadata?.name ?? '').toLowerCase();
    const inDefault = (e.metadata?.namespace ?? 'default') === 'default';
    if (e.kind === 'User') {
      const sub = e.metadata?.annotations?.[A.oidcSub];
      if (typeof sub === 'string' && sub)
        curated.subjects.add(sub.toLowerCase());
      if (inDefault && name) curated.users.add(name);
    } else if (e.kind === 'Group' && inDefault && name) {
      curated.groups.add(name);
    }
  }
  return curated;
}

export interface MappingOptions {
  /** Exact Authentik group names whose members (including members of child groups) are published. */
  userGroups: string[];
  curated: Curated;
}

export interface OrgMapping {
  /** Users first, then groups, each sorted by name. */
  entities: Entity[];
  users: number;
  groups: number;
  /** Authentik users and groups in the scope left to the reviewed org file. */
  curatedUsers: number;
  curatedGroups: number;
  /** Groups in the scope never published because Authentik makes their members superusers. */
  superuserGroups: number;
  /** Groups published without a parent because Authentik gives them several. */
  multiParentGroups: number;
  /** Problems an operator should fix in Authentik or in the configuration. */
  warnings: string[];
}

const byCodePoint = (a: string, b: string) => {
  if (a === b) return 0;
  return a < b ? -1 : 1;
};
/** Appends to the list under `key`, creating it on first use. */
const push = (map: Map<string, string[]>, key: string, value: string) => {
  const list = map.get(key);
  if (list) list.push(value);
  else map.set(key, [value]);
};
const sortedUnique = (values: string[]) =>
  [...new Set(values)].sort(byCodePoint);

function userEntity(
  user: AuthentikUser,
  name: string,
  memberOf: string[],
): Entity {
  const profile: Record<string, string> = {
    displayName: user.name.trim() || user.username,
  };
  if (user.email.trim()) profile.email = user.email.trim();
  const annotations: Record<string, string> = {
    'backstage.io/managed-by-location': LOCATION,
    'backstage.io/managed-by-origin-location': LOCATION,
    [A.oidcSub]: user.uid,
    [A.name]: user.username,
  };
  if (user.uuid) annotations[A.uuid] = user.uuid;
  return {
    apiVersion: 'backstage.io/v1alpha1',
    kind: 'User',
    metadata: { name, annotations },
    spec: { profile, memberOf },
  };
}

function groupEntity(
  group: AuthentikGroup,
  name: string,
  parent: string | undefined,
  children: string[],
  members: string[],
): Entity {
  return {
    apiVersion: 'backstage.io/v1alpha1',
    kind: 'Group',
    metadata: {
      name,
      annotations: {
        'backstage.io/managed-by-location': LOCATION,
        'backstage.io/managed-by-origin-location': LOCATION,
        [A.name]: group.name,
        [A.uuid]: group.pk,
      },
    },
    spec: {
      type: 'team',
      profile: { displayName: group.name },
      ...(parent ? { parent } : {}),
      children,
      members,
    },
  };
}

/**
 * Maps one Authentik snapshot to catalog Users and Groups.
 *
 * - Scope: the configured groups and every group below them (Authentik's parent links, all of them).
 * - Users: active people (internal or external) who are members of a group in the scope, as
 *   Authentik counts membership. `onemind.io/oidc-sub` is their `uid`.
 * - Groups: only groups in the scope, so the membership of other groups (argocd-admins, for
 *   example) is never published or listed in a user's memberOf. Within the scope, a superuser group
 *   or any group below one is not published either, because Authentik makes its members
 *   superusers. The hierarchy is copied only where Backstage can hold it (one parent per group).
 * - Names and subjects defined in the reviewed org file are never published, and nothing published
 *   points at them: admin membership (group:default/onemind-admins) stays Git-reviewed.
 * - Two records that map to the same entity name: the older user (lower pk) or the group whose name
 *   sorts first keeps it; the other is skipped with a warning.
 */
export function buildOrgEntities(
  snapshot: AuthentikSnapshot,
  { userGroups, curated }: MappingOptions,
): OrgMapping {
  const warnings: string[] = [];

  const childPks = new Map<string, string[]>();
  for (const group of snapshot.groups)
    for (const parent of group.parents) push(childPks, parent, group.pk);
  /** The given groups and every group below them; a parent cycle ends the walk. */
  const andBelow = (pks: string[]) => {
    const found = new Set<string>();
    const queue = [...pks];
    while (queue.length) {
      const pk = queue.pop()!;
      if (found.has(pk)) continue;
      found.add(pk);
      queue.push(...(childPks.get(pk) ?? []));
    }
    return found;
  };

  const wanted = new Set(userGroups);
  const roots = snapshot.groups.filter(g => wanted.has(g.name));
  const found = new Set(roots.map(g => g.name));
  const missing = userGroups.filter(n => !found.has(n));
  if (!roots.length)
    throw new Error(
      `None of the groups in authentik.orgSync.userGroups (${userGroups.join(
        ', ',
      )}) exists in Authentik; refusing to publish`,
    );
  if (missing.length)
    warnings.push(
      `authentik.orgSync.userGroups names groups that do not exist in Authentik: ${missing.join(
        ', ',
      )}`,
    );
  const scope = andBelow(roots.map(g => g.pk));
  // Authentik counts a member of any group below a superuser group as a superuser too.
  const superuser = andBelow(
    snapshot.groups.filter(g => g.isSuperuser).map(g => g.pk),
  );

  const groupNames = new Map<string, string>(); // Authentik pk -> published entity name
  const groupHolders = new Map<string, string>(); // entity name -> Authentik group name
  let curatedGroups = 0;
  let superuserGroups = 0;
  const ordered = snapshot.groups
    .filter(g => scope.has(g.pk))
    .sort((a, b) => byCodePoint(a.name, b.name));
  for (const group of ordered) {
    const name = entityName(group.name);
    if (curated.groups.has(name)) {
      curatedGroups++;
      continue;
    }
    if (superuser.has(group.pk)) {
      superuserGroups++;
      continue;
    }
    const holder = groupHolders.get(name);
    if (holder !== undefined) {
      warnings.push(
        `group "${group.name}" maps to the entity name ${name}, already used by group "${holder}"; skipped`,
      );
      continue;
    }
    groupHolders.set(name, group.name);
    groupNames.set(group.pk, name);
  }

  const users: { user: AuthentikUser; name: string }[] = [];
  const userHolders = new Map<string, string>(); // entity name -> Authentik username
  const subjects = new Set<string>();
  let curatedUsers = 0;
  for (const user of [...snapshot.users].sort((a, b) => a.pk - b.pk)) {
    if (!user.isActive || !PERSON_TYPES.includes(user.type)) continue;
    if (!user.groups.some(pk => scope.has(pk))) continue;
    const name = entityName(user.username);
    const subject = user.uid.toLowerCase();
    if (curated.users.has(name) || curated.subjects.has(subject)) {
      curatedUsers++;
      continue;
    }
    const holder = userHolders.get(name);
    if (holder !== undefined) {
      warnings.push(
        `user "${user.username}" maps to the entity name ${name}, already used by user "${holder}"; skipped`,
      );
      continue;
    }
    if (subjects.has(subject)) {
      warnings.push(
        `user "${user.username}" has the same subject as another user; skipped`,
      );
      continue;
    }
    userHolders.set(name, user.username);
    subjects.add(subject);
    users.push({ user, name });
  }

  const members = new Map<string, string[]>(); // group entity name -> user entity names
  const userEntities = users.map(({ user, name }) => {
    const memberOf = sortedUnique(
      user.groups.flatMap(pk => groupNames.get(pk) ?? []),
    );
    for (const group of memberOf) push(members, group, name);
    return userEntity(user, name, memberOf);
  });

  // Backstage allows one parent per group: copy a link only when Authentik has exactly one parent.
  const parents = new Map<string, string>(); // child entity name -> parent entity name
  let multiParentGroups = 0;
  for (const group of ordered) {
    const child = groupNames.get(group.pk);
    if (!child) continue;
    if (group.parents.length > 1) {
      multiParentGroups++;
      continue;
    }
    const parent =
      group.parents.length === 1 ? groupNames.get(group.parents[0]) : undefined;
    if (parent) parents.set(child, parent);
  }
  // Authentik should not allow a parent cycle; never publish one, since the org pages walk the tree.
  for (const child of [...parents.keys()]) {
    let next = parents.get(child);
    for (let step = 0; next !== undefined && step <= parents.size; step++) {
      if (next === child) {
        parents.delete(child);
        warnings.push(
          `group ${child} is in a parent cycle in Authentik; published without a parent`,
        );
        break;
      }
      next = parents.get(next);
    }
  }
  const children = new Map<string, string[]>();
  for (const [child, parent] of parents) push(children, parent, child);

  const groupEntities = ordered
    .filter(g => groupNames.get(g.pk) !== undefined)
    .map(group => {
      const name = groupNames.get(group.pk)!;
      return groupEntity(
        group,
        name,
        parents.get(name),
        sortedUnique(children.get(name) ?? []),
        sortedUnique(members.get(name) ?? []),
      );
    });

  const byName = (a: Entity, b: Entity) =>
    byCodePoint(a.metadata.name, b.metadata.name);
  return {
    entities: [...userEntities.sort(byName), ...groupEntities.sort(byName)],
    users: userEntities.length,
    groups: groupEntities.length,
    curatedUsers,
    curatedGroups,
    superuserGroups,
    multiParentGroups,
    warnings,
  };
}
