import type { CalendarSummary } from './CalendarPanel';
import type { CloudSummary } from './cloudWork';

export type TodayItem = {
  id: string;
  kind: 'event' | 'task';
  title: string;
  /** Short left-hand label: a time, "All day", "Late", "Due" or a date. */
  label: string;
  source: 'Calendar' | 'Projects';
  bucket: 'overdue' | 'today' | 'upcoming';
  tone: 'late' | 'normal';
  sortKey: number;
  /** Lower ranks sort first: overdue 0, today's events 1, tasks due today 2, later 3. */
  rank: number;
};

const pad = (n: number) => String(n).padStart(2, '0');
const dayKey = (d: Date) =>
  `${d.getFullYear()}-${pad(d.getMonth() + 1)}-${pad(d.getDate())}`;

/** Parses ISO, `YYYY-MM-DD` and Zoho's `MM-DD-YYYY` into a local date; undefined if unknown. */
export function parseWhen(value?: string | null): Date | undefined {
  if (!value) return undefined;
  const text = value.trim();
  let m = /^(\d{4})-(\d{2})-(\d{2})$/.exec(text);
  if (m) return new Date(+m[1], +m[2] - 1, +m[3]);
  m = /^(\d{2})-(\d{2})-(\d{4})$/.exec(text);
  if (m) return new Date(+m[3], +m[1] - 1, +m[2]);
  const t = Date.parse(text);
  return Number.isNaN(t) ? undefined : new Date(t);
}

const timeOf = (d: Date) => `${pad(d.getHours())}:${pad(d.getMinutes())}`;
const shortDate = (d: Date) =>
  d.toLocaleDateString(undefined, { month: 'short', day: 'numeric' });

function eventLabel(start: Date, allDay: boolean, startsToday: boolean) {
  if (allDay) return startsToday ? 'All day' : `${shortDate(start)} · all day`;
  return startsToday ? timeOf(start) : `${shortDate(start)} ${timeOf(start)}`;
}

/**
 * Merges Zoho Calendar events and open Zoho Projects tasks into one short
 * timeline: overdue tasks first, then today's remaining events, tasks due today,
 * then the next few upcoming items. Events that started more than 30 minutes
 * ago are dropped. Never invents data: sources that failed contribute nothing.
 */
export function mergeToday(
  events: CalendarSummary['events'] | undefined,
  tasks: NonNullable<CloudSummary['tasks']> | undefined,
  now = new Date(),
  cap = 6,
): TodayItem[] {
  const today = dayKey(now);
  const items: TodayItem[] = [];
  for (const task of tasks ?? []) {
    const due = parseWhen(task.dueDate);
    const dueToday = !!due && dayKey(due) === today;
    const overdue = task.overdue || (!!due && !dueToday && due < now);
    const base = {
      id: `task:${task.id}`,
      kind: 'task' as const,
      title: task.title,
      source: 'Projects' as const,
    };
    if (overdue) {
      items.push({
        ...base,
        label: 'Late',
        bucket: 'overdue',
        tone: 'late',
        rank: 0,
        sortKey: due ? due.getTime() : 0,
      });
    } else if (dueToday) {
      items.push({
        ...base,
        label: 'Due',
        bucket: 'today',
        tone: 'normal',
        rank: 2,
        sortKey: due!.getTime(),
      });
    } else {
      items.push({
        ...base,
        label: due ? shortDate(due) : 'Open',
        bucket: 'upcoming',
        tone: 'normal',
        rank: 3,
        sortKey: due ? due.getTime() : Number.MAX_SAFE_INTEGER,
      });
    }
  }
  for (const event of events ?? []) {
    const start = parseWhen(event.start);
    if (!start) continue;
    const startsToday = dayKey(start) === today;
    if (!event.allDay && start.getTime() < now.getTime() - 30 * 60000) continue;
    if (
      event.allDay &&
      start < new Date(now.getFullYear(), now.getMonth(), now.getDate())
    )
      continue;
    items.push({
      id: `event:${event.id}`,
      kind: 'event',
      title: event.title,
      source: 'Calendar',
      label: eventLabel(start, event.allDay, startsToday),
      bucket: startsToday ? 'today' : 'upcoming',
      tone: 'normal',
      rank: startsToday ? 1 : 3,
      sortKey: start.getTime(),
    });
  }
  return items
    .sort((a, b) => a.rank - b.rank || a.sortKey - b.sortKey)
    .slice(0, cap);
}
