---
id: so-01-20260412000012
title: "One Mind Image Style — AJCP Profile"
type: reference
status: active
created: 2026-04-12
owner: zeus
tags: [path/so, path/so/01, type/reference, prompt-engineering, json, image-generation, ajcp]
ai_summary: "Production AJCP profile for One Mind brand visuals — paste into Midjourney, DALL-E, Flux, or Firefly for consistent solarpunk cinematic style"
---
# One Mind Image Style — AJCP Profile

> **How to use:** Paste the JSON into your image generation prompt, or extract the values into native prompt syntax for Midjourney/DALL-E/Flux.
>
> **Midjourney shortcut:** Paste the `aesthetic_keywords` + `mood` as a comma-separated string, append `--ar 16:9 --style raw`

---

## Profile

```json
{
  "profile_type": "image_generation_ajcp",
  "brand": "One Mind",
  "style_name": "Solarpunk Sovereign",
  "style_description": "Cinematic solarpunk realism — technology and nature in harmony, sovereign human presence, structured and intentional",
  "lighting": {
    "primary": "dramatic side-lighting or golden hour rim light",
    "secondary": "deep shadows with selective highlights",
    "tech_contexts": "cool neon accent (cyan or amber) on dark backgrounds",
    "avoid": "flat even lighting, studio white, harsh overhead"
  },
  "composition": [
    "rule of thirds",
    "subject slightly off-center",
    "negative space on one side for text overlay",
    "leading lines toward subject",
    "depth layers: foreground element + subject + background"
  ],
  "color_palette": {
    "backgrounds": ["#0a0a0f", "#0f1923", "#1a1a2e"],
    "primary_accents": ["#4ade80", "#22d3ee"],
    "secondary_accents": ["#f59e0b", "#a78bfa"],
    "skin_tones": "natural, warm undertones",
    "avoid": ["pure white backgrounds", "flat corporate blues", "stock photo beige", "oversaturated greens"]
  },
  "aesthetic_keywords": [
    "solarpunk", "cinematic", "AI-native", "sovereign human",
    "clean technology", "structured clarity", "intentional living",
    "human-machine harmony", "architectural minimalism"
  ],
  "mood": [
    "calm intensity",
    "focused",
    "sovereign",
    "forward-moving",
    "quiet confidence"
  ],
  "subject_direction": {
    "human_subjects": "confident posture, direct or slightly off-camera gaze, purposeful expression",
    "environments": "clean workspaces, natural light + screens, outdoor-indoor blend",
    "tech_elements": "integrated naturally — not foregrounded, not hidden"
  },
  "negative_prompts": [
    "cartoon", "anime", "flat design", "stock photo", "corporate clip art",
    "blurry", "oversaturated", "cluttered backgrounds", "dystopian",
    "robotic dehumanization", "generic motivational poster"
  ],
  "text_overlay_space": true,
  "aspect_ratios": {
    "social_feed": "1:1 or 4:5",
    "cover_banner": "16:9",
    "ebook_cover": "6:9",
    "story": "9:16"
  },
  "use_cases": [
    "One Mind social media posts",
    "Manifesto and ebook covers",
    "Community banner and header",
    "Agent visualization graphics",
    "Zeus personal brand content"
  ]
}
```

---

## Midjourney Quick Prompt

Extract and paste directly:

```
solarpunk sovereign, cinematic realism, dramatic side-lighting, deep shadows, cyan neon accent, clean technology, human-machine harmony, calm intensity, focused, architectural minimalism, rule of thirds, negative space for text, dark background #0a0a0f --ar 16:9 --style raw --no cartoon, flat design, stock photo, blurry, dystopian
```


---

## Related

- [[00-24 SO (Sovereign Operations)/01 Intelligence (AI, Memory & Tools) (The Brain)/01 Intelligence|01 Intelligence]]
