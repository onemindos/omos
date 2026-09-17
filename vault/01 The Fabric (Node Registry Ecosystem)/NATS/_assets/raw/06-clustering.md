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

### Configure NATS Cluster Route Pool Size

Source: https://github.com/nats-io/nats.docs/blob/master/running-a-nats-service/configuration/clustering/v2_routes.md

Sets the number of route connections between NATS servers. A value greater than 1 enables connection pooling. Setting it to -1 disables pooling, reverting to a single connection. All servers in a cluster must have the same pool size.

```hcl
cluster {
  pool_size: 3
}
```

--------------------------------

### Configure NATS JetStream Cluster Nodes

Source: https://github.com/nats-io/nats.docs/blob/master/running-a-nats-service/configuration/clustering/jetstream_clustering/README.md

Configuration blocks for three NATS servers participating in a JetStream cluster. Each node defines a unique server name, a shared cluster name, and routing information to other cluster members.

```NATS Configuration
server_name=n1-c1
listen=4222

accounts {
  $SYS {
    users = [
      { user: "admin",
        pass: "$2a$11$DRh4C0KNbNnD8K/hb/buWe1zPxEHrLEiDmuq1Mi0rRJiH/W25Qidm"
      }
    ]
  }
}

jetstream {
   store_dir=/nats/storage
}

cluster {
  name: C1
  listen: 0.0.0.0:6222
  routes: [
    nats://host_b:6222
    nats://host_c:6222
  ]
}
```

```NATS Configuration
server_name=n2-c1
listen=4222

accounts {
  $SYS {
    users = [
      { user: "admin",
        pass: "$2a$11$DRh4C0KNbNnD8K/hb/buWe1zPxEHrLEiDmuq1Mi0rRJiH/W25Qidm"
      }
    ]
  }
}

jetstream {
   store_dir=/nats/storage
}

cluster {
  name: C1
  listen: 0.0.0.0:6222
  routes: [
    nats://host_a:6222
    nats://host_c:6222
  ]
}
```

--------------------------------

### NATS Cluster Route Pool Size Mismatch Error

Source: https://github.com/nats-io/nats.docs/blob/master/running-a-nats-service/configuration/clustering/v2_routes.md

Illustrates the error message received when NATS servers in a cluster are configured with different route pool sizes, preventing cluster formation.

```log
[ERR] 127.0.0.1:6222 - rid:6 - Mismatch route pool size: 3 vs 4
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

### NATS Cluster Host and Port Configuration

Source: https://github.com/nats-io/nats.docs/blob/master/running-a-nats-service/configuration/clustering/README.md

This configuration snippet demonstrates how to explicitly define the hostname and port for cluster communication. If the hostname is omitted, it defaults to '0.0.0.0', binding to all available interfaces.

```conf
cluster {
  host: 127.0.0.1
  port: 4248
}
```

--------------------------------

### Verify NATS Cluster Routes with curl

Source: https://github.com/nats-io/nats.docs/blob/master/running-a-nats-service/running/nats_docker/README.md

This snippet shows how to use curl to query the monitoring endpoint of a NATS server to verify that cluster routes are connected. It requests the /routez endpoint and expects to see the number of connected routes.

```bash
curl http://127.0.0.1:8222/routez
```

--------------------------------

### Create NATS Cluster with Docker

Source: https://github.com/nats-io/nats.docs/blob/master/running-a-nats-service/running/nats_docker/README.md

This snippet demonstrates how to create a NATS cluster using Docker. It starts a seed server and then additional servers that connect to the seed server using routes. Credentials for route protection are included as per the Docker image configuration.

```bash
docker run --name nats --network nats --rm -p 4222:4222 -p 8222:8222 nats --http_port 8222 --cluster_name NATS --cluster nats://0.0.0.0:6222
```

```bash
docker run --name nats-1 --network nats --rm nats --cluster_name NATS --cluster nats://0.0.0.0:6222 --routes=nats://ruser:T0pS3cr3t@nats:6222
docker run --name nats-2 --network nats --rm nats --cluster_name NATS --cluster nats://0.0.0.0:6222 --routes=nats://ruser:T0pS3cr3t@nats:6222
```

--------------------------------

### Starting Third NATS Server in Cluster

Source: https://github.com/nats-io/nats.docs/blob/master/running-a-nats-service/configuration/clustering/README.md

This command initiates the third NATS server in the cluster, configuring its client port, cluster address, and routes to connect to the seed server. This completes the multi-server cluster setup.

```bash
nats-server -p 6222 -cluster nats://localhost:6248 -routes nats://localhost:4248 -D
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

