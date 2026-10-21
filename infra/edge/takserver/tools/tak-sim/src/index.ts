// TAK Simulation Stack — main entry point
//
// Starts:
//   1. NATS connection (local nats-server)
//   2. CoT relay → FreeTAKServer (TCP) or UDP broadcast to ATAK clients
//   3. Drone simulator — 2 virtual UAVs flying real GPS routes
//   4. Geofence engine — evaluates all drones against defined zones
//   5. AI TAK agent — breach handler, NLP commands, telemetry relay to ATAK
//
// Usage:
//   npm run sim
//   NATS_URL=nats://localhost:4222 TAK_HOST=192.168.1.x TAK_PORT=8087 npm run sim

import { connect, StringCodec } from 'nats'
import { CotRelay } from './cot-relay.js'
import { DroneSimulator, DEFAULT_DRONES } from './drone-simulator.js'
import { GeofenceEngine } from './geofence-engine.js'
import { AITakAgent } from './ai-tak-agent.js'
import {
  geofenceCircleCot,
  waypointCot,
  friendlyUAVCot,
} from './cot-builder.js'

// ── Config from env ──────────────────────────────────────────────────────────
const NATS_URL     = process.env.NATS_URL      ?? 'nats://localhost:4222'
const TAK_HOST     = process.env.TAK_HOST      ?? '127.0.0.1'
const TAK_PORT     = parseInt(process.env.TAK_PORT ?? '8087')
const TAK_PROTO    = (process.env.TAK_PROTO    ?? 'tcp') as 'tcp' | 'udp'
const HERMES_URL   = process.env.HERMES_URL    ?? 'http://localhost:8642'
const HERMES_PROFILE = process.env.HERMES_PROFILE ?? 'default'
const USE_HERMES   = process.env.USE_HERMES    === 'true'

// ── Geofence definitions ─────────────────────────────────────────────────────
// Edit these to match your area of operations
const GEOFENCES = [
  {
    id: 'zone-alpha',
    name: 'ALPHA',
    lat: 38.8995,
    lon: -77.0380,
    radiusM: 250,
    action: 'launch_drone' as const,
    actionTarget: 'REAPER-2',   // REAPER-2 launches when ALPHA is breached
    severity: 'HIGH' as const,
  },
  {
    id: 'zone-bravo',
    name: 'BRAVO',
    lat: 38.9005,
    lon: -77.0330,
    radiusM: 150,
    action: 'alert' as const,
    severity: 'CRITICAL' as const,
  },
]

// ── Waypoints to pre-place on the map ────────────────────────────────────────
const WAYPOINTS = [
  { id: 'cp-alpha', name: 'CP-ALPHA', lat: 38.8977, lon: -77.0365 },
  { id: 'cp-bravo', name: 'CP-BRAVO', lat: 38.9010, lon: -77.0310 },
  { id: 'fob-main', name: 'FOB-MAIN', lat: 38.8950, lon: -77.0400 },
]

