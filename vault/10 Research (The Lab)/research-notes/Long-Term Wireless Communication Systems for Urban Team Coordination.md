---
id: ui-15-20260111053926
title: "Long Term Wireless Communication Systems For Urban Team Coordination"
type: concept
domain: SO
domain: 00-24
folder: "00-24 SO (Sovereign Operations)"
status: active
created: 2026-01-11
modified: 2026-01-11
owner: zeus
tags: [domain/so, folder/15]
todoist_project: "00-24 SO (Sovereign Operations)"
ai_summary: "Long Term Wireless Communication Systems For Urban Team Coordination"
---

Long-Term Wireless Communication Systems for
Urban Team Coordination
Introduction
Coordinating 20–50 people in a dense city environment requires a robust wireless communication system.
Key considerations include coverage range, reliability amid urban interference, support for one-way
broadcasts and two-way conversations, and long-term usability (battery life, easy maintenance). Below
we explore two broad classes of solutions – cellular network-based systems and independent radio
systems – and evaluate options like Push-to-Talk over Cellular (PoC), Digital Mobile Radio (DMR), mesh
networks, and wireless intercoms. All these systems can support a group of 50 or more users, with secure
communications and compatibility with earpiece headsets. We compare their ranges, infrastructure needs,
and pros/cons, then recommend top picks for scalable, reliable city-wide operation.


Cellular Network-Based Solutions (LTE/5G Push-to-Talk)

Push-to-Talk over Cellular (PoC)

Push-to-Talk over Cellular (PoC) systems use existing 4G/LTE/5G cellular networks (or Wi-Fi) to transmit
two-way radio-style voice communications. This means that as long as users have cellular coverage, they
can communicate at any distance – even across an entire city or country 1 . PoC devices or apps turn a
phone or specialized radio into a push-to-talk communicator, enabling instant one-to-many group calls at the
press of a button, similar to traditional walkie-talkies 2 3 . This technology essentially provides
nationwide radio service with low latency and high capacity by leveraging carrier networks 1 . Modern
PoC platforms often support not just voice, but also integrated features like GPS location tracking, text/
image sharing, and even video push-to-talk for situational awareness 3 .


One-way vs Two-way: PoC easily handles both modes – the user can broadcast instructions to the entire
group (one-way), and any team member can press their PTT button to respond or speak to the group (two-
way) as needed. Communication is typically half-duplex (one speaker at a time) by default, though some
PoC apps also support full-duplex voice calls when needed. The group size is essentially unlimited; dozens
or even hundreds of users can be in the same talk group since the server and network manage the voice
data distribution.


Range and Coverage: There is effectively no range limitation between devices – coverage is determined by
the cellular network. If a team is spread out over multiple city blocks or across town, PoC is ideal because
anyone with a signal (LTE/5G or Wi-Fi) can hear and transmit. This provides a far greater coverage area than
line-of-sight radios. For example, a PoC radio can allow communication “anywhere with 3G/4G/WLAN
coverage”, whether across a building or across the country 4 . In practice, urban cellular networks provide
near-continuous coverage, though underground areas or concrete structures might require using Wi-Fi as a
fallback (many PoC devices seamlessly switch to Wi-Fi if cellular signals are weak 5 ).




                                                     1
Reliability in Urban Environments: Because PoC relies on carrier infrastructure, it benefits from the
robustness of commercial networks. Interference is mitigated by the cellular network’s design – users don’t
have to worry about radio frequency coordination, since the carrier manages spectrum and quality of
service. LTE/5G networks are designed for high user capacity even in dense cities, and mission-critical PoC
services (MCPTT) can prioritize emergency group calls with QoS guarantees 6 . The trade-off is
dependency on the network – if the cellular network is congested or goes down (e.g. in a disaster), PoC
units would be affected 7 8 . That said, modern PoC is quite reliable for business use: it offers improved
security and efficiency for operations without requiring the user to build their own radio infrastructure 9 .
Latency is low (often under a second) for push-to-talk calls, though not quite as instantaneous as direct
radios in all cases. In summary, PoC provides wide-area coverage and scalability at the cost of being tied
to carrier service.


Devices and Earpiece Compatibility: PoC solutions can run on rugged smartphones, dedicated PoC
handhelds, or even adapt existing smartphones via apps. Many vendors offer purpose-built PoC radios
that look like traditional two-way radios (with a side PTT button and loudspeaker) but contain SIM cards and
LTE modems – for example, the Hytera PNC370 is a handheld PoC radio that operates over 3G/4G/Wi-Fi
networks 4 . These devices typically have standard audio jacks or Bluetooth, making them compatible with
covert acoustic tube earpieces, headsets, or earbud/mic combinations. In fact, Bluetooth earpiece support
is often built-in for hands-free use 10 , so the team can comfortably wear earpieces and receive commands.
Setting up PoC for the team may be as simple as distributing devices with the PoC app and configuring
them to the same talk group. No external antenna towers or repeaters are needed on-site – an internet
connection (via cellular or Wi-Fi) is the only requirement.


Battery Life: PoC devices must power cellular radios and data transmission, so battery life is a
consideration for long-term use. High-capacity batteries are common in dedicated PoC radios – for
example, a 3100 mAh battery in the PNC370 yields about 18 hours of operation on a charge 11 . In
practice, users report that PoC radios have slightly shorter battery endurance than analog radios if used
heavily 12 , but they can still last a full work shift or longer. Many devices charge via USB or docking stations
for convenience 11 . For multi-day operations, providing power banks or spare batteries is advisable.
Overall, with proper devices, PoC can support all-day communication without frequent recharging.


Security: Communications over PoC are generally very secure. Voice is transmitted as encrypted data over
the carrier network. Enterprise PoC platforms use advanced encryption and authentication to protect
calls 13 , and since traffic goes through a central server, administrators can enforce user permissions and
group membership. This means eavesdropping is virtually impossible for an outside party (unlike analog
radio scanners). However, since it is IP-based, PoC could theoretically be vulnerable to cyber threats if not
properly secured 14 – using reputable PoC services or private servers mitigates this. Overall, PoC offers
enterprise-grade security out of the box (and even meets public safety encryption standards in MCPTT
implementations).


Infrastructure and Maintenance: PoC’s big advantage is minimal infrastructure for the user. There’s no
need to deploy radio repeaters or obtain FCC licenses – you leverage the cellular network (typically a
subscription or data plan per device is required). Some PoC services are subscription-based (with monthly
fees), while others allow you to host your own PoC server for full control. Either way, initial setup is fast –
devices can be configured and in use within hours 2 15 . Maintenance is largely about managing the
devices (keeping them charged and updated) and the service subscription. If using a carrier’s Push-to-Talk
service or a cloud PoC provider, you may have recurring costs, but you avoid the higher upfront cost of




                                                       2
