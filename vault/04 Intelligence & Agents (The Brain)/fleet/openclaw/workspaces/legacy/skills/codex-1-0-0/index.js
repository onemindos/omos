/**
 * codex-1.0.0 — OpenClaw Skill
 *
 * Native interface to the OneMind Codex knowledge vault.
 * Provides structured access: capture, search, organize,
 * reflect, consolidate, and health audit.
 *
 * Reads _codex/ for conventions and protocols.
 *
 * Registers tools: codex_capture, codex_search, codex_organize,
 *   codex_reflect, codex_consolidate, codex_health, codex_tree
 */

"use strict";

const fs = require("fs");
const path = require("path");

// ---------------------------------------------------------------------------
// Configuration
// ---------------------------------------------------------------------------

// Resolve codex path: env var → workspace sibling → monorepo convention
const CODEX_PATH = process.env.CODEX_PATH || findCodexPath();

function findCodexPath() {
  // Check common locations
  const candidates = [
    path.resolve(__dirname, "../../../../codex"),                    // monorepo: skills/ → workspaces/ → openclaw/ → root/codex
    path.resolve(__dirname, "../../../../../codex"),                 // deeper nesting
    "/home/onemind/onemind-os/codex",                               // VPS
    path.join(process.env.HOME || "", "onemind-os/codex"),          // home dir
  ];
  for (const p of candidates) {
    if (fs.existsSync(p) && fs.existsSync(path.join(p, "_codex"))) return p;
  }
  // Fallback: just the first candidate that exists as a directory
  for (const p of candidates) {
    if (fs.existsSync(p)) return p;
  }
  return candidates[0]; // last resort
}

const CODEX_SYSTEM = path.join(CODEX_PATH, "_codex");
const INBOX_PATH = path.join(CODEX_PATH, "00-24 UI (Unified Intelligence)", "06 Inbox (Queue)");

// Domain mapping: code → folder name
const DOMAIN_MAP = {};

function loadDomainMap() {
  if (Object.keys(DOMAIN_MAP).length > 0) return;
  const quadrants = [
    "00-24 UI (Unified Intelligence)",
    "25-49 HP (Holistic Performance)",
    "50-74 LE (Legacy Evolution)",
    "75-99 GE (Generational Entrepreneurship)",
  ];
  for (const q of quadrants) {
    const qPath = path.join(CODEX_PATH, q);
    if (!fs.existsSync(qPath)) continue;
    for (const entry of fs.readdirSync(qPath, { withFileTypes: true })) {
      if (!entry.isDirectory()) continue;
      const match = entry.name.match(/^(\d{2})\s/);
      if (match) {
        DOMAIN_MAP[match[1]] = path.join(q, entry.name);
      }
    }
  }
}

function resolveDomain(code) {
  loadDomainMap();
  const rel = DOMAIN_MAP[code];
  if (!rel) throw new Error(`Unknown domain code "${code}". Valid: ${Object.keys(DOMAIN_MAP).sort().join(", ")}`);
  return path.join(CODEX_PATH, rel);
}

function today() {
  return new Date().toISOString().slice(0, 10);
}

function slugify(text) {
  return text
    .toLowerCase()
    .replace(/[^a-z0-9]+/g, "-")
    .replace(/^-|-$/g, "")
    .slice(0, 60);
}

// ---------------------------------------------------------------------------
// Read a _codex/ file (for protocol reference — not returned to user, used internally)
// ---------------------------------------------------------------------------

function readCodexSystem(relPath) {
  const full = path.join(CODEX_SYSTEM, relPath);
  if (!fs.existsSync(full)) return null;
  return fs.readFileSync(full, "utf-8");
}

// ---------------------------------------------------------------------------
// Tools
// ---------------------------------------------------------------------------

/**
 * TOOL: codex_capture
 * Create a new note in the codex with proper frontmatter and template.
 */
