### List Consumers for a Stream with NATS CLI

Source: https://github.com/nats-io/nats.docs/blob/master/running-a-nats-service/nats_admin/jetstream_admin/consumers.md

Lists all consumers associated with a specific NATS stream, in this case, the 'ORDERS' stream. This command provides a quick overview of the active consumers.

```shell
nats con ls ORDERS
```

--------------------------------

### Add a new Stream using NATS CLI

Source: https://github.com/nats-io/nats.docs/blob/master/running-a-nats-service/configuration/jetstream-config/configuration_mgmt/nats-admin-cli.md

Creates a new stream in NATS based on a provided JSON configuration file. The configuration file can be generated from existing streams using the stream info command.

```shell
nats str add ORDERS --config orders.json
```

--------------------------------

### Report NATS Stream Status with NATS CLI

Source: https://github.com/nats-io/nats.docs/blob/master/running-a-nats-service/configuration/leafnodes/jetstream_leafnodes.md

This command retrieves and displays statistics for all streams within the specified JS domain, including message counts, sizes, and consumer information.

```bash
nats  --server nats://acc:acc@localhost:4222 stream report --js-domain hub
Obtaining Stream stats

╭──────────────────────────────────────────────────────────────────────────────────────────╮
│                                      Stream Report                                       │
├─────────────────────┬─────────┬───────────┬──────────┬───────┬──────┬─────────┬──────────┤
│ Stream              │ Storage │ Consumers │ Messages │ Bytes │ Lost │ Deleted │ Replicas │
├─────────────────────┼─────────┼───────────┼──────────┼───────┼──────┼─────────┼──────────┤
│ backup-test-leaf    │ File    │ 0         │ 1        │ 45 B  │ 0    │ 0       │          │
│ test                │ File    │ 0         │ 1        │ 45 B  │ 0    │ 0       │          │
│ aggregate-test-leaf │ File    │ 1         │ 1        │ 98 B  │ 0    │ 0       │          │
╰─────────────────────┴─────────┴───────────┴──────────┴───────┴──────┴─────────┴──────────╯

╭────────────────────────────────────────────────────────────────────────────────────╮
│                                 Replication Report                                 │
├─────────────────────┬────────┬──────────────┬───────────────┬────────┬─────┬───────┤
│ Stream              │ Kind   │ API Prefix   │ Source Stream │ Active │ Lag │ Error │
├─────────────────────┼────────┼──────────────┼───────────────┼────────┼─────┼───────┤
│ backup-test-leaf    │ Mirror │ $JS.leaf.API │ test          │ 1.85s  │ 0   │       │
│ aggregate-test-leaf │ Source │ $JS.leaf.API │ test          │ 1.85s  │ 0   │       │
╰─────────────────────┴────────┴──────────────┴───────────────┴────────┴─────┴───────╯
```

--------------------------------

### Report NATS Consumer Status with NATS CLI

Source: https://github.com/nats-io/nats.docs/blob/master/running-a-nats-service/configuration/leafnodes/jetstream_leafnodes.md

This command provides detailed status information for consumers within a specified stream and JS domain, including mode, acknowledgment policy, and message processing status.

```bash
nats  --server nats://acc:acc@localhost:4222 consumer report --js-domain hub
? Select a Stream aggregate-test-leaf
╭─────────────────────────────────────────────────────────────────────────────────────────────────────────────╮
│                          Consumer report for aggregate-test-leaf with 1 consumers                           │
├──────────┬──────┬────────────┬──────────┬─────────────┬─────────────┬─────────────┬───────────┬─────────────┤
│ Consumer │ Mode │ Ack Policy │ Ack Wait │ Ack Pending │ Redelivered │ Unprocessed │ Ack Floor │ Cluster     │
├──────────┼──────┼────────────┼──────────┼─────────────┼─────────────┼─────────────┼───────────┼─────────────┤
│ dur      │ Pull │ Explicit   │ 30.00s   │ 0           │ 0           │ 1 / 100%    │ 0         │ hub-server* │
╰──────────┴──────┴────────────┴──────────┴─────────────┴─────────────┴─────────────┴───────────┴─────────────╯
```

--------------------------------

### Add a new Consumer using NATS CLI

Source: https://github.com/nats-io/nats.docs/blob/master/running-a-nats-service/configuration/jetstream-config/configuration_mgmt/nats-admin-cli.md

