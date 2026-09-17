### NATS Monitoring Endpoint with JSONP Callback

Source: https://github.com/nats-io/nats.docs/blob/master/running-a-nats-service/nats_admin/monitoring/README.md

This example demonstrates how to fetch monitoring data from a NATS endpoint (e.g., `/connz`) and utilize JSONP for cross-origin requests by specifying a `callback` query parameter.

```javascript
$.getJSON("https://demo.nats.io:8222/connz?callback=?", function (data) {
  console.log(data);
});
```

--------------------------------

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

### Enable NATS Monitoring via Configuration File

Source: https://github.com/nats-io/nats.docs/blob/master/running-a-nats-service/configuration/monitoring.md

Configures the NATS server to enable its HTTP monitoring endpoint using a YAML configuration file. The `http_port` directive specifies the port for monitoring. Binding to `localhost` enhances security by limiting access.

```yaml
http_port: 8222
```

```yaml
http: localhost:8222
```

--------------------------------

### Enable Monitoring via Server Configuration

Source: https://github.com/nats-io/nats.docs/blob/master/running-a-nats-service/nats_admin/monitoring/README.md

Configure the HTTP port for monitoring in the NATS server configuration file.

```yaml
http_port: 8222
```

--------------------------------

### General Information

Source: https://github.com/nats-io/nats.docs/blob/master/running-a-nats-service/nats_admin/monitoring/README.md

The `/varz` endpoint returns general information about the NATS server's state and configuration.

```APIDOC
## GET /varz

### Description
Retrieves general information about the NATS server's state and configuration.

### Method
GET

### Endpoint
/varz

### Parameters
N/A

### Response
#### Success Response (200)
Returns a JSON object containing server statistics and configuration details.

#### Response Example
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
```

--------------------------------

### Enable NATS Monitoring via Command Line

Source: https://github.com/nats-io/nats.docs/blob/master/running-a-nats-service/configuration/monitoring.md

Starts the NATS server with an HTTP monitoring endpoint on a specified port. This allows for real-time statistics and information retrieval. Ensure the port is not exposed publicly without proper security measures.

```bash
nats-server -m 8222
```

--------------------------------

### Verify NATS Cluster Routes with curl

Source: https://github.com/nats-io/nats.docs/blob/master/running-a-nats-service/running/nats_docker/README.md

This snippet shows how to use curl to query the monitoring endpoint of a NATS server to verify that cluster routes are connected. It requests the /routez endpoint and expects to see the number of connected routes.

```bash
curl http://127.0.0.1:8222/routez
```

--------------------------------

### Start NATS Server with Monitoring

Source: https://github.com/nats-io/nats.docs/blob/master/using-nats/nats-tools/nats_cli/natsbench.md

Starts the NATS server instance with the HTTP monitoring port enabled and JetStream support active. This is a prerequisite for running performance benchmarks.

```bash
nats-server -m 8222 -js
```

--------------------------------

### Subscription Routing Information

Source: https://github.com/nats-io/nats.docs/blob/master/running-a-nats-service/nats_admin/monitoring/README.md

The /subsz endpoint reports detailed information about the current subscriptions and the routing data structure. It is not normally used.

```APIDOC
## GET /subsz

### Description
Retrieves detailed information about current subscriptions and the routing data structure.

### Method
GET

### Endpoint
/subsz

### Response
#### Success Response (200)
(Response structure not detailed in source)

#### Error Response (400)
(Error details not detailed in source)
```

--------------------------------

### Connection Information

Source: https://github.com/nats-io/nats.docs/blob/master/running-a-nats-service/nats_admin/monitoring/README.md

The `/connz` endpoint reports detailed information on current and recently closed connections, supporting paging.

```APIDOC
## GET /connz

### Description
Retrieves detailed information about current and recently closed connections to the NATS server. Supports paging with a default of 1024 connections.

### Method
GET

### Endpoint
/connz

### Parameters
N/A

### Response
#### Success Response (200)
Returns a JSON object containing connection details.

#### Error Response (400)
Returned if the request is malformed.

#### Response Example
(Response example not provided in source text)
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

### Monitor NATS Connection Status and Statistics

Source: https://github.com/nats-io/nats.docs/blob/master/using-nats/developing-with-nats/events/README.md

Retrieves the current connection state, checks for reconnection events, and accesses client statistics to monitor the health of the NATS connection.

```go
nc, err := nats.Connect("demo.nats.io", nats.Name("API Example"))
if err != nil {
    log.Fatal(err)
}
defer nc.Close()

getStatusTxt := func(nc *nats.Conn) string {
    switch nc.Status() {
    case nats.CONNECTED:
        return "Connected"
    case nats.CLOSED:
        return "Closed"
    default:
        return "Other"
    }
}
log.Printf("The connection is %v\n", getStatusTxt(nc))

