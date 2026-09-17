"use strict";

/**
 * todoist-api — OpenClaw Skill
 *
 * Todoist API integration via Maton gateway with managed OAuth.
 * Manage tasks, projects, sections, labels, and comments.
 *
 * Requires: MATON_API_KEY environment variable
 *
 * Registered tools: todoist_tasks, todoist_projects, todoist_sections, todoist_labels
 */

const BASE_URL = "https://gateway.maton.ai/todoist/api/v1";

async function todoistFetch(endpoint, { method = "GET", body } = {}) {
  const apiKey = process.env.MATON_API_KEY;
  if (!apiKey) throw new Error("MATON_API_KEY environment variable is required");
  const opts = {
    method,
    headers: {
      "Authorization": `Bearer ${apiKey}`,
      "Content-Type": "application/json",
    },
    signal: AbortSignal.timeout(30000),
  };
  if (body) opts.body = JSON.stringify(body);
  const res = await fetch(`${BASE_URL}${endpoint}`, opts);
  if (!res.ok) {
    const text = await res.text();
    throw new Error(`Todoist API ${res.status}: ${text}`);
  }
  return res.json();
}

/**
 * TOOL: todoist_tasks
 * Manage Todoist tasks — list, get, create, update, close, reopen, delete.
 *
 * @param {string} action - One of: list, get, create, update, close, reopen, delete
 * @param {string} [id] - Task ID (required for get/update/close/reopen/delete)
 * @param {string} [filter] - Filter expression for list (e.g. "today", "overdue", "#ProjectName")
 * @param {string} [project_id] - Project ID to filter tasks
 * @param {string} [content] - Task content (for create/update)
 * @param {string} [description] - Task description
 * @param {string} [due_string] - Due date in natural language (e.g. "tomorrow", "every monday")
 * @param {number} [priority] - Priority 1-4 (4 = urgent)
 * @param {string[]} [labels] - Label names
 * @returns {Object} result
 */
async function todoist_tasks({ action, id, filter, project_id, content, description, due_string, priority, labels }) {
  if (!action) throw new Error("action is required (list, get, create, update, close, reopen, delete)");

  switch (action) {
    case "list": {
      const params = new URLSearchParams();
      if (filter) params.set("filter", filter);
      if (project_id) params.set("project_id", project_id);
      const qs = params.toString();
      return todoistFetch(`/tasks${qs ? `?${qs}` : ""}`);
    }
    case "get":
      if (!id) throw new Error("id is required for get");
      return todoistFetch(`/tasks/${id}`);
    case "create":
      if (!content) throw new Error("content is required for create");
      return todoistFetch("/tasks", { method: "POST", body: { content, description, due_string, priority, labels, project_id } });
    case "update":
      if (!id) throw new Error("id is required for update");
      return todoistFetch(`/tasks/${id}`, { method: "POST", body: { content, description, due_string, priority, labels } });
    case "close":
      if (!id) throw new Error("id is required for close");
      return todoistFetch(`/tasks/${id}/close`, { method: "POST" });
    case "reopen":
      if (!id) throw new Error("id is required for reopen");
      return todoistFetch(`/tasks/${id}/reopen`, { method: "POST" });
    case "delete":
      if (!id) throw new Error("id is required for delete");
      return todoistFetch(`/tasks/${id}`, { method: "DELETE" });
    default:
      throw new Error(`Unknown action: ${action}. Use: list, get, create, update, close, reopen, delete`);
  }
}

/**
 * TOOL: todoist_projects
 * Manage Todoist projects — list, get, create, update, delete.
 *
 * @param {string} action - One of: list, get, create, update, delete
 * @param {string} [id] - Project ID (for get/update/delete)
 * @param {string} [name] - Project name (for create/update)
 * @param {string} [color] - Project color
 * @param {string} [parent_id] - Parent project ID for nesting
 * @param {string} [view_style] - "list" or "board"
 * @returns {Object} result
 */
async function todoist_projects({ action, id, name, color, parent_id, view_style }) {
  if (!action) throw new Error("action is required (list, get, create, update, delete)");

  switch (action) {
    case "list":
      return todoistFetch("/projects");
    case "get":
      if (!id) throw new Error("id is required");
      return todoistFetch(`/projects/${id}`);
    case "create":
      if (!name) throw new Error("name is required for create");
      return todoistFetch("/projects", { method: "POST", body: { name, color, parent_id, view_style } });
    case "update":
      if (!id) throw new Error("id is required");
      return todoistFetch(`/projects/${id}`, { method: "POST", body: { name, color, view_style } });
    case "delete":
      if (!id) throw new Error("id is required");
      return todoistFetch(`/projects/${id}`, { method: "DELETE" });
    default:
      throw new Error(`Unknown action: ${action}`);
  }
}

/**
 * TOOL: todoist_sections
 * Manage Todoist sections within projects — list, create, update, delete.
 *
 * @param {string} action - One of: list, create, update, delete
 * @param {string} [id] - Section ID (for update/delete)
 * @param {string} [project_id] - Project ID (for list/create)
 * @param {string} [name] - Section name (for create/update)
 * @returns {Object} result
 */
async function todoist_sections({ action, id, project_id, name }) {
  if (!action) throw new Error("action is required (list, create, update, delete)");

  switch (action) {
    case "list": {
      const qs = project_id ? `?project_id=${project_id}` : "";
      return todoistFetch(`/sections${qs}`);
    }
    case "create":
      if (!name) throw new Error("name is required");
      if (!project_id) throw new Error("project_id is required");
      return todoistFetch("/sections", { method: "POST", body: { name, project_id } });
    case "update":
      if (!id) throw new Error("id is required");
      return todoistFetch(`/sections/${id}`, { method: "POST", body: { name } });
    case "delete":
      if (!id) throw new Error("id is required");
      return todoistFetch(`/sections/${id}`, { method: "DELETE" });
    default:
      throw new Error(`Unknown action: ${action}`);
  }
}

/**
 * TOOL: todoist_labels
 * Manage personal labels — list, create, update, delete.
 *
 * @param {string} action - One of: list, create, update, delete
 * @param {string} [id] - Label ID (for update/delete)
 * @param {string} [name] - Label name (for create/update)
 * @param {string} [color] - Label color
 * @returns {Object} result
 */
async function todoist_labels({ action, id, name, color }) {
  if (!action) throw new Error("action is required (list, create, update, delete)");

  switch (action) {
    case "list":
      return todoistFetch("/labels");
    case "create":
      if (!name) throw new Error("name is required");
      return todoistFetch("/labels", { method: "POST", body: { name, color } });
    case "update":
      if (!id) throw new Error("id is required");
      return todoistFetch(`/labels/${id}`, { method: "POST", body: { name, color } });
    case "delete":
      if (!id) throw new Error("id is required");
      return todoistFetch(`/labels/${id}`, { method: "DELETE" });
    default:
      throw new Error(`Unknown action: ${action}`);
  }
}

module.exports = { todoist_tasks, todoist_projects, todoist_sections, todoist_labels };
