/**
 * Generate TTS audio files using Hume AI Octave
 *
 * Usage: npx tsx src/scripts/generate-audio-hume.ts
 *
 * Requires HUME_API_KEY environment variable.
 * Legacy lines → Hume Octave (emotional, custom voice design)
 * Zeus lines → OpenAI TTS fallback (or skip if no key)
 *
 * Free tier: 10,000 chars/month (~10 min) — plenty for video narration
 */

import { HumeClient } from "hume";
import fs from "fs";
import path from "path";
import { NARRATION_SCRIPT, LEGACY_VOICE_DESCRIPTION, VOICE_MAP } from "./narration";

const AUDIO_DIR = path.join(__dirname, "../../public/audio");

async function generateWithHume() {
  const humeKey = process.env.HUME_API_KEY;
  if (!humeKey) {
    console.error("❌ Set HUME_API_KEY environment variable");
    console.error("   Get your free key at: https://app.hume.ai/sign-up");
    console.error("   Then: export HUME_API_KEY=...");
    process.exit(1);
  }

  const hume = new HumeClient({ apiKey: humeKey });

  // Optional: OpenAI for Zeus lines
  let openai: any = null;
  const openaiKey = process.env.OPENAI_API_KEY;
  if (openaiKey) {
    const OpenAI = (await import("openai")).default;
    openai = new OpenAI({ apiKey: openaiKey });
    console.log("✅ OpenAI available for Zeus lines\n");
  } else {
    console.log("⚠️  No OPENAI_API_KEY — Zeus lines will be skipped (record manually)\n");
  }

  fs.mkdirSync(AUDIO_DIR, { recursive: true });

  const legacyLines = NARRATION_SCRIPT.filter((l) => l.speaker === "legacy");
  const totalChars = legacyLines.reduce((sum, l) => sum + l.text.length, 0);

  console.log(`🎭 Hume Octave — Legacy's Voice`);
  console.log(`   ${legacyLines.length} lines, ${totalChars} characters (free tier: 10,000/month)`);
  console.log(`   Voice: Custom design from description`);
  console.log(`   Acting instructions: per-line emotional direction\n`);

  // Track the generation ID from the first Legacy line to reuse the voice
  let savedVoiceName: string | undefined;

  for (let i = 0; i < NARRATION_SCRIPT.length; i++) {
    const line = NARRATION_SCRIPT[i];
    const filename = `${String(i).padStart(2, "0")}-${line.speaker}-${line.scene}.mp3`;
    const filepath = path.join(AUDIO_DIR, filename);

    // Skip if already generated (delete files to regenerate)
    if (fs.existsSync(filepath)) {
      console.log(`⏭️  [${i}] Already exists: ${filename}`);
      continue;
    }

    if (line.speaker === "legacy") {
      // ── HUME OCTAVE for Legacy ──
      console.log(`🎭 [${i}] LEGACY (Hume Octave): "${line.text.slice(0, 55)}..."`);
      if (line.acting) {
        console.log(`   🎬 Acting: "${line.acting.slice(0, 65)}..."`);
      }

      try {
        // Build the utterance
        // - First call: description = voice design prompt (creates the voice)
        // - Subsequent calls with saved voice: description = acting instructions
        const utterance: any = {
          text: line.text,
          speed: 1.0,
          trailingSilence: 0.3,
        };

        if (savedVoiceName) {
          // Reuse the saved voice, use description for acting
          utterance.voice = { name: savedVoiceName, provider: "CUSTOM_VOICE" };
          if (line.acting) {
            utterance.description = line.acting;
          }
        } else {
          // First Legacy line — design the voice from description
          utterance.description = LEGACY_VOICE_DESCRIPTION;
        }

        const response = await hume.tts.synthesizeJson({
          utterances: [utterance],
          format: { type: "mp3" },
          numGenerations: 1,
        });

        // Extract audio from JSON response
        const generation = response.generations?.[0];
        if (!generation) {
          throw new Error("No generation returned");
        }

        // Save voice from first generation for consistency
        if (!savedVoiceName && generation.generationId) {
          console.log(`   💾 Saving voice for consistency (gen: ${generation.generationId})`);
          try {
            const voice = await hume.tts.voices.create({
              generationId: generation.generationId,
              name: `legacy-video-${Date.now()}`,
            });
            savedVoiceName = voice.name;
            console.log(`   ✅ Voice saved as: ${savedVoiceName}`);
          } catch (voiceErr: any) {
            console.log(`   ⚠️  Could not save voice: ${voiceErr.message} — will use description each time`);
          }
        }

        // Get audio data — generation.audio is the full output as base64
        if (!generation.audio) {
          throw new Error("No audio in response");
        }

        const buffer = Buffer.from(generation.audio, "base64");
        fs.writeFileSync(filepath, buffer);
        console.log(`   ✅ ${filename} (${(buffer.length / 1024).toFixed(0)} KB)`);
      } catch (err: any) {
        console.error(`   ❌ Hume error: ${err.body?.message || err.message || JSON.stringify(err)}`);

        // Fallback: try without acting, just voice description
        if (line.acting) {
          console.log(`   🔄 Retrying with description only...`);
          try {
            const response = await hume.tts.synthesizeJson({
              utterances: [{
                text: line.text,
                description: LEGACY_VOICE_DESCRIPTION,
                speed: 1.0,
              }],
              format: { type: "mp3" },
              numGenerations: 1,
            });
            const gen = response.generations?.[0];
            if (gen?.audio) {
              const buffer = Buffer.from(gen.audio, "base64");
              fs.writeFileSync(filepath, buffer);
              console.log(`   ✅ ${filename} (${(buffer.length / 1024).toFixed(0)} KB) [fallback]`);
            }
          } catch (retryErr: any) {
            console.error(`   ❌ Retry failed: ${retryErr.body?.message || retryErr.message}`);
          }
        }
      }
    } else if (line.speaker === "zeus" && openai) {
      // ── OPENAI TTS for Zeus ──
      const voice = VOICE_MAP[line.speaker];
      console.log(`🔊 [${i}] ZEUS (OpenAI ${voice}): "${line.text.slice(0, 55)}..."`);

      try {
        const response = await openai.audio.speech.create({
          model: "tts-1-hd",
          voice,
          input: line.text,
          speed: 1.0,
        });
        const buffer = Buffer.from(await response.arrayBuffer());
        fs.writeFileSync(filepath, buffer);
        console.log(`   ✅ ${filename} (${(buffer.length / 1024).toFixed(0)} KB)`);
      } catch (err: any) {
        console.error(`   ❌ OpenAI error: ${err.message || err}`);
      }
    } else {
      console.log(`⏩ [${i}] ZEUS: skipped (no OpenAI key — record manually)`);
    }
  }

  // Generate manifest
  const manifest = NARRATION_SCRIPT.map((line, i) => ({
    index: i,
    speaker: line.speaker,
    scene: line.scene,
    text: line.text,
    acting: line.acting || null,
    file: `${String(i).padStart(2, "0")}-${line.speaker}-${line.scene}.mp3`,
    engine: line.speaker === "legacy" ? "hume-octave" : "openai-tts",
    estimatedDuration: line.durationEstimate,
  }));

  fs.writeFileSync(
    path.join(AUDIO_DIR, "manifest.json"),
    JSON.stringify(manifest, null, 2)
  );

  console.log(`\n✅ Done! Audio in public/audio/`);
  console.log(`📋 Manifest: public/audio/manifest.json`);
  console.log(`\n💡 Tips:`);
  console.log(`   • Delete a .mp3 file to regenerate it`);
  console.log(`   • Replace Zeus .mp3 files with your real recordings`);
  console.log(`   • Free tier: ${totalChars} of 10,000 chars used this run\n`);
}

generateWithHume().catch(console.error);
