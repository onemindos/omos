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

### Disable NATS Reconnection (Python)

Source: https://github.com/nats-io/nats.docs/blob/master/_examples/reconnect_none.html

Connects to a NATS server using the asyncio client library with automatic reconnection disabled. The `allow_reconnect=False` parameter in the `connect` method ensures that the client will not try to re-establish a connection if it drops.

```Python
nc = NATS()
await nc.connect(
    servers=[
       "nats://demo.nats.io:1222",
       "nats://demo.nats.io:1223",
       "nats://demo.nats.io:1224"
       ],
    allow_reconnect=False,
    )

  # Do something with the connection

  await nc.close()
```

--------------------------------

### Gracefully Drain NATS Connection - Python

Source: https://github.com/nats-io/nats.docs/blob/master/_examples/drain_conn.html

This Python example uses the asyncio NATS client to demonstrate draining a connection. It sets up a request-reply pattern, sends multiple requests, and then gracefully drains the connection, printing the number of responses received.

```Python
import asyncio
from nats.aio.client import Client as NATS

async def example(loop):
    nc = NATS()

    await nc.connect("nats://127.0.0.1:4222", loop=loop)

    async def handler(msg):
        print("[Received] ", msg)
        await nc.publish(msg.reply, b'I can help')

        # Can check whether client is in draining state
        if nc.is_draining:
            print("Connection is draining")

    await nc.subscribe("help", "workers", cb=handler)
    await nc.flush()

    requests = []
    for i in range(0, 10):
        request = nc.request("help", b'help!', timeout=1)
        requests.append(request)

    # Wait for all the responses
    responses = []
    responses = await asyncio.gather(*requests)

    # Gracefully close the connection.
    await nc.drain()

    print("Received {} responses".format(len(responses)))
```

--------------------------------

### Handle NATS Connection Events in Python (Asyncio)

Source: https://github.com/nats-io/nats.docs/blob/master/_examples/connection_listener.html

Illustrates how to define and set asynchronous callback functions for connection events in the Python NATS client using asyncio. This includes callbacks for disconnected, reconnected, error, and closed events, configured via an options dictionary.

```Python
# Asyncio NATS client can be defined a number of event callbacks 
async def disconnected_cb(): 
    print("Got disconnected!") 

async def reconnected_cb(): 
    # See who we are connected to on reconnect. 
    print("Got reconnected to {url}".format(url=nc.connected_url.netloc))

async def error_cb(e): 
    print("There was an error: {}".format(e))

async def closed_cb(): 
    print("Connection is closed") 

# Setup callbacks to be notified on disconnects and reconnects 
options["disconnected_cb"] = disconnected_cb 
options["reconnected_cb"] = reconnected_cb 

# Setup callbacks to be notified when there is an error 
# or connection is closed. 
options["error_cb"] = error_cb 
options["closed_cb"] = closed_cb 

await nc.connect(**options)
```

--------------------------------

### Python (asyncio) NATS Client - Reconnect Time Wait Example

Source: https://github.com/nats-io/nats.docs/blob/master/_examples/reconnect_10s.html

Example of setting the reconnect time wait to 10 seconds for the Python asyncio NATS client.

```APIDOC
## Python (asyncio) NATS Client - Reconnect Time Wait

### Description
This example shows how to set the `reconnect_time_wait` parameter to 10 seconds when connecting with the Python asyncio NATS client.

### Method
`nc.connect`

### Endpoint
`nats://demo.nats.io:4222` (or your NATS server address)

### Parameters
#### Query Parameters
None

#### Request Body
None

### Request Example
```python
nc = NATS()
await nc.connect(
    servers=["nats://demo.nats.io:4222"],
    reconnect_time_wait=10,
    )

# Do something with the connection

await nc.close()
```

### Response
#### Success Response (200)
Connection established.

#### Response Example
None (connection is established programmatically)
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

### Check NATS Connection Status - Python (asyncio)

Source: https://github.com/nats-io/nats.docs/blob/master/_examples/connect_status.html

This Python snippet demonstrates connecting to NATS using `asyncio`, checking the connection status (`is_connected`, `is_reconnecting`), and closing the connection. It requires an asynchronous environment.

```Python
nc = NATS()

await nc.connect(
    servers=["nats://demo.nats.io:4222"],
    )

# Do something with the connection.

print("The connection is connected?", nc.is_connected)

while True:
  if nc.is_reconnecting:
    print("Reconnecting to NATS...")
    break
  await asyncio.sleep(1)

await nc.close()

print("The connection is closed?", nc.is_closed)
```

--------------------------------

### Note on Reconnect Buffer Support in Python

