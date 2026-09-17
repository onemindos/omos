### Start Second Queue Group Member

Source: https://github.com/nats-io/nats.docs/blob/master/nats-concepts/core-nats/queue-groups/queues_walkthrough.md

Starts another NATS subscriber to join the same queue group. This increases the number of available service instances for load balancing.

```bash
nats reply foo "service instance B Reply# {{Count}}"
```

--------------------------------

### SUBSCRIBE (Queue Group)

Source: https://github.com/nats-io/nats.docs/blob/master/_examples/subscribe_queue.html

Establishes a queue subscription to a subject, allowing load balancing across multiple instances of a service.

```APIDOC
## SUBSCRIBE (Queue Group)

### Description
Subscribes to a subject as part of a queue group. Only one member of the queue group will receive each message, enabling horizontal scaling and load balancing.

### Method
SUBSCRIBE

### Endpoint
[subject] (with queue group option)

### Parameters
#### Path Parameters
- **subject** (string) - Required - The NATS subject to subscribe to.

#### Request Body
- **queue** (string) - Required - The name of the queue group. All subscribers with the same queue name will share the message load.

### Request Example
{
  "subject": "updates",
  "queue": "workers"
}

### Response
#### Success Response (N/A - Async)
- **message** (object) - The received message payload.

#### Response Example
{
  "data": "All is Well",
  "subject": "updates"
}
```

--------------------------------

### Implement NATS Queue Subscriptions

Source: https://github.com/nats-io/nats.docs/blob/master/using-nats/developing-with-nats/receiving/queues.md

Demonstrates how to subscribe to a subject as part of a named queue group to enable load balancing across multiple application instances.

```go
nc, err := nats.Connect("demo.nats.io")
if err != nil {
    log.Fatal(err)
}
defer nc.Close()

wg := sync.WaitGroup{}
wg.Add(10)

if _, err := nc.QueueSubscribe("updates", "workers", func(m *nats.Msg) {
    wg.Done()
}); err != nil {
    log.Fatal(err)
}

wg.Wait()
```

```java
Connection nc = Nats.connect("nats://demo.nats.io:4222");
CountDownLatch latch = new CountDownLatch(10);
Dispatcher d = nc.createDispatcher((msg) -> {
    String str = new String(msg.getData(), StandardCharsets.UTF_8);
    System.out.println(str);
    latch.countDown();
});
d.subscribe("updates", "workers");
latch.await(); 
nc.close();
```

```javascript
nc.subscribe(subj, {
    queue: "workers",
    callback: (_err, _msg) => {
      t.log("worker1 got message");
    },
});

nc.subscribe(subj, {
    queue: "workers",
    callback: (_err, _msg) => {
      t.log("worker2 got message");
    },
});
```

```python
nc = NATS()
await nc.connect(servers=["nats://demo.nats.io:4222"])
future = asyncio.Future()

async def cb(msg):
  nonlocal future
  future.set_result(msg)

await nc.subscribe("updates", queue="workers", cb=cb)
await nc.publish("updates", b'All is Well')

msg = await asyncio.wait_for(future, 1)
print("Msg", msg)
```

```csharp
using NATS.Net;
await using var client = new NatsClient();
var count = 0;
await foreach (var msg in client.SubscribeAsync<string>(subject: "updates", queueGroup: "workers"))
{
    Console.WriteLine($"Received {++count}: {msg.Subject}: {msg.Data}");
    if (count == 10)
    {
        break;
    }
}
Console.WriteLine("Done");
```

```ruby
require 'nats/client'
require 'fiber'

NATS.start(servers:["nats://127.0.0.1:4222"]) do |nc|
  Fiber.new do
    f = Fiber.current
    nc.subscribe("updates", queue: "worker") do |msg, reply|
      f.resume Time.now
    end
    nc.publish("updates", "A")
    msg = Fiber.yield
    puts "Msg: #{msg}"
  end.resume
end
```

```c
static void onMsg(natsConnection *conn, natsSubscription *sub, natsMsg *msg, void *closure)
{
    printf("Received msg: %s - %.*s\n", natsMsg_GetSubject(msg), natsMsg_GetDataLength(msg), natsMsg_GetData(msg));
    natsMsg_Destroy(msg);
}

natsConnection *conn = NULL;
natsSubscription *sub = NULL;
natsStatus s;
s = natsConnection_ConnectTo(&conn, NATS_DEFAULT_URL);
if (s == NATS_OK)
    s = natsConnection_QueueSubscribe(&sub, conn, "updates", "workers", onMsg, NULL);
natsSubscription_Destroy(sub);
natsConnection_Destroy(conn);
```

--------------------------------

### Run NATS High-Load Throughput Benchmark

Source: https://github.com/nats-io/nats.docs/blob/master/using-nats/nats-tools/nats_cli/natsbench.md

