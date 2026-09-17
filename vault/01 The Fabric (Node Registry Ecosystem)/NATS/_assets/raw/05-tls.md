### Basic TLS Configuration

Source: https://github.com/nats-io/nats.docs/blob/master/running-a-nats-service/configuration/securing_nats/tls.md

This is the simplest configuration for enabling TLS by specifying the certificate and key files.

```yaml
tls: {
  cert_file: "./server-cert.pem"
  key_file: "./server-key.pem"
}
```

--------------------------------

### Configure TLS Client Certificate Verification

Source: https://github.com/nats-io/nats.docs/blob/master/running-a-nats-service/configuration/securing_nats/auth_intro/tls_mutual_auth.md

Demonstrates how to enable client certificate verification in the NATS server configuration file or via command line arguments. This ensures clients provide a certificate signed by a trusted CA.

```NATS Configuration
tls {
  cert_file: "server-cert.pem"
  key_file:  "server-key.pem"
  ca_file:   "rootCA.pem"
  verify:    true
}
```

```bash
nats-server --tlsverify --tlscert=server-cert.pem --tlskey=server-key.pem --tlscacert=rootCA.pem
```

--------------------------------

### Configure NATS Server for TLS

Source: https://github.com/nats-io/nats.docs/blob/master/using-nats/developing-with-nats/connecting/security/tls.md

Command line configuration to start a NATS server with TLS enabled, requiring client certificates and a CA for verification.

```bash
nats-server --tls --tlscert=server-cert.pem --tlskey=server-key.pem --tlscacert rootCA.pem --tlsverify
```

--------------------------------

### Configure TLS for Account Resolver

Source: https://github.com/nats-io/nats.docs/blob/master/release_notes/whats_new.md

Specify TLS configuration for the account resolver. Ensure cert_file, key_file, and ca_file are correctly set.

```hcl
resolver_tls {
  cert_file: ...
  key_file: ...
  ca_file: ...
}
```

--------------------------------

### Configure NATS TLS Connection

Source: https://github.com/nats-io/nats.docs/blob/master/using-nats/developing-with-nats/connecting/security/tls.md

Examples of initializing a NATS client connection with TLS options, including CA files, client certificates, and private keys. These snippets show the necessary setup to ensure secure communication with a NATS server.

```csharp
using NATS.Net;
using NATS.Client.Core;

await using var client = new NatsClient(new NatsOpts
{
    TlsOpts = new NatsTlsOpts
    {
        CaFile = "rootCA.pem",
        KeyFile = "client-key.pem",
        CertFile = "client-cert.pem",
    }
});
```

```ruby
EM.run do
  options = {
    :servers => ['nats://localhost:4222'],
    :tls => {
      :private_key_file => 'client-key.pem',
      :cert_chain_file  => 'client-cert.pem',
      :ca_file => 'rootCA.pem'
    }
  }
  NATS.connect(options) do |nc|
    # Connection logic here
  end
end
```

```c
natsConnection      *conn      = NULL;
natsOptions         *opts      = NULL;
natsStatus          s          = NATS_OK;

s = natsOptions_Create(&opts);
if (s == NATS_OK)
    s = natsOptions_LoadCertificatesChain(opts, "client-cert.pem", "client-key.pem");
if (s == NATS_OK)
    s = natsOptions_LoadCATrustedCertificates(opts, "rootCA.pem");
if (s == NATS_OK)
    s = natsConnection_Connect(&conn, opts);

natsConnection_Destroy(conn);
natsOptions_Destroy(opts);
```

--------------------------------

### TLS Configuration via Server Options

Source: https://github.com/nats-io/nats.docs/blob/master/running-a-nats-service/configuration/securing_nats/tls.md

Alternatively, TLS can be enabled using command-line flags when starting the NATS server.

```shell
nats-server --tls --tlscert=./server-cert.pem --tlskey=./server-key.pem
```

--------------------------------

### Connect to NATS via TLS

Source: https://github.com/nats-io/nats.docs/blob/master/_examples/connect_tls_url.html

Demonstrates how to establish a secure TLS connection to a NATS server. Each implementation handles certificate validation and protocol configuration specific to its language ecosystem.

