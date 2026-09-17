### Programmatic User Authentication and JWT Creation

Source: https://github.com/nats-io/nats.docs/blob/master/running-a-nats-service/nats_admin/jwt.md

Demonstrates how to generate a user JWT and NKEY pair programmatically, and how to use these artifacts to configure NATS client authentication. This includes creating a credentials file and connecting to a NATS server.

```go
func GetAccountSigningKey() nkeys.KeyPair {
    accSeed := []byte("SAAJGCAHPHHM6AVJJWQ2YAS3I4NETXMWVQSTCQMJ7VVTGAJF5UCN3IX7J4")
    accountSigningKey, err := nkeys.ParseDecoratedNKey(accSeed)
    if err != nil {
        panic(err)
    }
    return accountSigningKey
}

func RequestUser() {
    accountPublicKey := GetAccountPublicKey()
    accountSigningKey := GetAccountSigningKey()
    userPublicKey, userSeed, userKeyPair := generateUserKey()
    userJWT := generateUserJWT(userPublicKey, accountPublicKey, accountSigningKey)
    var jwtAuthOption nats.Option
    jwtAuthOption = nats.UserJWT(func() (string, error) {
            return userJWT, nil
        },
        func(bytes []byte) ([]byte, error) {
            return userKeyPair.Sign(bytes)
        },
    )
    credsContent, err := jwt.FormatUserConfig(userJWT, userSeed);
    if err != nil {
        panic(err)
    }
    ioutil.WriteFile("my.creds", credsContent, 0644)
    jwtAuthOption = nats.UserCredentials("my.creds")
    nc, err := nats.Connect("nats://localhost:4222", jwtAuthOption)
}
```

--------------------------------

### Connect to NATS using NKey Authentication

Source: https://github.com/nats-io/nats.docs/blob/master/_examples/connect_nkey.html

Demonstrates how to configure a NATS connection using NKey authentication. This involves providing a seed file or key pair to sign authentication nonces during the handshake process.

```go
opt, err := nats.NkeyOptionFromSeed("seed.txt")
if err != nil {
	log.Fatal(err)
}
nc, err := nats.Connect("127.0.0.1", opt)
if err != nil {
	log.Fatal(err)
}
defer nc.Close()
```

```java
NKey theNKey = NKey.createUser(null);
Options options = new Options.Builder().
            server("nats://localhost:4222").
            authHandler(new AuthHandler(){
                public char[] getID() {
                    try { return theNKey.getPublicKey(); } catch (Exception ex) { return null; }
                }
                public byte[] sign(byte[] nonce) {
                    try { return theNKey.sign(nonce); } catch (Exception ex) { return null; }
                }
                public char[] getJWT() { return null; }
            }).
            build();
Connection nc = Nats.connect(options);
nc.close();
```

```javascript
const seed = 'SUAEL6GG2L2HIF7DUGZJGMRUFKXELGGYFMHF76UO2AYBG3K4YLWR3FKC2Q';
let nc = NATS.connect({
    url: server.nats,
    nkey: 'UD6OU4D3CIOGIDZVL4ANXU3NWXOW5DCDE2YPZDBHPBXCVKHSODUA4FKI',
    sigCB: function (nonce) {
        const sk = nkeys.fromSeed(Buffer.from(seed));
        return sk.sign(nonce);
    }
});
```

```python
nc = NATS()
async def error_cb(e):
    print("Error:", e)
await nc.connect("nats://localhost:4222",
                 nkeys_seed="./path/to/nkeys/user.nk",
                 error_cb=error_cb)
await nc.close()
```

```typescript
const seed = 'SUAEL6GG2L2HIF7DUGZJGMRUFKXELGGYFMHF76UO2AYBG3K4YLWR3FKC2Q';
let nc = await connect({
    url: server.nats,
    nkey: 'UD6OU4D3CIOGIDZVL4ANXU3NWXOW5DCDE2YPZDBHPBXCVKHSODUA4FKI',
    nonceSigner: function (nonce) {
        const sk = fromSeed(Buffer.from(seed));
        return sk.sign(Buffer.from(nonce));
    }
});
```

--------------------------------

### JavaScript NKey Authentication

Source: https://github.com/nats-io/nats.docs/blob/master/using-nats/developing-with-nats/connecting/security/nkey.md

Authenticate with NATS using NKeys in JavaScript. This example uses the nkeyAuthenticator helper with a provided seed.

