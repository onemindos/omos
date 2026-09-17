# OneMind OS — Recon 01: Data Source Overlap
*Generated 2026-07-06 by Legacy (direct). Read-only analysis.*

## Source inventories (raw)

- **Crucix (27):** acled, adsb, bls, bluesky, cisa-kev, cloudflare-radar, comtrade, eia, epa, firms, fred, gdelt, gscpi, kiwisdr, noaa, ofac, opensanctions, opensky, patents, reddit, reliefweb, safecast, ships, space, telegram, treasury, usaspending, who, yfinance
- **Osiris (data routes):** air-quality, cctv, conflicts, country-risk, crypto, cyber-threats, earthquakes, fires, flights, frontlines, gdelt, health, infrastructure, live-news, malware, maritime, markets, news, radar, region-dossier, satellites, scm-suppliers, sentinel, space-weather, weather + OSINT toolkit (bgp, certs, cve, dns, github, ip, leaks, mac, phone, sanctions, shodan, sweep, threats, whois)
- **WorldMonitor:** 500+ news feeds + map layers in categories: Aviation, Conflict, Country Risk, Cyber, Maritime, Natural Disasters, News/Hotspots
- **ORP / om-unity:** 39 *protocol* adapters (AIS, ADS-B, MAVLink, OPC-UA, MQTT, Modbus, Zeek, syslog, GRIB, CoT, KLV, CCSDS+SGP4, HL7, Kafka, NATS…). Currently 5,938 live AIS ships in graph. This is protocol-level ingest, not curated feeds.

## 1. Overlap matrix (by domain)

| Domain | Crucix | Osiris | WorldMonitor | ORP/om-unity | Best pick |
|---|:--:|:--:|:--:|:--:|---|
| Flights / ADS-B | ✅ opensky, adsb | ✅ flights | ✅ Aviation | ✅ live protocol | **ORP** (live protocol) + Crucix opensky for enrich |
| Ships / AIS | ✅ ships | ✅ maritime | ✅ Maritime x3 | ✅ 5,938 live | **ORP** (already live) |
| Fires | ✅ firms | ✅ fires | — | via GRIB | **Crucix firms** (NASA FIRMS clean) |
| Earthquakes | — | ✅ earthquakes | ✅ Nat. Disasters | — | **Osiris earthquakes** (USGS) |
| Weather | ✅ noaa | ✅ weather | — | ✅ GRIB | **Crucix noaa** |
| Satellites / space | ✅ space | ✅ satellites, space-weather | — | ✅ CCSDS+SGP4 | **ORP** (SGP4 native) + Osiris space-weather |
| Sanctions | ✅ ofac, opensanctions | ✅ osint/sanctions | — | — | **Crucix opensanctions** (richest) |
| Conflict | ✅ acled | ✅ conflicts, frontlines | ✅ Conflict x2 | — | **Crucix acled** (event-level) + Osiris frontlines (geo) |
| Cyber | ✅ cisa-kev, cloudflare-radar | ✅ cyber-threats, malware, osint/cve, shodan, bgp | ✅ Cyber | ✅ Zeek/syslog | **Osiris** (deepest OSINT toolkit) |
| Finance | ✅ fred,bls,eia,yfinance,comtrade,treasury,usaspending,gscpi,patents | ✅ markets, scm-suppliers | — | — | **Crucix** (9 economic feeds) |
| News | ✅ gdelt | ✅ news, live-news, gdelt | ✅ 500+ feeds | — | **WorldMonitor** (500+ feeds unmatched) |
| Social / OSINT | ✅ bluesky, reddit, telegram | ✅ osint/* toolkit | — | — | **Split**: Crucix social feeds + Osiris OSINT recon |
| Environment | ✅ epa, safecast | ✅ air-quality | — | — | **Crucix** (epa + safecast radiation) |
| Health | ✅ who, reliefweb | ✅ health | — | — | **Crucix who** |
| CCTV / video | — | ✅ cctv (2K+ cams) | — | ✅ KLV | **Osiris cctv** (unique 2K cameras) |
| SIGINT / radio | ✅ kiwisdr | — | — | — | **Crucix kiwisdr** (unique) |
| Drones / MAVLink | — | — | — | ✅ live | **ORP/om-bridge-mavlink** (live) |

## 2. Unique sources (single-project only)

- **Crucix only:** kiwisdr (SDR radio), safecast (radiation), gscpi (supply-chain pressure), patents, comtrade, usaspending, bls, eia, reliefweb, bluesky
- **Osiris only:** cctv (2K+ cameras), full OSINT recon toolkit (shodan, bgp, certs, dns, whois, leaks, mac, phone, sweep), radar, sentinel, scm-suppliers, infrastructure
- **WorldMonitor only:** the 500+ curated news feed catalog + CII country-stress scoring algorithm
- **ORP only:** live protocol ingest (OPC-UA, Modbus, HL7, KLV, CCSDS, Zeek) — industrial/military protocols nobody else touches

## 3. Best-of-breed verdict per project

- **ORP/om-unity** → LIVE telemetry protocols (AIS, ADS-B, MAVLink, industrial). The real-time spine.
- **Crucix** → curated API feeds (fires, weather, sanctions, finance, health, environment, SDR). The intel enrichment layer. Cleanest adapter code — easiest to port to NATS.
- **Osiris** → OSINT recon toolkit + CCTV + cyber depth. The investigation layer.
- **WorldMonitor** → news feed catalog + CII scoring. The narrative/awareness layer.

## 4. Minimal port list (cover every domain, zero duplication)

**Tier 1 — from ORP (already live, just wire NATS):**
1. AIS ships → `tel.ais.*`
2. ADS-B flights → `tel.adsb.*`
3. MAVLink drones → `tel.drone.*` (already via om-bridge-mavlink)

**Tier 2 — port these Crucix adapters → NATS:**
4. firms → `det.fire.*`
5. noaa → `det.weather.*`
6. acled → `det.conflict.*`
7. opensanctions → `intel.sanctions.*`
8. space → `det.space.*`
9. who → `intel.health.*`
10. epa + safecast → `det.environment.*`
11. kiwisdr → `det.sigint.*` (unique)
12. Finance bundle (fred/yfinance/eia) → `intel.finance.*`

**Tier 3 — port these Osiris routes → NATS:**
13. earthquakes → `det.seismic.*`
14. cctv → `tel.cctv.*` (unique 2K cams)
15. satellites/space-weather → `det.space.*`
16. OSINT toolkit → on-demand `intel.osint.*` (not streaming — query service)

**Tier 4 — WorldMonitor:**
17. 500+ news feeds → `intel.news.*`
18. CII scoring algorithm → Legacy/Oracle country-risk service

**Result:** ~18 NATS publishers cover all 15 domains with no duplication. ORP handles 3, Crucix handles 9, Osiris handles 4, WorldMonitor handles 2.
