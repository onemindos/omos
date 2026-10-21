# PEAT — Defense Unicorns Mesh Protocol

**Repo:** github.com/defenseunicorns/peat  
**Language:** Rust  
**License:** Apache 2.0  
**Org:** Defense Unicorns (makers of Zarf, PEPR)  
**Role in OneMind:** DDIL transport mesh — connects nodes when NATS is unavailable

Mesh protocol that connects heterogeneous systems — phones, servers, sensors, embedded devices, AI models — into a coordinated whole, across any transport, even when the network is degraded or denied (DDIL).

## Components

- `peat` — core protocol (Rust)
- `peat-schema` — message schema definitions
- `peat-btle` — Bluetooth LE transport
- `peat-gateway` — protocol gateway/bridge
- `peat-lite` — lightweight embedded variant
- `peat-mesh` — mesh routing layer
- `peat-node` — node daemon

## Why it matters for OneMind

NATS requires connectivity. PEAT fills the gap when the network is degraded or denied:
- Field devices behind broken links
- Bluetooth/BLE mesh between nearby nodes
- Relay across heterogeneous transports (WiFi → BLE → LoRa → TCP)
- Sync when connectivity restores (store-and-forward)

## Integration point

PEAT sits BELOW NATS in the OneMind transport stack:

```
OneMind agents / services
        ↓
    NATS (primary)
        ↓ (DDIL fallback)
    PEAT mesh
        ↓
  BLE / LoRa / serial / degraded TCP
```

When connectivity restores, PEAT drains its buffer back into NATS JetStream.

## See Also

Defense Unicorns also makes:
- **Zarf** — K8s airgap packaging (how OneMind deploys to disconnected environments)
- **PEPR** — K8s policy enforcement
