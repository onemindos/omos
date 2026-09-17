# RethinkConn 2026 - NATS Fundementals Clustering Exercise

## Create Your Project Directory

Unzip the rethink-cluster-lab.zip file and change directory (`cd`) into the `rethink-cluster-lab` directory.

## Configure and Deploy The NATS Cluster

### Launch the NATS Servers

When you start the first server, you'll see some errors when trying to connect to routes. That's OK - the other servers haven't been started yet. Once all three servers are started, you'll see INFO messages in the console informing you that the routes have connected.

Launch Server n1:

```
nats-server -c ./config/n1.conf
```

Launch Server n2:

```
nats-server -c ./config/n2.conf
```

Launch Server n3:

```
nats-server -c ./config/n3.conf
```

### Create NATS Contexts

A NATS context is a named connection profile that stores server, authentication, and CLI settings so you can easily switch between NATS environments. We'll create and application context (`app`) to publish and subscribe NATS messages and exercise our cluster. We'll also create a system account context (`sys`) because some of the commands we'll be using (like `nats server ls`) require system privileges to run them.

```
nats context save sys \
  --server nats://127.0.0.1:4222 \
  --user sys \
  --password sys

nats context save app \
  --server nats://127.0.0.1:4222 \
  --user app \
  --password app
```

### Verify your contexts:

```
nats context ls
```

You should see an `app` context and `sys` context.

### Verify your cluster:

This command requires system privileges. Change your context to the `sys` context:

```
nats context select sys
```

And verify your cluster:

```
nats server ls
```

You should see all three servers and a cluster overview.

## Excercise the Cluster

### Change to the `app` Context

Use the System Account for administration and observability, not application messaging. Applications should always connect using their own account and credentials.

```
nats context select app
```

### Start a Listener on Each Server

Open listeners in their own terminal windows. You will need 3 for this exercise. We'll start listeners on subject "rethink.test".
In the first terminal window, start a listener using server `n1`:
The `app` context specifies the connection URL and port as `nats://localhost:4222` which points to the `n1` server. 

```
nats sub rethink.test
```

In the second terminal window, start a listener using server `n2`:
Here's where we take advantage of the modularity of NATS configuration files. Remember, we are using the `app` context which points to `nats://localhost:4222` (`n1`) but, we want to start this listener on `n2`. 

Fortunately, when we issue commands like `nats sub` we can use flags to override the values in our context. In this case we'll override the the URL with `nats://localhost:4223` and, as long as the account credentials match those of the server, the listener will start.

How do we know that the account credentials match? We only used a single accounts.conf file to define userids and passwords and used `include` to apply them to all three individual server configuration files. It's not 3 definitions, it's one definition used 3 times.

```
nats sub -s nats://localhost:4223 rethink.test
```

In the third terminal window, start a listener on `n3` using the `nats://localhost:4224` url in the -s flag to overide the context URL.

```
nats sub -s nats://localhost:4224 rethink.test
```

### Publish Messages

You'll need to open one more terminal window to publish messages.

Let's start by publishing a message on `n1`

```
nats pub rethink.test '1 - Hello RethinkConn!'
```

What just happened? You published a message to `n1` and not only did `n1` get the message but, `n2` and `n3` had active subscribers on the topic and they received the message too. The gossip protocol of the NATS servers propogated the interest in that subject to the other servers without you having to change your `nats pub` command at all.

NATS client libraries (like Go, Python, Java...) are also aware of multiple servers in the cluster and they learn about them the same way using the gossip protocol. So, if a client is connected to a server that goes down, it will reconnect to another server of which it was made aware by the gossip protocol. And it requires no additional application programming to implement that behavior, NATS handles it transparently.

Now, lets publish a message on `n3`. In the terminal window in which you just published the first message:

```
nats pub -s nats://localhost:4224 rethink.test '3 - Hello RethinkConn!!!'
```

## Monitoring Endpoints

Each NATS server exposes their monitoring endpoints on the specified monitoring port. Paste the links in your browser and explore the information available.

### Monitoring Endpoints for Server `n1`

```
http://localhost:8222
```

### Monitoring Endpoints for Server `n2`

```
http://localhost:8223
```

### Monitoring Endpoints for Server `n3`

```
http://localhost:8224
```
