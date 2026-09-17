/**
 * lob-1.0.0 — OpenClaw Skill
 *
 * Send physical mail (letters, postcards) and verify addresses via the Lob API.
 * USPS first-class mail, address verification, delivery tracking.
 *
 * Env vars: LOB_API_KEY
 *
 * Registers tools: lob
 */

"use strict";

const https = require("https");

// ---------------------------------------------------------------------------
// HTTP helper — Lob uses HTTP Basic auth (API key as username, empty password)
// ---------------------------------------------------------------------------

function lobRequest(method, path, { data, params } = {}) {
  return new Promise((resolve, reject) => {
    const apiKey = process.env.LOB_API_KEY || "";
    if (!apiKey) return resolve({ error: "LOB_API_KEY not set" });

    let qs = "";
    if (params) {
      qs = "?" + Object.entries(params).map(([k, v]) => `${encodeURIComponent(k)}=${encodeURIComponent(v)}`).join("&");
    }

    // Lob uses form-encoded POST data
    let payload = null;
    let contentType = "application/json";
    if (data) {
      payload = Object.entries(data)
        .map(([k, v]) => `${encodeURIComponent(k)}=${encodeURIComponent(v)}`)
        .join("&");
      contentType = "application/x-www-form-urlencoded";
    }

    const auth = Buffer.from(`${apiKey}:`).toString("base64");

    const opts = {
      hostname: "api.lob.com",
      path: `/v1${path}${qs}`,
      method: method.toUpperCase(),
      headers: {
        Authorization: `Basic ${auth}`,
        "Content-Type": contentType,
        "User-Agent": "OpenClaw-Lob/1.0",
      },
      timeout: 30000,
    };

    if (payload) opts.headers["Content-Length"] = Buffer.byteLength(payload);

    const req = https.request(opts, (res) => {
      let body = "";
      res.setEncoding("utf8");
      res.on("data", (c) => (body += c));
      res.on("end", () => {
        const code = res.statusCode;
        if (code >= 200 && code < 300) {
          return resolve(body ? JSON.parse(body) : { success: true });
        }
        let detail = body;
        try {
          const parsed = JSON.parse(body);
          detail = (parsed.error || {}).message || body;
        } catch {}
        resolve({ error: `API error: ${code}`, details: detail });
      });
    });
    req.on("timeout", () => { req.destroy(); resolve({ error: "Request timed out" }); });
    req.on("error", (e) => resolve({ error: e.message }));
    if (payload) req.write(payload);
    req.end();
  });
}

function isTestKey() {
  return (process.env.LOB_API_KEY || "").startsWith("test_");
}

// ---------------------------------------------------------------------------
// Tool: lob
// ---------------------------------------------------------------------------

