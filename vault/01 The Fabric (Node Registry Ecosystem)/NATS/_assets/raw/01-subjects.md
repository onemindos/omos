### Process Wildcard Messages with Subject Logic

Source: https://github.com/nats-io/nats.docs/blob/master/_examples/subscribe_arrow.html

Shows how to route or process messages based on specific subject patterns within a wildcard subscription. Useful for dynamic responses based on the received subject string.

```javascript
let nc = NATS.connect({url: "nats://demo.nats.io:4222"});
nc.subscribe('time.>', (msg, reply, subject) => {
    let time = "";
    switch (subject) {
        case 'time.us.east':
            time = new Date().toLocaleTimeString("en-us", {timeZone: "America/New_York"});
            break;
        default:
            time = "Unknown";
    }
    t.log(subject, time);
});
```

```typescript
await nc.subscribe('time.>', (err, msg) => {
    let time = "";
    switch (msg.subject) {
        case 'time.us.east':
            time = new Date().toLocaleTimeString("en-us", {timeZone: "America/New_York"});
            break;
        default:
            time = "Unknown";
    }
    t.log(msg.subject, time);
});
```

--------------------------------

### Map Subject with Partitioning on Multiple Tokens

Source: https://github.com/nats-io/nats.docs/blob/master/nats-concepts/subject_mapping.md

Maps a subject with multiple wildcards to a new subject, using a partition number derived from the union of multiple specified wildcard tokens. This distributes messages across a specified number of partitions.

```bash
nats server mapping "foo.*.*" "foo.{{wildcard(1)}}.{{wildcard(2)}}.{{partition(3,1,2)}}"
```

--------------------------------

### NATS Subject Token Reordering Example

Source: https://github.com/nats-io/nats.docs/blob/master/running-a-nats-service/configuration/configuring_subject_mapping.md

Demonstrates how to use wildcard tokens (`$<N>`) in subject mappings to reorder parts of a subject. This is useful for migrating to new subject namespaces or restructuring message flows.

```text
bar.*.*: baz.$2.$1
```

--------------------------------

### Map Subject with Partitioning

Source: https://github.com/nats-io/nats.docs/blob/master/nats-concepts/subject_mapping.md

Maps a wildcard subject to a new subject that includes a partition number derived from a wildcard token. This is useful for distributing messages across multiple partitions.

```bash
nats server mapping "neworders.*" "neworders.{{wildcard(1)}}.{{partition(3,1)}}" neworders.customerid1
> neworders.customerid1.0
```

--------------------------------

### Subscribe to Wildcard Subjects and Synchronize Messages

Source: https://github.com/nats-io/nats.docs/blob/master/using-nats/developing-with-nats/receiving/wildcards.md

Demonstrates subscribing to wildcard subjects to receive multiple messages and using synchronization mechanisms to wait for specific message counts before proceeding.

```Go
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

```Java
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

```JavaScript
let nc = NATS.connect({
    url: "nats://demo.nats.io:4222"
});

nc.subscribe('time.>', (msg, reply, subject) => {
    let time = "";
    switch (subject) {
        case 'time.us.east':
            time = new Date().toLocaleTimeString("en-us", {timeZone: "America/New_York"});
            break;
        default:
            time = "Unknown";
    }
    console.log(subject, time);
});
```

```Python
nc = NATS()
await nc.connect(servers=["nats://demo.nats.io:4222"])

queue = asyncio.Queue()
async def cb(msg):
  await queue.put(msg)

await nc.subscribe("time.>", cb=cb)
```

--------------------------------

### SUBSCRIBE Wildcard Subjects

Source: https://github.com/nats-io/nats.docs/blob/master/using-nats/developing-with-nats/receiving/wildcards.md

This endpoint/operation allows a client to subscribe to a subject pattern containing wildcards (e.g., *). When a message is published to a matching subject, the client receives the message and can inspect the specific subject.

