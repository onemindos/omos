/**
 * vision-1.0.0 — OpenClaw Skill
 *
 * AI camera vision via Home Assistant LLM Vision + Frigate NVR.
 * Analyze camera feeds with GPT-4o/Claude/Gemini, query motion events,
 * get snapshots, and monitor live streams.
 *
 * Env vars: HOME_ASSISTANT_URL, HOME_ASSISTANT_TOKEN, FRIGATE_URL
 *
 * Registers tools: vision
 */

"use strict";

const https = require("https");
const http = require("http");
const { URL } = require("url");

// ---------------------------------------------------------------------------
// HTTP helpers
// ---------------------------------------------------------------------------

function request(method, urlStr, { headers = {}, data, timeout = 30000 } = {}) {
  return new Promise((resolve) => {
    let parsed;
    try {
      parsed = new URL(urlStr);
    } catch {
      return resolve({ error: `Invalid URL: ${urlStr}` });
    }

    const mod = parsed.protocol === "https:" ? https : http;

    const opts = {
      hostname: parsed.hostname,
      port: parsed.port || (parsed.protocol === "https:" ? 443 : 80),
      path: parsed.pathname + parsed.search,
      method: method.toUpperCase(),
      headers: { "User-Agent": "OpenClaw-Vision/1.0", ...headers },
      timeout,
    };

    const payload = data ? JSON.stringify(data) : null;
    if (payload) {
      opts.headers["Content-Type"] = "application/json";
      opts.headers["Content-Length"] = Buffer.byteLength(payload);
    }

    const req = mod.request(opts, (res) => {
      const chunks = [];
      res.on("data", (c) => chunks.push(c));
      res.on("end", () => {
        const buf = Buffer.concat(chunks);
        const code = res.statusCode;
        const ct = (res.headers["content-type"] || "").toLowerCase();

        // Binary response (snapshots) → base64
        if (ct.includes("image/")) {
          return resolve({
            _status: code,
            _binary: true,
            base64: buf.toString("base64"),
            content_type: res.headers["content-type"],
            size_bytes: buf.length,
          });
        }

        // JSON / text
        const text = buf.toString("utf8");
        if (code >= 200 && code < 300) {
          try {
            return resolve(JSON.parse(text));
          } catch {
            return resolve({ response: text, _status: code });
          }
        }
        let detail = text;
        try { detail = JSON.parse(text); } catch {}
        resolve({ error: `HTTP ${code}`, details: detail });
      });
    });
    req.on("timeout", () => { req.destroy(); resolve({ error: "Request timed out" }); });
    req.on("error", (e) => resolve({ error: e.message }));
    if (payload) req.write(payload);
    req.end();
  });
}

// ---------------------------------------------------------------------------
// Config helpers
// ---------------------------------------------------------------------------

function haUrl() {
  return (process.env.HOME_ASSISTANT_URL || "").replace(/\/+$/, "");
}

function haHeaders() {
  const token = process.env.HOME_ASSISTANT_TOKEN || "";
  if (!token) return null;
  return { Authorization: `Bearer ${token}`, "Content-Type": "application/json" };
}

function frigateUrl() {
  return (process.env.FRIGATE_URL || "").replace(/\/+$/, "");
}

// ---------------------------------------------------------------------------
// Tool: vision
// ---------------------------------------------------------------------------

