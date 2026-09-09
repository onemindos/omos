# bachmach-pqc-iot-sentinel — PQC-Hardened IoT Telemetry

**Repo:** https://github.com/AuditorSEC-Initiative/bachmach-pqc-iot-sentinel
**Language:** C (ESP32/RP2040 firmware)
**License:** Unknown
**Topics:** esp32, pqc, post-quantum, iot, iot-security, nats, brave1, ukraine

## What it is

Post-quantum cryptography (PQC) hardened IoT telemetry system for edge nodes. Runs on ESP32/RP2040 microcontrollers. Publishes sensor data over NATS JetStream with ML-KEM (key encapsulation) and ML-DSA (digital signatures) — both NIST-standardized post-quantum algorithms. Developed as a Ukraine defense pilot (BRAVE1 / MaJoR EDF program).

## PQC Algorithms Used

| Algorithm | NIST Standard | Purpose |
|-----------|---------------|---------|
| ML-KEM (Kyber) | FIPS 203 | Key encapsulation — replaces ECDH in TLS handshake |
| ML-DSA (Dilithium) | FIPS 204 | Digital signatures — authenticates every NATS message |

## Why it's relevant to OneMind

| Use case | How |
|----------|-----|
| Edge node security | barn Pi, office box, robot, drone — all run `omd` and connect to NATS. PQC hardens those connections against future quantum attacks |
| ESP32 sensor nodes | OneMind already targets ESP32 for IoT telemetry (`tel.>` on FABRIC account). This shows how to do it with PQC |
| NATS JetStream from microcontrollers | Pattern for publishing telemetry from resource-constrained devices |
| Defense-grade security posture | OneMind targets defense/dual-use — PQC is increasingly a requirement for SBIR/SBIR/DoD contracts |
| BRAVE1 alignment | Ukraine defense tech ecosystem — relevant for OneMind's defense positioning |

## What's possible

1. **PQC-harden the NATS leaf node connections** — replace standard TLS with ML-KEM for key exchange on edge nodes (barn Pi, drones)
2. **Sign telemetry messages with ML-DSA** — every `tel.*` message carries a post-quantum signature, verifiable by ORP
3. **ESP32 → NATS JetStream firmware template** — drop-in starting point for custom sensor hardware publishing to the OneMind fabric
4. **Zero-trust edge identity** — each node has its own ML-DSA keypair, verified at the NATS callout auth layer (pairs with `callout.go`)

## Notes

- NATS JetStream transport confirmed working on ESP32/RP2040
- `liboqs` (Open Quantum Safe) is the C library used — also has Go bindings (`go-oqs`) for use in omd/bridge services
- BRAVE1 is Ukraine's defense tech accelerator — project may have further documentation not public
