### Implement Request-Reply with NATS

Source: https://github.com/nats-io/nats.docs/blob/master/using-nats/developing-with-nats/sending/replyto.md

This snippet shows how to send a request and handle a single reply using the NATS request-reply pattern. It involves creating a unique inbox for replies, subscribing to it, publishing the request with the inbox as the reply-to subject, and then waiting for and processing the response.

```Go
nc, err := nats.Connect("demo.nats.io")
if err != nil {
    log.Fatal(err)
}
deffer nc.Close()

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

```Java
Connection nc = Nats.connect("nats://demo.nats.io:4222");

// set up a listener for "time" requests
Dispatcher d = nc.createDispatcher(msg -> {
    System.out.println("Received time request");
    nc.publish(msg.getReplyTo(), ("" + System.currentTimeMillis()).getBytes());
});
d.subscribe("time");

// make a subject for replies and subscribe to that
String replyToThis = NUID.nextGlobal();
Subscription sub = nc.subscribe(replyToThis);

// publish to the "time" subject with reply-to subject that was set up
nc.publish("time", replyToThis, null);

// wait for a response
Message msg = sub.nextMessage(1000);

// look at the response
long time = Long.parseLong(new String(msg.getData()));
System.out.println(new Date(time));

nc.close();
```

```JavaScript
// set up a subscription to process the request
const sc = StringCodec();
nc.subscribe("time", {
  callback: (_err, msg) => {
    msg.respond(sc.encode(new Date().toLocaleTimeString()));
  },
});

// create a subscription subject that the responding send replies to
const inbox = createInbox();
const sub = nc.subscribe(inbox, {
  max: 1,
  callback: (_err, msg) => {
    t.log(`the time is ${sc.decode(msg.data)}`);
  },
});

nc.publish("time", Empty, { reply: inbox });
```

```Python
nc = NATS()

future = asyncio.Future()

async def sub(msg):
  nonlocal future
  future.set_result(msg)

await nc.connect(servers=["nats://demo.nats.io:4222"])
await nc.subscribe("time", cb=sub)

unique_reply_to = nc.new_inbox()
await nc.publish("time", b'', unique_reply_to)
```

--------------------------------

### Implement NATS Request-Reply Pattern

Source: https://github.com/nats-io/nats.docs/blob/master/using-nats/developing-with-nats/sending/replyto.md

Examples of sending a request and waiting for a reply using NATS. These snippets demonstrate the core pattern of publishing to a subject with a reply-to inbox and processing the incoming response.

```csharp
await using var client = new NatsClient();
await client.ConnectAsync();
string inbox = client.Connection.NewInbox();
await using INatsSub<string> timeSub = await client.Connection.SubscribeCoreAsync<string>("time");
Task responderTask = Task.Run(async () => {
    await foreach (var msg in timeSub.Msgs.ReadAllAsync()) {
        await msg.ReplyAsync<DateTimeOffset>(DateTimeOffset.UtcNow);
    }
});
await using INatsSub<DateTimeOffset> inboxSub = await client.Connection.SubscribeCoreAsync<DateTimeOffset>(inbox);
await client.PublishAsync(subject: "time", replyTo: inbox);
NatsMsg<DateTimeOffset> reply = await inboxSub.Msgs.ReadAsync();
Console.WriteLine($"The current date and time is: {reply.Data:R}");
await inboxSub.UnsubscribeAsync();
await timeSub.UnsubscribeAsync();
await responderTask;
```

```ruby
require 'nats/client'
require 'fiber'
NATS.start(servers:["nats://127.0.0.1:4222"]) do |nc|
  Fiber.new do
    f = Fiber.current
    nc.subscribe("time") do |msg, reply|
      f.resume msg
    end
    nc.publish("time", 'example', NATS.create_inbox)
    msg = Fiber.yield
    puts "Reply: #{msg}"
  end.resume
end
```

```c
natsConnection *conn = NULL;
natsStatus s = NATS_OK;
s = natsConnection_ConnectTo(&conn, NATS_DEFAULT_URL);
if (s == NATS_OK)
    s = natsConnection_PublishRequestString(conn, "request", "reply", "this is the request");
