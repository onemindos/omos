/**
 * skill-forge-1.0.0 — OpenClaw Skill
 *
 * Autonomous skill builder. Research APIs, generate OpenClaw skills,
 * and self-deploy them into the workspace. The MCP Builder pattern
 * adapted for OpenClaw's native skill system.
 *
 * Registers tools: skill_forge_create, skill_forge_validate,
 *                  skill_forge_list, skill_forge_from_openapi
 */

"use strict";

const fs = require("fs");
const path = require("path");
const https = require("https");
const http = require("http");
const { URL } = require("url");

// ---------------------------------------------------------------------------
// Constants
// ---------------------------------------------------------------------------

const FORGE_MARKER = ".forge-created";
const SKILLS_DIR = path.resolve(__dirname, "..");

// ---------------------------------------------------------------------------
// HTTP helper (GET only, for fetching OpenAPI specs)
// ---------------------------------------------------------------------------

function httpGet(urlStr, maxRedirects = 3) {
  return new Promise((resolve, reject) => {
    if (maxRedirects < 0) return reject(new Error("Too many redirects"));

    let parsed;
    try {
      parsed = new URL(urlStr);
    } catch {
      return reject(new Error(`Invalid URL: ${urlStr}`));
    }

    const lib = parsed.protocol === "https:" ? https : http;
    const opts = {
      hostname: parsed.hostname,
      port: parsed.port,
      path: parsed.pathname + parsed.search,
      method: "GET",
      headers: {
        "User-Agent": "OpenClaw-SkillForge/1.0",
        Accept: "application/json, text/plain, */*",
      },
      timeout: 30000,
    };

    const req = lib.request(opts, (res) => {
      if (res.statusCode >= 300 && res.statusCode < 400 && res.headers.location) {
        return resolve(httpGet(res.headers.location, maxRedirects - 1));
      }
      if (res.statusCode < 200 || res.statusCode >= 300) {
        return reject(new Error(`HTTP ${res.statusCode} fetching ${urlStr}`));
      }
      let data = "";
      res.setEncoding("utf8");
      res.on("data", (c) => (data += c));
      res.on("end", () => resolve(data));
    });
    req.on("timeout", () => {
      req.destroy();
      reject(new Error("Request timed out"));
    });
    req.on("error", reject);
    req.end();
  });
}

// ---------------------------------------------------------------------------
// Sanitization helpers
// ---------------------------------------------------------------------------

/** Enforce kebab-case, strip dangerous characters */
function sanitizeName(name) {
  return name
    .toLowerCase()
    .replace(/[^a-z0-9-]/g, "-")
    .replace(/-+/g, "-")
    .replace(/^-|-$/g, "")
    .slice(0, 60);
}

/** Enforce snake_case for tool function names */
function sanitizeToolName(name) {
  return name
    .toLowerCase()
    .replace(/[^a-z0-9_]/g, "_")
    .replace(/_+/g, "_")
    .replace(/^_|_$/g, "")
    .slice(0, 60);
}

