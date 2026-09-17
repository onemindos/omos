/**
 * Generate TTS audio files for all narration lines
 *
 * Usage: npx tsx src/scripts/generate-audio.ts
 *
 * Requires OPENAI_API_KEY environment variable.
 * Generates .mp3 files in public/audio/
 */

import OpenAI from "openai";
import fs from "fs";
import path from "path";
import { NARRATION_SCRIPT, VOICE_MAP } from "./narration";

const AUDIO_DIR = path.join(__dirname, "../../public/audio");

async function generateAudio() {
  const apiKey = process.env.OPENAI_API_KEY;
  if (!apiKey) {
    console.error("❌ Set OPENAI_API_KEY environment variable");
    console.error("   export OPENAI_API_KEY=sk-...");
    process.exit(1);
  }

  const openai = new OpenAI({ apiKey });

  // Ensure output directory exists
  fs.mkdirSync(AUDIO_DIR, { recursive: true });

  console.log(`\n🎙️  Generating ${NARRATION_SCRIPT.length} audio clips...\n`);

  for (let i = 0; i < NARRATION_SCRIPT.length; i++) {
    const line = NARRATION_SCRIPT[i];
    const filename = `${String(i).padStart(2, "0")}-${line.speaker}-${line.scene}.mp3`;
    const filepath = path.join(AUDIO_DIR, filename);

    // Skip if already generated
    if (fs.existsSync(filepath)) {
      console.log(`⏭️  [${i}] Already exists: ${filename}`);
      continue;
    }

    const voice = VOICE_MAP[line.speaker];
    console.log(`🔊 [${i}] ${line.speaker.toUpperCase()} (${voice}): "${line.text.slice(0, 60)}..."`);

    const response = await openai.audio.speech.create({
      model: "tts-1-hd",
      voice,
      input: line.text,
      speed: 1.0,
    });

    const buffer = Buffer.from(await response.arrayBuffer());
    fs.writeFileSync(filepath, buffer);
    console.log(`   ✅ ${filename} (${(buffer.length / 1024).toFixed(0)} KB)`);
  }

  // Generate a manifest with file paths and durations
  const manifest = NARRATION_SCRIPT.map((line, i) => ({
    index: i,
    speaker: line.speaker,
    scene: line.scene,
    text: line.text,
    file: `${String(i).padStart(2, "0")}-${line.speaker}-${line.scene}.mp3`,
    estimatedDuration: line.durationEstimate,
  }));

  fs.writeFileSync(
    path.join(AUDIO_DIR, "manifest.json"),
    JSON.stringify(manifest, null, 2)
  );

  console.log(`\n✅ All audio generated in public/audio/`);
  console.log(`📋 Manifest: public/audio/manifest.json`);
  console.log(`\n💡 To replace Zeus's lines with real recordings:`);
  console.log(`   1. Record your voiceover for each "zeus" line`);
  console.log(`   2. Replace the corresponding .mp3 file in public/audio/`);
  console.log(`   3. Keep the same filename\n`);
}

generateAudio().catch(console.error);