async function codex_capture({ title, type = "note", domain, content, tags = [] }) {
  if (!title) throw new Error("title is required");
  if (!content) throw new Error("content is required");

  const slug = slugify(title);
  const dateStr = today();
  const domainCode = domain || "06";
  const targetDir = domain ? resolveDomain(domain) : INBOX_PATH;
  const filename = `${slug}.md`;
  const filePath = path.join(targetDir, filename);

  if (fs.existsSync(filePath)) {
    throw new Error(`File already exists: ${filename} in ${domain ? `domain ${domain}` : "inbox"}`);
  }

  // Build frontmatter
  const tagsList = tags.length > 0 ? tags : [`type/${type}`];
  const frontmatter = [
    "---",
    `id: ${domainCode}-${dateStr.replace(/-/g, "")}-${slug}`,
    `title: "${title.replace(/"/g, '\\"')}"`,
    `type: ${type}`,
    `status: active`,
    `created: ${dateStr}`,
    `owner: zeus`,
    `tags: [${tagsList.join(", ")}]`,
    "---",
  ].join("\n");

  const fullContent = `${frontmatter}\n\n# ${title}\n\n${content}\n`;

  // Ensure directory exists
  if (!fs.existsSync(targetDir)) {
    fs.mkdirSync(targetDir, { recursive: true });
  }

  fs.writeFileSync(filePath, fullContent, "utf-8");

  return {
    status: "created",
    file: path.relative(CODEX_PATH, filePath),
    title,
    type,
    domain: domain || "inbox",
  };
}

/**
 * TOOL: codex_search
 * Search the codex using tiered loading protocol.
 */
async function codex_search({ query, domain, level = "auto" }) {
  if (!query) throw new Error("query is required");
  loadDomainMap();

  const results = { query, level_used: level, matches: [] };

  // Determine search scope
  let searchDirs = [];
  if (domain) {
    searchDirs.push(resolveDomain(domain));
  } else {
    // Search all quadrants
    const quadrants = [
      "00-24 UI (Unified Intelligence)",
      "25-49 HP (Holistic Performance)",
      "50-74 LE (Legacy Evolution)",
      "75-99 GE (Generational Entrepreneurship)",
    ];
    for (const q of quadrants) {
      const qp = path.join(CODEX_PATH, q);
      if (fs.existsSync(qp)) searchDirs.push(qp);
    }
  }

  const queryLower = query.toLowerCase();
  const useLevel = level === "auto" ? "L1" : level;

  if (useLevel === "L0") {
    // Just list domain folders
    for (const dir of searchDirs) {
      const entries = fs.readdirSync(dir, { withFileTypes: true })
        .filter(e => e.isDirectory() && !e.name.startsWith(".") && !e.name.startsWith("_"));
      results.matches.push({
        path: path.relative(CODEX_PATH, dir),
        type: "directory",
        subfolders: entries.map(e => e.name),
      });
    }
    results.level_used = "L0";
    return results;
  }

  // L1 and L2: scan files
  function scanDir(dir, depth = 0) {
    if (depth > 3) return;
    let entries;
    try { entries = fs.readdirSync(dir, { withFileTypes: true }); } catch { return; }

    for (const entry of entries) {
      if (entry.name.startsWith(".") || entry.name.startsWith("_")) continue;

      const fullPath = path.join(dir, entry.name);

      if (entry.isDirectory()) {
        // Check if folder name matches query
        if (entry.name.toLowerCase().includes(queryLower)) {
          results.matches.push({
            path: path.relative(CODEX_PATH, fullPath),
            type: "directory",
            match: "folder_name",
          });
        }
        scanDir(fullPath, depth + 1);
      } else if (entry.name.endsWith(".md")) {
        const relPath = path.relative(CODEX_PATH, fullPath);

        // L1: match on filename
        if (entry.name.toLowerCase().includes(queryLower)) {
          if (useLevel === "L2") {
            const content = fs.readFileSync(fullPath, "utf-8");
            results.matches.push({ path: relPath, type: "file", match: "filename", preview: content.slice(0, 500) });
          } else {
            results.matches.push({ path: relPath, type: "file", match: "filename" });
          }
          continue;
        }

        // L1/L2: match on content
        try {
          const content = fs.readFileSync(fullPath, "utf-8");
          if (content.toLowerCase().includes(queryLower)) {
            const lineIdx = content.toLowerCase().indexOf(queryLower);
            const snippet = content.slice(Math.max(0, lineIdx - 80), lineIdx + 120).trim();
            results.matches.push({
              path: relPath,
              type: "file",
              match: "content",
              preview: useLevel === "L2" ? content.slice(0, 1000) : snippet,
            });
          }
        } catch { /* skip files that can't be read */ }
      }
    }
  }

  for (const dir of searchDirs) {
    scanDir(dir);
  }

  results.level_used = useLevel;
  results.total = results.matches.length;

  // Limit results to prevent context overflow
  if (results.matches.length > 20) {
    results.matches = results.matches.slice(0, 20);
    results.truncated = true;
  }

  return results;
}

/**
 * TOOL: codex_organize
 * Move items from inbox to the correct domain folder.
 */
