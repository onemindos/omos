import { collectResources, toChange } from './resources';
import type { ChangeSet, ResourceType, Snapshots } from './types';

/**
 * Resources added or removed between two refreshes, by id. A type is compared only when its
 * section read successfully both times, so an outage never shows up as a mass removal.
 */
export function diffSnapshots(
  previous: Snapshots,
  current: Snapshots,
  at: string,
): ChangeSet | undefined {
  const before = collectResources(previous);
  const after = collectResources(current);
  const comparable = new Set<ResourceType>(
    [...before.complete].filter(t => after.complete.has(t)),
  );
  const key = (r: { type: string; id: string }) => `${r.type}:${r.id}`;
  const was = new Map(
    before.resources.filter(r => comparable.has(r.type)).map(r => [key(r), r]),
  );
  const now = new Map(
    after.resources.filter(r => comparable.has(r.type)).map(r => [key(r), r]),
  );
  const added = [...now]
    .filter(([k]) => !was.has(k))
    .map(([, r]) => toChange(r));
  const removed = [...was]
    .filter(([k]) => !now.has(k))
    .map(([, r]) => toChange(r));
  return added.length || removed.length ? { at, added, removed } : undefined;
}
