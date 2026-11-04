import { useSyncExternalStore } from 'react';

const KEY = 'onemind.theme.backgroundImage';
const listeners = new Set<() => void>();

function read(): boolean {
  try {
    return window.localStorage.getItem(KEY) === 'on';
  } catch {
    return false;
  }
}

let current = typeof window === 'undefined' ? false : read();

/** Whether the Legacy hero image is shown behind page headers. Default: off. */
export function getBackgroundImage(): boolean {
  return current;
}

export function setBackgroundImage(enabled: boolean): void {
  current = enabled;
  try {
    window.localStorage.setItem(KEY, enabled ? 'on' : 'off');
  } catch {
    // Storage can be blocked; the choice then lasts for this tab only.
  }
  listeners.forEach(l => l());
}

function subscribe(listener: () => void) {
  listeners.add(listener);
  return () => {
    listeners.delete(listener);
  };
}

/** `[enabled, setEnabled]`; re-renders when the setting changes. */
export function useBackgroundImage(): [boolean, (enabled: boolean) => void] {
  const enabled = useSyncExternalStore(
    subscribe,
    getBackgroundImage,
    () => false,
  );
  return [enabled, setBackgroundImage];
}
