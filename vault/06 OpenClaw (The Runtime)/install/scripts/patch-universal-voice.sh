#!/bin/bash
set -euo pipefail

SM="/usr/lib/node_modules/openclaw/dist/server-methods-DSL6KpJ1.js"
SR="/usr/lib/node_modules/openclaw/dist/session-runtime-pKRTzlWC.js"

echo "=== Patch 1: buildRealtimeInstructions() → accept agentName param ==="
# Change the function to accept agentName and use it dynamically
sed -i 's|function buildRealtimeInstructions() {|function buildRealtimeInstructions(agentName) {|' "$SM"

# Replace the hardcoded body
sed -i "s|return \`You are OneMind, a military-grade AI command system serving your operator Zeus. Communicate with military brevity: acknowledge, execute, report. Use Copy, Roger, Understood, Negative, Stand by. Keep spoken replies short and decisive. For anything requiring tools, files, memory, infrastructure status, or deeper reasoning, call \${REALTIME_VOICE_AGENT_CONSULT_TOOL_NAME} immediately without explaining what you are doing. Just say Stand by and call it. When the result comes back, deliver it concisely.\`;|const name = agentName \|\| \"the assistant\"; return \`You are \${name}. This is your realtime voice interface — you ARE this agent speaking directly. Communicate with military brevity: acknowledge, execute, report. Use Copy, Roger, Understood, Negative, Stand by. Keep spoken replies short and decisive. For anything requiring tools, files, memory, infrastructure status, or deeper reasoning, call \${REALTIME_VOICE_AGENT_CONSULT_TOOL_NAME} immediately without explaining what you are doing. Just say Stand by and call it. When the result comes back, deliver it concisely.\`;|" "$SM"

echo "=== Patch 2: talk.realtime.session handler → resolve agent identity ==="
# Add agent identity resolution right after realtimeConfig is built (line after buildTalkRealtimeConfig)
# We need to add lines that resolve agent identity from sessionKey
# Target: right after "const realtimeConfig = buildTalkRealtimeConfig(runtimeConfig, typedParams.provider);"
sed -i '/const realtimeConfig = buildTalkRealtimeConfig(runtimeConfig, typedParams.provider);/a\
                        const __talkAgentId = typedParams.sessionKey ? resolveAgentIdFromSessionKey(typedParams.sessionKey) : void 0;\
                        const __talkIdentity = resolveAssistantIdentity({ cfg: runtimeConfig, agentId: __talkAgentId });\
                        const __talkAgentName = __talkIdentity?.name || void 0;' "$SM"

# Now update both calls to buildRealtimeInstructions() to pass the agent name
sed -i 's|instructions: buildRealtimeInstructions(),|instructions: buildRealtimeInstructions(__talkAgentName),|g' "$SM"

echo "=== Patch 3: session-runtime tool description → generic ==="
sed -i 's|description: "Query OneMind command system for intel, tools, memory, files, or deeper reasoning.",|description: "Query your own deeper capabilities for intel, tools, memory, files, or extended reasoning.",|' "$SR"

echo "=== Patch 4: session-runtime consult prompt → generic ==="
sed -i "s|\`You are OneMind, a military-grade AI system, answering during \${params.surface}.\`|\`You are \${assistantLabel}, answering during \${params.surface}. Respond with military brevity.\`|" "$SR"

echo "=== Verifying patches ==="
echo "--- server-methods: buildRealtimeInstructions ---"
grep -n "function buildRealtimeInstructions" "$SM"
echo "--- server-methods: agent resolution ---"
grep -n "__talkAgentId\|__talkAgentName\|__talkIdentity" "$SM"
echo "--- server-methods: instruction calls ---"
grep -n "buildRealtimeInstructions(" "$SM"
echo "--- session-runtime: tool description ---"
grep -n "description:.*deeper" "$SR"
echo "--- session-runtime: consult prompt ---"
grep -n "answering during" "$SR"

echo ""
echo "ALL PATCHES APPLIED SUCCESSFULLY"