Source: https://github.com/nats-io/nats.docs/blob/master/using-nats/developing-with-nats/reconnect/buffer.md

States that the Asyncio NATS client currently does not implement a reconnect buffer feature.

```Python
# Asyncio NATS client currently does not implement a reconnect buffer
```

--------------------------------

### Synchronous Subscribe Limitation in Python (Asyncio)

Source: https://github.com/nats-io/nats.docs/blob/master/_examples/subscribe_sync.html

Indicates that the Asyncio NATS client for Python currently does not offer a synchronous subscribe API. Users must use asynchronous patterns for subscribing and receiving messages.

```Python
# Asyncio NATS client currently does not have a sync subscribe API
```

--------------------------------

### Connect to NATS with TLS using Python

Source: https://github.com/nats-io/nats.docs/blob/master/_examples/connect_tls.html

Connects to a NATS server securely using TLS in Python with the asyncio-nats library. It involves creating an SSL context, loading CA certificates, and loading client certificates and keys. The connection is then established using the configured SSL context.

```Python
nc = NATS()

ssl_ctx = ssl.create_default_context(purpose=ssl.Purpose.SERVER_AUTH)
ssl_ctx.load_verify_locations('ca.pem')
ssl_ctx.load_cert_chain(certfile='client-cert.pem',
                        keyfile='client-key.pem')
await nc.connect(io_loop=loop, tls=ssl_ctx)

await nc.connect(servers=["nats://demo.nats.io:4222"], tls=ssl_ctx)
# Do something with the connection.
```

--------------------------------

### Publish Messages Synchronously and Asynchronously

Source: https://github.com/nats-io/nats.docs/blob/master/using-nats/developing-with-nats/js/streams.md

Demonstrates how to publish messages to a JetStream subject using either synchronous blocking calls or asynchronous non-blocking calls.

```c
if (async)
    s = js_PublishAsync(js, subj, (const void*) payload, dataLen, NULL);
else
{
    jsPubAck *pa = NULL;
    s = js_Publish(&pa, js, subj, (const void*) payload, dataLen, NULL, &jerr);
    if (s == NATS_OK) {
        jsPubAck_Destroy(pa);
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

### Publish Bytes Message in Python (asyncio)

Source: https://github.com/nats-io/nats.docs/blob/master/_examples/publish_bytes.html

This Python code snippet utilizes the asyncio NATS client to connect to a NATS server and publish a byte string message to the 'updates' subject. It demonstrates asynchronous publishing.

```Python
nc = NATS()

await nc.connect(servers=["nats://demo.nats.io:4222"])

await nc.publish("updates", b'All is Well')
```

--------------------------------

### Python (Asyncio) - Subscribe

Source: https://github.com/nats-io/nats.docs/blob/master/_examples/subscribe_async.html

This Python example using asyncio demonstrates how to subscribe to a NATS subject and process messages asynchronously.

```APIDOC
## Python (Asyncio) - Subscribe

### Description
This Python example utilizes the asyncio library to connect to NATS, subscribe to the "updates" subject, and asynchronously process incoming messages. It also includes publishing a message to the same subject.

### Method
N/A (Client-side operation)

### Endpoint
N/A (Connects to a NATS server URL)

### Parameters
#### Query Parameters
- **nats_server_url** (string) - Required - The URL of the NATS server to connect to (e.g., "nats://demo.nats.io:4222").

### Request Example
```python
import asyncio
from nats.aio.client import Client as NATS

nc = NATS()

async def main():
    await nc.connect(servers=["nats://demo.nats.io:4222"])

    future = asyncio.Future()

    async def cb(msg):
        nonlocal future
        future.set_result(msg)

    await nc.subscribe("updates", cb=cb)
    await nc.publish("updates", b'All is Well')
    await nc.flush()

    msg = await asyncio.wait_for(future, 1)

asyncio.run(main())
```

### Response
#### Success Response (N/A)
This is a client-side example. The output will be the content of the received NATS message.

#### Response Example
(No direct response example, as it involves asynchronous message handling)
```

--------------------------------

### Set Connect Timeout in Python

Source: https://github.com/nats-io/nats.docs/blob/master/using-nats/developing-with-nats/connecting/connect_timeout.md

Set the `connect_timeout` to 10 seconds when establishing a connection with the Python NATS client. This parameter controls the maximum time to connect.

```python
nc = NATS()
await nc.connect(
  servers=["nats://demo.nats.io:4222"],
  connect_timeout=10)

# Do something with the connection

