/**
 * hitl-escalate-1.0.0 — OpenClaw Skill
 *
 * Human-in-the-Loop escalation for decisions requiring Zeus's approval.
 *
 * Implements the escalation ladder:
 *   1. Legacy posts approval request on Discord/Telegram (native channels — no skill needed)
 *   2. If no ack within timeout → hitl_escalate_sms  (this skill)
 *   3. If still no ack         → hitl_escalate_voice (this skill)
 *
 * Designed to be called by Legacy after she has already messaged Zeus on
 * Discord/Telegram and is waiting for a response.
 *
 * Env vars required:
 *   TWILIO_ACCOUNT_SID
 *   TWILIO_AUTH_TOKEN
 *   TWILIO_FROM_NUMBER
 *   ZEUS_PHONE_NUMBER
 *
 * Registers tools: hitl_escalate_sms, hitl_escalate_voice
 */

"use strict";

const https = require("https");

const ACCOUNT_SID = process.env.TWILIO_ACCOUNT_SID || "";
const AUTH_TOKEN  = process.env.TWILIO_AUTH_TOKEN  || "";
const FROM_NUMBER = process.env.TWILIO_FROM_NUMBER || "";
const ZEUS_PHONE  = process.env.ZEUS_PHONE_NUMBER  || "";

/** POST form-urlencoded to Twilio. */
function twilioPost(path, form) {
  return new Promise((resolve, reject) => {
    if (!ACCOUNT_SID || !AUTH_TOKEN) {
      return reject(new Error("Missing TWILIO_ACCOUNT_SID or TWILIO_AUTH_TOKEN"));
    }

    const auth = Buffer.from(`${ACCOUNT_SID}:${AUTH_TOKEN}`).toString("base64");
    const body = Object.entries(form)
      .map(([k, v]) => `${encodeURIComponent(k)}=${encodeURIComponent(v)}`)
      .join("&");

    const req = https.request(
      {
        hostname: "api.twilio.com",
        path,
        method: "POST",
        headers: {
          Authorization: `Basic ${auth}`,
          "Content-Type": "application/x-www-form-urlencoded",
          "Content-Length": Buffer.byteLength(body),
        },
      },
      (res) => {
        let data = "";
        res.on("data", (c) => (data += c));
        res.on("end", () => {
          try {
            const json = JSON.parse(data);
            if (res.statusCode >= 200 && res.statusCode < 300) resolve(json);
            else reject(new Error(`Twilio ${res.statusCode}: ${json.message || data}`));
          } catch {
            reject(new Error(`Twilio parse error: ${data}`));
          }
        });
      }
    );
    req.on("error", reject);
    req.write(body);
    req.end();
  });
}

/**
 * hitl_escalate_sms — Send an urgent SMS escalation to Zeus.
 *
 * Call this when Zeus has not responded to a HITL request on Discord/Telegram
 * within the expected window (typically 5–15 minutes for high-priority items).
 *
 * @param {string} decision     - What decision is pending (1 sentence)
 * @param {string} context      - Why it matters / what happens if ignored
 * @param {string} [options]    - Possible responses (e.g. "Reply YES or NO")
 * @param {string} [urgency]    - "high" | "critical" (default: "high")
 */
async function hitl_escalate_sms({ decision, context, options, urgency = "high" }) {
  if (!ZEUS_PHONE) throw new Error("ZEUS_PHONE_NUMBER not set");
  if (!FROM_NUMBER) throw new Error("TWILIO_FROM_NUMBER not set");

  const tag = urgency === "critical" ? "🔴 CRITICAL" : "🟡 ACTION NEEDED";
  const optionsLine = options ? `\n${options}` : "\nReply on Discord or here.";

  const message = [
    `${tag} — Legacy needs your input`,
    `Decision: ${decision}`,
    context,
    optionsLine,
  ]
    .join("\n")
    .slice(0, 1600);

  const result = await twilioPost(
    `/2010-04-01/Accounts/${ACCOUNT_SID}/Messages.json`,
    { Body: message, From: FROM_NUMBER, To: ZEUS_PHONE }
  );

  return { ok: true, sid: result.sid, channel: "sms", urgency };
}

/**
 * hitl_escalate_voice — Call Zeus to demand immediate attention.
 *
 * Use ONLY as last resort — Zeus has not responded to Discord + SMS within
 * the critical timeout window (typically 30+ minutes for L0/L1 decisions).
 *
 * @param {string} decision - One sentence: what decision is pending
 * @param {string} [urgency] - "critical" (only valid value for voice escalation)
 */
async function hitl_escalate_voice({ decision, urgency = "critical" }) {
  if (!ZEUS_PHONE) throw new Error("ZEUS_PHONE_NUMBER not set");
  if (!FROM_NUMBER) throw new Error("TWILIO_FROM_NUMBER not set");

  const speech = `This is Legacy. Urgent decision needed: ${decision}. Please check Discord immediately.`;
  const safeMsg = encodeURIComponent(speech.slice(0, 200));
  const twimlUrl = `https://twimlets.com/message?Message%5B0%5D=${safeMsg}`;

  const result = await twilioPost(
    `/2010-04-01/Accounts/${ACCOUNT_SID}/Calls.json`,
    { Url: twimlUrl, From: FROM_NUMBER, To: ZEUS_PHONE }
  );

  return { ok: true, sid: result.sid, channel: "voice", urgency };
}

// ─── OpenClaw tool registry ───────────────────────────────────────────────────
module.exports = {
  hitl_escalate_sms,
  hitl_escalate_voice,
};
