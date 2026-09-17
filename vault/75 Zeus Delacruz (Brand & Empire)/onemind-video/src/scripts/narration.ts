/**
 * Composio CLI Video — Narration Script
 *
 * Speaker: LEGACY = Azure Speech TTS (en-US-JennyNeural with SSML emotion styles)
 * Speaker: ZEUS   = Real voiceover (recorded separately), falls back to OpenAI TTS
 *
 * Total: ~60 seconds
 */

/**
 * Azure Speech SSML style values for en-US-JennyNeural:
 * angry, assistant, chat, cheerful, customerservice, excited, friendly,
 * hopeful, newscast, sad, shouting, terrified, unfriendly, whispering
 *
 * styleDegree: 0.01 (subtle) → 2.0 (intense), default 1.0
 */
export type AzureVoiceStyle =
  | "angry" | "assistant" | "chat" | "cheerful" | "customerservice"
  | "excited" | "friendly" | "hopeful" | "newscast" | "sad"
  | "shouting" | "terrified" | "unfriendly" | "whispering";

export interface NarrationLine {
  speaker: "legacy" | "zeus";
  text: string;
  scene: string;
  durationEstimate: number; // seconds
  /** Acting direction (human-readable) */
  acting?: string;
  /** Azure SSML emotion style */
  azureStyle?: AzureVoiceStyle;
  /** Azure SSML style intensity: 0.01 (subtle) → 2.0 (intense) */
  azureStyleDegree?: number;
}

export const NARRATION_SCRIPT: NarrationLine[] = [
  // ── SCENE 1: HOOK (0-8s) ──
  {
    speaker: "legacy",
    text: "Your AI agents are smart. But they can't do anything without tools.",
    scene: "hook",
    durationEstimate: 4,
    acting: "Confident and direct, like stating an obvious truth. Slight edge — almost challenging the listener.",
    azureStyle: "chat",
    azureStyleDegree: 1.4,
  },
  {
    speaker: "zeus",
    text: "What if you could connect them to a thousand apps — from your terminal?",
    scene: "hook",
    durationEstimate: 4,
  },

  // ── SCENE 2: WHAT IS COMPOSIO (8-20s) ──
  {
    speaker: "legacy",
    text: "Meet Composio CLI. One command to give your agents the tools they need — authenticated, scoped, and production-ready.",
    scene: "intro",
    durationEstimate: 7,
    acting: "Warm and inviting, like introducing a trusted ally. Build energy through the sentence — land hard on 'production-ready'.",
    azureStyle: "friendly",
    azureStyleDegree: 1.3,
  },
  {
    speaker: "zeus",
    text: "GitHub, Gmail, Slack, Notion, Stripe — over a thousand integrations, zero auth headaches.",
    scene: "intro",
    durationEstimate: 5,
  },

  // ── SCENE 3: HOW IT WORKS (20-35s) ──
  {
    speaker: "legacy",
    text: "Step one: Connect your apps. Composio handles all the OAuth and token refresh automatically.",
    scene: "how-it-works",
    durationEstimate: 5,
    acting: "Clear and instructional, like a tactical briefing. Calm authority — she's walked you through this a hundred times.",
    azureStyle: "assistant",
    azureStyleDegree: 1.0,
  },
  {
    speaker: "legacy",
    text: "Step two: Composio scopes the tools. Every call is scoped to the right user with the right permissions.",
    scene: "how-it-works",
    durationEstimate: 5,
    acting: "Precise and focused, emphasizing security. A hint of protectiveness — she takes permissions seriously.",
    azureStyle: "newscast",
    azureStyleDegree: 1.2,
  },
  {
    speaker: "legacy",
    text: "Step three: Your agent executes. Tool calls run in sandboxed environments. Every action is logged and traceable.",
    scene: "how-it-works",
    durationEstimate: 5,
    acting: "Building to a crescendo — this is the payoff. Energetic but controlled, like a mission commander giving the green light.",
    azureStyle: "excited",
    azureStyleDegree: 1.5,
  },

  // ── SCENE 4: KEY FEATURES (35-48s) ──
  {
    speaker: "zeus",
    text: "It's model-agnostic. Works with Claude, GPT, Gemini, LangChain, CrewAI — any framework.",
    scene: "features",
    durationEstimate: 5,
  },
  {
    speaker: "legacy",
    text: "Smart tool search resolves by intent, not configuration. Context-aware sessions mean your agent never starts from scratch.",
    scene: "features",
    durationEstimate: 6,
    acting: "Genuinely impressed and excited — this is the part she finds brilliant. Speak with conviction and a knowing smile.",
    azureStyle: "cheerful",
    azureStyleDegree: 1.6,
  },

  // ── SCENE 5: DEMO TEASER (48-55s) ──
  {
    speaker: "zeus",
    text: "One install. One login. And your agents go from asking questions to doing real work.",
    scene: "demo",
    durationEstimate: 5,
  },

  // ── SCENE 6: CTA (55-60s) ──
  {
    speaker: "legacy",
    text: "Composio CLI. Connect AI agents to a thousand apps from your terminal. Get started at composio.dev.",
    scene: "cta",
    durationEstimate: 6,
    acting: "Strong and inspiring close. Like a trusted guide sending you off on a mission — warm but with gravitas.",
    azureStyle: "hopeful",
    azureStyleDegree: 1.4,
  },
];

// Azure Speech voice config — consistent across all Legacy lines
export const AZURE_VOICE = {
  name: "en-US-JennyNeural",  // Warm, confident, multiple emotion styles — Cortana spiritual successor
  lang: "en-US",
  rate: "-3%",                // Slightly slower for natural conversational feel
  pitch: "+0%",               // Natural pitch — no artificial brightness
} as const;

// Legacy's voice description (used in docs/prompts, not by Azure directly)
export const LEGACY_VOICE_DESCRIPTION =
  "A warm, confident female AI companion with an edge of toughness. " +
  "Energetic but grounded — like a tactical AI advisor who genuinely cares about your success. " +
  "She's direct and commanding when the situation calls for it, but always with underlying warmth " +
  "and a spark of excitement about the future. Think Cortana from Halo — brilliant, loyal, " +
  "occasionally playful, never passive. Her voice carries authority without coldness.";

// OpenAI TTS voice mapping (fallback for Zeus lines)
export const VOICE_MAP = {
  legacy: "nova" as const,   // OpenAI fallback — not used when Azure is available
  zeus: "echo" as const,     // Placeholder for Zeus's real voice (will be replaced with recording)
} as const;
