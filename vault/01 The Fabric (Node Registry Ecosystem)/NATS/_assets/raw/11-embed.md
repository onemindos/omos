### Run NATS Client Application

Source: https://github.com/nats-io/nats.docs/blob/master/running-a-nats-service/configuration/leafnodes/jetstream_leafnodes.md

Executes a Go application that connects to a NATS server using the provided URL. This application is designed to subscribe to messages and process them.

```shell
./main nats://import_client:import_client@localhost:4111
```

--------------------------------

### Install NATS Server

Source: https://github.com/nats-io/nats.docs/blob/master/using-nats/nats-tools/nsc/basics.md

Use Go to install the NATS server binary.

```shell
go get github.com/nats-io/nats-server
```

--------------------------------

### Connect to NATS with TLS using Go

Source: https://github.com/nats-io/nats.docs/blob/master/_examples/connect_tls.html

Establishes a TLS connection to a NATS server using client certificates and a root CA. Requires specifying the server address and paths to certificate files.

```Go
nc, err := nats.Connect("localhost", 	nats.ClientCert("resources/certs/cert.pem", "resources/certs/key.pem"), 	nats.RootCAs("resources/certs/ca.pem")) if err != nil { 	log.Fatal(err) } defer nc.Close()  // Do something with the connection
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

### Connect to NATS with Token in URL

Source: https://github.com/nats-io/nats.docs/blob/master/_examples/connect_token_url.html

Demonstrates how to connect to a NATS server using a token embedded directly within the connection URL. This method is useful for simple authentication scenarios where the token can be safely included in the URL.

```Go
// Token in URL
nc, err := nats.Connect("mytoken@localhost")
if err != nil {
	log.Fatal(err)
}
defer nc.Close()

// Do something with the connection
```

```Java
Connection nc = Nats.connect("nats://mytoken@localhost:4222");//Token in URL

// Do something with the connection

nc.close();
```

```JavaScript
let url = `nats://mytoken!@127.0.0.1:${port}`;
let nc = NATS.connect({url: url});
```

```Python
nc = NATS()

await nc.connect(servers=["nats://mytoken@demo.nats.io:4222"])

# Do something with the connection.
```

```Ruby
NATS.start("deadbeef@127.0.0.1:4222") do |nc|
   puts "Connected using token!"
end
```

```TypeScript
let url = `nats://:mytoken!@127.0.0.1:${port}`;
let nc = await connect({url: url});
```

--------------------------------

### Connect to NATS using User/Password in URL

Source: https://github.com/nats-io/nats.docs/blob/master/using-nats/developing-with-nats/connecting/security/userpass.md

This method connects to a NATS server by embedding user credentials directly within the server URL. This format is supported by many NATS client libraries.

```Go
// Set a user and plain text password
nc, err := nats.Connect("myname:password@127.0.0.1")
if err != nil {
    log.Fatal(err)
}
deffer nc.Close()

// Do something with the connection
```

```Java
Connection nc = Nats.connect("nats://myname:password@localhost:4222");

// Do something with the connection

nc.close();
```

```JavaScript
// JavaScript clients don't support username/password in urls use `user` and `pass` options.
```

```Python
nc = NATS()

await nc.connect(servers=["nats://myname:password@demo.nats.io:4222"])
```

--------------------------------

### Connect to NATS with Options (Go)

Source: https://github.com/nats-io/nats.docs/blob/master/_examples/connect_options.html

Connects to the NATS server with a specified name and connection timeout. It ensures proper error handling and defers closing the connection.

```Go
nc, err := nats.Connect("demo.nats.io", nats.Name("API Options Example"), nats.Timeout(10*time.Second))
if err != nil {
	log.Fatal(err)
}
defers nc.Close()
```

--------------------------------

### Build NATS Server from Source

Source: https://github.com/nats-io/nats.docs/blob/master/using-nats/jetstream/getting_started/using_source.md

Commands to clone the NATS server repository, switch to the appropriate branch, and compile the binary using Go.

```shell
git clone https://github.com/nats-io/nats-server.git
cd nats-server
git checkout master
go build
ls -l nats-server
```

--------------------------------

### Install NATS CLI and Server Tools

Source: https://github.com/nats-io/nats.docs/blob/master/running-a-nats-service/nats_admin/jwt.md

Commands to install the NATS server, CLI, nkeys, and NSC tools using Go modules.

```shell
export GO111MODULE=on
go install github.com/nats-io/nats-server/v2@latest
go install github.com/nats-io/natscli/nats@latest
go install github.com/nats-io/nkeys/nk@latest
go install github.com/nats-io/nsc/v2@latest
```

--------------------------------

### Check NATS Connection Status - Go

Source: https://github.com/nats-io/nats.docs/blob/master/_examples/connect_status.html

Demonstrates how to connect to a NATS server using the Go client, check its connection status, and then close the connection. It uses the `nats.Connect` function and the `nc.Status()` method to retrieve the connection state.

```Go
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

