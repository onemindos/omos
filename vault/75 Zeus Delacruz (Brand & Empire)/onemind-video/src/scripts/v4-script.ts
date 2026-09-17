/**
 * One Mind Promo V4 — Cinematic Manifesto
 *
 * AI video clips (Nova Reel) + Legacy voice (OpenAI nova) + minimal text
 * Target: 60s main cut + 30s teaser
 *
 * Style: Cinematic manifesto — short, punchy, emotional
 * Single narrator: Legacy (nova voice)
 */

// ─── SCENE DEFINITIONS ──────────────────────────────────────────
// Each scene: AI video clip (6s) + voice line + optional text overlay
// Total: ~60s (10 scenes × 6s)

export interface Scene {
  id: string;
  /** Text prompt for Nova Reel video generation */
  videoPrompt: string;
  /** Legacy's voice line */
  voiceLine: string;
  /** Minimal text overlay (shown on screen) — keep to 3-6 words max */
  textOverlay?: string;
  /** Duration in seconds */
  duration: number;
  /** Whether to include in 30s teaser cut */
  inTeaser: boolean;
}

export const scenes: Scene[] = [
  // ── ACT 1: THE TENSION (0-18s) ────────────────────────────────
  {
    id: "01-disconnect",
    videoPrompt: "Cinematic aerial shot of a person surrounded by floating holographic screens and notifications in a dark room, overwhelmed, blue and orange light, moody atmosphere, volumetric fog, 4K cinematic",
    voiceLine: "You weren't built to live in fragments.",
    textOverlay: undefined, // let the visual + voice breathe
    duration: 6,
    inTeaser: true,
  },
  {
    id: "02-chaos",
    videoPrompt: "Abstract visualization of hundreds of app icons and browser tabs spiraling into chaos, data streams fragmenting, digital entropy, dark background with scattered neon lights, cinematic motion blur",
    voiceLine: "Forty-seven apps. Zero clarity. Your tools multiplied — but you didn't.",
    textOverlay: undefined,
    duration: 6,
    inTeaser: true,
  },
  {
    id: "03-question",
    videoPrompt: "Dramatic slow motion shot of a single glowing teal ember forming in complete darkness, macro lens, particles gathering, solarpunk energy emerging from void, cinematic depth of field",
    voiceLine: "What if there was a system — built for how you actually think?",
    textOverlay: undefined,
    duration: 6,
    inTeaser: true,
  },

  // ── ACT 2: THE ANSWER (18-42s) ────────────────────────────────
  {
    id: "04-onemind",
    videoPrompt: "Cinematic reveal of a glowing teal geometric symbol materializing from particles, solarpunk aesthetic, warm gold and teal light, volumetric rays, dark background transitioning to lush green, epic scale",
    voiceLine: "One Mind. You and AI — running your life together.",
    textOverlay: "ONE MIND",
    duration: 6,
    inTeaser: true,
  },
  {
    id: "05-codex",
    videoPrompt: "Beautiful slow pan across an infinite library of glowing markdown documents floating in space, organized in luminous teal columns, each page gently pulsing with light, cinematic solarpunk atmosphere",
    voiceLine: "One method. Five steps. Capture, Organize, Direct, Execute, Multiply.",
    textOverlay: "C.O.D.E.×",
    duration: 6,
    inTeaser: true,
  },
  {
    id: "06-domains",
    videoPrompt: "Three luminous pillars of light rising from a glowing teal platform, each pillar a different warm color — green gold and coral — representing life domains, solarpunk architecture, epic cinematic wide shot",
    voiceLine: "Three domains. Your self. Your foundation. Your endeavor. One operating layer beneath it all.",
    textOverlay: undefined,
    duration: 6,
    inTeaser: false,
  },
  {
    id: "07-vault",
    videoPrompt: "Hands interacting with a beautiful holographic vault interface, plain text files transforming into structured knowledge, warm teal glow, intimate close-up, solarpunk office environment, natural light mixing with digital elements",
    voiceLine: "Plain markdown. No lock-in. Every AI can read it. Every human can understand it.",
    textOverlay: undefined,
    duration: 6,
    inTeaser: false,
  },
  {
    id: "08-compound",
    videoPrompt: "Cinematic time-lapse of a single seed growing into a massive luminous tree, each branch representing expanding knowledge domains, solarpunk city visible in background, golden hour light, epic scale drone shot",
    voiceLine: "Every review compounds. Every cycle multiplies. The system gets smarter — because you do.",
    textOverlay: undefined,
    duration: 6,
    inTeaser: false,
  },

  // ── ACT 3: THE CALL (42-60s) ──────────────────────────────────
  {
    id: "09-future",
    videoPrompt: "Breathtaking solarpunk cityscape at golden hour, organic architecture with lush vegetation, holographic displays integrated with nature, a person walking confidently toward the horizon, cinematic lens flare",
    voiceLine: "This isn't productivity. This is enhancement. This is your operating system for life.",
    textOverlay: undefined,
    duration: 6,
    inTeaser: true,
  },
  {
    id: "10-cta",
    videoPrompt: "Elegant minimal animation of the One Mind logo forming from teal light particles against pure black, cinematic, the text 'codex.onemindos.com' appears below in clean white typography",
    voiceLine: "One Mind. Start building yours.",
    textOverlay: "codex.onemindos.com",
    duration: 6,
    inTeaser: true,
  },
];

// ── TEASER CUT (30s) ──────────────────────────────────────────
// Scenes marked inTeaser: 01, 02, 03, 04, 05, 09, 10 = 7 scenes × ~4.3s each
// We'll trim these to ~4s each for a tight 30s cut

// ── GENERATION PROMPTS ─────────────────────────────────────────
export function getVideoPrompts(): { id: string; prompt: string }[] {
  return scenes.map(s => ({ id: s.id, prompt: s.videoPrompt }));
}

export function getVoiceLines(): { id: string; text: string }[] {
  return scenes.map(s => ({ id: s.id, text: s.voiceLine }));
}

export function getTeaserScenes(): Scene[] {
  return scenes.filter(s => s.inTeaser);
}
