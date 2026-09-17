/**
 * nats-manager-1.0.0 — OpenClaw Skill
 *
 * NATS server monitoring + JetStream management.
 * Combines HTTP monitoring API (read-only stats) with
 * NATS client operations (stream/consumer management).
 *
 * Registers tools: nats_health, nats_connections, nats_streams,
 *   nats_stream_detail, nats_consumers, nats_subscriptions,
 *   nats_stream_create, nats_stream_purge
 */

"use strict";

const http = require("http");
const https = require("https");
const { URL } = require("url");
const { connect, JSONCodec, AckPolicy, RetentionPolicy, StorageType } = require("nats");

// ---------------------------------------------------------------------------
// Configuration
// ---------------------------------------------------------------------------

const NATS_URL = process.env.NATS_URL || "nats://localhost:4222";
const NATS_MONITOR_URL = process.env.NATS_MONITOR_URL || "http://localhost:8222";

const jc = JSONCodec();

// ---------------------------------------------------------------------------
// HTTP helper for NATS monitoring API
// ---------------------------------------------------------------------------

function monitorGet(path) {
  return new Promise((resolve, reject) => {
    const url = `${NATS_MONITOR_URL}${path}`;
    let parsed;
    try {
      parsed = new URL(url);
    } catch {
      return reject(new Error(`Invalid monitor URL: ${url}`));
    }

    const lib = parsed.protocol === "https:" ? https : http;
    const opts = {
      hostname: parsed.hostname,
      port: parsed.port || (parsed.protocol === "https:" ? 443 : 80),
      path: parsed.pathname + parsed.search,
      method: "GET",
      headers: { Accept: "application/json" },
      timeout: 10000,
    };

    const req = lib.request(opts, (res) => {
      let data = "";
      res.setEncoding("utf8");
      res.on("data", (c) => (data += c));
      res.on("end", () => {
        if (res.statusCode < 200 || res.statusCode >= 300) {
          return reject(new Error(`NATS monitor HTTP ${res.statusCode}: ${data.slice(0, 300)}`));
        }
        try {
          resolve(JSON.parse(data));
        } catch {
          resolve({ raw: data });
        }
      });
    });
    req.on("timeout", () => {
      req.destroy();
      reject(new Error("NATS monitor request timed out"));
    });
    req.on("error", (err) => reject(new Error(`Cannot reach NATS monitor at ${url}: ${err.message || err.code || "connection refused"}`)));
    req.end();
  });
}

// ---------------------------------------------------------------------------
// NATS client connection (for JetStream management)
// ---------------------------------------------------------------------------

let _nc = null;
async function getNats() {
  if (_nc && !_nc.isClosed()) return _nc;
  try {
    _nc = await connect({ servers: NATS_URL, timeout: 5000 });
    return _nc;
  } catch (err) {
    throw new Error(`NATS unavailable at ${NATS_URL}: ${err.message}`);
  }
}

function formatBytes(bytes) {
  if (bytes === 0) return "0 B";
  const units = ["B", "KiB", "MiB", "GiB", "TiB"];
  const i = Math.floor(Math.log(bytes) / Math.log(1024));
  return `${(bytes / Math.pow(1024, i)).toFixed(1)} ${units[i]}`;
}

function formatDuration(seconds) {
  if (!seconds || seconds <= 0) return "0s";
  const d = Math.floor(seconds / 86400);
  const h = Math.floor((seconds % 86400) / 3600);
  const m = Math.floor((seconds % 3600) / 60);
  const s = Math.floor(seconds % 60);
  const parts = [];
  if (d > 0) parts.push(`${d}d`);
  if (h > 0) parts.push(`${h}h`);
  if (m > 0) parts.push(`${m}m`);
  if (s > 0 || parts.length === 0) parts.push(`${s}s`);
  return parts.join(" ");
}

// ---------------------------------------------------------------------------
// TOOL: nats_health
// ---------------------------------------------------------------------------

async function nats_health() {
  try {
    const varz = await monitorGet("/varz");

    return JSON.stringify({
      server: {
        name: varz.server_name,
        id: varz.server_id,
        version: varz.version,
        go_version: varz.go,
        host: varz.host,
        port: varz.port,
        uptime: formatDuration(varz.now ? (Date.now() / 1000 - new Date(varz.start).getTime() / 1000) : 0),
        uptime_raw: varz.uptime,
      },
      resources: {
        cpu: `${(varz.cpu || 0).toFixed(1)}%`,
        memory: formatBytes(varz.mem || 0),
        goroutines: varz.routes || 0,
      },
      connections: {
        current: varz.connections || 0,
        total: varz.total_connections || 0,
      },
      messages: {
        in: varz.in_msgs || 0,
        out: varz.out_msgs || 0,
        in_bytes: formatBytes(varz.in_bytes || 0),
        out_bytes: formatBytes(varz.out_bytes || 0),
      },
      subscriptions: varz.subscriptions || 0,
      slow_consumers: varz.slow_consumers || 0,
      jetstream: varz.jetstream ? "enabled" : "disabled",
      leafnodes: {
        current: varz.leafnodes || 0,
      },
    }, null, 2);
  } catch (e) {
    return JSON.stringify({ error: e.message, hint: `Check NATS_MONITOR_URL (${NATS_MONITOR_URL})` });
  }
}

