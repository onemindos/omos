<template>
  <div class="onemind-ai-panel">
    <div class="ai-header">
      <span class="ai-title">⬡ Legacy AI</span>
      <span class="ai-status" :class="statusClass">{{ statusLabel }}</span>
    </div>

    <div class="ai-messages" ref="messagesEl">
      <div
        v-for="(msg, i) in messages"
        :key="i"
        class="ai-message"
        :class="msg.role"
      >
        <div class="ai-bubble">{{ msg.content }}</div>
        <div class="ai-meta">{{ msg.role === 'user' ? 'You' : 'Legacy' }} · {{ msg.time }}</div>
      </div>
      <div v-if="loading" class="ai-message assistant">
        <div class="ai-bubble ai-thinking">Legacy is thinking…</div>
      </div>
    </div>

    <div class="ai-input-row">
      <textarea
        v-model="draft"
        placeholder="Message Legacy…"
        rows="2"
        @keydown.enter.exact.prevent="send"
        :disabled="loading"
      />
      <button @click="send" :disabled="loading || !draft.trim()">Send</button>
    </div>

    <div class="ai-quick-row">
      <button v-for="q in quickCommands" :key="q" class="ai-quick" @click="quickSend(q)">
        {{ q }}
      </button>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, nextTick, onMounted } from 'vue';

interface Message {
  role: 'user' | 'assistant';
  content: string;
  time: string;
}

const BRIDGE_URL = (window as any).__ONEMIND_BRIDGE_URL__ || '/api/onemind/chat';

const messages = ref<Message[]>([]);
const draft = ref('');
const loading = ref(false);
const connected = ref(false);

const quickCommands = [
  'SA picture',
  'List missions',
  'Weather',
  'Ship count',
  'Air picture',
];

const statusClass = computed(() => connected.value ? 'online' : 'offline');
const statusLabel = computed(() => connected.value ? 'online' : 'connecting…');

const messagesEl = ref<HTMLElement | null>(null);

function now() {
  return new Date().toLocaleTimeString([], { hour: '2-digit', minute: '2-digit' });
}

async function send() {
  const text = draft.value.trim();
  if (!text || loading.value) return;
  draft.value = '';

  messages.value.push({ role: 'user', content: text, time: now() });
  loading.value = true;
  await scrollBottom();

  try {
    const res = await fetch(BRIDGE_URL, {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({ message: text }),
    });
    const data = await res.json();
    connected.value = true;
    messages.value.push({
      role: 'assistant',
      content: data.response || data.message || '(no response)',
      time: now(),
    });
  } catch (e: any) {
    messages.value.push({
      role: 'assistant',
      content: `⚠ Bridge unreachable: ${e.message}`,
      time: now(),
    });
    connected.value = false;
  } finally {
    loading.value = false;
    await scrollBottom();
  }
}

function quickSend(cmd: string) {
  draft.value = cmd;
  send();
}

async function scrollBottom() {
  await nextTick();
  if (messagesEl.value) {
    messagesEl.value.scrollTop = messagesEl.value.scrollHeight;
  }
}

onMounted(async () => {
  try {
    const res = await fetch(BRIDGE_URL.replace('/chat', '/health'));
    connected.value = res.ok;
  } catch {
    connected.value = false;
  }
});
</script>

<style scoped>
.onemind-ai-panel {
  display: flex;
  flex-direction: column;
  height: 100%;
  background: #0d1117;
  color: #e6edf3;
  font-family: 'Inter', sans-serif;
  font-size: 13px;
}
.ai-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 12px 16px 8px;
  border-bottom: 1px solid #21262d;
}
.ai-title { font-weight: 700; letter-spacing: 0.5px; color: #58a6ff; }
.ai-status { font-size: 11px; padding: 2px 8px; border-radius: 10px; }
.ai-status.online  { background: #1a4731; color: #3fb950; }
.ai-status.offline { background: #3d1f1f; color: #f85149; }
.ai-messages {
  flex: 1;
  overflow-y: auto;
  padding: 12px 16px;
  display: flex;
  flex-direction: column;
  gap: 8px;
}
.ai-message { display: flex; flex-direction: column; }
.ai-message.user  { align-items: flex-end; }
.ai-message.assistant { align-items: flex-start; }
.ai-bubble {
  max-width: 85%;
  padding: 8px 12px;
  border-radius: 12px;
  line-height: 1.5;
  white-space: pre-wrap;
  word-break: break-word;
}
.ai-message.user .ai-bubble      { background: #1f6feb; color: #fff; border-radius: 12px 12px 4px 12px; }
.ai-message.assistant .ai-bubble { background: #161b22; border: 1px solid #30363d; border-radius: 12px 12px 12px 4px; }
.ai-thinking { opacity: 0.6; font-style: italic; }
.ai-meta { font-size: 10px; color: #484f58; margin-top: 3px; padding: 0 4px; }
.ai-input-row {
  display: flex;
  gap: 8px;
  padding: 10px 16px;
  border-top: 1px solid #21262d;
}
.ai-input-row textarea {
  flex: 1;
  background: #161b22;
  border: 1px solid #30363d;
  border-radius: 6px;
  color: #e6edf3;
  padding: 8px 10px;
  resize: none;
  font-size: 13px;
  font-family: inherit;
}
.ai-input-row textarea:focus { outline: none; border-color: #58a6ff; }
.ai-input-row button {
  background: #238636;
  color: #fff;
  border: none;
  border-radius: 6px;
  padding: 0 16px;
  cursor: pointer;
  font-size: 13px;
  font-weight: 600;
}
.ai-input-row button:disabled { opacity: 0.4; cursor: default; }
.ai-quick-row {
  display: flex;
  flex-wrap: wrap;
  gap: 6px;
  padding: 8px 16px 12px;
}
.ai-quick {
  background: #21262d;
  color: #8b949e;
  border: 1px solid #30363d;
  border-radius: 14px;
  padding: 3px 10px;
  font-size: 11px;
  cursor: pointer;
}
.ai-quick:hover { background: #30363d; color: #e6edf3; }
</style>
