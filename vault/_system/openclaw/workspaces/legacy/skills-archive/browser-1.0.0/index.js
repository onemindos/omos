"use strict";

/**
 * browser — OpenClaw Skill
 *
 * Playwright-based browser automation: navigate pages, take screenshots,
 * extract structured data, fill forms, and execute JavaScript.
 *
 * Registered tools: browser_navigate, browser_screenshot, browser_extract, browser_fill_form, browser_execute_js
 */

let _browser = null;

async function getBrowser() {
  if (_browser && _browser.isConnected()) return _browser;
  const { chromium } = require("playwright");
  _browser = await chromium.launch({ headless: true });
  return _browser;
}

/**
 * TOOL: browser_navigate
 * Navigate to a URL and return the page title + text content.
 *
 * @param {string} url - The URL to navigate to
 * @param {number} [timeout=30000] - Navigation timeout in ms
 * @returns {Object} { title, url, textContent }
 */
async function browser_navigate({ url, timeout = 30000 }) {
  if (!url) throw new Error("url is required");
  const browser = await getBrowser();
  const page = await browser.newPage();
  try {
    await page.goto(url, { waitUntil: "domcontentloaded", timeout });
    const title = await page.title();
    const text = await page.innerText("body").catch(() => "");
    return { title, url: page.url(), textContent: text.slice(0, 10000) };
  } finally {
    await page.close();
  }
}

/**
 * TOOL: browser_screenshot
 * Take a screenshot of a URL. Returns base64-encoded PNG.
 *
 * @param {string} url - The URL to screenshot
 * @param {boolean} [fullPage=false] - Capture the full scrollable page
 * @returns {Object} { url, base64, mimeType }
 */
async function browser_screenshot({ url, fullPage = false }) {
  if (!url) throw new Error("url is required");
  const browser = await getBrowser();
  const page = await browser.newPage();
  try {
    await page.goto(url, { waitUntil: "networkidle", timeout: 30000 });
    const buffer = await page.screenshot({ fullPage });
    return { url: page.url(), base64: buffer.toString("base64"), mimeType: "image/png" };
  } finally {
    await page.close();
  }
}

/**
 * TOOL: browser_extract
 * Extract structured data from a page using CSS selectors.
 *
 * @param {string} url - The URL to extract from
 * @param {string[]} selectors - CSS selectors to extract text from
 * @returns {Object} { url, data }
 */
async function browser_extract({ url, selectors }) {
  if (!url) throw new Error("url is required");
  if (!selectors || !selectors.length) throw new Error("selectors array is required");
  const browser = await getBrowser();
  const page = await browser.newPage();
  try {
    await page.goto(url, { waitUntil: "domcontentloaded", timeout: 30000 });
    const data = {};
    for (const sel of selectors) {
      data[sel] = await page.$$eval(sel, els => els.map(e => e.textContent.trim())).catch(() => []);
    }
    return { url: page.url(), data };
  } finally {
    await page.close();
  }
}

/**
 * TOOL: browser_fill_form
 * Fill form fields on a page and optionally submit.
 *
 * @param {string} url - The URL with the form
 * @param {Object} fields - Map of CSS selector → value to fill
 * @param {string} [submitSelector] - CSS selector of submit button to click
 * @returns {Object} { url, filled, submitted }
 */
async function browser_fill_form({ url, fields, submitSelector }) {
  if (!url) throw new Error("url is required");
  if (!fields || !Object.keys(fields).length) throw new Error("fields object is required");
  const browser = await getBrowser();
  const page = await browser.newPage();
  try {
    await page.goto(url, { waitUntil: "domcontentloaded", timeout: 30000 });
    for (const [selector, value] of Object.entries(fields)) {
      await page.fill(selector, String(value));
    }
    let submitted = false;
    if (submitSelector) {
      await page.click(submitSelector);
      await page.waitForLoadState("domcontentloaded").catch(() => {});
      submitted = true;
    }
    return { url: page.url(), filled: Object.keys(fields), submitted };
  } finally {
    await page.close();
  }
}

/**
 * TOOL: browser_execute_js
 * Execute JavaScript on a page and return the result.
 *
 * @param {string} url - The URL to execute JS on
 * @param {string} script - JavaScript code to evaluate in the page context
 * @returns {Object} { url, result }
 */
async function browser_execute_js({ url, script }) {
  if (!url) throw new Error("url is required");
  if (!script) throw new Error("script is required");
  const browser = await getBrowser();
  const page = await browser.newPage();
  try {
    await page.goto(url, { waitUntil: "domcontentloaded", timeout: 30000 });
    const result = await page.evaluate(script);
    return { url: page.url(), result };
  } finally {
    await page.close();
  }
}

module.exports = { browser_navigate, browser_screenshot, browser_extract, browser_fill_form, browser_execute_js };
