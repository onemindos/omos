### Configure Cross-Account Message Imports and Exports

Source: https://github.com/nats-io/nats.docs/blob/master/running-a-nats-service/nats_admin/jwt.md

Configuration snippet demonstrating how to allow message flow between accounts using imports and exports.

```text
accounts: {
    A: {
        users: [{user: a, password: a}]
        imports: [{stream: {account: B, subject: "foo"}}]
    },
    B: {
        users: [{user: b, password: b}]
        exports: [{stream: "foo"}]
    },
}
```

--------------------------------

### Configure NATS Cross-Account JetStream Exports and Imports

Source: https://github.com/nats-io/nats.docs/blob/master/running-a-nats-service/configuration/leafnodes/jetstream_leafnodes.md

This configuration demonstrates how to define exports in a source account to allow other accounts to access JetStream consumers and streams. It also shows how to define imports in target accounts, including subject remapping for API calls.

```text
accounts {
    SYS: {
        users: [{user: admin, password: admin}]
    },
    ACC: {
        users: [{user: acc, password: acc}],
        jetstream: enabled
        exports: [
            {service: "$JS.hub.API.CONSUMER.CREATE.*", response_type: "stream"}
            {stream: "deliver.acc.hub.>"}
            {service: "$JS.hub.API.CONSUMER.MSG.NEXT.aggregate-test-leaf.dur", response_type: "stream"}
            {service: "$JS.ACK.aggregate-test-leaf.dur.>"}
            {service: "$JS.FC.aggregate-test-leaf.dur.>"}
        ]
    }
    IMPORT_MIRROR: {
        users: [{user: import_mirror, password: import_mirror}],
        jetstream: enabled
        imports: [
            {service: {account: ACC, subject: "$JS.hub.API.CONSUMER.CREATE.*"}, to: "JS.acc@hub.API.CONSUMER.CREATE.*" }
            {service: {account: ACC, subject: "$JS.FC.aggregate-test-leaf.dur.>"}}
            {stream: {account: ACC, subject: deliver.acc.hub.import_mirror.>}}
        ]
    }
    IMPORT_CLIENT: {
        users: [{user: import_client, password: import_client}],
        jetstream: enabled
        imports: [
            {service: {account: ACC, subject: "$JS.hub.API.CONSUMER.MSG.NEXT.aggregate-test-leaf.dur"}, to: "JS.acc@hub.API.CONSUMER.MSG.NEXT.aggregate-test-leaf.dur" }
            {service: {account: ACC, subject: "$JS.ACK.aggregate-test-leaf.dur.>"}}
        ]
    }
}
system_account: SYS
```

--------------------------------

### Add a Public Stream Export with NSC

Source: https://github.com/nats-io/nats.docs/blob/master/using-nats/nats-tools/nsc/streams.md

This command adds a new stream export to the current account. It uses a subject pattern with a wildcard to ensure all matching messages are shared with importing accounts.

```shell
nsc add export --name abc --subject "a.b.c.>"
```

--------------------------------

### Create Private Service Export with NSC

Source: https://github.com/nats-io/nats.docs/blob/master/using-nats/nats-tools/nsc/services.md

This command adds a private service export to an account, allowing for controlled access to services. The `--private` flag ensures that only authorized accounts with a valid JWT can import the service. Wildcards in the subject allow for flexible mapping of specific subjects to authorized accounts.

```shell
nsc add export --subject "private.help.*" --private --service --account A
```

```text
[ OK ] added private service export "private.help.*"
```

--------------------------------

### Remap Imported Subjects using NSC

Source: https://github.com/nats-io/nats.docs/blob/master/running-a-nats-service/nats_admin/jwt.md

This command demonstrates how to import a stream from a source account while remapping the subject name from the exporter's original subject to a local subject name.

```shell
nsc add import --account test --src-account ACJ6G45BE7LLOFCVAZSZR3RY4XELXQ32BOQRI7KQMQLICXXXJRP4P45Q --remote-subject foo --local-subject bar
```

