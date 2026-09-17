### Publish and Subscribe in NATS Cluster

Source: https://github.com/nats-io/nats.docs/blob/master/running-a-nats-service/running/nats_docker/README.md

Commands to subscribe to a subject on one node and publish messages to different cluster nodes to verify cross-node communication.

```shell
nats sub -s nats://nats:4222 hello &
nats pub -s "nats://nats-1:4222" hello first
nats pub -s "nats://nats-2:4222" hello second
```

--------------------------------

### Publish Request and Subscribe to Reply (Go)

Source: https://github.com/nats-io/nats.docs/blob/master/_examples/publish_with_reply.html

This Go code snippet demonstrates publishing a request to the 'time' subject and subscribing to a unique inbox subject to receive the reply. It handles connection, subscription, publishing, and reading the response. Errors are logged and the connection is deferred closed.

```Go
nc, err := nats.Connect("demo.nats.io")
if err != nil {
	log.Fatal(err)
}
deferr nc.Close()

// Create a unique subject name for replies.
uniqueReplyTo := nats.NewInbox()

// Listen for a single response
sub, err := nc.SubscribeSync(uniqueReplyTo)
if err != nil {
	log.Fatal(err)
}

// Send the request.
// If processing is synchronous, use Request() which returns the response message.
if err := nc.PublishRequest("time", uniqueReplyTo, nil); err != nil {
	log.Fatal(err)
}

// Read the reply
msg, err := sub.NextMsg(time.Second)
if err != nil {
	log.Fatal(err)
}

// Use the response
log.Printf("Reply: %s", msg.Data)
```

--------------------------------

### Subscribe and Publish in Ruby

Source: https://github.com/nats-io/nats.docs/blob/master/_examples/subscribe_star.html

This Ruby code snippet shows how to subscribe to a wildcard subject ('time.*.east') and publish messages to test the subscription. It uses Fibers for concurrency and the 'nats-client' gem. It yields to wait for published messages.

```Ruby
require 'nats/client'
require 'fiber'

NATS.start(servers:["nats://127.0.0.1:4222"]) do |nc|
  Fiber.new do
    f = Fiber.current
    nc.subscribe("time.*.east") do |msg, reply|
      f.resume Time.now
    end

    nc.publish("time.A.east", "A")
    nc.publish("time.B.east", "B")

    # Use the response
    msg_A = Fiber.yield
    puts "Msg A: #{msg_A}"

    msg_B = Fiber.yield
    puts "Msg B: #{msg_B}"
  end.resume
end
```

--------------------------------

### Publish and Subscribe with Structured Data in C#

Source: https://github.com/nats-io/nats.docs/blob/master/using-nats/developing-with-nats/receiving/structure.md

Demonstrates how to use the NATS .NET client to subscribe to typed messages and publish various data types including integers, strings, byte arrays, and custom records serialized as JSON.

```csharp
Task.Run(async () => {
    await foreach (var msg in nc.SubscribeAsync<MyData>("x.json", cancellationToken: cts.Token)) {
        Console.WriteLine($"Received data: {msg.Data}");
    }
});
await Task.Delay(1000);
await nc.PublishAsync<int>("x.int", 100);
await nc.PublishAsync<string>("x.string", "Hello, World!");
await nc.PublishAsync<byte[]>("x.bytes", new byte[] { 0x41, 0x42, 0x43 });
await nc.PublishAsync<MyData>("x.json", new MyData(30, "bar"));
await cts.CancelAsync();
public record MyData(int Id, string Name);
```

--------------------------------

### Subscribe with Reply - Python

Source: https://github.com/nats-io/nats.docs/blob/master/_examples/subscribe_w_reply.html

This Python code snippet illustrates how to subscribe to a 'time' subject and respond with the current time. It uses the asyncio NATS client to connect, subscribe with a callback, publish a request, and then publish the response.

