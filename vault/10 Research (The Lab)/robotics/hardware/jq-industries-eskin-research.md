---
company: JQ Industries
category: Robotics Hardware — Tactile Sensing / Electronic Skin
website: https://kr-asia.com/with-fresh-funding-jq-industries-bets-on-fiber-sensors-to-give-machines-a-sense-of-touch
status: research
added: 2026-05-24
tags: [robotics, eskin, tactile-sensing, humanoid, china, unitree, sensors, haptics]
---

# JQ Industries — Fabric-Based Electronic Skin for Humanoid Robots

## Overview

JQ Industries (威海巨桥工业技术 / Weihai Juqiao Industrial Technology) is a Shanghai-founded, Weihai-based Chinese precision sensing startup that developed a **fabric-based electronic skin (e-skin)** for humanoid robots. Founded in early 2024 by Zhang Jingqi, the company builds flexible fiber sensors using plant-derived textile materials, woven with conductive fibers to give robots a humanlike sense of touch.

The technology was publicly demonstrated on the **Unitree G1 humanoid robot** in early 2025 and attracted significant media coverage in April 2026 after a LinkedIn post by Mike Kalil went viral.

## Founder & Team

- **Founder:** Zhang Jingqi — holds an electronic engineering degree from the University of York; previously served as CTO at a Xiaomi ecosystem company
- **Team size:** ~30 people, including alumni from the Chinese Academy of Sciences and top hardware firms

## Technology

### Core Innovation: 4th-Generation Resistive Flexible Fiber Sensor

Unlike conventional electronic skin that uses thin-film pressure sensors (fragile, poor elasticity, limited breathability, hard to mass produce), JQ's architecture uses **fiber-based layers** within flexible yarn:

- **Structure:** Fiber nylon core + piezoresistive outer layer + electrode core — all inside a single flexible yarn
- **Fold radius:** < 0.1 mm (extreme flexibility)
- **Durability:** Passes automotive-grade impact tests; high elasticity and rebound without recalibration
- **Environmental resistance:** Strong resistance to interference; breathable and skin-comfortable when woven into textiles

### E-Skin Suite (Unitree G1 Demo — "G1 Skin Suite")

The G1 integration includes three components:
- **Gloves** — multi-point tactile feedback across the palm; supports motion capture and remote manipulation; accommodates joint folds to zero degrees; detects forces as low as **0.01 Newtons**; pressure sensitivity range up to **3,500 kPa**
- **Torso skin** — full-body distributed sensing
- **Soles** — foot-embedded sensors that detect surface hardness and balance for stable walking

**Simulation integration:** Compatible with NVIDIA Isaac Sim and MuJoCo for tactile data modeling and closed-loop AI training.

### Sensing Capability

- Detects gram-level pressure in real-time
- Maps pressure, contact, and force distribution across thousands of sensing points simultaneously
- Transmits data from thousands of points to the robot's AI for instant decision-making
- Full-body pressure sensing matrix visualization (real-time 3D heat maps)

### Manufacturing Advantage

JQ built its own factory, retrofitted textile machinery, and optimized dyeing/finishing processes:
- Production costs cut by **80%**
- Yield consistency improved by **200-300%**
- Industrial roll-to-roll line capacity: **200,000 sensors/year** at **95% yield**
- New sample prototyping: **within 24 hours**

Production method: Modified industrial weaving + roll-to-roll production (textile manufacturing techniques — not semiconductor fab)

## Funding

| Round | Amount | Date | Investors |
|-------|--------|------|-----------|
| Seed | ~$2.8M USD equivalent | 2024 | Undisclosed |
| Series A | Eight-figure RMB (≈ $1.4M–$14M USD) | May 2025 | JHF, New Oriental, Freyr Graphite (financial advisor: InvesTarget) |
| Series A+ | Underway as of May 2025 | 2025 | — |

**Total raised:** At least $2.8M USD (likely more post-Series A)

## Market & Applications

### Current
- **Humanoid robotics** — primary focus; negotiating with 10+ humanoid robotics companies; some units shipping in volume
- **Eldercare** — smart mattress that monitors pressure, posture, sleep quality, breathing (partners: Tian'e Daojia, Taikang, Yuecheng Senior Living)
- **Automotive** — cockpit comfort sensors, gamified car seat experiences entering mass production with major OEMs

### Product Portfolio
- Robotic skins (torso, limb, full-body)
- Fingertip sensors / sensor gloves
- Software platforms for algorithmic processing
- Smart eldercare monitoring systems
- In-car comfort monitoring

## Key Quotes

> "We're now heavily focused on humanoid robotics and scaling fiber sensors into large-area electronic skin. These robotic skins can help robots sense their environment more holistically, interact like humans, and ensure safer human-robot interaction. For example, foot-embedded sensors can detect surface hardness and balance to help robots walk more stably."
> — **Zhang Jingqi, Founder**, to 36Kr (May 2025)

> "JQ Industries delivers a compelling solution in embodied haptics. True full-body robotic skins require scalable, cost-effective, and reliable solutions. JQ's fiber-textile approach enables a rare combination of flexibility, sensitivity, accuracy, range, durability, and affordability."
> — **Zhou Lifeng, Partner at JHF** (lead Series A investor)

## OneMind Relevance

- **Unitree integration** — directly relevant to the OneMind robotics stack (see [[unitree-toborlife-partnership-research]] and [[unitree-openmind-integration]])
- **Embodied AI** — e-skin closes the tactile sensing gap; enables closed-loop learning where robots learn from touch the same way humans do
- **NVIDIA Isaac Sim compatibility** — aligned with simulation-to-real training pipelines
- **Domestic Chinese supplier advantage** — scales with the humanoid robotics supply chain ZeusDelacruz has been tracking

## Competitive Landscape

| Company | Approach | Weakness vs JQ |
|---------|----------|----------------|
| SynTouch (US) | Rigid BioTac fingertip sensors | Limited coverage area, expensive |
| Peratech (UK) | QTC thin-film pressure sensors | Fragility, elasticity limits |
| Touchence (JP) | Tactile film sensors | Durability, breathability |
| MIT/Stanford labs | Soft robotics research e-skin | Not commercially scaled |

**JQ's moat:** Textile manufacturing techniques (not semiconductor fab) → 80% lower cost, 95% yield, 200K units/year capacity. Plant-derived materials = sustainable + comfortable.

## Sources

- KrASIA (36Kr English): https://kr-asia.com/with-fresh-funding-jq-industries-bets-on-fiber-sensors-to-give-machines-a-sense-of-touch (published May 23, 2025)
- Mike Kalil Blog: https://mikekalil.com/blog/textile-robot-e-skin/ (published Apr 5, 2026)
- YouTube: Kalil 4.0 — "This Robot's Skin Can Actually Feel Touch" (1 month ago, 2.3K+ views)
- X/Twitter: @XRoboHub — coverage of Unitree G1 e-skin demo

## Related Notes

- [[unitree-toborlife-partnership-research]]
- [[unitree-openmind-integration]]
- [[Robotics]]
