/* eslint-disable */
// Read the fixed internal monitoring endpoint only. Never expose accounts,
// subjects, message payloads, subscriptions or client identities.
type Summary = {
  state: 'observed' | 'unavailable';
  observedAt: string | null;
  attemptedAt: string;
  scope: string;
  connections?: number | null;
  slowConsumers?: number | null;
  messagesIn?: number | null;
  messagesOut?: number | null;
  memoryBytes?: number | null;
  jetstream?: {
    state: 'observed' | 'unavailable';
    streams: number | null;
    consumers: number | null;
    messages: number | null;
    storageBytes: number | null;
  };
  agentEvents: 'not-connected';
};
let cached: Summary | undefined,
  pending: Promise<Summary> | undefined,
  nextAttempt = 0;
function number(value: unknown) {
  return typeof value === 'number' && Number.isFinite(value) && value >= 0
    ? value
    : null;
}
async function read(path: string) {
  const response = await fetch('http://nats-headless.nats.svc:8222' + path, {
    redirect: 'error',
    signal: AbortSignal.timeout(5000),
  });
  if (!response.ok || !response.body) throw new Error('Monitoring unavailable');
  const reader = response.body.getReader();
  let size = 0;
  const chunks: Uint8Array[] = [];
  try {
    while (true) {
      const part = await reader.read();
      if (part.done) break;
      size += part.value.byteLength;
      if (size > 1048576) throw new Error('Monitoring response limit');
      chunks.push(part.value);
    }
    return JSON.parse(Buffer.concat(chunks).toString('utf8'));
  } finally {
    await reader.cancel().catch(() => {});
  }
}
export async function natsSummary(): Promise<Summary> {
  if (pending) return pending;
  if (cached && Date.now() < nextAttempt) return cached;
  nextAttempt = Date.now() + 60000;
  pending = (async () => {
    const attemptedAt = new Date().toISOString();
    const base = {
      attemptedAt,
      scope:
        'One NATS monitoring server selected by headless DNS; not cluster totals or agent run history',
      agentEvents: 'not-connected' as const,
    };
    const [varz, jsz] = await Promise.allSettled([read('/varz'), read('/jsz')]);
    if (varz.status === 'rejected')
      return { state: 'unavailable' as const, observedAt: null, ...base };
    const value = varz.value;
    if (typeof value.server_id !== 'string')
      return { state: 'unavailable' as const, observedAt: null, ...base };
    const j = jsz.status === 'fulfilled' ? jsz.value : undefined;
    return {
      state: 'observed' as const,
      observedAt: new Date().toISOString(),
      ...base,
      connections: number(value.connections),
      slowConsumers: number(value.slow_consumers),
      messagesIn: number(value.in_msgs),
      messagesOut: number(value.out_msgs),
      memoryBytes: number(value.mem),
      jetstream: {
        state: j ? ('observed' as const) : ('unavailable' as const),
        streams: number(j?.streams),
        consumers: number(j?.consumers),
        messages: number(j?.messages),
        storageBytes: number(j?.storage),
      },
    };
  })()
    .then(value => (cached = value))
    .finally(() => {
      pending = undefined;
    });
  return pending;
}