natsConnection_Destroy(conn);
```

--------------------------------

### NATS Request-Reply (Ruby)

Source: https://github.com/nats-io/nats.docs/blob/master/using-nats/developing-with-nats/sending/request_reply.md

Shows a basic implementation of the request-reply pattern in Ruby with NATS. It subscribes to 'time' and sends a generic 'response' back, then makes a request and prints the reply.

```Ruby
require 'nats/client'
require 'fiber'

NATS.start(servers:["nats://127.0.0.1:4222"]) do |nc|
  nc.subscribe("time") do |msg, reply|
    nc.publish(reply, "response")
  end

  Fiber.new do
    # Use the response
    msg = nc.request("time", "")
    puts "Reply: #{msg}"
  end.resume
end
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

### Publish Request and Subscribe to Reply (Java)

Source: https://github.com/nats-io/nats.docs/blob/master/_examples/publish_with_reply.html

This Java code snippet shows how to send a request to the 'time' subject and listen for a reply on a generated unique inbox subject. It includes establishing a NATS connection, creating a subscription, publishing the request, and processing the received message. The connection is closed at the end.

```Java
Connection nc = Nats.connect("nats://demo.nats.io:4222");

// Create a unique subject name
String uniqueReplyTo = NUID.nextGlobal();

// Listen for a single response
Subscription sub = nc.subscribe(uniqueReplyTo);
sub.unsubscribe(1);

// Send the request
nc.publish("time", uniqueReplyTo, null);

// Read the reply
Message msg = sub.nextMessage(Duration.ofSeconds(1));

// Use the response
System.out.println(new String(msg.getData(), StandardCharsets.UTF_8));

// Close the connection
nc.close();
```

--------------------------------

### Implement NATS Request-Reply Pattern

Source: https://github.com/nats-io/nats.docs/blob/master/_examples/request_reply.html

Demonstrates how to send a request to a subject and handle the incoming response. These examples cover connection management, request submission with timeouts, and response processing.

```Go
nc, err := nats.Connect("demo.nats.io")
if err != nil {
	log.Fatal(err)
}
defer nc.Close()

// Send the request
msg, err := nc.Request("time", nil, time.Second)
if err != nil {
	log.Fatal(err)
}

// Use the response
log.Printf("Reply: %s", msg.Data)

// Close the connection
nc.Close()
```

```Java
Connection nc = Nats.connect("nats://demo.nats.io:4222");

// Send the request
Message msg = nc.request("time", null, Duration.ofSeconds(1));

// Use the response
System.out.println(new String(msg.getData(), StandardCharsets.UTF_8));

// Close the connection
nc.close();
```

```Python
nc = NATS()

async def sub(msg):
  await nc.publish(msg.reply, b'response')

await nc.connect(servers=["nats://demo.nats.io:4222"])
await nc.subscribe("time", cb=sub)

# Send the request
try:
  msg = await nc.request("time", b'', timeout=1)
  # Use the response
  print("Reply:", msg)
except asyncio.TimeoutError:
  print("Timed out waiting for response")
```

```TypeScript
let msg = await nc.request('time', 1000);
t.log('the time is', msg.data);
nc.close();
```

--------------------------------

### Publish Request and Subscribe to Reply (Python)

Source: https://github.com/nats-io/nats.docs/blob/master/_examples/publish_with_reply.html

This Python code snippet uses asyncio to publish a request to the 'time' subject and asynchronously wait for a reply. It defines an asynchronous callback function to handle the message and uses a Future to signal when the reply is received. A unique inbox subject is generated for the reply.

```Python
nc = NATS()

future = asyncio.Future()

async def sub(msg):
  nonlocal future
  future.set_result(msg)

await nc.connect(servers=["nats://demo.nats.io:4222"])
await nc.subscribe("time", cb=sub)

unique_reply_to = new_inbox()
await nc.publish_request("time", unique_reply_to, b'')

# Use the response
msg = await asyncio.wait_for(future, 1)
print("Reply:", msg)
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

### Respond to Time Request

Source: https://github.com/nats-io/nats.docs/blob/master/using-nats/developing-with-nats/receiving/reply.md

Listens for messages on the 'time' subject and responds with the current server time. This example showcases the request-reply pattern in NATS.

```Go
nc, err := nats.Connect("demo.nats.io")
if err != nil {
    log.Fatal(err)
}
deferr nc.Close()

