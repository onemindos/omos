### Gateway Information (`/gatewayz`)

Source: https://github.com/nats-io/nats.docs/blob/master/running-a-nats-service/nats_admin/monitoring/README.md

The `/gatewayz` endpoint provides information about gateways participating in a NATS supercluster. It details connections between NATS clusters, including gateway names and account information.

```APIDOC
## GET /gatewayz

### Description
Retrieves information about gateways used to create a NATS supercluster.

### Method
GET

### Endpoint
/gatewayz

### Query Parameters
- **accs** (boolean | integer) - Optional - Include account information. Accepts `true`, `1`, `false`, or `0`.
- **gw_name** (string) - Optional - Filter by a specific remote gateway name.
- **acc_name** (string) - Optional - Limit the list of accounts to this account name.

### Response
#### Success Response (200)
- **server_id** (string) - The ID of the NATS server.
- **now** (string) - The current timestamp.
- **num_gateways** (integer) - The total number of active gateways.
- **gateways** (array) - An array of gateway objects, each containing:
  - **rid** (integer) - The gateway connection ID.
  - **remote_id** (string) - The ID of the remote gateway.
  - **did_solicit** (boolean) - Indicates if this server solicited the connection.
  - **ip** (string) - The IP address of the remote gateway.
  - **port** (integer) - The port of the remote gateway.
  - **pending_size** (integer) - The size of pending messages to be sent.
  - **in_msgs** (integer) - The total number of messages received from the remote gateway.
  - **out_msgs** (integer) - The total number of messages sent to the remote gateway.
  - **in_bytes** (integer) - The total number of bytes received from the remote gateway.
  - **out_bytes** (integer) - The total number of bytes sent to the remote gateway.
  - **subscriptions** (integer) - The number of subscriptions advertised by the remote gateway.
  - **accounts** (array) - An array of account objects associated with the gateway (if `accs` is true).

#### Example
* Retrieve Gateway Information:
  `https://demo.nats.io:8222/gatewayz`
