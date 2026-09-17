"use strict";

/**
 * skill-creator — OpenClaw Skill
 *
 * Scaffolds new OpenClaw skills by creating the required directory structure:
 * SKILL.md (with frontmatter), package.json, and index.js stub.
 *
 * Registered tools: create_skill, list_skill_templates
 */

const fs = require("fs");
const path = require("path");

const SKILLS_DIR = path.resolve(__dirname, "..");

/**
 * TOOL: create_skill
 * Scaffold a new OpenClaw skill directory with all required files.
 *
 * @param {string} name - Skill name (e.g. "my-service")
 * @param {string} version - Skill version (e.g. "1.0.0")
 * @param {string} description - What the skill does (used in SKILL.md frontmatter + package.json)
 * @param {string[]} [tools] - Tool function names to scaffold (e.g. ["my_action", "my_query"])
 * @returns {Object} { skillDir, filesCreated }
 */
async function create_skill({ name, version = "1.0.0", description, tools = ["main_action"] }) {
  if (!name) throw new Error("name is required");
  if (!description) throw new Error("description is required");

  const dirName = `${name}-${version}`;
  const skillDir = path.join(SKILLS_DIR, dirName);

  if (fs.existsSync(skillDir)) throw new Error(`Skill directory already exists: ${dirName}`);
  fs.mkdirSync(skillDir, { recursive: true });

  // SKILL.md
  const skillMd = `---
name: ${name}
description: ${description}
---

# ${name.replace(/-/g, " ").replace(/\b\w/g, c => c.toUpperCase())}

## When to Use

${description}

## Available Tools

| Tool | Description |
|------|-------------|
${tools.map(t => `| \`${t}\` | TODO: describe |`).join("\n")}

## How to Use

${tools.map(t => `- \`${t}({ input: "..." })\``).join("\n")}
`;
  fs.writeFileSync(path.join(skillDir, "SKILL.md"), skillMd);

  // package.json
  const pkg = JSON.stringify({ name, version, description, main: "index.js", dependencies: {} }, null, 2) + "\n";
  fs.writeFileSync(path.join(skillDir, "package.json"), pkg);

  // index.js
  const toolFns = tools.map(t => `
/**
 * TOOL: ${t}
 * TODO: describe what this tool does.
 *
 * @param {string} input - Primary input
 * @returns {Object} result
 */
async function ${t}({ input }) {
  if (!input) throw new Error("input is required");
  // TODO: implement
  return { status: "ok", input, result: "not yet implemented" };
}`).join("\n");

  const indexJs = `"use strict";
${toolFns}

module.exports = { ${tools.join(", ")} };
`;
  fs.writeFileSync(path.join(skillDir, "index.js"), indexJs);

  return { skillDir: dirName, filesCreated: ["SKILL.md", "package.json", "index.js"] };
}

/**
 * TOOL: list_skill_templates
 * List available skill templates from the skill-templates directory.
 *
 * @returns {Object} { templates }
 */
async function list_skill_templates() {
  const templatesDir = path.resolve(__dirname, "..", "..", "..", "..", "skill-templates", "templates");
  if (!fs.existsSync(templatesDir)) return { templates: [], note: "skill-templates/templates/ not found" };
  const entries = fs.readdirSync(templatesDir, { withFileTypes: true });
  const templates = entries.filter(e => e.isDirectory()).map(e => e.name);
  return { templates };
}

module.exports = { create_skill, list_skill_templates };
