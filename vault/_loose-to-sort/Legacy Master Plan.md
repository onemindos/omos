# Legacy Master Plan

You type or speak to Legacy in Telegram (or Talk mode, or Frame glasses).

Legacy has ALL the power:
- Direct hardware access (plugins — in-process Python)
- All coding agents (fabric_prompt to CC/Pi over NATS)
- All edge devices (nats-fleet pub/sub)
- All IoT (HomeAssistant MCP)
- All file/knowledge (filesystem MCP + vault)
- All scheduling (cron)
- All background work (delegation)
- All rendering (fabric_prompt to OC → Canvas)
- All off-grid (meshtastic plugin)

---

## Daily Operations

**Morning:** Wake up → Telegram notification from Legacy:
> "Morning report ready. Weather: 32°C. Irrigation ran 4am-5am.
> Drone patrol found 2 deer near north fence. All systems normal."

**Working:** "Legacy, spray schedule for wheat field 3?"
→ Legacy checks weather, growth stage, previous applications
→ Responds with recommendation + adds to Todoist

**Emergency:** Frame glasses detect smoke → vision-detect plugin fires
→ Legacy: "⚠️ Smoke detected NE quadrant. Launching drone for visual.
Alert sent to your phone. Irrigation activated as precaution."

**Evening:** "Legacy, run evening summary"
→ Delegates: compile day's events, fleet status, tomorrow's plan
→ Renders on Canvas (OpenClaw) + sends Telegram digest

**Night:**
- 02:00 cron → backup all services
- 22:00 cron → security sweep (drone + cameras + sensors)

You sleep. Legacy handles it.
