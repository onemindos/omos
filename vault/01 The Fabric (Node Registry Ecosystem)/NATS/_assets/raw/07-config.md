### Start NATS Server with Configuration File

Source: https://github.com/nats-io/nats.docs/blob/master/running-a-nats-service/configuration/README.md

Specify a custom configuration file for the NATS server using the `-config` flag.

```shell
nats-server -config my-server.conf
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

### Define NATS Server Configuration

Source: https://github.com/nats-io/nats.docs/blob/master/using-nats/nats-tools/nsc/basics.md

Example configuration file for a local NATS server instance.

```text
server_name: servertest
listen: 127.0.0.1:4222
http: 8222

jetstream: enabled

include resolver.conf
```

--------------------------------

### Enable Monitoring via Server Configuration

Source: https://github.com/nats-io/nats.docs/blob/master/running-a-nats-service/nats_admin/monitoring/README.md

Configure the HTTP port for monitoring in the NATS server configuration file.

```yaml
http_port: 8222
```

--------------------------------

### NATS Server Configuration with Include Directive

Source: https://github.com/nats-io/nats.docs/blob/master/running-a-nats-service/configuration/README.md

Splits the NATS server configuration into multiple files using the `include` directive. The included file `auth.conf` is relative to the main configuration file `server.conf`.

```text
listen: 127.0.0.1:4222
include ./auth.conf
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

### NATS Server Leafnode Configuration

Source: https://github.com/nats-io/nats.docs/blob/master/running-a-nats-service/running/nats_docker/ngs-leafnodes-docker.md

A minimal NATS server configuration file for leaf nodes connecting to Synadia Cloud (NGS). It specifies the remote server URL and the path to the credentials file.

```NATS Config
leafnodes {
    remotes = [
        {
          url: "tls://connect.ngs.global"
          credentials: "ngs.creds"
        },
    ]
}
```

--------------------------------

### Generate Server Configuration for Memory Resolver

Source: https://github.com/nats-io/nats.docs/blob/master/running-a-nats-service/configuration/securing_nats/jwt/mem_resolver.md

This command generates a NATS server configuration file that utilizes the MEMORY resolver. It automatically configures the server to load account JWTs from memory, simplifying setup for smaller deployments. The output is saved to the specified file path.

```shell
nsc generate config --mem-resolver --config-file /tmp/server.conf
```

--------------------------------

### Start NATS Server with Included Configuration

Source: https://github.com/nats-io/nats.docs/blob/master/running-a-nats-service/configuration/README.md

Command to start the NATS server using a main configuration file (`server.conf`) that includes other configuration files.

```shell
nats-server -c server.conf
```

--------------------------------

### NATS Configuration with Environment Variable

Source: https://github.com/nats-io/nats.docs/blob/master/running-a-nats-service/configuration/README.md

Demonstrates how a NATS server configuration can reference a variable that is defined in the environment. Ensure the environment variable is exported before starting the server.

