/**
 * Generate Two-Voice clips for One Mind Intro V2 (Zeus + Legacy dialogue).
 *
 * - Parses the V2 script for ZEUS: and LEGACY: lines
 * - Generates individual TTS clips per line (onyx for Zeus, nova for Legacy)
 * - Concatenates per segment using ffmpeg with natural pauses between speakers
 * - Outputs: public/audio/onemind-intro-v2/segment-XX.mp3
 * - Also outputs: public/audio/onemind-intro-v2/timing.json (durations for composition)
 *
 * Usage: OPENAI_API_KEY=... npx tsx src/scripts/generate-intro-v2-voice.ts
 *    Or: npx tsx src/scripts/generate-intro-v2-voice.ts  (uses 1Password)
 */

import fs from "fs";
import path from "path";
import OpenAI from "openai";
import { execSync } from "child_process";

const ZEUS_VOICE = "onyx" as const;   // Deep, warm, energetic
const LEGACY_VOICE = "nova" as const; // Warm, expressive
const ZEUS_SPEED = 1.0;
const LEGACY_SPEED = 0.92;
const MODEL = "tts-1-hd";
const PAUSE_BETWEEN_SPEAKERS_MS = 350; // natural pause between speakers
const PAUSE_SAME_SPEAKER_MS = 150;     // short pause when same speaker continues

interface Line {
  speaker: "zeus" | "legacy";
  text: string;
}

interface Segment {
  number: number;
  name: string;
  lines: Line[];
}

