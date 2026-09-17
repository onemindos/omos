# LLM Vision

> **Multimodal vision capabilities for AI agents**

## Overview

LLM Vision refers to multimodal large language models that can process and understand images alongside text. Enables Legacy AI to "see" and analyze visual content from cameras, screenshots, documents, and more.

## Status

| Attribute   | Value                  |
| ----------- | ---------------------- |
| Status      | ✅ Active              |
| Category    | AI Agents              |
| Type        | Multimodal AI          |
| Integration | Claude, GPT-4V, Gemini |

## Use Cases in OneMind

1. **Camera Analysis** - Describe what's in security footage
2. **Document Understanding** - Extract info from images/PDFs
3. **Screenshot Analysis** - Understand UI, charts, diagrams
4. **Object Recognition** - Identify items, people, scenes
5. **Visual Q&A** - Answer questions about images

## Supported Models

| Model             | Provider    | Strengths                    |
| ----------------- | ----------- | ---------------------------- |
| Claude 3.5 Sonnet | Anthropic   | Detailed analysis, reasoning |
| Claude 3 Opus     | Anthropic   | Best overall quality         |
| GPT-4 Vision      | OpenAI      | Fast, good general vision    |
| Gemini Pro Vision | Google      | Multi-image, video frames    |
| LLaVA             | Open Source | Self-hosted option           |

## API Integration

### Claude Vision (Primary)

```python
import anthropic
import base64
import httpx

client = anthropic.Anthropic()

async def analyze_image_url(image_url: str, prompt: str) -> str:
    """Analyze image from URL using Claude Vision."""
    # Fetch image
    async with httpx.AsyncClient() as http:
        response = await http.get(image_url)
        image_data = base64.standard_b64encode(response.content).decode()

    media_type = "image/jpeg"
    if image_url.endswith(".png"):
        media_type = "image/png"

    message = client.messages.create(
        model="claude-sonnet-4-20250514",
        max_tokens=1024,
        messages=[
            {
                "role": "user",
                "content": [
                    {
                        "type": "image",
                        "source": {
                            "type": "base64",
                            "media_type": media_type,
                            "data": image_data
                        }
                    },
                    {
                        "type": "text",
                        "text": prompt
                    }
                ]
            }
        ]
    )

    return message.content[0].text

async def analyze_image_file(image_path: str, prompt: str) -> str:
    """Analyze local image file."""
    with open(image_path, "rb") as f:
        image_data = base64.standard_b64encode(f.read()).decode()

    message = client.messages.create(
        model="claude-sonnet-4-20250514",
        max_tokens=1024,
        messages=[
            {
                "role": "user",
                "content": [
                    {
                        "type": "image",
                        "source": {
                            "type": "base64",
                            "media_type": "image/jpeg",
                            "data": image_data
                        }
                    },
                    {
                        "type": "text",
                        "text": prompt
                    }
                ]
            }
        ]
    )

    return message.content[0].text
```

### OpenAI GPT-4 Vision

```python
from openai import OpenAI

client = OpenAI()

async def analyze_with_gpt4v(image_url: str, prompt: str) -> str:
    """Analyze image using GPT-4 Vision."""
    response = client.chat.completions.create(
        model="gpt-4-vision-preview",
        messages=[
            {
                "role": "user",
                "content": [
                    {"type": "text", "text": prompt},
                    {
                        "type": "image_url",
                        "image_url": {"url": image_url}
                    }
                ]
            }
        ],
        max_tokens=1024
    )

    return response.choices[0].message.content
```

### Multiple Images

```python
async def compare_images(images: list[str], prompt: str) -> str:
    """Analyze multiple images together."""
    content = []

    for image_url in images:
        async with httpx.AsyncClient() as http:
            response = await http.get(image_url)
            image_data = base64.standard_b64encode(response.content).decode()

        content.append({
            "type": "image",
            "source": {
                "type": "base64",
                "media_type": "image/jpeg",
                "data": image_data
            }
        })

    content.append({"type": "text", "text": prompt})

    message = client.messages.create(
        model="claude-sonnet-4-20250514",
        max_tokens=2048,
        messages=[{"role": "user", "content": content}]
    )

    return message.content[0].text
```