```text
# TOKEN is defined in the environment
authorization {
    token: $TOKEN
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

### Generate NATS Server Configuration and Start Server

Source: https://github.com/nats-io/nats.docs/blob/master/running-a-nats-service/nats_admin/jwt.md

Generates a NATS server configuration file using `nsc` and starts the NATS server in the background. This configuration includes operator details and system account JWTs. Ensure the `nsc` tool is installed and configured.

```shell
nsc generate config --nats-resolver > nats-res.cfg
nats-server -c nats-res.cfg --addr localhost --port 4222 &
```

--------------------------------

### Configure and Start NATS Server

Source: https://github.com/nats-io/nats.docs/blob/master/running-a-nats-service/configuration/sys_accounts/sys_accounts.md

Defines the server configuration file to link the operator JWT, system account public key, and the account resolver URL, then starts the server.

```text
operator: /Users/synadia/.nsc/nats/SAOP/SAOP.jwt
system_account: ADWJVSUSEVC2GHL5GRATN2LOEOQOY2E6Z2VXNU3JEIK6BDGPWNIW3AXF
resolver: URL(http://localhost:9090/jwt/v1/accounts/)
```

```shell
nats-server -c server.conf
```

--------------------------------

### Server Configuration: Simple Mappings

Source: https://github.com/nats-io/nats.docs/blob/master/nats-concepts/subject_mapping.md

Defines subject mappings directly in the NATS server configuration file. This applies mappings to the default account.

```yaml
server_name: "hub"
cluster: { name: "hub" }
mappings: {
    orders.flush  orders.central.flush 
	orders.* orders.central.{{wildcard(1)}}
}
```

--------------------------------

### Start NATS Server with Configuration

Source: https://github.com/nats-io/nats.docs/blob/master/running-a-nats-service/configuration/securing_nats/jwt/mem_resolver.md

This shell command demonstrates how to start the NATS server using a specified configuration file. It's a common command for initializing a NATS server instance with custom settings.

```shell
nats-server -c server.conf
```

--------------------------------

### NATS Server Configuration File Example

Source: https://github.com/nats-io/nats.docs/blob/master/running-a-nats-service/configuration/securing_nats/jwt/mem_resolver.md

This snippet shows the structure of a NATS server configuration file using memory resolver and preloaded account JWTs. It specifies the operator JWT path and the resolver settings, including a preloaded account JWT with its public key and contents.

```yaml
operator: /Users/synadia/.nsc/nats/memory/memory.jwt
resolver: MEMORY
resolver_preload: {
ACSU3Q6LTLBVLGAQUONAGXJHVNWGSKKAUA7IY5TB4Z7PLEKSR5O6JTGR: eyJ0eXAiOiJqd3QiLCJhbGciOiJlZDI1NTE5In0.eyJqdGkiOiJPRFhJSVI2Wlg1Q1AzMlFJTFczWFBENEtTSDYzUFNNSEZHUkpaT05DR1RLVVBISlRLQ0JBIiwiaWF0IjoxNTU2NjU1Njk0LCJpc3MiOiJPRFdaSjJLQVBGNzZXT1dNUENKRjZCWTRRSVBMVFVJWTRKSUJMVTRLM1lERzNHSElXQlZXQkhVWiIsIm5hbWUiOiJBIiwic3ViIjoiQUNTVTNRNkxUTEJWTEdBUVVPTkFHWEpIVk5XR1NLS0FVQTdJWTVUQjRaN1BMRUtTUjVPNkpUR1IiLCJ0eXBlIjoiYWNjb3VudCIsIm5hdHMiOnsibGltaXRzIjp7InN1YnMiOi0xLCJjb25uIjotMSwibGVhZiI6LTEsImltcG9ydHMiOi0xLCJleHBvcnRzIjotMSwiZGF0YSI6LTEsInBheWxvYWQiOi0xLCJ3aWxkY2FyZHMiOnRydWV9fX0._WW5C1triCh8a4jhyBxEZZP8RJ17pINS8qLzz-01o6zbz1uZfTOJGvwSTS6Yv2_849B9iUXSd-8kp1iMXHdoBA
}
```

--------------------------------

### Manage NATS Server Lifecycle

Source: https://github.com/nats-io/nats.docs/blob/master/running-a-nats-service/nats_admin/jwt.md

Commands to start the NATS server with a configuration file and reload the configuration during runtime.

```shell
nats-server -c server.conf
nats-server --signal reload
```

--------------------------------

### Enable NATS Monitoring via Configuration File

Source: https://github.com/nats-io/nats.docs/blob/master/running-a-nats-service/configuration/monitoring.md

Configures the NATS server to enable its HTTP monitoring endpoint using a YAML configuration file. The `http_port` directive specifies the port for monitoring. Binding to `localhost` enhances security by limiting access.

```yaml
http_port: 8222
```

```yaml
http: localhost:8222
```

--------------------------------

### NATS Server Configuration with NKeys Static Auth

Source: https://github.com/nats-io/nats.docs/blob/master/running-a-nats-service/configuration/securing_nats/jwt/jwt_nkey_auth.md

An example NATS server configuration file demonstrating the use of NKeys for static authentication. It defines accounts, users, and their associated NKeys, suitable for a centralized authentication model.

```hcl
port = 4222

cluster {
  port = 6222

  # We will bridge two different servers with different auth models via routes
  # routes [ nats://127.0.0.1:6223 ]
}

system_account = ABKOWIYVTHNEK5HELPWLAT2CF2CUPELIK4SZH2VCJHLFU22B5U2IIZUO

accounts {
  # Account A
  ADFB2JXYTXOJEL6LNAXDREUGRX35BOLZI3B4PFFAC7IRPR3OA4QNKBN2 {
    nkey: ADFB2JXYTXOJEL6LNAXDREUGRX35BOLZI3B4PFFAC7IRPR3OA4QNKBN2
    users = [
      {nkey: "UAPOK2P7EN3UFBL7SBJPQK3M3JMLALYRYKX5XWSVMVYK63ZMBHTOHVJR" }
    ]
  }

  # Account B
  ACWOMQA7PZTKJSBTR7BF6TBK3D776734PWHWDKO7HFMQOM5BIOYPSYZZ {
  }

  # Account SYS
  ABKOWIYVTHNEK5HELPWLAT2CF2CUPELIK4SZH2VCJHLFU22B5U2IIZUO {
  }
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

### Configure NATS with Bcrypted Token

Source: https://github.com/nats-io/nats.docs/blob/master/running-a-nats-service/configuration/securing_nats/auth_intro/tokens.md

Example of a server configuration file using a bcrypt-hashed token instead of a clear-text string.

```text
authorization {
    token: "$2a$11$PWIFAL8RsWyGI3jVZtO9Nu8.6jOxzxfZo7c/W0eLk017hjgUKWrhy"
}
```

--------------------------------

### Start NATS Server with Gateway Configuration

Source: https://github.com/nats-io/nats.docs/blob/master/running-a-nats-service/configuration/gateways/gateway.md

Command to start a NATS server using a specified configuration file that includes gateway settings.

```shell
nats-server -c A.conf
```