```JavaScript
// seed should be stored and treated like a secret
const seed = new TextEncoder().encode(
  "SUAEL6GG2L2HIF7DUGZJGMRUFKXELGGYFMHF76UO2AYBG3K4YLWR3FKC2Q",
);
const nc = await connect({
  port: ns.port,
  authenticator: nkeyAuthenticator(seed),
});
```

--------------------------------

### Client Authentication using Credentials File

Source: https://github.com/nats-io/nats.docs/blob/master/using-nats/developing-with-nats/connecting/security/creds.md

Demonstrates how to connect to a NATS server using a credentials file in various programming languages. The credentials file simplifies authentication by containing both the user JWT and the NKey private key.

```APIDOC
## NATS Client Authentication with Credentials File

This section provides examples of how to authenticate NATS clients using a credentials file across different programming languages. A credentials file consolidates the user JWT and the NKey private key, simplifying the connection process.

### Go

```go
nc, err := nats.Connect("127.0.0.1", nats.UserCredentials("path_to_creds_file"))
if err != nil {
    log.Fatal(err)
}
defer nc.Close()

// Do something with the connection
```

### Java

```java
Options options = new Options.Builder()
    .server("nats://localhost:4222")
    .authHandler(Nats.credentials("path_to_creds_file"))
    .build();
Connection nc = Nats.connect(options);

// Do something with the connection

nc.close();
```

### JavaScript

```javascript
// credentials file contains the JWT and the secret signing key
  const authenticator = credsAuthenticator(creds);
  const nc = await connect({
    port: ns.port,
    authenticator: authenticator,
});
```

### Python

```python
nc = NATS()

async def error_cb(e):
    print("Error:", e)

await nc.connect("nats://localhost:4222",
                 user_credentials="path_to_creds_file",
                 error_cb=error_cb,
                 )

# Do something with the connection

await nc.close()
```

### C#

```csharp
// dotnet add package NATS.Net
using NATS.Net;

await using var client = new NatsClient("127.0.0.1", credsFile: "/path/to/file.creds");
```

### C

```c
natsConnection      *conn      = NULL;
natsOptions         *opts      = NULL;
natsStatus          s          = NATS_OK;

s = natsOptions_Create(&opts);
if (s == NATS_OK)
    // Pass the credential file this way if the file contains both user JWT and seed.
    // Otherwise, if the content is split, the first file is the user JWT, the second
    // contains the seed.
    s = natsOptions_SetUserCredentialsFromFiles(opts, "path_to_creds_file", NULL);
if (s == NATS_OK)
    s = natsConnection_Connect(&conn, opts);

(...)

// Destroy objects that were created
natsConnection_Destroy(conn);
natsOptions_Destroy(opts);
```
```

--------------------------------

### C# NKey Authentication

Source: https://github.com/nats-io/nats.docs/blob/master/using-nats/developing-with-nats/connecting/security/nkey.md

Connect to NATS using NKey authentication in C#. This example configures NKey authentication via NatsAuthOpts, specifying the NKey file.

```C#
// dotnet add package NATS.Net
using NATS.Net;
using NATS.Client.Core;

await using var client = new NatsClient(new NatsOpts
{
    Url = "127.0.0.1",
    Name = "API NKey Example",
    AuthOpts = new NatsAuthOpts
    {
        NKeyFile = "/path/to/nkeys/user.nk"
    }
});
```

--------------------------------

### Configure Node.js Client for NKey Authentication

Source: https://github.com/nats-io/nats.docs/blob/master/running-a-nats-service/configuration/securing_nats/auth_intro/nkey_auth.md

Sets up a Node.js NATS client to authenticate using an NKey. It requires the public key for connection and a signature callback function that uses the private key (seed) to sign connection challenges.

```javascript
const NATS = require('nats');
const nkeys = require('ts-nkeys');

const nkey_seed = 'SUACSSL3UAHUDXKFSNVUZRF5UHPMWZ6BFDTJ7M6USDXIEDNPPQYYYCU3VY';
const nc = NATS.connect({
  port: PORT,
  nkey: 'UDXU4RCSJNZOIQHZNWXHXORDPRTGNJAHAHFRGZNEEJCPQTT2M7NLCNF4',
  sigCB: function (nonce) {
    // client loads seed safely from a file
    // or some constant like `nkey_seed` defined in
    // the program
    const sk = nkeys.fromSeed(Buffer.from(nkey_seed));
    return sk.sign(nonce);
   }
});
...
```