/** Escape a string for safe inclusion in generated JS */
function escapeJS(str) {
  if (typeof str !== "string") return "";
  return str.replace(/\\/g, "\\\\").replace(/"/g, '\\"').replace(/\n/g, "\\n");
}

/** Validate URL string */
function isValidUrl(str) {
  try {
    const u = new URL(str);
    return u.protocol === "http:" || u.protocol === "https:";
  } catch {
    return false;
  }
}

// ---------------------------------------------------------------------------
// Code generators
// ---------------------------------------------------------------------------

function generateAuthCode(authType, authEnvVar, authHeader) {
  const envVar = authEnvVar || "API_KEY";

  switch (authType) {
    case "bearer":
      return {
        envLine: `const API_TOKEN = process.env.${envVar} || "";`,
        headerLine: `"Authorization": \`Bearer \${API_TOKEN}\``,
        envDoc: `- \`${envVar}\` — Bearer token for API authentication`,
      };
    case "api-key-header": {
      const hdr = authHeader || "X-API-Key";
      return {
        envLine: `const API_KEY = process.env.${envVar} || "";`,
        headerLine: `"${hdr}": API_KEY`,
        envDoc: `- \`${envVar}\` — API key sent via \`${hdr}\` header`,
      };
    }
    case "api-key-query":
      return {
        envLine: `const API_KEY = process.env.${envVar} || "";`,
        headerLine: null, // handled in query string
        queryParam: `api_key=\${encodeURIComponent(API_KEY)}`,
        envDoc: `- \`${envVar}\` — API key sent as query parameter`,
      };
    case "basic":
      return {
        envLine: `const API_CREDENTIALS = process.env.${envVar} || ""; // user:pass`,
        headerLine: `"Authorization": \`Basic \${Buffer.from(API_CREDENTIALS).toString("base64")}\``,
        envDoc: `- \`${envVar}\` — Basic auth credentials in \`user:password\` format`,
      };
    case "none":
    default:
      return {
        envLine: "// No authentication required",
        headerLine: null,
        envDoc: "No authentication required.",
      };
  }
}

function generateToolFunction(tool, baseUrl, authType, authEnvVar, authHeader, extraHeaders) {
  const name = sanitizeToolName(tool.name);
  const method = (tool.method || "GET").toUpperCase();
  const toolPath = tool.path || "/";
  const params = tool.parameters || [];

  const pathParams = params.filter((p) => p.in === "path");
  const queryParams = params.filter((p) => p.in === "query");
  const bodyParams = params.filter((p) => p.in === "body");
  const allParamNames = params.map((p) => p.name);

  const auth = generateAuthCode(authType, authEnvVar, authHeader);

  // Build function signature
  const destructured = allParamNames.length > 0 ? `{ ${allParamNames.join(", ")} }` : "{}";

  // Build URL with path params
  let urlBuild = `let endpoint = \`${toolPath}\`;\n`;
  for (const pp of pathParams) {
    urlBuild += `    endpoint = endpoint.replace("{${pp.name}}", encodeURIComponent(${pp.name}));\n`;
  }

  // Build query string
  let queryBuild = "";
  if (queryParams.length > 0 || auth.queryParam) {
    queryBuild += `    const qs = new URLSearchParams();\n`;
    for (const qp of queryParams) {
      queryBuild += `    if (${qp.name} !== undefined && ${qp.name} !== null) qs.set("${qp.name}", String(${qp.name}));\n`;
    }
    if (auth.queryParam) {
      queryBuild += `    qs.set("api_key", API_KEY);\n`;
    }
    queryBuild += `    const query = qs.toString();\n`;
    queryBuild += `    if (query) endpoint += (endpoint.includes("?") ? "&" : "?") + query;\n`;
  }

  // Build body
  let bodyBuild = "";
  if (bodyParams.length > 0 && ["POST", "PUT", "PATCH"].includes(method)) {
    const bodyObj = bodyParams.map((p) => `"${p.name}": ${p.name}`).join(", ");
    bodyBuild = `    const body = JSON.stringify({ ${bodyObj} });\n`;
  } else if (["POST", "PUT", "PATCH"].includes(method) && bodyParams.length === 0) {
    // Allow passing through any non-path/query params as body
    bodyBuild = `    const body = null;\n`;
  }

  // Headers
  const headers = [`"Content-Type": "application/json"`];
  if (auth.headerLine) headers.push(auth.headerLine);
  if (extraHeaders) {
    for (const [k, v] of Object.entries(extraHeaders)) {
      headers.push(`"${escapeJS(k)}": "${escapeJS(v)}"`);
    }
  }

  const hasBody = ["POST", "PUT", "PATCH"].includes(method);

  return `
async function ${name}(${destructured}) {
  try {
    ${urlBuild}${queryBuild}${hasBody && bodyParams.length > 0 ? bodyBuild : ""}
    const url = BASE_URL + endpoint;
    const response = await apiRequest(url, "${method}"${hasBody && bodyParams.length > 0 ? ", body" : ""});
    return JSON.stringify(response, null, 2);
  } catch (e) {
    return JSON.stringify({ error: e.message });
  }
}`;
}

function generateIndexJS(skillDef) {
  const {
    name,
    description,
    base_url,
    auth_type = "bearer",
    auth_env_var,
    auth_header,
    tools = [],
    extra_headers = {},
  } = skillDef;

  const auth = generateAuthCode(auth_type, auth_env_var, auth_header);

  // Build headers object
  const defaultHeaders = [`"Content-Type": "application/json"`, `"User-Agent": "OpenClaw-${sanitizeName(name)}/1.0"`];
  if (auth.headerLine) defaultHeaders.push(auth.headerLine);
  for (const [k, v] of Object.entries(extra_headers)) {
    defaultHeaders.push(`"${escapeJS(k)}": "${escapeJS(v)}"`);
  }

  const toolFunctions = tools.map((t) =>
    generateToolFunction(t, base_url, auth_type, auth_env_var, auth_header, extra_headers)
  );

  const exportNames = tools.map((t) => sanitizeToolName(t.name));

  return `/**
 * ${sanitizeName(name)}-${skillDef.version || "1.0.0"} — OpenClaw Skill (auto-generated by skill-forge)
 *
 * ${escapeJS(description)}
 *
 * Generated: ${new Date().toISOString()}
 * Registers tools: ${exportNames.join(", ")}
 */

"use strict";

const https = require("https");
const http = require("http");
const { URL } = require("url");

// ---------------------------------------------------------------------------
// Configuration
// ---------------------------------------------------------------------------

const BASE_URL = "${escapeJS(base_url)}";
${auth.envLine}

// ---------------------------------------------------------------------------
// HTTP helper
// ---------------------------------------------------------------------------

function apiRequest(urlStr, method, body) {
  return new Promise((resolve, reject) => {
    let parsed;
    try { parsed = new URL(urlStr); } catch { return reject(new Error(\`Invalid URL: \${urlStr}\`)); }

    const lib = parsed.protocol === "https:" ? https : http;
    const opts = {
      hostname: parsed.hostname,
      port: parsed.port || (parsed.protocol === "https:" ? 443 : 80),
      path: parsed.pathname + parsed.search,
      method,
      headers: {
        ${defaultHeaders.join(",\n        ")},
      },
      timeout: 30000,
    };

    if (body) opts.headers["Content-Length"] = Buffer.byteLength(body);

    const req = lib.request(opts, (res) => {
      let data = "";
      res.setEncoding("utf8");
      res.on("data", (c) => (data += c));
      res.on("end", () => {
        if (res.statusCode < 200 || res.statusCode >= 300) {
          return reject(new Error(\`HTTP \${res.statusCode}: \${data.slice(0, 500)}\`));
        }
        try { resolve(JSON.parse(data)); }
        catch { resolve({ raw: data }); }
      });
    });
    req.on("timeout", () => { req.destroy(); reject(new Error("Request timed out")); });
    req.on("error", reject);
    if (body) req.write(body);
    req.end();
  });
}

// ---------------------------------------------------------------------------
// Tools
// ---------------------------------------------------------------------------
${toolFunctions.join("\n")}

// ---------------------------------------------------------------------------
// Exports
// ---------------------------------------------------------------------------

module.exports = { ${exportNames.join(", ")} };
`;
}

function generateSKILLMD(skillDef) {
  const { name, description, tools = [], auth_type, auth_env_var } = skillDef;
  const auth = generateAuthCode(auth_type, auth_env_var);

  const toolDocs = tools
    .map((t) => {
      const params = (t.parameters || []).map((p) => p.name).join(", ");
      return `- \`${sanitizeToolName(t.name)}(${params})\` — ${t.description || ""}`;
    })
    .join("\n");

  return `---
name: ${sanitizeName(name)}
description: "${escapeJS(description)}"
metadata: { "forge": true, "created": "${new Date().toISOString().slice(0, 10)}" }
---

# ${name
    .split("-")
    .map((w) => w.charAt(0).toUpperCase() + w.slice(1))
    .join(" ")}

${description}

## Tools

${toolDocs}

## Authentication

${auth.envDoc}

## Examples

${tools
  .slice(0, 3)
  .map((t) => `- "${t.description || `Use ${t.name}`}"`)
  .join("\n")}

---

*Auto-generated by skill-forge on ${new Date().toISOString().slice(0, 10)}*
`;
}

function generatePackageJSON(skillDef) {
  const { name, version = "1.0.0", description, tools = [], auth_env_var, npm_dependencies = {} } = skillDef;
  const safeName = sanitizeName(name);

  const toolDefs = tools.map((t) => {
    const params = (t.parameters || []).map((p) => ({
      name: p.name,
      type: p.type || "string",
      description: p.description || "",
      ...(p.required !== undefined ? { required: p.required } : {}),
    }));

    const required = params.filter((p) => p.required).map((p) => p.name);
    const properties = {};
    for (const p of params) {
      properties[p.name] = { type: p.type, description: p.description };
    }

    return {
      name: sanitizeToolName(t.name),
      description: t.description || "",
      parameters: {
        type: "object",
        properties,
        ...(required.length > 0 ? { required } : {}),
      },
    };
  });

  return JSON.stringify(
    {
      name: `${safeName}-${version}`,
      version,
      description,
      main: "index.js",
      tools: toolDefs,
      env: auth_env_var ? [auth_env_var] : [],
      dependencies: npm_dependencies,
    },
    null,
    2
  );
}

// ---------------------------------------------------------------------------
// TOOL: skill_forge_create
// ---------------------------------------------------------------------------

async function skill_forge_create(skillDef) {
  const { name, version = "1.0.0", tools = [] } = skillDef;

  // Validation
  if (!name || typeof name !== "string") {
    return JSON.stringify({ error: "name is required and must be a string" });
  }
  if (!skillDef.base_url || !isValidUrl(skillDef.base_url)) {
    return JSON.stringify({ error: "base_url is required and must be a valid HTTP(S) URL" });
  }
  if (!Array.isArray(tools) || tools.length === 0) {
    return JSON.stringify({ error: "tools array is required and must have at least one tool" });
  }
  if (tools.length > 30) {
    return JSON.stringify({ error: "Too many tools (max 30). Use include_paths to filter." });
  }

  const safeName = sanitizeName(name);
  const dirName = `${safeName}-${version}`;
  const skillDir = path.join(SKILLS_DIR, dirName);

  // Check for existing skill
  if (fs.existsSync(skillDir)) {
    return JSON.stringify({
      error: `Skill directory already exists: ${dirName}. Remove it first or bump the version.`,
    });
  }

  try {
    // Generate all files
    const indexJS = generateIndexJS(skillDef);
    const skillMD = generateSKILLMD(skillDef);
    const packageJSON = generatePackageJSON(skillDef);

    // Write to disk
    fs.mkdirSync(skillDir, { recursive: true });
    fs.writeFileSync(path.join(skillDir, "index.js"), indexJS, "utf8");
    fs.writeFileSync(path.join(skillDir, "SKILL.md"), skillMD, "utf8");
    fs.writeFileSync(path.join(skillDir, "package.json"), packageJSON, "utf8");

    // Mark as forge-created
    fs.writeFileSync(
      path.join(skillDir, FORGE_MARKER),
      JSON.stringify({ created: new Date().toISOString(), source: "skill-forge" }),
      "utf8"
    );

    // Install npm dependencies if any
    const deps = skillDef.npm_dependencies || {};
    let npmNote = "No npm dependencies.";
    if (Object.keys(deps).length > 0) {
      npmNote = `Run \`cd ${skillDir} && npm install\` to install dependencies: ${Object.keys(deps).join(", ")}`;
    }

    return JSON.stringify(
      {
        status: "created",
        skill: dirName,
        path: skillDir,
        tools_generated: tools.map((t) => sanitizeToolName(t.name)),
        files: ["index.js", "SKILL.md", "package.json"],
        npm: npmNote,
        next_steps: [
          deps && Object.keys(deps).length > 0 ? `cd ${skillDir} && npm install` : null,
          "Restart the gateway to load the new skill: openclaw reload",
          `Test with: "Use the ${sanitizeToolName(tools[0].name)} tool to ..."`,
        ].filter(Boolean),
      },
      null,
      2
    );
  } catch (e) {
    // Clean up on failure
    try {
      if (fs.existsSync(skillDir)) fs.rmSync(skillDir, { recursive: true });
    } catch {}
    return JSON.stringify({ error: `Failed to create skill: ${e.message}` });
  }
}

// ---------------------------------------------------------------------------
// TOOL: skill_forge_validate
// ---------------------------------------------------------------------------

async function skill_forge_validate({ skill_name }) {
  if (!skill_name) {
    return JSON.stringify({ error: "skill_name is required" });
  }

  // Prevent path traversal
  const safeName = path.basename(skill_name);
  const skillDir = path.join(SKILLS_DIR, safeName);

  if (!fs.existsSync(skillDir)) {
    return JSON.stringify({ error: `Skill not found: ${safeName}` });
  }

  const issues = [];
  const info = { name: safeName, path: skillDir, forge_created: false };

  // Check forge marker
  if (fs.existsSync(path.join(skillDir, FORGE_MARKER))) {
    info.forge_created = true;
  }

  // Check SKILL.md
  const skillMdPath = path.join(skillDir, "SKILL.md");
  if (!fs.existsSync(skillMdPath)) {
    issues.push("Missing SKILL.md");
  } else {
    const content = fs.readFileSync(skillMdPath, "utf8");
    if (!content.startsWith("---")) {
      issues.push("SKILL.md missing YAML frontmatter (must start with ---)");
    }
    if (!/^name:\s/m.test(content)) {
      issues.push("SKILL.md frontmatter missing 'name' field");
    }
    if (!/^description:\s/m.test(content)) {
      issues.push("SKILL.md frontmatter missing 'description' field");
    }
  }

  // Check package.json
  const pkgPath = path.join(skillDir, "package.json");
  if (!fs.existsSync(pkgPath)) {
    issues.push("Missing package.json");
  } else {
    try {
      const pkg = JSON.parse(fs.readFileSync(pkgPath, "utf8"));
      info.version = pkg.version;
      info.tool_count = (pkg.tools || []).length;

      if (!pkg.main) issues.push("package.json missing 'main' field");
      if (!pkg.tools || !Array.isArray(pkg.tools)) {
        issues.push("package.json missing 'tools' array");
      }

      // Check index.js exists
      const mainFile = path.join(skillDir, pkg.main || "index.js");
      if (!fs.existsSync(mainFile)) {
        issues.push(`Main file missing: ${pkg.main || "index.js"}`);
      } else {
        // Check exports match tools
        const indexContent = fs.readFileSync(mainFile, "utf8");
        const declared = (pkg.tools || []).map((t) => t.name);
        // Find the LAST module.exports in file (skip template strings)
        const exportMatches = [...indexContent.matchAll(/module\.exports\s*=\s*\{([^}]+)\}/g)];
        const exportMatch = exportMatches.length > 0 ? exportMatches[exportMatches.length - 1] : null;
        if (exportMatch) {
          const exported = exportMatch[1]
            .split(",")
            .map((s) => s.trim())
            .filter(Boolean);
          for (const t of declared) {
            if (!exported.includes(t)) {
              issues.push(`Tool '${t}' declared in package.json but not exported from index.js`);
            }
          }
        }
      }

      // Check npm deps installed
      if (pkg.dependencies && Object.keys(pkg.dependencies).length > 0) {
        const nmPath = path.join(skillDir, "node_modules");
        if (!fs.existsSync(nmPath)) {
          issues.push(`npm dependencies declared but node_modules/ missing — run: cd ${skillDir} && npm install`);
        }
      }
    } catch (e) {
      issues.push(`Invalid package.json: ${e.message}`);
    }
  }

  return JSON.stringify(
    {
      ...info,
      valid: issues.length === 0,
      issues,
    },
    null,
    2
  );
}

