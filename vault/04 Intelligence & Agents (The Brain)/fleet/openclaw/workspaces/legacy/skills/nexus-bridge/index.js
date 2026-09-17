/**
 * nexus-bridge — OpenClaw Skill
 * 
 * Connects Legacy to OneMind OS NATS nervous system.
 * Enables dispatch to Agno agents and OM1 robots.
 * 
 * Registers tools: nexus_dispatch, nexus_query, nexus_publish
 */

"use strict";

const { connect, StringCodec, JSONCodec } = require("nats");

// Configuration from environment
const NATS_URL = process.env.NATS_URL || "nats://localhost:4222";
const AGNO_BASE_URL = process.env.AGNO_BASE_URL || "http://localhost:7777";
const OPENCLAW_HOOKS_URL = process.env.OPENCLAW_HOOKS_URL || "http://localhost:18789";
const OPENCLAW_HOOKS_TOKEN = process.env.OPENCLAW_HOOKS_TOKEN || "";

const sc = StringCodec();
const jc = JSONCodec();

/** Get or create NATS connection */
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
 * TOOL: nexus_dispatch
 * Dispatch a task to an Agno agent or OM1 robot via NATS.
 */
async function nexus_dispatch({ agent, robot_id, task, command, params, priority = "normal", context = {}, wait_for_result = false }) {
  const nc = await getNats();
  const timestamp = new Date().toISOString();

  if (robot_id) {
    // Physical dispatch → OM1 robot
    const subject = `physical.robot.${robot_id}.command`;
    const payload = {
      command: command || task,
      params: params || {},
      priority,
      context,
      source: "legacy",
      timestamp,
    };
    await nc.publish(subject, jc.encode(payload));
    return {
      status: "dispatched",
      target: "robot",
      robot_id,
      subject,
      timestamp,
    };
  }

  if (agent) {
    // Try Agno HTTP API for synchronous dispatch (better for wait_for_result)
    if (wait_for_result) {
      try {
        const response = await fetch(`${AGNO_BASE_URL}/agents/${agent}/runs`, {
          method: "POST",
          headers: { "Content-Type": "application/json" },
          body: JSON.stringify({
            message: task,
            stream: false,
            session_id: `legacy-dispatch-${Date.now()}`,
          }),
          signal: AbortSignal.timeout(60000),
        });
        if (response.ok) {
          const result = await response.json();
          return {
            status: "completed",
            target: "agent",
            agent_id: agent,
            result: result.content || result.message || result,
            timestamp,
          };
        }
      } catch (err) {
        // Fall through to NATS async dispatch
      }
    }

    // Async NATS dispatch
    const subject = `digital.agent.${agent}.task`;
    const payload = {
      task,
      priority,
      context: { ...context, source: "legacy", timestamp },
      callback_subject: `digital.agent.${agent}.result.legacy`,
    };
    await nc.publish(subject, jc.encode(payload));

    return {
      status: "dispatched",
      target: "agent",
      agent_id: agent,
      subject,
      note: wait_for_result
        ? "NATS async dispatch used (Agno API unreachable). Result will arrive via webhook."
        : "Task dispatched. Agent will process asynchronously.",
      timestamp,
    };
  }

  throw new Error("Either 'agent' or 'robot_id' must be specified");
}

/**
 * TOOL: nexus_query
 * Query system state via NATS request-reply.
 */
async function nexus_query({ subject, timeout_ms = 5000 }) {
  const nc = await getNats();

  try {
    const msg = await nc.request(subject, sc.encode("query"), { timeout: timeout_ms });
    let data;
    try {
      data = jc.decode(msg.data);
    } catch {
      data = sc.decode(msg.data);
    }
    return { subject, response: data, timestamp: new Date().toISOString() };
  } catch (err) {
    if (err.message?.includes("TIMEOUT") || err.code === "TIMEOUT") {
      // Check Agno health as fallback for system.health queries
      if (subject.includes("health")) {
        try {
          const response = await fetch(`${AGNO_BASE_URL}/health`, {
            signal: AbortSignal.timeout(3000),
          });
          const agnoHealth = response.ok ? await response.json() : { status: "unreachable" };
          return {
            subject,
            response: {
              nats: "timeout",
              agno: agnoHealth,
              note: "NATS query timed out; Agno HTTP health returned instead",
            },
            timestamp: new Date().toISOString(),
          };
        } catch {
          // ignore
        }
      }
      return {
        subject,
        response: null,
        error: "timeout — no response within " + timeout_ms + "ms",
        timestamp: new Date().toISOString(),
      };
    }
    throw new Error(`NATS query failed for ${subject}: ${err.message}`);
  }
}

/**
 * TOOL: nexus_publish
 * Publish a raw event to any NATS subject.
 */
async function nexus_publish({ subject, payload }) {
  const nc = await getNats();
  await nc.publish(subject, jc.encode({ ...payload, source: "legacy", timestamp: new Date().toISOString() }));
  return { status: "published", subject, timestamp: new Date().toISOString() };
}

// OpenClaw skill tool registration
module.exports = {
  tools: {
    nexus_dispatch,
    nexus_query,
    nexus_publish,
  },

  // Called on skill load — verify NATS connectivity
  async onLoad() {
    try {
      await getNats();
      console.log(`[nexus-bridge] Connected to NATS at ${NATS_URL}`);
    } catch (err) {
      console.warn(`[nexus-bridge] NATS not reachable at startup (${err.message}). Will retry on first tool call.`);
    }
  },

  // Called on skill unload — close NATS connection
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
