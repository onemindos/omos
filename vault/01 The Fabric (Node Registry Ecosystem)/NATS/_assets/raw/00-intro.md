### Start NATS Server

Source: https://github.com/nats-io/nats.docs/blob/master/running-a-nats-service/configuration/leafnodes/README.md

Command to start the NATS server using a specified configuration file.

```bash
nats-server -c /tmp/server.conf
```

--------------------------------

### List Consumers for a Stream with NATS CLI

Source: https://github.com/nats-io/nats.docs/blob/master/running-a-nats-service/nats_admin/jetstream_admin/consumers.md

Lists all consumers associated with a specific NATS stream, in this case, the 'ORDERS' stream. This command provides a quick overview of the active consumers.

```shell
nats con ls ORDERS
```

--------------------------------

### Define NATS Message Structure

Source: https://github.com/nats-io/nats.docs/blob/master/reference/nats-protocol/nats-protocol/nats-client-dev.md

Defines the standard message structure used by NATS clients to handle incoming data, including subject, reply, and payload.

```go
type Msg struct {
    Subject string
    Reply   string
    Data    []byte
    Sub     *Subscription
}
```

--------------------------------

### Connect to Default NATS Server (Go)

Source: https://github.com/nats-io/nats.docs/blob/master/using-nats/developing-with-nats/connecting/default_server.md

