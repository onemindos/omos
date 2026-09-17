"use strict";

/**
 * advanced-skill — OpenClaw Skill
 *
 * Full-featured skill template with primary tools, script execution,
 * and deep reference patterns.
 *
 * Registered tools: advanced_process, advanced_analyze, advanced_report
 */

const { execFile } = require("child_process");
const path = require("path");
const { promisify } = require("util");

const execFileAsync = promisify(execFile);

// ─── Configuration ───────────────────────────────────────────────────────────
const BASE_URL = (process.env.ADVANCED_SKILL_URL || "https://api.example.com").replace(/\/$/, "");
const API_KEY = process.env.ADVANCED_SKILL_KEY;

// Path to scripts bundled with this skill
const SCRIPTS_DIR = path.join(__dirname, "scripts");

/** Shared authenticated fetch wrapper */
async function apiRequest(path_segment, { method = "GET", body = null } = {}) {
  if (!API_KEY) throw new Error("ADVANCED_SKILL_KEY env var is not configured");

  const res = await fetch(`${BASE_URL}${path_segment}`, {
    method,
    headers: {
      Authorization: `Bearer ${API_KEY}`,
      "Content-Type": "application/json",
      Accept: "application/json",
    },
    body: body ? JSON.stringify(body) : undefined,
    signal: AbortSignal.timeout(30000),
  });

  if (!res.ok) {
    const text = await res.text().catch(() => `HTTP ${res.status}`);
    throw new Error(`API error ${res.status} on ${method} ${path_segment}: ${text}`);
  }

  const ct = res.headers.get("content-type") || "";
  if (!ct.includes("application/json")) return { status: "ok" };
  return res.json();
}

/** Run a bundled Python/Bash script safely */
async function runScript(scriptName, args = []) {
  const scriptPath = path.join(SCRIPTS_DIR, scriptName);
  const ispython = scriptName.endsWith(".py");
  const interpreter = ispython ? "python3" : "bash";

  try {
    const { stdout, stderr } = await execFileAsync(interpreter, [scriptPath, ...args], {
      timeout: 30000,
      maxBuffer: 1024 * 1024, // 1MB output limit
    });
    if (stderr) console.warn(`[advanced-skill] script stderr: ${stderr}`);
    return stdout.trim();
  } catch (err) {
    throw new Error(`Script ${scriptName} failed: ${err.message}`);
  }
}

// ─── Tools ───────────────────────────────────────────────────────────────────

/**
 * TOOL: advanced_process
 * Primary processing action for this skill.
 *
 * @param {string} input - Input data or identifier
 * @param {string} [mode="standard"] - Processing mode: "standard", "fast", "thorough"
 * @returns {Object} Processing result
 */
async function advanced_process({ input, mode = "standard" }) {
  if (!input) throw new Error("input is required");

  const validModes = ["standard", "fast", "thorough"];
  if (!validModes.includes(mode)) {
    throw new Error(`mode must be one of: ${validModes.join(", ")}`);
  }

  // Option A: API call
  const result = await apiRequest("/process", {
    method: "POST",
    body: { input, mode },
  });

  // Option B: Run a script (comment out API call above and uncomment this)
  // const output = await runScript("process.py", ["--input", input, "--mode", mode]);
  // const result = JSON.parse(output);

  return {
    status: "ok",
    mode,
    input,
    ...result,
  };
}

/**
 * TOOL: advanced_analyze
 * Analyze a target, optionally with deep analysis via script.
 *
 * @param {string} target - Target identifier or data
 * @param {string} [depth="summary"] - Analysis depth: "summary" or "full"
 * @returns {Object} Analysis result with insights
 */
async function advanced_analyze({ target, depth = "summary" }) {
  if (!target) throw new Error("target is required");

  if (depth === "full") {
    // Use the bundled script for deep analysis (may take longer)
    const raw = await runScript("analyze.py", ["--target", target, "--full"]);
    return { depth: "full", target, analysis: JSON.parse(raw) };
  }

  // Quick API-based summary
  const result = await apiRequest(`/analyze/${encodeURIComponent(target)}`);
  return { depth: "summary", target, ...result };
}

/**
 * TOOL: advanced_report
 * Generate a structured report.
 *
 * @param {string} [format="markdown"] - Output format: "markdown" or "json"
 * @param {string} [timeframe="24h"] - Time window: "1h", "24h", "7d", "30d"
 * @returns {Object} Report content and metadata
 */
async function advanced_report({ format = "markdown", timeframe = "24h" }) {
  const validFormats = ["markdown", "json"];
  const validTimeframes = ["1h", "24h", "7d", "30d"];

  if (!validFormats.includes(format)) {
    throw new Error(`format must be one of: ${validFormats.join(", ")}`);
  }
  if (!validTimeframes.includes(timeframe)) {
    throw new Error(`timeframe must be one of: ${validTimeframes.join(", ")}`);
  }

  const result = await apiRequest(`/report?format=${format}&timeframe=${timeframe}`);
  return {
    format,
    timeframe,
    generated_at: new Date().toISOString(),
    ...result,
  };
}

// ─── Exports ─────────────────────────────────────────────────────────────────
module.exports = { advanced_process, advanced_analyze, advanced_report };
