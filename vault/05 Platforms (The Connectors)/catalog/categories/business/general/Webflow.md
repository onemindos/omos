# Webflow

> **Visual Website Builder**

## Overview

Webflow is a visual website builder for creating the OneMind marketing site, landing pages, and blog without code.

## Status

- **Active**: ✅ Yes
- **Plan**: CMS
- **Last Updated**: January 25, 2026

## Use Cases

1. **Marketing Site** - Main onemind.io website
2. **Landing Pages** - Campaign pages
3. **Blog** - Content marketing
4. **Documentation** - Public docs (optional)
5. **SEO** - Search optimization

## Key Features

- Visual design (no code)
- CMS for dynamic content
- Built-in hosting
- Custom animations
- Form handling
- SEO controls
- Responsive design

## Sites

| Site    | URL             | Purpose              |
| ------- | --------------- | -------------------- |
| OneMind | onemind.io      | Main marketing       |
| Blog    | onemind.io/blog | Content              |
| Docs    | docs.onemind.io | Public documentation |

## CMS Structure

### Blog Posts

```
Collection: Blog Posts
Fields:
- Name (Plain Text) - Title
- Slug (Slug)
- Thumbnail (Image)
- Summary (Plain Text)
- Body (Rich Text)
- Author (Reference → Authors)
- Category (Reference → Categories)
- Published Date (Date)
- Featured (Switch)
```

### Authors

```
Collection: Authors
Fields:
- Name (Plain Text)
- Slug (Slug)
- Photo (Image)
- Bio (Plain Text)
- Twitter (Link)
- LinkedIn (Link)
```

### Categories

```
Collection: Categories
Fields:
- Name (Plain Text)
- Slug (Slug)
- Description (Plain Text)
```

## API Integration

### Configuration

```python
import requests

WEBFLOW_API_KEY = os.getenv("WEBFLOW_API_KEY")
WEBFLOW_SITE_ID = os.getenv("WEBFLOW_SITE_ID")
WEBFLOW_COLLECTION_ID = os.getenv("WEBFLOW_BLOG_COLLECTION_ID")

headers = {
    "Authorization": f"Bearer {WEBFLOW_API_KEY}",
    "accept-version": "1.0.0",
    "Content-Type": "application/json"
}

base_url = "https://api.webflow.com"
```

### Get Collection Items

```python
def get_blog_posts(limit: int = 100):
    response = requests.get(
        f"{base_url}/collections/{WEBFLOW_COLLECTION_ID}/items",
        headers=headers,
        params={"limit": limit}
    )
    return response.json()
```

### Create CMS Item

```python
def create_blog_post(title: str, slug: str, body: str, summary: str):
    response = requests.post(
        f"{base_url}/collections/{WEBFLOW_COLLECTION_ID}/items",
        headers=headers,
        json={
            "fields": {
                "name": title,
                "slug": slug,
                "body": body,
                "summary": summary,
                "_archived": False,
                "_draft": True
            }
        }
    )
    return response.json()
```

### Publish Items

```python
def publish_items(item_ids: list):
    response = requests.post(
        f"{base_url}/sites/{WEBFLOW_SITE_ID}/publish",
        headers=headers,
        json={
            "domains": ["onemind.io", "www.onemind.io"]
        }
    )
    return response.json()
```

### Update Item

```python
def update_blog_post(item_id: str, fields: dict):
    response = requests.patch(
        f"{base_url}/collections/{WEBFLOW_COLLECTION_ID}/items/{item_id}",
        headers=headers,
        json={"fields": fields}
    )
    return response.json()
```

## Forms

### Form Submissions

```python
def get_form_submissions(form_id: str):
    response = requests.get(
        f"{base_url}/sites/{WEBFLOW_SITE_ID}/forms/{form_id}/submissions",
        headers=headers
    )
    return response.json()
```

### Form Webhook

Configure in Webflow:

1. Go to Project Settings → Integrations
2. Add webhook URL for form submissions

