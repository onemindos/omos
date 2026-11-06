import { useCallback, useEffect, useState } from 'react';
import { discoveryApiRef, fetchApiRef, useApi } from '@backstage/core-plugin-api';
import type { AppStatus } from '@internal/plugin-onemind-common';

/** Gatus statuses API shape (only the fields the card reads). `duration` is nanoseconds. */
export type GatusResult = {
  success: boolean;
  timestamp: string;
  duration?: number;
  status?: number;
  errors?: string[];
  conditionResults?: { condition: string; success: boolean }[];
};
export type GatusEndpoint = { name: string; group?: string; key?: string; results?: GatusResult[] };

export type GatusSummary =
  | {
      kind: 'checked';
      /** `unknown` when the last check is too old to trust, never green. */
      status: AppStatus;
      stale: boolean;
      lastCheck: string;
      responseMs?: number;
      httpStatus?: number;
      failing: string[];
      errors: string[];
      /** The 14 day condition: Gatus does not export the exact expiry in its API, only whether the condition held. */
      certificate: 'valid' | 'expiring' | 'not-checked';
    }
  | { kind: 'not-checked' }
  | { kind: 'unavailable'; reason: string };

/** A check older than this is not counted as current (the slowest interval is 5 minutes). */
export const STALE_AFTER_MS = 15 * 60 * 1000;
const CERT_CONDITION = 'CERTIFICATE_EXPIRATION';

/** Pure: the statuses list to what the card shows for one endpoint name. */
export function summariseEndpoint(statuses: GatusEndpoint[], name: string, now: number = Date.now()): GatusSummary {
  const endpoint = statuses.find(s => s.name === name);
  const last = endpoint?.results?.[endpoint.results.length - 1];
  if (!endpoint || !last) return { kind: 'not-checked' };
  const age = now - new Date(last.timestamp).getTime();
  const stale = !Number.isFinite(age) || age > STALE_AFTER_MS;
  let status: AppStatus = 'unknown';
  if (!stale) status = last.success ? 'ok' : 'down';
  const cert = last.conditionResults?.find(c => c.condition.includes(CERT_CONDITION));
  let certificate: 'valid' | 'expiring' | 'not-checked' = 'not-checked';
  if (cert) certificate = cert.success ? 'valid' : 'expiring';
  return {
    kind: 'checked',
    status,
    stale,
    lastCheck: last.timestamp,
    responseMs: typeof last.duration === 'number' ? Math.round(last.duration / 1e6) : undefined,
    httpStatus: last.status,
    failing: (last.conditionResults ?? []).filter(c => !c.success).map(c => c.condition),
    errors: last.errors ?? [],
    certificate,
  };
}

export type GatusState = { loading: boolean; summary?: GatusSummary; readAt?: string };

export class GatusError extends Error {}

type Fetcher = { fetch(url: string): Promise<Response> };
type Discovery = { getBaseUrl(pluginId: string): Promise<string> };

/** One read of the statuses list through the portal proxy. Throws on a failed or malformed answer. */
export async function readGatusStatuses(discovery: Discovery, fetchApi: Fetcher): Promise<GatusEndpoint[]> {
  const base = await discovery.getBaseUrl('proxy');
  const response = await fetchApi.fetch(`${base}/gatus/statuses?pageSize=1`);
  if (!response.ok) throw new GatusError(`Gatus is not reachable (HTTP ${response.status})`);
  const body = await response.json();
  if (!Array.isArray(body)) throw new GatusError('Gatus answered with an unexpected body');
  return body as GatusEndpoint[];
}

/** A stable function that reads the statuses with the portal's own APIs. */
export function useGatusReader() {
  const discovery = useApi(discoveryApiRef);
  const fetchApi = useApi(fetchApiRef);
  return useCallback(() => readGatusStatuses(discovery, fetchApi), [discovery, fetchApi]);
}

/** Reads through the portal proxy (/api/proxy/gatus/statuses, GET only). Refreshes every minute. */
export function useGatusStatus(name: string, refreshMs = 60_000, now: () => number = Date.now): GatusState {
  const readStatuses = useGatusReader();
  const [state, setState] = useState<GatusState>({ loading: true });
  useEffect(() => {
    let active = true;
    const read = async () => {
      let summary: GatusSummary;
      try {
        summary = summariseEndpoint(await readStatuses(), name, now());
      } catch (error) {
        summary = { kind: 'unavailable', reason: error instanceof GatusError ? error.message : 'Gatus is not reachable' };
      }
      if (active) setState({ loading: false, summary, readAt: new Date(now()).toISOString() });
    };
    void read();
    const timer = setInterval(read, refreshMs);
    return () => {
      active = false;
      clearInterval(timer);
    };
  }, [readStatuses, name, refreshMs, now]);
  return state;
}
