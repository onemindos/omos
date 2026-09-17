# Production Notes — Video Production

## Zeus's Real-Life Examples (for reference in scripts)

Use these REAL numbers and examples to make content concrete, not abstract:

| What           | Zeus's Actual Reality                                                           |
| -------------- | ------------------------------------------------------------------------------- |
| Vault size     | ~78 GB, ~111,000 files                                                          |
| Domains        | 100 numbered domains (00-99) across 4 groups                                    |
| Domain groups  | SO (00-24), HP (25-49), LE (50-74), GE (75-99)                                  |
| AI agents      | 4 primary: Sentinel (SO), Coach (HP), Haven (LE), Empire (GE) + OpenClaw agents |
| Main AI        | Legacy — runs on OpenClaw on a DigitalOcean VPS                                 |
| Sync           | Obsidian on Mac + iPhone + VPS, git for version control                         |
| Tools          | Obsidian, VS Code, CapCut, 1Password, Todoist, git                              |
| Deduplication  | 175,214 duplicate files removed, 94.5 GB reclaimed                              |
| Git cleanup    | .git went from 46.2 GB → 273 MB                                                 |
| Daily workflow | Capture → Organize → Direct → Execute → × Multiply                              |
| Public site    | codex.onemindos.com (Quartz on Cloudflare Pages)                                |

---

## Production Cost Per Episode

| Resource                          | V1 Cost      | V2 Cost      |
| --------------------------------- | ------------ | ------------ |
| Voice generation (11 segments)    | ~$0.09       | ~$0.13       |
| Script generation (GPT-4o)        | ~$0.02       | ~$0.02       |
| Metadata generation               | ~$0.01       | ~$0.01       |
| SFX generation (ffmpeg)           | Free         | Free         |
| Ambient music (ffmpeg synth)      | Free         | Free         |
| Remotion rendering                | Free (local) | Free (local) |
| CapCut editing (optional)         | Free tier    | Free tier    |
| **Total per episode**             | **~$0.12**   | **~$0.16**   |

---

## Known Issues & Future Improvements

- **Ambient music is synthetic** — `ambient-bg.mp3` is an ffmpeg-generated synth pad. Replace with Suno AI or royalty-free lo-fi for production.
- **KineticText unused** — Built but not yet used in any composition. Good for key phrase emphasis.
- **ImageSlide component missing** — No screenshot/photo slide component yet. Could add for real vault screenshots.
- **Speaker timing is estimated** — `activeSpeaker` frame numbers in compositions are based on cumulative line durations. May need manual adjustment after watching.
- **No face-cam integration yet** — V2 is pure Remotion. Future: intercut Zeus's face-cam footage in CapCut.
