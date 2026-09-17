# Episode Registry & npm Scripts — Video Production

## Episode Registry

| ID                 | Title                                   | Duration | Format | Status      |
| ------------------ | --------------------------------------- | -------- | ------ | ----------- |
| `OneMindIntroV2`   | One Mind — Zeus + Legacy Intro          | 5:04     | V2     | **Rendered** |
| `OneMindIntro`     | One Mind — You + AI = One Life          | 3:39     | V1     | Rendered (superseded by V2) |
| `OpenClawDreaming` | Your AI Learns While You Sleep          | 3:20     | V1     | Rendered     |
| `OneMindPromoV4`   | Promo (AI video clips)                  | 1:00     | —      | Rendered (rejected) |
| `OneMindPromoV3`   | Promo (motion graphics)                 | 1:27     | —      | Rendered     |
| `OneMindPromo`     | Promo (original)                        | 2:11     | —      | Rendered     |
| `ComposioCLI`      | Composio CLI promo                      | 1:00     | —      | Rendered     |

---

## npm Scripts Reference

```bash
# ── V2 Pipeline (current) ──
npm run voice:intro-v2     # Generate two-voice clips for OneMindIntroV2
npm run render:intro-v2    # Render OneMindIntroV2 to out/onemind-intro-v2.mp4

# ── V1 Pipeline ──
npm run voice:intro        # Generate Legacy-only clips for OneMindIntro
npm run render:intro       # Render OneMindIntro
npm run voice:dreaming     # Generate clips for OpenClaw Dreaming
npm run render:dreaming    # Render OpenClaw Dreaming

# ── General ──
npm run studio             # Preview at localhost:3000
npm run render:test        # Render TestVideo

# ── Content Generation ──
npm run script             # AI script generator (GPT-4o)
npm run metadata           # SEO title/description/tags
npm run captions           # Multi-platform social captions
```