--------------------------------

### Connect to NATS using Credentials File

Source: https://github.com/nats-io/nats.docs/blob/master/using-nats/developing-with-nats/connecting/security/creds.md

Demonstrates how to establish a connection to a NATS server using a credentials file. This file simplifies authentication by bundling the user JWT and the NKey private key. Ensure the path to the credentials file is correct.

```Go
nc, err := nats.Connect("127.0.0.1", nats.UserCredentials("path_to_creds_file"))
if err != nil {
    log.Fatal(err)
}
defer nc.Close()

// Do something with the connection
```

```Java
Options options = new Options.Builder()
    .server("nats://localhost:4222")
    .authHandler(Nats.credentials("path_to_creds_file"))
    .build();
Connection nc = Nats.connect(options);

// Do something with the connection

nc.close();
```

```JavaScript
// credentials file contains the JWT and the secret signing key
  const authenticator = credsAuthenticator(creds);
  const nc = await connect({
    port: ns.port,
    authenticator: authenticator,
});
```

```Python
nc = NATS()

async def error_cb(e):
    print("Error:", e)

await nc.connect("nats://localhost:4222",
                 user_credentials="path_to_creds_file",
                 error_cb=error_cb,
                 )

# Do something with the connection

await nc.close()
```

```C#
// dotnet add package NATS.Net
using NATS.Net;

await using var client = new NatsClient("127.0.0.1", credsFile: "/path/to/file.creds");
```

```C
natsConnection      *conn      = NULL;
natsOptions         *opts      = NULL;
natsStatus          s          = NATS_OK;

s = natsOptions_Create(&opts);
if (s == NATS_OK)
    // Pass the credential file this way if the file contains both user JWT and seed.
    // Otherwise, if the content is split, the first file is the user JWT, the second
    // contains the seed.
    s = natsOptions_SetUserCredentialsFromFiles(opts, "path_to_creds_file", NULL);
if (s == NATS_OK)
    s = natsConnection_Connect(&conn, opts);

(...)

// Destroy objects that were created
natsConnection_Destroy(conn);
natsOptions_Destroy(opts);
```

--------------------------------

### Go NKey Authentication

Source: https://github.com/nats-io/nats.docs/blob/master/using-nats/developing-with-nats/connecting/security/nkey.md

Connect to NATS using NKey authentication in Go. This example loads the NKey seed from a file and establishes a connection.

```Go
opt, err := nats.NkeyOptionFromSeed("seed.txt")
if err != nil {
    log.Fatal(err)
}
nc, err := nats.Connect("127.0.0.1", opt)
if err != nil {
    log.Fatal(err)
}
deferr nc.Close()

// Do something with the connection
```

--------------------------------

### Configure NATS Server for NKey Authentication

Source: https://github.com/nats-io/nats.docs/blob/master/running-a-nats-service/configuration/securing_nats/auth_intro/nkey_auth.md

Configures the NATS server to authorize a specific user based on their public NKey. This replaces the need for traditional username/password or token authentication for that user.

```yaml
authorization: {
  users: [
    { nkey: UDXU4RCSJNZOIQHZNWXHXORDPRTGNJAHAHFRGZNEEJCPQTT2M7NLCNF4 }
  ]
}
```

--------------------------------

### Connect to NATS with NKEY Authentication

Source: https://github.com/nats-io/nats.docs/blob/master/running-a-nats-service/configuration/securing_nats/jwt/jwt_nkey_auth.md

Demonstrates how to connect to a NATS server using an NKEY seed file for authentication. The code extracts the public key from the seed and applies it as a connection option.

```go
package main

import (
    "log"

    "github.com/nats-io/nats.go"
)

func main() {
    opts := make([]nats.Option, 0)

    // Extract public nkey from seed
    nkey, err := nats.NkeyOptionFromSeed("path/to/seed.nkey")
    if err != nil {
        log.Fatal(err)
    }
    opts = append(opts, nkey)
    nc, err := nats.Connect("127.0.0.1:4222", opts...)
    if err != nil {
        log.Fatal(err)
    }
}
```

--------------------------------

### Create Account JWT in Golang

Source: https://github.com/nats-io/nats.docs/blob/master/running-a-nats-service/nats_admin/security.md

This Go program demonstrates how to create an account JWT using the nkeys and jwt libraries. It takes an operator seed and an account name as input and outputs the signed account JWT. Dependencies include 'github.com/nats-io/jwt/v2' and 'github.com/nats-io/nkeys'.

