# EUD Hardware Selection

Which physical devices to buy for the OneMind stack, by mission tier. This is the **hardware**
companion to [loadouts/](loadouts/) (software builds) and [registry.md](registry.md) (roster).
Buy the right chassis first; the loadout goes on top.

> Rule of thumb: **phone = dismounted / on-body, tablet = command post / vehicle / planning.**
> Don't field a flagship in a case where a rugged unit belongs, and don't put a TOC planning
> surface in a chest rig.

## Why not a Samsung "Tactical Edition"?
S23 TE / XCover6 Pro TE are gov/mil-restricted (special-request, "not for civilian use") and
scalped on the secondary market. You get ~parity from the **retail** devices below without the
procurement lock. Skip the TE hunt.

---

## Phone tiers

| Tier | Device | Why | Tradeoff |
|---|---|---|---|
| **Primary (sovereign)** | **Pixel 9 Pro + GrapheneOS** | Top raw ATAK performance; degoogled OS = OPSEC fit for OneMind; hardware modem isolation | Smaller tactical accessory ecosystem than Samsung |
| **Field / rugged** | **Samsung Galaxy XCover 7 Pro** | True TE alternative you can buy retail. MIL-STD-810H, IP68, **user-replaceable battery**, programmable key (glove launch) | Mid-tier chip; fine for ATAK, not a flagship |
| **Performance + ecosystem** | **Samsung Galaxy S24 / S25** | SD 8-gen, 8–12GB RAM, widest mount/case support (Juggernaut, Bunker/MAAK); base of the actual TE | Glass slab — needs a case, no battery swap |
| **Turnkey** | **GoTAK FALCON / EUD V2/V3** | ATAK + plugins + server pre-paired, root-enabled, no build time | Costs more; vendor lock to their ecosystem |

**Go-to call:** Pixel 9 Pro (GrapheneOS) as primary EUD; XCover 7 Pro for anything that goes hard
in the field.

### Sovereignty note
GrapheneOS = degoogled, duress-wipe, per-app network/sensor control. This is the OPSEC-correct base
for the OneMind primary. **GrapheneOS is phone-only** — it does not support tablets (see below).

---

## Tablet tiers

| Tier | Device | Why | Tradeoff |
|---|---|---|---|
| **Field / rugged** | **Samsung Galaxy Tab Active 5** | Tablet XCover. MIL-STD-810H, IP68, **replaceable battery**, programmable key, S-Pen (glove/wet). GoTAK sells it as "EUD TA5" | 8" TFT, mid chip |
| **TOC / planning** | **Samsung Galaxy Tab S9 / S10 (or FE)** | Big bright AMOLED, flagship SoC + big RAM, DeX = planning workstation | Not rugged; stays inside the wire |
| **Budget field surface** | **Vortex ZTAB10** (owned) | 10.1", Android 13, LTE, multi-GNSS — fine spare / non-critical map node | Weak CPU (MT8766), 4GB RAM, low-res panel |
| **Turnkey** | **GoTAK Pro Tab** | 10.3" rugged, 4G, 66W, hotkey, pre-provisioned | Vendor cost/lock |

**Go-to call:** Tab Active 5 for a fielded tablet EUD; Tab S9/S10 for the command post.

### Sovereignty gap (tablets)
No clean degoogled tablet path — GrapheneOS is phone-only and Pixel Tablet isn't supported the same
way. Options: (a) run stock Android **hardened manually**, or (b) treat the tablet as a **dumb map
surface** and keep sensitive comms on the Pixel phone. For OneMind, (b) is the safer posture.

---

## OneMind EUD tier map

| Tier | Role | Hardware | OS | Loadout | Group/role |
|---|---|---|---|---|---|
| T1 — Dismounted primary | Team Lead on-body | Pixel 9 Pro | GrapheneOS | [android-full](loadouts/android-full.md) | Red / Team Lead |
| T2 — Field rugged | Dismounted member | XCover 7 Pro / Tab Active 5 | Stock (hardened) | [android-full](loadouts/android-full.md) | Red / Team Member |
| T3 — TOC / planning | Command post surface | Tab S9 / S10 | Stock | android-full (SA + planning plugins) | Red / Team Member |
| T4 — Spare / non-critical | Backup map node | Vortex ZTAB10 | Stock Android 13 | minimal (takchat, datasync) | Red / Team Member |

All tiers connect to TAK server over **Tailscale @ 100.74.182.117:8089 (SSL/TLS)** — see
[registry.md](registry.md). Provision = issue cert (`deploy/scripts/tak-add-device.sh`) → import
soft-cert package. Update registry status when live.

## Buy-priority (if starting fresh)
1. **Pixel 9 Pro** → flash GrapheneOS → T1 primary.
2. **XCover 7 Pro** → T2 field phone (battery swaps in the field = the whole point).
3. **Tab Active 5** → T2/T3 rugged tablet when you need screen real estate.
4. Tab S9/S10 only if you're standing up a real TOC.

## Selection checklist (any candidate device)
- [ ] Android 5.0+ / GLES 3.0 (ATAK minimum — all modern devices pass)
- [ ] ≥6GB RAM for stacked plugins (video/DF/ballistics); 4GB = basic SA only
- [ ] Multi-constellation GNSS (GPS + GLONASS + Galileo), not GPS-only
- [ ] MIL-STD-810H + IP68 **if it leaves the wire**
- [ ] Replaceable battery **if multi-day / no recharge in AO**
- [ ] Mount/case ecosystem exists for the model (check before buying)
- [ ] Carrier-unlocked; avoid locked MVNO SIMs (Total Wireless etc.)
