#!/usr/bin/env tsx
/**
 * Poll Nova Reel video generation jobs and download completed videos.
 *
 * Usage: npx tsx src/scripts/poll-v4-videos.ts
 */

import fs from "fs";
import path from "path";
import { execSync } from "child_process";

const VIDEO_DIR = path.join(__dirname, "../../public/video/v4");

interface PendingJob {
  arn: string;
  sceneId: string;
  startedAt: string;
}

function checkJob(job: PendingJob): "InProgress" | "Completed" | "Failed" {
  try {
    const result = execSync(
      `aws bedrock-runtime get-async-invoke ` +
        `--invocation-arn "${job.arn}" ` +
        `--region us-east-1 --output json`,
      { encoding: "utf-8", timeout: 15000 }
    );
    const data = JSON.parse(result);
    return data.status;
  } catch {
    return "Failed";
  }
}

function downloadVideo(job: PendingJob): boolean {
  try {
    // Get S3 output path
    const result = execSync(
      `aws bedrock-runtime get-async-invoke ` +
        `--invocation-arn "${job.arn}" ` +
        `--region us-east-1 --output json`,
      { encoding: "utf-8", timeout: 15000 }
    );
    const data = JSON.parse(result);
    const s3Uri = data.outputDataConfig?.s3OutputDataConfig?.s3Uri;

    if (!s3Uri) {
      console.error(`  No S3 URI for ${job.sceneId}`);
      return false;
    }

    // Download the output.mp4
    const localPath = path.join(VIDEO_DIR, `${job.sceneId}.mp4`);
    const s3Path = `${s3Uri}output.mp4`;

    // Try output.mp4 first, then list the bucket to find the file
    try {
      execSync(`aws s3 cp "${s3Path}" "${localPath}" --region us-east-1`, {
        encoding: "utf-8",
        timeout: 60000,
      });
    } catch {
      // List what's in the output dir and grab the first mp4
      const listing = execSync(`aws s3 ls "${s3Uri}" --region us-east-1`, {
        encoding: "utf-8",
        timeout: 15000,
      });
      console.log(`  Files at ${s3Uri}: ${listing.trim()}`);

      const mp4File = listing.split("\n").find(l => l.includes(".mp4"));
      if (mp4File) {
        const filename = mp4File.trim().split(/\s+/).pop();
        execSync(`aws s3 cp "${s3Uri}${filename}" "${localPath}" --region us-east-1`, {
          encoding: "utf-8",
          timeout: 60000,
        });
      } else {
        console.error(`  No mp4 found at ${s3Uri}`);
        return false;
      }
    }

    const size = fs.statSync(localPath).size;
    console.log(`  [downloaded] ${job.sceneId}.mp4 (${(size / 1024 / 1024).toFixed(1)} MB)`);
    return true;
  } catch (e: any) {
    console.error(`  [error] downloading ${job.sceneId}: ${e.message}`);
    return false;
  }
}

async function main() {
  const pendingFiles = fs.readdirSync(VIDEO_DIR).filter(f => f.endsWith(".pending"));

  if (pendingFiles.length === 0) {
    console.log("No pending video jobs found.");
    return;
  }

  console.log(`\nChecking ${pendingFiles.length} pending video jobs...\n`);

  let completed = 0;
  let inProgress = 0;
  let failed = 0;

  for (const file of pendingFiles) {
    const job: PendingJob = JSON.parse(fs.readFileSync(path.join(VIDEO_DIR, file), "utf-8"));
    const status = checkJob(job);

    if (status === "Completed") {
      console.log(`  ${job.sceneId}: DONE — downloading...`);
      if (downloadVideo(job)) {
        // Remove pending marker
        fs.unlinkSync(path.join(VIDEO_DIR, file));
        completed++;
      } else {
        failed++;
      }
    } else if (status === "InProgress") {
      const elapsed = ((Date.now() - new Date(job.startedAt).getTime()) / 1000).toFixed(0);
      console.log(`  ${job.sceneId}: generating (${elapsed}s elapsed)`);
      inProgress++;
    } else {
      console.log(`  ${job.sceneId}: FAILED`);
      failed++;
    }
  }

  console.log(`\n--- Summary: ${completed} downloaded, ${inProgress} in progress, ${failed} failed ---`);

  if (inProgress > 0) {
    console.log("Run this script again in 1-2 minutes to check remaining jobs.\n");
  } else if (completed > 0 && failed === 0) {
    console.log("\nAll videos downloaded! Ready to compose.\n");
  }
}

main().catch(console.error);
