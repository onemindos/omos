/**
 * Generate Legacy voice clips for the One Mind Intro episode.
 *
 * Reads the script file and generates one MP3 per segment using OpenAI TTS-HD.
 * Output: output/onemind-intro/segment-01.mp3, segment-02.mp3, etc.
 *
 * Usage: npx tsx src/scripts/generate-intro-voice.ts
 */

import fs from "fs";
import path from "path";
import OpenAI from "openai";
import { execSync } from "child_process";

const VOICE = "nova" as const;
const SPEED = 0.92;
const MODEL = "tts-1-hd";

async function main() {
  let apiKey = process.env.OPENAI_API_KEY;
  if (!apiKey) {
    apiKey = execSync(
      "op read 'op://aaeyris5z44pdqspk277u2oppi/OpenAI/api-key'"
    )
      .toString()
      .trim();
  }

  const openai = new OpenAI({ apiKey });

  const scriptPath = path.join(
    __dirname,
    "../../scripts/onemind-intro-legacy-lines.md"
  );
  const content = fs.readFileSync(scriptPath, "utf-8");

  // Parse segments
  const segments: { name: string; text: string }[] = [];
  const lines = content.split("\n");
  let currentName = "";
  let currentText = "";

  for (const line of lines) {
    if (line.startsWith("## Segment")) {
      if (currentName && currentText.trim()) {
        segments.push({ name: currentName, text: currentText.trim() });
      }
      currentName = line.replace("## ", "");
      currentText = "";
    } else if (currentName && !line.startsWith(">") && line.trim()) {
      // Strip leading/trailing quotes from spoken text
      const cleaned = line.trim().replace(/^[""]|[""]$/g, "");
      currentText += cleaned + " ";
    }
  }
  if (currentName && currentText.trim()) {
    segments.push({ name: currentName, text: currentText.trim() });
  }

  console.log(`Found ${segments.length} segments to generate.\n`);

  const outDir = path.join(__dirname, "../../output/onemind-intro");
  fs.mkdirSync(outDir, { recursive: true });

  const publicDir = path.join(__dirname, "../../public/audio/onemind-intro");
  fs.mkdirSync(publicDir, { recursive: true });

  let totalCost = 0;

  for (let i = 0; i < segments.length; i++) {
    const seg = segments[i];
    const filename = `segment-${String(i + 1).padStart(2, "0")}.mp3`;
    const outPath = path.join(outDir, filename);
    const pubPath = path.join(publicDir, filename);

    console.log(`[${i + 1}/${segments.length}] ${seg.name}`);
    console.log(`  Text: "${seg.text.slice(0, 80)}..."`);

    const response = await openai.audio.speech.create({
      model: MODEL,
      voice: VOICE,
      speed: SPEED,
      input: seg.text,
    });

    const buffer = Buffer.from(await response.arrayBuffer());
    fs.writeFileSync(outPath, buffer);
    fs.copyFileSync(outPath, pubPath);

    const charCount = seg.text.length;
    const cost = (charCount / 1_000_000) * 30;
    totalCost += cost;

    console.log(
      `  → ${filename} (${(buffer.length / 1024).toFixed(0)} KB, ${charCount} chars, ~$${cost.toFixed(4)})\n`
    );
  }

  console.log(`\nDone! ${segments.length} clips generated.`);
  console.log(`Output: ${outDir}`);
  console.log(`Public: ${publicDir}`);
  console.log(`Estimated cost: ~$${totalCost.toFixed(4)}`);
}

main().catch(console.error);
