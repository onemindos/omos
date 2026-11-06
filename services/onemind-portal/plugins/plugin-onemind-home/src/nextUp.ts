import type { CalendarSummary } from './CalendarPanel';
import type { CloudSummary } from './cloudWork';
import { parseWhen } from './todayItems';

type CalendarEvent = CalendarSummary['events'][number];
type Task = NonNullable<CloudSummary['tasks']>[number];

/**
 * The calendar summary carries a start time only. An event with no end time is
 * treated as running for this long when deciding whether it is still current.
 */
export const ASSUMED_EVENT_MINUTES = 30;
const MIN = 60000;

const startOfDay = (d: Date) =>
  new Date(d.getFullYear(), d.getMonth(), d.getDate()).getTime();

export type NextUp =
  | {
      kind: 'event';
      id: string;
      title: string;
      start: Date;
      /** Present only when the source sent an end time. */
      end?: Date;
      /** Minutes until the start; negative once it has started. */
      minutes: number;
    }
  | {
      kind: 'task';
      id: string;
      title: string;
      due?: Date;
      daysLate?: number;
    }
  | { kind: 'empty' };

function eventEnd(event: CalendarEvent): Date | undefined {
  return parseWhen((event as { end?: string | null }).end);
}

/** All timed events that start today, soonest first, each flagged when it has ended. */
export function todayTimedEvents(
  events: CalendarEvent[] | undefined,
  now = new Date(),
) {
  const day = startOfDay(now);
  return (events ?? [])
    .flatMap(event => {
      const start = parseWhen(event.start);
      if (!start || event.allDay) return [];
      if (startOfDay(start) !== day) return [];
      const end = eventEnd(event);
      const endMs = end
        ? end.getTime()
        : start.getTime() + ASSUMED_EVENT_MINUTES * MIN;
      return [{ event, start, end, ended: endMs <= now.getTime() }];
    })
    .sort((a, b) => a.start.getTime() - b.start.getTime());
}

/** Timed events that start today and have not ended yet, soonest first. */
export const remainingTodayEvents = (
  events: CalendarEvent[] | undefined,
  now = new Date(),
) => todayTimedEvents(events, now).filter(item => !item.ended);

/**
 * The one thing to show big: the next calendar event today that has not ended,
 * else the most overdue open task (oldest due date first), else empty.
 */
export function pickNextUp(
  events: CalendarEvent[] | undefined,
  tasks: Task[] | undefined,
  now = new Date(),
): NextUp {
  const next = remainingTodayEvents(events, now)[0];
  if (next)
    return {
      kind: 'event',
      id: next.event.id,
      title: next.event.title,
      start: next.start,
      end: next.end,
      minutes: Math.round((next.start.getTime() - now.getTime()) / MIN),
    };
  const day = startOfDay(now);
  const overdue = (tasks ?? [])
    .map(task => ({ task, due: parseWhen(task.dueDate) }))
    .filter(({ task, due }) => task.overdue || (!!due && startOfDay(due) < day))
    .sort(
      (a, b) =>
        (a.due?.getTime() ?? Number.MAX_SAFE_INTEGER) -
        (b.due?.getTime() ?? Number.MAX_SAFE_INTEGER),
    )[0];
  if (overdue)
    return {
      kind: 'task',
      id: overdue.task.id,
      title: overdue.task.title,
      due: overdue.due,
      daysLate: overdue.due
        ? Math.max(1, Math.round((day - startOfDay(overdue.due)) / 86400000))
        : undefined,
    };
  return { kind: 'empty' };
}

/** "in 48 min", "in 2 h 5 min", "now", "started 12 min ago". */
export function relativeStart(minutes: number): string {
  if (minutes <= -1) return `started ${Math.abs(minutes)} min ago`;
  if (minutes < 1) return 'now';
  if (minutes < 60) return `in ${minutes} min`;
  const h = Math.floor(minutes / 60);
  const m = minutes % 60;
  return m ? `in ${h} h ${m} min` : `in ${h} h`;
}

export const DEFAULT_RANGE = { startHour: 8, endHour: 20 } as const;

export type TimelineRange = { startHour: number; endHour: number };

/** 08:00 to 20:00, widened to whole hours when a remaining event falls outside. */
export function timelineRange(
  events: { start: Date; end?: Date }[],
): TimelineRange {
  let startHour: number = DEFAULT_RANGE.startHour;
  let endHour: number = DEFAULT_RANGE.endHour;
  for (const { start, end } of events) {
    startHour = Math.min(startHour, start.getHours());
    const last = end ?? new Date(start.getTime() + ASSUMED_EVENT_MINUTES * MIN);
    const sameDay = startOfDay(last) === startOfDay(start);
    endHour = Math.max(
      endHour,
      sameDay ? Math.ceil(last.getHours() + last.getMinutes() / 60) : 24,
    );
  }
  return { startHour: Math.max(0, startHour), endHour: Math.min(24, endHour) };
}

/** An event with no known end is drawn as a thin pin this wide (percent). */
export const PIN_WIDTH_PCT = 1.5;

export type TimelineBlock = {
  id: string;
  title: string;
  leftPct: number;
  widthPct: number;
  /** False when the source sent no end time: the block marks the start only. */
  hasEnd: boolean;
  startLabel: string;
};

const pad = (n: number) => String(n).padStart(2, '0');
export const hhmm = (d: Date) => `${pad(d.getHours())}:${pad(d.getMinutes())}`;

const clamp = (n: number, lo: number, hi: number) =>
  Math.max(lo, Math.min(hi, n));

/** Percent position of an instant along today's track, clamped to 0..100. */
export function percentAlong(date: Date, range: TimelineRange) {
  const span = (range.endHour - range.startHour) * 60;
  const minutes =
    date.getHours() * 60 + date.getMinutes() - range.startHour * 60;
  return clamp((minutes / span) * 100, 0, 100);
}

export function timelineBlocks(
  items: { id: string; title: string; start: Date; end?: Date }[],
  range: TimelineRange,
): TimelineBlock[] {
  return items.map(item => {
    const left = percentAlong(item.start, range);
    const sameDay =
      !!item.end && startOfDay(item.end) === startOfDay(item.start);
    let right = left + PIN_WIDTH_PCT;
    if (item.end) right = sameDay ? percentAlong(item.end, range) : 100;
    const widthPct = clamp(
      Math.max(right - left, PIN_WIDTH_PCT),
      0,
      100 - left,
    );
    return {
      id: item.id,
      title: item.title,
      leftPct: left,
      widthPct,
      hasEnd: !!item.end,
      startLabel: hhmm(item.start),
    };
  });
}

/** Position of the "now" marker, or undefined when now is outside the range. */
export function nowMarker(now: Date, range: TimelineRange) {
  const minutes = now.getHours() * 60 + now.getMinutes();
  if (minutes < range.startHour * 60 || minutes > range.endHour * 60)
    return undefined;
  return percentAlong(now, range);
}

export function timelineLabel(
  items: { title: string; start: Date; end?: Date }[],
  range: TimelineRange,
) {
  const span = `${pad(range.startHour)}:00 to ${pad(range.endHour)}:00`;
  if (!items.length) return `Today ${span}: no remaining events.`;
  const list = items
    .map(i =>
      i.end
        ? `${i.title} ${hhmm(i.start)} to ${hhmm(i.end)}`
        : `${i.title} starting ${hhmm(i.start)}`,
    )
    .join('; ');
  return `Today ${span}: ${items.length} remaining ${
    items.length === 1 ? 'event' : 'events'
  }. ${list}.`;
}