```Python
nc = NATS()

await nc.connect(servers=["nats://demo.nats.io:4222"])

future = asyncio.Future()

async def cb(msg):
  nonlocal future
  future.set_result(msg)

await nc.subscribe("time", cb=cb)

await nc.publish_request("time", new_inbox(), b'What is the time?')
await nc.flush()

# Read the message
msg = await asyncio.wait_for(future, 1)

# Send the time
time_as_bytes = "{}".format(datetime.now()).encode()
await nc.publish(msg.reply, time_as_bytes)
```

--------------------------------

### Java JetStream Publisher and Subscriber Example

Source: https://github.com/nats-io/nats.docs/blob/master/using-nats/developing-with-nats/js/consumers.md

This Java code demonstrates setting up multiple push subscribers and a publisher for a NATS JetStream stream. It includes managing threads, tracking received messages, and reporting results.

```Java
System.out.println();

            // Setup the subscribers
            // - the PushSubscribeOptions can be re-used since all the subscribers are the same
            // - use a concurrent integer to track all the messages received
            // - have a list of subscribers and threads so I can track them
            PushSubscribeOptions pso = PushSubscribeOptions.builder().durable(exArgs.durable).build();
            AtomicInteger allReceived = new AtomicInteger();
            List<JsQueueSubscriber> subscribers = new ArrayList<>();
            List<Thread> subThreads = new ArrayList<>();
            for (int id = 1; id <= exArgs.subCount; id++) {
                // setup the subscription
                JetStreamSubscription sub = js.subscribe(exArgs.subject, exArgs.queue, pso);
                // create and track the runnable
                JsQueueSubscriber qs = new JsQueueSubscriber(id, exArgs, js, sub, allReceived);
                subscribers.add(qs);
                // create, track and start the thread
                Thread t = new Thread(qs);
                subThreads.add(t);
                t.start();
            }
            nc.flush(Duration.ofSeconds(1)); // flush outgoing communication with/to the server

            // create and start the publishing
            Thread pubThread = new Thread(new JsPublisher(js, exArgs));
            pubThread.start();

            // wait for all threads to finish
            pubThread.join();
            for (Thread t : subThreads) {
                t.join();
            }

            // report
            for (JsQueueSubscriber qs : subscribers) {
                qs.report();
            }

            System.out.println();

            // delete the stream since we are done with it.
            jsm.deleteStream(exArgs.stream);
        }
        catch (Exception e) {
            e.printStackTrace();
        }
    }

    static class JsPublisher implements Runnable {
        JetStream js;
        ExampleArgs exArgs;

        public JsPublisher(JetStream js, ExampleArgs exArgs) {
            this.js = js;
            this.exArgs = exArgs;
        }

        @Override
        public void run() {
            for (int x = 1; x <= exArgs.msgCount; x++) {
                try {
                    PublishAck pa = js.publish(exArgs.subject, ("Data # " + x).getBytes(StandardCharsets.US_ASCII));
                } catch (IOException | JetStreamApiException e) {
                    // something pretty wrong here
                    e.printStackTrace();
                    System.exit(-1);
                }
            }
        }
    }

    static class JsQueueSubscriber implements Runnable {
        int id;
        int thisReceived;
        List<String> datas;

        ExampleArgs exArgs;
        JetStream js;
        JetStreamSubscription sub;
        AtomicInteger allReceived;

        public JsQueueSubscriber(int id, ExampleArgs exArgs, JetStream js, JetStreamSubscription sub, AtomicInteger allReceived) {
            this.id = id;
            thisReceived = 0;
            datas = new ArrayList<>();
            this.exArgs = exArgs;
            this.js = js;
            this.sub = sub;
            this.allReceived = allReceived;
        }

        public void report() {
            System.out.printf("Sub # %d handled %d messages.\n", id, thisReceived);
        }

        @Override
        public void run() {
            while (allReceived.get() < exArgs.msgCount) {
                try {
                    Message msg = sub.nextMessage(Duration.ofMillis(500));
                    while (msg != null) {
                        thisReceived++;
                        allReceived.incrementAndGet();
                        String data = new String(msg.getData(), StandardCharsets.US_ASCII);
                        datas.add(data);
                        System.out.printf("QS # %d message # %d %s\n", id, thisReceived, data);
                        msg.ack();

                        msg = sub.nextMessage(Duration.ofMillis(500));
                    }
                } catch (InterruptedException e) {
                    // just try again
                }
            }
            System.out.printf("QS # %d completed.\n", id);
        }
    }
}
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

--------------------------------

### Python JetStream Publish and Pull Subscribe

Source: https://github.com/nats-io/nats.docs/blob/master/using-nats/developing-with-nats/js/consumers.md

Shows how to publish messages to a JetStream stream and then create a pull-based consumer to fetch and acknowledge them. Includes basic error handling for connection.

```python
import asyncio

