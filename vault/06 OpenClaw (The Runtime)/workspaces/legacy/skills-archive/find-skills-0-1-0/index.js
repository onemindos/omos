"use strict";

/**
 * find-skills — OpenClaw Skill
 *
 * Discover and install agent skills from the open skills ecosystem.
 * Wraps `npx skills find` and `npx skills add` CLI commands.
 *
 * Registered tools: skills_find, skills_add, skills_check
 */

const { execFile } = require("child_process");

function runNpx(args) {
  return new Promise((resolve, reject) => {
    execFile("npx", ["skills", ...args], { timeout: 60000 }, (err, stdout, stderr) => {
      if (err) return reject(new Error(`skills CLI error: ${stderr || err.message}`));
      resolve(stdout.trim());
    });
  });
}

/**
 * TOOL: skills_find
 * Search the open skills ecosystem for skills matching a query.
 *
 * @param {string} query - Search terms (e.g. "react performance", "pr review", "changelog")
 * @returns {Object} { query, results }
 */
async function skills_find({ query }) {
  if (!query) throw new Error("query is required");
  const output = await runNpx(["find", query]);
  return { query, results: output };
}

/**
 * TOOL: skills_add
 * Install a skill from the ecosystem.
 *
 * @param {string} package_name - Skill package to install (e.g. "owner/repo@skill-name")
 * @returns {Object} { package_name, output }
 */
async function skills_add({ package_name }) {
  if (!package_name) throw new Error("package_name is required");
  const output = await runNpx(["add", package_name]);
  return { package_name, output };
}

/**
 * TOOL: skills_check
 * Check installed skills for available updates.
 *
 * @returns {Object} { output }
 */
async function skills_check() {
  const output = await runNpx(["check"]);
  return { output };
}

module.exports = { skills_find, skills_add, skills_check };