Adds a new consumer to a specific stream using a JSON configuration file. The configuration can be derived from existing consumer metadata.

```shell
nats con add ORDERS NEW --config orders_new.json
```

--------------------------------

### Consuming Messages with a Pull-Based Consumer using `nats` CLI

Source: https://github.com/nats-io/nats.docs/blob/master/running-a-nats-service/nats_admin/jetstream_admin/consumers.md

Demonstrates how to fetch and acknowledge messages from a pull-based consumer named 'DISPATCH' on the 'ORDERS' stream using the `nats` CLI. It shows the output of received messages and the acknowledgment.

```shell
nats con next ORDERS DISPATCH
```

--------------------------------

### Create Pull-Based Consumer with NATS CLI

Source: https://github.com/nats-io/nats.docs/blob/master/running-a-nats-service/nats_admin/jetstream_admin/consumers.md

Demonstrates creating a pull-based consumer named 'NEW' for the 'ORDERS' stream using the NATS CLI. It specifies a sample rate of 100, filters messages on 'ORDERS.received', and sets a maximum of 20 deliveries. This consumer requires explicit message acknowledgments.

```shell
nats con add --sample 100
```

--------------------------------

### Create Advanced Pull-Based Consumer with NATS CLI

Source: https://github.com/nats-io/nats.docs/blob/master/running-a-nats-service/nats_admin/jetstream_admin/consumers.md

Creates an advanced pull-based consumer named 'DISPATCH' for the 'ORDERS' stream. This command specifies filtering on 'ORDERS.processed', explicit acknowledgments, pull mode, delivering all messages, a sample rate of 100, and a maximum of 20 deliveries.

```shell
nats con add ORDERS DISPATCH --filter ORDERS.processed --ack explicit --pull --deliver all --sample 100 --max-deliver 20
```

--------------------------------

### Create Consumer from JSON Configuration with NATS CLI

Source: https://github.com/nats-io/nats.docs/blob/master/running-a-nats-service/nats_admin/jetstream_admin/consumers.md

Adds a consumer to the 'ORDERS' stream using a configuration defined in a 'monitor.json' file. This method allows for complex configurations to be managed externally.

```shell
nats con add ORDERS --config monitor.json
```

--------------------------------

### Retrieve Messages from Durable Pull Consumer with NATS CLI

Source: https://github.com/nats-io/nats.docs/blob/master/running-a-nats-service/configuration/leafnodes/jetstream_leafnodes.md

This command fetches the next available message from a durable pull consumer named 'dur' on the 'aggregate-test-leaf' stream. It requires specifying the correct JS API prefix for the connection.

```bash
nats --server nats://import_client:import_client@localhost:4111 consumer next aggregate-test-leaf dur --js-api-prefix JS.acc@hub.API
[17:44:16] subj: test / tries: 1 / cons seq: 1 / str seq: 1 / pending: 0

Headers:

  Nats-Stream-Source: test:mSx7q4yJ 1

Data:


hello world

Acknowledged message
```

--------------------------------

### Create Push-Based Consumer with NATS CLI

Source: https://github.com/nats-io/nats.docs/blob/master/running-a-nats-service/nats_admin/jetstream_admin/consumers.md

Demonstrates creating a push-based consumer named 'MONITOR' for the 'ORDERS' stream. This consumer is configured to deliver messages to 'monitor.ORDERS', acknowledges messages with 'none', and replays the last message.

