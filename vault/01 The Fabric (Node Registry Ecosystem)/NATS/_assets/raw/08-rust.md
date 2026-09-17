### Connect to NATS with Options (Ruby)

Source: https://github.com/nats-io/nats.docs/blob/master/_examples/connect_options.html

Simulates a connection timeout in the Ruby NATS client by using an `EM.add_timer`. The NATS client itself does not directly support connection timeouts in its API.

```Ruby
# There is currently no connect timeout as part of the Ruby NATS client API, but you can use a timer to mimic it.
require 'nats/client'

timer = EM.add_timer(5) do
  NATS.connect do |nc|
    # Do something with the connection
    # Close the connection
    nc.close
  end
end
EM.cancel_timer(timer)
```

--------------------------------

### Disable NATS Reconnection (Python)

Source: https://github.com/nats-io/nats.docs/blob/master/_examples/reconnect_none.html

Connects to a NATS server using the asyncio client library with automatic reconnection disabled. The `allow_reconnect=False` parameter in the `connect` method ensures that the client will not try to re-establish a connection if it drops.

```Python
nc = NATS()
await nc.connect(
    servers=[
       "nats://demo.nats.io:1222",
       "nats://demo.nats.io:1223",
       "nats://demo.nats.io:1224"
       ],
    allow_reconnect=False,
    )

  # Do something with the connection

  await nc.close()
```

--------------------------------

### Set Connect Timeout in Java

Source: https://github.com/nats-io/nats.docs/blob/master/using-nats/developing-with-nats/connecting/connect_timeout.md

Use the `connectionTimeout` option with `Duration.ofSeconds(10)` to set the connection timeout for the Java NATS client. This ensures the client does not wait indefinitely to connect.

```java
Options options = new Options.Builder()
    .server("nats://demo.nats.io:4222")
    .connectionTimeout(Duration.ofSeconds(10)) // Set timeout
    .build();
Connection nc = Nats.connect(options);

// Do something with the connection

nc.close();
```

--------------------------------

### Handle NATS Connection Events in JavaScript

Source: https://github.com/nats-io/nats.docs/blob/master/using-nats/developing-with-nats/events/events.md

Listen for connection status updates using async iterators in the NATS JavaScript client. This provides real-time information about connection state changes like disconnects, reconnects, and errors. Requires the 'nats' npm package.

```JavaScript
const nc = await connect({ servers: ["demo.nats.io"] });
nc.closed().then(() => {
  t.log("the connection closed!");
});

(async () => {
    for await (const s of nc.status()) {
      switch (s.type) {
        case Events.Disconnect:
          t.log(`client disconnected - ${s.data}`);
          break;
        case Events.LDM:
          t.log("client has been requested to reconnect");
          break;
        case Events.Update:
          t.log(`client received a cluster update - ${s.data}`);
          break;
        case Events.Reconnect:
          t.log(`client reconnected - ${s.data}`);
          break;
        case Events.Error:
          t.log("client got a permissions error");
          break;
        case DebugEvents.Reconnecting:
          t.log("client is attempting to reconnect");
          break;
        case DebugEvents.StaleConnection:
          t.log("client has a stale connection");
          break;
        default:
          t.log(`got an unknown status ${s.type}`);
      }
    }
})().then();
```

--------------------------------

### TypeScript NATS Client - Reconnect Time Wait Example

Source: https://github.com/nats-io/nats.docs/blob/master/_examples/reconnect_10s.html

Example of setting the reconnect time wait to 10 seconds for the TypeScript NATS client.

```APIDOC
## TypeScript NATS Client - Reconnect Time Wait

### Description
This example shows how to configure the TypeScript NATS client with a `reconnectTimeWait` of 10 seconds.

### Method
`connect`

### Endpoint
`nats://demo.nats.io:4222` (or your NATS server address)

### Parameters
#### Query Parameters
None

#### Request Body
None

### Request Example
```typescript
// will throw an exception if connection fails
let nc = await connect({
    reconnectTimeWait: 10*1000, //10s
    servers: ["nats://demo.nats.io:4222"]
});
nc.close();
```

### Response
#### Success Response (200)
Connection established.

#### Response Example
None (connection is established programmatically)
```

--------------------------------

### Connect NATS Client with Token

Source: https://github.com/nats-io/nats.docs/blob/master/running-a-nats-service/configuration/securing_nats/auth_intro/tokens.md

Shows how a client connects to a NATS server by embedding the authentication token directly into the connection URL.

```shell
nats sub -s nats://s3cr3t@localhost:4222 ">"
```

--------------------------------

### Publish Request and Subscribe to Reply (TypeScript)

Source: https://github.com/nats-io/nats.docs/blob/master/_examples/publish_with_reply.html

This TypeScript example shows how to publish a request and handle the reply using async/await. It sets up a subscription to the 'time' subject that publishes the current time to the reply subject. A subscription to a unique inbox subject is created to receive the response, after which the connection is closed.

```TypeScript
// set up a subscription to process the request
await nc.subscribe('time', (err, msg) => {
    if (err) {
        // this example is running inside of a promise
        reject();
        return;
    }
    if (msg.reply) {
        nc.publish(msg.reply, new Date().toLocaleTimeString());
    }
});

