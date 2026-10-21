// Geofence engine — stores zones, evaluates breach, fires NATS events
// Reuses the same Haversine + bearing logic as APEX-SENTINEL's ProtectedZoneBreachDetector

import type { NatsConnection } from 'nats'
import { StringCodec } from 'nats'

export type BreachType = 'INSIDE' | 'ENTERING' | 'APPROACHING'
export type ZoneAction = 'launch_drone' | 'alert' | 'mark_hostile' | 'custom'

export interface GeofenceZone {
  id: string
  name: string
  lat: number
  lon: number
  radiusM: number
  action: ZoneAction
  actionTarget?: string  // drone id for launch_drone
  severity: 'LOW' | 'MEDIUM' | 'HIGH' | 'CRITICAL'
}

export interface BreachEvent {
  zoneId: string
  zoneName: string
  entityId: string
  breachType: BreachType
  distanceM: number
  entityLat: number
  entityLon: number
  action: ZoneAction
  actionTarget?: string
  severity: string
  timestamp: string
}

const R_EARTH = 6371000

function haversineM(lat1: number, lon1: number, lat2: number, lon2: number): number {
  const toRad = (d: number) => d * Math.PI / 180
  const dLat = toRad(lat2 - lat1)
  const dLon = toRad(lon2 - lon1)
  const a = Math.sin(dLat/2)**2 + Math.cos(toRad(lat1)) * Math.cos(toRad(lat2)) * Math.sin(dLon/2)**2
  return R_EARTH * 2 * Math.atan2(Math.sqrt(a), Math.sqrt(1-a))
}

function bearingDeg(lat1: number, lon1: number, lat2: number, lon2: number): number {
  const toRad = (d: number) => d * Math.PI / 180
  const toDeg = (r: number) => r * 180 / Math.PI
  const dLon = toRad(lon2 - lon1)
  const y = Math.sin(dLon) * Math.cos(toRad(lat2))
  const x = Math.cos(toRad(lat1)) * Math.sin(toRad(lat2)) -
            Math.sin(toRad(lat1)) * Math.cos(toRad(lat2)) * Math.cos(dLon)
  return (toDeg(Math.atan2(y, x)) + 360) % 360
}

function angleDiff(a: number, b: number): number {
  const diff = Math.abs(a - b) % 360
  return diff > 180 ? 360 - diff : diff
}

export class GeofenceEngine {
  private zones = new Map<string, GeofenceZone>()
  private breachState = new Map<string, Set<string>>()  // entityId → Set<zoneId> currently breached
  private nc: NatsConnection
  private sc = StringCodec()

  constructor(nc: NatsConnection) {
    this.nc = nc
  }

  addZone(zone: GeofenceZone): void {
    this.zones.set(zone.id, zone)
    console.log(`[GeofenceEngine] Zone added: ${zone.name} @ (${zone.lat.toFixed(4)}, ${zone.lon.toFixed(4)}) r=${zone.radiusM}m action=${zone.action}`)
  }

  removeZone(id: string): void {
    this.zones.delete(id)
  }

  getZones(): GeofenceZone[] {
    return Array.from(this.zones.values())
  }

  // Call this every time a telemetry update arrives
  evaluate(entityId: string, lat: number, lon: number, heading?: number): BreachEvent[] {
    const events: BreachEvent[] = []

    if (!this.breachState.has(entityId)) {
      this.breachState.set(entityId, new Set())
    }
    const activeBreaches = this.breachState.get(entityId)!

    for (const zone of this.zones.values()) {
      const distM = haversineM(lat, lon, zone.lat, zone.lon)
      const alertRadiusM = zone.radiusM * 1.1   // 10% buffer for ENTERING/APPROACHING

      let breachType: BreachType | null = null

      if (distM < zone.radiusM) {
        breachType = 'INSIDE'
      } else if (distM <= alertRadiusM && heading !== undefined) {
        const brngToZone = bearingDeg(lat, lon, zone.lat, zone.lon)
        breachType = angleDiff(heading, brngToZone) <= 45 ? 'ENTERING' : 'APPROACHING'
      }

      const breachKey = `${entityId}:${zone.id}`

      if (breachType) {
        const alreadyBreached = activeBreaches.has(zone.id)

        if (!alreadyBreached || breachType === 'INSIDE') {
          // New breach — fire event
          activeBreaches.add(zone.id)
          const ev: BreachEvent = {
            zoneId: zone.id,
            zoneName: zone.name,
            entityId,
            breachType,
            distanceM: distM,
            entityLat: lat,
            entityLon: lon,
            action: zone.action,
            actionTarget: zone.actionTarget,
            severity: zone.severity,
            timestamp: new Date().toISOString(),
          }
          events.push(ev)

          // Publish to NATS
          const subject = `ops.breach.${zone.id}`
          this.nc.publish(subject, this.sc.encode(JSON.stringify(ev)))
          console.log(`[GeofenceEngine] BREACH ${breachType} — entity=${entityId} zone=${zone.name} dist=${distM.toFixed(0)}m`)
        }
      } else {
        // Entity has left — clear breach state
        if (activeBreaches.has(zone.id)) {
          activeBreaches.delete(zone.id)
          console.log(`[GeofenceEngine] CLEARED — entity=${entityId} exited zone=${zone.name}`)
        }
      }
    }

    return events
  }
}