// ---------------------------------------------------------------------------
// TOOL: skill_forge_list
// ---------------------------------------------------------------------------

async function skill_forge_list({ filter = "all" } = {}) {
  const entries = [];

  let dirs;
  try {
    dirs = fs.readdirSync(SKILLS_DIR, { withFileTypes: true }).filter((d) => d.isDirectory());
  } catch (e) {
    return JSON.stringify({ error: `Cannot read skills directory: ${e.message}` });
  }

  for (const dir of dirs) {
    const dirPath = path.join(SKILLS_DIR, dir.name);
    const forgeCreated = fs.existsSync(path.join(dirPath, FORGE_MARKER));
    const hasPkg = fs.existsSync(path.join(dirPath, "package.json"));
    const hasSkillMd = fs.existsSync(path.join(dirPath, "SKILL.md"));
    const hasIndex = fs.existsSync(path.join(dirPath, "index.js"));

    let toolCount = 0;
    let version = "?";
    let description = "";

    if (hasPkg) {
      try {
        const pkg = JSON.parse(fs.readFileSync(path.join(dirPath, "package.json"), "utf8"));
        toolCount = (pkg.tools || []).length;
        version = pkg.version || "?";
        description = pkg.description || "";
      } catch {}
    }

    const hasIssues = !hasSkillMd || !hasPkg || (hasPkg && !hasIndex);

    const entry = {
      name: dir.name,
      version,
      tools: toolCount,
      forge_created: forgeCreated,
      has_issues: hasIssues,
      type: hasIndex ? "node" : "markdown-only",
      description: description.slice(0, 100),
    };

    if (filter === "all") entries.push(entry);
    else if (filter === "forge-created" && forgeCreated) entries.push(entry);
    else if (filter === "has-issues" && hasIssues) entries.push(entry);
  }

  return JSON.stringify(
    {
      total: dirs.length,
      showing: entries.length,
      filter,
      skills: entries,
    },
    null,
    2
  );
}

