/**
 * cloudflare-api-1.0.0 — OpenClaw Skill
 *
 * Cloudflare DNS and zone management. Legacy can list zones,
 * create/update/delete DNS records, and inspect zone settings.
 *
 * Env: CLOUDFLARE_API_TOKEN (scoped: Zone:Read, DNS:Edit)
 */

const https = require("https");

// ── Cloudflare API helper ──────────────────────────────────────────

function cfRequest(method, path, body) {
  const token = process.env.CLOUDFLARE_API_TOKEN;
  if (!token) return Promise.reject(new Error("CLOUDFLARE_API_TOKEN not set"));

  return new Promise((resolve, reject) => {
    const payload = body ? JSON.stringify(body) : null;
    const opts = {
      hostname: "api.cloudflare.com",
      path: `/client/v4${path}`,
      method,
      headers: {
        Authorization: `Bearer ${token}`,
        "Content-Type": "application/json",
        ...(payload ? { "Content-Length": Buffer.byteLength(payload) } : {}),
      },
    };

    const req = https.request(opts, (res) => {
      const chunks = [];
      res.on("data", (c) => chunks.push(c));
      res.on("end", () => {
        try {
          const data = JSON.parse(Buffer.concat(chunks).toString());
          if (!data.success) {
            const msg = data.errors?.[0]?.message || "Unknown Cloudflare error";
            reject(new Error(msg));
          } else {
            resolve(data);
          }
        } catch (e) {
          reject(new Error("Failed to parse Cloudflare response"));
        }
      });
    });

    req.on("error", reject);
    if (payload) req.write(payload);
    req.end();
  });
}

// ── Zone ID resolver (cached per call) ────────────────────────────

async function getZoneId(domain) {
  const data = await cfRequest("GET", `/zones?name=${encodeURIComponent(domain)}`);
  const zone = data.result?.[0];
  if (!zone) throw new Error(`Zone not found: ${domain}`);
  return zone.id;
}

// ── Tools ─────────────────────────────────────────────────────────

/**
 * List all zones (domains) in the account
 */
async function cf_zones() {
  const data = await cfRequest("GET", "/zones?per_page=50");
  const zones = data.result.map((z) => ({
    id: z.id,
    name: z.name,
    status: z.status,
    plan: z.plan?.name || "unknown",
    nameservers: z.name_servers,
    paused: z.paused,
  }));
  return JSON.stringify({ zones, count: zones.length }, null, 2);
}

/**
 * List DNS records for a zone
 */
async function cf_dns_list({ zone, record_type, name }) {
  const zoneId = await getZoneId(zone);
  let path = `/zones/${zoneId}/dns_records?per_page=100`;
  if (record_type) path += `&type=${encodeURIComponent(record_type)}`;
  if (name) {
    const fqdn = name.includes(".") ? name : `${name}.${zone}`;
    path += `&name=${encodeURIComponent(fqdn)}`;
  }

  const data = await cfRequest("GET", path);
  const records = data.result.map((r) => ({
    id: r.id,
    type: r.type,
    name: r.name,
    content: r.content,
    proxied: r.proxied,
    ttl: r.ttl,
  }));
  return JSON.stringify({ zone, records, count: records.length }, null, 2);
}

/**
 * Create a DNS record
 */
async function cf_dns_create({ zone, type, name, content, proxied, ttl }) {
  const zoneId = await getZoneId(zone);
  const body = {
    type: type.toUpperCase(),
    name,
    content,
    proxied: proxied ?? false,
    ttl: ttl || 1,
  };

  const data = await cfRequest("POST", `/zones/${zoneId}/dns_records`, body);
  const r = data.result;
  return JSON.stringify(
    {
      created: true,
      record: {
        id: r.id,
        type: r.type,
        name: r.name,
        content: r.content,
        proxied: r.proxied,
        ttl: r.ttl,
      },
    },
    null,
    2
  );
}

/**
 * Update an existing DNS record
 */
async function cf_dns_update({ zone, record_id, type, name, content, proxied }) {
  const zoneId = await getZoneId(zone);
  const body = {
    type: type.toUpperCase(),
    name,
    content,
    proxied: proxied ?? false,
  };

  const data = await cfRequest(
    "PUT",
    `/zones/${zoneId}/dns_records/${record_id}`,
    body
  );
  const r = data.result;
  return JSON.stringify(
    {
      updated: true,
      record: {
        id: r.id,
        type: r.type,
        name: r.name,
        content: r.content,
        proxied: r.proxied,
      },
    },
    null,
    2
  );
}

/**
 * Delete a DNS record
 */
async function cf_dns_delete({ zone, record_id }) {
  const zoneId = await getZoneId(zone);
  await cfRequest("DELETE", `/zones/${zoneId}/dns_records/${record_id}`);
  return JSON.stringify({ deleted: true, record_id });
}

/**
 * Get zone settings
 */
async function cf_zone_settings({ zone }) {
  const zoneId = await getZoneId(zone);
  const data = await cfRequest("GET", `/zones/${zoneId}/settings`);
  const settings = {};
  for (const s of data.result) {
    settings[s.id] = s.value;
  }
  return JSON.stringify(
    {
      zone,
      ssl: settings.ssl,
      security_level: settings.security_level,
      min_tls_version: settings.min_tls_version,
      always_use_https: settings.always_use_https,
      cache_level: settings.cache_level,
      browser_cache_ttl: settings.browser_cache_ttl,
    },
    null,
    2
  );
}

// ── Exports ───────────────────────────────────────────────────────

module.exports = {
  cf_zones,
  cf_dns_list,
  cf_dns_create,
  cf_dns_update,
  cf_dns_delete,
  cf_zone_settings,
};