```go
nc, err := nats.Connect("tls://localhost", nats.RootCAs("resources/certs/ca.pem"))
if err != nil {
	log.Fatal(err)
}
defer nc.Close()
```

```java
SSLContext ctx = SSLContext.getInstance(Options.DEFAULT_SSL_PROTOCOL);
ctx.init(createTestKeyManagers(), createTestTrustManagers(), new SecureRandom());
SSLContext.setDefault(ctx);
Options options = new Options.Builder().server("tls://localhost:4222").build();
Connection nc = Nats.connect(options);
```

```javascript
let nc = NATS.connect({
    url: "tls://demo.nats.io:4443",
    tls: true
});
```

```python
nc = NATS()
await nc.connect("tls://demo.nats.io:4443", loop=loop)
```

```ruby
options = {
  :servers => ['nats://demo.nats.io:4443'],
  :tls => {}
}
NATS.connect(options) do |nc|
  puts :connected
end
```

```typescript
let nc = await connect({
    url: "tls://demo.nats.io:4443"
});
nc.close();
```

--------------------------------

### Connect to NATS with TLS using Node.js

Source: https://github.com/nats-io/nats.docs/blob/master/_examples/connect_tls.html

Establishes a TLS connection to a NATS server using Node.js. It requires the paths to CA certificate, client certificate, and client key files. The connection options specify the TLS configuration.

```JavaScript
let caCert = fs.readFileSync(caCertPath);
let clientCert = fs.readFileSync(clientCertPath);
let clientKey = fs.readFileSync(clientKeyPath);
let nc = NATS.connect({
    url: url,
    tls: {
        ca: [caCert],
        key: [clientKey],
        cert: [clientCert]
    }
});
```

--------------------------------

### Map TLS Certificates to User Identities

Source: https://github.com/nats-io/nats.docs/blob/master/running-a-nats-service/configuration/securing_nats/auth_intro/tls_mutual_auth.md

Shows how to use the verify_and_map option to authenticate users based on certificate attributes like SANs or Distinguished Names. This allows for passwordless authentication by linking certificate data to defined NATS users.

```NATS Configuration
tls {
  cert_file: "server-cert.pem"
  key_file:  "server-key.pem"
  ca_file:   "rootCA.pem"
  verify_and_map: true
}
```

```bash
openssl x509 -noout -text -in client-cert.pem
```

```NATS Configuration
authorization {
  users = [
    {user: "email@localhost"}
  ]
}
```

```NATS Configuration
authorization {
  users = [
    {user: "OU=testuser@MacBook-Pro.local (Test User),O=mkcert development certificate"}
  ]
}
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

### Connect to NATS with TLS using Go

Source: https://github.com/nats-io/nats.docs/blob/master/_examples/connect_tls.html

Establishes a TLS connection to a NATS server using client certificates and a root CA. Requires specifying the server address and paths to certificate files.

```Go
nc, err := nats.Connect("localhost", 	nats.ClientCert("resources/certs/cert.pem", "resources/certs/key.pem"), 	nats.RootCAs("resources/certs/ca.pem")) if err != nil { 	log.Fatal(err) } defer nc.Close()  // Do something with the connection
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

### Connect to NATS with TLS using Ruby

Source: https://github.com/nats-io/nats.docs/blob/master/_examples/connect_tls.html

Demonstrates establishing a TLS connection to a NATS server using the Ruby NATS client. It configures the connection with server details and TLS options, including paths to the private key and certificate chain files. The example also includes setting up various event callbacks.

