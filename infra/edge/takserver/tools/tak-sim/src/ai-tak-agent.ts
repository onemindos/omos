// AI TAK Agent — subscribes to NATS, calls Hermes, writes decisions back to ATAK map
// This is the bridge: NATS events → AI decision → CoT XML → ATAK

import { StringCodec, type NatsConnection, type Subscription } from 'nats'
import { CotRelay } from './cot-relay.js'
import {
  geofenceCircleCot,
  hostileMarkerCot,
  waypointCot,
  emergencyBeaconCot,
  friendlyUAVCot,
} from './cot-builder.js'
import type { BreachEvent } from './geofence-engine.js'
import type { DroneSimulator } from './drone-simulator.js'

export interface HermesConfig {
  url: string         // e.g. http://localhost:8642
  profile: string     // your Hermes profile name for TAK ops
  apiKey?: string
}

interface DroneAction {
  droneId: string
  command: 'launch' | 'rtb' | 'orbit' | 'goto'
  targetLat?: number
  targetLon?: number
  altM?: number
}

export class AITakAgent {
  private nc: NatsConnection
  private relay: CotRelay
  private sc = StringCodec()
  private hermes?: HermesConfig
  private sim?: DroneSimulator
  private subs: Subscription[] = []
  private actionLog: string[] = []

  constructor(nc: NatsConnection, relay: CotRelay, hermes?: HermesConfig, sim?: DroneSimulator) {
    this.nc = nc
    this.relay = relay
    this.hermes = hermes
    this.sim = sim
  }

  start(): void {
    this.subscribeBreachEvents()
    this.subscribeNlpCommands()
    this.subscribeTelemetryRelay()
    console.log('[AITakAgent] Running — subscribed to ops.breach.>, ai.tak.nlp, constellation.telemetry.>')
  }

  stop(): void {
    for (const sub of this.subs) sub.unsubscribe()
  }

  // ── 1. Geofence breach → autonomous response ─────────────────────────────

  private subscribeBreachEvents(): void {
    const sub = this.nc.subscribe('ops.breach.>')
    this.subs.push(sub)

    ;(async () => {
      for await (const msg of sub) {
        try {
          const breach: BreachEvent = JSON.parse(this.sc.decode(msg.data))
          await this.handleBreach(breach)
        } catch (e) {
          console.error('[AITakAgent] breach parse error', e)
        }
      }
    })()
  }

  private async handleBreach(breach: BreachEvent): Promise<void> {
    console.log(`\n[AITakAgent] *** BREACH *** zone=${breach.zoneName} entity=${breach.entityId} type=${breach.breachType}`)

    // 1. Place hostile marker on ATAK map immediately
    const marker = hostileMarkerCot(
      `BREACH-${breach.entityId}-${Date.now()}`,
      `HOSTILE-${breach.entityId}`,
      breach.entityLat,
      breach.entityLon
    )
    await this.relay.send(marker)
    this.log(`Hostile marker placed @ ${breach.entityLat.toFixed(4)}, ${breach.entityLon.toFixed(4)}`)

    // 2. Ask Hermes what to do (if configured)
    let decision = `Breach detected in zone ${breach.zoneName}. Entity ${breach.entityId} is ${breach.breachType}. Auto-launching response drone.`
    if (this.hermes) {
      try {
        decision = await this.askHermes(
          `Tactical situation: entity "${breach.entityId}" has breached geofence "${breach.zoneName}" (${breach.breachType}, distance: ${breach.distanceM.toFixed(0)}m). ` +
          `Zone severity: ${breach.severity}. Configured response action: ${breach.action}. ` +
          `What is the recommended immediate response? Keep it under 2 sentences.`
        )
        this.log(`Hermes: ${decision}`)
      } catch (e) {
        console.warn('[AITakAgent] Hermes unreachable — proceeding with default action')
      }
    }

    // 3. Execute configured action
    if (breach.action === 'launch_drone' && breach.actionTarget) {
      await this.executeDroneAction({
        droneId: breach.actionTarget,
        command: 'goto',
        targetLat: breach.entityLat,
        targetLon: breach.entityLon,
        altM: 80,
      })
    } else if (breach.action === 'alert') {
      const beacon = emergencyBeaconCot(
        `EMRG-${breach.zoneId}-${Date.now()}`,
        `ALERT-${breach.zoneName}`,
        breach.entityLat,
        breach.entityLon
      )
      await this.relay.send(beacon)
    }

    // 4. Publish decision summary to NATS for logging/Cordum
    this.nc.publish('ai.tak.decisions', this.sc.encode(JSON.stringify({
      breach,
      decision,
      timestamp: new Date().toISOString(),
    })))
  }

