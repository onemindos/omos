#!/usr/bin/env tsx
/**
 * Generate all V4 assets:
 * 1. Voice clips via OpenAI TTS (nova)
 * 2. Video clips via AWS Bedrock Nova Reel
 *
 * Usage: OPENAI_API_KEY=... npx tsx src/scripts/generate-v4-assets.ts
 */

import { scenes } from "./v4-script";
import fs from "fs";
import path from "path";
import { execSync } from "child_process";

const AUDIO_DIR = path.join(__dirname, "../../public/audio/v4");
const VIDEO_DIR = path.join(__dirname, "../../public/video/v4");
const S3_BUCKET = "s3://onemind-video-gen/v4";

async function generateVoice(sceneId: string, text: string): Promise<string> {
  const outPath = path.join(AUDIO_DIR, `${sceneId}.mp3`);
  if (fs.existsSync(outPath)) {
    console.log(`  [skip] ${sceneId}.mp3 already exists`);
    return outPath;
  }

  const apiKey = process.env.OPENAI_API_KEY;
  if (!apiKey) throw new Error("OPENAI_API_KEY not set");

  const res = await fetch("https://api.openai.com/v1/audio/speech", {
    method: "POST",
    headers: {
      Authorization: `Bearer ${apiKey}`,
      "Content-Type": "application/json",
    },
    body: JSON.stringify({
      model: "tts-1-hd",
      voice: "nova",
      input: text,
      speed: 0.92,
      response_format: "mp3",
    }),
  });

  if (!res.ok) throw new Error(`TTS failed for ${sceneId}: ${res.status} ${await res.text()}`);

  const buffer = Buffer.from(await res.arrayBuffer());
  fs.writeFileSync(outPath, buffer);
  console.log(`  [done] ${sceneId}.mp3 (${(buffer.length / 1024).toFixed(0)} KB)`);
  return outPath;
}

async function generateVideo(sceneId: string, prompt: string): Promise<string> {
  const markerPath = path.join(VIDEO_DIR, `${sceneId}.pending`);
  const mp4Path = path.join(VIDEO_DIR, `${sceneId}.mp4`);

  if (fs.existsSync(mp4Path)) {
    console.log(`  [skip] ${sceneId}.mp4 already exists`);
    return mp4Path;
  }

  // Start async Nova Reel generation
  const s3Output = `${S3_BUCKET}/${sceneId}/`;
  const modelInput = JSON.stringify({
    taskType: "TEXT_VIDEO",
    textToVideoParams: { text: prompt },
    videoGenerationConfig: {
      durationSeconds: 6,
      fps: 24,
      dimension: "1280x720",
      seed: 42,
    },
  });

  try {
    const result = execSync(
      `aws bedrock-runtime start-async-invoke ` +
        `--model-id amazon.nova-reel-v1:1 ` +
        `--region us-east-1 ` +
        `--model-input '${modelInput.replace(/'/g, "'\\''")}' ` +
        `--output-data-config '{"s3OutputDataConfig":{"s3Uri":"${s3Output}"}}' ` +
        `--output json`,
      { encoding: "utf-8", timeout: 30000 }
    );

    const { invocationArn } = JSON.parse(result);
    // Write marker file with ARN for polling
    fs.writeFileSync(markerPath, JSON.stringify({ arn: invocationArn, sceneId, startedAt: new Date().toISOString() }));
    console.log(`  [started] ${sceneId} → ${invocationArn.split("/").pop()}`);
    return invocationArn;
  } catch (e: any) {
    console.error(`  [error] ${sceneId}: ${e.message}`);
    return "";
  }
}

async function main() {
  // Ensure output dirs
  fs.mkdirSync(AUDIO_DIR, { recursive: true });
  fs.mkdirSync(VIDEO_DIR, { recursive: true });

  console.log("\n=== GENERATING VOICE CLIPS (OpenAI TTS) ===\n");
  for (const scene of scenes) {
    await generateVoice(scene.id, scene.voiceLine);
  }

  console.log("\n=== GENERATING VIDEO CLIPS (Nova Reel) ===\n");
  // Nova Reel can handle concurrent requests
  const videoPromises = scenes.map(scene => generateVideo(scene.id, scene.videoPrompt));
  await Promise.all(videoPromises);

  console.log("\n=== GENERATION LAUNCHED ===");
  console.log(`Voice clips: ${AUDIO_DIR}`);
  console.log(`Video jobs: check with 'npx tsx src/scripts/poll-v4-videos.ts'`);
  console.log("Nova Reel takes 2-4 minutes per clip.\n");
}

main().catch(console.error);
