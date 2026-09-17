/**
 * zoho-mail-1.0.0 — OpenClaw Skill
 *
 * Read, search, move, and send emails via Zoho Mail's MCP endpoint.
 * All operations go through the Zoho MCP HTTP API (JSON-RPC 2.0).
 *
 * Required env: ZOHO_MCP_URL
 * Optional env: ZOHO_ACCOUNT_ID (defaults to Zeus's account)
 *
 * Tools: zoho_list_folders, zoho_list_emails, zoho_search_emails,
 *        zoho_read_email, zoho_move_emails, zoho_mark_read,
 *        zoho_send_email, zoho_get_account_info
 */

"use strict";

const https = require("https");
const { URL } = require("url");

// ---------------------------------------------------------------------------
// Config
// ---------------------------------------------------------------------------

const DEFAULT_ACCOUNT_ID = "8698543000000009002";

// Known folder name → ID mapping (pre-seeded, refreshed by zoho_list_folders)
const FOLDER_MAP = {
  inbox:        "8698543000000009014",
  drafts:       "8698543000000009016",
  templates:    "8698543000000009018",
  snoozed:      "8698543000000009020",
  sent:         "8698543000000009022",
  spam:         "8698543000000009024",
  trash:        "8698543000000009026",
  notification: "8698543000000590137",
  newsletter:   "8698543000000590147",
  hp:           "8698543000000596001",
  "holistic performance": "8698543000000596001",
  le:           "8698543000000597001",
  "legacy evolution": "8698543000000597001",
  ui:           "8698543000000598001",
  "unified intelligence": "8698543000000598001",
  ge:           "8698543000000599001",
  "generational entrepreneurship": "8698543000000599001",
};

// Send-as identity map → sendMailId
const SEND_IDENTITY_MAP = {
  hp:       "8698543000000594001",
  le:       "8698543000000594003",
  ui:       "8698543000000594005",
  ge:       "8698543000000594007",
  personal: "8698543000000594009",
  onemind:  "8698543000000594011",
};

const SEND_FROM_EMAIL = {
  hp: "zeus@onemindcodex.com",
  le: "zeus@onemindcodex.com",
  ui: "zeus@onemindcodex.com",
  ge: "zeus@onemindcodex.com",
  personal: "zeus@onemindcodex.com",
  onemind: "finance@legacyhq.ai",
};

// ---------------------------------------------------------------------------
// MCP HTTP client
// ---------------------------------------------------------------------------

function getMcpUrl() {
  const url = process.env.ZOHO_MCP_URL;
  if (!url) throw new Error("ZOHO_MCP_URL environment variable is not set");
  return url;
}

function getAccountId() {
  return process.env.ZOHO_ACCOUNT_ID || DEFAULT_ACCOUNT_ID;
}

function mcpRequest(mcpUrl, payload) {
  return new Promise((resolve, reject) => {
    const parsed = new URL(mcpUrl);
    const body = JSON.stringify(payload);

    const opts = {
      hostname: parsed.hostname,
      path: parsed.pathname + parsed.search,
      method: "POST",
      headers: {
        "Content-Type": "application/json",
        "Content-Length": Buffer.byteLength(body),
      },
      timeout: 20000,
    };

    const req = https.request(opts, (res) => {
      let data = "";
      res.setEncoding("utf8");
      res.on("data", (c) => (data += c));
      res.on("end", () => {
        try {
          const parsed = JSON.parse(data);
          if (parsed.error) return reject(new Error(parsed.error.message || JSON.stringify(parsed.error)));
          resolve(parsed);
        } catch {
          reject(new Error(`Invalid JSON response: ${data.slice(0, 200)}`));
        }
      });
    });

    req.on("timeout", () => { req.destroy(); reject(new Error("MCP request timed out")); });
    req.on("error", reject);
    req.write(body);
    req.end();
  });
}

let _reqId = 0;
async function toolCall(toolName, args) {
  const mcpUrl = getMcpUrl();
  const payload = {
    jsonrpc: "2.0",
    id: ++_reqId,
    method: "tools/call",
    params: { name: toolName, arguments: args },
  };
  const response = await mcpRequest(mcpUrl, payload);
  const content = response.result?.content?.[0]?.text;
  if (!content) throw new Error(`Empty response from ${toolName}`);
  try {
    return JSON.parse(content);
  } catch {
    throw new Error(`Non-JSON response from ${toolName}: ${content.slice(0, 200)}`);
  }
}

// ---------------------------------------------------------------------------
// Helpers
// ---------------------------------------------------------------------------

function resolveFolderId(folder) {
  if (!folder) return null;
  // If it's already a numeric ID
  if (/^\d{10,}$/.test(folder)) return folder;
  // Lookup by name
  const key = folder.toLowerCase().trim();
  return FOLDER_MAP[key] || null;
}

