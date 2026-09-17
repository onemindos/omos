"use strict";

/**
 * basic-skill — OpenClaw Skill
 *
 * [Describe what this skill does in 1-2 sentences.]
 *
 * Registered tools: basic_action
 */

// ─── Configuration ───────────────────────────────────────────────────────────
// Load env vars at the top. Fail fast if critical vars are missing.
// const API_KEY = process.env.MY_SERVICE_API_KEY;
// const BASE_URL = process.env.MY_SERVICE_URL || "https://api.example.com";

// ─── Tools ───────────────────────────────────────────────────────────────────

/**
 * TOOL: basic_action
 * Does the primary action for this skill.
 *
 * @param {string} input - The primary input for this action
 * @returns {Object} result
 */
async function basic_action({ input }) {
  if (!input) throw new Error("input is required");

  // TODO: implement real logic here
  // Example:
  // const response = await fetch(`${BASE_URL}/endpoint`, {
  //   method: "POST",
  //   headers: { Authorization: `Bearer ${API_KEY}`, "Content-Type": "application/json" },
  //   body: JSON.stringify({ data: input }),
  //   signal: AbortSignal.timeout(30000),
  // });
  // if (!response.ok) {
  //   const err = await response.text();
  //   throw new Error(`API error ${response.status}: ${err}`);
  // }
  // return response.json();

  return {
    status: "ok",
    input,
    result: `Processed: ${input}`,
  };
}

// ─── Exports ─────────────────────────────────────────────────────────────────
// Only exported functions become OpenClaw tools.
module.exports = { basic_action };