### Route Information (`/routez`)

Source: https://github.com/nats-io/nats.docs/blob/master/running-a-nats-service/nats_admin/monitoring/README.md

The `/routez` endpoint provides information about active routes within a NATS cluster. It details the connections between NATS servers, including their IDs, IP addresses, and message transfer statistics.

```APIDOC
## GET /routez

### Description
Retrieves information about active routes (connections between NATS servers) in the cluster.

### Method
GET

### Endpoint
/routez

### Query Parameters
- **subs** (boolean | integer | string) - Optional - Include subscription details. Accepts `true`, `1`, `false`, `0`, or `detail` for more verbose subscription information.

### Response
#### Success Response (200)
- **server_id** (string) - The ID of the NATS server.
- **now** (string) - The current timestamp.
- **num_routes** (integer) - The total number of active routes.
- **routes** (array) - An array of route objects, each containing:
  - **rid** (integer) - The route ID.
  - **remote_id** (string) - The ID of the remote NATS server.
  - **did_solicit** (boolean) - Indicates if this server solicited the connection.
  - **ip** (string) - The IP address of the remote server.
  - **port** (integer) - The port of the remote server.
  - **pending_size** (integer) - The size of pending messages to be sent.
  - **in_msgs** (integer) - The total number of messages received from the remote server.
  - **out_msgs** (integer) - The total number of messages sent to the remote server.
  - **in_bytes** (integer) - The total number of bytes received from the remote server.
  - **out_bytes** (integer) - The total number of bytes sent to the remote server.
  - **subscriptions** (integer) - The number of subscriptions advertised by the remote server.

#### Response Example
```json
{
  "server_id": "NACDVKFBUW4C4XA24OOT6L4MDP56MW76J5RJDFXG7HLABSB46DCMWCOW",
  "now": "2019-06-24T14:29:16.046656-07:00",
  "num_routes": 1,
  "routes": [
    {
      "rid": 1,
      "remote_id": "de475c0041418afc799bccf0fdd61b47",
      "did_solicit": true,
      "ip": "127.0.0.1",
      "port": 61791,
      "pending_size": 0,
      "in_msgs": 0,
      "out_msgs": 0,
      "in_bytes": 0,
      "out_bytes": 0,
      "subscriptions": 0
    }
  ]
}
```
```

--------------------------------

### Initialize NATS Cluster Infrastructure

Source: https://github.com/nats-io/nats.docs/blob/master/running-a-nats-service/running/nats_docker/docker_swarm.md

Commands to create a Docker overlay network and launch the initial seed NATS server. The seed server listens on port 6222 for cluster routing.

```bash
docker network create --driver overlay nats-cluster-example
docker service create --network nats-cluster-example --name nats-cluster-node-1 nats:1.0.0 -cluster nats://0.0.0.0:6222 -DV
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

--------------------------------

### NATS Seed Node Configuration

Source: https://github.com/nats-io/nats.docs/blob/master/running-a-nats-service/configuration/clustering/README.md

This configuration block defines the listening ports for client connections and HTTP monitoring, as well as the cluster communication port for a NATS server. It's essential for bootstrapping a NATS cluster.

```conf
listen: 127.0.0.1:4222
http: 8222

