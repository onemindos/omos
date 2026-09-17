### Configure NATS for Docker WebSocket Support

Source: https://github.com/nats-io/nats.docs/blob/master/running-a-nats-service/configuration/websocket/websocket_conf.md

Provides the minimal configuration required to enable WebSocket support in a NATS Docker container and the corresponding command to execute the container with that configuration.

```text
websocket 
{
     port: 8080
     no_tls: true
}
```

```bash
docker run -it --rm  -v /tmp:/container -p 8080:8080 nats -c /container/nats.conf
```

--------------------------------

### Configure NATS WebSocket Server

Source: https://github.com/nats-io/nats.docs/blob/master/running-a-nats-service/configuration/websocket/websocket_conf.md

This configuration block enables WebSocket support for a NATS server. It allows specifying the listen host and port, TLS certificates, CORS settings, and various authorization parameters. Ensure that TLS is configured correctly for secure connections, or disable it for testing environments. The `allowed_origins` and `same_origin` directives control cross-origin resource sharing.

```text
websocket {
    # Specify a host and port to listen for websocket connections
    #
    # listen: "host:port"

    # It can also be configured with individual parameters,
    # namely host and port.
    #
    # host: "hostname"
    port: 443

    # This will optionally specify what host:port for websocket
    # connections to be advertised in the cluster.
    #
    # advertise: "host:port"

    # TLS configuration is required by default
    #
    tls {
      cert_file: "/path/to/cert.pem"
      key_file: "/path/to/key.pem"
    }

    # For test environments, you can disable the need for TLS
    # by explicitly setting this option to `true`
    #
    # no_tls: true

    # [Cross-origin resource sharing option](https://developer.mozilla.org/en-US/docs/Web/HTTP/CORS).
    #
    # IMPORTANT! This option is used only when the http request presents an Origin
    # header, which is the case for web browsers. If no Origin header is present,
    # this check will not be performed.
    #
    # When set to `true`, the HTTP origin header must match the request’s hostname.
    # The default is `false`.
    #
    # same_origin: true

    # [Cross-origin resource sharing option](https://developer.mozilla.org/en-US/docs/Web/HTTP/CORS).
    #
    # IMPORTANT! This option is used only when the http request presents an Origin
    # header, which is the case for web browsers. If no Origin header is present,
    # this check will not be performed.
    #
    # List of accepted origins. When empty, and `same_origin` is `false`, clients from any origin are allowed to connect.
    # This list specifies the only accepted values for the client's request Origin header. The scheme,
    # host and port must match. By convention, the absence of TCP port in the URL will be port 80
    # for an "http://" scheme, and 443 for "https://".
    #
    # allowed_origins [
    #    "http://www.example.com"
    #    "https://www.other-example.com"
    # ]

    # This enables support for compressed websocket frames
    # in the server. For compression to be used, both server
    # and client have to support it.
    #
    # compression: true

    # This is the total time allowed for the server to
    # read the client request and write the response back
    # to the client. This includes the time needed for the
    # TLS handshake.
    #
    # handshake_timeout: "2s"

    # Name for an HTTP cookie, that if present will be used as a client JWT. 
    # If the client specifies a JWT in the CONNECT protocol, this option is ignored.
    # The cookie should be set by the HTTP server as described [here](https://developer.mozilla.org/en-US/docs/Web/HTTP/Cookies#restrict_access_to_cookies).
    # This setting is useful when generating NATS `Bearer` client JWTs as the
    # result of some authentication mechanism. The HTTP server after correct
    # authentication can issue a JWT for the user, that is set securely preventing
    # access by unintended scripts. Note these JWTs must be [NATS JWTs](https://docs.nats.io/nats-server/configuration/securing_nats/jwt).
    #
    # jwt_cookie: "my_jwt_cookie_name"

    # If no user name is provided when a websocket client connects, will default
    # this user name in the authentication phase. If specified, this will
    # override, for websocket clients, any `no_auth_user` value defined in the
    # main configuration file.
    # Note that this is not compatible with running the server in operator mode.
    #
    # no_auth_user: "my_username_for_apps_not_providing_credentials"

    # See below to know what is the normal way of limiting websocket clients
    # to specific users.
    # If there are no users specified in the configuration, this simple authorization
    # block allows you to override the values that would be configured in the
    # equivalent block in the main section.
    #
    # authorization {
    #     # If this is specified, the client has to provide the same username
    #     # and password to be able to connect.
    #     # username: "my_user_name"
    #     # password: "my_password"
    #
    #     # If this is specified, the password field in the CONNECT has to
    #     # match this token.
    #     # token: "my_token"
    #
    #     # This overrides the main's authorization timeout. For consistency
    #     # with the main's authorization configuration block, this is expressed
    #     # as a number of seconds.
    #     # timeout: 2.0
    #}
}
```

