/**
 * Generate TTS audio files using Azure Speech (Jenny Neural + SSML emotions)
 *
 * Usage: npx tsx src/scripts/generate-audio-azure.ts
 *
 * Requires:
 *   AZURE_SPEECH_KEY    — Azure Speech resource key
 *   AZURE_SPEECH_REGION — Azure region (e.g., eastus)
 *
 * Legacy lines → Azure Speech (en-US-JennyNeural with per-line emotion styles)
 * Zeus lines   → OpenAI TTS fallback (or skip if no key)
 *
 * Cost: ~$16 per 1M characters. This video ≈ $0.01.
 */

import * as sdk from "microsoft-cognitiveservices-speech-sdk";
import fs from "fs";
import path from "path";
import { NARRATION_SCRIPT, AZURE_VOICE, VOICE_MAP } from "./narration";

const AUDIO_DIR = path.join(__dirname, "../../public/audio");

/**
 * Build Azure SSML for a narration line with emotion styling.
 * en-US-JennyNeural supports: angry, assistant, chat, cheerful, customerservice,
 * excited, friendly, hopeful, newscast, sad, shouting, terrified, unfriendly, whispering
 */
function buildSSML(text: string, style?: string, styleDegree?: number): string {
  const voice = AZURE_VOICE.name;
  const lang = AZURE_VOICE.lang;
  const rate = AZURE_VOICE.rate;
  const pitch = AZURE_VOICE.pitch;

  // Escape XML entities in text
  const escapedText = text
    .replace(/&/g, "&amp;")
    .replace(/</g, "&lt;")
    .replace(/>/g, "&gt;")
    .replace(/"/g, "&quot;");

  let innerContent: string;
  if (style) {
    const degree = styleDegree ?? 1.0;
    innerContent = `<mstts:express-as style="${style}" styledegree="${degree.toFixed(2)}">${escapedText}</mstts:express-as>`;
  } else {
    innerContent = escapedText;
  }

  return `<speak version="1.0" xmlns="http://www.w3.org/2001/10/synthesis"
  xmlns:mstts="https://www.w3.org/2001/mstts"
  xml:lang="${lang}">
  <voice name="${voice}">
    <prosody rate="${rate}" pitch="${pitch}">
      ${innerContent}
    </prosody>
  </voice>
</speak>`;
}

/**
 * Synthesize speech from SSML and save as MP3
 */
function synthesizeToFile(
  speechConfig: sdk.SpeechConfig,
  ssml: string,
  filepath: string
): Promise<void> {
  return new Promise((resolve, reject) => {
    // Use pull stream to collect audio in memory, then save
    const audioConfig = sdk.AudioConfig.fromAudioFileOutput(filepath);
    const synthesizer = new sdk.SpeechSynthesizer(speechConfig, audioConfig);

    synthesizer.speakSsmlAsync(
      ssml,
      (result) => {
        synthesizer.close();
        if (result.reason === sdk.ResultReason.SynthesizingAudioCompleted) {
          resolve();
        } else if (result.reason === sdk.ResultReason.Canceled) {
          const cancellation = sdk.CancellationDetails.fromResult(result);
          reject(
            new Error(
              `Synthesis canceled: ${cancellation.reason} — ${cancellation.errorDetails}`
            )
          );
        } else {
          reject(new Error(`Synthesis failed: ${result.reason}`));
        }
      },
      (err) => {
        synthesizer.close();
        reject(new Error(`Synthesis error: ${err}`));
      }
    );
  });
}

async function generateWithAzure() {
  const speechKey = process.env.AZURE_SPEECH_KEY;
  const speechRegion = process.env.AZURE_SPEECH_REGION;

  if (!speechKey || !speechRegion) {
    console.error("❌ Set AZURE_SPEECH_KEY and AZURE_SPEECH_REGION");
    console.error("   Azure Portal → Speech resource → Keys and Endpoint");
    console.error("   Example:");
    console.error("     export AZURE_SPEECH_KEY=abc123...");
    console.error("     export AZURE_SPEECH_REGION=eastus");
    process.exit(1);
  }

  // Configure Azure Speech
  const speechConfig = sdk.SpeechConfig.fromSubscription(speechKey, speechRegion);
  speechConfig.speechSynthesisOutputFormat =
    sdk.SpeechSynthesisOutputFormat.Audio16Khz128KBitRateMonoMp3;

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
  const estimatedCost = (totalChars / 1_000_000) * 16;

  console.log(`🎭 Azure Speech — Legacy's Voice (${AZURE_VOICE.name})`);
  console.log(`   ${legacyLines.length} lines, ${totalChars} characters`);
  console.log(`   Estimated cost: $${estimatedCost.toFixed(4)}`);
  console.log(`   Voice: ${AZURE_VOICE.name} | Rate: ${AZURE_VOICE.rate} | Pitch: ${AZURE_VOICE.pitch}`);
  console.log(`   SSML styles: per-line emotion control\n`);

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
      // ── AZURE SPEECH for Legacy ──
      const styleLabel = line.azureStyle
        ? `${line.azureStyle} @${(line.azureStyleDegree ?? 1.0).toFixed(1)}`
        : "default";
      console.log(`🎭 [${i}] LEGACY (Azure ${styleLabel}): "${line.text.slice(0, 55)}..."`);

      try {
        const ssml = buildSSML(line.text, line.azureStyle, line.azureStyleDegree);
        await synthesizeToFile(speechConfig, ssml, filepath);

        const fileSize = fs.statSync(filepath).size;
        console.log(`   ✅ ${filename} (${(fileSize / 1024).toFixed(0)} KB)`);
      } catch (err: any) {
        console.error(`   ❌ Azure error: ${err.message}`);

        // Fallback: try without emotion style
        if (line.azureStyle) {
          console.log(`   🔄 Retrying without emotion style...`);
          try {
            const ssml = buildSSML(line.text);
            await synthesizeToFile(speechConfig, ssml, filepath);
            const fileSize = fs.statSync(filepath).size;
            console.log(`   ✅ ${filename} (${(fileSize / 1024).toFixed(0)} KB) [no style]`);
          } catch (retryErr: any) {
            console.error(`   ❌ Retry failed: ${retryErr.message}`);
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
    azureStyle: line.azureStyle || null,
    azureStyleDegree: line.azureStyleDegree || null,
    file: `${String(i).padStart(2, "0")}-${line.speaker}-${line.scene}.mp3`,
    engine: line.speaker === "legacy" ? "azure-speech" : "openai-tts",
    voice: line.speaker === "legacy" ? AZURE_VOICE.name : VOICE_MAP.zeus,
    estimatedDuration: line.durationEstimate,
  }));

  fs.writeFileSync(
    path.join(AUDIO_DIR, "manifest.json"),
    JSON.stringify(manifest, null, 2)
  );

  console.log(`\n✅ Done! Audio in public/audio/`);
  console.log(`📋 Manifest: public/audio/manifest.json`);
  console.log(`💰 Total Legacy chars: ${totalChars} (~$${estimatedCost.toFixed(4)})`);
  console.log(`\n💡 Tips:`);
  console.log(`   • Delete a .mp3 to regenerate it`);
  console.log(`   • Replace Zeus .mp3 files with your real recordings`);
  console.log(`   • Tweak azureStyle and azureStyleDegree in narration.ts to adjust emotion\n`);
}

generateWithAzure().catch(console.error);
