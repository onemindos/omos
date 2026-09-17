### Connect to NATS with Token Authentication

Source: https://github.com/nats-io/nats.docs/blob/master/_examples/connect_token.html

Establishes a connection to a NATS server using a provided authentication token. This is supported across Go, Java, JavaScript, Python, Ruby, and TypeScript.

```APIDOC
## CONNECT /nats/connect

### Description
Establishes a client connection to the NATS server using token-based authentication.

### Method
CONNECT

### Endpoint
nats://[token]@[host]:[port]

### Parameters
#### Path Parameters
- **token** (string) - Required - The authentication token required by the NATS server.
- **host** (string) - Required - The server address.
- **port** (integer) - Required - The server port (default 4222).

### Request Example
{ "url": "nats://mytoken@127.0.0.1:4222" }

### Response
#### Success Response (200)
- **status** (string) - Connection established successfully.

#### Response Example
{ "status": "connected", "server": "nats://127.0.0.1:4222" }
```

--------------------------------

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

### Configure Multiple User Authentication

Source: https://github.com/nats-io/nats.docs/blob/master/running-a-nats-service/configuration/securing_nats/auth_intro/username_password.md

Allows defining an array of multiple authorized users in the configuration file to support different credentials for different clients.

```hcl
authorization: {
    users: [
        {user: a, password: b},
        {user: b, password: a}
    ]
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

### Connect to NATS with Token Authentication

Source: https://github.com/nats-io/nats.docs/blob/master/_examples/connect_token.html

Demonstrates how to connect to a NATS server using a provided token for authentication. This is a common security practice for NATS connections. Ensure the token is valid and the server is accessible.

```Go
// Set a token nc, err := nats.Connect("127.0.0.1", nats.Name("API Token Example"), nats.Token("mytoken")) if err != nil { 
	log.Fatal(err) 
}
defer nc.Close() 

// Do something with the connection
```

```Java
Options options = new Options.Builder().                             server("nats://localhost:4222").                             token("mytoken"). // Set a token                             build(); Connection nc = Nats.connect(options);

// Do something with the connection

nc.close();
```

```JavaScript
let nc = NATS.connect({url: `nats://127.0.0.1:${port}`, token: "mytoken!"});
```

```Python
nc = NATS()

await nc.connect(servers=["nats://mytoken@demo.nats.io:4222"])

# Do something with the connection.
```

```Ruby
NATS.start(token: "deadbeef") do |nc|
   puts "Connected using token"
end
```

```TypeScript
let nc = await connect({url: server.nats, token: "mytoken"});
```

--------------------------------

### NATS Leafnode Multi-User Authorization Configuration

Source: https://github.com/nats-io/nats.docs/blob/master/running-a-nats-service/configuration/leafnodes/leafnode_conf.md

Configures multiple users for leafnode authentication, allowing different leafnodes to connect with distinct credentials and be bound to different accounts. This provides more granular control in complex environments.

```hcl
leafnodes {
  port: 6222
  authorization {
    users = [
      {user: leaf1, password: secret, account: account1}
      {user: leaf2, password: secret, account: account2}
    ]
  }
}
```

--------------------------------

### Go: User Provisioning Client Process

Source: https://github.com/nats-io/nats.docs/blob/master/running-a-nats-service/nats_admin/jwt.md

Initiates a user provisioning request to the service. It generates a user key pair, requests a JWT from the provisioning service, and then connects to NATS using the obtained JWT and private key for authentication.

```go
func startUserProcess(userRequestChan chan userRequest, obtainAuthorizationCb func() interface{}) {
    requestUser := func(userRequestChan chan userRequest, authInfo interface{}) (jwtAuthOption nats.Option) {
        userPublicKey, _, userKeyPair := generateUserKey()
        respChan := make(chan string)
        // request jwt
        userRequestChan <- userRequest{
            respChan,
            userPublicKey,
            authInfo,
        }
        userJWT := <-respChan // wait for response
        // userJWT and userKeyPair can be used in conjunction with this nats.Option
        jwtAuthOption = nats.UserJWT(func() (string, error) {
            return userJWT, nil
        },
            func(bytes []byte) ([]byte, error) {
                return userKeyPair.Sign(bytes)
            },
        )
        // Alternatively you can create a creds file and use it as nats.Option
        return
    }
    go func() {
        jwtAuthOption := requestUser(userRequestChan, obtainAuthorizationCb())
        nc, err := nats.Connect("nats://localhost:4222", jwtAuthOption)
        if err != nil {
            return
        }
        defer nc.Close()
        time.Sleep(time.Second) // simulate work one would want to do
    }()
}
```

--------------------------------

### NATS Server Authorization Configuration with Bcrypt

Source: https://github.com/nats-io/nats.docs/blob/master/using-nats/nats-tools/nats_cli/README.md

Example of how to configure user authorization in a NATS server, using a bcrypt hash for the password.

```hcl
authorization {
    user: derek
    password: $2a$11$3kIDaCxw.Glsl1.u5nKa6eUnNDLV5HV9tIuUp7EHhMt6Nm9myW1aS
  }
