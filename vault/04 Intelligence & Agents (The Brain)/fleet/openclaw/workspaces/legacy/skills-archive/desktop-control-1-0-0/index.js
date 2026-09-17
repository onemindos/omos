"use strict";

/**
 * desktop-control — OpenClaw Skill
 *
 * Desktop automation via Python bridge (pyautogui + pillow).
 * Provides mouse control, keyboard input, screenshots, and window management.
 * Shells out to the bundled ai_agent.py for actual desktop operations.
 *
 * Registered tools: desktop_mouse, desktop_keyboard, desktop_screenshot, desktop_window
 */

const { execFile } = require("child_process");
const path = require("path");

const PYTHON = process.env.PYTHON_PATH || "python3";
const AGENT_SCRIPT = path.join(__dirname, "ai_agent.py");

function runPython(action, args) {
  return new Promise((resolve, reject) => {
    const payload = JSON.stringify({ action, ...args });
    execFile(PYTHON, [AGENT_SCRIPT, "--json", payload], { timeout: 30000 }, (err, stdout, stderr) => {
      if (err) return reject(new Error(`desktop-control error: ${stderr || err.message}`));
      try {
        resolve(JSON.parse(stdout));
      } catch {
        resolve({ output: stdout.trim() });
      }
    });
  });
}

/**
 * TOOL: desktop_mouse
 * Control mouse movement and clicks.
 *
 * @param {string} action - One of: move, click, double_click, right_click, drag, scroll
 * @param {number} [x] - X coordinate
 * @param {number} [y] - Y coordinate
 * @param {number} [toX] - Drag destination X
 * @param {number} [toY] - Drag destination Y
 * @param {number} [scrollAmount] - Scroll amount (positive=up, negative=down)
 * @returns {Object} result
 */
async function desktop_mouse({ action, x, y, toX, toY, scrollAmount }) {
  if (!action) throw new Error("action is required (move, click, double_click, right_click, drag, scroll)");
  return runPython("mouse", { action, x, y, toX, toY, scrollAmount });
}

/**
 * TOOL: desktop_keyboard
 * Type text, press keys, or execute keyboard shortcuts.
 *
 * @param {string} action - One of: type, press, hotkey
 * @param {string} [text] - Text to type (for action=type)
 * @param {string} [key] - Key name (for action=press, e.g. "enter", "tab")
 * @param {string[]} [keys] - Key combo (for action=hotkey, e.g. ["ctrl", "c"])
 * @returns {Object} result
 */
async function desktop_keyboard({ action, text, key, keys }) {
  if (!action) throw new Error("action is required (type, press, hotkey)");
  return runPython("keyboard", { action, text, key, keys });
}

/**
 * TOOL: desktop_screenshot
 * Capture a screenshot of the current screen or a region.
 *
 * @param {number} [x] - Region top-left X
 * @param {number} [y] - Region top-left Y
 * @param {number} [width] - Region width
 * @param {number} [height] - Region height
 * @returns {Object} { base64, mimeType }
 */
async function desktop_screenshot({ x, y, width, height } = {}) {
  return runPython("screenshot", { x, y, width, height });
}

/**
 * TOOL: desktop_window
 * List windows or activate/minimize/maximize a specific window.
 *
 * @param {string} action - One of: list, activate, minimize, maximize
 * @param {string} [title] - Window title to target (for activate/minimize/maximize)
 * @returns {Object} result
 */
async function desktop_window({ action, title }) {
  if (!action) throw new Error("action is required (list, activate, minimize, maximize)");
  return runPython("window", { action, title });
}

module.exports = { desktop_mouse, desktop_keyboard, desktop_screenshot, desktop_window };
