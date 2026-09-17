# Seedance

ByteDance's video generation models — compete with Sora and Kling for video gen quality.

| Model ID | Resolution | Notes |
|----------|-----------|-------|
| `seedance-1-0-lite` | Up to 1080p | Fast video gen, lower compute |
| `seedance-1-0-pro` | Up to 1080p | Flagship video quality |

**Best for:** Brand content (GE/75), marketing videos, course demos. Seedance edges Kling on speed; Kling edges on some quality metrics.

**API:** Separate REST endpoint from the LLM API — returns a `task_id`, poll for completion.

**Stack use case:** GE/Empire agent — video content for brand/marketing pipelines.


---

## Related

- [[00-24 SO (Sovereign Operations)/01 Intelligence (AI, Memory & Tools) (The Brain)/01 Intelligence|01 Intelligence]]