--------------------------------

### Describe NATS Account Details

Source: https://github.com/nats-io/nats.docs/blob/master/using-nats/nats-tools/nsc/services.md

Displays the configuration and import/export status of a specific NATS account. Useful for verifying that private imports were successfully registered.

```shell
nsc describe account B
```

--------------------------------

### Describe Account Details and Exports with NSC

Source: https://github.com/nats-io/nats.docs/blob/master/using-nats/nats-tools/nsc/services.md

This command displays detailed information about a specific NATS account, including its configuration, limits, and existing exports. It is useful for verifying the setup of private services and understanding which subjects are exported and their accessibility.

```bash
nsc describe account A
```

```text
╭──────────────────────────────────────────────────────────────────────────────────────╮
│                                   Account Details                                    │
├───────────────────────────┬──────────────────────────────────────────────────────────┤
│ Name                      │ A                                                        │
│ Account ID                │ ADETPT36WBIBUKM3IBCVM4A5YUSDXFEJPW4M6GGVBYCBW7RRNFTV5NGE │
│ Issuer ID                 │ OAFEEYZSYYVI4FXLRXJTMM32PQEI3RGOWZJT7Y3YFM4HB7ACPE4RTJPG │
│ Issued                    │ 2019-12-04 20:19:19 UTC                                  │
│ Expires                   │                                                          │
├───────────────────────────┼──────────────────────────────────────────────────────────┤
│ Max Connections           │ Unlimited                                                │
│ Max Leaf Node Connections │ Unlimited                                                │
│ Max Data                  │ Unlimited                                                │
│ Max Exports               │ Unlimited                                                │
│ Max Imports               │ Unlimited                                                │
│ Max Msg Payload           │ Unlimited                                                │
│ Max Subscriptions         │ Unlimited                                                │
│ Exports Allows Wildcards  │ True                                                     │
├───────────────────────────┼──────────────────────────────────────────────────────────┤
│ Imports                   │ None                                                     │
╰───────────────────────────┴──────────────────────────────────────────────────────────╯

╭─────────────────────────────────────────────────────────────────────────────╮
│                                   Exports                                   │
├────────────────┬─────────┬────────────────┬────────┬─────────────┬──────────┤
│ Name           │ Type    │ Subject        │ Public │ Revocations │ Tracking │
├────────────────┼─────────┼────────────────┼────────┼─────────────┼──────────┤
│ help           │ Service │ help           │ Yes    │ 0           │ -        │
│ private.help.* │ Service │ private.help.* │ No     │ 0           │ -        │
╰────────────────┴─────────┴────────────────┴────────┴─────────────┴──────────╯
```

--------------------------------

### NATS Export Configuration Example

Source: https://github.com/nats-io/nats.docs/blob/master/running-a-nats-service/configuration/securing_nats/accounts.md

Demonstrates the structure of an export configuration map in NATS. It shows how to define streams and services, and optionally restrict access to specific accounts. The example includes public exports and account-specific exports.

```text
accounts: {
    A: {
        users: [
            {user: a, password: a}
        ]
        exports: [
            {stream: puba.>}
            {service: pubq.>}
            {stream: b.>, accounts: [B]}
            {service: q.b, accounts: [B]}
        ]
    }
    ...
}
```

--------------------------------

### Import Private Stream using nsc

Source: https://github.com/nats-io/nats.docs/blob/master/using-nats/nats-tools/nsc/streams.md

This command imports a private stream into account B using a provided activation token. The token contains all necessary details for the import. The output confirms the successful addition of the stream import.

```shell
nsc add import --account B --token /tmp/activation.jwt
```

--------------------------------

### Add Public Service Export using NSC CLI

Source: https://github.com/nats-io/nats.docs/blob/master/using-nats/nats-tools/nsc/services.md

