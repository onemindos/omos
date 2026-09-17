### Configure JetStream via Configuration File

Source: https://github.com/nats-io/nats.docs/blob/master/using-nats/jetstream/getting_started/README.md

Starts the NATS server using a custom configuration file and defines the storage directory for JetStream data.

```shell
nats-server -c js.conf
```

```text
# js.conf
jetstream {
   store_dir=nats
}
```

--------------------------------

### Obtain JetStream Context (Go)

Source: https://github.com/nats-io/nats.docs/blob/master/using-nats/developing-with-nats/js/context.md

Demonstrates how to obtain a JetStream context in Go. It shows how to get a context for publishing/subscribing only, for managing streams only, or a combined context for all operations. Requires the 'nats' library.

```Go
func ExampleJetStreamContext() {
nc, _ := nats.Connect("localhost")

var js nats.JetStream
var jsm nats.JetStreamManager
var jsctx nats.JetStreamContext

// JetStream that can publish/subscribe but cannot manage streams.
js, _ = nc.JetStream()
js.Publish("foo", []byte("hello"))

// JetStream context that can manage streams/consumers but cannot produce messages.
jsm, _ = nc.JetStream()
jsm.AddStream(&nats.StreamConfig{Name: "FOO"})

// JetStream context that can both manage streams/consumers
// as well as publish/subscribe.
jsctx, _ = nc.JetStream()
jsctx.AddStream(&nats.StreamConfig{Name: "BAR"})
jsctx.Publish("bar", []byte("hello world"))
}
```

--------------------------------

### Enable JetStream with Custom Storage and Limits

Source: https://github.com/nats-io/nats.docs/blob/master/running-a-nats-service/configuration/README.md

Configure JetStream to use a local directory for storage and set memory and file storage limits. This example demonstrates a minimal configuration for JetStream.

```text
jetstream {
  store_dir: nats

  # 1GB
  max_memory_store: 1073741824

  # 10GB
  max_file_store: 10737418240
}
```

--------------------------------

### Obtain JetStream Context (Python)

Source: https://github.com/nats-io/nats.docs/blob/master/using-nats/developing-with-nats/js/context.md

Provides an example of how to get the JetStream context in Python using the NATS client library. This context is essential for performing JetStream operations like publishing and subscribing.

```Python
async def main():
    nc = await nats.connect("localhost")

    # Create JetStream context.
    js = nc.jetstream()
    
if __name__ == '__main__':
asyncio.run(main())
```

--------------------------------

### Start NATS Server with JetStream

Source: https://github.com/nats-io/nats.docs/blob/master/using-nats/jetstream/getting_started/using_source.md

Commands to launch the NATS server with JetStream enabled, including options for default and custom storage directories.

```shell
nats-server -js
nats-server -js -sd /tmp/test
```

--------------------------------

### Go JetStream Consumer Examples

Source: https://github.com/nats-io/nats.docs/blob/master/using-nats/developing-with-nats/js/consumers.md

Demonstrates various ways to create and use JetStream consumers in Go, including asynchronous and synchronous subscriptions, manual acknowledgments, and queue groups.

```go
func ExampleJetStream() {
	nc, err := nats.Connect("localhost")
	if err != nil {
		log.Fatal(err)
	}

	// Use the JetStream context to produce and consume messages
	// that have been persisted.
	js, err := nc.JetStream(nats.PublishAsyncMaxPending(256))
	if err != nil {
		log.Fatal(err)
	}

	js.AddStream(&nats.StreamConfig{
		Name:     "FOO",
		Subjects: []string{"foo"},
	})

	js.Publish("foo", []byte("Hello JS!"))

	// Publish messages asynchronously.
	for i := 0; i < 500; i++ {
		js.PublishAsync("foo", []byte("Hello JS Async!"))
	}
	select {
	case <-js.PublishAsyncComplete():
	case <-time.After(5 * time.Second):
		fmt.Println("Did not resolve in time")
	}

	// Create async consumer on subject 'foo'. Async subscribers
	// ack a message once exiting the callback.
	js.Subscribe("foo", func(msg *nats.Msg) {
		meta, _ := msg.Metadata()
		fmt.Printf("Stream Sequence  : %v\n", meta.Sequence.Stream)
		fmt.Printf("Consumer Sequence: %v\n", meta.Sequence.Consumer)
	})

	// Async subscriber with manual acks.
	js.Subscribe("foo", func(msg *nats.Msg) {
		msg.Ack()
	}, nats.ManualAck())

	// Async queue subscription where members load balance the
	// received messages together.
	// If no consumer name is specified, either with nats.Bind() 
	// or nats.Durable() options, the queue name is used as the
	// durable name (that is, as if you were passing the
	// nats.Durable(<queue group name>) option.
	// It is recommended to use nats.Bind() or nats.Durable() 
	// and preferably create the JetStream consumer beforehand
	// (using js.AddConsumer) so that the JS consumer is not
	// deleted on an Unsubscribe() or Drain() when the member
	// that created the consumer goes away first.
	// Check Godoc for the QueueSubscribe() API for more details.
	js.QueueSubscribe("foo", "group", func(msg *nats.Msg) {
		msg.Ack()
	}, nats.ManualAck())

	// Subscribe to consume messages synchronously.
	sub, _ := js.SubscribeSync("foo")
	msg, _ := sub.NextMsg(2 * time.Second)
	msg.Ack()

	// We can add a member to the group, with this member using
	// the synchronous version of the QueueSubscribe.
	sub, _ = js.QueueSubscribeSync("foo", "group")
	msg, _ = sub.NextMsg(2 * time.Second)
	msg.Ack()

	// ChanSubscribe
	msgCh := make(chan *nats.Msg, 8192)
	sub, _ = js.ChanSubscribe("foo", msgCh)

	select {
	case msg := <-msgCh:
		fmt.Println("[Received]", msg)
	case <-time.After(1 * time.Second):
	}
}

```

