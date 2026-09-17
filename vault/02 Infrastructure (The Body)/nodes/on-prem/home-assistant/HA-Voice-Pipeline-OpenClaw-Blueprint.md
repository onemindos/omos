# HA Voice Pipeline → OpenClaw Blueprint
**Status:** Draft for Zeus review  
**Date:** 2026-04-22  
**Context:** Atom Echo (M5Stack) + HA Assist → replace HA Cloud LLM with OpenClaw chatCompletions API  

---

## Overview

Route Home Assistant voice pipeline (Atom Echo wake word + Assist) through OpenClaw's `chatCompletions` API endpoint instead of HA Cloud. Enables:
- Local LLM inference (no HA Cloud dependency)
- Natural language control of all HA entities (lights, switches, scenes, automations)
- Full conversation history in OpenClaw sessions
- Fallback to HA Assist if OpenClaw is unavailable

**Why this matters:**
- HA Cloud is optional + cost
- OpenClaw is already running on VPS (chatCompletions enabled 2026-04-16)
- Voice becomes part of the integrated One Mind OS flow

---

## Architecture

```
Atom Echo (M5Stack)
    ↓ (wake word + STT)
HA Assist (local speech-to-text)
    ↓ (text intent)
OpenClaw chatCompletions
    ↓ (LLM response: intent classification + entity control)
HA Service Calls (light.turn_on, switch.toggle, etc.)
    ↓
Smart Home Action
```

**Current state:**
- Atom Echo: ESPHome, in HA as `assist_satellite.m5stack_atom_echo_30dc34_assist_satellite` (unavailable — not powered)
- Assist Microphone: separate device, also configured
- Extended OpenAI Conversation: already in HA config (unavailable) — legacy integration we'll repurpose

---

## Implementation Steps

### 1. Verify HA Assist Pipeline (local STT)

```yaml
# configuration.yaml
assist_pipeline:
  debug: true  # optional, logs speech intents
```

Check HA Settings → Devices & Services → Voice Assistants → Assist — should show "local" for all components.

### 2. Configure OpenClaw chatCompletions Endpoint

**VPS gateway info:**
- Base URL: `https://omos-core.tail717ba0.ts.net` (Tailscale Funnel)
- chatCompletions port: `18789`
- chatCompletions endpoint: `/api/v1/chat/completions`
- Auth: Bearer token (OpenClaw gateway auth token from 1Password SO vault)

**Example request:**
```bash
curl -X POST https://omos-core.tail717ba0.ts.net/api/v1/chat/completions \
  -H "Authorization: Bearer <gateway-token>" \
  -H "Content-Type: application/json" \
  -d '{
    "messages": [
      {"role": "user", "content": "turn on the kitchen lights"}
    ],
    "model": "github-copilot/claude-sonnet-4.6",
    "max_tokens": 500
  }'
```

### 3. Create Custom Conversation Agent (Option A: REST Command)

**Add to `configuration.yaml`:**
```yaml
# REST command to forward conversation to OpenClaw
rest_command:
  openclaw_conversation:
    url: "https://omos-core.tail717ba0.ts.net/api/v1/chat/completions"
    method: POST
    headers:
      Authorization: !secret openclaw_gateway_token
      Content-Type: "application/json"
    payload: |
      {
        "messages": [
          {"role": "user", "content": "{{ message }}"}
        ],
        "model": "github-copilot/claude-sonnet-4.6",
        "max_tokens": 500
      }
    timeout: 30
```

**Automation to route Assist → OpenClaw:**
```yaml
automation:
  - alias: "Voice Intent → OpenClaw"
    trigger:
      platform: state
      entity_id: conversation.home_assistant
    action:
      - service: rest_command.openclaw_conversation
        data:
          message: "{{ trigger.to_state.state }}"
      - service: tts.google_translate_say
        target:
          entity_id: media_player.m5stack_atom_echo_30dc34
        data:
          message: "{{ state_attr('rest_command.openclaw_conversation', 'response') }}"
```

### 4. Alternative: HACS Custom Integration (Option B: More Native)

Install `conversation-openai-compatible` or `openai-custom-conversation` HACS integration:
1. HACS → Integrations → Search "openai custom"
2. Install → Restart HA
3. Settings → Devices & Services → Create Integration → OpenAI Compatible
4. Set URL: `https://omos-core.tail717ba0.ts.net/api/v1/chat/completions`
5. Auth token: `<openclaw-gateway-token>`
6. Model: `github-copilot/claude-sonnet-4.6`

