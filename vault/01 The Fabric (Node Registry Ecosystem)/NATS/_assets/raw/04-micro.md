### General Server Information Endpoint (/varz)

Source: https://github.com/nats-io/nats.docs/blob/master/running-a-nats-service/nats_admin/monitoring/README.md

Access the /varz endpoint to retrieve general information about the NATS server's state and configuration. This endpoint returns a JSON object with details such as server version, host, port, connection counts, message statistics, and configuration load time.

```json
{
  "server_id": "NACDVKFBUW4C4XA24OOT6L4MDP56MW76J5RJDFXG7HLABSB46DCMWCOW",
  "version": "2.0.0",
  "proto": 1,
  "go": "go1.12",
  "host": "0.0.0.0",
  "port": 4222,
  "max_connections": 65536,
  "ping_interval": 120000000000,
  "ping_max": 2,
  "http_host": "0.0.0.0",
  "http_port": 8222,
  "https_port": 0,
  "auth_timeout": 1,
  "max_control_line": 4096,
  "max_payload": 1048576,
  "max_pending": 67108864,
  "cluster": {},
  "gateway": {},
  "leaf": {},
  "tls_timeout": 0.5,
  "write_deadline": 2000000000,
  "start": "2019-06-24T14:24:43.928582-07:00",
  "now": "2019-06-24T14:24:46.894852-07:00",
  "uptime": "2s",
  "mem": 9617408,
  "cores": 4,
  "gomaxprocs": 4,
  "cpu": 0,
  "connections": 0,
  "total_connections": 0,
  "routes": 0,
  "remotes": 0,
  "leafnodes": 0,
  "in_msgs": 0,
  "out_msgs": 0,
  "in_bytes": 0,
  "out_bytes": 0,
  "slow_consumers": 2,
  "subscriptions": 0,
  "http_req_stats": {
    "/": 0,
    "/connz": 0,
    "/gatewayz": 0,
    "/routez": 0,
    "/subsz": 0,
    "/varz": 1
  },
  "config_load_time": "2019-06-24T14:24:43.928582-07:00",
  "slow_consumer_stats": {
    "clients": 1,
    "routes": 1,
    "gateways": 0,
    "leafs": 0
  }
}
```

--------------------------------

### Health Check Endpoint

Source: https://github.com/nats-io/nats.docs/blob/master/running-a-nats-service/nats_admin/monitoring/README.md

The `/healthz` endpoint checks if the NATS server is operational and can accept connections. It supports query parameters to filter health checks, particularly for JetStream.

```APIDOC
## GET /healthz

### Description
Checks if the NATS server is able to accept connections. Returns OK if the server is healthy.

### Method
GET

### Endpoint
/healthz

### Parameters
#### Query Parameters
- **js-enabled-only** (boolean) - Optional - Returns an error if JetStream is disabled.
- **js-server-only** (boolean) - Optional - Skips health check of accounts, streams, and consumers.
- **js-enabled** (boolean) - Optional - Returns an error if JetStream is disabled. (Deprecated: use `js-enabled-only` instead).

### Request Example
* Default: `https://demo.nats.io:8222/healthz`
* Expect JetStream: `https://demo.nats.io:8222/healthz?js-enabled-only=true`

### Response
#### Success Response (200 OK)
- **status** (string) - Indicates the server status, typically "ok".

