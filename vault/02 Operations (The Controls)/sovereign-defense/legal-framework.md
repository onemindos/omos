---
domain: 01
name: Sovereign Defense — Legal Framework (US Drone Law)
status: active
owner: zeus
tags: [sovereign-defense, legal, faa, drone-law, privacy, use-of-force]
created: 2026-05-13
updated: 2026-05-13
---

# Legal Framework — US Drone Law by Threat Level

**Know exactly where the legal line is before building the system. This is not legal advice — consult an attorney for your specific state and situation.**

---

## Summary: What's Legal at Each Level

| Threat Level | Drone Ops | Thermal/Facial Rec | Deterrence (Siren/Strobe) | Autonomous Behaviors |
| --- | --- | --- | --- | --- |
| **LEVEL 5-4** (patrol, surveillance) | Legal with FAA compliance | Legal on own property | Legal on own property | Legal (supervised) |
| **LEVEL 3** (threat present) | Legal with FAA compliance | Legal on own property | Legal on own property | Legal (supervised) |
| **LEVEL 2** (active intruder) | Legal with FAA compliance | Legal on own property | Legal on own property | Legal (supervised) |
| **LEVEL 1** (WARFIGHTER) | Legal — same FAA rules apply | Legal on own property | Legal on own property | Legal if non-harmful |
| **Weapons on drones** | **ILLEGAL in all 50 states** | — | — | — |
| **Ramming with drone** | **ILLEGAL — assault/battery** | — | — | — |

Key principle: **You can surveil, deter, and document. You cannot use a drone as a weapon.**

---

## FAA Regulations — What Always Applies

### The Hard Limits (All Threat Levels)

These apply 24/7, regardless of threat level. Violating them is a federal crime.

| Rule | Limit | Penalty |
| --- | --- | --- |
| **Maximum altitude** | 400ft AGL (above ground level) | Up to $27,500 civil / criminal charges |
| **Visual line of sight** | Must keep drone in visual sight | Same |
| **Controlled airspace** | Cannot enter Class B/C/D without authorization | Same |
| **Emergency aircraft** | Must give way to manned aircraft | Same |
| **Airports** | 5-mile radius = notify/permission required | Same |
| **Night flight** | Legal if drone has anti-collision lights | No extra penalty |
| **Over people** | Legal if Part 107 compliant or Category 1-4 | Same |

### FAA Part 107 — The Drone Operator License

**Recommendation: Get Part 107 certified.** It costs $175, takes ~3 hours of studying, and removes most recreational limitations. Allows commercial ops, over-people flight, and establishes you as a competent operator for any legal dispute.

Without Part 107: You are a recreational flier. You cannot fly for any commercial purpose.  
With Part 107: You can fly commercially, over people (within limits), at night, in most airspace with LAANC approval.

**LAANC:** FAA's Low Altitude Authorization and Notification Capability. Free app (Aloft/AirMap). Get authorization in minutes for controlled airspace near airports.

---

## Property Surveillance — What's Legal

### On Your Own Property

| Activity | Legal? | Conditions |
| --- | --- | --- |
| Thermal imaging on own property | YES | No restrictions |
| Facial recognition on own property | YES | No restrictions (private property) |
| Recording video/audio on own property | YES | Subject to state wiretapping laws for audio |
| Flying over own property | YES | FAA altitude limits still apply |
| Pointing cameras at public roads/spaces | YES (from private property) | Cannot intercept private comms |
| Automated patrols (scheduled flights) | YES | Must comply with FAA rules |

### Near or Off Your Property

| Activity | Legal? | Conditions |
| --- | --- | --- |
| Flying over neighboring property | GRAY AREA | Below ~400ft = likely trespass; above = federal airspace |
| Recording neighbors' yards | GRAY AREA | Expectation of privacy varies by state |
| Following a vehicle on public roads | GRAY AREA | Surveillance from airspace, not from the vehicle's space |
| Flying over public streets | YES | FAA rules apply, no state obstruction |

---

## Deterrence Systems — What's Legal

### Siren / Strobe / Spotlights

| System | Legal? | Conditions |
| --- | --- | --- |
| Loud siren on drone (120dB) | YES | No federal prohibition. May violate local noise ordinances after hours. |
| Strobe lights on drone | YES | Actually FAA-required for night flight (anti-collision lights) |
| Spotlights / floodlights | YES | Directing at a person's eyes intentionally = potential assault |
| Audio challenge messages | YES | No restriction on pre-recorded audio |
| Laser dazzlers | **NO** | Illegal under 18 U.S.C. § 39A — federal felony to aim laser at aircraft or persons |

### Physical Deterrence

| System | Legal? | Notes |
| --- | --- | --- |
| Drone blocking someone's path | GRAY AREA | Could be construed as false imprisonment |
| Drone hovering over an intruder | YES | Surveillance only, not contact |
| Drone ramming a person | **ILLEGAL** | Assault and battery, federal charges (aircraft weapon) |
| Drone dropping items on a person | **ILLEGAL** | Federal crime + assault |
| Drone carrying a firearm or weapon | **ILLEGAL — ALL 50 STATES** | Criminal charges regardless of use |

---

## Use of Force and Self-Defense Considerations

### The Drone Is Not a Weapon