### Install NATS Server Binary from Source

Source: https://github.com/nats-io/nats.docs/blob/master/running-a-nats-service/installation.md

Uses the Go toolchain to fetch and install the latest development build of the NATS server from the main branch.

```shell
go install github.com/nats-io/nats-server/v2@main
```

--------------------------------

### Disable NATS Reconnection (Go)

Source: https://github.com/nats-io/nats.docs/blob/master/_examples/reconnect_none.html

Connects to a NATS server without enabling automatic reconnection attempts. This is useful for scenarios where explicit control over reconnections is required. The `nats.NoReconnect()` option is used.

```Go
// Disable reconnect attempts nc, err := nats.Connect("demo.nats.io", nats.NoReconnect()) if err != nil { 
	log.Fatal(err) }
defers nc.Close()  // Do something with the connection
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

### Gracefully Drain NATS Connection - Go

Source: https://github.com/nats-io/nats.docs/blob/master/_examples/drain_conn.html

This Go example demonstrates how to drain a NATS connection. It subscribes to a subject, publishes a message, and then drains the connection, ensuring all messages are processed before closing. It includes error handling and waits for the connection to be closed.

```Go
wg := sync.WaitGroup{}
wg.Add(1)
errCh := make(chan error, 1)

// To simulate a timeout, you would set the DrainTimeout()
// to a value less than the time spent in the message callback,
// so say: nats.DrainTimeout(10*time.Millisecond).

nc, err := nats.Connect("demo.nats.io", 
	nats.DrainTimeout(10*time.Second),
	nats.ErrorHandler(func(_ *nats.Conn, _ *nats.Subscription, err error) {
		errCh <- err
	}),
	nats.ClosedHandler(func(_ *nats.Conn) {
		wg.Done()
	}))
if err != nil {
	log.Fatal(err)
}

// Just to not collide using the demo server with other users.
subject := nats.NewInbox()

// Subscribe, but add some delay while processing.
if _, err := nc.Subscribe(subject, func(_ *nats.Msg) {
	time.Sleep(200 * time.Millisecond)
}); err != nil {
	log.Fatal(err)
}

// Publish a message
if err := nc.Publish(subject, []byte("hello")); err != nil {
	log.Fatal(err)
}

// Drain the connection, which will close it when done.
if err := nc.Drain(); err != nil {
	log.Fatal(err)
}

// Wait for the connection to be closed.
wg.Wait()

// Check if there was an error
select {
case e := <-errCh:
	log.Fatal(e)
default:
}
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

### NATS Request-Reply Implementation in Go

Source: https://github.com/nats-io/nats.docs/blob/master/running-a-nats-service/configuration/securing_nats/jwt/jwt_nkey_auth.md

This snippet demonstrates a NATS responder that listens for messages on a subject and replies with a payload, and a requester that sends messages using user credentials for authentication. It utilizes the nats.go library to manage connections and message exchange.

```go
nc.Subscribe("test", func(m *nats.Msg){
    log.Printf("[Received] %q, replying... \n", string(m.Data))
    m.Respond([]byte("pong from nkeys based server"))
})
```

```go
nc, err := nats.Connect("127.0.0.1:4223", nats.UserCredentials("path/to/user.creds"))
if err != nil {
    log.Fatal(err)
}

for range time.NewTicker(1 * time.Second).C {
    resp, err := nc.Request("test", []byte("test"), 1*time.Second)
    if err != nil {
        log.Println("[Error]", err)
        continue
    }
    log.Println("[Received]", string(resp.Data))
}
```

--------------------------------

### Handle Server Discovery Events in NATS Clients

Source: https://github.com/nats-io/nats.docs/blob/master/_examples/servers_added.html

This snippet demonstrates how to register handlers or listeners to be notified when new NATS servers are discovered in a cluster. It covers Go, Java, JavaScript, and TypeScript implementations.

```go
nc, err := nats.Connect("demo.nats.io", nats.DiscoveredServersHandler(func(nc *nats.Conn) {
	log.Printf("Known servers: %v\n", nc.Servers())
	log.Printf("Discovered servers: %v\n", nc.DiscoveredServers())
}))
```

```java
class ServersAddedListener implements ConnectionListener {
    public void connectionEvent(Connection nc, Events event) {
        if (event == Events.DISCOVERED_SERVERS) {
            for (String server : nc.getServers()) {
                System.out.println("Known server: "+server);
            }
        }
    }
}
Options options = new Options.Builder().server("nats://demo.nats.io:4222").connectionListener(new ServersAddedListener()).build();
```