```APIDOC
## SUBSCRIBE [Subject Pattern]

### Description
Subscribes to a NATS subject pattern that includes wildcards. The wildcard '*' matches any single token in the subject hierarchy.

### Method
SUBSCRIBE

### Endpoint
[Subject Pattern (e.g., time.*.east)]

### Parameters
#### Path Parameters
- **subject** (string) - Required - The subject string including wildcards to subscribe to.

### Request Example
{
  "subject": "time.*.east"
}

### Response
#### Success Response (200)
- **msg** (object) - The message object containing the actual subject and data payload.

#### Response Example
{
  "subject": "time.us.east",
  "data": "12:00:00 PM"
}
```

--------------------------------

### Subject Token Reordering

Source: https://github.com/nats-io/nats.docs/blob/master/nats-concepts/subject_mapping.md

Remaps subject tokens by referencing wildcard positions using `{{wildcard(position)}}`. This allows for reordering parts of a subject.

```bash
nats server mapping "bar.*.*"  "baz.{{wildcard(2)}}.{{wildcard(1)}}" bar.a.b
> baz.b.a
```

--------------------------------

### Dropping Subject Tokens

Source: https://github.com/nats-io/nats.docs/blob/master/nats-concepts/subject_mapping.md

Allows dropping subject tokens by not using all wildcard tokens in the destination transform. Note: all wildcards must be used in import/export mappings.

```bash
nats server mapping "orders.*.*" "foo.{{wildcard(2)}}" orders.local.order1
> orders.order1
```

--------------------------------

### Implement NATS Wildcard Subscriptions

Source: https://github.com/nats-io/nats.docs/blob/master/using-nats/developing-with-nats/receiving/wildcards.md

Demonstrates subscribing to wildcard subjects in NATS across different languages. The examples show how to connect to a server, define a subscription pattern, and process messages based on the subject.

```go
nc, err := nats.Connect("demo.nats.io")
if err != nil {
    log.Fatal(err)
}
defer nc.Close()

wg := sync.WaitGroup{}
wg.Add(2)

if _, err := nc.Subscribe("time.*.east", func(m *nats.Msg) {
    log.Printf("%s: %s", m.Subject, m.Data)
    wg.Done()
}); err != nil {
    log.Fatal(err)
}

wg.Wait()
```

```java
Connection nc = Nats.connect("nats://demo.nats.io:4222");
CountDownLatch latch = new CountDownLatch(2);

Dispatcher d = nc.createDispatcher((msg) -> {
    String subject = msg.getSubject();
    String str = new String(msg.getData(), StandardCharsets.UTF_8);
    System.out.println(subject + ": " + str);
    latch.countDown();
});

d.subscribe("time.*.east");
latch.await();
nc.close();
```

```javascript
nc.subscribe("time.us.*", (_err, msg) => {
    let time;
    switch (msg.subject) {
      case "time.us.east":
        time = new Date().toLocaleTimeString("en-us", { timeZone: "America/New_York" });
        break;
      case "time.us.central":
        time = new Date().toLocaleTimeString("en-us", { timeZone: "America/Chicago" });
        break;
      case "time.us.mountain":
        time = new Date().toLocaleTimeString("en-us", { timeZone: "America/Denver" });
        break;
      case "time.us.west":
        time = new Date().toLocaleTimeString("en-us", { timeZone: "America/Los_Angeles" });
        break;
      default:
        time = "I don't know what you are talking about Willis";
    }
    t.log(subject, time);
});
```

