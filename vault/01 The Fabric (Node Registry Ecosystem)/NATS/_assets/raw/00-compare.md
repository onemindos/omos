### Publish Messages with NATS CLI

Source: https://github.com/nats-io/nats.docs/blob/master/howto-examples/overview.md

Demonstrates how to use the NATS CLI publish command to send messages. It includes examples of using Go templates for dynamic message content like counters and random strings.

```shell
nats publish
nats pub test --count 10 "Message {{Count}} @ {{Time}}"
nats pub test --count 10 "Message {{Count}}: {{ Random 10 100 }}"
```

--------------------------------

### NATS Wildcard Subscription and Message Handling

Source: https://github.com/nats-io/nats.docs/blob/master/using-nats/developing-with-nats/receiving/wildcards.md

Demonstrates how to subscribe to wildcard subjects (e.g., 'time.>') and process incoming messages. Each implementation handles message reception and cleanup according to the language-specific NATS client library.

```python
await nc.publish("time.A.east", b'A')
await nc.publish("time.B.east", b'B')
await nc.publish("time.C.west", b'C')
await nc.publish("time.D.west", b'D')

for i in range(0, 4):
  msg = await queue.get()
  print("Msg:", msg)

await nc.close()
```

```csharp
await using var client = new NatsClient();
var count = 0;
await foreach (var msg in client.SubscribeAsync<string>("time.>"))
{
    Console.WriteLine($"Received {++count}: {msg.Subject}: {msg.Data}");
    if (count == 4) break;
}
Console.WriteLine("Done");
```

```ruby
NATS.start(servers:["nats://127.0.0.1:4222"]) do |nc|
  Fiber.new do
    f = Fiber.current
    nc.subscribe("time.>") { |msg, reply| f.resume Time.now.to_f }
    nc.publish("time.A.east", "A")
    nc.publish("time.B.east", "B")
    nc.publish("time.C.west", "C")
    nc.publish("time.D.west", "D")
    4.times { msg = Fiber.yield; puts "Msg: #{msg}" }
  end.resume
end
```

```typescript
await nc.subscribe('time.>', (err, msg) => {
    let time = "";
    switch (msg.subject) {
        case 'time.us.east': time = new Date().toLocaleTimeString("en-us", {timeZone: "America/New_York"}); break;
        default: time = "Unknown";
    }
    t.log(msg.subject, time);
});
```

```c
static void onMsg(natsConnection *conn, natsSubscription *sub, natsMsg *msg, void *closure) {
    printf("Received msg: %s - %.*s\n", natsMsg_GetSubject(msg), natsMsg_GetDataLength(msg), natsMsg_GetData(msg));
    natsMsg_Destroy(msg);
}

natsConnection *conn = NULL;
natsSubscription *sub = NULL;
natsConnection_ConnectTo(&conn, NATS_DEFAULT_URL);
natsConnection_Subscribe(&sub, conn, "time.>", onMsg, NULL);
```

--------------------------------

### Consuming Messages with a Pull-Based Consumer using `nats` CLI

Source: https://github.com/nats-io/nats.docs/blob/master/running-a-nats-service/nats_admin/jetstream_admin/consumers.md

Demonstrates how to fetch and acknowledge messages from a pull-based consumer named 'DISPATCH' on the 'ORDERS' stream using the `nats` CLI. It shows the output of received messages and the acknowledgment.

```shell
nats con next ORDERS DISPATCH
```

--------------------------------

### NATS Messaging Operations

Source: https://github.com/nats-io/nats.docs/blob/master/using-nats/developing-with-nats/receiving/unsub_after.md

Examples of publishing messages, subscribing to subjects, and managing subscription lifecycles using NATS client libraries.

