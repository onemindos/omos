### Configure NATS JetStream Leaf Node Connections

Source: https://github.com/nats-io/nats.docs/blob/master/running-a-nats-service/configuration/leafnodes/jetstream_leafnodes.md

Configuration files for setting up a hub and leaf node architecture with JetStream enabled. These examples demonstrate how to define domains to isolate JetStream namespaces and configure leaf node remotes for account synchronization.

```text
accounts {
    SYS: {
        users: [{user: admin, password: admin}]
    },
    ACC: {
        users: [{user: acc, password: acc}],
        jetstream: enabled
    }
}
system_account: SYS
```

```text
port: 4222
server_name: hub-server
jetstream {
    store_dir="./store_server"
    domain=hub
}
leafnodes {
    port: 7422
}
include ./accounts.conf
```

```text
port: 4111
server_name: leaf-server
jetstream {
    store_dir="./store_leaf"
    domain=leaf
}
leafnodes {
    remotes = [
        {
            urls: ["nats://admin:admin@0.0.0.0:7422"]
            account: "SYS"
        },
        {
            urls: ["nats://acc:acc@0.0.0.0:7422"]
            account: "ACC"
        }
    ]
}
include ./accounts.conf
```

--------------------------------

### NATS Server Leafnode Configuration

Source: https://github.com/nats-io/nats.docs/blob/master/running-a-nats-service/running/nats_docker/ngs-leafnodes-docker.md

A minimal NATS server configuration file for leaf nodes connecting to Synadia Cloud (NGS). It specifies the remote server URL and the path to the credentials file.

```NATS Config
leafnodes {
    remotes = [
        {
          url: "tls://connect.ngs.global"
          credentials: "ngs.creds"
        },
    ]
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

### Configure NATS Server for Leaf Node Connections

Source: https://github.com/nats-io/nats.docs/blob/master/running-a-nats-service/configuration/leafnodes/README.md

This configuration snippet sets up a standard NATS server to accept leaf node connections on a specified port (default 7422). It also defines accounts and users for authentication.

```hcl
leafnodes {
    port: 7422
}
accounts: {
    app: { 
        users: [
            {user: appuser, password: s3cr3t},
            {user: leafuser, password: s3cr3t}
           ]
    }
}
```

--------------------------------

### Configure Leaf Node Remotes in Non-Operator Mode

Source: https://github.com/nats-io/nats.docs/blob/master/running-a-nats-service/nats_admin/jwt.md

This configuration example demonstrates how to set up leaf node remotes in a non-operator mode environment. It shows the explicit definition of both a user account and the system account ($SYS) to ensure proper connectivity.

```hcl
leafnodes {
    remotes = [
        {
          url: "nats://localhost:4222"
          credentials: "./your-account.creds"
        },
        {
          url: "nats://localhost:4222"
          account: "$SYS"
          credentials: "./system-account.creds"
        },
    ]
}
```

--------------------------------

### Start Leaf Node Server

Source: https://github.com/nats-io/nats.docs/blob/master/running-a-nats-service/configuration/leafnodes/README.md

Command to start the NATS leaf node server using a specified configuration file.

```bash
nats-server -c /tmp/leaf.conf
```

--------------------------------

### NATS Leafnode Singleton Authorization Configuration

Source: https://github.com/nats-io/nats.docs/blob/master/running-a-nats-service/configuration/leafnodes/leafnode_conf.md

Configures a single user for leafnode authentication. The leafnode connects using the provided username and password, and is bound to the specified account. This configuration is useful for simple deployments.

```hcl
leafnodes {
  port: 6222
  authorization {
    user: leaf
    password: secret
    account: TheAccount
  }
}
```

--------------------------------

### LeafNode `remotes` Entry Block Configuration

Source: https://github.com/nats-io/nats.docs/blob/master/running-a-nats-service/configuration/leafnodes/leafnode_conf.md

Defines the properties for configuring remote LeafNode connections.

```APIDOC
## LeafNode `remotes` Entry Block

This block configures remote LeafNode connections, specifying how a NATS server should connect to and interact with a remote LeafNode.

### Properties

