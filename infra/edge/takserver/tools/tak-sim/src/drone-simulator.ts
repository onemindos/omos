// Simulates one or more UAVs flying realistic waypoint-to-waypoint routes
// Publishes telemetry to NATS (same subject pattern as Aero Arc Relay)
// and CoT to TAK via CotRelay so the drone moves live on the ATAK map

import { connect, StringCodec, type NatsConnection } from 'nats'
import { CotRelay } from './cot-relay.js'
import { friendlyUAVCot } from './cot-builder.js'

export interface Waypoint {
  lat: number
  lon: number
  altM: number
  label?: string
}

export interface DroneConfig {
  id: string
  callsign: string
  speedMs: number       // typical cruise speed m/s (e.g. 15 for a quad)
  updateHz: number      // position updates per second
  waypoints: Waypoint[] // route — loops when complete
  battery: number       // starting %, drains at 0.01%/s
}

interface DroneState {
  lat: number
  lon: number
  altM: number
  course: number
  battery: number
  wpIndex: number
}

const R_EARTH = 6371000  // metres

function toRad(d: number): number { return d * Math.PI / 180 }
function toDeg(r: number): number { return r * 180 / Math.PI }

function haversineM(lat1: number, lon1: number, lat2: number, lon2: number): number {
  const dLat = toRad(lat2 - lat1)
  const dLon = toRad(lon2 - lon1)
  const a = Math.sin(dLat/2)**2 + Math.cos(toRad(lat1)) * Math.cos(toRad(lat2)) * Math.sin(dLon/2)**2
  return R_EARTH * 2 * Math.atan2(Math.sqrt(a), Math.sqrt(1-a))
}

function bearing(lat1: number, lon1: number, lat2: number, lon2: number): number {
  const dLon = toRad(lon2 - lon1)
  const y = Math.sin(dLon) * Math.cos(toRad(lat2))
  const x = Math.cos(toRad(lat1)) * Math.sin(toRad(lat2)) -
            Math.sin(toRad(lat1)) * Math.cos(toRad(lat2)) * Math.cos(dLon)
  return (toDeg(Math.atan2(y, x)) + 360) % 360
}

function movePoint(lat: number, lon: number, distM: number, bearingDeg: number): [number, number] {
  const d = distM / R_EARTH
  const b = toRad(bearingDeg)
  const lat2 = Math.asin(Math.sin(toRad(lat)) * Math.cos(d) +
               Math.cos(toRad(lat)) * Math.sin(d) * Math.cos(b))
  const lon2 = toRad(lon) + Math.atan2(
    Math.sin(b) * Math.sin(d) * Math.cos(toRad(lat)),
    Math.cos(d) - Math.sin(toRad(lat)) * Math.sin(lat2)
  )
  return [toDeg(lat2), toDeg(lon2)]
}

export class DroneSimulator {
  private nc: NatsConnection
  private relay: CotRelay
  private sc = StringCodec()
  private drones = new Map<string, { config: DroneConfig; state: DroneState; timer?: ReturnType<typeof setInterval> }>()

  constructor(nc: NatsConnection, relay: CotRelay) {
    this.nc = nc
    this.relay = relay
  }

  addDrone(cfg: DroneConfig): void {
    const start = cfg.waypoints[0]
    const state: DroneState = {
      lat: start.lat,
      lon: start.lon,
      altM: start.altM,
      course: 0,
      battery: cfg.battery,
      wpIndex: 0,
    }
    this.drones.set(cfg.id, { config: cfg, state })
  }

  startAll(): void {
    for (const [id, drone] of this.drones) {
      const intervalMs = Math.round(1000 / drone.config.updateHz)
      drone.timer = setInterval(() => this.tick(id), intervalMs)
      console.log(`[DroneSimulator] ${id} started — ${drone.config.waypoints.length} waypoints, ${drone.config.speedMs}m/s`)
    }
  }

