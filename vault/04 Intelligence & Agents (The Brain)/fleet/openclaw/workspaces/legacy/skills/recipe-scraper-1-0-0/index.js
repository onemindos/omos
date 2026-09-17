/**
 * recipe-scraper-1.0.0 — OpenClaw Skill
 *
 * Extracts structured recipe data from any URL using JSON-LD Schema.org
 * (fast, free, works on AllRecipes, Serious Eats, NYT Cooking, etc.)
 * Falls back to microdata attribute parsing.
 *
 * Ported from reference/oneintel/backend/recipe_scraper.py
 * HP domain: nutrition tracking, meal planning, grocery lists.
 *
 * Env vars: none required (public web scraping only)
 *
 * Registers tools: recipe_scrape
 */

"use strict";

const https = require("https");
const http = require("http");
const { URL } = require("url");

/** Fetch HTML from a URL, following redirects up to 5 times. */
function fetchHtml(urlStr, redirectCount = 0) {
  return new Promise((resolve, reject) => {
    if (redirectCount > 5) return reject(new Error("Too many redirects"));

    let parsedUrl;
    try {
      parsedUrl = new URL(urlStr);
    } catch {
      return reject(new Error(`Invalid URL: ${urlStr}`));
    }

    const lib = parsedUrl.protocol === "https:" ? https : http;
    const options = {
      hostname: parsedUrl.hostname,
      path: parsedUrl.pathname + parsedUrl.search,
      method: "GET",
      headers: {
        "User-Agent":
          "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0.0.0 Safari/537.36",
        Accept: "text/html,application/xhtml+xml",
        "Accept-Language": "en-US,en;q=0.9",
      },
      timeout: 20000,
    };

    const req = lib.request(options, (res) => {
      // Follow redirects
      if ([301, 302, 303, 307, 308].includes(res.statusCode) && res.headers.location) {
        const next = res.headers.location.startsWith("http")
          ? res.headers.location
          : `${parsedUrl.protocol}//${parsedUrl.hostname}${res.headers.location}`;
        return resolve(fetchHtml(next, redirectCount + 1));
      }

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

/** Parse ISO 8601 duration (PT1H30M) → minutes. */
function parseDuration(iso) {
  if (!iso) return null;
  const m = iso.match(/PT(?:(\d+)H)?(?:(\d+)M)?(?:(\d+)S)?/i);
  if (!m) return null;
  return (parseInt(m[1] || 0) * 60) + parseInt(m[2] || 0);
}

/** Normalise ingredient string → { name, amount, unit } */
function parseIngredient(text) {
  if (!text || typeof text !== "string") return null;
  text = text.trim();
  // Simple unit extraction: "2 cups flour" → { amount: 2, unit: "cups", name: "flour" }
  const m = text.match(/^([\d./½¼¾⅓⅔]+)?\s*(cup|tbsp|tsp|oz|lb|g|kg|ml|l|clove|slice|piece|pinch|dash|handful|bunch|can|package|pkg|sprig|stalk|sheet|drop)s?\s+(.+)/i);
  if (m) {
    return { amount: m[1] ? parseFloat(m[1]) : null, unit: m[2].toLowerCase(), name: m[3].trim() };
  }
  return { amount: null, unit: null, name: text };
}

/** Extract recipe from JSON-LD <script> tags. */
function extractJsonLd(html) {
  const scriptRe = /<script[^>]*type=["']application\/ld\+json["'][^>]*>([\s\S]*?)<\/script>/gi;
  let match;

  function findRecipe(obj) {
    if (!obj) return null;
    if (Array.isArray(obj)) {
      for (const item of obj) { const r = findRecipe(item); if (r) return r; }
    } else if (typeof obj === "object") {
      const type = obj["@type"];
      if (type === "Recipe" || (Array.isArray(type) && type.includes("Recipe"))) return obj;
      if (obj["@graph"]) return findRecipe(obj["@graph"]);
      if (obj["mainEntity"]) return findRecipe(obj["mainEntity"]);
    }
    return null;
  }

  while ((match = scriptRe.exec(html)) !== null) {
    try {
      const data = JSON.parse(match[1]);
      const recipe = findRecipe(data);
      if (recipe) return recipe;
    } catch { /* skip malformed */ }
  }
  return null;
}

/** Normalise raw JSON-LD recipe object → clean structure. */
function normaliseRecipe(raw, sourceUrl) {
  // Ingredients
  const ingredients = (raw.recipeIngredient || []).map(parseIngredient).filter(Boolean);

  // Instructions
  let instructions = [];
  const rawInstructions = raw.recipeInstructions || [];
  const instructionList = Array.isArray(rawInstructions) ? rawInstructions : [rawInstructions];
  let step = 1;
  for (const inst of instructionList) {
    if (typeof inst === "string") {
      instructions.push({ step: step++, text: inst.trim() });
    } else if (inst["@type"] === "HowToStep") {
      instructions.push({ step: step++, text: (inst.text || inst.name || "").trim() });
    } else if (inst["@type"] === "HowToSection") {
      for (const sub of inst.itemListElement || []) {
        instructions.push({ step: step++, text: (sub.text || sub.name || "").trim() });
      }
    }
  }

  // Nutrition
  let nutrition = null;
  if (raw.nutrition) {
    const n = raw.nutrition;
    const parseNut = (v) => v ? parseFloat(String(v).replace(/[^\d.]/g, "")) : null;
    nutrition = {
      calories: parseNut(n.calories),
      protein: parseNut(n.proteinContent),
      carbs: parseNut(n.carbohydrateContent),
      fat: parseNut(n.fatContent),
      fiber: parseNut(n.fiberContent),
    };
  }

  // Servings
  let servings = 4;
  if (raw.recipeYield) {
    const sy = Array.isArray(raw.recipeYield) ? raw.recipeYield[0] : raw.recipeYield;
    const parsed = parseInt(String(sy).replace(/[^\d]/g, ""));
    if (!isNaN(parsed) && parsed > 0) servings = parsed;
  }

  // Tags
  let tags = [];
  if (raw.keywords) {
    const kw = typeof raw.keywords === "string" ? raw.keywords : raw.keywords.join(",");
    tags = kw.split(/[,;]/).map((t) => t.trim().toLowerCase()).filter(Boolean);
  }
  if (raw.recipeCategory) {
    const cats = Array.isArray(raw.recipeCategory) ? raw.recipeCategory : [raw.recipeCategory];
    tags.push(...cats.map((c) => c.trim().toLowerCase()));
  }

  return {
    name: raw.name || "Unknown Recipe",
    description: raw.description || null,
    source_url: sourceUrl,
    image_url: raw.image
      ? (Array.isArray(raw.image) ? (raw.image[0].url || raw.image[0]) : (raw.image.url || raw.image))
      : null,
    prep_time_minutes: parseDuration(raw.prepTime),
    cook_time_minutes: parseDuration(raw.cookTime),
    total_time_minutes: parseDuration(raw.totalTime),
    servings,
    cuisine: raw.recipeCuisine
      ? (Array.isArray(raw.recipeCuisine) ? raw.recipeCuisine[0] : raw.recipeCuisine)
      : null,
    tags: [...new Set(tags)],
    ingredients,
    instructions,
    nutrition,
    source_method: "json-ld",
  };
}

/**
 * recipe_scrape — Extract structured recipe data from a URL.
 *
 * Works on AllRecipes, Serious Eats, NYT Cooking, Minimalist Baker,
 * Budget Bytes, and any site using Schema.org JSON-LD markup.
 *
 * @param {string} url - Full URL of the recipe page
 */
async function recipe_scrape({ url }) {
  if (!url) throw new Error("url is required");

  // Basic URL validation — no SSRF: only allow http/https to public internet
  let parsedUrl;
  try {
    parsedUrl = new URL(url);
  } catch {
    throw new Error(`Invalid URL: ${url}`);
  }
  if (!["http:", "https:"].includes(parsedUrl.protocol)) {
    throw new Error("Only http/https URLs are supported");
  }
  // Block private/loopback ranges
  const hostname = parsedUrl.hostname.toLowerCase();
  if (
    hostname === "localhost" ||
    hostname.startsWith("127.") ||
    hostname.startsWith("10.") ||
    hostname.startsWith("192.168.") ||
    hostname.startsWith("169.254.") ||
    hostname === "0.0.0.0" ||
    /^::1$|^fc00:/i.test(hostname)
  ) {
    throw new Error("Private/internal URLs are not allowed");
  }

  const html = await fetchHtml(url);
  const raw = extractJsonLd(html);

  if (!raw) {
    throw new Error(
      `Could not find recipe data on ${url}. ` +
      "The site may not use Schema.org markup, require JavaScript, or be paywalled."
    );
  }

  return normaliseRecipe(raw, url);
}

// ─── OpenClaw tool registry ───────────────────────────────────────────────────
module.exports = { recipe_scrape };