// create a subscription subject that the responding send replies to
let inbox = createInbox();
await nc.subscribe(inbox, (err, msg) => {
    t.log('the time is', msg.data);
    // this example is running inside of a promise
    nc.close();
    resolve();
}, {max: 1});

nc.publish('time', "", inbox);
```

--------------------------------

### Note on Reconnect Buffer Support in Python

Source: https://github.com/nats-io/nats.docs/blob/master/using-nats/developing-with-nats/reconnect/buffer.md

States that the Asyncio NATS client currently does not implement a reconnect buffer feature.

```Python
# Asyncio NATS client currently does not implement a reconnect buffer
```

--------------------------------

### Connect to NATS with Options (TypeScript)

Source: https://github.com/nats-io/nats.docs/blob/master/_examples/connect_options.html

Connects to the NATS server using the `connect` function with a specified URL and timeout. This example highlights the TypeScript client's configuration options.

```TypeScript
let nc = await connect({
    url: "nats://demo.nats.io:4222",
    timeout: 1000
});
```

--------------------------------

### Java Client Connection Listener

Source: https://github.com/nats-io/nats.docs/blob/master/using-nats/developing-with-nats/events/events.md

Illustrates how to implement and set a `ConnectionListener` in the NATS Java client to receive notifications about connection events.

```APIDOC
## Java Client Connection Listener

### Description
This example demonstrates how to use the `ConnectionListener` interface in the NATS Java client to handle various connection events. You create a custom class that implements this interface and then register it with the NATS connection options.

### Method
N/A (Configuration during connection)

### Endpoint
N/A

### Parameters
#### Connection Options
- **server** (String) - The NATS server URL.
- **connectionListener** (ConnectionListener) - An instance of a class implementing the `ConnectionListener` interface.

### Request Example
```java
class MyConnectionListener implements ConnectionListener {
    @Override
    public void connectionEvent(Connection natsConnection, Events event) {
        System.out.println("Connection event - " + event);
    }
}

// ... inside main or another method
Options options = new Options.Builder()
    .server("nats://demo.nats.io:4222")
    .connectionListener(new MyConnectionListener()) // Set the listener
    .build();
Connection nc = Nats.connect(options);

// Do something with the connection

nc.close();
```

### Response
#### Success Response (Connection Established)
- **Connection** (Connection) - The established NATS connection object.

#### Response Example
(Output will be printed to console via `System.out.println` in the listener)
```
Connection event - CONNECTED
```
(Or other event types like DISCONNECTED, RECONNECTED, etc.)
```

--------------------------------

### Asynchronous Subscription Behavior in TypeScript

Source: https://github.com/nats-io/nats.docs/blob/master/_examples/subscribe_sync.html

Highlights that subscriptions in the ts-nats (TypeScript) client are always asynchronous. This necessitates the use of asynchronous programming constructs like Promises or async/await for managing subscriptions and messages.

```TypeScript
// ts-nats subscriptions are always async.
```

--------------------------------

### Python (asyncio) NATS Client - Reconnect Time Wait Example

Source: https://github.com/nats-io/nats.docs/blob/master/_examples/reconnect_10s.html

Example of setting the reconnect time wait to 10 seconds for the Python asyncio NATS client.

```APIDOC
## Python (asyncio) NATS Client - Reconnect Time Wait

### Description
This example shows how to set the `reconnect_time_wait` parameter to 10 seconds when connecting with the Python asyncio NATS client.

### Method
`nc.connect`

### Endpoint
`nats://demo.nats.io:4222` (or your NATS server address)

### Parameters
#### Query Parameters
None

#### Request Body
None

### Request Example
```python
nc = NATS()
await nc.connect(
    servers=["nats://demo.nats.io:4222"],
    reconnect_time_wait=10,
    )

# Do something with the connection

await nc.close()
```

### Response
#### Success Response (200)
Connection established.

