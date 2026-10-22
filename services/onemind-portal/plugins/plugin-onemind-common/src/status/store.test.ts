import { createStatusStore } from './store';
import type { StatusSnapshot } from './model';

const snap = (n: number): StatusSnapshot => ({ observedAt: `t${n}` });

beforeEach(() => jest.useFakeTimers());
afterEach(() => jest.useRealTimers());

test('one poller serves every subscriber and stops with the last one', async () => {
  let n = 0;
  const read = jest.fn(async () => snap(++n));
  const store = createStatusStore({ read }, 60000);
  const a = store.subscribe(() => {});
  const b = store.subscribe(() => {});
  await Promise.resolve();
  expect(read).toHaveBeenCalledTimes(1);
  jest.advanceTimersByTime(60000);
  await Promise.resolve();
  expect(read).toHaveBeenCalledTimes(2);
  a();
  jest.advanceTimersByTime(60000);
  expect(read).toHaveBeenCalledTimes(3);
  b();
  jest.advanceTimersByTime(300000);
  expect(read).toHaveBeenCalledTimes(3);
});

test('a failed read keeps the last snapshot and flags the error', async () => {
  const read = jest.fn().mockResolvedValueOnce(snap(1)).mockRejectedValueOnce(new Error('down'));
  const store = createStatusStore({ read });
  const seen: boolean[] = [];
  const off = store.subscribe(() => seen.push(store.getState().error));
  await Promise.resolve(); await Promise.resolve();
  expect(store.getState().snapshot?.observedAt).toBe('t1');
  await store.refresh();
  expect(store.getState().snapshot?.observedAt).toBe('t1');
  expect(store.getState().error).toBe(true);
  off();
});

test('overlapping refreshes do not stack requests', async () => {
  let release: (s: StatusSnapshot) => void = () => {};
  const read = jest.fn(() => new Promise<StatusSnapshot>(r => { release = r; }));
  const store = createStatusStore({ read });
  void store.refresh(); void store.refresh(); void store.refresh();
  expect(read).toHaveBeenCalledTimes(1);
  release(snap(1));
});
