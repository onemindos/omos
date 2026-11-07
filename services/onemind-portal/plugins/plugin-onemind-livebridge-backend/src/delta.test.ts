import { LiveBridge } from './bridge';
import { ClusterState } from './cluster';
import { LIVE_CHANNELS, type Track } from './contract';
import { MAX_MESSAGE_BYTES, splitDelta } from './delta';
import { TrackStore } from './trackStore';

const env = { origin: 'o', seq: 0, ts: '2026-11-07T00:00:00Z' };

describe('splitDelta', () => {
  it('returns nothing for an empty delta', () => {
    expect(splitDelta(env, { upserts: [], removes: [] })).toEqual([]);
  });

  it('keeps a small delta in one message', () => {
    const out = splitDelta(env, { upserts: [{ a: 1 }], removes: ['x'] });
    expect(out).toHaveLength(1);
    expect(out[0]).toMatchObject({ upserts: [{ a: 1 }], removes: ['x'] });
  });

  it('splits a large delta below the limit, loses nothing and never sends an empty message', () => {
    const upserts = Array.from({ length: 2000 }, (_, i) => ({
      uid: `uid-${i}`,
      pad: 'x'.repeat(100),
    }));
    const removes = Array.from({ length: 3000 }, (_, i) => `removed-uid-${i}`);
    const out = splitDelta(env, { upserts, removes });
    expect(out.length).toBeGreaterThan(2);
    let u = 0;
    let r = 0;
    for (const m of out) {
      expect(Buffer.byteLength(JSON.stringify(m))).toBeLessThan(60 * 1024);
      expect(Buffer.byteLength(JSON.stringify(m))).toBeLessThanOrEqual(
        MAX_MESSAGE_BYTES,
      );
      expect(m.upserts.length + m.removes.length).toBeGreaterThan(0);
      u += m.upserts.length;
      r += m.removes.length;
    }
    expect([u, r]).toEqual([2000, 3000]);
  });

  it('drops an item that can never fit', () => {
    const out = splitDelta(
      env,
      { items: [{ big: 'x'.repeat(100) }, { ok: 1 }] },
      120,
    );
    expect(out.flatMap(m => m.items)).toEqual([{ ok: 1 }]);
  });
});

describe('TrackStore', () => {
  const track = (uid: string, stale: string): Track => ({
    uid,
    type: 'a-f-G',
    lat: 1,
    lon: 2,
    time: '2026-11-07T10:00:00Z',
    stale,
    subject: 's',
  });
  const now = Date.parse('2026-11-07T10:00:00Z');

  it('caps at the maximum and evicts the least recently updated', () => {
    const s = new TrackStore(3);
    ['a', 'b', 'c'].forEach(u =>
      s.upsert(track(u, '2026-11-07T10:10:00Z'), now),
    );
    s.upsert(track('a', '2026-11-07T10:10:00Z'), now);
    s.upsert(track('d', '2026-11-07T10:10:00Z'), now);
    expect(
      s
        .all()
        .map(t => t.uid)
        .sort(),
    ).toEqual(['a', 'c', 'd']);
    expect(s.drain().removes).toEqual(['b']);
  });

  it('keeps tracks at least 2 minutes and at most 30 minutes', () => {
    const s = new TrackStore();
    s.upsert(track('short', '2026-11-07T10:00:05Z'), now);
    s.upsert(track('long', '2026-11-08T10:00:00Z'), now);
    s.sweep(now + 100_000);
    expect(
      s
        .all()
        .map(t => t.uid)
        .sort(),
    ).toEqual(['long', 'short']);
    s.sweep(now + 130_000);
    expect(s.all().map(t => t.uid)).toEqual(['long']);
    s.sweep(now + 31 * 60_000);
    expect(s.size).toBe(0);
    expect(s.drain().removes.sort()).toEqual(['long', 'short']);
  });
});

describe('LiveBridge', () => {
  it('publishes bounded deltas with increasing seq and never an empty one', () => {
    let t = Date.parse('2026-11-07T10:00:00Z');
    const sent: Array<{ channel: string; message: any }> = [];
    const bridge = new LiveBridge({
      origin: 'host-abc',
      logger: { warn() {}, info() {}, debug() {}, error() {} } as never,
      publish: (channel, message) => void sent.push({ channel, message }),
      cluster: new ClusterState(),
      providers: {
        nats: () => ({ state: 'connected' }),
        kubernetes: () => ({ state: 'disabled' }),
        argo: () => ({ state: 'disabled' }),
      },
      now: () => t,
    });
    bridge.flush();
    expect(sent).toHaveLength(0);

    const enc = new TextEncoder();
    const cot = (uid: string) =>
      enc.encode(
        `<event uid="${uid}" type="a-n-A-C-F" time="2026-11-07T10:00:00Z" stale="2026-11-07T10:01:00Z"><point lat="1" lon="2"/><detail><contact callsign="${uid}"/></detail></event>`,
      );
    bridge.ingest(
      'intel.conflict.acled.event',
      enc.encode('{"title":"Clash near X","n":1}'),
    );
    for (let i = 0; i < 300; i++)
      bridge.ingest(`tak.cot.out.${i}`, cot(`u${i}`));
    t += 1000;
    bridge.flush();

    const tracks = sent.filter(s => s.channel === LIVE_CHANNELS.tracks);
    const events = sent.filter(s => s.channel === LIVE_CHANNELS.events);
    expect(tracks.flatMap(s => s.message.upserts)).toHaveLength(300);
    expect(events).toHaveLength(1);
    expect(events[0].message.events.length).toBeLessThanOrEqual(50);
    expect(events[0].message.dropped).toBeGreaterThan(200);
    expect(events[0].message.events[0].summary).toBe('Clash near X');
    const seqs = sent.map(s => s.message.seq);
    expect(seqs).toEqual([...seqs].sort((a, b) => a - b));
    expect(new Set(seqs).size).toBe(seqs.length);
    for (const s of sent)
      expect(Buffer.byteLength(JSON.stringify(s.message))).toBeLessThan(
        60 * 1024,
      );
    expect(bridge.snapshot().events.length).toBeLessThanOrEqual(200);
    expect(bridge.snapshot().tracks).toHaveLength(300);
    expect(bridge.busRate()).toBeGreaterThan(0);
  });
});
