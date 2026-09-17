# AI Intelligence Feed Stack

> *25 RSS feeds curated for One Mind OS intelligence briefings.*
> *Import the OPML file into Miniflux (or any RSS reader).*

---

## Setup

1. Install Miniflux on VPS: `docker run -d -p 8080:8080 miniflux/miniflux:latest`
2. Open Miniflux → Feeds → Import → paste the OPML below
3. Configure Trinity to pull from Miniflux API for intelligence briefings

Miniflux API for Trinity: `GET /v1/entries?status=unread&limit=50`

---

## OPML — Import This

```xml
<?xml version="1.0" encoding="UTF-8"?>
<opml version="2.0">
  <head>
    <title>One Mind OS Intelligence Feed Stack</title>
  </head>
  <body>

    <!-- AI & Agents -->
    <outline text="AI & Agents" title="AI & Agents">
      <outline type="rss" text="Anthropic News" title="Anthropic News" xmlUrl="https://www.anthropic.com/rss.xml" htmlUrl="https://www.anthropic.com"/>
      <outline type="rss" text="OpenAI Blog" title="OpenAI Blog" xmlUrl="https://openai.com/blog/rss" htmlUrl="https://openai.com/blog"/>
      <outline type="rss" text="Google DeepMind" title="Google DeepMind" xmlUrl="https://deepmind.google/blog/rss.xml" htmlUrl="https://deepmind.google"/>
      <outline type="rss" text="Hugging Face Blog" title="Hugging Face Blog" xmlUrl="https://huggingface.co/blog/feed.xml" htmlUrl="https://huggingface.co/blog"/>
      <outline type="rss" text="Simon Willison" title="Simon Willison" xmlUrl="https://simonwillison.net/atom/everything/" htmlUrl="https://simonwillison.net"/>
      <outline type="rss" text="Ethan Mollick (One Useful Thing)" title="Ethan Mollick" xmlUrl="https://www.oneusefulthing.org/feed" htmlUrl="https://www.oneusefulthing.org"/>
    </outline>

    <!-- Security & Infrastructure -->
    <outline text="Security & Infrastructure" title="Security & Infrastructure">
      <outline type="rss" text="Krebs on Security" title="Krebs on Security" xmlUrl="https://krebsonsecurity.com/feed/" htmlUrl="https://krebsonsecurity.com"/>
      <outline type="rss" text="Schneier on Security" title="Schneier on Security" xmlUrl="https://www.schneier.com/feed/atom/" htmlUrl="https://www.schneier.com"/>
      <outline type="rss" text="SANS Internet Stormcast" title="SANS ISC" xmlUrl="https://isc.sans.edu/rssfeed.xml" htmlUrl="https://isc.sans.edu"/>
      <outline type="rss" text="Hacker News Front Page" title="Hacker News" xmlUrl="https://hnrss.org/frontpage" htmlUrl="https://news.ycombinator.com"/>
      <outline type="rss" text="The Register" title="The Register" xmlUrl="https://www.theregister.com/headlines.atom" htmlUrl="https://www.theregister.com"/>
    </outline>

    <!-- Business & Empire -->
    <outline text="Business & Empire" title="Business & Empire">
      <outline type="rss" text="Paul Graham Essays" title="Paul Graham" xmlUrl="https://paulgraham.com/rss.html" htmlUrl="https://paulgraham.com"/>
      <outline type="rss" text="Naval Ravikant Blog" title="Naval" xmlUrl="https://nav.al/feed" htmlUrl="https://nav.al"/>
      <outline type="rss" text="Stratechery" title="Stratechery" xmlUrl="https://stratechery.com/feed/" htmlUrl="https://stratechery.com"/>
      <outline type="rss" text="The Generalist" title="The Generalist" xmlUrl="https://www.readthegeneralist.com/feed" htmlUrl="https://www.readthegeneralist.com"/>
      <outline type="rss" text="Lenny's Newsletter" title="Lenny's Newsletter" xmlUrl="https://www.lennysnewsletter.com/feed" htmlUrl="https://www.lennysnewsletter.com"/>
    </outline>

    <!-- Personal Development & Performance -->
    <outline text="Personal Development" title="Personal Development">
      <outline type="rss" text="James Clear" title="James Clear" xmlUrl="https://jamesclear.com/feed" htmlUrl="https://jamesclear.com"/>
      <outline type="rss" text="Derek Sivers" title="Derek Sivers" xmlUrl="https://sive.rs/en.atom" htmlUrl="https://sive.rs"/>
      <outline type="rss" text="Scott H Young" title="Scott H Young" xmlUrl="https://www.scotthyoung.com/blog/feed/" htmlUrl="https://www.scotthyoung.com"/>
      <outline type="rss" text="Andrew Huberman Lab" title="Huberman Lab" xmlUrl="https://feeds.buzzsprout.com/1800277.rss" htmlUrl="https://hubermanlab.com"/>
    </outline>

    <!-- Finance & Real Estate -->
    <outline text="Finance & Real Estate" title="Finance & Real Estate">
      <outline type="rss" text="Bigger Pockets Blog" title="BiggerPockets" xmlUrl="https://www.biggerpockets.com/blog/feed" htmlUrl="https://www.biggerpockets.com"/>
      <outline type="rss" text="A Wealth of Common Sense" title="A Wealth of Common Sense" xmlUrl="https://awealthofcommonsense.com/feed/" htmlUrl="https://awealthofcommonsense.com"/>
      <outline type="rss" text="Collaborative Fund Blog" title="Collaborative Fund" xmlUrl="https://www.collaborativefund.com/blog/feed.xml" htmlUrl="https://www.collaborativefund.com"/>
    </outline>

    <!-- Emerging Tech & Research -->
    <outline text="Research & Emerging Tech" title="Research & Emerging Tech">
      <outline type="rss" text="MIT Technology Review" title="MIT Technology Review" xmlUrl="https://www.technologyreview.com/feed/" htmlUrl="https://www.technologyreview.com"/>
      <outline type="rss" text="Ars Technica" title="Ars Technica" xmlUrl="https://feeds.arstechnica.com/arstechnica/index/" htmlUrl="https://arstechnica.com"/>
    </outline>

  </body>
</opml>
```

---

## Feed Categories

| Category | Count | Primary Use |
|----------|-------|-------------|
| AI & Agents | 6 | Agent stack updates, model releases |
| Security & Infrastructure | 5 | Threat intelligence, CVEs, VPS news |
| Business & Empire | 5 | Strategy, venture insights |
| Personal Development | 4 | HP inputs, performance research |
| Finance & Real Estate | 3 | GE investment intelligence |
| Research & Emerging Tech | 2 | Trinity research inputs |
| **Total** | **25** | |

---

## Intelligence Brief Integration

Ask Trinity weekly:
> *"Trinity, pull the last 7 days of unread RSS items from Miniflux. Surface the 5 most relevant items for my current goals: [goals]. Format as: title + source + 2-sentence summary + why it matters to One Mind OS."*

---

## Adding to the Stack

Before adding a new feed, ask:
1. Does it signal something I act on, or just inform?
2. Is the signal-to-noise ratio worth it? (Test for 2 weeks)
3. Does it duplicate an existing feed?

Prune quarterly: remove any feed you haven't clicked in 3 months.

---

*Part of the One Mind OS Intelligence & Research — Course 8*