Under federal law (18 U.S.C. § 32), using an aircraft to harm a person is a federal crime. Your PX4 drone is legally classified as an "aircraft." Ramming, dropping objects, or weaponizing it is a federal felony regardless of self-defense claims.

### What the Drone CAN Do in Self-Defense Situations

| Action | Legal | Why |
| --- | --- | --- |
| Document the threat on video | YES | Evidence preservation |
| Illuminate threat with spotlight | YES | Visibility aid |
| Sound siren near threat | YES | Deterrence, no contact |
| Broadcast GPS position for 911 | YES | Emergency coordination |
| Relay live video to law enforcement | YES | Cooperating with authorities |
| Fly between you and threat | GRAY | No contact — surveillance positioning |
| Fly aggressively near threat (not contacting) | GRAY | Possible harassment claim |

### Recording and Evidence

Always record. SD card footage from your drones is admissible evidence in criminal and civil proceedings. This is one of the most valuable things the system does — independent, timestamped, GPS-tagged video of any incident.

---

## State-Specific Considerations

### States with Drone-Specific Privacy Laws (as of 2026)

| State | Notable Restriction |
| --- | --- |
| **California** | Drone paparazzi law (SB 142) — no harassment with drones for media. Doesn't restrict private property surveillance. |
| **Texas** | H.B. 1643 — restricts drone surveillance over certain private properties (farms, petrochemical). Own property is exempt. |
| **Florida** | § 934.50 — no drone surveillance of people without consent, except law enforcement. Own property use disputed. |
| **North Carolina** | Cannot photograph private property (with reasonable expectation of privacy) without consent. Own property use disputed. |
| **Nevada** | Criminal penalty for harassing livestock with drones. |
| **All others** | Generally follow federal law only. Check your specific state. |

**Best practice:** Stay above the property you own, point cameras inward (toward your own property), and never pursue a fleeing person off your land.

### Audio Recording (Wiretapping)

Drone microphones are subject to state wiretapping laws:

| State Type | Law | Implication |
| --- | --- | --- |
| One-party consent states (majority) | Only Zeus needs to consent to recording conversations | Audio recording on own property = legal |
| Two-party consent states (CA, FL, IL, MD, MI, NV, NH, OR, PA, WA) | All parties must consent | Recording conversations of intruders outdoors = gray area |

Thermal cameras and visual recording have no audio wiretapping issues. The concern is only microphones picking up private conversations.

---

## Facial Recognition — Legal Landscape

### Current Status (2026)

| Context | Legal? | Notes |
| --- | --- | --- |
| On own private property | YES | No federal prohibition. No state prohibition except Illinois (BIPA) and Texas (CUBI) for commercial use. |
| Identifying family / employees you enrolled | YES | Consented or reasonable expectation |
| Storing face templates of intruders | GRAY | No consent obtained, data retention laws vary |
| Sharing facial data with third parties | NO | Privacy laws likely apply |

**Practical approach:** Use facial recognition to identify known people (green list) and flag unknown faces. Do not store templates of people you cannot identify. Delete footage of non-incidents after 30 days.

---

## Engaging Law Enforcement

### During an Incident

If WARFIGHTER activates and police arrive:

1. **Land all drones immediately** or put them in supervised hover out of the way
2. **Stop all autonomous behaviors** — "STAND DOWN"
3. **Offer SD card footage** — law enforcement will appreciate the independent video evidence
4. **Do not obstruct** law enforcement airspace or operations
5. **Identify yourself** as the operator if asked

### Pre-Registration

- Register all drones with the FAA (faa.gov/uas/getting_started) — required for drones >250g ($5 fee, covers all your drones under one registration)
- Label each drone with your FAA registration number
- Keep pilot certificate (Part 107) on your person when flying

---

## Liability Exposure Summary

| Scenario | Your Exposure |
| --- | --- |
| Drone crashes on intruder during autonomous flight | Potential criminal (reckless endangerment) + civil liability |
| Drone spotlight causes someone to trip and fall | Civil liability — negligence |
| Drone siren causes hearing damage | Civil liability — potentially |
| Drone records intruder committing crime | Zero — helps you legally |
| Drone records neighbor without incident | Privacy lawsuit in two-party consent states |
| Drone flies over neighbor's property | FAA airspace violation + civil trespass claim |
| Drone captures evidence used in criminal prosecution | Beneficial — evidence of intruder's crime |

**Insurance:** Standard homeowner's or renter's insurance typically does NOT cover drone incidents. Consider dedicated drone liability insurance (around $500-1,000/year for $1M coverage via companies like Verifly or SkyWatch.AI).

---

## Recommended Legal Preparations

- [ ] Get FAA Part 107 certification ($175 test)
- [ ] Register all drones with FAA (faa.gov/uas)
- [ ] Get drone liability insurance (Verifly / SkyWatch.AI)
- [ ] Consult a local attorney about your state's specific drone and privacy laws
- [ ] Post "Private Property — Drone Surveillance in Use" signage at property entries
- [ ] Document your drone program: serial numbers, registration, insurance, pilot cert
- [ ] Set a 30-day auto-delete policy for non-incident footage
- [ ] Know your local sheriff's department — introduce yourself and your property security system

---

*Full suite: [[README]] | Operating procedures: [[wartime-weapon-protocol]]*
