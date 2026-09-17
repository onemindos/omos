/**
 * transit-1.0.0 — OpenClaw Skill
 *
 * Multi-modal transit routing via Google Maps Directions API.
 * Plan trips by subway, bus, train, walking, driving, bicycling.
 * Find nearby transit stations.
 *
 * Env vars: GOOGLE_MAPS_API_KEY
 *
 * Registers tools: plan_trip, get_directions, nearby_transit
 */

"use strict";

const https = require("https");
const { URL } = require("url");

// ---------------------------------------------------------------------------
// HTTP helper
// ---------------------------------------------------------------------------

function httpsGet(urlStr) {
  return new Promise((resolve, reject) => {
    let parsed;
    try { parsed = new URL(urlStr); } catch { return reject(new Error(`Invalid URL: ${urlStr}`)); }

    const opts = {
      hostname: parsed.hostname,
      path: parsed.pathname + parsed.search,
      method: "GET",
      headers: { "User-Agent": "OpenClaw-Transit/1.0" },
      timeout: 15000,
    };

    const req = https.request(opts, (res) => {
      if (res.statusCode < 200 || res.statusCode >= 300) {
        return reject(new Error(`HTTP ${res.statusCode}`));
      }
      let data = "";
      res.setEncoding("utf8");
      res.on("data", (c) => (data += c));
      res.on("end", () => resolve(data));
    });
    req.on("timeout", () => { req.destroy(); reject(new Error("Request timed out")); });
    req.on("error", reject);
    req.end();
  });
}

function getApiKey() {
  return process.env.GOOGLE_MAPS_API_KEY || "";
}

// ---------------------------------------------------------------------------
// Helpers
// ---------------------------------------------------------------------------

function cleanHtml(html) {
  return (html || "").replace(/<[^>]+>/g, " ").replace(/\s+/g, " ").trim();
}

function parseStep(step) {
  const mode = (step.travel_mode || "UNKNOWN").toLowerCase();
  const result = {
    mode,
    duration: (step.duration || {}).text || "",
    distance: (step.distance || {}).text || "",
    instruction: cleanHtml(step.html_instructions),
  };

  if (mode === "transit") {
    const td = step.transit_details || {};
    const line = td.line || {};
    const vehicle = line.vehicle || {};
    result.transit = {
      type: vehicle.type || "",
      name: line.short_name || line.name || "",
      headsign: td.headsign || "",
      departure_stop: (td.departure_stop || {}).name || "",
      arrival_stop: (td.arrival_stop || {}).name || "",
      departure_time: (td.departure_time || {}).text || "",
      arrival_time: (td.arrival_time || {}).text || "",
      num_stops: td.num_stops || 0,
    };
  }
  return result;
}

function buildSummary(steps) {
  const parts = [];
  for (const s of steps) {
    if (s.mode === "transit" && s.transit) {
      const t = s.transit;
      const type = (t.type || "").toLowerCase();
      if (["subway", "metro", "heavy_rail"].includes(type)) {
        parts.push(`${t.name} train to ${t.headsign}`);
      } else if (type === "bus") {
        parts.push(`Bus ${t.name} to ${t.headsign}`);
      } else {
        parts.push(`${t.name} to ${t.headsign}`);
      }
    }
  }
  return parts.length ? parts.join(" → ") : "Walk to destination";
}

// ---------------------------------------------------------------------------
// Tool: plan_trip
// ---------------------------------------------------------------------------

