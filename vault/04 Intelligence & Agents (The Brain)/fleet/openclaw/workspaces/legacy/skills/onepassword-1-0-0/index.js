/**
 * onepassword-1.0.0 — OpenClaw Skill
 *
 * Read, write, and update items in 1Password vaults via the
 * 1Password Service Accounts REST API (api.1password.com).
 * No Connect server required — just a service account token.
 *
 * Primary use: store VPS config backups (openclaw.json, .env) as
 * secure notes so secrets are never in git.
 *
 * Registers tools:
 *   op_get_item      — retrieve an item (including secure notes)
 *   op_save_note     — create or overwrite a secure note
 *   op_list_items    — list items in a vault (name + id)
 *   op_get_field     — get a single field value from an item
 *   op_update_field  — update a single field value on an item
 *
 * Requires env vars:
 *   OP_SERVICE_ACCOUNT_TOKEN  — 1Password service account token
 *   OP_VAULT_ID               — default vault ID (OneMind-UI vault)
 */

"use strict";

const https = require("https");

const BASE_URL = "https://api.1password.com/v1";
const TOKEN = process.env.OP_SERVICE_ACCOUNT_TOKEN;
const DEFAULT_VAULT = process.env.OP_VAULT_ID;

// ---------------------------------------------------------------------------
// HTTP helper
// ---------------------------------------------------------------------------

function request(method, path, body = null) {
  return new Promise((resolve, reject) => {
    if (!TOKEN) return reject(new Error("OP_SERVICE_ACCOUNT_TOKEN not set"));
    const url = new URL(BASE_URL + path);
    const payload = body ? JSON.stringify(body) : null;
    const opts = {
      hostname: url.hostname,
      path: url.pathname + url.search,
      method,
      headers: {
        "Authorization": `Bearer ${TOKEN}`,
        "Content-Type": "application/json",
        "User-Agent": "OpenClaw-1Password/1.0",
      },
      timeout: 15000,
    };
    if (payload) opts.headers["Content-Length"] = Buffer.byteLength(payload);

    const req = https.request(opts, (res) => {
      let data = "";
      res.setEncoding("utf8");
      res.on("data", (c) => (data += c));
      res.on("end", () => {
        if (res.statusCode < 200 || res.statusCode >= 300) {
          return reject(new Error(`HTTP ${res.statusCode}: ${data}`));
        }
        try { resolve(data ? JSON.parse(data) : {}); }
        catch { resolve(data); }
      });
    });
    req.on("timeout", () => { req.destroy(); reject(new Error("Request timed out")); });
    req.on("error", reject);
    if (payload) req.write(payload);
    req.end();
  });
}

// ---------------------------------------------------------------------------
// Helpers
// ---------------------------------------------------------------------------

function resolveVault(vaultId) {
  const v = vaultId || DEFAULT_VAULT;
  if (!v) throw new Error("No vault ID provided and OP_VAULT_ID not set");
  return v;
}

function findField(item, label) {
  const fields = item.fields || [];
  return fields.find(
    (f) => f.label?.toLowerCase() === label.toLowerCase() ||
           f.id?.toLowerCase() === label.toLowerCase()
  );
}

// ---------------------------------------------------------------------------
// Tool: op_get_item
// ---------------------------------------------------------------------------

async function opGetItem({ item_id, vault_id }) {
  const vault = resolveVault(vault_id);
  const item = await request("GET", `/vaults/${vault}/items/${item_id}`);
  // Surface the most useful fields
  const fields = (item.fields || []).map((f) => ({
    label: f.label,
    value: f.value,
    type: f.type,
  }));
  return {
    id: item.id,
    title: item.title,
    category: item.category,
    vault: item.vault?.id,
    fields,
    updated_at: item.updated_at,
  };
}

// ---------------------------------------------------------------------------
// Tool: op_save_note
// Creates a new secure note or updates the "notesPlain" field of an existing one.
// Pass item_id to update, omit to create.
// ---------------------------------------------------------------------------

async function opSaveNote({ title, content, vault_id, item_id, tags }) {
  const vault = resolveVault(vault_id);

  if (item_id) {
    // Update existing item — patch the notesPlain field
    const item = await request("GET", `/vaults/${vault}/items/${item_id}`);
    const fields = item.fields || [];
    const noteField = fields.find((f) => f.purpose === "NOTES") ||
                      fields.find((f) => f.id === "notesPlain");
    const updatedFields = noteField
      ? fields.map((f) =>
          (f.purpose === "NOTES" || f.id === "notesPlain")
            ? { ...f, value: content }
            : f
        )
      : [...fields, { id: "notesPlain", purpose: "NOTES", value: content }];

    await request("PUT", `/vaults/${vault}/items/${item_id}`, {
      ...item,
      title: title || item.title,
      fields: updatedFields,
    });
    return { success: true, action: "updated", item_id, title: title || item.title };
  }

  // Create new secure note
  const payload = {
    vault: { id: vault },
    title,
    category: "SECURE_NOTE",
    fields: [{ id: "notesPlain", purpose: "NOTES", value: content }],
  };
  if (tags?.length) payload.tags = tags;
  const created = await request("POST", `/vaults/${vault}/items`, payload);
  return { success: true, action: "created", item_id: created.id, title };
}