async function codex_organize({ file, target_domain, process_all = false }) {
  loadDomainMap();

  if (process_all) {
    // List inbox items
    if (!fs.existsSync(INBOX_PATH)) return { status: "ok", message: "Inbox folder does not exist", items: [] };
    const items = fs.readdirSync(INBOX_PATH)
      .filter(f => f.endsWith(".md") && f !== "README.md");
    return {
      status: "inbox_listed",
      count: items.length,
      items: items.map(f => ({
        file: path.join("06 Inbox (Queue)", f),
        name: f,
      })),
      instruction: "Review each item and call codex_organize with file + target_domain to move them.",
    };
  }

  if (!file) throw new Error("file is required (relative to codex root)");
  if (!target_domain) throw new Error("target_domain is required (e.g. '08' for Infrastructure)");

  const srcPath = path.join(CODEX_PATH, file);
  if (!fs.existsSync(srcPath)) throw new Error(`File not found: ${file}`);

  const targetDir = resolveDomain(target_domain);
  const destPath = path.join(targetDir, path.basename(file));

  if (fs.existsSync(destPath)) {
    throw new Error(`Target already exists: ${path.relative(CODEX_PATH, destPath)}`);
  }

  // Read file, update domain tag in frontmatter
  let content = fs.readFileSync(srcPath, "utf-8");
  // Simple tag update: add domain tag if frontmatter exists
  const quadrant = target_domain < "25" ? "ui" : target_domain < "50" ? "hp" : target_domain < "75" ? "le" : "ge";
  if (content.includes("tags: [") && !content.includes(`domain/${quadrant}`)) {
    content = content.replace(/tags: \[/, `tags: [domain/${quadrant}, `);
  }

  fs.writeFileSync(destPath, content, "utf-8");
  fs.unlinkSync(srcPath);

  return {
    status: "moved",
    from: file,
    to: path.relative(CODEX_PATH, destPath),
    domain: target_domain,
  };
}

/**
 * TOOL: codex_reflect
 * Run a reflection cycle. Returns the protocol + current data for the agent to act on.
 */
async function codex_reflect({ cadence }) {
  if (!cadence) throw new Error("cadence is required (daily, weekly, monthly)");

  const protocol = readCodexSystem(`skills/reflect.md`);
  const template = readCodexSystem(`templates/${cadence === "daily" ? "daily" : "review"}.md`);
  const profile = cadence === "weekly" ? readCodexSystem("profiles/review.md") : null;

  // Gather context based on cadence
  const context = { cadence, date: today() };

  if (cadence === "daily") {
    context.instruction = "Create a daily log following the protocol. Use the template provided.";
  } else if (cadence === "weekly") {
    // Count inbox items
    if (fs.existsSync(INBOX_PATH)) {
      context.inbox_count = fs.readdirSync(INBOX_PATH).filter(f => f.endsWith(".md") && f !== "README.md").length;
    }
    context.instruction = "Follow the weekly review protocol. Process inbox items, review goals, create review note.";
  } else if (cadence === "monthly") {
    context.instruction = "Follow the monthly review protocol. Audit each quadrant, consolidate knowledge, check vault health.";
  }

  return {
    status: "protocol_loaded",
    cadence,
    protocol: protocol ? protocol.slice(0, 3000) : "Protocol file not found at _codex/skills/reflect.md",
    template: template || "Template not found",
    profile: profile ? profile.slice(0, 1500) : undefined,
    context,
  };
}

/**
 * TOOL: codex_consolidate
 * Find and return scattered notes on a topic for consolidation.
 */
async function codex_consolidate({ topic, source_files }) {
  const protocol = readCodexSystem("skills/consolidate.md");

  let notes = [];

  if (source_files && source_files.length > 0) {
    // Read specified files
    for (const f of source_files) {
      const full = path.join(CODEX_PATH, f);
      if (fs.existsSync(full)) {
        notes.push({ path: f, content: fs.readFileSync(full, "utf-8").slice(0, 2000) });
      }
    }
  } else if (topic) {
    // Search for topic across the vault
    const searchResult = await codex_search({ query: topic, level: "L1" });
    notes = searchResult.matches.map(m => ({ path: m.path, preview: m.preview || m.match }));
  }

  return {
    status: "ready_to_consolidate",
    topic: topic || "specified files",
    notes_found: notes.length,
    notes: notes.slice(0, 10),
    protocol: protocol ? protocol.slice(0, 2000) : "Protocol not found",
    instruction: "Review these notes and merge related content into one curated document. Use codex_capture to create the consolidated note, then delete or archive the originals.",
  };
}

/**
 * TOOL: codex_health
 * Audit vault health.
 */
async function codex_health() {
  loadDomainMap();

  const report = {
    status: "ok",
    codex_path: CODEX_PATH,
    codex_system: fs.existsSync(CODEX_SYSTEM),
    timestamp: new Date().toISOString(),
    issues: [],
  };

  // 1. Inbox count
  if (fs.existsSync(INBOX_PATH)) {
    const inboxItems = fs.readdirSync(INBOX_PATH).filter(f => f.endsWith(".md") && f !== "README.md");
    report.inbox_count = inboxItems.length;
    if (inboxItems.length > 20) {
      report.issues.push({ severity: "high", issue: `Inbox overflow: ${inboxItems.length} items (threshold: 20)` });
    }
  } else {
    report.inbox_count = 0;
    report.issues.push({ severity: "medium", issue: "Inbox folder does not exist" });
  }

  // 2. Domain folder coverage
  const expectedDomains = Object.keys(DOMAIN_MAP).length;
  report.domains_found = expectedDomains;

  // 3. Check for missing READMEs in domain folders
  let missingReadmes = [];
  for (const [code, rel] of Object.entries(DOMAIN_MAP)) {
    const readmePath = path.join(CODEX_PATH, rel, "README.md");
    if (!fs.existsSync(readmePath)) {
      missingReadmes.push(`${code} (${path.basename(rel)})`);
    }
  }
  if (missingReadmes.length > 0) {
    report.issues.push({
      severity: "low",
      issue: `Missing README.md in ${missingReadmes.length} domain folders`,
      details: missingReadmes.slice(0, 10),
    });
  }

  // 4. Check _codex/ completeness
  const expectedSystemFiles = [
    "INTERFACE.md", "CONVENTIONS.md",
    "skills/search.md", "skills/capture.md", "skills/organize.md",
    "skills/reflect.md", "skills/consolidate.md", "skills/maintain.md",
  ];
  const missingSystem = expectedSystemFiles.filter(f => !fs.existsSync(path.join(CODEX_SYSTEM, f)));
  if (missingSystem.length > 0) {
    report.issues.push({
      severity: "medium",
      issue: `Missing _codex/ system files: ${missingSystem.join(", ")}`,
    });
  }

  // 5. Count total files
  let totalFiles = 0;
  function countFiles(dir) {
    let entries;
    try { entries = fs.readdirSync(dir, { withFileTypes: true }); } catch { return; }
    for (const e of entries) {
      if (e.name.startsWith(".")) continue;
      if (e.isDirectory()) {
        if (e.name !== "node_modules") countFiles(path.join(dir, e.name));
      } else if (e.name.endsWith(".md")) {
        totalFiles++;
      }
    }
  }
  countFiles(CODEX_PATH);
  report.total_markdown_files = totalFiles;

  // Summary
  if (report.issues.length === 0) {
    report.summary = "Vault is healthy. No issues found.";
  } else {
    const highCount = report.issues.filter(i => i.severity === "high").length;
    const medCount = report.issues.filter(i => i.severity === "medium").length;
    report.summary = `Found ${report.issues.length} issues (${highCount} high, ${medCount} medium). See details.`;
  }

  return report;
}

/**
 * TOOL: codex_tree
 * Show folder structure at a given depth.
 */
async function codex_tree({ path: subPath = "", depth = 2 }) {
  const basePath = subPath ? path.join(CODEX_PATH, subPath) : CODEX_PATH;

  if (!fs.existsSync(basePath)) {
    throw new Error(`Path not found: ${subPath || "codex root"}`);
  }

  function buildTree(dir, currentDepth) {
    if (currentDepth >= depth) return null;
    let entries;
    try { entries = fs.readdirSync(dir, { withFileTypes: true }); } catch { return null; }

    const result = {};
    const sorted = entries
      .filter(e => !e.name.startsWith(".") && e.name !== "node_modules")
      .sort((a, b) => {
        // Directories first, then files
        if (a.isDirectory() && !b.isDirectory()) return -1;
        if (!a.isDirectory() && b.isDirectory()) return 1;
        return a.name.localeCompare(b.name);
      });

    for (const entry of sorted) {
      if (entry.isDirectory()) {
        const children = buildTree(path.join(dir, entry.name), currentDepth + 1);
        result[entry.name + "/"] = children || "...";
      } else if (entry.name.endsWith(".md")) {
        result[entry.name] = "file";
      }
    }
    return result;
  }

  return {
    root: subPath || "/",
    depth,
    tree: buildTree(basePath, 0),
  };
}

// ---------------------------------------------------------------------------
// Exports
// ---------------------------------------------------------------------------

module.exports = {
  codex_capture,
  codex_search,
  codex_organize,
  codex_reflect,
  codex_consolidate,
  codex_health,
  codex_tree,
};
