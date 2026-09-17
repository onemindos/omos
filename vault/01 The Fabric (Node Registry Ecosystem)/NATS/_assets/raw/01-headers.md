### NATS HPUB Message Description

Source: https://github.com/nats-io/nats.docs/blob/master/reference/nats-protocol/nats-protocol/README.md

The HPUB message extends the PUB message by allowing the inclusion of NATS headers along with the message payload. Headers are structured as 'name: value' pairs and support multi-value headers. The total message size calculation must account for both headers and payload. This is useful for adding metadata to messages.

```text
HPUB <subject> [reply-to] <#bytes>
[headers]
[payload]

```

--------------------------------

### HMSG Protocol Message

Source: https://github.com/nats-io/nats.docs/blob/master/reference/nats-protocol/nats-protocol/README.md

The HMSG message extends the standard MSG protocol by including headers alongside the message payload.

```APIDOC
## HMSG

### Description
The HMSG message is used to deliver messages that include headers. It follows the format: HMSG <subject> <sid> [reply-to] <#header bytes> <#total bytes>\r\n[headers]\r\n\r\n[payload]\r\n

### Parameters
- **subject** (string) - Required - Subject name this message was received on.
- **sid** (string) - Required - The unique alphanumeric subscription ID.
- **reply-to** (string) - Optional - The subject on which the publisher is listening for responses.
- **#header bytes** (int) - Required - The size of the headers section in bytes.
- **#total bytes** (int) - Required - The total size of headers and payload sections in bytes.

### Request Example
HMSG FOO.BAR 34 45\r\nNATS/1.0\r\nFoodGroup: vegetable\r\n\r\nHello World\r\n
```

--------------------------------

### HPUB Message

Source: https://github.com/nats-io/nats.docs/blob/master/reference/nats-protocol/nats-protocol/README.md

The HPUB message extends the PUB functionality by allowing the inclusion of NATS headers in the message payload.

```APIDOC
## HPUB

### Description
The `HPUB` message is the same as `PUB` but extends the message payload to include NATS headers. Headers are structured as name: value pairs.

### Method
N/A (Protocol Message)

### Endpoint
HPUB <subject> [reply-to] <header_bytes> <total_bytes>\r\n[headers]\r\n[payload]\r\n

### Parameters
#### Path Parameters
- **subject** (string) - Required - The destination subject.
- **reply-to** (string) - Optional - The reply subject.
- **header_bytes** (int) - Required - Size of headers in bytes.
- **total_bytes** (int) - Required - Total size of headers plus payload.

### Request Example
HPUB FOO 10 20\r\nNats-Msg-Id: 1\r\n\r\nHello NATS!\r\n

### Response
#### Success Response (200)
- **N/A** - Asynchronous protocol message.
```

--------------------------------

### Publish Message (HPUB) - NATS Protocol

Source: https://github.com/nats-io/nats.docs/blob/master/reference/nats-protocol/nats-protocol/README.md

The HPUB command is used to publish messages to a specified subject. It supports optional reply subjects, headers, and payloads. The command requires specifying the byte sizes for headers and the total message.

```nats
HPUB <subject> [reply-to] <#header bytes> <#total bytes>\r\n[headers]\r\n\r\n[payload]\r\n
```

```nats
HPUB FOO 22 33\r\nNATS/1.0\r\nBar: Baz\r\n\r\nHello NATS!\r\n
```

```nats
HPUB FRONT.DOOR JOKE.22 45 56\r\nNATS/1.0\r\nBREAKFAST: donut\r\nLUNCH: burger\r\n\r\nKnock Knock\r\n
```

```nats
HPUB NOTIFY 22 22\r\nNATS/1.0\r\nBar: Baz\r\n\r\n\r\n
```

```nats
HPUB MORNING.MENU 47 51\r\nNATS/1.0\r\nBREAKFAST: donut\r\nBREAKFAST: eggs\r\n\r\nYum!\r\n
```

--------------------------------

### Configure Message Tracing for Account

Source: https://github.com/nats-io/nats.docs/blob/master/running-a-nats-service/configuration/README.md

