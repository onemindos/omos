import { Box, Typography } from '@material-ui/core';
import { AppLink } from './HomeCard';
import { useCalendarSummary } from './CalendarPanel';
import { useCloudWork } from './cloudWork';
import { DayTimeline } from './DayTimeline';
import { Eyebrow, Surface, useNow } from './bento';
import { useHomeTokens } from './homeTokens';
import {
  hhmm,
  pickNextUp,
  relativeStart,
  timelineRange,
  todayTimedEvents,
} from './nextUp';
import { calendarStatus, hasData, tasksStatus } from './sources';

const shortDate = (d: Date) =>
  d.toLocaleDateString(undefined, { month: 'short', day: 'numeric' });

/** Row A, wide: the single most important upcoming thing, with today's timeline. */
export function NextUpHero() {
  const t = useHomeTokens();
  const now = useNow(30000);
  const calendar = useCalendarSummary();
  const cloud = useCloudWork();
  const calStatus = calendarStatus(
    calendar.data,
    calendar.error,
    now.getTime(),
  );
  const workStatus = tasksStatus(cloud.data, cloud.error, now.getTime());
  const calOk = hasData(calStatus);
  const workOk = hasData(workStatus);

  const events = calOk ? calendar.data?.events : undefined;
  const tasks = workOk ? cloud.data?.tasks : undefined;
  const next = pickNextUp(events, tasks, now);
  const today = todayTimedEvents(events, now);
  const range = timelineRange(today.filter(item => !item.ended));
  const items = today.map(item => ({
    id: item.event.id,
    title: item.event.title,
    start: item.start,
    end: item.end,
    ended: item.ended,
  }));
  const allDay = calOk
    ? (calendar.data?.events ?? []).filter(e => e.allDay).slice(0, 3)
    : [];

  const loading =
    calStatus.state === 'loading' || workStatus.state === 'loading';
  const notes = [calStatus, workStatus].filter(
    s => s.message && s.state !== 'loading',
  );

  let emptyTitle = 'Nothing scheduled — your day is open';
  if (loading) emptyTitle = 'Loading your day…';
  else if (notes.length) emptyTitle = 'Your day is not fully visible';

  let eyebrow = 'Next up';
  let eyebrowColor: string | undefined;
  let title: string | undefined;
  let secondary: string | undefined;
  let href: string | undefined;
  if (next.kind === 'event') {
    eyebrow = `Next up · ${relativeStart(next.minutes)}`;
    eyebrowColor = t.status.info.fg;
    title = next.title;
    secondary = `${hhmm(next.start)}${
      next.end ? ` to ${hhmm(next.end)}` : ''
    } · Calendar`;
    href = calendar.data?.link;
  } else if (next.kind === 'task') {
    eyebrow = next.daysLate
      ? `Overdue · ${next.daysLate} ${
          next.daysLate === 1 ? 'day' : 'days'
        } late`
      : 'Overdue';
    eyebrowColor = t.status.critical.fg;
    title = next.title;
    secondary = `${
      next.due ? `Due ${shortDate(next.due)} · ` : ''
    }Zoho Projects`;
    href = cloud.data?.link;
  }

  return (
    <Surface label="Next up" hero>
      <Eyebrow color={eyebrowColor}>{eyebrow}</Eyebrow>
      {title ? (
        <Box>
          <Typography
            variant="h4"
            component="h2"
            style={{
              margin: 0,
              fontSize: 28,
              lineHeight: 1.2,
              fontWeight: 700,
              overflowWrap: 'anywhere',
              color: t.text[1],
            }}
          >
            {href ? (
              <AppLink href={href} color="inherit" underline="none">
                {title}
              </AppLink>
            ) : (
              title
            )}
          </Typography>
          <Typography
            variant="body1"
            style={{
              color: t.text[2],
              marginTop: 4,
              fontVariantNumeric: 'tabular-nums',
            }}
          >
            {secondary}
          </Typography>
        </Box>
      ) : (
        <Typography
          variant="h5"
          component="h2"
          style={{ margin: 0, fontWeight: 600, color: t.text[2] }}
        >
          {emptyTitle}
        </Typography>
      )}
      {calOk && (
        <DayTimeline
          items={items}
          range={range}
          now={now}
          highlightId={next.kind === 'event' ? next.id : undefined}
        />
      )}
      {calOk && !!allDay.length && (
        <Typography variant="caption" style={{ color: t.text[2] }}>
          All day: {allDay.map(e => e.title).join(', ')}
        </Typography>
      )}
      {calOk && items.some(item => !item.end) && (
        <Typography variant="caption" style={{ color: t.text[3] }}>
          Calendar gives start times only; blocks mark when each event starts.
        </Typography>
      )}
      {notes.map(note => (
        <Typography
          key={note.message}
          variant="caption"
          role={note.state === 'unavailable' ? 'alert' : undefined}
          style={{
            color:
              note.state === 'unavailable' || note.state === 'stale'
                ? t.status.warn.fg
                : t.text[2],
          }}
        >
          {note.message}
        </Typography>
      ))}
    </Surface>
  );
}