```

--------------------------------

### Go: Orchestrating Distributed User Creation

Source: https://github.com/nats-io/nats.docs/blob/master/running-a-nats-service/nats_admin/jwt.md

Orchestrates the distributed user creation process by starting the provisioning service and multiple user client processes. It demonstrates the overall flow of requesting and obtaining JWTs for NATS authentication.

```go
func RequestUserDistributed() {
    reqChan := startUserProvisioningService(IsTokenAuthorized)
    defer close(reqChan)
    // start multiple user processes
    for i := 0; i < 4; i++ {
        startUserProcess(reqChan, ObtainAuthorizationToken)
    }
    time.Sleep(5 * time.Second)
}
```

--------------------------------

### Connect to NATS with Token in URL

Source: https://github.com/nats-io/nats.docs/blob/master/_examples/connect_token_url.html

Demonstrates how to connect to a NATS server using a token embedded directly within the connection URL. This method is useful for simple authentication scenarios where the token can be safely included in the URL.

```Go
// Token in URL
nc, err := nats.Connect("mytoken@localhost")
if err != nil {
	log.Fatal(err)
}
defer nc.Close()

// Do something with the connection
```

```Java
Connection nc = Nats.connect("nats://mytoken@localhost:4222");//Token in URL

// Do something with the connection

nc.close();
```

```JavaScript
let url = `nats://mytoken!@127.0.0.1:${port}`;
let nc = NATS.connect({url: url});
```

```Python
nc = NATS()

await nc.connect(servers=["nats://mytoken@demo.nats.io:4222"])

# Do something with the connection.
```

```Ruby
NATS.start("deadbeef@127.0.0.1:4222") do |nc|
   puts "Connected using token!"
end
```

```TypeScript
let url = `nats://:mytoken!@127.0.0.1:${port}`;
let nc = await connect({url: url});
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

### Issue User JWT in C#

Source: https://github.com/nats-io/nats.docs/blob/master/running-a-nats-service/nats_admin/jwt.md

This C# code snippet demonstrates how to issue a NATS user JWT. It requires the NATS.NKeys and SimpleBase packages. The function takes a user's public key and uses account signing keys to generate a valid JWT, including calculating the JWT ID (jti) based on the claim content. It also formats the output into a standard .creds file.

