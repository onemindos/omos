# Authentik org sync: catalog Users and Groups from Authentik

Status: **built, disabled by default.** Nothing is read from Authentik or published until `authentik.orgSync.enabled` is `true`.
Package: `services/onemind-portal/plugins/plugin-onemind-org-backend` (catalog module `onemind-authentik-org`). Config: `platform/backstage/app-config.cluster.yaml`, schema in the package's `config.d.ts`.

## What it does

Every 10 minutes (configurable), on one portal replica, the module lists Authentik's users and groups (`GET /api/v3/core/users/` and `/api/v3/core/groups/`, every page) with a read-only token and publishes them in one full catalog mutation:

| Authentik | Catalog |
|---|---|
| Active `internal` or `external` user who is a member of a group in `authentik.orgSync.userGroups`, or of a group below one (Authentik counts members of child groups as members of the parent) | `User` named from the username. `onemind.io/oidc-sub` is the user's `uid`; `profile.displayName` is the name (or the username), `profile.email` the email; `memberOf` lists the published groups the user is directly in |
| Any other user: service accounts, deactivated users, people outside the configured groups | Nothing |
| Every group | `Group` named from the group name, `type: team`, `displayName` the Authentik name, `members` the published users in it, `parent`/`children` where the group has exactly one Authentik parent (Backstage allows one; a group with several parents is published without one) |

Each entity also carries `onemind.io/authentik-name` (the name before it was made a valid entity name), `onemind.io/authentik-uuid` and `backstage.io/managed-by-location: provider:onemind-authentik-org`.

Entity names are lower case, accents dropped, every run of other characters becomes one `-` (`daniel@example.com` becomes `daniel-example.com`, `Field Ops` becomes `field-ops`); names over 63 characters are cut to 54 plus `-` and 8 hex of a hash. When two records map to the same name, the older user (lower pk) or the group whose name sorts first keeps it, and the other is skipped with a warning in the log.

## Sign-in

The sign-in resolver is unchanged: it looks up the User whose `onemind.io/oidc-sub` equals the verified OIDC `sub`. Authentik issues `sub` = the user's `uid` (sha256 hashed user id) because the `onemind-backstage` provider uses `sub_mode: hashed_user_id` (`platform/ecosystem/resources/authentik-bootstrap-provider.yaml`). If that provider's subject mode is ever changed, synced users stop matching, exactly as the reviewed `akadmin` entry would.

Enabled, every published user can sign in with read access (the OneMind permission policy) if Authentik admits them to the `onemind-backstage` application. That is why `userGroups` is required: bind the same group to the application in Authentik, so the catalog and Authentik agree on who may use the portal.

## What stays Git-owned

Users and Groups defined in `catalog/org.yaml` (today `akadmin` and `onemind-admins`) win: the module never publishes their names or a user holding one of their subjects, and nothing it publishes refers to them. Admin rights (membership of `group:default/onemind-admins`) therefore remain a reviewed Git change, even if Authentik has a group called "OneMind Admins".

## Safety rails

Any of these fails the run, and the catalog keeps the last good state: a non-2xx answer, an unexpected response shape, a listing that changed while it was read, no active person or no group at all (Authentik answers an empty list when the token lacks a view permission), none of the `userGroups` existing, or an unreadable `catalog/org.yaml`. Requests send the token only as a bearer header to `authentik.baseUrl`, follow no redirects, have a 20 second deadline and a 10 MB response cap. The token is a `@visibility secret` config value and never appears in a log line or error.

## Turning it on (owner)

1. Release an image that contains this package (a new backend module needs an image; the deployment digest bump is its own reviewed change).
2. In Authentik, create a service account (for example `backstage-catalog-sync`), a role with only the global permissions **Can view User** and **Can view Group**, assigned to that account, and an API token for it.
3. Seal the token as Secret `backstage-authentik-catalog`, key `AUTHENTIK_CATALOG_TOKEN`, namespace `backstage` (procedure in the root `CLAUDE.md`, "Secrets"), save it as `platform/backstage/authentik-catalog.sealedsecret.yaml` and add it to `platform/backstage/kustomization.yaml`. The Deployment already reads it as an optional variable.
4. Choose the Authentik group(s) whose members may use the portal, bind them to the `onemind-backstage` application, and list them in `authentik.orgSync.userGroups`.
5. Set `authentik.orgSync.enabled: true`. The ConfigMap change rolls both replicas, which then read the token and settings.
6. Check the logs for `Authentik org sync enabled: members of ...` and, a minute later, `Authentik org sync: N users and M groups published; ...`, then sign in as one synced user.

Missing settings while enabled (no token, no URL, empty `userGroups`) fail startup with a message naming the key; the rolling update then keeps the old replicas serving.

## Turning it off

Set `enabled: false` (or remove the block). On restart the catalog removes every entity of a provider that is no longer installed (`catalog.orphanProviderStrategy` defaults to `delete`), so synced users can no longer sign in; `catalog/org.yaml` is untouched. If an old replica happened to run one more sync during the rollout, the next restart removes what it published.
