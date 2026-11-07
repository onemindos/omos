import type { BusEvent } from './contract';
import { EventSampler } from './sampler';
import { RateMeter } from './rate';

const ev = (n: number): BusEvent => ({
  id: `o-e${n}`,
  subject: 's',
  kind: 's',
  ts: 't',
  bytes: 1,
  summary: String(n),
});

describe('EventSampler', () => {
  it('admits at most 2 per subject per second', () => {
    const s = new EventSampler();
    const t = 1_000_000;
    const admitted = [1, 2, 3, 4, 5].filter(() => s.admit('tak.cot.out', t));
    expect(admitted).toHaveLength(2);
    expect(s.drain().dropped).toBe(3);
  });

  it('admits at most 50 per second overall', () => {
    const s = new EventSampler();
    const t = 5_000_000;
    let ok = 0;
    for (let i = 0; i < 200; i++) if (s.admit(`tel.n${i % 100}.x`, t)) ok++;
    expect(ok).toBe(50);
    expect(s.drain().dropped).toBe(150);
  });

  it('resets in the next second', () => {
    const s = new EventSampler();
    expect(s.admit('a', 1000)).toBe(true);
    expect(s.admit('a', 1100)).toBe(true);
    expect(s.admit('a', 1200)).toBe(false);
    expect(s.admit('a', 2000)).toBe(true);
  });

  it('keeps the last 200 in the ring, newest first, and never sends more than 50 per drain', () => {
    const s = new EventSampler();
    for (let i = 1; i <= 260; i++) s.record(ev(i));
    const recent = s.recent();
    expect(recent).toHaveLength(200);
    expect(recent[0].summary).toBe('260');
    expect(recent[199].summary).toBe('61');
    const drained = s.drain();
    expect(drained.events).toHaveLength(50);
    expect(drained.events[49].summary).toBe('260');
    expect(drained.dropped).toBe(210);
    expect(s.drain()).toEqual({ events: [], dropped: 0 });
  });
});

describe('RateMeter', () => {
  it('averages over 10 seconds and forgets old buckets', () => {
    const r = new RateMeter();
    for (let s = 100; s < 110; s++) r.hit(s * 1000, 20);
    expect(r.rate(109_500)).toBe(20);
    expect(r.rate(125_000)).toBe(0);
  });
});