```APIDOC
## POST /publish

### Description
Publishes a message to a specific subject on the NATS server.

### Method
PUBLISH

### Parameters
#### Path Parameters
- **subject** (string) - Required - The subject to publish the message to.

#### Request Body
- **payload** (bytes/string) - Required - The data to be sent.

### Request Example
{"subject": "updates", "payload": "data"}

## POST /subscribe

### Description
Subscribes to a subject and optionally sets limits on the number of messages received before auto-unsubscribing.

### Method
SUBSCRIBE

### Parameters
#### Query Parameters
- **max_msgs** (integer) - Optional - Automatically unsubscribe after receiving this many messages.

### Response
#### Success Response (200)
- **data** (string) - The received message payload.

#### Response Example
{"data": "Received message content"}
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

### Publish Message using NATS CLI

Source: https://github.com/nats-io/nats.docs/blob/master/running-a-nats-service/configuration/leafnodes/jetstream_leafnodes.md

Publishes a message 'hello world 2' to the 'test' subject using the NATS CLI. This is a simple way to send messages to a NATS server for testing purposes.

```bash
nats  --server nats://acc:acc@localhost:4222 pub test "hello world 2"
```

--------------------------------

### NATS JetStream Consumer Message Flow Subjects

Source: https://github.com/nats-io/nats.docs/blob/master/using-nats/jetstream/nats_api_reference.md

Subjects related to the flow of messages to and from NATS JetStream consumers, including next message retrieval and acknowledgments.

```text
$JS.API.CONSUMER.MSG.NEXT.<stream>.<consumer>
$JS.SNAPSHOT.RESTORE.<stream>.<msg id>
$JS.ACK.<stream>.<consumer>.x.x.x
$JS.SNAPSHOT.ACK.<stream>.<msg id>
$JS.FC.<stream>.>
```

--------------------------------

### Benchmark NATS JetStream Publication and Consumption

Source: https://github.com/nats-io/nats.docs/blob/master/using-nats/nats-tools/nats_cli/natsbench.md

This snippet demonstrates how to benchmark both the publication and consumption of messages in NATS JetStream. It involves purging the stream, starting consumers, and then initiating asynchronous publishers. The output provides detailed statistics on message rates and latencies for both publishers and consumers.

```bash
nats bench js ordered --purge --clients 8 --no-progress
```

```bash
nats bench js pub async jsfoo --clients 8 --no-progress
```

--------------------------------

### Publish NATS Message to Queue Group

Source: https://github.com/nats-io/nats.docs/blob/master/nats-concepts/core-nats/queue-groups/queues_walkthrough.md

Publishes a single NATS message to the 'foo' subject. One of the subscribers in the queue group will receive and process this message.

```bash
nats request foo "Simple request"
```

--------------------------------

### Deliver Message (MSG) - NATS Protocol

Source: https://github.com/nats-io/nats.docs/blob/master/reference/nats-protocol/nats-protocol/README.md

The MSG protocol message is sent by the NATS server to deliver an application message to a subscribed client. It includes the subject, subscription ID, optional reply-to subject, payload size, and the payload itself.

```nats
MSG <subject> <sid> [reply-to] <#bytes>\r\n[payload]\r\n
```

```nats
MSG FOO.BAR 9 11\r\nHello World\r\n
```

```nats
MSG FOO.BAR 9 GREETING.34 11\r\nHello World\r\n
```

--------------------------------

### Checking Stream Information for Deduplication

Source: https://github.com/nats-io/nats.docs/blob/master/using-nats/jetstream/model_deep_dive.md

After publishing messages with `Nats-Msg-Id`, you can inspect the stream to confirm that only the first occurrence of a message was stored.

```shell
nats stream info ORDERS
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

### Test NATS Connectivity

Source: https://github.com/nats-io/nats.docs/blob/master/using-nats/nats-tools/nsc/basics.md

Use the NATS CLI to subscribe to and publish messages using credential files.

```shell
nats sub --creds ~/.nkeys/creds/MyOperator/MyAccount/MyUser.creds ">"
```

```shell
nats pub --creds ~/.nkeys/creds/MyOperator/MyAccount/MyUser.creds hello NATS
```

--------------------------------

### NATS JetStream Message Acknowledgment Subject Format

Source: https://github.com/nats-io/nats.docs/blob/master/using-nats/jetstream/nats_api_reference.md

Explains the format of the reply subject for acknowledging messages in NATS JetStream.

```text
$JS.ACK.<stream>.<consumer>.<delivered count>.<stream sequence>.<consumer sequence>.<timestamp>.<pending messages>
```

--------------------------------

### Implement Slow Consumer Detection in NATS Clients

Source: https://github.com/nats-io/nats.docs/blob/master/_examples/slow_listener.html

Demonstrates how to configure error listeners or callbacks to handle slow consumer events in Go, Java, and Python NATS clients. This allows applications to react when the consumer cannot keep up with the message stream.

```go
nc, err := nats.Connect("demo.nats.io", nats.ErrorHandler(logSlowConsumer))
if err != nil {
	log.Fatal(err)
}
defer nc.Close()
```

```java
class SlowConsumerReporter implements ErrorListener {
    public void slowConsumerDetected(Connection conn, Consumer consumer) {
        System.out.println("A slow consumer dropped messages: "+ consumer.getDroppedCount());
    }
}
Options options = new Options.Builder().server("nats://demo.nats.io:4222").errorListener(new SlowConsumerReporter()).build();
Connection nc = Nats.connect(options);
```