physical radio infrastructure. This scalability and low barrier to entry make PoC very attractive: it’s easy to
add more users (just give them a device and login) and even interconnect with other systems (some PoC
services can bridge to conventional radio networks via gateways 16 17 ).


Summary (PoC): Push-to-Talk over Cellular is a highly scalable, wide-area solution ideal for coordinating a
dispersed team. It enables instant group voice communication across a city with minimal setup, using
either smartphones or specialized PoC radios. Strengths include virtually unlimited range, easy group
scaling (50+ users is no problem), and rich features like location tracking. Drawbacks include reliance on
cellular coverage and potential subscription costs, as well as slightly shorter device battery life compared to
some traditional radios 18 19 . For many urban operations with good carrier coverage, PoC is a top choice
for long-term team comms due to its flexibility and broad reach.


Independent Radio-Based Solutions
If a truly independent communication line is needed (for example, to avoid relying on public networks or
where infrastructure may be unavailable), several radio-based systems can support 50+ users with reliable
performance. These operate on private radio frequencies or dedicated wireless links. Below, we examine
Digital Mobile Radio (DMR), mesh network radios, and wireless intercom systems as options. All of
these can function without carrier networks – an important advantage for mission-critical or disaster
scenarios. They do require varying degrees of hardware setup (radios, repeaters, base stations), but offer
proven reliability, low latency, and control over your communication network.


Digital Mobile Radio (DMR)

Digital Mobile Radio (DMR) is a modern two-way radio standard widely used in public safety and industry.
DMR handheld radios look like traditional “walkie-talkies” and operate on UHF/VHF frequencies that you
control (usually via a license). They provide crystal-clear digital voice, group call capability, and
encryption support, making them a prime choice for secure team communication. DMR technology uses
TDMA (time-division multiple access) to allow two simultaneous voice channels on one frequency, doubling
capacity over analog systems 20 . For a team of 50, this means you could even have two separate talk
groups on one channel, or handle two calls at once, if needed. In practice, most teams would share one
talkgroup for broadcast-style comms, which DMR handles very well – one person speaks, all others instantly
hear, with sub-second push-to-talk latency.




                                                      3
High-power DMR radios can cover large areas even in dense buildings. Modern UHF DMR portables with strong
output (5–10 W) and good antennas can penetrate high-rise structures – for example, some are rated to cover up
to 700,000 sq.ft (roughly 65,000 m²) or about 50 floors of a building with a single radio, given its improved signal
penetration 21 .


Range and Coverage: In direct mode (radio-to-radio), handheld DMR units typically reach 3–5 km in urban
areas under real-world conditions 22 23 . This range is plenty to cover a large campus or event venue.
However, dense city blocks with buildings can limit range to a kilometer or two unless users are in relatively
open areas. To extend coverage, DMR can use repeaters or trunked networks. A single high-power DMR
repeater (often placed on a rooftop) can blanket a wide area – by one estimate, a properly sited repeater
can allow handheld coverage over “dozens of kilometers” or an entire city 22 24 . For example, a UHF DMR
repeater on a tall building could let all 50 team members communicate across town reliably. DMR also
supports multi-site networks: multiple repeaters linked via IP can cover different city regions, handing off
calls (this is DMR Tier III trunking). In short, DMR is very scalable in range – it can be as simple as a few
radios for local use, or scaled up with infrastructure to cover a metropolis. The trade-off is the upfront
setup: deploying repeaters and obtaining frequency licenses requires coordination and cost (which, given
no budget constraints, might be acceptable). Once set up, though, a DMR system gives you independent,
always-on communication with no reliance on external networks 25 .


Reliability and Interference: A well-configured DMR system offers high reliability and excellent audio
clarity in urban environments. Because you control the frequencies (usually licensed channels), there is
minimal interference from outsiders. Digital modulation and error correction in DMR also mean that even if
there is some radio frequency noise, users will hear clear audio or nothing (no static as in analog). In fact,
DMR is preferred for mission-critical use due to its immediate connectivity and independence from any
public network 26 27 . If one user transmits, all others in the group hear it nearly instantly with negligible
delay. Additionally, DMR devices are typically ruggedized for heavy use (many meet MIL-STD durability).
With dedicated frequencies and high-grade hardware, DMR communications remain solid even amid
urban RF congestion. In scenarios where cellular networks might become overloaded or fail, a DMR radio
network will keep running as long as your devices and repeaters have power. The only interference to watch
for is other radio users on nearby frequencies, but proper licensing and frequency coordination mitigate
that. In summary, DMR provides low-latency, interference-resistant comms ideal for dense city
conditions.




                                                         4
One-way vs Two-way: DMR natively supports one-to-many group calls. This works like a classic push-to-talk
– the team leader can press the PTT button and deliver a message to all 50 users’ earpieces simultaneously.
If operating in a strictly one-way mode, others can stay in listen-only mode (some radios can be configured
to disable PTT for certain users or simply instruct everyone that only the leader speaks). For two-way needs,
any team member can transmit in turn; DMR systems manage this by one user speaking at a time per
talkgroup. The transition is virtually seamless – when the leader releases PTT, another user can press theirs
to reply. There is no dial-up delay or network negotiation needed, which makes conversations snappy.
Although only one person is heard at a time (half-duplex), coordination is straightforward with push-to-talk
etiquette. If simultaneous separate conversations are needed, DMR’s dual time-slot nature allows multiple
channels. For example, you could split into two sub-teams of 25 each on different time-slot channels and still
only use one frequency. This flexibility in group organization is a plus for DMR when scaling to 50+ users.


Battery Life: DMR handhelds are known for long battery life, which is a boon for long-term operations.
Digital modulation is more power-efficient than analog – by using TDMA (only transmitting half the time in a
time slot), DMR radios can achieve up to ~40% longer run time than analog radios 28 . For instance, one
comparison showed ~13 hours talk time in digital mode vs 9 hours in analog on the same device 29 . High-
capacity batteries (e.g. 2000–3000 mAh) can often sustain 12–15+ hours of active use on DMR portables.
Under light usage (mostly standby listening with occasional talk), some radios boast up to 20–24 hours on a
charge 30 . In practice, a standard shift or full day event can be covered on one battery; for continuous
24/7 operations, batteries can be swapped or units recharged in docks. Thus, DMR is very suitable for long-
term deployments – battery endurance is generally not a limiting factor. (It’s wise to have spare battery
packs or chargers for all radios just in case.)


