/**
 * uber-1.0.0 — OpenClaw Skill
 *
 * Request and manage Uber rides via the Uber Ride Request API v1.2.
 * Fare estimates, ride requests, status tracking, cancellation.
 *
 * Env vars: UBER_CLIENT_ID, UBER_CLIENT_SECRET, UBER_SERVER_TOKEN,
 *           UBER_ACCESS_TOKEN, UBER_USE_PRODUCTION (default: false = sandbox)
 *
 * Registers tools: uber
 */

"use strict";

const https = require("https");
const { URL } = require("url");

// ---------------------------------------------------------------------------
// Configuration
// ---------------------------------------------------------------------------

const UBER_API_BASE = "https://api.uber.com/v1.2";
const UBER_SANDBOX_BASE = "https://sandbox-api.uber.com/v1.2";

function getBase() {
  return (process.env.UBER_USE_PRODUCTION || "").toLowerCase() === "true"
    ? UBER_API_BASE
    : UBER_SANDBOX_BASE;
}

function isSandbox() {
  return (process.env.UBER_USE_PRODUCTION || "").toLowerCase() !== "true";
}

// ---------------------------------------------------------------------------
// HTTP helpers
// ---------------------------------------------------------------------------

function makeRequest(method, path, { body, params, useServerToken } = {}) {
  return new Promise((resolve, reject) => {
    const token = useServerToken
      ? process.env.UBER_SERVER_TOKEN
      : process.env.UBER_ACCESS_TOKEN;

    if (!token) {
      const envVar = useServerToken ? "UBER_SERVER_TOKEN" : "UBER_ACCESS_TOKEN";
      return resolve({ error: `${envVar} not set. Run OAuth flow first.` });
    }

    const authHeader = useServerToken ? `Token ${token}` : `Bearer ${token}`;
    const base = getBase();

    let qs = "";
    if (params) {
      qs = "?" + Object.entries(params).map(([k, v]) => `${encodeURIComponent(k)}=${encodeURIComponent(v)}`).join("&");
    }

    const parsed = new URL(`${base}${path}${qs}`);
    const payload = body ? JSON.stringify(body) : null;

    const opts = {
      hostname: parsed.hostname,
      path: parsed.pathname + parsed.search,
      method: method.toUpperCase(),
      headers: {
        Authorization: authHeader,
        "Content-Type": "application/json",
        "Accept-Language": "en_US",
        "User-Agent": "OpenClaw-Uber/1.0",
      },
      timeout: 30000,
    };

    if (payload) opts.headers["Content-Length"] = Buffer.byteLength(payload);

    const req = https.request(opts, (res) => {
      let data = "";
      res.setEncoding("utf8");
      res.on("data", (c) => (data += c));
      res.on("end", () => {
        const code = res.statusCode;
        if (code >= 200 && code < 300) {
          return resolve(data ? JSON.parse(data) : { success: true });
        }
        const errMap = {
          401: "Unauthorized. Token may be expired.",
          403: "Forbidden. Check scopes.",
          404: "Not found.",
          409: "Conflict. May already have an active ride.",
          422: "Unprocessable. Check parameters.",
        };
        resolve({ error: errMap[code] || `API error: ${code}`, details: data });
      });
    });
    req.on("timeout", () => { req.destroy(); resolve({ error: "Request timed out" }); });
    req.on("error", (e) => resolve({ error: e.message }));
    if (payload) req.write(payload);
    req.end();
  });
}

// ---------------------------------------------------------------------------
// Tool: uber
// ---------------------------------------------------------------------------

