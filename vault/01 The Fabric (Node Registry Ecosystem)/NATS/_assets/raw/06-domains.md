### Mirror a JetStream across domains

Source: https://github.com/nats-io/nats.docs/blob/master/running-a-nats-service/configuration/leafnodes/jetstream_leafnodes.md

Creates a mirror of a stream from a foreign domain into the local hub domain. This ensures data is available locally even if the connection to the leaf node is lost.

```bash
nats --server nats://acc:acc@localhost:4222 stream add --js-domain hub --mirror test
```

--------------------------------

### Create JetStream Stream in Specific Domain

Source: https://github.com/nats-io/nats.docs/blob/master/running-a-nats-service/configuration/leafnodes/jetstream_leafnodes.md

Creates a stream within a target domain by providing the --js-domain flag. This allows for cross-domain management while maintaining a single connection point.

```bash
nats --server nats://acc:acc@localhost:4222 stream add --js-domain leaf
```

--------------------------------

### Create JetStream Stream

Source: https://github.com/nats-io/nats.docs/blob/master/running-a-nats-service/configuration/leafnodes/jetstream_leafnodes.md

Adds a new stream to the JetStream domain associated with the current connection. The command prompts for configuration details like storage backend and retention policies.

```bash
nats --server nats://acc:acc@localhost:4222 stream add
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

### Aggregate JetStream sources across domains

Source: https://github.com/nats-io/nats.docs/blob/master/running-a-nats-service/configuration/leafnodes/jetstream_leafnodes.md

Aggregates streams from multiple leaf nodes into a single stream in the hub domain using the source feature. This is ideal for centralized processing of data originating from various locations.

```bash
nats --server nats://acc:acc@localhost:4222 stream add --js-domain hub --source test
```

--------------------------------

### Configure NATS JetStream Leaf Node Connections

Source: https://github.com/nats-io/nats.docs/blob/master/running-a-nats-service/configuration/leafnodes/jetstream_leafnodes.md

Configuration files for setting up a hub and leaf node architecture with JetStream enabled. These examples demonstrate how to define domains to isolate JetStream namespaces and configure leaf node remotes for account synchronization.

```text
accounts {
    SYS: {
        users: [{user: admin, password: admin}]
    },
    ACC: {
        users: [{user: acc, password: acc}],
        jetstream: enabled
    }
}
system_account: SYS
```

```text
port: 4222
server_name: hub-server
jetstream {
    store_dir="./store_server"
    domain=hub
}
leafnodes {
    port: 7422
}
include ./accounts.conf
```

```text
port: 4111
server_name: leaf-server
jetstream {
    store_dir="./store_leaf"
    domain=leaf
}
leafnodes {
    remotes = [
        {
            urls: ["nats://admin:admin@0.0.0.0:7422"]
            account: "SYS"
        },
        {
            urls: ["nats://acc:acc@0.0.0.0:7422"]
            account: "ACC"
        }
    ]
}
include ./accounts.conf
```

--------------------------------

### JetStream Stream Management and Publishing

Source: https://github.com/nats-io/nats.docs/blob/master/using-nats/developing-with-nats/js/streams.md

Endpoints and functions for managing JetStream streams and publishing messages to subjects.

```APIDOC
## POST /js/streams

### Description
Creates a new JetStream stream with a specific configuration, including subjects and storage type.

### Method
POST

### Endpoint
/js/streams

### Parameters
#### Request Body
- **Name** (string) - Required - The name of the stream.
- **Subjects** (array) - Required - List of subjects to bind to the stream.
- **Storage** (enum) - Required - Storage type (e.g., js_MemoryStorage).

### Request Example
{
  "Name": "mystream",
  "Subjects": ["events.*"],
  "Storage": "memory"
}

### Response
#### Success Response (200)
- **si** (object) - Stream information object containing state and configuration.

---

## POST /js/publish

### Description
Publishes a message to a subject within a JetStream context. Supports both synchronous and asynchronous modes.

### Method
POST

### Endpoint
/js/publish

### Parameters
#### Query Parameters
- **async** (boolean) - Optional - If true, performs asynchronous publishing.

#### Request Body
- **subject** (string) - Required - The target subject.
- **payload** (bytes) - Required - The message content.

### Request Example
{
  "subject": "events.test",
  "payload": "hello world"
}

### Response
#### Success Response (200)
- **pa** (object) - Publication acknowledgment containing sequence and duplicate status.

---

## DELETE /js/streams/{stream}

### Description
Deletes an existing JetStream stream by name.

### Method
DELETE

### Endpoint
/js/streams/{stream}

### Parameters
#### Path Parameters
- **stream** (string) - Required - The name of the stream to delete.

### Response
#### Success Response (200)
- **status** (string) - Confirmation of successful deletion.
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

### NATS JetStream Stream Administration Subjects

Source: https://github.com/nats-io/nats.docs/blob/master/using-nats/jetstream/nats_api_reference.md

Subjects for managing NATS JetStream streams, including creation, update, deletion, and information retrieval.

