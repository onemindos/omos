# Pillars and domains

The organising idea behind the overseers: Zeus's whole life and operation is divided into four **pillars**, each owning a numbered range of **domains** (00 to 99). Every request belongs to one pillar. Legacy is the front door and routes work to the right one.

This is recorded from the May 2026 design (the `omos-core` era, when five domain agents ran on GitHub Copilot models). Those agents no longer run; their pillars live on as the four overseers (SO, HP, LE, GE). The original files are in the vault archive. Whether the numeric ranges and routing rules below still apply as written is a decision for Zeus; they are written down here so the idea is not lost.

## The four pillars

| Pillar | Domains | Covers | May 2026 domain agent | Today's overseer |
|--------|---------|--------|----------------------|------------------|
| **SO** Sovereign Operations | 00–24 | Infrastructure, VPS, Docker, security, credentials, cybersecurity, automation, n8n, GitHub, CI/CD, development, research, drones, robotics, compliance and legal deadlines, the operating layer of OneMind itself | Sovereignty | `agents/so` |
| **HP** Holistic Performance | 25–49 | Body, fitness, nutrition, sleep, recovery, mindset, habits, learning, personal finance, investments, purpose, identity, intelligence | Zenith | `agents/hp` |
| **LE** Legacy Evolution | 50–74 | Family, marriage, parenting, elder care, home, land, homestead, off-grid systems, heritage, estate planning, community | Unity | `agents/le` |
| **GE** Generational Empire | 75–99 | Business ventures, brand, products, marketing, capital, content, community, partnerships, Legacy Productions, investor relations, compliance, generational wealth | Eternal | `agents/ge` |

HP is Zeus the person. LE is the home and land as one living system. GE works on the longest time horizon: decades and generations, not quarters.

## Specialist skills each pillar loads on demand

In the May design a domain agent loaded specialist skills for deep expertise. Several of those names became today's planned servants (`agents/planned/`).

| Pillar | Skills |
|--------|--------|
| SO | grid, guardian, overwatch, eagleeye |
| HP | spartan, oracle, identity, librarian, pulse (genesis also tagged HP) |
| LE | haven, helix, eden, nomad |
| GE | empire, creed, trinity, apex, finance |

A separate watchdog, Sentinel, ran on its own node outside the core to monitor it and handle failover. It belonged to no pillar.

## Routing rules

1. A domain agent acts only on tasks inside its own domain range.
2. Anything outside the range is routed back to Legacy ("This is outside GE. Routing to Legacy.").
3. Agents never contact another pillar's agent directly. All cross-domain work goes through Legacy.
4. Decisions that need Zeus (the highest-stakes levels) escalate to him through Legacy.
5. SO coordinates with the watchdog for monitoring and failover scenarios.