This command adds a new public service export named 'help' with the subject 'help'. Services are associated with the account performing replies and are advertised in the exporting account's JWT. This functionality is crucial for enabling inter-account communication via request-reply patterns.

```bash
nsc add export --name help --subject help --service
```

--------------------------------

### Import Private Service with NSC

Source: https://github.com/nats-io/nats.docs/blob/master/using-nats/nats-tools/nsc/services.md

Adds a private service import to an account using an activation token file. This command maps a local subject to a remote service definition.

```shell
nsc add import --account B -u /tmp/activation.jwt --local-subject private.help --name private.help
```

--------------------------------

### Account Information (/accountz)

Source: https://github.com/nats-io/nats.docs/blob/master/running-a-nats-service/nats_admin/monitoring/README.md

The /accountz endpoint provides information on the server's active accounts. It can list all accounts or provide details for a specific account.

```APIDOC
## GET /accountz

### Description
Retrieves information about NATS accounts. By default, it lists all accounts known to the server. It can also provide detailed metrics for a specific account when the `acc` parameter is provided.

### Method
GET

### Endpoint
/accountz

### Parameters
#### Query Parameters
- **acc** (account name) - Optional - Include metrics for the specified account. If not set, a list of all accounts is included.

### Response
#### Success Response (200)
- **server_id** (string) - The ID of the NATS server.
- **now** (string) - The current timestamp.
- **system_account** (string) - The system account identifier.
- **accounts** (array of strings) - A list of account names (when `acc` is not specified).
- **account_detail** (object) - Detailed information for a specific account (when `acc` is specified).
  - **account_name** (string) - The name of the account.
  - **update_time** (string) - The last update time for the account.
  - **is_system** (boolean) - Indicates if the account is a system account.
  - **expired** (boolean) - Indicates if the account has expired.
  - **complete** (boolean) - Indicates if the account information is complete.
  - **jetstream_enabled** (boolean) - Indicates if JetStream is enabled for the account.
  - **leafnode_connections** (integer) - Number of leafnode connections.
  - **client_connections** (integer) - Number of client connections.
  - **subscriptions** (integer) - Number of subscriptions.
  - **exports** (array of objects) - List of account exports.
  - **jwt** (string) - The account JWT.
  - **issuer_key** (string) - The issuer key of the JWT.
  - **name_tag** (string) - The name tag of the account.
  - **decoded_jwt** (object) - The decoded JWT payload.
  - **sublist_stats** (object) - Statistics related to subscriptions within the account.

#### Response Example (Default behavior)
```json
{
  "server_id": "NAB2EEQ3DLS2BHU4K2YMXMPIOOOAOFOAQAC5NQRIEUI4BHZKFBI4ZU4A",
  "now": "2021-02-08T17:31:29.551146-05:00",
  "system_account": "AAAXAUVSGK7TCRHFIRAS4SYXVJ76EWDMNXZM6ARFGXP7BASNDGLKU7A5",
  "accounts": ["AAAXAUVSGK7TCRHFIRAS4SYXVJ76EWDMNXZM6ARFGXP7BASNDGLKU7A5", "$G"]
}
```

#### Response Example (Specific account)
```json
{
  "server_id": "NAB2EEQ3DLS2BHU4K2YMXMPIOOOAOFOAQAC5NQRIEUI4BHZKFBI4ZU4A",
  "now": "2021-02-08T17:37:55.80856-05:00",
  "system_account": "AAAXAUVSGK7TCRHFIRAS4SYXVJ76EWDMNXZM6ARFGXP7BASNDGLKU7A5",
  "account_detail": {
    "account_name": "AAAXAUVSGK7TCRHFIRAS4SYXVJ76EWDMNXZM6ARFGXP7BASNDGLKU7A5",
    "update_time": "2021-02-08T17:31:22.390334-05:00",
    "is_system": true,
    "expired": false,
    "complete": true,
    "jetstream_enabled": false,
    "leafnode_connections": 0,
    "client_connections": 0,
    "subscriptions": 42,
    "exports": [
      {
        "subject": "$SYS.DEBUG.SUBSCRIBERS",
        "type": "service",
        "response_type": "Singleton"
      }
    ],
    "jwt": "eyJ0eXAiOiJqd3QiLCJhbGciOiJlZDI1NTE5In0.eyJqdGkiOiJVVlU2VEpXRU8zS0hYWTZVMkgzM0RCVklET1A3U05DTkJPMlM0M1dPNUM2T1RTTDNVSUxBIiwiaWF0IjoxNjAzNDczNzg4LCJpc3MiOiJPQlU1T05GSjMyNFVEUFJCSVZSR0Y3Q05FT0hHTFBTN0VZUEJUVlFaS1NCSElJWklCNkhENjZKRiIsIm5hbWUiOiJTWVMiLCJzdWIiOiJBQUFYQVVWU0dLN1RDUkhGSVJBUzRTWVhWSjc2RVdETU5YWk02QVJGR1hQN0JBU05ER0xLVTdBNSIsInR5cGUiOiJhY2NvdW50IiwibmF0cyI6eyJsaW1pdHMiOnsic3VicyI6LTEsImNvbm4iOi0xLCJsZWFmIjotMSwiaW1wb3J0cyI6LTEsImV4cG9ydHMiOi0xLCJkYXRhIjotMSwicGF5bG9hZCI6LTEsIndpbGRjYXJkcyI6dHJ1ZX19fQ.CeGo16i5oD0b1uBJ8UdGmLH-l9dL8yNqXHggkAt2T5c88fM7k4G08wLguMAnlvzrdlYvdZvOx_5tHLuDZmGgCg",
    "issuer_key": "OBU5O5FJ324UDPRBIVRGF7CNEOHGLPS7EYPBTVQZKSBHIIZIB6HD66JF",
    "name_tag": "SYS",
    "decoded_jwt": {
      "jti": "UVU6TJWEO3KHXY6U2H33DBVIDOP7SNCNBO2S43WO5C6OTSL3UILA",
      "iat": 1603473788,
      "iss": "OBU5O5FJ324UDPRBIVRGF7CNEOHGLPS7EYPBTVQZKSBHIIZIB6HD66JF",
      "name": "SYS",
      "sub": "AAAXAUVSGK7TCRHFIRAS4SYXVJ76EWDMNXZM6ARFGXP7BASNDGLKU7A5",
      "nats": {
        "limits": {
          "subs": -1,
          "data": -1,
          "payload": -1,
          "imports": -1,
          "exports": -1,
          "wildcards": true,
          "conn": -1,
          "leaf": -1
        },
        "default_permissions": {
          "pub": {},
          "sub": {}
        },
        "type": "account",
        "version": 1
      }
    },
    "sublist_stats": {
      "num_subscriptions": 42,
      "num_cache": 6,
      "num_inserts": 42,
      "num_removes": 0,
      "num_matches": 6,
      "cache_hit_rate": 0,
      "max_fanout": 1,
      "avg_fanout": 0.8333333333333334
    }
  }
}
```
```