- **`url`** (string) - The URL of the Leafnode. The protocol must be `nats-leaf`.
- **`urls`** (array of strings) - An array of Leafnode URLs for discovery. Supports multiple URLs for connection attempts.
- **`no_randomize`** (boolean) - If `true`, connection attempts will follow the order of URLs provided. If `false` (default), URLs are shuffled for connection attempts.
- **`account`** (string) - The name or JWT public key of the local account to bind to this remote server. Traffic on this account will be forwarded to the remote server.
- **`deny_imports`** (array of strings) - A list of subjects that will not be imported from the remote server. Subscriptions to these subjects will not be propagated to the hub.
- **`deny_exports`** (array of strings) - A list of subjects that will not be exported to the remote server. Subscriptions to these subjects will not be propagated into the leaf node.
- **`credentials`** (string) - Path to the credential file for connecting to the leafnode server.
- **`nkey`** (string) - The Nkey used for connecting to the leafnode server.
- **`tls`** (object) - A TLS configuration block. Specifies certificates for connecting and authenticating.
- **`ws_compression`** (boolean) - If using the WebSocket protocol, this indicates a desire to use compression. Defaults to `false`.
- **`ws_no_masking`** (boolean) - If using the WebSocket protocol, this indicates a desire not to mask outbound WebSocket frames. Defaults to `false`.
- **`compression`** (string) - Configures compression for leafnode connections. Defaults to `s2_auto`. See [here](../clustering/v2_routes.md#compression) for details.
- **`hub`** (boolean) - If `true`, reverses the roles of leaf node and hub, allowing the hub to initiate the connection.
- **`first_info_timeout`** (duration) - The maximum time to wait for initial server metadata. Defaults to `1s`. Useful for detecting non-NATS ports.
```

--------------------------------

### Run NATS Leafnode in Docker (Red User)

Source: https://github.com/nats-io/nats.docs/blob/master/running-a-nats-service/running/nats_docker/ngs-leafnodes-docker.md

Launches a NATS server as a Docker container configured to act as a leaf node. It maps ports, mounts configuration and credential files, and installs necessary SSL certificates for secure connection to NGS.

```Shell
docker run  -p 4222:4222 -v leafnode.conf:/leafnode.conf -v /etc/ssl/cert.pem:/etc/ssl/cert.pem -v default-red.creds:/ngs.creds  nats:latest -c /leafnode.conf
```

--------------------------------

### Configure NATS Leafnodes in Operator Mode

Source: https://github.com/nats-io/nats.docs/blob/master/running-a-nats-service/nats_admin/jwt.md

This configuration snippet demonstrates how to define operator and system account settings for NATS leafnodes, including remote connection details and credentials for specific accounts.

```text
operator: ./trustedOperator.jwt
system_account: AAAXAUVSGK7TCRHFIRAS4SYXVJ76EWDMNXZM6ARFGXP7BASNDGLKU7A5
leafnodes {
    remotes = [
        {
          url: "nats://localhost:4222"
          account: "ADKGAJU55CHYOIF5H432K2Z2ME3NPSJ5S3VY5Q42Q3OTYOCYRRG7WOWV"
          credentials: "./your-account.creds"
        },
        {
          url: "nats://localhost:4222"
          account: "AAAXAUVSGK7TCRHFIRAS4SYXVJ76EWDMNXZM6ARFGXP7BASNDGLKU7A5"
          credentials: "./system-account.creds"
        },
    ]
}
```

--------------------------------

### NATS Leafnode Multi-User Authorization Configuration

Source: https://github.com/nats-io/nats.docs/blob/master/running-a-nats-service/configuration/leafnodes/leafnode_conf.md

Configures multiple users for leafnode authentication, allowing different leafnodes to connect with distinct credentials and be bound to different accounts. This provides more granular control in complex environments.

```hcl
leafnodes {
  port: 6222
  authorization {
    users = [
      {user: leaf1, password: secret, account: account1}
      {user: leaf2, password: secret, account: account2}
    ]
  }
}
```

--------------------------------

### Start NATS Service on Leafnode (Red)

Source: https://github.com/nats-io/nats.docs/blob/master/running-a-nats-service/running/nats_docker/ngs-leafnodes-docker.md

Starts a NATS service on the 'red' leaf node, configured to listen on localhost:4222. This service acts as a replier, responding to requests with a formatted message including the time and the received request.

```Shell
nats -s localhost:4222 reply docker-leaf-test "At {{Time}}, I received your request: {{Request}}"
```

--------------------------------

### Define Leaf Node Remote Connection

Source: https://github.com/nats-io/nats.docs/blob/master/running-a-nats-service/configuration/leafnodes/README.md

Configuration snippet for the NATS server to connect to a remote NGS global service using TLS and user credentials.

```hcl
leafnodes {
    remotes = [
        {
          url: "tls://connect.ngs.global"
          credentials: "/Users/alberto/.nkeys/creds/synadia/leaftest/leaftestuser.creds"
        },
    ]
}
```

--------------------------------

### Configure Leafnode TLS Handshake First (Accepting Side)

Source: https://github.com/nats-io/nats.docs/blob/master/running-a-nats-service/configuration/leafnodes/leafnode_conf.md

This configuration enables a TLS-first handshake for incoming leafnode connections on the accepting server. It requires that remote leafnodes also have the `handshake_first` setting enabled for a successful connection. If not, the connection will fail due to protocol mismatch.

```go
leafnodes {
  port: 7422
  tls {
    handshake_first: true
    # other TLS fields...
  }
}
```

--------------------------------

### Run NATS Leafnode in Docker (Blue User)

Source: https://github.com/nats-io/nats.docs/blob/master/running-a-nats-service/running/nats_docker/ngs-leafnodes-docker.md

Starts a second NATS leaf node in a Docker container, similar to the first but using different credentials and a different host port to avoid conflicts. This allows for multiple isolated leaf node connections.

```Shell
docker run  -p 4333:4222 -v leafnode.conf:/leafnode.conf -v /etc/ssl/cert.pem:/etc/ssl/cert.pem -v default-blue.creds:/ngs.creds  nats:latest -c /leafnode.conf
```

--------------------------------

### Send NATS Request via Leafnode (Blue)

Source: https://github.com/nats-io/nats.docs/blob/master/running-a-nats-service/running/nats_docker/ngs-leafnodes-docker.md

Sends a request to the 'docker-leaf-test' subject via the 'blue' leaf node, connected to localhost:4333. It expects a reply from the service running on the 'red' leaf node.

```Shell
$ nats -s localhost:4333 request docker-leaf-test "Hello World"
```

--------------------------------

### LeafNode Remote Entry Configuration

Source: https://github.com/nats-io/nats.docs/blob/master/running-a-nats-service/configuration/leafnodes/leafnode_conf.md

Defines the properties for configuring a remote LeafNode connection, including URLs, account information, denial lists for imports/exports, credentials, TLS settings, and WebSocket options.

```go
type RemoteLeafOpts struct {
	Url              string
	Urls             []string
	NoRandomize      bool
	Account          string
	DenyImports      []string
	DenyExports      []string
	Credentials      string
	Nkey             string
	Tls              *TlsConfig
	WsCompression    bool
	WsNoMasking      bool
	Compression      string
	Hub              bool
	FirstInfoTimeout time.Duration
}
```

--------------------------------

### Leaf Node Information

Source: https://github.com/nats-io/nats.docs/blob/master/running-a-nats-service/nats_admin/monitoring/README.md

The /leafz endpoint reports detailed information about the leaf node connections. It supports an optional 'subs' argument to include internal subscriptions.

```APIDOC
## GET /leafz

### Description
Retrieves detailed information about leaf node connections.

### Method
GET

### Endpoint
/leafz

### Parameters
#### Query Parameters
- **subs** (boolean) - Optional - Include internal subscriptions. Defaults to false.

### Request Example
GET https://demo.nats.io:8222/leafz?subs=1

### Response
#### Success Response (200)
- **server_id** (string) - The ID of the NATS server.
- **now** (string) - The current timestamp.
- **leafnodes** (integer) - The total number of leaf nodes.
- **leafs** (array) - An array of leaf node objects.
  - **account** (string) - The account associated with the leaf node.
  - **ip** (string) - The IP address of the leaf node.
  - **port** (integer) - The port of the leaf node.
  - **rtt** (string) - The round-trip time to the leaf node.
  - **in_msgs** (integer) - The number of messages received by the leaf node.
  - **out_msgs** (integer) - The number of messages sent by the leaf node.
  - **in_bytes** (integer) - The number of bytes received by the leaf node.
  - **out_bytes** (integer) - The number of bytes sent by the leaf node.
  - **subscriptions** (integer) - The number of subscriptions the leaf node has.
  - **subscriptions_list** (array) - A list of subjects the leaf node is subscribed to.

#### Response Example
{
  "server_id": "NC2FJCRMPBE5RI5OSRN7TKUCWQONCKNXHKJXCJIDVSAZ6727M7MQFVT3",
  "now": "2019-08-27T09:07:05.841132-06:00",
  "leafnodes": 1,
  "leafs": [
    {
      "account": "$G",
      "ip": "127.0.0.1",
      "port": 6223,
      "rtt": "200µs",
      "in_msgs": 0,
      "out_msgs": 10000,
      "in_bytes": 0,
      "out_bytes": 1280000,
      "subscriptions": 1,
      "subscriptions_list": ["foo"]
    }
  ]
}
```

--------------------------------

### Configure Leafnode TLS Handshake First (Remote Side)

Source: https://github.com/nats-io/nats.docs/blob/master/running-a-nats-service/configuration/leafnodes/leafnode_conf.md

This configuration specifies that a leafnode client should initiate a TLS-first handshake when connecting to a remote server. If the remote server does not also have `handshake_first` enabled in its configuration, the connection will fail because the client will attempt a TLS handshake while the server expects the NATS INFO protocol.

```go
leafnodes {
  remotes [
    {
      urls: ["tls://example:7422"]
      tls: {
         handshake_first: true
         # other TLS fields...
      }
    }
  ]
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

### Start NATS Leaf Server

Source: https://github.com/nats-io/nats.docs/blob/master/running-a-nats-service/configuration/leafnodes/README.md

Command to launch the NATS server using a specific configuration file to initiate the leaf node connection.

```bash
nats-server -c /tmp/ngs_leaf.conf
```

--------------------------------

### NATS Leafnode Account Binding Without Credentials

Source: https://github.com/nats-io/nats.docs/blob/master/running-a-nats-service/configuration/leafnodes/leafnode_conf.md

Configures leafnode connections to be bound to a specific account without requiring explicit username and password authentication. This is useful when other authentication methods are in place or when anonymous access to a specific account is desired.

```hcl
leafnodes {
  port: 6222
  authorization {
    account: TheAccount
  }
}
```

--------------------------------

### Configure NATS Account and User with NSC

Source: https://github.com/nats-io/nats.docs/blob/master/running-a-nats-service/configuration/leafnodes/README.md

Commands to set the active account and create a new user for leaf node authentication using the nsc CLI tool.

```bash
nsc env -a DELETE_ME
nsc describe account
nsc add user leaftestuser
```