```

--------------------------------

### Configure NATS Gateways with Multiple Endpoints

Source: https://github.com/nats-io/nats.docs/blob/master/running-a-nats-service/configuration/gateways/gateway.md

Define gateway connections to remote NATS clusters using the `gateways` configuration block. Specify multiple `urls` for each remote gateway to ensure high availability and load balancing.

```hcl
gateway {
    name: "DC-A"
    listen: "localhost:7222"

    gateways: [
        {name: "DC-A", urls: ["nats://localhost:7222", "nats://localhost:7223", "nats://localhost:7224"]},
        {name: "DC-B", urls: ["nats://localhost:7332", "nats://localhost:7333", "nats://localhost:7334"]},
        {name: "DC-C", urls: ["nats://localhost:7442", "nats://localhost:7333", "nats://localhost:7335"]}
    ]
}
```

--------------------------------

### NATS Server Cluster Configuration

Source: https://github.com/nats-io/nats.docs/blob/master/nats-concepts/jetstream/streams.md

Define the `name` field within the `cluster` block in the NATS server configuration to enable cluster-based stream placement. Each server in a supercluster must have a unique cluster name.

```hcl
cluster {
  name: aws-us-east1-c1
  # etc..
}
```

--------------------------------

### Cluster-Scoped Subject Mapping

Source: https://github.com/nats-io/nats.docs/blob/master/nats-concepts/subject_mapping.md

Define mappings that apply only to messages published from a specific cluster. This example maps `foo` to `foo.cluster_name` based on the cluster (`west`, `central`, `east`).

```json
mappings = {
        "foo":[
               {destination:"foo.west", weight: 100%, cluster: "west"},
               {destination:"foo.central", weight: 100%, cluster: "central"},
               {destination:"foo.east", weight: 100%, cluster: "east"}
        ]
}
```

--------------------------------

### Gateway Configuration Block

Source: https://github.com/nats-io/nats.docs/blob/master/running-a-nats-service/configuration/gateways/gateway.md

Defines the configuration for a NATS gateway, including its name, listen address, authorization credentials, and a list of other gateways to connect to.

```hcl
gateway {
    name: "A"
    listen: "localhost:7222"
    authorization {
        user: gwu
        password: gwp
    }

    gateways: [
        {name: "A", url: "nats://gwu:gwp@localhost:7222"},
        {name: "B", url: "nats://gwu:gwp@localhost:7333"},
        {name: "C", url: "nats://gwu:gwp@localhost:7444"},
    ]
}
```

--------------------------------

### NATS Server Gateway Startup Logs

Source: https://github.com/nats-io/nats.docs/blob/master/running-a-nats-service/configuration/gateways/gateway.md

Example log output from a NATS server starting up with gateway configuration, showing listening addresses, gateway connections, and server readiness.

```log
[85803] 2019/05/07 10:50:55.902474 [INF] Starting nats-server version 2.0.0
[85803] 2019/05/07 10:50:55.903669 [INF] Gateway name is A
[85803] 2019/05/07 10:50:55.903684 [INF] Listening for gateways connections on localhost:7222
[85803] 2019/05/07 10:50:55.903696 [INF] Address for gateway "A" is localhost:7222
[85803] 2019/05/07 10:50:55.903909 [INF] Listening for client connections on 0.0.0.0:4222
[85803] 2019/05/07 10:50:55.903914 [INF] Server id is NBHUDBF3TVJSWCDPG2HSKI4I2SBSPDTNYEXEMOFAZUZYXVA2IYRUGPZU
[85803] 2019/05/07 10:50:55.903917 [INF] Server is ready
[85803] 2019/05/07 10:50:56.830669 [INF] 127.0.0.1:50892 - gid:2 - Processing inbound gateway connection
[85803] 2019/05/07 10:50:56.830673 [INF] 127.0.0.1:50891 - gid:1 - Processing inbound gateway connection
[85803] 2019/05/07 10:50:56.831079 [INF] 127.0.0.1:50892 - gid:2 - Inbound gateway connection from "C" (NBHWDFO3KHANNI6UCEUL27VNWL7NWD2MC4BI4L2C7VVLFBSMZ3CRD7HE) registered
[85803] 2019/05/07 10:50:56.831211 [INF] 127.0.0.1:50891 - gid:1 - Inbound gateway connection from "B" (ND2UJB3GFUHXOQ2UUMZQGOCL4QVR2LRJODPZH7MIPGLWCQRARJBU27C3) registered
[85803] 2019/05/07 10:50:56.906103 [INF] Connecting to explicit gateway "B" (localhost:7333) at 127.0.0.1:7333
[85803] 2019/05/07 10:50:56.906104 [INF] Connecting to explicit gateway "C" (localhost:7444) at 127.0.0.1:7444
[85803] 2019/05/07 10:50:56.906404 [INF] 127.0.0.1:7333 - gid:3 - Creating outbound gateway connection to "B"
[85803] 2019/05/07 10:50:56.906444 [INF] 127.0.0.1:7444 - gid:4 - Creating outbound gateway connection to "C"
[85803] 2019/05/07 10:50:56.906647 [INF] 127.0.0.1:7444 - gid:4 - Outbound gateway connection to "C" (NBHWDFO3KHANNI6UCEUL27VNWL7NWD2MC4BI4L2C7VVLFBSMZ3CRD7HE) registered
[85803] 2019/05/07 10:50:56.906772 [INF] 127.0.0.1:7333 - gid:3 - Outbound gateway connection to "B" (ND2UJB3GFUHXOQ2UUMZQGOCL4QVR2LRJODPZH7MIPGLWCQRARJBU27C3) registered
```

--------------------------------

### Clustering and Connection Options

Source: https://github.com/nats-io/nats.docs/blob/master/running-a-nats-service/configuration/clustering/cluster_config.md

Configuration options related to NATS server clustering and connection behavior.

```APIDOC
## NATS Server Configuration Options

### Routes

#### Description
A list of other servers (URLs) to cluster with. Self-routes are ignored. If authentication via `token` or `username`/`password` is required, specify them as part of the URL.

### Connect Retries

#### Description
After how many failed connect attempts to give up establishing a connection to a discovered route. Default is `0`, do not retry. When enabled, attempts will be made once a second. This does not apply to explicitly configured routes.

### Connect Backoff

#### Description
Enable exponential backoff on reconnection attempts. If `true`, will start exponential backoff at 1 second up to 30 seconds.
Version Introduced: 2.12.0

### Authorization

#### Description
Authorization map for configuring cluster routes. When a single `username`/`password` is used, it defines the authentication mechanism this server expects, and how this server will authenticate itself when establishing a connection to a _discovered_ route. This will not be used for routes explicitly listed in `routes` and therefore have to be provided as part of the URL. With this authentication mode, either use the same credentials throughout the system or list every route explicitly on every server. If the `tls` configuration map specifies `verify_and_map` only provide the expected `username`. Here different certificates can be used, but they have to map to the same `username`. The authorization map also allows for `timeout` which is honored but `users` and `token` configuration are not supported and will prevent the server from starting. The `permissions` block is ignored.

### Pool Size