Enables tracing for messages with the 'traceparent' header. Configure the destination subject and sampling rate.

```hcl
accounts {
  A {
      msg_trace {
        dest: "trace-ot.accountA.samples"
        sampling: 5
      }
    }
}

```

--------------------------------

### Received Message Output

Source: https://github.com/nats-io/nats.docs/blob/master/running-a-nats-service/configuration/gateways/gateway.md

Example output showing a message received by a subscriber, indicating successful message forwarding through the gateway.

```text
[#1] Received on [foo] : 'bar'
```

--------------------------------

### Java Message Consumer Example

Source: https://github.com/nats-io/nats.docs/blob/master/using-nats/developing-with-nats/js/consumers.md

Demonstrates a simplified message consumer in Java using JetStream. It sets up a stream, publishes messages, and consumes them with a handler that acknowledges messages and stops after a certain count.

```Java
package io.nats.examples.jetstream.simple;

import io.nats.client.*;
import io.nats.client.api.ConsumerConfiguration;
import io.nats.examples.jetstream.ResilientPublisher;
import java.io.IOException;
import java.util.concurrent.CountDownLatch;
import java.util.concurrent.atomic.AtomicInteger;
import static io.nats.examples.jetstream.NatsJsUtils.createOrReplaceStream;

/**
* This example will demonstrate simplified consume with a handler
*/
public class MessageConsumerExample {
 private static final String STREAM = "consume-stream";
 private static final String SUBJECT = "consume-subject";
 private static final String CONSUMER_NAME = "consume-consumer";
 private static final String MESSAGE_PREFIX = "consume";
 private static final int STOP_COUNT = 500;
 private static final int REPORT_EVERY = 100;
 private static final String SERVER = "nats://localhost:4222";

 public static void main(String[] args) {
     Options options = Options.builder().server(SERVER).build();
     try (Connection nc = Nats.connect(options)) {
         JetStreamManagement jsm = nc.jetStreamManagement();
         createOrReplaceStream(jsm, STREAM, SUBJECT);

         //Utility for filling the stream with some messages
         System.out.println("Starting publish...");
         ResilientPublisher publisher = new ResilientPublisher(nc, jsm, STREAM, SUBJECT).basicDataPrefix(MESSAGE_PREFIX).jitter(10);
         Thread pubThread = new Thread(publisher);
         pubThread.start();

         // get stream context, create consumer and get the consumer context
         StreamContext streamContext;
         ConsumerContext consumerContext;
         CountDownLatch latch = new CountDownLatch(1);
         AtomicInteger atomicCount = new AtomicInteger();
         long start = System.nanoTime();

         streamContext = nc.getStreamContext(STREAM);
         streamContext.createOrUpdateConsumer(ConsumerConfiguration.builder().durable(CONSUMER_NAME).build());
         consumerContext = streamContext.getConsumerContext(CONSUMER_NAME);

         MessageHandler handler = msg -> {
             msg.ack();
             int count = atomicCount.incrementAndGet();
             if (count % REPORT_EVERY == 0) {
            	 System.out.println("Handler" + ": Received " + count + " messages in " + (System.nanoTime() - start) / 1_000_000 + "ms.");
             }
             if (count == STOP_COUNT) {
                 latch.countDown();
             }
         };

     	 // create the consumer and install handler
     	 MessageConsumer consumer = consumerContext.consume(handler);
     	 //Waiting for the handler signalling us to stop
         latch.await();
         // When stop is called, no more pull requests will be made, but messages already requested
         // will still come across the wire to the client.
         System.out.println("Stopping the consumer...");
         consumer.stop();
         // wait until the consumer is finished processing backlog
         while (!consumer.isFinished()) {
             Thread.sleep(10);
         }
         System.out.println("Final" + ": Received " + atomicCount.get() + " messages in " + (System.nanoTime() - start) / 1_000_000 + "ms.");

         publisher.stop(); // otherwise the ConsumerContext background thread will complain when the connection goes away
         pubThread.join();
     }
     catch (JetStreamApiException | IOException e) {
         // JetStreamApiException:
         //      1. the stream or consumer did not exist
         //      2. api calls under the covers theoretically this could fail, but practically it won't.
         // IOException:

```

