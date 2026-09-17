/**
 * Generate AI background images for OneMind Promo scenes
 *
 * Usage: npx tsx src/scripts/generate-images.ts
 *
 * Uses GPT-Image-1 (gpt-image-1) via OpenAI API
 * Generates 1920x1080 solarpunk backgrounds for each scene
 *
 * Requires: OPENAI_API_KEY environment variable
 * Cost: ~$0.04-0.17 per image depending on quality/size
 */

import OpenAI from "openai";
import fs from "fs";
import path from "path";
import { SCENE_IMAGES } from "./promo-scenes";

const IMAGE_DIR = path.join(__dirname, "../../public/images");

async function generateImages() {
  const apiKey = process.env.OPENAI_API_KEY;
  if (!apiKey) {
    console.error("❌ Set OPENAI_API_KEY environment variable");
    process.exit(1);
  }

  const openai = new OpenAI({ apiKey });
  fs.mkdirSync(IMAGE_DIR, { recursive: true });

  console.log(`🎨 GPT-Image-1 — Solarpunk Scene Backgrounds`);
  console.log(`   ${SCENE_IMAGES.length} scenes to generate`);
  console.log(`   Output: public/images/\n`);

  for (const scene of SCENE_IMAGES) {
    const filename = `scene-${scene.id}.png`;
    const filepath = path.join(IMAGE_DIR, filename);

    // Skip if already generated
    if (fs.existsSync(filepath)) {
      console.log(`⏭️  Already exists: ${filename}`);
      continue;
    }

    console.log(`🖼️  Generating: ${scene.id}`);
    console.log(`   Prompt: "${scene.prompt.slice(0, 80)}..."`);

    try {
      const response = await openai.images.generate({
        model: "gpt-image-1",
        prompt: scene.prompt,
        n: 1,
        size: "1536x1024", // closest to 16:9 widescreen
        quality: "high",
      });

      const imageData = response.data?.[0];
      if (!imageData) {
        throw new Error("No image data returned");
      }

      // gpt-image-1 returns base64 by default
      if (imageData.b64_json) {
        const buffer = Buffer.from(imageData.b64_json, "base64");
        fs.writeFileSync(filepath, buffer);
        console.log(`   ✅ ${filename} (${(buffer.length / 1024 / 1024).toFixed(1)} MB)`);
      } else if (imageData.url) {
        // Fallback: download from URL
        const res = await fetch(imageData.url);
        const buffer = Buffer.from(await res.arrayBuffer());
        fs.writeFileSync(filepath, buffer);
        console.log(`   ✅ ${filename} (${(buffer.length / 1024 / 1024).toFixed(1)} MB)`);
      } else {
        throw new Error("No image data (b64 or url) in response");
      }
    } catch (err: any) {
      console.error(`   ❌ Error: ${err.message}`);

      // If gpt-image-1 not available, fall back to dall-e-3
      if (err.message?.includes("model") || err.status === 404) {
        console.log(`   🔄 Falling back to dall-e-3...`);
        try {
          const response = await openai.images.generate({
            model: "dall-e-3",
            prompt: scene.prompt,
            n: 1,
            size: "1792x1024",
            quality: "hd",
            style: "vivid",
          });

          const url = response.data?.[0]?.url;
          if (url) {
            const res = await fetch(url);
            const buffer = Buffer.from(await res.arrayBuffer());
            fs.writeFileSync(filepath, buffer);
            console.log(`   ✅ ${filename} (${(buffer.length / 1024 / 1024).toFixed(1)} MB) [dall-e-3]`);
          }
        } catch (fallbackErr: any) {
          console.error(`   ❌ Fallback failed: ${fallbackErr.message}`);
        }
      }
    }
  }

  // Write scene manifest
  const manifest = SCENE_IMAGES.map((scene) => ({
    id: scene.id,
    file: `scene-${scene.id}.png`,
    kenBurns: scene.kenBurns,
    prompt: scene.prompt.slice(0, 100) + "...",
  }));

  fs.writeFileSync(
    path.join(IMAGE_DIR, "manifest.json"),
    JSON.stringify(manifest, null, 2)
  );

  console.log(`\n✅ Done! Images in public/images/`);
  console.log(`📋 Manifest: public/images/manifest.json`);
  console.log(`\n💡 Delete a .png to regenerate it`);
}

generateImages().catch(console.error);
