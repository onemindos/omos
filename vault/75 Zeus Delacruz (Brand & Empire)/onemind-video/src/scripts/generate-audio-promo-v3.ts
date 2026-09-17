/**
 * Generate TTS for OneMind Promo v3 — ALL OpenAI voices
 *
 * Legacy = nova (warm, natural female)
 * Zeus   = echo (deep, grounded male)
 *
 * Usage: npx tsx src/scripts/generate-audio-promo-v3.ts
 */

import fs from "fs";
import path from "path";
import OpenAI from "openai";
import { PROMO_NARRATION } from "./promo-scenes";

const AUDIO_DIR = path.join(__dirname, "../../public/audio/promo");

const VOICES = {
  legacy: "nova" as const,   // warm, natural — NOT robotic Azure Jenny
  zeus: "echo" as const,     // deep, grounded
};

const SPEEDS = {
  legacy: 0.92,  // slightly slower for dramatic weight
  zeus: 1.0,
};

async function main() {
  const apiKey = process.env.OPENAI_API_KEY;
  if (!apiKey) {
    console.error("❌ Set OPENAI_API_KEY");
    process.exit(1);
  }

  const openai = new OpenAI({ apiKey });
  fs.mkdirSync(AUDIO_DIR, { recursive: true });

  console.log(`\n🎬 OneMind Promo v3 — Voice Generation (OpenAI only)`);
  console.log(`   Legacy: ${VOICES.legacy} @ ${SPEEDS.legacy}x`);
  console.log(`   Zeus:   ${VOICES.zeus} @ ${SPEEDS.zeus}x`);
  console.log(`   Lines:  ${PROMO_NARRATION.length}\n`);

  for (let i = 0; i < PROMO_NARRATION.length; i++) {
    const line = PROMO_NARRATION[i];
    const filename = `${String(i).padStart(2, "0")}-${line.speaker}-${line.scene}.mp3`;
    const filepath = path.join(AUDIO_DIR, filename);

    if (fs.existsSync(filepath)) {
      console.log(`⏭️  [${i}] Already exists: ${filename}`);
      continue;
    }

    const voice = VOICES[line.speaker];
    const speed = SPEEDS[line.speaker];

    console.log(`🎙️  [${i}] ${line.speaker.toUpperCase()} (${voice} @ ${speed}x): "${line.text.slice(0, 60)}..."`);

    try {
      const response = await openai.audio.speech.create({
        model: "tts-1-hd",
        voice,
        input: line.text,
        speed,
      });

      const buffer = Buffer.from(await response.arrayBuffer());
      fs.writeFileSync(filepath, buffer);
      console.log(`   ✅ ${filename} (${(buffer.length / 1024).toFixed(0)} KB)`);
    } catch (err: any) {
      console.error(`   ❌ ${err.message}`);
    }
  }

  // Write manifest
  const manifest = PROMO_NARRATION.map((line, i) => ({
    index: i,
    speaker: line.speaker,
    scene: line.scene,
    text: line.text,
    file: `${String(i).padStart(2, "0")}-${line.speaker}-${line.scene}.mp3`,
    voice: VOICES[line.speaker],
    speed: SPEEDS[line.speaker],
  }));

  fs.writeFileSync(path.join(AUDIO_DIR, "manifest.json"), JSON.stringify(manifest, null, 2));
  console.log(`\n✅ Done! Audio in ${AUDIO_DIR}`);
}

main().catch(console.error);