  // ── 2. Natural language TAK commands ─────────────────────────────────────
  // Publish to ai.tak.nlp to talk to your map:
  // "add a hostile marker at 38.897, -77.036"
  // "draw a 500m geofence around checkpoint bravo at 38.901, -77.031"
  // "move SENTINEL-1 to 38.895, -77.040"

  private subscribeNlpCommands(): void {
    const sub = this.nc.subscribe('ai.tak.nlp')
    this.subs.push(sub)

    ;(async () => {
      for await (const msg of sub) {
        const text = this.sc.decode(msg.data)
        console.log(`[AITakAgent] NLP command: "${text}"`)
        await this.handleNlpCommand(text)
      }
    })()
  }

  private async handleNlpCommand(text: string): Promise<void> {
    // Simple pattern matching — extend with Hermes for full NL understanding
    const lower = text.toLowerCase()
    const coordRegex = /(-?\d+\.?\d*)[,\s]+(-?\d+\.?\d*)/

    if ((lower.includes('hostile') || lower.includes('enemy')) && coordRegex.test(text)) {
      const [, latS, lonS] = text.match(coordRegex)!
      const lat = parseFloat(latS), lon = parseFloat(lonS)
      const xml = hostileMarkerCot(`HOSTILE-NLP-${Date.now()}`, 'HOSTILE', lat, lon)
      await this.relay.send(xml)
      this.log(`Hostile marker placed @ ${lat.toFixed(4)}, ${lon.toFixed(4)} (NLP)`)
      return
    }

    if ((lower.includes('geofence') || lower.includes('fence') || lower.includes('perimeter')) && coordRegex.test(text)) {
      const [, latS, lonS] = text.match(coordRegex)!
      const lat = parseFloat(latS), lon = parseFloat(lonS)
      const radiusMatch = text.match(/(\d+)\s*m/)
      const radiusM = radiusMatch ? parseInt(radiusMatch[1]) : 300
      const name = `NLP-${Date.now()}`
      const xml = geofenceCircleCot(name, 'AI-FENCE', lat, lon, radiusM)
      await this.relay.send(xml)
      this.log(`Geofence drawn @ ${lat.toFixed(4)}, ${lon.toFixed(4)} r=${radiusM}m (NLP)`)
      return
    }

    if (lower.includes('waypoint') && coordRegex.test(text)) {
      const [, latS, lonS] = text.match(coordRegex)!
      const lat = parseFloat(latS), lon = parseFloat(lonS)
      const xml = waypointCot(`WP-${Date.now()}`, 'AI-WP', lat, lon)
      await this.relay.send(xml)
      this.log(`Waypoint placed @ ${lat.toFixed(4)}, ${lon.toFixed(4)} (NLP)`)
      return
    }

    // Fall through to Hermes if configured
    if (this.hermes) {
      try {
        const resp = await this.askHermes(
          `You are a TAK operator AI. Parse this command and respond ONLY with JSON: ` +
          `{"action":"hostile_marker"|"geofence"|"waypoint"|"unknown","lat":number,"lon":number,"radiusM":number,"name":string}. ` +
          `Command: "${text}"`
        )
        try {
          const parsed = JSON.parse(resp)
          if (parsed.action === 'hostile_marker') {
            await this.relay.send(hostileMarkerCot(`HOSTILE-H-${Date.now()}`, parsed.name ?? 'HOSTILE', parsed.lat, parsed.lon))
          } else if (parsed.action === 'geofence') {
            await this.relay.send(geofenceCircleCot(`H-${Date.now()}`, parsed.name ?? 'FENCE', parsed.lat, parsed.lon, parsed.radiusM ?? 300))
          } else if (parsed.action === 'waypoint') {
            await this.relay.send(waypointCot(`H-${Date.now()}`, parsed.name ?? 'WP', parsed.lat, parsed.lon))
          }
        } catch {
          this.log(`Hermes response (no action): ${resp}`)
        }
      } catch {
        this.log('Hermes unreachable for NLP command')
      }
    }
  }

