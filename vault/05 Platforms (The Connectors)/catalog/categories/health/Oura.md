# Oura Ring

> **Smart ring tracking sleep, HRV, readiness, and activity**

## Overview

Oura Ring is the primary biometric wearable in the One Mind health stack. Tracks sleep stages, heart rate variability (HRV), resting heart rate, body temperature, and activity. Powers the HP domain readiness and performance scoring. Ring Gen 3/4 with REST API access.

## Status

| Attribute   | Value                        |
| ----------- | ---------------------------- |
| Status      | ✅ Active                     |
| Category    | Health / Wearable            |
| Type        | Smart Ring                   |
| Integration | REST API, iOS App            |

## Use Cases in OneMind

1. **Sleep Scoring** — Nightly sleep quality data feeds into HP domain metrics
2. **Readiness Score** — Daily readiness informs training and work intensity
3. **HRV Tracking** — Nervous system recovery indicator for HP optimization
4. **Body Temp** — Illness detection and circadian rhythm tracking
5. **Activity Rings** — Steps, active calories, training volume

## Key Metrics

| Metric           | What It Measures                      |
| ---------------- | ------------------------------------- |
| Sleep Score      | Overall sleep quality (0-100)         |
| Readiness Score  | Daily recovery and energy (0-100)     |
| HRV              | Heart rate variability (ms)           |
| RHR              | Resting heart rate (bpm)              |
| Body Temp Dev    | Nightly temperature deviation (°C)    |

## API

```bash
# Get today's readiness
curl -H "Authorization: Bearer OURA_TOKEN" \
  https://api.ouraring.com/v2/usercollection/daily_readiness

# Get sleep data
curl -H "Authorization: Bearer OURA_TOKEN" \
  "https://api.ouraring.com/v2/usercollection/sleep?start_date=2026-04-10"
```

## Links

- Website: [ouraring.com](https://ouraring.com)
- API Docs: [cloud.ouraring.com/v2/docs](https://cloud.ouraring.com/v2/docs)


---

## Related

- [[00-24 SO (Sovereign Operations)/02 Infrastructure (VPS, Servers & Networking) (The Body)/02 Infrastructure|02 Infrastructure]]
