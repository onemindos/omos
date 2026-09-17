import OpenAI from "openai";
import * as fs from "fs";
import * as path from "path";

/**
 * generate-legacy-voice.ts — Generate Legacy AI voice clips from text lines
 *
 * Usage:
 *   OPENAI_API_KEY=... npx tsx src/scripts/generate-legacy-voice.ts <lines-file> [--output-dir=path]
 *
 * Example:
 *   OPENAI_API_KEY=$(op read 'op://aaeyris5z44pdqspk277u2oppi/OpenAI/api-key') \
 *     npx tsx src/scripts/generate-legacy-voice.ts scripts/2026-04-08-topic-legacy-lines.txt
 *
 * Input: text file with one line per Legacy speech segment
 * Output: numbered MP3 files in output/YYYY-MM-DD/ directory
 */

const VOICE = "nova"; // Legacy's voice
const MODEL = "tts-1-hd"; // High quality
const SPEED = 0.92; // Slightly slower for gravitas

async function generateVoice(linesFile: string, outputDir?: string): Promise<void> {
  const apiKey = process.env.OPENAI_API_KEY;
  if (!apiKey) {
    console.error("Error: OPENAI_API_KEY not set");
    process.exit(1);
  }

  if (!fs.existsSync(linesFile)) {
    console.error(`Error: File not found: ${linesFile}`);
    process.exit(1);
  }

  const openai = new OpenAI({ apiKey });

  // Read lines
  const lines = fs
    .readFileSync(linesFile, "utf-8")
    .split("\n")
    .map((l) => l.trim())
    .filter((l) => l.length > 0);

  if (lines.length === 0) {
    console.error("Error: No lines found in file");
    process.exit(1);
  }

  console.log(`Found ${lines.length} Legacy lines to generate`);

  // Create output directory
  const today = new Date().toISOString().split("T")[0];
  const outDir = outputDir || path.join(__dirname, "../../output", today);
  fs.mkdirSync(outDir, { recursive: true });

  let totalCost = 0;

  for (let i = 0; i < lines.length; i++) {
    const line = lines[i];
    const padded = String(i + 1).padStart(2, "0");
    const outPath = path.join(outDir, `legacy-${padded}.mp3`);

    console.log(`\n[${padded}/${String(lines.length).padStart(2, "0")}] Generating: "${line.slice(0, 60)}${line.length > 60 ? "..." : ""}"`);

    const response = await openai.audio.speech.create({
      model: MODEL,
      voice: VOICE,
      input: line,
      speed: SPEED,
      response_format: "mp3",
    });

    const buffer = Buffer.from(await response.arrayBuffer());
    fs.writeFileSync(outPath, buffer);

    // Estimate cost: $0.030 per 1K characters for tts-1-hd
    const charCost = (line.length / 1000) * 0.03;
    totalCost += charCost;

    console.log(`  Saved: ${outPath} (${(buffer.length / 1024).toFixed(0)} KB, ~$${charCost.toFixed(4)})`);
  }

  console.log(`\n--- DONE ---`);
  console.log(`Generated: ${lines.length} voice clips`);
  console.log(`Output dir: ${outDir}`);
  console.log(`Estimated cost: ~$${totalCost.toFixed(4)}`);
  console.log(`\nNext step: Import these into CapCut and place on timeline where Legacy speaks.`);
}

// Parse args
const args = process.argv.slice(2);
const linesFile = args.find((a) => !a.startsWith("--"));
const outputArg = args.find((a) => a.startsWith("--output-dir="));
const outputDir = outputArg?.split("=")[1];

if (!linesFile) {
  console.log("Usage: npx tsx src/scripts/generate-legacy-voice.ts <lines-file> [--output-dir=path]");
  console.log("\nExample:");
  console.log("  npx tsx src/scripts/generate-legacy-voice.ts scripts/2026-04-08-my-topic-legacy-lines.txt");
  console.log("\nThe lines file should have one Legacy speech segment per line.");
  process.exit(0);
}

generateVoice(linesFile, outputDir);