```python
nc = NATS()
await nc.connect(servers=["nats://demo.nats.io:4222"])

queue = asyncio.Queue()
async def cb(msg):
  await queue.put_nowait(msg)

await nc.subscribe("time.*.east", cb=cb)
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

### Map Subject with Composite Partition Key

Source: https://github.com/nats-io/nats.docs/blob/master/nats-concepts/subject_mapping.md

Maps a subject with multiple wildcards to a new subject, using a composite partition key derived from multiple wildcard tokens. This allows for finer-grained partitioning.

```bash
nats server mapping "foo.*.*" "foo.{{wildcard(1)}}.{{wildcard(2)}}.{{partition(5,1,2)}}" foo.us.customerid 
> foo.us.customerid.0
```

--------------------------------

### Subscribe to a NATS Subject

Source: https://github.com/nats-io/nats.docs/blob/master/reference/nats-protocol/nats-protocol-demo.md

Subscribes to a NATS subject using the SUB command. This example subscribes to the wildcard subject 'foo.*' with a subscription ID of '90'. A successful subscription is confirmed with an '+OK' message.

```text
SUB foo.* 90
```

```text
+OK
```

--------------------------------

### Configure NATS Subject Mappings (Server Config)

Source: https://github.com/nats-io/nats.docs/blob/master/running-a-nats-service/configuration/configuring_subject_mapping.md

Defines subject mappings within the NATS server configuration file for static authentication. Supports simple direct mapping, token reordering using wildcards, weighted mappings for canary/A/B testing, and introducing message loss for chaos testing.

```text
mappings = {

  # Simple direct mapping.  Messages published to foo are mapped to bar.
  foo: bar

  # remapping tokens can be done with $<N> representing token position.
  # In this example bar.a.b would be mapped to baz.b.a.
  bar.*.*: baz.$2.$1

  # You can scope mappings to a particular cluster
  foo.cluster.scoped : [
    { destination: bar.cluster.scoped, weight:100%, cluster: us-west-1 }
  ]

  # Use weighted mapping for canary testing or A/B testing.  Change dynamically
  # at any time with a server reload.
  myservice.request: [
    { destination: myservice.request.v1, weight: 90% },
    { destination: myservice.request.v2, weight: 10% }
  ]

  # A testing example of wildcard mapping balanced across two subjects.
  # 20% of the traffic is mapped to a service in QA coded to fail.
  myservice.test.*: [
    { destination: myservice.test.$1, weight: 80% },
    { destination: myservice.test.fail.$1, weight: 20% }
  ]

  # A chaos testing trick that introduces 50% artificial message loss of
  # messages published to foo.loss
  foo.loss.>: [ { destination: foo.loss.>, weight: 50% } ]
}
```

--------------------------------

### Publishing to Hierarchical Subjects with NATS

Source: https://github.com/nats-io/nats.docs/blob/master/_examples/wildcard_tester.html

Demonstrates connecting to a NATS server and publishing messages to specific hierarchical subjects. This pattern allows for flexible subscription models using NATS wildcards.

```go
nc, err := nats.Connect("demo.nats.io")
if err != nil {
	log.Fatal(err)
}
defer nc.Close()

zoneID, err := time.LoadLocation("America/New_York")
if err != nil {
	log.Fatal(err)
}
now := time.Now()
zoneDateTime := now.In(zoneID)
formatted := zoneDateTime.String()

nc.Publish("time.us.east", []byte(formatted))
nc.Publish("time.us.east.atlanta", []byte(formatted))

zoneID, err = time.LoadLocation("Europe/Warsaw")
if err != nil {
	log.Fatal(err)
}
zoneDateTime = now.In(zoneID)
formatted = zoneDateTime.String()

nc.Publish("time.eu.east", []byte(formatted))
nc.Publish("time.eu.east.warsaw", []byte(formatted))
```

```java
Connection nc = Nats.connect("nats://demo.nats.io:4222");
ZoneId zoneId = ZoneId.of("America/New_York");
ZonedDateTime zonedDateTime = ZonedDateTime.ofInstant(Instant.now(), zoneId);
String formatted = zonedDateTime.format(DateTimeFormatter.ISO_ZONED_DATE_TIME);

nc.publish("time.us.east", formatted.getBytes(StandardCharsets.UTF_8));
nc.publish("time.us.east.atlanta", formatted.getBytes(StandardCharsets.UTF_8));