Security: Security is another strong point for DMR. Many commercial DMR radios support 128-bit or 256-
bit AES encryption for voice calls, preventing anyone from eavesdropping on the communications 31 .
Unlike analog “privacy codes” (which only filter chatter but don’t truly encrypt), DMR encryption actually
scrambles the digital audio such that only radios with the correct key can decode it. This means a team’s
conversations stay confidential – an important factor if coordinating sensitive operations. Even without
encryption, DMR uses digital encoding, so casual scanners won’t pick up intelligible audio. With encryption
enabled, DMR meets a high security standard (some models are used by law enforcement with strict
encryption requirements). The system is also closed: only the distributed radios on your frequencies can
join the call, which adds to privacy. Overall, a properly programmed DMR system offers secure
communications on par with military/law-enforcement grade radios for the team.


Device Compatibility and Earpieces: DMR radios come in handheld, mobile (vehicle), and even dispatcher
console forms, but for a team of people on foot, you’d issue handheld portable radios. These typically have
universal accessory connectors for audio. It is straightforward to attach surveillance earpiece kits,
speaker-microphones, or headsets to each radio so that users can hear and talk discreetly. For example, a
clear acoustic tube earpiece with a lapel microphone is commonly used in security teams – this plugs into
the radio’s accessory port. Many DMR radios also now support wireless Bluetooth earpieces, allowing
users to go truly wire-free. Earpiece audio will be as clear as the radio’s speaker (DMR digital voice has very
good clarity). The leader can thus transmit commands, and everyone hears it in their earpiece even in noisy
city streets. Users can also transmit by using their headset’s inline PTT button or the radio’s button. In terms
of device management, each radio needs to be programmed with the group’s frequency and settings
beforehand, but once set up, operation for the end-user is as simple as push-to-talk. Distributing 50 DMR
radios might involve a bit of initial programming and labeling, but there are services and software to clone
configurations easily. Maintenance over time would include keeping batteries charged, and occasional




                                                       5
radio firmware/programming updates or license renewals. These radios are built for daily use over years
(often rugged and weather-proof, IP-rated for dust/water), so they are well-suited to long-term projects.


Summary (DMR): Digital Mobile Radio is a top choice for independent, secure group communication in
an urban setting. With a dedicated DMR system, your team gets instant, reliable push-to-talk
communication that you fully control. Pros: extremely reliable low-latency audio, works off-network, good
building penetration and range (expandable with repeaters), strong security (encryption), long battery life,
and compatibility with a variety of earpiece/headset accessories. It’s proven to scale to large groups – 50
users can easily share a channel, and more can be accommodated with additional channels or talkgroups.
Cons: Upfront equipment costs and the need for FCC frequency licensing and repeater setup (which can be
significant, though a “no budget constraint” scenario lessens this issue) 32 . Also, without a repeater,
handheld-to-handheld range is limited to a few kilometers, so additional infrastructure may be needed to
cover an entire city (unlike PoC which uses existing towers). Overall, for an organization that values
independence and reliability, a DMR radio network is a highly scalable solution delivering secure comms for
50+ people in a city environment.


Mesh Network Radios (Ad-Hoc Mesh Communication)

Mesh network radios represent an advanced solution where each device not only transmits and receives,
but also acts as a node that relays data for others. In a mesh, units automatically form an ad-hoc network
– if two team members are out of direct range or blocked by buildings, intermediate members’ radios will
pass along the message. This creates a self-healing, robust web of connectivity 33 34 . For a group of 20–
50, a mesh system can be very powerful: as long as there are enough nodes in between, the group’s
coverage can extend well beyond the range of a single radio and can circumvent urban obstacles.


How it Works: Each mesh radio is essentially a smart transceiver (often using MANET – Mobile Ad hoc
Network – protocols). When a user talks, the voice is digitized and sent out; whichever radios are in range
will automatically forward that data to others, hopping until all nodes in the network receive it. This
happens in real-time and dynamically – routes are found on the fly, and if one radio drops out or moves, the
network reroutes around it 35 33 . There is no central base station required (infrastructureless operation).
This is ideal for long-term deployments where the team may spread out or move through a city: the
network coverage expands with the team. For example, if some members go down one street and others
are a few blocks away, the radios between them act as repeaters to maintain connectivity.


Range and Coverage: Individually, a mesh radio node has a certain transmit range (similar to other radios).
Many mesh devices use relatively high power to maximize reach – e.g. the Defensor-TS1 handheld mesh
PTT radio can output 2–25 W, adjustable 36 . At maximum power, line-of-sight range of a couple kilometers
per hop in urban terrain is achievable. The real advantage is multi-hop: if each hop covers 1–2 km in the
city, 5 or 6 hops can cover 5–10+ km, effectively linking a large area 36 37 . In one use-case, a single mesh
node on the 1st floor was able to cover from the basement (-2F) up to the 80th floor of a skyscraper by
relaying through mesh nodes placed at intervals 37 – demonstrating how well a mesh can penetrate
vertical structures. These systems often support a fixed maximum number of hops (e.g. 6 hops is common)
before the signal needs a fresh source, but that is usually enough to cover an expansive zone when nodes
are spaced out. Additionally, the more nodes active, the more alternate paths exist. In practical terms, a
mesh of 50 radios scattered around a city sector could create a wide coverage blanket, especially if a few
high-powered nodes are strategically placed. It’s not unusual for advanced mesh networks (like military
MANET radios) to link dozens or even hundreds of nodes across a city with no fixed towers 38 . Thus, for




                                                     6
scalability, mesh networks shine in environments where the team might not be all in one place – they extend
coverage dynamically as users move.


Reliability and Interference: Mesh networks are built for robustness. The self-healing nature means if
one radio fails or leaves the area, others can still communicate by finding new routes. In a noisy urban RF
environment, mesh radios often use sophisticated waveforms (e.g. COFDM, frequency hopping, MIMO
antennas) to combat interference 39 40 . They may operate in bands less congested by common radios
(some mesh systems use 900 MHz, others 2.4 GHz or even L-band depending on design). Because they form
a closed network, they can also employ strong encryption – many tactical mesh radios have built-in AES
encryption similar to DMR, keeping communications secure. Overall, a well-designed mesh will be very
reliable for mission-critical use, as it was originally developed for military and emergency responder
scenarios where infrastructure is down. The main challenge in an urban setting is that very long multi-hop
paths can introduce slight latency or reduced voice quality if bandwidth is low. However, modern mesh
radios are optimized for voice. They also can adjust routes on the fly if one path becomes noisy. By having
multiple redundant paths, mesh systems are resistant to interference – even if one link is blocked by a
building, an alternate via a different neighbor might succeed. This redundancy greatly increases
communication reliability in a city’s unpredictable conditions.