```csharp
// dotnet add package NATS.NKeys --prerelease
// dotnet add package SimpleBase
using NATS.NKeys;
using System.Security.Cryptography;
using System.Text;

string creds = IssueUserCreds();
Console.WriteLine(creds);

static string IssueUserJwt(string userKeyPub)
{
    // Load account signing key and account identity for
    // the account you wish to issue users for
    const string accSeed = "SAANWFZ3JINNPERWT3ALE45U7GYT2ZDW6GJUIVPDKUF6GKAX6AISZJMAS4";
    const string accId = "ACV63DGCZGOIT3P5ZA7PQT3KYJ6UDFFHZ7KETHYMDMZ4N44KYAQ2ZZ5F";
    KeyPair accountSigningKey = KeyPair.FromSeed(accSeed);
    string accSigningKeyPub = accountSigningKey.GetPublicKey();

    // Use nsc to create a user any way you like.
    // Export the user as json using:
    // nsc describe user --name <user name> --account <account name> --json
    // Turn the output into a format string and replace values you want replaced.
    // Fields that need to be replaced are:
    // iat (issued at), iss (issuer), sub (subject) and jti (claim hash)
    const string claimFmt = @"{ 
  \"iat\": {0},
  \"iss\": \"{1}\",
  \"jti\": \"{2}\",
  \"name\": \"{3}\",
  \"nats\": {
    \"data\": -1,
    \"issuer_account\": \"{4}\",
    \"payload\": -1,
    \"pub\": {{}},
    \"sub\": {{}},
    \"subs\": -1,
    \"type\": \"user\",
    \"version\": 2
  },
  \"sub\": \"{3}\" 
}";
    const string header = @"{ 
  \"typ\":\"JWT\",
  \"alg\":\"ed25519-nkey\"
}";

    // Issue At time is stored in unix seconds
    long issuedAt = DateTimeOffset.Now.ToUnixTimeSeconds();
    
    // Generate a claim without jti so we can compute jti off of it
    string claim = string.Format(
        claimFmt,
        issuedAt,
        accSigningKeyPub,
        "", /* blank jti */
        userKeyPub,
        accId);
    
    // Compute jti, a base32 encoded sha256 hash
    string jti = SimpleBase.Base32.Rfc4648.Encode(
        SHA256.Create().ComputeHash(Encoding.UTF8.GetBytes(claim)),
        false);
    
    // recreate full claim with jti set
    claim = string.Format(
        claimFmt,
        issuedAt,
        accSigningKeyPub,
        jti,
        userKeyPub,
        accId
    );
    
    // all three components (header/body/signature) are base64url encoded
    string encHeader = ToBase64Url(Encoding.UTF8.GetBytes(header));
    string encBody = ToBase64Url(Encoding.UTF8.GetBytes(claim));
    
    // compute the signature off of header + body (. included on purpose)
    byte[] sig = Encoding.UTF8.GetBytes($"{encHeader}.{encBody}");
    var signature = new byte[64];
    accountSigningKey.Sign(sig, signature);
    string encSig = ToBase64Url(signature);
    
    // append signature to header and body and return it
    return $"{encHeader}.{encBody}.{encSig}";
}

static string IssueUserCreds()
{
    // Generate a user NKEY for the new user.
    // The private portion of the NKEY is not needed when issuing the jwt.
    // Therefore generating the key can also be done separately from the JWT.
    // Say by the requester.
    KeyPair userSeed = KeyPair.CreatePair(PrefixByte.User);
    string userKeyPub = userSeed.GetPublicKey();
    string jwt = IssueUserJwt(userKeyPub);
    
    // return jwt and corresponding user seed as creds
    return $
----BEGIN NATS USER JWT-----
{jwt}
------END NATS USER JWT------

************************* IMPORTANT *************************
    NKEY Seed printed below can be used to sign and prove identity.
    NKEYs are sensitive and should be treated as secrets.

-----BEGIN USER NKEY SEED-----
{userSeed.GetSeed()}
------END USER NKEY SEED------

*************************************************************";
}

static string ToBase64Url(byte[] input)
{
    var stringBuilder = new StringBuilder(Convert.ToBase64String(input).TrimEnd('='));
    stringBuilder.Replace('+', '-');
    stringBuilder.Replace('/', '_');
    return stringBuilder.ToString();
}

```

--------------------------------

### CLI: Generate Activation Token