// ---------------------------------------------------------------------------
// Tool: op_list_items
// ---------------------------------------------------------------------------

async function opListItems({ vault_id, filter_tag }) {
  const vault = resolveVault(vault_id);
  let path = `/vaults/${vault}/items`;
  if (filter_tag) path += `?filter=tag%3D"${encodeURIComponent(filter_tag)}"`;
  const items = await request("GET", path);
  return (Array.isArray(items) ? items : []).map((i) => ({
    id: i.id,
    title: i.title,
    category: i.category,
    updated_at: i.updated_at,
    tags: i.tags,
  }));
}

// ---------------------------------------------------------------------------
// Tool: op_get_field
// ---------------------------------------------------------------------------

async function opGetField({ item_id, field_label, vault_id }) {
  const vault = resolveVault(vault_id);
  const item = await request("GET", `/vaults/${vault}/items/${item_id}`);
  const field = findField(item, field_label);
  if (!field) throw new Error(`Field "${field_label}" not found in item ${item_id}`);
  return { label: field.label, value: field.value, type: field.type };
}

// ---------------------------------------------------------------------------
// Tool: op_update_field
// ---------------------------------------------------------------------------

async function opUpdateField({ item_id, field_label, value, vault_id }) {
  const vault = resolveVault(vault_id);
  const item = await request("GET", `/vaults/${vault}/items/${item_id}`);
  const fields = item.fields || [];
  let found = false;
  const updatedFields = fields.map((f) => {
    if (f.label?.toLowerCase() === field_label.toLowerCase() ||
        f.id?.toLowerCase() === field_label.toLowerCase()) {
      found = true;
      return { ...f, value };
    }
    return f;
  });
  if (!found) throw new Error(`Field "${field_label}" not found in item ${item_id}`);
  await request("PUT", `/vaults/${vault}/items/${item_id}`, {
    ...item,
    fields: updatedFields,
  });
  return { success: true, item_id, field_label, updated: true };
}

// ---------------------------------------------------------------------------
// Register tools
// ---------------------------------------------------------------------------

module.exports = {
  tools: [
    {
      name: "op_get_item",
      description: "Get a 1Password item by ID, returning all fields. Use to retrieve secure notes, credentials, or config backups.",
      inputSchema: {
        type: "object",
        properties: {
          item_id: { type: "string", description: "The 1Password item ID" },
          vault_id: { type: "string", description: "Vault ID (defaults to OP_VAULT_ID env var)" },
        },
        required: ["item_id"],
      },
      handler: opGetItem,
    },
    {
      name: "op_save_note",
      description: "Create a new secure note in 1Password, or update an existing one. Use to back up VPS configs (openclaw.json, .env contents) as secrets.",
      inputSchema: {
        type: "object",
        properties: {
          title: { type: "string", description: "Note title, e.g. 'VPS openclaw.json backup'" },
          content: { type: "string", description: "Full text content of the note" },
          vault_id: { type: "string", description: "Vault ID (defaults to OP_VAULT_ID)" },
          item_id: { type: "string", description: "Existing item ID to update (omit to create new)" },
          tags: {
            type: "array",
            items: { type: "string" },
            description: "Tags to apply, e.g. ['path:ui', 'category:infrastructure']",
          },
        },
        required: ["title", "content"],
      },
      handler: opSaveNote,
    },
    {
      name: "op_list_items",
      description: "List items in a 1Password vault. Returns title, id, category, and tags for each item.",
      inputSchema: {
        type: "object",
        properties: {
          vault_id: { type: "string", description: "Vault ID (defaults to OP_VAULT_ID)" },
          filter_tag: { type: "string", description: "Optional: filter by tag, e.g. 'category:infrastructure'" },
        },
      },
      handler: opListItems,
    },
    {
      name: "op_get_field",
      description: "Get a single field value from a 1Password item by label.",
      inputSchema: {
        type: "object",
        properties: {
          item_id: { type: "string", description: "The 1Password item ID" },
          field_label: { type: "string", description: "Field label to retrieve, e.g. 'username', 'password', 'API Key'" },
          vault_id: { type: "string", description: "Vault ID (defaults to OP_VAULT_ID)" },
        },
        required: ["item_id", "field_label"],
      },
      handler: opGetField,
    },
    {
      name: "op_update_field",
      description: "Update a single field value on an existing 1Password item.",
      inputSchema: {
        type: "object",
        properties: {
          item_id: { type: "string", description: "The 1Password item ID" },
          field_label: { type: "string", description: "Field label to update" },
          value: { type: "string", description: "New value" },
          vault_id: { type: "string", description: "Vault ID (defaults to OP_VAULT_ID)" },
        },
        required: ["item_id", "field_label", "value"],
      },
      handler: opUpdateField,
    },
  ],
};