```python
@app.post("/webhooks/webflow/form")
async def webflow_form_webhook(request: Request):
    data = await request.json()

    form_name = data.get("name")
    form_data = data.get("data", {})

    if form_name == "Contact Form":
        await handle_contact_form(form_data)
    elif form_name == "Newsletter":
        await handle_newsletter_signup(form_data)

    return {"status": "ok"}
```

## Custom Code

### Head Code

```html
<!-- Google Analytics -->
<script
  async
  src="https://www.googletagmanager.com/gtag/js?id=G-XXXXXX"
></script>
<script>
  window.dataLayer = window.dataLayer || [];
  function gtag() {
    dataLayer.push(arguments);
  }
  gtag("js", new Date());
  gtag("config", "G-XXXXXX");
</script>

<!-- Meta Pixel -->
<script>
  !function(f,b,e,v,n,t,s)
  ...
</script>
```

### Footer Code

```html
<!-- Intercom -->
<script>
  window.intercomSettings = {
    api_base: "https://api-iam.intercom.io",
    app_id: "xxxxx",
  };
</script>
```

## SEO Configuration

### Per-Page Settings

- **Title Tag**: `Feature Name | OneMind OS`
- **Meta Description**: 150-160 characters
- **Open Graph Image**: 1200x630px
- **Canonical URL**: Set for duplicate content

### Sitemap

Auto-generated at `/sitemap.xml`

### Robots.txt

```
User-agent: *
Allow: /
Disallow: /admin/
Sitemap: https://onemind.io/sitemap.xml
```

## Integrations

### Connected Services

| Service                   | Purpose          |
| ------------------------- | ---------------- |
| **Google Analytics**      | Traffic tracking |
| **Google Search Console** | SEO monitoring   |
| **Mailchimp**             | Newsletter       |
| **Zapier**                | Form automation  |
| **Intercom**              | Chat widget      |

## Hosting

### Custom Domain

1. Add domain in Webflow
2. Update DNS at Cloudflare
3. SSL auto-provisioned

### DNS Records

```
Type: A
Name: @
Value: 75.2.70.75

Type: CNAME
Name: www
Value: proxy-ssl.webflow.com
```

## Legacy AI Integration

### Blog Publishing

```python
@Tool
def publish_blog_post(title: str, content: str, summary: str) -> str:
    """Publish a blog post to the OneMind website."""

    slug = generate_slug(title)

    # Create draft
    post = create_blog_post(
        title=title,
        slug=slug,
        body=content,
        summary=summary
    )

    return f"Draft created: https://onemind.io/blog/{slug}"
```

### Get Recent Posts

```python
@Tool
def get_recent_blog_posts(limit: int = 5) -> list:
    """Get recent blog posts from the website."""
    posts = get_blog_posts(limit)
    return [
        {
            "title": p["name"],
            "url": f"https://onemind.io/blog/{p['slug']}",
            "published": p.get("published-date")
        }
        for p in posts.get("items", [])
    ]
```

## Style Guide

### Colors

| Name      | Hex     | Usage       |
| --------- | ------- | ----------- |
| Primary   | #6366F1 | CTAs, links |
| Secondary | #06B6D4 | Accents     |
| Dark      | #0F172A | Text        |
| Light     | #F8FAFC | Backgrounds |

### Typography

- **Headings**: Inter, Bold
- **Body**: Inter, Regular
- **Code**: JetBrains Mono

## Best Practices

- ✅ Use CMS for repeating content
- ✅ Optimize images (WebP)
- ✅ Set proper SEO on all pages
- ✅ Use interactions sparingly
- ✅ Test on mobile devices
- ✅ Keep custom code minimal
- ✅ Back up before major changes

## Resources

- [Webflow University](https://university.webflow.com/)
- [Webflow API Docs](https://developers.webflow.com/)
- [Webflow Showcase](https://webflow.com/made-in-webflow)

---

**Related**: [Cloudflare](./Cloudflare.md) | [Circle](Circle.md)


---

## Related

- [[00-24 SO (Sovereign Operations)/02 Infrastructure (VPS, Servers & Networking) (The Body)/02 Infrastructure|02 Infrastructure]]
