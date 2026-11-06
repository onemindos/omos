import { mergeToday, parseWhen } from './todayItems';

const now = new Date(2026, 10, 6, 10, 0, 0);
const iso = (h: number, m = 0, dayOffset = 0) =>
  new Date(2026, 10, 6 + dayOffset, h, m).toISOString();

test('parses Zoho and ISO date shapes, rejects junk', () => {
  expect(parseWhen('11-06-2026')?.getDate()).toBe(6);
  expect(parseWhen('2026-11-06')?.getMonth()).toBe(10);
  expect(parseWhen('nonsense')).toBeUndefined();
  expect(parseWhen(null)).toBeUndefined();
});

test('orders overdue tasks, then today events, tasks due today, then upcoming', () => {
  const items = mergeToday(
    [
      { id: 'e2', title: 'Later today', start: iso(15), allDay: false },
      { id: 'e1', title: 'Soon', start: iso(11), allDay: false },
      { id: 'e3', title: 'Tomorrow', start: iso(9, 0, 1), allDay: false },
      { id: 'old', title: 'Already over', start: iso(7), allDay: false },
    ],
    [
      { id: 't1', title: 'Due today', dueDate: '2026-11-06', overdue: false },
      { id: 't2', title: 'Late one', dueDate: '2026-11-01', overdue: true },
      { id: 't3', title: 'Next week', dueDate: '2026-11-12', overdue: false },
      { id: 't4', title: 'No date', dueDate: null, overdue: false },
    ],
    now,
    20,
  );
  expect(items.map(i => i.title)).toEqual([
    'Late one',
    'Soon',
    'Later today',
    'Due today',
    'Tomorrow',
    'Next week',
    'No date',
  ]);
  expect(items[0]).toMatchObject({ label: 'Late', tone: 'late' });
  expect(items[3].label).toBe('Due');
});

test('caps the list and treats missing sources as empty', () => {
  const tasks = Array.from({ length: 10 }, (_, i) => ({
    id: String(i),
    title: `T${i}`,
    dueDate: '2026-10-01',
    overdue: true,
  }));
  expect(mergeToday(undefined, tasks, now)).toHaveLength(6);
  expect(mergeToday(undefined, undefined, now)).toEqual([]);
});
