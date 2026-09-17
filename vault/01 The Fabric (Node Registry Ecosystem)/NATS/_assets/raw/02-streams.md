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

### Create JetStream Stream

Source: https://github.com/nats-io/nats.docs/blob/master/running-a-nats-service/configuration/leafnodes/jetstream_leafnodes.md

Adds a new stream to the JetStream domain associated with the current connection. The command prompts for configuration details like storage backend and retention policies.

```bash
nats --server nats://acc:acc@localhost:4222 stream add
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

### JetStream Server Limits Configuration

Source: https://github.com/nats-io/nats.docs/blob/master/running-a-nats-service/configuration/README.md

Configure JetStream server limits such as max file store size, max acknowledgments pending, and duplicate window. Ensure these settings align with your expected load and data retention policies.

```hcl
jetstream {
  store_dir: nats
  max_file_store: 10G
  limits {
      max_ack_pending: 10000
      duplicate_window: 600s
  }
}
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

### C JetStream Consumer Configuration and Usage

Source: https://github.com/nats-io/nats.docs/blob/master/using-nats/developing-with-nats/js/consumers.md

Demonstrates C NATS JetStream consumer setup, including options for global delivery, synchronous/asynchronous reception, pull subscriptions, flow control, and message counting. It shows how to initialize JetStream context, configure streams and consumers, and handle message delivery callbacks.

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
    printf("Async error: %u - %s\n", err, natsStatus_GetText(err));

    natsSubscription_GetDropped(sub, (int64_t*) &dropped);
}

int main(int argc, char **argv)
{
    natsConnection      *conn  = NULL;
    natsStatistics      *stats = NULL;
    natsOptions         *opts  = NULL;
    natsSubscription    *sub   = NULL;
    natsMsg             *msg   = NULL;
    jsCtx               *js    = NULL;
    jsErrCode           jerr   = 0;
    jsOptions           jsOpts;
    jsSubOptions        so;
    natsStatus          s;
    bool                delStream = false;

    opts = parseArgs(argc, argv, usage);

    printf("Created %s subscription on '%s'.\n",
        (pull ? "pull" : (async ? "asynchronous" : "synchronous")), subj);

    s = natsOptions_SetErrorHandler(opts, asyncCb, NULL);

    if (s == NATS_OK)
        s = natsConnection_Connect(&conn, opts);

    if (s == NATS_OK)
        s = jsOptions_Init(&jsOpts);

    if (s == NATS_OK)
        s = jsSubOptions_Init(&so);
    if (s == NATS_OK)
    {
        so.Stream = stream;
        so.Consumer = durable;
        if (flowctrl)
        {
            so.Config.FlowControl = true;
            so.Config.Heartbeat = (int64_t)1E9;
        }
    }

    if (s == NATS_OK)
        s = natsConnection_JetStream(&js, conn, &jsOpts);

    if (s == NATS_OK)
    {
        jsStreamInfo    *si = NULL;

        // First check if the stream already exists.
        s = js_GetStreamInfo(&si, js, stream, NULL, &jerr);
        if (s == NATS_NOT_FOUND)
        {
            jsStreamConfig  cfg;

            // Since we are the one creating this stream, we can delete at the end.
            delStream = true;

            // Initialize the configuration structure.
            jsStreamConfig_Init(&cfg);
            cfg.Name = stream;
            // Set the subject
            cfg.Subjects = (const char*[1]){subj};
            cfg.SubjectsLen = 1;
            // Make it a memory stream.

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

### Configure JetStream for Always Sync

Source: https://github.com/nats-io/nats.docs/blob/master/nats-concepts/jetstream/README.md

Configures a NATS server to use the 'always' sync interval for JetStream, ensuring data is fsynced after every message. This is applied to servers tagged with 'sync:always'.

```hcl
server_tags: ["sync:always"]

jetstream {
    sync_interval: always
}
```

--------------------------------

### JetStream Stream Subject Transform Configuration

Source: https://github.com/nats-io/nats.docs/blob/master/nats-concepts/subject_mapping.md

Defines subject transformation rules within a JetStream stream configuration. It maps 'orders.local.*' to 'orders.{{wildcard(1)}}' for messages imported into the stream and also configures a republish rule.

```json
{
  "name": "orders",
  "subjects": [ "orders.local.*"],
  "subject_transform":{"src":"orders.local.*","dest":"orders.{{wildcard(1)}}"},
  "retention": "limits",
  ...
  "republish": {
    "src": "orders.*",
    "dest": "orders.trace.{{wildcard(1)}}"
  },

```

--------------------------------

### Add Replicated Stream with Placement Tags

Source: https://github.com/nats-io/nats.docs/blob/master/nats-concepts/jetstream/README.md

Creates a new JetStream with 3 replicas, specifically targeting servers that have the 'sync:always' tag.

```bash
nats stream add --replicas 3 --tag sync:always
```

--------------------------------

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

### Get NATS JetStream Stream Info

Source: https://github.com/nats-io/nats.docs/blob/master/nats-concepts/jetstream/js_walkthrough.md

Retrieves and displays detailed information about the specified JetStream stream, including its subjects, configuration, limits, and current state.

```shell
nats stream info my_stream
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

--------------------------------

### Report NATS JetStream Status

Source: https://github.com/nats-io/nats.docs/blob/master/running-a-nats-service/nats_admin/jetstream_admin/replication.md

Displays the current status and statistics of all configured JetStream streams.

```shell
nats s report
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

--------------------------------

### Report JetStream status and replication

Source: https://github.com/nats-io/nats.docs/blob/master/running-a-nats-service/configuration/leafnodes/jetstream_leafnodes.md

Retrieves a report of all streams and their replication status within a specific JetStream domain. This helps verify that mirrors and sources are active and synchronized.

```bash
nats --server nats://acc:acc@localhost:4222 stream report --js-domain hub
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

### Create a NATS JetStream Stream

Source: https://github.com/nats-io/nats.docs/blob/master/running-a-nats-service/nats_admin/jetstream_admin/streams.md

Demonstrates how to create a new stream using interactive prompts, command-line flags, or a JSON configuration file.

```shell
nats str add ORDERS
nats str add ORDERS --subjects "ORDERS.*" --ack --max-msgs=-1 --max-bytes=-1 --max-age=1y --storage file --retention limits --max-msg-size=-1 --discard old --dupe-window="0s" --replicas 1
nats str add ORDERS --config orders.json
```