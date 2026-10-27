# Authentik org sync: catalog Users and Groups from Authentik

Status: **built, disabled by default.** Nothing is read from Authentik or published until `authentik.orgSync.enabled` is `true`. While it is off, the module only removes entities it published earlier (see "Turning it off").
Package: `services/onemind-portal/plugins/plugin-onemind-org-backend` (catalog module `onemind-authentik-org`). Config: `platform/backstage/app-config.cluster.yaml`, schema in the package's `config.d.ts`.

## What it does

Every 10 minutes (configurable), on one portal replica, the module lists Authentik's users and groups (`GET /api/v3/core/users/` and `/api/v3/core/groups/`, every page) with a read-only token and publishes them in one full catalog mutation:

| Authentik | Catalog |
|---|---|
| Active `internal` or `external` user who is a member of a group in `authentik.orgSync.userGroups`, or of a group below one (Authentik counts members of child groups as members of the parent) | `User` named from the username. `onemind.io/oidc-sub` is the user's `uid`; `profile.displayName` is the name (or the username), `profile.email` the email; `memberOf` lists the published groups the user is directly in |
| Any other user: service accounts, deactivated users, people outside the configured groups | Nothing |
| A group in `userGroups`, or below one, unless it is an Authentik superuser group or below one | `Group` named from the group name, `type: team`, `displayName` the Authentik name, `members` the published users in it, `parent`/`children` where the group has exactly one Authentik parent (Backstage allows one; a group with several parents is published without one) |
| Any other group: outside the configured groups (for example `argocd-admins`, `fedhub-admins`, `authentik Admins`), a superuser group or a group below one (Authentik makes their members superusers, even inside the configured groups), or a group defined in `catalog/org.yaml` | Nothing, and no published user lists it in `memberOf`. A member of such a group who is also in the configured groups is still published, without that membership |

Each entity also carries `onemind.io/authentik-name` (the name before it was made a valid entity name), `onemind.io/authentik-uuid` and `backstage.io/managed-by-location: provider:onemind-authentik-org`.

## Who can read it

The OneMind permission policy allows every read, so every signed-in portal user, and the MCP machine token that agents use, can read every synced user (entity name, display name, email) and every published group with its synced members. Groups outside `userGroups` are never published, so who is in `argocd-admins`, `fedhub-admins` or `authentik Admins` stays in Authentik, and so does the membership of superuser groups placed inside the configured groups. Any other group placed below a configured group is published with its members: choose `userGroups` so that no group below them has a membership that should stay private.

Entity names are lower case, accents dropped, every run of other characters becomes one `-` (`daniel@example.com` becomes `daniel-example.com`, `Field Ops` becomes `field-ops`); names over 63 characters are cut to 54 plus `-` and 8 hex of a hash. When two records map to the same name, the older user (lower pk) or the group whose name sorts first keeps it, and the other is skipped with a warning in the log.

## Sign-in

The sign-in resolver is unchanged: it looks up the User whose `onemind.io/oidc-sub` equals the verified OIDC `sub`. Authentik issues `sub` = the user's `uid` (sha256 hashed user id) because the `onemind-backstage` provider uses `sub_mode: hashed_user_id` (`platform/ecosystem/resources/authentik-bootstrap-provider.yaml`). If that provider's subject mode is ever changed, synced users stop matching, exactly as the reviewed `akadmin` entry would.

Enabled, every published user can sign in with read access (the OneMind permission policy) if Authentik admits them to the `onemind-backstage` application. Two gates therefore decide who gets in, and they do not react at the same speed:

- **The Authentik application binding** applies at every login: Authentik refuses anyone outside the bound groups, whatever the catalog holds.
- **The catalog** follows Authentik only when a sync succeeds. Someone removed from the portal group keeps their catalog User until the next successful sync, and indefinitely while syncs fail (see "Safety rails").

An application without bindings admits every Authentik user. The bootstrap job creates `onemind-backstage` with `policy_engine_mode: any` and no bindings, and whether the live application has bindings is not recorded in Git. Binding the `userGroups` group(s) to the application, and verifying that binding, is therefore a precondition of turning the sync on (step 4 below), not a recommendation.

## What stays Git-owned

Users and Groups defined in `catalog/org.yaml` (today `akadmin` and `onemind-admins`) win: the module never publishes their names or a user holding one of their subjects, and nothing it publishes refers to them. Admin rights (membership of `group:default/onemind-admins`) therefore remain a reviewed Git change, even if Authentik has a group called "OneMind Admins". The module reads the same org file the catalog reads: the one `file` location in `catalog.locations` named `org.yaml` (`/app/catalog/org.yaml` in the image, `../../catalog/org.yaml` from `packages/backend` under `yarn start`), unless `authentik.orgSync.curatedFile` names another.