async function lob({
  action,
  // letter fields
  to_name, to_address_line1, to_address_line2, to_city, to_state, to_zip, to_country,
  from_name, from_address_line1, from_address_line2, from_city, from_state, from_zip, from_country,
  content, file_url, color, double_sided, mail_type, use_type,
  // verify fields
  address_line1, city, state, zip_code,
  // tracking
  letter_id,
}) {
  to_country = to_country || "US";
  from_country = from_country || "US";
  color = color || false;
  double_sided = double_sided !== false;
  mail_type = mail_type || "usps_first_class";
  use_type = use_type || "operational";

  try {
    // --- SEND LETTER ---
    if (action === "send_letter") {
      if (!to_name || !to_address_line1 || !to_city || !to_state || !to_zip) {
        return JSON.stringify({ error: "Missing recipient address fields (to_name, to_address_line1, to_city, to_state, to_zip)" });
      }
      if (!from_name || !from_address_line1 || !from_city || !from_state || !from_zip) {
        return JSON.stringify({ error: "Missing sender address fields (from_name, from_address_line1, from_city, from_state, from_zip)" });
      }
      if (!content && !file_url) {
        return JSON.stringify({ error: "Either content (HTML string) or file_url (PDF URL) required" });
      }

      const formData = {
        "to[name]": to_name,
        "to[address_line1]": to_address_line1,
        "to[address_city]": to_city,
        "to[address_state]": to_state,
        "to[address_zip]": to_zip,
        "to[address_country]": to_country,
        "from[name]": from_name,
        "from[address_line1]": from_address_line1,
        "from[address_city]": from_city,
        "from[address_state]": from_state,
        "from[address_zip]": from_zip,
        "from[address_country]": from_country,
        color: String(color),
        double_sided: String(double_sided),
        mail_type,
        use_type,
      };

      if (to_address_line2) formData["to[address_line2]"] = to_address_line2;
      if (from_address_line2) formData["from[address_line2]"] = from_address_line2;

      if (file_url) {
        formData.file = file_url;
      } else {
        formData.file = content.trim().startsWith("<html") ? content : `<html><body>${content}</body></html>`;
      }

      const result = await lobRequest("POST", "/letters", { data: formData });
      if (result.error) return JSON.stringify(result, null, 2);

      return JSON.stringify({
        action: "send_letter",
        success: true,
        test_mode: isTestKey(),
        letter_id: result.id,
        url: result.url,
        to: (result.to || {}).name,
        expected_delivery: result.expected_delivery_date,
        tracking_number: result.tracking_number,
      }, null, 2);
    }

    // --- SEND POSTCARD ---
    if (action === "send_postcard") {
      if (!to_name || !to_address_line1 || !to_city || !to_state || !to_zip) {
        return JSON.stringify({ error: "Missing recipient address fields" });
      }

      const formData = {
        "to[name]": to_name,
        "to[address_line1]": to_address_line1,
        "to[address_city]": to_city,
        "to[address_state]": to_state,
        "to[address_zip]": to_zip,
        "to[address_country]": to_country,
      };

      if (from_name) {
        formData["from[name]"] = from_name;
        formData["from[address_line1]"] = from_address_line1 || "";
        formData["from[address_city]"] = from_city || "";
        formData["from[address_state]"] = from_state || "";
        formData["from[address_zip]"] = from_zip || "";
        formData["from[address_country]"] = from_country;
      }

      if (file_url) {
        formData.front = file_url;
      }
      if (content) {
        formData.back = content.trim().startsWith("<html") ? content : `<html><body>${content}</body></html>`;
      }

      const result = await lobRequest("POST", "/postcards", { data: formData });
      if (result.error) return JSON.stringify(result, null, 2);

      return JSON.stringify({
        action: "send_postcard",
        success: true,
        test_mode: isTestKey(),
        postcard_id: result.id,
        url: result.url,
        expected_delivery: result.expected_delivery_date,
      }, null, 2);
    }

    // --- VERIFY ADDRESS ---
    if (action === "verify_address") {
      if (!address_line1 || !city || !state || !zip_code) {
        return JSON.stringify({ error: "address_line1, city, state, zip_code required" });
      }

      const result = await lobRequest("POST", "/us_verifications", {
        data: { primary_line: address_line1, city, state, zip_code },
      });
      if (result.error) return JSON.stringify(result, null, 2);

      const components = result.components || {};
      return JSON.stringify({
        action: "verify_address",
        deliverability: result.deliverability,
        primary_line: result.primary_line,
        city: components.city,
        state: components.state,
        zip_code: components.zip_code,
      }, null, 2);
    }

    // --- GET LETTER ---
    if (action === "get_letter") {
      if (!letter_id) return JSON.stringify({ error: "letter_id required" });
      const result = await lobRequest("GET", `/letters/${encodeURIComponent(letter_id)}`);
      if (result.error) return JSON.stringify(result, null, 2);

      return JSON.stringify({
        action: "get_letter",
        letter_id: result.id,
        to: (result.to || {}).name,
        tracking_number: result.tracking_number,
        tracking_events: result.tracking_events || [],
        expected_delivery: result.expected_delivery_date,
      }, null, 2);
    }

    // --- LIST LETTERS ---
    if (action === "list_letters") {
      const result = await lobRequest("GET", "/letters", { params: { limit: 10 } });
      if (result.error) return JSON.stringify(result, null, 2);

      return JSON.stringify({
        action: "list_letters",
        count: (result.data || []).length,
        letters: (result.data || []).map((l) => ({
          id: l.id,
          to: (l.to || {}).name,
          expected_delivery: l.expected_delivery_date,
        })),
      }, null, 2);
    }

    // --- CANCEL LETTER ---
    if (action === "cancel_letter") {
      if (!letter_id) return JSON.stringify({ error: "letter_id required" });
      const result = await lobRequest("DELETE", `/letters/${encodeURIComponent(letter_id)}`);
      if (result.error) return JSON.stringify(result, null, 2);
      return JSON.stringify({ action: "cancel_letter", success: true, letter_id }, null, 2);
    }

    return JSON.stringify({
      error: `Unknown action: ${action}`,
      valid_actions: ["send_letter", "send_postcard", "verify_address", "get_letter", "list_letters", "cancel_letter"],
    }, null, 2);
  } catch (e) {
    return JSON.stringify({ error: e.message });
  }
}

// ---------------------------------------------------------------------------
// Exports
// ---------------------------------------------------------------------------

module.exports = { lob };
