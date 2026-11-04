import { useMemo, useSyncExternalStore } from 'react';
import { useApi } from '@backstage/core-plugin-api';
import { useEntity } from '@backstage/plugin-catalog-react';
import { combineStatus, targetForEntity, AppStatusResult, EntityLike } from './model';
import { appStatusApiRef } from './source';
import { createStatusStore, StatusStore } from './store';

// One store per source instance (the app has one source; tests pass their own).
const stores = new WeakMap<object, StatusStore>();
function storeFor(source: object, make: () => StatusStore) {
  let store = stores.get(source);
  if (!store) {
    store = make();
    stores.set(source, store);
  }
  return store;
}

/** Shared snapshot state for the page. */
export function useStatusSnapshot() {
  const source = useApi(appStatusApiRef);
  const store = storeFor(source, () => createStatusStore(source));
  const state = useSyncExternalStore(store.subscribe, store.getState);
  return { ...state, refresh: store.refresh };
}

/** Status of one entity. */
export function useAppStatus(entity: EntityLike): AppStatusResult & { loading: boolean; error: boolean } {
  const { snapshot, error } = useStatusSnapshot();
  const target = useMemo(() => targetForEntity(entity), [entity]);
  const result = useMemo(() => combineStatus(target, snapshot), [target, snapshot]);
  return { ...result, loading: !snapshot && !error, error };
}

/** For the entity page: reads `useEntity()`. */
export function useEntityStatus() {
  const { entity } = useEntity();
  return useAppStatus(entity as EntityLike);
}

/** Statuses of several entities from the same shared snapshot (System and Domain summaries). */
export function useAppStatuses(entities: EntityLike[]) {
  const { snapshot, error } = useStatusSnapshot();
  const results = useMemo(
    () => entities.map(entity => ({ entity, result: combineStatus(targetForEntity(entity), snapshot) })),
    [entities, snapshot],
  );
  return { results, loading: !snapshot && !error, error };
}