nc.Close()

log.Printf("The connection is %v\n", getStatusTxt(nc))
```

```java
Connection nc = Nats.connect("nats://demo.nats.io:4222");

System.out.println("The Connection is: " + nc.getStatus()); // CONNECTED

nc.close();

System.out.println("The Connection is: " + nc.getStatus()); // CLOSED
```

```javascript
// you can find out where you connected:
t.log(`connected to a nats server version ${nc.info.version}`);

// or information about the data in/out of the client:
const stats = nc.stats();
t.log(`client sent ${stats.outMsgs} messages and received ${stats.inMsgs}`);
```

```python
nc = NATS()

await nc.connect(
   servers=["nats://demo.nats.io:4222"],
   )

# Do something with the connection.

print("The connection is connected?", nc.is_connected)

while True:
  if nc.is_reconnecting:
    print("Reconnecting to NATS...")
    break
  await asyncio.sleep(1)

await nc.close()

print("The connection is closed?", nc.is_closed)
```

```csharp
// dotnet add package NATS.Net
using NATS.Net;

await using var client = new NatsClient();

Console.WriteLine($"{client.Connection.ConnectionState}"); // Closed

await client.ConnectAsync();

Console.WriteLine($"{client.Connection.ConnectionState}"); // Open
```

```ruby
NATS.start(max_reconnect_attempts: 2) do |nc|
  puts "Connect is connected?: #{nc.connected?}"

  timer = EM.add_periodic_timer(1) do
    if nc.closing?
      puts "Connection closed..."
      EM.cancel_timer(timer)
      NATS.stop
    end

    if nc.reconnecting?
      puts "Reconnecting to NATS..."
      next
    end
  end
end
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

### NATS System Monitoring Requests

Source: https://github.com/nats-io/nats.docs/blob/master/running-a-nats-service/configuration/sys_accounts/README.md

Request-reply patterns to query server and account-specific monitoring data.

```APIDOC
## REQUEST $SYS.REQ.SERVER.<id>.<endpoint>

### Description
Sends a request to a specific server to retrieve monitoring data such as stats, connections, or health status.

### Method
REQUEST

### Endpoint
$SYS.REQ.SERVER.<id>.<endpoint>

### Parameters
#### Path Parameters
- **id** (string) - Required - The unique server identifier.
- **endpoint** (string) - Required - The monitoring endpoint name (e.g., VARZ, CONNZ, ROUTEZ, GATEWAYZ, LEAFZ, SUBSZ, JSZ, ACCOUNTZ, HEALTHZ).

### Response
#### Success Response (200)
- **Payload** (JSON) - The requested monitoring data corresponding to the endpoint name.
```

--------------------------------

### Enable Monitoring via Command-Line Option

Source: https://github.com/nats-io/nats.docs/blob/master/running-a-nats-service/nats_admin/monitoring/README.md

Start the NATS server with a command-line flag to specify the monitoring port.

```bash
nats-server -m 8222
```

--------------------------------

### Subscription Information (/subsz)

Source: https://github.com/nats-io/nats.docs/blob/master/running-a-nats-service/nats_admin/monitoring/README.md

The /subsz endpoint provides information about the server's subscriptions, including counts and routing details. It supports filtering and pagination.

```APIDOC
## GET /subsz

### Description
Retrieves information about NATS subscriptions on the server. Supports filtering by subscription status and pagination.

### Method
GET

### Endpoint
/subsz

### Parameters
#### Query Parameters
- **subs** (boolean) - Optional - Include subscriptions. Default is false.
- **offset** (integer) - Optional - Pagination offset. Default is 0.
- **limit** (integer) - Optional - Number of results to return. Default is 1024.
- **test** (subject) - Optional - Test whether a subscription exists.

### Response
#### Success Response (200)
- **num_subscriptions** (integer) - The total number of subscriptions.
- **num_cache** (integer) - The number of cached subscriptions.
- **num_inserts** (integer) - The number of inserted subscriptions.
- **num_removes** (integer) - The number of removed subscriptions.
- **num_matches** (integer) - The number of matched subscriptions.
- **cache_hit_rate** (number) - The cache hit rate.
- **max_fanout** (integer) - The maximum fanout.
- **avg_fanout** (number) - The average fanout.

#### Response Example
```json
{
  "num_subscriptions": 2,
  "num_cache": 0,
  "num_inserts": 2,
  "num_removes": 0,
  "num_matches": 0,
  "cache_hit_rate": 0,
  "max_fanout": 0,
  "avg_fanout": 0
}
```
```