#### Response Example
None (connection is established programmatically)
```

--------------------------------

### NATS Client CONNECT Message with JWT

Source: https://github.com/nats-io/nats.docs/blob/master/running-a-nats-service/nats_admin/jwt.md

Example of a NATS client's CONNECT message, including the JWT for authentication and a signature for the nonce.

```json
[98019] 2020/10/26 16:07:53.861612 [TRC] 127.0.0.1:56830 - cid:4 - <<- [CONNECT {
    "echo": true,
    "headers": true,
    "jwt": "eyJ0eXAiOiJqd3QiLCJhbGciOiJlZDI1NTE5In0.eyJqdGkiOiJXNkFYSFlSS1RHVTNFUklQM0dSRDdNV0FQTzQ2VzQ2Vzc3R1JNMk5SWFFIQ0VRQ0tCRjJRIiwiaWF0IjoxNjAzNDczNzg4LCJpc3MiOiJBQUFYQVVWU0dLN1RDUkhGSVJBUzRTWVhWSjc2RVdETU5YWk02QVJGR1hQN0JBU05ER0xLVTdBNSIsIm5hbWUiOiJzeXMiLCJzdWIiOiJVRE5ZMktLUFRJQVBQTk9OT0xBVE5SWlBHTVBMTkZXSFFQS1VYSjZBMllUQTQ3Tk41Vk5GSU80NSIsInR5cGUiOiJ1c2VyIiwibmF0cyI6eyJwdWIiOnt9LCJzdWIiOnt9fX0.ae3OvcapjQgbXhI2QbgIs32AWr3iBb2UFRZbXzIg0duFHNPQI5LsprR0OQoSlc2tic6e3sn8YM5x0Rt34FryDA",
    "lang": "go",
    "name": "NATS CLI",
    "no_responders": true,
    "pedantic": false,
    "protocol": 1,
    "sig": "VirwM--xq5i2RI9VEQiFYv_6JBs-IR4oObypglR7qVxYtXDUtIKIr1qXW_M54iHFB6Afu698J_in5CfBRjuVBg",
    "tls_required": true,
    "verbose": false,
    "version": "1.11.0"
}]
```

--------------------------------

### Ruby Client Connection Events

Source: https://github.com/nats-io/nats.docs/blob/master/using-nats/developing-with-nats/events/events.md

Notes that the NATS Ruby client does not have a single listener for connection events, implying that specific handling might require different approaches or is not directly exposed in the same manner as other clients.

```APIDOC
## Ruby Client Connection Events

### Description
This section indicates that the NATS Ruby client does not provide a single, unified listener for all connection events in the same way some other client libraries do. Users may need to consult the specific API documentation for the Ruby client for alternative methods or direct event handling mechanisms.

### Method
N/A

### Endpoint
N/A

### Parameters
N/A

### Request Example
```ruby
# There is not a single listener for connection events in the Ruby NATS Client.
# Consult the specific NATS Ruby client library documentation for details on handling connection events.
```

### Response
N/A
```

--------------------------------

### Mimic Connect Timeout in Ruby

Source: https://github.com/nats-io/nats.docs/blob/master/using-nats/developing-with-nats/connecting/connect_timeout.md

The Ruby NATS client does not have a direct connect timeout option. This example demonstrates how to mimic this behavior using an `EM.add_timer` to cancel the connection attempt after 10 seconds.

```ruby
# There is currently no connect timeout as part of the Ruby NATS client API, but you can use a timer to mimic it.
require 'nats/client'

timer = EM.add_timer(10) do
  NATS.connect do |nc|
    # Do something with the connection

    # Close the connection
    nc.close
  end
end
EM.cancel_timer(timer)
```

--------------------------------

### NATS Client CONNECT Message Example (Go)

Source: https://github.com/nats-io/nats.docs/blob/master/reference/nats-protocol/nats-protocol/README.md

An example of the initial CONNECT message sent by a Go NATS client to the server. This message includes client configuration options such as verbosity, TLS requirements, client name, language, version, and protocol version. Most clients default 'verbose' to false.

```go
CONNECT {"verbose":false,"pedantic":false,"tls_required":false,"name":"","lang":"go","version":"1.2.2","protocol":1}

```

--------------------------------

### Connect to NATS with Options (Python)

Source: https://github.com/nats-io/nats.docs/blob/master/_examples/connect_options.html

Connects to the NATS server asynchronously with a specified connection timeout. It demonstrates the use of the `connect_timeout` parameter.

```Python
nc = NATS()
await nc.connect(connect_timeout=2)  # Do something with the connection

await nc.close()
```

--------------------------------

### Disable NATS Reconnection (JavaScript)

Source: https://github.com/nats-io/nats.docs/blob/master/_examples/reconnect_none.html

Connects to a NATS server using the Node.js client library with automatic reconnection disabled. The `reconnect: false` option in the connection configuration prevents the client from attempting to reconnect if the connection is lost.

```JavaScript
let nc = NATS.connect({
    reconnect: false,
    servers: ["nats://demo.nats.io:4222"]
});
```

--------------------------------

### C# Client Connection Event Handlers

Source: https://github.com/nats-io/nats.docs/blob/master/using-nats/developing-with-nats/events/events.md

Demonstrates how to subscribe to connection events like disconnection and opening using event handlers in the NATS C# client.

```APIDOC
## C# Client Connection Event Handlers

