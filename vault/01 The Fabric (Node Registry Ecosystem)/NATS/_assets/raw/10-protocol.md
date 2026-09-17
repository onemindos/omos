### SUB Protocol

Source: https://github.com/nats-io/nats.docs/blob/master/reference/nats-protocol/nats-protocol/README.md

The SUB command initiates a subscription to a subject.

```APIDOC
## SUB <subject> [queue group] <sid>

### Description
Initiates a subscription to a subject, optionally joining a distributed queue group.

### Parameters
- **subject** (string) - Required - The subject name.
- **queue group** (string) - Optional - The queue group name.
- **sid** (string) - Required - A unique alphanumeric subscription ID.

### Request Example
SUB FOO 1\r\n
```

--------------------------------

### INFO Protocol Message

Source: https://github.com/nats-io/nats.docs/blob/master/reference/nats-protocol/nats-protocol/README.md

The INFO message is sent by the NATS server to the client upon initial connection to provide server configuration and security requirements.

```APIDOC
## INFO

### Description
A client starts as a plain TCP connection. Upon acceptance, the server sends an INFO message containing server configuration, security requirements, and authentication details necessary for the client to proceed.

### Method
SERVER_PUSH

### Endpoint
N/A (Protocol Level Message)

### Parameters
#### Path Parameters
- None

#### Query Parameters
- None

#### Request Body
- **server_id** (string) - Server unique identifier
- **version** (string) - Server version
- **auth_required** (boolean) - Indicates if authentication is required

### Request Example
N/A (Sent by Server)

### Response
#### Success Response (N/A)
- **INFO** (object) - Server configuration details

#### Response Example
INFO {"server_id":"my-server","version":"2.9.0","auth_required":true}
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

### Java: Subscribe and Receive Single Message

Source: https://github.com/nats-io/nats.docs/blob/master/using-nats/developing-with-nats/receiving/sync.md

Illustrates how to create a synchronous subscription to the 'updates' subject using the Java NATS client. It connects to the NATS server, subscribes, fetches a single message, prints its content, and closes the connection.

```Java
Connection nc = Nats.connect("nats://demo.nats.io:4222")

// Subscribe
Subscription sub = nc.subscribe("updates")

// Read a message
Message msg = sub.nextMessage(Duration.ZERO)

String str = new String(msg.getData(), StandardCharsets.UTF_8)
System.out.println(str)

// Close the connection
nc.close()
```

--------------------------------

### NATS INFO Protocol Message

Source: https://github.com/nats-io/nats.docs/blob/master/reference/nats-protocol/nats-protocol/README.md

Details the structure and fields of the INFO message sent by the NATS server to clients.

```APIDOC
## INFO Protocol Message

### Description
The INFO message is a JSON-encoded object sent by the NATS server to the client immediately after the connection is established. It provides the client with necessary server configuration and capabilities.

### Method
N/A (Server-to-Client Push)

### Endpoint
N/A (Connection Initialization)

### Parameters
#### Request Body (JSON Object)
- **server_id** (string) - Always - The unique identifier of the NATS server.
- **server_name** (string) - Always - The name of the NATS server.
- **version** (string) - Always - The version of NATS.
- **go** (string) - Always - The version of golang the NATS server was built with.
- **host** (string) - Always - The IP address used to start the NATS server.
- **port** (int) - Always - The port number the NATS server is configured to listen on.
- **headers** (bool) - Always - Whether the server supports headers.
- **max_payload** (int) - Always - Maximum payload size in bytes.
- **proto** (int) - Always - Protocol version of the server.
- **client_id** (uint64) - Optional - Internal client identifier.
- **auth_required** (bool) - Optional - If true, client must authenticate.
- **tls_required** (bool) - Optional - If true, client must perform TLS handshake.
- **jetstream** (bool) - Optional - Whether the server supports JetStream.

### Response Example
{
  "server_id": "S12345",
  "server_name": "NATS-Server-01",
  "version": "2.9.0",
  "go": "go1.19",
  "host": "0.0.0.0",
  "port": 4222,
  "headers": true,
  "max_payload": 1048576,
  "proto": 1
}
```

--------------------------------

### Subscribe with Reply - Java

Source: https://github.com/nats-io/nats.docs/blob/master/_examples/subscribe_w_reply.html

This Java code snippet shows how to subscribe to a 'time' subject and send the current time as a reply. It uses the NATS Java client to connect, subscribe, read a message, and publish the response to the reply-to subject.

```Java
Connection nc = Nats.connect("nats://demo.nats.io:4222");

