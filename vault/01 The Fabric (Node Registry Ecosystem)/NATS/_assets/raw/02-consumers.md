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

### C# JetStream Push Consumer Note

Source: https://github.com/nats-io/nats.docs/blob/master/using-nats/developing-with-nats/js/consumers.md

NATS .NET client abstracts push consumer mechanics. All consumers are treated generically, with message delivery handled by the library.

```csharp
// NATS .NET doesn't publicly support push consumers and treats all consumers
// as just consumers. The mechanics of the consuming messages are abstracted
// away from the applications and are handled by the library.
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

### Benchmark JetStream Durable Consumer (Fetch) Performance

Source: https://github.com/nats-io/nats.docs/blob/master/using-nats/nats-tools/nats_cli/natsbench.md

Measures the performance of retrieving messages in batches using the `Fetch()` function on a durable consumer. This benchmark tests consumption without explicit acknowledgements.

```bash
nats bench js fetch --acks none --clients 2
```

```text
14:09:10 Starting JetStream durable consumer (fetch) benchmark [acks=none, batch=500, clients=2, consumer=nats-bench, double-acked=false, msg-size=128 B, msgs=100,000, purge=false, sleep=0s, stream=benchstream]
14:09:10 [1] Starting JetStream durable consumer (fetch), expecting 50,000 messages
14:09:10 [2] Starting JetStream durable consumer (fetch), expecting 50,000 messages
Finished      0s [================================================================] 100%
Finished      0s [================================================================] 100%

  [1] 567,330 msgs/sec ~ 69 MiB/sec ~ 1.76us (50,000 msgs)
  [2] 567,067 msgs/sec ~ 69 MiB/sec ~ 1.76us (50,000 msgs)

 NATS JetStream durable consumer (fetch) aggregated stats: 1,128,932 msgs/sec ~ 138 MiB/sec
 message rates min 567,067 | avg 567,198 | max 567,330 | stddev 131 msgs
 avg latencies min 1.76us | avg 1.76us | max 1.76us | stddev 0.00us
```

--------------------------------

### Create Durable Consumer

Source: https://github.com/nats-io/nats.docs/blob/master/using-nats/jetstream/nats_api_reference.md

Creates a durable consumer for a given stream. Durable consumers retain their state across server restarts.

```APIDOC
## POST $JS.API.CONSUMER.DURABLE.CREATE.<stream>.<consumer>

### Description
Creates a durable consumer for a given stream.

### Method
POST

### Endpoint
`$JS.API.CONSUMER.DURABLE.CREATE.<stream>.<consumer>`

### Parameters
#### Path Parameters
- **stream** (string) - Required - The name of the stream.
- **consumer** (string) - Required - The name of the durable consumer.

#### Request Body
- **stream** (string) - Required - The name of the stream.
- **config** (object) - Required - The consumer configuration.
  - **name** (string) - Optional - The name of the consumer.
  - **durable** (string) - Optional - The name of the durable consumer (if applicable).
  - **filter_subject** (string) - Optional - A filter subject for the consumer.
  - **deliver_subject** (string) - Optional - The subject to deliver messages to.
  - **ack_policy** (string) - Optional - The acknowledgment policy (e.g., 'explicit', 'all').
  - **ack_wait** (integer) - Optional - The acknowledgment wait time in milliseconds.
  - **max_deliver** (integer) - Optional - The maximum number of deliveries for a message.
  - **backoff** (object) - Optional - The backoff configuration for retries.
    - **max_interval** (integer) - Optional - The maximum interval in milliseconds.
    - **max_elapsed** (integer) - Optional - The maximum elapsed time in milliseconds.
  - **rate_limit** (integer) - Optional - The rate limit for message delivery in bytes per second.
  - **sample_freq** (string) - Optional - The sample frequency for message delivery.
  - **headers_only** (boolean) - Optional - Whether to only send headers.
  - **replay_policy** (string) - Optional - The replay policy (e.g., 'instant', 'original').
  - **deliver_policy** (string) - Optional - The delivery policy (e.g., 'all', 'last', 'new', 'workq').
  - **deliver_group** (string) - Optional - The deliver group for the consumer.
  - **start_sequence** (integer) - Optional - The starting sequence number.
  - **start_time** (integer) - Optional - The starting time in milliseconds since epoch.
  - **start_delta** (string) - Optional - The starting time delta (e.g., '10s').

