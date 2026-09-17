# FLUX (Image Generation)

FLUX by Black Forest Labs — current best-in-class image generation models, hosted on fal.ai.

| Model ID | Speed | Quality | Price | Notes |
|----------|-------|---------|-------|-------|
| `fal-ai/flux/schnell` | ~1s | Good | ~$0.003/img | Fastest — real-time use |
| `fal-ai/flux/dev` | ~3s | Excellent | ~$0.025/img | High quality, no watermark |
| `fal-ai/flux-pro` | ~8s | Best | ~$0.055/img | Flagship quality |
| `fal-ai/flux-pro/v1.1-ultra` | ~6s | Flagship | ~$0.06/img | Highest resolution output |
| `fal-ai/flux-realism` | ~4s | Photorealistic | ~$0.03/img | Best for realistic portraits |
| `fal-ai/flux-lora` | ~3s | Custom | ~$0.025/img | Run custom LoRA fine-tunes |

## FLUX vs Midjourney vs DALL-E

FLUX-Pro matches or exceeds Midjourney quality on photorealistic and artistic subjects. DALL-E 3 (via OpenAI) is slower and more restricted. FLUX has no content policy restrictions for most use cases.

## API Usage

```python
import fal_client
result = fal_client.subscribe("fal-ai/flux/dev",
    arguments={
        "prompt": "professional headshot, Zeus, confident, dark background",
        "image_size": "square_hd",
        "num_images": 1
    }
)
print(result["images"][0]["url"])
```

## Stack Use Cases

- `flux/schnell`: Real-time UI image generation, bulk content
- `flux-pro`: Final marketing assets, brand content
- `flux-realism`: Headshots, product photography, realistic scenes
- `flux-lora`: Consistent character/brand identity across images


---

## Related

- [[00-24 SO (Sovereign Operations)/01 Intelligence (AI, Memory & Tools) (The Brain)/01 Intelligence|01 Intelligence]]