// Subscribe
Subscription sub = nc.subscribe("time");

// Read a message
Message msg = sub.nextMessage(Duration.ZERO);

// Get the time
Calendar cal = Calendar.getInstance();
SimpleDateFormat sdf = new SimpleDateFormat("HH:mm:ss");
byte[] timeAsBytes = sdf.format(cal.getTime()).getBytes(StandardCharsets.UTF_8);

// Send the time
nc.publish(msg.getReplyTo(), timeAsBytes);

// Flush and close the connection
nc.flush(Duration.ZERO);
nc.close();
```

--------------------------------

### NATS Errors Requiring Connection Closure

Source: https://github.com/nats-io/nats.docs/blob/master/reference/nats-protocol/nats-protocol/README.md

Lists NATS protocol errors that necessitate the client closing the connection due to the invalid state of the connection.

```APIDOC
## NATS Errors Requiring Connection Closure

### Description

Some protocol errors result in the server closing the connection. Upon receiving these errors, the connection is no longer valid and the client should clean up relevant resources.

### Errors

- **Unknown Protocol Operation**: `-ERR 'Unknown Protocol Operation'`
  - **Description**: The client sent an unrecognized protocol command.

- **Attempted To Connect To Route Port**: `-ERR 'Attempted To Connect To Route Port'`
  - **Description**: The client attempted to connect to a route port instead of the client port.

- **Authorization Violation**: `-ERR 'Authorization Violation'`
  - **Description**: The client failed to authenticate to the server with the credentials provided in the `CONNECT` message.

- **Authorization Timeout**: `-ERR 'Authorization Timeout'`
  - **Description**: The client took too long to authenticate after establishing a connection (default 1 second).

- **Invalid Client Protocol**: `-ERR 'Invalid Client Protocol'`
  - **Description**: The client specified an invalid protocol version in the `CONNECT` message.

- **Maximum Control Line Exceeded**: `-ERR 'Maximum Control Line Exceeded'`
  - **Description**: The combined length of the message destination subject and reply subject exceeded the `max_control_line` server option (default 1024 bytes).

- **Parser Error**: `-ERR 'Parser Error'`
  - **Description**: The server could not parse the protocol message sent by the client.

- **Secure Connection - TLS Required**: `-ERR 'Secure Connection - TLS Required'`
  - **Description**: The server requires TLS, but the client did not enable it.

- **Stale Connection**: `-ERR 'Stale Connection'`
  - **Description**: The server has not received a message (including a `PONG`) from the client within the expected timeout period.

- **Maximum Connections Exceeded**: `-ERR 'Maximum Connections Exceeded'`
  - **Description**: The server has reached its maximum connection limit (`max_connections`, default 64k) when attempting to create a new connection.

- **Slow Consumer**: `-ERR 'Slow Consumer'`
  - **Description**: The server's pending data size for the connection has reached the maximum limit (default 10MB).

- **Maximum Payload Violation**: `-ERR 'Maximum Payload Violation'`
  - **Description**: The client attempted to publish a message with a payload size exceeding the `max_payload` configured on the server. This value is provided in the initial `INFO` message.
```

--------------------------------

### Test Private Stream Subscription

Source: https://github.com/nats-io/nats.docs/blob/master/using-nats/nats-tools/nsc/streams.md

This command subscribes to a specific private stream ('private.abc.AAM46E3YF5WOZSE5WNYWHN3YYISVZOSI6XHTF2Q64ECPXSFQZROJMP2H') from account B using user 'b'. This is the first step in testing the stream's functionality.

```bash
nsc tools sub --account B --user b private.abc.AAM46E3YF5WOZSE5WNYWHN3YYISVZOSI6XHTF2Q64ECPXSFQZROJMP2H
```

--------------------------------

### NATS Client CONNECT Message Example (Go)

Source: https://github.com/nats-io/nats.docs/blob/master/reference/nats-protocol/nats-protocol/README.md

An example of the initial CONNECT message sent by a Go NATS client to the server. This message includes client configuration options such as verbosity, TLS requirements, client name, language, version, and protocol version. Most clients default 'verbose' to false.

```go
CONNECT {"verbose":false,"pedantic":false,"tls_required":false,"name":"","lang":"go","version":"1.2.2","protocol":1}

