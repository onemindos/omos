"use strict";

/**
 * simple-weather — OpenClaw Skill (Complete Working Example)
 *
 * Fetches weather data from Open-Meteo (free, no API key required).
 * Demonstrates: geocoding, chained API calls, unit conversion, error handling.
 *
 * Registered tools: weather_current, weather_forecast
 */

// ─── WMO Weather Condition Codes ─────────────────────────────────────────────
const WMO_CODES = {
  0: "Clear sky",
  1: "Mainly clear",
  2: "Partly cloudy",
  3: "Overcast",
  45: "Foggy",
  48: "Icy fog",
  51: "Light drizzle",
  53: "Moderate drizzle",
  55: "Dense drizzle",
  61: "Slight rain",
  63: "Moderate rain",
  65: "Heavy rain",
  71: "Slight snow",
  73: "Moderate snow",
  75: "Heavy snow",
  77: "Snow grains",
  80: "Slight rain showers",
  81: "Moderate rain showers",
  82: "Violent rain showers",
  85: "Slight snow showers",
  86: "Heavy snow showers",
  95: "Thunderstorm",
  96: "Thunderstorm with slight hail",
  99: "Thunderstorm with heavy hail",
};

/** Resolve a location name to coordinates using Open-Meteo geocoding */
async function geocode(location) {
  const url = `https://geocoding-api.open-meteo.com/v1/search?name=${encodeURIComponent(location)}&count=1&language=en&format=json`;
  const res = await fetch(url, { signal: AbortSignal.timeout(10000) });
  if (!res.ok) throw new Error(`Geocoding request failed: HTTP ${res.status}`);

  const data = await res.json();
  if (!data.results || data.results.length === 0) {
    throw new Error(`Location not found: "${location}". Try a more specific name or provide latitude/longitude.`);
  }

  const place = data.results[0];
  return {
    latitude: place.latitude,
    longitude: place.longitude,
    name: place.name,
    country: place.country,
    display: `${place.name}, ${place.country}`,
  };
}

/** Convert Celsius to Fahrenheit */
function toF(c) {
  return Math.round((c * 9) / 5 + 32);
}

/** Convert WMO code to human-readable string */
function wmoCondition(code) {
  return WMO_CODES[code] || `Condition code ${code}`;
}

// ─── Tools ───────────────────────────────────────────────────────────────────

/**
 * TOOL: weather_current
 * Get current weather conditions at a location.
 *
 * @param {string} [location] - City/location name (e.g., "Lagos, Nigeria")
 * @param {number} [latitude] - Latitude (if not using location name)
 * @param {number} [longitude] - Longitude (if not using location name)
 * @returns {Object} Current weather conditions
 */
async function weather_current({ location, latitude, longitude }) {
  let coords;
  let displayName = location;

  if (location && !latitude) {
    coords = await geocode(location);
    displayName = coords.display;
    latitude = coords.latitude;
    longitude = coords.longitude;
  } else if (!latitude || !longitude) {
    throw new Error("Provide either a location name or both latitude and longitude");
  }

  const url =
    `https://api.open-meteo.com/v1/forecast` +
    `?latitude=${latitude}&longitude=${longitude}` +
    `&current=temperature_2m,apparent_temperature,relative_humidity_2m,` +
    `precipitation,weather_code,wind_speed_10m,wind_direction_10m,uv_index` +
    `&wind_speed_unit=kmh&temperature_unit=celsius&timezone=auto`;

  const res = await fetch(url, { signal: AbortSignal.timeout(15000) });
  if (!res.ok) throw new Error(`Weather API error: HTTP ${res.status}`);
  const data = await res.json();

  const c = data.current;
  const tempC = Math.round(c.temperature_2m);
  const feelsC = Math.round(c.apparent_temperature);

  // Convert wind direction degrees to compass
  const dirs = ["N", "NE", "E", "SE", "S", "SW", "W", "NW"];
  const windDir = dirs[Math.round(c.wind_direction_10m / 45) % 8];

  return {
    location: displayName || `${latitude}, ${longitude}`,
    temperature_c: tempC,
    temperature_f: toF(tempC),
    feels_like_c: feelsC,
    feels_like_f: toF(feelsC),
    humidity_pct: c.relative_humidity_2m,
    wind_speed_kmh: Math.round(c.wind_speed_10m),
    wind_direction: windDir,
    condition: wmoCondition(c.weather_code),
    uv_index: c.uv_index,
    precipitation_mm: c.precipitation,
    observed_at: c.time,
  };
}

/**
 * TOOL: weather_forecast
 * Get a multi-day weather forecast for a location.
 *
 * @param {string} [location] - City/location name
 * @param {number} [latitude] - Latitude (if not using location name)
 * @param {number} [longitude] - Longitude (if not using location name)
 * @param {number} [days=7] - Number of forecast days (1–16)
 * @returns {Object} Forecast with daily conditions
 */
async function weather_forecast({ location, latitude, longitude, days = 7 }) {
  if (days < 1 || days > 16) throw new Error("days must be between 1 and 16");

  let displayName = location;

  if (location && !latitude) {
    const coords = await geocode(location);
    displayName = coords.display;
    latitude = coords.latitude;
    longitude = coords.longitude;
  } else if (!latitude || !longitude) {
    throw new Error("Provide either a location name or both latitude and longitude");
  }

  const url =
    `https://api.open-meteo.com/v1/forecast` +
    `?latitude=${latitude}&longitude=${longitude}` +
    `&daily=weather_code,temperature_2m_max,temperature_2m_min,` +
    `precipitation_sum,wind_speed_10m_max,uv_index_max` +
    `&forecast_days=${days}&temperature_unit=celsius&wind_speed_unit=kmh&timezone=auto`;

  const res = await fetch(url, { signal: AbortSignal.timeout(15000) });
  if (!res.ok) throw new Error(`Forecast API error: HTTP ${res.status}`);
  const data = await res.json();

  const d = data.daily;
  const forecast = d.time.map((date, i) => ({
    date,
    condition: wmoCondition(d.weather_code[i]),
    high_c: Math.round(d.temperature_2m_max[i]),
    high_f: toF(Math.round(d.temperature_2m_max[i])),
    low_c: Math.round(d.temperature_2m_min[i]),
    low_f: toF(Math.round(d.temperature_2m_min[i])),
    precipitation_mm: Math.round(d.precipitation_sum[i] * 10) / 10,
    wind_max_kmh: Math.round(d.wind_speed_10m_max[i]),
    uv_index_max: d.uv_index_max[i],
  }));

  return {
    location: displayName || `${latitude}, ${longitude}`,
    days: forecast.length,
    forecast,
  };
}

// ─── Exports ─────────────────────────────────────────────────────────────────
module.exports = { weather_current, weather_forecast };
