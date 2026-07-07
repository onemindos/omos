<template>
    <div class='onemind-ai-panel'>
        <div class='ai-header'>
            <span class='ai-title'>⬡ Legacy AI</span>
            <span
                class='ai-status'
                :class='statusClass'
            >{{ statusLabel }}</span>
        </div>

        <div
            ref='messagesEl'
            class='ai-messages'
        >
            <div
                v-for='(msg, i) in messages'
                :key='i'
                class='ai-message'
                :class='msg.role'
            >
                <div class='ai-bubble'>{{ msg.content }}</div>
                <div class='ai-meta'>{{ msg.role === 'user' ? 'You' : 'Legacy' }} · {{ msg.time }}</div>
            </div>

            <div
                v-if='streaming'
                class='ai-message assistant'
            >
                <div class='ai-bubble'>
                    <template v-if='streamContent'>{{ streamContent }}</template>
                    <span
                        v-else
                        class='cursor'
                    >▌</span>
                </div>
            </div>
        </div>

        <div class='ai-input-row'>
            <textarea
                v-model='draft'
                placeholder='Message Legacy…'
                rows='2'
                :disabled='streaming'
                @keydown.enter.exact.prevent='send'
            />
            <button
                :class='{ "btn-stop": streaming }'
                :disabled='!streaming && !draft.trim()'
                @click='streaming ? cancelStream() : send()'
            >{{ streaming ? '■' : 'Send' }}</button>
        </div>

        <div class='ai-quick-row'>
            <button
                v-for='q in quickCommands'
                :key='q'
                class='ai-quick'
                :disabled='streaming'
                @click='quickSend(q)'
            >{{ q }}</button>
        </div>
    </div>
</template>

<script setup lang='ts'>
import { ref, computed, nextTick, onMounted, onUnmounted } from 'vue';

interface Message {
    role: 'user' | 'assistant';
    content: string;
    time: string;
}

interface BridgeChunk {
    kind: string;
    id?: string;
    text?: string;
    response?: string;
    message?: string;
}

// Injected by api.yaml Patch 4 into index.html at pod startup.
// nginx Patch 5 proxies /onemind/ → the bridge pod cluster-internally.
const BRIDGE_BASE: string = (window as Window & { __ONEMIND_BRIDGE_URL__?: string }).__ONEMIND_BRIDGE_URL__ ?? '/onemind';

function getWsUrl(): string {
    if (BRIDGE_BASE.startsWith('/')) {
        const proto = location.protocol === 'https:' ? 'wss:' : 'ws:';
        return `${proto}//${location.host}${BRIDGE_BASE}/ws`;
    }
    try {
        const u = new URL(BRIDGE_BASE);
        u.protocol = u.protocol === 'https:' ? 'wss:' : 'ws:';
        u.pathname = u.pathname.replace(/\/$/, '') + '/ws';
        return u.toString();
    } catch {
        return `ws://${location.host}/onemind/ws`;
    }
}

const messages      = ref<Message[]>([]);
const draft         = ref('');
const streaming     = ref(false);
const streamContent = ref('');
const wsStatus      = ref<'connecting' | 'open' | 'closed'>('connecting');
const messagesEl    = ref<HTMLElement | null>(null);

const quickCommands = ['SA picture', 'List missions', 'Weather', 'Ship count', 'Air picture'];

const statusClass = computed(() => ({
    online:     wsStatus.value === 'open',
    connecting: wsStatus.value === 'connecting',
    offline:    wsStatus.value === 'closed',
}));

const statusLabel = computed(() => {
    if (wsStatus.value === 'open') return 'online';
    if (wsStatus.value === 'connecting') return 'connecting…';
    return 'offline';
});

let ws: WebSocket | null = null;
let reconnectTimer: ReturnType<typeof setTimeout> | null = null;
let reconnectAttempts = 0;
let activeId: string | null = null;
const streams = new Map<string, (msg: BridgeChunk) => void>();

function connectWs(): void {
    if (ws) return;
    wsStatus.value = 'connecting';
    ws = new WebSocket(getWsUrl());

    ws.addEventListener('open', () => {
        wsStatus.value = 'open';
        reconnectAttempts = 0;
    });

    ws.addEventListener('message', (ev) => {
        try {
            const msg = JSON.parse(String(ev.data)) as BridgeChunk;
            if (msg.id) streams.get(msg.id)?.(msg);
        } catch { /* ignore malformed */ }
    });

    ws.addEventListener('close', () => {
        wsStatus.value = 'closed';
        ws = null;
        for (const [, handler] of streams) {
            handler({ kind: 'error', message: 'connection closed' });
        }
        streams.clear();
        scheduleReconnect();
    });
}