```shell
nats con add
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

### Add Stream with Source Configuration using NATS CLI

Source: https://github.com/nats-io/nats.docs/blob/master/nats-concepts/jetstream/source_and_mirror_example.md

Use the NATS CLI to add a stream defined by a JSON configuration file. This is the recommended approach for managing complex stream configurations like those involving sources.

```bash
nats stream add --config stream_with_sources.json
```

--------------------------------

### Create NATS Context

Source: https://github.com/nats-io/nats.docs/blob/master/using-nats/nats-tools/nsc/basics.md

Create a named context to simplify authentication when using the NATS CLI.

```shell
nats context add myuser --creds ~/.nkeys/creds/MyOperator/MyAccount/MyUser.creds
```

--------------------------------

### Requesting Messages from a Pull-Based Consumer Programmatically

Source: https://github.com/nats-io/nats.docs/blob/master/running-a-nats-service/nats_admin/jetstream_admin/consumers.md

Shows how to request the next message from a specific consumer ('DISPATCH' on 'ORDERS' stream) using the `nats req` command. This simulates a programmatic request and displays the received message.

```shell
nats req '$JS.API.CONSUMER.MSG.NEXT.ORDERS.DISPATCH' ''
```

--------------------------------

### Create Durable Pull Consumer with NATS CLI

Source: https://github.com/nats-io/nats.docs/blob/master/running-a-nats-service/configuration/leafnodes/jetstream_leafnodes.md

This command adds a durable pull consumer named 'dur' to the 'aggregate-test-leaf' stream within the 'hub' JS domain. It prompts for configuration details and confirms creation.

```bash
nats  --server nats://acc:acc@localhost:4222 consumer add  --js-domain hub
? Consumer name dur
? Delivery target (empty for Pull Consumers)
? Start policy (all, new, last, 1h, msg sequence) all
? Replay policy instant
? Filter Stream by subject (blank for all)
? Maximum Allowed Deliveries -1
? Maximum Acknowledgements Pending 0
? Select a Stream aggregate-test-leaf
Information for Consumer aggregate-test-leaf > dur created 2021-06-28T17:16:51-04:00

Configuration:

        Durable Name: dur
           Pull Mode: true
         Deliver All: true
          Ack Policy: Explicit
            Ack Wait: 30s
       Replay Policy: Instant
     Max Ack Pending: 20,000
   Max Waiting Pulls: 512

State:

   Last Delivered Message: Consumer sequence: 0 Stream sequence: 0
     Acknowledgment floor: Consumer sequence: 0 Stream sequence: 0
         Outstanding Acks: 0 out of maximum 20000
     Redelivered Messages: 0
     Unprocessed Messages: 1
            Waiting Pulls: 0 of maximum 512
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

### Edit an existing Stream using NATS CLI

Source: https://github.com/nats-io/nats.docs/blob/master/running-a-nats-service/configuration/jetstream-config/configuration_mgmt/nats-admin-cli.md

Updates an existing stream configuration to match the specifications defined in the provided JSON file. This ensures the stream complies with the desired settings.

```shell
nats str edit ORDERS --config orders.json
```

--------------------------------

### Create Advanced Push-Based Consumer with NATS CLI

Source: https://github.com/nats-io/nats.docs/blob/master/running-a-nats-service/nats_admin/jetstream_admin/consumers.md

Creates a push-based consumer named 'MONITOR' for the 'ORDERS' stream with specific configurations. It sets the delivery target to 'monitor.ORDERS', uses 'none' for acknowledgment policy, replays the last message, and filters are empty.

```shell
nats con add ORDERS MONITOR --ack none --target monitor.ORDERS --deliver last --replay instant --filter ''
```

--------------------------------

### Enable Monitoring via Command-Line Option

Source: https://github.com/nats-io/nats.docs/blob/master/running-a-nats-service/nats_admin/monitoring/README.md

Start the NATS server with a command-line flag to specify the monitoring port.

```bash
nats-server -m 8222
```

--------------------------------

### Force NATS Consumer Leader Election (Step Down)

Source: https://github.com/nats-io/nats.docs/blob/master/running-a-nats-service/configuration/clustering/jetstream_clustering/administration.md

This command initiates a leader election for a specified NATS JetStream consumer. Similar to stream leader election, it forces the current consumer leader to step down, allowing the RAFT group to elect a new leader. This is typically used for maintenance or to manage leadership.

```shell
nats consumer cluster step-down ORDERS NEW
```

--------------------------------

### Subscribing to a Push-Based Consumer's Delivery Subject

Source: https://github.com/nats-io/nats.docs/blob/master/running-a-nats-service/nats_admin/jetstream_admin/consumers.md

Illustrates how to subscribe to the delivery subject ('monitor.ORDERS') of a push-based consumer ('MONITOR' on 'ORDERS' stream) using the `nats sub` command. This allows real-time reception of messages pushed by the consumer.

```shell
nats sub monitor.ORDERS
```

--------------------------------

### Enable NATS Monitoring via Command Line

Source: https://github.com/nats-io/nats.docs/blob/master/running-a-nats-service/configuration/monitoring.md

Starts the NATS server with an HTTP monitoring endpoint on a specified port. This allows for real-time statistics and information retrieval. Ensure the port is not exposed publicly without proper security measures.

```bash
nats-server -m 8222
```