// ---------------------------------------------------------------------------
// TOOL: skill_forge_from_openapi
// ---------------------------------------------------------------------------

async function skill_forge_from_openapi({
  name,
  spec_url,
  spec_json,
  auth_env_var,
  include_paths = [],
  exclude_methods = [],
  max_tools = 20,
}) {
  if (!name) {
    return JSON.stringify({ error: "name is required" });
  }

  // Get the spec
  let spec;
  try {
    if (spec_url) {
      if (!isValidUrl(spec_url)) {
        return JSON.stringify({ error: "spec_url must be a valid HTTP(S) URL" });
      }
      const raw = await httpGet(spec_url);
      spec = JSON.parse(raw);
    } else if (spec_json) {
      spec = JSON.parse(spec_json);
    } else {
      return JSON.stringify({ error: "Either spec_url or spec_json is required" });
    }
  } catch (e) {
    return JSON.stringify({ error: `Failed to load/parse OpenAPI spec: ${e.message}` });
  }

  // Extract info
  const info = spec.info || {};
  const description = info.description || info.title || `${name} API integration`;
  const version = "1.0.0";

  // Determine base URL
  let baseUrl = "";
  if (spec.servers && spec.servers.length > 0) {
    baseUrl = spec.servers[0].url || "";
  } else if (spec.host) {
    // Swagger 2.0
    const scheme = (spec.schemes || ["https"])[0];
    baseUrl = `${scheme}://${spec.host}${spec.basePath || ""}`;
  }
  // Handle relative URLs by resolving against spec_url origin
  if (baseUrl && !isValidUrl(baseUrl) && spec_url) {
    try {
      const origin = new URL(spec_url);
      baseUrl = `${origin.protocol}//${origin.host}${baseUrl}`;
    } catch {}
  }
  // Strip trailing slash
  baseUrl = baseUrl.replace(/\/+$/, "");
  if (!baseUrl || !isValidUrl(baseUrl)) {
    return JSON.stringify({ error: "Could not determine base URL from spec. Set servers[0].url in the spec." });
  }

  // Determine auth
  let authType = "none";
  const securitySchemes = (spec.components && spec.components.securitySchemes) || spec.securityDefinitions || {};
  for (const [, scheme] of Object.entries(securitySchemes)) {
    if (scheme.type === "http" && scheme.scheme === "bearer") {
      authType = "bearer";
      break;
    }
    if (scheme.type === "apiKey") {
      authType = scheme.in === "header" ? "api-key-header" : "api-key-query";
      break;
    }
  }

  // Extract paths and operations
  const tools = [];
  const paths = spec.paths || {};
  const excludeSet = new Set(exclude_methods.map((m) => m.toUpperCase()));

  for (const [pathStr, pathItem] of Object.entries(paths)) {
    // Filter by include_paths
    if (include_paths.length > 0 && !include_paths.some((p) => pathStr.startsWith(p))) {
      continue;
    }

    for (const [method, op] of Object.entries(pathItem)) {
      if (["parameters", "summary", "description", "$ref"].includes(method)) continue;
      if (excludeSet.has(method.toUpperCase())) continue;
      if (tools.length >= max_tools) break;

      const opId = op.operationId || `${method}_${pathStr.replace(/[^a-zA-Z0-9]/g, "_")}`;
      const toolName = sanitizeToolName(opId);
      const desc = op.summary || op.description || `${method.toUpperCase()} ${pathStr}`;

      // Extract parameters
      const parameters = [];
      const allParams = [...(pathItem.parameters || []), ...(op.parameters || [])];

      for (const param of allParams) {
        parameters.push({
          name: param.name,
          type: param.schema ? param.schema.type || "string" : "string",
          description: param.description || "",
          required: param.required || false,
          in: param.in === "path" ? "path" : param.in === "query" ? "query" : "query",
        });
      }

      // Extract request body params (OpenAPI 3.0)
      if (op.requestBody && op.requestBody.content) {
        const jsonContent = op.requestBody.content["application/json"];
        if (jsonContent && jsonContent.schema) {
          const schema = jsonContent.schema;
          const props = schema.properties || {};
          const reqFields = schema.required || [];
          for (const [propName, propDef] of Object.entries(props)) {
            parameters.push({
              name: propName,
              type: propDef.type || "string",
              description: propDef.description || "",
              required: reqFields.includes(propName),
              in: "body",
            });
          }
        }
      }

      tools.push({
        name: toolName,
        description: desc.slice(0, 200),
        method: method.toUpperCase(),
        path: pathStr,
        parameters,
      });
    }
    if (tools.length >= max_tools) break;
  }

  if (tools.length === 0) {
    return JSON.stringify({
      error: "No tools could be extracted from the spec. Check include_paths or the spec format.",
    });
  }

  // Build skill definition and delegate to skill_forge_create
  const skillDef = {
    name,
    version,
    description: description.slice(0, 200),
    base_url: baseUrl,
    auth_type: authType,
    auth_env_var: auth_env_var || `${sanitizeName(name).toUpperCase().replace(/-/g, "_")}_API_KEY`,
    tools,
  };

  return skill_forge_create(skillDef);
}

// ---------------------------------------------------------------------------
// Exports
// ---------------------------------------------------------------------------

module.exports = { skill_forge_create, skill_forge_validate, skill_forge_list, skill_forge_from_openapi };