import nats
from nats.errors import TimeoutError

async def main():
    nc = await nats.connect("localhost")

    # Create JetStream context.
    js = nc.jetstream()

    # Persist messages on 'foo's subject.
    await js.add_stream(name="sample-stream", subjects=["foo"])

    for i in range(0, 10):
        ack = await js.publish("foo", f"hello world: {i}".encode())
        print(ack)

    # Create pull based consumer on 'foo'.
    psub = await js.pull_subscribe("foo", "psub")

    # Fetch and ack messagess from consumer.
    for i in range(0, 10):
        msgs = await psub.fetch(1)
        for msg in msgs:
            print(msg)

    await nc.close()

if __name__ == '__main__':
    asyncio.run(main())
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

### Publishing and Subscribing to Updates (C#)

Source: https://github.com/nats-io/nats.docs/blob/master/using-nats/developing-with-nats/receiving/unsub_after.md

This C# snippet demonstrates subscribing to the 'updates' subject with a maximum of 10 messages. It uses the NATS.Net library and shows how to iterate over received messages and print their content. The subscription automatically unsubscribes after the specified limit is reached.

```csharp
// dotnet add package NATS.Net
using NATS.Net;
using NATS.Client.Core;

await using var client = new NatsClient();

// Unsubscribe after 10 messages
var opts = new NatsSubOpts { MaxMsgs = 10 };

var count = 0;

// Subscribe to updates with options
await foreach (var msg in client.SubscribeAsync<string>("updates", opts: opts))
{
    Console.WriteLine($"Received[{++count}]: {msg.Data}");
}

Console.WriteLine("Unsubscribed from updates");
```

--------------------------------

### Publish and Subscribe with NATS Clients

Source: https://github.com/nats-io/nats.docs/blob/master/using-nats/developing-with-nats/receiving/wildcards.md

Basic examples of publishing messages to specific subjects and subscribing to wildcard patterns using various NATS client libraries.

```Python
await nc.publish("time.A.east", b'A')
await nc.publish("time.B.east", b'B')

msg_A = await queue.get()
msg_B = await queue.get()

print("Msg A:", msg_A)
print("Msg B:", msg_B)
```

```C#
using NATS.Net;

await using var client = new NatsClient();

var count = 0;
await foreach (var msg in client.SubscribeAsync<string>("time.*.east"))
{
    Console.WriteLine($"Received {++count}: {msg.Subject}: {msg.Data}");
    if (count == 2) break;
}
Console.WriteLine("Done");
```

```Ruby
NATS.start(servers:["nats://127.0.0.1:4222"]) do |nc|
  Fiber.new do
    f = Fiber.current
    nc.subscribe("time.*.east") { |msg, reply| f.resume Time.now }
    nc.publish("time.A.east", "A")
    nc.publish("time.B.east", "B")
    msg_A = Fiber.yield
    puts "Msg A: #{msg_A}"
  end.resume
end
```

