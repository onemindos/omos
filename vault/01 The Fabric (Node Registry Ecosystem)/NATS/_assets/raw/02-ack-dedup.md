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

### Fetch Next Message Without Acknowledging

Source: https://github.com/nats-io/nats.docs/blob/master/using-nats/jetstream/model_deep_dive.md

Fetches the next available message from a NATS JetStream consumer but does not acknowledge it. This leaves the message pending acknowledgement and increases redelivery counts on subsequent fetches.

```shell
nats consumer next ORDERS DISPATCH --no-ack
```

--------------------------------

### Fetch Next Message and Acknowledge

Source: https://github.com/nats-io/nats.docs/blob/master/using-nats/jetstream/model_deep_dive.md

Fetches the next available message from a NATS JetStream consumer and acknowledges its receipt. This updates the consumer's acknowledgement floor.

```shell
nats con next ORDERS DISPATCH
```

--------------------------------

### Fetch Next Message and Acknowledge (CLI)

Source: https://github.com/nats-io/nats.docs/blob/master/using-nats/jetstream/model_deep_dive.md

Fetches the next available message from a NATS JetStream consumer and acknowledges its receipt using the CLI. This updates the consumer's acknowledgement floor.

```shell
nats consumer next ORDERS DISPATCH
```

--------------------------------

### NATS JetStream Message Acknowledgment Subject Format

Source: https://github.com/nats-io/nats.docs/blob/master/using-nats/jetstream/nats_api_reference.md

Explains the format of the reply subject for acknowledging messages in NATS JetStream.

```text
$JS.ACK.<stream>.<consumer>.<delivered count>.<stream sequence>.<consumer sequence>.<timestamp>.<pending messages>
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

### View NATS JetStream Stream Messages

Source: https://github.com/nats-io/nats.docs/blob/master/nats-concepts/jetstream/js_walkthrough.md

Allows viewing the messages stored within a specific JetStream stream. This is useful for inspecting the data that has been captured.

```shell
nats stream view my_stream
```

--------------------------------

### Get Specific Message from NATS JetStream Stream

Source: https://github.com/nats-io/nats.docs/blob/master/nats-concepts/jetstream/js_walkthrough.md

Retrieves a specific message from a JetStream stream, identified by its sequence number or other criteria. This command is used to fetch individual messages.

```shell
nats stream get my_stream
```

--------------------------------

### NATS JetStream Flow Control Message Reply

Source: https://github.com/nats-io/nats.docs/blob/master/using-nats/jetstream/nats_api_reference.md

Demonstrates how to respond to a flow control message in NATS JetStream. This is crucial to prevent consumer stalling.

```go
msg.Respond(nil)
```

--------------------------------

### Checking Stream Information for Deduplication

Source: https://github.com/nats-io/nats.docs/blob/master/using-nats/jetstream/model_deep_dive.md

After publishing messages with `Nats-Msg-Id`, you can inspect the stream to confirm that only the first occurrence of a message was stored.

```shell
nats stream info ORDERS
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

### Enable Domain-Aware Acknowledgements and Flow Control (v2)

Source: https://github.com/nats-io/nats.docs/blob/master/release_notes/whats_new_214.md

Configure JetStream to use the v2 format for acknowledgements and flow control subjects by setting the `js_ack_fc_v2` feature flag to `true`. This enables domain and account hashing for deconfliction. Not specifying the flag uses the server default (v1 for 2.14).

```hcl
feature_flags {
  js_ack_fc_v2: true
}
```

--------------------------------

### Report JetStream status and replication

Source: https://github.com/nats-io/nats.docs/blob/master/running-a-nats-service/configuration/leafnodes/jetstream_leafnodes.md

Retrieves a report of all streams and their replication status within a specific JetStream domain. This helps verify that mirrors and sources are active and synchronized.

```bash
nats --server nats://acc:acc@localhost:4222 stream report --js-domain hub
```

--------------------------------

### C# JetStream Publish and Fetch Consumer

Source: https://github.com/nats-io/nats.docs/blob/master/using-nats/developing-with-nats/js/consumers.md

Demonstrates creating a JetStream stream, publishing messages (including concurrently), and setting up a consumer with a maximum inflight message limit. It then fetches messages and acknowledges them.

