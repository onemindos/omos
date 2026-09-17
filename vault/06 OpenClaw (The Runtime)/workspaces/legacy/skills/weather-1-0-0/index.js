/**
 * weather-1.0.0 — OpenClaw Skill
 *
 * Get current weather conditions and forecasts for any location.
 * Uses wttr.in (primary) and Open-Meteo (structured JSON fallback).
 * No API keys required.
 *
 * Registers tools: get_weather, get_forecast
 */

"use strict";

const https = require("https");
const { URL } = require("url");

// ---------------------------------------------------------------------------
// HTTP helper
// ---------------------------------------------------------------------------

function httpsGet(urlStr, options = {}) {
  return new Promise((resolve, reject) => {
    let parsed;
    try { parsed = new URL(urlStr); } catch { return reject(new Error(`Invalid URL: ${urlStr}`)); }

    const opts = {
      hostname: parsed.hostname,
      path: parsed.pathname + parsed.search,
      method: "GET",
      headers: { "User-Agent": "OpenClaw-Weather/1.0", ...options.headers },
      timeout: 15000,
    };

    const req = https.request(opts, (res) => {
      if (res.statusCode < 200 || res.statusCode >= 300) {
        return reject(new Error(`HTTP ${res.statusCode} for ${urlStr}`));
      }
      let data = "";
      res.setEncoding("utf8");
      res.on("data", (c) => (data += c));
      res.on("end", () => resolve(data));
    });
    req.on("timeout", () => { req.destroy(); reject(new Error("Request timed out")); });
    req.on("error", reject);
    req.end();
  });
}

// ---------------------------------------------------------------------------
// Weather code descriptions (WMO)
// ---------------------------------------------------------------------------

const WMO_CODES = {
  0: "Clear sky", 1: "Mainly clear", 2: "Partly cloudy", 3: "Overcast",
  45: "Fog", 48: "Rime fog",
  51: "Light drizzle", 53: "Moderate drizzle", 55: "Dense drizzle",
  61: "Light rain", 63: "Moderate rain", 65: "Heavy rain",
  66: "Light freezing rain", 67: "Heavy freezing rain",
  71: "Light snow", 73: "Moderate snow", 75: "Heavy snow",
  77: "Snow grains",
  80: "Light rain showers", 81: "Moderate rain showers", 82: "Violent rain showers",
  85: "Light snow showers", 86: "Heavy snow showers",
  95: "Thunderstorm", 96: "Thunderstorm with slight hail", 99: "Thunderstorm with heavy hail",
};

function fmtTemp(c) {
  const f = (c * 9 / 5) + 32;
  return `${c.toFixed(1)}°C (${f.toFixed(0)}°F)`;
}

// ---------------------------------------------------------------------------
// Geocode via Open-Meteo Geocoding API (free, no key)
// ---------------------------------------------------------------------------

async function geocode(location) {
  const encoded = encodeURIComponent(location);
  const url = `https://geocoding-api.open-meteo.com/v1/search?name=${encoded}&count=1&language=en&format=json`;
  const raw = await httpsGet(url);
  const data = JSON.parse(raw);
  const results = data.results || [];
  if (results.length === 0) return null;
  const r = results[0];
  return { lat: r.latitude, lng: r.longitude, name: r.name, country: r.country, admin1: r.admin1 || "" };
}

// ---------------------------------------------------------------------------
// Tool: get_weather
// ---------------------------------------------------------------------------