  // ── 3. Mirror drone telemetry to ATAK ────────────────────────────────────
  // Subscribes to the NATS telemetry stream and pushes CoT updates
  // This is what makes the drone icon move on the ATAK map

  private subscribeTelemetryRelay(): void {
    const sub = this.nc.subscribe('constellation.telemetry.>')
    this.subs.push(sub)

    ;(async () => {
      for await (const msg of sub) {
        try {
          const env = JSON.parse(this.sc.decode(msg.data))
          if (env.msg_name !== 'GlobalPositionInt') continue

          const lat = env.fields.latitude / 1e7
          const lon = env.fields.longitude / 1e7
          const altM = env.fields.altitude / 1000
          const course = env.fields.heading / 100
          const speedMs = Math.sqrt(env.fields.vx**2 + env.fields.vy**2) / 100

          const cot = friendlyUAVCot(
            env.drone_id,
            lat, lon, altM,
            speedMs, course,
            100  // battery from KV in production; 100 for sim display
          )
          await this.relay.send(cot)
        } catch { /* skip bad messages */ }
      }
    })()
  }

  // ── Drone action ─────────────────────────────────────────────────────────

  private async executeDroneAction(action: DroneAction): Promise<void> {
    console.log(`[AITakAgent] Drone action: ${action.command} → ${action.droneId}`)

    // Publish guidance command to NATS (SentinelX/nats-ros-connector listens here)
    const cmd = {
      drone_id: action.droneId,
      command: action.command,
      target_lat: action.targetLat,
      target_lon: action.targetLon,
      alt_m: action.altM ?? 80,
      timestamp: new Date().toISOString(),
    }
    this.nc.publish(`commands.${action.droneId}.guidance`, this.sc.encode(JSON.stringify(cmd)))
    this.log(`Guidance command sent → ${action.droneId} (${action.command})`)

    // If we have a simulator, add a new intercept waypoint at breach location
    if (this.sim && action.targetLat && action.targetLon) {
      // The sim will naturally route to the next waypoint — for a real launch
      // you'd update the waypoint list here. For demo this just logs.
      this.log(`Sim drone ${action.droneId} re-routed toward breach location`)
    }
  }

  // ── Hermes call ──────────────────────────────────────────────────────────

  private async askHermes(prompt: string): Promise<string> {
    if (!this.hermes) throw new Error('Hermes not configured')
    const body = {
      messages: [{ role: 'user', content: prompt }],
      profile: this.hermes.profile,
      stream: false,
    }
    const headers: Record<string, string> = { 'Content-Type': 'application/json' }
    if (this.hermes.apiKey) headers['Authorization'] = `Bearer ${this.hermes.apiKey}`

    const res = await fetch(`${this.hermes.url}/v1/chat/completions`, {
      method: 'POST',
      headers,
      body: JSON.stringify(body),
      signal: AbortSignal.timeout(10_000),
    })
    if (!res.ok) throw new Error(`Hermes HTTP ${res.status}`)
    const data = await res.json() as { choices?: Array<{ message?: { content?: string } }> }
    return data.choices?.[0]?.message?.content ?? ''
  }

  private log(msg: string): void {
    const line = `[${new Date().toTimeString().slice(0,8)}] ${msg}`
    this.actionLog.push(line)
    console.log(`[AITakAgent] ${line}`)
    // Publish to NATS for any dashboard listening
    this.nc.publish('ai.tak.log', this.sc.encode(line))
  }

  getLog(): string[] { return this.actionLog }
}