#### Response Example
```json
{
  "status": "ok"
}
```
```

--------------------------------

### NATS Monitoring Endpoint with JSONP Callback

Source: https://github.com/nats-io/nats.docs/blob/master/running-a-nats-service/nats_admin/monitoring/README.md

This example demonstrates how to fetch monitoring data from a NATS endpoint (e.g., `/connz`) and utilize JSONP for cross-origin requests by specifying a `callback` query parameter.

```javascript
$.getJSON("https://demo.nats.io:8222/connz?callback=?", function (data) {
  console.log(data);
});
```

--------------------------------

### Describe Service User

Source: https://github.com/nats-io/nats.docs/blob/master/using-nats/nats-tools/nsc/basics.md

Displays the configuration details for the service user.

```shell
nsc describe user s
```

--------------------------------

### Import Private Service with NSC

Source: https://github.com/nats-io/nats.docs/blob/master/using-nats/nats-tools/nsc/services.md

Adds a private service import to an account using an activation token file. This command maps a local subject to a remote service definition.

```shell
nsc add import --account B -u /tmp/activation.jwt --local-subject private.help --name private.help
```

--------------------------------

### Account Statistics Endpoint

Source: https://github.com/nats-io/nats.docs/blob/master/running-a-nats-service/nats_admin/monitoring/README.md

Retrieve statistics for NATS accounts. Use the `unused=1` argument to include accounts without active connections.

```json
{
  "server_id": "NDJ5M4F5WAIBUA26NJ3QMH532AQPN7QNTJP3Y4SBHSHL4Y7QUAKNJEAF",
  "now": "2022-10-19T17:16:20.881296749Z",
  "account_statz": [
    {
      "acc": "default",
      "conns": 31,
      "leafnodes": 2,
      "total_conns": 33,
      "num_subscriptions": 45,
      "sent": {
        "msgs": 1876970,
        "bytes": 246705616
      },
      "received": {
        "msgs": 1347454,
        "bytes": 219438308
      },
      "slow_consumers": 29
    },
    {
      "acc": "$G",
      "conns": 1,
      "leafnodes": 0,
      "total_conns": 1,
      "num_subscriptions": 3,
      "sent": {
        "msgs": 0,
        "bytes": 0
      },
      "received": {
        "msgs": 107,
        "bytes": 1094
      },
      "slow_consumers": 0
    }
  ]
}
```

--------------------------------

### JetStream Information Endpoint

Source: https://github.com/nats-io/nats.docs/blob/master/running-a-nats-service/nats_admin/monitoring/README.md

Retrieve detailed JetStream information. Supports filtering by account, streams, consumers, and configuration. Pagination is available via `limit` and `offset`.

```http
GET /jsz HTTP/1.1
Host: demo.nats.io:8222


```

```http
GET /jsz?accounts=true&limit=16&offset=128 HTTP/1.1
Host: demo.nats.io:8222


```

```http
GET /jsz?consumers=true HTTP/1.1
Host: demo.nats.io:8222


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

### Start NATS Service on Leafnode (Red)

Source: https://github.com/nats-io/nats.docs/blob/master/running-a-nats-service/running/nats_docker/ngs-leafnodes-docker.md

Starts a NATS service on the 'red' leaf node, configured to listen on localhost:4222. This service acts as a replier, responding to requests with a formatted message including the time and the received request.

```Shell
nats -s localhost:4222 reply docker-leaf-test "At {{Time}}, I received your request: {{Request}}"
```

--------------------------------

### Configure Operator Service URL

Source: https://github.com/nats-io/nats.docs/blob/master/using-nats/nats-tools/nsc/basics.md

Sets the service URL for an operator to specify where the NATS server is listening.

```bash
nsc edit operator --service-url nats://localhost:4222
```

--------------------------------

### JetStream Information (/jsz)

Source: https://github.com/nats-io/nats.docs/blob/master/running-a-nats-service/nats_admin/monitoring/README.md

The /jsz endpoint provides detailed information about JetStream, including account-specific metrics, streams, and consumers. It supports pagination and filtering.

```APIDOC
## GET /jsz

### Description
Retrieves detailed information about JetStream, including account statistics, streams, and consumers. This endpoint supports pagination and can be filtered by various parameters. In a clustered environment, it's recommended to query the leader for the most accurate data.

### Method
GET

### Endpoint
/jsz

### Query Parameters
- **acc** (string) - Optional - Include metrics for the specified account. Default is unset.
- **accounts** (boolean) - Optional - Include account-specific JetStream information. Default is false.
- **streams** (boolean) - Optional - Include streams. When set, implies `accounts=true`. Default is false.
- **consumers** (boolean) - Optional - Include consumer information. When set, implies `streams=true`. Default is false.
- **config** (boolean) - Optional - When stream or consumer are requested, include their respective configuration. Default is false.
- **leader-only** (boolean) - Optional - Only the leader responds. Default is false.
- **offset** (number) - Optional - Pagination offset. Must be greater than 0. Default is 0.
- **limit** (number) - Optional - Number of results to return. Must be greater than 0. Default is 1024.
- **raft** (boolean) - Optional - Include information details about the Raft group. Default is false.

### Response
#### Success Response (200)
- The response structure varies based on the query parameters used. It can include account information, stream details, consumer metrics, and Raft group information.

#### Error Response (400)
- Bad Request

### Request Example
- Basic JetStream information:
  `https://demo.nats.io:8222/jsz`
- Requesting accounts with custom limit and offset:
  `https://demo.nats.io:8222/jsz?accounts=true&limit=16&offset=128`
- Requesting detailed consumer information:
  `https://demo.nats.io:8222/jsz?consumers=true`
