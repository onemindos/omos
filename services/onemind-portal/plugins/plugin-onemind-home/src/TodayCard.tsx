import { Box, Typography, useTheme } from '@material-ui/core';
import { AppLink, HomeCard } from './HomeCard';
import { useCalendarSummary } from './CalendarPanel';
import { useCloudWork } from './cloudWork';
import { mergeToday } from './todayItems';

const observed = (iso: string | null | undefined) =>
  iso
    ? new Date(iso).toLocaleTimeString([], {
        hour: '2-digit',
        minute: '2-digit',
        hour12: false,
      })
    : undefined;

export function TodayCard() {
  const theme = useTheme();
  const calendar = useCalendarSummary();
  const cloud = useCloudWork();
  const items = mergeToday(calendar.data?.events, cloud.data?.tasks);
  const notes: { text: string; error?: boolean }[] = [];

  const cal = calendar.data;
  const calOk =
    !!cal &&
    cal.state !== 'configuration-required' &&
    cal.state !== 'unavailable';
  if (calendar.error)
    notes.push({ text: `Calendar: ${calendar.error}.`, error: true });
  else if (!cal) notes.push({ text: 'Loading calendar…' });
  else if (cal.state === 'configuration-required')
    notes.push({
      text: 'Connect read-only Calendar access and select a calendar UID to enable upcoming events.',
    });
  else if (cal.state === 'unavailable')
    notes.push({
      text: 'Calendar source unavailable; an empty list does not mean no events.',
      error: true,
    });
  const calStale =
    !!cal &&
    (!!calendar.error ||
      cal.state === 'stale' ||
      (!!cal.expiresAt && Date.parse(cal.expiresAt) < Date.now()));

  const work = cloud.data;
  const workOk =
    !!work &&
    work.state !== 'configuration-required' &&
    work.state !== 'unavailable';
  if (cloud.error)
    notes.push({
      text: 'Tasks: summary unavailable; previous observations may be stale.',
      error: true,
    });
  else if (!work) notes.push({ text: 'Loading tasks…' });
  else if (work.state === 'configuration-required')
    notes.push({
      text: 'Connect read-only Zoho access to enable task summaries.',
    });
  else if (work.state === 'unavailable')
    notes.push({ text: 'Tasks: provider unavailable.', error: true });
  const workStale =
    !!work &&
    (!!cloud.error ||
      work.state === 'stale' ||
      (!!work.expiresAt && Date.parse(work.expiresAt) < Date.now()));

  const allLoaded = (!!cal || !!calendar.error) && (!!work || !!cloud.error);
  const footer = [
    calOk && cal?.observedAt
      ? `${calStale ? 'Stale · ' : ''}Calendar ${observed(cal.observedAt)}`
      : undefined,
    workOk && work?.observedAt
      ? `${workStale ? 'Stale · ' : ''}Tasks ${observed(work.observedAt)}${
          work.partial ? ' (sample)' : ''
        }`
      : undefined,
  ].filter(Boolean);

  return (
    <HomeCard
      label="Today"
      title="Today"
      action={
        <AppLink href="/today" variant="caption">
          See all
        </AppLink>
      }
    >
      {items.length > 0 && (
        <Box
          component="ul"
          style={{ listStyle: 'none', margin: 0, padding: 0 }}
        >
          {items.map(item => (
            <Box
              component="li"
              key={item.id}
              display="flex"
              alignItems="baseline"
              style={{ gap: 8, padding: theme.spacing(0.5, 0) }}
            >
              <Typography
                variant="caption"
                style={{
                  width: 76,
                  flexShrink: 0,
                  fontWeight: 600,
                  fontVariantNumeric: 'tabular-nums',
                  color:
                    item.tone === 'late'
                      ? theme.palette.error.main
                      : theme.palette.text.secondary,
                }}
              >
                {item.label}
              </Typography>
              <Typography
                variant="body2"
                style={{ flex: 1, minWidth: 0, overflowWrap: 'anywhere' }}
              >
                {item.title}
              </Typography>
              <Typography
                variant="caption"
                color="textSecondary"
                style={{ flexShrink: 0 }}
              >
                {item.source}
              </Typography>
            </Box>
          ))}
        </Box>
      )}
      {!items.length && allLoaded && !notes.some(n => n.error) && (
        <Typography variant="body2" color="textSecondary">
          Nothing due or scheduled in the observed sources.
        </Typography>
      )}
      {notes.map(note => (
        <Typography
          key={note.text}
          variant="caption"
          display="block"
          role={note.error ? 'alert' : undefined}
          color={note.error ? 'error' : 'textSecondary'}
        >
          {note.text}
        </Typography>
      ))}
      {!!footer.length && (
        <Typography
          variant="caption"
          color={calStale || workStale ? 'error' : 'textSecondary'}
          display="block"
          style={{ marginTop: 6 }}
        >
          Observed: {footer.join(' · ')}
        </Typography>
      )}
    </HomeCard>
  );
}
