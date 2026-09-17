import OpenAI from "openai";
import * as fs from "fs";
import * as path from "path";

/**
 * generate-script.ts — AI-powered daily video script generator
 *
 * Usage:
 *   OPENAI_API_KEY=... npx tsx src/scripts/generate-script.ts "topic" [--pillar so|hp|le|ge]
 *
 * Example:
 *   OPENAI_API_KEY=$(op read 'op://aaeyris5z44pdqspk277u2oppi/OpenAI/api-key') \
 *     npx tsx src/scripts/generate-script.ts "How I use AI to manage my entire life"
 */

const SYSTEM_PROMPT = `You are a YouTube script writer for Zeus Delacruz's channel "One Mind."

Zeus is building One Mind — an AI-native human enhancement framework. His tagline: "You + AI = One Mind."
His AI assistant is named Legacy (voice: female, confident, data-driven).

The channel format:
- Zeus talks to the camera explaining concepts in a casual, confident, direct style
- Legacy (AI voice) drops in to provide data, breakdowns, counterpoints, or recaps
- Visual cues mark where screen recordings, demos, or B-roll should appear

Write scripts in this EXACT format:

---
title: "[VIDEO TITLE — attention-grabbing, YouTube-optimized]"
hook_text: "[One-line hook for thumbnails and first 3 seconds]"
pillar: [so|hp|le|ge|method]
estimated_length: [e.g. "7 minutes"]
---

## HOOK (0:00-0:30)
[ZEUS — FACE CAM]
[Bold opening statement or question]

## SECTION 1: [Section Name] (0:30-3:00)
[ZEUS — FACE CAM]
[Main explanation]

[VISUAL: description of what to show on screen]

## LEGACY DROP-IN (3:00-4:00)
[LEGACY — AI VOICE]
"Legacy here. [Legacy's commentary]"

[VISUAL: supporting data/graphics while Legacy talks]

## SECTION 2: [Section Name] (4:00-7:00)
[ZEUS — FACE CAM]
[Deeper explanation, examples]

[VISUAL: demo/screen recording description]

## LEGACY RECAP (7:00-7:45)
[LEGACY — AI VOICE]
"Three things to remember: [key takeaways]"

## CTA (7:45-8:15)
[ZEUS — FACE CAM]
[Call to action — subscribe, link, "Legacy and I will see you tomorrow"]

---

Rules:
- Zeus's tone: casual, direct, no fluff, slightly intense, like talking to a smart friend
- Legacy's tone: calm, data-driven, slightly formal but warm, concise
- Mark EVERY visual cue with [VISUAL: ...]
- Mark speaker switches clearly with [ZEUS — FACE CAM] or [LEGACY — AI VOICE]
- Write natural spoken language, NOT essay language
- Include specific numbers, examples, tools when possible
- The hook must be scroll-stopping — bold claim, surprising stat, or contrarian take
- Keep Legacy's sections SHORT (30-60 seconds each)
- Total script should be 5-10 minutes when spoken`;

async function generateScript(topic: string, pillar?: string): Promise<void> {
  const apiKey = process.env.OPENAI_API_KEY;
  if (!apiKey) {
    console.error("Error: OPENAI_API_KEY not set");
    process.exit(1);
  }

  const openai = new OpenAI({ apiKey });

  const userPrompt = `Write a YouTube video script about: "${topic}"${pillar ? `\nContent pillar: ${pillar}` : ""}

Make it engaging, specific, and actionable. The hook should make someone stop scrolling.`;

  console.log(`Generating script for: "${topic}"...`);

  const response = await openai.chat.completions.create({
    model: "gpt-4o",
    messages: [
      { role: "system", content: SYSTEM_PROMPT },
      { role: "user", content: userPrompt },
    ],
    temperature: 0.8,
    max_tokens: 3000,
  });

  const script = response.choices[0]?.message?.content;
  if (!script) {
    console.error("Error: No script generated");
    process.exit(1);
  }

  // Create output directory
  const today = new Date().toISOString().split("T")[0];
  const slug = topic
    .toLowerCase()
    .replace(/[^a-z0-9]+/g, "-")
    .replace(/^-|-$/g, "")
    .slice(0, 60);

  const scriptsDir = path.join(__dirname, "../../scripts");
  fs.mkdirSync(scriptsDir, { recursive: true });

  // Save full script
  const scriptPath = path.join(scriptsDir, `${today}-${slug}.md`);
  fs.writeFileSync(scriptPath, script, "utf-8");
  console.log(`Script saved: ${scriptPath}`);

  // Extract Legacy lines
  const legacyLines: string[] = [];
  const lines = script.split("\n");
  let inLegacy = false;

  for (const line of lines) {
    if (line.includes("[LEGACY")) {
      inLegacy = true;
      continue;
    }
    if (line.startsWith("## ") && !line.includes("LEGACY")) {
      inLegacy = false;
    }
    if (inLegacy) {
      // Extract quoted speech or content lines
      const quoted = line.match(/"([^"]+)"/);
      if (quoted) {
        legacyLines.push(quoted[1]);
      } else if (
        line.trim() &&
        !line.startsWith("[") &&
        !line.startsWith("#")
      ) {
        legacyLines.push(line.trim());
      }
    }
  }

  if (legacyLines.length > 0) {
    const legacyPath = path.join(
      scriptsDir,
      `${today}-${slug}-legacy-lines.txt`
    );
    fs.writeFileSync(legacyPath, legacyLines.join("\n"), "utf-8");
    console.log(`Legacy lines saved: ${legacyPath} (${legacyLines.length} lines)`);
  }

  // Print summary
  console.log("\n--- SCRIPT PREVIEW ---");
  console.log(script.slice(0, 500) + "...\n");
  console.log(`Total length: ~${Math.round(script.length / 150)} minutes when spoken`);
  console.log(`Legacy segments: ${legacyLines.length} lines`);
}

// Parse args
const args = process.argv.slice(2);
const topic = args.find((a) => !a.startsWith("--"));
const pillarArg = args.find((a) => a.startsWith("--pillar="));
const pillar = pillarArg?.split("=")[1];

if (!topic) {
  console.log("Usage: npx tsx src/scripts/generate-script.ts \"topic\" [--pillar=so|hp|le|ge|method]");
  console.log("\nExamples:");
  console.log('  npx tsx src/scripts/generate-script.ts "How I use AI to manage my entire life"');
  console.log('  npx tsx src/scripts/generate-script.ts "My morning routine" --pillar=hp');
  console.log('  npx tsx src/scripts/generate-script.ts "Building a second brain with One Mind" --pillar=method');
  process.exit(0);
}

generateScript(topic, pillar);