function stripHtml(html) {
  if (!html) return "";
  return html
    .replace(/<br\s*\/?>/gi, "\n")
    .replace(/<\/p>/gi, "\n\n")
    .replace(/<[^>]+>/g, "")
    .replace(/&nbsp;/g, " ")
    .replace(/&amp;/g, "&")
    .replace(/&lt;/g, "<")
    .replace(/&gt;/g, ">")
    .replace(/&quot;/g, '"')
    .replace(/&#39;/g, "'")
    .replace(/\n{3,}/g, "\n\n")
    .trim();
}

function truncate(text, max = 2000) {
  if (!text || text.length <= max) return text || "";
  return text.slice(0, max) + "…";
}

// ---------------------------------------------------------------------------
// Tool: zoho_list_folders
// ---------------------------------------------------------------------------

async function zoho_list_folders() {
  try {
    const accountId = getAccountId();
    const data = await toolCall("ZohoMail_getAllFolders", {
      path_variables: { accountId },
    });

    const folders = (data.data || []).map((f) => ({
      id: f.folderId,
      name: f.folderName || f.path,
      path: f.path,
      // messageCount not returned by getAllFolders — use zoho_list_emails to get counts
    }));

    // Refresh local map
    for (const f of folders) {
      const key = (f.name || "").toLowerCase().replace(/^\/inbox\//, "").trim();
      if (key) FOLDER_MAP[key] = f.id;
    }

    return JSON.stringify({ folders, total_folders: folders.length }, null, 2);
  } catch (e) {
    return JSON.stringify({ error: e.message });
  }
}

// ---------------------------------------------------------------------------
// Tool: zoho_list_emails
// ---------------------------------------------------------------------------

async function zoho_list_emails({ folder, limit, offset, unread_only } = {}) {
  try {
    const accountId = getAccountId();
    const folderId = resolveFolderId(folder);
    if (!folderId) {
      return JSON.stringify({ error: `Unknown folder: '${folder}'. Use zoho_list_folders to see available folders.` });
    }

    const count = Math.min(limit || 20, 200);
    const start = (offset || 0) + 1; // Zoho is 1-based

    const args = {
      path_variables: { accountId },
      query_params: {
        folderId,
        limit: count,
        start,
        ...(unread_only ? { status: "unread" } : {}),
      },
    };

    const data = await toolCall("ZohoMail_listEmails", args);
    const emails = (data.data || []).map((m) => ({
      id: m.messageId,
      subject: m.subject || "(no subject)",
      from: m.fromAddress,
      from_name: m.sender || "",
      date: m.receivedTime || m.sentDateInGMT,
      read: m.isRead === "true" || m.isRead === true,
      has_attachment: m.hasAttachment === "true" || m.hasAttachment === true,
      size: m.size,
    }));

    return JSON.stringify({
      folder,
      folder_id: folderId,
      count: emails.length,
      offset: start - 1,
      emails,
    }, null, 2);
  } catch (e) {
    return JSON.stringify({ error: e.message });
  }
}

// ---------------------------------------------------------------------------
// Tool: zoho_search_emails
// ---------------------------------------------------------------------------

async function zoho_search_emails({ query, limit } = {}) {
  try {
    const accountId = getAccountId();
    const count = Math.min(limit || 20, 200);

    const data = await toolCall("ZohoMail_SearchEmails", {
      path_variables: { accountId },
      query_params: {
        searchKey: query,
        limit: count,
        start: 1,
      },
    });

    const emails = (data.data || []).map((m) => ({
      id: m.messageId,
      subject: m.subject || "(no subject)",
      from: m.fromAddress,
      from_name: m.sender || "",
      folder: m.folderId,
      date: m.receivedTime || m.sentDateInGMT,
      read: m.isRead === "true" || m.isRead === true,
    }));

    return JSON.stringify({
      query,
      count: emails.length,
      emails,
    }, null, 2);
  } catch (e) {
    return JSON.stringify({ error: e.message });
  }
}

// ---------------------------------------------------------------------------
// Tool: zoho_read_email
// ---------------------------------------------------------------------------

async function zoho_read_email({ message_id } = {}) {
  try {
    const accountId = getAccountId();

    const data = await toolCall("ZohoMail_getMessageContent", {
      path_variables: { accountId, messageId: message_id },
    });

    const msg = data.data || data;
    // Zoho returns content as a string field or nested object
    const bodyHtml = msg.htmlBody || msg.content || (typeof msg === "string" ? msg : "");
    const bodyText = msg.textBody || msg.text || stripHtml(bodyHtml);

    // Also try to fetch headers for from/to/subject
    let meta = {};
    try {
      const hdr = await toolCall("ZohoMail_getMessageDetails", {
        path_variables: { accountId, messageId: message_id },
      });
      meta = hdr.data || {};
    } catch { /* ok if headers fail */ }

    return JSON.stringify({
      id: message_id,
      subject: meta.subject || msg.subject || "",
      from: meta.fromAddress || msg.fromAddress || "",
      from_name: meta.sender || msg.sender || "",
      to: meta.toAddress || msg.toAddress || "",
      cc: meta.ccAddress || msg.ccAddress || "",
      date: meta.receivedTime || msg.receivedTime || "",
      has_attachment: meta.hasAttachment || msg.hasAttachment,
      body: truncate(bodyText || JSON.stringify(msg), 3000),
    }, null, 2);
  } catch (e) {
    return JSON.stringify({ error: e.message });
  }
}

// ---------------------------------------------------------------------------
// Tool: zoho_move_emails
// ---------------------------------------------------------------------------

async function zoho_move_emails({ message_ids, to_folder } = {}) {
  try {
    if (!message_ids || message_ids.length === 0) {
      return JSON.stringify({ error: "message_ids is required and must be non-empty" });
    }

    const accountId = getAccountId();
    const folderId = resolveFolderId(to_folder);
    if (!folderId) {
      return JSON.stringify({ error: `Unknown folder: '${to_folder}'. Use zoho_list_folders to see available folders.` });
    }

    const data = await toolCall("ZohoMail_moveMessages", {
      path_variables: { accountId },
      body: {
        messageIds: message_ids,
        toFolderId: folderId,
        mode: "move",
      },
    });

    return JSON.stringify({
      success: true,
      moved: message_ids.length,
      to_folder,
      to_folder_id: folderId,
      response: data.status || data,
    }, null, 2);
  } catch (e) {
    return JSON.stringify({ error: e.message });
  }
}

// ---------------------------------------------------------------------------
// Tool: zoho_mark_read
// ---------------------------------------------------------------------------

async function zoho_mark_read({ message_ids, read = true } = {}) {
  try {
    if (!message_ids || message_ids.length === 0) {
      return JSON.stringify({ error: "message_ids is required" });
    }
    const accountId = getAccountId();

    const toolName = read ? "ZohoMail_readMessages" : "ZohoMail_unreadMessage";
    const data = await toolCall(toolName, {
      path_variables: { accountId },
      body: { messageIds: message_ids },
    });

    return JSON.stringify({
      success: true,
      marked_as: read ? "read" : "unread",
      count: message_ids.length,
      response: data.status || data,
    }, null, 2);
  } catch (e) {
    return JSON.stringify({ error: e.message });
  }
}

// ---------------------------------------------------------------------------
// Tool: zoho_send_email
// ---------------------------------------------------------------------------

async function zoho_send_email({ to, subject, body, from_identity, cc } = {}) {
  try {
    if (!to || !subject || !body) {
      return JSON.stringify({ error: "to, subject, and body are required" });
    }

    const accountId = getAccountId();
    const fromAddress = SEND_FROM_EMAIL[from_identity] || "zeus@zeusdelacruz.com";

    const bodyArgs = {
      fromAddress,
      toAddress: to,
      subject,
      content: body,
      mailFormat: body.includes("<") ? "html" : "plaintext",
      ...(cc ? { ccAddress: cc } : {}),
    };

    const data = await toolCall("ZohoMail_sendEmail", {
      path_variables: { accountId },
      body: bodyArgs,
    });

    return JSON.stringify({
      success: true,
      from: fromAddress,
      to,
      subject,
      identity: from_identity || "primary",
      response: data.status || data,
    }, null, 2);
  } catch (e) {
    return JSON.stringify({ error: e.message });
  }
}

// ---------------------------------------------------------------------------
// Tool: zoho_get_account_info
// ---------------------------------------------------------------------------

async function zoho_get_account_info() {
  // Static account info (getMailAccounts requires org scope not available on this key)
  return JSON.stringify({
    primary_email: "zeus@zeusdelacruz.com",
    account_id: getAccountId(),
    account_name: "legacyplatform",
    display_name: "Zeus",
    role: "super_admin",
    tfa_enabled: true,
    imap_enabled: true,
    timezone: "America/New_York",
    aliases: [
      { address: "zeus@zeusdelacruz.com", primary: true },
      { address: "zeus@onemindcodex.com", alias: true },
      { address: "hp@zeusdelacruz.me", alias: true },
      { address: "ui@zeusdelacruz.me", alias: true },
      { address: "le@zeusdelacruz.me", alias: true },
      { address: "ge@zeusdelacruz.me", alias: true },
      { address: "zeus.vaults@pm.me", external: true },
    ],
    send_identities: [
      { key: "hp", name: "Holistic Performance", from: "zeus@onemindcodex.com" },
      { key: "le", name: "Legacy Evolution", from: "zeus@onemindcodex.com" },
      { key: "ui", name: "Unified Intelligence", from: "zeus@onemindcodex.com" },
      { key: "ge", name: "Generational Entrepreneurship", from: "zeus@onemindcodex.com" },
      { key: "personal", name: "Zeus (Personal Brand)", from: "zeus@onemindcodex.com" },
      { key: "onemind", name: "Zeus (OneMind)", from: "finance@legacyhq.ai" },
    ],
    note: "getMailAccounts requires org-admin scope. Account details above are pre-seeded from last audit (2026-03-12).",
  }, null, 2);
}

// ---------------------------------------------------------------------------
// Exports
// ---------------------------------------------------------------------------

module.exports = {
  zoho_list_folders,
  zoho_list_emails,
  zoho_search_emails,
  zoho_read_email,
  zoho_move_emails,
  zoho_mark_read,
  zoho_send_email,
  zoho_get_account_info,
};