// Subscribe
sub, err := nc.SubscribeSync("time")
if err != nil {
    log.Fatal(err)
}

// Read a message
msg, err := sub.NextMsg(10 * time.Second)
if err != nil {
    log.Fatal(err)
}

// Get the time
timeAsBytes := []byte(time.Now().String())

// Send the time as the response.
msg.Respond(timeAsBytes)
```

```Java
Connection nc = Nats.connect("nats://demo.nats.io:4222")

// Subscribe to the "time" subject and reply with the current time
Subscription sub = nc.subscribe("time")

// Read a message
Message msg = sub.nextMessage(Duration.ZERO)

// Get the time
Calendar cal = Calendar.getInstance()
SimpleDateFormat sdf = new SimpleDateFormat("HH:mm:ss")
byte[] timeAsBytes = sdf.format(cal.getTime()).getBytes(StandardCharsets.UTF_8)

// Send the time to the reply to subject
nc.publish(msg.getReplyTo(), timeAsBytes)

// Flush and close the connection
nc.flush(Duration.ZERO)
nc.close()
```

```JavaScript
const sc = StringCodec()
// set up a subscription to process a request
const sub = nc.subscribe("time")
for await (const m of sub) {
  m.respond(sc.encode(new Date().toLocaleDateString()))
}
```

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

```C#
// dotnet add package NATS.Net
using NATS.Net

await using var client = new NatsClient()

// Subscribe to the "time" subject and reply with the current time
await foreach (var msg in client.SubscribeAsync<string>("time"))
{
    await msg.ReplyAsync(DateTime.Now)
}
```

```Ruby
require 'nats/client'
require 'fiber'

NATS.start(servers:["nats://127.0.0.1:4222"]) do |nc|
  Fiber.new do
    f = Fiber.current

    nc.subscribe("time") do |msg, reply|
      f.resume Time.now
    end

    nc.publish("time", 'What is the time?', NATS.create_inbox)

    # Use the response
    msg = Fiber.yield
    puts "Reply: #{msg}"

  end.resume
end
```

```C
natsConnection      *conn      = NULL;
natsSubscription    *sub       = NULL;
natsMsg             *msg       = NULL;
natsStatus          s          = NATS_OK;

s = natsConnection_ConnectTo(&conn, NATS_DEFAULT_URL)

// Subscribe
if (s == NATS_OK)
    s = natsConnection_SubscribeSync(&sub, conn, "time")

// Wait for messages
if (s == NATS_OK)
    s = natsSubscription_NextMsg(&msg, sub, 10000)

if (s == NATS_OK)
{
    char buf[64]

    snprintf(buf, sizeof(buf), "%lld", nats_Now())

    // Send the time as a response
    s = natsConnection_Publish(conn, natsMsg_GetReply(msg), buf, (int) strlen(buf))

    // Destroy message that was received
    natsMsg_Destroy(msg)
}

(...)

// Destroy objects that were created
natsSubscription_Destroy(sub)
natsConnection_Destroy(conn)
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

### Publish Request and Subscribe to Reply (Node.js)

Source: https://github.com/nats-io/nats.docs/blob/master/_examples/publish_with_reply.html

This Node.js (JavaScript) example demonstrates publishing a request to the 'time' subject and setting up a subscription to handle replies. It utilizes NATS.createInbox() for the reply subject and closes the connection after receiving the response. The example includes a listener for the 'time' subject that publishes the current time to the reply subject.

```JavaScript
let nc = NATS.connect({url: "nats://demo.nats.io:4222"});
// set up a subscription to process the request
nc.subscribe('time', (msg, reply) => {
    if(reply) {
        nc.publish(reply, new Date().toLocaleTimeString());
    }
});

// create a subscription subject that the responding send replies to
let inbox = NATS.createInbox();
nc.subscribe(inbox, {max: 1}, (msg) => {
    t.log('the time is', msg);
    nc.close();
});

nc.publish('time', "", inbox);
```

--------------------------------

### NATS Request-Reply with Timeout (Java)