--------------------------------

### Create a NATS JetStream Stream

Source: https://github.com/nats-io/nats.docs/blob/master/running-a-nats-service/nats_admin/jetstream_admin/streams.md

Demonstrates how to create a new stream using interactive prompts, command-line flags, or a JSON configuration file.

```shell
nats str add ORDERS
nats str add ORDERS --subjects "ORDERS.*" --ack --max-msgs=-1 --max-bytes=-1 --max-age=1y --storage file --retention limits --max-msg-size=-1 --discard old --dupe-window="0s" --replicas 1
nats str add ORDERS --config orders.json
```

--------------------------------

### Get NATS JetStream Stream Info

Source: https://github.com/nats-io/nats.docs/blob/master/nats-concepts/jetstream/js_walkthrough.md

Retrieves and displays detailed information about the specified JetStream stream, including its subjects, configuration, limits, and current state.

```shell
nats stream info my_stream
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

### Start NATS Server with Monitoring

Source: https://github.com/nats-io/nats.docs/blob/master/using-nats/nats-tools/nats_cli/natsbench.md

Starts the NATS server instance with the HTTP monitoring port enabled and JetStream support active. This is a prerequisite for running performance benchmarks.

```bash
nats-server -m 8222 -js
```

--------------------------------

### Manage JetStream Streams and Consumers with GitHub Actions

Source: https://github.com/nats-io/nats.docs/blob/master/running-a-nats-service/configuration/jetstream-config/configuration_mgmt/github_actions.md

This snippet demonstrates how to use GitHub Actions to manage JetStream streams and consumers. It includes steps to delete existing streams and consumers if they exist, create new streams from configuration files, and create new consumers associated with those streams. The actions require the `nats-io/jetstream-gh-action` and depend on `actions/checkout` for accessing configuration files.

```yaml
on: push
name: orders
jobs:

  # First we delete the ORDERS stream and consumer if they already exist
  clean_orders:
    runs-on: ubuntu-latest
    steps:
      - name: orders_stream
        uses: nats-io/jetstream-gh-action/delete/stream@master
        with:
          missing_ok: 1
          stream: ORDERS
          server: js.example.net

  # Now we create the Stream and Consumers using the same configuration files the 
  # nats CLI utility would use as shown above
  create_orders:
    runs-on: ubuntu-latest
    needs: clean_orders
    steps:
      - uses: actions/checkout@master
      - name: orders_stream
        uses: nats-io/jetstream-gh-action/create/stream@master
        with:
          config: ORDERS.json
          server: js.example.net
      - name: orders_new_consumer
        uses: nats-io/jetstream-gh-action/create/consumer@master
        with:
          config: ORDERS_NEW.json
          stream: ORDERS
          server: js.example.net
