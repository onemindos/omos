/**
 * OneMind Promo — Scene definitions for AI image generation + narration
 *
 * Each scene has:
 *   - Image prompt for GPT-Image-1 (solarpunk style)
 *   - Narration lines with Azure SSML styles
 *   - Ken Burns camera motion (zoom, pan direction)
 */

import type { AzureVoiceStyle } from "./narration";

export interface SceneImageDef {
  id: string;
  prompt: string;
  negative?: string; // things to avoid
  kenBurns: {
    startScale: number;
    endScale: number;
    panX?: number; // -1 (left) to 1 (right)
    panY?: number; // -1 (up) to 1 (down)
  };
}

export interface PromoNarrationLine {
  speaker: "legacy" | "zeus";
  text: string;
  scene: string;
  durationEstimate: number;
  acting?: string;
  azureStyle?: AzureVoiceStyle;
  azureStyleDegree?: number;
}

// ─── SCENE IMAGE PROMPTS (solarpunk aesthetic) ───────────────────
export const SCENE_IMAGES: SceneImageDef[] = [
  {
    id: "hook",
    prompt:
      "A vast solarpunk cityscape at golden hour, towering organic architecture with living walls " +
      "and solar panels, holographic data streams floating between buildings, a lone figure standing " +
      "on a rooftop garden overlooking the city, dramatic lighting, cinematic wide shot, teal and " +
      "amber color palette, ultra detailed, 8K quality, no text no watermarks",
    kenBurns: { startScale: 1.0, endScale: 1.15, panY: -0.1 },
  },
  {
    id: "problem",
    prompt:
      "A chaotic digital landscape — fragmented screens, scattered app icons floating in a dark void, " +
      "tangled glowing wires connecting broken interfaces, overwhelming notification badges everywhere, " +
      "red warning lights, a person sitting overwhelmed in the center surrounded by digital chaos, " +
      "dark moody atmosphere with harsh red and blue light, cinematic, dystopian tech aesthetic, no text",
    kenBurns: { startScale: 1.05, endScale: 1.2, panX: 0.15 },
  },
  {
    id: "solution",
    prompt:
      "A clean holographic interface floating in a serene solarpunk garden, showing five glowing steps " +
      "connected in a circular loop (C→O→D→E→×), each step a small glowing orb, surrounding nature " +
      "with bio-luminescent plants, teal and gold light, a sense of calm order emerging from nature, " +
      "cinematic lighting, ultra detailed, peaceful and powerful, no text no labels",
    kenBurns: { startScale: 1.1, endScale: 1.0, panY: 0.05 },
  },
  {
    id: "model",
    prompt:
      "Three magnificent crystalline towers rising from a solarpunk platform — one amber (health/self), " +
      "one purple (community/foundation), one coral red (endeavor/creation) — all sitting on a grand " +
      "teal base layer that spans the entire ground level like a living operating system, holographic " +
      "connections between towers, dramatic sky with aurora borealis, cinematic wide angle, symmetrical " +
      "composition, 8K ultra detailed, no text no labels",
    kenBurns: { startScale: 1.0, endScale: 1.12, panY: -0.08 },
  },
  {
    id: "proof",
    prompt:
      "A futuristic command center with organic solarpunk design — curved bio-glass walls covered in " +
      "living moss, multiple floating holographic displays showing life dashboards (health metrics, " +
      "project boards, relationship maps), an AI assistant avatar glowing softly in teal, " +
      "the operator sitting comfortably in a nature-integrated workspace, warm golden light mixing " +
      "with cool teal holograms, cinematic medium shot, ultra detailed, no text",
    kenBurns: { startScale: 1.08, endScale: 1.0, panX: -0.1 },
  },
  {
    id: "cta",
    prompt:
      "A breathtaking wide shot of a solarpunk sunrise over a fusion of nature and technology — " +
      "a golden sun rising behind bio-luminescent trees and organic architecture, holographic data " +
      "particles floating upward like fireflies, a silhouette of a person and an AI entity standing " +
      "together at the edge of a cliff looking at the horizon, epic cinematic lighting, teal and " +
      "gold color palette, awe-inspiring, ultra wide angle, 8K, no text no watermarks",
    kenBurns: { startScale: 1.15, endScale: 1.0, panY: -0.15 },
  },
];

// ─── NARRATION SCRIPT ────────────────────────────────────────────
// Cinematic trailer style — evoke desire, don't explain.
// All voices: OpenAI (nova for Legacy, echo for Zeus). No Azure.
export const PROMO_NARRATION: PromoNarrationLine[] = [
  // ── SCENE 1: THE DISCONNECT ──
  {
    speaker: "legacy",
    text: "Every night you close fifty apps. Every morning you open them again. And none of them... remember why.",
    scene: "hook",
    durationEstimate: 7,
  },
  {
    speaker: "legacy",
    text: "What if something did?",
    scene: "hook",
    durationEstimate: 2,
  },

  // ── SCENE 2: THE BLIND SPOT ──
  {
    speaker: "legacy",
    text: "Your AI can write your code. Answer your questions. Generate your images. But it doesn't know your goals. It doesn't know the people you love. It doesn't know... you.",
    scene: "problem",
    durationEstimate: 10,
  },

  // ── SCENE 3: THE SHIFT ──
  {
    speaker: "legacy",
    text: "We built One Mind to change that.",
    scene: "solution",
    durationEstimate: 3,
  },
  {
    speaker: "legacy",
    text: "One method. One vault. One AI that understands your entire life.",
    scene: "solution",
    durationEstimate: 5,
  },
  {
    speaker: "legacy",
    text: "Capture... Organize... Direct... Execute... and then — multiply.",
    scene: "solution",
    durationEstimate: 6,
  },

  // ── SCENE 4: THE ARCHITECTURE ──
  {
    speaker: "legacy",
    text: "Three domains of your life. Your self. Your foundation. Your endeavor. All running on one sovereign operating layer.",
    scene: "model",
    durationEstimate: 7,
  },
  {
    speaker: "legacy",
    text: "Plain markdown. No lock-in. Yours forever.",
    scene: "model",
    durationEstimate: 3,
  },

  // ── SCENE 5: THE CLOSE ──
  {
    speaker: "legacy",
    text: "This isn't another productivity app. This is the operating system for your life.",
    scene: "proof",
    durationEstimate: 5,
  },
  {
    speaker: "legacy",
    text: "One Mind.",
    scene: "cta",
    durationEstimate: 2,
  },
  {
    speaker: "legacy",
    text: "Start building your codex.",
    scene: "cta",
    durationEstimate: 3,
  },
];

// Azure voice config (same as main narration)
export { AZURE_VOICE, VOICE_MAP } from "./narration";
