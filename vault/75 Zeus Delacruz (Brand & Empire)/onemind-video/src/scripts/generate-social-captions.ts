import OpenAI from "openai";

/**
 * generate-social-captions.ts — Generate platform-specific captions for video clips
 *
 * Usage:
 *   OPENAI_API_KEY=... npx tsx src/scripts/generate-social-captions.ts "Video Title" "Brief description"
 *
 * Generates captions optimized for:
 *   - YouTube Shorts
 *   - TikTok
 *   - Instagram Reels
 *   - Twitter/X (thread + clip post)
 *   - LinkedIn
 */

const SYSTEM_PROMPT = `You are a social media manager for Zeus Delacruz's brand "One Mind."

Brand: "You + AI = One Mind" — AI-native human enhancement framework
Tone: Direct, confident, no fluff, slightly provocative, thought-leader energy
Zeus's AI: Legacy (female voice, data-driven)

Generate captions for ALL platforms in this exact format. Each platform has different rules:

---

## YOUTUBE SHORTS
[Short, punchy caption. 1-2 sentences max. Include 3 hashtags.]
#OneMind #AI #[relevant]

## TIKTOK
[Hook in first line — must stop the scroll. 2-3 lines max. Casual/bold tone. 3-5 hashtags.]
#oneminds #ai #[relevant] #fyp #techtok

## INSTAGRAM REELS
[Storytelling caption. 3-5 lines. Emoji-friendly. End with a question or CTA. 20-30 hashtags in comment block.]

Caption:
[main caption]

Hashtags (paste as first comment):
#[30 relevant hashtags]

## TWITTER/X — CLIP POST
[Bold one-liner for the clip post. Must be tweetable standalone. Under 200 chars.]

## TWITTER/X — THREAD (5 tweets)
1/ [Hook — contrarian or surprising take]
2/ [Context / problem statement]
3/ [The insight / solution]
4/ [Specific example or data point]
5/ [CTA + link to full video]

## LINKEDIN
[Professional tone. 3-5 short paragraphs. Start with a hook line. End with a question. No hashtags in body — add 3-5 at the very end.]

---

Rules:
- Each platform's tone should match its culture
- TikTok = casual, bold, Gen Z energy
- LinkedIn = professional, thought-leader, no cringe
- Twitter = sharp, concise, quotable
- Instagram = visual storytelling, community
- YouTube Shorts = minimal, let the video speak
- Always reference the full YouTube video where relevant
- Never use "guys" — use "you" or "y'all"`;

async function generateCaptions(title: string, description: string): Promise<void> {
  const apiKey = process.env.OPENAI_API_KEY;
  if (!apiKey) {
    console.error("Error: OPENAI_API_KEY not set");
    process.exit(1);
  }

  const openai = new OpenAI({ apiKey });

  console.log(`Generating social captions for: "${title}"...`);

  const response = await openai.chat.completions.create({
    model: "gpt-4o",
    messages: [
      { role: "system", content: SYSTEM_PROMPT },
      {
        role: "user",
        content: `Generate platform-specific social media captions for this video:\n\nTitle: ${title}\nContent: ${description}\n\nThe full video is on YouTube. Clips will be posted across all platforms.`,
      },
    ],
    temperature: 0.8,
    max_tokens: 2000,
  });

  const captions = response.choices[0]?.message?.content;
  if (!captions) {
    console.error("Error: No captions generated");
    process.exit(1);
  }

  console.log("\n" + "=".repeat(60));
  console.log(captions);
  console.log("=".repeat(60));
  console.log("\nCopy each platform's caption when posting.");
}

// Parse args
const args = process.argv.slice(2);

if (args.length < 2) {
  console.log('Usage: npx tsx src/scripts/generate-social-captions.ts "Video Title" "Brief description"');
  console.log("\nExample:");
  console.log(
    '  npx tsx src/scripts/generate-social-captions.ts "I Built an AI That Runs My Life" "Zeus shows how Legacy AI manages his daily tasks using the CODEX method"'
  );
  process.exit(0);
}

generateCaptions(args[0], args[1]);
