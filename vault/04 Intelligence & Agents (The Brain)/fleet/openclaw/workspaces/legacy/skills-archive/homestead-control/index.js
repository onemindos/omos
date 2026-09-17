/**
 * homestead-control — OpenClaw Skill
 *
 * Aggregation tools for the 6-building homestead edge grid.
 * Uses native OpenClaw `nodes` tool + nexus-bridge NATS queries.
 *
 * Registers tools: homestead_status, homestead_alert
 */

"use strict";

const { connect, StringCodec, JSONCodec } = require("nats");

const NATS_URL = process.env.NATS_URL || "nats://localhost:4222";

const sc = StringCodec();
const jc = JSONCodec();

// Canonical building registry
const BUILDINGS = [
  { id: "main-house", name: "Main House", sensors: ["motion", "temp", "door"], voice: true },
  { id: "office", name: "Office", sensors: ["motion", "temp", "power"], voice: true },
  { id: "dojo", name: "Dojo", sensors: ["motion", "temp", "humidity", "sound"], voice: true },
  { id: "greenhouse", name: "Greenhouse", sensors: ["temp", "humidity", "soil"], voice: false, zones: ["zone-a", "zone-b", "zone-c"] },
  { id: "guest-house", name: "Guest House", sensors: ["motion", "temp", "door"], voice: true },
  { id: "barn", name: "Barn", sensors: ["temp", "motion", "air-quality"], voice: true, stalls: ["stall-1", "stall-2", "stall-3"] },
];

// Alert thresholds
const THRESHOLDS = {
  "barn.temp.critical": 90,
  "barn.temp.warning": 80,
  "greenhouse.soil.critical": 25,
  "greenhouse.soil.warning": 35,
  "greenhouse.temp.critical": 100,
  "bedroom.temp.sleep": 76,
  "office.power.warning": 500,
  "dojo.sound.session": 70,
};

let _nc = null;
async function getNats() {
  if (_nc) return _nc;
  try {
    _nc = await connect({ servers: NATS_URL, timeout: 3000 });
    return _nc;
  } catch (err) {
    throw new Error(`NATS unavailable at ${NATS_URL}: ${err.message}`);
  }
}

/**
 * Query a single NATS subject with timeout.
 * Returns null on timeout instead of throwing.
 */
async function safeQuery(nc, subject, timeoutMs = 3000) {
  try {
    const msg = await nc.request(subject, sc.encode("query"), { timeout: timeoutMs });
    try {
      return jc.decode(msg.data);
    } catch {
      return sc.decode(msg.data);
    }
  } catch {
    return null;
  }
}

/**
 * TOOL: homestead_status
 *
 * Get aggregated status for one or all buildings.
 * Queries node status + WireClaw sensors in parallel.
 *
 * @param {string} [building] - Specific building ID, or omit for all 6.
 * @param {boolean} [include_sensors=true] - Include WireClaw sensor data.
 * @param {boolean} [include_diagnostics=false] - Include disk/uptime/network info.
 * @returns {object} Aggregated status object.
 */