--------------------------------

### Generate Activation Token

Source: https://github.com/nats-io/nats.docs/blob/master/using-nats/nats-tools/nsc/services.md

Creates an activation token JWT for a target account. It maps the exporting account, the target account's public key, and the specific subject to an output file.

```shell
nsc generate activation --account A --target-account AAM46E3YF5WOZSE5WNYWHN3YYISVZOSI6XHTF2Q64ECPXSFQZROJMP2H --subject private.help.AAM46E3YF5WOZSE5WNYWHN3YYISVZOSI6XHTF2Q64ECPXSFQZROJMP2H -o /tmp/activation.jwt
```

--------------------------------

### Create an Account

Source: https://github.com/nats-io/nats.docs/blob/master/using-nats/nats-tools/nsc/basics.md

Adds a new account entity to the current operator.

```bash
nsc add account MyAccount
```

--------------------------------

### Test Private Service Communication

Source: https://github.com/nats-io/nats.docs/blob/master/using-nats/nats-tools/nsc/services.md

Demonstrates testing a private service by setting up a responder on the source account and sending a request from the importing account.

```bash
nsc reply --account A --user U "private.help.*" "help is here"
nsc req --account B --user b private.help help_me
```

--------------------------------

### Push and Pull Account Changes

Source: https://github.com/nats-io/nats.docs/blob/master/using-nats/nats-tools/nsc/basics.md

