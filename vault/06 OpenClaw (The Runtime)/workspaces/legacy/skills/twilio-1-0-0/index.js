/**
 * twilio-1.0.0 — OpenClaw Skill
 *
 * SMS and voice call delivery via Twilio REST API.
 * Ported from reference/notifications/channels/twilio.py
 *
 * Env vars required:
 *   TWILIO_ACCOUNT_SID
 *   TWILIO_AUTH_TOKEN
 *   TWILIO_FROM_NUMBER      ← your Twilio number e.g. +12015551234
 *   ZEUS_PHONE_NUMBER       ← Zeus's mobile e.g. +13475559876
 *
 * Registers tools: sms_send, voice_call
 */

"use strict";

const https = require("https");

const ACCOUNT_SID  = process.env.TWILIO_ACCOUNT_SID  || "";
const AUTH_TOKEN   = process.env.TWILIO_AUTH_TOKEN   || "";
const FROM_NUMBER  = process.env.TWILIO_FROM_NUMBER  || "";
const ZEUS_PHONE   = process.env.ZEUS_PHONE_NUMBER   || "";

/** POST form-urlencoded to Twilio REST API. */
function twilioRequest(path, formBody) {
  return new Promise((resolve, reject) => {
    if (!ACCOUNT_SID || !AUTH_TOKEN) {
      return reject(new Error("Missing TWILIO_ACCOUNT_SID or TWILIO_AUTH_TOKEN"));
    }

    const auth = Buffer.from(`${ACCOUNT_SID}:${AUTH_TOKEN}`).toString("base64");
    const bodyStr = Object.entries(formBody)
      .map(([k, v]) => `${encodeURIComponent(k)}=${encodeURIComponent(v)}`)
      .join("&");

    const options = {
      hostname: "api.twilio.com",
      path,
      method: "POST",
      headers: {
        Authorization: `Basic ${auth}`,
        "Content-Type": "application/x-www-form-urlencoded",
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
            reject(new Error(`Twilio ${res.statusCode}: ${json.message || data}`));
          }
        } catch {
          reject(new Error(`Twilio parse error: ${data}`));
        }
      });
    });

    req.on("error", reject);
    req.write(bodyStr);
    req.end();
  });
}

/**
 * sms_send — Send an SMS to Zeus (or a custom number).
 * @param {string} message - Text to send (max 1600 chars)
 * @param {string} [to]    - Override recipient (defaults to ZEUS_PHONE_NUMBER)
 */
async function sms_send({ message, to }) {
  const toNumber = to || ZEUS_PHONE;
  if (!toNumber) throw new Error("No recipient number — set ZEUS_PHONE_NUMBER or pass 'to'");
  if (!FROM_NUMBER) throw new Error("TWILIO_FROM_NUMBER not set");

  const body = message.length > 1600 ? message.slice(0, 1597) + "..." : message;

  const result = await twilioRequest(
    `/2010-04-01/Accounts/${ACCOUNT_SID}/Messages.json`,
    { Body: body, From: FROM_NUMBER, To: toNumber }
  );

  return { ok: true, sid: result.sid, to: toNumber };
}

/**
 * voice_call — Initiate a Twilio voice call to Zeus.
 * Uses TwiML to say a short message.
 * @param {string} message - Message to speak (max 200 chars for URL safety)
 * @param {string} [to]    - Override recipient (defaults to ZEUS_PHONE_NUMBER)
 */
async function voice_call({ message, to }) {
  const toNumber = to || ZEUS_PHONE;
  if (!toNumber) throw new Error("No recipient number — set ZEUS_PHONE_NUMBER or pass 'to'");
  if (!FROM_NUMBER) throw new Error("TWILIO_FROM_NUMBER not set");

  // Build inline TwiML
  const safeMsg = encodeURIComponent(message.slice(0, 200));
  const twimlUrl = `https://twimlets.com/message?Message%5B0%5D=${safeMsg}`;

  const result = await twilioRequest(
    `/2010-04-01/Accounts/${ACCOUNT_SID}/Calls.json`,
    { Url: twimlUrl, From: FROM_NUMBER, To: toNumber }
  );

  return { ok: true, sid: result.sid, to: toNumber };
}

// ─── OpenClaw tool registry ───────────────────────────────────────────────────
module.exports = {
  sms_send,
  voice_call,
};