```csharp
// dotnet add package NATS.Net
using NATS.Net;
using NATS.Client.JetStream;
using NATS.Client.JetStream.Models;

await using var client = new NatsClient();

INatsJSContext js = client.CreateJetStreamContext();

// Create a stream
var streamConfig = new StreamConfig(name: "FOO", subjects: ["foo"]);
await js.CreateStreamAsync(streamConfig);

// Publish a message
{
    PubAckResponse ack = await js.PublishAsync("foo", "Hello, JetStream!");
    ack.EnsureSuccess();
}

// Publish messages concurrently
List<NatsJSPublishConcurrentFuture> futures = new();
for (var i = 0; i < 500; i++)
{
    NatsJSPublishConcurrentFuture future
        = await js.PublishConcurrentAsync("foo", "Hello, JetStream 1!");
    futures.Add(future);
}

foreach (var future in futures)
{
    await using (future)
    {
        PubAckResponse ack = await future.GetResponseAsync();
        ack.EnsureSuccess();
    }
}


// Create a consumer with a maximum 128 inflight messages
INatsJSConsumer consumer = await js.CreateConsumerAsync("FOO", new ConsumerConfig(name: "foo")
{
    MaxWaiting = 128,
});

using CancellationTokenSource cts = new(TimeSpan.FromSeconds(10));

while (cts.IsCancellationRequested == false)
{
    var opts = new NatsJSFetchOpts { MaxMsgs = 10 };
    await foreach (NatsJSMsg<string> msg in consumer.FetchAsync<string>(opts, cancellationToken: cts.Token))
    {
        await msg.AckAsync(cancellationToken: cts.Token);
    }
}
```

--------------------------------

### Python JetStream Consumer Example

Source: https://github.com/nats-io/nats.docs/blob/master/using-nats/developing-with-nats/js/consumers.md

Demonstrates creating a JetStream context, adding a stream, publishing messages, creating pull and push consumers, and subscribing to a durable queue group. Includes message fetching, acknowledgment, and handling of load-balanced messages.

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

    # Create single push based subscriber that is durable across restarts.
    sub = await js.subscribe("foo", durable="myapp")
    msg = await sub.next_msg()
    await msg.ack()

    # Create deliver group that will be have load balanced messages.
    async def qsub_a(msg):
        print("QSUB A:", msg)
        await msg.ack()

    async def qsub_b(msg):
        print("QSUB B:", msg)
        await msg.ack()
    await js.subscribe("foo", "workers", cb=qsub_a)
    await js.subscribe("foo", "workers", cb=qsub_b)

    for i in range(0, 10):
        ack = await js.publish("foo", f"hello world: {i}".encode())
        print("\t", ack)

    await nc.close()

if __name__ == '__main__':
    asyncio.run(main())
```

--------------------------------

### Purge Stream Messages (JavaScript)

Source: https://github.com/nats-io/nats.docs/blob/master/using-nats/developing-with-nats/js/streams.md

Purges messages from a NATS JetStream stream based on various criteria like subject filters, sequence numbers, and keeping a certain number of messages.

```javascript
await jsm.streams.purge(stream, { filter: "a.b" });
await jsm.streams.purge(stream, { filter: "a.c", keep: 5 });
await jsm.streams.purge(stream, { seq: 100 });
await jsm.streams.purge(stream, { filter: "a.d", seq: 100 });
```

--------------------------------

### Add NATS JetStream Stream

Source: https://github.com/nats-io/nats.docs/blob/master/nats-concepts/jetstream/js_walkthrough.md

Creates a new JetStream stream named 'my_stream' that will capture messages published to the 'foo' subject. It uses default settings for storage, retention, and other stream configurations.

```shell
nats stream add my_stream
? Subjects foo
? Storage file
? Replication 1
? Retention Policy Limits
? Discard Policy Old
? Stream Messages Limit -1
? Per Subject Messages Limit -1
? Total Stream Size -1
? Message TTL -1
? Max Message Size -1
? Duplicate tracking time window 2m0s
? Allow message Roll-ups No
? Allow message deletion Yes
? Purge subjects or the entire stream Yes
Stream my_stream was created
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

### Get Consumer Info

Source: https://github.com/nats-io/nats.docs/blob/master/using-nats/jetstream/model_deep_dive.md

Retrieves detailed state information for a NATS JetStream consumer. This includes tracking message delivery and acknowledgements.

```shell
nats con info ORDERS DISPATCH
```

--------------------------------

### Benchmark JetStream Durable Consumer (Callback) Performance

Source: https://github.com/nats-io/nats.docs/blob/master/using-nats/nats-tools/nats_cli/natsbench.md

Measures the performance of consuming messages using a durable consumer with the `Consume()` callback function. This method is suitable for distributing message consumption across multiple clients with explicit acknowledgements.

```bash
nats bench js consume --clients 4 --no-progress
```

```text
13:46:04 Starting JetStream durable consumer (callback) benchmark [acks=explicit, batch=500, clients=4, consumer=nats-bench, double-acked=false, msg-size=128 B, msgs=100,000, purge=false, sleep=0s, stream=benchstream]
13:46:04 [1] Starting JetStream durable consumer (callback), expecting 25,000 messages
13:46:04 [2] Starting JetStream durable consumer (callback), expecting 25,000 messages
13:46:04 [3] Starting JetStream durable consumer (callback), expecting 25,000 messages
13:46:04 [4] Starting JetStream durable consumer (callback), expecting 25,000 messages

  [1] 73,230 msgs/sec ~ 8.9 MiB/sec ~ 13.66us (25,000 msgs)
  [2] 72,921 msgs/sec ~ 8.9 MiB/sec ~ 13.71us (25,000 msgs)
  [3] 72,696 msgs/sec ~ 8.9 MiB/sec ~ 13.76us (25,000 msgs)
  [4] 72,687 msgs/sec ~ 8.9 MiB/sec ~ 13.76us (25,000 msgs)

 NATS JetStream durable consumer (callback) aggregated stats: 290,438 msgs/sec ~ 36 MiB/sec
 message rates min 72,687 | avg 72,883 | max 73,230 | stddev 220 msgs
 avg latencies min 13.66us | avg 13.72us | max 13.76us | stddev 0.04us
```

