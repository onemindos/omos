/**
 * home-assistant-1.0.0 — OpenClaw Skill
 *
 * General-purpose Home Assistant control via REST API.
 * Supports: lights, switches, locks, climate, scenes, automations, sensors.
 *
 * Env vars:
 *   HOME_ASSISTANT_URL   — e.g. http://homeassistant.local:8123
 *   HOME_ASSISTANT_TOKEN — Long-lived access token
 *
 * Tools: ha_control, ha_status, ha_state, ha_scene, ha_automation, ha_list
 */

"use strict";

const HA_URL = (process.env.HOME_ASSISTANT_URL || "http://homeassistant.local:8123").replace(/\/$/, "");
const HA_TOKEN = process.env.HOME_ASSISTANT_TOKEN;

function headers() {
  if (!HA_TOKEN) throw new Error("HOME_ASSISTANT_TOKEN not set");
  return {
    Authorization: `Bearer ${HA_TOKEN}`,
    "Content-Type": "application/json",
  };
}

async function haFetch(path, opts = {}) {
  const url = `${HA_URL}${path}`;
  const res = await fetch(url, { ...opts, headers: headers() });
  if (!res.ok) {
    const text = await res.text().catch(() => "");
    throw new Error(`HA API ${res.status}: ${text.slice(0, 200)}`);
  }
  return res.json();
}

// ─── ha_control ─────────────────────────────────────────────────────────────

async function ha_control({ action, entity_id, brightness, temperature }) {
  if (!entity_id) throw new Error("entity_id required");
  if (!action) throw new Error("action required");

  const domain = entity_id.split(".")[0];
  let service;
  const data = { entity_id };

  switch (action) {
    case "turn_on":
      service = "turn_on";
      if (brightness != null && domain === "light") data.brightness = brightness;
      if (temperature != null && domain === "climate") data.temperature = temperature;
      break;
    case "turn_off":
      service = "turn_off";
      break;
    case "toggle":
      service = "toggle";
      break;
    case "lock":
      return haFetch("/api/services/lock/lock", {
        method: "POST",
        body: JSON.stringify({ entity_id }),
      }).then(() => ({ action: "lock", entity_id, success: true }));
    case "unlock":
      return haFetch("/api/services/lock/unlock", {
        method: "POST",
        body: JSON.stringify({ entity_id }),
      }).then(() => ({ action: "unlock", entity_id, success: true }));
    default:
      throw new Error(`Unknown action: ${action}. Use turn_on, turn_off, toggle, lock, unlock.`);
  }

  await haFetch(`/api/services/${domain}/${service}`, {
    method: "POST",
    body: JSON.stringify(data),
  });

  return { action, entity_id, success: true };
}

// ─── ha_status ──────────────────────────────────────────────────────────────

async function ha_status({ domain } = {}) {
  const states = await haFetch("/api/states");

  if (domain) {
    const filtered = states
      .filter((s) => s.entity_id.startsWith(`${domain}.`))
      .map((s) => ({
        entity_id: s.entity_id,
        state: s.state,
        name: s.attributes?.friendly_name || s.entity_id,
      }));
    return { domain, count: filtered.length, entities: filtered.slice(0, 50) };
  }

  // Full summary
  let lightsOn = 0, lightsTotal = 0;
  let locksLocked = 0, locksTotal = 0;
  let doorsOpen = 0, doorsTotal = 0;
  const climate = [];

  for (const s of states) {
    const eid = s.entity_id;
    if (eid.startsWith("light.")) {
      lightsTotal++;
      if (s.state === "on") lightsOn++;
    } else if (eid.startsWith("lock.")) {
      locksTotal++;
      if (s.state === "locked") locksLocked++;
    } else if (eid.startsWith("binary_sensor.") && s.attributes?.device_class === "door") {
      doorsTotal++;
      if (s.state === "on") doorsOpen++;
    } else if (eid.startsWith("climate.")) {
      climate.push({
        entity_id: eid,
        name: s.attributes?.friendly_name || eid,
        state: s.state,
        current_temp: s.attributes?.current_temperature,
        target_temp: s.attributes?.temperature,
      });
    }
  }

  return {
    lights: { on: lightsOn, total: lightsTotal },
    locks: { locked: locksLocked, total: locksTotal },
    doors: { open: doorsOpen, total: doorsTotal },
    climate,
  };
}

// ─── ha_state ───────────────────────────────────────────────────────────────

async function ha_state({ entity_id }) {
  if (!entity_id) throw new Error("entity_id required");
  const s = await haFetch(`/api/states/${entity_id}`);
  return {
    entity_id: s.entity_id,
    state: s.state,
    name: s.attributes?.friendly_name || s.entity_id,
    unit: s.attributes?.unit_of_measurement || null,
    attributes: s.attributes || {},
    last_changed: s.last_changed,
  };
}

// ─── ha_scene ───────────────────────────────────────────────────────────────

async function ha_scene({ scene }) {
  if (!scene) throw new Error("scene name required");
  const sceneId = scene.startsWith("scene.") ? scene : `scene.${scene}`;
  await haFetch("/api/services/scene/turn_on", {
    method: "POST",
    body: JSON.stringify({ entity_id: sceneId }),
  });
  return { action: "scene", scene_id: sceneId, success: true };
}

// ─── ha_automation ──────────────────────────────────────────────────────────

async function ha_automation({ automation_id }) {
  if (!automation_id) throw new Error("automation_id required");
  const aid = automation_id.startsWith("automation.") ? automation_id : `automation.${automation_id}`;
  await haFetch("/api/services/automation/trigger", {
    method: "POST",
    body: JSON.stringify({ entity_id: aid }),
  });
  return { action: "trigger", automation_id: aid, success: true };
}

// ─── ha_list ────────────────────────────────────────────────────────────────

async function ha_list({ domain }) {
  if (!domain) throw new Error("domain required (light, lock, sensor, climate, switch, etc.)");
  const states = await haFetch("/api/states");
  const filtered = states
    .filter((s) => s.entity_id.startsWith(`${domain}.`))
    .map((s) => ({
      entity_id: s.entity_id,
      state: s.state,
      name: s.attributes?.friendly_name || s.entity_id,
    }));
  return { domain, count: filtered.length, entities: filtered.slice(0, 100) };
}

// ─── Exports ────────────────────────────────────────────────────────────────

module.exports = { ha_control, ha_status, ha_state, ha_scene, ha_automation, ha_list };