// ---------------------------------------------------------------------------
// TOOL: nats_connections
// ---------------------------------------------------------------------------

async function nats_connections({ sort = "msgs_to", limit = 50 } = {}) {
  try {
    const connz = await monitorGet(`/connz?sort=${sort}&limit=${limit}`);

    const connections = (connz.connections || []).map((c) => ({
      cid: c.cid,
      name: c.name || "(unnamed)",
      ip: c.ip,
      port: c.port,
      lang: c.lang,
      version: c.version,
      subscriptions: c.subscriptions || 0,
      msgs_to: c.in_msgs || 0,
      msgs_from: c.out_msgs || 0,
      bytes_to: formatBytes(c.in_bytes || 0),
      bytes_from: formatBytes(c.out_bytes || 0),
      idle: c.idle || "0s",
      uptime: c.uptime || "0s",
    }));

    return JSON.stringify({
      server: connz.server_id,
      total: connz.total || 0,
      showing: connections.length,
      sorted_by: sort,
      connections,
    }, null, 2);
  } catch (e) {
    return JSON.stringify({ error: e.message });
  }
}

// ---------------------------------------------------------------------------
// TOOL: nats_streams
// ---------------------------------------------------------------------------

async function nats_streams() {
  try {
    const nc = await getNats();
    const jsm = await nc.jetstreamManager();
    const streams = await jsm.streams.list().next();

    if (!streams || streams.length === 0) {
      return JSON.stringify({
        total: 0,
        streams: [],
        note: "No JetStream streams found. Use nats_stream_create to set up event channels.",
      }, null, 2);
    }

    const result = streams.map((s) => ({
      name: s.config.name,
      subjects: s.config.subjects,
      messages: s.state.messages,
      bytes: formatBytes(s.state.bytes),
      consumers: s.state.consumer_count,
      first_seq: s.state.first_seq,
      last_seq: s.state.last_seq,
      storage: s.config.storage === 2 ? "memory" : "file",
      retention: s.config.retention === 2 ? "interest" : s.config.retention === 3 ? "workqueue" : "limits",
      max_age: s.config.max_age ? formatDuration(s.config.max_age / 1e9) : "unlimited",
      max_bytes: s.config.max_bytes > 0 ? formatBytes(s.config.max_bytes) : "unlimited",
    }));

    return JSON.stringify({ total: result.length, streams: result }, null, 2);
  } catch (e) {
    return JSON.stringify({ error: e.message });
  }
}

// ---------------------------------------------------------------------------
// TOOL: nats_stream_detail
// ---------------------------------------------------------------------------

async function nats_stream_detail({ stream }) {
  if (!stream) return JSON.stringify({ error: "stream name is required" });

  try {
    const nc = await getNats();
    const jsm = await nc.jetstreamManager();
    const info = await jsm.streams.info(stream);

    // Get consumers
    const consumers = [];
    const consumerList = await jsm.consumers.list(stream).next();
    if (consumerList) {
      for (const c of consumerList) {
        consumers.push({
          name: c.name,
          durable: c.config.durable_name || null,
          deliver_subject: c.config.deliver_subject || "(pull)",
          ack_pending: c.num_ack_pending || 0,
          num_pending: c.num_pending || 0,
          delivered: c.delivered ? c.delivered.stream_seq : 0,
          filter: c.config.filter_subject || "(all)",
        });
      }
    }

    return JSON.stringify({
      name: info.config.name,
      description: info.config.description || "",
      subjects: info.config.subjects,
      config: {
        retention: info.config.retention === 2 ? "interest" : info.config.retention === 3 ? "workqueue" : "limits",
        storage: info.config.storage === 2 ? "memory" : "file",
        max_age: info.config.max_age ? formatDuration(info.config.max_age / 1e9) : "unlimited",
        max_bytes: info.config.max_bytes > 0 ? formatBytes(info.config.max_bytes) : "unlimited",
        max_msg_size: info.config.max_msg_size > 0 ? formatBytes(info.config.max_msg_size) : "unlimited",
        replicas: info.config.num_replicas,
        discard: info.config.discard === 1 ? "old" : "new",
        duplicates: info.config.duplicate_window ? formatDuration(info.config.duplicate_window / 1e9) : "2m",
      },
      state: {
        messages: info.state.messages,
        bytes: formatBytes(info.state.bytes),
        first_seq: info.state.first_seq,
        last_seq: info.state.last_seq,
        first_ts: info.state.first_ts,
        last_ts: info.state.last_ts,
        consumer_count: info.state.consumer_count,
      },
      consumers,
    }, null, 2);
  } catch (e) {
    return JSON.stringify({ error: e.message });
  }
}

// ---------------------------------------------------------------------------
// TOOL: nats_consumers
// ---------------------------------------------------------------------------

