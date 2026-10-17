# Autonats — Go NATS Service Mesh Code Generator

**Repo:** github.com/zyra/autonats  
**Language:** Go  
**Role in OneMind:** Generate NATS pub/sub service mesh boilerplate from Go interface declarations

Parses Go interface declarations annotated with `@nats:server` and generates the full publisher/subscriber service mesh. Fire-and-forget for void methods, request/reply for methods that return values.

## Usage

```go
// @nats:server User
type UserService interface {
    GetById(ctx context.Context, id string) (*User, error)  // → request/reply
    NotifyEvent(ctx context.Context, event Event)            // → fire-and-forget
}
```

Generates: NATS subject routing, serialization, error handling, publisher + subscriber boilerplate.

## Why it matters for OneMind

OneMind services on NATS all follow the same request/reply + fire-and-forget pattern. Autonats eliminates the boilerplate for Go services in omos-services. Define the interface, generate the NATS wiring, implement the logic.

## Integration point

Used during `omos-services` Go service development. Annotate service interfaces, run `go generate`, get NATS micro-service scaffolding that matches the `svc.{service}.{endpoint}` subject grammar.