--------------------------------

### Publish Message Through Gateway

Source: https://github.com/nats-io/nats.docs/blob/master/running-a-nats-service/configuration/gateways/gateway.md

Command to publish a message to a subject on a specific port, demonstrating message forwarding via a gateway.

```shell
nats pub -s localhost:4444 foo bar
```

--------------------------------

### PUB Message

Source: https://github.com/nats-io/nats.docs/blob/master/reference/nats-protocol/nats-protocol/README.md

The PUB message is used to publish a message payload to a specific subject, with an optional reply subject.

```APIDOC
## PUB

### Description
The `PUB` message publishes the message payload to the given subject name, optionally supplying a reply subject.

### Method
N/A (Protocol Message)

### Endpoint
PUB <subject> [reply-to] <#bytes>\r\n[payload]\r\n

### Parameters
#### Path Parameters
- **subject** (string) - Required - The destination subject to publish to.
- **reply-to** (string) - Optional - The reply subject that subscribers can use to send a response.
- **#bytes** (int) - Required - The payload size in bytes.
- **payload** (string) - Optional - The message payload data.

### Request Example
PUB FOO 11\r\nHello NATS!\r\n

### Response
#### Success Response (200)
- **N/A** - NATS protocol messages are asynchronous; success is implied unless an ERR is received.
```

--------------------------------

### Publish Messages with Delay

Source: https://github.com/nats-io/nats.docs/blob/master/using-nats/jetstream/model_deep_dive.md

This script publishes messages with a specified delay between each, simulating original message arrival times for testing replay functionality.

```shell
for i in 1 2 3                                                                                                                                                      <15:15:35
do
  nats pub ORDERS.processed "order ${i}"
  sleep 10
done
```

--------------------------------

### Fetch Messages with Time-Based Consumer

Source: https://github.com/nats-io/nats.docs/blob/master/using-nats/jetstream/model_deep_dive.md

Retrieves messages using a time-based consumer configured to start from 2 minutes ago. This example shows receiving a message published after that time.

```shell
nats consumer next ORDERS 2MIN
```

--------------------------------

### Publishing with Message Deduplication

Source: https://github.com/nats-io/nats.docs/blob/master/using-nats/jetstream/model_deep_dive.md

Use the `Nats-Msg-Id` header to ensure idempotent message writes. JetStream will ignore duplicate messages with the same ID within a configured window.

```shell
nats req -H Nats-Msg-Id:1 ORDERS.new hello1
nats req -H Nats-Msg-Id:1 ORDERS.new hello2
nats req -H Nats-Msg-Id:1 ORDERS.new hello3
nats req -H Nats-Msg-Id:1 ORDERS.new hello4
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

### HPUB Protocol

Source: https://github.com/nats-io/nats.docs/blob/master/reference/nats-protocol/nats-protocol/README.md

The HPUB command is used to publish messages with headers to a specific subject.

```APIDOC
## HPUB <subject> [reply-to] <#header bytes> <#total bytes>

### Description
Publishes a message with headers to a destination subject.

### Parameters
- **subject** (string) - Required - The destination subject.
- **reply-to** (string) - Optional - The reply subject for responses.
- **#header bytes** (int) - Required - Size of the headers section.
- **#total bytes** (int) - Required - Total size of headers and payload.

### Request Example
HPUB FOO 22 33\r\nNATS/1.0\r\nBar: Baz\r\n\r\nHello NATS!\r\n
```

--------------------------------

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

### Subscribe to Gateway Messages

Source: https://github.com/nats-io/nats.docs/blob/master/running-a-nats-service/configuration/gateways/gateway.md

Command to subscribe to messages on a specific port, demonstrating message reception through a gateway.

```shell
nats sub -s localhost:4333 ">"
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

### Fetch Messages with DeliverAll Consumer

Source: https://github.com/nats-io/nats.docs/blob/master/using-nats/jetstream/model_deep_dive.md

Retrieves messages from a DeliverAll consumer. This example shows receiving and acknowledging a message.