One-way vs Two-way: Mesh radios support the same push-to-talk group functionality as conventional
radios – one user’s voice is broadcast to all. From the users’ perspective, it works like a normal group radio
call. So one-way leader broadcasts and two-way exchanges are supported similarly (half-duplex push-to-talk
voice). The mesh aspect is mostly “behind the scenes” in how the signal reaches everyone. Some mesh
systems might allow multiple simultaneous talkers if the network has enough channels, but typically the
limiting factor is human coordination (to avoid garble). So, like DMR, it’s usually one speaker at a time in a
given group. Most mesh radios also carry data, so they could allow text or telemetry concurrently without
interrupting voice. For team coordination, the voice group call is the primary use – all 50 can listen, and
anyone can respond by pressing PTT. Notably, because mesh nodes may need to transmit for others, the
network often is engineered to handle a lot of traffic. Everyone on the mesh can be in one talkgroup or
divided into subgroups logically. Scalability in users is generally high: many systems advertise handling
30+ active nodes easily, and some up to 100+. More nodes can actually strengthen the network (up to a
point) by providing more pathways. Therefore, a mesh network can comfortably accommodate 50 users
and is inherently designed to scale.


Battery Life: Mesh radios tend to be more power-hungry due to constant networking overhead and often
higher transmit power. However, high-end devices mitigate this with large battery capacity and smart power
management. For example, the Defensor-TS1 mesh radio has a 31-hour battery life on a single charge,
and can even be extended to 5 days (120 hours) with an external power bank attachment 41 . This is
exceptionally good – it indicates the device likely operates at lower duty cycle most of the time and uses
efficient routing. In active use (relaying others’ messages), battery drain will be higher, but even so, these
are designed for long field deployments. To ensure long-term operation, one can issue each user a portable
charger or spare battery as well. The bottom line: battery life can be sufficient for multi-day continuous
ops if planned properly, though it’s wise to have charging routines since mesh units do more work than
simple radios. Many mesh radios allow operation while charging (say, via USB or external battery) to
facilitate 24/7 uptime. Given no budget constraint, one could equip some mesh nodes at fixed locations
with continuous power (to act as dedicated relays) and have others portable. This hybrid approach can
guarantee the network stays up indefinitely.




                                                      7
Ease of Use and Maintenance: For the end-user, a mesh radio can be as easy to use as any walkie-talkie.
They don’t need to manage the network – it forms automatically. Users just need to power on their unit and
it will connect to the mesh. The complexity is mostly in initial setup and network planning. Frequency
selection and device configuration need to be done (often, these are pre-configured to all operate on the
same network ID or waveform). If the mesh radios come from a single vendor as a kit, they typically arrive
ready to mesh; otherwise, some technical integration is needed. Maintenance involves ensuring all units
have up-to-date firmware and possibly managing any encryption keys. Since each device is essentially a
mini base station, there’s a bit more to monitor (like making sure no device overheats or drains
unexpectedly from heavy routing). However, mesh systems are built for field robustness – they often
operate in harsh conditions with minimal user intervention. One consideration is training: the team should
be briefed that these radios will automatically relay for others, so they should remain powered on and
within range as much as possible to maintain the network. With 50 nodes, there is a lot of redundancy, so
losing a few won’t break communications, but leaving them all on ensures best performance. In terms of
deployment, no external infrastructure (towers) is required, but sometimes teams will deploy a few fixed
relay nodes (e.g. mount a couple in high locations or vehicles) to boost coverage. Compared to setting up
multiple repeaters for DMR, setting up a mesh is relatively quick – it’s mostly deploy and go, as the network
self-configures.


Summary (Mesh Radios): Mesh network radios offer a cutting-edge solution for team communication
with maximum resiliency and coverage. They are particularly well-suited for scenarios where the team may
disperse or move through an urban area, or where existing infrastructure is untrusted. Pros: No fixed
infrastructure needed (self-forming network), extended range via multi-hop, resilient to single-point
failures, secure encrypted links, and adequate throughput for voice and data. They can scale to 50+ nodes
with improving coverage as more nodes join. Cons: Higher complexity and device cost, potentially heavier
units, and dependency on having a sufficient density of radios to maintain hops. Also, while voice is near-
real-time, a long chain of hops can introduce slight delays or require careful network tuning. Nonetheless,
with proper deployment, mesh radios provide ultimate reliability – even if cellular towers fail and
conventional repeaters go down, the mesh keeps the team connected. This makes it a top choice for long-
term operations where self-sufficiency and coverage flexibility are paramount (e.g. disaster response
teams, large-scale events, or security operations spread over a city).


Wireless Intercom Systems (Full-Duplex)

A wireless intercom system is another viable solution, especially if your team operates in a confined area
(like a building, campus, or event venue) and you value full-duplex communication (i.e. everyone can talk
and listen simultaneously, like a phone call). These systems are commonly used in film production crews,
event staff, and security details for instantaneous, hands-free chatter among many people. Wireless
intercoms typically use the DECT 6.0 standard (1.9 GHz frequency band) or similar, providing a secure, low-
latency voice link. They consist of beltpack transceivers with headsets (earpiece + mic) for each user and
one or more base stations that manage the conference.


One-way vs Two-way: Unlike push-to-talk radios, intercom beltpacks are often open-mic or push-to-talk in
a full-duplex audio channel. This means multiple people can speak at the same time and be heard by all,
much like a conference call. In practice, users still take turns to avoid chaos, but the system doesn’t hard-
enforce one speaker at a time. If desired, you can certainly operate it in a one-way mode by designating
only the leader to speak; however, the hardware itself allows others to chime in without a PTT button. Many
intercoms have a “push-to-transmit” button that users press when they want to talk (to mute when not




                                                     8
talking), so it can function similarly to PTT radios but with the benefit that others don’t get cut off if someone
presses at the same time. The audio is continuous and real-time, so it’s great for scenarios where quick back-
and-forth or interruptions are needed (“conversation” style communication). For purely broadcast use (one-
way), the leader could just keep their microphone live and everyone else stays muted. The flexibility is there
to have simultaneous two-way conversations among the group with no perceived delay.