Configures multiple service instances and multiple concurrent clients to measure system throughput under load. It utilizes queue groups to distribute requests across service instances.

```bash
nats bench service serve foo --size 16 --clients 2
nats bench service request foo --size 16 --clients 50 --no-progress
```

--------------------------------

### Publish Multiple NATS Messages

Source: https://github.com/nats-io/nats.docs/blob/master/nats-concepts/core-nats/queue-groups/queues_walkthrough.md

Publishes multiple NATS messages to the 'foo' subject with a specified count. This demonstrates the load balancing behavior across queue group members.

```bash
nats request foo --count 10 "Request {{Count}}"
```

--------------------------------

### TypeScript: Subscribe with Queue Group

Source: https://github.com/nats-io/nats.docs/blob/master/_examples/drain_sub.html

This snippet demonstrates subscribing to a subject with a queue group in TypeScript. While not directly related to draining, queue groups are often used in conjunction with message processing patterns where draining might be applied during shutdown.

```TypeScript
let sub = await nc.subscribe('updates', (err, msg) => {
    t.log('worker got message', msg.data);
}, {queue: "workers"});
```

--------------------------------

### Start First Queue Group Member

Source: https://github.com/nats-io/nats.docs/blob/master/nats-concepts/core-nats/queue-groups/queues_walkthrough.md

Starts a NATS subscriber that automatically joins a queue group. This is the first instance of a service in the group.

```bash
nats reply foo "service instance A Reply# {{Count}}"
```

--------------------------------

### Java JetStream Push Subscribe with Durable Consumer and Queue

Source: https://github.com/nats-io/nats.docs/blob/master/using-nats/developing-with-nats/js/consumers.md

This Java example demonstrates setting up a JetStream push subscriber that uses a durable consumer and joins a queue group for load balancing.

```java
package io.nats.examples.jetstream;

import io.nats.client.*;
import io.nats.client.api.PublishAck;
import io.nats.examples.ExampleArgs;
import io.nats.examples.ExampleUtils;

import java.io.IOException;
import java.nio.charset.StandardCharsets;
import java.time.Duration;
import java.util.ArrayList;
import java.util.List;
import java.util.concurrent.atomic.AtomicInteger;

import static io.nats.examples.jetstream.NatsJsUtils.createStreamExitWhenExists;

/**
 * This example will demonstrate JetStream push subscribing using a durable consumer and a queue
 */
public class NatsJsPushSubQueueDurable {
    static final String usageString = 
        "\nUsage: java -cp <classpath> NatsJsPushSubQueueDurable [-s server] [-strm stream] [-sub subject] [-q queue] [-dur durable] [-mcnt msgCount] [-scnt subCount]" 
            + "\n\nDefault Values:" 
            + "\n   [-strm stream]   qdur-stream" 
            + "\n   [-sub subject]   qdur-subject" 
            + "\n   [-q queue]       qdur-queue" 
            + "\n   [-dur durable]   qdur-durable" 
            + "\n   [-mcnt msgCount] 100" 
            + "\n   [-scnt subCount] 5" 
            + "\n\nUse tls:// or opentls:// to require tls, via the Default SSLContext\n" 
            + "\nSet the environment variable NATS_NKEY to use challenge response authentication by setting a file containing your private key.\n" 
            + "\nSet the environment variable NATS_CREDS to use JWT/NKey authentication by setting a file containing your user creds.\n" 
            + "\nUse the URL in the -s server parameter for user/pass/token authentication.\n";

    public static void main(String[] args) {
        ExampleArgs exArgs = ExampleArgs.builder("Push Subscribe, Durable Consumer, Queue", args, usageString)
                .defaultStream("qdur-stream")
                .defaultSubject("qdur-subject")
                .defaultQueue("qdur-queue")
                .defaultDurable("qdur-durable")
                .defaultMsgCount(100)
                .defaultSubCount(5)
                .build();

        try (Connection nc = Nats.connect(ExampleUtils.createExampleOptions(exArgs.server, true))) {

            // Create a JetStreamManagement context.
            JetStreamManagement jsm = nc.jetStreamManagement();

            // Use the utility to create a stream stored in memory.
            createStreamExitWhenExists(jsm, exArgs.stream, exArgs.subject);

            // Create our JetStream context
            JetStream js = nc.jetStream();

```

--------------------------------

### Queue Subscribe

Source: https://github.com/nats-io/nats.docs/blob/master/using-nats/developing-with-nats/receiving/queues.md

Establishes a subscription to a subject within a specific queue group to enable load-balanced message delivery.