function parseScript(content: string): Segment[] {
  const segments: Segment[] = [];
  const rawLines = content.split("\n");
  let currentSeg: Segment | null = null;

  for (const rawLine of rawLines) {
    const segMatch = rawLine.match(/^## Segment (\d+)\s*[—–-]\s*(.+)/);
    if (segMatch) {
      if (currentSeg && currentSeg.lines.length > 0) {
        segments.push(currentSeg);
      }
      currentSeg = {
        number: parseInt(segMatch[1], 10),
        name: segMatch[2].trim(),
        lines: [],
      };
      continue;
    }

    if (!currentSeg) continue;

    // Match ZEUS: "text" or LEGACY: "text"
    const lineMatch = rawLine.match(/^(ZEUS|LEGACY):\s*"(.+)"$/);
    if (lineMatch) {
      const speaker = lineMatch[1].toLowerCase() as "zeus" | "legacy";
      // Strip smart quotes and clean up
      const text = lineMatch[2]
        .replace(/[""]/g, '"')
        .replace(/['']/g, "'")
        .replace(/\s+/g, " ")
        .trim();
      currentSeg.lines.push({ speaker, text });
    }
  }
  if (currentSeg && currentSeg.lines.length > 0) {
    segments.push(currentSeg);
  }

  return segments;
}

async function generateClip(
  openai: OpenAI,
  text: string,
  speaker: "zeus" | "legacy",
  outPath: string
): Promise<number> {
  const voice = speaker === "zeus" ? ZEUS_VOICE : LEGACY_VOICE;
  const speed = speaker === "zeus" ? ZEUS_SPEED : LEGACY_SPEED;

  const response = await openai.audio.speech.create({
    model: MODEL,
    voice,
    speed,
    input: text,
  });

  const buffer = Buffer.from(await response.arrayBuffer());
  fs.writeFileSync(outPath, buffer);

  // Get duration via ffprobe
  const duration = parseFloat(
    execSync(
      `ffprobe -v error -show_entries format=duration -of csv=p=0 "${outPath}"`
    )
      .toString()
      .trim()
  );

  return duration;
}

function generateSilence(durationMs: number, outPath: string): void {
  const durationSec = durationMs / 1000;
  // Match TTS output format: 24kHz mono mp3
  execSync(
    `ffmpeg -y -f lavfi -i anullsrc=r=24000:cl=mono -t ${durationSec} -c:a libmp3lame -b:a 48k "${outPath}" 2>&1`
  );
}

function concatenateClips(clipPaths: string[], outPath: string): void {
  // Use ffmpeg filter_complex to concatenate (more robust than concat demuxer)
  const inputs = clipPaths.map((p) => `-i "${p}"`).join(" ");
  const filterParts = clipPaths.map((_, i) => `[${i}:a]`).join("");
  const filter = `${filterParts}concat=n=${clipPaths.length}:v=0:a=1[out]`;

  execSync(
    `ffmpeg -y ${inputs} -filter_complex "${filter}" -map "[out]" -c:a libmp3lame -q:a 2 "${outPath}" 2>&1`
  );
}

async function main() {
  let apiKey = process.env.OPENAI_API_KEY;
  if (!apiKey) {
    console.log("No OPENAI_API_KEY in env, trying 1Password...");
    apiKey = execSync(
      "op read 'op://aaeyris5z44pdqspk277u2oppi/OpenAI/api-key'"
    )
      .toString()
      .trim();
  }

  const openai = new OpenAI({ apiKey });

  // Parse script
  const scriptPath = path.join(
    __dirname,
    "../../scripts/onemind-intro-v2-script.md"
  );
  const content = fs.readFileSync(scriptPath, "utf-8");
  const segments = parseScript(content);

  console.log(`Found ${segments.length} segments.\n`);

  // Setup dirs
  const tmpDir = path.join(__dirname, "../../output/onemind-intro-v2/tmp");
  const outDir = path.join(__dirname, "../../output/onemind-intro-v2");
  const publicDir = path.join(
    __dirname,
    "../../public/audio/onemind-intro-v2"
  );
  fs.mkdirSync(tmpDir, { recursive: true });
  fs.mkdirSync(publicDir, { recursive: true });

  // Generate silence clips for pauses
  const pauseBetween = path.join(tmpDir, "pause-between.mp3");
  const pauseSame = path.join(tmpDir, "pause-same.mp3");
  generateSilence(PAUSE_BETWEEN_SPEAKERS_MS, pauseBetween);
  generateSilence(PAUSE_SAME_SPEAKER_MS, pauseSame);

  const timing: Record<
    string,
    { duration: number; frames: number; lines: { speaker: string; duration: number }[] }
  > = {};
  let totalCost = 0;
  let totalLines = 0;

  for (const seg of segments) {
    const segKey = `s${String(seg.number).padStart(2, "0")}`;
    const segFilename = `segment-${String(seg.number).padStart(2, "0")}.mp3`;
    console.log(`\n━━━ Segment ${seg.number}: ${seg.name} (${seg.lines.length} lines) ━━━`);

    const clipPaths: string[] = [];
    const lineDurations: { speaker: string; duration: number }[] = [];
    let prevSpeaker: string | null = null;

    for (let li = 0; li < seg.lines.length; li++) {
      const line = seg.lines[li];
      const clipFile = `seg${seg.number}-line${li}-${line.speaker}.mp3`;
      const clipPath = path.join(tmpDir, clipFile);

      // Add pause between lines
      if (li > 0) {
        const pauseFile =
          line.speaker !== prevSpeaker ? pauseBetween : pauseSame;
        clipPaths.push(pauseFile);
      }

      console.log(
        `  [${line.speaker.toUpperCase()}] "${line.text.slice(0, 60)}${line.text.length > 60 ? "..." : ""}"`
      );

      const duration = await generateClip(openai, line.text, line.speaker, clipPath);
      clipPaths.push(clipPath);
      lineDurations.push({ speaker: line.speaker, duration });

      const charCount = line.text.length;
      const cost = (charCount / 1_000_000) * 30;
      totalCost += cost;
      totalLines++;

      console.log(`    → ${(duration).toFixed(2)}s, ${charCount} chars`);

      prevSpeaker = line.speaker;
    }

    // Concatenate all clips for this segment
    const segOutPath = path.join(outDir, segFilename);
    const segPubPath = path.join(publicDir, segFilename);
    concatenateClips(clipPaths, segOutPath);
    fs.copyFileSync(segOutPath, segPubPath);

    // Get final segment duration
    const segDuration = parseFloat(
      execSync(
        `ffprobe -v error -show_entries format=duration -of csv=p=0 "${segOutPath}"`
      )
        .toString()
        .trim()
    );

    const segFrames = Math.ceil(segDuration * 30);
    timing[segKey] = {
      duration: segDuration,
      frames: segFrames,
      lines: lineDurations,
    };

    console.log(`  ✓ ${segFilename}: ${segDuration.toFixed(2)}s (${segFrames} frames)`);
  }

  // Write timing data
  const timingPath = path.join(publicDir, "timing.json");
  fs.writeFileSync(timingPath, JSON.stringify(timing, null, 2));
  fs.copyFileSync(timingPath, path.join(outDir, "timing.json"));

  // Cleanup tmp
  console.log("\nCleaning up temp files...");
  for (const f of fs.readdirSync(tmpDir)) {
    fs.unlinkSync(path.join(tmpDir, f));
  }
  fs.rmdirSync(tmpDir);

  console.log(`\n━━━ DONE ━━━`);
  console.log(`Segments: ${segments.length}`);
  console.log(`Total lines: ${totalLines}`);
  console.log(`Output: ${outDir}`);
  console.log(`Public: ${publicDir}`);
  console.log(`Timing: ${timingPath}`);
  console.log(`Estimated cost: ~$${totalCost.toFixed(4)}`);

  // Print segment summary for composition
  console.log(`\n// Paste into OneMindIntroV2.tsx:`);
  console.log(`const SEG = {`);
  for (const [key, val] of Object.entries(timing)) {
    const padding = 30; // extra frames for visual breathing room
    console.log(
      `  ${key}: { audio: ${val.frames}, total: ${val.frames + padding} },  // ${val.duration.toFixed(1)}s`
    );
  }
  console.log(`};`);
}

main().catch(console.error);
