/* eslint-disable no-restricted-syntax */
/* Test-only helpers: Authentik API records shaped like the 2026.8 API (never imported by index.ts). */
import { createHash } from 'node:crypto';
import { resolve } from 'node:path';
import type { FetchJson } from './authentik';

/** A stable stand-in for Authentik's hashed user id (really sha256("<user id>-<install id>")). */
export const uidOf = (pk: number) =>
  createHash('sha256').update(`fixture-user-${pk}`).digest('hex');
export const groupPk = (n: number) =>
  `00000000-0000-4000-8000-${String(n).padStart(12, '0')}`;

export const portalPath = (relative: string) =>
  resolve(__dirname, '../../..', relative);
export const orgYamlPath = portalPath('catalog/org.yaml');

/** One GET /api/v3/core/users/ record with every field the API returns (include_groups=false). */
export function apiUser(
  pk: number,
  username: string,
  groups: number[],
  extra: Record<string, unknown> = {},
) {
  return {
    pk,
    username,
    name: `Person ${pk}`,
    is_active: true,
    last_login: null,
    date_joined: '2026-01-01T00:00:00Z',
    is_superuser: false,
    groups: groups.map(groupPk),
    groups_obj: null,
    roles: [],
    roles_obj: null,
    email: `person${pk}@example.test`,
    avatar: 'data:image/svg+xml;base64,PHN2Zy8+',
    attributes: {},
    uid: uidOf(pk),
    path: 'users',
    type: 'internal',
    uuid: `10000000-0000-4000-8000-${String(pk).padStart(12, '0')}`,
    password_change_date: '2026-01-01T00:00:00Z',
    last_updated: '2026-01-01T00:00:00Z',
    ...extra,
  };
}

/** One GET /api/v3/core/groups/ record with every field the API returns (include_users=false). */
export function apiGroup(n: number, name: string, parents: number[] = []) {
  return {
    pk: groupPk(n),
    num_pk: 10000 + n,
    name,
    is_superuser: false,
    parents: parents.map(groupPk),
    parents_obj: null,
    users: [],
    users_obj: null,
    attributes: {},
    roles: [],
    roles_obj: [],
    inherited_roles_obj: null,
    children: [],
    children_obj: null,
  };
}

/**
 * Serves users and groups the way Authentik pages them: `page`/`page_size` honoured (or a smaller
 * server page size), `pagination.next` 0 on the last page. Server-side filters are NOT applied, so
 * the mapping's own filtering is what the tests see.
 */
export function fakeAuthentik(
  data: { users: unknown[]; groups: unknown[] },
  serverPageSize?: number,
) {
  const urls: string[] = [];
  const fetchJson: FetchJson = async url => {
    urls.push(url);
    const u = new URL(url);
    let items: unknown[];
    if (u.pathname === '/api/v3/core/users/') items = data.users;
    else if (u.pathname === '/api/v3/core/groups/') items = data.groups;
    else throw new Error(`unexpected path ${u.pathname}`);
    const size = Math.min(
      Number(u.searchParams.get('page_size') ?? 20),
      serverPageSize ?? 100,
    );
    const page = Number(u.searchParams.get('page') ?? 1);
    const totalPages = Math.max(1, Math.ceil(items.length / size));
    const start = (page - 1) * size;
    const results = items.slice(start, start + size);
    return {
      pagination: {
        next: page < totalPages ? page + 1 : 0,
        previous: page > 1 ? page - 1 : 0,
        count: items.length,
        current: page,
        total_pages: totalPages,
        start_index: results.length ? start + 1 : 0,
        end_index: start + results.length,
      },
      results,
      autocomplete: {},
    };
  };
  return { fetchJson, urls };
}

/**
 * A small directory: a portal scope group with a child and grandchild, a group with two parents,
 * groups outside the scope, a group whose name collides with the reviewed onemind-admins, and users
 * that must or must not be published.
 */
export const directory = {
  groups: [
    apiGroup(1, 'OneMind Portal Users'),
    apiGroup(2, 'Field Ops', [1]),
    apiGroup(3, 'Field Ops / Night Shift', [2]),
    apiGroup(4, 'authentik Admins'),
    apiGroup(5, 'OneMind Admins'),
    apiGroup(6, 'Two Parents', [1, 4]),
    apiGroup(7, 'TAK Users'),
  ],
  users: [
    apiUser(1, 'akadmin', [4, 1, 5], { is_superuser: true }),
    apiUser(2, 'daniel@example.test', [1, 5], { name: 'Daniel Example' }),
    apiUser(3, 'night.owl', [3]),
    apiUser(4, 'ak-outpost-ldap', [1], { type: 'internal_service_account' }),
    apiUser(5, 'former', [1], { is_active: false }),
    apiUser(6, 'tak-only', [7]),
    apiUser(7, 'Daniel@Example.test', [1]),
    apiUser(8, 'two-parents-member', [6], { email: '' }),
    apiUser(9, 'external.partner', [2], { type: 'external', name: '' }),
  ],
};
export const SCOPE = ['OneMind Portal Users'];