async function vision({
  action,
  camera_name,
  prompt,
  provider,
  duration,
  event_id,
  label,
  hours_back,
  limit,
  quality,
}) {
  try {
    // =================================================================
    // LLM VISION — analyze camera feeds with AI
    // =================================================================

    if (action === "analyze_camera") {
      const ha = haUrl();
      const hd = haHeaders();
      if (!ha || !hd) return JSON.stringify({ error: "HOME_ASSISTANT_URL and HOME_ASSISTANT_TOKEN required" });
      if (!camera_name) return JSON.stringify({ error: "camera_name required" });

      const entity = camera_name.startsWith("camera.") ? camera_name : `camera.${camera_name}`;
      const svc = {
        image_entity: entity,
        message: prompt || "Describe what you see in detail. Include any people, vehicles, objects, or activity.",
        max_tokens: 500,
      };
      if (provider) svc.provider = provider;

      const res = await request("POST", `${ha}/api/services/llmvision/image_analyzer`, {
        headers: hd,
        data: svc,
        timeout: 60000,
      });

      if (res.error) {
        if (String(res.details).includes("404"))
          return JSON.stringify({ error: "LLM Vision not installed in Home Assistant", hint: "Install via HACS: https://llmvision.org" });
        return JSON.stringify(res, null, 2);
      }

      return JSON.stringify({
        camera: entity,
        analysis: res.response || "Analysis complete",
        provider: provider || "default",
        timestamp: new Date().toISOString(),
        success: true,
      }, null, 2);
    }

    if (action === "analyze_video_event") {
      const ha = haUrl();
      const hd = haHeaders();
      const fUrl = frigateUrl();
      if (!ha || !hd) return JSON.stringify({ error: "HOME_ASSISTANT_URL and HOME_ASSISTANT_TOKEN required" });
      if (!fUrl) return JSON.stringify({ error: "FRIGATE_URL required for video events" });
      if (!event_id) return JSON.stringify({ error: "event_id required" });

      const clipUrl = `${fUrl}/api/events/${encodeURIComponent(event_id)}/clip.mp4`;
      const svc = {
        video_file: clipUrl,
        message: prompt || "Describe what happened in this event. Include any people, vehicles, or notable activity.",
        max_tokens: 500,
      };

      const res = await request("POST", `${ha}/api/services/llmvision/video_analyzer`, {
        headers: hd,
        data: svc,
        timeout: 120000,
      });

      if (res.error) return JSON.stringify(res, null, 2);
      return JSON.stringify({
        event_id,
        analysis: res.response || "Analysis complete",
        timestamp: new Date().toISOString(),
        success: true,
      }, null, 2);
    }

    if (action === "watch_stream") {
      const ha = haUrl();
      const hd = haHeaders();
      if (!ha || !hd) return JSON.stringify({ error: "HOME_ASSISTANT_URL and HOME_ASSISTANT_TOKEN required" });
      if (!camera_name) return JSON.stringify({ error: "camera_name required" });

      const entity = camera_name.startsWith("camera.") ? camera_name : `camera.${camera_name}`;
      const dur = Math.min(duration || 10, 60);

      const svc = {
        image_entity: entity,
        message: prompt || "Describe any activity or changes you observe.",
        duration: dur,
        max_tokens: 500,
      };

      const res = await request("POST", `${ha}/api/services/llmvision/stream_analyzer`, {
        headers: hd,
        data: svc,
        timeout: (dur + 30) * 1000,
      });

      if (res.error) return JSON.stringify(res, null, 2);
      return JSON.stringify({
        camera: entity,
        duration: dur,
        analysis: res.response || "Stream analysis complete",
        timestamp: new Date().toISOString(),
        success: true,
      }, null, 2);
    }

    // =================================================================
    // FRIGATE — direct NVR access
    // =================================================================

    if (action === "get_snapshot") {
      const fUrl = frigateUrl();
      if (!fUrl) return JSON.stringify({ error: "FRIGATE_URL required" });
      if (!camera_name) return JSON.stringify({ error: "camera_name required" });

      const q = quality ? `?quality=${quality}` : "?quality=70";
      const res = await request("GET", `${fUrl}/api/${encodeURIComponent(camera_name)}/latest.jpg${q}`, {
        timeout: 10000,
      });

      if (res.error) {
        if (String(res.details).includes("404"))
          return JSON.stringify({ error: `Camera '${camera_name}' not found`, hint: "Use action 'list_cameras' to see available cameras" });
        return JSON.stringify(res, null, 2);
      }

      if (res._binary) {
        return JSON.stringify({
          camera: camera_name,
          timestamp: new Date().toISOString(),
          format: "jpeg",
          size_bytes: res.size_bytes,
          image_base64: res.base64,
          success: true,
        }, null, 2);
      }

      return JSON.stringify({ error: "Unexpected response from Frigate" });
    }

    if (action === "get_events") {
      const fUrl = frigateUrl();
      if (!fUrl) return JSON.stringify({ error: "FRIGATE_URL required" });

      const hb = hours_back || 4;
      const lim = limit || 10;
      const after = Math.floor(Date.now() / 1000 - hb * 3600);

      const params = [`after=${after}`, `limit=${lim}`, "has_snapshot=1"];
      if (camera_name) params.push(`camera=${encodeURIComponent(camera_name)}`);
      if (label) params.push(`label=${encodeURIComponent(label)}`);

      const res = await request("GET", `${fUrl}/api/events?${params.join("&")}`, {
        timeout: 10000,
      });

      if (res.error) return JSON.stringify(res, null, 2);

      const events = (Array.isArray(res) ? res : []).map((e) => ({
        id: e.id,
        camera: e.camera,
        label: e.label,
        confidence: e.top_score ? Math.round(e.top_score * 1000) / 10 : null,
        time: e.start_time ? new Date(e.start_time * 1000).toISOString() : null,
        end_time: e.end_time ? new Date(e.end_time * 1000).toISOString() : null,
        zones: e.zones || [],
        has_clip: e.has_clip || false,
        thumbnail: `${fUrl}/api/events/${e.id}/thumbnail.jpg`,
      }));

      return JSON.stringify({
        events,
        count: events.length,
        query: { camera: camera_name || "all", label: label || "all", hours_back: hb },
        success: true,
      }, null, 2);
    }

    if (action === "get_thumbnail") {
      const fUrl = frigateUrl();
      if (!fUrl) return JSON.stringify({ error: "FRIGATE_URL required" });
      if (!event_id) return JSON.stringify({ error: "event_id required" });

      const res = await request("GET", `${fUrl}/api/events/${encodeURIComponent(event_id)}/thumbnail.jpg`, {
        timeout: 10000,
      });

      if (res.error) return JSON.stringify(res, null, 2);
      if (res._binary) {
        return JSON.stringify({
          event_id,
          format: "jpeg",
          size_bytes: res.size_bytes,
          image_base64: res.base64,
          success: true,
        }, null, 2);
      }
      return JSON.stringify({ error: "Unexpected response" });
    }

    if (action === "list_cameras") {
      const fUrl = frigateUrl();
      if (!fUrl) return JSON.stringify({ error: "FRIGATE_URL required" });

      const configRes = await request("GET", `${fUrl}/api/config`, { timeout: 10000 });
      if (configRes.error) return JSON.stringify(configRes, null, 2);

      const statsRes = await request("GET", `${fUrl}/api/stats`, { timeout: 10000 });
      const stats = statsRes.error ? {} : statsRes;

      const cameras = Object.entries(configRes.cameras || {}).map(([name, cfg]) => {
        const cs = (stats.cameras || {})[name] || {};
        return {
          name,
          enabled: cfg.enabled !== false,
          detect: (cfg.detect || {}).enabled !== false,
          record: (cfg.record || {}).enabled === true,
          fps: cs.camera_fps || 0,
          zones: Object.keys(cfg.zones || {}),
          objects: (cfg.objects || {}).track || [],
        };
      });

      return JSON.stringify({ cameras, count: cameras.length, success: true }, null, 2);
    }

    if (action === "camera_status") {
      const fUrl = frigateUrl();
      if (!fUrl) return JSON.stringify({ error: "FRIGATE_URL required" });
      if (!camera_name) return JSON.stringify({ error: "camera_name required" });

      const res = await request("GET", `${fUrl}/api/stats`, { timeout: 10000 });
      if (res.error) return JSON.stringify(res, null, 2);

      const cs = (res.cameras || {})[camera_name];
      if (!cs) {
        return JSON.stringify({
          error: `Camera '${camera_name}' not found`,
          available: Object.keys(res.cameras || {}),
        });
      }

      return JSON.stringify({
        camera: camera_name,
        status: (cs.camera_fps || 0) > 0 ? "online" : "offline",
        camera_fps: cs.camera_fps || 0,
        detection_fps: cs.detection_fps || 0,
        process_fps: cs.process_fps || 0,
        success: true,
      }, null, 2);
    }

    // =================================================================
    // UNKNOWN ACTION
    // =================================================================

    return JSON.stringify({
      error: `Unknown action: ${action}`,
      available: [
        "analyze_camera",
        "analyze_video_event",
        "watch_stream",
        "get_snapshot",
        "get_events",
        "get_thumbnail",
        "list_cameras",
        "camera_status",
      ],
    });
  } catch (err) {
    return JSON.stringify({ error: err.message });
  }
}

module.exports = { vision };