```

--------------------------------

### Configure and Publish to NATS JetStream

Source: https://github.com/nats-io/nats.docs/blob/master/using-nats/developing-with-nats/js/publish.md

Demonstrates initializing a JetStream context, creating a memory-backed stream if it does not exist, and publishing messages either synchronously or asynchronously.

```c
s = natsConnection_JetStream(&js, conn, &jsOpts);
if (s == NATS_OK) {
    jsStreamConfig cfg;
    jsStreamConfig_Init(&cfg);
    cfg.Name = stream;
    cfg.Subjects = (const char*[1]){subj};
    cfg.SubjectsLen = 1;
    cfg.Storage = js_MemoryStorage;
    s = js_AddStream(&si, js, &cfg, NULL, &jerr);
}
// Publishing logic
if (async) {
    s = js_PublishAsync(js, subj, (const void*) payload, dataLen, NULL);
} else {
    jsPubAck *pa = NULL;
    s = js_Publish(&pa, js, subj, (const void*) payload, dataLen, NULL, &jerr);
    jsPubAck_Destroy(pa);
}
```

--------------------------------

### Get Specific Message from NATS JetStream Stream

Source: https://github.com/nats-io/nats.docs/blob/master/nats-concepts/jetstream/js_walkthrough.md

Retrieves a specific message from a JetStream stream, identified by its sequence number or other criteria. This command is used to fetch individual messages.

```shell
nats stream get my_stream
```

--------------------------------

### Configure and Add a JetStream Memory Stream

Source: https://github.com/nats-io/nats.docs/blob/master/using-nats/developing-with-nats/js/streams.md

Initializes a stream configuration, sets the storage type to memory, and adds the stream to the NATS JetStream context.

```c
jsStreamConfig_Init(&cfg);
cfg.Name = stream;
cfg.Subjects = (const char*[1]){subj};
cfg.SubjectsLen = 1;
cfg.Storage = js_MemoryStorage;
s = js_AddStream(&si, js, &cfg, NULL, &jerr);
```

--------------------------------

### Get Stream Info

Source: https://github.com/nats-io/nats.docs/blob/master/using-nats/jetstream/model_deep_dive.md

Retrieves information about a NATS JetStream stream. Useful for checking message counts and active consumers.

```shell
nats str info ORDERS
```

--------------------------------

### Get Consumer Info

Source: https://github.com/nats-io/nats.docs/blob/master/using-nats/jetstream/model_deep_dive.md

Retrieves detailed state information for a NATS JetStream consumer. This includes tracking message delivery and acknowledgements.

```shell
nats con info ORDERS DISPATCH
```

--------------------------------

### Create JetStream Stream

Source: https://github.com/nats-io/nats.docs/blob/master/running-a-nats-service/configuration/leafnodes/jetstream_leafnodes.md

Adds a new stream to the JetStream domain associated with the current connection. The command prompts for configuration details like storage backend and retention policies.

```bash
nats --server nats://acc:acc@localhost:4222 stream add
```

--------------------------------

### Manage JetStream Consumers (JavaScript)

Source: https://github.com/nats-io/nats.docs/blob/master/using-nats/developing-with-nats/js/streams.md

Demonstrates how to list, add, retrieve information about, and delete consumers for a NATS JetStream stream.

```javascript
// list all consumers for a stream:
const consumers = await jsm.consumers.list(stream).next();
consumers.forEach((ci) => {
    console.log(ci);
});

// add a new durable pull consumer
await jsm.consumers.add(stream, {
    durable_name: "me",
    ack_policy: AckPolicy.Explicit,
});

// retrieve a consumer's configuration
const ci = await jsm.consumers.info(stream, "me");
console.log(ci);

// delete a particular consumer
await jsm.consumers.delete(stream, "me");
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

### Manage NATS JetStream Resources with kubectl

Source: https://github.com/nats-io/nats.docs/blob/master/running-a-nats-service/configuration/jetstream-config/configuration_mgmt/kubernetes_controller.md

These bash commands demonstrate how to interact with NATS JetStream Streams and Consumers managed by Kubernetes CRDs. You can list existing streams and consumers, and use 'kubectl describe' to get detailed information about specific resources, which is helpful for troubleshooting.

```bash
$ kubectl get streams
NAME       STATE     STREAM NAME   SUBJECTS
mystream   Created   mystream      [orders.*]

$ kubectl get consumers
NAME               STATE     STREAM     CONSUMER           ACK POLICY
my-pull-consumer   Created   mystream   my-pull-consumer   explicit
my-push-consumer   Created   mystream   my-push-consumer   none

# If you end up in an Errored state, run kubectl describe for more info.
#     kubectl describe streams mystream
#     kubectl describe consumers my-pull-consumer

```

--------------------------------

### Benchmark JetStream Synchronous Getter Performance

Source: https://github.com/nats-io/nats.docs/blob/master/using-nats/nats-tools/nats_cli/natsbench.md

Measures the latency of directly retrieving messages by sequence number using synchronous gets, without involving a consumer.

```bash
nats bench js get sync
```

```text
14:13:30 Starting JetStream synchronous getter benchmark [clients=1, msg-size=128 B, msgs=100,000, sleep=0s, stream=benchstream]
14:13:30 [1] Starting JetStream synchronous getter, expecting 100,000 messages
Finished      3s [================================================================] 100%

NATS JetStream synchronous getter stats: 33,244 msgs/sec ~ 4.1 MiB/sec ~ 30.08us
```

--------------------------------

### Configure JetStream in NATS Configuration File

Source: https://github.com/nats-io/nats.docs/blob/master/using-nats/jetstream/getting_started/using_source.md

Example configuration block to enable JetStream and define storage directory and resource limits within the NATS server configuration file.

```text
jetstream {
    store_dir: "/data/nats-server"
    max_memory_store: 1073741824
    max_file_store: 10737418240
}
```