Source: https://github.com/nats-io/nats.docs/blob/master/using-nats/nats-tools/nsc/streams.md

Generates a JWT-based activation token that grants a target account permission to import a specific stream subject.

```APIDOC
## NSC generate activation

### Description
Generates an activation token for a foreign account to import a private stream. This token defines the permissions and subject mapping for the imported stream.

### Method
CLI Command

### Parameters
#### Flags
- **--account** (string) - Required - The account name that owns the stream export.
- **--target-account** (string) - Required - The public key of the account importing the service.
- **--subject** (string) - Required - The subject where the stream will be published to the target account.
- **-o** (string) - Required - The output file path for the generated JWT.

### Request Example
`nsc generate activation --account A --target-account AAM46E3YF5WOZSE5WNYWHN3YYISVZOSI6XHTF2Q64ECPXSFQZROJMP2H --subject private.abc.AAM46E3YF5WOZSE5WNYWHN3YYISVZOSI6XHTF2Q64ECPXSFQZROJMP2H -o /tmp/activation.jwt`

### Response
#### Success Response
- **Output** (file) - A JWT file containing the activation details, which can be shared with the importing account.
```

--------------------------------

### NATS Leafnode Singleton Authorization Configuration

Source: https://github.com/nats-io/nats.docs/blob/master/running-a-nats-service/configuration/leafnodes/leafnode_conf.md

Configures a single user for leafnode authentication. The leafnode connects using the provided username and password, and is bound to the specified account. This configuration is useful for simple deployments.

```hcl
leafnodes {
  port: 6222
  authorization {
    user: leaf
    password: secret
    account: TheAccount
  }
}
```

--------------------------------

### Generate NATS User JWT and Credentials

Source: https://github.com/nats-io/nats.docs/blob/master/running-a-nats-service/nats_admin/jwt.md

This snippet demonstrates how to create a user keypair, generate user claims, sign them with an account key, and format the output into a NATS credentials file using the C# JWT library.

```csharp
KeyPair accountSigningKey = KeyPair.FromSeed(accSeed);
KeyPair ukp = KeyPair.CreatePair(PrefixByte.User);
string upk = ukp.GetPublicKey();
NatsUserClaims uc = jwt.NewUserClaims(upk);
uc.User.IssuerAccount = accId;
string userJwt = jwt.EncodeUserClaims(uc, accountSigningKey);
string userSeed = ukp.GetSeed();
string creds = jwt.FormatUserConfig(userJwt, userSeed);
Console.WriteLine(creds);
```

--------------------------------

### Start NATS Server with Token Authentication Flag

Source: https://github.com/nats-io/nats.docs/blob/master/running-a-nats-service/configuration/securing_nats/auth_intro/tokens.md

Demonstrates starting the NATS server process using the command-line interface with a specified authentication token.

```shell
nats-server --auth s3cr3t
```

--------------------------------

### Generate User JWT in Go

Source: https://github.com/nats-io/nats.docs/blob/master/running-a-nats-service/nats_admin/jwt.md

This Go function generates a user JWT, allowing configuration of publish/subscribe permissions, expiration, and issuer account. It requires the user's public key, the account's public key, and the account's signing key. The generated JWT is then encoded using the account's signing key. Validation is performed before encoding to ensure the claims are valid.

```go
func generateUserJWT(userPublicKey, accountPublicKey string, accountSigningKey nkeys.KeyPair) (userJWT string) {
    uc := jwt.NewUserClaims(userPublicKey)
    uc.Pub.Allow.Add("subject.foo") // only allow publishing to subject.foo
    uc.Expires = time.Now().Add(time.Hour).Unix() // expire in an hour
    uc.IssuerAccount = accountPublicKey
    vr := jwt.ValidationResults{}
    uc.Validate(&vr)
    if vr.IsBlocking(true) {
        panic("Generated user claim is invalid")
    }
    var err error
    userJWT, err = uc.Encode(accountSigningKey)
    if err != nil {
        return ""
    }
    return
}
```