# Pick-up: personal daily-driver device + PTT

Parking the thread from 2026-07-03 so it's resumable. Hardware tiers live in
[eud-hardware.md](eud-hardware.md); this is the **personal daily-driver decision** + the **PTT plan**
that came out of it.

---

## 1. Personal Android — daily driver that converts to deep work

Goal: **one Android for business + life** that scales up into a **bigger surface for deep work**.
Different job than a field EUD — durability/field objections don't apply here.

**Decided:**
- Run **stock Android, hardened** — NOT GrapheneOS. Degoogling a daily business phone friction-loads
  banking/notifications/business apps. Sovereignty (Graphene) stays on a **separate dedicated Pixel node**.
- ATAK still runs fine on this as a *personal/planning* tool — just not the dismounted field device.

**Open fork — the whole decision hinges on one question: what does "bigger" mean?**

| If "bigger" = ... | Buy | Why | Cost |
|---|---|---|---|
| **Tablet in-hand, anywhere, no gear** | **Samsung Galaxy Z Fold7** | Self-contained: pocket phone → unfold 8" tablet → DeX to a monitor. 3 modes, 1 device. SD 8 Elite, 12–16GB RAM | ~$2000 |
| **Bigger only at a desk** | **Samsung Galaxy S25 Ultra + DeX** | More device for ~$700 less. Tougher (IP68), S-Pen, bigger battery. DeX → full desktop on any monitor/TV. No tablet-in-hand | ~$1300 |

Fold caveats (fine for daily, NOT field): IP48 hinge (sand/dust), no MIL-STD, no battery swap,
soft inner screen. All livable at pocket/desk.

**Next action:** answer the "bigger = in-hand vs at-a-desk" question → that picks the device.

---

## 2. PTT on the Android

**Backend already exists and is live** — Murmur (Mumble) sovereign voice server on Tailscale
(`100.74.182.117:64738`, LAN `192.168.1.183`). See [../voice/README.md](../voice/README.md).

**Plan:**
- Client: **Mumla** (maintained Mumble Android client). Point at the Murmur server over Tailscale.
  Register each operator as a Mumble user (not SuperUser).
- Trigger: **Bluetooth PTT button** — because Fold7 / S25 Ultra have **no 3.5mm jack** (wired PTT
  would need a USB-C dongle). On-screen PTT or volume-key mapping work as fallback.
- Skip **Zello** (cloud/SaaS — breaks local-first). Skip **TAKTalk** unless we want PTT fused into
  ATAK's UI later — it needs a SIP/VoIP backend we'd have to stand up; Mumble is simpler and already running.

**Hard limit:** Mumble is **IP-based** (needs cell/wifi/Tailscale). **Grid-down = no Mumble voice.**
True off-grid voice = **RF radio** (separate layer). Meshtastic/goTenna = data only, not live voice.

**Next action (optional):** write a `voice/` client-setup note — Mumla config + BT PTT pairing —
alongside the server README, matching the device-loadout pattern.

---

## Quick resume checklist
- [ ] Decide "bigger" = in-hand (Fold7) or at-a-desk (S25 Ultra + DeX)
- [ ] Buy the daily driver; set up **stock hardened**, not Graphene
- [ ] Install **Mumla** → connect to Murmur over Tailscale
- [ ] Get a **Bluetooth PTT button** (no 3.5mm jack on either flagship)
- [ ] (Optional) document RF/off-grid voice layer for grid-down