--------------------------------

### Configure NATS User Connection Restrictions

Source: https://github.com/nats-io/nats.docs/blob/master/running-a-nats-service/configuration/websocket/websocket_conf.md

Demonstrates how to use the allowed_connection_types field in the NATS configuration to limit which protocols a specific user can use to connect. This example restricts the user 'bar' to only WebSocket connections.

```text
authorization {
  users [
    {user: foo password: foopwd, permission: {...}}
    {user: bar password: barpwd, permission: {...}, allowed_connection_types: ["WEBSOCKET"]}
  ]
}
```

--------------------------------

### LeafNode Remote Entry Configuration

Source: https://github.com/nats-io/nats.docs/blob/master/running-a-nats-service/configuration/leafnodes/leafnode_conf.md

Defines the properties for configuring a remote LeafNode connection, including URLs, account information, denial lists for imports/exports, credentials, TLS settings, and WebSocket options.

```go
type RemoteLeafOpts struct {
	Url              string
	Urls             []string
	NoRandomize      bool
	Account          string
	DenyImports      []string
	DenyExports      []string
	Credentials      string
	Nkey             string
	Tls              *TlsConfig
	WsCompression    bool
	WsNoMasking      bool
	Compression      string
	Hub              bool
	FirstInfoTimeout time.Duration
}
```

--------------------------------

### Enable Monitoring via Server Configuration

Source: https://github.com/nats-io/nats.docs/blob/master/running-a-nats-service/nats_admin/monitoring/README.md

Configure the HTTP port for monitoring in the NATS server configuration file.

```yaml
http_port: 8222
```

--------------------------------

### Configure allow_responses for Request-Reply

Source: https://github.com/nats-io/nats.docs/blob/master/running-a-nats-service/configuration/securing_nats/authorization.md

This snippet demonstrates how to configure user permissions to allow responses to requests. It includes examples for unrestricted access, single-reply limits, and time-bound reply constraints.

```hcl
authorization: {
    users: [
        { user: a, password: a },
        { user: b, password: b, permissions: {subscribe: "q", allow_responses: true } },
        { user: c, password: c, permissions: {subscribe: "q", allow_responses: { max: 5, expires: "1m" } } }
        { user: d, password: d, permissions: {subscribe: "q", publish: "x", allow_responses: true } }
    ]
}
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

### Configure Leafnode TLS Handshake First (Accepting Side)

Source: https://github.com/nats-io/nats.docs/blob/master/running-a-nats-service/configuration/leafnodes/leafnode_conf.md

This configuration enables a TLS-first handshake for incoming leafnode connections on the accepting server. It requires that remote leafnodes also have the `handshake_first` setting enabled for a successful connection. If not, the connection will fail due to protocol mismatch.

```go
leafnodes {
  port: 7422
  tls {
    handshake_first: true
    # other TLS fields...
  }
}
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

### Configure Leafnode TLS Handshake First (Remote Side)

Source: https://github.com/nats-io/nats.docs/blob/master/running-a-nats-service/configuration/leafnodes/leafnode_conf.md