```python
async def error_cb(e):
    if type(e) is nats.aio.errors.ErrSlowConsumer:
        print("Slow consumer error, unsubscribing...")
        await nc.unsubscribe(e.sid)
await nc.connect(servers=["nats://demo.nats.io:4222"], error_cb=error_cb)
```

--------------------------------

### NATS HPUB Message Description

Source: https://github.com/nats-io/nats.docs/blob/master/reference/nats-protocol/nats-protocol/README.md

The HPUB message extends the PUB message by allowing the inclusion of NATS headers along with the message payload. Headers are structured as 'name: value' pairs and support multi-value headers. The total message size calculation must account for both headers and payload. This is useful for adding metadata to messages.

```text
HPUB <subject> [reply-to] <#bytes>
[headers]
[payload]

```

--------------------------------

### Run NATS Request-Reply Latency Benchmark

Source: https://github.com/nats-io/nats.docs/blob/master/using-nats/nats-tools/nats_cli/natsbench.md

Demonstrates how to start a service and a requester to measure basic request-reply latency. This uses the nats bench service command to simulate a single-client synchronous interaction.

```bash
nats bench service serve foo
nats bench service request foo
```

--------------------------------

### NATS Request-Reply Manual Implementation (Go)

Source: https://github.com/nats-io/nats.docs/blob/master/using-nats/developing-with-nats/sending/request_reply.md

Illustrates a manual implementation of the request-reply pattern in Go, showing the underlying steps of subscribing, publishing a request with a reply-to subject, and waiting for a single message. Includes error handling and unsubscribing.

```Go
sub, err := nc.SubscribeSync(replyTo)
if err != nil {
    log.Fatal(err)
}

// Send the request immediately
nc.PublishRequest(subject, replyTo, []byte(input))
nc.Flush()

// Wait for a single response
for {
    msg, err := sub.NextMsg(1 * time.Second)
    if err != nil {
        log.Fatal(err)
    }

    response = string(msg.Data)
    break
}
sub.Unsubscribe()
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

### NATS PING/PONG for Liveness

Source: https://github.com/nats-io/nats.docs/blob/master/reference/nats-protocol/nats-protocol-demo.md

Demonstrates the PING and PONG commands used for maintaining connection liveness between the NATS client and server. The server sends PING, and the client responds with PONG. Failure to do so can result in connection closure.

```text
PONG
```

--------------------------------

### Subscribe to Wildcard Subjects in NATS

Source: https://github.com/nats-io/nats.docs/blob/master/_examples/subscribe_arrow.html

Demonstrates how to connect to a NATS server and subscribe to wildcard subjects using the '>' operator. These examples show how to handle message callbacks and synchronize message receipt using language-specific primitives like WaitGroups, Latches, or Async Queues.

```go
nc, err := nats.Connect("demo.nats.io")
if err != nil {
	log.Fatal(err)
}
defer nc.Close()
wg := sync.WaitGroup{}
wg.Add(4)
if _, err := nc.Subscribe("time.>", func(m *nats.Msg) {
	log.Printf("%s: %s", m.Subject, m.Data)
	wg.Done()
}); err != nil {
	log.Fatal(err)
}
wg.Wait()
nc.Close()
```

```java
Connection nc = Nats.connect("nats://demo.nats.io:4222");
CountDownLatch latch = new CountDownLatch(4);
Dispatcher d = nc.createDispatcher((msg) -> {
    String subject = msg.getSubject();
    String str = new String(msg.getData(), StandardCharsets.UTF_8);
    System.out.println(subject + ": " + str);
    latch.countDown();
});
d.subscribe("time.>");
latch.await();
nc.close();
```

```python
nc = NATS()
await nc.connect(servers=["nats://demo.nats.io:4222"])
queue = asyncio.Queue()
async def cb(msg):
  await queue.put(msg)
await nc.subscribe("time.>", cb=cb)
await nc.publish("time.A.east", b'A')
await nc.publish("time.B.east", b'B')
await nc.publish("time.C.west", b'C')
await nc.publish("time.D.west", b'D')
for i in range(0, 4):
  msg = await queue.get()
  print("Msg:", msg)
await nc.close()
```

--------------------------------

### NATS JetStream Flow Control Message Reply

Source: https://github.com/nats-io/nats.docs/blob/master/using-nats/jetstream/nats_api_reference.md

Demonstrates how to respond to a flow control message in NATS JetStream. This is crucial to prevent consumer stalling.

```go
msg.Respond(nil)
```