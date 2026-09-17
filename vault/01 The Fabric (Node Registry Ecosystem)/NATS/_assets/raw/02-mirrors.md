### Report JetStream status and replication

Source: https://github.com/nats-io/nats.docs/blob/master/running-a-nats-service/configuration/leafnodes/jetstream_leafnodes.md

Retrieves a report of all streams and their replication status within a specific JetStream domain. This helps verify that mirrors and sources are active and synchronized.

```bash
nats --server nats://acc:acc@localhost:4222 stream report --js-domain hub
```

--------------------------------

### Mirror a JetStream across domains

Source: https://github.com/nats-io/nats.docs/blob/master/running-a-nats-service/configuration/leafnodes/jetstream_leafnodes.md

Creates a mirror of a stream from a foreign domain into the local hub domain. This ensures data is available locally even if the connection to the leaf node is lost.

```bash
nats --server nats://acc:acc@localhost:4222 stream add --js-domain hub --mirror test
```

--------------------------------

### Create Mirrored NATS JetStream

Source: https://github.com/nats-io/nats.docs/blob/master/running-a-nats-service/nats_admin/jetstream_admin/replication.md

Creates a new stream named REPORT that mirrors the contents of the ARCHIVE stream.

```shell
nats s add REPORT --mirror ARCHIVE
```

--------------------------------

### Minimal Stream Configuration with Mirror

Source: https://github.com/nats-io/nats.docs/blob/master/nats-concepts/jetstream/source_and_mirror_example.md

A basic JSON configuration for a NATS JetStream stream that mirrors another stream. It specifies the target stream name and the name of the stream to be mirrored.

```json
{
  "name": "MIRROR_TARGET"
  "discard": "old",
  "mirror": {
    "name": "MIRROR_ORIGIN"
  },
  "deny_delete": false,
  "sealed": false,
  "max_msg_size": -1,
  "allow_rollup_hdrs": false,
  "max_bytes": -1,
  "storage": "file",
  "allow_direct": false,
  "max_age": 0,
  "max_consumers": -1,
  "max_msgs_per_subject": -1,
  "num_replicas": 1,
  "name": "MIRROR_TARGET",
  "deny_purge": false,
  "compression": "none",
  "max_msgs": -1,
  "retention": "limits",
  "mirror_direct": false
}
```

--------------------------------

### Advanced Stream Configuration with Mirror Options

Source: https://github.com/nats-io/nats.docs/blob/master/nats-concepts/jetstream/source_and_mirror_example.md

An advanced JSON configuration for a NATS JetStream stream mirroring another stream. This includes options for starting from a specific time, configuring external APIs, and other stream settings.

```json
{
  "name": "MIRROR_TARGET"
  "discard": "old",
  "mirror": {
    "opt_start_time": "2024-07-11T08:57:20.4441646Z",
    "external": {
      "deliver": "",
      "api": "$JS.domainB.API"
    },
    "name": "MIRROR_ORIGIN"
  },
  "consumer_limits": {
    
  },
  "deny_delete": false,
  "sealed": false,
  "max_msg_size": -1,
  "allow_rollup_hdrs": false,
  "max_bytes": -1,
  "storage": "file",
  "allow_direct": false,
  "max_age": 0,
  "max_consumers": -1,
  "max_msgs_per_subject": -1,
  "num_replicas": 1,
  "name": "MIRROR_TARGET",
  "deny_purge": false,
  "compression": "none",
  "max_msgs": -1,
  "retention": "limits",
  "mirror_direct": false
}
```

--------------------------------

### Add Replicated Stream with Placement Tags

Source: https://github.com/nats-io/nats.docs/blob/master/nats-concepts/jetstream/README.md

Creates a new JetStream with 3 replicas, specifically targeting servers that have the 'sync:always' tag.

```bash
nats stream add --replicas 3 --tag sync:always
```

--------------------------------

### Aggregate JetStream sources across domains

Source: https://github.com/nats-io/nats.docs/blob/master/running-a-nats-service/configuration/leafnodes/jetstream_leafnodes.md

Aggregates streams from multiple leaf nodes into a single stream in the hub domain using the source feature. This is ideal for centralized processing of data originating from various locations.