```

--------------------------------

### Account Statistics (/accstatz)

Source: https://github.com/nats-io/nats.docs/blob/master/running-a-nats-service/nats_admin/monitoring/README.md

The /accstatz endpoint provides statistics for NATS accounts, including connection counts and data transfer metrics. It can optionally include accounts without active connections.

```APIDOC
## GET /accstatz

### Description
Retrieves statistics for NATS accounts. By default, it only includes accounts with active connections. The `unused` query parameter can be used to include accounts that do not have any current connections.

### Method
GET

### Endpoint
/accstatz

### Query Parameters
- **unused** (boolean) - Optional - If true, include accounts that do not have any current connections. Default is false.

### Response
#### Success Response (200)
- **server_id** (string) - The ID of the NATS server.
- **now** (string) - The timestamp when the statistics were generated.
- **account_statz** (array) - An array of account statistics objects.
  - **acc** (string) - The name of the account.
  - **conns** (integer) - The number of active connections for the account.
  - **leafnodes** (integer) - The number of active leafnodes for the account.
  - **total_conns** (integer) - The total number of connections (including leafnodes).
  - **num_subscriptions** (integer) - The number of subscriptions in the account.
  - **sent** (object) - Data sent by the account.
    - **msgs** (integer) - Number of messages sent.
    - **bytes** (integer) - Number of bytes sent.
  - **received** (object) - Data received by the account.
    - **msgs** (integer) - Number of messages received.
    - **bytes** (integer) - Number of bytes received.
  - **slow_consumers** (integer) - The number of slow consumers.

### Response Example
```json
{
  "server_id": "NDJ5M4F5WAIBUA26NJ3QMH532AQPN7QNTJP3Y4SBHSHL4Y7QUAKNJEAF",
  "now": "2022-10-19T17:16:20.881296749Z",
  "account_statz": [
    {
      "acc": "default",
      "conns": 31,
      "leafnodes": 2,
      "total_conns": 33,
      "num_subscriptions": 45,
      "sent": {
        "msgs": 1876970,
        "bytes": 246705616
      },
      "received": {
        "msgs": 1347454,
        "bytes": 219438308
      },
      "slow_consumers": 29
    }
  ]
}
```
```

--------------------------------

### Control NATS Windows Service with Signals

Source: https://github.com/nats-io/nats.docs/blob/master/running-a-nats-service/running/windows_srv.md

Demonstrates how to send signals to the NATS server service for common operations like reloading configuration, rotating logs, or stopping the server. The service name can be specified if it differs from the default 'nats-server'.

```shell
REM Reload server configuration
nats-server.exe --signal reload

REM Reopen log file for log rotation
nats-server.exe --signal reopen

REM Stop the server
nats-server.exe --signal stop

REM Stop a specific service
nats-server.exe --signal stop=<service name>
```

--------------------------------

### Connection Information (`/connz`)

Source: https://github.com/nats-io/nats.docs/blob/master/running-a-nats-service/nats_admin/monitoring/README.md

The `/connz` endpoint reports information on active client connections to a NATS server. It provides details about each connection, including its ID, type, IP address, start time, last activity, and resource usage.

```APIDOC
## GET /connz

### Description
Retrieves information about active client connections to the NATS server.

### Method
GET

### Endpoint
/connz

### Query Parameters
- **subs** (boolean | integer | string) - Optional - Controls the inclusion of subscription details. Can be `true`, `1`, `false`, `0`, or `detail` for more verbose subscription information.

### Response
#### Success Response (200)
- **server_id** (string) - The ID of the NATS server.
- **now** (string) - The current timestamp.
- **num_connections** (integer) - The total number of active connections.
- **total** (integer) - The total number of connections (including potentially inactive ones).
- **offset** (integer) - The offset for pagination (if applicable).
- **limit** (integer) - The limit for pagination (if applicable).
- **connections** (array) - An array of connection objects, each containing:
  - **cid** (integer) - The client connection ID.
  - **kind** (string) - The type of client (e.g., 'Client').
  - **type** (string) - The protocol type (e.g., 'nats', 'mqtt').
  - **ip** (string) - The IP address of the client.
  - **port** (integer) - The port the client is connected on.
  - **start** (string) - The timestamp when the connection started.
  - **last_activity** (string) - The timestamp of the last client activity.
  - **rtt** (string) - The round-trip time to the client.
  - **uptime** (string) - The duration the client has been connected.
  - **idle** (string) - The duration since the last client activity.
  - **pending_bytes** (integer) - The number of bytes pending to be sent to the client.
  - **in_msgs** (integer) - The total number of messages received from the client.
  - **out_msgs** (integer) - The total number of messages sent to the client.
  - **in_bytes** (integer) - The total number of bytes received from the client.
  - **out_bytes** (integer) - The total number of bytes sent to the client.
  - **subscriptions** (integer) - The number of subscriptions the client has.
  - **name** (string) - The name of the client (if set).
  - **lang** (string) - The programming language of the client library.
  - **version** (string) - The version of the client library.
  - **mqtt_client** (string) - The MQTT client ID (if applicable).