zoneId = ZoneId.of("Europe/Warsaw");
zonedDateTime = ZonedDateTime.ofInstant(Instant.now(), zoneId);
formatted = zonedDateTime.format(DateTimeFormatter.ISO_ZONED_DATE_TIME);
nc.publish("time.eu.east", formatted.getBytes(StandardCharsets.UTF_8));
nc.publish("time.eu.east.warsaw", formatted.getBytes(StandardCharsets.UTF_8));

nc.flush(Duration.ZERO);
nc.close();
```

```python
nc = NATS()
await nc.connect(servers=["nats://demo.nats.io:4222"])

await nc.publish("time.us.east", b'...')
await nc.publish("time.us.east.atlanta", b'...')

await nc.publish("time.eu.east", b'...')
await nc.publish("time.eu.east.warsaw", b'...')

await nc.close()
```

```ruby
NATS.start do |nc|
  nc.publish("time.us.east", '...')
  nc.publish("time.us.east.atlanta", '...')
  nc.publish("time.eu.east", '...')
  nc.publish("time.eu.east.warsaw", '...')
  nc.drain
end
```

--------------------------------

### Subscribe to Wildcard Subject in Java

Source: https://github.com/nats-io/nats.docs/blob/master/_examples/subscribe_star.html

This Java code snippet shows how to subscribe to a wildcard subject like 'time.*.east' using the NATS client library for Java. It uses a CountDownLatch to wait for two messages. Dependencies include the 'io.nats:nats-client' library.

```Java
Connection nc = Nats.connect("nats://demo.nats.io:4222");

// Use a latch to wait for 2 messages to arrive
CountDownLatch latch = new CountDownLatch(2);

// Create a dispatcher and inline message handler
Dispatcher d = nc.createDispatcher((msg) -> {
    String subject = msg.getSubject();
    String str = new String(msg.getData(), StandardCharsets.UTF_8);
    System.out.println(subject + ": " + str);
    latch.countDown();
});

// Subscribe
d.subscribe("time.*.east");

// Wait for messages to come in
latch.await();

// Close the connection
nc.close();
```

--------------------------------

### Subscribe to Wildcard Subject in Go

Source: https://github.com/nats-io/nats.docs/blob/master/_examples/subscribe_star.html

This Go code snippet demonstrates how to subscribe to a wildcard subject like 'time.*.east' using the NATS client library. It waits for two messages to arrive before closing the connection. Dependencies include the 'nats-io/nats.go' library.

```Go
nc, err := nats.Connect("demo.nats.io")
if err != nil {
	log.Fatal(err)
}
deferr nc.Close()

// Use a WaitGroup to wait for 2 messages to arrive
wg := sync.WaitGroup{}
wg.Add(2)

// Subscribe
if _, err := nc.Subscribe("time.*.east", func(m *nats.Msg) {
	log.Printf("%s: %s", m.Subject, m.Data)
	wg.Done()
}); err != nil {
	log.Fatal(err)
}

// Wait for the 2 messages to come in
wg.Wait()
```

--------------------------------

### Subscribe and Publish in Python

Source: https://github.com/nats-io/nats.docs/blob/master/_examples/subscribe_star.html

This Python code snippet demonstrates subscribing to a wildcard subject ('time.*.east') and also publishing messages to trigger the subscription. It uses asyncio for asynchronous operations and the 'nats-python' library. It waits for two messages to be received.

```Python
nc = NATS()

await nc.connect(servers=["nats://demo.nats.io:4222"])

# Use queue to wait for 2 messages to arrive
queue = asyncio.Queue()
async def cb(msg):
  await queue.put_nowait(msg)

await nc.subscribe("time.*.east", cb=cb)

# Send 2 messages and wait for them to come in
await nc.publish("time.A.east", b'A')
await nc.publish("time.B.east", b'B')

msg_A = await queue.get()
msg_B = await queue.get()

print("Msg A:", msg_A)
print("Msg B:", msg_B)
```