```

--------------------------------

### Connect to NATS with Pedantic Mode

Source: https://github.com/nats-io/nats.docs/blob/master/_examples/connect_pedantic.html

Configures the NATS client connection to operate in pedantic mode, which performs extra checks on protocol messages.

```APIDOC
## CONNECT NATS (Pedantic Mode)

### Description
Establishes a connection to a NATS server with the pedantic flag enabled. Pedantic mode forces the server to perform additional validation on protocol messages sent by the client.

### Method
CONNECT

### Endpoint
nats://demo.nats.io:4222

### Parameters
#### Request Body
- **url** (string) - Required - The NATS server URL.
- **pedantic** (boolean) - Required - Set to true to enable strict protocol checking.

### Request Example
{
  "url": "nats://demo.nats.io:4222",
  "pedantic": true
}

### Response
#### Success Response (200)
- **connection** (object) - The active NATS connection object.

#### Response Example
{
  "status": "connected",
  "server": "nats://demo.nats.io:4222"
}
```

--------------------------------

### Publishing and Subscribing with Auto-Unsubscribe (C)

Source: https://github.com/nats-io/nats.docs/blob/master/using-nats/developing-with-nats/receiving/unsub_after.md

This C code snippet shows how to connect to a NATS server, subscribe to the 'updates' subject, and automatically unsubscribe after receiving one message. It demonstrates synchronous message retrieval using NATS_Client and includes error handling for connection and subscription operations. It also shows how to process received message data and clean up resources.

```c
natsConnection      *conn      = NULL;
natsSubscription    *sub       = NULL;
natsMsg             *msg       = NULL;
natsStatus          s          = NATS_OK;

s = natsConnection_ConnectTo(&conn, NATS_DEFAULT_URL);

// Subscribe
if (s == NATS_OK)
    s = natsConnection_SubscribeSync(&sub, conn, "updates");

// Unsubscribe after 1 message is received
if (s == NATS_OK)
    s = natsSubscription_AutoUnsubscribe(sub, 1);

// Wait for messages
if (s == NATS_OK)
    s = natsSubscription_NextMsg(&msg, sub, 10000);

if (s == NATS_OK)
{
    printf("Received msg: %s - %.*s\n",
            natsMsg_GetSubject(msg),
            natsMsg_GetDataLength(msg),
            natsMsg_GetData(msg));

    // Destroy message that was received
    natsMsg_Destroy(msg);
}

(...)

// Destroy objects that were created
natsSubscription_Destroy(sub);
natsConnection_Destroy(conn);
```

--------------------------------

### Subscribe to Subject (SUB) - NATS Protocol

Source: https://github.com/nats-io/nats.docs/blob/master/reference/nats-protocol/nats-protocol/README.md

The SUB command establishes a subscription to a given subject. It can optionally include a queue group name for distributed message consumption and requires a unique subscription ID (sid).

```nats
SUB <subject> [queue group] <sid>\r\n
```

```nats
SUB FOO 1\r\n
```

```nats
SUB BAR G1 44\r\n
```

--------------------------------

### NATS Errors Where Connection Remains Open

Source: https://github.com/nats-io/nats.docs/blob/master/reference/nats-protocol/nats-protocol/README.md

Lists NATS protocol errors where the server does not close the connection, allowing the client to continue operating after handling the error.

```APIDOC
## NATS Errors Where Connection Remains Open

### Description

Protocol error messages where the connection remains open are listed below. The client should not close the connection in these cases.

### Errors

- **Invalid Subject**: `-ERR 'Invalid Subject'`
  - **Description**: The client sent a malformed subject (e.g., `sub foo. 90`).

- **Permissions Violation for Subscription to <subject>**: `-ERR 'Permissions Violation for Subscription to <subject>'`
  - **Description**: The user specified in the `CONNECT` message does not have permission to subscribe to the specified subject.

- **Permissions Violation for Publish to <subject>**: `-ERR 'Permissions Violation for Publish to <subject>'`
  - **Description**: The user specified in the `CONNECT` message does not have permissions to publish to the specified subject.
