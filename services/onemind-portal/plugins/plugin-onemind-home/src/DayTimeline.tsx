import { Box, Typography } from '@material-ui/core';
import { useHomeTokens } from './homeTokens';
import {
  nowMarker,
  timelineBlocks,
  timelineLabel,
  TimelineRange,
} from './nextUp';

export type TimelineItem = {
  id: string;
  title: string;
  start: Date;
  end?: Date;
  ended: boolean;
};

function tickShift(h: number, range: TimelineRange) {
  if (h === range.startHour) return 'none';
  if (h === range.endHour) return 'translateX(-100%)';
  return 'translateX(-50%)';
}

const hourLabel = (h: number) => `${String(h % 24).padStart(2, '0')}:00`;

/**
 * Today's remaining and past timed events as blocks on a horizontal track,
 * with a "now" line. Events without an end time are thin pins at their start.
 */
export function DayTimeline({
  items,
  range,
  now,
  highlightId,
}: {
  items: TimelineItem[];
  range: TimelineRange;
  now: Date;
  highlightId?: string;
}) {
  const t = useHomeTokens();
  const blocks = timelineBlocks(items, range);
  const marker = nowMarker(now, range);
  const mid = (range.startHour + range.endHour) / 2;
  const ticks = [
    range.startHour,
    ...(Number.isInteger(mid) ? [mid] : []),
    range.endHour,
  ];
  return (
    <Box>
      <div
        role="img"
        aria-label={timelineLabel(items, range)}
        style={{
          position: 'relative',
          height: 18,
          borderRadius: 9,
          background: t.surface[3],
          overflow: 'hidden',
        }}
      >
        {blocks.map((block, i) => {
          const item = items[i];
          return (
            <span
              key={block.id}
              title={`${block.startLabel} ${block.title}`}
              style={{
                position: 'absolute',
                top: 3,
                bottom: 3,
                left: `${block.leftPct}%`,
                width: `${block.widthPct}%`,
                borderRadius: 6,
                background:
                  block.id === highlightId ? t.brand.base : t.status.info.fg,
                opacity: item.ended ? 0.35 : 1,
                outline: `1px solid ${t.surface[2]}`,
              }}
            />
          );
        })}
        {marker !== undefined && (
          <span
            style={{
              position: 'absolute',
              top: 0,
              bottom: 0,
              left: `${marker}%`,
              width: 2,
              marginLeft: -1,
              background: t.text[1],
            }}
          />
        )}
      </div>
      <Box
        aria-hidden
        style={{ position: 'relative', marginTop: 4, minHeight: 16 }}
      >
        {ticks.map(h => (
          <Typography
            key={h}
            variant="caption"
            style={{
              position: 'absolute',
              left: `${
                ((h - range.startHour) / (range.endHour - range.startHour)) *
                100
              }%`,
              transform: tickShift(h, range),
              color: t.text[3],
              fontVariantNumeric: 'tabular-nums',
            }}
          >
            {hourLabel(h)}
          </Typography>
        ))}
      </Box>
    </Box>
  );
}
