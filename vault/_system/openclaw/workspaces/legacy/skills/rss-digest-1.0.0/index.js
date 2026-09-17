/**
 * rss-digest-1.0.0 — OpenClaw Skill
 *
 * Fetch and digest RSS feeds from Substack newsletters and other sources.
 * Manages a persistent feed list, fetches recent articles, and returns
 * structured digests for Legacy to summarize.
 *
 * Tools: fetch_rss_digest, list_rss_feeds, add_rss_feed, remove_rss_feed
 */

"use strict";

const fs = require("fs");
const path = require("path");
const RssParser = require("rss-parser");

// ---------------------------------------------------------------------------
// Feed store — persisted as JSON in skill directory
// ---------------------------------------------------------------------------

const FEEDS_FILE = path.join(__dirname, "feeds.json");

const DEFAULT_FEEDS = [
  // Homestead & Gardening
  { name: "houseofgreen", url: "https://houseofgreen.substack.com/feed", category: "homestead" },
  { name: "holisticfarming", url: "https://holisticfarming.substack.com/feed", category: "homestead" },
  { name: "charlesdowding", url: "https://charlesdowding549359.substack.com/feed", category: "homestead" },
  { name: "acornlandlabs", url: "https://acornlandlabs.substack.com/feed", category: "homestead" },
  { name: "peasantpolis", url: "https://peasantpolis.substack.com/feed", category: "homestead" },
  // Parenting
  { name: "dralizapressman", url: "https://dralizapressman.substack.com/feed", category: "parenting" },
  { name: "theparentingguide", url: "https://theparentingguide.substack.com/feed", category: "parenting" },
  { name: "screenlessdad", url: "https://screenlessdad.substack.com/feed", category: "parenting" },
  // Health
  { name: "dremilyleeming", url: "https://dremilyleeming.substack.com/feed", category: "health" },
  { name: "vajenda", url: "https://vajenda.substack.com/feed", category: "health" },
  { name: "joshturknett", url: "https://joshturknett.substack.com/feed", category: "health" },
  // Philosophy & Mindset
  { name: "stoicwisdoms", url: "https://stoicwisdoms.substack.com/feed", category: "philosophy" },
  { name: "johannadebiase", url: "https://johannadebiase.substack.com/feed", category: "philosophy" },
  { name: "yungpueblo", url: "https://yungpueblo.substack.com/feed", category: "philosophy" },
];

function loadFeeds() {
  try {
    if (fs.existsSync(FEEDS_FILE)) {
      return JSON.parse(fs.readFileSync(FEEDS_FILE, "utf8"));
    }
  } catch {
    // corrupted file — reset to defaults
  }
  saveFeeds(DEFAULT_FEEDS);
  return [...DEFAULT_FEEDS];
}

function saveFeeds(feeds) {
  fs.writeFileSync(FEEDS_FILE, JSON.stringify(feeds, null, 2), "utf8");
}

// ---------------------------------------------------------------------------
// HTML → plain text (lightweight, no dependency)
// ---------------------------------------------------------------------------

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

function truncate(text, maxLen = 500) {
  if (!text || text.length <= maxLen) return text || "";
  return text.slice(0, maxLen).replace(/\s+\S*$/, "") + "…";
}

// ---------------------------------------------------------------------------
// RSS fetcher
// ---------------------------------------------------------------------------

const parser = new RssParser({
  timeout: 15000,
  headers: { "User-Agent": "OpenClaw-RSS-Digest/1.0" },
});

async function fetchFeed(feed, hoursBack, maxItems) {
  const cutoff = new Date(Date.now() - hoursBack * 60 * 60 * 1000);

  try {
    const parsed = await parser.parseURL(feed.url);
    const items = (parsed.items || [])
      .map((item) => {
        const pubDate = item.pubDate ? new Date(item.pubDate) : null;
        const snippet = item.contentSnippet || item.content || item.summary || "";
        return {
          title: item.title || "(untitled)",
          link: item.link || "",
          date: pubDate ? pubDate.toISOString() : null,
          summary: truncate(stripHtml(snippet), 500),
          _ts: pubDate ? pubDate.getTime() : 0,
        };
      })
      .filter((a) => !a.date || new Date(a.date) >= cutoff)
      .sort((a, b) => b._ts - a._ts)
      .slice(0, maxItems)
      .map(({ _ts, ...rest }) => rest);

    return {
      feed_name: feed.name,
      feed_title: parsed.title || feed.name,
      category: feed.category || "uncategorized",
      article_count: items.length,
      articles: items,
    };
  } catch (e) {
    return {
      feed_name: feed.name,
      category: feed.category || "uncategorized",
      article_count: 0,
      articles: [],
      error: e.message,
    };
  }
}