```bash
nats --server nats://acc:acc@localhost:4222 stream add --js-domain hub --source test
```

--------------------------------

### Advanced Stream Configuration with Multiple Sources

Source: https://github.com/nats-io/nats.docs/blob/master/nats-concepts/jetstream/source_and_mirror_example.md

An enhanced JSON configuration for a NATS JetStream stream with multiple sources. This example includes filtering subjects, starting from specific message sequence numbers, and external API configurations for sources.

```json
{
  "name": "SOURCE_TARGET",
  "subjects": [
    "foo1.ext.*",
    "foo2.ext.*"
  ],
  "discard": "old",
  "duplicate_window": 120000000000,
  "sources": [
    {
      "name": "SOURCE1_ORIGIN",
      "filter_subject": "foo1.bar",
      "opt_start_seq": 42,
      "external": {
        "deliver": "",
        "api": "$JS.domainA.API"
      }
    },
    {
      "name": "SOURCE2_ORIGIN",
      "filter_subject": "foo2.bar"
    }
  ],
  "consumer_limits": {
    
  },
  "deny_delete": false,
  "sealed": false,
  "max_msg_size": -1,
  "allow_rollup_hdrs": false,
  "max_bytes": -1,
  "storage": "file",
  "allow_direct": false,
  "max_age": 0,
  "max_consumers": -1,
  "max_msgs_per_subject": -1,
  "num_replicas": 1,
  "name": "SOURCE_TARGET",
  "deny_purge": false,
  "compression": "none",
  "max_msgs": -1,
  "retention": "limits",
  "mirror_direct": false
}
```

--------------------------------

### Minimal Stream Configuration with Sources

Source: https://github.com/nats-io/nats.docs/blob/master/nats-concepts/jetstream/source_and_mirror_example.md

A basic JSON configuration for a NATS JetStream stream that ingests messages from other streams. It defines the stream name, subjects to consume, and a single source stream.

```json
{
  "name": "SOURCE_TARGET",
  "subjects": [
    "foo1.ext.*",
    "foo2.ext.*"
  ],
  "discard": "old",
  "duplicate_window": 120000000000,
  "sources": [
    {
      "name": "SOURCE1_ORIGIN"
    }
  ],
  "deny_delete": false,
  "sealed": false,
  "max_msg_size": -1,
  "allow_rollup_hdrs": false,
  "max_bytes": -1,
  "storage": "file",
  "allow_direct": false,
  "max_age": 0,
  "max_consumers": -1,
  "max_msgs_per_subject": -1,
  "num_replicas": 1,
  "name": "SOURCE_TARGET",
  "deny_purge": false,
  "compression": "none",
  "max_msgs": -1,
  "retention": "limits",
  "mirror_direct": false
}
```

--------------------------------

### Create Sourced NATS JetStream

Source: https://github.com/nats-io/nats.docs/blob/master/running-a-nats-service/nats_admin/jetstream_admin/replication.md

Creates a new stream named ARCHIVE that aggregates data from the ORDERS and RETURNS streams.

```shell
nats s add ARCHIVE --source ORDERS --source RETURNS
```

--------------------------------

### Inspect NATS JetStream Information

Source: https://github.com/nats-io/nats.docs/blob/master/running-a-nats-service/nats_admin/jetstream_admin/replication.md

Retrieves detailed configuration and state information for a specific stream, including source or mirror lag details.

```shell
nats stream info ARCHIVE
nats stream info REPORT
```

--------------------------------

### Create Clustered NATS Stream with Replicas

Source: https://github.com/nats-io/nats.docs/blob/master/running-a-nats-service/configuration/clustering/jetstream_clustering/administration.md

This command creates a new stream in NATS JetStream with a specified number of replicas. Choosing more than one replica ensures data redundancy across multiple nodes in the cluster using the RAFT protocol. The output shows stream information and cluster details, including leader and replica status.

```shell
nats stream add ORDERS --replicas 3
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

### Benchmark JetStream Ordered Ephemeral Consumer Performance

Source: https://github.com/nats-io/nats.docs/blob/master/using-nats/nats-tools/nats_cli/natsbench.md

Measures the replay performance of messages using an ordered ephemeral consumer. This method ensures each client receives its own copy of the messages.

```bash
nats bench js ordered
```

```text
13:33:48 Starting JetStream ordered ephemeral consumer benchmark [clients=1, msg-size=128 B, msgs=100,000, purge=false, sleep=0s, stream=benchstream]
13:33:48 [1] Starting JetStream ordered ephemeral consumer, expecting 100,000 messages
Finished      0s [================================================================] 100%