async function main(): Promise<void> {
  console.log(`
╔══════════════════════════════════════════════════════════════╗
║  TAK-SIM  —  Autonomous Operations Simulation Stack          ║
╠══════════════════════════════════════════════════════════════╣
║  NATS    → ${NATS_URL.padEnd(47)}║
║  TAK     → ${TAK_PROTO.toUpperCase()}://${TAK_HOST}:${TAK_PORT}${' '.repeat(Math.max(0, 41 - TAK_HOST.length - TAK_PORT.toString().length))}║
║  Hermes  → ${(USE_HERMES ? HERMES_URL : 'disabled').padEnd(47)}║
╚══════════════════════════════════════════════════════════════╝
`)

  // 1. NATS
  console.log('[TAK-SIM] Connecting to NATS...')
  const nc = await connect({ servers: NATS_URL })
  console.log(`[TAK-SIM] NATS connected → ${NATS_URL}`)
  const sc = StringCodec()

  // 2. CoT relay
  const relay = new CotRelay({
    host: TAK_HOST,
    port: TAK_PORT,
    protocol: TAK_PROTO,
    nats: nc,
    natsCotSubject: 'tak.cot.outbound',
  })
  await relay.connect()
  console.log(`[TAK-SIM] CoT relay ready → ${TAK_PROTO.toUpperCase()}://${TAK_HOST}:${TAK_PORT}`)

  // 3. Drone simulator
  const sim = new DroneSimulator(nc, relay)
  for (const cfg of DEFAULT_DRONES) {
    sim.addDrone(cfg)
  }

  // 4. Geofence engine
  const fenceEngine = new GeofenceEngine(nc)
  for (const z of GEOFENCES) {
    fenceEngine.addZone(z)
  }

  // 5. AI TAK agent
  const agent = new AITakAgent(
    nc,
    relay,
    USE_HERMES ? { url: HERMES_URL, profile: HERMES_PROFILE } : undefined,
    sim
  )

  // ── Wire geofence evaluation to telemetry stream ─────────────────────────
  // Every telemetry update → geofence check → breach event → AI action
  const telSub = nc.subscribe('constellation.telemetry.>')
  ;(async () => {
    for await (const msg of telSub) {
      try {
        const env = JSON.parse(sc.decode(msg.data))
        if (env.msg_name !== 'GlobalPositionInt') continue
        const lat = env.fields.latitude / 1e7
        const lon = env.fields.longitude / 1e7
        const heading = env.fields.heading / 100
        fenceEngine.evaluate(env.drone_id, lat, lon, heading)
      } catch { /* skip */ }
    }
  })()

  // ── Pre-load map: draw geofences and waypoints on ATAK at startup ─────────
  console.log('\n[TAK-SIM] Pre-loading map objects...')
  await new Promise(r => setTimeout(r, 500))  // brief pause for TAK connection

  for (const z of GEOFENCES) {
    const xml = geofenceCircleCot(z.id, z.name, z.lat, z.lon, z.radiusM)
    await relay.send(xml)
    console.log(`  ◯ Geofence drawn: ${z.name} (r=${z.radiusM}m) action=${z.action}`)
  }

  for (const wp of WAYPOINTS) {
    const xml = waypointCot(wp.id, wp.name, wp.lat, wp.lon)
    await relay.send(xml)
    console.log(`  ◆ Waypoint placed: ${wp.name}`)
  }

  // ── Start everything ─────────────────────────────────────────────────────
  agent.start()
  sim.startAll()

  console.log(`
[TAK-SIM] ALL SYSTEMS RUNNING
  • ${DEFAULT_DRONES.length} virtual drones flying — watch them move on your ATAK map
  • ${GEOFENCES.length} geofences active — breach triggers AI autonomous response
  • ${WAYPOINTS.length} waypoints loaded on map
  • AI agent listening on ops.breach.> and ai.tak.nlp

  To send a manual command (in another terminal):
    nats pub ai.tak.nlp "add hostile marker at 38.899, -77.037"
    nats pub ai.tak.nlp "draw 400m geofence at 38.903, -77.032"

  To manually trigger a breach simulation:
    nats pub ops.breach.zone-alpha '{"zoneId":"zone-alpha","zoneName":"ALPHA","entityId":"TEST-1","breachType":"INSIDE","distanceM":0,"entityLat":38.8995,"entityLon":-77.038,"action":"launch_drone","actionTarget":"REAPER-2","severity":"HIGH","timestamp":"${new Date().toISOString()}"}'

  Press Ctrl+C to stop.
`)

  // ── Graceful shutdown ────────────────────────────────────────────────────
  process.on('SIGINT', async () => {
    console.log('\n[TAK-SIM] Shutting down...')
    sim.stopAll()
    agent.stop()
    await nc.drain()
    process.exit(0)
  })
}

main().catch(err => {
  console.error('[TAK-SIM] Fatal error:', err)
  process.exit(1)
})