// ---------------------------------------------------------------------------
// Tool: fetch_rss_digest
// ---------------------------------------------------------------------------

async function fetch_rss_digest({ feed_name, hours, max_per_feed } = {}) {
  const hoursBack = Math.max(hours || 24, 1);
  const maxItems = Math.min(Math.max(max_per_feed || 5, 1), 20);
  const feeds = loadFeeds();

  let targets = feeds;
  if (feed_name) {
    const match = feeds.find(
      (f) => f.name.toLowerCase() === feed_name.toLowerCase()
    );
    if (!match) {
      return JSON.stringify({
        error: `Feed '${feed_name}' not found. Use list_rss_feeds to see available feeds.`,
      });
    }
    targets = [match];
  }

  const results = await Promise.all(
    targets.map((f) => fetchFeed(f, hoursBack, maxItems))
  );

  const totalArticles = results.reduce((s, r) => s + r.article_count, 0);
  const errors = results.filter((r) => r.error);

  // Group by category
  const byCategory = {};
  for (const r of results) {
    if (r.article_count === 0 && !r.error) continue;
    const cat = r.category || "uncategorized";
    if (!byCategory[cat]) byCategory[cat] = [];
    byCategory[cat].push(r);
  }

  return JSON.stringify(
    {
      digest_period: `Last ${hoursBack} hours`,
      total_articles: totalArticles,
      feeds_checked: targets.length,
      feeds_with_errors: errors.length,
      categories: byCategory,
    },
    null,
    2
  );
}

// ---------------------------------------------------------------------------
// Tool: list_rss_feeds
// ---------------------------------------------------------------------------

async function list_rss_feeds() {
  const feeds = loadFeeds();
  const byCategory = {};
  for (const f of feeds) {
    const cat = f.category || "uncategorized";
    if (!byCategory[cat]) byCategory[cat] = [];
    byCategory[cat].push({ name: f.name, url: f.url });
  }
  return JSON.stringify(
    { total: feeds.length, categories: byCategory },
    null,
    2
  );
}

// ---------------------------------------------------------------------------
// Tool: add_rss_feed
// ---------------------------------------------------------------------------

async function add_rss_feed({ name, url, category }) {
  if (!name || typeof name !== "string") {
    return JSON.stringify({ error: "name is required" });
  }

  const feedName = name.toLowerCase().replace(/[^a-z0-9_-]/g, "");
  if (!feedName) {
    return JSON.stringify({ error: "Invalid feed name — use alphanumeric characters" });
  }

  const feeds = loadFeeds();
  if (feeds.find((f) => f.name === feedName)) {
    return JSON.stringify({ error: `Feed '${feedName}' already exists` });
  }

  // Auto-generate Substack URL if no URL provided
  const feedUrl = url || `https://${feedName}.substack.com/feed`;

  // Validate the feed is reachable
  try {
    await parser.parseURL(feedUrl);
  } catch (e) {
    return JSON.stringify({
      error: `Could not fetch feed at ${feedUrl}: ${e.message}`,
    });
  }

  const entry = {
    name: feedName,
    url: feedUrl,
    category: category || "uncategorized",
  };
  feeds.push(entry);
  saveFeeds(feeds);

  return JSON.stringify({
    success: true,
    message: `Added feed '${feedName}' (${feedUrl})`,
    feed: entry,
    total_feeds: feeds.length,
  });
}

// ---------------------------------------------------------------------------
// Tool: remove_rss_feed
// ---------------------------------------------------------------------------

async function remove_rss_feed({ name }) {
  if (!name) {
    return JSON.stringify({ error: "name is required" });
  }

  const feeds = loadFeeds();
  const idx = feeds.findIndex(
    (f) => f.name.toLowerCase() === name.toLowerCase()
  );
  if (idx === -1) {
    return JSON.stringify({
      error: `Feed '${name}' not found. Use list_rss_feeds to see available feeds.`,
    });
  }

  const removed = feeds.splice(idx, 1)[0];
  saveFeeds(feeds);

  return JSON.stringify({
    success: true,
    message: `Removed feed '${removed.name}'`,
    remaining_feeds: feeds.length,
  });
}

// ---------------------------------------------------------------------------
// Exports
// ---------------------------------------------------------------------------

module.exports = {
  fetch_rss_digest,
  list_rss_feeds,
  add_rss_feed,
  remove_rss_feed,
};