```shell
nats consumer next ORDERS ALL
```

--------------------------------

### Configure Cross-Account Message Imports and Exports

Source: https://github.com/nats-io/nats.docs/blob/master/running-a-nats-service/nats_admin/jwt.md

Configuration snippet demonstrating how to allow message flow between accounts using imports and exports.

```text
accounts: {
    A: {
        users: [{user: a, password: a}]
        imports: [{stream: {account: B, subject: "foo"}}]
    },
    B: {
        users: [{user: b, password: b}]
        exports: [{stream: "foo"}]
    },
}
```

--------------------------------

### C Async Message Callback

Source: https://github.com/nats-io/nats.docs/blob/master/using-nats/developing-with-nats/js/consumers.md

Illustrates the asynchronous callback function for receiving messages in C. It includes basic message printing and a mechanism to track the total number of received messages for performance measurement.

```c
#include "examples.h"

static const char *usage = "\
-gd            use global message delivery thread pool\
-sync          receive synchronously (default is asynchronous)\
-pull          use pull subscription\
-fc            enable flow control\
-count         number of expected messages\
";

static void
onMsg(natsConnection *nc, natsSubscription *sub, natsMsg *msg, void *closure)
{
    if (print)
        printf("Received msg: %s - %.*s\n",
               natsMsg_GetSubject(msg),
               natsMsg_GetDataLength(msg),
               natsMsg_GetData(msg));

    if (start == 0)
        start = nats_Now();

    // We should be using a mutex to protect those variables since
    // they are used from the subscription's delivery and the main
    // threads. For demo purposes, this is fine.
    if (++count == total)
        elapsed = nats_Now() - start;

    // Since this is auto-ack callback, we don't need to ack here.
    natsMsg_Destroy(msg);
}

static void
asyncCb(natsConnection *nc, natsSubscription *sub, natsStatus err, void *closure)
{

```

--------------------------------

### Publish String Message to 'updates' Subject

Source: https://github.com/nats-io/nats.docs/blob/master/using-nats/developing-with-nats/sending/README.md

Demonstrates sending a UTF-8 encoded string message to a specified subject using different NATS client libraries. Ensure the NATS server is running and accessible.

```Go
nc, err := nats.Connect("demo.nats.io", nats.Name("API PublishBytes Example"))
if err != nil {
    log.Fatal(err)
}
deferr nc.Close()

if err := nc.Publish("updates", []byte("All is Well")); err != nil {
    log.Fatal(err)
}
```

```Java
Connection nc = Nats.connect("nats://demo.nats.io:4222");

nc.publish("updates", "All is Well".getBytes(StandardCharsets.UTF_8));
```

```JavaScript
const sc = StringCodec();
nc.publish("updates", sc.encode("All is Well"));
```

```Python
nc = NATS()

await nc.connect(servers=["nats://demo.nats.io:4222"])

await nc.publish("updates", b'All is Well')
```

```C#
// dotnet add package NATS.Net
using NATS.Net;

await using var client = new NatsClient(url: "demo.nats.io", name: "API Publish String Example");

// The default serializer uses UTF-8 encoding for strings
await client.PublishAsync<string>(subject: "updates", data: "All is Well");
```

```Ruby
require 'nats/client'

NATS.start(servers:["nats://127.0.0.1:4222"]) do |nc|
  nc.publish("updates", "All is Well")
end
```

--------------------------------

### Fetch Next Message (Single)

Source: https://github.com/nats-io/nats.docs/blob/master/using-nats/jetstream/nats_api_reference.md

Use this to request a single message from a pull-based consumer. The request is sent to a specific API subject.

```shell
nats req '$JS.API.CONSUMER.MSG.NEXT.ORDERS.test' '1'
```

--------------------------------

### Fetch Messages with DeliverLast Consumer

Source: https://github.com/nats-io/nats.docs/blob/master/using-nats/jetstream/model_deep_dive.md

Retrieves the latest message using a DeliverLast consumer. This example shows receiving and acknowledging the message.

```shell
nats consumer next ORDERS LAST
```