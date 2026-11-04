import { useEffect, useRef } from 'react';
import { useSignal } from '@backstage/plugin-signals-react';

/** F3's Signals channel: a message arrives when an alert fires or resolves. */
export const ALERT_SIGNAL_CHANNEL = 'onemind:alerts';
/** Fallback when Signals is unavailable or a message is missed. */
export const POLL_INTERVAL_MS = 30000;
// Several alerts often arrive together; fetch once for the burst.
const BURST_MS = 400;

/**
 * Calls `refresh` once on mount, every 30 s, and shortly after every alert Signal.
 * Returns whether live push is connected, so the UI can say how it is updating.
 */
export function useAlertRefresh(
  refresh: () => void,
  intervalMs: number = POLL_INTERVAL_MS,
): { push: boolean } {
  const { lastSignal, isSignalsAvailable } = useSignal(ALERT_SIGNAL_CHANNEL);
  const latest = useRef(refresh);
  latest.current = refresh;

  useEffect(() => {
    latest.current();
    const timer = setInterval(() => latest.current(), intervalMs);
    return () => clearInterval(timer);
  }, [intervalMs]);

  useEffect(() => {
    if (!lastSignal) return undefined;
    const timer = setTimeout(() => latest.current(), BURST_MS);
    return () => clearTimeout(timer);
  }, [lastSignal]);

  return { push: isSignalsAvailable };
}