```Ruby
EM.run do
   options = {
    :servers => [
      'nats://localhost:4222',
    ],
    :tls => {
      :private_key_file => './spec/configs/certs/key.pem',
      :cert_chain_file  => './spec/configs/certs/server.pem'
    }
  }

   NATS.connect(options) do |nc|
     puts "#{Time.now.to_f} - Connected to NATS at #{nc.connected_server}"

     nc.subscribe("hello") do |msg|
       puts "#{Time.now.to_f} - Received: #{msg}"
     end

     nc.flush do
       nc.publish("hello", "world")
     end

     EM.add_periodic_timer(0.1) do
       next unless nc.connected?
       nc.publish("hello", "hello")
     end

     # Set default callbacks
     nc.on_error do |e|
       puts "#{Time.now.to_f } - Error: #{e}"
     end

     nc.on_disconnect do |reason|
       puts "#{Time.now.to_f} - Disconnected: #{reason}"
     end

     nc.on_reconnect do |nc|
       puts "#{Time.now.to_f} - Reconnected to NATS server at #{nc.connected_server}"
     end

     nc.on_close do
       puts "#{Time.now.to_f} - Connection to NATS closed"
       EM.stop
     end
  end
end
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

### Configure TLS Timeout

Source: https://github.com/nats-io/nats.docs/blob/master/running-a-nats-service/configuration/securing_nats/tls.md

Control the time allowed for a client to upgrade its connection to TLS. A low value may cause handshake errors, while a high value exposes the server to resource exhaustion attacks. The example shows a value too low, which will cause clients to fail connection.

```text
tls: {
  cert_file: "./server-cert.pem"
  key_file: "./server-key.pem"
  # clients will fail to connect (value is too low)
  timeout: 0.0001
}
```

--------------------------------

### Configure NATS Authentication Timeout

Source: https://github.com/nats-io/nats.docs/blob/master/running-a-nats-service/configuration/securing_nats/auth_intro/auth_timeout.md

This configuration snippet shows how to set the authentication timeout within the NATS server configuration. The timeout value is specified in seconds and should be longer than the TLS timeout. If not specified, it defaults to 1 second more than the tls_timeout.

```json
{
    "authorization": {
        "timeout": 3,
        "users": [
            {"user": "a", "password": "b"},
            {"user": "b", "password": "a"}
        ]
    }
}
```

--------------------------------

### Connect to NATS with TLS using TypeScript

Source: https://github.com/nats-io/nats.docs/blob/master/_examples/connect_tls.html

This TypeScript code snippet shows how to connect to a NATS server using TLS. It reads CA certificate, client certificate, and client key files from the filesystem and configures the `connect` function with these TLS options. Ensure the `ts-nats` library is installed.

```typescript
let caCertPath = "path/to/ca.crt";
let clientCertPath = "path/to/client.crt";
let clientKeyPath = "path/to/client.key";
let url = "nats://localhost:4222";

let caCert = readFileSync(caCertPath);
let clientCert = readFileSync(clientCertPath);
let clientKey = readFileSync(clientKeyPath);

let nc = await connect({
    url: url,
    tls: {
        ca: [caCert],
        key: [clientKey],
        cert: [clientCert]
    }
});
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

### Configure NATS for TLS-Terminating Reverse Proxy

Source: https://github.com/nats-io/nats.docs/blob/master/running-a-nats-service/configuration/securing_nats/tls.md

Configure the NATS server to be TLS-available and allow non-TLS connections. This setup is necessary when using a TLS-terminating reverse proxy to handle client connections.

```yaml
tls {}
allow_non_tls: true
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

### Connect to NATS with TLS

Source: https://github.com/nats-io/nats.docs/blob/master/using-nats/developing-with-nats/connecting/security/tls.md

Demonstrates how to establish a secure NATS connection using client certificates and CA trust stores across different client libraries.

```go
nc, err := nats.Connect("localhost",
    nats.ClientCert("client-cert.pem", "client-key.pem"),
    nats.RootCAs("rootCA.pem"))
if err != nil {
    log.Fatal(err)
}
defer nc.Close()
```

```java
SSLContext ctx = SSLUtils.createSSLContext();
Options options = new Options.Builder()
    .server("nats://localhost:4222")
    .sslContext(ctx)
    .build();
Connection nc = Nats.connect(options);
```

```javascript
const nc = await connect({
  port: ns.port,
  debug: true,
  tls: {
    caFile: caCertPath,
    keyFile: clientKeyPath,
    certFile: clientCertPath,
  },
});
```

```python
nc = NATS()
ssl_ctx = ssl.create_default_context(purpose=ssl.Purpose.SERVER_AUTH)
ssl_ctx.load_verify_locations('rootCA.pem')
ssl_ctx.load_cert_chain(certfile='client-cert.pem', keyfile='client-key.pem')
await nc.connect(servers=["nats://demo.nats.io:4222"], tls=ssl_ctx)
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