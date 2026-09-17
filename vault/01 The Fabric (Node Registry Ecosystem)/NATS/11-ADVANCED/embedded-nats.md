---
title: "Embedded NATS"
section: "11-ADVANCED"
type: note
tags: [nats, embedded, golang, in-process, edge, server, testing]
source: https://docs.nats.io/using-nats/developer/embedding
created: 2026-06-14
---

# Embedded NATS

Because both `nats-server` and the Go client are Go libraries, you can **run a full NATS server inside your own Go program** — no separate process, no external dependency. One binary *is* the broker, the client, and your application logic. This is unique among major messaging systems and is the foundation of zero-dependency edge nodes.

See also: [[NATS/08-CLIENTS-SDKS/go|Go Client]] · [[NATS/06-TOPOLOGY-SCALING/leaf-nodes|Leaf Nodes]] · [[NATS/11-ADVANCED/benchmarking|Benchmarking]]

---

## Why embed

- **One artifact** — ship a single static binary that contains the broker. No "install NATS first."
- **Edge / offline-first** — a barn Pi or robot runs its own in-process bus; local clients work even with no network, then a leaf link syncs to the hub.
- **Tests** — spin a real server in a unit test, no Docker, no fixtures.
- **CLIs / desktop apps** — bundle a bus without asking the user to run infrastructure.

---

## Minimal embedded server (Go)

```go
import (
    "github.com/nats-io/nats-server/v2/server"
    "github.com/nats-io/nats.go"
)

func main() {
    // 1. Configure & start the server in-process
    opts := &server.Options{
        ServerName: "omd-barn",
        Host:       "127.0.0.1",
        Port:       4222,
        JetStream:  true,
        StoreDir:   "/var/lib/omd/js",
        // LeafNode: server.LeafNodeOpts{ Remotes: []*server.RemoteLeafOpts{...} },
    }
    ns, err := server.NewServer(opts)
    if err != nil { panic(err) }

    go ns.Start()
    if !ns.ReadyForConnections(5 * time.Second) {
        panic("nats-server not ready")
    }
    defer ns.Shutdown()

    // 2. Connect in-process (no network hop with InProcessServer)
    nc, _ := nats.Connect("", nats.InProcessServer(ns))
    defer nc.Drain()

    nc.Publish("hello", []byte("from inside one binary"))
    nc.Flush()
}
```

Two things make this powerful:
- **`server.NewServer(opts)`** — the entire broker as a library object you control (start, shutdown, reload).
- **`nats.InProcessServer(ns)`** — the client talks to the server **without a TCP round-trip** (in-memory pipe), so embedded pub/sub is essentially free.

---

## Embed as a leaf node (the edge pattern)

Give the embedded server a **leaf remote** pointing at the hub, plus a local **JetStream domain**, and you get the full OneMind edge node in one binary:

```go
opts := &server.Options{
    ServerName: "omd-vehicle",
    JetStream:  true,
    JetStreamDomain: "vehicle",            // isolated local JS domain
    StoreDir:   "/var/lib/omd/js",
}
opts.LeafNode.Remotes = []*server.RemoteLeafOpts{{
    URLs:        parseURLs("nats-leaf://hub:7422"),
    Credentials: "/etc/omd/leaf.creds",
}}
```

Local producers publish to **local** JetStream (durable offline); a hub **source** pulls it up on reconnect. See [[NATS/06-TOPOLOGY-SCALING/leaf-nodes|Leaf Nodes]] and [[NATS/06-TOPOLOGY-SCALING/domains|JetStream Domains]].

---

## Embedded server for tests

```go
func RunTestServer(t *testing.T) (*server.Server, *nats.Conn) {
    opts := &server.Options{Port: -1, JetStream: true,  // Port -1 = random free port
        StoreDir: t.TempDir()}
    ns, err := server.NewServer(opts)
    if err != nil { t.Fatal(err) }
    go ns.Start()
    if !ns.ReadyForConnections(2 * time.Second) { t.Fatal("not ready") }
    t.Cleanup(ns.Shutdown)
    nc, _ := nats.Connect(ns.ClientURL())
    return ns, nc
}
```

`Port: -1` grabs a random free port — perfect for parallel tests. The official servers use exactly this for their own test suites.

---

## Gotchas
- **Version coupling** — embedding pins you to a specific `nats-server` Go module version; upgrade deliberately and test.
- **Lifecycle** — always `Shutdown()` the server and `Drain()` the client; on edge devices wire these to signal handlers.
- **Resource accounting** — the broker shares your process's memory/FDs; size JetStream limits for the device.
- **Only Go** — embedding is a Go capability. Other languages run `nats-server` as a child process instead.
- **`ReadyForConnections`** — always gate startup on it; don't connect before the server is listening.

---

## OneMind Relevance (central)
- **`omd --embed`** edge nodes are exactly this: one Go binary that runs a leaf NATS server + local JetStream + the device's agent logic, with zero extra processes to supervise. This is how "hundreds of nodes" stays operationally sane.
- In-process transport means a robot's local control loop pays no network cost for on-device messaging, while still being a first-class fabric participant via its leaf link.
- Embedded test servers back our Go service test suites — real bus, no Docker.

## Sources
- Embedding NATS — https://docs.nats.io/using-nats/developer/embedding
- nats-server (server pkg) — https://pkg.go.dev/github.com/nats-io/nats-server/v2/server
- nats.go InProcessServer — https://pkg.go.dev/github.com/nats-io/nats.go
- java-nats-server-runner (JVM analogue) — https://github.com/nats-io/java-nats-server-runner