```C
static void onMsg(natsConnection *conn, natsSubscription *sub, natsMsg *msg, void *closure) {
    printf("Received msg: %s - %.*s\n", natsMsg_GetSubject(msg), natsMsg_GetDataLength(msg), natsMsg_GetData(msg));
    natsMsg_Destroy(msg);
}

natsConnection *conn = NULL;
natsSubscription *sub = NULL;
natsConnection_ConnectTo(&conn, NATS_DEFAULT_URL);
natsConnection_Subscribe(&sub, conn, "time.*.east", onMsg, NULL);
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

### NATS CLI: Subscribe and Publish to Single Server

Source: https://github.com/nats-io/nats.docs/blob/master/running-a-nats-service/configuration/clustering/README.md

This example demonstrates how to subscribe to a subject on a specific NATS server and then publish a message to that same server. It verifies basic message flow within a single server instance.

```bash
nats sub -s "nats://127.0.0.1:4222" hello &
nats pub -s "nats://127.0.0.1:4222" hello world_4222
```

--------------------------------

### Subscribe with Reply - TypeScript

Source: https://github.com/nats-io/nats.docs/blob/master/_examples/subscribe_w_reply.html

This TypeScript code snippet demonstrates subscribing to a 'time' subject and publishing the current time to the reply subject if it exists. It uses the NATS TypeScript client to handle subscriptions and message publishing.

```TypeScript
// set up a subscription to process a request
await nc.subscribe('time', (err, msg) => {
    if (msg.reply) {
        nc.publish(msg.reply, new Date().toLocaleTimeString());
    } else {
        t.log('got a request for the time, but no reply subject was set.');
    }
});
```

--------------------------------

### Subscribe with Reply - Node.js

Source: https://github.com/nats-io/nats.docs/blob/master/_examples/subscribe_w_reply.html

This Node.js (JavaScript) code snippet demonstrates subscribing to a 'time' subject and publishing the current time to the reply subject if it exists. It uses the NATS Node.js client to handle subscriptions and message publishing.

```JavaScript
let nc = NATS.connect({
    url: "nats://demo.nats.io:4222"
});

// set up a subscription to process a request
nc.subscribe('time', (msg, reply) => {
    if (msg.reply) {
        nc.publish(msg.reply, new Date().toLocaleTimeString());
    }
});
```

--------------------------------

### TypeScript - Subscribe

Source: https://github.com/nats-io/nats.docs/blob/master/_examples/subscribe_async.html

This TypeScript example demonstrates how to subscribe to a NATS subject and log incoming messages.

```APIDOC
## TypeScript - Subscribe

### Description
This TypeScript example shows how to subscribe to the "updates" subject in NATS and log received messages or any errors that occur.

### Method
N/A (Client-side operation)

### Endpoint
N/A (Assumes a connection to a NATS server has already been established)

### Parameters
(This snippet assumes a `nc` object representing an active NATS connection is available.)

### Request Example
```typescript
nc.subscribe("updates", (err, msg) => {
    if(err) {
        console.log('error', err);
    } else {
        t.log(msg.data);
    }
});
```

### Response
#### Success Response (N/A)
This is a client-side example. The output will be the content of the received NATS message or an error log.

#### Response Example
(No direct response example, as it involves asynchronous message handling)
```

--------------------------------

### Publishing and Subscribing with Response Handling (Ruby)

Source: https://github.com/nats-io/nats.docs/blob/master/using-nats/developing-with-nats/receiving/unsub_after.md

This Ruby example demonstrates publishing a request to the 'time' subject and handling the response. It uses Fibers for asynchronous operations and subscribes to the 'time' subject with a maximum of one message. A second publish to the same subject is included to show it won't be received due to the max limit.

```ruby
require 'nats/client'
require 'fiber'

NATS.start(servers:["nats://127.0.0.1:4222"]) do |nc|
  Fiber.new do
    f = Fiber.current

    nc.subscribe("time", max: 1) do |msg, reply|
      f.resume Time.now
    end

    nc.publish("time", 'What is the time?', NATS.create_inbox)

    # Use the response
    msg = Fiber.yield
    puts "Reply: #{msg}"

    # Won't be received
    nc.publish("time", 'What is the time?', NATS.create_inbox)

  end.resume
end
```