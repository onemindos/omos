import type { CalendarSummary } from './CalendarPanel';
import type { CloudSummary } from './cloudWork';

export type SourceState =
  | 'loading'
  | 'ok'
  | 'stale'
  | 'unconfigured'
  | 'unavailable';

export type SourceStatus = { state: SourceState; message?: string };

const expired = (expiresAt: string | null | undefined, now: number) =>
  !!expiresAt && Date.parse(expiresAt) < now;

/** Calendar wording matches the Today card; `ok`/`stale` mean there is data to show. */
export function calendarStatus(
  data: CalendarSummary | undefined,
  error: string,
  now = Date.now(),
): SourceStatus {
  if (!data) {
    return error
      ? { state: 'unavailable', message: `Calendar: ${error}.` }
      : { state: 'loading', message: 'Loading calendar…' };
  }
  if (data.state === 'configuration-required')
    return {
      state: 'unconfigured',
      message:
        'Connect read-only Calendar access and select a calendar UID to enable upcoming events.',
    };
  if (data.state === 'unavailable')
    return {
      state: 'unavailable',
      message:
        'Calendar source unavailable; an empty list does not mean no events.',
    };
  if (error || data.state === 'stale' || expired(data.expiresAt, now))
    return { state: 'stale', message: 'Calendar data may be stale.' };
  return { state: 'ok' };
}

/** Zoho Projects wording matches the Today card. */
export function tasksStatus(
  data: CloudSummary | undefined,
  error: boolean,
  now = Date.now(),
): SourceStatus {
  if (!data) {
    return error
      ? {
          state: 'unavailable',
          message:
            'Tasks: summary unavailable; previous observations may be stale.',
        }
      : { state: 'loading', message: 'Loading tasks…' };
  }
  if (data.state === 'configuration-required')
    return {
      state: 'unconfigured',
      message: 'Connect read-only Zoho access to enable task summaries.',
    };
  if (data.state === 'unavailable')
    return { state: 'unavailable', message: 'Tasks: provider unavailable.' };
  if (error || data.state === 'stale' || expired(data.expiresAt, now))
    return { state: 'stale', message: 'Tasks may be stale.' };
  return { state: 'ok' };
}

/** True when the source has real observations to draw. */
export const hasData = (status: SourceStatus) =>
  status.state === 'ok' || status.state === 'stale';