Then set Atom Echo Assist to use this as conversation endpoint.

### 5. Test Flow

**Manually:**
```bash
# From Mac node or OpenClaw:
curl -X POST https://omos-core.tail717ba0.ts.net/api/v1/chat/completions \
  -H "Authorization: Bearer <token>" \
  -H "Content-Type: application/json" \
  -d '{
    "messages": [{"role": "user", "content": "what lights do i have"}],
    "model": "github-copilot/claude-sonnet-4.6"
  }' | jq '.choices[0].message.content'
```

**Via HA UI:**
1. Power on Atom Echo
2. Say wake word (default: "Alexa")
3. Request: "turn on the kitchen lights"
4. Expected: Lights turn on + Echo responds with confirmation

### 6. Entity Mapping for Intent Recognition

OpenClaw LLM needs to know what entities HA has. Include in the system prompt:

```yaml
# HA Template to feed available entities
assistant:
  homeassistant_entities: |
    Available smart home entities:
    - Lights: {{ states | selectattr('entity_id', 'match', 'light\..*') | map(attribute='entity_id') | list }}
    - Switches: {{ states | selectattr('entity_id', 'match', 'switch\..*') | map(attribute='entity_id') | list }}
    - Scenes: {{ states | selectattr('entity_id', 'match', 'scene\..*') | map(attribute='entity_id') | list }}
    - Automations: {{ states | selectattr('entity_id', 'match', 'automation\..*') | map(attribute='entity_id') | list }}
```

Or simpler: just send the current entity list as context in each request.

---

## Challenges & Mitigations

| Challenge | Mitigation |
|-----------|-----------|
| OpenClaw endpoint not accessible from HA (LAN IP 192.168.1.164) | Use Tailscale IP `100.98.25.10` or Funnel URL; HA is on tailnet |
| Token expiry / rotation | Store in HA secrets; refresh monthly |
| LLM doesn't understand intent → no entity control | Fine-tune system prompt with entity names; test with `hass-cli state list` |
| Atom Echo unavailable (not powered) | Power supply + check Tailscale connection |
| Fallback to HA Cloud if OpenClaw is down | Keep `conversation.home_assistant` as fallback; OR use HA Assist alone (local LLM inference not implemented yet) |

---

## Success Criteria

- [ ] Atom Echo shows `assist_satellite.m5stack_atom_echo_30dc34_assist_satellite: idle` (not unavailable)
- [ ] Say wake word + request results in OpenClaw processing (check logs)
- [ ] Entity control works: "turn on lights" → HA service call executed
- [ ] Conversation history appears in OpenClaw session logs
- [ ] TTS response plays back on Atom Echo

---

## Next Phase: HA-to-OpenClaw MCP Server

Once voice pipeline is working, implement Home Assistant as an MCP (Model Context Protocol) server inside OpenClaw. This allows:
- OpenClaw agents to directly query HA state (not just via chatCompletions)
- Bidirectional sync: HA automations can trigger OpenClaw workflows
- Full "smart home as a compute platform" integration

**Sketch:**
```
OpenClaw Gateway
  ├── MCP Server: Home Assistant
  │   ├── Tools: state_get, state_set, service_call, automation_trigger
  │   └── Context: All HA entities + history
  └── Tools: (existing exec, web, etc.)
```

---

## References

- **OpenClaw chatCompletions API:** https://docs.openclaw.ai/extensions/chat-completions
- **HA Assist:** https://www.home-assistant.io/docs/assist/
- **HA REST API:** https://developers.home-assistant.io/docs/api/rest/
- **HACS OpenAI Integration:** https://github.com/vit-lebeda/home-assistant-openai-custom-conversation
- **Atom Echo + ESPHome:** https://github.com/esphome/home-assistant-addon

---

## Decision Points for Zeus

1. **Implementation path:** REST command (Option A, simpler) or HACS integration (Option B, more native)?
2. **Fallback behavior:** Keep HA Cloud as backup or rely on Assist-only mode?
3. **Entity mapping:** Auto-generate system prompt from HA state or manually maintain?
4. **Next priority:** Get voice working first, then move to MCP server phase?



---

## Related

- [[00-24 SO (Sovereign Operations)/02 Infrastructure (VPS, Servers & Networking) (The Body)/02 Infrastructure|02 Infrastructure]]
