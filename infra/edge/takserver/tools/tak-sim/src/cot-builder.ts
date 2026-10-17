// CoT XML builder — covers every map operation ATAK understands

export interface CotPoint {
  lat: number
  lon: number
  hae?: number  // height above ellipsoid (metres)
  ce?: number   // circular error (metres)
  le?: number   // linear error (metres)
}

export interface CotEvent {
  uid: string
  type: string
  how?: string
  callsign: string
  point: CotPoint
  staleMinutes?: number
  remarks?: string
  color?: number       // ARGB integer — fills shape or marker color
  strokeColor?: number
  strokeWeight?: number
  fillOpacity?: number  // 0.0–1.0
  // for geocircle / ellipse shapes
  shape?: {
    radiusM: number
  }
  // for drone/UAS tracks
  speed?: number   // m/s
  course?: number  // degrees true
  battery?: number // 0–100
}

function isoNow(): string {
  return new Date().toISOString()
}

function isoStale(minutes: number): string {
  return new Date(Date.now() + minutes * 60_000).toISOString()
}

function escXml(s: string): string {
  return s
    .replace(/&/g, '&amp;')
    .replace(/</g, '&lt;')
    .replace(/>/g, '&gt;')
    .replace(/"/g, '&quot;')
    .replace(/'/g, '&apos;')
}

export function buildCot(ev: CotEvent): string {
  const stale = isoStale(ev.staleMinutes ?? 5)
  const now   = isoNow()
  const point = ev.point
  const hae   = (point.hae ?? 0).toFixed(1)
  const ce    = (point.ce  ?? 9999999).toString()
  const le    = (point.le  ?? 9999999).toString()

  let detail = `<contact callsign="${escXml(ev.callsign)}"/>`

  if (ev.remarks) {
    detail += `\n    <remarks>${escXml(ev.remarks)}</remarks>`
  }

  // Drone/UAS specific
  if (ev.speed !== undefined || ev.course !== undefined || ev.battery !== undefined) {
    detail += `\n    <track speed="${(ev.speed ?? 0).toFixed(1)}" course="${(ev.course ?? 0).toFixed(1)}"/>`
    if (ev.battery !== undefined) {
      detail += `\n    <status battery="${ev.battery}"/>`
    }
  }

  // Shape (geocircle for geofences)
  if (ev.shape) {
    const r = ev.shape.radiusM.toFixed(0)
    detail += `\n    <shape>
      <ellipse semiMajor="${r}" semiMinor="${r}" angle="0"/>
    </shape>`
    const fill   = ev.color      ?? 0x20FF0000  // semi-transparent red
    const stroke = ev.strokeColor ?? 0xFFFF0000  // solid red
    const weight = ev.strokeWeight ?? 3
    detail += `\n    <strokeColor value="${stroke}"/>`
    detail += `\n    <strokeWeight value="${weight}.0"/>`
    detail += `\n    <fillColor value="${fill}"/>`
  }

  return [
    `<?xml version="1.0" encoding="UTF-8"?>`,
    `<event version="2.0"`,
    `       uid="${escXml(ev.uid)}"`,
    `       type="${escXml(ev.type)}"`,
    `       how="${ev.how ?? 'm-g'}"`,
    `       time="${now}"`,
    `       start="${now}"`,
    `       stale="${stale}">`,
    `  <point lat="${point.lat.toFixed(7)}"`,
    `         lon="${point.lon.toFixed(7)}"`,
    `         hae="${hae}"`,
    `         ce="${ce}"`,
    `         le="${le}"/>`,
    `  <detail>`,
    `    ${detail}`,
    `  </detail>`,
    `</event>`,
  ].join('\n')
}

// ── Convenience factories ────────────────────────────────────────────────────

export function friendlyUAVCot(droneId: string, lat: number, lon: number, altM: number, speedMs: number, courseDeg: number, battery: number): string {
  return buildCot({
    uid: `UAV-${droneId}`,
    type: 'a-f-A-M-F-Q',   // friendly airborne military fixed-wing UAS
    callsign: droneId,
    point: { lat, lon, hae: altM, ce: 10, le: 10 },
    staleMinutes: 1,
    speed: speedMs,
    course: courseDeg,
    battery,
    remarks: `Simulated UAV | Alt ${altM.toFixed(0)}m | Spd ${speedMs.toFixed(0)}m/s`,
  })
}

export function hostileMarkerCot(uid: string, callsign: string, lat: number, lon: number): string {
  return buildCot({
    uid,
    type: 'a-h-G',   // hostile ground
    callsign,
    point: { lat, lon, hae: 0, ce: 50, le: 9999 },
    staleMinutes: 30,
    remarks: 'BREACH DETECTED — AI autonomous response triggered',
  })
}

export function geofenceCircleCot(zoneId: string, name: string, lat: number, lon: number, radiusM: number): string {
  return buildCot({
    uid: `GEOFENCE-${zoneId}`,
    type: 'u-d-c-c',   // geocircle shape
    how: 'h-e',
    callsign: `FENCE-${name}`,
    point: { lat, lon, hae: 0, ce: 9999999, le: 9999999 },
    staleMinutes: 1440,   // 24 hours — persists on map
    shape: { radiusM },
    color: 0x20FF0000,
    strokeColor: -65536,  // ATAK uses signed ARGB: 0xFFFF0000 = -65536 = solid red
    strokeWeight: 3,
    remarks: `AI Geofence: ${name} | r=${radiusM}m | breach triggers autonomous response`,
  })
}

export function waypointCot(wpId: string, name: string, lat: number, lon: number): string {
  return buildCot({
    uid: `WP-${wpId}`,
    type: 'b-m-p-s-m',  // waypoint / navigation point
    how: 'h-e',
    callsign: name,
    point: { lat, lon, hae: 0, ce: 5, le: 9999 },
    staleMinutes: 1440,
    remarks: `Nav waypoint: ${name}`,
  })
}

export function emergencyBeaconCot(uid: string, callsign: string, lat: number, lon: number): string {
  return buildCot({
    uid,
    type: 'b-a-o-tbl',  // emergency / 9-line
    callsign,
    point: { lat, lon, hae: 0, ce: 50, le: 9999 },
    staleMinutes: 60,
    remarks: 'EMERGENCY — AI-triggered alert',
  })
}