Range and Coverage: DECT-based wireless intercoms have a moderate range – typically about 100–
300 meters from the base station in open conditions. In indoor urban environments, a single base can
cover a building floor or two, or a small outdoor area (100m is roughly 330 feet). This is shorter-range than
radios, but intercom systems can be extended by using multiple antennas or repeaters. For example, Clear-
Com’s FreeSpeak II system allows multiple distributed transceiver antennas connected to the base station
to widen coverage; users can roam between cells of these antennas. Using several antennas, you could
cover an entire stadium or building complex. Still, wireless intercoms are best suited when the team
operates within a contained zone – unlike PoC or DMR, they are not meant for kilometers of range away
from the base. If the user’s operation is within, say, a large building or an event footprint, intercoms can
work very well. The advantage in range is that you get coverage where you need it with high voice quality
and not much concern about external interference (DECT is license-free but relatively isolated in its own
band).


Capacity (Scaling to 50+ users): Professional wireless intercom systems can support a surprising number
of users by deploying multiple base stations or a matrix. A single base station often supports on the order
of 15–25 beltpacks. For instance, one popular system’s base can handle “up to 20 wireless beltpacks” in a
group 42 , and another claims “maximum 25 beltpacks in the same RF space on one base” 43 . To scale beyond
that, some systems allow linking of two bases or using higher-capacity modes. Newer intercom solutions
(like Clear-Com FreeSpeak Edge in 5 GHz) can achieve up to 50–100 wireless beltpacks on a unified
system by leveraging multiple transceivers and frequency bands 44 . In one configuration, 16 beltpacks per
base is standard, but with a matrix connecting multiple base stations, “up to 100 beltpacks [can be]
operational over a single system” 44 . So, it is feasible to have 50 people on the same conversation. They
would all hear each other and can speak as needed. In practice, having dozens of open mics can get
unwieldy, so often only a handful talk frequently while others mostly listen. But the system capacity will not
be a limiting factor with the right hardware. Thus, scalability to 50+ users is achievable by high-end
wireless intercom setups (often used in broadcast and large events). Keep in mind that as you add more
users and more area, you may need additional infrastructure (antennas) to maintain coverage and call
quality.


Audio Quality and Interference: Wireless intercoms usually provide clear, full-duplex audio comparable
to telephone quality. DECT systems are engineered to reject interference – DECT 6.0 uses automatic channel
hopping and encryption so that even in a busy RF environment, the audio stays stable 45 . It also operates
in a band (around 1.9 GHz) that is typically reserved for such communications, so it avoids interference from
Wi-Fi (2.4/5 GHz) or two-way radio bands. In a dense city building, there might be other DECT systems (e.g.
office cordless phones), but these use dynamic channel selection to minimize conflict. The result is a secure
and stable transmission that outsiders cannot eavesdrop on 45 . In fact, DECT uses built-in encryption to
prevent anyone from listening in on the intercom channel 45 . Even newer products explicitly advertise that
the DECT intercom “prevents people outside the system from hearing your conversation with encryption
technology” 45 . So security and interference are well-handled. The intercom base manages time slots for
each beltpack, and the system can “double up” users on channels if needed with TDMA (some advertise
double capacity modes) 46 . Bottom line: audio is real-time and reliable, with no push-to-talk delay, which




                                                        9
can enhance team responsiveness. The only potential issue in a city is if you try to cover too large an area
with too few antennas, users at the edge might drop out – but that’s solved by adding more coverage
points.


Earpiece/Headset Use: By design, wireless intercom beltpacks are meant to be used with headsets or
earpieces. Typically, each belt unit has a headset jack (like 4-pin XLR or 3.5mm) and is paired with a
lightweight headset (single-ear or dual-ear) that includes a microphone boom. This gives excellent audio in
loud environments and leaves the user’s hands free. There are also in-ear earpiece options for intercoms,
and even Bluetooth headset support on some packs 47 . For covert or low-profile needs, small in-ear
monitors can be used instead of the standard headset. The leader of the team might choose a heavier-duty
headset with a noise-cancelling mic if they’re in a noisy command post, while others can use simple
earpieces. The key benefit is comfort for long wear and the ability to hear and speak without holding a
radio. Distribution is simple: hand out the beltpacks and headsets, and each person turns on their pack and
is live on the network. Because everything is wireless, there are no cables tethering users to a base, which
is great for mobility in the field.


Battery Life: Wireless intercom beltpacks use rechargeable Li-ion batteries that typically last 8–12 hours
per charge, depending on usage. For example, a common beltpack model provides about “18+ hours of
battery life” with a high-capacity battery 48 , while another yields around 12 hours with its battery 49 50 .
This is usually sufficient for a full day event or shift. Users can recharge packs during breaks or swap
batteries if needed. Many systems come with multi-unit charging cradles so you can charge all packs
overnight. In continuous operation scenarios, you might have two sets of batteries per pack (one charging,
one in use). The batteries are often quick to charge (e.g. 3-4 hours for full charge). So for long-term
operations, power cycling the devices is manageable. It’s also worth noting that because intercom is full-
duplex, the transmitter is active whenever someone is talking (which can be frequent), so battery drain can
be steady. Still, the numbers above show that double-digit hours of talk time are standard, which covers
most needs. With planned charging, a wireless intercom setup can run indefinitely day after day.


Infrastructure and Maintenance: Setting up a wireless intercom for 50 users is more involved than a
simple two-way radio deployment but quite feasible. You would need to install one or more base stations
at the site. These bases need AC power (and possibly a network connection if linking multiple or integrating
with other systems). They act as the hub that all beltpacks connect to via DECT radio. For larger coverage, as
mentioned, you’d set up several remote antennas/transceivers connected to the base via cables or IP
network. For example, you might place a transceiver on each floor of a building or around a venue to
ensure seamless coverage. Once installed, though, the system requires little intervention – it will manage
the wireless connections automatically as users move around. Maintenance includes monitoring battery
health, updating firmware occasionally, and ensuring the base station is functioning (it’s wise to have a
backup power or UPS for the base station to survive any power outages, since if the base loses power, the
wireless beltpacks can’t communicate with each other except in systems that support a fallback direct
mode). Compared to a mesh or DMR with multiple repeaters, the intercom might actually be simpler: a
single base can cover the area with a few antennas, and there’s no need for frequency licenses or
programming channels. Everything is plug-and-play after initial setup. The ease of operation for users is
arguably highest with wireless intercom – they essentially just talk as if on a group phone call, with no PTT
needed (unless using the mute/talk button for etiquette). Training is minimal: it’s very intuitive (“wear this
earpiece and talk to the group as needed”).




                                                     10
