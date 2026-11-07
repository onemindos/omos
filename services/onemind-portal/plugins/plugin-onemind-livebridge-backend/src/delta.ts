/** Every Signals message must stay below this many bytes of JSON (the limit is 60 KB). */
export const MAX_MESSAGE_BYTES = 56 * 1024;

const size = (value: unknown) => Buffer.byteLength(JSON.stringify(value));

/**
 * Split a delta into messages that each serialise below `maxBytes`. `base` holds the envelope and
 * scalar fields; `arrays` holds the list fields (every message carries every key, empty when not
 * used). One message when it fits. Otherwise each list is chunked on its own, and a message never
 * has all lists empty. Items that alone exceed the budget are dropped (cannot be sent).
 */
export function splitDelta<B extends object, K extends string>(
  base: B,
  arrays: Record<K, unknown[]>,
  maxBytes = MAX_MESSAGE_BYTES,
): Array<B & Record<K, unknown[]>> {
  const keys = Object.keys(arrays) as K[];
  const empty = () =>
    Object.fromEntries(keys.map(k => [k, []])) as unknown as Record<
      K,
      unknown[]
    >;
  if (keys.every(k => arrays[k].length === 0)) return [];
  const whole = { ...base, ...arrays };
  if (size(whole) <= maxBytes) return [whole];

  const budget = maxBytes - size({ ...base, ...empty() });
  const out: Array<B & Record<K, unknown[]>> = [];
  for (const key of keys) {
    let chunk: unknown[] = [];
    let used = 0;
    const flush = () => {
      if (chunk.length > 0) out.push({ ...base, ...empty(), [key]: chunk });
      chunk = [];
      used = 0;
    };
    for (const item of arrays[key]) {
      const bytes = size(item) + 1;
      if (bytes > budget) continue;
      if (used + bytes > budget) flush();
      chunk.push(item);
      used += bytes;
    }
    flush();
  }
  return out;
}