async function nats_consumers({ stream }) {
  if (!stream) return JSON.stringify({ error: "stream name is required" });

  try {
    const nc = await getNats();
    const jsm = await nc.jetstreamManager();
    const list = await jsm.consumers.list(stream).next();

    if (!list || list.length === 0) {
      return JSON.stringify({
        stream,
        total: 0,
        consumers: [],
        note: "No consumers on this stream.",
      }, null, 2);
    }

    const consumers = list.map((c) => ({
      name: c.name,
      durable: c.config.durable_name || null,
      deliver_subject: c.config.deliver_subject || "(pull-based)",
      filter: c.config.filter_subject || "(all subjects)",
      ack_policy: c.config.ack_policy === 1 ? "all" : c.config.ack_policy === 2 ? "explicit" : "none",
      delivered: {
        stream_seq: c.delivered ? c.delivered.stream_seq : 0,
        consumer_seq: c.delivered ? c.delivered.consumer_seq : 0,
      },
      ack_pending: c.num_ack_pending || 0,
      num_pending: c.num_pending || 0,
      num_redelivered: c.num_redelivered || 0,
      created: c.created,
    }));

    return JSON.stringify({ stream, total: consumers.length, consumers }, null, 2);
  } catch (e) {
    return JSON.stringify({ error: e.message });
  }
}

// ---------------------------------------------------------------------------
// TOOL: nats_subscriptions
// ---------------------------------------------------------------------------

async function nats_subscriptions({ test_subject } = {}) {
  try {
    const url = test_subject ? `/subsz?subs=1&test=${encodeURIComponent(test_subject)}` : "/subsz?subs=1";
    const subsz = await monitorGet(url);

    return JSON.stringify({
      total: subsz.num_subscriptions || 0,
      cache_hits: subsz.num_cache || 0,
      inserts: subsz.num_inserts || 0,
      removes: subsz.num_removes || 0,
      matches: subsz.num_matches || 0,
      test_subject: test_subject || null,
      subscriptions: (subsz.subscriptions_list || []).slice(0, 100),
    }, null, 2);
  } catch (e) {
    return JSON.stringify({ error: e.message });
  }
}

// ---------------------------------------------------------------------------
// TOOL: nats_stream_create
// ---------------------------------------------------------------------------

async function nats_stream_create({ stream, subjects, max_age_hours = 168, max_bytes = 104857600, storage = "file" }) {
  if (!stream) return JSON.stringify({ error: "stream name is required" });
  if (!subjects || !Array.isArray(subjects) || subjects.length === 0) {
    return JSON.stringify({ error: "subjects array is required (e.g. ['system.>', 'nexus.>'])" });
  }

  // Validate stream name (UPPER_CASE, alphanumeric + underscores)
  if (!/^[A-Z][A-Z0-9_]*$/.test(stream)) {
    return JSON.stringify({ error: "Stream name must be UPPER_CASE alphanumeric with underscores (e.g. SYSTEM, DIGITAL_AGENTS)" });
  }

  try {
    const nc = await getNats();
    const jsm = await nc.jetstreamManager();

    const config = {
      name: stream,
      subjects,
      retention: RetentionPolicy.Limits,
      storage: storage === "memory" ? StorageType.Memory : StorageType.File,
      max_age: max_age_hours * 3600 * 1e9, // nanoseconds
      max_bytes: max_bytes,
      num_replicas: 1,
      discard: 1, // old
      duplicate_window: 120 * 1e9, // 2 minutes
    };

    const info = await jsm.streams.add(config);

    return JSON.stringify({
      status: "created",
      stream: info.config.name,
      subjects: info.config.subjects,
      max_age: formatDuration(max_age_hours * 3600),
      max_bytes: formatBytes(max_bytes),
      storage,
    }, null, 2);
  } catch (e) {
    if (e.message && e.message.includes("already in use")) {
      return JSON.stringify({
        status: "exists",
        stream,
        note: "Stream already exists. Use nats_stream_detail to inspect it.",
      }, null, 2);
    }
    return JSON.stringify({ error: e.message });
  }
}

// ---------------------------------------------------------------------------
// TOOL: nats_stream_purge
// ---------------------------------------------------------------------------

async function nats_stream_purge({ stream }) {
  if (!stream) return JSON.stringify({ error: "stream name is required" });

  try {
    const nc = await getNats();
    const jsm = await nc.jetstreamManager();

    // Get current state before purge
    const before = await jsm.streams.info(stream);
    const msgCount = before.state.messages;
    const byteCount = before.state.bytes;

    await jsm.streams.purge(stream);

    return JSON.stringify({
      status: "purged",
      stream,
      messages_removed: msgCount,
      bytes_freed: formatBytes(byteCount),
      warning: "All messages permanently deleted.",
    }, null, 2);
  } catch (e) {
    return JSON.stringify({ error: e.message });
  }
}

// ---------------------------------------------------------------------------
// Exports
// ---------------------------------------------------------------------------

module.exports = { nats_health, nats_connections, nats_streams, nats_stream_detail, nats_consumers, nats_subscriptions, nats_stream_create, nats_stream_purge };