```text
$JS.API.STREAM.CREATE.<stream>
$JS.API.STREAM.UPDATE.<stream>
$JS.API.STREAM.DELETE.<stream>
$JS.API.STREAM.INFO.<stream>
$JS.API.STREAM.PURGE.<stream>
$JS.API.STREAM.LIST
$JS.API.STREAM.NAMES
$JS.API.STREAM.MSG.DELETE.<stream>
$JS.API.STREAM.MSG.GET.<stream>
$JS.API.STREAM.SNAPSHOT.<stream>
$JS.API.STREAM.RESTORE.<stream>
```

--------------------------------

### JetStream Information (/jsz)

Source: https://github.com/nats-io/nats.docs/blob/master/running-a-nats-service/nats_admin/monitoring/README.md

The /jsz endpoint provides detailed information about JetStream, including account-specific metrics, streams, and consumers. It supports pagination and filtering.

```APIDOC
## GET /jsz

### Description
Retrieves detailed information about JetStream, including account statistics, streams, and consumers. This endpoint supports pagination and can be filtered by various parameters. In a clustered environment, it's recommended to query the leader for the most accurate data.

### Method
GET

### Endpoint
/jsz

### Query Parameters
- **acc** (string) - Optional - Include metrics for the specified account. Default is unset.
- **accounts** (boolean) - Optional - Include account-specific JetStream information. Default is false.
- **streams** (boolean) - Optional - Include streams. When set, implies `accounts=true`. Default is false.
- **consumers** (boolean) - Optional - Include consumer information. When set, implies `streams=true`. Default is false.
- **config** (boolean) - Optional - When stream or consumer are requested, include their respective configuration. Default is false.
- **leader-only** (boolean) - Optional - Only the leader responds. Default is false.
- **offset** (number) - Optional - Pagination offset. Must be greater than 0. Default is 0.
- **limit** (number) - Optional - Number of results to return. Must be greater than 0. Default is 1024.
- **raft** (boolean) - Optional - Include information details about the Raft group. Default is false.

### Response
#### Success Response (200)
- The response structure varies based on the query parameters used. It can include account information, stream details, consumer metrics, and Raft group information.

#### Error Response (400)
- Bad Request

### Request Example
- Basic JetStream information:
  `https://demo.nats.io:8222/jsz`
- Requesting accounts with custom limit and offset:
  `https://demo.nats.io:8222/jsz?accounts=true&limit=16&offset=128`
- Requesting detailed consumer information:
  `https://demo.nats.io:8222/jsz?consumers=true`
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

### NATS JetStream Request Command Help

Source: https://github.com/nats-io/nats.docs/blob/master/running-a-nats-service/configuration/clustering/jetstream_clustering/administration.md

This command displays the help information for the `nats server request jetstream` command, outlining available flags and arguments for querying JetStream details. It shows options for server connection, authentication, and filtering the requested JetStream information.

```shell
nats server req jetstream --user admin --password s3cr3t! --help
```

--------------------------------

### JetStream Information Endpoint

Source: https://github.com/nats-io/nats.docs/blob/master/running-a-nats-service/nats_admin/monitoring/README.md

Retrieve detailed JetStream information. Supports filtering by account, streams, consumers, and configuration. Pagination is available via `limit` and `offset`.

```http
GET /jsz HTTP/1.1
Host: demo.nats.io:8222


```

```http
GET /jsz?accounts=true&limit=16&offset=128 HTTP/1.1
Host: demo.nats.io:8222


```

```http
GET /jsz?consumers=true HTTP/1.1
Host: demo.nats.io:8222


```

--------------------------------

### View JetStream Advisories with nats CLI

Source: https://github.com/nats-io/nats.docs/blob/master/running-a-nats-service/nats_admin/monitoring/monitoring_jetstream.md

Use the `nats event --js-advisory` command to view all JetStream advisories on your console. This command helps in real-time monitoring of JetStream events.