Connects to the default NATS server URL (nats://localhost:4222) using the Go NATS client library. Handles potential connection errors and ensures the connection is closed when done.

```Go
nc, err := nats.Connect(nats.DefaultURL)
if err != nil {
    log.Fatal(err)
}
defer nc.Close()

// Do something with the connection
```

--------------------------------

### NATS CLI Connection with Creds File

Source: https://github.com/nats-io/nats.docs/blob/master/running-a-nats-service/nats_admin/jwt.md

Command to connect to a NATS server using the NATS CLI and a user credentials file.

```bash
> nats -s localhost:4222 "--creds=user.creds" pub "foo" "hello world"
```

--------------------------------

### Start NATS Server

Source: https://github.com/nats-io/nats.docs/blob/master/howto-examples/overview.md

Commands to launch the NATS server. The first command starts the server with default settings, while the second enables debug and trace logging for development purposes.

```shell
nats-server
```

```shell
nats-server -DV
```

--------------------------------

### Connect to NATS Server

Source: https://github.com/nats-io/nats.docs/blob/master/_examples/connect_url.html

Examples show how to establish a connection to a NATS server using the provided URL. It covers different client libraries and includes basic connection setup and closing.

```Go
// If connecting to the default port, the URL can be simplified // to just the hostname/IP. // That is, the connect below is equivalent to: // nats.Connect("nats://demo.nats.io:4222") nc, err := nats.Connect("demo.nats.io") if err != nil { 
	log.Fatal(err) }
defer nc.Close() 
// Do something with the connection
```

```Java
Connection nc = Nats.connect("nats://demo.nats.io:4222"); 
// Do something with the connection 

nc.close();
```

```JavaScript
let nc = NATS.connect("nats://demo.nats.io:4222"); nc.on('connect', (c) => { 
    // Do something with the connection 
    doSomething(); 
    // When done close it 
    nc.close() 
}); nc.on('error', (err) => { 
    failed(err); 
});
```

```Python
nc = NATS() 
await nc.connect(servers=["nats://demo.nats.io:4222"])

# Do something with the connection 

await nc.close()
```

```Ruby
require 'nats/client' 

NATS.start(servers: ["nats://demo.nats.io:4222"]) do |nc| 
   # Do something with the connection 
   # Close the connection 
   nc.close 
end
```

```TypeScript
// will throw an exception if connection fails 
    let nc = await connect("nats://demo.nats.io:4222"); 
    // Do something with the connection 

    // Close the connection 
    nc.close();
```

--------------------------------

### Install NATS Helm Chart

Source: https://github.com/nats-io/nats.docs/blob/master/running-a-nats-service/nats-on-kubernetes/nats-kubernetes.md

Install the NATS chart from the added repository. This command deploys NATS with default configuration values.

```sh
helm install nats nats/nats
```

--------------------------------

### NATS Server Configuration with Subject Mappings

Source: https://github.com/nats-io/nats.docs/blob/master/running-a-nats-service/configuration/README.md

Example NATS server configuration demonstrating global and account-level subject mappings. Useful for setting up aliases and pattern-based translations for subjects.

```text
host: 0.0.0.0
port:4222

mappings: {
	foo: bar
}

accounts: {
    accountA: {
	mappings: {
	    orders.acme.*: orders.$1
	}
        users: [
            {user: admin, password: admin},
            {user: user, password: user}
           ]
    },
}
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

### Connect to NATS with Ruby

Source: https://github.com/nats-io/nats.docs/blob/master/using-nats/developing-with-nats/connecting/security/userpass.md

Connects to a NATS server using the Ruby NATS client. It includes setup for connection, error handling, and reconnection events.

```ruby
require 'nats/client'

NATS.start(servers:["nats://myname:password@127.0.0.1:4222"], name: "my-connection") do |nc|
   nc.on_error do |e|
    puts "Error: #{e}"
  end

   nc.on_reconnect do
    puts "Got reconnected to #{nc.connected_server}"
  end

  nc.on_disconnect do |reason|
    puts "Got disconnected! #{reason}"
  end

  nc.close
end
```

--------------------------------

### Initiate NATS Connection with Telnet

Source: https://github.com/nats-io/nats.docs/blob/master/reference/nats-protocol/nats-protocol-demo.md

Establishes a connection to a NATS server using Telnet. This is the first step to interact with NATS using the protocol directly. It shows the command to connect and the expected server response including connection details.

```shell
telnet demo.nats.io 4222
```

```text
Trying 107.170.221.32...
Connected to demo.nats.io.
Escape character is '^]'.
INFO {"server_id":"NCXMJZYQEWUDJFLYLSTTE745I2WUNCVG3LJJ3NRKSFJXEG6RGK7753DJ","version":"2.0.0","proto":1,"go":"go1.11.10","host":"0.0.0.0","port":4222,"max_payload":1048576,"client_id":5089}
```

--------------------------------

### Accessing NATS CLI Help

Source: https://github.com/nats-io/nats.docs/blob/master/using-nats/nats-tools/nats_cli/README.md

Displays the help documentation for the NATS CLI tool to understand available commands and global flags.

```shell
nats --help
```

--------------------------------

### Connect to Default NATS Server

Source: https://github.com/nats-io/nats.docs/blob/master/using-nats/developing-with-nats/connecting/default_server.md

Examples demonstrating how to connect to the default NATS server URL (nats://localhost:4222) across multiple programming languages.

```APIDOC
## Connect to Default NATS Server

### Description
This section shows how to connect to the default NATS server, typically located at `nats://localhost:4222`, using various client libraries.

### Method
N/A (Client-side connection)

### Endpoint
N/A (Client-side connection to `nats://localhost:4222`)

### Parameters
N/A

### Request Example
N/A

### Response
N/A

## Language-Specific Connection Examples

### Go
```go
nc, err := nats.Connect(nats.DefaultURL)
if err != nil {
    log.Fatal(err)
}
defer nc.Close()

// Do something with the connection
```

### Java
```java
Connection nc = Nats.connect();

// Do something with the connection

nc.close();
```

### JavaScript
```javascript
const nc = await connect();
// Do something with the connection
doSomething();
// When done close it
await nc.close();
```

### Python
```python
nc = NATS()
await nc.connect()

# Do something with the connection

await nc.close()
```

### C#
```csharp
// dotnet add package NATS.Net
using NATS.Net;

await using var client = new NatsClient();

// It's optional to call ConnectAsync() 
as it will be called when needed automatically
await client.ConnectAsync();
```

### Ruby
```ruby
require 'nats/client'

NATS.start do |nc|
   # Do something with the connection

   # Close the connection
   nc.close
end
```

### C
```c
natsConnection      *conn = NULL;
natsStatus          s;

s = natsConnection_ConnectTo(&conn, NATS_DEFAULT_URL);
if (s != NATS_OK)
  // handle error

// Destroy connection, no-op if conn is NULL.
natsConnection_Destroy(conn);
```
```

--------------------------------

### Manage NATS Server with Docker

Source: https://github.com/nats-io/nats.docs/blob/master/running-a-nats-service/installation.md

Commands to pull the latest NATS image and run the server container. It exposes the default NATS port 4222.

```shell
docker pull nats:latest
docker run -p 4222:4222 -ti nats:latest
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

### Define NATS Server Configuration

Source: https://github.com/nats-io/nats.docs/blob/master/using-nats/nats-tools/nsc/basics.md

Example configuration file for a local NATS server instance.

```text
server_name: servertest
listen: 127.0.0.1:4222
http: 8222

jetstream: enabled

include resolver.conf
```

--------------------------------

### Connect to Default NATS Server

Source: https://github.com/nats-io/nats.docs/blob/master/_examples/connect_default.html

Demonstrates how to connect to the default NATS server URL. This is a basic connection setup. Ensure the NATS server is running.

```Go
nc, err := nats.Connect(nats.DefaultURL)
if err != nil {
	log.Fatal(err)
}
defer nc.Close()

// Do something with the connection
```

```Java
Connection nc = Nats.connect();

// Do something with the connection

nc.close();
```

```JavaScript
let nc = NATS.connect();
nc.on('connect', (c) => {
    // Do something with the connection
    doSomething();
    // When done close it
    nc.close();
});
nc.on('error', (err) => {
    failed(err);
});
```

```Python
nc = NATS()
await nc.connect()

// Do something with the connection

await nc.close()
```

```Ruby
require 'nats/client'

NATS.start do |nc|
   # Do something with the connection
    # Close the connection
    nc.close
end
```

```TypeScript
// will throw an exception if connection fails
let nc = await connect();
// Do something with the connection

// When done close it
nc.close();

// alternatively, you can use the Promise pattern
let nc1: Client;
connect()
    .then((c) => {
        nc1 = c;
        // Do something with the connection
        nc1.close();
    });
    // add a .catch/.finally
```

--------------------------------

### Install and Start NATS Windows Service

Source: https://github.com/nats-io/nats.docs/blob/master/running-a-nats-service/running/windows_srv.md

Commands to create and start the NATS server as a Windows service using `sc.exe`. Ensure to replace `[nats-server flags]` with your desired server configuration flags.

```shell
sc.exe create nats-server binPath= "%NATS_PATH%\nats-server.exe [nats-server flags]"
sc.exe start nats-server
```

--------------------------------

### Connect to Default NATS Server (Java)

Source: https://github.com/nats-io/nats.docs/blob/master/using-nats/developing-with-nats/connecting/default_server.md

Establishes a connection to the default NATS server using the Java NATS client. The connection is explicitly closed after use.

```Java
Connection nc = Nats.connect();

// Do something with the connection

nc.close();
```

--------------------------------

### Watch NATS Object Bucket Changes

Source: https://github.com/nats-io/nats.docs/blob/master/nats-concepts/jetstream/object-store/obj_walkthrough.md

Monitors a specified NATS object storage bucket for changes, such as PUT or DELETE operations. This command streams real-time updates about objects within the bucket. It requires the NATS CLI to be installed and configured to connect to a NATS server.

```shell
nats object watch myobjbucket
```

--------------------------------

### Start NATS Server with Included Configuration

Source: https://github.com/nats-io/nats.docs/blob/master/running-a-nats-service/configuration/README.md

Command to start the NATS server using a main configuration file (`server.conf`) that includes other configuration files.

```shell
nats-server -c server.conf
```

--------------------------------

### Connect to NATS with C

Source: https://github.com/nats-io/nats.docs/blob/master/using-nats/developing-with-nats/connecting/security/userpass.md

Demonstrates connecting to a NATS server using the C NATS client library. It covers creating options, setting the URL, and establishing the connection.

```c
natsConnection      *conn      = NULL;
natsOptions         *opts      = NULL;
natsStatus          s          = NATS_OK;

s = natsOptions_Create(&opts);
if (s == NATS_OK)
    s = natsOptions_SetURL(opts, "nats://myname:password@127.0.0.1:4222");
if (s == NATS_OK)
    s = natsConnection_Connect(&conn, opts);

(...)

// Destroy objects that were created
natsConnection_Destroy(conn);
natsOptions_Destroy(opts);
```