Synchronize account JWTs between the local NSC environment and the NATS server.

```shell
nsc push -a MyAccount -u nats://localhost
```

```shell
nsc pull -u nats://localhost
```

--------------------------------

### Server Configuration: Account-Specific Mappings

Source: https://github.com/nats-io/nats.docs/blob/master/nats-concepts/subject_mapping.md

Defines subject mappings within a specific account configuration. This allows for granular control over subject remapping per account.

```yaml
server_name: "hub"
cluster: { name: "hub" }

accounts {
    accountA: { 
        mappings: {
            orders.flush  orders.central.flush 
        	orders.* orders.central.{{wildcard(1)}}
        }
    }
}
```

--------------------------------

### CLI: Generate Activation Token

Source: https://github.com/nats-io/nats.docs/blob/master/using-nats/nats-tools/nsc/streams.md

Generates a JWT-based activation token that grants a target account permission to import a specific stream subject.

```APIDOC
## NSC generate activation

### Description
Generates an activation token for a foreign account to import a private stream. This token defines the permissions and subject mapping for the imported stream.

### Method
CLI Command

### Parameters
#### Flags
- **--account** (string) - Required - The account name that owns the stream export.
- **--target-account** (string) - Required - The public key of the account importing the service.
- **--subject** (string) - Required - The subject where the stream will be published to the target account.
- **-o** (string) - Required - The output file path for the generated JWT.

### Request Example
`nsc generate activation --account A --target-account AAM46E3YF5WOZSE5WNYWHN3YYISVZOSI6XHTF2Q64ECPXSFQZROJMP2H --subject private.abc.AAM46E3YF5WOZSE5WNYWHN3YYISVZOSI6XHTF2Q64ECPXSFQZROJMP2H -o /tmp/activation.jwt`

### Response
#### Success Response
- **Output** (file) - A JWT file containing the activation details, which can be shared with the importing account.
```

--------------------------------

### Add NATS Users with Team Tags for Templated Permissions

Source: https://github.com/nats-io/nats.docs/blob/master/using-nats/nats-tools/nsc/signing_keys.md

These commands add two NATS users, 'pam' and 'joe', to the 'sales' account using the 'team-service' signing key. Each user is assigned a 'team' tag ('support' for pam, 'leads' for joe), which will be used by the template functions in the signing key to define their specific subscription subjects.

```shell
nsc add user pam -K team-service --tag team:support
nsc add user joe -K team-service --tag team:leads
```

--------------------------------

### Generate Activation Token

Source: https://github.com/nats-io/nats.docs/blob/master/using-nats/nats-tools/nsc/streams.md

Creates an activation token for a target account, allowing it to access a specific exported stream subject. The output is saved as a JWT file.

```bash
nsc generate activation --account A --target-account AAM46E3YF5WOZSE5WNYWHN3YYISVZOSI6XHTF2Q64ECPXSFQZROJMP2H --subject private.abc.AAM46E3YF5WOZSE5WNYWHN3YYISVZOSI6XHTF2Q64ECPXSFQZROJMP2H -o /tmp/activation.jwt
```

--------------------------------

### Initialize and Add Managed Operators

Source: https://github.com/nats-io/nats.docs/blob/master/using-nats/nats-tools/nsc/managed.md

Commands to add a managed operator manually or initialize a new account with a well-known operator.

```bash
nsc add operator -i
nsc init -o synadia -n MyFirstAccount
```