async function get_weather({ location }) {
  // Try wttr.in first for a quick human-readable summary
  try {
    const encoded = encodeURIComponent(location);
    const wttrUrl = `https://wttr.in/${encoded}?format=%l:+%c+%t+%h+%w+UV:%u`;
    const quick = (await httpsGet(wttrUrl)).trim();

    // Also get structured data from Open-Meteo
    const geo = await geocode(location);
    if (!geo) {
      return JSON.stringify({ location, summary: quick, source: "wttr.in" });
    }

    const omUrl = `https://api.open-meteo.com/v1/forecast?latitude=${geo.lat}&longitude=${geo.lng}`
      + `&current=temperature_2m,relative_humidity_2m,apparent_temperature,weather_code,`
      + `wind_speed_10m,wind_direction_10m,uv_index,cloud_cover,precipitation`
      + `&temperature_unit=celsius&wind_speed_unit=kmh`;

    const raw = await httpsGet(omUrl);
    const om = JSON.parse(raw);
    const cur = om.current || {};

    return JSON.stringify({
      location: `${geo.name}, ${geo.admin1}, ${geo.country}`,
      coordinates: { lat: geo.lat, lng: geo.lng },
      summary: quick,
      current: {
        temperature: fmtTemp(cur.temperature_2m || 0),
        feels_like: fmtTemp(cur.apparent_temperature || 0),
        humidity: `${cur.relative_humidity_2m || 0}%`,
        wind: `${(cur.wind_speed_10m || 0).toFixed(1)} km/h from ${cur.wind_direction_10m || 0}°`,
        uv_index: cur.uv_index || 0,
        cloud_cover: `${cur.cloud_cover || 0}%`,
        precipitation_mm: cur.precipitation || 0,
        conditions: WMO_CODES[cur.weather_code] || "Unknown",
      },
      source: "wttr.in + Open-Meteo",
    }, null, 2);
  } catch (e) {
    return JSON.stringify({ error: e.message });
  }
}

// ---------------------------------------------------------------------------
// Tool: get_forecast
// ---------------------------------------------------------------------------

async function get_forecast({ location, days }) {
  days = Math.min(Math.max(days || 5, 1), 10);

  try {
    const geo = await geocode(location);
    if (!geo) return JSON.stringify({ error: `Could not find location: ${location}` });

    const url = `https://api.open-meteo.com/v1/forecast?latitude=${geo.lat}&longitude=${geo.lng}`
      + `&daily=weather_code,temperature_2m_max,temperature_2m_min,apparent_temperature_max,`
      + `apparent_temperature_min,precipitation_sum,precipitation_probability_max,`
      + `wind_speed_10m_max,uv_index_max,sunrise,sunset`
      + `&temperature_unit=celsius&wind_speed_unit=kmh&forecast_days=${days}&timezone=auto`;

    const raw = await httpsGet(url);
    const data = JSON.parse(raw);
    const daily = data.daily || {};
    const times = daily.time || [];

    const forecast = times.map((date, i) => ({
      date,
      conditions: WMO_CODES[(daily.weather_code || [])[i]] || "Unknown",
      high: fmtTemp((daily.temperature_2m_max || [])[i] || 0),
      low: fmtTemp((daily.temperature_2m_min || [])[i] || 0),
      feels_like_high: fmtTemp((daily.apparent_temperature_max || [])[i] || 0),
      feels_like_low: fmtTemp((daily.apparent_temperature_min || [])[i] || 0),
      precipitation_mm: (daily.precipitation_sum || [])[i] || 0,
      precipitation_chance: `${(daily.precipitation_probability_max || [])[i] || 0}%`,
      wind_max_kmh: ((daily.wind_speed_10m_max || [])[i] || 0).toFixed(1),
      uv_index_max: (daily.uv_index_max || [])[i] || 0,
      sunrise: ((daily.sunrise || [])[i] || "").split("T")[1] || null,
      sunset: ((daily.sunset || [])[i] || "").split("T")[1] || null,
    }));

    return JSON.stringify({
      location: `${geo.name}, ${geo.admin1}, ${geo.country}`,
      coordinates: { lat: geo.lat, lng: geo.lng },
      forecast,
      source: "Open-Meteo",
    }, null, 2);
  } catch (e) {
    return JSON.stringify({ error: e.message });
  }
}

// ---------------------------------------------------------------------------
// Exports
// ---------------------------------------------------------------------------

module.exports = { get_weather, get_forecast };