This configuration specifies that a leafnode client should initiate a TLS-first handshake when connecting to a remote server. If the remote server does not also have `handshake_first` enabled in its configuration, the connection will fail because the client will attempt a TLS handshake while the server expects the NATS INFO protocol.

```go
leafnodes {
  remotes [
    {
      urls: ["tls://example:7422"]
      tls: {
         handshake_first: true
         # other TLS fields...
      }
    }
  ]
}
```

--------------------------------

### Configure Reconnect Buffer Size in Go

Source: https://github.com/nats-io/nats.docs/blob/master/using-nats/developing-with-nats/reconnect/buffer.md

Demonstrates how to set the reconnect buffer size in bytes for the Go NATS client. This allows buffering messages during short reconnects. The buffer size is specified in bytes.

```Go
// Set reconnect buffer size in bytes (5 MB)
nc, err := nats.Connect("demo.nats.io", nats.ReconnectBufSize(5*1024*1024))
if err != nil {
    log.Fatal(err)
}
deffer nc.Close()

// Do something with the connection
```

--------------------------------

### Enable TLS-first Handshake

Source: https://github.com/nats-io/nats.docs/blob/master/running-a-nats-service/configuration/securing_nats/tls.md

Configure the server to perform a TLS handshake before sending the INFO protocol message. This ensures no clear text is sent initially. Set to 'true' for strict TLS-first behavior, or 'auto' (default 50ms) or a duration (e.g., '300ms') to allow fallback for older clients.

```text
tls: {
  cert_file: "./server-cert.pem"
  key_file: "./server-key.pem"
  handshake_first: true
}
```

```text
tls: {
  cert_file: "./server-cert.pem"
  key_file: "./server-key.pem"
  handshake_first: auto
}
```

```text
tls {
    cert_file: ...
    key_file: ...

    handshake_first: "300ms"
}
```

--------------------------------

### Configure NATS System Account

Source: https://github.com/nats-io/nats.docs/blob/master/running-a-nats-service/nats_admin/jwt.md

Shows the configuration syntax for defining a system account in the NATS server configuration file using a public NKEY.

```text
system_account: AAAXAUVSGK7TCRHFIRAS4SYXVJ76EWDMNXZM6ARFGXP7BASNDGLKU7A5
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

### Configure NATS Token Authentication in Configuration File

Source: https://github.com/nats-io/nats.docs/blob/master/running-a-nats-service/configuration/securing_nats/auth_intro/tokens.md

Defines how to set a clear-text token within the NATS server configuration file under the authorization section.

```text
authorization {
    token: "s3cr3t"
}
```

--------------------------------

### Configure Reconnect Buffer Size

Source: https://github.com/nats-io/nats.docs/blob/master/_examples/reconnect_5mb.html

How to set the reconnect buffer size for NATS clients in Go and Java.

```APIDOC
## Configure Reconnect Buffer Size

### Description
Sets the size of the buffer used for messages during a reconnection event. This ensures messages are not lost if the client is temporarily disconnected.

### Method
Client Connection Configuration

### Parameters
#### Options
- **reconnectBufferSize** (integer) - Required - The size of the buffer in bytes.

### Request Example (Go)
```go
nc, err := nats.Connect("demo.nats.io", nats.ReconnectBufSize(5*1024*1024))
```

### Request Example (Java)
```java
Options options = new Options.Builder().
    server("nats://demo.nats.io:4222").
    reconnectBufferSize(5 * 1024 * 1024).
    build();
Connection nc = Nats.connect(options);
```

### Note
This feature is currently not supported or configurable in the Node.js, Python (asyncio), Ruby, or TypeScript NATS clients.
```

--------------------------------

### Start NATS Server with Included Configuration

Source: https://github.com/nats-io/nats.docs/blob/master/running-a-nats-service/configuration/README.md

Command to start the NATS server using a main configuration file (`server.conf`) that includes other configuration files.

```shell
nats-server -c server.conf
```

--------------------------------

### Configure NATS Connection with Ping Interval

Source: https://github.com/nats-io/nats.docs/blob/master/_examples/ping_20s.html

This snippet shows how to establish a NATS connection and configure the ping interval to 20 seconds. It covers implementations in multiple languages, demonstrating common patterns for setting connection options.

```Go
// Set Ping Interval to 20 seconds
nc, err := nats.Connect("demo.nats.io", nats.Name("API Ping Example"), nats.PingInterval(20*time.Second))
if err != nil {
	log.Fatal(err)
}
defers nc.Close()