async function plan_trip({ origin, destination, mode, arrive_by, depart_at }) {
  const key = getApiKey();
  if (!key) return JSON.stringify({ error: "GOOGLE_MAPS_API_KEY not set" });

  mode = mode || "transit";

  try {
    let url = `https://maps.googleapis.com/maps/api/directions/json`
      + `?origin=${encodeURIComponent(origin)}`
      + `&destination=${encodeURIComponent(destination)}`
      + `&mode=${encodeURIComponent(mode)}`
      + `&alternatives=true`
      + `&key=${encodeURIComponent(key)}`;

    if (mode === "transit") {
      if (arrive_by) {
        const ts = Math.floor(new Date(arrive_by).getTime() / 1000);
        if (!isNaN(ts)) url += `&arrival_time=${ts}`;
      } else if (depart_at && depart_at !== "now") {
        const ts = Math.floor(new Date(depart_at).getTime() / 1000);
        if (!isNaN(ts)) url += `&departure_time=${ts}`;
      } else {
        url += `&departure_time=now`;
      }
    }

    const raw = await httpsGet(url);
    const data = JSON.parse(raw);

    if (data.status !== "OK" || !data.routes || data.routes.length === 0) {
      return JSON.stringify({ error: `No routes found (${data.status})`, origin, destination });
    }

    const routes = data.routes.slice(0, 3).map((route) => {
      const leg = route.legs[0];
      const steps = (leg.steps || []).map(parseStep);
      const transitSteps = steps.filter((s) => s.mode === "transit");
      return {
        summary: route.summary || "",
        duration: (leg.duration || {}).text || "",
        duration_seconds: (leg.duration || {}).value || 0,
        distance: (leg.distance || {}).text || "",
        departure_time: (leg.departure_time || {}).text || "",
        arrival_time: (leg.arrival_time || {}).text || "",
        start_address: leg.start_address || "",
        end_address: leg.end_address || "",
        transfers: Math.max(0, transitSteps.length - 1),
        steps,
      };
    });

    const best = routes[0];
    return JSON.stringify({
      origin,
      destination,
      mode,
      best_route: {
        duration: best.duration,
        departure: best.departure_time,
        arrival: best.arrival_time,
        transfers: best.transfers,
      },
      summary: buildSummary(best.steps),
      routes,
      alternatives: routes.length - 1,
    }, null, 2);
  } catch (e) {
    return JSON.stringify({ error: e.message });
  }
}

// ---------------------------------------------------------------------------
// Tool: get_directions
// ---------------------------------------------------------------------------

async function get_directions({ origin, destination, mode }) {
  const key = getApiKey();
  if (!key) return JSON.stringify({ error: "GOOGLE_MAPS_API_KEY not set" });

  mode = mode || "transit";

  try {
    let url = `https://maps.googleapis.com/maps/api/directions/json`
      + `?origin=${encodeURIComponent(origin)}`
      + `&destination=${encodeURIComponent(destination)}`
      + `&mode=${encodeURIComponent(mode)}`
      + `&key=${encodeURIComponent(key)}`;

    if (mode === "transit") url += `&departure_time=now`;

    const raw = await httpsGet(url);
    const data = JSON.parse(raw);

    if (data.status !== "OK" || !data.routes || data.routes.length === 0) {
      return JSON.stringify({ error: "No route found" });
    }

    const leg = data.routes[0].legs[0];
    return JSON.stringify({
      origin: leg.start_address || origin,
      destination: leg.end_address || destination,
      mode,
      duration: (leg.duration || {}).text || "",
      distance: (leg.distance || {}).text || "",
      steps: (leg.steps || []).map(parseStep),
    }, null, 2);
  } catch (e) {
    return JSON.stringify({ error: e.message });
  }
}

// ---------------------------------------------------------------------------
// Tool: nearby_transit
// ---------------------------------------------------------------------------

async function nearby_transit({ location, radius }) {
  const key = getApiKey();
  if (!key) return JSON.stringify({ error: "GOOGLE_MAPS_API_KEY not set" });

  radius = radius || 1000;

  try {
    // Geocode first
    const geoUrl = `https://maps.googleapis.com/maps/api/geocode/json`
      + `?address=${encodeURIComponent(location)}&key=${encodeURIComponent(key)}`;
    const geoRaw = await httpsGet(geoUrl);
    const geoData = JSON.parse(geoRaw);

    if (!geoData.results || geoData.results.length === 0) {
      return JSON.stringify({ error: `Could not find location: ${location}` });
    }

    const { lat, lng } = geoData.results[0].geometry.location;

    // Find nearby transit stations
    const placesUrl = `https://maps.googleapis.com/maps/api/place/nearbysearch/json`
      + `?location=${lat},${lng}&radius=${radius}&type=transit_station`
      + `&key=${encodeURIComponent(key)}`;
    const placesRaw = await httpsGet(placesUrl);
    const placesData = JSON.parse(placesRaw);

    const stations = (placesData.results || []).slice(0, 10).map((p) => ({
      name: p.name,
      address: p.vicinity || "",
      types: p.types || [],
    }));

    return JSON.stringify({
      location,
      coordinates: { lat, lng },
      radius_meters: radius,
      stations,
    }, null, 2);
  } catch (e) {
    return JSON.stringify({ error: e.message });
  }
}

// ---------------------------------------------------------------------------
// Exports
// ---------------------------------------------------------------------------

module.exports = { plan_trip, get_directions, nearby_transit };
