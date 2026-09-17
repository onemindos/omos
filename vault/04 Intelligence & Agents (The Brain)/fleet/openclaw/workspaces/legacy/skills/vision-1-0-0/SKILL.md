# Vision Skill

AI-powered camera vision combining Home Assistant LLM Vision (GPT-4o / Claude / Gemini) and Frigate NVR for direct camera access.

## When to Use

- User asks "what's happening in the garage / front door / driveway?"
- User wants to check camera feeds or snapshots
- User asks about recent detections (people, cars, packages)
- User mentions security cameras, Frigate, or motion events
- User wants AI analysis of what a camera sees

## Architecture

```
LLM Vision (AI analysis)        Frigate NVR (direct access)
  analyze_camera                   get_snapshot
  analyze_video_event              get_events
  watch_stream                     get_thumbnail
                                   list_cameras
                                   camera_status
```

**LLM Vision** calls Home Assistant's LLM Vision integration to analyze camera feeds with a vision model (GPT-4o, Claude, Gemini).

**Frigate** calls the Frigate NVR API directly for raw snapshots, event queries, and camera status.

## Quick Reference

### "What's in the garage?"

```
vision({ action: "analyze_camera", camera_name: "garage" })
```

### "Is anyone at the front door?"

```
vision({ action: "analyze_camera", camera_name: "front_door", prompt: "Are there any people or packages visible?" })
```

### "Any people detected recently?"

```
vision({ action: "get_events", label: "person", hours_back: 4 })
```

### "Show me the driveway camera"

```
vision({ action: "get_snapshot", camera_name: "driveway" })
```

### "Watch the backyard for 30 seconds"

```
vision({ action: "watch_stream", camera_name: "backyard", duration: 30, prompt: "Watch for any animals" })
```

### "What triggered that alert?"

```
vision({ action: "analyze_video_event", event_id: "1234567890.abcdef" })
```

## Env Vars

| Variable             | Required For       | Example                           |
| -------------------- | ------------------ | --------------------------------- |
| HOME_ASSISTANT_URL   | LLM Vision actions | `http://homeassistant.local:8123` |
| HOME_ASSISTANT_TOKEN | LLM Vision actions | HA long-lived access token        |
| FRIGATE_URL          | Frigate actions    | `http://frigate:5000`             |

Either side works independently — you can use just Frigate for snapshots/events or just LLM Vision for AI analysis.


---

## Related

- [[00-24 SO (Sovereign Operations)/02 Infrastructure (VPS, Servers & Networking) (The Body)/02 Infrastructure|02 Infrastructure]]