## Safety rails

Any of these fails the run, and the catalog keeps the last good state: a non-2xx answer, an unexpected response shape, a listing that changed while it was read, no active person or no group at all (Authentik answers an empty list when the token lacks a view permission), none of the `userGroups` existing, or an unreadable `catalog/org.yaml`. Requests send the token only as a bearer header to `authentik.baseUrl`, follow no redirects, have a 20 second deadline and a 10 MB response cap. The token is a `@visibility secret` config value and never appears in a log line or error.

A sync that keeps failing (an expired, revoked or rotated token, Authentik down) therefore leaves the last published users in place, including people since removed from the portal group, until a sync succeeds. The catalog's task scheduler logs each failed run at error level with the cause (for example `Authentik API answered 403 for /api/v3/core/users/ (check that authentik.token is valid and may view users and groups)`). There is no metric or alert: the portal exports no metrics. The Authentik binding above is what keeps a stale catalog from letting a removed person in.

If `authentik.token` is unset while the sync is enabled (the secret was deleted, or could not be decrypted after a restore), startup does not fail. The module logs an error naming the secret, calls nothing, and removes the users and groups it published, exactly as when it is off. Pods that started with the token keep syncing until they restart, so until every pod has restarted the published set can come and go. Missing or invalid reviewed settings (no URL, empty `userGroups`, a bad interval, no `org.yaml` file location and no `curatedFile`) still fail startup with a message naming the key, and the rolling update then keeps the old replicas serving.

## Turning it on (owner)

1. Release an image that contains this package (a new backend module needs an image; the deployment digest bump is its own reviewed change).
2. In Authentik, create a service account (for example `backstage-catalog-sync`), a role with only the global permissions **Can view User** and **Can view Group**, assigned to that account, and an API token for it. Make it non-expiring, or record its rotation date.
3. Seal the token as Secret `backstage-authentik-catalog`, key `AUTHENTIK_CATALOG_TOKEN`, namespace `backstage` (procedure in the root `CLAUDE.md`, "Secrets"), save it as `platform/backstage/authentik-catalog.sealedsecret.yaml` and add it to `platform/backstage/kustomization.yaml`. The Deployment already reads it as an optional variable. The variable is read when a pod starts: after re-sealing a rotated token, roll the Deployment through Git in the same PR, for example by adding or bumping a `kubectl.kubernetes.io/restartedAt` annotation on the pod template in `platform/backstage/deployment.yaml`.
4. Choose the Authentik group(s) whose members may use the portal and bind them to the `onemind-backstage` application (Applications, `onemind-backstage`, Policy / Group / User Bindings). **Verify the binding before going on:** sign in to the portal as an Authentik user outside those groups; Authentik must refuse access to the application. Then list the groups in `authentik.orgSync.userGroups`.
5. Set `authentik.orgSync.enabled: true`, and in the same PR update the reviewed-state assertions in `scripts/tests/test_gitops_portal_org_sync.py` and in the plugin's `module.test.ts`. The ConfigMap change rolls both replicas, which then read the token and settings. During the rollout an old replica may remove what a new one published once more; the next sync publishes it again.
6. Check the logs for `Authentik org sync enabled: members of ...` and, a minute later, `Authentik org sync: N users and M groups published; ...`, then sign in as one synced user. After any later token change, check that the second line keeps appearing every interval.

## Turning it off

Set `enabled: false` (or remove the block). The module stays registered and calls nothing: each new replica removes the synced entities when its catalog starts, and every run of the global task `onemind-authentik-org-refresh` removes them again. An old replica that runs the task once more during the rolling update can republish them only until the next run after it is gone. Runs of a global task are spaced by the task's interval, so the catalog is clear within 10 minutes (or the enabled `frequencySeconds`, if it was set higher) of the last old replica's last run. `catalog/org.yaml` is untouched.

To confirm, list `/api/catalog/entities?filter=metadata.annotations.backstage.io/managed-by-location=provider:onemind-authentik-org` once the rollout has finished and that interval has passed: it must return nothing.

Removing the module from `packages/backend/src/index.ts` is different: the catalog then deletes its entities once at each pod start (`catalog.orphanProviderStrategy`, default `delete`), and nothing removes what an old replica republishes afterwards. Turn the sync off first, and remove the module only after the check above returns nothing.