// Do something with the connection
```

```Java
Options options = new Options.Builder().                             server("nats://demo.nats.io:4222").                             pingInterval(Duration.ofSeconds(20)). // Set Ping Interval                             build();
Connection nc = Nats.connect(options);

// Do something with the connection

nc.close();
```

```JavaScript
let nc = NATS.connect({
    pingInterval: 20*1000, //20s
    url: "nats://demo.nats.io:4222"
});
```

```Python
nc = NATS()

await nc.connect(
    servers=["nats://demo.nats.io:4222"],
    # Set Ping Interval to 20 seconds
    ping_interval=20,
    )

# Do something with the connection.
```

```Ruby
require 'nats/client'

NATS.start(ping_interval: 20) do |nc|
  nc.on_reconnect do
    puts "Got reconnected to #{nc.connected_server}"
  end
  nc.on_disconnect do |reason|
    puts "Got disconnected! #{reason}"
  end
  # Do something with the connection
end
```

```TypeScript
// will throw an exception if connection fails
let nc = await connect({
    pingInterval: 20*1000, //20s
    url: "nats://demo.nats.io:4222"
});
nc.close();
```

--------------------------------

### Included NATS Configuration File

Source: https://github.com/nats-io/nats.docs/blob/master/running-a-nats-service/configuration/README.md

An example of an included configuration file (`auth.conf`) that defines authorization settings. This file is referenced by the main server configuration.

```text
authorization {
    token: "f0oBar"
}
```

--------------------------------

### NATS Configuration with Variable Reference

Source: https://github.com/nats-io/nats.docs/blob/master/running-a-nats-service/configuration/README.md

Defines a variable `TOKEN` in the configuration and then references it within the `authorization` block. This demonstrates how to use variables for dynamic configuration.

```text
# Define a variable in the config
TOKEN: "secret"

# Reference the variable
authorization {
    token: $TOKEN
}
```

--------------------------------

### Configure Reconnect Buffer Size in Java

Source: https://github.com/nats-io/nats.docs/blob/master/using-nats/developing-with-nats/reconnect/buffer.md

Shows how to configure the reconnect buffer size in bytes for the Java NATS client. This feature buffers messages when the connection is down and sends them upon reconnection.

```Java
Options options = new Options.Builder()
    .server("nats://demo.nats.io:4222")
    .reconnectBufferSize(5 * 1024 * 1024)  // Set buffer in bytes
    .build();
Connection nc = Nats.connect(options);

// Do something with the connection

nc.close();
```

--------------------------------

### Configure Reconnect Buffer Size in Go

Source: https://github.com/nats-io/nats.docs/blob/master/_examples/reconnect_5mb.html

Sets the reconnection buffer size for the Go NATS client to 5MB. This allows the client to buffer messages during temporary disconnections up to the specified size. Errors during connection are logged and the program exits.

```Go
// Set reconnect buffer size in bytes (5 MB) nc, err := nats.Connect("demo.nats.io", nats.ReconnectBufSize(5*1024*1024)) if err != nil { 
	log.Fatal(err) 
}
defer nc.Close() 

// Do something with the connection
```

--------------------------------

### Generate NATS Server Configuration and Start Server

Source: https://github.com/nats-io/nats.docs/blob/master/running-a-nats-service/nats_admin/jwt.md

Generates a NATS server configuration file using `nsc` and starts the NATS server in the background. This configuration includes operator details and system account JWTs. Ensure the `nsc` tool is installed and configured.

```shell
nsc generate config --nats-resolver > nats-res.cfg
nats-server -c nats-res.cfg --addr localhost --port 4222 &
```