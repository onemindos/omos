"use strict";

/**
 * webhook-notifier — OpenClaw Skill (Complete Working Example)
 *
 * Sends webhook notifications to any HTTP endpoint.
 * Includes specialized helpers for Discord and Slack webhooks.
 * Demonstrates: retry logic, pre-formatted payloads, flexible headers.
 *
 * Registered tools: webhook_send, webhook_discord, webhook_slack, webhook_ping
 */

// ─── Helpers ─────────────────────────────────────────────────────────────────

const DISCORD_COLORS = {
  green: 0x57f287,
  red: 0xed4245,
  yellow: 0xfee75c,
  blue: 0x5865f2,
  purple: 0x9b59b6,
  gray: 0x95a5a6,
};

/**
 * Core POST function with one retry on server error or timeout.
 *
 * @param {string} url - Target URL
 * @param {Object} body - JSON payload
 * @param {Object} [extraHeaders={}] - Additional headers
 * @returns {Object} { status, ok, response_code, body }
 */
async function postWithRetry(url, body, extraHeaders = {}) {
  const headers = {
    "Content-Type": "application/json",
    "User-Agent": "OpenClaw-Legacy/1.0",
    ...extraHeaders,
  };

  for (let attempt = 1; attempt <= 2; attempt++) {
    let res;
    try {
      res = await fetch(url, {
        method: "POST",
        headers,
        body: JSON.stringify(body),
        signal: AbortSignal.timeout(15000),
      });
    } catch (err) {
      if (attempt === 2) throw new Error(`Webhook request failed after retry: ${err.message}`);
      await new Promise((r) => setTimeout(r, 1000)); // wait 1s before retry
      continue;
    }

    // 4xx = client error, don't retry
    if (res.status >= 400 && res.status < 500) {
      const text = await res.text().catch(() => "");
      return {
        ok: false,
        status: "client_error",
        response_code: res.status,
        error: text || `HTTP ${res.status}`,
      };
    }

    // 5xx = server error, retry once
    if (res.status >= 500 && attempt < 2) {
      await new Promise((r) => setTimeout(r, 1500));
      continue;
    }

    const ct = res.headers.get("content-type") || "";
    let responseBody;
    try {
      responseBody = ct.includes("application/json") ? await res.json() : await res.text();
    } catch {
      responseBody = null;
    }

    return {
      ok: res.ok,
      status: res.ok ? "sent" : "server_error",
      response_code: res.status,
      body: responseBody,
    };
  }
}

// ─── Tools ───────────────────────────────────────────────────────────────────

/**
 * TOOL: webhook_send
 * Send a POST request to any webhook URL with a JSON payload.
 *
 * @param {string} url - Webhook URL
 * @param {Object} payload - JSON payload to send
 * @param {Object} [headers={}] - Optional extra headers (e.g., auth headers)
 * @returns {Object} { ok, status, response_code, body }
 */
async function webhook_send({ url, payload = {}, headers = {} }) {
  if (!url) throw new Error("url is required");
  if (typeof payload !== "object") throw new Error("payload must be a JSON object");

  return postWithRetry(url, payload, headers);
}

/**
 * TOOL: webhook_discord
 * Send a formatted embed message to a Discord webhook.
 *
 * @param {string} url - Discord webhook URL
 * @param {string} title - Embed title
 * @param {string} [description] - Embed description / body text
 * @param {string} [color="blue"] - Color: green | red | yellow | blue | purple | gray
 * @param {string} [username="Legacy"] - Bot display name
 * @param {Array}  [fields] - Optional array of { name, value, inline } field objects
 * @param {string} [footer] - Optional footer text
 * @returns {Object} Send result
 */
async function webhook_discord({ url, title, description, color = "blue", username = "Legacy", fields = [], footer }) {
  if (!url) throw new Error("url is required");
  if (!title) throw new Error("title is required");

  const colorHex = DISCORD_COLORS[color] ?? DISCORD_COLORS.blue;

  const embed = {
    title,
    description,
    color: colorHex,
    timestamp: new Date().toISOString(),
  };

  if (fields.length > 0) embed.fields = fields;
  if (footer) embed.footer = { text: footer };

  const payload = {
    username,
    embeds: [embed],
  };

  return postWithRetry(url, payload);
}

/**
 * TOOL: webhook_slack
 * Send a message to a Slack incoming webhook.
 *
 * @param {string} url - Slack webhook URL
 * @param {string} text - Message text (supports Slack mrkdwn)
 * @param {string} [username="Legacy"] - Bot display name
 * @param {string} [icon_emoji=":robot_face:"] - Emoji for the bot icon
 * @param {Array}  [attachments=[]] - Optional Slack attachment blocks
 * @returns {Object} Send result
 */
async function webhook_slack({ url, text, username = "Legacy", icon_emoji = ":robot_face:", attachments = [] }) {
  if (!url) throw new Error("url is required");
  if (!text) throw new Error("text is required");

  const payload = {
    text,
    username,
    icon_emoji,
  };

  if (attachments.length > 0) payload.attachments = attachments;

  return postWithRetry(url, payload);
}

/**
 * TOOL: webhook_ping
 * Ping a URL to check if it's alive and responding.
 *
 * @param {string} url - URL to ping
 * @param {string} [method="GET"] - HTTP method to use (GET or HEAD)
 * @returns {Object} { reachable, response_code, latency_ms, url }
 */
async function webhook_ping({ url, method = "GET" }) {
  if (!url) throw new Error("url is required");

  const start = Date.now();
  try {
    const res = await fetch(url, {
      method,
      signal: AbortSignal.timeout(10000),
    });
    return {
      reachable: true,
      response_code: res.status,
      ok: res.ok,
      latency_ms: Date.now() - start,
      url,
    };
  } catch (err) {
    return {
      reachable: false,
      error: err.message,
      latency_ms: Date.now() - start,
      url,
    };
  }
}

// ─── Exports ─────────────────────────────────────────────────────────────────
module.exports = { webhook_send, webhook_discord, webhook_slack, webhook_ping };
