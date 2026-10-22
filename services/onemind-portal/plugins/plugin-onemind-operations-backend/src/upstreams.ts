// The real upstream readers. Fixed in-cluster URLs only: nothing here accepts a URL,
// query or path from a caller.
import { health, list } from '@internal/plugin-onemind-common-backend';
import { AmAlert } from './build';
import { CollectorDeps } from './collect';
import { parseVector } from './liveQueries';

const PROMETHEUS = 'http://monitoring-kube-prometheus-prometheus.monitoring.svc:9090';
const ALERTMANAGER = 'http://monitoring-kube-prometheus-alertmanager.monitoring.svc:9093';
const MAX_BYTES = 8_000_000;

async function getJson(url: string, signal: AbortSignal): Promise<unknown> {
  const response = await fetch(url, { signal, redirect: 'error' });
  if (!response.ok) throw new Error(`Upstream returned ${response.status}`);
  const text = await response.text();
  if (text.length > MAX_BYTES) throw new Error('Upstream response too large');
  return JSON.parse(text);
}

export const realDeps: CollectorDeps = {
  list,
  prometheus: async (query, signal) =>
    parseVector(await getJson(`${PROMETHEUS}/api/v1/query?query=${encodeURIComponent(query)}`, signal)),
  alertmanager: async signal => {
    const body = await getJson(
      `${ALERTMANAGER}/api/v2/alerts?active=true&silenced=false&inhibited=false`,
      signal,
    );
    if (!Array.isArray(body)) throw new Error('Invalid alert response');
    return body as AmAlert[];
  },
  uptime: () => health() as never,
};
