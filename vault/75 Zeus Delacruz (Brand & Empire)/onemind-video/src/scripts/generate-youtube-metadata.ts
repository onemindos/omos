import OpenAI from "openai";

/**
 * generate-youtube-metadata.ts — Generate SEO-optimized YouTube metadata
 *
 * Usage:
 *   OPENAI_API_KEY=... npx tsx src/scripts/generate-youtube-metadata.ts "Video Title" "Brief description of content"
 *
 * Generates:
 *   - Optimized title (with keyword front-loading)
 *   - Full description with timestamps template
 *   - 30 relevant tags
 *   - Thumbnail text/prompt idea
 *   - Hashtags
 */

const SYSTEM_PROMPT = `You are a YouTube SEO expert for the channel "One Mind" by Zeus Delacruz.

Channel niche: AI + productivity + personal development + life systems
Audience: Tech-savvy 25-45 year olds interested in AI, self-improvement, and building systems
Brand: "You + AI = One Mind" — AI-native human enhancement framework

Generate YouTube metadata in this EXACT format:

TITLE: [60 chars max, keyword-loaded, attention-grabbing. Use patterns like "I Built...", "How I...", "Stop Doing X...", "The X Nobody Talks About"]

DESCRIPTION:
[First 2 lines are CRITICAL — they show in search. Hook + value promise.]

[Empty line]

In this video, [2-3 sentence summary of what viewer will learn]

[Empty line]

TIMESTAMPS:
0:00 - [Hook/Intro]
0:30 - [Section 1 name]
3:00 - [Legacy AI Commentary]
4:00 - [Section 2 name]
7:00 - [Legacy Recap]
7:45 - [CTA]

[Empty line]

RESOURCES:
- One Mind Codex: https://codex.onemindos.com
- [Any other relevant links]

[Empty line]

CONNECT:
- Twitter/X: @zeuslegacy
- GitHub: github.com/OneMind-OS

#OneMind #AI #[3-5 more relevant hashtags]

TAGS (comma-separated, 30 tags):
[Mix of broad + niche + long-tail keywords]

THUMBNAIL:
Text: [2-4 words max, large bold text for thumbnail]
Visual: [Brief description of thumbnail composition]
Colors: [Suggest colors that pop — yellows, reds, or contrasting]

CATEGORY: [Science & Technology | Education | People & Blogs | Howto & Style]`;

async function generateMetadata(title: string, description: string): Promise<void> {
  const apiKey = process.env.OPENAI_API_KEY;
  if (!apiKey) {
    console.error("Error: OPENAI_API_KEY not set");
    process.exit(1);
  }

  const openai = new OpenAI({ apiKey });

  console.log(`Generating YouTube metadata for: "${title}"...`);

  const response = await openai.chat.completions.create({
    model: "gpt-4o",
    messages: [
      { role: "system", content: SYSTEM_PROMPT },
      {
        role: "user",
        content: `Generate optimized YouTube metadata for this video:\n\nTitle/Topic: ${title}\nContent: ${description}`,
      },
    ],
    temperature: 0.7,
    max_tokens: 1500,
  });

  const metadata = response.choices[0]?.message?.content;
  if (!metadata) {
    console.error("Error: No metadata generated");
    process.exit(1);
  }

  console.log("\n" + "=".repeat(60));
  console.log(metadata);
  console.log("=".repeat(60));
  console.log("\nCopy the above into YouTube Studio. Adjust timestamps after editing.");
}

// Parse args
const args = process.argv.slice(2);

if (args.length < 2) {
  console.log('Usage: npx tsx src/scripts/generate-youtube-metadata.ts "Video Title" "Brief description"');
  console.log("\nExample:");
  console.log(
    '  npx tsx src/scripts/generate-youtube-metadata.ts "How I Use AI to Manage My Entire Life" "Walkthrough of the One Mind CODEX method and how Legacy AI helps me organize everything"'
  );
  process.exit(0);
}

generateMetadata(args[0], args[1]);
