import type { AppStatusSource } from './source';
import type { StatusSnapshot } from './model';

export type StoreState = { snapshot?: StatusSnapshot; error: boolean; loadedAt?: number };

/**
 * One shared poller per source: every dot, header and System summary on the page subscribes to the
 * same snapshot, so a page of fifty apps costs one request per minute, not fifty.
 */
export function createStatusStore(source: AppStatusSource, intervalMs = 60000) {
  let state: StoreState = { error: false };
  const listeners = new Set<() => void>();
  let timer: ReturnType<typeof setInterval> | undefined;
  let pending = false;
  const publish = (next: StoreState) => {
    state = next;
    listeners.forEach(l => l());
  };
  const refresh = async () => {
    if (pending) return;
    pending = true;
    try {
      const snapshot = await source.read();
      publish({ snapshot, error: false, loadedAt: Date.now() });
    } catch {
      // Keep the last snapshot: the model marks it stale by age; the error flag explains why.
      publish({ ...state, error: true });
    } finally {
      pending = false;
    }
  };
  return {
    getState: () => state,
    refresh,
    subscribe(listener: () => void) {
      listeners.add(listener);
      if (listeners.size === 1) {
        void refresh();
        timer = setInterval(() => void refresh(), intervalMs);
      }
      return () => {
        listeners.delete(listener);
        if (!listeners.size && timer) {
          clearInterval(timer);
          timer = undefined;
        }
      };
    },
  };
}
export type StatusStore = ReturnType<typeof createStatusStore>;
