import { affiliation } from './cot';
import type { LiveSnapshot } from './contract';

const AFFILIATIONS: Record<string, string> = {
  f: 'friendly',
  h: 'hostile',
  n: 'neutral',
  u: 'unknown',
  a: 'assumed friend',
  s: 'suspect',
  j: 'joker',
  k: 'faker',
  p: 'pending',
  o: 'none specified',
  other: 'other',
};

/** Compact digest of a snapshot for the read-only MCP action. */
export function summarizeSnapshot(snap: LiveSnapshot) {
  const namespaceHealth = { ok: 0, warn: 0, down: 0 };
  for (const c of snap.cluster) namespaceHealth[c.health]++;
  const trackCounts: Record<string, number> = {};
  for (const t of snap.tracks) {
    const key = AFFILIATIONS[affiliation(t.type)] ?? 'other';
    trackCounts[key] = (trackCounts[key] ?? 0) + 1;
  }
  return {
    generatedAt: snap.generatedAt,
    origin: snap.origin,
    status: Object.fromEntries(
      Object.entries(snap.status).map(([k, s]) => [
        k,
        {
          state: s.state,
          rate: s.rate ?? 0,
          lastMessageAt: s.lastMessageAt,
          detail: s.detail,
        },
      ]),
    ),
    namespaceHealth,
    unhealthyNamespaces: snap.cluster
      .filter(c => c.health !== 'ok')
      .map(c => ({
        namespace: c.namespace,
        health: c.health,
        ready: c.ready,
        total: c.total,
        pending: c.pending,
        failed: c.failed,
        restarts: c.restarts,
        lastChange: c.lastChange,
      })),
    argoNotSyncedHealthy: snap.argo
      .filter(a => a.sync !== 'Synced' || a.health !== 'Healthy')
      .map(a => ({
        name: a.name,
        sync: a.sync,
        health: a.health,
        changedAt: a.changedAt,
      })),
    argoTotal: snap.argo.length,
    trackTotal: snap.tracks.length,
    trackCountsByAffiliation: trackCounts,
    busMessagesPerSecond: snap.status.nats.rate ?? 0,
    lastEvents: snap.events.slice(0, 20),
  };
}