### Response
#### Success Response (200)
- **consumer** (object) - Information about the created consumer.
  - **stream** (string) - The name of the stream.
  - **name** (string) - The name of the consumer.
  - **durable** (boolean) - Whether the consumer is durable.
  - **api_version** (integer) - The API version.
  - **config** (object) - The consumer configuration.
  - **created** (string) - The creation timestamp.
  - **delivered** (object) - Information about delivered messages.
  - **ack_floor** (object) - Information about acknowledged messages.
  - **num_pending** (integer) - The number of pending messages.
  - **num_outstanding** (integer) - The number of outstanding messages.
  - **num_redelivered** (integer) - The number of redelivered messages.
  - **num_waiting** (integer) - The number of waiting messages.
  - **num_ack_pending** (integer) - The number of acknowledgment pending messages.
  - **num_messages_delivered** (integer) - The total number of messages delivered.
  - **num_messages_fetched** (integer) - The total number of messages fetched.
  - **num_messages_acked** (integer) - The total number of messages acknowledged.
  - **num_messages_redelivered** (integer) - The total number of messages redelivered.
  - **num_messages_pending** (integer) - The total number of messages pending.
  - **last_active** (string) - The last active timestamp.
  - **inactive_threshold** (integer) - The inactive threshold in seconds.
  - **memory_max_bytes** (integer) - The maximum memory usage in bytes.
  - **disk_max_bytes** (integer) - The maximum disk usage in bytes.
  - **pause_remaining** (integer) - The remaining pause time in milliseconds.

### Request Example
```json
{
  "stream": "my_stream",
  "config": {
    "name": "my_durable_consumer",
    "deliver_policy": "all"
  }
}
```

### Response Example
```json
{
  "consumer": {
    "stream": "my_stream",
    "name": "my_durable_consumer",
    "durable": true,
    "api_version": 1,
    "config": {
      "name": "my_durable_consumer",
      "durable": "my_durable_consumer",
      "deliver_policy": "all"
    },
    "created": "2023-10-27T10:00:00Z",
    "delivered": {
      "offset": 0,
      "msgs": 0,
      "bytes": 0
    },
    "ack_floor": {
      "offset": 0,
      "msgs": 0,
      "bytes": 0
    },
    "num_pending": 0,
    "num_outstanding": 0,
    "num_redelivered": 0,
    "num_waiting": 0,
    "num_ack_pending": 0,
    "num_messages_delivered": 0,
    "num_messages_fetched": 0,
    "num_messages_acked": 0,
    "num_messages_redelivered": 0,
    "num_messages_pending": 0,
    "last_active": "2023-10-27T10:00:00Z",
    "inactive_threshold": 300,
    "memory_max_bytes": 1048576,
    "disk_max_bytes": 1073741824,
    "pause_remaining": 0
  }
}
```
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

### NATS JetStream Consumer Administration Subjects

Source: https://github.com/nats-io/nats.docs/blob/master/using-nats/jetstream/nats_api_reference.md

Subjects for managing NATS JetStream consumers, including creation, durable consumer setup, deletion, resetting, and pausing.

```text
$JS.API.CONSUMER.CREATE.<stream>
$JS.API.CONSUMER.DURABLE.CREATE.<stream>.<consumer>
$JS.API.CONSUMER.DELETE.<stream>.<consumer>
$JS.API.CONSUMER.RESET.<stream>.<consumer>
$JS.API.CONSUMER.PAUSE.<stream>.<consumer>
$JS.API.CONSUMER.INFO.<stream>.<consumer>
$JS.API.CONSUMER.LIST.<stream>
$JS.API.CONSUMER.NAMES.<stream>
```