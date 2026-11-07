import type { BusEvent } from './contract';

export const MAX_EVENTS_PER_SECOND = 50;
export const MAX_EVENTS_PER_SUBJECT_PER_SECOND = 2;
export const RING_SIZE = 200;

/**
 * Decides which bus messages become BusEvents: at most 50 per second overall and 2 per second
 * per subject (fixed one-second windows). Everything else is counted in `dropped`. The sampler
 * never limits the rate statistics; callers count every message separately.
 */
export class EventSampler {
  private window = -1;
  private total = 0;
  private perSubject = new Map<string, number>();
  private pending: BusEvent[] = [];
  private droppedSinceFlush = 0;
  private readonly ring: BusEvent[] = [];

  /** True when this message may become an event. */
  admit(subject: string, nowMs: number): boolean {
    const w = Math.floor(nowMs / 1000);
    if (w !== this.window) {
      this.window = w;
      this.total = 0;
      this.perSubject.clear();
    }
    if (this.total >= MAX_EVENTS_PER_SECOND) {
      this.droppedSinceFlush++;
      return false;
    }
    const n = this.perSubject.get(subject) ?? 0;
    if (n >= MAX_EVENTS_PER_SUBJECT_PER_SECOND) {
      this.droppedSinceFlush++;
      return false;
    }
    this.perSubject.set(subject, n + 1);
    this.total++;
    return true;
  }

  record(event: BusEvent): void {
    this.pending.push(event);
    this.ring.push(event);
    if (this.ring.length > RING_SIZE)
      this.ring.splice(0, this.ring.length - RING_SIZE);
  }

  /** Events for the next delta (oldest first, at most 50) and the dropped count since the last flush. */
  drain(): { events: BusEvent[]; dropped: number } {
    let events = this.pending;
    let dropped = this.droppedSinceFlush;
    if (events.length > MAX_EVENTS_PER_SECOND) {
      dropped += events.length - MAX_EVENTS_PER_SECOND;
      events = events.slice(events.length - MAX_EVENTS_PER_SECOND);
    }
    this.pending = [];
    this.droppedSinceFlush = 0;
    return { events, dropped };
  }

  /** Newest first, for the snapshot. */
  recent(limit = RING_SIZE): BusEvent[] {
    return this.ring.slice(-limit).reverse();
  }
}
