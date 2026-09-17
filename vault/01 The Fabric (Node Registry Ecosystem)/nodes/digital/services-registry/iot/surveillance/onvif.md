---
title: "ONVIF"
type: reference
status: active
tags: [platform, iot, surveillance, protocol, standard]
---
# ONVIF

> Universal IP camera protocol standard. Any ONVIF camera works with any ONVIF-compatible NVR/software.

## What It Is
Open Network Video Interface Forum — industry standard for IP-based security products. If a camera supports ONVIF, it works with Frigate, Scrypted, Blue Iris, HA, and most NVRs.

## Why It Matters
- **Vendor independence** — don't get locked into one brand
- **Auto-discovery** — find cameras on network automatically
- **Standard features** — PTZ, events, profiles, streaming all standardized
- **Buy criteria** — always confirm ONVIF Profile S/T support before purchasing

## Profiles
| Profile | Purpose |
|---------|---------|
| S | Streaming (video/audio) |
| T | Advanced streaming (H.265, imaging, PTZ) |
| G | Recording + storage |
| C | Access control (doors) |
| A | Access control config |

## Discovery
```bash
# Python ONVIF discovery
pip install onvif-zeep
python -c "from onvif import ONVIFCamera; cam = ONVIFCamera('IP', 80, 'user', 'pass')"

# Home Assistant: auto-discovers ONVIF cameras on network
```

## Buying Rule
> Always verify **ONVIF Profile S** minimum. Prefer Profile T (H.265 + PTZ). Avoid cameras that only work with proprietary apps.