### Description
This example shows how to handle connection events in the NATS C# client by subscribing to events provided by the `NatsClient.Connection` object. Events include disconnection, opening, and failed reconnections.

### Method
N/A (Event subscription)

### Endpoint
N/A

### Parameters
#### Event Handlers
- **ConnectionDisconnected** (event handler) - Subscribes to the event triggered when the client disconnects.
- **ConnectionOpened** (event handler) - Subscribes to the event triggered when the client successfully connects or reconnects.
- **ReconnectFailed** (event handler) - Subscribes to the event triggered when the client fails to reconnect after multiple attempts.

### Request Example
```csharp
// dotnet add package NATS.Net
using NATS.Net;

await using var client = new NatsClient();

client.Connection.ConnectionDisconnected += async (sender, args) =>
{
    Console.WriteLine($"Disconnected: {args.Message}");
};

client.Connection.ConnectionOpened += async (sender, args) =>
{
    Console.WriteLine($"Connected: {args.Message}");
};

client.Connection.ReconnectFailed += async (sender, args) =>
{
    Console.WriteLine($"Reconnect Failed: {args.Message}");
};

// Assuming client is configured with server details
// await client.ConnectAsync(); 
```

### Response
No direct response body. Event handlers are executed asynchronously when the corresponding connection events occur.
```

--------------------------------

### Set Connect Timeout in C#

Source: https://github.com/nats-io/nats.docs/blob/master/using-nats/developing-with-nats/connecting/connect_timeout.md

Configure the `ConnectTimeout` to `TimeSpan.FromSeconds(10)` using `NatsOpts` for the C# NATS client. The connection is established automatically on the first operation if `ConnectAsync` is not explicitly called.

```csharp
// dotnet add package NATS.Net
using NATS.Net;
using NATS.Client.Core;

await using var client = new NatsClient(new NatsOpts
{
    Url = "nats://demo.nats.io:4222",
    ConnectTimeout = TimeSpan.FromSeconds(10)
});

// You don't have to call ConnectAsync() explicitly,
// first operation will make the connection otherwise.
await client.ConnectAsync();
```

--------------------------------

### Connect to NATS with TLS

Source: https://github.com/nats-io/nats.docs/blob/master/using-nats/developing-with-nats/connecting/security/tls.md

Demonstrates how to establish a secure NATS connection using client certificates and CA trust stores across different client libraries.

```go
nc, err := nats.Connect("localhost",
    nats.ClientCert("client-cert.pem", "client-key.pem"),
    nats.RootCAs("rootCA.pem"))
if err != nil {
    log.Fatal(err)
}
defer nc.Close()
```

```java
SSLContext ctx = SSLUtils.createSSLContext();
Options options = new Options.Builder()
    .server("nats://localhost:4222")
    .sslContext(ctx)
    .build();
Connection nc = Nats.connect(options);
```

```javascript
const nc = await connect({
  port: ns.port,
  debug: true,
  tls: {
    caFile: caCertPath,
    keyFile: clientKeyPath,
    certFile: clientCertPath,
  },
});
```

```python
nc = NATS()
ssl_ctx = ssl.create_default_context(purpose=ssl.Purpose.SERVER_AUTH)
ssl_ctx.load_verify_locations('rootCA.pem')
ssl_ctx.load_cert_chain(certfile='client-cert.pem', keyfile='client-key.pem')
await nc.connect(servers=["nats://demo.nats.io:4222"], tls=ssl_ctx)
```

--------------------------------

### Note on Reconnect Buffer Support in Ruby

Source: https://github.com/nats-io/nats.docs/blob/master/using-nats/developing-with-nats/reconnect/buffer.md

Clarifies that the Ruby NATS client does not currently include a reconnect pending buffer.

```Ruby
# There is currently no reconnect pending buffer as part of the Ruby NATS client
```

--------------------------------

### Go NATS Client - Reconnect Wait Example

Source: https://github.com/nats-io/nats.docs/blob/master/_examples/reconnect_10s.html

Example of how to set the reconnect wait interval to 10 seconds when connecting to a NATS server using the Go client library.

```APIDOC
## Go NATS Client - Reconnect Wait

### Description
This example demonstrates how to configure the Go NATS client to wait for 10 seconds before attempting to reconnect to the server.

### Method
`nats.Connect`

### Endpoint
`demo.nats.io` (or your NATS server address)

### Parameters
#### Query Parameters
None

#### Request Body
None

### Request Example
```go
// Set reconnect interval to 10 seconds
nc, err := nats.Connect("demo.nats.io", nats.ReconnectWait(10*time.Second))
if err != nil {
	log.Fatal(err)
}
defer nc.Close()

// Do something with the connection
```

### Response
#### Success Response (200)
Connection established.

#### Response Example
None (connection is established programmatically)
```