```javascript
let nc = NATS.connect("nats://demo.nats.io:4222");
nc.on('serversDiscovered', (urls) => {
    t.log('serversDiscovered', urls);
});
```

```typescript
nc.on('serversChanged', (ce) => {
    t.log('servers changed\n', 'added: ',ce.added, 'removed:', ce.removed);
});
```

--------------------------------

### NATS Request-Reply with Timeout (Go)

Source: https://github.com/nats-io/nats.docs/blob/master/using-nats/developing-with-nats/sending/request_reply.md

Demonstrates sending a request message with a specified timeout and handling the response using the Go NATS client library. It connects to a NATS server, sends a 'time' request, and logs the received reply.

```Go
nc, err := nats.Connect("demo.nats.io")
if err != nil {
    log.Fatal(err)
}
deferr nc.Close()

// Send the request
msg, err := nc.Request("time", nil, time.Second)
if err != nil {
    log.Fatal(err)
}

// Use the response
log.Printf("Reply: %s", msg.Data)

// Close the connection
nc.Close()
```

--------------------------------

### Configure NATS Connection with Ping Interval

Source: https://github.com/nats-io/nats.docs/blob/master/_examples/ping_20s.html

This snippet shows how to establish a NATS connection and configure the ping interval to 20 seconds. It covers implementations in multiple languages, demonstrating common patterns for setting connection options.

```Go
// Set Ping Interval to 20 seconds
nc, err := nats.Connect("demo.nats.io", nats.Name("API Ping Example"), nats.PingInterval(20*time.Second))
if err != nil {
	log.Fatal(err)
}
defers nc.Close()

// Do something with the connection
```

```Java
Options options = new Options.Builder().                             server("nats://demo.nats.io:4222").                             pingInterval(Duration.ofSeconds(20)). // Set Ping Interval                             build();
Connection nc = Nats.connect(options);

// Do something with the connection

nc.close();
```

```JavaScript
let nc = NATS.connect({
    pingInterval: 20*1000, //20s
    url: "nats://demo.nats.io:4222"
});
```

```Python
nc = NATS()

await nc.connect(
    servers=["nats://demo.nats.io:4222"],
    # Set Ping Interval to 20 seconds
    ping_interval=20,
    )

# Do something with the connection.
```

```Ruby
require 'nats/client'

NATS.start(ping_interval: 20) do |nc|
  nc.on_reconnect do
    puts "Got reconnected to #{nc.connected_server}"
  end
  nc.on_disconnect do |reason|
    puts "Got disconnected! #{reason}"
  end
  # Do something with the connection
end
```

```TypeScript
// will throw an exception if connection fails
let nc = await connect({
    pingInterval: 20*1000, //20s
    url: "nats://demo.nats.io:4222"
});
nc.close();
```

--------------------------------

### Go - Discovered Servers Handler

Source: https://github.com/nats-io/nats.docs/blob/master/_examples/servers_added.html

This Go example demonstrates how to use the `nats.DiscoveredServersHandler` to be notified when new servers join the cluster and to log all known and discovered servers.

```APIDOC
## Go - Discovered Servers Handler

### Description
This example shows how to set up a handler in the Go NATS client to be notified when new servers are discovered in the cluster. It logs the list of all known servers and the newly discovered ones.

### Method
N/A (Event-driven callback)

### Endpoint
N/A (Client-side configuration)

### Parameters
N/A

### Request Example
N/A

### Response
N/A (Callback function is invoked)

#### Success Response (Callback Invocation)
- **Known servers**: `[]string` - A slice of strings representing all servers the client is aware of.
- **Discovered servers**: `[]string` - A slice of strings representing the servers that were just discovered.

#### Response Example
```
Known servers: ["nats://demo.nats.io:4222"]
Discovered servers: ["nats://demo.nats.io:4222"]
```
```

--------------------------------

### NATS Auto-Unsubscribe in Go

Source: https://github.com/nats-io/nats.docs/blob/master/_examples/unsubscribe_auto.html

This Go code snippet shows how to set up both synchronous and asynchronous subscriptions with auto-unsubscribe. After the specified number of messages (1 in this case), the subscription is automatically closed by the NATS server.

```Go
nc, err := nats.Connect("demo.nats.io")
if err != nil {
	log.Fatal(err)
}
deferr nc.Close()

// Sync Subscription
sub, err := nc.SubscribeSync("updates")
if err != nil {
	log.Fatal(err)
}
if err := sub.AutoUnsubscribe(1); err != nil {
	log.Fatal(err)
}

// Async Subscription
sub, err = nc.Subscribe("updates", func(_ *nats.Msg) {}) 
if err != nil {
	log.Fatal(err)
}
if err := sub.AutoUnsubscribe(1); err != nil {
	log.Fatal(err)
}
```