```bash
nats event --js-advisory
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

### Manage JetStream Streams and Consumers

Source: https://github.com/nats-io/nats.docs/blob/master/using-nats/developing-with-nats/js/streams.md

Demonstrates common JetStream management operations including adding, updating, and deleting streams and consumers, as well as retrieving metadata.

```go
func ExampleJetStreamManager() {
	nc, _ := nats.Connect("localhost")

	js, _ := nc.JetStream()

	// Create a stream
	js.AddStream(&nats.StreamConfig{
		Name:     "example-stream",
		Subjects: []string{"example-subject"},
		MaxBytes: 1024,
	})

	// Update a stream
	js.UpdateStream(&nats.StreamConfig{
		Name:     "example-stream",
		MaxBytes: 2048,
	})

	// Create a durable consumer
	js.AddConsumer("example-stream", &nats.ConsumerConfig{
		Durable: "example-consumer-name",
	})

	// Get information about all streams (with Context JSOpt)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	for info := range js.StreamsInfo(nats.Context(ctx)) {
		fmt.Println("stream name: ", info.Config.Name)
	}

	// Get information about all consumers (with MaxWait JSOpt)
	for info := range js.ConsumersInfo("example-stream", nats.MaxWait(10*time.Second)) {
		fmt.Println("consumer name: ", info.Name)
	}

	// Delete a consumer
	js.DeleteConsumer("example-stream", "example-consumer-name")

	// Delete a stream
	js.DeleteStream("example-stream")
}
```

```java
try (Connection nc = Nats.connect("localhost")) {
    JetStreamManagement jsm = nc.jetStreamManagement();

    // Create a stream
    StreamInfo si = jsm.addStream(StreamConfiguration.builder()
        .name("example-stream")
        .subjects("example-subject")
        .maxBytes(1024)
        .build());
    StreamConfiguration config = si.getConfiguration();
    System.out.println("stream name: " + config.getName() + ", max_bytes: " + config.getMaxBytes());

    // Update a stream
    si = jsm.updateStream(StreamConfiguration.builder()
        .name("example-stream")
        .maxBytes(2048)
        .build());
    config = si.getConfiguration();
    System.out.println("stream name: " + config.getName() + ", max_bytes: " + config.getMaxBytes());

    // Create a durable consumer
    jsm.createConsumer("example-stream", ConsumerConfiguration.builder()
        .durable("example-consumer-name")
        .build());

    // Get information about all streams
    List<StreamInfo> streams = jsm.getStreams();
    for (StreamInfo info : streams) {
        System.out.println("stream name: " + info.getConfiguration().getName());
    }

    // Get information about all consumers
    List<ConsumerInfo> consumers = jsm.getConsumers("example-stream");
    for (ConsumerInfo ci : consumers) {
        System.out.println("consumer name: " + ci.getName());
    }

    // Delete a consumer
    jsm.deleteConsumer("example-stream", "example-consumer-name");

    // Delete a stream
    jsm.deleteStream("example-stream");
}
```

```javascript
import { AckPolicy, connect, Empty } from "../../src/mod.ts";

const nc = await connect();
const jsm = await nc.jetstreamManager();

// list all the streams, the `next()` function
// retrieves a paged result.
const streams = await jsm.streams.list().next();
streams.forEach((si) => {
    console.log(si);
});

// add a stream
const stream = "mystream";
const subj = `mystream.*`;
await jsm.streams.add({ name: stream, subjects: [subj] });

// publish a reg nats message directly to the stream
for (let i = 0; i < 10; i++) {
    nc.publish(`${subj}.a`, Empty);
}

// find a stream that stores a specific subject:
const name = await jsm.streams.find("mystream.A");

// retrieve info about the stream by its name
const si = await jsm.streams.info(name);

// update a stream configuration
si.config.subjects?.push("a.b");
await jsm.streams.update(name, si.config);

// get a particular stored message in the stream by sequence
// this is not associated with a consumer
const sm = await jsm.streams.getMessage(stream, { seq: 1 });
console.log(sm.seq);

// delete the 5th message in the stream, securely erasing it
await jsm.streams.deleteMessage(stream, 5);

// purge all messages in the stream, the stream itself
// remains.
await jsm.streams.purge(stream);
```

--------------------------------

### Add Replicated Stream with Placement Tags

Source: https://github.com/nats-io/nats.docs/blob/master/nats-concepts/jetstream/README.md

Creates a new JetStream with 3 replicas, specifically targeting servers that have the 'sync:always' tag.

```bash
nats stream add --replicas 3 --tag sync:always
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

### Obtain JetStream Context (Java)

Source: https://github.com/nats-io/nats.docs/blob/master/using-nats/developing-with-nats/js/context.md

Shows how to obtain the JetStream context and the JetStream Management context from a NATS connection object in Java. These contexts are used for JetStream-enabled operations and stream management, respectively.

```Java
// Getting the JetStream context
JetStream js = nc.jetStream();
// Getting the JetStream management context
JetStreamManagement jsm = nc.jetStreamManagement();
```

--------------------------------

### Obtain JetStream Context (C)

Source: https://github.com/nats-io/nats.docs/blob/master/using-nats/developing-with-nats/js/context.md

Demonstrates the initialization of a JetStream context in C using the NATS client library. It includes setting up JetStream options, such as asynchronous publish error handlers, and connecting to the NATS server.

```C
int main(int argc, char **argv)
{
    natsConnection      *conn  = NULL;
    natsOptions         *opts  = NULL;
    jsCtx               *js    = NULL;
    jsOptions           jsOpts;
    jsErrCode           jerr   = 0;
    volatile int        errors = 0;

    opts = parseArgs(argc, argv, usage);
    dataLen = (int) strlen(payload);

    s = natsConnection_Connect(&conn, opts);

    if (s == NATS_OK)
        s = jsOptions_Init(&jsOpts);

    if (s == NATS_OK)
    {
        if (async)
        {
            jsOpts.PublishAsync.ErrHandler           = _jsPubErr;
            jsOpts.PublishAsync.ErrHandlerClosure    = (void*) &errors;
        }
        s = natsConnection_JetStream(&js, conn, &jsOpts);
    }
    
    // Destroy all our objects to avoid report of memory leak
    jsCtx_Destroy(js);
    natsConnection_Destroy(conn);
    natsOptions_Destroy(opts);

    // To silence reports of memory still in used with valgrind
    nats_Close();

    return 0;
}
```