```APIDOC
## Queue Subscribe

### Description
Subscribes to a subject as part of a queue group. The NATS server load balances messages among all subscribers sharing the same queue name.

### Method
SUBSCRIBE (NATS Protocol)

### Parameters
#### Path Parameters
- **subject** (string) - Required - The subject to listen for messages on.
- **queue** (string) - Required - The name of the queue group to join.

### Request Example
// Example using Go
nc.QueueSubscribe("updates", "workers", func(m *nats.Msg) { ... });

### Response
#### Success Response (200)
- **message** (object) - The received message payload and metadata.

#### Response Example
{
  "subject": "updates",
  "data": "All is Well"
}
```

--------------------------------

### Start Third Queue Group Member

Source: https://github.com/nats-io/nats.docs/blob/master/nats-concepts/core-nats/queue-groups/queues_walkthrough.md

Starts a third NATS subscriber, further expanding the queue group. This provides more redundancy and distribution for incoming messages.

```bash
nats reply foo "service instance C Reply# {{Count}}"
```

--------------------------------

### Define Queue Subscription Permissions

Source: https://github.com/nats-io/nats.docs/blob/master/running-a-nats-service/configuration/securing_nats/authorization.md

This snippet shows how to restrict queue subscriptions using allow and deny lists. It demonstrates how to permit specific queue groups while preventing access to sensitive production groups.

```hcl
users = [
  {
    user: "a", password: "a", permissions: {
      sub: {
        allow: ["foo queue"]
     }
  }
  {
    user: "b", password: "b", permissions: {
      sub: {
        allow: ["foo", "foo v1", "foo v1.>", "foo *.dev"],
        deny: ["> *.prod"]
     }
  }
]
```

--------------------------------

### NATS Queue Subscribe - Go

Source: https://github.com/nats-io/nats.docs/blob/master/_examples/subscribe_queue.html

Subscribes to the 'updates' subject with a queue group 'worker' in Go. It waits for 10 messages using a WaitGroup before closing the connection. This ensures that messages are distributed among multiple subscribers in the same queue group.

```Go
nc, err := nats.Connect("demo.nats.io")
if err != nil {
	log.Fatal(err)
}
deferr nc.Close()

wg := sync.WaitGroup{}
wg.Add(10)

if _, err := nc.QueueSubscribe("updates", "worker", func(m *nats.Msg) {
	wg.Done()
}); err != nil {
	log.Fatal(err)
}

wg.Wait()
```

--------------------------------

### NATS Queue Subscribe - Node.js

Source: https://github.com/nats-io/nats.docs/blob/master/_examples/subscribe_queue.html

Shows how to subscribe to a subject with a queue group using the NATS client library for Node.js. It connects to the NATS server and sets up a subscription for 'updates' in the 'workers' queue group.

```JavaScript
let nc = NATS.connect({
    url: "nats://demo.nats.io:4222"});

nc.subscribe('updates', {queue: "workers"}, (msg) => {
    t.log('worker got message', msg);
});
```

--------------------------------

### NATS Queue Subscribe - Java

Source: https://github.com/nats-io/nats.docs/blob/master/_examples/subscribe_queue.html

Demonstrates creating a queue subscription in Java. It connects to the NATS server, sets up a dispatcher to handle incoming messages on the 'updates' subject with the queue group 'workers', and uses a CountDownLatch to wait for 10 messages.

```Java
Connection nc = Nats.connect("nats://demo.nats.io:4222");

CountDownLatch latch = new CountDownLatch(10);

Dispatcher d = nc.createDispatcher((msg) -> {
    String str = new String(msg.getData(), StandardCharsets.UTF_8);
    System.out.println(str);
    latch.countDown();
});

d.subscribe("updates", "workers");

latch.await();

nc.close();
```

--------------------------------

### Publish NATS Message to Queue Group

Source: https://github.com/nats-io/nats.docs/blob/master/nats-concepts/core-nats/queue-groups/queues_walkthrough.md

Publishes a single NATS message to the 'foo' subject. One of the subscribers in the queue group will receive and process this message.

```bash
nats request foo "Simple request"
```

--------------------------------

### NATS Queue Subscribe - TypeScript

Source: https://github.com/nats-io/nats.docs/blob/master/_examples/subscribe_queue.html

A TypeScript example for creating a queue subscription. It connects to the NATS server and subscribes to the 'updates' subject within the 'workers' queue group, logging received messages.

```TypeScript
await nc.subscribe('updates', (err, msg) => {
    t.log('worker got message', msg.data);
}, {queue: "workers"});
```

--------------------------------

### NATS Queue Subscribe - Python

Source: https://github.com/nats-io/nats.docs/blob/master/_examples/subscribe_queue.html

An asynchronous example in Python for creating a queue subscription. It connects to the NATS server, defines a callback function to handle messages, subscribes to 'updates' with the queue 'workers', and publishes a test message.

