# Offline Services - Quick Reference Card

> **1-page cheat sheet for when you need answers fast**

---

## Two Systems, Different Purposes

| | Meshtastic | Emergency Kit |
|---|------------|---------------|
| **Icon** | 📡 | 💾 |
| **For** | THE FLEET (drones/robots) | YOU (human operator) |
| **Purpose** | Communication backup | Knowledge backup |
| **Form** | LoRa radios (fixed) | USB drive (portable) |
| **When** | Internet down | Everything down |
| **Cost** | $180-1,415 | $30-75 |
| **Issues** | OMOS-132, 133, 134 | OMOS-135, 136 |

---

## Meshtastic (📡 Fleet Communication)

**What:** LoRa mesh network for drones/robots/sensors

**Key Features:**
- 5-10km range between devices
- MQTT bridge to NATS
- Automatic failover
- Zero monthly costs

**Use When:**
- Internet goes down
- Fleet needs to keep operating
- Remote sensors reporting
- Emergency commands to drones

**Hardware:**
- 2× RAK19007 kits ($130) - starter
- 10+ devices ($1,415) - full deployment

**Docs:**
- `MESHTASTIC-MESH-INTEGRATION.md` - Architecture
- `MESHTASTIC-STARTER-GUIDE.md` - Setup
- `MESHTASTIC-USE-CASES.md` - Examples

---

## Emergency Kit (💾 Personal Knowledge)

**What:** Portable AI + knowledge on USB drive

**Key Features:**
- 3× AI models (Phi-3, Llama, Qwen)
- 4GB survival manuals
- OneMind OS docs
- Works on any laptop
- 100% offline

**Use When:**
- All systems down
- Need to troubleshoot
- Bug-out scenario
- Learning/training
- Gift for preppers

**Hardware:**
- 1× USB drive ($30-75)
- That's it!

**Docs:**
- `ONEMIND-EMERGENCY-KIT.md` - Complete guide

---

## The 4-Layer Stack

```
Layer 1: Cloud NATS ────────────→ Primary
            ↓ (internet fails)
Layer 2: Meshtastic ────────────→ Fleet Backup
            ↓ (devices fail)
Layer 3: Edge Ollama ───────────→ Local AI
            ↓ (everything fails)
Layer 4: Emergency USB ─────────→ Personal Backup
```

---

## Common Questions

**Q: Which do I need first?**
A: Meshtastic (fleet more critical)

**Q: Can I use just one?**
A: Yes, but together = complete resilience

**Q: What's Bunker AI?**
A: Competitor ($89). Buying to learn from.

**Q: Do they overlap?**
A: No - different purposes, work together

---

## Quick Start Commands

```bash
# Create all Linear issues
python3 scripts/linear/create_meshtastic_issues.py
python3 scripts/linear/create_emergency_kit_issues.py

# Order hardware
# Meshtastic: RAK Wireless Store ($130)
# Emergency Kit: Amazon USB drive ($30-75)
```

---

## Total Investment

| Item | Cost | Monthly |
|------|------|---------|
| Meshtastic Starter | $180 | $0 |
| Meshtastic Full | $1,415 | $0 |
| Emergency Kit | $30-75 | $0 |
| **Grand Total** | **$210-1,490** | **$0** |

**Timeline:** 14 weeks (can do in phases)

---

## When to Use What

| Situation | Meshtastic | Emergency Kit |
|-----------|:----------:|:-------------:|
| Internet down | ✅ | - |
| Devices offline | - | ✅ |
| In the field | ✅ | - |
| Bug-out | - | ✅ |
| Troubleshoot | - | ✅ |
| Learning | - | ✅ |
| Fleet ops | ✅ | - |

---

## Master Document

**Start here:** `OFFLINE-SERVICES-OVERVIEW.md`
- Complete explanation of everything
- All documentation links
- Full cost breakdown
- Implementation timeline


---

## Related

- [[00-24 SO (Sovereign Operations)/02 Infrastructure (VPS, Servers & Networking) (The Body)/02 Infrastructure|02 Infrastructure]]