NATS JetStream ordered ephemeral consumer stats: 1,201,540 msgs/sec ~ 147 MiB/sec ~ 0.83us
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

### Configure NATS Cross-Account JetStream Exports and Imports

Source: https://github.com/nats-io/nats.docs/blob/master/running-a-nats-service/configuration/leafnodes/jetstream_leafnodes.md

This configuration demonstrates how to define exports in a source account to allow other accounts to access JetStream consumers and streams. It also shows how to define imports in target accounts, including subject remapping for API calls.

```text
accounts {
    SYS: {
        users: [{user: admin, password: admin}]
    },
    ACC: {
        users: [{user: acc, password: acc}],
        jetstream: enabled
        exports: [
            {service: "$JS.hub.API.CONSUMER.CREATE.*", response_type: "stream"}
            {stream: "deliver.acc.hub.>"}
            {service: "$JS.hub.API.CONSUMER.MSG.NEXT.aggregate-test-leaf.dur", response_type: "stream"}
            {service: "$JS.ACK.aggregate-test-leaf.dur.>"}
            {service: "$JS.FC.aggregate-test-leaf.dur.>"}
        ]
    }
    IMPORT_MIRROR: {
        users: [{user: import_mirror, password: import_mirror}],
        jetstream: enabled
        imports: [
            {service: {account: ACC, subject: "$JS.hub.API.CONSUMER.CREATE.*"}, to: "JS.acc@hub.API.CONSUMER.CREATE.*" }
            {service: {account: ACC, subject: "$JS.FC.aggregate-test-leaf.dur.>"}}
            {stream: {account: ACC, subject: deliver.acc.hub.import_mirror.>}}
        ]
    }
    IMPORT_CLIENT: {
        users: [{user: import_client, password: import_client}],
        jetstream: enabled
        imports: [
            {service: {account: ACC, subject: "$JS.hub.API.CONSUMER.MSG.NEXT.aggregate-test-leaf.dur"}, to: "JS.acc@hub.API.CONSUMER.MSG.NEXT.aggregate-test-leaf.dur" }
            {service: {account: ACC, subject: "$JS.ACK.aggregate-test-leaf.dur.>"}}
        ]
    }
}
system_account: SYS
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

### Enable JetStream in Account Configuration

Source: https://github.com/nats-io/nats.docs/blob/master/running-a-nats-service/configuration/README.md

Minimal configuration to enable JetStream for a specific account. JetStream is disabled by default.

```hcl
accounts {
  A {
      jetstream {
      }
    }
}

```

--------------------------------

### Obtain JetStream Context (JavaScript)

Source: https://github.com/nats-io/nats.docs/blob/master/using-nats/developing-with-nats/js/context.md

Illustrates how to retrieve the JetStream context and JetStream Management context using the NATS client library in JavaScript. The JetStream context is for message operations, while the management context is for stream and consumer administration.

```JavaScript
const nc = await connect();
// Getting the JetStream context
const js = nc.jetstream();
// Getting the JetStream management context
const jsm = await nc.jetstreamManager();
```

--------------------------------

### Basic JetStream Connection and Stream Creation (Python)

Source: https://github.com/nats-io/nats.docs/blob/master/using-nats/developing-with-nats/js/streams.md

Establishes a connection to NATS and creates a JetStream context, then adds a stream for message persistence.

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

    await nc.close()

if __name__ == '__main__':
    asyncio.run(main())
```

--------------------------------

### NATS JetStream Request Command Help

Source: https://github.com/nats-io/nats.docs/blob/master/running-a-nats-service/configuration/clustering/jetstream_clustering/administration.md

This command displays the help information for the `nats server request jetstream` command, outlining available flags and arguments for querying JetStream details. It shows options for server connection, authentication, and filtering the requested JetStream information.

```shell
nats server req jetstream --user admin --password s3cr3t! --help
```