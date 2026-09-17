// Entity Registry API
// Reads markdown files from 05 Environments/ and serves them as structured JSON.
// Does NOT touch or modify any markdown files.
//
// Endpoints:
//   GET /api/entities               → all entities, all types
//   GET /api/entities/:type         → e.g. /api/entities/machines
//   GET /api/entities/:type/:name   → e.g. /api/entities/machines/vps
//   GET /api/health                 → liveness check

const express = require("express");
const fs = require("fs");
const path = require("path");
const matter = require("gray-matter");

const app = express();
const PORT = process.env.PORT || 18810;

// Root of the entity registry — one level up from this api/ folder
const REGISTRY_ROOT = path.resolve(__dirname, "..");

// Entity type → folder mapping
const ENTITY_FOLDERS = {
  users: "users",
  locations: "locations",
  machines: "machines",
  devices: "devices",
  assets: "assets",
  "digital-assets": "digital-assets",
};

// ─── Helpers ─────────────────────────────────────────────────────────────────

/**
 * Parse a single markdown file into a structured entity object.
 * Returns frontmatter data merged with { _name, _type, _file }.
 */
function parseEntity(filePath, type) {
  try {
    const raw = fs.readFileSync(filePath, "utf8");
    const { data, content } = matter(raw);
    const name = path.basename(filePath, ".md");
    return {
      _name: name,
      _type: type,
      _file: path.relative(REGISTRY_ROOT, filePath),
      ...data,
      // Attach trimmed body for agents that want the full context
      _body: content.trim(),
    };
  } catch (err) {
    return { _name: path.basename(filePath, ".md"), _type: type, _error: err.message };
  }
}

/**
 * Load all entities of a given type from its folder.
 */
function loadType(type) {
  const folder = ENTITY_FOLDERS[type];
  if (!folder) return null;

  const dir = path.join(REGISTRY_ROOT, folder);
  if (!fs.existsSync(dir)) return [];

  return fs
    .readdirSync(dir)
    .filter((f) => f.endsWith(".md") && !f.startsWith("_"))
    .map((f) => parseEntity(path.join(dir, f), type));
}

/**
 * Load all entities across all types.
 */
function loadAll() {
  const result = {};
  for (const type of Object.keys(ENTITY_FOLDERS)) {
    result[type] = loadType(type) || [];
  }
  return result;
}

// ─── Routes ──────────────────────────────────────────────────────────────────

app.get("/api/health", (req, res) => {
  res.json({ status: "ok", registry: REGISTRY_ROOT, timestamp: new Date().toISOString() });
});

// GET /api/entities → all types
app.get("/api/entities", (req, res) => {
  res.json(loadAll());
});

// GET /api/entities/:type → e.g. /api/entities/machines
app.get("/api/entities/:type", (req, res) => {
  const { type } = req.params;
  const entities = loadType(type);
  if (entities === null) {
    return res.status(404).json({ error: `Unknown entity type: ${type}`, available: Object.keys(ENTITY_FOLDERS) });
  }
  res.json(entities);
});

// GET /api/entities/:type/:name → e.g. /api/entities/machines/vps
app.get("/api/entities/:type/:name", (req, res) => {
  const { type, name } = req.params;
  const folder = ENTITY_FOLDERS[type];
  if (!folder) {
    return res.status(404).json({ error: `Unknown entity type: ${type}` });
  }

  const filePath = path.join(REGISTRY_ROOT, folder, `${name}.md`);
  if (!fs.existsSync(filePath)) {
    return res.status(404).json({ error: `Entity not found: ${type}/${name}` });
  }

  res.json(parseEntity(filePath, type));
});

// ─── Live data stub hooks (not wired yet) ────────────────────────────────────
// These endpoints exist as stubs for future integrations.
// POST /api/entities/:type/:name/update → write live data back to an entity
// POST /api/entities/:type/:name/location → update GPS/HA location
// GET  /api/entities/presence → who/what is home right now

app.post("/api/entities/:type/:name/update", (req, res) => {
  res.status(501).json({ status: "not_implemented", message: "Live write-back not wired yet." });
});

app.get("/api/entities/presence", (req, res) => {
  res.status(501).json({ status: "not_implemented", message: "Home Assistant presence not wired yet." });
});

// ─── Start ───────────────────────────────────────────────────────────────────

app.listen(PORT, () => {
  console.log(`Entity Registry API running on port ${PORT}`);
  console.log(`Registry root: ${REGISTRY_ROOT}`);
  console.log(`Endpoints: /api/health  /api/entities  /api/entities/:type  /api/entities/:type/:name`);
});