Source: https://github.com/nats-io/nats.docs/blob/master/using-nats/developing-with-nats/sending/request_reply.md

Illustrates the request-reply pattern in Java using the NATS client. It sets up a dispatcher to listen for 'time' requests, publishes the current time as a response, and makes a request with a one-second timeout.

```Java
Connection nc = Nats.connect("nats://demo.nats.io:4222");

// set up a listener for "time" requests
Dispatcher d = nc.createDispatcher(msg -> {
    System.out.println("Received time request");
    nc.publish(msg.getReplyTo(), ("" + System.currentTimeMillis()).getBytes());
});
d.subscribe("time");

// make a request to the "time" subject and wait 1 second for a response
Message msg = nc.request("time", null, Duration.ofSeconds(1));

// look at the response
long time = Long.parseLong(new String(msg.getData()));
System.out.println(new Date(time));

nc.close()
```

--------------------------------

### NATS Request-Reply with Timeout (C#)

Source: https://github.com/nats-io/nats.docs/blob/master/using-nats/developing-with-nats/sending/request_reply.md

Demonstrates the request-reply pattern in C# using the NATS.Net library. It sets up a subscription to the 'time' subject, responds with the current time, and then makes a request, printing the received data.

```C#
// dotnet add package NATS.Net
using NATS.Net;

await using var client = new NatsClient();

using CancellationTokenSource cts = new();

// Process the time messages in a separate task
Task subscription = Task.Run(async () =>
{
    await foreach (var msg in client.SubscribeAsync<string>("time", cancellationToken: cts.Token))
    {
        await msg.ReplyAsync(DateTimeOffset.Now);
    }
});

// Wait for the subscription task to be ready
await Task.Delay(1000);

var reply = await client.RequestAsync<DateTimeOffset>("time");

Console.WriteLine($"Reply: {reply.Data:O}");

await cts.CancelAsync();
await subscription;

// Output:
// Reply: 2024-10-23T05:20:55.0000000+01:00
```

--------------------------------

### Subscribe with Reply - Ruby

Source: https://github.com/nats-io/nats.docs/blob/master/_examples/subscribe_w_reply.html

This Ruby code snippet shows how to subscribe to a 'time' subject and respond with the current time. It uses fibers to handle asynchronous operations, subscribing to a subject, publishing a request, and then processing the reply.

```Ruby
require 'nats/client'
require 'fiber'

NATS.start(servers=["nats://127.0.0.1:4222"]) do |nc|
  Fiber.new do
    f = Fiber.current

    nc.subscribe("time") do |msg, reply|
      f.resume Time.now
    end

    nc.publish("time", 'What is the time?', NATS.create_inbox)
    # Use the response
    msg = Fiber.yield
    puts "Reply: #{msg}"
  end.resume
end
```

--------------------------------

### NATS Request-Reply with Timeout (Python)

Source: https://github.com/nats-io/nats.docs/blob/master/using-nats/developing-with-nats/sending/request_reply.md

Provides a Python example for the NATS request-reply pattern using asyncio. It defines a handler for 'time' requests, connects to the NATS server, subscribes to 'time', and sends a request with a 1-second timeout, handling potential timeouts.

```Python
nc = NATS()

async def sub(msg):
  await nc.publish(msg.reply, b'response')

await nc.connect(servers=["nats://demo.nats.io:4222"])
await nc.subscribe("time", cb=sub)

# Send the request
try:
  msg = await nc.request("time", b'', timeout=1)
  # Use the response
  print("Reply:", msg)
except asyncio.TimeoutError:
  print("Timed out waiting for response")
```

--------------------------------

### Request-Reply Pattern Note

Source: https://github.com/nats-io/nats.docs/blob/master/running-a-nats-service/configuration/securing_nats/authorization.md

Important considerations for maintaining request-reply patterns when configuring authorization, especially with wildcards.

```APIDOC
**Important Note** It is important to not break request-reply patterns. In some cases (as shown [below](authorization.md#variables)) you need to add rules for the `_INBOX.>` pattern. If an unauthorized client publishes or attempts to subscribe to a subject that has not been _allow listed_, the action fails and is logged at the server, and an error message is returned to the client. The [allow responses](authorization.md#allow-responses-map) option can simplify this.
```