await nc.close()
```

--------------------------------

### Benchmark JetStream Throughput with Asynchronous Publication

Source: https://github.com/nats-io/nats.docs/blob/master/using-nats/nats-tools/nats_cli/natsbench.md

Measures throughput using asynchronous publication to file storage. The stream is removed first to ensure a clean test environment.

```bash
nats stream rm -f benchstream
nats bench js pub async jsfoo --create
```

--------------------------------

### Connect to NATS Cluster in Python

Source: https://github.com/nats-io/nats.docs/blob/master/using-nats/developing-with-nats/connecting/cluster.md

Illustrates how to connect to a NATS cluster using Python. The code initializes the NATS client, connects using a list of servers, and ensures the connection is closed.

```Python
nc = NATS()
await nc.connect(servers=[
   "nats://127.0.0.1:1222",
   "nats://127.0.0.1:1223",
   "nats://127.0.0.1:1224"
   ])

# Do something with the connection

await nc.close()
```

--------------------------------

### Python Client Connection Event Callbacks

Source: https://github.com/nats-io/nats.docs/blob/master/using-nats/developing-with-nats/events/events.md

Explains how to define and assign callback functions for connection events like disconnection, reconnection, errors, and closure in the NATS Python client.

```APIDOC
## Python Client Connection Event Callbacks

### Description
This example shows how to set up asynchronous callback functions for various connection events in the NATS Python client. Callbacks can be defined for disconnection, successful reconnection, errors, and connection closure.

### Method
N/A (Configuration via options dictionary)

### Endpoint
N/A

### Parameters
#### Connection Options (`options` dictionary)
- **disconnected_cb** (async function) - Callback executed upon disconnection.
- **reconnected_cb** (async function) - Callback executed upon successful reconnection.
- **error_cb** (async function) - Callback executed when an error occurs.
- **closed_cb** (async function) - Callback executed when the connection is closed.

### Request Example
```python
async def disconnected_cb():
    print("Got disconnected!")

async def reconnected_cb():
    # See who we are connected to on reconnect.
    print("Got reconnected to {url}".format(url=nc.connected_url.netloc))

async def error_cb(e):
    print("There was an error: {}".format(e))

async def closed_cb():
    print("Connection is closed")

# Setup callbacks to be notified on disconnects and reconnects
options["disconnected_cb"] = disconnected_cb
options["reconnected_cb"] = reconnected_cb

# Setup callbacks to be notified when there is an error
# or connection is closed.
options["error_cb"] = error_cb
options["closed_cb"] = closed_cb

# Assuming 'nc' is initialized elsewhere and 'options' is a dictionary
# await nc.connect(**options)
```

### Response
No direct response body. The callbacks are executed asynchronously when the corresponding events occur.
```

--------------------------------

### Retrieve JetStream Server Report

Source: https://github.com/nats-io/nats.docs/blob/master/running-a-nats-service/configuration/leafnodes/jetstream_leafnodes.md

Fetches a summary report of JetStream activity across connected servers. This command requires administrative credentials to view domain-wide metrics.

```bash
nats --server nats://admin:admin@localhost:4222 server report jetstream
```

--------------------------------

### Connect to Default NATS Server (Python)

Source: https://github.com/nats-io/nats.docs/blob/master/using-nats/developing-with-nats/connecting/default_server.md

Connects to the default NATS server using Python's NATS client library. The connection is managed asynchronously and closed upon completion.

```Python
nc = NATS()
await nc.connect()

# Do something with the connection

await nc.close()
```

--------------------------------

### Python Asyncio NATS - Set Pending Limits

Source: https://github.com/nats-io/nats.docs/blob/master/_examples/slow_pending_limits.html

Illustrates setting pending message and byte limits when subscribing in Python using the asyncio NATS client.

```APIDOC
## POST /api/subscriptions/pending_limits

### Description
Allows setting pending message and byte limits for a subscription in the Python asyncio NATS client. These limits help manage memory usage and prevent excessive buffering.

### Method
POST

### Endpoint
/api/subscriptions

### Parameters
#### Query Parameters
- **subject** (string) - Required - The subject to subscribe to.
- **cb** (callback function) - Required - The callback function to execute when a message is received.
- **pending_msgs_limit** (integer) - Optional - The maximum number of messages to buffer. Defaults to system limits if not specified.
- **pending_bytes_limit** (integer) - Optional - The maximum number of bytes to buffer. Defaults to system limits if not specified.

### Request Example
```python
await nc.subscribe("updates", cb=callback, pending_msgs_limit=1000, pending_bytes_limit=5*1024*1024)
```

### Response
#### Success Response (200)
- **subscription_id** (string) - The ID of the created subscription.

#### Response Example
```json
{
  "subscription_id": "sub-12345"
}
```
```