#### Description
The size of the connection pool used to distribute load across non-pinned accounts. Default is `3`. Refer to [v2 Routes](./v2_routes.md#connection-pooling) for details.
```

--------------------------------

### Subscribe to Gateway Messages

Source: https://github.com/nats-io/nats.docs/blob/master/running-a-nats-service/configuration/gateways/gateway.md

Command to subscribe to messages on a specific port, demonstrating message reception through a gateway.

```shell
nats sub -s localhost:4333 ">"
```

--------------------------------

### Publish Message Through Gateway

Source: https://github.com/nats-io/nats.docs/blob/master/running-a-nats-service/configuration/gateways/gateway.md

Command to publish a message to a subject on a specific port, demonstrating message forwarding via a gateway.

```shell
nats pub -s localhost:4444 foo bar
```

--------------------------------

### Central Cluster Subject Mapping

Source: https://github.com/nats-io/nats.docs/blob/master/nats-concepts/subject_mapping.md

Configures subject mapping on a central NATS server. It demonstrates mapping 'orders.*' to 'orders.central.{{wildcard(1)}}' or 'orders.>' to 'orders.central.>' using different wildcard syntaxes.

```yaml
server_name: "hub"
cluster: { name: "hub" }
mappings: {
	orders.* orders.central.{{wildcard(1)}}
}
```

```yaml
server_name: "hub"
cluster: { name: "hub" }
mappings: {
	orders.> orders.central.>
}
```

--------------------------------

### Start NATS Server with Gateway Configuration

Source: https://github.com/nats-io/nats.docs/blob/master/running-a-nats-service/configuration/gateways/gateway.md

Command to start a NATS server using a specified configuration file that includes gateway settings.

```shell
nats-server -c A.conf
```

--------------------------------

### Leaf Cluster Subject Mapping

Source: https://github.com/nats-io/nats.docs/blob/master/nats-concepts/subject_mapping.md

Applies subject mapping on a leaf NATS cluster, translating subjects from 'orders.central.*' to 'orders.local.{{wildcard(1)}}'.

```yaml
server_name: "store1"
cluster: { name: "store1" }
mappings: {
	orders.central.* orders.local.{{wildcard(1)}}
}
```

--------------------------------

### Received Message Output

Source: https://github.com/nats-io/nats.docs/blob/master/running-a-nats-service/configuration/gateways/gateway.md

Example output showing a message received by a subscriber, indicating successful message forwarding through the gateway.

```text
[#1] Received on [foo] : 'bar'
```

--------------------------------

### Remove Offline Peer from Cluster

Source: https://github.com/nats-io/nats.docs/blob/master/running-a-nats-service/configuration/clustering/jetstream_clustering/administration.md

Removes a node from the cluster that is expected to remain offline permanently. This triggers JetStream to re-replicate data to a new node based on placement rules.

```shell
nats server cluster peer-remove n4-c1 --user admin --password s3cr3t!
```

--------------------------------

### Verify NATS Cluster Routes with curl

Source: https://github.com/nats-io/nats.docs/blob/master/running-a-nats-service/running/nats_docker/README.md

This snippet shows how to use curl to query the monitoring endpoint of a NATS server to verify that cluster routes are connected. It requests the /routez endpoint and expects to see the number of connected routes.

```bash
curl http://127.0.0.1:8222/routez
```

--------------------------------

### Configure NATS Server Cluster and Routes

Source: https://github.com/nats-io/nats.docs/blob/master/running-a-nats-service/configuration/clustering/cluster_config.md

This configuration snippet defines a NATS cluster block, specifying the listening port for inbound routes, authorization credentials for peer servers, and a list of static routes to connect to other NATS servers.

```NATS Configuration
cluster {
  name: example

  listen: localhost:4244

  authorization {
    user: route_user
    password: pwd
    timeout: 0.5
  }

  routes = [
    nats://route_user:pwd@127.0.0.1:4245
    nats://route_user:pwd@127.0.0.1:4246
  ]
}
```

--------------------------------

### Connect to NATS Cluster in JavaScript

Source: https://github.com/nats-io/nats.docs/blob/master/using-nats/developing-with-nats/connecting/cluster.md

Provides an example of connecting to a NATS cluster using JavaScript. It utilizes the `connect` function with an array of server URIs and demonstrates closing the connection.

```JavaScript
const nc = await connect({
    servers: [
      "nats://demo.nats.io:4222",
      "nats://localhost:4222",
    ],
});
// Do something with the connection
doSomething();
// When done close it
await nc.close();
```

--------------------------------

### Connect to NATS Cluster in Go

Source: https://github.com/nats-io/nats.docs/blob/master/using-nats/developing-with-nats/connecting/cluster.md

Demonstrates how to connect to a NATS cluster using a list of server URLs in Go. It handles potential connection errors and ensures the connection is closed properly.

```Go
servers := []string{"nats://127.0.0.1:1222", "nats://127.0.0.1:1223", "nats://127.0.0.1:1224"}

nc, err := nats.Connect(strings.Join(servers, ","))
if err != nil {
    log.Fatal(err)
}
defers nc.Close()

// Do something with the connection
```

--------------------------------

### Connecting Additional NATS Servers to Cluster

Source: https://github.com/nats-io/nats.docs/blob/master/running-a-nats-service/configuration/clustering/README.md

This command starts a new NATS server, specifying its client port, cluster address, and routes to connect to the seed server. It's used to expand the NATS cluster.

```bash
nats-server -p 5222 -cluster nats://localhost:5248 -routes nats://localhost:4248 -D
```

--------------------------------

### Listen for New Servers in NATS Cluster

Source: https://github.com/nats-io/nats.docs/blob/master/using-nats/developing-with-nats/events/events.md

Enables clients to be notified when new servers are added or changed within a NATS cluster. This is essential for applications that need to be aware of the cluster topology.

```Go
// Be notified if a new server joins the cluster.
// Print all the known servers and the only the ones that were discovered.
nc, err := nats.Connect("demo.nats.io",
    nats.DiscoveredServersHandler(func(nc *nats.Conn) {
        log.Printf("Known servers: %v\n", nc.Servers())
        log.Printf("Discovered servers: %v\n", nc.DiscoveredServers())
    }))
if err != nil {
    log.Fatal(err)
}
deffer nc.Close()

// Do something with the connection
```

```Java
class ServersAddedListener implements ConnectionListener {
    public void connectionEvent(Connection nc, Events event) {
        if (event == Events.DISCOVERED_SERVERS) {
            for (String server : nc.getServers()) {
                System.out.println("Known server: "+server);
            }
        }
    }
}

public class ListenForNewServers {
    public static void main(String[] args) {

        try {
            Options options = new Options.Builder().
                                        server("nats://demo.nats.io:4222").
                                        connectionListener(new ServersAddedListener()). // Set the listener
                                        build();
            Connection nc = Nats.connect(options);

            // Do something with the connection

            nc.close();
        } catch (Exception e) {
            e.printStackTrace();
        }
    }
}
```

```JavaScript
const nc = await connect({ servers: ["demo.nats.io:4222"] });
(async () => {
  for await (const s of nc.status()) {
    switch (s.type) {
      case Status.Update:
        t.log(`servers added - ${s.data.added}`);
        t.log(`servers deleted - ${s.data.deleted}`);
        break;
      default:
    }
  }
})().then();
```

```C
static void
discoveredServersCB(natsConnection *conn, void *closure)
{
    natsStatus  s         = NATS_OK;
    char        **servers = NULL;
    int         count     = 0;

    s = natsConnection_GetDiscoveredServers(conn, &servers, &count);
    if (s == NATS_OK)
    {
        int i;

        // Do something...
        for (i=0; i<count; i++)
            printf("Discovered server: %s\n", servers[i]);

        // Free allocated memory
        for (i=0; i<count; i++)
            free(servers[i]);
        free(servers);
    }
}

(...)

natsConnection      *conn      = NULL;
natsOptions         *opts      = NULL;
natsStatus          s          = NATS_OK;

s = natsOptions_Create(&opts);
if (s == NATS_OK)
    s = natsOptions_SetDiscoveredServersCB(opts, discoveredServersCB, NULL);
if (s == NATS_OK)
    s = natsConnection_Connect(&conn, opts);

(...)


// Destroy objects that were created
natsConnection_Destroy(conn);
natsOptions_Destroy(opts);
```

--------------------------------

### NATS Cluster Route Pool Size Mismatch Error

Source: https://github.com/nats-io/nats.docs/blob/master/running-a-nats-service/configuration/clustering/v2_routes.md

Illustrates the error message received when NATS servers in a cluster are configured with different route pool sizes, preventing cluster formation.

```log
[ERR] 127.0.0.1:6222 - rid:6 - Mismatch route pool size: 3 vs 4
```

--------------------------------

### Scale and Expand NATS Cluster

Source: https://github.com/nats-io/nats.docs/blob/master/running-a-nats-service/running/nats_docker/docker_swarm.md

Commands to add additional NATS server nodes to the cluster by routing to the seed node and scaling the number of client service replicas.

```bash
docker service create --network nats-cluster-example --name nats-cluster-node-2 nats:1.0.0 -cluster nats://0.0.0.0:6222 -routes nats://nats-cluster-node-1:6222 -DV
docker service scale ruby-nats=3
docker service ps ruby-nats
```