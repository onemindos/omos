import {
  nowMarker,
  percentAlong,
  pickNextUp,
  PIN_WIDTH_PCT,
  relativeStart,
  remainingTodayEvents,
  timelineBlocks,
  timelineLabel,
  timelineRange,
} from './nextUp';

const now = new Date(2026, 10, 6, 10, 0, 0);
const at = (h: number, m = 0) => new Date(2026, 10, 6, h, m).toISOString();
const event = (id: string, h: number, m = 0, extra: object = {}) => ({
  id,
  title: `Event ${id}`,
  start: at(h, m),
  allDay: false,
  ...extra,
});
const task = (id: string, dueDate: string | null, overdue = true) => ({
  id,
  title: `Task ${id}`,
  dueDate,
  overdue,
});

describe('pickNextUp', () => {
  test('prefers the soonest remaining event today', () => {
    const next = pickNextUp(
      [event('b', 14), event('a', 10, 48)],
      [task('t', '2026-11-01')],
      now,
    );
    expect(next).toMatchObject({ kind: 'event', id: 'a', minutes: 48 });
  });

  test('skips ended events, all-day events and other days', () => {
    const events = [
      event('past', 8),
      event('end', 9, 0, { end: at(9, 45) }),
      event('allday', 0, 0, { allDay: true }),
      {
        ...event('tomorrow', 9),
        start: new Date(2026, 10, 7, 9).toISOString(),
      },
    ];
    expect(remainingTodayEvents(events, now)).toEqual([]);
    expect(pickNextUp(events, [], now)).toEqual({ kind: 'empty' });
  });

  test('an event without an end time is current for 30 minutes', () => {
    expect(pickNextUp([event('now', 9, 40)], [], now)).toMatchObject({
      kind: 'event',
      minutes: -20,
    });
    expect(pickNextUp([event('old', 9, 25)], [], now).kind).toBe('empty');
  });

  test('falls back to the most overdue task', () => {
    const next = pickNextUp(
      [],
      [
        task('recent', '2026-11-04'),
        task('oldest', '2026-10-30'),
        task('open', '2026-11-20', false),
      ],
      now,
    );
    expect(next).toMatchObject({ kind: 'task', id: 'oldest', daysLate: 7 });
  });

  test('empty when nothing is due or scheduled', () => {
    expect(pickNextUp(undefined, undefined, now)).toEqual({ kind: 'empty' });
    expect(pickNextUp([], [task('x', '2026-11-20', false)], now)).toEqual({
      kind: 'empty',
    });
  });
});

test('relativeStart', () => {
  expect(relativeStart(48)).toBe('in 48 min');
  expect(relativeStart(125)).toBe('in 2 h 5 min');
  expect(relativeStart(120)).toBe('in 2 h');
  expect(relativeStart(0)).toBe('now');
  expect(relativeStart(-12)).toBe('started 12 min ago');
});

describe('timeline', () => {
  const range = { startHour: 8, endHour: 20 };
  const d = (h: number, m = 0) => new Date(2026, 10, 6, h, m);

  test('range defaults to 08:00 to 20:00 and widens for outliers', () => {
    expect(timelineRange([])).toEqual(range);
    expect(timelineRange([{ start: d(10), end: d(11) }])).toEqual(range);
    expect(
      timelineRange([
        { start: d(6, 30), end: d(7) },
        { start: d(21), end: d(22, 15) },
      ]),
    ).toEqual({
      startHour: 6,
      endHour: 23,
    });
  });

  test('blocks are positioned in percent of the range', () => {
    const [block] = timelineBlocks(
      [{ id: 'a', title: 'A', start: d(11), end: d(12) }],
      range,
    );
    expect(block.leftPct).toBeCloseTo(25);
    expect(block.widthPct).toBeCloseTo(100 / 12);
    expect(block.hasEnd).toBe(true);
  });

  test('a block with no end is a pin; positions clamp to the track', () => {
    const [pin, early, late] = timelineBlocks(
      [
        { id: 'p', title: 'P', start: d(14) },
        { id: 'e', title: 'E', start: d(6), end: d(9) },
        { id: 'l', title: 'L', start: d(19, 30), end: d(23) },
      ],
      range,
    );
    expect(pin.widthPct).toBe(PIN_WIDTH_PCT);
    expect(pin.hasEnd).toBe(false);
    expect(early.leftPct).toBe(0);
    expect(early.widthPct).toBeCloseTo(100 / 12);
    expect(late.leftPct + late.widthPct).toBeCloseTo(100);
  });

  test('now marker is hidden outside the range', () => {
    expect(nowMarker(d(14), range)).toBeCloseTo(50);
    expect(nowMarker(d(7), range)).toBeUndefined();
    expect(nowMarker(d(21), range)).toBeUndefined();
    expect(percentAlong(d(7), range)).toBe(0);
  });

  test('label describes the events, or says there are none', () => {
    expect(timelineLabel([], range)).toContain('no remaining events');
    const label = timelineLabel(
      [
        { title: 'Standup', start: d(10), end: d(10, 30) },
        { title: 'Review', start: d(15) },
      ],
      range,
    );
    expect(label).toContain('2 remaining events');
    expect(label).toContain('Standup 10:00 to 10:30');
    expect(label).toContain('Review starting 15:00');
  });
});
