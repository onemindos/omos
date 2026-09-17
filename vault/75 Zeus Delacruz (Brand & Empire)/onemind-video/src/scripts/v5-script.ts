/**
 * One Mind Promo V5 — Short & Sweet
 *
 * 5 scenes × ~6s = ~30s
 * Style: Fast, punchy, emotional. Every word earns its place.
 * Single narrator: Legacy (OpenAI nova @ 0.92x)
 */

export interface Scene {
  id: string;
  videoPrompt: string;
  voiceLine: string;
  textOverlay?: string;
  duration: number;
}

export const scenes: Scene[] = [
  // ── HOOK ──────────────────────────────────────────────────────
  {
    id: "01-hook",
    videoPrompt:
      "Cinematic close-up of a person in a dark room, face lit by the glow of dozens of phone notifications and app screens reflecting in their eyes, shallow depth of field, moody blue light, film grain",
    voiceLine: "You weren't built to live in fragments.",
    duration: 6,
  },

  // ── TURN ──────────────────────────────────────────────────────
  {
    id: "02-turn",
    videoPrompt:
      "A single teal ember igniting in darkness, expanding into a warm geometric light structure, particles coalescing into order, cinematic macro shot, solarpunk energy, golden and teal tones",
    voiceLine: "One mind. One method. You and AI — together.",
    textOverlay: "ONE MIND",
    duration: 6,
  },

  // ── METHOD ────────────────────────────────────────────────────
  {
    id: "03-method",
    videoPrompt:
      "Elegant top-down shot of five glowing steps forming a loop on a dark surface — each step pulses with teal light as the camera slowly orbits, minimal and clean, cinematic depth of field",
    voiceLine: "Five steps. Capture. Organize. Direct. Execute. Multiply.",
    textOverlay: "C.O.D.E.×",
    duration: 6,
  },

  // ── PROMISE ───────────────────────────────────────────────────
  {
    id: "04-promise",
    videoPrompt:
      "Breathtaking solarpunk cityscape at golden hour, a person walking confidently through lush organic architecture, holographic interfaces blending with nature, epic cinematic drone pull-back revealing the full city",
    voiceLine: "Not another app. Your operating system for life.",
    duration: 6,
  },

  // ── CTA ───────────────────────────────────────────────────────
  {
    id: "05-cta",
    videoPrompt:
      "One Mind logo materializing from teal light particles against pure black, minimal and elegant, the URL appears below in clean white type, cinematic",
    voiceLine: "Start building yours.",
    textOverlay: "codex.onemindos.com",
    duration: 6,
  },
];