  stopAll(): void {
    for (const drone of this.drones.values()) {
      if (drone.timer) clearInterval(drone.timer)
    }
  }

  private tick(id: string): void {
    const drone = this.drones.get(id)!
    const { config, state } = drone

    const target = config.waypoints[state.wpIndex]
    const distToTarget = haversineM(state.lat, state.lon, target.lat, target.lon)
    const stepM = config.speedMs / config.updateHz

    if (distToTarget < stepM * 2) {
      // Arrived at waypoint — advance to next
      state.wpIndex = (state.wpIndex + 1) % config.waypoints.length
      const next = config.waypoints[state.wpIndex]
      console.log(`[DroneSimulator] ${id} → WP${state.wpIndex} (${next.label ?? `${next.lat.toFixed(4)},${next.lon.toFixed(4)}`})`)
    }

    const brng = bearing(state.lat, state.lon, target.lat, target.lon)
    state.course = brng
    const [newLat, newLon] = movePoint(state.lat, state.lon, stepM, brng)
    state.lat = newLat
    state.lon = newLon

    // Smooth altitude interpolation
    const altDiff = target.altM - state.altM
    state.altM += altDiff * 0.05

    // Battery drain
    state.battery = Math.max(0, state.battery - 0.01 / config.updateHz)

    // Publish to NATS (same envelope shape as Aero Arc Relay)
    const envelope = {
      drone_id: id,
      source: 'tak-sim',
      timestamp_relay: new Date().toISOString(),
      msg_name: 'GlobalPositionInt',
      fields: {
        latitude:     Math.round(state.lat * 1e7),
        longitude:    Math.round(state.lon * 1e7),
        altitude:     Math.round(state.altM * 1000),
        relative_alt: Math.round(state.altM * 1000),
        heading:      Math.round(state.course * 100),
        vx: Math.round(Math.sin(toRad(state.course)) * config.speedMs * 100),
        vy: Math.round(Math.cos(toRad(state.course)) * config.speedMs * 100),
        vz: 0,
      },
    }
    this.nc.publish(
      `constellation.telemetry.${id}`,
      this.sc.encode(JSON.stringify(envelope))
    )

    // Send CoT to ATAK map (drone moves live)
    const cot = friendlyUAVCot(
      config.callsign,
      state.lat,
      state.lon,
      state.altM,
      config.speedMs,
      state.course,
      Math.round(state.battery)
    )
    this.relay.send(cot).catch(() => {})
  }

  getState(id: string): DroneState | undefined {
    return this.drones.get(id)?.state
  }
}

// ── Predefined routes (US-centric, easy to change) ──────────────────────────
// Default patrol: box pattern around a 1km area — change lat/lon to your AO

export const DEFAULT_DRONES: DroneConfig[] = [
  {
    id: 'SENTINEL-1',
    callsign: 'SENTINEL-1',
    speedMs: 12,
    updateHz: 2,
    battery: 100,
    waypoints: [
      { lat: 41.01865, lon: -73.64153, altM: 120, label: 'WP-Alpha' },
      { lat: 41.02000, lon: -73.64153, altM: 120, label: 'WP-Bravo' },
      { lat: 41.02000, lon: -73.63900, altM: 150, label: 'WP-Charlie' },
      { lat: 41.01865, lon: -73.63900, altM: 120, label: 'WP-Delta' },
    ],
  },
  {
    id: 'REAPER-2',
    callsign: 'REAPER-2',
    speedMs: 20,
    updateHz: 2,
    battery: 85,
    waypoints: [
      { lat: 41.01700, lon: -73.64300, altM: 200, label: 'WP-East' },
      { lat: 41.01900, lon: -73.64500, altM: 200, label: 'WP-North' },
      { lat: 41.02100, lon: -73.64300, altM: 220, label: 'WP-West' },
      { lat: 41.01900, lon: -73.64100, altM: 200, label: 'WP-South' },
    ],
  },
]
