# OneMind role modes (FEATURES J2)

`onemind-roles.json` is the single source for the five roles. Not applied by any Application
(`platform-security` does not recurse into this directory).

| Field | Meaning |
|---|---|
| `client_account` | Account the auth callout places users in (`CLIENT`, A19) |
| `name_prefix` | Prefix of callout-issued NATS user names; the name is `ak:<username>;<group>,...` (services/auth-callout `PrincipalName`, read by services/policy `resolveCallout`) |
| `group_tag_prefix`, `role_tag_prefix` | Tags written into the user JWT (`group:`, `role:`). nats-server ignores them in config mode; they are kept for operator mode |
| `inbox_prefix` | Clients use `<inbox_prefix>.<sub>` as their inbox prefix; the user JWT allows only that |
| `max_ttl_seconds` | Upper bound on a callout user's lifetime (also bounded by the token's `exp`) |
| `connection_types` | WEBSOCKET and STANDARD only (never LEAFNODE) |
| `roles[].group` | The Authentik group granting the role (created by `platform/auth-callout/authentik-blueprint.yaml`) |
| `roles[].rank` | Higher wins when a user is in several groups; every role must keep all rights of the role below it (`TestRoleMapSourceRules`) |
| `roles[].ui_modes` | Client modes to show (explore, watch, operate, command, admin). Display only; the bus and the policy service enforce |
| `roles[].publish`, `roles[].subscribe` | NATS allow lists for the role. The parser refuses anything outside `svc.{service}.*`, `agents.prompt.hermes.>` and `cmd.map.>` (no data subjects, no `$JS`, no `_INBOX` wildcard) |
| `roles[].openfga` | `inherits` (nested group, lower role), `operate_sites`, and the highest action level the role may propose, approve and second-approve. `platform/security/openfga/onemind-roles.tuples.json` is generated from it |

After editing: copy the file into `platform/auth-callout/configmap.yaml` (`roles.json`), and
regenerate the tuples:

```bash
cd services/policy && go test ./internal/policy -run TestRoleTuplesGolden -update-role-tuples
cd ../auth-callout && go test ./...      # TestConfigMapMatchesSource, TestRoleMapSourceRules
```
