"use strict";

/**
 * api-integration — OpenClaw Skill
 *
 * REST API integration template. Provides standard CRUD operations
 * against a JSON REST API with Bearer token authentication.
 *
 * Registered tools: api_list, api_get, api_create, api_update, api_delete
 */

// ─── Configuration ───────────────────────────────────────────────────────────
const BASE_URL = (process.env.MY_SERVICE_URL || "https://api.example.com").replace(/\/$/, "");
const API_KEY = process.env.MY_SERVICE_API_KEY;

/** Shared fetch wrapper with auth, timeout, and error handling */
async function apiRequest(path, { method = "GET", body = null } = {}) {
  if (!API_KEY) throw new Error("MY_SERVICE_API_KEY env var is not configured");

  const options = {
    method,
    headers: {
      Authorization: `Bearer ${API_KEY}`,
      "Content-Type": "application/json",
      Accept: "application/json",
    },
    signal: AbortSignal.timeout(30000),
  };

  if (body) options.body = JSON.stringify(body);

  const res = await fetch(`${BASE_URL}${path}`, options);

  if (!res.ok) {
    let errorText;
    try {
      errorText = await res.text();
    } catch {
      errorText = `HTTP ${res.status}`;
    }
    throw new Error(`API error ${res.status} on ${method} ${path}: ${errorText}`);
  }

  // Handle empty responses (e.g. 204 No Content)
  const contentType = res.headers.get("content-type") || "";
  if (!contentType.includes("application/json")) return { status: "ok" };

  return res.json();
}

// ─── Tools ───────────────────────────────────────────────────────────────────

/**
 * TOOL: api_list
 * List resources with optional filters and pagination.
 *
 * @param {string} [filter] - Optional filter string (e.g., "status=active")
 * @param {number} [limit=20] - Max results to return
 * @param {number} [offset=0] - Pagination offset
 * @returns {Object} { items: [], total: number }
 */
async function api_list({ filter = null, limit = 20, offset = 0 }) {
  // TODO: adjust the endpoint and query param names to match your API
  const params = new URLSearchParams({ limit: String(limit), offset: String(offset) });
  if (filter) params.set("filter", filter);

  return apiRequest(`/items?${params}`);
}

/**
 * TOOL: api_get
 * Get a single resource by ID.
 *
 * @param {string} id - Resource ID
 * @returns {Object} The resource
 */
async function api_get({ id }) {
  if (!id) throw new Error("id is required");
  return apiRequest(`/items/${encodeURIComponent(id)}`);
}

/**
 * TOOL: api_create
 * Create a new resource.
 *
 * @param {string} name - Resource name (required)
 * @param {string} [description] - Optional description
 * @param {Object} [extra] - Any additional fields for the resource
 * @returns {Object} The created resource
 */
async function api_create({ name, description, ...extra }) {
  if (!name) throw new Error("name is required");
  return apiRequest("/items", {
    method: "POST",
    body: { name, description, ...extra },
  });
}

/**
 * TOOL: api_update
 * Update an existing resource by ID.
 *
 * @param {string} id - Resource ID (required)
 * @param {Object} updates - Fields to update (name, description, status, etc.)
 * @returns {Object} The updated resource
 */
async function api_update({ id, ...updates }) {
  if (!id) throw new Error("id is required");
  if (Object.keys(updates).length === 0) throw new Error("No update fields provided");
  return apiRequest(`/items/${encodeURIComponent(id)}`, { method: "PATCH", body: updates });
}

/**
 * TOOL: api_delete
 * Delete a resource by ID. This is irreversible.
 *
 * @param {string} id - Resource ID (required)
 * @returns {Object} { status: "deleted", id }
 */
async function api_delete({ id }) {
  if (!id) throw new Error("id is required");
  await apiRequest(`/items/${encodeURIComponent(id)}`, { method: "DELETE" });
  return { status: "deleted", id };
}

// ─── Exports ─────────────────────────────────────────────────────────────────
module.exports = { api_list, api_get, api_create, api_update, api_delete };