function scheduleReconnect(): void {
    if (reconnectTimer !== null) return;
    const delay = Math.min(15000, 500 * 2 ** Math.min(reconnectAttempts, 6));
    reconnectAttempts++;
    reconnectTimer = setTimeout(() => {
        reconnectTimer = null;
        connectWs();
    }, delay);
}

function now(): string {
    return new Date().toLocaleTimeString([], { hour: '2-digit', minute: '2-digit' });
}

async function scrollBottom(): Promise<void> {
    await nextTick();
    if (messagesEl.value) messagesEl.value.scrollTop = messagesEl.value.scrollHeight;
}

function makeId(): string {
    return crypto.randomUUID();
}

async function send(): Promise<void> {
    const text = draft.value.trim();
    if (!text || streaming.value) return;
    draft.value = '';

    messages.value.push({ role: 'user', content: text, time: now() });
    streaming.value = true;
    streamContent.value = '';
    await scrollBottom();

    const id = makeId();
    activeId = id;

    if (!ws || ws.readyState !== WebSocket.OPEN) {
        await sendHttp(text);
        return;
    }

    streams.set(id, (msg: BridgeChunk) => {
        if (msg.kind === 'chunk') {
            streamContent.value += msg.text ?? '';
            scrollBottom();
        } else if (msg.kind === 'done') {
            const content = streamContent.value || msg.response || '(no response)';
            messages.value.push({ role: 'assistant', content, time: now() });
            streaming.value = false;
            streamContent.value = '';
            activeId = null;
            streams.delete(id);
            scrollBottom();
        } else if (msg.kind === 'error') {
            messages.value.push({
                role: 'assistant',
                content: `⚠ ${msg.message ?? 'Bridge error'}`,
                time: now(),
            });
            streaming.value = false;
            streamContent.value = '';
            activeId = null;
            streams.delete(id);
            scrollBottom();
        }
    });

    ws.send(JSON.stringify({ kind: 'prompt', id, text, callsign: 'Operator' }));
}

async function sendHttp(text: string): Promise<void> {
    try {
        const res = await fetch(`${BRIDGE_BASE}/chat`, {
            method: 'POST',
            headers: { 'Content-Type': 'application/json' },
            body: JSON.stringify({ message: text, callsign: 'Operator' }),
        });
        const data = await res.json() as { response?: string };
        messages.value.push({
            role: 'assistant',
            content: data.response ?? '(no response)',
            time: now(),
        });
    } catch (e: unknown) {
        messages.value.push({
            role: 'assistant',
            content: `⚠ ${e instanceof Error ? e.message : 'Bridge unreachable'}`,
            time: now(),
        });
    } finally {
        streaming.value = false;
        streamContent.value = '';
        await scrollBottom();
    }
}

function cancelStream(): void {
    if (activeId && ws && ws.readyState === WebSocket.OPEN) {
        ws.send(JSON.stringify({ kind: 'cancel', id: activeId }));
    }
    if (activeId) streams.delete(activeId);
    streaming.value = false;
    streamContent.value = '';
    activeId = null;
}

function quickSend(cmd: string): void {
    draft.value = cmd;
    send();
}

onMounted(() => connectWs());
onUnmounted(() => {
    if (reconnectTimer !== null) clearTimeout(reconnectTimer);
    ws?.close();
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
.ai-status.online     { background: #1a4731; color: #3fb950; }
.ai-status.connecting { background: #2d2a1e; color: #d29922; }
.ai-status.offline    { background: #3d1f1f; color: #f85149; }
.ai-messages {
    flex: 1;
    overflow-y: auto;
    padding: 12px 16px;
    display: flex;
    flex-direction: column;
    gap: 8px;
}
.ai-message { display: flex; flex-direction: column; }
.ai-message.user      { align-items: flex-end; }
.ai-message.assistant { align-items: flex-start; }
.ai-bubble {
    max-width: 85%;
    padding: 8px 12px;
    border-radius: 12px;
    line-height: 1.5;
    white-space: pre-wrap;
    word-break: break-word;
}
.ai-message.user .ai-bubble {
    background: #1f6feb;
    color: #fff;
    border-radius: 12px 12px 4px 12px;
}
.ai-message.assistant .ai-bubble {
    background: #161b22;
    border: 1px solid #30363d;
    border-radius: 12px 12px 12px 4px;
}
.cursor {
    display: inline-block;
    animation: blink 1s step-end infinite;
    color: #58a6ff;
}
@keyframes blink { 50% { opacity: 0; } }
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
    min-width: 56px;
}
.ai-input-row button:disabled { opacity: 0.4; cursor: default; }
.ai-input-row button.btn-stop { background: #da3633; }
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
.ai-quick:hover:not(:disabled) { background: #30363d; color: #e6edf3; }
.ai-quick:disabled { opacity: 0.4; cursor: default; }
</style>
