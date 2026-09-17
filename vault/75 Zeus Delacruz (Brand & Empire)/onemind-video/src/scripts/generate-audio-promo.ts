/**
 * Generate TTS for OneMind Promo narration (Azure Speech + OpenAI)
 *
 * Usage: npx tsx src/scripts/generate-audio-promo.ts
 *
 * Same Azure Jenny Neural voice as Composio video for consistency.
 */

import * as sdk from "microsoft-cognitiveservices-speech-sdk";
import fs from "fs";
import path from "path";
import { PROMO_NARRATION, AZURE_VOICE, VOICE_MAP } from "./promo-scenes";

const AUDIO_DIR = path.join(__dirname, "../../public/audio/promo");

function buildSSML(text: string, style?: string, styleDegree?: number): string {
  const voice = AZURE_VOICE.name;
  const lang = AZURE_VOICE.lang;
  const rate = AZURE_VOICE.rate;
  const pitch = AZURE_VOICE.pitch;

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

function synthesizeToFile(
  speechConfig: sdk.SpeechConfig,
  ssml: string,
  filepath: string
): Promise<void> {
  return new Promise((resolve, reject) => {
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
          reject(new Error(`Canceled: ${cancellation.reason} — ${cancellation.errorDetails}`));
        } else {
          reject(new Error(`Failed: ${result.reason}`));
        }
      },
      (err) => {
        synthesizer.close();
        reject(new Error(`Error: ${err}`));
      }
    );
  });
}

async function generatePromoAudio() {
  const speechKey = process.env.AZURE_SPEECH_KEY;
  const speechRegion = process.env.AZURE_SPEECH_REGION;

  if (!speechKey || !speechRegion) {
    console.error("❌ Set AZURE_SPEECH_KEY and AZURE_SPEECH_REGION");
    process.exit(1);
  }

  const speechConfig = sdk.SpeechConfig.fromSubscription(speechKey, speechRegion);
  speechConfig.speechSynthesisOutputFormat =
    sdk.SpeechSynthesisOutputFormat.Audio16Khz128KBitRateMonoMp3;

  let openai: any = null;
  const openaiKey = process.env.OPENAI_API_KEY;
  if (openaiKey) {
    const OpenAI = (await import("openai")).default;
    openai = new OpenAI({ apiKey: openaiKey });
    console.log("✅ OpenAI available for Zeus lines\n");
  } else {
    console.log("⚠️  No OPENAI_API_KEY — Zeus lines skipped\n");
  }

  fs.mkdirSync(AUDIO_DIR, { recursive: true });

  const legacyLines = PROMO_NARRATION.filter((l) => l.speaker === "legacy");
  const totalChars = legacyLines.reduce((sum, l) => sum + l.text.length, 0);

  console.log(`🎭 OneMind Promo — Voice Generation`);
  console.log(`   Voice: ${AZURE_VOICE.name} | ${legacyLines.length} Legacy lines, ${totalChars} chars\n`);

  for (let i = 0; i < PROMO_NARRATION.length; i++) {
    const line = PROMO_NARRATION[i];
    const filename = `${String(i).padStart(2, "0")}-${line.speaker}-${line.scene}.mp3`;
    const filepath = path.join(AUDIO_DIR, filename);

    if (fs.existsSync(filepath)) {
      console.log(`⏭️  [${i}] Already exists: ${filename}`);
      continue;
    }

    if (line.speaker === "legacy") {
      const styleLabel = line.azureStyle
        ? `${line.azureStyle} @${(line.azureStyleDegree ?? 1.0).toFixed(1)}`
        : "default";
      console.log(`🎭 [${i}] LEGACY (${styleLabel}): "${line.text.slice(0, 60)}..."`);

      try {
        const ssml = buildSSML(line.text, line.azureStyle, line.azureStyleDegree);
        await synthesizeToFile(speechConfig, ssml, filepath);
        const size = fs.statSync(filepath).size;
        console.log(`   ✅ ${filename} (${(size / 1024).toFixed(0)} KB)`);
      } catch (err: any) {
        console.error(`   ❌ ${err.message}`);
      }
    } else if (line.speaker === "zeus" && openai) {
      console.log(`🔊 [${i}] ZEUS: "${line.text.slice(0, 60)}..."`);
      try {
        const response = await openai.audio.speech.create({
          model: "tts-1-hd",
          voice: VOICE_MAP.zeus,
          input: line.text,
          speed: 1.0,
        });
        const buffer = Buffer.from(await response.arrayBuffer());
        fs.writeFileSync(filepath, buffer);
        console.log(`   ✅ ${filename} (${(buffer.length / 1024).toFixed(0)} KB)`);
      } catch (err: any) {
        console.error(`   ❌ ${err.message}`);
      }
    } else {
      console.log(`⏩ [${i}] ZEUS: skipped`);
    }
  }

  const manifest = PROMO_NARRATION.map((line, i) => ({
    index: i,
    speaker: line.speaker,
    scene: line.scene,
    text: line.text,
    file: `${String(i).padStart(2, "0")}-${line.speaker}-${line.scene}.mp3`,
    engine: line.speaker === "legacy" ? "azure-speech" : "openai-tts",
    estimatedDuration: line.durationEstimate,
  }));

  fs.writeFileSync(path.join(AUDIO_DIR, "manifest.json"), JSON.stringify(manifest, null, 2));
  console.log(`\n✅ Done! Audio in public/audio/promo/`);
}

generatePromoAudio().catch(console.error);