cluster {
  listen: 127.0.0.1:4248
}
```

--------------------------------

### Connecting Additional NATS Servers to Cluster

Source: https://github.com/nats-io/nats.docs/blob/master/running-a-nats-service/configuration/clustering/README.md

This command starts a new NATS server, specifying its client port, cluster address, and routes to connect to the seed server. It's used to expand the NATS cluster.

```bash
nats-server -p 5222 -cluster nats://localhost:5248 -routes nats://localhost:4248 -D
```

--------------------------------

### Enable Compression on NATS Server Routes

Source: https://github.com/nats-io/nats.docs/blob/master/running-a-nats-service/configuration/clustering/v2_routes.md

This configuration snippet enables compression for route connections between NATS servers. The 'accept' mode allows the server to accept the compression mode offered by the peer. By default, compression is off.

```hcl
cluster {
  compression {
    mode: accept
  }
}
```

--------------------------------

### Configure Leaf Node Server

Source: https://github.com/nats-io/nats.docs/blob/master/running-a-nats-service/configuration/leafnodes/README.md

This configuration sets up a leaf node server, specifying the local listening port and defining remote NATS clusters to connect to. It includes options for standard and TLS connections to remotes.

```hcl
listen: "127.0.0.1:4111"
leafnodes {
    remotes = [
        {
          url: "nats://leafuser:s3cr3t@localhost"
        },
    ]
}
```

```hcl
listen: "127.0.0.1:4111"
leafnodes {
    remotes = [
        {
          url: "tls://leafuser:s3cr3t@localhost"
        },
    ]
}
```

--------------------------------

### Start NATS Cluster Nodes via CLI

Source: https://github.com/nats-io/nats.docs/blob/master/running-a-nats-service/configuration/clustering/README.md

Commands to initialize a NATS cluster by starting a seed server and connecting additional nodes using the -cluster and -routes flags.

```bash
nats-server -p 4222 -cluster nats://localhost:4248 --cluster_name test-cluster
nats-server -p 5222 -cluster nats://localhost:5248 -routes nats://localhost:4248 --cluster_name test-cluster
nats-server -p 6222 -cluster nats://localhost:6248 -routes nats://localhost:4248 --cluster_name test-cluster
```

--------------------------------

### Add Stream to Specific Cluster (CLI)

Source: https://github.com/nats-io/nats.docs/blob/master/nats-concepts/jetstream/streams.md

Use the CLI to add a stream and explicitly assign it to a cluster by its name. Ensure servers in the target cluster have their `name` field defined in the server configuration.

```bash
nats stream add --cluster aws-us-east1-c1
```

--------------------------------

### Connect to Multiple NATS Servers

Source: https://github.com/nats-io/nats.docs/blob/master/_examples/connect_multiple.html

Establishes a connection to a NATS cluster by providing an array of server URLs. This ensures the client can attempt connections to multiple endpoints for high availability.

```Go
servers := []string{"nats://127.0.0.1:1222", "nats://127.0.0.1:1223", "nats://127.0.0.1:1224"}
nc, err := nats.Connect(strings.Join(servers, ","))
if err != nil {
	log.Fatal(err)
}
defer nc.Close()
```

```Java
Options options = new Options.Builder().
	server("nats://localhost:1222").
	server("nats://localhost:1223").
	server("nats://localhost:1224").
	build();
Connection nc = Nats.connect(options);
nc.close();
```

```JavaScript
let nc = NATS.connect({
	servers: [
		"nats://demo.nats.io:4222",
		"nats://localhost:4222"
	]
});
nc.on('connect', (c) => {
	nc.close();
});
nc.on('error', (err) => {
	console.error(err);
});
```

```Python
nc = NATS()
await nc.connect(servers=[
	"nats://127.0.0.1:1222",
	"nats://127.0.0.1:1223",
	"nats://127.0.0.1:1224"
])
await nc.close()
```

```Ruby
require 'nats/client'
NATS.start(servers: ["nats://127.0.0.1:1222", "nats://127.0.0.1:1223", "nats://127.0.0.1:1224"]) do |nc|
	nc.close
end
```

```TypeScript
let nc = await connect({
	servers: [
		"nats://demo.nats.io:4222",
		"nats://localhost:4222"
	]
});
nc.close();
```

--------------------------------

### Connect to NATS Cluster in Java

Source: https://github.com/nats-io/nats.docs/blob/master/using-nats/developing-with-nats/connecting/cluster.md

Shows how to establish a connection to a NATS cluster in Java using the NATS client library. It configures the connection with server URLs and includes proper connection closing.

```Java
Options options = new Options.Builder()
    .server("nats://127.0.0.1:1222,nats://127.0.0.1:1223,nats://127.0.0.1:1224")
    .build();
Connection nc = Nats.connect(options);

// Do something with the connection

nc.close();
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