```Python
nc = NATS()

await nc.connect(servers=["nats://demo.nats.io:4222"])

future = asyncio.Future()

async def cb(msg):
  nonlocal future
  future.set_result(msg)

await nc.subscribe("updates", queue="workers", cb=cb)
await nc.publish("updates", b'All is Well')

msg = await asyncio.wait_for(future, 1)
print("Msg", msg)
```

--------------------------------

### NATS Scatter-Gather with Response Count and Timeout (Go)

Source: https://github.com/nats-io/nats.docs/blob/master/using-nats/developing-with-nats/sending/request_reply.md

Illustrates a scatter-gather approach in Go that aims to receive at least a minimum number of responses within a given time limit. It subscribes, sends a request, and collects messages until either the minimum response count is met or the timeout expires. This pattern is useful when a certain number of acknowledgments are required.

```Go
sub, err := nc.SubscribeSync(replyTo)
if err != nil {
    log.Fatal(err)
}
nc.Flush()

// Send the request
nc.PublishRequest(subject, replyTo, []byte(input))

// Wait for a single response
max := 500 * time.Millisecond
start := time.Now()
for time.Now().Sub(start) < max {
    msg, err := sub.NextMsg(1 * time.Second)
    if err != nil {
        break
    }

    responses = append(responses, string(msg.Data))

    if len(responses) >= minResponses {
        break
    }
}
sub.Unsubscribe()
```

--------------------------------

### Configure Pending Message Limits in NATS Clients

Source: https://github.com/nats-io/nats.docs/blob/master/using-nats/developing-with-nats/events/slow.md

This snippet demonstrates how to set limits on the incoming message queue (by count and by bytes) for different NATS client libraries. These settings help manage bursty traffic and prevent slow consumer errors.

```go
nc, err := nats.Connect("demo.nats.io")
if err != nil {
    log.Fatal(err)
}
defer nc.Close()

sub1, err := nc.Subscribe("updates", func(m *nats.Msg) {})
if err != nil {
    log.Fatal(err)
}

sub1.SetPendingLimits(1000, 5*1024*1024)

sub2, err := nc.Subscribe("updates", func(m *nats.Msg) {})
if err != nil {
    log.Fatal(err)
}

sub2.SetPendingLimits(-1, -1)

nc.Close()
```

```java
Connection nc = Nats.connect("nats://demo.nats.io:4222");

Dispatcher d = nc.createDispatcher((msg) -> {});
d.subscribe("updates");
d.setPendingLimits(1_000, 5 * 1024 * 1024);

Subscription sub = nc.subscribe("updates");
sub.setPendingLimits(1_000, 5 * 1024 * 1024);

nc.close();
```

```python
nc = NATS()
await nc.connect(servers=["nats://demo.nats.io:4222"])

future = asyncio.Future()

async def cb(msg):
  nonlocal future
  future.set_result(msg)

await nc.subscribe("updates", cb=cb, pending_bytes_limit=5*1024*1024, pending_msgs_limit=1000)
```

```csharp
await using var client = new NatsClient();

var subOpts = new NatsSubOpts
{
    ChannelOpts = new NatsSubChannelOpts
    {
        Capacity = 1000,
        FullMode = BoundedChannelFullMode.DropOldest
    }
};
await foreach (var msg in client.SubscribeAsync<string>(subject: "updates", opts: subOpts))
{
    Console.WriteLine($"Received: {msg.Subject}: {msg.Data}");    
}
```

--------------------------------

### NATS Scatter-Gather with Time-Based Timeout (Go)

Source: https://github.com/nats-io/nats.docs/blob/master/using-nats/developing-with-nats/sending/request_reply.md

Demonstrates a scatter-gather implementation in Go where the loop for receiving messages is limited by time. It sends a request and collects responses until a specified timeout is reached or an error occurs. Dependencies include the NATS client library and the time package.

```Go
sub, err := nc.SubscribeSync(replyTo)
if err != nil {
    log.Fatal(err)
}
nc.Flush()

// Send the request
nc.PublishRequest(subject, replyTo, []byte(input))

// Wait for a single response
max := 100 * time.Millisecond
start := time.Now()
for time.Now().Sub(start) < max {
    msg, err := sub.NextMsg(1 * time.Second)
    if err != nil {
        break
    }

    responses = append(responses, string(msg.Data))
}
sub.Unsubscribe()
```

--------------------------------

### NATS Combined Traffic Splitting and Loss Simulation

Source: https://github.com/nats-io/nats.docs/blob/master/running-a-nats-service/configuration/configuring_subject_mapping.md

Combine weighted mappings to split traffic between service versions, simulate failure conditions, and introduce message loss. This example sends 90% to the primary service, 8% to a failure simulation, and the remaining 2% is lost.

```text
myservice.requests: [
  { destination: myservice.requests.v3, weight: 90% },
  { destination: myservice.requests.v3.fail, weight: 8% }
  # the remaining 2% is "lost"
]
```