```

--------------------------------

### Inspect NATS Connection Handshake

Source: https://github.com/nats-io/nats.docs/blob/master/running-a-nats-service/nats_admin/jwt.md

Example of the NATS server INFO message containing the nonce, and the client CONNECT message containing the signature.

```json
{
    "auth_required": true,
    "nonce": "-QPTE1Jsk8kI3rE",
    "server_id": "NBSHIXACRHUODC4FY2Z3OYXSZSRUBRH6VWIKQNGVPKOTA7H4YTXWJRTO"
}
```

```json
{
    "nkey": "UC435ZYS52HF72E2VMQF4GO6CUJOCHDUUPEBU7XDXW5AQLIC6JZ46PO5",
    "sig": "lopzgs98JBQYyRdw1zT_BoBpSFRDCfTvT4le5MYSKrt0IqGWZ2OXhPW1J_zo2_sBod8XaWgQc9oWohWBN0NdDg"
}
```

--------------------------------

### Synchronous Subscribe and Receive Message in Java

Source: https://github.com/nats-io/nats.docs/blob/master/_examples/subscribe_sync.html

Shows how to connect to a NATS server, create a synchronous subscription to a subject, and retrieve the next available message. This example uses the Java NATS client and includes closing the connection.

```Java
Connection nc = Nats.connect("nats://demo.nats.io:4222");

// Subscribe
Subscription sub = nc.subscribe("updates");

// Read a message
Message msg = sub.nextMessage(Duration.ZERO);

String str = new String(msg.getData(), StandardCharsets.UTF_8);
System.out.println(str);

// Close the connection
nc.close();
```

--------------------------------

### NATS PUB Message Syntax and Examples

Source: https://github.com/nats-io/nats.docs/blob/master/reference/nats-protocol/nats-protocol/README.md

The PUB message is used to publish a message payload to a specified subject. It can optionally include a reply subject for request-reply patterns. The syntax requires the subject, an optional reply-to subject, the payload size in bytes, and the payload itself. Examples demonstrate publishing text, requests with reply subjects, and empty messages.

```text
PUB <subject> [reply-to] <#bytes>
[payload]

```

```text
PUB FOO 11
Hello NATS!

```

```text
PUB FRONT.DOOR JOKE.22 11
Knock Knock

```

```text
PUB NOTIFY 0


```

--------------------------------

### NATS Server INFO Message

Source: https://github.com/nats-io/nats.docs/blob/master/running-a-nats-service/nats_admin/jwt.md

Example of the INFO message sent by a NATS server to a connecting client. It includes server details and a nonce for authentication.

```json
INFO {
    "auth_required": true,
    "client_id": 5,
    "client_ip": "127.0.0.1",
    "go": "go1.14.1",
    "headers": true,
    "host": "localhost",
    "max_payload": 1048576,
    "nonce": "aN9-ZtS7taDoAZk",
    "port": 4222,
    "proto": 1,
    "server_id": "NCIK6FX5MRIEPMEK22YL2ECLIWVJBH2SWFD5EQWSI5XRDQPKZXWKX3VP",
    "server_name": "NCIK6FX5MRIEPMEK22YL2ECLIWVJBH2SWFD5EQWSI5XRDQPKZXWKX3VP",
    "tls_required": true,
    "version": "2.2.0-beta.26"
}
```

--------------------------------

### Publish Request and Subscribe to Reply (Ruby)

Source: https://github.com/nats-io/nats.docs/blob/master/_examples/publish_with_reply.html

This Ruby code snippet demonstrates the publish-request pattern using NATS client and Fibers for asynchronous handling. It subscribes to the 'time' subject, publishes a request with a generated inbox for the reply, and then yields to wait for the response. The connection is managed within NATS.start block.

```Ruby
require 'nats/client'
require 'fiber'

NATS.start(servers:["nats://127.0.0.1:4222"])
do |nc|
  Fiber.new do
    f = Fiber.current

    nc.subscribe("time") do |msg, reply|
      f.resume msg
    end

    nc.publish("time", 'example', NATS.create_inbox)

    # Use the response
    msg = Fiber.yield
    puts "Reply: #{msg}"
  end.resume
end
```