```Go
package main

import (
	"flag"
	"fmt"
	"log"

	"github.com/nats-io/jwt/v2"
	"github.com/nats-io/nkeys"
)

func main() {
	log.SetFlags(0)

	var (
		accountSeed  string
		operatorSeed string
		name         string
	)

	flag.StringVar(&operatorSeed, "operator", "", "Operator seed for creating an account.")
	flag.StringVar(&accountSeed, "account", "", "Account seed for creating a user.")
	flag.StringVar(&name, "name", "", "Account or user name to be created.")

	flag.Parse()

	if accountSeed != "" && operatorSeed != "" {
		log.Fatal("operator and account cannot both be provided")
	}

	var (
		jwt string
		err error
	)

	if operatorSeed != "" {
		jwt, err = createAccount(operatorSeed, name)
	} else if accountSeed != "" {
		jwt, err = createUser(accountSeed, name)
	} else {
		flag.PrintDefaults()
		return
	}
	if err != nil {
		log.Fatalf("error creating account JWT: %v", err)
	}

	fmt.Println(jwt)
}

func createAccount(operatorSeed, accountName string) (string, error) {
	akp, err := nkeys.CreateAccount()
	if err != nil {
		return "", fmt.Errorf("unable to create account using nkeys: %w", err)
	}

	apub, err := akp.PublicKey()
	if err != nil {
		return "", fmt.Errorf("unable to retrieve public key: %w", err)
	}

	ac := jwt.NewAccountClaims(apub)
	ac.Name = accountName

	// Load operator key pair
	okp, err := nkeys.FromSeed([]byte(operatorSeed))
	if err != nil {
		return "", fmt.Errorf("unable to create operator key pair from seed: %w", err)
	}

	// Sign the account claims and convert it into a JWT string
	ajwt, err := ac.Encode(okp)
	if err != nil {
		return "", fmt.Errorf("unable to sign the claims: %w", err)
	}

	return ajwt, nil
}

func createUser(accountSeed, userName string) (string, error) {
	ukp, err := nkeys.CreateUser()
	if err != nil {
		return "", fmt.Errorf("unable to create user using nkeys: %w", err)
	}

	upub, err := ukp.PublicKey()
	if err != nil {
		return "", fmt.Errorf("unable to retrieve public key: %w", err)
	}

	uc := jwt.NewUserClaims(upub)
	uc.Name = userName

	// Load account key pair
	akp, err := nkeys.FromSeed([]byte(accountSeed))
	if err != nil {
		return "", fmt.Errorf("unable to create account key pair from seed: %w", err)
	}

	// Sign the user claims and convert it into a JWT string
	ujwt, err := uc.Encode(akp)
	if err != nil {
		return "", fmt.Errorf("unable to sign the claims: %w", err)
	}

	return ujwt, nil
}

```

--------------------------------

### Python NKey Authentication

Source: https://github.com/nats-io/nats.docs/blob/master/using-nats/developing-with-nats/connecting/security/nkey.md

Connect to NATS using NKey authentication in Python. This example specifies the path to the NKey seed file.

```Python
nc = NATS()

async def error_cb(e):
    print("Error:", e)

await nc.connect("nats://localhost:4222",
                 nkeys_seed="./path/to/nkeys/user.nk",
                 error_cb=error_cb,
                 )

# Do something with the connection

await nc.close()
```

--------------------------------

### Java NKey Authentication

Source: https://github.com/nats-io/nats.docs/blob/master/using-nats/developing-with-nats/connecting/security/nkey.md

Implement NKey authentication in Java using a custom AuthHandler. This approach allows the client to provide its public key ID and sign challenges.

```Java
NKey theNKey = NKey.createUser(null); // really should load from somewhere
Options options = new Options.Builder()
    .server("nats://localhost:4222")
    .authHandler(new AuthHandler(){
        public char[] getID() {
            try {
                return theNKey.getPublicKey();
            } catch (GeneralSecurityException|IOException|NullPointerException ex) {
                return null;
            }
        }

        public byte[] sign(byte[] nonce) {
            try {
                return theNKey.sign(nonce);
            } catch (GeneralSecurityException|IOException|NullPointerException ex) {
                return null;
            }
        }

        public char[] getJWT() {
            return null;
        }
    })
    .build();
Connection nc = Nats.connect(options);

// Do something with the connection

nc.close();
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