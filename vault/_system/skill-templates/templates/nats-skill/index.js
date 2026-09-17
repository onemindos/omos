"use strict";

/**
 * nats-skill — OpenClaw Skill
 *
 * Publishes messages to and makes requests on OneMind OS NATS subjects.
 * Template for any skill that needs cross-runtime communication.
 *
 * Registered tools: nats_publish, nats_request, nats_health
 */

const { connect, JSONCodec } = require("nats");

// ─── Configuration ───────────────────────────────────────────────────────────
const NATS_URL = process.env.NATS_URL || "nats://localhost:4222";
const jc = JSONCodec();

/** Singleton NATS connection — reconnects if dropped */
let _nc = null;

async function getNats() {
  if (_nc && !_nc.isClosed()) return _nc;
  try {
    _nc = await connect({ servers: NATS_URL, timeout: 4000 });
    return _nc;
  } catch (err) {
    throw new Error(`Cannot reach NATS at ${NATS_URL}: ${err.message}`);
  }
}

// ─── Tools ───────────────────────────────────────────────────────────────────

/**
 * TOOL: nats_publish
 * Publish a fire-and-forget message to a NATS subject.
 *
 * @param {string} subject - NATS subject (e.g., "digital.agent.researcher.task")
 * @param {Object} data - JSON-serializable payload
 * @returns {Object} { published: true, subject, timestamp }
 */
async function nats_publish({ subject, data = {} }) {
  if (!subject) throw new Error("subject is required");

  const nc = await getNats();
  const payload = {
    ...data,
    _meta: {
      source: "legacy",
      timestamp: new Date().toISOString(),
    },
  };

  nc.publish(subject, jc.encode(payload));

  return {
    published: true,
    subject,
    timestamp: payload._meta.timestamp,
  };
}

/**
 * TOOL: nats_request
 * Publish a message and wait for a reply (request/reply pattern).
 *
 * @param {string} subject - NATS subject to publish to
 * @param {Object} data - JSON-serializable payload
 * @param {number} [timeout_ms=5000] - How long to wait for reply (milliseconds)
 * @returns {Object} The reply payload
 */
async function nats_request({ subject, data = {}, timeout_ms = 5000 }) {
  if (!subject) throw new Error("subject is required");

  const nc = await getNats();
  const payload = {
    ...data,
    _meta: {
      source: "legacy",
      timestamp: new Date().toISOString(),
    },
  };

  try {
    const msg = await nc.request(subject, jc.encode(payload), { timeout: timeout_ms });
    return jc.decode(msg.data);
  } catch (err) {
    if (err.message?.includes("timeout")) {
      throw new Error(`No reply from "${subject}" within ${timeout_ms}ms. Is any service listening?`);
    }
    throw err;
  }
}

/**
 * TOOL: nats_health
 * Check the NATS connection health.
 *
 * @returns {Object} { connected: boolean, server: string, status: string }
 */
async function nats_health() {
  try {
    const nc = await getNats();
    const info = nc.info;
    return {
      connected: true,
      server: NATS_URL,
      server_id: info?.server_id || "unknown",
      version: info?.version || "unknown",
      status: "healthy",
    };
  } catch (err) {
    return {
      connected: false,
      server: NATS_URL,
      status: "unreachable",
      error: err.message,
    };
  }
}

// ─── Exports ─────────────────────────────────────────────────────────────────
module.exports = { nats_publish, nats_request, nats_health };
