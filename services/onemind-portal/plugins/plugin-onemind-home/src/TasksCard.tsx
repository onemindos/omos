import { Box, Typography } from '@material-ui/core';
import { AppLink } from './HomeCard';
import { Eyebrow, Surface } from './bento';
import { useCloudWork } from './cloudWork';
import { useHomeTokens } from './homeTokens';
import { hasData, tasksStatus } from './sources';
import { mergeToday, TodayItem } from './todayItems';

/** Overdue share of open tasks as a percentage; undefined without both counts. */
export function overduePercent(open?: number, overdue?: number) {
  if (typeof open !== 'number' || typeof overdue !== 'number' || open <= 0)
    return undefined;
  return Math.max(0, Math.min(100, Math.round((overdue / open) * 100)));
}

export const topTasks = (
  tasks: Parameters<typeof mergeToday>[1],
  now = new Date(),
): TodayItem[] => mergeToday(undefined, tasks, now, 3);

/** Row B: open and overdue Zoho Projects tasks. The API returns open tasks only. */
export function TasksCard() {
  const t = useHomeTokens();
  const { data, error } = useCloudWork();
  const status = tasksStatus(data, error);
  const ok = hasData(status);
  const open = ok ? data?.open : undefined;
  const overdue = ok ? data?.overdue : undefined;
  const pct = overduePercent(open, overdue);
  const items = ok ? topTasks(data?.tasks) : [];
  const stale = status.state === 'stale';
  return (
    <Surface label="Tasks">
      <Eyebrow>Tasks · Zoho Projects</Eyebrow>
      {ok ? (
        <Typography
          component="p"
          style={{
            margin: 0,
            fontSize: 28,
            fontWeight: 700,
            lineHeight: 1.15,
            fontVariantNumeric: 'tabular-nums',
            color: stale ? t.status.warn.fg : t.text[1],
          }}
        >
          {typeof open === 'number' ? open : '—'} open
          <span style={{ color: t.text[3] }}> · </span>
          <span
            style={{
              color: overdue ? t.status.critical.fg : t.text[1],
            }}
          >
            {typeof overdue === 'number' ? overdue : '—'} overdue
          </span>
        </Typography>
      ) : (
        <Typography
          component="p"
          style={{ margin: 0, fontSize: 28, fontWeight: 700, color: t.text[3] }}
        >
          —
        </Typography>
      )}
      {pct !== undefined && (
        <div
          role="img"
          aria-label={`${overdue} of ${open} open tasks are overdue (${pct} percent)`}
          style={{
            display: 'flex',
            height: 8,
            borderRadius: 4,
            overflow: 'hidden',
            background: t.surface[3],
          }}
        >
          <span
            style={{ width: `${pct}%`, background: t.status.critical.fg }}
          />
        </div>
      )}
      {!!items.length && (
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
              style={{ gap: 8, padding: '4px 0' }}
            >
              <Typography
                variant="caption"
                style={{
                  width: 52,
                  flexShrink: 0,
                  fontWeight: 700,
                  fontVariantNumeric: 'tabular-nums',
                  color:
                    item.tone === 'late' ? t.status.critical.fg : t.text[2],
                }}
              >
                {item.tone === 'late' ? 'late' : item.label}
              </Typography>
              <Typography
                variant="body2"
                style={{
                  minWidth: 0,
                  overflowWrap: 'anywhere',
                  color: t.text[1],
                }}
              >
                {item.title}
              </Typography>
            </Box>
          ))}
        </Box>
      )}
      {ok && !items.length && (
        <Typography variant="body2" style={{ color: t.text[2] }}>
          No open tasks in the observed sample.
        </Typography>
      )}
      {status.message && (
        <Typography
          variant="caption"
          role={status.state === 'unavailable' ? 'alert' : undefined}
          style={{
            color:
              status.state === 'unavailable' || stale
                ? t.status.warn.fg
                : t.text[2],
          }}
        >
          {status.message}
        </Typography>
      )}
      {ok && data?.partial && (
        <Typography variant="caption" style={{ color: t.text[3] }}>
          Counts come from a sample of the project list.
        </Typography>
      )}
      <AppLink href="/today" variant="caption">
        See all in Today
      </AppLink>
    </Surface>
  );
}
