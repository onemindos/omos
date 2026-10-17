# APEX-SENTINEL

**Repo:** (third-party — no GitHub URL in source)  
**Language:** TypeScript (34 modules)  
**Stack:** Acoustic · RF · RTL-SDR · TDoA · EKF · LSTM  
**Status:** W4 complete, W5 in progress | 397 tests, 92% coverage

Distributed civilian drone detection network targeting FPV combat drones ($400 each) and Shahed-136 loitering munitions ($20K each) on the NATO Eastern Flank. Detects via acoustic sensors, RF fingerprinting, RTL-SDR radio, and Time Difference of Arrival (TDoA) triangulation. Extended Kalman Filter + LSTM for tracking.

## Why it matters for OneMind

Direct C-UAS integration target. The detection pipeline feeds track data that maps naturally onto:
- `evt.defense.drone.detected` / `evt.defense.drone.tracked`
- TAK CoT `a-h-A-*` (hostile aircraft) events via om-overwatch-tak
- Overwatch COP via cot2constellation

## Integration path

APEX-SENTINEL detection output → NATS FABRIC → om-overwatch-tak (CoT) → TAK clients  
or → cot2constellation → Constellation Overwatch COP

## Hardware

RTL-SDR V4, directional microphones, Raspberry Pi edge nodes. NATO Eastern Flank deployment focus but architecture is geography-agnostic.
