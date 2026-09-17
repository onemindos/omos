### Configure MQTT in NATS Server

Source: https://github.com/nats-io/nats.docs/blob/master/running-a-nats-service/configuration/mqtt/mqtt_conf.md

This configuration block enables MQTT support in the NATS server. It includes settings for the listening port, TLS security, authentication overrides, and QoS acknowledgment parameters.

```hcl
mqtt {
    port: 1883
    tls {
        cert_file: "/path/to/cert.pem"
        key_file: "/path/to/key.pem"
    }
    ack_wait: "1m"
    max_ack_pending: 100
}
```

--------------------------------

### NATS Server Configuration with MQTT User Permissions

Source: https://github.com/nats-io/nats.docs/blob/master/running-a-nats-service/configuration/mqtt/mqtt_conf.md

This example shows a complete NATS server configuration that includes MQTT user authorization. It defines specific publish and subscribe permissions for an MQTT user, including the necessary permission for JetStream durable subscriptions (`$MQTT.sub.>"). The MQTT listener is also configured.

```HCL
    listen: 127.0.0.1:4222
    jetstream: enabled
    authorization {
        mqtt_perms = {
            publish = ["baz"]
            subscribe = ["foo", "bar", "$MQTT.sub.>"]
        }
        users = [
            {user: mqtt, password: pass, permissions: $mqtt_perms, allowed_connection_types: ["MQTT"]}
        ]
    }
    mqtt {
        listen: 127.0.0.1:1883
    }
```

--------------------------------

### Configure MQTT User Permissions with JetStream Support

Source: https://github.com/nats-io/nats.docs/blob/master/running-a-nats-service/configuration/mqtt/mqtt_config.md

Shows a complete NATS server configuration for an MQTT user. It includes the required subscription permission for '$MQTT.sub.>' to support JetStream-backed QoS 1 subscriptions.

```text
listen: 127.0.0.1:4222
jetstream: enabled
authorization {
    mqtt_perms = {
        publish = ["baz"]
        subscribe = ["foo", "bar", "$MQTT.sub.>"]
    }
    users = [
        {user: mqtt, password: pass, permissions: $mqtt_perms, allowed_connection_types: ["MQTT"]}
    ]
}
mqtt {
    listen: 127.0.0.1:1883
}
```

--------------------------------

### Configure MQTT User Authorization with Connection Types

Source: https://github.com/nats-io/nats.docs/blob/master/running-a-nats-service/configuration/mqtt/mqtt_conf.md

This configuration snippet demonstrates how to set up user authorization for MQTT connections in NATS. It shows how to restrict specific users to only accept MQTT connections using the 'allowed_connection_types' field. The absence of this field allows all connection types by default.

```HCL
authorization {
  users [
    {user: foo password: foopwd, permission: {...}}
    {user: bar password: barpwd, permission: {...}, allowed_connection_types: ["MQTT"]}
  ]
}
```

```HCL
authorization {
  users [
    {user: foo password: foopwd, permission: {...}}
    {user: bar password: barpwd, permission: {...}, allowed_connection_types: ["STANDARD", "MQTT"]}
  ]
}
```

--------------------------------

### Configure Allowed Connection Types for NATS Users

Source: https://github.com/nats-io/nats.docs/blob/master/running-a-nats-service/configuration/mqtt/mqtt_config.md

Demonstrates how to restrict a user to specific connection types like MQTT or a combination of STANDARD and MQTT. When the allowed_connection_types field is omitted, all connection types are permitted by default.

```text
authorization {
  users [
    {user: foo password: foopwd, permission: {...}}
    {user: bar password: barpwd, permission: {...}, allowed_connection_types: ["STANDARD", "MQTT"]}
  ]
}
```

--------------------------------

### Set NATS Server Name for MQTT

Source: https://github.com/nats-io/nats.docs/blob/master/running-a-nats-service/configuration/mqtt/mqtt_config.md

Shows how to define a unique server name at the top level of the NATS configuration, which is a requirement for MQTT functionality.

```text
server_name: "my_mqtt_server"
mqtt {
    port: 1883
}
```

--------------------------------

### Configure User for Operator Mode

Source: https://github.com/nats-io/nats.docs/blob/master/running-a-nats-service/configuration/mqtt/mqtt_config.md

Command to set the bearer flag on a user using nsc, which is necessary for MQTT clients to authenticate in operator mode.

```shell
nsc edit user --name U --account A --bearer
```

--------------------------------

### Enable Monitoring via Server Configuration

Source: https://github.com/nats-io/nats.docs/blob/master/running-a-nats-service/nats_admin/monitoring/README.md

Configure the HTTP port for monitoring in the NATS server configuration file.

```yaml
http_port: 8222
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

### Configure NATS System Account

Source: https://github.com/nats-io/nats.docs/blob/master/running-a-nats-service/nats_admin/jwt.md

Shows the configuration syntax for defining a system account in the NATS server configuration file using a public NKEY.

```text
system_account: AAAXAUVSGK7TCRHFIRAS4SYXVJ76EWDMNXZM6ARFGXP7BASNDGLKU7A5
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

### Start NATS Server with Included Configuration

Source: https://github.com/nats-io/nats.docs/blob/master/running-a-nats-service/configuration/README.md

Command to start the NATS server using a main configuration file (`server.conf`) that includes other configuration files.

```shell
nats-server -c server.conf
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

### Included NATS Configuration File

Source: https://github.com/nats-io/nats.docs/blob/master/running-a-nats-service/configuration/README.md

An example of an included configuration file (`auth.conf`) that defines authorization settings. This file is referenced by the main server configuration.

```text
authorization {
    token: "f0oBar"
}
```

--------------------------------

### Configure Reconnect Buffer Size in Java

Source: https://github.com/nats-io/nats.docs/blob/master/_examples/reconnect_5mb.html

Configures the reconnection buffer size for the Java NATS client to 5MB using Options. This buffer helps manage messages during network interruptions. The connection is established with the specified options and closed afterwards.

```Java
Options options = new Options.Builder().                             server("nats://demo.nats.io:4222").                             reconnectBufferSize(5 * 1024 * 1024).  // Set buffer in bytes                             build(); 
Connection nc = Nats.connect(options);

// Do something with the connection

nc.close();
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

### Generate NATS Server Configuration and Start Server

Source: https://github.com/nats-io/nats.docs/blob/master/running-a-nats-service/nats_admin/jwt.md

Generates a NATS server configuration file using `nsc` and starts the NATS server in the background. This configuration includes operator details and system account JWTs. Ensure the `nsc` tool is installed and configured.

```shell
nsc generate config --nats-resolver > nats-res.cfg
nats-server -c nats-res.cfg --addr localhost --port 4222 &
```