async function uber({
  action, start_lat, start_lng, end_lat, end_lng,
  start_address, end_address, product_id, ride_id, seat_count,
}) {
  seat_count = seat_count || 1;

  try {
    // --- AUTH URL ---
    if (action === "auth_url") {
      const clientId = process.env.UBER_CLIENT_ID || "";
      const redirectUri = process.env.UBER_REDIRECT_URI || "http://localhost:8080/callback/uber";
      const scopes = "profile request places ride_widgets";
      const url = `https://login.uber.com/oauth/v2/authorize`
        + `?response_type=code&client_id=${encodeURIComponent(clientId)}`
        + `&redirect_uri=${encodeURIComponent(redirectUri)}`
        + `&scope=${encodeURIComponent(scopes)}`;
      return JSON.stringify({
        action: "auth_url", url,
        instructions: "Visit this URL to authorize Uber access. After authorizing, you'll get a code to exchange for a token.",
      }, null, 2);
    }

    // --- PRODUCTS ---
    if (action === "products") {
      if (!start_lat || !start_lng) return JSON.stringify({ error: "start_lat and start_lng required" });
      const result = await makeRequest("GET", "/products", {
        params: { latitude: start_lat, longitude: start_lng },
        useServerToken: true,
      });
      if (result.error) return JSON.stringify(result, null, 2);

      return JSON.stringify({
        action: "products",
        location: { lat: start_lat, lng: start_lng },
        count: (result.products || []).length,
        sandbox: isSandbox(),
        products: (result.products || []).map((p) => ({
          id: p.product_id,
          name: p.display_name,
          description: p.description,
          capacity: p.capacity,
        })),
      }, null, 2);
    }

    // --- ESTIMATE ---
    if (action === "estimate") {
      if (!start_lat || !start_lng || !end_lat || !end_lng) {
        return JSON.stringify({ error: "start_lat, start_lng, end_lat, end_lng required" });
      }

      const [priceResult, timeResult] = await Promise.all([
        makeRequest("GET", "/estimates/price", {
          params: { start_latitude: start_lat, start_longitude: start_lng, end_latitude: end_lat, end_longitude: end_lng, seat_count },
          useServerToken: true,
        }),
        makeRequest("GET", "/estimates/time", {
          params: { start_latitude: start_lat, start_longitude: start_lng },
          useServerToken: true,
        }),
      ]);

      const prices = priceResult.error ? [] : (priceResult.prices || []);
      const times = timeResult.error ? [] : (timeResult.times || []);
      const timeMap = {};
      for (const t of times) timeMap[t.product_id] = t.estimate;

      return JSON.stringify({
        action: "estimate",
        from: { lat: start_lat, lng: start_lng, address: start_address || null },
        to: { lat: end_lat, lng: end_lng, address: end_address || null },
        sandbox: isSandbox(),
        estimates: prices.map((p) => ({
          product_id: p.product_id,
          name: p.display_name,
          estimate: p.estimate,
          low: p.low_estimate,
          high: p.high_estimate,
          currency: p.currency_code,
          duration_minutes: Math.round((p.duration || 0) / 60 * 10) / 10,
          distance_miles: Math.round((p.distance || 0) * 10) / 10,
          pickup_eta_seconds: timeMap[p.product_id] || null,
          surge: p.surge_multiplier || 1.0,
        })),
      }, null, 2);
    }

    // --- REQUEST RIDE ---
    if (action === "request") {
      if (!start_lat || !start_lng || !end_lat || !end_lng) {
        return JSON.stringify({ error: "start_lat, start_lng, end_lat, end_lng required" });
      }

      const estimateBody = {
        start_latitude: start_lat, start_longitude: start_lng,
        end_latitude: end_lat, end_longitude: end_lng, seat_count,
      };
      if (product_id) estimateBody.product_id = product_id;

      const estResult = await makeRequest("POST", "/requests/estimate", { body: estimateBody });
      if (estResult.error) return JSON.stringify({ action: "request", error: estResult.error, step: "estimate" }, null, 2);

      const fareId = (estResult.fare || {}).fare_id;
      if (!fareId) return JSON.stringify({ action: "request", error: "Could not get fare_id", estimate: estResult }, null, 2);

      const reqBody = {
        fare_id: fareId,
        start_latitude: start_lat, start_longitude: start_lng,
        end_latitude: end_lat, end_longitude: end_lng, seat_count,
      };
      if (product_id) reqBody.product_id = product_id;

      const result = await makeRequest("POST", "/requests", { body: reqBody });
      if (result.error) return JSON.stringify({ action: "request", error: result.error, step: "request" }, null, 2);

      return JSON.stringify({
        action: "request", success: true, sandbox: isSandbox(),
        request_id: result.request_id, status: result.status,
        eta_minutes: result.eta ? Math.round(result.eta / 60 * 10) / 10 : null,
        surge: result.surge_multiplier || 1.0,
        driver: result.driver ? { name: result.driver.name, phone: result.driver.phone_number, rating: result.driver.rating } : null,
        vehicle: result.vehicle ? { make: result.vehicle.make, model: result.vehicle.model, license: result.vehicle.license_plate } : null,
      }, null, 2);
    }

    // --- STATUS ---
    if (action === "status") {
      const endpoint = ride_id ? `/requests/${encodeURIComponent(ride_id)}` : "/requests/current";
      const result = await makeRequest("GET", endpoint);
      if (result.error) return JSON.stringify({ action: "status", has_active_ride: false, error: result.error }, null, 2);

      return JSON.stringify({
        action: "status", has_active_ride: true,
        request_id: result.request_id, status: result.status,
        eta_minutes: result.eta ? Math.round(result.eta / 60 * 10) / 10 : null,
        driver: result.driver ? { name: result.driver.name, phone: result.driver.phone_number, rating: result.driver.rating } : null,
        vehicle: result.vehicle ? { make: result.vehicle.make, model: result.vehicle.model, license: result.vehicle.license_plate } : null,
      }, null, 2);
    }

    // --- CANCEL ---
    if (action === "cancel") {
      const endpoint = ride_id ? `/requests/${encodeURIComponent(ride_id)}` : "/requests/current";
      const result = await makeRequest("DELETE", endpoint);
      if (result.error) return JSON.stringify({ action: "cancel", success: false, error: result.error }, null, 2);
      return JSON.stringify({ action: "cancel", success: true, message: "Ride cancelled" }, null, 2);
    }

    // --- HISTORY ---
    if (action === "history") {
      const result = await makeRequest("GET", "/history", { params: { limit: 10 } });
      if (result.error) return JSON.stringify(result, null, 2);
      return JSON.stringify({
        action: "history",
        rides: (result.history || []).map((h) => ({
          request_id: h.request_id, status: h.status,
          start_time: h.start_time, end_time: h.end_time,
          distance_miles: h.distance,
        })),
      }, null, 2);
    }

    return JSON.stringify({
      error: `Unknown action: ${action}`,
      valid_actions: ["products", "estimate", "request", "status", "cancel", "history", "auth_url"],
    }, null, 2);
  } catch (e) {
    return JSON.stringify({ error: e.message });
  }
}

// ---------------------------------------------------------------------------
// Exports
// ---------------------------------------------------------------------------

module.exports = { uber };
