/**
 * resend-1.0.0 — OpenClaw Skill
 *
 * Email delivery via Resend API (https://resend.com).
 * Ported from reference/notifications/channels/resend.py
 *
 * Env vars required:
 *   RESEND_API_KEY          ← from resend.com dashboard
 *   RESEND_FROM_EMAIL       ← e.g. "Legacy <notifications@onemindos.dev>"
 *   ZEUS_EMAIL              ← Zeus's email address
 *
 * Registers tools: email_send
 */

"use strict";

const https = require("https");

const API_KEY    = process.env.RESEND_API_KEY    || "";
const FROM_EMAIL = process.env.RESEND_FROM_EMAIL || "Legacy <notifications@onemindos.dev>";
const ZEUS_EMAIL = process.env.ZEUS_EMAIL        || "";

/** POST JSON to Resend API. */
function resendRequest(payload) {
  return new Promise((resolve, reject) => {
    if (!API_KEY) return reject(new Error("RESEND_API_KEY not set"));

    const bodyStr = JSON.stringify(payload);
    const options = {
      hostname: "api.resend.com",
      path: "/emails",
      method: "POST",
      headers: {
        Authorization: `Bearer ${API_KEY}`,
        "Content-Type": "application/json",
        "Content-Length": Buffer.byteLength(bodyStr),
      },
    };

    const req = https.request(options, (res) => {
      let data = "";
      res.on("data", (c) => (data += c));
      res.on("end", () => {
        try {
          const json = JSON.parse(data);
          if (res.statusCode >= 200 && res.statusCode < 300) {
            resolve(json);
          } else {
            reject(new Error(`Resend ${res.statusCode}: ${json.message || data}`));
          }
        } catch {
          reject(new Error(`Resend parse error: ${data}`));
        }
      });
    });

    req.on("error", reject);
    req.write(bodyStr);
    req.end();
  });
}

/** Build a priority prefix for the subject line. */
function priorityPrefix(priority) {
  if (priority === "critical") return "[CRITICAL] ";
  if (priority === "high")     return "[HIGH] ";
  return "";
}

/** Build a minimal HTML email body. */
function buildHtml(subject, body, actionUrl) {
  const buttonHtml = actionUrl
    ? `<div style="margin-top:20px;">
        <a href="${actionUrl}" style="display:inline-block;padding:10px 20px;background-color:#000;color:#fff;text-decoration:none;border-radius:4px;">Open</a>
       </div>`
    : "";

  return `<!DOCTYPE html>
<html>
<head><meta charset="utf-8"><meta name="viewport" content="width=device-width,initial-scale=1.0"></head>
<body style="font-family:-apple-system,BlinkMacSystemFont,'Segoe UI',Roboto,sans-serif;margin:0;padding:0;background:#f5f5f5;">
  <div style="max-width:600px;margin:0 auto;padding:20px;">
    <div style="background:#fff;border-radius:8px;padding:24px;box-shadow:0 2px 4px rgba(0,0,0,0.1);">
      <h1 style="margin:0 0 16px;font-size:20px;color:#000;">${subject}</h1>
      ${body ? `<p style="margin:0 0 16px;color:#333;line-height:1.5;">${body}</p>` : ""}
      ${buttonHtml}
      <div style="margin-top:24px;padding-top:16px;border-top:1px solid #eee;font-size:12px;color:#666;">
        Sent by Legacy · OneMind OS
      </div>
    </div>
  </div>
</body>
</html>`;
}

/**
 * email_send — Send an email via Resend.
 * @param {string} subject      - Email subject (priority prefix auto-added)
 * @param {string} body         - Plain-text body (rendered in HTML)
 * @param {string} [to]         - Recipient email (defaults to ZEUS_EMAIL)
 * @param {string} [priority]   - "critical" | "high" | "normal" | "low"
 * @param {string} [action_url] - Optional CTA button URL
 */
async function email_send({ subject, body, to, priority = "normal", action_url }) {
  const toEmail = to || ZEUS_EMAIL;
  if (!toEmail) throw new Error("No recipient — set ZEUS_EMAIL or pass 'to'");

  const fullSubject = `${priorityPrefix(priority)}${subject}`;
  const html = buildHtml(subject, body || "", action_url || "");

  const result = await resendRequest({
    from: FROM_EMAIL,
    to: [toEmail],
    subject: fullSubject,
    html,
  });

  return { ok: true, id: result.id, to: toEmail };
}

// ─── OpenClaw tool registry ───────────────────────────────────────────────────
module.exports = {
  email_send,
};