--------------------------------

### C Example: NATS JetStream Consumer Logic

Source: https://github.com/nats-io/nats.docs/blob/master/using-nats/developing-with-nats/js/consumers.md

This C code demonstrates the core logic for a NATS JetStream consumer, handling stream creation, subscription, message fetching, acknowledgment, and cleanup. It supports both pull and push subscription models.

```c
            cfg.Storage = js_MemoryStorage;
            // Add the stream.
            s = js_AddStream(&si, js, &cfg, NULL, &jerr);
        }
        if (s == NATS_OK)
        {
            printf("Stream %s has %" PRIu64 " messages (%
```

```c
PRIu64 " bytes)\n",
                si->Config->Name, si->State.Msgs, si->State.Bytes);

            // Need to destroy the returned stream object.
            jsStreamInfo_Destroy(si);
        }
    }

    if (s == NATS_OK)
    {
        if (pull)
            s = js_PullSubscribe(&sub, js, subj, durable, &jsOpts, &so, &jerr);
        else if (async)
            s = js_Subscribe(&sub, js, subj, onMsg, NULL, &jsOpts, &so, &jerr);
        else
            s = js_SubscribeSync(&sub, js, subj, &jsOpts, &so, &jerr);
    }
    if (s == NATS_OK)
        s = natsSubscription_SetPendingLimits(sub, -1, -1);

    if (s == NATS_OK)
        s = natsStatistics_Create(&stats);

    if ((s == NATS_OK) && pull)
    {
        natsMsgList list;
        int         i;

        for (count = 0; (s == NATS_OK) && (count < total); )
        {
            s = natsSubscription_Fetch(&list, sub, 1024, 5000, &jerr);
            if (s != NATS_OK)
                break;

            if (start == 0)
                start = nats_Now();

            count += (int64_t) list.Count;
            for (i=0; (s == NATS_OK) && (i<list.Count); i++)
                s = natsMsg_Ack(list.Msgs[i], &jsOpts);

            natsMsgList_Destroy(&list);
        }
    }
    else if ((s == NATS_OK) && async)
    {
        while (s == NATS_OK)
        {
            if (count + dropped == total)
                break;

            nats_Sleep(1000);
        }
    }
    else if (s == NATS_OK)
    {
        for (count = 0; (s == NATS_OK) && (count < total); count++)
        {
            s = natsSubscription_NextMsg(&msg, sub, 5000);
            if (s != NATS_OK)
                break;

            if (start == 0)
                start = nats_Now();

            s = natsMsg_Ack(msg, &jsOpts);
            natsMsg_Destroy(msg);
        }
    }

    if (s == NATS_OK)
    {
        printStats(STATS_IN|STATS_COUNT, conn, sub, stats);
        printPerf("Received");
    }
    if (s == NATS_OK)
    {
        jsStreamInfo *si = NULL;

        // Let's report some stats after the run
        s = js_GetStreamInfo(&si, js, stream, NULL, &jerr);
        if (s == NATS_OK)
        {
            printf("\nStream %s has %" PRIu64 " messages (%
```

```c
PRIu64 " bytes)\n",
                si->Config->Name, si->State.Msgs, si->State.Bytes);

            jsStreamInfo_Destroy(si);
        }
        if (delStream)
        {
            printf("\nDeleting stream %s: ", stream);
            s = js_DeleteStream(js, stream, NULL, &jerr);
            if (s == NATS_OK)
                printf("OK!");
            printf("\n");
        }
    }
    else
    {
        printf("Error: %u - %s - jerr=%u\n", s, natsStatus_GetText(s), jerr);
        nats_PrintLastErrorStack(stderr);
    }

    // Destroy all our objects to avoid report of memory leak
    jsCtx_Destroy(js);
    natsStatistics_Destroy(stats);
    natsSubscription_Destroy(sub);
    natsConnection_Destroy(conn);
    natsOptions_Destroy(opts);

    // To silence reports of memory still in used with valgrind
    nats_Close();.

    return 0;
}
```

```

--------------------------------

### Copy and Edit NATS JetStream Streams

Source: https://github.com/nats-io/nats.docs/blob/master/running-a-nats-service/nats_admin/jetstream_admin/streams.md

Commands for duplicating an existing stream with modified configurations and updating the properties of an existing stream.

```shell
nats str cp ORDERS ARCHIVE --subjects "ORDERS_ARCHIVE.*" --max-age 2y
nats str info ORDERS -j | jq .config.subjects
nats str edit ORDERS --subjects "ORDERS.*"
nats str edit ORDERS --config orders.json
```