async function homestead_status({ building, include_sensors = true, include_diagnostics = false } = {}) {
  const nc = await getNats();
  const timestamp = new Date().toISOString();
  const targets = building
    ? BUILDINGS.filter(b => b.id === building)
    : BUILDINGS;

  if (targets.length === 0) {
    throw new Error(`Unknown building: ${building}. Valid: ${BUILDINGS.map(b => b.id).join(", ")}`);
  }

  const results = await Promise.all(targets.map(async (bld) => {
    const status = {
      id: bld.id,
      name: bld.name,
      node: null,
      sensors: {},
    };

    // Query node heartbeat via NATS
    const heartbeat = await safeQuery(nc, `nexus.heartbeat.node.${bld.id}`, 2000);
    status.node = heartbeat || { status: "unknown", note: "No heartbeat response" };

    // Query WireClaw sensors
    if (include_sensors) {
      const sensorQueries = [];

      // Temperature — always query
      sensorQueries.push(
        safeQuery(nc, `wireclaw.${bld.id}.*.temp`, 2000)
          .then(data => { status.sensors.temp = data; })
      );

      // Building-specific sensors
      if (bld.sensors.includes("humidity")) {
        sensorQueries.push(
          safeQuery(nc, `wireclaw.${bld.id}.*.humidity`, 2000)
            .then(data => { status.sensors.humidity = data; })
        );
      }
      if (bld.sensors.includes("motion")) {
        sensorQueries.push(
          safeQuery(nc, `wireclaw.${bld.id}.*.motion`, 2000)
            .then(data => { status.sensors.motion = data; })
        );
      }
      if (bld.sensors.includes("door")) {
        sensorQueries.push(
          safeQuery(nc, `wireclaw.${bld.id}.*.door`, 2000)
            .then(data => { status.sensors.door = data; })
        );
      }
      if (bld.sensors.includes("soil") && bld.zones) {
        for (const zone of bld.zones) {
          sensorQueries.push(
            safeQuery(nc, `wireclaw.${bld.id}.${zone}.soil`, 2000)
              .then(data => {
                if (!status.sensors.soil) status.sensors.soil = {};
                status.sensors.soil[zone] = data;
              })
          );
        }
      }
      if (bld.sensors.includes("power")) {
        sensorQueries.push(
          safeQuery(nc, `wireclaw.${bld.id}.*.power`, 2000)
            .then(data => { status.sensors.power = data; })
        );
      }
      if (bld.sensors.includes("sound")) {
        sensorQueries.push(
          safeQuery(nc, `wireclaw.${bld.id}.*.sound`, 2000)
            .then(data => { status.sensors.sound = data; })
        );
      }
      if (bld.sensors.includes("air-quality")) {
        sensorQueries.push(
          safeQuery(nc, `wireclaw.${bld.id}.*.air-quality`, 2000)
            .then(data => { status.sensors.airQuality = data; })
        );
      }
      if (bld.stalls) {
        for (const stall of bld.stalls) {
          sensorQueries.push(
            safeQuery(nc, `wireclaw.${bld.id}.${stall}.temp`, 2000)
              .then(data => {
                if (!status.sensors.stalls) status.sensors.stalls = {};
                status.sensors.stalls[stall] = data;
              })
          );
        }
      }

      await Promise.all(sensorQueries);
    }

    return status;
  }));

  // Check thresholds and attach advisories
  const advisories = [];
  for (const r of results) {
    if (r.sensors.stalls) {
      for (const [stall, data] of Object.entries(r.sensors.stalls)) {
        if (data && typeof data.value === "number") {
          if (data.value >= THRESHOLDS["barn.temp.critical"]) {
            advisories.push({ building: r.id, severity: "critical", message: `${stall} temp ${data.value}°F — exceeds ${THRESHOLDS["barn.temp.critical"]}°F` });
          } else if (data.value >= THRESHOLDS["barn.temp.warning"]) {
            advisories.push({ building: r.id, severity: "warning", message: `${stall} temp ${data.value}°F — approaching critical` });
          }
        }
      }
    }
    if (r.sensors.soil) {
      for (const [zone, data] of Object.entries(r.sensors.soil)) {
        if (data && typeof data.value === "number") {
          if (data.value <= THRESHOLDS["greenhouse.soil.critical"]) {
            advisories.push({ building: r.id, severity: "high", message: `${zone} soil moisture ${data.value}% — needs irrigation NOW` });
          } else if (data.value <= THRESHOLDS["greenhouse.soil.warning"]) {
            advisories.push({ building: r.id, severity: "warning", message: `${zone} soil moisture ${data.value}% — irrigation may be needed soon` });
          }
        }
      }
    }
  }

  return {
    timestamp,
    buildings: results,
    advisories,
    summary: {
      total: results.length,
      online: results.filter(r => r.node && r.node.status === "online").length,
      advisories: advisories.length,
    },
  };
}

/**
 * TOOL: homestead_alert
 *
 * Publish a homestead alert to NATS for logging, notification, and agent routing.
 *
 * @param {string} building - Building ID where the alert originated.
 * @param {string} severity - critical | high | warning | info
 * @param {string} message - Human-readable alert message.
 * @param {string} [route_to] - Agent to route to: legacy | le | hp | ge
 * @param {object} [data] - Additional structured data.
 */
async function homestead_alert({ building, severity, message, route_to = "legacy", data = {} }) {
  const nc = await getNats();
  const timestamp = new Date().toISOString();
  const validBuildings = BUILDINGS.map(b => b.id);

  if (building && !validBuildings.includes(building)) {
    throw new Error(`Unknown building: ${building}. Valid: ${validBuildings.join(", ")}`);
  }

  const payload = {
    type: "homestead_alert",
    building,
    severity,
    message,
    route_to,
    data,
    source: "homestead-control",
    timestamp,
  };

  // Publish to alert stream
  const subject = `homestead.alert.${building}.${severity}`;
  await nc.publish(subject, jc.encode(payload));

  // Also publish to the routed agent's channel if not legacy (legacy sees all alerts)
  if (route_to && route_to !== "legacy") {
    const agentSubject = `comms.text.agent.${route_to}`;
    await nc.publish(agentSubject, jc.encode({
      type: "homestead_alert_forwarded",
      original: payload,
      note: `Forwarded from homestead-control: ${message}`,
    }));
  }

  return {
    status: "published",
    subject,
    severity,
    building,
    route_to,
    timestamp,
  };
}

// OpenClaw skill tool registration
module.exports = {
  tools: {
    homestead_status,
    homestead_alert,
  },

  // Return building registry for agent context
  get registry() {
    return BUILDINGS;
  },

  get thresholds() {
    return THRESHOLDS;
  },

  async onLoad() {
    try {
      await getNats();
      console.log(`[homestead-control] Connected to NATS at ${NATS_URL}`);
      console.log(`[homestead-control] Monitoring ${BUILDINGS.length} buildings`);
    } catch (err) {
      console.warn(`[homestead-control] NATS not reachable at startup (${err.message}). Will retry on first tool call.`);
    }
  },

  async onUnload() {
    if (_nc) {
      try {
        await _nc.drain();
        _nc = null;
      } catch {
        // ignore
      }
    }
  },
};