#### Response Example
```json
{
  "server_id": "NACDVKFBUW4C4XA24OOT6L4MDP56MW76J5RJDFXG7HLABSB46DCMWCOW",
  "now": "2019-06-24T14:28:16.520365-07:00",
  "num_connections": 2,
  "total": 2,
  "offset": 0,
  "limit": 1024,
  "connections": [
    {
      "cid": 5,
      "kind": "Client",
      "type": "nats",
      "ip": "127.0.0.1",
      "port": 62714,
      "start": "2021-09-09T23:16:43.040862Z",
      "last_activity": "2021-09-09T23:16:43.042364Z",
      "rtt": "95µs",
      "uptime": "5s",
      "idle": "5s",
      "pending_bytes": 0,
      "in_msgs": 0,
      "out_msgs": 0,
      "in_bytes": 0,
      "out_bytes": 0,
      "subscriptions": 1,
      "name": "NATS Benchmark",
      "lang": "go",
      "version": "1.12.1"
    }
  ]
}
```
```

--------------------------------

### Configure Service User Permissions

Source: https://github.com/nats-io/nats.docs/blob/master/using-nats/nats-tools/nsc/basics.md

Grants a service user permission to subscribe to a request subject and publish to inbox subjects.

```bash
nsc add user s --allow-pub "_INBOX.>" --allow-sub q
```

--------------------------------

### Configure Operator Service URLs

Source: https://github.com/nats-io/nats.docs/blob/master/using-nats/nats-tools/nsc/managed.md

Setting service URLs for an operator to enable tool actions like publishing messages.

```bash
nsc edit operator -n nats://localhost:4222
nsc tool pub hello world
```

--------------------------------

### Account Statistics (/accstatz)

Source: https://github.com/nats-io/nats.docs/blob/master/running-a-nats-service/nats_admin/monitoring/README.md

The /accstatz endpoint provides per-account statistics, such as connection counts and message/byte transfer rates.

```APIDOC
## GET /accstatz

### Description
Retrieves per-account statistics, including connection counts and message/byte transfer rates.

### Method
GET

### Endpoint
/accstatz

### Parameters
This endpoint does not appear to accept any query parameters based on the provided documentation.

### Response
#### Success Response (200)
- The response structure for this endpoint is not detailed in the provided documentation. It is expected to contain various statistical metrics for accounts.

#### Response Example
(No example provided in the source text)
```

--------------------------------

### Go: Distributed User Provisioning Service

Source: https://github.com/nats-io/nats.docs/blob/master/running-a-nats-service/nats_admin/jwt.md

Sets up a service to provision users by handling JWT requests. It uses Go channels for communication and closures for authorization logic. The service generates user JWTs based on provided public keys and authorization information.

```go
func ObtainAuthorizationToken() interface{} {
    // whatever you want, 3rd party token/username&password
    return ""
}

func IsTokenAuthorized(token interface{}) bool {
    // whatever logic to determine if the input authorizes the requester to obtain a user jwt
    return token.(string) == ""
}

// request struct to exchange data
type userRequest struct {
    UserJWTResponseChan chan string
    UserPublicKey       string
    AuthInfo            interface{}
}

func startUserProvisioningService(isAuthorizedCb func(token interface{}) bool) chan userRequest {
    userRequestChan := make(chan userRequest) // channel to send requests for jwt to
    go func() {
        accountSigningKey := GetAccountSigningKey() // Setup, obtain account signing key
        for {
            req := <-userRequestChan // receive request
            if !isAuthorizedCb(req.AuthInfo) {
                fmt.Printf("Request is not authorized to receive a JWT, timeout on purpose")
            } else if userJWT := generateUserJWT(req.UserPublicKey, accountSigningKey); userJWT != "" {
                req.UserJWTResponseChan <- userJWT // respond with jwt
            }
        }
    }()
    return userRequestChan
}
```