## Legacy AI Integration

```python
from agno import Tool

@Tool
async def analyze_image(
    image_source: str,
    question: str = "Describe what you see in detail."
) -> str:
    """Analyze an image and answer questions about it.

    Args:
        image_source: URL or file path to image
        question: What to analyze or ask about the image
    """
    if image_source.startswith("http"):
        return await analyze_image_url(image_source, question)
    else:
        return await analyze_image_file(image_source, question)

@Tool
async def analyze_camera_snapshot(
    camera: str,
    question: str = "Describe what's happening in this camera view."
) -> str:
    """Get and analyze current camera snapshot.

    Args:
        camera: Camera name (front_door, backyard, etc.)
        question: What to look for or describe
    """
    snapshot_url = f"{FRIGATE_URL}/api/{camera}/latest.jpg"
    return await analyze_image_url(snapshot_url, question)

@Tool
async def read_document_image(image_path: str) -> str:
    """Extract and read text from a document image.

    Args:
        image_path: Path to document image or screenshot
    """
    prompt = """Extract all text from this document image.
    Preserve the structure and formatting as much as possible.
    If there are tables, represent them clearly."""

    return await analyze_image_file(image_path, prompt)

@Tool
async def identify_person(camera: str) -> dict:
    """Try to identify who is at a camera location.

    Args:
        camera: Camera name
    """
    snapshot_url = f"{FRIGATE_URL}/api/{camera}/latest.jpg"

    analysis = await analyze_image_url(
        snapshot_url,
        """Describe any people in this image:
        - How many people
        - General appearance (clothing, build)
        - What they appear to be doing
        - Any packages or items they're carrying"""
    )

    return {
        "camera": camera,
        "analysis": analysis,
        "snapshot_url": snapshot_url
    }
```

## Video Frame Analysis

```python
import cv2
import tempfile

async def analyze_video_clip(video_path: str, prompt: str, frame_count: int = 4) -> str:
    """Extract frames from video and analyze."""
    cap = cv2.VideoCapture(video_path)
    total_frames = int(cap.get(cv2.CAP_PROP_FRAME_COUNT))

    frames = []
    for i in range(frame_count):
        frame_pos = int(total_frames * i / frame_count)
        cap.set(cv2.CAP_PROP_POS_FRAMES, frame_pos)
        ret, frame = cap.read()
        if ret:
            # Save frame temporarily
            with tempfile.NamedTemporaryFile(suffix=".jpg", delete=False) as f:
                cv2.imwrite(f.name, frame)
                frames.append(f.name)

    cap.release()

    # Analyze all frames together
    return await compare_images(frames, prompt)
```

## Environment Variables

```bash
ANTHROPIC_API_KEY=your_key
OPENAI_API_KEY=your_key
GOOGLE_API_KEY=your_key

# Default vision model
VISION_MODEL=claude-sonnet-4-20250514
```

## Best Practices

1. **Optimize image size** - Resize large images (max 1568px)
2. **Be specific in prompts** - Guide what to look for
3. **Use appropriate model** - Claude for detail, GPT-4V for speed
4. **Handle errors** - Images may fail to load
5. **Cache results** - Vision calls are expensive

## Use Case Examples

```python
# Security check
await analyze_camera_snapshot("front_door", "Is there a package on the porch?")

# Document reading
await read_document_image("/path/to/receipt.jpg")

# Chart understanding
await analyze_image(chart_url, "What trends does this chart show?")

# Scene description
await analyze_image(photo_url, "Describe the setting and mood of this photo")
```

## Resources

- [Claude Vision Docs](https://docs.anthropic.com/claude/docs/vision)
- [OpenAI Vision Guide](https://platform.openai.com/docs/guides/vision)
- [Gemini Vision](https://ai.google.dev/gemini-api/docs/vision)

---

Last Updated: January 25, 2026


---

## Related

- [[00-24 SO (Sovereign Operations)/02 Infrastructure (VPS, Servers & Networking) (The Body)/02 Infrastructure|02 Infrastructure]]