Summary (Wireless Intercom): A wireless DECT intercom system is an excellent solution for local, high-
coordination team communications. It provides a continuous, full-duplex audio net for the team, which
can improve reaction times and allow a more natural flow of communication. Strengths: real-time
simultaneous conversation, secure encrypted channel, easy to use (hands-free), and good audio quality. It
can support large teams (with the right system, 50 or more users can be in the same conference) and
typically has sufficient range for a building or event site, with options to extend coverage. It’s also license-
free and self-contained (no reliance on telecom networks). Weaknesses: limited range per base station (not
ideal for an entire city unless you deploy multiple systems), and it requires a central base infrastructure that
must remain within range of all users. If team members spread out over many kilometers, this is not the
right choice without deploying several linked bases. Another consideration is cost: professional intercom
systems for 50 users can be expensive, but given the lack of budget constraints, one could invest in a high-
end system (like Clear-Com or similar) to achieve the goal. For scenarios like event production, campus
security, or anywhere the team stays within a relatively defined area, wireless intercoms offer unparalleled
ease of operation and group talk capability. They are a top pick when constant open communication is
needed and the coverage area is manageable.


A Note on Traditional Analog Options

For completeness, it’s worth noting simpler analog radio options (like FRS/GMRS walkie-talkies or analog
wireless intercoms). These license-free radios are inexpensive and easy to deploy, but they come with
major drawbacks in a city environment. FRS/GMRS handhelds have very limited power (0.5 W for FRS, a
couple watts for GMRS) and suffer heavy interference in urban areas – in fact, they “have little or no
interference so they are not reliable for use in a city setting” according to CERT communications guidance 51 .
Their range is often much less than advertised once buildings and obstructions are present (often just a few
blocks at best). They also lack any encryption (anyone with a similar radio can listen in), and with only 22
channels, a busy city can have many users on the same channels causing crosstalk. While FRS/GMRS could
technically equip 50 people cheaply, the lack of security, limited scalability, and poor urban reliability
make them a poor choice for a serious long-term operation – especially since budget is not a concern here.
Similarly, older analog wireless intercoms (or simple FM transmitters) for tour-guide style one-way
communication might be used in niche cases (e.g. a tour guide broadcasting to 50 listeners with FM
receivers), but these too can be subject to interference and offer no two-way capability. Therefore, the focus
remains on the digital, high-reliability systems above (PoC, DMR, mesh, DECT), which are far better suited
to the user’s needs.


Comparison of Key Features
The table below summarizes key features of the discussed systems, comparing their typical range,
communication mode, network requirements, device compatibility, and high-level pros/cons:




                                                      11
                         Typical                            Network          Devices &
