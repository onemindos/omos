<template>
    <div
        class='onemind-bottom'
        @click='openPanel'
    >
        <span :class='dotClass' />
        <span>Legacy AI</span>
    </div>
</template>

<script setup lang='ts'>
import { ref, computed, onMounted, onUnmounted } from 'vue';
import { useRouter } from 'vue-router';

const alive = ref(false);
const router = useRouter();
const dotClass = computed(() => alive.value ? 'dot green' : 'dot red');
const BRIDGE_BASE: string = (window as Window & { __ONEMIND_BRIDGE_URL__?: string }).__ONEMIND_BRIDGE_URL__ ?? '/onemind';

async function ping(): Promise<void> {
    try {
        const r = await fetch(`${BRIDGE_BASE}/health`, { signal: AbortSignal.timeout(3000) });
        alive.value = r.ok;
    } catch {
        alive.value = false;
    }
}

function openPanel(): void {
    router.push({ name: 'home-menu-onemind-ai' });
}

let interval: ReturnType<typeof setInterval>;
onMounted(() => {
    ping();
    interval = setInterval(ping, 15000);
});
onUnmounted(() => clearInterval(interval));
</script>

<style scoped>
.onemind-bottom {
    display: flex;
    align-items: center;
    gap: 5px;
    font-size: 11px;
    color: #8b949e;
    cursor: pointer;
    padding: 0 8px;
}
.onemind-bottom:hover { color: #e6edf3; }
.dot { width: 6px; height: 6px; border-radius: 50%; display: inline-block; }
.dot.green { background: #3fb950; }
.dot.red   { background: #f85149; }
</style>