System                                  Communication                                       Pros                Cons
                         Range                              Used             Earpieces

                                                                                                                • Net
                                                                                            • Wide-area         depe
                                                                                            coverage (not       (cellu
                                                                                            distance-limited)   or co
                                                                             Smartphones      1 .<br>•          can d
                         Unlimited
                                        Half-duplex PTT                      or dedicated   Highly scalable     comm
                         (within                            4G/5G cellular
                                        (one-at-a-time                       PoC radios.    (50+ users             7 .
                         cellular/Wi-                       data or Wi-Fi
Push-to-Talk over                       speaking),                           Standard       easily) 18 .<br>•   Subs
                         Fi                                 (requires
Cellular<br>(LTE/5G                     supports one-                        wired or       Low                 data
                         coverage)                          carrier
PoC)                                    way broadcast                        Bluetooth      infrastructure      devic
                            1 . City-                       network or
                                        or group two-                        earpieces      setup (uses         Sligh
                         wide and                           internet) 52 .
                                        way.                                 supported      existing            comp
                         beyond.
                                                                              10 .          networks).<br>•     direc
                                                                                            Advanced            Devi
                                                                                            features (GPS,      life s
                                                                                            text, recording).   unde
                                                                                                                  12 .


                                                                                            • Very reliable,    • Lim
                                                                                            instant             with
                                                                                            communication       infra
                         3–5 km                                              Handheld       (independent of     (repe
                         handheld-                          Private UHF/     two-way        public networks)    need
                         to-                                VHF radio        radios           26 .<br>• Clear   city)
                         handheld                           frequencies.     (rugged).      digital audio &     Requ
                                        Half-duplex PTT.
Digital Mobile Radio     in urban                           Repeaters        Supports       encryption          frequ
                                        Group calls (one-
(DMR)<br>(Conventional   terrain 22 ;                       (optional) for   wired          support             licen
                                        way or two-way
or Trunked)              extended                           wide coverage.   earpieces,       31 .<br>• Good    setup
                                        sequenced).
                         to city-wide                       No external      speaker-       battery life (12+   High
                         with                               network          mics, or       hrs) 29 .<br>•      cost
                         repeaters                          needed 54 .      Bluetooth      Scales to large     repe
                           53  24 .                                          headsets.      groups; proven      One-
                                                                                            tech for            talki
                                                                                            security/public     simu
                                                                                            safety.             as in




                                               12
                             Typical                            Network            Devices &
 System                                     Communication                                           Pros                Cons
                             Range                              Used               Earpieces

                                                                                                    •
                                                                                                    Infrastructure-     • Com
                                                                                                    free – network      cost
                                                                                                    forms               (adva
                             ~1–2 km
                                                                                                    automatically       techn
                             per hop in
                                                                                   Specialized        55 .<br>•         Mus
                             city (node-                        Self-forming
                                                                                   mesh radio       Extended            suffi
                             to-node);                          mesh network
                                                                                   units            coverage via        dens
                             multi-hop      Half-duplex PTT     (no fixed
                                                                                   (handheld or     multi-hop           depe
                             links up to    voice (group        infrastructure).
 Mesh Network                                                                      wearable).       (network relays     inter
                             5–10+ km       broadcast via       Uses
 Radios<br>(Ad-hoc                                                                 Earpieces or     around              units
                             with 6 hops    mesh).              unlicensed or
 MANET)                                                                            headsets         buildings)          Shor
                               36  37 .     Concurrent data     dedicated
                                                                                   attach             37 .<br>•         per h
                             Can cover      capabilities.       band
                                                                                   similarly to     Resilient & self-   high
                             large areas                        depending on
                                                                                   two-way          healing (nodes      sing
                             if nodes                           system.
                                                                                   radios.          route around        (nee
                             are
                                                                                                    failures)           relay
                             distributed.
                                                                                                      35 .<br>•         Pote
                                                                                                    Secure, high-       laten
                                                                                                    end encryption      over
                                                                                                    (often AES).

                                                                                                                        • Lim
                                                                                                                        – not
                             ~100–                                                                  • Real-time,        area
                             300 m from                                                             overlapping         cove
                             base                                                                   conversation        base
                             (indoors                           Dedicated          Beltpack         (no PTT             Requ
                                            Full-duplex voice
                             ~floor/                            wireless           transceivers     needed).<br>•       base
                                            (multiple can
                             building                           intercom           + headsets/      Encrypted and       and l
                                            talk and listen
 Wireless                    coverage).                         system on          earpieces        secure (DECT        ish
                                            simultaneously).
 Intercom<br>(DECT-          Expandable                         1.9 GHz DECT         45 . Usually   6.0 encryption)     cove
                                            One-way
 based full-duplex)          with                               (no external       worn on            45 .<br>• Easy    Capa
                                            broadcast
                             multiple                           network; base      body with an     to use, hands-      base
                                            possible by
                             base                               station on-        earpiece and     free                (need
                                            muting others.
                             antennas                           site).             mic.             operation.<br>•     base
                             (roaming                                                               High audio           42

                             across                                                                 quality, low        Base
                             cells).                                                                latency.            sing
                                                                                                                        failu
                                                                                                                        back

(Table References: Range and coverage estimates from vendor data and practical experience 22 36 .
Communication modes and network needs from technology standards 1 54 . Device compatibility and features
from product specifications 10 31 . Pros/cons summarized from comparative analyses 18 56 .)




                                                   13
Recommendations
Considering scalability, dense-city reliability, and ease of operation, here are the top recommended
communication systems for the user’s scenario:


     • Primary Recommendation – Digital Mobile Radio (DMR) System with Repeaters: For an
       independent, long-term solution that can reliably cover a city environment, a DMR-based radio
       network is a top choice. DMR strikes an excellent balance of range, capacity, and robustness. With
       one or two strategically placed repeaters, 50 handheld DMR radios can maintain city-wide coverage,
       including inside buildings (UHF frequencies penetrate urban structures well, especially with a strong
       repeater signal) 53 24 . Each team member gets a durable radio with an earpiece, enabling clear
       reception of the user’s commands. The system supports one-to-many broadcasts perfectly, and any
       team member can respond via push-to-talk. Scalability: 50+ users are easily handled on one
       channel, and the system could be expanded with additional channels if needed. Reliability: Because
       it doesn’t rely on public networks, it will keep working during emergencies or cell network outages
         25 . Interference is minimal on a licensed channel, and digital encoding ensures clear audio even in

       RF-dense areas. Security: With AES encryption enabled, communications will be secure from
       eavesdropping 31 . Maintenance: Once set up, a DMR system is straightforward to maintain (mostly
       charging batteries and occasional equipment checks). Given the long battery life and proven
       longevity of DMR radios, this solution is very suited for continuous, long-term use. Overall, a DMR
       system (for example, a Motorola MOTOTRBO or Hytera DMR package) would provide the team
       with a private, robust radio network that is scalable, secure, and effective in urban conditions. It is
      the best choice if the user prioritizes having full control and guaranteed uptime of the
      communication system.


     • Secondary Recommendation – Push-to-Talk over Cellular (PoC) Solution: If leveraging existing
       infrastructure is acceptable, a PoC system offers maximum range and easy scalability. This would
       involve equipping each team member with a PoC radio or smartphone (plus headset) and using a
       PoC service (e.g. an app like ESChat or a carrier’s PTT platform). Scalability: Adding users is trivial –
       50 or even 500 users can be in the group, since the heavy lifting is done by servers. Coverage: As
       long as there is 4G/5G coverage, the team stays connected across the entire city (and even beyond).
       Urban cellular networks are quite reliable, and given that the team is in Greenwich, CT (a developed
       area), coverage should be strong. Reliability: In normal operations, PoC will work flawlessly with low
       latency; however, it’s contingent on the cell network. To mitigate that, one could choose a carrier
       known for good service locally, and possibly use a PTT service with priority access (for example,
       FirstNet or a mission-critical PTT if the use-case qualifies for it) 7 . Alternatively, a hybrid device (like
       a dual-mode DMR + LTE radio) could be used to failover to direct radio if cellular fails – but that
       adds complexity. Ease of use: PoC devices can be as user-friendly as a standard two-way radio with a
       PTT button. There is very little infrastructure to set up (no towers or repeaters of your own – just
       account setup on the PoC service). Security: Communications are encrypted and can be recorded or
       logged if needed for review. Considering budget is not an issue, subscription costs for a high-quality
       PoC service can be accommodated. In summary, PoC is ideal if the team might be geographically
       spread and you want a hassle-free deployment. It is effectively outsourcing the radio network to a
       carrier. The trade-off is dependence on that carrier, but the benefit is an instant, wide-area solution
       with minimal hardware fuss. If the user’s operations are during times when networks are stable (and
       not during, say, disaster outages), PoC could be the most convenient solution.




                                                        14
     • Additional Recommendation – Wireless Intercom for Localized Operations: If the user’s scenario
       often involves all team members being in a specific venue or within a few blocks (for example,
       coordinating an event or a campus security team), a professional wireless intercom system is
       highly recommended. A system like Clear-Com FreeSpeak II/Edge or similar would allow the leader
       to talk to everyone with open-mic convenience, and team members can speak up immediately as
       needed. Scalability: With the right configuration (multiple transceivers), 50 users can be
       accommodated on one conference line 44 . Reliability: The closed DECT network will be
       interference-free and encrypted 45 . Ease: Users will appreciate not having to push-to-talk for every
       comment – it creates a very natural communication flow, which can be crucial during fast-paced
       operations. The main limitation is that if the team disperses beyond the coverage bubble, they’d lose
       contact; so this recommendation shines for contained areas or events. It could even be used in
       conjunction with another system: e.g. during a large event, the core team on-site uses the intercom,
       but if some members go off-site, they switch to DMR or PoC. In purely local deployments, however,
       the wireless intercom provides a level of interactivity that push-to-talk systems cannot (true
       simultaneous communication).


     • Niche Recommendation – Mesh Network Radios for Ultimate Resilience: If the user anticipates
       scenarios with little to no infrastructure (for example, responding in disasters where cell networks
       and power might be out, or covering a very challenging radio environment with high-rise density), a
       mesh radio solution is worth considering. For instance, a goTenna Pro X or a military-grade
       MANET radio system could keep the 50-person team connected even if they split up and move
       through the city’s nooks and crannies. The mesh will automatically extend as far as the team
       spreads. This is somewhat experimental for civilian use at this scale, but it promises unparalleled
       independence. Scalability: Many mesh systems can handle 50 nodes and continue expanding, and
       the network actually improves with more nodes providing alternate routes. Battery and Portability:
       Modern mesh units are manageable in size and have good battery life (as noted, some can run 24–
       30 hours on a charge) 41 . The downside is the configuration and cost – these are advanced devices
       that may require expert setup and are pricier per unit than typical radios. Therefore, for routine use,
       one might choose DMR or PoC, but keep a mesh system as a backup or for special operations. Mesh
       shines when you absolutely require a self-healing network that adapts to the team’s movements
       and cannot afford a single point of failure.


Final Thoughts: In conclusion, the best approach may involve a combination of these systems, tailored to
specific needs. However, if we prioritize scalability and reliability for 50 users in a dense city, a DMR
radio network with possible LTE/PoC integration is the top pick. DMR gives you a rock-solid foundation
that you control, and PoC can complement it for extended reach or convenience (indeed, some
organizations use dual-mode radios that utilize DMR on-site and PoC off-site). For day-to-day ease of
operation, ensure everyone has comfortable earpieces and clear instructions on device use. With the
recommended systems, the user’s commands will be heard loud and clear by all team members, with
minimal downtime and interference. Each option above is scalable, secure, and proven – the final choice
depends on how much one values total independence (DMR/mesh) versus leveraging existing networks
(PoC), and whether full-duplex communication (intercom) is required for the team’s workflow. By
implementing the appropriate solution, the user can confidently coordinate their 20–50 person team across
the urban landscape, knowing the communication link will remain solid in the long term.




                                                     15
 1    2    3    6   15   17   Push-to-Talk Over Cellular: Choosing the Right Communication Solution - Hytera Blogs
https://www.hytera.com/en/connect/blog/push-to-talk-over-cellular-choosing-the-right-communication-solution

 4    5   10   11   Get Hytera PNC370, PoC System Portable Radios For Sale - Hytera
https://www.hytera.com/en/product-new/lte-broadband/poc-radio/pnc370.html

 7    8   12   13   14   18   19   25   26   27   32   52   54   56   Comparison of Push-to-Talk over Cellular (PoC) vs.
Traditional Radio C
https://www.peakptt.com/pages/comparison-of-push-to-talk-over-cellular-poc-vs-traditional-radio-communication-methods?
srsltid=AfmBOoqrXSs2u1HxHx1IXFNCMN0TctUhbIZUzBfZWg58rri14uTMm6Hi

 9   16   Push-to-Talk over Cellular – The NEXT Revolution -- Occupational Health & Safety
https://ohsonline.com/webcasts/2024/11/push-to-talk-over-cellular.aspx

20   28   29   Longer battery life and greater power efficiency - Benefits of DMR Rad
https://talkpod.com/blogs/dmr-wiki/longer-battery-life-and-greater-power-efficiency-benefits-of-dmr-radios-talkpod

21   31   Retevis NR30D Long Range Waterproof DMR Radio with Earpieces and Mics 2pack
https://www.retevis.com/retevis-nr30d-long-range-waterproof-dmr-radio-with-earpieces-and-mics-2pack-us?
srsltid=AfmBOopveXcS01wcFCm1jDMGWd4kaqJNrsrM_VxwIsQa2Tgc-vC5Q5yN

22   23   24   53   The Communication Distance of Two Way Radios - Hytera Blogs
https://www.hytera.com/en/connect/blog/the-communication-distance-of-two-way-radios

30   How to Charge Two-Way Radios and Extend Battery Life | Hytera US
https://www.hytera.us/resources/how-to-charge-two-way-radios-and-extend-battery-life/

33   34   35   38   39   40   Wireless MANET Mesh Networking | Stream Voice, Data & Video
https://silvustechnologies.com/why-silvus/technology/introduction-to-manet/

36   37   41   55   China Handheld PTT MESH Radio Base Station manufacturers and suppliers | Iwave
https://www.iwavecomms.com/handheld-ptt-mesh-radio-base-station-product/

42   Clear-Com FreeSpeak II Intercom Products - G&G Technologies
https://www.ggvideo.com/clear_freespeakii.php

43   Clear-Com FreeSpeak II Installation Video 1 of 4 (Instructional)
https://clearcom.com/News-Events/Details/clear-com-freespeak-ii-installation-video-1-of-4-instructional

44   Clear-Com Freespeak Edge 5GHz Wireless Beltpack - Sweetwater
https://www.sweetwater.com/store/detail/FSEBP50X4--clear-com-freespeak-edge-5ghz-wireless-beltpack

45   46   Hollyland Solidcom C1-8S Full-Duplex Wireless DECT ... - B&H Photo
https://www.bhphotovideo.com/c/product/1729120-REG/hollyland_hl_solidcom_c1_7s_solidcom_c1_7s_full_duplex_wireless.html

47   FSII-BP19 - Clear-Com
https://www.clearcom.com/product/1-9ghz-digital-beltpack/

48   [PDF] FreeSpeak II 1.9 GHz Beltpack - Clear-Com
https://clearcom.com/DownloadCenter/datasheets/FreeSpeakII/FSII-BP19_FreeSpeakII_19_Beltpack_Datasheet.pdf

49   Beltpack Batteries - Clear-Com
https://www.clearcom.com/DownloadCenter/manuals/FreeSpeakEdge_Online_Manual/Content/Edge_for_Eclipse/Topics/
Beltpack_batteries.htm

50   [PDF] AC80 Battery Charger - Clear-Com
https://clearcom.com/DownloadCenter/datasheets/FreeSpeakEdge/AC80_FreeSpeakEdge_Charger_Datasheet.pdf




                                                                        16
51   Microsoft Word - FRS - GMRS Radios by KI6DGB 101910.doc
https://www.cert-la.com/downloads/radio/FRS-GMRS-Radios-by-KI6DGB-10-19-10.pdf




                                                          17


---

## Related

- [[00-24 SO (Sovereign Operations)/06 Research (R&D) (The Lab)/06 Research|06 Research]]
